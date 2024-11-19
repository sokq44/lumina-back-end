<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/handshake_client_tls13.go - Go Documentation Server</title>

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
<a href="handshake_client_tls13.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">handshake_client_tls13.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/tls">crypto/tls</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2018 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package tls
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/ecdh&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/hmac&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/rsa&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type clientHandshakeStateTLS13 struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	c           *Conn
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	ctx         context.Context
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	serverHello *serverHelloMsg
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	hello       *clientHelloMsg
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	ecdheKey    *ecdh.PrivateKey
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	session     *SessionState
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	earlySecret []byte
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	binderKey   []byte
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	certReq       *certificateRequestMsgTLS13
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	usingPSK      bool
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	sentDummyCCS  bool
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	suite         *cipherSuiteTLS13
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	transcript    hash.Hash
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	masterSecret  []byte
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	trafficSecret []byte <span class="comment">// client_application_traffic_secret_0</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// handshake requires hs.c, hs.hello, hs.serverHello, hs.ecdheKey, and,</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// optionally, hs.session, hs.earlySecret and hs.binderKey to be set.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) handshake() error {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	c := hs.c
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if needFIPS() {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		return errors.New(&#34;tls: internal error: TLS 1.3 reached in FIPS mode&#34;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// The server must not select TLS 1.3 in a renegotiation. See RFC 8446,</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// sections 4.1.2 and 4.1.3.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if c.handshakes &gt; 0 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		c.sendAlert(alertProtocolVersion)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected TLS 1.3 in a renegotiation&#34;)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// Consistency check on the presence of a keyShare and its parameters.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if hs.ecdheKey == nil || len(hs.hello.keyShares) != 1 {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if err := hs.checkServerHelloOrHRR(); err != nil {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return err
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	hs.transcript = hs.suite.hash.New()
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	if err := transcriptMsg(hs.hello, hs.transcript); err != nil {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return err
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	if bytes.Equal(hs.serverHello.random, helloRetryRequestRandom) {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if err := hs.sendDummyChangeCipherSpec(); err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return err
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if err := hs.processHelloRetryRequest(); err != nil {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return err
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if err := transcriptMsg(hs.serverHello, hs.transcript); err != nil {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return err
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	c.buffering = true
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if err := hs.processServerHello(); err != nil {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return err
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	if err := hs.sendDummyChangeCipherSpec(); err != nil {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		return err
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	if err := hs.establishHandshakeKeys(); err != nil {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return err
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if err := hs.readServerParameters(); err != nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return err
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if err := hs.readServerCertificate(); err != nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		return err
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if err := hs.readServerFinished(); err != nil {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		return err
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if err := hs.sendClientCertificate(); err != nil {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		return err
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if err := hs.sendClientFinished(); err != nil {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		return err
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if _, err := c.flush(); err != nil {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return err
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	c.isHandshakeComplete.Store(true)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	return nil
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// checkServerHelloOrHRR does validity checks that apply to both ServerHello and</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// HelloRetryRequest messages. It sets hs.suite.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) checkServerHelloOrHRR() error {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	c := hs.c
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if hs.serverHello.supportedVersion == 0 {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		c.sendAlert(alertMissingExtension)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected TLS 1.3 using the legacy version field&#34;)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if hs.serverHello.supportedVersion != VersionTLS13 {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected an invalid version after a HelloRetryRequest&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if hs.serverHello.vers != VersionTLS12 {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent an incorrect legacy version&#34;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	if hs.serverHello.ocspStapling ||
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		hs.serverHello.ticketSupported ||
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		hs.serverHello.extendedMasterSecret ||
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		hs.serverHello.secureRenegotiationSupported ||
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		len(hs.serverHello.secureRenegotiation) != 0 ||
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		len(hs.serverHello.alpnProtocol) != 0 ||
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		len(hs.serverHello.scts) != 0 {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		c.sendAlert(alertUnsupportedExtension)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent a ServerHello extension forbidden in TLS 1.3&#34;)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if !bytes.Equal(hs.hello.sessionId, hs.serverHello.sessionId) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server did not echo the legacy session ID&#34;)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if hs.serverHello.compressionMethod != compressionNone {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected unsupported compression format&#34;)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	selectedSuite := mutualCipherSuiteTLS13(hs.hello.cipherSuites, hs.serverHello.cipherSuite)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if hs.suite != nil &amp;&amp; selectedSuite != hs.suite {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server changed cipher suite after a HelloRetryRequest&#34;)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if selectedSuite == nil {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server chose an unconfigured cipher suite&#34;)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	hs.suite = selectedSuite
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	c.cipherSuite = hs.suite.id
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	return nil
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// sendDummyChangeCipherSpec sends a ChangeCipherSpec record for compatibility</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// with middleboxes that didn&#39;t implement TLS correctly. See RFC 8446, Appendix D.4.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) sendDummyChangeCipherSpec() error {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if hs.c.quic != nil {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		return nil
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if hs.sentDummyCCS {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return nil
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	hs.sentDummyCCS = true
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return hs.c.writeChangeCipherRecord()
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// processHelloRetryRequest handles the HRR in hs.serverHello, modifies and</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// resends hs.hello, and reads the new ServerHello into hs.serverHello.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) processHelloRetryRequest() error {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	c := hs.c
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// The first ClientHello gets double-hashed into the transcript upon a</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// HelloRetryRequest. (The idea is that the server might offload transcript</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// storage to the client in the cookie.) See RFC 8446, Section 4.4.1.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	chHash := hs.transcript.Sum(nil)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	hs.transcript.Reset()
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	hs.transcript.Write([]byte{typeMessageHash, 0, 0, uint8(len(chHash))})
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	hs.transcript.Write(chHash)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	if err := transcriptMsg(hs.serverHello, hs.transcript); err != nil {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		return err
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// The only HelloRetryRequest extensions we support are key_share and</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// cookie, and clients must abort the handshake if the HRR would not result</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// in any change in the ClientHello.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	if hs.serverHello.selectedGroup == 0 &amp;&amp; hs.serverHello.cookie == nil {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent an unnecessary HelloRetryRequest message&#34;)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if hs.serverHello.cookie != nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		hs.hello.cookie = hs.serverHello.cookie
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	if hs.serverHello.serverShare.group != 0 {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		c.sendAlert(alertDecodeError)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		return errors.New(&#34;tls: received malformed key_share extension&#34;)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// If the server sent a key_share extension selecting a group, ensure it&#39;s</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// a group we advertised but did not send a key share for, and send a key</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// share for it this time.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if curveID := hs.serverHello.selectedGroup; curveID != 0 {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		curveOK := false
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		for _, id := range hs.hello.supportedCurves {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			if id == curveID {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				curveOK = true
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				break
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if !curveOK {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			c.sendAlert(alertIllegalParameter)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server selected unsupported group&#34;)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		if sentID, _ := curveIDForCurve(hs.ecdheKey.Curve()); sentID == curveID {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			c.sendAlert(alertIllegalParameter)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server sent an unnecessary HelloRetryRequest key_share&#34;)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		if _, ok := curveForCurveID(curveID); !ok {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			return errors.New(&#34;tls: CurvePreferences includes unsupported curve&#34;)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		key, err := generateECDHEKey(c.config.rand(), curveID)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		if err != nil {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			return err
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		hs.ecdheKey = key
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		hs.hello.keyShares = []keyShare{{group: curveID, data: key.PublicKey().Bytes()}}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	hs.hello.raw = nil
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if len(hs.hello.pskIdentities) &gt; 0 {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		pskSuite := cipherSuiteTLS13ByID(hs.session.cipherSuite)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		if pskSuite == nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			return c.sendAlert(alertInternalError)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		if pskSuite.hash == hs.suite.hash {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			<span class="comment">// Update binders and obfuscated_ticket_age.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			ticketAge := c.config.time().Sub(time.Unix(int64(hs.session.createdAt), 0))
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			hs.hello.pskIdentities[0].obfuscatedTicketAge = uint32(ticketAge/time.Millisecond) + hs.session.ageAdd
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			transcript := hs.suite.hash.New()
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			transcript.Write([]byte{typeMessageHash, 0, 0, uint8(len(chHash))})
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			transcript.Write(chHash)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			if err := transcriptMsg(hs.serverHello, transcript); err != nil {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>				return err
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			helloBytes, err := hs.hello.marshalWithoutBinders()
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			if err != nil {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>				return err
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			transcript.Write(helloBytes)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			pskBinders := [][]byte{hs.suite.finishedHash(hs.binderKey, transcript)}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			if err := hs.hello.updateBinders(pskBinders); err != nil {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				return err
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		} else {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			<span class="comment">// Server selected a cipher suite incompatible with the PSK.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			hs.hello.pskIdentities = nil
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			hs.hello.pskBinders = nil
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	if hs.hello.earlyData {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		hs.hello.earlyData = false
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		c.quicRejectedEarlyData()
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(hs.hello, hs.transcript); err != nil {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		return err
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// serverHelloMsg is not included in the transcript</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	msg, err := c.readHandshake(nil)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	if err != nil {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		return err
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	serverHello, ok := msg.(*serverHelloMsg)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	if !ok {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		return unexpectedMessageError(serverHello, msg)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	hs.serverHello = serverHello
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if err := hs.checkServerHelloOrHRR(); err != nil {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		return err
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	return nil
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) processServerHello() error {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	c := hs.c
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	if bytes.Equal(hs.serverHello.random, helloRetryRequestRandom) {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent two HelloRetryRequest messages&#34;)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if len(hs.serverHello.cookie) != 0 {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		c.sendAlert(alertUnsupportedExtension)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent a cookie in a normal ServerHello&#34;)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if hs.serverHello.selectedGroup != 0 {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		c.sendAlert(alertDecodeError)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return errors.New(&#34;tls: malformed key_share extension&#34;)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	if hs.serverHello.serverShare.group == 0 {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server did not send a key share&#34;)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if sentID, _ := curveIDForCurve(hs.ecdheKey.Curve()); hs.serverHello.serverShare.group != sentID {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected unsupported group&#34;)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if !hs.serverHello.selectedIdentityPresent {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return nil
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if int(hs.serverHello.selectedIdentity) &gt;= len(hs.hello.pskIdentities) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected an invalid PSK&#34;)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	if len(hs.hello.pskIdentities) != 1 || hs.session == nil {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	pskSuite := cipherSuiteTLS13ByID(hs.session.cipherSuite)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if pskSuite == nil {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	if pskSuite.hash != hs.suite.hash {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected an invalid PSK and cipher suite pair&#34;)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	hs.usingPSK = true
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	c.didResume = true
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	c.peerCertificates = hs.session.peerCertificates
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	c.activeCertHandles = hs.session.activeCertHandles
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	c.verifiedChains = hs.session.verifiedChains
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	c.ocspResponse = hs.session.ocspResponse
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	c.scts = hs.session.scts
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	return nil
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) establishHandshakeKeys() error {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	c := hs.c
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	peerKey, err := hs.ecdheKey.Curve().NewPublicKey(hs.serverHello.serverShare.data)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	if err != nil {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid server key share&#34;)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	sharedKey, err := hs.ecdheKey.ECDH(peerKey)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	if err != nil {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid server key share&#34;)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	earlySecret := hs.earlySecret
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	if !hs.usingPSK {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		earlySecret = hs.suite.extract(nil, nil)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	handshakeSecret := hs.suite.extract(sharedKey,
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		hs.suite.deriveSecret(earlySecret, &#34;derived&#34;, nil))
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	clientSecret := hs.suite.deriveSecret(handshakeSecret,
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		clientHandshakeTrafficLabel, hs.transcript)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	c.out.setTrafficSecret(hs.suite, QUICEncryptionLevelHandshake, clientSecret)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	serverSecret := hs.suite.deriveSecret(handshakeSecret,
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		serverHandshakeTrafficLabel, hs.transcript)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	c.in.setTrafficSecret(hs.suite, QUICEncryptionLevelHandshake, serverSecret)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		if c.hand.Len() != 0 {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			c.sendAlert(alertUnexpectedMessage)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		c.quicSetWriteSecret(QUICEncryptionLevelHandshake, hs.suite.id, clientSecret)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		c.quicSetReadSecret(QUICEncryptionLevelHandshake, hs.suite.id, serverSecret)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelClientHandshake, hs.hello.random, clientSecret)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	if err != nil {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		return err
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelServerHandshake, hs.hello.random, serverSecret)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if err != nil {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		return err
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	hs.masterSecret = hs.suite.extract(nil,
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		hs.suite.deriveSecret(handshakeSecret, &#34;derived&#34;, nil))
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	return nil
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) readServerParameters() error {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	c := hs.c
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	msg, err := c.readHandshake(hs.transcript)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if err != nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		return err
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	encryptedExtensions, ok := msg.(*encryptedExtensionsMsg)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	if !ok {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		return unexpectedMessageError(encryptedExtensions, msg)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if err := checkALPN(hs.hello.alpnProtocols, encryptedExtensions.alpnProtocol, c.quic != nil); err != nil {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446 specifies that no_application_protocol is sent by servers, but</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		<span class="comment">// does not specify how clients handle the selection of an incompatible protocol.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		<span class="comment">// RFC 9001 Section 8.1 specifies that QUIC clients send no_application_protocol</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		<span class="comment">// in this case. Always sending no_application_protocol seems reasonable.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		c.sendAlert(alertNoApplicationProtocol)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		return err
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	c.clientProtocol = encryptedExtensions.alpnProtocol
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		if encryptedExtensions.quicTransportParameters == nil {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			<span class="comment">// RFC 9001 Section 8.2.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			c.sendAlert(alertMissingExtension)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server did not send a quic_transport_parameters extension&#34;)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		c.quicSetTransportParameters(encryptedExtensions.quicTransportParameters)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	} else {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		if encryptedExtensions.quicTransportParameters != nil {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			c.sendAlert(alertUnsupportedExtension)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server sent an unexpected quic_transport_parameters extension&#34;)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if !hs.hello.earlyData &amp;&amp; encryptedExtensions.earlyData {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		c.sendAlert(alertUnsupportedExtension)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server sent an unexpected early_data extension&#34;)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	if hs.hello.earlyData &amp;&amp; !encryptedExtensions.earlyData {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		c.quicRejectedEarlyData()
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if encryptedExtensions.earlyData {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		if hs.session.cipherSuite != c.cipherSuite {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			c.sendAlert(alertHandshakeFailure)
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server accepted 0-RTT with the wrong cipher suite&#34;)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		if hs.session.alpnProtocol != c.clientProtocol {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			c.sendAlert(alertHandshakeFailure)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			return errors.New(&#34;tls: server accepted 0-RTT with the wrong ALPN&#34;)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	return nil
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) readServerCertificate() error {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	c := hs.c
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">// Either a PSK or a certificate is always used, but not both.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	<span class="comment">// See RFC 8446, Section 4.1.1.</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if hs.usingPSK {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// Make sure the connection is still being verified whether or not this</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">// is a resumption. Resumptions currently don&#39;t reverify certificates so</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// they don&#39;t call verifyServerCertificate. See Issue 31641.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		if c.config.VerifyConnection != nil {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			if err := c.config.VerifyConnection(c.connectionStateLocked()); err != nil {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				c.sendAlert(alertBadCertificate)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				return err
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		return nil
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	msg, err := c.readHandshake(hs.transcript)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if err != nil {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		return err
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	certReq, ok := msg.(*certificateRequestMsgTLS13)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	if ok {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		hs.certReq = certReq
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		msg, err = c.readHandshake(hs.transcript)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		if err != nil {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			return err
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	certMsg, ok := msg.(*certificateMsgTLS13)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	if !ok {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		return unexpectedMessageError(certMsg, msg)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if len(certMsg.certificate.Certificate) == 0 {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		c.sendAlert(alertDecodeError)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		return errors.New(&#34;tls: received empty certificates message&#34;)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	c.scts = certMsg.certificate.SignedCertificateTimestamps
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	c.ocspResponse = certMsg.certificate.OCSPStaple
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if err := c.verifyServerCertificate(certMsg.certificate.Certificate); err != nil {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		return err
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">// certificateVerifyMsg is included in the transcript, but not until</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// after we verify the handshake signature, since the state before</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// this message was sent is used.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	msg, err = c.readHandshake(nil)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	if err != nil {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		return err
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	}
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	certVerify, ok := msg.(*certificateVerifyMsg)
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	if !ok {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		return unexpectedMessageError(certVerify, msg)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">// See RFC 8446, Section 4.4.3.</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	if !isSupportedSignatureAlgorithm(certVerify.signatureAlgorithm, supportedSignatureAlgorithms()) {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		return errors.New(&#34;tls: certificate used with invalid signature algorithm&#34;)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	sigType, sigHash, err := typeAndHashFromSignatureScheme(certVerify.signatureAlgorithm)
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	if err != nil {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	if sigType == signaturePKCS1v15 || sigHash == crypto.SHA1 {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		return errors.New(&#34;tls: certificate used with invalid signature algorithm&#34;)
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	signed := signedMessage(sigHash, serverSignatureContext, hs.transcript)
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if err := verifyHandshakeSignature(sigType, c.peerCertificates[0].PublicKey,
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		sigHash, signed, certVerify.signature); err != nil {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		c.sendAlert(alertDecryptError)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid signature by the server certificate: &#34; + err.Error())
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if err := transcriptMsg(certVerify, hs.transcript); err != nil {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		return err
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	return nil
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) readServerFinished() error {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	c := hs.c
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// finishedMsg is included in the transcript, but not until after we</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// check the client version, since the state before this message was</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// sent is used during verification.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	msg, err := c.readHandshake(nil)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if err != nil {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return err
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	finished, ok := msg.(*finishedMsg)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	if !ok {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		return unexpectedMessageError(finished, msg)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	expectedMAC := hs.suite.finishedHash(c.in.trafficSecret, hs.transcript)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	if !hmac.Equal(expectedMAC, finished.verifyData) {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		c.sendAlert(alertDecryptError)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid server finished hash&#34;)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	if err := transcriptMsg(finished, hs.transcript); err != nil {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		return err
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// Derive secrets that take context through the server Finished.</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	hs.trafficSecret = hs.suite.deriveSecret(hs.masterSecret,
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		clientApplicationTrafficLabel, hs.transcript)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	serverSecret := hs.suite.deriveSecret(hs.masterSecret,
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		serverApplicationTrafficLabel, hs.transcript)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	c.in.setTrafficSecret(hs.suite, QUICEncryptionLevelApplication, serverSecret)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelClientTraffic, hs.hello.random, hs.trafficSecret)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	if err != nil {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		return err
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelServerTraffic, hs.hello.random, serverSecret)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	if err != nil {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		return err
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	c.ekm = hs.suite.exportKeyingMaterial(hs.masterSecret, hs.transcript)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	return nil
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) sendClientCertificate() error {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	c := hs.c
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	if hs.certReq == nil {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		return nil
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	cert, err := c.getClientCertificate(&amp;CertificateRequestInfo{
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		AcceptableCAs:    hs.certReq.certificateAuthorities,
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		SignatureSchemes: hs.certReq.supportedSignatureAlgorithms,
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		Version:          c.vers,
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		ctx:              hs.ctx,
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	})
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	if err != nil {
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		return err
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	certMsg := new(certificateMsgTLS13)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	certMsg.certificate = *cert
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	certMsg.scts = hs.certReq.scts &amp;&amp; len(cert.SignedCertificateTimestamps) &gt; 0
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	certMsg.ocspStapling = hs.certReq.ocspStapling &amp;&amp; len(cert.OCSPStaple) &gt; 0
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(certMsg, hs.transcript); err != nil {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		return err
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">// If we sent an empty certificate message, skip the CertificateVerify.</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	if len(cert.Certificate) == 0 {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		return nil
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	certVerifyMsg := new(certificateVerifyMsg)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	certVerifyMsg.hasSignatureAlgorithm = true
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	certVerifyMsg.signatureAlgorithm, err = selectSignatureScheme(c.vers, cert, hs.certReq.supportedSignatureAlgorithms)
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	if err != nil {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		<span class="comment">// getClientCertificate returned a certificate incompatible with the</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		<span class="comment">// CertificateRequestInfo supported signature algorithms.</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		c.sendAlert(alertHandshakeFailure)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		return err
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	sigType, sigHash, err := typeAndHashFromSignatureScheme(certVerifyMsg.signatureAlgorithm)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	if err != nil {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	signed := signedMessage(sigHash, clientSignatureContext, hs.transcript)
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	signOpts := crypto.SignerOpts(sigHash)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	if sigType == signatureRSAPSS {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		signOpts = &amp;rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: sigHash}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	sig, err := cert.PrivateKey.(crypto.Signer).Sign(c.config.rand(), signed, signOpts)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if err != nil {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		return errors.New(&#34;tls: failed to sign handshake: &#34; + err.Error())
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	certVerifyMsg.signature = sig
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(certVerifyMsg, hs.transcript); err != nil {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		return err
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	}
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	return nil
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>func (hs *clientHandshakeStateTLS13) sendClientFinished() error {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	c := hs.c
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	finished := &amp;finishedMsg{
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		verifyData: hs.suite.finishedHash(c.out.trafficSecret, hs.transcript),
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(finished, hs.transcript); err != nil {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		return err
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	c.out.setTrafficSecret(hs.suite, QUICEncryptionLevelApplication, hs.trafficSecret)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	if !c.config.SessionTicketsDisabled &amp;&amp; c.config.ClientSessionCache != nil {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		c.resumptionSecret = hs.suite.deriveSecret(hs.masterSecret,
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			resumptionLabel, hs.transcript)
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		if c.hand.Len() != 0 {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			c.sendAlert(alertUnexpectedMessage)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		c.quicSetWriteSecret(QUICEncryptionLevelApplication, hs.suite.id, hs.trafficSecret)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	return nil
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>func (c *Conn) handleNewSessionTicket(msg *newSessionTicketMsgTLS13) error {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	if !c.isClient {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		return errors.New(&#34;tls: received new session ticket from a client&#34;)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	if c.config.SessionTicketsDisabled || c.config.ClientSessionCache == nil {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		return nil
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	<span class="comment">// See RFC 8446, Section 4.6.1.</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	if msg.lifetime == 0 {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		return nil
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	lifetime := time.Duration(msg.lifetime) * time.Second
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if lifetime &gt; maxSessionTicketLifetime {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		return errors.New(&#34;tls: received a session ticket with invalid lifetime&#34;)
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	<span class="comment">// RFC 9001, Section 4.6.1</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	if c.quic != nil &amp;&amp; msg.maxEarlyData != 0 &amp;&amp; msg.maxEarlyData != 0xffffffff {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid early data for QUIC connection&#34;)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	cipherSuite := cipherSuiteTLS13ByID(c.cipherSuite)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	if cipherSuite == nil || c.resumptionSecret == nil {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	psk := cipherSuite.expandLabel(c.resumptionSecret, &#34;resumption&#34;,
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		msg.nonce, cipherSuite.hash.Size())
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	session, err := c.sessionState()
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	if err != nil {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		return err
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	session.secret = psk
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	session.useBy = uint64(c.config.time().Add(lifetime).Unix())
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	session.ageAdd = msg.ageAdd
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	session.EarlyData = c.quic != nil &amp;&amp; msg.maxEarlyData == 0xffffffff <span class="comment">// RFC 9001, Section 4.6.1</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	cs := &amp;ClientSessionState{ticket: msg.label, session: session}
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	if cacheKey := c.clientSessionCacheKey(); cacheKey != &#34;&#34; {
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		c.config.ClientSessionCache.Put(cacheKey, cs)
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	}
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	return nil
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>
</pre><p><a href="handshake_client_tls13.go?m=text">View as plain text</a></p>

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
