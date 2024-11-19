<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/hash/adler32/adler32.go - Go Documentation Server</title>

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
<a href="adler32.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/hash">hash</a>/<a href="http://localhost:8080/src/hash/adler32">adler32</a>/<span class="text-muted">adler32.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/hash/adler32">hash/adler32</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package adler32 implements the Adler-32 checksum.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// It is defined in RFC 1950:</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//	Adler-32 is composed of two sums accumulated per byte: s1 is</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//	the sum of all bytes, s2 is the sum of all s1 values. Both sums</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//	are done modulo 65521. s1 is initialized to 1, s2 to zero.  The</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//	Adler-32 checksum is stored as s2*65536 + s1 in most-</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//	significant-byte first (network) order.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>package adler32
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>import (
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>const (
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// mod is the largest prime that is less than 65536.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	mod = 65521
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// nmax is the largest n such that</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// 255 * n * (n+1) / 2 + (n+1) * (mod-1) &lt;= 2^32-1.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// It is mentioned in RFC 1950 (search for &#34;5552&#34;).</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	nmax = 5552
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// The size of an Adler-32 checksum in bytes.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>const Size = 4
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// digest represents the partial evaluation of a checksum.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// The low 16 bits are s1, the high 16 bits are s2.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>type digest uint32
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func (d *digest) Reset() { *d = 1 }
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// New returns a new hash.Hash32 computing the Adler-32 checksum. Its</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// Sum method will lay the value out in big-endian byte order. The</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// returned Hash32 also implements [encoding.BinaryMarshaler] and</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// [encoding.BinaryUnmarshaler] to marshal and unmarshal the internal</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// state of the hash.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func New() hash.Hash32 {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	d := new(digest)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	d.Reset()
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	return d
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func (d *digest) Size() int { return Size }
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>func (d *digest) BlockSize() int { return 4 }
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>const (
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	magic         = &#34;adl\x01&#34;
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	marshaledSize = len(magic) + 4
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (d *digest) MarshalBinary() ([]byte, error) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	b = append(b, magic...)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	b = appendUint32(b, uint32(*d))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return b, nil
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func (d *digest) UnmarshalBinary(b []byte) error {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if len(b) &lt; len(magic) || string(b[:len(magic)]) != magic {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		return errors.New(&#34;hash/adler32: invalid hash state identifier&#34;)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	if len(b) != marshaledSize {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		return errors.New(&#34;hash/adler32: invalid hash state size&#34;)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	*d = digest(readUint32(b[len(magic):]))
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return nil
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// appendUint32 is semantically the same as [binary.BigEndian.AppendUint32]</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func appendUint32(b []byte, x uint32) []byte {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	return append(b,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		byte(x&gt;&gt;24),
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		byte(x&gt;&gt;16),
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		byte(x&gt;&gt;8),
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		byte(x),
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// readUint32 is semantically the same as [binary.BigEndian.Uint32]</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func readUint32(b []byte) uint32 {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	_ = b[3]
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	return uint32(b[3]) | uint32(b[2])&lt;&lt;8 | uint32(b[1])&lt;&lt;16 | uint32(b[0])&lt;&lt;24
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// Add p to the running checksum d.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func update(d digest, p []byte) digest {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	s1, s2 := uint32(d&amp;0xffff), uint32(d&gt;&gt;16)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	for len(p) &gt; 0 {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		var q []byte
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		if len(p) &gt; nmax {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			p, q = p[:nmax], p[nmax:]
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		for len(p) &gt;= 4 {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			s1 += uint32(p[0])
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			s2 += s1
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			s1 += uint32(p[1])
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			s2 += s1
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			s1 += uint32(p[2])
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			s2 += s1
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			s1 += uint32(p[3])
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			s2 += s1
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			p = p[4:]
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		for _, x := range p {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			s1 += uint32(x)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			s2 += s1
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		s1 %= mod
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		s2 %= mod
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		p = q
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	return digest(s2&lt;&lt;16 | s1)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>func (d *digest) Write(p []byte) (nn int, err error) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	*d = update(*d, p)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	return len(p), nil
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func (d *digest) Sum32() uint32 { return uint32(*d) }
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>func (d *digest) Sum(in []byte) []byte {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	s := uint32(*d)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	return append(in, byte(s&gt;&gt;24), byte(s&gt;&gt;16), byte(s&gt;&gt;8), byte(s))
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// Checksum returns the Adler-32 checksum of data.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func Checksum(data []byte) uint32 { return uint32(update(1, data)) }
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
</pre><p><a href="adler32.go?m=text">View as plain text</a></p>

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
