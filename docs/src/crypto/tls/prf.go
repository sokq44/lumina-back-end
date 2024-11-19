<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/prf.go - Go Documentation Server</title>

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
<a href="prf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">prf.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;crypto/hmac&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto/md5&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/sha1&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/sha256&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/sha512&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Split a premaster secret in two as specified in RFC 4346, Section 5.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func splitPreMasterSecret(secret []byte) (s1, s2 []byte) {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	s1 = secret[0 : (len(secret)+1)/2]
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	s2 = secret[len(secret)/2:]
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	return
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// pHash implements the P_hash function, as defined in RFC 4346, Section 5.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func pHash(result, secret, seed []byte, hash func() hash.Hash) {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	h := hmac.New(hash, secret)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	h.Write(seed)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	a := h.Sum(nil)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	j := 0
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	for j &lt; len(result) {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		h.Reset()
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		h.Write(a)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		h.Write(seed)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		b := h.Sum(nil)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		copy(result[j:], b)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		j += len(b)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		h.Reset()
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		h.Write(a)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		a = h.Sum(nil)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// prf10 implements the TLS 1.0 pseudo-random function, as defined in RFC 2246, Section 5.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func prf10(result, secret, label, seed []byte) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	hashSHA1 := sha1.New
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	hashMD5 := md5.New
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	labelAndSeed := make([]byte, len(label)+len(seed))
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	copy(labelAndSeed, label)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	copy(labelAndSeed[len(label):], seed)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	s1, s2 := splitPreMasterSecret(secret)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	pHash(result, s1, labelAndSeed, hashMD5)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	result2 := make([]byte, len(result))
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	pHash(result2, s2, labelAndSeed, hashSHA1)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	for i, b := range result2 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		result[i] ^= b
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// prf12 implements the TLS 1.2 pseudo-random function, as defined in RFC 5246, Section 5.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func prf12(hashFunc func() hash.Hash) func(result, secret, label, seed []byte) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return func(result, secret, label, seed []byte) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		labelAndSeed := make([]byte, len(label)+len(seed))
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		copy(labelAndSeed, label)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		copy(labelAndSeed[len(label):], seed)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		pHash(result, secret, labelAndSeed, hashFunc)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>const (
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	masterSecretLength   = 48 <span class="comment">// Length of a master secret in TLS 1.1.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	finishedVerifyLength = 12 <span class="comment">// Length of verify_data in a Finished message.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>var masterSecretLabel = []byte(&#34;master secret&#34;)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>var extendedMasterSecretLabel = []byte(&#34;extended master secret&#34;)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>var keyExpansionLabel = []byte(&#34;key expansion&#34;)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>var clientFinishedLabel = []byte(&#34;client finished&#34;)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>var serverFinishedLabel = []byte(&#34;server finished&#34;)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func prfAndHashForVersion(version uint16, suite *cipherSuite) (func(result, secret, label, seed []byte), crypto.Hash) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	switch version {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	case VersionTLS10, VersionTLS11:
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return prf10, crypto.Hash(0)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	case VersionTLS12:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if suite.flags&amp;suiteSHA384 != 0 {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			return prf12(sha512.New384), crypto.SHA384
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		return prf12(sha256.New), crypto.SHA256
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	default:
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		panic(&#34;unknown version&#34;)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func prfForVersion(version uint16, suite *cipherSuite) func(result, secret, label, seed []byte) {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	prf, _ := prfAndHashForVersion(version, suite)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	return prf
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// masterFromPreMasterSecret generates the master secret from the pre-master</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// secret. See RFC 5246, Section 8.1.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func masterFromPreMasterSecret(version uint16, suite *cipherSuite, preMasterSecret, clientRandom, serverRandom []byte) []byte {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	seed := make([]byte, 0, len(clientRandom)+len(serverRandom))
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	seed = append(seed, clientRandom...)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	seed = append(seed, serverRandom...)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	masterSecret := make([]byte, masterSecretLength)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	prfForVersion(version, suite)(masterSecret, preMasterSecret, masterSecretLabel, seed)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	return masterSecret
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// extMasterFromPreMasterSecret generates the extended master secret from the</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// pre-master secret. See RFC 7627.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func extMasterFromPreMasterSecret(version uint16, suite *cipherSuite, preMasterSecret, transcript []byte) []byte {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	masterSecret := make([]byte, masterSecretLength)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	prfForVersion(version, suite)(masterSecret, preMasterSecret, extendedMasterSecretLabel, transcript)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	return masterSecret
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// keysFromMasterSecret generates the connection keys from the master</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// secret, given the lengths of the MAC key, cipher key and IV, as defined in</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// RFC 2246, Section 6.3.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func keysFromMasterSecret(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte, macLen, keyLen, ivLen int) (clientMAC, serverMAC, clientKey, serverKey, clientIV, serverIV []byte) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	seed := make([]byte, 0, len(serverRandom)+len(clientRandom))
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	seed = append(seed, serverRandom...)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	seed = append(seed, clientRandom...)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	n := 2*macLen + 2*keyLen + 2*ivLen
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	keyMaterial := make([]byte, n)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	prfForVersion(version, suite)(keyMaterial, masterSecret, keyExpansionLabel, seed)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	clientMAC = keyMaterial[:macLen]
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	keyMaterial = keyMaterial[macLen:]
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	serverMAC = keyMaterial[:macLen]
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	keyMaterial = keyMaterial[macLen:]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	clientKey = keyMaterial[:keyLen]
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	keyMaterial = keyMaterial[keyLen:]
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	serverKey = keyMaterial[:keyLen]
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	keyMaterial = keyMaterial[keyLen:]
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	clientIV = keyMaterial[:ivLen]
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	keyMaterial = keyMaterial[ivLen:]
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	serverIV = keyMaterial[:ivLen]
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	return
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func newFinishedHash(version uint16, cipherSuite *cipherSuite) finishedHash {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	var buffer []byte
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if version &gt;= VersionTLS12 {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		buffer = []byte{}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	prf, hash := prfAndHashForVersion(version, cipherSuite)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if hash != 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return finishedHash{hash.New(), hash.New(), nil, nil, buffer, version, prf}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	return finishedHash{sha1.New(), sha1.New(), md5.New(), md5.New(), buffer, version, prf}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// A finishedHash calculates the hash of a set of handshake messages suitable</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span><span class="comment">// for including in a Finished message.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>type finishedHash struct {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	client hash.Hash
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	server hash.Hash
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Prior to TLS 1.2, an additional MD5 hash is required.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	clientMD5 hash.Hash
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	serverMD5 hash.Hash
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// In TLS 1.2, a full buffer is sadly required.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	buffer []byte
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	version uint16
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	prf     func(result, secret, label, seed []byte)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>func (h *finishedHash) Write(msg []byte) (n int, err error) {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	h.client.Write(msg)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	h.server.Write(msg)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if h.version &lt; VersionTLS12 {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		h.clientMD5.Write(msg)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		h.serverMD5.Write(msg)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if h.buffer != nil {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		h.buffer = append(h.buffer, msg...)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return len(msg), nil
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>func (h finishedHash) Sum() []byte {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if h.version &gt;= VersionTLS12 {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return h.client.Sum(nil)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	out := make([]byte, 0, md5.Size+sha1.Size)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	out = h.clientMD5.Sum(out)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return h.client.Sum(out)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// clientSum returns the contents of the verify_data member of a client&#39;s</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// Finished message.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (h finishedHash) clientSum(masterSecret []byte) []byte {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	out := make([]byte, finishedVerifyLength)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	h.prf(out, masterSecret, clientFinishedLabel, h.Sum())
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	return out
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// serverSum returns the contents of the verify_data member of a server&#39;s</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// Finished message.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (h finishedHash) serverSum(masterSecret []byte) []byte {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	out := make([]byte, finishedVerifyLength)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	h.prf(out, masterSecret, serverFinishedLabel, h.Sum())
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	return out
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// hashForClientCertificate returns the handshake messages so far, pre-hashed if</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// necessary, suitable for signing by a TLS client certificate.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func (h finishedHash) hashForClientCertificate(sigType uint8, hashAlg crypto.Hash) []byte {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if (h.version &gt;= VersionTLS12 || sigType == signatureEd25519) &amp;&amp; h.buffer == nil {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		panic(&#34;tls: handshake hash for a client certificate requested after discarding the handshake buffer&#34;)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if sigType == signatureEd25519 {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return h.buffer
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if h.version &gt;= VersionTLS12 {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		hash := hashAlg.New()
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		hash.Write(h.buffer)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		return hash.Sum(nil)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if sigType == signatureECDSA {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return h.server.Sum(nil)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	return h.Sum()
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// discardHandshakeBuffer is called when there is no more need to</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// buffer the entirety of the handshake messages.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (h *finishedHash) discardHandshakeBuffer() {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	h.buffer = nil
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// noEKMBecauseRenegotiation is used as a value of</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// ConnectionState.ekm when renegotiation is enabled and thus</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// we wish to fail all key-material export requests.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func noEKMBecauseRenegotiation(label string, context []byte, length int) ([]byte, error) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	return nil, errors.New(&#34;crypto/tls: ExportKeyingMaterial is unavailable when renegotiation is enabled&#34;)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// noEKMBecauseNoEMS is used as a value of ConnectionState.ekm when Extended</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// Master Secret is not negotiated and thus we wish to fail all key-material</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// export requests.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>func noEKMBecauseNoEMS(label string, context []byte, length int) ([]byte, error) {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	return nil, errors.New(&#34;crypto/tls: ExportKeyingMaterial is unavailable when neither TLS 1.3 nor Extended Master Secret are negotiated; override with GODEBUG=tlsunsafeekm=1&#34;)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// ekmFromMasterSecret generates exported keying material as defined in RFC 5705.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>func ekmFromMasterSecret(version uint16, suite *cipherSuite, masterSecret, clientRandom, serverRandom []byte) func(string, []byte, int) ([]byte, error) {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	return func(label string, context []byte, length int) ([]byte, error) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		switch label {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		case &#34;client finished&#34;, &#34;server finished&#34;, &#34;master secret&#34;, &#34;key expansion&#34;:
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			<span class="comment">// These values are reserved and may not be used.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;crypto/tls: reserved ExportKeyingMaterial label: %s&#34;, label)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		seedLen := len(serverRandom) + len(clientRandom)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		if context != nil {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			seedLen += 2 + len(context)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		seed := make([]byte, 0, seedLen)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		seed = append(seed, clientRandom...)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		seed = append(seed, serverRandom...)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		if context != nil {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			if len(context) &gt;= 1&lt;&lt;16 {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>				return nil, fmt.Errorf(&#34;crypto/tls: ExportKeyingMaterial context too long&#34;)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			seed = append(seed, byte(len(context)&gt;&gt;8), byte(len(context)))
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			seed = append(seed, context...)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		keyMaterial := make([]byte, length)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		prfForVersion(version, suite)(keyMaterial, masterSecret, []byte(label), seed)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		return keyMaterial, nil
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
</pre><p><a href="prf.go?m=text">View as plain text</a></p>

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
