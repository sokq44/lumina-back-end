<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/hash/fnv/fnv.go - Go Documentation Server</title>

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
<a href="fnv.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/hash">hash</a>/<a href="http://localhost:8080/src/hash/fnv">fnv</a>/<span class="text-muted">fnv.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/hash/fnv">hash/fnv</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package fnv implements FNV-1 and FNV-1a, non-cryptographic hash functions</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// See</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// All the hash.Hash implementations returned by this package also</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// implement encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// marshal and unmarshal the internal state of the hash.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>package fnv
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>import (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type (
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	sum32   uint32
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	sum32a  uint32
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	sum64   uint64
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	sum64a  uint64
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	sum128  [2]uint64
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	sum128a [2]uint64
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>const (
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	offset32        = 2166136261
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	offset64        = 14695981039346656037
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	offset128Lower  = 0x62b821756295c58d
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	offset128Higher = 0x6c62272e07bb0142
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	prime32         = 16777619
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	prime64         = 1099511628211
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	prime128Lower   = 0x13b
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	prime128Shift   = 24
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// New32 returns a new 32-bit FNV-1 [hash.Hash].</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func New32() hash.Hash32 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	var s sum32 = offset32
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	return &amp;s
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// New32a returns a new 32-bit FNV-1a [hash.Hash].</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func New32a() hash.Hash32 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	var s sum32a = offset32
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	return &amp;s
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// New64 returns a new 64-bit FNV-1 [hash.Hash].</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>func New64() hash.Hash64 {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	var s sum64 = offset64
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return &amp;s
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// New64a returns a new 64-bit FNV-1a [hash.Hash].</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func New64a() hash.Hash64 {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	var s sum64a = offset64
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return &amp;s
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// New128 returns a new 128-bit FNV-1 [hash.Hash].</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func New128() hash.Hash {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	var s sum128
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	s[0] = offset128Higher
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	s[1] = offset128Lower
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	return &amp;s
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// New128a returns a new 128-bit FNV-1a [hash.Hash].</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// Its Sum method will lay the value out in big-endian byte order.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func New128a() hash.Hash {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	var s sum128a
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	s[0] = offset128Higher
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	s[1] = offset128Lower
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	return &amp;s
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func (s *sum32) Reset()   { *s = offset32 }
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func (s *sum32a) Reset()  { *s = offset32 }
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>func (s *sum64) Reset()   { *s = offset64 }
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func (s *sum64a) Reset()  { *s = offset64 }
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (s *sum128) Reset()  { s[0] = offset128Higher; s[1] = offset128Lower }
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (s *sum128a) Reset() { s[0] = offset128Higher; s[1] = offset128Lower }
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func (s *sum32) Sum32() uint32  { return uint32(*s) }
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func (s *sum32a) Sum32() uint32 { return uint32(*s) }
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (s *sum64) Sum64() uint64  { return uint64(*s) }
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func (s *sum64a) Sum64() uint64 { return uint64(*s) }
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>func (s *sum32) Write(data []byte) (int, error) {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	hash := *s
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		hash *= prime32
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		hash ^= sum32(c)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	*s = hash
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return len(data), nil
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func (s *sum32a) Write(data []byte) (int, error) {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	hash := *s
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		hash ^= sum32a(c)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		hash *= prime32
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	*s = hash
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	return len(data), nil
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>func (s *sum64) Write(data []byte) (int, error) {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	hash := *s
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		hash *= prime64
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		hash ^= sum64(c)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	*s = hash
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return len(data), nil
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (s *sum64a) Write(data []byte) (int, error) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	hash := *s
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		hash ^= sum64a(c)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		hash *= prime64
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	*s = hash
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	return len(data), nil
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func (s *sum128) Write(data []byte) (int, error) {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		<span class="comment">// Compute the multiplication</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		s0, s1 := bits.Mul64(prime128Lower, s[1])
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		s0 += s[1]&lt;&lt;prime128Shift + prime128Lower*s[0]
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// Update the values</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		s[1] = s1
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		s[0] = s0
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		s[1] ^= uint64(c)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	return len(data), nil
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (s *sum128a) Write(data []byte) (int, error) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	for _, c := range data {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		s[1] ^= uint64(c)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		<span class="comment">// Compute the multiplication</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		s0, s1 := bits.Mul64(prime128Lower, s[1])
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		s0 += s[1]&lt;&lt;prime128Shift + prime128Lower*s[0]
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// Update the values</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		s[1] = s1
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		s[0] = s0
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	return len(data), nil
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (s *sum32) Size() int   { return 4 }
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>func (s *sum32a) Size() int  { return 4 }
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func (s *sum64) Size() int   { return 8 }
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>func (s *sum64a) Size() int  { return 8 }
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>func (s *sum128) Size() int  { return 16 }
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>func (s *sum128a) Size() int { return 16 }
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>func (s *sum32) BlockSize() int   { return 1 }
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>func (s *sum32a) BlockSize() int  { return 1 }
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func (s *sum64) BlockSize() int   { return 1 }
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func (s *sum64a) BlockSize() int  { return 1 }
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>func (s *sum128) BlockSize() int  { return 1 }
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func (s *sum128a) BlockSize() int { return 1 }
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>func (s *sum32) Sum(in []byte) []byte {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	v := uint32(*s)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	return append(in, byte(v&gt;&gt;24), byte(v&gt;&gt;16), byte(v&gt;&gt;8), byte(v))
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func (s *sum32a) Sum(in []byte) []byte {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	v := uint32(*s)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return append(in, byte(v&gt;&gt;24), byte(v&gt;&gt;16), byte(v&gt;&gt;8), byte(v))
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func (s *sum64) Sum(in []byte) []byte {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	v := uint64(*s)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	return append(in, byte(v&gt;&gt;56), byte(v&gt;&gt;48), byte(v&gt;&gt;40), byte(v&gt;&gt;32), byte(v&gt;&gt;24), byte(v&gt;&gt;16), byte(v&gt;&gt;8), byte(v))
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>func (s *sum64a) Sum(in []byte) []byte {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	v := uint64(*s)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	return append(in, byte(v&gt;&gt;56), byte(v&gt;&gt;48), byte(v&gt;&gt;40), byte(v&gt;&gt;32), byte(v&gt;&gt;24), byte(v&gt;&gt;16), byte(v&gt;&gt;8), byte(v))
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>func (s *sum128) Sum(in []byte) []byte {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	return append(in,
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		byte(s[0]&gt;&gt;56), byte(s[0]&gt;&gt;48), byte(s[0]&gt;&gt;40), byte(s[0]&gt;&gt;32), byte(s[0]&gt;&gt;24), byte(s[0]&gt;&gt;16), byte(s[0]&gt;&gt;8), byte(s[0]),
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		byte(s[1]&gt;&gt;56), byte(s[1]&gt;&gt;48), byte(s[1]&gt;&gt;40), byte(s[1]&gt;&gt;32), byte(s[1]&gt;&gt;24), byte(s[1]&gt;&gt;16), byte(s[1]&gt;&gt;8), byte(s[1]),
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func (s *sum128a) Sum(in []byte) []byte {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	return append(in,
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		byte(s[0]&gt;&gt;56), byte(s[0]&gt;&gt;48), byte(s[0]&gt;&gt;40), byte(s[0]&gt;&gt;32), byte(s[0]&gt;&gt;24), byte(s[0]&gt;&gt;16), byte(s[0]&gt;&gt;8), byte(s[0]),
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		byte(s[1]&gt;&gt;56), byte(s[1]&gt;&gt;48), byte(s[1]&gt;&gt;40), byte(s[1]&gt;&gt;32), byte(s[1]&gt;&gt;24), byte(s[1]&gt;&gt;16), byte(s[1]&gt;&gt;8), byte(s[1]),
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>const (
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	magic32          = &#34;fnv\x01&#34;
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	magic32a         = &#34;fnv\x02&#34;
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	magic64          = &#34;fnv\x03&#34;
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	magic64a         = &#34;fnv\x04&#34;
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	magic128         = &#34;fnv\x05&#34;
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	magic128a        = &#34;fnv\x06&#34;
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	marshaledSize32  = len(magic32) + 4
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	marshaledSize64  = len(magic64) + 8
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	marshaledSize128 = len(magic128) + 8*2
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>func (s *sum32) MarshalBinary() ([]byte, error) {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize32)
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	b = append(b, magic32...)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	b = appendUint32(b, uint32(*s))
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	return b, nil
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>func (s *sum32a) MarshalBinary() ([]byte, error) {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize32)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	b = append(b, magic32a...)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	b = appendUint32(b, uint32(*s))
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	return b, nil
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>func (s *sum64) MarshalBinary() ([]byte, error) {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize64)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	b = append(b, magic64...)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	b = appendUint64(b, uint64(*s))
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	return b, nil
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func (s *sum64a) MarshalBinary() ([]byte, error) {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize64)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	b = append(b, magic64a...)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	b = appendUint64(b, uint64(*s))
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	return b, nil
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>func (s *sum128) MarshalBinary() ([]byte, error) {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize128)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	b = append(b, magic128...)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	b = appendUint64(b, s[0])
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	b = appendUint64(b, s[1])
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	return b, nil
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func (s *sum128a) MarshalBinary() ([]byte, error) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	b := make([]byte, 0, marshaledSize128)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	b = append(b, magic128a...)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	b = appendUint64(b, s[0])
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	b = appendUint64(b, s[1])
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	return b, nil
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func (s *sum32) UnmarshalBinary(b []byte) error {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	if len(b) &lt; len(magic32) || string(b[:len(magic32)]) != magic32 {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	if len(b) != marshaledSize32 {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	*s = sum32(readUint32(b[4:]))
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	return nil
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>func (s *sum32a) UnmarshalBinary(b []byte) error {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if len(b) &lt; len(magic32a) || string(b[:len(magic32a)]) != magic32a {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if len(b) != marshaledSize32 {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	*s = sum32a(readUint32(b[4:]))
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	return nil
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>func (s *sum64) UnmarshalBinary(b []byte) error {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if len(b) &lt; len(magic64) || string(b[:len(magic64)]) != magic64 {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	if len(b) != marshaledSize64 {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	*s = sum64(readUint64(b[4:]))
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	return nil
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (s *sum64a) UnmarshalBinary(b []byte) error {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	if len(b) &lt; len(magic64a) || string(b[:len(magic64a)]) != magic64a {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	if len(b) != marshaledSize64 {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	*s = sum64a(readUint64(b[4:]))
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	return nil
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>func (s *sum128) UnmarshalBinary(b []byte) error {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	if len(b) &lt; len(magic128) || string(b[:len(magic128)]) != magic128 {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	if len(b) != marshaledSize128 {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	s[0] = readUint64(b[4:])
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	s[1] = readUint64(b[12:])
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	return nil
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func (s *sum128a) UnmarshalBinary(b []byte) error {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if len(b) &lt; len(magic128a) || string(b[:len(magic128a)]) != magic128a {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state identifier&#34;)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	if len(b) != marshaledSize128 {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return errors.New(&#34;hash/fnv: invalid hash state size&#34;)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	s[0] = readUint64(b[4:])
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	s[1] = readUint64(b[12:])
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	return nil
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// readUint32 is semantically the same as [binary.BigEndian.Uint32]</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>func readUint32(b []byte) uint32 {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	_ = b[3]
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	return uint32(b[3]) | uint32(b[2])&lt;&lt;8 | uint32(b[1])&lt;&lt;16 | uint32(b[0])&lt;&lt;24
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// appendUint32 is semantically the same as [binary.BigEndian.AppendUint32]</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func appendUint32(b []byte, x uint32) []byte {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	return append(b,
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		byte(x&gt;&gt;24),
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		byte(x&gt;&gt;16),
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		byte(x&gt;&gt;8),
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		byte(x),
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// appendUint64 is semantically the same as [binary.BigEndian.AppendUint64]</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func appendUint64(b []byte, x uint64) []byte {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	return append(b,
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		byte(x&gt;&gt;56),
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		byte(x&gt;&gt;48),
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		byte(x&gt;&gt;40),
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		byte(x&gt;&gt;32),
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		byte(x&gt;&gt;24),
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		byte(x&gt;&gt;16),
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		byte(x&gt;&gt;8),
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		byte(x),
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// readUint64 is semantically the same as [binary.BigEndian.Uint64]</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// We copied this function because we can not import &#34;encoding/binary&#34; here.</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>func readUint64(b []byte) uint64 {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	_ = b[7]
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	return uint64(b[7]) | uint64(b[6])&lt;&lt;8 | uint64(b[5])&lt;&lt;16 | uint64(b[4])&lt;&lt;24 |
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		uint64(b[3])&lt;&lt;32 | uint64(b[2])&lt;&lt;40 | uint64(b[1])&lt;&lt;48 | uint64(b[0])&lt;&lt;56
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
</pre><p><a href="fnv.go?m=text">View as plain text</a></p>

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
