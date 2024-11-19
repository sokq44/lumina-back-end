<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/x509/pkix/pkix.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="pkix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/x509">x509</a>/<a href="http://localhost:8080/src/crypto/x509/pkix">pkix</a>/<span class="text-muted">pkix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/x509/pkix">crypto/x509/pkix</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package pkix contains shared, low level structures used for ASN.1 parsing</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// and serialization of X.509 certificates, CRL and OCSP.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package pkix
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;encoding/asn1&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;encoding/hex&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;math/big&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// AlgorithmIdentifier represents the ASN.1 structure of the same name. See RFC</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// 5280, section 4.1.1.2.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type AlgorithmIdentifier struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	Algorithm  asn1.ObjectIdentifier
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	Parameters asn1.RawValue `asn1:&#34;optional&#34;`
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type RDNSequence []RelativeDistinguishedNameSET
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>var attributeTypeNames = map[string]string{
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;2.5.4.6&#34;:  &#34;C&#34;,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;2.5.4.10&#34;: &#34;O&#34;,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;2.5.4.11&#34;: &#34;OU&#34;,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	&#34;2.5.4.3&#34;:  &#34;CN&#34;,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	&#34;2.5.4.5&#34;:  &#34;SERIALNUMBER&#34;,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	&#34;2.5.4.7&#34;:  &#34;L&#34;,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	&#34;2.5.4.8&#34;:  &#34;ST&#34;,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	&#34;2.5.4.9&#34;:  &#34;STREET&#34;,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	&#34;2.5.4.17&#34;: &#34;POSTALCODE&#34;,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// String returns a string representation of the sequence r,</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// roughly following the RFC 2253 Distinguished Names syntax.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>func (r RDNSequence) String() string {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	s := &#34;&#34;
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	for i := 0; i &lt; len(r); i++ {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		rdn := r[len(r)-1-i]
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		if i &gt; 0 {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			s += &#34;,&#34;
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		for j, tv := range rdn {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			if j &gt; 0 {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				s += &#34;+&#34;
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			oidString := tv.Type.String()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			typeName, ok := attributeTypeNames[oidString]
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			if !ok {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>				derBytes, err := asn1.Marshal(tv.Value)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>				if err == nil {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>					s += oidString + &#34;=#&#34; + hex.EncodeToString(derBytes)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>					continue <span class="comment">// No value escaping necessary.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>				}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>				typeName = oidString
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			valueString := fmt.Sprint(tv.Value)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			escaped := make([]rune, 0, len(valueString))
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			for k, c := range valueString {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>				escape := false
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>				switch c {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>				case &#39;,&#39;, &#39;+&#39;, &#39;&#34;&#39;, &#39;\\&#39;, &#39;&lt;&#39;, &#39;&gt;&#39;, &#39;;&#39;:
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>					escape = true
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>				case &#39; &#39;:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>					escape = k == 0 || k == len(valueString)-1
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>				case &#39;#&#39;:
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>					escape = k == 0
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>				}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>				if escape {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>					escaped = append(escaped, &#39;\\&#39;, c)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				} else {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>					escaped = append(escaped, c)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>				}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			s += typeName + &#34;=&#34; + string(escaped)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	return s
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>type RelativeDistinguishedNameSET []AttributeTypeAndValue
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// AttributeTypeAndValue mirrors the ASN.1 structure of the same name in</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// RFC 5280, Section 4.1.2.4.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>type AttributeTypeAndValue struct {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	Type  asn1.ObjectIdentifier
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	Value any
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// AttributeTypeAndValueSET represents a set of ASN.1 sequences of</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// [AttributeTypeAndValue] sequences from RFC 2986 (PKCS #10).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>type AttributeTypeAndValueSET struct {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	Type  asn1.ObjectIdentifier
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	Value [][]AttributeTypeAndValue `asn1:&#34;set&#34;`
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// Extension represents the ASN.1 structure of the same name. See RFC</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// 5280, section 4.2.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>type Extension struct {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	Id       asn1.ObjectIdentifier
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	Critical bool `asn1:&#34;optional&#34;`
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	Value    []byte
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// Name represents an X.509 distinguished name. This only includes the common</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// elements of a DN. Note that Name is only an approximation of the X.509</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// structure. If an accurate representation is needed, asn1.Unmarshal the raw</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// subject or issuer as an [RDNSequence].</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>type Name struct {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	Country, Organization, OrganizationalUnit []string
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	Locality, Province                        []string
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	StreetAddress, PostalCode                 []string
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	SerialNumber, CommonName                  string
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// Names contains all parsed attributes. When parsing distinguished names,</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// this can be used to extract non-standard attributes that are not parsed</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// by this package. When marshaling to RDNSequences, the Names field is</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// ignored, see ExtraNames.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	Names []AttributeTypeAndValue
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// ExtraNames contains attributes to be copied, raw, into any marshaled</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// distinguished names. Values override any attributes with the same OID.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// The ExtraNames field is not populated when parsing, see Names.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	ExtraNames []AttributeTypeAndValue
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// FillFromRDNSequence populates n from the provided [RDNSequence].</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// Multi-entry RDNs are flattened, all entries are added to the</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// relevant n fields, and the grouping is not preserved.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func (n *Name) FillFromRDNSequence(rdns *RDNSequence) {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	for _, rdn := range *rdns {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if len(rdn) == 0 {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			continue
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		for _, atv := range rdn {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			n.Names = append(n.Names, atv)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			value, ok := atv.Value.(string)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			if !ok {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				continue
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			t := atv.Type
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>			if len(t) == 4 &amp;&amp; t[0] == 2 &amp;&amp; t[1] == 5 &amp;&amp; t[2] == 4 {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>				switch t[3] {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>				case 3:
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>					n.CommonName = value
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>				case 5:
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>					n.SerialNumber = value
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				case 6:
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>					n.Country = append(n.Country, value)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>				case 7:
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>					n.Locality = append(n.Locality, value)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				case 8:
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>					n.Province = append(n.Province, value)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				case 9:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>					n.StreetAddress = append(n.StreetAddress, value)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				case 10:
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>					n.Organization = append(n.Organization, value)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				case 11:
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>					n.OrganizationalUnit = append(n.OrganizationalUnit, value)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				case 17:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>					n.PostalCode = append(n.PostalCode, value)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>var (
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	oidCountry            = []int{2, 5, 4, 6}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	oidOrganization       = []int{2, 5, 4, 10}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	oidOrganizationalUnit = []int{2, 5, 4, 11}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	oidCommonName         = []int{2, 5, 4, 3}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	oidSerialNumber       = []int{2, 5, 4, 5}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	oidLocality           = []int{2, 5, 4, 7}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	oidProvince           = []int{2, 5, 4, 8}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	oidStreetAddress      = []int{2, 5, 4, 9}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	oidPostalCode         = []int{2, 5, 4, 17}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// appendRDNs appends a relativeDistinguishedNameSET to the given RDNSequence</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// and returns the new value. The relativeDistinguishedNameSET contains an</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// attributeTypeAndValue for each of the given values. See RFC 5280, A.1, and</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// search for AttributeTypeAndValue.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (n Name) appendRDNs(in RDNSequence, values []string, oid asn1.ObjectIdentifier) RDNSequence {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if len(values) == 0 || oidInAttributeTypeAndValue(oid, n.ExtraNames) {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return in
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	s := make([]AttributeTypeAndValue, len(values))
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	for i, value := range values {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		s[i].Type = oid
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		s[i].Value = value
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	return append(in, s)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// ToRDNSequence converts n into a single [RDNSequence]. The following</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// attributes are encoded as multi-value RDNs:</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">//   - Country</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">//   - Organization</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">//   - OrganizationalUnit</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">//   - Locality</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">//   - Province</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">//   - StreetAddress</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">//   - PostalCode</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// Each ExtraNames entry is encoded as an individual RDN.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func (n Name) ToRDNSequence() (ret RDNSequence) {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.Country, oidCountry)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.Province, oidProvince)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.Locality, oidLocality)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.StreetAddress, oidStreetAddress)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.PostalCode, oidPostalCode)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.Organization, oidOrganization)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	ret = n.appendRDNs(ret, n.OrganizationalUnit, oidOrganizationalUnit)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	if len(n.CommonName) &gt; 0 {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		ret = n.appendRDNs(ret, []string{n.CommonName}, oidCommonName)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if len(n.SerialNumber) &gt; 0 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		ret = n.appendRDNs(ret, []string{n.SerialNumber}, oidSerialNumber)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	for _, atv := range n.ExtraNames {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		ret = append(ret, []AttributeTypeAndValue{atv})
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	return ret
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// String returns the string form of n, roughly following</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span><span class="comment">// the RFC 2253 Distinguished Names syntax.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func (n Name) String() string {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	var rdns RDNSequence
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// If there are no ExtraNames, surface the parsed value (all entries in</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// Names) instead.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if n.ExtraNames == nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		for _, atv := range n.Names {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			t := atv.Type
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			if len(t) == 4 &amp;&amp; t[0] == 2 &amp;&amp; t[1] == 5 &amp;&amp; t[2] == 4 {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>				switch t[3] {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				case 3, 5, 6, 7, 8, 9, 10, 11, 17:
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>					<span class="comment">// These attributes were already parsed into named fields.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>					continue
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			<span class="comment">// Place non-standard parsed values at the beginning of the sequence</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			<span class="comment">// so they will be at the end of the string. See Issue 39924.</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			rdns = append(rdns, []AttributeTypeAndValue{atv})
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	rdns = append(rdns, n.ToRDNSequence()...)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	return rdns.String()
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// oidInAttributeTypeAndValue reports whether a type with the given OID exists</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">// in atv.</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>func oidInAttributeTypeAndValue(oid asn1.ObjectIdentifier, atv []AttributeTypeAndValue) bool {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	for _, a := range atv {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		if a.Type.Equal(oid) {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			return true
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	return false
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// CertificateList represents the ASN.1 structure of the same name. See RFC</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// 5280, section 5.1. Use Certificate.CheckCRLSignature to verify the</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// signature.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// Deprecated: x509.RevocationList should be used instead.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>type CertificateList struct {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	TBSCertList        TBSCertificateList
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	SignatureAlgorithm AlgorithmIdentifier
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	SignatureValue     asn1.BitString
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// HasExpired reports whether certList should have been updated by now.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>func (certList *CertificateList) HasExpired(now time.Time) bool {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	return !now.Before(certList.TBSCertList.NextUpdate)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">// TBSCertificateList represents the ASN.1 structure of the same name. See RFC</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// 5280, section 5.1.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// Deprecated: x509.RevocationList should be used instead.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>type TBSCertificateList struct {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	Raw                 asn1.RawContent
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	Version             int `asn1:&#34;optional,default:0&#34;`
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	Signature           AlgorithmIdentifier
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	Issuer              RDNSequence
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	ThisUpdate          time.Time
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	NextUpdate          time.Time            `asn1:&#34;optional&#34;`
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	RevokedCertificates []RevokedCertificate `asn1:&#34;optional&#34;`
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	Extensions          []Extension          `asn1:&#34;tag:0,optional,explicit&#34;`
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// RevokedCertificate represents the ASN.1 structure of the same name. See RFC</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// 5280, section 5.1.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>type RevokedCertificate struct {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	SerialNumber   *big.Int
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	RevocationTime time.Time
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	Extensions     []Extension `asn1:&#34;optional&#34;`
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
</pre><p><a href="pkix.go?m=text">View as plain text</a></p>

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
