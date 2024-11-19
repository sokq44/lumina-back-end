<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/handshake_server_tls13.go - Go Documentation Server</title>

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
<a href="handshake_server_tls13.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">handshake_server_tls13.go</span>
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
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/hmac&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/rsa&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// maxClientPSKIdentities is the number of client PSK identities the server will</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// attempt to validate. It will ignore the rest not to let cheap ClientHello</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// messages cause too much work in session ticket decryption attempts.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>const maxClientPSKIdentities = 5
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>type serverHandshakeStateTLS13 struct {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	c               *Conn
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	ctx             context.Context
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	clientHello     *clientHelloMsg
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	hello           *serverHelloMsg
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	sentDummyCCS    bool
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	usingPSK        bool
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	earlyData       bool
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	suite           *cipherSuiteTLS13
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	cert            *Certificate
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	sigAlg          SignatureScheme
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	earlySecret     []byte
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	sharedKey       []byte
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	handshakeSecret []byte
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	masterSecret    []byte
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	trafficSecret   []byte <span class="comment">// client_application_traffic_secret_0</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	transcript      hash.Hash
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	clientFinished  []byte
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) handshake() error {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	c := hs.c
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if needFIPS() {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		return errors.New(&#34;tls: internal error: TLS 1.3 reached in FIPS mode&#34;)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// For an overview of the TLS 1.3 handshake, see RFC 8446, Section 2.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if err := hs.processClientHello(); err != nil {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return err
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	if err := hs.checkForResumption(); err != nil {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		return err
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if err := hs.pickCertificate(); err != nil {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return err
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	c.buffering = true
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if err := hs.sendServerParameters(); err != nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		return err
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	if err := hs.sendServerCertificate(); err != nil {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return err
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if err := hs.sendServerFinished(); err != nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return err
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// Note that at this point we could start sending application data without</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// waiting for the client&#39;s second flight, but the application might not</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// expect the lack of replay protection of the ClientHello parameters.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	if _, err := c.flush(); err != nil {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		return err
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if err := hs.readClientCertificate(); err != nil {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		return err
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if err := hs.readClientFinished(); err != nil {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return err
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	c.isHandshakeComplete.Store(true)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	return nil
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) processClientHello() error {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	c := hs.c
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	hs.hello = new(serverHelloMsg)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// TLS 1.3 froze the ServerHello.legacy_version field, and uses</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// supported_versions instead. See RFC 8446, sections 4.1.3 and 4.2.1.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	hs.hello.vers = VersionTLS12
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	hs.hello.supportedVersion = c.vers
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	if len(hs.clientHello.supportedVersions) == 0 {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return errors.New(&#34;tls: client used the legacy version field to negotiate TLS 1.3&#34;)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// Abort if the client is doing a fallback and landing lower than what we</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// support. See RFC 7507, which however does not specify the interaction</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// with supported_versions. The only difference is that with</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// supported_versions a client has a chance to attempt a [TLS 1.2, TLS 1.4]</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// handshake in case TLS 1.3 is broken but 1.2 is not. Alas, in that case,</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// it will have to drop the TLS_FALLBACK_SCSV protection if it falls back to</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// TLS 1.2, because a TLS 1.3 server would abort here. The situation before</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// supported_versions was not better because there was just no way to do a</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// TLS 1.4 handshake without risking the server selecting TLS 1.3.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	for _, id := range hs.clientHello.cipherSuites {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if id == TLS_FALLBACK_SCSV {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			<span class="comment">// Use c.vers instead of max(supported_versions) because an attacker</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			<span class="comment">// could defeat this by adding an arbitrary high version otherwise.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			if c.vers &lt; c.config.maxSupportedVersion(roleServer) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				c.sendAlert(alertInappropriateFallback)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>				return errors.New(&#34;tls: client using inappropriate protocol fallback&#34;)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			break
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if len(hs.clientHello.compressionMethods) != 1 ||
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		hs.clientHello.compressionMethods[0] != compressionNone {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return errors.New(&#34;tls: TLS 1.3 client supports illegal compression methods&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	hs.hello.random = make([]byte, 32)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if _, err := io.ReadFull(c.config.rand(), hs.hello.random); err != nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return err
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if len(hs.clientHello.secureRenegotiation) != 0 {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		c.sendAlert(alertHandshakeFailure)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return errors.New(&#34;tls: initial handshake had non-empty renegotiation extension&#34;)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if hs.clientHello.earlyData &amp;&amp; c.quic != nil {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		if len(hs.clientHello.pskIdentities) == 0 {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			c.sendAlert(alertIllegalParameter)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			return errors.New(&#34;tls: early_data without pre_shared_key&#34;)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	} else if hs.clientHello.earlyData {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		<span class="comment">// See RFC 8446, Section 4.2.10 for the complicated behavior required</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// here. The scenario is that a different server at our address offered</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// to accept early data in the past, which we can&#39;t handle. For now, all</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		<span class="comment">// 0-RTT enabled session tickets need to expire before a Go server can</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		<span class="comment">// replace a server or join a pool. That&#39;s the same requirement that</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		<span class="comment">// applies to mixing or replacing with any TLS 1.2 server.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		c.sendAlert(alertUnsupportedExtension)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		return errors.New(&#34;tls: client sent unexpected early data&#34;)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	hs.hello.sessionId = hs.clientHello.sessionId
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	hs.hello.compressionMethod = compressionNone
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	preferenceList := defaultCipherSuitesTLS13
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if !hasAESGCMHardwareSupport || !aesgcmPreferred(hs.clientHello.cipherSuites) {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		preferenceList = defaultCipherSuitesTLS13NoAES
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	for _, suiteID := range preferenceList {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		hs.suite = mutualCipherSuiteTLS13(hs.clientHello.cipherSuites, suiteID)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		if hs.suite != nil {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			break
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	if hs.suite == nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		c.sendAlert(alertHandshakeFailure)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return errors.New(&#34;tls: no cipher suite supported by both client and server&#34;)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	c.cipherSuite = hs.suite.id
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	hs.hello.cipherSuite = hs.suite.id
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	hs.transcript = hs.suite.hash.New()
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// Pick the ECDHE group in server preference order, but give priority to</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// groups with a key share, to avoid a HelloRetryRequest round-trip.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	var selectedGroup CurveID
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	var clientKeyShare *keyShare
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>GroupSelection:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	for _, preferredGroup := range c.config.curvePreferences() {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		for _, ks := range hs.clientHello.keyShares {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			if ks.group == preferredGroup {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				selectedGroup = ks.group
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				clientKeyShare = &amp;ks
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>				break GroupSelection
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		if selectedGroup != 0 {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			continue
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		for _, group := range hs.clientHello.supportedCurves {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			if group == preferredGroup {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>				selectedGroup = group
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				break
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if selectedGroup == 0 {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		c.sendAlert(alertHandshakeFailure)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		return errors.New(&#34;tls: no ECDHE curve supported by both client and server&#34;)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if clientKeyShare == nil {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		if err := hs.doHelloRetryRequest(selectedGroup); err != nil {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			return err
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		clientKeyShare = &amp;hs.clientHello.keyShares[0]
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if _, ok := curveForCurveID(selectedGroup); !ok {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		return errors.New(&#34;tls: CurvePreferences includes unsupported curve&#34;)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	key, err := generateECDHEKey(c.config.rand(), selectedGroup)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	if err != nil {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return err
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	hs.hello.serverShare = keyShare{group: selectedGroup, data: key.PublicKey().Bytes()}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	peerKey, err := key.Curve().NewPublicKey(clientKeyShare.data)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	if err != nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid client key share&#34;)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	hs.sharedKey, err = key.ECDH(peerKey)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if err != nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid client key share&#34;)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	selectedProto, err := negotiateALPN(c.config.NextProtos, hs.clientHello.alpnProtocols, c.quic != nil)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if err != nil {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		c.sendAlert(alertNoApplicationProtocol)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		return err
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	c.clientProtocol = selectedProto
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		<span class="comment">// RFC 9001 Section 4.2: Clients MUST NOT offer TLS versions older than 1.3.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		for _, v := range hs.clientHello.supportedVersions {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			if v &lt; VersionTLS13 {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>				c.sendAlert(alertProtocolVersion)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				return errors.New(&#34;tls: client offered TLS version older than TLS 1.3&#34;)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		<span class="comment">// RFC 9001 Section 8.2.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		if hs.clientHello.quicTransportParameters == nil {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			c.sendAlert(alertMissingExtension)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			return errors.New(&#34;tls: client did not send a quic_transport_parameters extension&#34;)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		c.quicSetTransportParameters(hs.clientHello.quicTransportParameters)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	} else {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		if hs.clientHello.quicTransportParameters != nil {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			c.sendAlert(alertUnsupportedExtension)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			return errors.New(&#34;tls: client sent an unexpected quic_transport_parameters extension&#34;)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	c.serverName = hs.clientHello.serverName
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return nil
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) checkForResumption() error {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	c := hs.c
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if c.config.SessionTicketsDisabled {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		return nil
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	modeOK := false
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	for _, mode := range hs.clientHello.pskModes {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		if mode == pskModeDHE {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			modeOK = true
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			break
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if !modeOK {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return nil
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	if len(hs.clientHello.pskIdentities) != len(hs.clientHello.pskBinders) {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid or missing PSK binders&#34;)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if len(hs.clientHello.pskIdentities) == 0 {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		return nil
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	for i, identity := range hs.clientHello.pskIdentities {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		if i &gt;= maxClientPSKIdentities {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			break
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		var sessionState *SessionState
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		if c.config.UnwrapSession != nil {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			var err error
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			sessionState, err = c.config.UnwrapSession(identity.label, c.connectionStateLocked())
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			if err != nil {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				return err
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			if sessionState == nil {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				continue
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		} else {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>			plaintext := c.config.decryptTicket(identity.label, c.ticketKeys)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			if plaintext == nil {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				continue
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>			}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			var err error
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			sessionState, err = ParseSessionState(plaintext)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			if err != nil {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>				continue
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		if sessionState.version != VersionTLS13 {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			continue
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		createdAt := time.Unix(int64(sessionState.createdAt), 0)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		if c.config.time().Sub(createdAt) &gt; maxSessionTicketLifetime {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			continue
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		pskSuite := cipherSuiteTLS13ByID(sessionState.cipherSuite)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		if pskSuite == nil || pskSuite.hash != hs.suite.hash {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			continue
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		<span class="comment">// PSK connections don&#39;t re-establish client certificates, but carry</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		<span class="comment">// them over in the session ticket. Ensure the presence of client certs</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		<span class="comment">// in the ticket is consistent with the configured requirements.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		sessionHasClientCerts := len(sessionState.peerCertificates) != 0
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		needClientCerts := requiresClientCert(c.config.ClientAuth)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		if needClientCerts &amp;&amp; !sessionHasClientCerts {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			continue
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if sessionHasClientCerts &amp;&amp; c.config.ClientAuth == NoClientCert {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			continue
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		if sessionHasClientCerts &amp;&amp; c.config.time().After(sessionState.peerCertificates[0].NotAfter) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			continue
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		if sessionHasClientCerts &amp;&amp; c.config.ClientAuth &gt;= VerifyClientCertIfGiven &amp;&amp;
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			len(sessionState.verifiedChains) == 0 {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			continue
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		hs.earlySecret = hs.suite.extract(sessionState.secret, nil)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		binderKey := hs.suite.deriveSecret(hs.earlySecret, resumptionBinderLabel, nil)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		<span class="comment">// Clone the transcript in case a HelloRetryRequest was recorded.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		transcript := cloneHash(hs.transcript, hs.suite.hash)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		if transcript == nil {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			return errors.New(&#34;tls: internal error: failed to clone hash&#34;)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		clientHelloBytes, err := hs.clientHello.marshalWithoutBinders()
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		if err != nil {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			return err
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		transcript.Write(clientHelloBytes)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		pskBinder := hs.suite.finishedHash(binderKey, transcript)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		if !hmac.Equal(hs.clientHello.pskBinders[i], pskBinder) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			c.sendAlert(alertDecryptError)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			return errors.New(&#34;tls: invalid PSK binder&#34;)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		if c.quic != nil &amp;&amp; hs.clientHello.earlyData &amp;&amp; i == 0 &amp;&amp;
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			sessionState.EarlyData &amp;&amp; sessionState.cipherSuite == hs.suite.id &amp;&amp;
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			sessionState.alpnProtocol == c.clientProtocol {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			hs.earlyData = true
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			transcript := hs.suite.hash.New()
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			if err := transcriptMsg(hs.clientHello, transcript); err != nil {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>				return err
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			earlyTrafficSecret := hs.suite.deriveSecret(hs.earlySecret, clientEarlyTrafficLabel, transcript)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			c.quicSetReadSecret(QUICEncryptionLevelEarly, hs.suite.id, earlyTrafficSecret)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		c.didResume = true
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		c.peerCertificates = sessionState.peerCertificates
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		c.ocspResponse = sessionState.ocspResponse
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		c.scts = sessionState.scts
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		c.verifiedChains = sessionState.verifiedChains
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		hs.hello.selectedIdentityPresent = true
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		hs.hello.selectedIdentity = uint16(i)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		hs.usingPSK = true
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		return nil
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	return nil
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">// cloneHash uses the encoding.BinaryMarshaler and encoding.BinaryUnmarshaler</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// interfaces implemented by standard library hashes to clone the state of in</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// to a new instance of h. It returns nil if the operation fails.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>func cloneHash(in hash.Hash, h crypto.Hash) hash.Hash {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	<span class="comment">// Recreate the interface to avoid importing encoding.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	type binaryMarshaler interface {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		MarshalBinary() (data []byte, err error)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		UnmarshalBinary(data []byte) error
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	marshaler, ok := in.(binaryMarshaler)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if !ok {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		return nil
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	state, err := marshaler.MarshalBinary()
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if err != nil {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		return nil
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	out := h.New()
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	unmarshaler, ok := out.(binaryMarshaler)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if !ok {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		return nil
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	if err := unmarshaler.UnmarshalBinary(state); err != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		return nil
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	return out
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) pickCertificate() error {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	c := hs.c
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">// Only one of PSK and certificates are used at a time.</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	if hs.usingPSK {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return nil
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// signature_algorithms is required in TLS 1.3. See RFC 8446, Section 4.2.3.</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	if len(hs.clientHello.supportedSignatureAlgorithms) == 0 {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		return c.sendAlert(alertMissingExtension)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	certificate, err := c.config.getCertificate(clientHelloInfo(hs.ctx, c, hs.clientHello))
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if err != nil {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		if err == errNoCertificates {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			c.sendAlert(alertUnrecognizedName)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		} else {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		return err
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	hs.sigAlg, err = selectSignatureScheme(c.vers, certificate, hs.clientHello.supportedSignatureAlgorithms)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	if err != nil {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		<span class="comment">// getCertificate returned a certificate that is unsupported or</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		<span class="comment">// incompatible with the client&#39;s signature algorithms.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		c.sendAlert(alertHandshakeFailure)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		return err
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	hs.cert = certificate
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	return nil
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// sendDummyChangeCipherSpec sends a ChangeCipherSpec record for compatibility</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// with middleboxes that didn&#39;t implement TLS correctly. See RFC 8446, Appendix D.4.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) sendDummyChangeCipherSpec() error {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	if hs.c.quic != nil {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		return nil
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	if hs.sentDummyCCS {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		return nil
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	hs.sentDummyCCS = true
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	return hs.c.writeChangeCipherRecord()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) doHelloRetryRequest(selectedGroup CurveID) error {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	c := hs.c
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// The first ClientHello gets double-hashed into the transcript upon a</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// HelloRetryRequest. See RFC 8446, Section 4.4.1.</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	if err := transcriptMsg(hs.clientHello, hs.transcript); err != nil {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		return err
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	chHash := hs.transcript.Sum(nil)
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	hs.transcript.Reset()
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	hs.transcript.Write([]byte{typeMessageHash, 0, 0, uint8(len(chHash))})
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	hs.transcript.Write(chHash)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	helloRetryRequest := &amp;serverHelloMsg{
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		vers:              hs.hello.vers,
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		random:            helloRetryRequestRandom,
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		sessionId:         hs.hello.sessionId,
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		cipherSuite:       hs.hello.cipherSuite,
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		compressionMethod: hs.hello.compressionMethod,
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		supportedVersion:  hs.hello.supportedVersion,
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		selectedGroup:     selectedGroup,
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(helloRetryRequest, hs.transcript); err != nil {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		return err
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	if err := hs.sendDummyChangeCipherSpec(); err != nil {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		return err
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">// clientHelloMsg is not included in the transcript.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	msg, err := c.readHandshake(nil)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if err != nil {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		return err
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	clientHello, ok := msg.(*clientHelloMsg)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if !ok {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		return unexpectedMessageError(clientHello, msg)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	if len(clientHello.keyShares) != 1 || clientHello.keyShares[0].group != selectedGroup {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		return errors.New(&#34;tls: client sent invalid key share in second ClientHello&#34;)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	if clientHello.earlyData {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		return errors.New(&#34;tls: client indicated early data in second ClientHello&#34;)
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	if illegalClientHelloChange(clientHello, hs.clientHello) {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		c.sendAlert(alertIllegalParameter)
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		return errors.New(&#34;tls: client illegally modified second ClientHello&#34;)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	hs.clientHello = clientHello
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	return nil
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// illegalClientHelloChange reports whether the two ClientHello messages are</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// different, with the exception of the changes allowed before and after a</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// HelloRetryRequest. See RFC 8446, Section 4.1.2.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>func illegalClientHelloChange(ch, ch1 *clientHelloMsg) bool {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if len(ch.supportedVersions) != len(ch1.supportedVersions) ||
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		len(ch.cipherSuites) != len(ch1.cipherSuites) ||
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		len(ch.supportedCurves) != len(ch1.supportedCurves) ||
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		len(ch.supportedSignatureAlgorithms) != len(ch1.supportedSignatureAlgorithms) ||
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		len(ch.supportedSignatureAlgorithmsCert) != len(ch1.supportedSignatureAlgorithmsCert) ||
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		len(ch.alpnProtocols) != len(ch1.alpnProtocols) {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		return true
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	for i := range ch.supportedVersions {
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		if ch.supportedVersions[i] != ch1.supportedVersions[i] {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			return true
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	}
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	for i := range ch.cipherSuites {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		if ch.cipherSuites[i] != ch1.cipherSuites[i] {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			return true
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	for i := range ch.supportedCurves {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		if ch.supportedCurves[i] != ch1.supportedCurves[i] {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			return true
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	for i := range ch.supportedSignatureAlgorithms {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		if ch.supportedSignatureAlgorithms[i] != ch1.supportedSignatureAlgorithms[i] {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			return true
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	for i := range ch.supportedSignatureAlgorithmsCert {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		if ch.supportedSignatureAlgorithmsCert[i] != ch1.supportedSignatureAlgorithmsCert[i] {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>			return true
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	for i := range ch.alpnProtocols {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		if ch.alpnProtocols[i] != ch1.alpnProtocols[i] {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			return true
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	return ch.vers != ch1.vers ||
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		!bytes.Equal(ch.random, ch1.random) ||
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		!bytes.Equal(ch.sessionId, ch1.sessionId) ||
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		!bytes.Equal(ch.compressionMethods, ch1.compressionMethods) ||
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		ch.serverName != ch1.serverName ||
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		ch.ocspStapling != ch1.ocspStapling ||
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		!bytes.Equal(ch.supportedPoints, ch1.supportedPoints) ||
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		ch.ticketSupported != ch1.ticketSupported ||
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		!bytes.Equal(ch.sessionTicket, ch1.sessionTicket) ||
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		ch.secureRenegotiationSupported != ch1.secureRenegotiationSupported ||
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		!bytes.Equal(ch.secureRenegotiation, ch1.secureRenegotiation) ||
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		ch.scts != ch1.scts ||
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		!bytes.Equal(ch.cookie, ch1.cookie) ||
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		!bytes.Equal(ch.pskModes, ch1.pskModes)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) sendServerParameters() error {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	c := hs.c
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	if err := transcriptMsg(hs.clientHello, hs.transcript); err != nil {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		return err
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(hs.hello, hs.transcript); err != nil {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		return err
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	if err := hs.sendDummyChangeCipherSpec(); err != nil {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		return err
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	earlySecret := hs.earlySecret
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	if earlySecret == nil {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		earlySecret = hs.suite.extract(nil, nil)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	hs.handshakeSecret = hs.suite.extract(hs.sharedKey,
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		hs.suite.deriveSecret(earlySecret, &#34;derived&#34;, nil))
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	clientSecret := hs.suite.deriveSecret(hs.handshakeSecret,
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		clientHandshakeTrafficLabel, hs.transcript)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	c.in.setTrafficSecret(hs.suite, QUICEncryptionLevelHandshake, clientSecret)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	serverSecret := hs.suite.deriveSecret(hs.handshakeSecret,
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		serverHandshakeTrafficLabel, hs.transcript)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	c.out.setTrafficSecret(hs.suite, QUICEncryptionLevelHandshake, serverSecret)
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		if c.hand.Len() != 0 {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			c.sendAlert(alertUnexpectedMessage)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		c.quicSetWriteSecret(QUICEncryptionLevelHandshake, hs.suite.id, serverSecret)
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>		c.quicSetReadSecret(QUICEncryptionLevelHandshake, hs.suite.id, clientSecret)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	err := c.config.writeKeyLog(keyLogLabelClientHandshake, hs.clientHello.random, clientSecret)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	if err != nil {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		return err
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelServerHandshake, hs.clientHello.random, serverSecret)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	if err != nil {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		return err
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	encryptedExtensions := new(encryptedExtensionsMsg)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	encryptedExtensions.alpnProtocol = c.clientProtocol
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		p, err := c.quicGetTransportParameters()
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		if err != nil {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>			return err
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		encryptedExtensions.quicTransportParameters = p
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		encryptedExtensions.earlyData = hs.earlyData
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(encryptedExtensions, hs.transcript); err != nil {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		return err
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	return nil
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) requestClientCert() bool {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	return hs.c.config.ClientAuth &gt;= RequestClientCert &amp;&amp; !hs.usingPSK
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) sendServerCertificate() error {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	c := hs.c
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// Only one of PSK and certificates are used at a time.</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	if hs.usingPSK {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		return nil
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	if hs.requestClientCert() {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		<span class="comment">// Request a client certificate</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		certReq := new(certificateRequestMsgTLS13)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		certReq.ocspStapling = true
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		certReq.scts = true
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		certReq.supportedSignatureAlgorithms = supportedSignatureAlgorithms()
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		if c.config.ClientCAs != nil {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			certReq.certificateAuthorities = c.config.ClientCAs.Subjects()
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		if _, err := hs.c.writeHandshakeRecord(certReq, hs.transcript); err != nil {
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			return err
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	certMsg := new(certificateMsgTLS13)
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	certMsg.certificate = *hs.cert
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	certMsg.scts = hs.clientHello.scts &amp;&amp; len(hs.cert.SignedCertificateTimestamps) &gt; 0
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	certMsg.ocspStapling = hs.clientHello.ocspStapling &amp;&amp; len(hs.cert.OCSPStaple) &gt; 0
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(certMsg, hs.transcript); err != nil {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		return err
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	certVerifyMsg := new(certificateVerifyMsg)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	certVerifyMsg.hasSignatureAlgorithm = true
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	certVerifyMsg.signatureAlgorithm = hs.sigAlg
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	sigType, sigHash, err := typeAndHashFromSignatureScheme(hs.sigAlg)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	if err != nil {
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		return c.sendAlert(alertInternalError)
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	}
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	signed := signedMessage(sigHash, serverSignatureContext, hs.transcript)
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	signOpts := crypto.SignerOpts(sigHash)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	if sigType == signatureRSAPSS {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		signOpts = &amp;rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: sigHash}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	sig, err := hs.cert.PrivateKey.(crypto.Signer).Sign(c.config.rand(), signed, signOpts)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	if err != nil {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		public := hs.cert.PrivateKey.(crypto.Signer).Public()
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		if rsaKey, ok := public.(*rsa.PublicKey); ok &amp;&amp; sigType == signatureRSAPSS &amp;&amp;
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			rsaKey.N.BitLen()/8 &lt; sigHash.Size()*2+2 { <span class="comment">// key too small for RSA-PSS</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			c.sendAlert(alertHandshakeFailure)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		} else {
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		return errors.New(&#34;tls: failed to sign handshake: &#34; + err.Error())
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	certVerifyMsg.signature = sig
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(certVerifyMsg, hs.transcript); err != nil {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		return err
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	return nil
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) sendServerFinished() error {
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	c := hs.c
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	finished := &amp;finishedMsg{
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		verifyData: hs.suite.finishedHash(c.out.trafficSecret, hs.transcript),
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	if _, err := hs.c.writeHandshakeRecord(finished, hs.transcript); err != nil {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		return err
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>	}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// Derive secrets that take context through the server Finished.</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	hs.masterSecret = hs.suite.extract(nil,
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		hs.suite.deriveSecret(hs.handshakeSecret, &#34;derived&#34;, nil))
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	hs.trafficSecret = hs.suite.deriveSecret(hs.masterSecret,
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>		clientApplicationTrafficLabel, hs.transcript)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	serverSecret := hs.suite.deriveSecret(hs.masterSecret,
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		serverApplicationTrafficLabel, hs.transcript)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	c.out.setTrafficSecret(hs.suite, QUICEncryptionLevelApplication, serverSecret)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	if c.quic != nil {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		if c.hand.Len() != 0 {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			<span class="comment">// TODO: Handle this in setTrafficSecret?</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>			c.sendAlert(alertUnexpectedMessage)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		c.quicSetWriteSecret(QUICEncryptionLevelApplication, hs.suite.id, serverSecret)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	err := c.config.writeKeyLog(keyLogLabelClientTraffic, hs.clientHello.random, hs.trafficSecret)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	if err != nil {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		return err
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	}
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	err = c.config.writeKeyLog(keyLogLabelServerTraffic, hs.clientHello.random, serverSecret)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	if err != nil {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		c.sendAlert(alertInternalError)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		return err
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	c.ekm = hs.suite.exportKeyingMaterial(hs.masterSecret, hs.transcript)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	<span class="comment">// If we did not request client certificates, at this point we can</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	<span class="comment">// precompute the client finished and roll the transcript forward to send</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	<span class="comment">// session tickets in our first flight.</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	if !hs.requestClientCert() {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		if err := hs.sendSessionTickets(); err != nil {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			return err
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	}
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	return nil
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>}
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) shouldSendSessionTickets() bool {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	if hs.c.config.SessionTicketsDisabled {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>		return false
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	<span class="comment">// QUIC tickets are sent by QUICConn.SendSessionTicket, not automatically.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	if hs.c.quic != nil {
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		return false
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t send tickets the client wouldn&#39;t use. See RFC 8446, Section 4.2.9.</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	for _, pskMode := range hs.clientHello.pskModes {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		if pskMode == pskModeDHE {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			return true
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	}
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	return false
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) sendSessionTickets() error {
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	c := hs.c
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	hs.clientFinished = hs.suite.finishedHash(c.in.trafficSecret, hs.transcript)
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	finishedMsg := &amp;finishedMsg{
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		verifyData: hs.clientFinished,
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	if err := transcriptMsg(finishedMsg, hs.transcript); err != nil {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		return err
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	c.resumptionSecret = hs.suite.deriveSecret(hs.masterSecret,
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		resumptionLabel, hs.transcript)
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	if !hs.shouldSendSessionTickets() {
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		return nil
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	}
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	return c.sendSessionTicket(false)
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>func (c *Conn) sendSessionTicket(earlyData bool) error {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	suite := cipherSuiteTLS13ByID(c.cipherSuite)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	if suite == nil {
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		return errors.New(&#34;tls: internal error: unknown cipher suite&#34;)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	<span class="comment">// ticket_nonce, which must be unique per connection, is always left at</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	<span class="comment">// zero because we only ever send one ticket per connection.</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	psk := suite.expandLabel(c.resumptionSecret, &#34;resumption&#34;,
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		nil, suite.hash.Size())
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	m := new(newSessionTicketMsgTLS13)
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	state, err := c.sessionState()
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	if err != nil {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		return err
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	state.secret = psk
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	state.EarlyData = earlyData
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	if c.config.WrapSession != nil {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		m.label, err = c.config.WrapSession(c.connectionStateLocked(), state)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		if err != nil {
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			return err
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	} else {
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		stateBytes, err := state.Bytes()
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		if err != nil {
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>			c.sendAlert(alertInternalError)
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			return err
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		m.label, err = c.config.encryptTicket(stateBytes, c.ticketKeys)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		if err != nil {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>			return err
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	m.lifetime = uint32(maxSessionTicketLifetime / time.Second)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	<span class="comment">// ticket_age_add is a random 32-bit value. See RFC 8446, section 4.6.1</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	<span class="comment">// The value is not stored anywhere; we never need to check the ticket age</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// because 0-RTT is not supported.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	ageAdd := make([]byte, 4)
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	_, err = c.config.rand().Read(ageAdd)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	if err != nil {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		return err
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	m.ageAdd = binary.LittleEndian.Uint32(ageAdd)
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	if earlyData {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		<span class="comment">// RFC 9001, Section 4.6.1</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		m.maxEarlyData = 0xffffffff
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	}
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	if _, err := c.writeHandshakeRecord(m, nil); err != nil {
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		return err
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	return nil
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) readClientCertificate() error {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	c := hs.c
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>	if !hs.requestClientCert() {
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		<span class="comment">// Make sure the connection is still being verified whether or not</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		<span class="comment">// the server requested a client certificate.</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		if c.config.VerifyConnection != nil {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			if err := c.config.VerifyConnection(c.connectionStateLocked()); err != nil {
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>				c.sendAlert(alertBadCertificate)
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>				return err
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>			}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		return nil
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	}
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	<span class="comment">// If we requested a client certificate, then the client must send a</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	<span class="comment">// certificate message. If it&#39;s empty, no CertificateVerify is sent.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	msg, err := c.readHandshake(hs.transcript)
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	if err != nil {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		return err
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	certMsg, ok := msg.(*certificateMsgTLS13)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	if !ok {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		return unexpectedMessageError(certMsg, msg)
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	}
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	if err := c.processCertsFromClient(certMsg.certificate); err != nil {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		return err
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	if c.config.VerifyConnection != nil {
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		if err := c.config.VerifyConnection(c.connectionStateLocked()); err != nil {
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>			c.sendAlert(alertBadCertificate)
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			return err
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>		}
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	}
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	if len(certMsg.certificate.Certificate) != 0 {
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		<span class="comment">// certificateVerifyMsg is included in the transcript, but not until</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		<span class="comment">// after we verify the handshake signature, since the state before</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		<span class="comment">// this message was sent is used.</span>
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		msg, err = c.readHandshake(nil)
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		if err != nil {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			return err
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		certVerify, ok := msg.(*certificateVerifyMsg)
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		if !ok {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>			c.sendAlert(alertUnexpectedMessage)
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>			return unexpectedMessageError(certVerify, msg)
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		}
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		<span class="comment">// See RFC 8446, Section 4.4.3.</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		if !isSupportedSignatureAlgorithm(certVerify.signatureAlgorithm, supportedSignatureAlgorithms()) {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			c.sendAlert(alertIllegalParameter)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			return errors.New(&#34;tls: client certificate used with invalid signature algorithm&#34;)
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>		sigType, sigHash, err := typeAndHashFromSignatureScheme(certVerify.signatureAlgorithm)
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		if err != nil {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			return c.sendAlert(alertInternalError)
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		}
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		if sigType == signaturePKCS1v15 || sigHash == crypto.SHA1 {
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			c.sendAlert(alertIllegalParameter)
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			return errors.New(&#34;tls: client certificate used with invalid signature algorithm&#34;)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>		}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>		signed := signedMessage(sigHash, clientSignatureContext, hs.transcript)
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		if err := verifyHandshakeSignature(sigType, c.peerCertificates[0].PublicKey,
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			sigHash, signed, certVerify.signature); err != nil {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			c.sendAlert(alertDecryptError)
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			return errors.New(&#34;tls: invalid signature by the client certificate: &#34; + err.Error())
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		}
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>		if err := transcriptMsg(certVerify, hs.transcript); err != nil {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			return err
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>		}
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>	}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>	<span class="comment">// If we waited until the client certificates to send session tickets, we</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>	<span class="comment">// are ready to do it now.</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>	if err := hs.sendSessionTickets(); err != nil {
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>		return err
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>	}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>	return nil
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>}
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>func (hs *serverHandshakeStateTLS13) readClientFinished() error {
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	c := hs.c
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	<span class="comment">// finishedMsg is not included in the transcript.</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>	msg, err := c.readHandshake(nil)
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	if err != nil {
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		return err
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	}
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	finished, ok := msg.(*finishedMsg)
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	if !ok {
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>		c.sendAlert(alertUnexpectedMessage)
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>		return unexpectedMessageError(finished, msg)
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	if !hmac.Equal(hs.clientFinished, finished.verifyData) {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		c.sendAlert(alertDecryptError)
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid client finished hash&#34;)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	c.in.setTrafficSecret(hs.suite, QUICEncryptionLevelApplication, hs.trafficSecret)
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	return nil
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>
</pre><p><a href="handshake_server_tls13.go?m=text">View as plain text</a></p>

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
