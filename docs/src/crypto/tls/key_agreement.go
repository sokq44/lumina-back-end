<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/key_agreement.go - Go Documentation Server</title>

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
<a href="key_agreement.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">key_agreement.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/tls">crypto/tls</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package tls
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;crypto/ecdh&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto/md5&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/rsa&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/sha1&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/x509&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// a keyAgreement implements the client and server side of a TLS key agreement</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// protocol by generating and processing key exchange messages.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type keyAgreement interface {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// On the server side, the first two methods are called in order.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// In the case that the key agreement protocol doesn&#39;t use a</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// ServerKeyExchange message, generateServerKeyExchange can return nil,</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// nil.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	generateServerKeyExchange(*Config, *Certificate, *clientHelloMsg, *serverHelloMsg) (*serverKeyExchangeMsg, error)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	processClientKeyExchange(*Config, *Certificate, *clientKeyExchangeMsg, uint16) ([]byte, error)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// On the client side, the next two methods are called in order.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// This method may not be called if the server doesn&#39;t send a</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// ServerKeyExchange message.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	processServerKeyExchange(*Config, *clientHelloMsg, *serverHelloMsg, *x509.Certificate, *serverKeyExchangeMsg) error
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	generateClientKeyExchange(*Config, *clientHelloMsg, *x509.Certificate) ([]byte, *clientKeyExchangeMsg, error)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>var errClientKeyExchange = errors.New(&#34;tls: invalid ClientKeyExchange message&#34;)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>var errServerKeyExchange = errors.New(&#34;tls: invalid ServerKeyExchange message&#34;)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// rsaKeyAgreement implements the standard TLS key agreement where the client</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// encrypts the pre-master secret to the server&#39;s public key.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>type rsaKeyAgreement struct{}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (ka rsaKeyAgreement) generateServerKeyExchange(config *Config, cert *Certificate, clientHello *clientHelloMsg, hello *serverHelloMsg) (*serverKeyExchangeMsg, error) {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	return nil, nil
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (ka rsaKeyAgreement) processClientKeyExchange(config *Config, cert *Certificate, ckx *clientKeyExchangeMsg, version uint16) ([]byte, error) {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if len(ckx.ciphertext) &lt; 2 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		return nil, errClientKeyExchange
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	ciphertextLen := int(ckx.ciphertext[0])&lt;&lt;8 | int(ckx.ciphertext[1])
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	if ciphertextLen != len(ckx.ciphertext)-2 {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		return nil, errClientKeyExchange
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	ciphertext := ckx.ciphertext[2:]
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	priv, ok := cert.PrivateKey.(crypto.Decrypter)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if !ok {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: certificate private key does not implement crypto.Decrypter&#34;)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// Perform constant time RSA PKCS #1 v1.5 decryption</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	preMasterSecret, err := priv.Decrypt(config.rand(), ciphertext, &amp;rsa.PKCS1v15DecryptOptions{SessionKeyLen: 48})
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if err != nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return nil, err
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// We don&#39;t check the version number in the premaster secret. For one,</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// by checking it, we would leak information about the validity of the</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// encrypted pre-master secret. Secondly, it provides only a small</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// benefit against a downgrade attack and some implementations send the</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// wrong version anyway. See the discussion at the end of section</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// 7.4.7.1 of RFC 4346.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return preMasterSecret, nil
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func (ka rsaKeyAgreement) processServerKeyExchange(config *Config, clientHello *clientHelloMsg, serverHello *serverHelloMsg, cert *x509.Certificate, skx *serverKeyExchangeMsg) error {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	return errors.New(&#34;tls: unexpected ServerKeyExchange&#34;)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func (ka rsaKeyAgreement) generateClientKeyExchange(config *Config, clientHello *clientHelloMsg, cert *x509.Certificate) ([]byte, *clientKeyExchangeMsg, error) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	preMasterSecret := make([]byte, 48)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	preMasterSecret[0] = byte(clientHello.vers &gt;&gt; 8)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	preMasterSecret[1] = byte(clientHello.vers)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	_, err := io.ReadFull(config.rand(), preMasterSecret[2:])
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if err != nil {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return nil, nil, err
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if !ok {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		return nil, nil, errors.New(&#34;tls: server certificate contains incorrect key type for selected ciphersuite&#34;)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	encrypted, err := rsa.EncryptPKCS1v15(config.rand(), rsaKey, preMasterSecret)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	if err != nil {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		return nil, nil, err
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	ckx := new(clientKeyExchangeMsg)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	ckx.ciphertext = make([]byte, len(encrypted)+2)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	ckx.ciphertext[0] = byte(len(encrypted) &gt;&gt; 8)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	ckx.ciphertext[1] = byte(len(encrypted))
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	copy(ckx.ciphertext[2:], encrypted)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	return preMasterSecret, ckx, nil
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// sha1Hash calculates a SHA1 hash over the given byte slices.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func sha1Hash(slices [][]byte) []byte {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	hsha1 := sha1.New()
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	for _, slice := range slices {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		hsha1.Write(slice)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return hsha1.Sum(nil)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// md5SHA1Hash implements TLS 1.0&#39;s hybrid hash function which consists of the</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// concatenation of an MD5 and SHA1 hash.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>func md5SHA1Hash(slices [][]byte) []byte {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	md5sha1 := make([]byte, md5.Size+sha1.Size)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	hmd5 := md5.New()
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	for _, slice := range slices {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		hmd5.Write(slice)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	copy(md5sha1, hmd5.Sum(nil))
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	copy(md5sha1[md5.Size:], sha1Hash(slices))
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return md5sha1
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// hashForServerKeyExchange hashes the given slices and returns their digest</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// using the given hash function (for &gt;= TLS 1.2) or using a default based on</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// the sigType (for earlier TLS versions). For Ed25519 signatures, which don&#39;t</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// do pre-hashing, it returns the concatenation of the slices.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>func hashForServerKeyExchange(sigType uint8, hashFunc crypto.Hash, version uint16, slices ...[]byte) []byte {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if sigType == signatureEd25519 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		var signed []byte
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		for _, slice := range slices {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			signed = append(signed, slice...)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		return signed
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if version &gt;= VersionTLS12 {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		h := hashFunc.New()
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		for _, slice := range slices {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			h.Write(slice)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		digest := h.Sum(nil)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		return digest
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if sigType == signatureECDSA {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		return sha1Hash(slices)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	return md5SHA1Hash(slices)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// ecdheKeyAgreement implements a TLS key agreement where the server</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// generates an ephemeral EC public/private key pair and signs it. The</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// pre-master secret is then calculated using ECDH. The signature may</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// be ECDSA, Ed25519 or RSA.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>type ecdheKeyAgreement struct {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	version uint16
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	isRSA   bool
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	key     *ecdh.PrivateKey
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// ckx and preMasterSecret are generated in processServerKeyExchange</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// and returned in generateClientKeyExchange.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	ckx             *clientKeyExchangeMsg
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	preMasterSecret []byte
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>func (ka *ecdheKeyAgreement) generateServerKeyExchange(config *Config, cert *Certificate, clientHello *clientHelloMsg, hello *serverHelloMsg) (*serverKeyExchangeMsg, error) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	var curveID CurveID
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	for _, c := range clientHello.supportedCurves {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if config.supportsCurve(c) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			curveID = c
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			break
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if curveID == 0 {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: no supported elliptic curves offered&#34;)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	if _, ok := curveForCurveID(curveID); !ok {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: CurvePreferences includes unsupported curve&#34;)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	key, err := generateECDHEKey(config.rand(), curveID)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if err != nil {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		return nil, err
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	ka.key = key
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// See RFC 4492, Section 5.4.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	ecdhePublic := key.PublicKey().Bytes()
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	serverECDHEParams := make([]byte, 1+2+1+len(ecdhePublic))
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	serverECDHEParams[0] = 3 <span class="comment">// named curve</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	serverECDHEParams[1] = byte(curveID &gt;&gt; 8)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	serverECDHEParams[2] = byte(curveID)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	serverECDHEParams[3] = byte(len(ecdhePublic))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	copy(serverECDHEParams[4:], ecdhePublic)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	priv, ok := cert.PrivateKey.(crypto.Signer)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if !ok {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;tls: certificate private key of type %T does not implement crypto.Signer&#34;, cert.PrivateKey)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	var signatureAlgorithm SignatureScheme
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	var sigType uint8
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	var sigHash crypto.Hash
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if ka.version &gt;= VersionTLS12 {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		signatureAlgorithm, err = selectSignatureScheme(ka.version, cert, clientHello.supportedSignatureAlgorithms)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		if err != nil {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			return nil, err
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		sigType, sigHash, err = typeAndHashFromSignatureScheme(signatureAlgorithm)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		if err != nil {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			return nil, err
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	} else {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		sigType, sigHash, err = legacyTypeAndHashFromPublicKey(priv.Public())
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if err != nil {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			return nil, err
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if (sigType == signaturePKCS1v15 || sigType == signatureRSAPSS) != ka.isRSA {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: certificate cannot be used with the selected cipher suite&#34;)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	signed := hashForServerKeyExchange(sigType, sigHash, ka.version, clientHello.random, hello.random, serverECDHEParams)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	signOpts := crypto.SignerOpts(sigHash)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if sigType == signatureRSAPSS {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		signOpts = &amp;rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: sigHash}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	sig, err := priv.Sign(config.rand(), signed, signOpts)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if err != nil {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: failed to sign ECDHE parameters: &#34; + err.Error())
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	skx := new(serverKeyExchangeMsg)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	sigAndHashLen := 0
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if ka.version &gt;= VersionTLS12 {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		sigAndHashLen = 2
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	skx.key = make([]byte, len(serverECDHEParams)+sigAndHashLen+2+len(sig))
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	copy(skx.key, serverECDHEParams)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	k := skx.key[len(serverECDHEParams):]
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	if ka.version &gt;= VersionTLS12 {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		k[0] = byte(signatureAlgorithm &gt;&gt; 8)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		k[1] = byte(signatureAlgorithm)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		k = k[2:]
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	k[0] = byte(len(sig) &gt;&gt; 8)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	k[1] = byte(len(sig))
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	copy(k[2:], sig)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	return skx, nil
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func (ka *ecdheKeyAgreement) processClientKeyExchange(config *Config, cert *Certificate, ckx *clientKeyExchangeMsg, version uint16) ([]byte, error) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	if len(ckx.ciphertext) == 0 || int(ckx.ciphertext[0]) != len(ckx.ciphertext)-1 {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return nil, errClientKeyExchange
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	peerKey, err := ka.key.Curve().NewPublicKey(ckx.ciphertext[1:])
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if err != nil {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		return nil, errClientKeyExchange
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	preMasterSecret, err := ka.key.ECDH(peerKey)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if err != nil {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return nil, errClientKeyExchange
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	return preMasterSecret, nil
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func (ka *ecdheKeyAgreement) processServerKeyExchange(config *Config, clientHello *clientHelloMsg, serverHello *serverHelloMsg, cert *x509.Certificate, skx *serverKeyExchangeMsg) error {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if len(skx.key) &lt; 4 {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if skx.key[0] != 3 { <span class="comment">// named curve</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected unsupported curve&#34;)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	curveID := CurveID(skx.key[1])&lt;&lt;8 | CurveID(skx.key[2])
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	publicLen := int(skx.key[3])
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	if publicLen+4 &gt; len(skx.key) {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	serverECDHEParams := skx.key[:4+publicLen]
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	publicKey := serverECDHEParams[4:]
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	sig := skx.key[4+publicLen:]
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if len(sig) &lt; 2 {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	if _, ok := curveForCurveID(curveID); !ok {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		return errors.New(&#34;tls: server selected unsupported curve&#34;)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	key, err := generateECDHEKey(config.rand(), curveID)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	if err != nil {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		return err
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	ka.key = key
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	peerKey, err := key.Curve().NewPublicKey(publicKey)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if err != nil {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	ka.preMasterSecret, err = key.ECDH(peerKey)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if err != nil {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	ourPublicKey := key.PublicKey().Bytes()
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	ka.ckx = new(clientKeyExchangeMsg)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	ka.ckx.ciphertext = make([]byte, 1+len(ourPublicKey))
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	ka.ckx.ciphertext[0] = byte(len(ourPublicKey))
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	copy(ka.ckx.ciphertext[1:], ourPublicKey)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	var sigType uint8
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	var sigHash crypto.Hash
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	if ka.version &gt;= VersionTLS12 {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		signatureAlgorithm := SignatureScheme(sig[0])&lt;&lt;8 | SignatureScheme(sig[1])
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		sig = sig[2:]
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		if len(sig) &lt; 2 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			return errServerKeyExchange
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		if !isSupportedSignatureAlgorithm(signatureAlgorithm, clientHello.supportedSignatureAlgorithms) {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			return errors.New(&#34;tls: certificate used with invalid signature algorithm&#34;)
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		sigType, sigHash, err = typeAndHashFromSignatureScheme(signatureAlgorithm)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		if err != nil {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			return err
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	} else {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		sigType, sigHash, err = legacyTypeAndHashFromPublicKey(cert.PublicKey)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		if err != nil {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			return err
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if (sigType == signaturePKCS1v15 || sigType == signatureRSAPSS) != ka.isRSA {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	sigLen := int(sig[0])&lt;&lt;8 | int(sig[1])
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	if sigLen+2 != len(sig) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		return errServerKeyExchange
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	sig = sig[2:]
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	signed := hashForServerKeyExchange(sigType, sigHash, ka.version, clientHello.random, serverHello.random, serverECDHEParams)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	if err := verifyHandshakeSignature(sigType, cert.PublicKey, sigHash, signed, sig); err != nil {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		return errors.New(&#34;tls: invalid signature by the server certificate: &#34; + err.Error())
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	return nil
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>func (ka *ecdheKeyAgreement) generateClientKeyExchange(config *Config, clientHello *clientHelloMsg, cert *x509.Certificate) ([]byte, *clientKeyExchangeMsg, error) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	if ka.ckx == nil {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		return nil, nil, errors.New(&#34;tls: missing ServerKeyExchange message&#34;)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	return ka.preMasterSecret, ka.ckx, nil
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
</pre><p><a href="key_agreement.go?m=text">View as plain text</a></p>

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
