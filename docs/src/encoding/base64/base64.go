<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/base64/base64.go - Go Documentation Server</title>

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
<a href="base64.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/base64">base64</a>/<span class="text-muted">base64.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/encoding/base64">encoding/base64</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package base64 implements base64 encoding as specified by RFC 4648.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>package base64
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>import (
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;slices&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">/*
<span id="L16" class="ln">    16&nbsp;&nbsp;</span> * Encodings
<span id="L17" class="ln">    17&nbsp;&nbsp;</span> */</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// An Encoding is a radix 64 encoding/decoding scheme, defined by a</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// 64-character alphabet. The most common encoding is the &#34;base64&#34;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// encoding defined in RFC 4648 and used in MIME (RFC 2045) and PEM</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// (RFC 1421).  RFC 4648 also defines an alternate encoding, which is</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// the standard encoding with - and _ substituted for + and /.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type Encoding struct {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	encode    [64]byte   <span class="comment">// mapping of symbol index to symbol byte value</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	decodeMap [256]uint8 <span class="comment">// mapping of symbol byte value to symbol index</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	padChar   rune
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	strict    bool
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>const (
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	StdPadding rune = &#39;=&#39; <span class="comment">// Standard padding character</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	NoPadding  rune = -1  <span class="comment">// No padding</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>const (
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	decodeMapInitialize = &#34;&#34; +
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34; +
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		&#34;\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff&#34;
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	invalidIndex = &#39;\xff&#39;
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// NewEncoding returns a new padded Encoding defined by the given alphabet,</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// which must be a 64-byte string that contains unique byte values and</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// does not contain the padding character or CR / LF (&#39;\r&#39;, &#39;\n&#39;).</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// The alphabet is treated as a sequence of byte values</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// without any special treatment for multi-byte UTF-8.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// The resulting Encoding uses the default padding character (&#39;=&#39;),</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// which may be changed or disabled via [Encoding.WithPadding].</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func NewEncoding(encoder string) *Encoding {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if len(encoder) != 64 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		panic(&#34;encoding alphabet is not 64-bytes long&#34;)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	e := new(Encoding)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	e.padChar = StdPadding
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	copy(e.encode[:], encoder)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	copy(e.decodeMap[:], decodeMapInitialize)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	for i := 0; i &lt; len(encoder); i++ {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// Note: While we document that the alphabet cannot contain</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// the padding character, we do not enforce it since we do not know</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		<span class="comment">// if the caller intends to switch the padding from StdPadding later.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		switch {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		case encoder[i] == &#39;\n&#39; || encoder[i] == &#39;\r&#39;:
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			panic(&#34;encoding alphabet contains newline character&#34;)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		case e.decodeMap[encoder[i]] != invalidIndex:
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			panic(&#34;encoding alphabet includes duplicate symbols&#34;)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		e.decodeMap[encoder[i]] = uint8(i)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	return e
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// WithPadding creates a new encoding identical to enc except</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// with a specified padding character, or [NoPadding] to disable padding.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// The padding character must not be &#39;\r&#39; or &#39;\n&#39;,</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// must not be contained in the encoding&#39;s alphabet,</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// must not be negative, and must be a rune equal or below &#39;\xff&#39;.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Padding characters above &#39;\x7f&#39; are encoded as their exact byte value</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// rather than using the UTF-8 representation of the codepoint.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (enc Encoding) WithPadding(padding rune) *Encoding {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	switch {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	case padding &lt; NoPadding || padding == &#39;\r&#39; || padding == &#39;\n&#39; || padding &gt; 0xff:
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		panic(&#34;invalid padding&#34;)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	case padding != NoPadding &amp;&amp; enc.decodeMap[byte(padding)] != invalidIndex:
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		panic(&#34;padding contained in alphabet&#34;)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	enc.padChar = padding
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	return &amp;enc
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// Strict creates a new encoding identical to enc except with</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// strict decoding enabled. In this mode, the decoder requires that</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// trailing padding bits are zero, as described in RFC 4648 section 3.5.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// Note that the input is still malleable, as new line characters</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// (CR and LF) are still ignored.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func (enc Encoding) Strict() *Encoding {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	enc.strict = true
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	return &amp;enc
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// StdEncoding is the standard base64 encoding, as defined in RFC 4648.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>var StdEncoding = NewEncoding(&#34;ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/&#34;)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// URLEncoding is the alternate base64 encoding defined in RFC 4648.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// It is typically used in URLs and file names.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>var URLEncoding = NewEncoding(&#34;ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_&#34;)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// RawStdEncoding is the standard raw, unpadded base64 encoding,</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// as defined in RFC 4648 section 3.2.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// This is the same as [StdEncoding] but omits padding characters.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>var RawStdEncoding = StdEncoding.WithPadding(NoPadding)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// RawURLEncoding is the unpadded alternate base64 encoding defined in RFC 4648.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// It is typically used in URLs and file names.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// This is the same as [URLEncoding] but omits padding characters.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>var RawURLEncoding = URLEncoding.WithPadding(NoPadding)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">/*
<span id="L136" class="ln">   136&nbsp;&nbsp;</span> * Encoder
<span id="L137" class="ln">   137&nbsp;&nbsp;</span> */</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// Encode encodes src using the encoding enc,</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// writing [Encoding.EncodedLen](len(src)) bytes to dst.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// The encoding pads the output to a multiple of 4 bytes,</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// so Encode is not appropriate for use on individual blocks</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// of a large data stream. Use [NewEncoder] instead.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func (enc *Encoding) Encode(dst, src []byte) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if len(src) == 0 {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// enc is a pointer receiver, so the use of enc.encode within the hot</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// loop below means a nil check at every operation. Lift that nil check</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// outside of the loop to speed up the encoder.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	_ = enc.encode
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	di, si := 0, 0
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	n := (len(src) / 3) * 3
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	for si &lt; n {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// Convert 3x 8bit source bytes into 4 bytes</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		val := uint(src[si+0])&lt;&lt;16 | uint(src[si+1])&lt;&lt;8 | uint(src[si+2])
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		dst[di+0] = enc.encode[val&gt;&gt;18&amp;0x3F]
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		dst[di+1] = enc.encode[val&gt;&gt;12&amp;0x3F]
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		dst[di+2] = enc.encode[val&gt;&gt;6&amp;0x3F]
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		dst[di+3] = enc.encode[val&amp;0x3F]
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		si += 3
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		di += 4
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	remain := len(src) - si
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if remain == 0 {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// Add the remaining small block</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	val := uint(src[si+0]) &lt;&lt; 16
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if remain == 2 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		val |= uint(src[si+1]) &lt;&lt; 8
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	dst[di+0] = enc.encode[val&gt;&gt;18&amp;0x3F]
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	dst[di+1] = enc.encode[val&gt;&gt;12&amp;0x3F]
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	switch remain {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	case 2:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		dst[di+2] = enc.encode[val&gt;&gt;6&amp;0x3F]
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		if enc.padChar != NoPadding {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			dst[di+3] = byte(enc.padChar)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	case 1:
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		if enc.padChar != NoPadding {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			dst[di+2] = byte(enc.padChar)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			dst[di+3] = byte(enc.padChar)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// AppendEncode appends the base64 encoded src to dst</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// and returns the extended buffer.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>func (enc *Encoding) AppendEncode(dst, src []byte) []byte {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	n := enc.EncodedLen(len(src))
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	dst = slices.Grow(dst, n)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	enc.Encode(dst[len(dst):][:n], src)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	return dst[:len(dst)+n]
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// EncodeToString returns the base64 encoding of src.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func (enc *Encoding) EncodeToString(src []byte) string {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	buf := make([]byte, enc.EncodedLen(len(src)))
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	enc.Encode(buf, src)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	return string(buf)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>type encoder struct {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	err  error
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	enc  *Encoding
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	w    io.Writer
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	buf  [3]byte    <span class="comment">// buffered data waiting to be encoded</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	nbuf int        <span class="comment">// number of bytes in buf</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	out  [1024]byte <span class="comment">// output buffer</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>func (e *encoder) Write(p []byte) (n int, err error) {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	if e.err != nil {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		return 0, e.err
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// Leading fringe.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if e.nbuf &gt; 0 {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		var i int
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		for i = 0; i &lt; len(p) &amp;&amp; e.nbuf &lt; 3; i++ {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			e.buf[e.nbuf] = p[i]
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			e.nbuf++
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		n += i
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		p = p[i:]
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		if e.nbuf &lt; 3 {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			return
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		e.enc.Encode(e.out[:], e.buf[:])
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		if _, e.err = e.w.Write(e.out[:4]); e.err != nil {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			return n, e.err
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		e.nbuf = 0
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// Large interior chunks.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	for len(p) &gt;= 3 {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		nn := len(e.out) / 4 * 3
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		if nn &gt; len(p) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			nn = len(p)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			nn -= nn % 3
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		e.enc.Encode(e.out[:], p[:nn])
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		if _, e.err = e.w.Write(e.out[0 : nn/3*4]); e.err != nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			return n, e.err
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		n += nn
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		p = p[nn:]
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// Trailing fringe.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	copy(e.buf[:], p)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	e.nbuf = len(p)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	n += len(p)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	return
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">// Close flushes any pending output from the encoder.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// It is an error to call Write after calling Close.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func (e *encoder) Close() error {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	<span class="comment">// If there&#39;s anything left in the buffer, flush it out</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if e.err == nil &amp;&amp; e.nbuf &gt; 0 {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		e.enc.Encode(e.out[:], e.buf[:e.nbuf])
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		_, e.err = e.w.Write(e.out[:e.enc.EncodedLen(e.nbuf)])
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		e.nbuf = 0
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	return e.err
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// NewEncoder returns a new base64 stream encoder. Data written to</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// the returned writer will be encoded using enc and then written to w.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// Base64 encodings operate in 4-byte blocks; when finished</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// writing, the caller must Close the returned encoder to flush any</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// partially written blocks.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return &amp;encoder{enc: enc, w: w}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// EncodedLen returns the length in bytes of the base64 encoding</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// of an input buffer of length n.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>func (enc *Encoding) EncodedLen(n int) int {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	if enc.padChar == NoPadding {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		return n/3*4 + (n%3*8+5)/6 <span class="comment">// minimum # chars at 6 bits per char</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	return (n + 2) / 3 * 4 <span class="comment">// minimum # 4-char quanta, 3 bytes each</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">/*
<span id="L298" class="ln">   298&nbsp;&nbsp;</span> * Decoder
<span id="L299" class="ln">   299&nbsp;&nbsp;</span> */</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>type CorruptInputError int64
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>func (e CorruptInputError) Error() string {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	return &#34;illegal base64 data at input byte &#34; + strconv.FormatInt(int64(e), 10)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// decodeQuantum decodes up to 4 base64 bytes. The received parameters are</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// the destination buffer dst, the source buffer src and an index in the</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// source buffer si.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// It returns the number of bytes read from src, the number of bytes written</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// to dst, and an error, if any.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func (enc *Encoding) decodeQuantum(dst, src []byte, si int) (nsi, n int, err error) {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	<span class="comment">// Decode quantum using the base64 alphabet</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	var dbuf [4]byte
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	dlen := 4
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// Lift the nil check outside of the loop.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	_ = enc.decodeMap
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	for j := 0; j &lt; len(dbuf); j++ {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		if len(src) == si {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			switch {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			case j == 0:
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				return si, 0, nil
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			case j == 1, enc.padChar != NoPadding:
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>				return si, 0, CorruptInputError(si - j)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			dlen = j
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			break
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		in := src[si]
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		si++
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		out := enc.decodeMap[in]
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		if out != 0xff {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			dbuf[j] = out
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			continue
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		if in == &#39;\n&#39; || in == &#39;\r&#39; {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			j--
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			continue
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		if rune(in) != enc.padChar {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			return si, 0, CorruptInputError(si - 1)
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		<span class="comment">// We&#39;ve reached the end and there&#39;s padding</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		switch j {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		case 0, 1:
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			<span class="comment">// incorrect padding</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			return si, 0, CorruptInputError(si - 1)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		case 2:
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			<span class="comment">// &#34;==&#34; is expected, the first &#34;=&#34; is already consumed.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			<span class="comment">// skip over newlines</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			for si &lt; len(src) &amp;&amp; (src[si] == &#39;\n&#39; || src[si] == &#39;\r&#39;) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>				si++
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			if si == len(src) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>				<span class="comment">// not enough padding</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>				return si, 0, CorruptInputError(len(src))
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			if rune(src[si]) != enc.padChar {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				<span class="comment">// incorrect padding</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				return si, 0, CorruptInputError(si - 1)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			si++
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		<span class="comment">// skip over newlines</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		for si &lt; len(src) &amp;&amp; (src[si] == &#39;\n&#39; || src[si] == &#39;\r&#39;) {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			si++
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if si &lt; len(src) {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			<span class="comment">// trailing garbage</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			err = CorruptInputError(si)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		dlen = j
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		break
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// Convert 4x 6bit source bytes into 3 bytes</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	val := uint(dbuf[0])&lt;&lt;18 | uint(dbuf[1])&lt;&lt;12 | uint(dbuf[2])&lt;&lt;6 | uint(dbuf[3])
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	dbuf[2], dbuf[1], dbuf[0] = byte(val&gt;&gt;0), byte(val&gt;&gt;8), byte(val&gt;&gt;16)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	switch dlen {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	case 4:
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		dst[2] = dbuf[2]
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		dbuf[2] = 0
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		fallthrough
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	case 3:
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		dst[1] = dbuf[1]
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		if enc.strict &amp;&amp; dbuf[2] != 0 {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			return si, 0, CorruptInputError(si - 1)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		dbuf[1] = 0
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		fallthrough
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	case 2:
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		dst[0] = dbuf[0]
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		if enc.strict &amp;&amp; (dbuf[1] != 0 || dbuf[2] != 0) {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			return si, 0, CorruptInputError(si - 2)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		}
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	return si, dlen - 1, err
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// AppendDecode appends the base64 decoded src to dst</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// and returns the extended buffer.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// If the input is malformed, it returns the partially decoded src and an error.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>func (enc *Encoding) AppendDecode(dst, src []byte) ([]byte, error) {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	<span class="comment">// Compute the output size without padding to avoid over allocating.</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	n := len(src)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	for n &gt; 0 &amp;&amp; rune(src[n-1]) == enc.padChar {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		n--
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	n = decodedLen(n, NoPadding)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	dst = slices.Grow(dst, n)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	n, err := enc.Decode(dst[len(dst):][:n], src)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	return dst[:len(dst)+n], err
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">// DecodeString returns the bytes represented by the base64 string s.</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>func (enc *Encoding) DecodeString(s string) ([]byte, error) {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	dbuf := make([]byte, enc.DecodedLen(len(s)))
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	n, err := enc.Decode(dbuf, []byte(s))
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	return dbuf[:n], err
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>type decoder struct {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	err     error
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	readErr error <span class="comment">// error from r.Read</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	enc     *Encoding
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	r       io.Reader
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	buf     [1024]byte <span class="comment">// leftover input</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	nbuf    int
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	out     []byte <span class="comment">// leftover decoded output</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	outbuf  [1024 / 4 * 3]byte
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>func (d *decoder) Read(p []byte) (n int, err error) {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">// Use leftover decoded output from last read.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	if len(d.out) &gt; 0 {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		n = copy(p, d.out)
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		d.out = d.out[n:]
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		return n, nil
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	if d.err != nil {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		return 0, d.err
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">// This code assumes that d.r strips supported whitespace (&#39;\r&#39; and &#39;\n&#39;).</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	<span class="comment">// Refill buffer.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	for d.nbuf &lt; 4 &amp;&amp; d.readErr == nil {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		nn := len(p) / 3 * 4
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		if nn &lt; 4 {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			nn = 4
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		if nn &gt; len(d.buf) {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			nn = len(d.buf)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		nn, d.readErr = d.r.Read(d.buf[d.nbuf:nn])
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		d.nbuf += nn
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	if d.nbuf &lt; 4 {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		if d.enc.padChar == NoPadding &amp;&amp; d.nbuf &gt; 0 {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			<span class="comment">// Decode final fragment, without padding.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			var nw int
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			nw, d.err = d.enc.Decode(d.outbuf[:], d.buf[:d.nbuf])
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			d.nbuf = 0
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			d.out = d.outbuf[:nw]
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			n = copy(p, d.out)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			d.out = d.out[n:]
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			if n &gt; 0 || len(p) == 0 &amp;&amp; len(d.out) &gt; 0 {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>				return n, nil
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			if d.err != nil {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>				return 0, d.err
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			}
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		d.err = d.readErr
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		if d.err == io.EOF &amp;&amp; d.nbuf &gt; 0 {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			d.err = io.ErrUnexpectedEOF
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		return 0, d.err
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	<span class="comment">// Decode chunk into p, or d.out and then p if p is too small.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	nr := d.nbuf / 4 * 4
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	nw := d.nbuf / 4 * 3
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	if nw &gt; len(p) {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		nw, d.err = d.enc.Decode(d.outbuf[:], d.buf[:nr])
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		d.out = d.outbuf[:nw]
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		n = copy(p, d.out)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		d.out = d.out[n:]
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	} else {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		n, d.err = d.enc.Decode(p, d.buf[:nr])
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	d.nbuf -= nr
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	copy(d.buf[:d.nbuf], d.buf[nr:])
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	return n, d.err
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">// Decode decodes src using the encoding enc. It writes at most</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span><span class="comment">// [Encoding.DecodedLen](len(src)) bytes to dst and returns the number of bytes</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span><span class="comment">// written. If src contains invalid base64 data, it will return the</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// number of bytes successfully written and [CorruptInputError].</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// New line characters (\r and \n) are ignored.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if len(src) == 0 {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		return 0, nil
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	<span class="comment">// Lift the nil check outside of the loop. enc.decodeMap is directly</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	<span class="comment">// used later in this function, to let the compiler know that the</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	<span class="comment">// receiver can&#39;t be nil.</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	_ = enc.decodeMap
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	si := 0
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	for strconv.IntSize &gt;= 64 &amp;&amp; len(src)-si &gt;= 8 &amp;&amp; len(dst)-n &gt;= 8 {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		src2 := src[si : si+8]
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		if dn, ok := assemble64(
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			enc.decodeMap[src2[0]],
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			enc.decodeMap[src2[1]],
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>			enc.decodeMap[src2[2]],
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			enc.decodeMap[src2[3]],
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>			enc.decodeMap[src2[4]],
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			enc.decodeMap[src2[5]],
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			enc.decodeMap[src2[6]],
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			enc.decodeMap[src2[7]],
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		); ok {
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			binary.BigEndian.PutUint64(dst[n:], dn)
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			n += 6
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			si += 8
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		} else {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			var ninc int
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>			n += ninc
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			if err != nil {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				return n, err
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	for len(src)-si &gt;= 4 &amp;&amp; len(dst)-n &gt;= 4 {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		src2 := src[si : si+4]
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		if dn, ok := assemble32(
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			enc.decodeMap[src2[0]],
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			enc.decodeMap[src2[1]],
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			enc.decodeMap[src2[2]],
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			enc.decodeMap[src2[3]],
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		); ok {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			binary.BigEndian.PutUint32(dst[n:], dn)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			n += 3
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			si += 4
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		} else {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			var ninc int
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			n += ninc
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			if err != nil {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>				return n, err
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	for si &lt; len(src) {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		var ninc int
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		n += ninc
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		if err != nil {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			return n, err
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	return n, err
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// assemble32 assembles 4 base64 digits into 3 bytes.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">// Each digit comes from the decode map, and will be 0xff</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span><span class="comment">// if it came from an invalid character.</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>func assemble32(n1, n2, n3, n4 byte) (dn uint32, ok bool) {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	<span class="comment">// Check that all the digits are valid. If any of them was 0xff, their</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// bitwise OR will be 0xff.</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	if n1|n2|n3|n4 == 0xff {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		return 0, false
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	return uint32(n1)&lt;&lt;26 |
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>			uint32(n2)&lt;&lt;20 |
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>			uint32(n3)&lt;&lt;14 |
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>			uint32(n4)&lt;&lt;8,
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		true
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span><span class="comment">// assemble64 assembles 8 base64 digits into 6 bytes.</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span><span class="comment">// Each digit comes from the decode map, and will be 0xff</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span><span class="comment">// if it came from an invalid character.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>func assemble64(n1, n2, n3, n4, n5, n6, n7, n8 byte) (dn uint64, ok bool) {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	<span class="comment">// Check that all the digits are valid. If any of them was 0xff, their</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	<span class="comment">// bitwise OR will be 0xff.</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	if n1|n2|n3|n4|n5|n6|n7|n8 == 0xff {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		return 0, false
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	return uint64(n1)&lt;&lt;58 |
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			uint64(n2)&lt;&lt;52 |
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>			uint64(n3)&lt;&lt;46 |
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			uint64(n4)&lt;&lt;40 |
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			uint64(n5)&lt;&lt;34 |
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			uint64(n6)&lt;&lt;28 |
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>			uint64(n7)&lt;&lt;22 |
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			uint64(n8)&lt;&lt;16,
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		true
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>type newlineFilteringReader struct {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	wrapped io.Reader
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>func (r *newlineFilteringReader) Read(p []byte) (int, error) {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	n, err := r.wrapped.Read(p)
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		offset := 0
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		for i, b := range p[:n] {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			if b != &#39;\r&#39; &amp;&amp; b != &#39;\n&#39; {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>				if i != offset {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>					p[offset] = b
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>				}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				offset++
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>		}
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>		if offset &gt; 0 {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			return offset, err
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>		<span class="comment">// Previous buffer entirely whitespace, read again</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>		n, err = r.wrapped.Read(p)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	return n, err
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span><span class="comment">// NewDecoder constructs a new base64 stream decoder.</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	return &amp;decoder{enc: enc, r: &amp;newlineFilteringReader{r}}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span><span class="comment">// DecodedLen returns the maximum length in bytes of the decoded data</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span><span class="comment">// corresponding to n bytes of base64-encoded data.</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>func (enc *Encoding) DecodedLen(n int) int {
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	return decodedLen(n, enc.padChar)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>func decodedLen(n int, padChar rune) int {
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	if padChar == NoPadding {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		<span class="comment">// Unpadded data may end with partial block of 2-3 characters.</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		return n/4*3 + n%4*6/8
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// Padded base64 should always be a multiple of 4 characters in length.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	return n / 4 * 3
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>
</pre><p><a href="base64.go?m=text">View as plain text</a></p>

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
