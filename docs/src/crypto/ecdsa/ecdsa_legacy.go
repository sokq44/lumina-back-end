<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/ecdsa/ecdsa_legacy.go - Go Documentation Server</title>

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
<a href="ecdsa_legacy.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/ecdsa">ecdsa</a>/<span class="text-muted">ecdsa_legacy.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/ecdsa">crypto/ecdsa</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package ecdsa
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto/elliptic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;math/big&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/cryptobyte&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/cryptobyte/asn1&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// This file contains a math/big implementation of ECDSA that is only used for</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// deprecated custom curves.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>func generateLegacy(c elliptic.Curve, rand io.Reader) (*PrivateKey, error) {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	k, err := randFieldElement(c, rand)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	if err != nil {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		return nil, err
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	priv := new(PrivateKey)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	priv.PublicKey.Curve = c
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	priv.D = k
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	return priv, nil
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// hashToInt converts a hash value to an integer. Per FIPS 186-4, Section 6.4,</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// we use the left-most bits of the hash to match the bit-length of the order of</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// the curve. This also performs Step 5 of SEC 1, Version 2.0, Section 4.1.3.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	orderBits := c.Params().N.BitLen()
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	orderBytes := (orderBits + 7) / 8
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if len(hash) &gt; orderBytes {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		hash = hash[:orderBytes]
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	ret := new(big.Int).SetBytes(hash)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	excess := len(hash)*8 - orderBits
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	if excess &gt; 0 {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		ret.Rsh(ret, uint(excess))
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	return ret
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>var errZeroParam = errors.New(&#34;zero parameter&#34;)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// Sign signs a hash (which should be the result of hashing a larger message)</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// using the private key, priv. If the hash is longer than the bit-length of the</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// private key&#39;s curve order, the hash will be truncated to that length. It</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// returns the signature as a pair of integers. Most applications should use</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// [SignASN1] instead of dealing directly with r, s.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	sig, err := SignASN1(rand, priv, hash)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if err != nil {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return nil, nil, err
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	r, s = new(big.Int), new(big.Int)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	var inner cryptobyte.String
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	input := cryptobyte.String(sig)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if !input.ReadASN1(&amp;inner, asn1.SEQUENCE) ||
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		!input.Empty() ||
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		!inner.ReadASN1Integer(r) ||
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		!inner.ReadASN1Integer(s) ||
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		!inner.Empty() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return nil, nil, errors.New(&#34;invalid ASN.1 from SignASN1&#34;)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return r, s, nil
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func signLegacy(priv *PrivateKey, csprng io.Reader, hash []byte) (sig []byte, err error) {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	c := priv.Curve
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// SEC 1, Version 2.0, Section 4.1.3</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	N := c.Params().N
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	if N.Sign() == 0 {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		return nil, errZeroParam
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	var k, kInv, r, s *big.Int
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	for {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		for {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			k, err = randFieldElement(c, csprng)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			if err != nil {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				return nil, err
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			kInv = new(big.Int).ModInverse(k, N)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			r, _ = c.ScalarBaseMult(k.Bytes())
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			r.Mod(r, N)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			if r.Sign() != 0 {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>				break
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		e := hashToInt(hash, c)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		s = new(big.Int).Mul(priv.D, r)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		s.Add(s, e)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		s.Mul(s, kInv)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		s.Mod(s, N) <span class="comment">// N != 0</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if s.Sign() != 0 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			break
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return encodeSignature(r.Bytes(), s.Bytes())
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// Verify verifies the signature in r, s of hash using the public key, pub. Its</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// return value records whether the signature is valid. Most applications should</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// use VerifyASN1 instead of dealing directly with r, s.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	if r.Sign() &lt;= 0 || s.Sign() &lt;= 0 {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		return false
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	sig, err := encodeSignature(r.Bytes(), s.Bytes())
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if err != nil {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return false
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return VerifyASN1(pub, hash, sig)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func verifyLegacy(pub *PublicKey, hash []byte, sig []byte) bool {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	rBytes, sBytes, err := parseSignature(sig)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if err != nil {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return false
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	r, s := new(big.Int).SetBytes(rBytes), new(big.Int).SetBytes(sBytes)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	c := pub.Curve
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	N := c.Params().N
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if r.Sign() &lt;= 0 || s.Sign() &lt;= 0 {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		return false
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	if r.Cmp(N) &gt;= 0 || s.Cmp(N) &gt;= 0 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		return false
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// SEC 1, Version 2.0, Section 4.1.4</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	e := hashToInt(hash, c)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	w := new(big.Int).ModInverse(s, N)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	u1 := e.Mul(e, w)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	u1.Mod(u1, N)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	u2 := w.Mul(r, w)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	u2.Mod(u2, N)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	x1, y1 := c.ScalarBaseMult(u1.Bytes())
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	x2, y2 := c.ScalarMult(pub.X, pub.Y, u2.Bytes())
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	x, y := c.Add(x1, y1, x2, y2)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if x.Sign() == 0 &amp;&amp; y.Sign() == 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return false
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	x.Mod(x, N)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	return x.Cmp(r) == 0
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>var one = new(big.Int).SetInt64(1)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span><span class="comment">// randFieldElement returns a random element of the order of the given</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// curve using the procedure given in FIPS 186-4, Appendix B.5.2.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// See randomPoint for notes on the algorithm. This has to match, or s390x</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// signatures will come out different from other architectures, which will</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// break TLS recorded tests.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	for {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		N := c.Params().N
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		b := make([]byte, (N.BitLen()+7)/8)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if _, err = io.ReadFull(rand, b); err != nil {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			return
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		if excess := len(b)*8 - N.BitLen(); excess &gt; 0 {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			b[0] &gt;&gt;= excess
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		k = new(big.Int).SetBytes(b)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		if k.Sign() != 0 &amp;&amp; k.Cmp(N) &lt; 0 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			return
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
</pre><p><a href="ecdsa_legacy.go?m=text">View as plain text</a></p>

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
