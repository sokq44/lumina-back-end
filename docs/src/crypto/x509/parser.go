<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/x509/parser.go - Go Documentation Server</title>

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
<a href="parser.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/x509">x509</a>/<span class="text-muted">parser.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/x509">crypto/x509</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package x509
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;crypto/dsa&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto/ecdh&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/ecdsa&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/ed25519&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/elliptic&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;crypto/rsa&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;crypto/x509/pkix&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;encoding/asn1&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;math/big&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;net&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;net/url&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;unicode/utf16&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;golang.org/x/crypto/cryptobyte&#34;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	cryptobyte_asn1 &#34;golang.org/x/crypto/cryptobyte/asn1&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// isPrintable reports whether the given b is in the ASN.1 PrintableString set.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// This is a simplified version of encoding/asn1.isPrintable.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func isPrintable(b byte) bool {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	return &#39;a&#39; &lt;= b &amp;&amp; b &lt;= &#39;z&#39; ||
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		&#39;A&#39; &lt;= b &amp;&amp; b &lt;= &#39;Z&#39; ||
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		&#39;0&#39; &lt;= b &amp;&amp; b &lt;= &#39;9&#39; ||
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		&#39;\&#39;&#39; &lt;= b &amp;&amp; b &lt;= &#39;)&#39; ||
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		&#39;+&#39; &lt;= b &amp;&amp; b &lt;= &#39;/&#39; ||
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		b == &#39; &#39; ||
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		b == &#39;:&#39; ||
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		b == &#39;=&#39; ||
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		b == &#39;?&#39; ||
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		<span class="comment">// This is technically not allowed in a PrintableString.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// However, x509 certificates with wildcard strings don&#39;t</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// always use the correct string type so we permit it.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		b == &#39;*&#39; ||
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		<span class="comment">// This is not technically allowed either. However, not</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		<span class="comment">// only is it relatively common, but there are also a</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		<span class="comment">// handful of CA certificates that contain it. At least</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		<span class="comment">// one of which will not expire until 2027.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		b == &#39;&amp;&#39;
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// parseASN1String parses the ASN.1 string types T61String, PrintableString,</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// UTF8String, BMPString, IA5String, and NumericString. This is mostly copied</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// from the respective encoding/asn1.parse... methods, rather than just</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// increasing the API surface of that package.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func parseASN1String(tag cryptobyte_asn1.Tag, value []byte) (string, error) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	switch tag {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	case cryptobyte_asn1.T61String:
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return string(value), nil
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	case cryptobyte_asn1.PrintableString:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		for _, b := range value {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			if !isPrintable(b) {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				return &#34;&#34;, errors.New(&#34;invalid PrintableString&#34;)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return string(value), nil
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	case cryptobyte_asn1.UTF8String:
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if !utf8.Valid(value) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return &#34;&#34;, errors.New(&#34;invalid UTF-8 string&#34;)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		return string(value), nil
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	case cryptobyte_asn1.Tag(asn1.TagBMPString):
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		if len(value)%2 != 0 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			return &#34;&#34;, errors.New(&#34;invalid BMPString&#34;)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// Strip terminator if present.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		if l := len(value); l &gt;= 2 &amp;&amp; value[l-1] == 0 &amp;&amp; value[l-2] == 0 {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			value = value[:l-2]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		s := make([]uint16, 0, len(value)/2)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		for len(value) &gt; 0 {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			s = append(s, uint16(value[0])&lt;&lt;8+uint16(value[1]))
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			value = value[2:]
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return string(utf16.Decode(s)), nil
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	case cryptobyte_asn1.IA5String:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		s := string(value)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		if isIA5String(s) != nil {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			return &#34;&#34;, errors.New(&#34;invalid IA5String&#34;)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		return s, nil
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	case cryptobyte_asn1.Tag(asn1.TagNumericString):
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		for _, b := range value {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			if !(&#39;0&#39; &lt;= b &amp;&amp; b &lt;= &#39;9&#39; || b == &#39; &#39;) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>				return &#34;&#34;, errors.New(&#34;invalid NumericString&#34;)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		return string(value), nil
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return &#34;&#34;, fmt.Errorf(&#34;unsupported string type: %v&#34;, tag)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// parseName parses a DER encoded Name as defined in RFC 5280. We may</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// want to export this function in the future for use in crypto/tls.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func parseName(raw cryptobyte.String) (*pkix.RDNSequence, error) {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if !raw.ReadASN1(&amp;raw, cryptobyte_asn1.SEQUENCE) {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: invalid RDNSequence&#34;)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	var rdnSeq pkix.RDNSequence
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	for !raw.Empty() {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		var rdnSet pkix.RelativeDistinguishedNameSET
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		var set cryptobyte.String
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		if !raw.ReadASN1(&amp;set, cryptobyte_asn1.SET) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid RDNSequence&#34;)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		for !set.Empty() {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			var atav cryptobyte.String
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			if !set.ReadASN1(&amp;atav, cryptobyte_asn1.SEQUENCE) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: invalid RDNSequence: invalid attribute&#34;)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			var attr pkix.AttributeTypeAndValue
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			if !atav.ReadASN1ObjectIdentifier(&amp;attr.Type) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: invalid RDNSequence: invalid attribute type&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			var rawValue cryptobyte.String
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			var valueTag cryptobyte_asn1.Tag
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			if !atav.ReadAnyASN1(&amp;rawValue, &amp;valueTag) {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: invalid RDNSequence: invalid attribute value&#34;)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			var err error
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			attr.Value, err = parseASN1String(valueTag, rawValue)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if err != nil {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				return nil, fmt.Errorf(&#34;x509: invalid RDNSequence: invalid attribute value: %s&#34;, err)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			rdnSet = append(rdnSet, attr)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		rdnSeq = append(rdnSeq, rdnSet)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	return &amp;rdnSeq, nil
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>func parseAI(der cryptobyte.String) (pkix.AlgorithmIdentifier, error) {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	ai := pkix.AlgorithmIdentifier{}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if !der.ReadASN1ObjectIdentifier(&amp;ai.Algorithm) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return ai, errors.New(&#34;x509: malformed OID&#34;)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if der.Empty() {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return ai, nil
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	var params cryptobyte.String
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	var tag cryptobyte_asn1.Tag
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if !der.ReadAnyASN1Element(&amp;params, &amp;tag) {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		return ai, errors.New(&#34;x509: malformed parameters&#34;)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	ai.Parameters.Tag = int(tag)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	ai.Parameters.FullBytes = params
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	return ai, nil
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>func parseTime(der *cryptobyte.String) (time.Time, error) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	var t time.Time
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	switch {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	case der.PeekASN1Tag(cryptobyte_asn1.UTCTime):
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		if !der.ReadASN1UTCTime(&amp;t) {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			return t, errors.New(&#34;x509: malformed UTCTime&#34;)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	case der.PeekASN1Tag(cryptobyte_asn1.GeneralizedTime):
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if !der.ReadASN1GeneralizedTime(&amp;t) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			return t, errors.New(&#34;x509: malformed GeneralizedTime&#34;)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	default:
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		return t, errors.New(&#34;x509: unsupported time format&#34;)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	return t, nil
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>func parseValidity(der cryptobyte.String) (time.Time, time.Time, error) {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	notBefore, err := parseTime(&amp;der)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	if err != nil {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		return time.Time{}, time.Time{}, err
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	notAfter, err := parseTime(&amp;der)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if err != nil {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		return time.Time{}, time.Time{}, err
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return notBefore, notAfter, nil
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>func parseExtension(der cryptobyte.String) (pkix.Extension, error) {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	var ext pkix.Extension
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if !der.ReadASN1ObjectIdentifier(&amp;ext.Id) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return ext, errors.New(&#34;x509: malformed extension OID field&#34;)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if der.PeekASN1Tag(cryptobyte_asn1.BOOLEAN) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		if !der.ReadASN1Boolean(&amp;ext.Critical) {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			return ext, errors.New(&#34;x509: malformed extension critical field&#34;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	var val cryptobyte.String
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	if !der.ReadASN1(&amp;val, cryptobyte_asn1.OCTET_STRING) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		return ext, errors.New(&#34;x509: malformed extension value field&#34;)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	ext.Value = val
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	return ext, nil
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func parsePublicKey(keyData *publicKeyInfo) (any, error) {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	oid := keyData.Algorithm.Algorithm
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	params := keyData.Algorithm.Parameters
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	der := cryptobyte.String(keyData.PublicKey.RightAlign())
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	switch {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	case oid.Equal(oidPublicKeyRSA):
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// RSA public keys must have a NULL in the parameters.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">// See RFC 3279, Section 2.3.1.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		if !bytes.Equal(params.FullBytes, asn1.NullBytes) {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: RSA key missing NULL parameters&#34;)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		p := &amp;pkcs1PublicKey{N: new(big.Int)}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		if !der.ReadASN1(&amp;der, cryptobyte_asn1.SEQUENCE) {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid RSA public key&#34;)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		if !der.ReadASN1Integer(p.N) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid RSA modulus&#34;)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		if !der.ReadASN1Integer(&amp;p.E) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid RSA public exponent&#34;)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		if p.N.Sign() &lt;= 0 {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: RSA modulus is not a positive number&#34;)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		if p.E &lt;= 0 {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: RSA public exponent is not a positive number&#34;)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		pub := &amp;rsa.PublicKey{
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			E: p.E,
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			N: p.N,
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		return pub, nil
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	case oid.Equal(oidPublicKeyECDSA):
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		paramsDer := cryptobyte.String(params.FullBytes)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		namedCurveOID := new(asn1.ObjectIdentifier)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		if !paramsDer.ReadASN1ObjectIdentifier(namedCurveOID) {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid ECDSA parameters&#34;)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		namedCurve := namedCurveFromOID(*namedCurveOID)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		if namedCurve == nil {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: unsupported elliptic curve&#34;)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		x, y := elliptic.Unmarshal(namedCurve, der)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		if x == nil {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: failed to unmarshal elliptic curve point&#34;)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		pub := &amp;ecdsa.PublicKey{
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			Curve: namedCurve,
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			X:     x,
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			Y:     y,
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		return pub, nil
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	case oid.Equal(oidPublicKeyEd25519):
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		<span class="comment">// RFC 8410, Section 3</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		<span class="comment">// &gt; For all of the OIDs, the parameters MUST be absent.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		if len(params.FullBytes) != 0 {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: Ed25519 key encoded with illegal parameters&#34;)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		if len(der) != ed25519.PublicKeySize {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: wrong Ed25519 public key size&#34;)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return ed25519.PublicKey(der), nil
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	case oid.Equal(oidPublicKeyX25519):
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// RFC 8410, Section 3</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">// &gt; For all of the OIDs, the parameters MUST be absent.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		if len(params.FullBytes) != 0 {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: X25519 key encoded with illegal parameters&#34;)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		return ecdh.X25519().NewPublicKey(der)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	case oid.Equal(oidPublicKeyDSA):
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		y := new(big.Int)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		if !der.ReadASN1Integer(y) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid DSA public key&#34;)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		pub := &amp;dsa.PublicKey{
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			Y: y,
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			Parameters: dsa.Parameters{
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>				P: new(big.Int),
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				Q: new(big.Int),
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				G: new(big.Int),
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			},
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		paramsDer := cryptobyte.String(params.FullBytes)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		if !paramsDer.ReadASN1(&amp;paramsDer, cryptobyte_asn1.SEQUENCE) ||
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			!paramsDer.ReadASN1Integer(pub.Parameters.P) ||
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			!paramsDer.ReadASN1Integer(pub.Parameters.Q) ||
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			!paramsDer.ReadASN1Integer(pub.Parameters.G) {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid DSA parameters&#34;)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		if pub.Y.Sign() &lt;= 0 || pub.Parameters.P.Sign() &lt;= 0 ||
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			pub.Parameters.Q.Sign() &lt;= 0 || pub.Parameters.G.Sign() &lt;= 0 {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: zero or negative DSA parameter&#34;)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		return pub, nil
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	default:
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: unknown public key algorithm&#34;)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>func parseKeyUsageExtension(der cryptobyte.String) (KeyUsage, error) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	var usageBits asn1.BitString
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	if !der.ReadASN1BitString(&amp;usageBits) {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		return 0, errors.New(&#34;x509: invalid key usage&#34;)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	var usage int
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	for i := 0; i &lt; 9; i++ {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		if usageBits.At(i) != 0 {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			usage |= 1 &lt;&lt; uint(i)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	return KeyUsage(usage), nil
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>func parseBasicConstraintsExtension(der cryptobyte.String) (bool, int, error) {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	var isCA bool
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	if !der.ReadASN1(&amp;der, cryptobyte_asn1.SEQUENCE) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		return false, 0, errors.New(&#34;x509: invalid basic constraints&#34;)
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if der.PeekASN1Tag(cryptobyte_asn1.BOOLEAN) {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		if !der.ReadASN1Boolean(&amp;isCA) {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			return false, 0, errors.New(&#34;x509: invalid basic constraints&#34;)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	maxPathLen := -1
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if der.PeekASN1Tag(cryptobyte_asn1.INTEGER) {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		if !der.ReadASN1Integer(&amp;maxPathLen) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			return false, 0, errors.New(&#34;x509: invalid basic constraints&#34;)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// TODO: map out.MaxPathLen to 0 if it has the -1 default value? (Issue 19285)</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	return isCA, maxPathLen, nil
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>func forEachSAN(der cryptobyte.String, callback func(tag int, data []byte) error) error {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	if !der.ReadASN1(&amp;der, cryptobyte_asn1.SEQUENCE) {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		return errors.New(&#34;x509: invalid subject alternative names&#34;)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	for !der.Empty() {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		var san cryptobyte.String
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		var tag cryptobyte_asn1.Tag
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		if !der.ReadAnyASN1(&amp;san, &amp;tag) {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			return errors.New(&#34;x509: invalid subject alternative name&#34;)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		if err := callback(int(tag^0x80), san); err != nil {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			return err
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	return nil
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func parseSANExtension(der cryptobyte.String) (dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL, err error) {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	err = forEachSAN(der, func(tag int, data []byte) error {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		switch tag {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		case nameTypeEmail:
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			email := string(data)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			if err := isIA5String(email); err != nil {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>				return errors.New(&#34;x509: SAN rfc822Name is malformed&#34;)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			emailAddresses = append(emailAddresses, email)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		case nameTypeDNS:
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			name := string(data)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			if err := isIA5String(name); err != nil {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				return errors.New(&#34;x509: SAN dNSName is malformed&#34;)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			dnsNames = append(dnsNames, string(name))
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		case nameTypeURI:
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			uriStr := string(data)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			if err := isIA5String(uriStr); err != nil {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				return errors.New(&#34;x509: SAN uniformResourceIdentifier is malformed&#34;)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			uri, err := url.Parse(uriStr)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			if err != nil {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;x509: cannot parse URI %q: %s&#34;, uriStr, err)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			if len(uri.Host) &gt; 0 {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				if _, ok := domainToReverseLabels(uri.Host); !ok {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>					return fmt.Errorf(&#34;x509: cannot parse URI %q: invalid domain&#34;, uriStr)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>				}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			uris = append(uris, uri)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		case nameTypeIP:
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			switch len(data) {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			case net.IPv4len, net.IPv6len:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>				ipAddresses = append(ipAddresses, data)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			default:
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>				return errors.New(&#34;x509: cannot parse IP address of length &#34; + strconv.Itoa(len(data)))
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		return nil
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	})
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	return
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>func parseExtKeyUsageExtension(der cryptobyte.String) ([]ExtKeyUsage, []asn1.ObjectIdentifier, error) {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	var extKeyUsages []ExtKeyUsage
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	var unknownUsages []asn1.ObjectIdentifier
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if !der.ReadASN1(&amp;der, cryptobyte_asn1.SEQUENCE) {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		return nil, nil, errors.New(&#34;x509: invalid extended key usages&#34;)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	for !der.Empty() {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		var eku asn1.ObjectIdentifier
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		if !der.ReadASN1ObjectIdentifier(&amp;eku) {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			return nil, nil, errors.New(&#34;x509: invalid extended key usages&#34;)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		if extKeyUsage, ok := extKeyUsageFromOID(eku); ok {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			extKeyUsages = append(extKeyUsages, extKeyUsage)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		} else {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			unknownUsages = append(unknownUsages, eku)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	return extKeyUsages, unknownUsages, nil
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>func parseCertificatePoliciesExtension(der cryptobyte.String) ([]OID, error) {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	var oids []OID
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	if !der.ReadASN1(&amp;der, cryptobyte_asn1.SEQUENCE) {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: invalid certificate policies&#34;)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	for !der.Empty() {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		var cp cryptobyte.String
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		var OIDBytes cryptobyte.String
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		if !der.ReadASN1(&amp;cp, cryptobyte_asn1.SEQUENCE) || !cp.ReadASN1(&amp;OIDBytes, cryptobyte_asn1.OBJECT_IDENTIFIER) {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid certificate policies&#34;)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		oid, ok := newOIDFromDER(OIDBytes)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		if !ok {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: invalid certificate policies&#34;)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		oids = append(oids, oid)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	return oids, nil
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span><span class="comment">// isValidIPMask reports whether mask consists of zero or more 1 bits, followed by zero bits.</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>func isValidIPMask(mask []byte) bool {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	seenZero := false
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	for _, b := range mask {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		if seenZero {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			if b != 0 {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				return false
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>			}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			continue
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		switch b {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		case 0x00, 0x80, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfe:
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			seenZero = true
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		case 0xff:
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		default:
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			return false
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	return true
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>func parseNameConstraintsExtension(out *Certificate, e pkix.Extension) (unhandled bool, err error) {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	<span class="comment">// RFC 5280, 4.2.1.10</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	<span class="comment">// NameConstraints ::= SEQUENCE {</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	<span class="comment">//      permittedSubtrees       [0]     GeneralSubtrees OPTIONAL,</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">//      excludedSubtrees        [1]     GeneralSubtrees OPTIONAL }</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	<span class="comment">// GeneralSubtrees ::= SEQUENCE SIZE (1..MAX) OF GeneralSubtree</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// GeneralSubtree ::= SEQUENCE {</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	<span class="comment">//      base                    GeneralName,</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	<span class="comment">//      minimum         [0]     BaseDistance DEFAULT 0,</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	<span class="comment">//      maximum         [1]     BaseDistance OPTIONAL }</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	<span class="comment">// BaseDistance ::= INTEGER (0..MAX)</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	outer := cryptobyte.String(e.Value)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	var toplevel, permitted, excluded cryptobyte.String
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	var havePermitted, haveExcluded bool
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	if !outer.ReadASN1(&amp;toplevel, cryptobyte_asn1.SEQUENCE) ||
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		!outer.Empty() ||
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		!toplevel.ReadOptionalASN1(&amp;permitted, &amp;havePermitted, cryptobyte_asn1.Tag(0).ContextSpecific().Constructed()) ||
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		!toplevel.ReadOptionalASN1(&amp;excluded, &amp;haveExcluded, cryptobyte_asn1.Tag(1).ContextSpecific().Constructed()) ||
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		!toplevel.Empty() {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		return false, errors.New(&#34;x509: invalid NameConstraints extension&#34;)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if !havePermitted &amp;&amp; !haveExcluded || len(permitted) == 0 &amp;&amp; len(excluded) == 0 {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		<span class="comment">// From RFC 5280, Section 4.2.1.10:</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		<span class="comment">//   “either the permittedSubtrees field</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		<span class="comment">//   or the excludedSubtrees MUST be</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		<span class="comment">//   present”</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		return false, errors.New(&#34;x509: empty name constraints extension&#34;)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	getValues := func(subtrees cryptobyte.String) (dnsNames []string, ips []*net.IPNet, emails, uriDomains []string, err error) {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		for !subtrees.Empty() {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			var seq, value cryptobyte.String
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			var tag cryptobyte_asn1.Tag
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			if !subtrees.ReadASN1(&amp;seq, cryptobyte_asn1.SEQUENCE) ||
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>				!seq.ReadAnyASN1(&amp;value, &amp;tag) {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				return nil, nil, nil, nil, fmt.Errorf(&#34;x509: invalid NameConstraints extension&#34;)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			var (
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				dnsTag   = cryptobyte_asn1.Tag(2).ContextSpecific()
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				emailTag = cryptobyte_asn1.Tag(1).ContextSpecific()
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				ipTag    = cryptobyte_asn1.Tag(7).ContextSpecific()
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				uriTag   = cryptobyte_asn1.Tag(6).ContextSpecific()
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			switch tag {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			case dnsTag:
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>				domain := string(value)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>				if err := isIA5String(domain); err != nil {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>					return nil, nil, nil, nil, errors.New(&#34;x509: invalid constraint value: &#34; + err.Error())
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>				}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>				trimmedDomain := domain
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>				if len(trimmedDomain) &gt; 0 &amp;&amp; trimmedDomain[0] == &#39;.&#39; {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>					<span class="comment">// constraints can have a leading</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>					<span class="comment">// period to exclude the domain</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>					<span class="comment">// itself, but that&#39;s not valid in a</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>					<span class="comment">// normal domain name.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>					trimmedDomain = trimmedDomain[1:]
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>				}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>				if _, ok := domainToReverseLabels(trimmedDomain); !ok {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>					return nil, nil, nil, nil, fmt.Errorf(&#34;x509: failed to parse dnsName constraint %q&#34;, domain)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>				}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>				dnsNames = append(dnsNames, domain)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			case ipTag:
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				l := len(value)
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				var ip, mask []byte
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>				switch l {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>				case 8:
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>					ip = value[:4]
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>					mask = value[4:]
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>				case 32:
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>					ip = value[:16]
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>					mask = value[16:]
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>				default:
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>					return nil, nil, nil, nil, fmt.Errorf(&#34;x509: IP constraint contained value of length %d&#34;, l)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>				if !isValidIPMask(mask) {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>					return nil, nil, nil, nil, fmt.Errorf(&#34;x509: IP constraint contained invalid mask %x&#34;, mask)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>				}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>				ips = append(ips, &amp;net.IPNet{IP: net.IP(ip), Mask: net.IPMask(mask)})
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			case emailTag:
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>				constraint := string(value)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>				if err := isIA5String(constraint); err != nil {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>					return nil, nil, nil, nil, errors.New(&#34;x509: invalid constraint value: &#34; + err.Error())
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>				}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>				<span class="comment">// If the constraint contains an @ then</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>				<span class="comment">// it specifies an exact mailbox name.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>				if strings.Contains(constraint, &#34;@&#34;) {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>					if _, ok := parseRFC2821Mailbox(constraint); !ok {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>						return nil, nil, nil, nil, fmt.Errorf(&#34;x509: failed to parse rfc822Name constraint %q&#34;, constraint)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>					}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>				} else {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>					<span class="comment">// Otherwise it&#39;s a domain name.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>					domain := constraint
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>					if len(domain) &gt; 0 &amp;&amp; domain[0] == &#39;.&#39; {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>						domain = domain[1:]
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>					}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>					if _, ok := domainToReverseLabels(domain); !ok {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>						return nil, nil, nil, nil, fmt.Errorf(&#34;x509: failed to parse rfc822Name constraint %q&#34;, constraint)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>					}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>				}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				emails = append(emails, constraint)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			case uriTag:
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>				domain := string(value)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>				if err := isIA5String(domain); err != nil {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>					return nil, nil, nil, nil, errors.New(&#34;x509: invalid constraint value: &#34; + err.Error())
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>				}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>				if net.ParseIP(domain) != nil {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>					return nil, nil, nil, nil, fmt.Errorf(&#34;x509: failed to parse URI constraint %q: cannot be IP address&#34;, domain)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>				}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>				trimmedDomain := domain
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>				if len(trimmedDomain) &gt; 0 &amp;&amp; trimmedDomain[0] == &#39;.&#39; {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>					<span class="comment">// constraints can have a leading</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>					<span class="comment">// period to exclude the domain itself,</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>					<span class="comment">// but that&#39;s not valid in a normal</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>					<span class="comment">// domain name.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>					trimmedDomain = trimmedDomain[1:]
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>				}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>				if _, ok := domainToReverseLabels(trimmedDomain); !ok {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>					return nil, nil, nil, nil, fmt.Errorf(&#34;x509: failed to parse URI constraint %q&#34;, domain)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>				}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				uriDomains = append(uriDomains, domain)
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			default:
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>				unhandled = true
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>			}
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		return dnsNames, ips, emails, uriDomains, nil
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	if out.PermittedDNSDomains, out.PermittedIPRanges, out.PermittedEmailAddresses, out.PermittedURIDomains, err = getValues(permitted); err != nil {
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		return false, err
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	if out.ExcludedDNSDomains, out.ExcludedIPRanges, out.ExcludedEmailAddresses, out.ExcludedURIDomains, err = getValues(excluded); err != nil {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		return false, err
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	out.PermittedDNSDomainsCritical = e.Critical
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	return unhandled, nil
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>func processExtensions(out *Certificate) error {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	var err error
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	for _, e := range out.Extensions {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		unhandled := false
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		if len(e.Id) == 4 &amp;&amp; e.Id[0] == 2 &amp;&amp; e.Id[1] == 5 &amp;&amp; e.Id[2] == 29 {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			switch e.Id[3] {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>			case 15:
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>				out.KeyUsage, err = parseKeyUsageExtension(e.Value)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>				if err != nil {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>					return err
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>				}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>			case 19:
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>				out.IsCA, out.MaxPathLen, err = parseBasicConstraintsExtension(e.Value)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>				if err != nil {
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>					return err
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>				}
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>				out.BasicConstraintsValid = true
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>				out.MaxPathLenZero = out.MaxPathLen == 0
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			case 17:
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>				out.DNSNames, out.EmailAddresses, out.IPAddresses, out.URIs, err = parseSANExtension(e.Value)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>				if err != nil {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>					return err
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>				}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>				if len(out.DNSNames) == 0 &amp;&amp; len(out.EmailAddresses) == 0 &amp;&amp; len(out.IPAddresses) == 0 &amp;&amp; len(out.URIs) == 0 {
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>					<span class="comment">// If we didn&#39;t parse anything then we do the critical check, below.</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>					unhandled = true
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>				}
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>			case 30:
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>				unhandled, err = parseNameConstraintsExtension(out, e)
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>				if err != nil {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>					return err
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>			case 31:
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>				<span class="comment">// RFC 5280, 4.2.1.13</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>				<span class="comment">// CRLDistributionPoints ::= SEQUENCE SIZE (1..MAX) OF DistributionPoint</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>				<span class="comment">// DistributionPoint ::= SEQUENCE {</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>				<span class="comment">//     distributionPoint       [0]     DistributionPointName OPTIONAL,</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>				<span class="comment">//     reasons                 [1]     ReasonFlags OPTIONAL,</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				<span class="comment">//     cRLIssuer               [2]     GeneralNames OPTIONAL }</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>				<span class="comment">// DistributionPointName ::= CHOICE {</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>				<span class="comment">//     fullName                [0]     GeneralNames,</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>				<span class="comment">//     nameRelativeToCRLIssuer [1]     RelativeDistinguishedName }</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>				val := cryptobyte.String(e.Value)
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>				if !val.ReadASN1(&amp;val, cryptobyte_asn1.SEQUENCE) {
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid CRL distribution points&#34;)
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>				}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>				for !val.Empty() {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>					var dpDER cryptobyte.String
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>					if !val.ReadASN1(&amp;dpDER, cryptobyte_asn1.SEQUENCE) {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>						return errors.New(&#34;x509: invalid CRL distribution point&#34;)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>					}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>					var dpNameDER cryptobyte.String
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>					var dpNamePresent bool
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>					if !dpDER.ReadOptionalASN1(&amp;dpNameDER, &amp;dpNamePresent, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>						return errors.New(&#34;x509: invalid CRL distribution point&#34;)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>					}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>					if !dpNamePresent {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>						continue
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>					}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>					if !dpNameDER.ReadASN1(&amp;dpNameDER, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>						return errors.New(&#34;x509: invalid CRL distribution point&#34;)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>					}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>					for !dpNameDER.Empty() {
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>						if !dpNameDER.PeekASN1Tag(cryptobyte_asn1.Tag(6).ContextSpecific()) {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>							break
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>						}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>						var uri cryptobyte.String
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>						if !dpNameDER.ReadASN1(&amp;uri, cryptobyte_asn1.Tag(6).ContextSpecific()) {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>							return errors.New(&#34;x509: invalid CRL distribution point&#34;)
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>						}
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>						out.CRLDistributionPoints = append(out.CRLDistributionPoints, string(uri))
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>					}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>				}
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			case 35:
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>				<span class="comment">// RFC 5280, 4.2.1.1</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>				val := cryptobyte.String(e.Value)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>				var akid cryptobyte.String
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>				if !val.ReadASN1(&amp;akid, cryptobyte_asn1.SEQUENCE) {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid authority key identifier&#34;)
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>				}
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>				if akid.PeekASN1Tag(cryptobyte_asn1.Tag(0).ContextSpecific()) {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>					if !akid.ReadASN1(&amp;akid, cryptobyte_asn1.Tag(0).ContextSpecific()) {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>						return errors.New(&#34;x509: invalid authority key identifier&#34;)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>					}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>					out.AuthorityKeyId = akid
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>				}
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			case 37:
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>				out.ExtKeyUsage, out.UnknownExtKeyUsage, err = parseExtKeyUsageExtension(e.Value)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>				if err != nil {
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>					return err
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>				}
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			case 14:
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>				<span class="comment">// RFC 5280, 4.2.1.2</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>				val := cryptobyte.String(e.Value)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>				var skid cryptobyte.String
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>				if !val.ReadASN1(&amp;skid, cryptobyte_asn1.OCTET_STRING) {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid subject key identifier&#34;)
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>				}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>				out.SubjectKeyId = skid
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			case 32:
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>				out.Policies, err = parseCertificatePoliciesExtension(e.Value)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>				if err != nil {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>					return err
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>				}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>				out.PolicyIdentifiers = make([]asn1.ObjectIdentifier, 0, len(out.Policies))
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>				for _, oid := range out.Policies {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>					if oid, ok := oid.toASN1OID(); ok {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>						out.PolicyIdentifiers = append(out.PolicyIdentifiers, oid)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>					}
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>				}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>			default:
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>				<span class="comment">// Unknown extensions are recorded if critical.</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>				unhandled = true
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		} else if e.Id.Equal(oidExtensionAuthorityInfoAccess) {
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			<span class="comment">// RFC 5280 4.2.2.1: Authority Information Access</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			val := cryptobyte.String(e.Value)
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			if !val.ReadASN1(&amp;val, cryptobyte_asn1.SEQUENCE) {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>				return errors.New(&#34;x509: invalid authority info access&#34;)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			for !val.Empty() {
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>				var aiaDER cryptobyte.String
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				if !val.ReadASN1(&amp;aiaDER, cryptobyte_asn1.SEQUENCE) {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid authority info access&#34;)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>				}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>				var method asn1.ObjectIdentifier
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>				if !aiaDER.ReadASN1ObjectIdentifier(&amp;method) {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid authority info access&#34;)
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				}
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>				if !aiaDER.PeekASN1Tag(cryptobyte_asn1.Tag(6).ContextSpecific()) {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>					continue
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>				if !aiaDER.ReadASN1(&amp;aiaDER, cryptobyte_asn1.Tag(6).ContextSpecific()) {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>					return errors.New(&#34;x509: invalid authority info access&#34;)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>				}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>				switch {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>				case method.Equal(oidAuthorityInfoAccessOcsp):
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>					out.OCSPServer = append(out.OCSPServer, string(aiaDER))
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				case method.Equal(oidAuthorityInfoAccessIssuers):
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>					out.IssuingCertificateURL = append(out.IssuingCertificateURL, string(aiaDER))
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>				}
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>			}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		} else {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			<span class="comment">// Unknown extensions are recorded if critical.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>			unhandled = true
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		}
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		if e.Critical &amp;&amp; unhandled {
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			out.UnhandledCriticalExtensions = append(out.UnhandledCriticalExtensions, e.Id)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		}
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	}
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	return nil
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>}
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>func parseCertificate(der []byte) (*Certificate, error) {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	cert := &amp;Certificate{}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	input := cryptobyte.String(der)
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	<span class="comment">// we read the SEQUENCE including length and tag bytes so that</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	<span class="comment">// we can populate Certificate.Raw, before unwrapping the</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	<span class="comment">// SEQUENCE so it can be operated on</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	if !input.ReadASN1Element(&amp;input, cryptobyte_asn1.SEQUENCE) {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed certificate&#34;)
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>	}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	cert.Raw = input
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	if !input.ReadASN1(&amp;input, cryptobyte_asn1.SEQUENCE) {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed certificate&#34;)
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	}
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	var tbs cryptobyte.String
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// do the same trick again as above to extract the raw</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	<span class="comment">// bytes for Certificate.RawTBSCertificate</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	if !input.ReadASN1Element(&amp;tbs, cryptobyte_asn1.SEQUENCE) {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed tbs certificate&#34;)
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	cert.RawTBSCertificate = tbs
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	if !tbs.ReadASN1(&amp;tbs, cryptobyte_asn1.SEQUENCE) {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed tbs certificate&#34;)
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	}
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	if !tbs.ReadOptionalASN1Integer(&amp;cert.Version, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific(), 0) {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed version&#34;)
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	}
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	if cert.Version &lt; 0 {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed version&#34;)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	<span class="comment">// for backwards compat reasons Version is one-indexed,</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	<span class="comment">// rather than zero-indexed as defined in 5280</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	cert.Version++
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	if cert.Version &gt; 3 {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: invalid version&#34;)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	serial := new(big.Int)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	if !tbs.ReadASN1Integer(serial) {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed serial number&#34;)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	}
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>	<span class="comment">// we ignore the presence of negative serial numbers because</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	<span class="comment">// of their prevalence, despite them being invalid</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	<span class="comment">// TODO(rolandshoemaker): revisit this decision, there are currently</span>
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	<span class="comment">// only 10 trusted certificates with negative serial numbers</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	<span class="comment">// according to censys.io.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	cert.SerialNumber = serial
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	var sigAISeq cryptobyte.String
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	if !tbs.ReadASN1(&amp;sigAISeq, cryptobyte_asn1.SEQUENCE) {
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed signature algorithm identifier&#34;)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	<span class="comment">// Before parsing the inner algorithm identifier, extract</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	<span class="comment">// the outer algorithm identifier and make sure that they</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	<span class="comment">// match.</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	var outerSigAISeq cryptobyte.String
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	if !input.ReadASN1(&amp;outerSigAISeq, cryptobyte_asn1.SEQUENCE) {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed algorithm identifier&#34;)
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	}
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	if !bytes.Equal(outerSigAISeq, sigAISeq) {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: inner and outer signature algorithm identifiers don&#39;t match&#34;)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	}
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	sigAI, err := parseAI(sigAISeq)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	if err != nil {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		return nil, err
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	cert.SignatureAlgorithm = getSignatureAlgorithmFromAI(sigAI)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	var issuerSeq cryptobyte.String
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	if !tbs.ReadASN1Element(&amp;issuerSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed issuer&#34;)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	}
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	cert.RawIssuer = issuerSeq
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	issuerRDNs, err := parseName(issuerSeq)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	if err != nil {
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		return nil, err
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	cert.Issuer.FillFromRDNSequence(issuerRDNs)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	var validity cryptobyte.String
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	if !tbs.ReadASN1(&amp;validity, cryptobyte_asn1.SEQUENCE) {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed validity&#34;)
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	}
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	cert.NotBefore, cert.NotAfter, err = parseValidity(validity)
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	if err != nil {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		return nil, err
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	}
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	var subjectSeq cryptobyte.String
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>	if !tbs.ReadASN1Element(&amp;subjectSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed issuer&#34;)
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	}
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	cert.RawSubject = subjectSeq
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	subjectRDNs, err := parseName(subjectSeq)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	if err != nil {
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		return nil, err
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	cert.Subject.FillFromRDNSequence(subjectRDNs)
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	var spki cryptobyte.String
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	if !tbs.ReadASN1Element(&amp;spki, cryptobyte_asn1.SEQUENCE) {
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed spki&#34;)
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	}
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	cert.RawSubjectPublicKeyInfo = spki
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	if !spki.ReadASN1(&amp;spki, cryptobyte_asn1.SEQUENCE) {
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed spki&#34;)
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	}
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	var pkAISeq cryptobyte.String
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	if !spki.ReadASN1(&amp;pkAISeq, cryptobyte_asn1.SEQUENCE) {
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed public key algorithm identifier&#34;)
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	}
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	pkAI, err := parseAI(pkAISeq)
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	if err != nil {
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		return nil, err
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	}
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	cert.PublicKeyAlgorithm = getPublicKeyAlgorithmFromOID(pkAI.Algorithm)
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	var spk asn1.BitString
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	if !spki.ReadASN1BitString(&amp;spk) {
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed subjectPublicKey&#34;)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	}
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	if cert.PublicKeyAlgorithm != UnknownPublicKeyAlgorithm {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		cert.PublicKey, err = parsePublicKey(&amp;publicKeyInfo{
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			Algorithm: pkAI,
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>			PublicKey: spk,
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		})
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		if err != nil {
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>			return nil, err
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>	}
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	if cert.Version &gt; 1 {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		if !tbs.SkipOptionalASN1(cryptobyte_asn1.Tag(1).ContextSpecific()) {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: malformed issuerUniqueID&#34;)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		}
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		if !tbs.SkipOptionalASN1(cryptobyte_asn1.Tag(2).ContextSpecific()) {
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: malformed subjectUniqueID&#34;)
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		}
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		if cert.Version == 3 {
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>			var extensions cryptobyte.String
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			var present bool
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>			if !tbs.ReadOptionalASN1(&amp;extensions, &amp;present, cryptobyte_asn1.Tag(3).Constructed().ContextSpecific()) {
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed extensions&#34;)
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			}
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			if present {
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>				seenExts := make(map[string]bool)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>				if !extensions.ReadASN1(&amp;extensions, cryptobyte_asn1.SEQUENCE) {
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>					return nil, errors.New(&#34;x509: malformed extensions&#34;)
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>				for !extensions.Empty() {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>					var extension cryptobyte.String
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>					if !extensions.ReadASN1(&amp;extension, cryptobyte_asn1.SEQUENCE) {
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>						return nil, errors.New(&#34;x509: malformed extension&#34;)
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>					}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>					ext, err := parseExtension(extension)
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>					if err != nil {
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>						return nil, err
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>					}
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>					oidStr := ext.Id.String()
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>					if seenExts[oidStr] {
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>						return nil, errors.New(&#34;x509: certificate contains duplicate extensions&#34;)
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>					}
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>					seenExts[oidStr] = true
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>					cert.Extensions = append(cert.Extensions, ext)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>				}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>				err = processExtensions(cert)
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>				if err != nil {
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>					return nil, err
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>				}
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>			}
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>		}
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	var signature asn1.BitString
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	if !input.ReadASN1BitString(&amp;signature) {
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed signature&#34;)
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	}
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	cert.Signature = signature.RightAlign()
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	return cert, nil
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>}
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span><span class="comment">// ParseCertificate parses a single certificate from the given ASN.1 DER data.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>func ParseCertificate(der []byte) (*Certificate, error) {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	cert, err := parseCertificate(der)
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	if err != nil {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		return nil, err
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	}
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>	if len(der) != len(cert.Raw) {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: trailing data&#34;)
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>	}
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	return cert, err
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>}
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span><span class="comment">// ParseCertificates parses one or more certificates from the given ASN.1 DER</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span><span class="comment">// data. The certificates must be concatenated with no intermediate padding.</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>func ParseCertificates(der []byte) ([]*Certificate, error) {
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	var certs []*Certificate
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	for len(der) &gt; 0 {
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>		cert, err := parseCertificate(der)
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>		if err != nil {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>			return nil, err
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>		}
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		certs = append(certs, cert)
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>		der = der[len(cert.Raw):]
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	}
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	return certs, nil
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>}
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span><span class="comment">// The X.509 standards confusingly 1-indexed the version names, but 0-indexed</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span><span class="comment">// the actual encoded version, so the version for X.509v2 is 1.</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>const x509v2Version = 1
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span><span class="comment">// ParseRevocationList parses a X509 v2 [Certificate] Revocation List from the given</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span><span class="comment">// ASN.1 DER data.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>func ParseRevocationList(der []byte) (*RevocationList, error) {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	rl := &amp;RevocationList{}
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	input := cryptobyte.String(der)
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	<span class="comment">// we read the SEQUENCE including length and tag bytes so that</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	<span class="comment">// we can populate RevocationList.Raw, before unwrapping the</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	<span class="comment">// SEQUENCE so it can be operated on</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	if !input.ReadASN1Element(&amp;input, cryptobyte_asn1.SEQUENCE) {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	}
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	rl.Raw = input
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	if !input.ReadASN1(&amp;input, cryptobyte_asn1.SEQUENCE) {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	var tbs cryptobyte.String
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	<span class="comment">// do the same trick again as above to extract the raw</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	<span class="comment">// bytes for Certificate.RawTBSCertificate</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	if !input.ReadASN1Element(&amp;tbs, cryptobyte_asn1.SEQUENCE) {
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed tbs crl&#34;)
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>	}
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	rl.RawTBSRevocationList = tbs
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>	if !tbs.ReadASN1(&amp;tbs, cryptobyte_asn1.SEQUENCE) {
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed tbs crl&#34;)
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>	}
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>	var version int
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	if !tbs.PeekASN1Tag(cryptobyte_asn1.INTEGER) {
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: unsupported crl version&#34;)
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	}
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>	if !tbs.ReadASN1Integer(&amp;version) {
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	}
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	if version != x509v2Version {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;x509: unsupported crl version: %d&#34;, version)
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	}
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	var sigAISeq cryptobyte.String
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>	if !tbs.ReadASN1(&amp;sigAISeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed signature algorithm identifier&#34;)
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	}
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	<span class="comment">// Before parsing the inner algorithm identifier, extract</span>
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>	<span class="comment">// the outer algorithm identifier and make sure that they</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	<span class="comment">// match.</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	var outerSigAISeq cryptobyte.String
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>	if !input.ReadASN1(&amp;outerSigAISeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed algorithm identifier&#34;)
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	}
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	if !bytes.Equal(outerSigAISeq, sigAISeq) {
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: inner and outer signature algorithm identifiers don&#39;t match&#34;)
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	}
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	sigAI, err := parseAI(sigAISeq)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	if err != nil {
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>		return nil, err
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>	}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	rl.SignatureAlgorithm = getSignatureAlgorithmFromAI(sigAI)
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	var signature asn1.BitString
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>	if !input.ReadASN1BitString(&amp;signature) {
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed signature&#34;)
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>	}
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>	rl.Signature = signature.RightAlign()
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	var issuerSeq cryptobyte.String
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	if !tbs.ReadASN1Element(&amp;issuerSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed issuer&#34;)
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	}
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>	rl.RawIssuer = issuerSeq
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	issuerRDNs, err := parseName(issuerSeq)
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	if err != nil {
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>		return nil, err
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	}
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>	rl.Issuer.FillFromRDNSequence(issuerRDNs)
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>	rl.ThisUpdate, err = parseTime(&amp;tbs)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>	if err != nil {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>		return nil, err
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	if tbs.PeekASN1Tag(cryptobyte_asn1.GeneralizedTime) || tbs.PeekASN1Tag(cryptobyte_asn1.UTCTime) {
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>		rl.NextUpdate, err = parseTime(&amp;tbs)
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>		if err != nil {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>			return nil, err
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		}
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	}
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	if tbs.PeekASN1Tag(cryptobyte_asn1.SEQUENCE) {
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		var revokedSeq cryptobyte.String
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>		if !tbs.ReadASN1(&amp;revokedSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>		}
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>		for !revokedSeq.Empty() {
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			rce := RevocationListEntry{}
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			var certSeq cryptobyte.String
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			if !revokedSeq.ReadASN1Element(&amp;certSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			}
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>			rce.Raw = certSeq
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>			if !certSeq.ReadASN1(&amp;certSeq, cryptobyte_asn1.SEQUENCE) {
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed crl&#34;)
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>			rce.SerialNumber = new(big.Int)
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			if !certSeq.ReadASN1Integer(rce.SerialNumber) {
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed serial number&#34;)
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			}
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>			rce.RevocationTime, err = parseTime(&amp;certSeq)
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>			if err != nil {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>				return nil, err
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>			}
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>			var extensions cryptobyte.String
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>			var present bool
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>			if !certSeq.ReadOptionalASN1(&amp;extensions, &amp;present, cryptobyte_asn1.SEQUENCE) {
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed extensions&#34;)
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>			}
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>			if present {
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>				for !extensions.Empty() {
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>					var extension cryptobyte.String
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>					if !extensions.ReadASN1(&amp;extension, cryptobyte_asn1.SEQUENCE) {
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>						return nil, errors.New(&#34;x509: malformed extension&#34;)
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>					}
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>					ext, err := parseExtension(extension)
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>					if err != nil {
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>						return nil, err
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>					}
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>					if ext.Id.Equal(oidExtensionReasonCode) {
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>						val := cryptobyte.String(ext.Value)
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>						if !val.ReadASN1Enum(&amp;rce.ReasonCode) {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>							return nil, fmt.Errorf(&#34;x509: malformed reasonCode extension&#34;)
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>						}
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>					}
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>					rce.Extensions = append(rce.Extensions, ext)
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>				}
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>			}
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>			rl.RevokedCertificateEntries = append(rl.RevokedCertificateEntries, rce)
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>			rcDeprecated := pkix.RevokedCertificate{
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>				SerialNumber:   rce.SerialNumber,
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>				RevocationTime: rce.RevocationTime,
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>				Extensions:     rce.Extensions,
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>			}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>			rl.RevokedCertificates = append(rl.RevokedCertificates, rcDeprecated)
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		}
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	}
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	var extensions cryptobyte.String
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>	var present bool
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>	if !tbs.ReadOptionalASN1(&amp;extensions, &amp;present, cryptobyte_asn1.Tag(0).Constructed().ContextSpecific()) {
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>		return nil, errors.New(&#34;x509: malformed extensions&#34;)
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>	}
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>	if present {
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>		if !extensions.ReadASN1(&amp;extensions, cryptobyte_asn1.SEQUENCE) {
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>			return nil, errors.New(&#34;x509: malformed extensions&#34;)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>		}
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		for !extensions.Empty() {
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>			var extension cryptobyte.String
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>			if !extensions.ReadASN1(&amp;extension, cryptobyte_asn1.SEQUENCE) {
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>				return nil, errors.New(&#34;x509: malformed extension&#34;)
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>			}
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>			ext, err := parseExtension(extension)
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>			if err != nil {
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>				return nil, err
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>			}
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			if ext.Id.Equal(oidExtensionAuthorityKeyId) {
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>				rl.AuthorityKeyId = ext.Value
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			} else if ext.Id.Equal(oidExtensionCRLNumber) {
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>				value := cryptobyte.String(ext.Value)
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>				rl.Number = new(big.Int)
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>				if !value.ReadASN1Integer(rl.Number) {
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>					return nil, errors.New(&#34;x509: malformed crl number&#34;)
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>				}
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>			}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			rl.Extensions = append(rl.Extensions, ext)
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>		}
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	}
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>	return rl, nil
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>
</pre><p><a href="parser.go?m=text">View as plain text</a></p>

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
