<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/auth.go - Go Documentation Server</title>

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
<a href="auth.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">auth.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/tls">crypto/tls</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package tls
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;crypto&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto/ecdsa&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/ed25519&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/elliptic&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/rsa&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// verifyHandshakeSignature verifies a signature against pre-hashed</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// (if required) handshake contents.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>func verifyHandshakeSignature(sigType uint8, pubkey crypto.PublicKey, hashFunc crypto.Hash, signed, sig []byte) error {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	switch sigType {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	case signatureECDSA:
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		pubKey, ok := pubkey.(*ecdsa.PublicKey)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		if !ok {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;expected an ECDSA public key, got %T&#34;, pubkey)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		if !ecdsa.VerifyASN1(pubKey, signed, sig) {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			return errors.New(&#34;ECDSA verification failure&#34;)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	case signatureEd25519:
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		pubKey, ok := pubkey.(ed25519.PublicKey)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		if !ok {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;expected an Ed25519 public key, got %T&#34;, pubkey)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		if !ed25519.Verify(pubKey, signed, sig) {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			return errors.New(&#34;Ed25519 verification failure&#34;)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	case signaturePKCS1v15:
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		pubKey, ok := pubkey.(*rsa.PublicKey)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		if !ok {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;expected an RSA public key, got %T&#34;, pubkey)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		if err := rsa.VerifyPKCS1v15(pubKey, hashFunc, signed, sig); err != nil {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			return err
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	case signatureRSAPSS:
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		pubKey, ok := pubkey.(*rsa.PublicKey)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if !ok {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;expected an RSA public key, got %T&#34;, pubkey)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		signOpts := &amp;rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		if err := rsa.VerifyPSS(pubKey, hashFunc, signed, sig, signOpts); err != nil {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			return err
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	default:
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		return errors.New(&#34;internal error: unknown signature type&#34;)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	return nil
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>const (
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	serverSignatureContext = &#34;TLS 1.3, server CertificateVerify\x00&#34;
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	clientSignatureContext = &#34;TLS 1.3, client CertificateVerify\x00&#34;
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>var signaturePadding = []byte{
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// signedMessage returns the pre-hashed (if necessary) message to be signed by</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// certificate keys in TLS 1.3. See RFC 8446, Section 4.4.3.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func signedMessage(sigHash crypto.Hash, context string, transcript hash.Hash) []byte {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	if sigHash == directSigning {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		b := &amp;bytes.Buffer{}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		b.Write(signaturePadding)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		io.WriteString(b, context)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		b.Write(transcript.Sum(nil))
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return b.Bytes()
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	h := sigHash.New()
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	h.Write(signaturePadding)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	io.WriteString(h, context)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	h.Write(transcript.Sum(nil))
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return h.Sum(nil)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// typeAndHashFromSignatureScheme returns the corresponding signature type and</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// crypto.Hash for a given TLS SignatureScheme.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func typeAndHashFromSignatureScheme(signatureAlgorithm SignatureScheme) (sigType uint8, hash crypto.Hash, err error) {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	switch signatureAlgorithm {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	case PKCS1WithSHA1, PKCS1WithSHA256, PKCS1WithSHA384, PKCS1WithSHA512:
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		sigType = signaturePKCS1v15
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	case PSSWithSHA256, PSSWithSHA384, PSSWithSHA512:
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		sigType = signatureRSAPSS
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	case ECDSAWithSHA1, ECDSAWithP256AndSHA256, ECDSAWithP384AndSHA384, ECDSAWithP521AndSHA512:
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		sigType = signatureECDSA
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	case Ed25519:
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		sigType = signatureEd25519
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	default:
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		return 0, 0, fmt.Errorf(&#34;unsupported signature algorithm: %v&#34;, signatureAlgorithm)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	switch signatureAlgorithm {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	case PKCS1WithSHA1, ECDSAWithSHA1:
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		hash = crypto.SHA1
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	case PKCS1WithSHA256, PSSWithSHA256, ECDSAWithP256AndSHA256:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		hash = crypto.SHA256
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	case PKCS1WithSHA384, PSSWithSHA384, ECDSAWithP384AndSHA384:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		hash = crypto.SHA384
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	case PKCS1WithSHA512, PSSWithSHA512, ECDSAWithP521AndSHA512:
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		hash = crypto.SHA512
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	case Ed25519:
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		hash = directSigning
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	default:
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		return 0, 0, fmt.Errorf(&#34;unsupported signature algorithm: %v&#34;, signatureAlgorithm)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return sigType, hash, nil
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// legacyTypeAndHashFromPublicKey returns the fixed signature type and crypto.Hash for</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// a given public key used with TLS 1.0 and 1.1, before the introduction of</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// signature algorithm negotiation.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func legacyTypeAndHashFromPublicKey(pub crypto.PublicKey) (sigType uint8, hash crypto.Hash, err error) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	switch pub.(type) {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	case *rsa.PublicKey:
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return signaturePKCS1v15, crypto.MD5SHA1, nil
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case *ecdsa.PublicKey:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return signatureECDSA, crypto.SHA1, nil
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	case ed25519.PublicKey:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		<span class="comment">// RFC 8422 specifies support for Ed25519 in TLS 1.0 and 1.1,</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		<span class="comment">// but it requires holding on to a handshake transcript to do a</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// full signature, and not even OpenSSL bothers with the</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// complexity, so we can&#39;t even test it properly.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		return 0, 0, fmt.Errorf(&#34;tls: Ed25519 public keys are not supported before TLS 1.2&#34;)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	default:
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return 0, 0, fmt.Errorf(&#34;tls: unsupported public key: %T&#34;, pub)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>var rsaSignatureSchemes = []struct {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	scheme          SignatureScheme
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	minModulusBytes int
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	maxVersion      uint16
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>}{
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// RSA-PSS is used with PSSSaltLengthEqualsHash, and requires</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">//    emLen &gt;= hLen + sLen + 2</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	{PSSWithSHA256, crypto.SHA256.Size()*2 + 2, VersionTLS13},
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	{PSSWithSHA384, crypto.SHA384.Size()*2 + 2, VersionTLS13},
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	{PSSWithSHA512, crypto.SHA512.Size()*2 + 2, VersionTLS13},
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// PKCS #1 v1.5 uses prefixes from hashPrefixes in crypto/rsa, and requires</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">//    emLen &gt;= len(prefix) + hLen + 11</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// TLS 1.3 dropped support for PKCS #1 v1.5 in favor of RSA-PSS.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	{PKCS1WithSHA256, 19 + crypto.SHA256.Size() + 11, VersionTLS12},
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	{PKCS1WithSHA384, 19 + crypto.SHA384.Size() + 11, VersionTLS12},
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	{PKCS1WithSHA512, 19 + crypto.SHA512.Size() + 11, VersionTLS12},
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	{PKCS1WithSHA1, 15 + crypto.SHA1.Size() + 11, VersionTLS12},
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// signatureSchemesForCertificate returns the list of supported SignatureSchemes</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// for a given certificate, based on the public key and the protocol version,</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// and optionally filtered by its explicit SupportedSignatureAlgorithms.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// This function must be kept in sync with supportedSignatureAlgorithms.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// FIPS filtering is applied in the caller, selectSignatureScheme.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>func signatureSchemesForCertificate(version uint16, cert *Certificate) []SignatureScheme {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	priv, ok := cert.PrivateKey.(crypto.Signer)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if !ok {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		return nil
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	var sigAlgs []SignatureScheme
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	switch pub := priv.Public().(type) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	case *ecdsa.PublicKey:
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if version != VersionTLS13 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			<span class="comment">// In TLS 1.2 and earlier, ECDSA algorithms are not</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			<span class="comment">// constrained to a single curve.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			sigAlgs = []SignatureScheme{
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				ECDSAWithP256AndSHA256,
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>				ECDSAWithP384AndSHA384,
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				ECDSAWithP521AndSHA512,
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				ECDSAWithSHA1,
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			break
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		switch pub.Curve {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		case elliptic.P256():
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			sigAlgs = []SignatureScheme{ECDSAWithP256AndSHA256}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		case elliptic.P384():
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			sigAlgs = []SignatureScheme{ECDSAWithP384AndSHA384}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		case elliptic.P521():
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			sigAlgs = []SignatureScheme{ECDSAWithP521AndSHA512}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		default:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			return nil
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	case *rsa.PublicKey:
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		size := pub.Size()
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		sigAlgs = make([]SignatureScheme, 0, len(rsaSignatureSchemes))
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		for _, candidate := range rsaSignatureSchemes {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			if size &gt;= candidate.minModulusBytes &amp;&amp; version &lt;= candidate.maxVersion {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				sigAlgs = append(sigAlgs, candidate.scheme)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	case ed25519.PublicKey:
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		sigAlgs = []SignatureScheme{Ed25519}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	default:
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return nil
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if cert.SupportedSignatureAlgorithms != nil {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		var filteredSigAlgs []SignatureScheme
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		for _, sigAlg := range sigAlgs {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			if isSupportedSignatureAlgorithm(sigAlg, cert.SupportedSignatureAlgorithms) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>				filteredSigAlgs = append(filteredSigAlgs, sigAlg)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		return filteredSigAlgs
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	return sigAlgs
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// selectSignatureScheme picks a SignatureScheme from the peer&#39;s preference list</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// that works with the selected certificate. It&#39;s only called for protocol</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// versions that support signature algorithms, so TLS 1.2 and 1.3.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>func selectSignatureScheme(vers uint16, c *Certificate, peerAlgs []SignatureScheme) (SignatureScheme, error) {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	supportedAlgs := signatureSchemesForCertificate(vers, c)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if len(supportedAlgs) == 0 {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		return 0, unsupportedCertificateError(c)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if len(peerAlgs) == 0 &amp;&amp; vers == VersionTLS12 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// For TLS 1.2, if the client didn&#39;t send signature_algorithms then we</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		<span class="comment">// can assume that it supports SHA1. See RFC 5246, Section 7.4.1.4.1.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		peerAlgs = []SignatureScheme{PKCS1WithSHA1, ECDSAWithSHA1}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// Pick signature scheme in the peer&#39;s preference order, as our</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// preference order is not configurable.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	for _, preferredAlg := range peerAlgs {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		if needFIPS() &amp;&amp; !isSupportedSignatureAlgorithm(preferredAlg, fipsSupportedSignatureAlgorithms) {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			continue
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if isSupportedSignatureAlgorithm(preferredAlg, supportedAlgs) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			return preferredAlg, nil
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	return 0, errors.New(&#34;tls: peer doesn&#39;t support any of the certificate&#39;s signature algorithms&#34;)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// unsupportedCertificateError returns a helpful error for certificates with</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// an unsupported private key.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>func unsupportedCertificateError(cert *Certificate) error {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	switch cert.PrivateKey.(type) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	case rsa.PrivateKey, ecdsa.PrivateKey:
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: unsupported certificate: private key is %T, expected *%T&#34;,
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			cert.PrivateKey, cert.PrivateKey)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	case *ed25519.PrivateKey:
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: unsupported certificate: private key is *ed25519.PrivateKey, expected ed25519.PrivateKey&#34;)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	signer, ok := cert.PrivateKey.(crypto.Signer)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if !ok {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: certificate private key (%T) does not implement crypto.Signer&#34;,
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			cert.PrivateKey)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	switch pub := signer.Public().(type) {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	case *ecdsa.PublicKey:
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		switch pub.Curve {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		case elliptic.P256():
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		case elliptic.P384():
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		case elliptic.P521():
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		default:
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;tls: unsupported certificate curve (%s)&#34;, pub.Curve.Params().Name)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	case *rsa.PublicKey:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: certificate RSA key size too small for supported signature algorithms&#34;)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	case ed25519.PublicKey:
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	default:
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: unsupported certificate key (%T)&#34;, pub)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	if cert.SupportedSignatureAlgorithms != nil {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;tls: peer doesn&#39;t support the certificate custom signature algorithms&#34;)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	return fmt.Errorf(&#34;tls: internal error: unsupported key (%T)&#34;, cert.PrivateKey)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
</pre><p><a href="auth.go?m=text">View as plain text</a></p>

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
