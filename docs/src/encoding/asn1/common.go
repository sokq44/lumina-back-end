<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/asn1/common.go - Go Documentation Server</title>

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
<a href="common.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/asn1">asn1</a>/<span class="text-muted">common.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/encoding/asn1">encoding/asn1</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package asn1
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// ASN.1 objects have metadata preceding them:</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//   the tag: the type of the object</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//   a flag denoting if this object is compound or not</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//   the class type: the namespace of the tag</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//   the length of the object, in bytes</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Here are some standard tags and classes</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// ASN.1 tags represent the type of the following object.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const (
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	TagBoolean         = 1
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	TagInteger         = 2
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	TagBitString       = 3
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	TagOctetString     = 4
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	TagNull            = 5
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	TagOID             = 6
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	TagEnum            = 10
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	TagUTF8String      = 12
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	TagSequence        = 16
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	TagSet             = 17
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	TagNumericString   = 18
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	TagPrintableString = 19
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	TagT61String       = 20
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	TagIA5String       = 22
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	TagUTCTime         = 23
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	TagGeneralizedTime = 24
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	TagGeneralString   = 27
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	TagBMPString       = 30
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// ASN.1 class types represent the namespace of the tag.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>const (
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	ClassUniversal       = 0
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	ClassApplication     = 1
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	ClassContextSpecific = 2
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	ClassPrivate         = 3
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>type tagAndLength struct {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	class, tag, length int
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	isCompound         bool
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// ASN.1 has IMPLICIT and EXPLICIT tags, which can be translated as &#34;instead</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// of&#34; and &#34;in addition to&#34;. When not specified, every primitive type has a</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// default tag in the UNIVERSAL class.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// For example: a BIT STRING is tagged [UNIVERSAL 3] by default (although ASN.1</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// doesn&#39;t actually have a UNIVERSAL keyword). However, by saying [IMPLICIT</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// CONTEXT-SPECIFIC 42], that means that the tag is replaced by another.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// On the other hand, if it said [EXPLICIT CONTEXT-SPECIFIC 10], then an</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// /additional/ tag would wrap the default tag. This explicit tag will have the</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// compound flag set.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// (This is used in order to remove ambiguity with optional elements.)</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// You can layer EXPLICIT and IMPLICIT tags to an arbitrary depth, however we</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// don&#39;t support that here. We support a single layer of EXPLICIT or IMPLICIT</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// tagging with tag strings on the fields of a structure.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// fieldParameters is the parsed representation of tag string from a structure field.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>type fieldParameters struct {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	optional     bool   <span class="comment">// true iff the field is OPTIONAL</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	explicit     bool   <span class="comment">// true iff an EXPLICIT tag is in use.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	application  bool   <span class="comment">// true iff an APPLICATION tag is in use.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	private      bool   <span class="comment">// true iff a PRIVATE tag is in use.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	defaultValue *int64 <span class="comment">// a default value for INTEGER typed fields (maybe nil).</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	tag          *int   <span class="comment">// the EXPLICIT or IMPLICIT tag (maybe nil).</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	stringType   int    <span class="comment">// the string tag to use when marshaling.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	timeType     int    <span class="comment">// the time tag to use when marshaling.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	set          bool   <span class="comment">// true iff this should be encoded as a SET</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	omitEmpty    bool   <span class="comment">// true iff this should be omitted if empty when marshaling.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// Invariants:</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">//   if explicit is set, tag is non-nil.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// Given a tag string with the format specified in the package comment,</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// parseFieldParameters will parse it into a fieldParameters structure,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// ignoring unknown parts of the string.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func parseFieldParameters(str string) (ret fieldParameters) {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	var part string
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	for len(str) &gt; 0 {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		part, str, _ = strings.Cut(str, &#34;,&#34;)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		switch {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		case part == &#34;optional&#34;:
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>			ret.optional = true
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		case part == &#34;explicit&#34;:
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			ret.explicit = true
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			if ret.tag == nil {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				ret.tag = new(int)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		case part == &#34;generalized&#34;:
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			ret.timeType = TagGeneralizedTime
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		case part == &#34;utc&#34;:
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			ret.timeType = TagUTCTime
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		case part == &#34;ia5&#34;:
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			ret.stringType = TagIA5String
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		case part == &#34;printable&#34;:
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			ret.stringType = TagPrintableString
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		case part == &#34;numeric&#34;:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			ret.stringType = TagNumericString
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		case part == &#34;utf8&#34;:
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			ret.stringType = TagUTF8String
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		case strings.HasPrefix(part, &#34;default:&#34;):
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			i, err := strconv.ParseInt(part[8:], 10, 64)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			if err == nil {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				ret.defaultValue = new(int64)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>				*ret.defaultValue = i
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		case strings.HasPrefix(part, &#34;tag:&#34;):
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			i, err := strconv.Atoi(part[4:])
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			if err == nil {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				ret.tag = new(int)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				*ret.tag = i
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		case part == &#34;set&#34;:
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			ret.set = true
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		case part == &#34;application&#34;:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			ret.application = true
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			if ret.tag == nil {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				ret.tag = new(int)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		case part == &#34;private&#34;:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			ret.private = true
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if ret.tag == nil {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				ret.tag = new(int)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		case part == &#34;omitempty&#34;:
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			ret.omitEmpty = true
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// Given a reflected Go type, getUniversalType returns the default tag number</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// and expected compound flag.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>func getUniversalType(t reflect.Type) (matchAny bool, tagNumber int, isCompound, ok bool) {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	switch t {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	case rawValueType:
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return true, -1, false, true
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	case objectIdentifierType:
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		return false, TagOID, false, true
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	case bitStringType:
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return false, TagBitString, false, true
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	case timeType:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return false, TagUTCTime, false, true
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	case enumeratedType:
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		return false, TagEnum, false, true
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	case bigIntType:
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return false, TagInteger, false, true
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	switch t.Kind() {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return false, TagBoolean, false, true
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		return false, TagInteger, false, true
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		return false, TagSequence, true, true
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		if t.Elem().Kind() == reflect.Uint8 {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			return false, TagOctetString, false, true
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if strings.HasSuffix(t.Name(), &#34;SET&#34;) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			return false, TagSet, true, true
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return false, TagSequence, true, true
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	case reflect.String:
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		return false, TagPrintableString, false, true
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return false, 0, false, false
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
</pre><p><a href="common.go?m=text">View as plain text</a></p>

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
