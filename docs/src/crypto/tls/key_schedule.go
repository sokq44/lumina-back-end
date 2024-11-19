<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/tls/key_schedule.go - Go Documentation Server</title>

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
<a href="key_schedule.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/tls">tls</a>/<span class="text-muted">key_schedule.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto/ecdh&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;crypto/hmac&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/cryptobyte&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/hkdf&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// This file contains the functions necessary to compute the TLS 1.3 key</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// schedule. See RFC 8446, Section 7.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const (
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	resumptionBinderLabel         = &#34;res binder&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	clientEarlyTrafficLabel       = &#34;c e traffic&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	clientHandshakeTrafficLabel   = &#34;c hs traffic&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	serverHandshakeTrafficLabel   = &#34;s hs traffic&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	clientApplicationTrafficLabel = &#34;c ap traffic&#34;
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	serverApplicationTrafficLabel = &#34;s ap traffic&#34;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	exporterLabel                 = &#34;exp master&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	resumptionLabel               = &#34;res master&#34;
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	trafficUpdateLabel            = &#34;traffic upd&#34;
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// expandLabel implements HKDF-Expand-Label from RFC 8446, Section 7.1.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) expandLabel(secret []byte, label string, context []byte, length int) []byte {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	var hkdfLabel cryptobyte.Builder
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	hkdfLabel.AddUint16(uint16(length))
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	hkdfLabel.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		b.AddBytes([]byte(&#34;tls13 &#34;))
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		b.AddBytes([]byte(label))
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	})
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	hkdfLabel.AddUint8LengthPrefixed(func(b *cryptobyte.Builder) {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		b.AddBytes(context)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	})
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	hkdfLabelBytes, err := hkdfLabel.Bytes()
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	if err != nil {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		<span class="comment">// Rather than calling BytesOrPanic, we explicitly handle this error, in</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// order to provide a reasonable error message. It should be basically</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		<span class="comment">// impossible for this to panic, and routing errors back through the</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		<span class="comment">// tree rooted in this function is quite painful. The labels are fixed</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		<span class="comment">// size, and the context is either a fixed-length computed hash, or</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		<span class="comment">// parsed from a field which has the same length limitation. As such, an</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		<span class="comment">// error here is likely to only be caused during development.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		<span class="comment">// NOTE: another reasonable approach here might be to return a</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		<span class="comment">// randomized slice if we encounter an error, which would break the</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		<span class="comment">// connection, but avoid panicking. This would perhaps be safer but</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		<span class="comment">// significantly more confusing to users.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		panic(fmt.Errorf(&#34;failed to construct HKDF label: %s&#34;, err))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	out := make([]byte, length)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	n, err := hkdf.Expand(c.hash.New, secret, hkdfLabelBytes).Read(out)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if err != nil || n != length {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		panic(&#34;tls: HKDF-Expand-Label invocation failed unexpectedly&#34;)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return out
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// deriveSecret implements Derive-Secret from RFC 8446, Section 7.1.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) deriveSecret(secret []byte, label string, transcript hash.Hash) []byte {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if transcript == nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		transcript = c.hash.New()
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return c.expandLabel(secret, label, transcript.Sum(nil), c.hash.Size())
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// extract implements HKDF-Extract with the cipher suite hash.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) extract(newSecret, currentSecret []byte) []byte {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if newSecret == nil {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		newSecret = make([]byte, c.hash.Size())
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	return hkdf.Extract(c.hash.New, newSecret, currentSecret)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// nextTrafficSecret generates the next traffic secret, given the current one,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// according to RFC 8446, Section 7.2.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) nextTrafficSecret(trafficSecret []byte) []byte {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return c.expandLabel(trafficSecret, trafficUpdateLabel, nil, c.hash.Size())
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// trafficKey generates traffic keys according to RFC 8446, Section 7.3.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) trafficKey(trafficSecret []byte) (key, iv []byte) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	key = c.expandLabel(trafficSecret, &#34;key&#34;, nil, c.keyLen)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	iv = c.expandLabel(trafficSecret, &#34;iv&#34;, nil, aeadNonceLength)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	return
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// finishedHash generates the Finished verify_data or PskBinderEntry according</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// to RFC 8446, Section 4.4.4. See sections 4.4 and 4.2.11.2 for the baseKey</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// selection.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) finishedHash(baseKey []byte, transcript hash.Hash) []byte {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	finishedKey := c.expandLabel(baseKey, &#34;finished&#34;, nil, c.hash.Size())
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	verifyData := hmac.New(c.hash.New, finishedKey)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	verifyData.Write(transcript.Sum(nil))
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	return verifyData.Sum(nil)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// exportKeyingMaterial implements RFC5705 exporters for TLS 1.3 according to</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// RFC 8446, Section 7.5.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>func (c *cipherSuiteTLS13) exportKeyingMaterial(masterSecret []byte, transcript hash.Hash) func(string, []byte, int) ([]byte, error) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	expMasterSecret := c.deriveSecret(masterSecret, exporterLabel, transcript)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return func(label string, context []byte, length int) ([]byte, error) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		secret := c.deriveSecret(expMasterSecret, label, nil)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		h := c.hash.New()
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		h.Write(context)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		return c.expandLabel(secret, &#34;exporter&#34;, h.Sum(nil), length), nil
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// generateECDHEKey returns a PrivateKey that implements Diffie-Hellman</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// according to RFC 8446, Section 4.2.8.2.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func generateECDHEKey(rand io.Reader, curveID CurveID) (*ecdh.PrivateKey, error) {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	curve, ok := curveForCurveID(curveID)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	if !ok {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		return nil, errors.New(&#34;tls: internal error: unsupported curve&#34;)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	return curve.GenerateKey(rand)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func curveForCurveID(id CurveID) (ecdh.Curve, bool) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	switch id {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	case X25519:
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return ecdh.X25519(), true
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case CurveP256:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return ecdh.P256(), true
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	case CurveP384:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		return ecdh.P384(), true
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	case CurveP521:
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return ecdh.P521(), true
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	default:
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		return nil, false
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func curveIDForCurve(curve ecdh.Curve) (CurveID, bool) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	switch curve {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	case ecdh.X25519():
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		return X25519, true
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	case ecdh.P256():
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		return CurveP256, true
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	case ecdh.P384():
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return CurveP384, true
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	case ecdh.P521():
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		return CurveP521, true
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	default:
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return 0, false
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
</pre><p><a href="key_schedule.go?m=text">View as plain text</a></p>

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
