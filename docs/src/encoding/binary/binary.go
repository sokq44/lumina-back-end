<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/binary/binary.go - Go Documentation Server</title>

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
<a href="binary.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/binary">binary</a>/<span class="text-muted">binary.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/encoding/binary">encoding/binary</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package binary implements simple translation between numbers and byte</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// sequences and encoding and decoding of varints.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// Numbers are translated by reading and writing fixed-size values.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// A fixed-size value is either a fixed-size arithmetic</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// type (bool, int8, uint8, int16, float32, complex64, ...)</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// or an array or struct containing only fixed-size values.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// The varint functions encode and decode single integer values using</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// a variable-length encoding; smaller values require fewer bytes.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// For a specification, see</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// https://developers.google.com/protocol-buffers/docs/encoding.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// This package favors simplicity over efficiency. Clients that require</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// high-performance serialization, especially for large data structures,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// should look at more advanced solutions such as the [encoding/gob]</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// package or [google.golang.org/protobuf] for protocol buffers.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>package binary
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>import (
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// A ByteOrder specifies how to convert byte slices into</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// 16-, 32-, or 64-bit unsigned integers.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// It is implemented by [LittleEndian], [BigEndian], and [NativeEndian].</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>type ByteOrder interface {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	Uint16([]byte) uint16
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	Uint32([]byte) uint32
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	Uint64([]byte) uint64
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	PutUint16([]byte, uint16)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	PutUint32([]byte, uint32)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	PutUint64([]byte, uint64)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	String() string
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// AppendByteOrder specifies how to append 16-, 32-, or 64-bit unsigned integers</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// into a byte slice.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// It is implemented by [LittleEndian], [BigEndian], and [NativeEndian].</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>type AppendByteOrder interface {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	AppendUint16([]byte, uint16) []byte
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	AppendUint32([]byte, uint32) []byte
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	AppendUint64([]byte, uint64) []byte
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	String() string
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// LittleEndian is the little-endian implementation of [ByteOrder] and [AppendByteOrder].</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>var LittleEndian littleEndian
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// BigEndian is the big-endian implementation of [ByteOrder] and [AppendByteOrder].</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>var BigEndian bigEndian
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>type littleEndian struct{}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func (littleEndian) Uint16(b []byte) uint16 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	_ = b[1] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return uint16(b[0]) | uint16(b[1])&lt;&lt;8
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (littleEndian) PutUint16(b []byte, v uint16) {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	_ = b[1] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	b[0] = byte(v)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	b[1] = byte(v &gt;&gt; 8)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func (littleEndian) AppendUint16(b []byte, v uint16) []byte {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	return append(b,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		byte(v),
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (littleEndian) Uint32(b []byte) uint32 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	_ = b[3] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	return uint32(b[0]) | uint32(b[1])&lt;&lt;8 | uint32(b[2])&lt;&lt;16 | uint32(b[3])&lt;&lt;24
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func (littleEndian) PutUint32(b []byte, v uint32) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	_ = b[3] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	b[0] = byte(v)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	b[1] = byte(v &gt;&gt; 8)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	b[2] = byte(v &gt;&gt; 16)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	b[3] = byte(v &gt;&gt; 24)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (littleEndian) AppendUint32(b []byte, v uint32) []byte {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	return append(b,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		byte(v),
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		byte(v&gt;&gt;16),
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		byte(v&gt;&gt;24),
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (littleEndian) Uint64(b []byte) uint64 {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	_ = b[7] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	return uint64(b[0]) | uint64(b[1])&lt;&lt;8 | uint64(b[2])&lt;&lt;16 | uint64(b[3])&lt;&lt;24 |
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		uint64(b[4])&lt;&lt;32 | uint64(b[5])&lt;&lt;40 | uint64(b[6])&lt;&lt;48 | uint64(b[7])&lt;&lt;56
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func (littleEndian) PutUint64(b []byte, v uint64) {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	_ = b[7] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	b[0] = byte(v)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	b[1] = byte(v &gt;&gt; 8)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	b[2] = byte(v &gt;&gt; 16)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	b[3] = byte(v &gt;&gt; 24)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	b[4] = byte(v &gt;&gt; 32)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	b[5] = byte(v &gt;&gt; 40)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	b[6] = byte(v &gt;&gt; 48)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	b[7] = byte(v &gt;&gt; 56)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func (littleEndian) AppendUint64(b []byte, v uint64) []byte {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	return append(b,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		byte(v),
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		byte(v&gt;&gt;16),
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		byte(v&gt;&gt;24),
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		byte(v&gt;&gt;32),
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		byte(v&gt;&gt;40),
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		byte(v&gt;&gt;48),
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		byte(v&gt;&gt;56),
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (littleEndian) String() string { return &#34;LittleEndian&#34; }
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func (littleEndian) GoString() string { return &#34;binary.LittleEndian&#34; }
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>type bigEndian struct{}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func (bigEndian) Uint16(b []byte) uint16 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	_ = b[1] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return uint16(b[1]) | uint16(b[0])&lt;&lt;8
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>func (bigEndian) PutUint16(b []byte, v uint16) {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	_ = b[1] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	b[0] = byte(v &gt;&gt; 8)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	b[1] = byte(v)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>func (bigEndian) AppendUint16(b []byte, v uint16) []byte {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	return append(b,
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		byte(v),
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	)
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (bigEndian) Uint32(b []byte) uint32 {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	_ = b[3] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	return uint32(b[3]) | uint32(b[2])&lt;&lt;8 | uint32(b[1])&lt;&lt;16 | uint32(b[0])&lt;&lt;24
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func (bigEndian) PutUint32(b []byte, v uint32) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	_ = b[3] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	b[0] = byte(v &gt;&gt; 24)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	b[1] = byte(v &gt;&gt; 16)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	b[2] = byte(v &gt;&gt; 8)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	b[3] = byte(v)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>func (bigEndian) AppendUint32(b []byte, v uint32) []byte {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	return append(b,
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		byte(v&gt;&gt;24),
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		byte(v&gt;&gt;16),
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		byte(v),
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>func (bigEndian) Uint64(b []byte) uint64 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	_ = b[7] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return uint64(b[7]) | uint64(b[6])&lt;&lt;8 | uint64(b[5])&lt;&lt;16 | uint64(b[4])&lt;&lt;24 |
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		uint64(b[3])&lt;&lt;32 | uint64(b[2])&lt;&lt;40 | uint64(b[1])&lt;&lt;48 | uint64(b[0])&lt;&lt;56
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>func (bigEndian) PutUint64(b []byte, v uint64) {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	_ = b[7] <span class="comment">// early bounds check to guarantee safety of writes below</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	b[0] = byte(v &gt;&gt; 56)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	b[1] = byte(v &gt;&gt; 48)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	b[2] = byte(v &gt;&gt; 40)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	b[3] = byte(v &gt;&gt; 32)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	b[4] = byte(v &gt;&gt; 24)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	b[5] = byte(v &gt;&gt; 16)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	b[6] = byte(v &gt;&gt; 8)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	b[7] = byte(v)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (bigEndian) AppendUint64(b []byte, v uint64) []byte {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	return append(b,
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		byte(v&gt;&gt;56),
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		byte(v&gt;&gt;48),
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		byte(v&gt;&gt;40),
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		byte(v&gt;&gt;32),
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		byte(v&gt;&gt;24),
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		byte(v&gt;&gt;16),
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		byte(v&gt;&gt;8),
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		byte(v),
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func (bigEndian) String() string { return &#34;BigEndian&#34; }
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func (bigEndian) GoString() string { return &#34;binary.BigEndian&#34; }
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func (nativeEndian) String() string { return &#34;NativeEndian&#34; }
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (nativeEndian) GoString() string { return &#34;binary.NativeEndian&#34; }
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// Read reads structured binary data from r into data.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// Data must be a pointer to a fixed-size value or a slice</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// of fixed-size values.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// Bytes read from r are decoded using the specified byte order</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// and written to successive fields of the data.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// When decoding boolean values, a zero byte is decoded as false, and</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// any other non-zero byte is decoded as true.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// When reading into structs, the field data for fields with</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// blank (_) field names is skipped; i.e., blank field names</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">// may be used for padding.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// When reading into a struct, all non-blank fields must be exported</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// or Read may panic.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">// The error is [io.EOF] only if no bytes were read.</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// If an [io.EOF] happens after reading some but not all the bytes,</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// Read returns [io.ErrUnexpectedEOF].</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>func Read(r io.Reader, order ByteOrder, data any) error {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// Fast path for basic types and slices.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	if n := intDataSize(data); n != 0 {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		bs := make([]byte, n)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		if _, err := io.ReadFull(r, bs); err != nil {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			return err
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		switch data := data.(type) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		case *bool:
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			*data = bs[0] != 0
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		case *int8:
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			*data = int8(bs[0])
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		case *uint8:
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			*data = bs[0]
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		case *int16:
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			*data = int16(order.Uint16(bs))
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		case *uint16:
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			*data = order.Uint16(bs)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		case *int32:
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			*data = int32(order.Uint32(bs))
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		case *uint32:
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			*data = order.Uint32(bs)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		case *int64:
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			*data = int64(order.Uint64(bs))
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		case *uint64:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			*data = order.Uint64(bs)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		case *float32:
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			*data = math.Float32frombits(order.Uint32(bs))
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		case *float64:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			*data = math.Float64frombits(order.Uint64(bs))
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		case []bool:
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			for i, x := range bs { <span class="comment">// Easier to loop over the input for 8-bit values.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				data[i] = x != 0
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		case []int8:
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			for i, x := range bs {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				data[i] = int8(x)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		case []uint8:
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			copy(data, bs)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		case []int16:
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			for i := range data {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>				data[i] = int16(order.Uint16(bs[2*i:]))
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		case []uint16:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			for i := range data {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				data[i] = order.Uint16(bs[2*i:])
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		case []int32:
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			for i := range data {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				data[i] = int32(order.Uint32(bs[4*i:]))
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		case []uint32:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			for i := range data {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				data[i] = order.Uint32(bs[4*i:])
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		case []int64:
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			for i := range data {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				data[i] = int64(order.Uint64(bs[8*i:]))
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		case []uint64:
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			for i := range data {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				data[i] = order.Uint64(bs[8*i:])
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		case []float32:
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			for i := range data {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				data[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		case []float64:
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			for i := range data {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				data[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		default:
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			n = 0 <span class="comment">// fast path doesn&#39;t apply</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		}
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		if n != 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			return nil
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// Fallback to reflect-based decoding.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	v := reflect.ValueOf(data)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	size := -1
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	switch v.Kind() {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	case reflect.Pointer:
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		v = v.Elem()
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		size = dataSize(v)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		size = dataSize(v)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if size &lt; 0 {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		return errors.New(&#34;binary.Read: invalid type &#34; + reflect.TypeOf(data).String())
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	d := &amp;decoder{order: order, buf: make([]byte, size)}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	if _, err := io.ReadFull(r, d.buf); err != nil {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		return err
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	d.value(v)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	return nil
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// Write writes the binary representation of data into w.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// Data must be a fixed-size value or a slice of fixed-size</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// values, or a pointer to such data.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// Boolean values encode as one byte: 1 for true, and 0 for false.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// Bytes written to w are encoded using the specified byte order</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// and read from successive fields of the data.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// When writing structs, zero values are written for fields</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// with blank (_) field names.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func Write(w io.Writer, order ByteOrder, data any) error {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	<span class="comment">// Fast path for basic types and slices.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	if n := intDataSize(data); n != 0 {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		bs := make([]byte, n)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		switch v := data.(type) {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		case *bool:
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			if *v {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>				bs[0] = 1
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			} else {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>				bs[0] = 0
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		case bool:
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			if v {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>				bs[0] = 1
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			} else {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>				bs[0] = 0
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		case []bool:
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				if x {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>					bs[i] = 1
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>				} else {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>					bs[i] = 0
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>				}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		case *int8:
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			bs[0] = byte(*v)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		case int8:
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			bs[0] = byte(v)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		case []int8:
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>				bs[i] = byte(x)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		case *uint8:
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			bs[0] = *v
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		case uint8:
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			bs[0] = v
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		case []uint8:
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			bs = v
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		case *int16:
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			order.PutUint16(bs, uint16(*v))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		case int16:
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			order.PutUint16(bs, uint16(v))
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		case []int16:
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>				order.PutUint16(bs[2*i:], uint16(x))
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		case *uint16:
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			order.PutUint16(bs, *v)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		case uint16:
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			order.PutUint16(bs, v)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		case []uint16:
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				order.PutUint16(bs[2*i:], x)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		case *int32:
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			order.PutUint32(bs, uint32(*v))
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		case int32:
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			order.PutUint32(bs, uint32(v))
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		case []int32:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>				order.PutUint32(bs[4*i:], uint32(x))
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		case *uint32:
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			order.PutUint32(bs, *v)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		case uint32:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			order.PutUint32(bs, v)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		case []uint32:
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				order.PutUint32(bs[4*i:], x)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		case *int64:
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			order.PutUint64(bs, uint64(*v))
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		case int64:
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			order.PutUint64(bs, uint64(v))
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		case []int64:
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>				order.PutUint64(bs[8*i:], uint64(x))
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		case *uint64:
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			order.PutUint64(bs, *v)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		case uint64:
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			order.PutUint64(bs, v)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		case []uint64:
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>				order.PutUint64(bs[8*i:], x)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		case *float32:
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			order.PutUint32(bs, math.Float32bits(*v))
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		case float32:
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			order.PutUint32(bs, math.Float32bits(v))
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		case []float32:
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				order.PutUint32(bs[4*i:], math.Float32bits(x))
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		case *float64:
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			order.PutUint64(bs, math.Float64bits(*v))
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		case float64:
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			order.PutUint64(bs, math.Float64bits(v))
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		case []float64:
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			for i, x := range v {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				order.PutUint64(bs[8*i:], math.Float64bits(x))
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		_, err := w.Write(bs)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return err
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// Fallback to reflect-based encoding.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	v := reflect.Indirect(reflect.ValueOf(data))
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	size := dataSize(v)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	if size &lt; 0 {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		return errors.New(&#34;binary.Write: some values are not fixed-sized in type &#34; + reflect.TypeOf(data).String())
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	buf := make([]byte, size)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	e := &amp;encoder{order: order, buf: buf}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	e.value(v)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	_, err := w.Write(buf)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	return err
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">// Size returns how many bytes [Write] would generate to encode the value v, which</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">// must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// If v is neither of these, Size returns -1.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>func Size(v any) int {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	return dataSize(reflect.Indirect(reflect.ValueOf(v)))
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>var structSize sync.Map <span class="comment">// map[reflect.Type]int</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">// dataSize returns the number of bytes the actual data represented by v occupies in memory.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// it returns the length of the slice times the element size and does not count the memory</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// occupied by the header. If the type of v is not acceptable, dataSize returns -1.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>func dataSize(v reflect.Value) int {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	switch v.Kind() {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		if s := sizeof(v.Type().Elem()); s &gt;= 0 {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			return s * v.Len()
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		t := v.Type()
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		if size, ok := structSize.Load(t); ok {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			return size.(int)
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		size := sizeof(t)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		structSize.Store(t, size)
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		return size
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	default:
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		if v.IsValid() {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			return sizeof(v.Type())
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	return -1
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">// sizeof returns the size &gt;= 0 of variables for the given type or -1 if the type is not acceptable.</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>func sizeof(t reflect.Type) int {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	switch t.Kind() {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	case reflect.Array:
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		if s := sizeof(t.Elem()); s &gt;= 0 {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			return s * t.Len()
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		sum := 0
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		for i, n := 0, t.NumField(); i &lt; n; i++ {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			s := sizeof(t.Field(i).Type)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			if s &lt; 0 {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>				return -1
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			sum += s
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		return sum
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	case reflect.Bool,
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		return int(t.Size())
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	return -1
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>type coder struct {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	order  ByteOrder
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	buf    []byte
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	offset int
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>}
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>type decoder coder
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>type encoder coder
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>func (d *decoder) bool() bool {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	x := d.buf[d.offset]
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	d.offset++
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	return x != 0
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>func (e *encoder) bool(x bool) {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	if x {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		e.buf[e.offset] = 1
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	} else {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		e.buf[e.offset] = 0
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	e.offset++
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>func (d *decoder) uint8() uint8 {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	x := d.buf[d.offset]
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	d.offset++
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	return x
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>func (e *encoder) uint8(x uint8) {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	e.buf[e.offset] = x
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	e.offset++
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>func (d *decoder) uint16() uint16 {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	x := d.order.Uint16(d.buf[d.offset : d.offset+2])
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	d.offset += 2
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	return x
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>func (e *encoder) uint16(x uint16) {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	e.order.PutUint16(e.buf[e.offset:e.offset+2], x)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	e.offset += 2
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>func (d *decoder) uint32() uint32 {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	d.offset += 4
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	return x
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>func (e *encoder) uint32(x uint32) {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	e.order.PutUint32(e.buf[e.offset:e.offset+4], x)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	e.offset += 4
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>func (d *decoder) uint64() uint64 {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	x := d.order.Uint64(d.buf[d.offset : d.offset+8])
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	d.offset += 8
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	return x
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>}
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>func (e *encoder) uint64(x uint64) {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	e.order.PutUint64(e.buf[e.offset:e.offset+8], x)
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	e.offset += 8
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>func (d *decoder) int8() int8 { return int8(d.uint8()) }
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>func (d *decoder) int16() int16 { return int16(d.uint16()) }
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>func (d *decoder) int32() int32 { return int32(d.uint32()) }
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>func (d *decoder) int64() int64 { return int64(d.uint64()) }
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>func (d *decoder) value(v reflect.Value) {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	switch v.Kind() {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	case reflect.Array:
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		l := v.Len()
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>			d.value(v.Index(i))
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		t := v.Type()
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		l := v.NumField()
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			<span class="comment">// Note: Calling v.CanSet() below is an optimization.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>			<span class="comment">// It would be sufficient to check the field name,</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			<span class="comment">// but creating the StructField info for each field is</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			<span class="comment">// costly (run &#34;go test -bench=ReadStruct&#34; and compare</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			<span class="comment">// results when making changes to this code).</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			if v := v.Field(i); v.CanSet() || t.Field(i).Name != &#34;_&#34; {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>				d.value(v)
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			} else {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>				d.skip(v)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		l := v.Len()
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			d.value(v.Index(i))
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		v.SetBool(d.bool())
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	case reflect.Int8:
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		v.SetInt(int64(d.int8()))
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	case reflect.Int16:
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		v.SetInt(int64(d.int16()))
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	case reflect.Int32:
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		v.SetInt(int64(d.int32()))
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	case reflect.Int64:
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		v.SetInt(d.int64())
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	case reflect.Uint8:
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		v.SetUint(uint64(d.uint8()))
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	case reflect.Uint16:
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		v.SetUint(uint64(d.uint16()))
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	case reflect.Uint32:
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		v.SetUint(uint64(d.uint32()))
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	case reflect.Uint64:
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		v.SetUint(d.uint64())
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	case reflect.Float32:
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		v.SetFloat(float64(math.Float32frombits(d.uint32())))
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	case reflect.Float64:
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		v.SetFloat(math.Float64frombits(d.uint64()))
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	case reflect.Complex64:
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		v.SetComplex(complex(
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			float64(math.Float32frombits(d.uint32())),
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>			float64(math.Float32frombits(d.uint32())),
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		))
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	case reflect.Complex128:
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>		v.SetComplex(complex(
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>			math.Float64frombits(d.uint64()),
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			math.Float64frombits(d.uint64()),
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		))
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>func (e *encoder) value(v reflect.Value) {
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	switch v.Kind() {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	case reflect.Array:
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		l := v.Len()
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>			e.value(v.Index(i))
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		}
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		t := v.Type()
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		l := v.NumField()
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			<span class="comment">// see comment for corresponding code in decoder.value()</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			if v := v.Field(i); v.CanSet() || t.Field(i).Name != &#34;_&#34; {
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>				e.value(v)
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			} else {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>				e.skip(v)
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		l := v.Len()
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		for i := 0; i &lt; l; i++ {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>			e.value(v.Index(i))
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		e.bool(v.Bool())
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		switch v.Type().Kind() {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		case reflect.Int8:
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>			e.int8(int8(v.Int()))
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		case reflect.Int16:
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			e.int16(int16(v.Int()))
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		case reflect.Int32:
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			e.int32(int32(v.Int()))
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		case reflect.Int64:
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			e.int64(v.Int())
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		switch v.Type().Kind() {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		case reflect.Uint8:
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			e.uint8(uint8(v.Uint()))
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		case reflect.Uint16:
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			e.uint16(uint16(v.Uint()))
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		case reflect.Uint32:
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			e.uint32(uint32(v.Uint()))
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		case reflect.Uint64:
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			e.uint64(v.Uint())
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	case reflect.Float32, reflect.Float64:
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		switch v.Type().Kind() {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		case reflect.Float32:
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			e.uint32(math.Float32bits(float32(v.Float())))
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		case reflect.Float64:
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			e.uint64(math.Float64bits(v.Float()))
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	case reflect.Complex64, reflect.Complex128:
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		switch v.Type().Kind() {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		case reflect.Complex64:
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>			x := v.Complex()
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			e.uint32(math.Float32bits(float32(real(x))))
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>			e.uint32(math.Float32bits(float32(imag(x))))
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		case reflect.Complex128:
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			x := v.Complex()
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			e.uint64(math.Float64bits(real(x)))
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>			e.uint64(math.Float64bits(imag(x)))
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		}
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	}
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>func (d *decoder) skip(v reflect.Value) {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	d.offset += dataSize(v)
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>func (e *encoder) skip(v reflect.Value) {
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	n := dataSize(v)
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	zero := e.buf[e.offset : e.offset+n]
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	for i := range zero {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		zero[i] = 0
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	e.offset += n
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>}
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>
<span id="L775" class="ln">   775&nbsp;&nbsp;</span><span class="comment">// intDataSize returns the size of the data required to represent the data when encoded.</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span><span class="comment">// It returns zero if the type cannot be implemented by the fast path in Read or Write.</span>
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>func intDataSize(data any) int {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	switch data := data.(type) {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	case bool, int8, uint8, *bool, *int8, *uint8:
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>		return 1
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	case []bool:
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>		return len(data)
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	case []int8:
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		return len(data)
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	case []uint8:
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		return len(data)
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	case int16, uint16, *int16, *uint16:
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		return 2
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	case []int16:
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		return 2 * len(data)
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	case []uint16:
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		return 2 * len(data)
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	case int32, uint32, *int32, *uint32:
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		return 4
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	case []int32:
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		return 4 * len(data)
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	case []uint32:
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		return 4 * len(data)
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	case int64, uint64, *int64, *uint64:
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		return 8
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	case []int64:
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		return 8 * len(data)
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	case []uint64:
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		return 8 * len(data)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	case float32, *float32:
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		return 4
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	case float64, *float64:
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		return 8
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	case []float32:
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		return 4 * len(data)
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	case []float64:
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		return 8 * len(data)
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	}
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	return 0
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>}
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
</pre><p><a href="binary.go?m=text">View as plain text</a></p>

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
