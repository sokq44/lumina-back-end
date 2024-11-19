<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/elliptic/nistec.go - Go Documentation Server</title>

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
<a href="nistec.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/elliptic">elliptic</a>/<span class="text-muted">nistec.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/elliptic">crypto/elliptic</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package elliptic
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto/internal/nistec&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;math/big&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>var p224 = &amp;nistCurve[*nistec.P224Point]{
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	newPoint: nistec.NewP224Point,
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>}
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>func initP224() {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	p224.params = &amp;CurveParams{
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>		Name:    &#34;P-224&#34;,
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		BitSize: 224,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		<span class="comment">// FIPS 186-4, section D.1.2.2</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>		P:  bigFromDecimal(&#34;26959946667150639794667015087019630673557916260026308143510066298881&#34;),
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		N:  bigFromDecimal(&#34;26959946667150639794667015087019625940457807714424391721682722368061&#34;),
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>		B:  bigFromHex(&#34;b4050a850c04b3abf54132565044b0b7d7bfd8ba270b39432355ffb4&#34;),
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		Gx: bigFromHex(&#34;b70e0cbd6bb4bf7f321390b94a03c1d356c21122343280d6115c1d21&#34;),
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		Gy: bigFromHex(&#34;bd376388b5f723fb4c22dfe6cd4375a05a07476444d5819985007e34&#34;),
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>type p256Curve struct {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	nistCurve[*nistec.P256Point]
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>var p256 = &amp;p256Curve{nistCurve[*nistec.P256Point]{
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	newPoint: nistec.NewP256Point,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func initP256() {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	p256.params = &amp;CurveParams{
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		Name:    &#34;P-256&#34;,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		BitSize: 256,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		<span class="comment">// FIPS 186-4, section D.1.2.3</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		P:  bigFromDecimal(&#34;115792089210356248762697446949407573530086143415290314195533631308867097853951&#34;),
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		N:  bigFromDecimal(&#34;115792089210356248762697446949407573529996955224135760342422259061068512044369&#34;),
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		B:  bigFromHex(&#34;5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b&#34;),
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		Gx: bigFromHex(&#34;6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296&#34;),
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		Gy: bigFromHex(&#34;4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5&#34;),
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>var p384 = &amp;nistCurve[*nistec.P384Point]{
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	newPoint: nistec.NewP384Point,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func initP384() {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	p384.params = &amp;CurveParams{
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		Name:    &#34;P-384&#34;,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		BitSize: 384,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		<span class="comment">// FIPS 186-4, section D.1.2.4</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		P: bigFromDecimal(&#34;394020061963944792122790401001436138050797392704654&#34; +
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			&#34;46667948293404245721771496870329047266088258938001861606973112319&#34;),
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		N: bigFromDecimal(&#34;394020061963944792122790401001436138050797392704654&#34; +
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			&#34;46667946905279627659399113263569398956308152294913554433653942643&#34;),
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		B: bigFromHex(&#34;b3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088&#34; +
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			&#34;f5013875ac656398d8a2ed19d2a85c8edd3ec2aef&#34;),
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		Gx: bigFromHex(&#34;aa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741&#34; +
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			&#34;e082542a385502f25dbf55296c3a545e3872760ab7&#34;),
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		Gy: bigFromHex(&#34;3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da31&#34; +
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			&#34;13b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f&#34;),
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>var p521 = &amp;nistCurve[*nistec.P521Point]{
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	newPoint: nistec.NewP521Point,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>func initP521() {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	p521.params = &amp;CurveParams{
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		Name:    &#34;P-521&#34;,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		BitSize: 521,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// FIPS 186-4, section D.1.2.5</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		P: bigFromDecimal(&#34;68647976601306097149819007990813932172694353001433&#34; +
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			&#34;0540939446345918554318339765605212255964066145455497729631139148&#34; +
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			&#34;0858037121987999716643812574028291115057151&#34;),
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		N: bigFromDecimal(&#34;68647976601306097149819007990813932172694353001433&#34; +
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			&#34;0540939446345918554318339765539424505774633321719753296399637136&#34; +
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			&#34;3321113864768612440380340372808892707005449&#34;),
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		B: bigFromHex(&#34;0051953eb9618e1c9a1f929a21a0b68540eea2da725b99b315f3b8&#34; +
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			&#34;b489918ef109e156193951ec7e937b1652c0bd3bb1bf073573df883d2c34f1ef&#34; +
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			&#34;451fd46b503f00&#34;),
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		Gx: bigFromHex(&#34;00c6858e06b70404e9cd9e3ecb662395b4429c648139053fb521f8&#34; +
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			&#34;28af606b4d3dbaa14b5e77efe75928fe1dc127a2ffa8de3348b3c1856a429bf9&#34; +
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			&#34;7e7e31c2e5bd66&#34;),
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		Gy: bigFromHex(&#34;011839296a789a3bc0045c8a5fb42c7d1bd998f54449579b446817&#34; +
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			&#34;afbd17273e662c97ee72995ef42640c550b9013fad0761353c7086a272c24088&#34; +
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			&#34;be94769fd16650&#34;),
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// nistCurve is a Curve implementation based on a nistec Point.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// It&#39;s a wrapper that exposes the big.Int-based Curve interface and encodes the</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// legacy idiosyncrasies it requires, such as invalid and infinity point</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// handling.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// To interact with the nistec package, points are encoded into and decoded from</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// properly formatted byte slices. All big.Int use is limited to this package.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// Encoding and decoding is 1/1000th of the runtime of a scalar multiplication,</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// so the overhead is acceptable.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>type nistCurve[Point nistPoint[Point]] struct {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	newPoint func() Point
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	params   *CurveParams
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// nistPoint is a generic constraint for the nistec Point types.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>type nistPoint[T any] interface {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	Bytes() []byte
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	SetBytes([]byte) (T, error)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	Add(T, T) T
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	Double(T) T
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	ScalarMult(T, []byte) (T, error)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	ScalarBaseMult([]byte) (T, error)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) Params() *CurveParams {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return curve.params
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) IsOnCurve(x, y *big.Int) bool {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// IsOnCurve is documented to reject (0, 0), the conventional point at</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// infinity, which however is accepted by pointFromAffine.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if x.Sign() == 0 &amp;&amp; y.Sign() == 0 {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return false
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	_, err := curve.pointFromAffine(x, y)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	return err == nil
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) pointFromAffine(x, y *big.Int) (p Point, err error) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// (0, 0) is by convention the point at infinity, which can&#39;t be represented</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// in affine coordinates. See Issue 37294.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	if x.Sign() == 0 &amp;&amp; y.Sign() == 0 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		return curve.newPoint(), nil
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// Reject values that would not get correctly encoded.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if x.Sign() &lt; 0 || y.Sign() &lt; 0 {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return p, errors.New(&#34;negative coordinate&#34;)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if x.BitLen() &gt; curve.params.BitSize || y.BitLen() &gt; curve.params.BitSize {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		return p, errors.New(&#34;overflowing coordinate&#34;)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// Encode the coordinates and let SetBytes reject invalid points.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	byteLen := (curve.params.BitSize + 7) / 8
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	buf := make([]byte, 1+2*byteLen)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	buf[0] = 4 <span class="comment">// uncompressed point</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	x.FillBytes(buf[1 : 1+byteLen])
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	y.FillBytes(buf[1+byteLen : 1+2*byteLen])
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	return curve.newPoint().SetBytes(buf)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) pointToAffine(p Point) (x, y *big.Int) {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	out := p.Bytes()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if len(out) == 1 &amp;&amp; out[0] == 0 {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// This is the encoding of the point at infinity, which the affine</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// coordinates API represents as (0, 0) by convention.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		return new(big.Int), new(big.Int)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	byteLen := (curve.params.BitSize + 7) / 8
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	x = new(big.Int).SetBytes(out[1 : 1+byteLen])
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	y = new(big.Int).SetBytes(out[1+byteLen:])
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	return x, y
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	p1, err := curve.pointFromAffine(x1, y1)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if err != nil {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: Add was called on an invalid point&#34;)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	p2, err := curve.pointFromAffine(x2, y2)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if err != nil {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: Add was called on an invalid point&#34;)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	return curve.pointToAffine(p1.Add(p1, p2))
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	p, err := curve.pointFromAffine(x1, y1)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if err != nil {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: Double was called on an invalid point&#34;)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	return curve.pointToAffine(p.Double(p))
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// normalizeScalar brings the scalar within the byte size of the order of the</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">// curve, as expected by the nistec scalar multiplication functions.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) normalizeScalar(scalar []byte) []byte {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	byteSize := (curve.params.N.BitLen() + 7) / 8
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if len(scalar) == byteSize {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		return scalar
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	s := new(big.Int).SetBytes(scalar)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if len(scalar) &gt; byteSize {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		s.Mod(s, curve.params.N)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	out := make([]byte, byteSize)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return s.FillBytes(out)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) ScalarMult(Bx, By *big.Int, scalar []byte) (*big.Int, *big.Int) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	p, err := curve.pointFromAffine(Bx, By)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if err != nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: ScalarMult was called on an invalid point&#34;)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	scalar = curve.normalizeScalar(scalar)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	p, err = p.ScalarMult(p, scalar)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if err != nil {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: nistec rejected normalized scalar&#34;)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	return curve.pointToAffine(p)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) ScalarBaseMult(scalar []byte) (*big.Int, *big.Int) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	scalar = curve.normalizeScalar(scalar)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	p, err := curve.newPoint().ScalarBaseMult(scalar)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	if err != nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: nistec rejected normalized scalar&#34;)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	return curve.pointToAffine(p)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// CombinedMult returns [s1]G + [s2]P where G is the generator. It&#39;s used</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// through an interface upgrade in crypto/ecdsa.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) CombinedMult(Px, Py *big.Int, s1, s2 []byte) (x, y *big.Int) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	s1 = curve.normalizeScalar(s1)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	q, err := curve.newPoint().ScalarBaseMult(s1)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	if err != nil {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: nistec rejected normalized scalar&#34;)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	p, err := curve.pointFromAffine(Px, Py)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if err != nil {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: CombinedMult was called on an invalid point&#34;)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	s2 = curve.normalizeScalar(s2)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	p, err = p.ScalarMult(p, s2)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if err != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: nistec rejected normalized scalar&#34;)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	return curve.pointToAffine(p.Add(p, q))
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) Unmarshal(data []byte) (x, y *big.Int) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if len(data) == 0 || data[0] != 4 {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		return nil, nil
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	<span class="comment">// Use SetBytes to check that data encodes a valid point.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	_, err := curve.newPoint().SetBytes(data)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	if err != nil {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		return nil, nil
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// We don&#39;t use pointToAffine because it involves an expensive field</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// inversion to convert from Jacobian to affine coordinates, which we</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// already have.</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	byteLen := (curve.params.BitSize + 7) / 8
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	y = new(big.Int).SetBytes(data[1+byteLen:])
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	return x, y
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func (curve *nistCurve[Point]) UnmarshalCompressed(data []byte) (x, y *big.Int) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if len(data) == 0 || (data[0] != 2 &amp;&amp; data[0] != 3) {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		return nil, nil
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	p, err := curve.newPoint().SetBytes(data)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	if err != nil {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		return nil, nil
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	return curve.pointToAffine(p)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>func bigFromDecimal(s string) *big.Int {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	b, ok := new(big.Int).SetString(s, 10)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	if !ok {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: internal error: invalid encoding&#34;)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return b
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>func bigFromHex(s string) *big.Int {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	b, ok := new(big.Int).SetString(s, 16)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if !ok {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		panic(&#34;crypto/elliptic: internal error: invalid encoding&#34;)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	return b
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
</pre><p><a href="nistec.go?m=text">View as plain text</a></p>

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
