<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/dnsclient_unix.go - Go Documentation Server</title>

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
<a href="dnsclient_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">dnsclient_unix.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// DNS client: see RFC 1035.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// Has to be linked into package net for Dial.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// TODO(rsc):</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//	Could potentially handle many outstanding lookups faster.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//	Random UDP source port (net.Dial should do that for us).</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//	Random request IDs.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>package net
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>import (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;internal/bytealg&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;internal/itoa&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;sync/atomic&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;golang.org/x/net/dns/dnsmessage&#34;
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>const (
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// to be used as a useTCP parameter to exchange</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	useTCPOnly  = true
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	useUDPOrTCP = false
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// Maximum DNS packet size.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// Value taken from https://dnsflagday.net/2020/.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	maxDNSPacketSize = 1232
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>var (
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	errLameReferral              = errors.New(&#34;lame referral&#34;)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	errCannotUnmarshalDNSMessage = errors.New(&#34;cannot unmarshal DNS message&#34;)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	errCannotMarshalDNSMessage   = errors.New(&#34;cannot marshal DNS message&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	errServerMisbehaving         = errors.New(&#34;server misbehaving&#34;)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	errInvalidDNSResponse        = errors.New(&#34;invalid DNS response&#34;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	errNoAnswerFromDNSServer     = errors.New(&#34;no answer from DNS server&#34;)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// errServerTemporarilyMisbehaving is like errServerMisbehaving, except</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// that when it gets translated to a DNSError, the IsTemporary field</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// gets set to true.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	errServerTemporarilyMisbehaving = errors.New(&#34;server misbehaving&#34;)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func newRequest(q dnsmessage.Question, ad bool) (id uint16, udpReq, tcpReq []byte, err error) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	id = uint16(randInt())
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	b := dnsmessage.NewBuilder(make([]byte, 2, 514), dnsmessage.Header{ID: id, RecursionDesired: true, AuthenticData: ad})
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if err := b.StartQuestions(); err != nil {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if err := b.Question(q); err != nil {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Accept packets up to maxDNSPacketSize.  RFC 6891.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if err := b.StartAdditionals(); err != nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	var rh dnsmessage.ResourceHeader
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if err := rh.SetEDNS0(maxDNSPacketSize, dnsmessage.RCodeSuccess, false); err != nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if err := b.OPTResource(rh, dnsmessage.OPTResource{}); err != nil {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	tcpReq, err = b.Finish()
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if err != nil {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		return 0, nil, nil, err
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	udpReq = tcpReq[2:]
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	l := len(tcpReq) - 2
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	tcpReq[0] = byte(l &gt;&gt; 8)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	tcpReq[1] = byte(l)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	return id, udpReq, tcpReq, nil
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func checkResponse(reqID uint16, reqQues dnsmessage.Question, respHdr dnsmessage.Header, respQues dnsmessage.Question) bool {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if !respHdr.Response {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		return false
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if reqID != respHdr.ID {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		return false
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if reqQues.Type != respQues.Type || reqQues.Class != respQues.Class || !equalASCIIName(reqQues.Name, respQues.Name) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		return false
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	return true
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func dnsPacketRoundTrip(c Conn, id uint16, query dnsmessage.Question, b []byte) (dnsmessage.Parser, dnsmessage.Header, error) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if _, err := c.Write(b); err != nil {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	b = make([]byte, maxDNSPacketSize)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	for {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		n, err := c.Read(b)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		if err != nil {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		var p dnsmessage.Parser
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// Ignore invalid responses as they may be malicious</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// forgery attempts. Instead continue waiting until</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		<span class="comment">// timeout. See golang.org/issue/13281.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		h, err := p.Start(b[:n])
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		if err != nil {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			continue
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		q, err := p.Question()
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		if err != nil || !checkResponse(id, query, h, q) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			continue
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		return p, h, nil
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func dnsStreamRoundTrip(c Conn, id uint16, query dnsmessage.Question, b []byte) (dnsmessage.Parser, dnsmessage.Header, error) {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if _, err := c.Write(b); err != nil {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	b = make([]byte, 1280) <span class="comment">// 1280 is a reasonable initial size for IP over Ethernet, see RFC 4035</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if _, err := io.ReadFull(c, b[:2]); err != nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	l := int(b[0])&lt;&lt;8 | int(b[1])
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	if l &gt; len(b) {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		b = make([]byte, l)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	n, err := io.ReadFull(c, b[:l])
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if err != nil {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	var p dnsmessage.Parser
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	h, err := p.Start(b[:n])
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if err != nil {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, errCannotUnmarshalDNSMessage
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	q, err := p.Question()
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if err != nil {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, errCannotUnmarshalDNSMessage
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if !checkResponse(id, query, h, q) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, errInvalidDNSResponse
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	return p, h, nil
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// exchange sends a query on the connection and hopes for a response.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (r *Resolver) exchange(ctx context.Context, server string, q dnsmessage.Question, timeout time.Duration, useTCP, ad bool) (dnsmessage.Parser, dnsmessage.Header, error) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	q.Class = dnsmessage.ClassINET
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	id, udpReq, tcpReq, err := newRequest(q, ad)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if err != nil {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, dnsmessage.Header{}, errCannotMarshalDNSMessage
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	var networks []string
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if useTCP {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		networks = []string{&#34;tcp&#34;}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	} else {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		networks = []string{&#34;udp&#34;, &#34;tcp&#34;}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	for _, network := range networks {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		defer cancel()
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		c, err := r.dial(ctx, network, server)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if err != nil {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			return dnsmessage.Parser{}, dnsmessage.Header{}, err
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		if d, ok := ctx.Deadline(); ok &amp;&amp; !d.IsZero() {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			c.SetDeadline(d)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		var p dnsmessage.Parser
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		var h dnsmessage.Header
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		if _, ok := c.(PacketConn); ok {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			p, h, err = dnsPacketRoundTrip(c, id, q, udpReq)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		} else {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			p, h, err = dnsStreamRoundTrip(c, id, q, tcpReq)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		c.Close()
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		if err != nil {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			return dnsmessage.Parser{}, dnsmessage.Header{}, mapErr(err)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		if err := p.SkipQuestion(); err != dnsmessage.ErrSectionDone {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			return dnsmessage.Parser{}, dnsmessage.Header{}, errInvalidDNSResponse
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		if h.Truncated { <span class="comment">// see RFC 5966</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			continue
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		return p, h, nil
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	return dnsmessage.Parser{}, dnsmessage.Header{}, errNoAnswerFromDNSServer
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// checkHeader performs basic sanity checks on the header.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func checkHeader(p *dnsmessage.Parser, h dnsmessage.Header) error {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	rcode := extractExtendedRCode(*p, h)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if rcode == dnsmessage.RCodeNameError {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		return errNoSuchHost
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	_, err := p.AnswerHeader()
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if err != nil &amp;&amp; err != dnsmessage.ErrSectionDone {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		return errCannotUnmarshalDNSMessage
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	<span class="comment">// libresolv continues to the next server when it receives</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	<span class="comment">// an invalid referral response. See golang.org/issue/15434.</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	if rcode == dnsmessage.RCodeSuccess &amp;&amp; !h.Authoritative &amp;&amp; !h.RecursionAvailable &amp;&amp; err == dnsmessage.ErrSectionDone {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return errLameReferral
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if rcode != dnsmessage.RCodeSuccess &amp;&amp; rcode != dnsmessage.RCodeNameError {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// None of the error codes make sense</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		<span class="comment">// for the query we sent. If we didn&#39;t get</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// a name error and we didn&#39;t get success,</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// the server is behaving incorrectly or</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// having temporary trouble.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		if rcode == dnsmessage.RCodeServerFailure {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			return errServerTemporarilyMisbehaving
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return errServerMisbehaving
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	return nil
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>func skipToAnswer(p *dnsmessage.Parser, qtype dnsmessage.Type) error {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	for {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		h, err := p.AnswerHeader()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		if err == dnsmessage.ErrSectionDone {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			return errNoSuchHost
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if err != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			return errCannotUnmarshalDNSMessage
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if h.Type == qtype {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			return nil
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if err := p.SkipAnswer(); err != nil {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			return errCannotUnmarshalDNSMessage
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// extractExtendedRCode extracts the extended RCode from the OPT resource (EDNS(0))</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// If an OPT record is not found, the RCode from the hdr is returned.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func extractExtendedRCode(p dnsmessage.Parser, hdr dnsmessage.Header) dnsmessage.RCode {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	p.SkipAllAnswers()
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	p.SkipAllAuthorities()
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	for {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		ahdr, err := p.AdditionalHeader()
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		if err != nil {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			return hdr.RCode
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		if ahdr.Type == dnsmessage.TypeOPT {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			return ahdr.ExtendedRCode(hdr.RCode)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		if err := p.SkipAdditional(); err != nil {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			return hdr.RCode
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// Do a lookup for a single name, which must be rooted</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// (otherwise answer will not find the answers).</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>func (r *Resolver) tryOneName(ctx context.Context, cfg *dnsConfig, name string, qtype dnsmessage.Type) (dnsmessage.Parser, string, error) {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	var lastErr error
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	serverOffset := cfg.serverOffset()
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	sLen := uint32(len(cfg.servers))
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	n, err := dnsmessage.NewName(name)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if err != nil {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, &#34;&#34;, errCannotMarshalDNSMessage
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	q := dnsmessage.Question{
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		Name:  n,
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		Type:  qtype,
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		Class: dnsmessage.ClassINET,
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	for i := 0; i &lt; cfg.attempts; i++ {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		for j := uint32(0); j &lt; sLen; j++ {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			server := cfg.servers[(serverOffset+j)%sLen]
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			p, h, err := r.exchange(ctx, server, q, cfg.timeout, cfg.useTCP, cfg.trustAD)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if err != nil {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				dnsErr := &amp;DNSError{
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>					Err:    err.Error(),
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>					Name:   name,
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>					Server: server,
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				if nerr, ok := err.(Error); ok &amp;&amp; nerr.Timeout() {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>					dnsErr.IsTimeout = true
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				<span class="comment">// Set IsTemporary for socket-level errors. Note that this flag</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				<span class="comment">// may also be used to indicate a SERVFAIL response.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				if _, ok := err.(*OpError); ok {
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					dnsErr.IsTemporary = true
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				lastErr = dnsErr
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>				continue
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			if err := checkHeader(&amp;p, h); err != nil {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>				dnsErr := &amp;DNSError{
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>					Err:    err.Error(),
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>					Name:   name,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>					Server: server,
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>				}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				if err == errServerTemporarilyMisbehaving {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>					dnsErr.IsTemporary = true
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				if err == errNoSuchHost {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>					<span class="comment">// The name does not exist, so trying</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>					<span class="comment">// another server won&#39;t help.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>					dnsErr.IsNotFound = true
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>					return p, server, dnsErr
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>				}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				lastErr = dnsErr
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				continue
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			err = skipToAnswer(&amp;p, qtype)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			if err == nil {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				return p, server, nil
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			lastErr = &amp;DNSError{
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>				Err:    err.Error(),
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				Name:   name,
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				Server: server,
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			if err == errNoSuchHost {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				<span class="comment">// The name does not exist, so trying another</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>				<span class="comment">// server won&#39;t help.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>				lastErr.(*DNSError).IsNotFound = true
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				return p, server, lastErr
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	return dnsmessage.Parser{}, &#34;&#34;, lastErr
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// A resolverConfig represents a DNS stub resolver configuration.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>type resolverConfig struct {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	initOnce sync.Once <span class="comment">// guards init of resolverConfig</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// ch is used as a semaphore that only allows one lookup at a</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// time to recheck resolv.conf.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	ch          chan struct{} <span class="comment">// guards lastChecked and modTime</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	lastChecked time.Time     <span class="comment">// last time resolv.conf was checked</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	dnsConfig atomic.Pointer[dnsConfig] <span class="comment">// parsed resolv.conf structure used in lookups</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>var resolvConf resolverConfig
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>func getSystemDNSConfig() *dnsConfig {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	resolvConf.tryUpdate(&#34;/etc/resolv.conf&#34;)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	return resolvConf.dnsConfig.Load()
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span><span class="comment">// init initializes conf and is only called via conf.initOnce.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>func (conf *resolverConfig) init() {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// Set dnsConfig and lastChecked so we don&#39;t parse</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// resolv.conf twice the first time.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	conf.dnsConfig.Store(dnsReadConfig(&#34;/etc/resolv.conf&#34;))
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	conf.lastChecked = time.Now()
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// Prepare ch so that only one update of resolverConfig may</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// run at once.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	conf.ch = make(chan struct{}, 1)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// tryUpdate tries to update conf with the named resolv.conf file.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">// The name variable only exists for testing. It is otherwise always</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// &#34;/etc/resolv.conf&#34;.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>func (conf *resolverConfig) tryUpdate(name string) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	conf.initOnce.Do(conf.init)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	if conf.dnsConfig.Load().noReload {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		return
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// Ensure only one update at a time checks resolv.conf.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	if !conf.tryAcquireSema() {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		return
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	defer conf.releaseSema()
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	now := time.Now()
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if conf.lastChecked.After(now.Add(-5 * time.Second)) {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		return
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	conf.lastChecked = now
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	switch runtime.GOOS {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	case &#34;windows&#34;:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s no file on disk, so don&#39;t bother checking</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		<span class="comment">// and failing.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// The Windows implementation of dnsReadConfig (called</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		<span class="comment">// below) ignores the name.</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	default:
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		var mtime time.Time
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		if fi, err := os.Stat(name); err == nil {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			mtime = fi.ModTime()
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if mtime.Equal(conf.dnsConfig.Load().mtime) {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			return
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	dnsConf := dnsReadConfig(name)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	conf.dnsConfig.Store(dnsConf)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>func (conf *resolverConfig) tryAcquireSema() bool {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	select {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	case conf.ch &lt;- struct{}{}:
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return true
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	default:
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		return false
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>func (conf *resolverConfig) releaseSema() {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	&lt;-conf.ch
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>func (r *Resolver) lookup(ctx context.Context, name string, qtype dnsmessage.Type, conf *dnsConfig) (dnsmessage.Parser, string, error) {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	if !isDomainName(name) {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// We used to use &#34;invalid domain name&#34; as the error,</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		<span class="comment">// but that is a detail of the specific lookup mechanism.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		<span class="comment">// Other lookups might allow broader name syntax</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		<span class="comment">// (for example Multicast DNS allows UTF-8; see RFC 6762).</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		<span class="comment">// For consistency with libc resolvers, report no such host.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return dnsmessage.Parser{}, &#34;&#34;, &amp;DNSError{Err: errNoSuchHost.Error(), Name: name, IsNotFound: true}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	if conf == nil {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		conf = getSystemDNSConfig()
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	var (
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		p      dnsmessage.Parser
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		server string
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		err    error
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	for _, fqdn := range conf.nameList(name) {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		p, server, err = r.tryOneName(ctx, conf, fqdn, qtype)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		if err == nil {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			break
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		if nerr, ok := err.(Error); ok &amp;&amp; nerr.Temporary() &amp;&amp; r.strictErrors() {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			<span class="comment">// If we hit a temporary error with StrictErrors enabled,</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			<span class="comment">// stop immediately instead of trying more names.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			break
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if err == nil {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		return p, server, nil
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	if err, ok := err.(*DNSError); ok {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		<span class="comment">// Show original name passed to lookup, not suffixed one.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		<span class="comment">// In general we might have tried many suffixes; showing</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		<span class="comment">// just one is misleading. See also golang.org/issue/6324.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		err.Name = name
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	return dnsmessage.Parser{}, &#34;&#34;, err
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">// avoidDNS reports whether this is a hostname for which we should not</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span><span class="comment">// use DNS. Currently this includes only .onion, per RFC 7686. See</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span><span class="comment">// golang.org/issue/13705. Does not cover .local names (RFC 6762),</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">// see golang.org/issue/16739.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>func avoidDNS(name string) bool {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	if name == &#34;&#34; {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		return true
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	if name[len(name)-1] == &#39;.&#39; {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		name = name[:len(name)-1]
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	return stringsHasSuffixFold(name, &#34;.onion&#34;)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">// nameList returns a list of names for sequential DNS queries.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>func (conf *dnsConfig) nameList(name string) []string {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	<span class="comment">// Check name length (see isDomainName).</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	l := len(name)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	rooted := l &gt; 0 &amp;&amp; name[l-1] == &#39;.&#39;
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if l &gt; 254 || l == 254 &amp;&amp; !rooted {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		return nil
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// If name is rooted (trailing dot), try only that name.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if rooted {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		if avoidDNS(name) {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			return nil
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		return []string{name}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	hasNdots := bytealg.CountString(name, &#39;.&#39;) &gt;= conf.ndots
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	name += &#34;.&#34;
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	l++
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	<span class="comment">// Build list of search choices.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	names := make([]string, 0, 1+len(conf.search))
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	<span class="comment">// If name has enough dots, try unsuffixed first.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	if hasNdots &amp;&amp; !avoidDNS(name) {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		names = append(names, name)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">// Try suffixes that are not too long (see isDomainName).</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	for _, suffix := range conf.search {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>		fqdn := name + suffix
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		if !avoidDNS(fqdn) &amp;&amp; len(fqdn) &lt;= 254 {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			names = append(names, fqdn)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	<span class="comment">// Try unsuffixed, if not tried first above.</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	if !hasNdots &amp;&amp; !avoidDNS(name) {
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		names = append(names, name)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	return names
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// hostLookupOrder specifies the order of LookupHost lookup strategies.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// It is basically a simplified representation of nsswitch.conf.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// &#34;files&#34; means /etc/hosts.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>type hostLookupOrder int
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>const (
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// hostLookupCgo means defer to cgo.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	hostLookupCgo      hostLookupOrder = iota
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	hostLookupFilesDNS                 <span class="comment">// files first</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	hostLookupDNSFiles                 <span class="comment">// dns first</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	hostLookupFiles                    <span class="comment">// only files</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	hostLookupDNS                      <span class="comment">// only DNS</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>var lookupOrderName = map[hostLookupOrder]string{
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	hostLookupCgo:      &#34;cgo&#34;,
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	hostLookupFilesDNS: &#34;files,dns&#34;,
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	hostLookupDNSFiles: &#34;dns,files&#34;,
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	hostLookupFiles:    &#34;files&#34;,
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	hostLookupDNS:      &#34;dns&#34;,
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>func (o hostLookupOrder) String() string {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	if s, ok := lookupOrderName[o]; ok {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		return s
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	return &#34;hostLookupOrder=&#34; + itoa.Itoa(int(o)) + &#34;??&#34;
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>func (r *Resolver) goLookupHostOrder(ctx context.Context, name string, order hostLookupOrder, conf *dnsConfig) (addrs []string, err error) {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	if order == hostLookupFilesDNS || order == hostLookupFiles {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">// Use entries from /etc/hosts if they match.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		addrs, _ = lookupStaticHost(name)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		if len(addrs) &gt; 0 {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			return
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		if order == hostLookupFiles {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			return nil, &amp;DNSError{Err: errNoSuchHost.Error(), Name: name, IsNotFound: true}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	ips, _, err := r.goLookupIPCNAMEOrder(ctx, &#34;ip&#34;, name, order, conf)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if err != nil {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	addrs = make([]string, 0, len(ips))
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	for _, ip := range ips {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		addrs = append(addrs, ip.String())
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	return
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span><span class="comment">// lookup entries from /etc/hosts</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>func goLookupIPFiles(name string) (addrs []IPAddr, canonical string) {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	addr, canonical := lookupStaticHost(name)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	for _, haddr := range addr {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		haddr, zone := splitHostZone(haddr)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		if ip := ParseIP(haddr); ip != nil {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>			addr := IPAddr{IP: ip, Zone: zone}
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			addrs = append(addrs, addr)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	sortByRFC6724(addrs)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	return addrs, canonical
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">// goLookupIP is the native Go implementation of LookupIP.</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// The libc versions are in cgo_*.go.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>func (r *Resolver) goLookupIP(ctx context.Context, network, host string, order hostLookupOrder, conf *dnsConfig) (addrs []IPAddr, err error) {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	addrs, _, err = r.goLookupIPCNAMEOrder(ctx, network, host, order, conf)
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	return
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func (r *Resolver) goLookupIPCNAMEOrder(ctx context.Context, network, name string, order hostLookupOrder, conf *dnsConfig) (addrs []IPAddr, cname dnsmessage.Name, err error) {
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	if order == hostLookupFilesDNS || order == hostLookupFiles {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		var canonical string
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		addrs, canonical = goLookupIPFiles(name)
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		if len(addrs) &gt; 0 {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>			var err error
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			cname, err = dnsmessage.NewName(canonical)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			if err != nil {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				return nil, dnsmessage.Name{}, err
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			return addrs, cname, nil
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		if order == hostLookupFiles {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			return nil, dnsmessage.Name{}, &amp;DNSError{Err: errNoSuchHost.Error(), Name: name, IsNotFound: true}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	if !isDomainName(name) {
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		<span class="comment">// See comment in func lookup above about use of errNoSuchHost.</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		return nil, dnsmessage.Name{}, &amp;DNSError{Err: errNoSuchHost.Error(), Name: name, IsNotFound: true}
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	type result struct {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		p      dnsmessage.Parser
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		server string
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		error
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	if conf == nil {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		conf = getSystemDNSConfig()
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	lane := make(chan result, 1)
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	qtypes := []dnsmessage.Type{dnsmessage.TypeA, dnsmessage.TypeAAAA}
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	if network == &#34;CNAME&#34; {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		qtypes = append(qtypes, dnsmessage.TypeCNAME)
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	switch ipVersion(network) {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	case &#39;4&#39;:
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		qtypes = []dnsmessage.Type{dnsmessage.TypeA}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	case &#39;6&#39;:
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		qtypes = []dnsmessage.Type{dnsmessage.TypeAAAA}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	var queryFn func(fqdn string, qtype dnsmessage.Type)
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	var responseFn func(fqdn string, qtype dnsmessage.Type) result
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	if conf.singleRequest {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		queryFn = func(fqdn string, qtype dnsmessage.Type) {}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		responseFn = func(fqdn string, qtype dnsmessage.Type) result {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			dnsWaitGroup.Add(1)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			defer dnsWaitGroup.Done()
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			p, server, err := r.tryOneName(ctx, conf, fqdn, qtype)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			return result{p, server, err}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		}
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	} else {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		queryFn = func(fqdn string, qtype dnsmessage.Type) {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			dnsWaitGroup.Add(1)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>			go func(qtype dnsmessage.Type) {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>				p, server, err := r.tryOneName(ctx, conf, fqdn, qtype)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>				lane &lt;- result{p, server, err}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>				dnsWaitGroup.Done()
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			}(qtype)
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		responseFn = func(fqdn string, qtype dnsmessage.Type) result {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			return &lt;-lane
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	var lastErr error
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	for _, fqdn := range conf.nameList(name) {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		for _, qtype := range qtypes {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>			queryFn(fqdn, qtype)
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		hitStrictError := false
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		for _, qtype := range qtypes {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			result := responseFn(fqdn, qtype)
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			if result.error != nil {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>				if nerr, ok := result.error.(Error); ok &amp;&amp; nerr.Temporary() &amp;&amp; r.strictErrors() {
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>					<span class="comment">// This error will abort the nameList loop.</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>					hitStrictError = true
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>					lastErr = result.error
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>				} else if lastErr == nil || fqdn == name+&#34;.&#34; {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>					<span class="comment">// Prefer error for original name.</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>					lastErr = result.error
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>				}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>				continue
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			<span class="comment">// Presotto says it&#39;s okay to assume that servers listed in</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			<span class="comment">// /etc/resolv.conf are recursive resolvers.</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			<span class="comment">// We asked for recursion, so it should have included all the</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			<span class="comment">// answers we need in this one packet.</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			<span class="comment">// Further, RFC 1034 section 4.3.1 says that &#34;the recursive</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			<span class="comment">// response to a query will be... The answer to the query,</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			<span class="comment">// possibly preface by one or more CNAME RRs that specify</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			<span class="comment">// aliases encountered on the way to an answer.&#34;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			<span class="comment">// Therefore, we should be able to assume that we can ignore</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			<span class="comment">// CNAMEs and that the A and AAAA records we requested are</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			<span class="comment">// for the canonical name.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		loop:
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			for {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>				h, err := result.p.AnswerHeader()
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>				if err != nil &amp;&amp; err != dnsmessage.ErrSectionDone {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>					lastErr = &amp;DNSError{
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>						Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>						Name:   name,
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>						Server: result.server,
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>					}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				}
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>				if err != nil {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>					break
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>				}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>				switch h.Type {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>				case dnsmessage.TypeA:
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>					a, err := result.p.AResource()
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>					if err != nil {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>						lastErr = &amp;DNSError{
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>							Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>							Name:   name,
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>							Server: result.server,
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>						}
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>						break loop
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>					}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>					addrs = append(addrs, IPAddr{IP: IP(a.A[:])})
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>					if cname.Length == 0 &amp;&amp; h.Name.Length != 0 {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>						cname = h.Name
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>					}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>				case dnsmessage.TypeAAAA:
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>					aaaa, err := result.p.AAAAResource()
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>					if err != nil {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>						lastErr = &amp;DNSError{
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>							Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>							Name:   name,
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>							Server: result.server,
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>						}
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>						break loop
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>					}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>					addrs = append(addrs, IPAddr{IP: IP(aaaa.AAAA[:])})
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>					if cname.Length == 0 &amp;&amp; h.Name.Length != 0 {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>						cname = h.Name
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>					}
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>				case dnsmessage.TypeCNAME:
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>					c, err := result.p.CNAMEResource()
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>					if err != nil {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>						lastErr = &amp;DNSError{
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>							Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>							Name:   name,
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>							Server: result.server,
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>						}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>						break loop
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>					}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>					if cname.Length == 0 &amp;&amp; c.CNAME.Length &gt; 0 {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>						cname = c.CNAME
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>					}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>				default:
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>					if err := result.p.SkipAnswer(); err != nil {
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>						lastErr = &amp;DNSError{
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>							Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>							Name:   name,
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>							Server: result.server,
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>						}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>						break loop
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>					}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>					continue
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>				}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		if hitStrictError {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>			<span class="comment">// If either family hit an error with StrictErrors enabled,</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>			<span class="comment">// discard all addresses. This ensures that network flakiness</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>			<span class="comment">// cannot turn a dualstack hostname IPv4/IPv6-only.</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			addrs = nil
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			break
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		if len(addrs) &gt; 0 || network == &#34;CNAME&#34; &amp;&amp; cname.Length &gt; 0 {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			break
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		}
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	if lastErr, ok := lastErr.(*DNSError); ok {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		<span class="comment">// Show original name passed to lookup, not suffixed one.</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		<span class="comment">// In general we might have tried many suffixes; showing</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		<span class="comment">// just one is misleading. See also golang.org/issue/6324.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		lastErr.Name = name
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	sortByRFC6724(addrs)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	if len(addrs) == 0 &amp;&amp; !(network == &#34;CNAME&#34; &amp;&amp; cname.Length &gt; 0) {
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		if order == hostLookupDNSFiles {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			var canonical string
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			addrs, canonical = goLookupIPFiles(name)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>			if len(addrs) &gt; 0 {
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>				var err error
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>				cname, err = dnsmessage.NewName(canonical)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>				if err != nil {
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>					return nil, dnsmessage.Name{}, err
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>				}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>				return addrs, cname, nil
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>			}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		if lastErr != nil {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>			return nil, dnsmessage.Name{}, lastErr
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	}
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	return addrs, cname, nil
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>}
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span><span class="comment">// goLookupCNAME is the native Go (non-cgo) implementation of LookupCNAME.</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>func (r *Resolver) goLookupCNAME(ctx context.Context, host string, order hostLookupOrder, conf *dnsConfig) (string, error) {
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	_, cname, err := r.goLookupIPCNAMEOrder(ctx, &#34;CNAME&#34;, host, order, conf)
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	return cname.String(), err
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>}
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span><span class="comment">// goLookupPTR is the native Go implementation of LookupAddr.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>func (r *Resolver) goLookupPTR(ctx context.Context, addr string, order hostLookupOrder, conf *dnsConfig) ([]string, error) {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	if order == hostLookupFiles || order == hostLookupFilesDNS {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		names := lookupStaticAddr(addr)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		if len(names) &gt; 0 {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>			return names, nil
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		if order == hostLookupFiles {
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			return nil, &amp;DNSError{Err: errNoSuchHost.Error(), Name: addr, IsNotFound: true}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		}
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	arpa, err := reverseaddr(addr)
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	if err != nil {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		return nil, err
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	}
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	p, server, err := r.lookup(ctx, arpa, dnsmessage.TypePTR, conf)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	if err != nil {
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		var dnsErr *DNSError
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		if errors.As(err, &amp;dnsErr) &amp;&amp; dnsErr.IsNotFound {
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>			if order == hostLookupDNSFiles {
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>				names := lookupStaticAddr(addr)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>				if len(names) &gt; 0 {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>					return names, nil
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>				}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			}
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		}
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		return nil, err
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	var ptrs []string
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	for {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		h, err := p.AnswerHeader()
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		if err == dnsmessage.ErrSectionDone {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>			break
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		if err != nil {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			return nil, &amp;DNSError{
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>				Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>				Name:   addr,
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>				Server: server,
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>			}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		if h.Type != dnsmessage.TypePTR {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			err := p.SkipAnswer()
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>			if err != nil {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>				return nil, &amp;DNSError{
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>					Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>					Name:   addr,
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>					Server: server,
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>				}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			continue
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		ptr, err := p.PTRResource()
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>		if err != nil {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			return nil, &amp;DNSError{
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>				Err:    errCannotUnmarshalDNSMessage.Error(),
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>				Name:   addr,
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>				Server: server,
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			}
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		}
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		ptrs = append(ptrs, ptr.PTR.String())
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	}
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	return ptrs, nil
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>
</pre><p><a href="dnsclient_unix.go?m=text">View as plain text</a></p>

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
