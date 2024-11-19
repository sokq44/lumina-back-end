<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/handshake_messages.go - Go Documentation Server</title>

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
<a href="handshake_messages.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">handshake_messages.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/tls">crypto/tls</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package tls
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/cryptobyte&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The marshalingFunction type is an adapter to allow the use of ordinary</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// functions as cryptobyte.MarshalingValue.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type marshalingFunction func(b *cryptobyte.Builder) error
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>func (f marshalingFunction) Marshal(b *cryptobyte.Builder) error {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	return f(b)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// addBytesWithLength appends a sequence of bytes to the cryptobyte.Builder. If</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// the length of the sequence is not the value specified, it produces an error.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func addBytesWithLength(b *cryptobyte.Builder, v []byte, n int) {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	b.AddValue(marshalingFunction(func(b *cryptobyte.Builder) error {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		if len(v) != n {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;invalid value length: expected %d, got %d&#34;, n, len(v))
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		b.AddBytes(v)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		return nil
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}))
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// addUint64 appends a big-endian, 64-bit value to the cryptobyte.Builder.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func addUint64(b *cryptobyte.Builder, v uint64) {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	b.AddUint32(uint32(v &gt;&gt; 32))
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	b.AddUint32(uint32(v))
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// readUint64 decodes a big-endian, 64-bit value into out and advances over it.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// It reports whether the read was successful.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func readUint64(s *cryptobyte.String, out *uint64) bool {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	var hi, lo uint32
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	if !s.ReadUint32(&amp;hi) || !s.ReadUint32(&amp;lo) {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		return false
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	*out = uint64(hi)&lt;&lt;32 | uint64(lo)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	return true
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// readUint8LengthPrefixed acts like s.ReadUint8LengthPrefixed, but targets a</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// []byte instead of a cryptobyte.String.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func readUint8LengthPrefixed(s *cryptobyte.String, out *[]byte) bool {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return s.ReadUint8LengthPrefixed((*cryptobyte.String)(out))
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// readUint16LengthPrefixed acts like s.ReadUint16LengthPrefixed, but targets a</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// []byte instead of a cryptobyte.String.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func readUint16LengthPrefixed(s *cryptobyte.String, out *[]byte) bool {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	return s.ReadUint16LengthPrefixed((*cryptobyte.String)(out))
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// readUint24LengthPrefixed acts like s.ReadUint24LengthPrefixed, but targets a</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// []byte instead of a cryptobyte.String.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func readUint24LengthPrefixed(s *cryptobyte.String, out *[]byte) bool {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return s.ReadUint24LengthPrefixed((*cryptobyte.String)(out))
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>type clientHelloMsg struct {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	raw                              []byte
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	vers                             uint16
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	random                           []byte
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	sessionId                        []byte
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	cipherSuites                     []uint16
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	compressionMethods               []uint8
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	serverName                       string
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	ocspStapling                     bool
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	supportedCurves                  []CurveID
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	supportedPoints                  []uint8
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	ticketSupported                  bool
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	sessionTicket                    []uint8
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	supportedSignatureAlgorithms     []SignatureScheme
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	supportedSignatureAlgorithmsCert []SignatureScheme
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	secureRenegotiationSupported     bool
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	secureRenegotiation              []byte
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	extendedMasterSecret             bool
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	alpnProtocols                    []string
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	scts                             bool
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	supportedVersions                []uint16
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	cookie                           []byte
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	keyShares                        []keyShare
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	earlyData                        bool
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	pskModes                         []uint8
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	pskIdentities                    []pskIdentity
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	pskBinders                       [][]byte
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	quicTransportParameters          []byte
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (m *clientHelloMsg) marshal() ([]byte, error) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	var exts cryptobyte.Builder
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if len(m.serverName) &gt; 0 {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// RFC 6066, Section 3</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		exts.AddUint16(extensionServerName)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>				exts.AddUint8(0) <span class="comment">// name_type = host_name</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>				exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>					exts.AddBytes([]byte(m.serverName))
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				})
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			})
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		})
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if m.ocspStapling {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		<span class="comment">// RFC 4366, Section 3.6</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		exts.AddUint16(extensionStatusRequest)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			exts.AddUint8(1)  <span class="comment">// status_type = ocsp</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			exts.AddUint16(0) <span class="comment">// empty responder_id_list</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			exts.AddUint16(0) <span class="comment">// empty request_extensions</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		})
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if len(m.supportedCurves) &gt; 0 {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		<span class="comment">// RFC 4492, sections 5.1.1 and RFC 8446, Section 4.2.7</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		exts.AddUint16(extensionSupportedCurves)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>				for _, curve := range m.supportedCurves {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>					exts.AddUint16(uint16(curve))
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>				}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			})
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		})
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if len(m.supportedPoints) &gt; 0 {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// RFC 4492, Section 5.1.2</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		exts.AddUint16(extensionSupportedPoints)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				exts.AddBytes(m.supportedPoints)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			})
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		})
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if m.ticketSupported {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		<span class="comment">// RFC 5077, Section 3.2</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		exts.AddUint16(extensionSessionTicket)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			exts.AddBytes(m.sessionTicket)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		})
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if len(m.supportedSignatureAlgorithms) &gt; 0 {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// RFC 5246, Section 7.4.1.4.1</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		exts.AddUint16(extensionSignatureAlgorithms)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>				for _, sigAlgo := range m.supportedSignatureAlgorithms {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>					exts.AddUint16(uint16(sigAlgo))
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>				}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			})
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		})
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if len(m.supportedSignatureAlgorithmsCert) &gt; 0 {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.3</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		exts.AddUint16(extensionSignatureAlgorithmsCert)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				for _, sigAlgo := range m.supportedSignatureAlgorithmsCert {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>					exts.AddUint16(uint16(sigAlgo))
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			})
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		})
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if m.secureRenegotiationSupported {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		<span class="comment">// RFC 5746, Section 3.2</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		exts.AddUint16(extensionRenegotiationInfo)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>				exts.AddBytes(m.secureRenegotiation)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			})
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		})
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if m.extendedMasterSecret {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">// RFC 7627</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		exts.AddUint16(extensionExtendedMasterSecret)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if len(m.alpnProtocols) &gt; 0 {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		<span class="comment">// RFC 7301, Section 3.1</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		exts.AddUint16(extensionALPN)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				for _, proto := range m.alpnProtocols {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>					exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>						exts.AddBytes([]byte(proto))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					})
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			})
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		})
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if m.scts {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		<span class="comment">// RFC 6962, Section 3.3.1</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		exts.AddUint16(extensionSCT)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if len(m.supportedVersions) &gt; 0 {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.1</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		exts.AddUint16(extensionSupportedVersions)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				for _, vers := range m.supportedVersions {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>					exts.AddUint16(vers)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			})
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		})
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	if len(m.cookie) &gt; 0 {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.2</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		exts.AddUint16(extensionCookie)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				exts.AddBytes(m.cookie)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			})
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		})
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if len(m.keyShares) &gt; 0 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.8</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		exts.AddUint16(extensionKeyShare)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				for _, ks := range m.keyShares {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>					exts.AddUint16(uint16(ks.group))
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>					exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>						exts.AddBytes(ks.data)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>					})
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			})
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		})
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if m.earlyData {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.10</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		exts.AddUint16(extensionEarlyData)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	if len(m.pskModes) &gt; 0 {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.9</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		exts.AddUint16(extensionPSKModes)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>				exts.AddBytes(m.pskModes)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			})
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		})
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	if m.quicTransportParameters != nil { <span class="comment">// marshal zero-length parameters when present</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// RFC 9001, Section 8.2</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		exts.AddUint16(extensionQUICTransportParameters)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			exts.AddBytes(m.quicTransportParameters)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		})
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	if len(m.pskIdentities) &gt; 0 { <span class="comment">// pre_shared_key must be the last extension</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		<span class="comment">// RFC 8446, Section 4.2.11</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		exts.AddUint16(extensionPreSharedKey)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				for _, psk := range m.pskIdentities {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>					exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>						exts.AddBytes(psk.label)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>					})
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>					exts.AddUint32(psk.obfuscatedTicketAge)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			})
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>				for _, binder := range m.pskBinders {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>					exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>						exts.AddBytes(binder)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>					})
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			})
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		})
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	extBytes, err := exts.Bytes()
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	if err != nil {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		return nil, err
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	b.AddUint8(typeClientHello)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		b.AddUint16(m.vers)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		addBytesWithLength(b, m.random, 32)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			b.AddBytes(m.sessionId)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		})
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			for _, suite := range m.cipherSuites {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				b.AddUint16(suite)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		})
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			b.AddBytes(m.compressionMethods)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		})
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		if len(extBytes) &gt; 0 {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				b.AddBytes(extBytes)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>			})
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	})
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	return m.raw, err
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// marshalWithoutBinders returns the ClientHello through the</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// PreSharedKeyExtension.identities field, according to RFC 8446, Section</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// 4.2.11.2. Note that m.pskBinders must be set to slices of the correct length.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func (m *clientHelloMsg) marshalWithoutBinders() ([]byte, error) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	bindersLen := 2 <span class="comment">// uint16 length prefix</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	for _, binder := range m.pskBinders {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		bindersLen += 1 <span class="comment">// uint8 length prefix</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		bindersLen += len(binder)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	fullMessage, err := m.marshal()
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if err != nil {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return nil, err
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	return fullMessage[:len(fullMessage)-bindersLen], nil
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// updateBinders updates the m.pskBinders field, if necessary updating the</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// cached marshaled representation. The supplied binders must have the same</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// length as the current m.pskBinders.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>func (m *clientHelloMsg) updateBinders(pskBinders [][]byte) error {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	if len(pskBinders) != len(m.pskBinders) {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		return errors.New(&#34;tls: internal error: pskBinders length mismatch&#34;)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	for i := range m.pskBinders {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		if len(pskBinders[i]) != len(m.pskBinders[i]) {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			return errors.New(&#34;tls: internal error: pskBinders length mismatch&#34;)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	m.pskBinders = pskBinders
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		helloBytes, err := m.marshalWithoutBinders()
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		if err != nil {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			return err
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		lenWithoutBinders := len(helloBytes)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		b := cryptobyte.NewFixedBuilder(m.raw[:lenWithoutBinders])
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			for _, binder := range m.pskBinders {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>					b.AddBytes(binder)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>				})
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		})
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		if out, err := b.Bytes(); err != nil || len(out) != len(m.raw) {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			return errors.New(&#34;tls: internal error: failed to update binders&#34;)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	return nil
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (m *clientHelloMsg) unmarshal(data []byte) bool {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	*m = clientHelloMsg{raw: data}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		!s.ReadUint16(&amp;m.vers) || !s.ReadBytes(&amp;m.random, 32) ||
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		!readUint8LengthPrefixed(&amp;s, &amp;m.sessionId) {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		return false
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	var cipherSuites cryptobyte.String
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	if !s.ReadUint16LengthPrefixed(&amp;cipherSuites) {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		return false
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	m.cipherSuites = []uint16{}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	m.secureRenegotiationSupported = false
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	for !cipherSuites.Empty() {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		var suite uint16
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		if !cipherSuites.ReadUint16(&amp;suite) {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			return false
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		if suite == scsvRenegotiation {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			m.secureRenegotiationSupported = true
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		m.cipherSuites = append(m.cipherSuites, suite)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	if !readUint8LengthPrefixed(&amp;s, &amp;m.compressionMethods) {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		return false
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if s.Empty() {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		<span class="comment">// ClientHello is optionally followed by extension data</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		return true
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	var extensions cryptobyte.String
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if !s.ReadUint16LengthPrefixed(&amp;extensions) || !s.Empty() {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		return false
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	seenExts := make(map[uint16]bool)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	for !extensions.Empty() {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		var extension uint16
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		var extData cryptobyte.String
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		if !extensions.ReadUint16(&amp;extension) ||
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			return false
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		if seenExts[extension] {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			return false
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		seenExts[extension] = true
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		switch extension {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		case extensionServerName:
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			<span class="comment">// RFC 6066, Section 3</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			var nameList cryptobyte.String
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;nameList) || nameList.Empty() {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				return false
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			for !nameList.Empty() {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>				var nameType uint8
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>				var serverName cryptobyte.String
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>				if !nameList.ReadUint8(&amp;nameType) ||
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>					!nameList.ReadUint16LengthPrefixed(&amp;serverName) ||
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>					serverName.Empty() {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>					return false
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>				}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				if nameType != 0 {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>					continue
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>				}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>				if len(m.serverName) != 0 {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>					<span class="comment">// Multiple names of the same name_type are prohibited.</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>					return false
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>				}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>				m.serverName = string(serverName)
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				<span class="comment">// An SNI value may not include a trailing dot.</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>				if strings.HasSuffix(m.serverName, &#34;.&#34;) {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>					return false
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		case extensionStatusRequest:
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			<span class="comment">// RFC 4366, Section 3.6</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			var statusType uint8
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			var ignored cryptobyte.String
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			if !extData.ReadUint8(&amp;statusType) ||
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>				!extData.ReadUint16LengthPrefixed(&amp;ignored) ||
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>				!extData.ReadUint16LengthPrefixed(&amp;ignored) {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				return false
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			m.ocspStapling = statusType == statusTypeOCSP
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		case extensionSupportedCurves:
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			<span class="comment">// RFC 4492, sections 5.1.1 and RFC 8446, Section 4.2.7</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			var curves cryptobyte.String
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;curves) || curves.Empty() {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				return false
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			for !curves.Empty() {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>				var curve uint16
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>				if !curves.ReadUint16(&amp;curve) {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>					return false
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>				m.supportedCurves = append(m.supportedCurves, CurveID(curve))
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		case extensionSupportedPoints:
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			<span class="comment">// RFC 4492, Section 5.1.2</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			if !readUint8LengthPrefixed(&amp;extData, &amp;m.supportedPoints) ||
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>				len(m.supportedPoints) == 0 {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>				return false
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		case extensionSessionTicket:
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			<span class="comment">// RFC 5077, Section 3.2</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			m.ticketSupported = true
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			extData.ReadBytes(&amp;m.sessionTicket, len(extData))
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		case extensionSignatureAlgorithms:
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			<span class="comment">// RFC 5246, Section 7.4.1.4.1</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			var sigAndAlgs cryptobyte.String
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;sigAndAlgs) || sigAndAlgs.Empty() {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				return false
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>			for !sigAndAlgs.Empty() {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>				var sigAndAlg uint16
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				if !sigAndAlgs.ReadUint16(&amp;sigAndAlg) {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>					return false
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>				}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				m.supportedSignatureAlgorithms = append(
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>					m.supportedSignatureAlgorithms, SignatureScheme(sigAndAlg))
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		case extensionSignatureAlgorithmsCert:
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.3</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			var sigAndAlgs cryptobyte.String
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;sigAndAlgs) || sigAndAlgs.Empty() {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>				return false
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			for !sigAndAlgs.Empty() {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>				var sigAndAlg uint16
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				if !sigAndAlgs.ReadUint16(&amp;sigAndAlg) {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>					return false
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>				}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>				m.supportedSignatureAlgorithmsCert = append(
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>					m.supportedSignatureAlgorithmsCert, SignatureScheme(sigAndAlg))
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		case extensionRenegotiationInfo:
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			<span class="comment">// RFC 5746, Section 3.2</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			if !readUint8LengthPrefixed(&amp;extData, &amp;m.secureRenegotiation) {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				return false
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			m.secureRenegotiationSupported = true
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		case extensionExtendedMasterSecret:
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			<span class="comment">// RFC 7627</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			m.extendedMasterSecret = true
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		case extensionALPN:
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>			<span class="comment">// RFC 7301, Section 3.1</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>			var protoList cryptobyte.String
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;protoList) || protoList.Empty() {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>				return false
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			for !protoList.Empty() {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				var proto cryptobyte.String
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				if !protoList.ReadUint8LengthPrefixed(&amp;proto) || proto.Empty() {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>					return false
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>				m.alpnProtocols = append(m.alpnProtocols, string(proto))
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		case extensionSCT:
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			<span class="comment">// RFC 6962, Section 3.3.1</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			m.scts = true
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		case extensionSupportedVersions:
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.1</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			var versList cryptobyte.String
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			if !extData.ReadUint8LengthPrefixed(&amp;versList) || versList.Empty() {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>				return false
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>			}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			for !versList.Empty() {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				var vers uint16
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>				if !versList.ReadUint16(&amp;vers) {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>					return false
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>				}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>				m.supportedVersions = append(m.supportedVersions, vers)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>			}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		case extensionCookie:
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.2</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			if !readUint16LengthPrefixed(&amp;extData, &amp;m.cookie) ||
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				len(m.cookie) == 0 {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				return false
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		case extensionKeyShare:
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.8</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			var clientShares cryptobyte.String
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;clientShares) {
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>				return false
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			for !clientShares.Empty() {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>				var ks keyShare
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>				if !clientShares.ReadUint16((*uint16)(&amp;ks.group)) ||
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>					!readUint16LengthPrefixed(&amp;clientShares, &amp;ks.data) ||
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>					len(ks.data) == 0 {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>					return false
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>				m.keyShares = append(m.keyShares, ks)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		case extensionEarlyData:
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.10</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			m.earlyData = true
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		case extensionPSKModes:
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.9</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			if !readUint8LengthPrefixed(&amp;extData, &amp;m.pskModes) {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>				return false
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		case extensionQUICTransportParameters:
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>			m.quicTransportParameters = make([]byte, len(extData))
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			if !extData.CopyBytes(m.quicTransportParameters) {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>				return false
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		case extensionPreSharedKey:
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.11</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			if !extensions.Empty() {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				return false <span class="comment">// pre_shared_key must be the last extension</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>			}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>			var identities cryptobyte.String
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;identities) || identities.Empty() {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>				return false
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			for !identities.Empty() {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>				var psk pskIdentity
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				if !readUint16LengthPrefixed(&amp;identities, &amp;psk.label) ||
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>					!identities.ReadUint32(&amp;psk.obfuscatedTicketAge) ||
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>					len(psk.label) == 0 {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>					return false
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>				}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>				m.pskIdentities = append(m.pskIdentities, psk)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>			}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			var binders cryptobyte.String
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;binders) || binders.Empty() {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>				return false
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			for !binders.Empty() {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>				var binder []byte
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>				if !readUint8LengthPrefixed(&amp;binders, &amp;binder) ||
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>					len(binder) == 0 {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>					return false
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>				}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>				m.pskBinders = append(m.pskBinders, binder)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			}
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		default:
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			<span class="comment">// Ignore unknown extensions.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			continue
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		if !extData.Empty() {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			return false
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	return true
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>type serverHelloMsg struct {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	raw                          []byte
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	vers                         uint16
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	random                       []byte
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	sessionId                    []byte
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	cipherSuite                  uint16
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	compressionMethod            uint8
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	ocspStapling                 bool
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	ticketSupported              bool
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	secureRenegotiationSupported bool
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	secureRenegotiation          []byte
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	extendedMasterSecret         bool
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	alpnProtocol                 string
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	scts                         [][]byte
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	supportedVersion             uint16
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	serverShare                  keyShare
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	selectedIdentityPresent      bool
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	selectedIdentity             uint16
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	supportedPoints              []uint8
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// HelloRetryRequest extensions</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	cookie        []byte
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	selectedGroup CurveID
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>func (m *serverHelloMsg) marshal() ([]byte, error) {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	}
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	var exts cryptobyte.Builder
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	if m.ocspStapling {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		exts.AddUint16(extensionStatusRequest)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	if m.ticketSupported {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		exts.AddUint16(extensionSessionTicket)
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	if m.secureRenegotiationSupported {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		exts.AddUint16(extensionRenegotiationInfo)
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>				exts.AddBytes(m.secureRenegotiation)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>			})
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		})
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	if m.extendedMasterSecret {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		exts.AddUint16(extensionExtendedMasterSecret)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		exts.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	if len(m.alpnProtocol) &gt; 0 {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		exts.AddUint16(extensionALPN)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>					exts.AddBytes([]byte(m.alpnProtocol))
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>				})
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			})
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		})
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	if len(m.scts) &gt; 0 {
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		exts.AddUint16(extensionSCT)
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>				for _, sct := range m.scts {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>					exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>						exts.AddBytes(sct)
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>					})
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>				}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			})
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		})
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	if m.supportedVersion != 0 {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		exts.AddUint16(extensionSupportedVersions)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			exts.AddUint16(m.supportedVersion)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		})
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	if m.serverShare.group != 0 {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		exts.AddUint16(extensionKeyShare)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			exts.AddUint16(uint16(m.serverShare.group))
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>				exts.AddBytes(m.serverShare.data)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>			})
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		})
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	if m.selectedIdentityPresent {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		exts.AddUint16(extensionPreSharedKey)
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			exts.AddUint16(m.selectedIdentity)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		})
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	if len(m.cookie) &gt; 0 {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		exts.AddUint16(extensionCookie)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>				exts.AddBytes(m.cookie)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			})
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		})
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	if m.selectedGroup != 0 {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		exts.AddUint16(extensionKeyShare)
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>			exts.AddUint16(uint16(m.selectedGroup))
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		})
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if len(m.supportedPoints) &gt; 0 {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		exts.AddUint16(extensionSupportedPoints)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		exts.AddUint16LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			exts.AddUint8LengthPrefixed(func(exts *cryptobyte.Builder) {
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>				exts.AddBytes(m.supportedPoints)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			})
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		})
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	extBytes, err := exts.Bytes()
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	if err != nil {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		return nil, err
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	b.AddUint8(typeServerHello)
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		b.AddUint16(m.vers)
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		addBytesWithLength(b, m.random, 32)
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			b.AddBytes(m.sessionId)
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		})
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		b.AddUint16(m.cipherSuite)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		b.AddUint8(m.compressionMethod)
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		if len(extBytes) &gt; 0 {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>				b.AddBytes(extBytes)
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			})
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	})
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	return m.raw, err
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>func (m *serverHelloMsg) unmarshal(data []byte) bool {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	*m = serverHelloMsg{raw: data}
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		!s.ReadUint16(&amp;m.vers) || !s.ReadBytes(&amp;m.random, 32) ||
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		!readUint8LengthPrefixed(&amp;s, &amp;m.sessionId) ||
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		!s.ReadUint16(&amp;m.cipherSuite) ||
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		!s.ReadUint8(&amp;m.compressionMethod) {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		return false
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	if s.Empty() {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		<span class="comment">// ServerHello is optionally followed by extension data</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		return true
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	var extensions cryptobyte.String
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	if !s.ReadUint16LengthPrefixed(&amp;extensions) || !s.Empty() {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		return false
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	seenExts := make(map[uint16]bool)
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	for !extensions.Empty() {
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		var extension uint16
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		var extData cryptobyte.String
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		if !extensions.ReadUint16(&amp;extension) ||
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			return false
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		if seenExts[extension] {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>			return false
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		seenExts[extension] = true
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		switch extension {
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		case extensionStatusRequest:
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>			m.ocspStapling = true
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		case extensionSessionTicket:
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>			m.ticketSupported = true
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		case extensionRenegotiationInfo:
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			if !readUint8LengthPrefixed(&amp;extData, &amp;m.secureRenegotiation) {
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>				return false
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>			}
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>			m.secureRenegotiationSupported = true
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		case extensionExtendedMasterSecret:
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			m.extendedMasterSecret = true
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>		case extensionALPN:
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>			var protoList cryptobyte.String
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;protoList) || protoList.Empty() {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>				return false
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			}
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>			var proto cryptobyte.String
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			if !protoList.ReadUint8LengthPrefixed(&amp;proto) ||
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>				proto.Empty() || !protoList.Empty() {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>				return false
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>			}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>			m.alpnProtocol = string(proto)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		case extensionSCT:
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			var sctList cryptobyte.String
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;sctList) || sctList.Empty() {
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>				return false
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>			}
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			for !sctList.Empty() {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>				var sct []byte
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>				if !readUint16LengthPrefixed(&amp;sctList, &amp;sct) ||
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>					len(sct) == 0 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>					return false
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>				}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>				m.scts = append(m.scts, sct)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			}
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>		case extensionSupportedVersions:
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			if !extData.ReadUint16(&amp;m.supportedVersion) {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>				return false
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		case extensionCookie:
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>			if !readUint16LengthPrefixed(&amp;extData, &amp;m.cookie) ||
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>				len(m.cookie) == 0 {
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>				return false
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>			}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		case extensionKeyShare:
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			<span class="comment">// This extension has different formats in SH and HRR, accept either</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			<span class="comment">// and let the handshake logic decide. See RFC 8446, Section 4.2.8.</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			if len(extData) == 2 {
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>				if !extData.ReadUint16((*uint16)(&amp;m.selectedGroup)) {
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>					return false
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>				}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			} else {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>				if !extData.ReadUint16((*uint16)(&amp;m.serverShare.group)) ||
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>					!readUint16LengthPrefixed(&amp;extData, &amp;m.serverShare.data) {
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>					return false
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>				}
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		case extensionPreSharedKey:
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			m.selectedIdentityPresent = true
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			if !extData.ReadUint16(&amp;m.selectedIdentity) {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>				return false
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			}
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		case extensionSupportedPoints:
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			<span class="comment">// RFC 4492, Section 5.1.2</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>			if !readUint8LengthPrefixed(&amp;extData, &amp;m.supportedPoints) ||
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>				len(m.supportedPoints) == 0 {
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>				return false
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			}
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		default:
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>			<span class="comment">// Ignore unknown extensions.</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			continue
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		}
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		if !extData.Empty() {
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>			return false
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		}
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	}
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	return true
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>}
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>type encryptedExtensionsMsg struct {
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	raw                     []byte
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	alpnProtocol            string
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	quicTransportParameters []byte
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	earlyData               bool
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>func (m *encryptedExtensionsMsg) marshal() ([]byte, error) {
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	b.AddUint8(typeEncryptedExtensions)
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>			if len(m.alpnProtocol) &gt; 0 {
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>				b.AddUint16(extensionALPN)
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>						b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>							b.AddBytes([]byte(m.alpnProtocol))
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>						})
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>					})
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>				})
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			}
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>			if m.quicTransportParameters != nil { <span class="comment">// marshal zero-length parameters when present</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>				<span class="comment">// draft-ietf-quic-tls-32, Section 8.2</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>				b.AddUint16(extensionQUICTransportParameters)
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>					b.AddBytes(m.quicTransportParameters)
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>				})
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			}
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			if m.earlyData {
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>				<span class="comment">// RFC 8446, Section 4.2.10</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>				b.AddUint16(extensionEarlyData)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>				b.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>			}
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		})
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	})
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>	var err error
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	return m.raw, err
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>func (m *encryptedExtensionsMsg) unmarshal(data []byte) bool {
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	*m = encryptedExtensionsMsg{raw: data}
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	var extensions cryptobyte.String
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		!s.ReadUint16LengthPrefixed(&amp;extensions) || !s.Empty() {
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		return false
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>	}
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	for !extensions.Empty() {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		var extension uint16
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		var extData cryptobyte.String
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>		if !extensions.ReadUint16(&amp;extension) ||
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			return false
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>		}
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		switch extension {
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		case extensionALPN:
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			var protoList cryptobyte.String
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;protoList) || protoList.Empty() {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>				return false
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>			}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>			var proto cryptobyte.String
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			if !protoList.ReadUint8LengthPrefixed(&amp;proto) ||
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>				proto.Empty() || !protoList.Empty() {
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>				return false
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>			}
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			m.alpnProtocol = string(proto)
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>		case extensionQUICTransportParameters:
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			m.quicTransportParameters = make([]byte, len(extData))
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			if !extData.CopyBytes(m.quicTransportParameters) {
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>				return false
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>		case extensionEarlyData:
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>			<span class="comment">// RFC 8446, Section 4.2.10</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			m.earlyData = true
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		default:
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>			<span class="comment">// Ignore unknown extensions.</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			continue
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		if !extData.Empty() {
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>			return false
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		}
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	}
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	return true
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>type endOfEarlyDataMsg struct{}
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>func (m *endOfEarlyDataMsg) marshal() ([]byte, error) {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	x := make([]byte, 4)
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	x[0] = typeEndOfEarlyData
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	return x, nil
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>}
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>func (m *endOfEarlyDataMsg) unmarshal(data []byte) bool {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	return len(data) == 4
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>type keyUpdateMsg struct {
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	raw             []byte
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	updateRequested bool
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>}
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>func (m *keyUpdateMsg) marshal() ([]byte, error) {
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>	b.AddUint8(typeKeyUpdate)
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>		if m.updateRequested {
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>			b.AddUint8(1)
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		} else {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>			b.AddUint8(0)
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>		}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	})
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	var err error
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>func (m *keyUpdateMsg) unmarshal(data []byte) bool {
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	m.raw = data
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	var updateRequested uint8
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>		!s.ReadUint8(&amp;updateRequested) || !s.Empty() {
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>		return false
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	}
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	switch updateRequested {
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	case 0:
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>		m.updateRequested = false
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	case 1:
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		m.updateRequested = true
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	default:
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		return false
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	}
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	return true
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>}
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>type newSessionTicketMsgTLS13 struct {
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	raw          []byte
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	lifetime     uint32
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	ageAdd       uint32
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	nonce        []byte
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	label        []byte
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	maxEarlyData uint32
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>}
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>func (m *newSessionTicketMsgTLS13) marshal() ([]byte, error) {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	}
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	b.AddUint8(typeNewSessionTicket)
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		b.AddUint32(m.lifetime)
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>		b.AddUint32(m.ageAdd)
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		b.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			b.AddBytes(m.nonce)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		})
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>			b.AddBytes(m.label)
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		})
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>			if m.maxEarlyData &gt; 0 {
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>				b.AddUint16(extensionEarlyData)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>					b.AddUint32(m.maxEarlyData)
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>				})
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>			}
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>		})
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	})
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	var err error
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>}
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>func (m *newSessionTicketMsgTLS13) unmarshal(data []byte) bool {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	*m = newSessionTicketMsgTLS13{raw: data}
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	var extensions cryptobyte.String
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		!s.ReadUint32(&amp;m.lifetime) ||
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>		!s.ReadUint32(&amp;m.ageAdd) ||
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		!readUint8LengthPrefixed(&amp;s, &amp;m.nonce) ||
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>		!readUint16LengthPrefixed(&amp;s, &amp;m.label) ||
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>		!s.ReadUint16LengthPrefixed(&amp;extensions) ||
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		!s.Empty() {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		return false
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	for !extensions.Empty() {
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		var extension uint16
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>		var extData cryptobyte.String
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		if !extensions.ReadUint16(&amp;extension) ||
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>			!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>			return false
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		}
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>		switch extension {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		case extensionEarlyData:
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>			if !extData.ReadUint32(&amp;m.maxEarlyData) {
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>				return false
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			}
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>		default:
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			<span class="comment">// Ignore unknown extensions.</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			continue
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>		}
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>		if !extData.Empty() {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>			return false
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>		}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	return true
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>}
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>type certificateRequestMsgTLS13 struct {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	raw                              []byte
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	ocspStapling                     bool
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>	scts                             bool
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	supportedSignatureAlgorithms     []SignatureScheme
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	supportedSignatureAlgorithmsCert []SignatureScheme
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	certificateAuthorities           [][]byte
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>}
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>func (m *certificateRequestMsgTLS13) marshal() ([]byte, error) {
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	}
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	b.AddUint8(typeCertificateRequest)
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>		<span class="comment">// certificate_request_context (SHALL be zero length unless used for</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>		<span class="comment">// post-handshake authentication)</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>		b.AddUint8(0)
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>			if m.ocspStapling {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>				b.AddUint16(extensionStatusRequest)
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>				b.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>			}
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>			if m.scts {
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>				<span class="comment">// RFC 8446, Section 4.4.2.1 makes no mention of</span>
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>				<span class="comment">// signed_certificate_timestamp in CertificateRequest, but</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>				<span class="comment">// &#34;Extensions in the Certificate message from the client MUST</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>				<span class="comment">// correspond to extensions in the CertificateRequest message</span>
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>				<span class="comment">// from the server.&#34; and it appears in the table in Section 4.2.</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>				b.AddUint16(extensionSCT)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>				b.AddUint16(0) <span class="comment">// empty extension_data</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>			}
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>			if len(m.supportedSignatureAlgorithms) &gt; 0 {
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>				b.AddUint16(extensionSignatureAlgorithms)
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>						for _, sigAlgo := range m.supportedSignatureAlgorithms {
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>							b.AddUint16(uint16(sigAlgo))
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>						}
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>					})
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>				})
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>			}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>			if len(m.supportedSignatureAlgorithmsCert) &gt; 0 {
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>				b.AddUint16(extensionSignatureAlgorithmsCert)
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>						for _, sigAlgo := range m.supportedSignatureAlgorithmsCert {
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>							b.AddUint16(uint16(sigAlgo))
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>						}
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>					})
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>				})
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>			}
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>			if len(m.certificateAuthorities) &gt; 0 {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>				b.AddUint16(extensionCertificateAuthorities)
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>				b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>						for _, ca := range m.certificateAuthorities {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>							b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>								b.AddBytes(ca)
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>							})
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>						}
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>					})
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>				})
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>			}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>		})
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>	})
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	var err error
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>}
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>func (m *certificateRequestMsgTLS13) unmarshal(data []byte) bool {
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>	*m = certificateRequestMsgTLS13{raw: data}
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>	var context, extensions cryptobyte.String
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>		!s.ReadUint8LengthPrefixed(&amp;context) || !context.Empty() ||
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>		!s.ReadUint16LengthPrefixed(&amp;extensions) ||
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>		!s.Empty() {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		return false
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>	}
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>	for !extensions.Empty() {
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		var extension uint16
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>		var extData cryptobyte.String
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		if !extensions.ReadUint16(&amp;extension) ||
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>			!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			return false
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		}
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>		switch extension {
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>		case extensionStatusRequest:
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>			m.ocspStapling = true
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>		case extensionSCT:
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>			m.scts = true
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>		case extensionSignatureAlgorithms:
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>			var sigAndAlgs cryptobyte.String
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;sigAndAlgs) || sigAndAlgs.Empty() {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>				return false
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>			}
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>			for !sigAndAlgs.Empty() {
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>				var sigAndAlg uint16
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				if !sigAndAlgs.ReadUint16(&amp;sigAndAlg) {
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>					return false
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>				}
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>				m.supportedSignatureAlgorithms = append(
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>					m.supportedSignatureAlgorithms, SignatureScheme(sigAndAlg))
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>			}
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>		case extensionSignatureAlgorithmsCert:
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>			var sigAndAlgs cryptobyte.String
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;sigAndAlgs) || sigAndAlgs.Empty() {
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>				return false
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>			}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>			for !sigAndAlgs.Empty() {
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>				var sigAndAlg uint16
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>				if !sigAndAlgs.ReadUint16(&amp;sigAndAlg) {
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>					return false
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>				}
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>				m.supportedSignatureAlgorithmsCert = append(
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>					m.supportedSignatureAlgorithmsCert, SignatureScheme(sigAndAlg))
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>			}
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>		case extensionCertificateAuthorities:
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>			var auths cryptobyte.String
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>			if !extData.ReadUint16LengthPrefixed(&amp;auths) || auths.Empty() {
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>				return false
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			}
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>			for !auths.Empty() {
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>				var ca []byte
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>				if !readUint16LengthPrefixed(&amp;auths, &amp;ca) || len(ca) == 0 {
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>					return false
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>				}
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>				m.certificateAuthorities = append(m.certificateAuthorities, ca)
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>			}
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>		default:
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>			<span class="comment">// Ignore unknown extensions.</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>			continue
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>		}
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>		if !extData.Empty() {
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>			return false
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>		}
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	}
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	return true
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>}
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>type certificateMsg struct {
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	raw          []byte
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	certificates [][]byte
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>}
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>func (m *certificateMsg) marshal() ([]byte, error) {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	}
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	var i int
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	for _, slice := range m.certificates {
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>		i += len(slice)
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>	length := 3 + 3*len(m.certificates) + i
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>	x := make([]byte, 4+length)
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	x[0] = typeCertificate
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>	x[1] = uint8(length &gt;&gt; 16)
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	x[2] = uint8(length &gt;&gt; 8)
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>	x[3] = uint8(length)
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	certificateOctets := length - 3
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	x[4] = uint8(certificateOctets &gt;&gt; 16)
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	x[5] = uint8(certificateOctets &gt;&gt; 8)
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>	x[6] = uint8(certificateOctets)
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>	y := x[7:]
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	for _, slice := range m.certificates {
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>		y[0] = uint8(len(slice) &gt;&gt; 16)
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		y[1] = uint8(len(slice) &gt;&gt; 8)
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>		y[2] = uint8(len(slice))
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		copy(y[3:], slice)
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		y = y[3+len(slice):]
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>	}
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>	m.raw = x
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>	return m.raw, nil
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>}
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>func (m *certificateMsg) unmarshal(data []byte) bool {
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>	if len(data) &lt; 7 {
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		return false
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>	}
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	m.raw = data
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	certsLen := uint32(data[4])&lt;&lt;16 | uint32(data[5])&lt;&lt;8 | uint32(data[6])
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>	if uint32(len(data)) != certsLen+7 {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>		return false
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	}
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	numCerts := 0
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>	d := data[7:]
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>	for certsLen &gt; 0 {
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>		if len(d) &lt; 4 {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			return false
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		certLen := uint32(d[0])&lt;&lt;16 | uint32(d[1])&lt;&lt;8 | uint32(d[2])
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>		if uint32(len(d)) &lt; 3+certLen {
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>			return false
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>		}
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>		d = d[3+certLen:]
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>		certsLen -= 3 + certLen
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		numCerts++
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>	}
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>	m.certificates = make([][]byte, numCerts)
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	d = data[7:]
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>	for i := 0; i &lt; numCerts; i++ {
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		certLen := uint32(d[0])&lt;&lt;16 | uint32(d[1])&lt;&lt;8 | uint32(d[2])
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		m.certificates[i] = d[3 : 3+certLen]
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>		d = d[3+certLen:]
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	}
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	return true
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>}
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>type certificateMsgTLS13 struct {
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	raw          []byte
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>	certificate  Certificate
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>	ocspStapling bool
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	scts         bool
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>}
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>func (m *certificateMsgTLS13) marshal() ([]byte, error) {
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>	}
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>	b.AddUint8(typeCertificate)
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		b.AddUint8(0) <span class="comment">// certificate_request_context</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>		certificate := m.certificate
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>		if !m.ocspStapling {
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>			certificate.OCSPStaple = nil
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>		}
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>		if !m.scts {
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>			certificate.SignedCertificateTimestamps = nil
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		}
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		marshalCertificate(b, certificate)
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>	})
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>	var err error
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>}
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>func marshalCertificate(b *cryptobyte.Builder, certificate Certificate) {
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>		for i, cert := range certificate.Certificate {
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>			b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>				b.AddBytes(cert)
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>			})
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>			b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>				if i &gt; 0 {
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>					<span class="comment">// This library only supports OCSP and SCT for leaf certificates.</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>					return
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>				}
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>				if certificate.OCSPStaple != nil {
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>					b.AddUint16(extensionStatusRequest)
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>						b.AddUint8(statusTypeOCSP)
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>						b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>							b.AddBytes(certificate.OCSPStaple)
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>						})
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>					})
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>				}
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>				if certificate.SignedCertificateTimestamps != nil {
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>					b.AddUint16(extensionSCT)
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>					b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>						b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>							for _, sct := range certificate.SignedCertificateTimestamps {
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>								b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>									b.AddBytes(sct)
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>								})
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>							}
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>						})
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>					})
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>				}
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>			})
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>		}
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>	})
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>}
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>func (m *certificateMsgTLS13) unmarshal(data []byte) bool {
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>	*m = certificateMsgTLS13{raw: data}
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>	var context cryptobyte.String
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>		!s.ReadUint8LengthPrefixed(&amp;context) || !context.Empty() ||
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>		!unmarshalCertificate(&amp;s, &amp;m.certificate) ||
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>		!s.Empty() {
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>		return false
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>	}
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>	m.scts = m.certificate.SignedCertificateTimestamps != nil
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>	m.ocspStapling = m.certificate.OCSPStaple != nil
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>	return true
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>}
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>func unmarshalCertificate(s *cryptobyte.String, certificate *Certificate) bool {
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	var certList cryptobyte.String
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>	if !s.ReadUint24LengthPrefixed(&amp;certList) {
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>		return false
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>	}
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>	for !certList.Empty() {
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>		var cert []byte
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>		var extensions cryptobyte.String
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>		if !readUint24LengthPrefixed(&amp;certList, &amp;cert) ||
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>			!certList.ReadUint16LengthPrefixed(&amp;extensions) {
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>			return false
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>		}
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>		certificate.Certificate = append(certificate.Certificate, cert)
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>		for !extensions.Empty() {
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>			var extension uint16
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>			var extData cryptobyte.String
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>			if !extensions.ReadUint16(&amp;extension) ||
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>				!extensions.ReadUint16LengthPrefixed(&amp;extData) {
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>				return false
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>			}
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>			if len(certificate.Certificate) &gt; 1 {
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>				<span class="comment">// This library only supports OCSP and SCT for leaf certificates.</span>
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>				continue
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>			}
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>			switch extension {
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>			case extensionStatusRequest:
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>				var statusType uint8
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>				if !extData.ReadUint8(&amp;statusType) || statusType != statusTypeOCSP ||
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>					!readUint24LengthPrefixed(&amp;extData, &amp;certificate.OCSPStaple) ||
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>					len(certificate.OCSPStaple) == 0 {
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>					return false
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>				}
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>			case extensionSCT:
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>				var sctList cryptobyte.String
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>				if !extData.ReadUint16LengthPrefixed(&amp;sctList) || sctList.Empty() {
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>					return false
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>				}
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>				for !sctList.Empty() {
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>					var sct []byte
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>					if !readUint16LengthPrefixed(&amp;sctList, &amp;sct) ||
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>						len(sct) == 0 {
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>						return false
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>					}
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>					certificate.SignedCertificateTimestamps = append(
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>						certificate.SignedCertificateTimestamps, sct)
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>				}
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>			default:
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>				<span class="comment">// Ignore unknown extensions.</span>
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>				continue
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>			}
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>			if !extData.Empty() {
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>				return false
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>			}
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>		}
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>	}
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>	return true
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>}
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>type serverKeyExchangeMsg struct {
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>	raw []byte
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>	key []byte
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>}
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>func (m *serverKeyExchangeMsg) marshal() ([]byte, error) {
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>	}
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>	length := len(m.key)
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>	x := make([]byte, length+4)
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>	x[0] = typeServerKeyExchange
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>	x[1] = uint8(length &gt;&gt; 16)
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	x[2] = uint8(length &gt;&gt; 8)
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>	x[3] = uint8(length)
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>	copy(x[4:], m.key)
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>	m.raw = x
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>	return x, nil
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>}
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>func (m *serverKeyExchangeMsg) unmarshal(data []byte) bool {
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>	m.raw = data
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>	if len(data) &lt; 4 {
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>		return false
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>	}
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>	m.key = data[4:]
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>	return true
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>}
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>type certificateStatusMsg struct {
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>	raw      []byte
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>	response []byte
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>}
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>func (m *certificateStatusMsg) marshal() ([]byte, error) {
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>	}
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>	b.AddUint8(typeCertificateStatus)
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>		b.AddUint8(statusTypeOCSP)
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>		b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>			b.AddBytes(m.response)
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>		})
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>	})
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>	var err error
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>}
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>func (m *certificateStatusMsg) unmarshal(data []byte) bool {
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>	m.raw = data
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>	var statusType uint8
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>	if !s.Skip(4) || <span class="comment">// message type and uint24 length field</span>
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>		!s.ReadUint8(&amp;statusType) || statusType != statusTypeOCSP ||
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>		!readUint24LengthPrefixed(&amp;s, &amp;m.response) ||
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>		len(m.response) == 0 || !s.Empty() {
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>		return false
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>	}
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>	return true
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>}
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>type serverHelloDoneMsg struct{}
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>func (m *serverHelloDoneMsg) marshal() ([]byte, error) {
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>	x := make([]byte, 4)
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>	x[0] = typeServerHelloDone
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>	return x, nil
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>}
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>func (m *serverHelloDoneMsg) unmarshal(data []byte) bool {
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>	return len(data) == 4
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>}
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>type clientKeyExchangeMsg struct {
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>	raw        []byte
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>	ciphertext []byte
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>}
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>func (m *clientKeyExchangeMsg) marshal() ([]byte, error) {
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>	}
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>	length := len(m.ciphertext)
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>	x := make([]byte, length+4)
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>	x[0] = typeClientKeyExchange
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>	x[1] = uint8(length &gt;&gt; 16)
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>	x[2] = uint8(length &gt;&gt; 8)
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>	x[3] = uint8(length)
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>	copy(x[4:], m.ciphertext)
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>	m.raw = x
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>	return x, nil
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>func (m *clientKeyExchangeMsg) unmarshal(data []byte) bool {
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>	m.raw = data
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>	if len(data) &lt; 4 {
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>		return false
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>	}
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>	l := int(data[1])&lt;&lt;16 | int(data[2])&lt;&lt;8 | int(data[3])
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>	if l != len(data)-4 {
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>		return false
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>	}
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>	m.ciphertext = data[4:]
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>	return true
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>}
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>type finishedMsg struct {
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>	raw        []byte
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>	verifyData []byte
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>}
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>func (m *finishedMsg) marshal() ([]byte, error) {
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>	}
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>	b.AddUint8(typeFinished)
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>		b.AddBytes(m.verifyData)
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>	})
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>	var err error
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>}
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>func (m *finishedMsg) unmarshal(data []byte) bool {
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>	m.raw = data
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>	return s.Skip(1) &amp;&amp;
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>		readUint24LengthPrefixed(&amp;s, &amp;m.verifyData) &amp;&amp;
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>		s.Empty()
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>}
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>type certificateRequestMsg struct {
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>	raw []byte
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>	<span class="comment">// hasSignatureAlgorithm indicates whether this message includes a list of</span>
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>	<span class="comment">// supported signature algorithms. This change was introduced with TLS 1.2.</span>
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>	hasSignatureAlgorithm bool
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>	certificateTypes             []byte
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>	supportedSignatureAlgorithms []SignatureScheme
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>	certificateAuthorities       [][]byte
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>}
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>func (m *certificateRequestMsg) marshal() ([]byte, error) {
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>	}
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>	<span class="comment">// See RFC 4346, Section 7.4.4.</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>	length := 1 + len(m.certificateTypes) + 2
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>	casLength := 0
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>	for _, ca := range m.certificateAuthorities {
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>		casLength += 2 + len(ca)
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>	}
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>	length += casLength
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>	if m.hasSignatureAlgorithm {
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>		length += 2 + 2*len(m.supportedSignatureAlgorithms)
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>	}
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>	x := make([]byte, 4+length)
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>	x[0] = typeCertificateRequest
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>	x[1] = uint8(length &gt;&gt; 16)
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>	x[2] = uint8(length &gt;&gt; 8)
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>	x[3] = uint8(length)
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	x[4] = uint8(len(m.certificateTypes))
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	copy(x[5:], m.certificateTypes)
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>	y := x[5+len(m.certificateTypes):]
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>	if m.hasSignatureAlgorithm {
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>		n := len(m.supportedSignatureAlgorithms) * 2
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>		y[0] = uint8(n &gt;&gt; 8)
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>		y[1] = uint8(n)
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>		y = y[2:]
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>		for _, sigAlgo := range m.supportedSignatureAlgorithms {
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>			y[0] = uint8(sigAlgo &gt;&gt; 8)
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>			y[1] = uint8(sigAlgo)
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>			y = y[2:]
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>		}
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>	}
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>	y[0] = uint8(casLength &gt;&gt; 8)
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>	y[1] = uint8(casLength)
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>	y = y[2:]
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>	for _, ca := range m.certificateAuthorities {
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>		y[0] = uint8(len(ca) &gt;&gt; 8)
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>		y[1] = uint8(len(ca))
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>		y = y[2:]
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>		copy(y, ca)
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>		y = y[len(ca):]
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>	}
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>	m.raw = x
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>	return m.raw, nil
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>}
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>func (m *certificateRequestMsg) unmarshal(data []byte) bool {
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>	m.raw = data
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>	if len(data) &lt; 5 {
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>		return false
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>	}
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>	length := uint32(data[1])&lt;&lt;16 | uint32(data[2])&lt;&lt;8 | uint32(data[3])
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>	if uint32(len(data))-4 != length {
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>		return false
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>	}
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>	numCertTypes := int(data[4])
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>	data = data[5:]
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>	if numCertTypes == 0 || len(data) &lt;= numCertTypes {
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>		return false
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>	}
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>	m.certificateTypes = make([]byte, numCertTypes)
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>	if copy(m.certificateTypes, data) != numCertTypes {
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>		return false
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>	}
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>
<span id="L1734" class="ln">  1734&nbsp;&nbsp;</span>	data = data[numCertTypes:]
<span id="L1735" class="ln">  1735&nbsp;&nbsp;</span>
<span id="L1736" class="ln">  1736&nbsp;&nbsp;</span>	if m.hasSignatureAlgorithm {
<span id="L1737" class="ln">  1737&nbsp;&nbsp;</span>		if len(data) &lt; 2 {
<span id="L1738" class="ln">  1738&nbsp;&nbsp;</span>			return false
<span id="L1739" class="ln">  1739&nbsp;&nbsp;</span>		}
<span id="L1740" class="ln">  1740&nbsp;&nbsp;</span>		sigAndHashLen := uint16(data[0])&lt;&lt;8 | uint16(data[1])
<span id="L1741" class="ln">  1741&nbsp;&nbsp;</span>		data = data[2:]
<span id="L1742" class="ln">  1742&nbsp;&nbsp;</span>		if sigAndHashLen&amp;1 != 0 {
<span id="L1743" class="ln">  1743&nbsp;&nbsp;</span>			return false
<span id="L1744" class="ln">  1744&nbsp;&nbsp;</span>		}
<span id="L1745" class="ln">  1745&nbsp;&nbsp;</span>		if len(data) &lt; int(sigAndHashLen) {
<span id="L1746" class="ln">  1746&nbsp;&nbsp;</span>			return false
<span id="L1747" class="ln">  1747&nbsp;&nbsp;</span>		}
<span id="L1748" class="ln">  1748&nbsp;&nbsp;</span>		numSigAlgos := sigAndHashLen / 2
<span id="L1749" class="ln">  1749&nbsp;&nbsp;</span>		m.supportedSignatureAlgorithms = make([]SignatureScheme, numSigAlgos)
<span id="L1750" class="ln">  1750&nbsp;&nbsp;</span>		for i := range m.supportedSignatureAlgorithms {
<span id="L1751" class="ln">  1751&nbsp;&nbsp;</span>			m.supportedSignatureAlgorithms[i] = SignatureScheme(data[0])&lt;&lt;8 | SignatureScheme(data[1])
<span id="L1752" class="ln">  1752&nbsp;&nbsp;</span>			data = data[2:]
<span id="L1753" class="ln">  1753&nbsp;&nbsp;</span>		}
<span id="L1754" class="ln">  1754&nbsp;&nbsp;</span>	}
<span id="L1755" class="ln">  1755&nbsp;&nbsp;</span>
<span id="L1756" class="ln">  1756&nbsp;&nbsp;</span>	if len(data) &lt; 2 {
<span id="L1757" class="ln">  1757&nbsp;&nbsp;</span>		return false
<span id="L1758" class="ln">  1758&nbsp;&nbsp;</span>	}
<span id="L1759" class="ln">  1759&nbsp;&nbsp;</span>	casLength := uint16(data[0])&lt;&lt;8 | uint16(data[1])
<span id="L1760" class="ln">  1760&nbsp;&nbsp;</span>	data = data[2:]
<span id="L1761" class="ln">  1761&nbsp;&nbsp;</span>	if len(data) &lt; int(casLength) {
<span id="L1762" class="ln">  1762&nbsp;&nbsp;</span>		return false
<span id="L1763" class="ln">  1763&nbsp;&nbsp;</span>	}
<span id="L1764" class="ln">  1764&nbsp;&nbsp;</span>	cas := make([]byte, casLength)
<span id="L1765" class="ln">  1765&nbsp;&nbsp;</span>	copy(cas, data)
<span id="L1766" class="ln">  1766&nbsp;&nbsp;</span>	data = data[casLength:]
<span id="L1767" class="ln">  1767&nbsp;&nbsp;</span>
<span id="L1768" class="ln">  1768&nbsp;&nbsp;</span>	m.certificateAuthorities = nil
<span id="L1769" class="ln">  1769&nbsp;&nbsp;</span>	for len(cas) &gt; 0 {
<span id="L1770" class="ln">  1770&nbsp;&nbsp;</span>		if len(cas) &lt; 2 {
<span id="L1771" class="ln">  1771&nbsp;&nbsp;</span>			return false
<span id="L1772" class="ln">  1772&nbsp;&nbsp;</span>		}
<span id="L1773" class="ln">  1773&nbsp;&nbsp;</span>		caLen := uint16(cas[0])&lt;&lt;8 | uint16(cas[1])
<span id="L1774" class="ln">  1774&nbsp;&nbsp;</span>		cas = cas[2:]
<span id="L1775" class="ln">  1775&nbsp;&nbsp;</span>
<span id="L1776" class="ln">  1776&nbsp;&nbsp;</span>		if len(cas) &lt; int(caLen) {
<span id="L1777" class="ln">  1777&nbsp;&nbsp;</span>			return false
<span id="L1778" class="ln">  1778&nbsp;&nbsp;</span>		}
<span id="L1779" class="ln">  1779&nbsp;&nbsp;</span>
<span id="L1780" class="ln">  1780&nbsp;&nbsp;</span>		m.certificateAuthorities = append(m.certificateAuthorities, cas[:caLen])
<span id="L1781" class="ln">  1781&nbsp;&nbsp;</span>		cas = cas[caLen:]
<span id="L1782" class="ln">  1782&nbsp;&nbsp;</span>	}
<span id="L1783" class="ln">  1783&nbsp;&nbsp;</span>
<span id="L1784" class="ln">  1784&nbsp;&nbsp;</span>	return len(data) == 0
<span id="L1785" class="ln">  1785&nbsp;&nbsp;</span>}
<span id="L1786" class="ln">  1786&nbsp;&nbsp;</span>
<span id="L1787" class="ln">  1787&nbsp;&nbsp;</span>type certificateVerifyMsg struct {
<span id="L1788" class="ln">  1788&nbsp;&nbsp;</span>	raw                   []byte
<span id="L1789" class="ln">  1789&nbsp;&nbsp;</span>	hasSignatureAlgorithm bool <span class="comment">// format change introduced in TLS 1.2</span>
<span id="L1790" class="ln">  1790&nbsp;&nbsp;</span>	signatureAlgorithm    SignatureScheme
<span id="L1791" class="ln">  1791&nbsp;&nbsp;</span>	signature             []byte
<span id="L1792" class="ln">  1792&nbsp;&nbsp;</span>}
<span id="L1793" class="ln">  1793&nbsp;&nbsp;</span>
<span id="L1794" class="ln">  1794&nbsp;&nbsp;</span>func (m *certificateVerifyMsg) marshal() ([]byte, error) {
<span id="L1795" class="ln">  1795&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1796" class="ln">  1796&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1797" class="ln">  1797&nbsp;&nbsp;</span>	}
<span id="L1798" class="ln">  1798&nbsp;&nbsp;</span>
<span id="L1799" class="ln">  1799&nbsp;&nbsp;</span>	var b cryptobyte.Builder
<span id="L1800" class="ln">  1800&nbsp;&nbsp;</span>	b.AddUint8(typeCertificateVerify)
<span id="L1801" class="ln">  1801&nbsp;&nbsp;</span>	b.AddUint24LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1802" class="ln">  1802&nbsp;&nbsp;</span>		if m.hasSignatureAlgorithm {
<span id="L1803" class="ln">  1803&nbsp;&nbsp;</span>			b.AddUint16(uint16(m.signatureAlgorithm))
<span id="L1804" class="ln">  1804&nbsp;&nbsp;</span>		}
<span id="L1805" class="ln">  1805&nbsp;&nbsp;</span>		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L1806" class="ln">  1806&nbsp;&nbsp;</span>			b.AddBytes(m.signature)
<span id="L1807" class="ln">  1807&nbsp;&nbsp;</span>		})
<span id="L1808" class="ln">  1808&nbsp;&nbsp;</span>	})
<span id="L1809" class="ln">  1809&nbsp;&nbsp;</span>
<span id="L1810" class="ln">  1810&nbsp;&nbsp;</span>	var err error
<span id="L1811" class="ln">  1811&nbsp;&nbsp;</span>	m.raw, err = b.Bytes()
<span id="L1812" class="ln">  1812&nbsp;&nbsp;</span>	return m.raw, err
<span id="L1813" class="ln">  1813&nbsp;&nbsp;</span>}
<span id="L1814" class="ln">  1814&nbsp;&nbsp;</span>
<span id="L1815" class="ln">  1815&nbsp;&nbsp;</span>func (m *certificateVerifyMsg) unmarshal(data []byte) bool {
<span id="L1816" class="ln">  1816&nbsp;&nbsp;</span>	m.raw = data
<span id="L1817" class="ln">  1817&nbsp;&nbsp;</span>	s := cryptobyte.String(data)
<span id="L1818" class="ln">  1818&nbsp;&nbsp;</span>
<span id="L1819" class="ln">  1819&nbsp;&nbsp;</span>	if !s.Skip(4) { <span class="comment">// message type and uint24 length field</span>
<span id="L1820" class="ln">  1820&nbsp;&nbsp;</span>		return false
<span id="L1821" class="ln">  1821&nbsp;&nbsp;</span>	}
<span id="L1822" class="ln">  1822&nbsp;&nbsp;</span>	if m.hasSignatureAlgorithm {
<span id="L1823" class="ln">  1823&nbsp;&nbsp;</span>		if !s.ReadUint16((*uint16)(&amp;m.signatureAlgorithm)) {
<span id="L1824" class="ln">  1824&nbsp;&nbsp;</span>			return false
<span id="L1825" class="ln">  1825&nbsp;&nbsp;</span>		}
<span id="L1826" class="ln">  1826&nbsp;&nbsp;</span>	}
<span id="L1827" class="ln">  1827&nbsp;&nbsp;</span>	return readUint16LengthPrefixed(&amp;s, &amp;m.signature) &amp;&amp; s.Empty()
<span id="L1828" class="ln">  1828&nbsp;&nbsp;</span>}
<span id="L1829" class="ln">  1829&nbsp;&nbsp;</span>
<span id="L1830" class="ln">  1830&nbsp;&nbsp;</span>type newSessionTicketMsg struct {
<span id="L1831" class="ln">  1831&nbsp;&nbsp;</span>	raw    []byte
<span id="L1832" class="ln">  1832&nbsp;&nbsp;</span>	ticket []byte
<span id="L1833" class="ln">  1833&nbsp;&nbsp;</span>}
<span id="L1834" class="ln">  1834&nbsp;&nbsp;</span>
<span id="L1835" class="ln">  1835&nbsp;&nbsp;</span>func (m *newSessionTicketMsg) marshal() ([]byte, error) {
<span id="L1836" class="ln">  1836&nbsp;&nbsp;</span>	if m.raw != nil {
<span id="L1837" class="ln">  1837&nbsp;&nbsp;</span>		return m.raw, nil
<span id="L1838" class="ln">  1838&nbsp;&nbsp;</span>	}
<span id="L1839" class="ln">  1839&nbsp;&nbsp;</span>
<span id="L1840" class="ln">  1840&nbsp;&nbsp;</span>	<span class="comment">// See RFC 5077, Section 3.3.</span>
<span id="L1841" class="ln">  1841&nbsp;&nbsp;</span>	ticketLen := len(m.ticket)
<span id="L1842" class="ln">  1842&nbsp;&nbsp;</span>	length := 2 + 4 + ticketLen
<span id="L1843" class="ln">  1843&nbsp;&nbsp;</span>	x := make([]byte, 4+length)
<span id="L1844" class="ln">  1844&nbsp;&nbsp;</span>	x[0] = typeNewSessionTicket
<span id="L1845" class="ln">  1845&nbsp;&nbsp;</span>	x[1] = uint8(length &gt;&gt; 16)
<span id="L1846" class="ln">  1846&nbsp;&nbsp;</span>	x[2] = uint8(length &gt;&gt; 8)
<span id="L1847" class="ln">  1847&nbsp;&nbsp;</span>	x[3] = uint8(length)
<span id="L1848" class="ln">  1848&nbsp;&nbsp;</span>	x[8] = uint8(ticketLen &gt;&gt; 8)
<span id="L1849" class="ln">  1849&nbsp;&nbsp;</span>	x[9] = uint8(ticketLen)
<span id="L1850" class="ln">  1850&nbsp;&nbsp;</span>	copy(x[10:], m.ticket)
<span id="L1851" class="ln">  1851&nbsp;&nbsp;</span>
<span id="L1852" class="ln">  1852&nbsp;&nbsp;</span>	m.raw = x
<span id="L1853" class="ln">  1853&nbsp;&nbsp;</span>
<span id="L1854" class="ln">  1854&nbsp;&nbsp;</span>	return m.raw, nil
<span id="L1855" class="ln">  1855&nbsp;&nbsp;</span>}
<span id="L1856" class="ln">  1856&nbsp;&nbsp;</span>
<span id="L1857" class="ln">  1857&nbsp;&nbsp;</span>func (m *newSessionTicketMsg) unmarshal(data []byte) bool {
<span id="L1858" class="ln">  1858&nbsp;&nbsp;</span>	m.raw = data
<span id="L1859" class="ln">  1859&nbsp;&nbsp;</span>
<span id="L1860" class="ln">  1860&nbsp;&nbsp;</span>	if len(data) &lt; 10 {
<span id="L1861" class="ln">  1861&nbsp;&nbsp;</span>		return false
<span id="L1862" class="ln">  1862&nbsp;&nbsp;</span>	}
<span id="L1863" class="ln">  1863&nbsp;&nbsp;</span>
<span id="L1864" class="ln">  1864&nbsp;&nbsp;</span>	length := uint32(data[1])&lt;&lt;16 | uint32(data[2])&lt;&lt;8 | uint32(data[3])
<span id="L1865" class="ln">  1865&nbsp;&nbsp;</span>	if uint32(len(data))-4 != length {
<span id="L1866" class="ln">  1866&nbsp;&nbsp;</span>		return false
<span id="L1867" class="ln">  1867&nbsp;&nbsp;</span>	}
<span id="L1868" class="ln">  1868&nbsp;&nbsp;</span>
<span id="L1869" class="ln">  1869&nbsp;&nbsp;</span>	ticketLen := int(data[8])&lt;&lt;8 + int(data[9])
<span id="L1870" class="ln">  1870&nbsp;&nbsp;</span>	if len(data)-10 != ticketLen {
<span id="L1871" class="ln">  1871&nbsp;&nbsp;</span>		return false
<span id="L1872" class="ln">  1872&nbsp;&nbsp;</span>	}
<span id="L1873" class="ln">  1873&nbsp;&nbsp;</span>
<span id="L1874" class="ln">  1874&nbsp;&nbsp;</span>	m.ticket = data[10:]
<span id="L1875" class="ln">  1875&nbsp;&nbsp;</span>
<span id="L1876" class="ln">  1876&nbsp;&nbsp;</span>	return true
<span id="L1877" class="ln">  1877&nbsp;&nbsp;</span>}
<span id="L1878" class="ln">  1878&nbsp;&nbsp;</span>
<span id="L1879" class="ln">  1879&nbsp;&nbsp;</span>type helloRequestMsg struct {
<span id="L1880" class="ln">  1880&nbsp;&nbsp;</span>}
<span id="L1881" class="ln">  1881&nbsp;&nbsp;</span>
<span id="L1882" class="ln">  1882&nbsp;&nbsp;</span>func (*helloRequestMsg) marshal() ([]byte, error) {
<span id="L1883" class="ln">  1883&nbsp;&nbsp;</span>	return []byte{typeHelloRequest, 0, 0, 0}, nil
<span id="L1884" class="ln">  1884&nbsp;&nbsp;</span>}
<span id="L1885" class="ln">  1885&nbsp;&nbsp;</span>
<span id="L1886" class="ln">  1886&nbsp;&nbsp;</span>func (*helloRequestMsg) unmarshal(data []byte) bool {
<span id="L1887" class="ln">  1887&nbsp;&nbsp;</span>	return len(data) == 4
<span id="L1888" class="ln">  1888&nbsp;&nbsp;</span>}
<span id="L1889" class="ln">  1889&nbsp;&nbsp;</span>
<span id="L1890" class="ln">  1890&nbsp;&nbsp;</span>type transcriptHash interface {
<span id="L1891" class="ln">  1891&nbsp;&nbsp;</span>	Write([]byte) (int, error)
<span id="L1892" class="ln">  1892&nbsp;&nbsp;</span>}
<span id="L1893" class="ln">  1893&nbsp;&nbsp;</span>
<span id="L1894" class="ln">  1894&nbsp;&nbsp;</span><span class="comment">// transcriptMsg is a helper used to marshal and hash messages which typically</span>
<span id="L1895" class="ln">  1895&nbsp;&nbsp;</span><span class="comment">// are not written to the wire, and as such aren&#39;t hashed during Conn.writeRecord.</span>
<span id="L1896" class="ln">  1896&nbsp;&nbsp;</span>func transcriptMsg(msg handshakeMessage, h transcriptHash) error {
<span id="L1897" class="ln">  1897&nbsp;&nbsp;</span>	data, err := msg.marshal()
<span id="L1898" class="ln">  1898&nbsp;&nbsp;</span>	if err != nil {
<span id="L1899" class="ln">  1899&nbsp;&nbsp;</span>		return err
<span id="L1900" class="ln">  1900&nbsp;&nbsp;</span>	}
<span id="L1901" class="ln">  1901&nbsp;&nbsp;</span>	h.Write(data)
<span id="L1902" class="ln">  1902&nbsp;&nbsp;</span>	return nil
<span id="L1903" class="ln">  1903&nbsp;&nbsp;</span>}
<span id="L1904" class="ln">  1904&nbsp;&nbsp;</span>
</pre><p><a href="handshake_messages.go?m=text">View as plain text</a></p>

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
