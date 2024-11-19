<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/fmt/format.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../index.html">GoDoc</a></div>
<a href="format.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/fmt">fmt</a>/<span class="text-muted">format.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/fmt">fmt</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package fmt
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>const (
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	ldigits = &#34;0123456789abcdefx&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	udigits = &#34;0123456789ABCDEFX&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	signed   = true
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	unsigned = false
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// flags placed in a separate struct for easy clearing.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>type fmtFlags struct {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	widPresent  bool
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	precPresent bool
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	minus       bool
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	plus        bool
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	sharp       bool
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	space       bool
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	zero        bool
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// For the formats %+v %#v, we set the plusV/sharpV flags</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// and clear the plus/sharp flags since %+v and %#v are in effect</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// different, flagless formats set at the top level.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	plusV  bool
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	sharpV bool
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// A fmt is the raw formatter used by Printf etc.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// It prints into a buffer that must be set up separately.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>type fmt struct {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	buf *buffer
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	fmtFlags
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	wid  int <span class="comment">// width</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	prec int <span class="comment">// precision</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// intbuf is large enough to store %b of an int64 with a sign and</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// avoids padding at the end of the struct on 32 bit architectures.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	intbuf [68]byte
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (f *fmt) clearflags() {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	f.fmtFlags = fmtFlags{}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func (f *fmt) init(buf *buffer) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	f.buf = buf
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	f.clearflags()
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// writePadding generates n bytes of padding.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (f *fmt) writePadding(n int) {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if n &lt;= 0 { <span class="comment">// No padding bytes needed.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	buf := *f.buf
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	oldLen := len(buf)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	newLen := oldLen + n
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// Make enough room for padding.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if newLen &gt; cap(buf) {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		buf = make(buffer, cap(buf)*2+n)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		copy(buf, *f.buf)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// Decide which byte the padding should be filled with.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	padByte := byte(&#39; &#39;)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if f.zero {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		padByte = byte(&#39;0&#39;)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// Fill padding with padByte.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	padding := buf[oldLen:newLen]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	for i := range padding {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		padding[i] = padByte
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	*f.buf = buf[:newLen]
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// pad appends b to f.buf, padded on left (!f.minus) or right (f.minus).</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>func (f *fmt) pad(b []byte) {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if !f.widPresent || f.wid == 0 {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		f.buf.write(b)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	width := f.wid - utf8.RuneCount(b)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if !f.minus {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		<span class="comment">// left padding</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		f.writePadding(width)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		f.buf.write(b)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	} else {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		<span class="comment">// right padding</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		f.buf.write(b)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		f.writePadding(width)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// padString appends s to f.buf, padded on left (!f.minus) or right (f.minus).</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func (f *fmt) padString(s string) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	if !f.widPresent || f.wid == 0 {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		f.buf.writeString(s)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	width := f.wid - utf8.RuneCountInString(s)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if !f.minus {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		<span class="comment">// left padding</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		f.writePadding(width)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		f.buf.writeString(s)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	} else {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		<span class="comment">// right padding</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		f.buf.writeString(s)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		f.writePadding(width)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// fmtBoolean formats a boolean.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func (f *fmt) fmtBoolean(v bool) {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if v {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		f.padString(&#34;true&#34;)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	} else {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		f.padString(&#34;false&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// fmtUnicode formats a uint64 as &#34;U+0078&#34; or with f.sharp set as &#34;U+0078 &#39;x&#39;&#34;.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func (f *fmt) fmtUnicode(u uint64) {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	buf := f.intbuf[0:]
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// With default precision set the maximum needed buf length is 18</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// for formatting -1 with %#U (&#34;U+FFFFFFFFFFFFFFFF&#34;) which fits</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// into the already allocated intbuf with a capacity of 68 bytes.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	prec := 4
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	if f.precPresent &amp;&amp; f.prec &gt; 4 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		prec = f.prec
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// Compute space needed for &#34;U+&#34; , number, &#34; &#39;&#34;, character, &#34;&#39;&#34;.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		width := 2 + prec + 2 + utf8.UTFMax + 1
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if width &gt; len(buf) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			buf = make([]byte, width)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Format into buf, ending at buf[i]. Formatting numbers is easier right-to-left.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	i := len(buf)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// For %#U we want to add a space and a quoted character at the end of the buffer.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if f.sharp &amp;&amp; u &lt;= utf8.MaxRune &amp;&amp; strconv.IsPrint(rune(u)) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		i--
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		buf[i] = &#39;\&#39;&#39;
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		i -= utf8.RuneLen(rune(u))
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		utf8.EncodeRune(buf[i:], rune(u))
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		i--
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		buf[i] = &#39;\&#39;&#39;
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		i--
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		buf[i] = &#39; &#39;
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// Format the Unicode code point u as a hexadecimal number.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	for u &gt;= 16 {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		i--
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		buf[i] = udigits[u&amp;0xF]
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		prec--
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		u &gt;&gt;= 4
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	i--
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	buf[i] = udigits[u]
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	prec--
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// Add zeros in front of the number until requested precision is reached.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	for prec &gt; 0 {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		i--
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		buf[i] = &#39;0&#39;
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		prec--
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Add a leading &#34;U+&#34;.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	i--
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	buf[i] = &#39;+&#39;
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	i--
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	buf[i] = &#39;U&#39;
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	oldZero := f.zero
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	f.zero = false
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	f.pad(buf[i:])
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	f.zero = oldZero
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// fmtInteger formats signed and unsigned integers.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>func (f *fmt) fmtInteger(u uint64, base int, isSigned bool, verb rune, digits string) {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	negative := isSigned &amp;&amp; int64(u) &lt; 0
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if negative {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		u = -u
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	buf := f.intbuf[0:]
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// The already allocated f.intbuf with a capacity of 68 bytes</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// is large enough for integer formatting when no precision or width is set.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if f.widPresent || f.precPresent {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		<span class="comment">// Account 3 extra bytes for possible addition of a sign and &#34;0x&#34;.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		width := 3 + f.wid + f.prec <span class="comment">// wid and prec are always positive.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if width &gt; len(buf) {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re going to need a bigger boat.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			buf = make([]byte, width)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Two ways to ask for extra leading zero digits: %.3d or %03d.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// If both are specified the f.zero flag is ignored and</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// padding with spaces is used instead.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	prec := 0
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if f.precPresent {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		prec = f.prec
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// Precision of 0 and value of 0 means &#34;print nothing&#34; but padding.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if prec == 0 &amp;&amp; u == 0 {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			oldZero := f.zero
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			f.zero = false
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			f.writePadding(f.wid)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			f.zero = oldZero
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			return
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	} else if f.zero &amp;&amp; f.widPresent {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		prec = f.wid
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		if negative || f.plus || f.space {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			prec-- <span class="comment">// leave room for sign</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// Because printing is easier right-to-left: format u into buf, ending at buf[i].</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">// We could make things marginally faster by splitting the 32-bit case out</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// into a separate block but it&#39;s not worth the duplication, so u has 64 bits.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	i := len(buf)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Use constants for the division and modulo for more efficient code.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// Switch cases ordered by popularity.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	switch base {
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	case 10:
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		for u &gt;= 10 {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			i--
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			next := u / 10
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			buf[i] = byte(&#39;0&#39; + u - next*10)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			u = next
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	case 16:
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		for u &gt;= 16 {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			i--
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			buf[i] = digits[u&amp;0xF]
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			u &gt;&gt;= 4
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	case 8:
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		for u &gt;= 8 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			i--
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			buf[i] = byte(&#39;0&#39; + u&amp;7)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			u &gt;&gt;= 3
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	case 2:
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		for u &gt;= 2 {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			i--
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			buf[i] = byte(&#39;0&#39; + u&amp;1)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			u &gt;&gt;= 1
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	default:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		panic(&#34;fmt: unknown base; can&#39;t happen&#34;)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	i--
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	buf[i] = digits[u]
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	for i &gt; 0 &amp;&amp; prec &gt; len(buf)-i {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		i--
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		buf[i] = &#39;0&#39;
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// Various prefixes: 0x, -, etc.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if f.sharp {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		switch base {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		case 2:
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			<span class="comment">// Add a leading 0b.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			i--
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			buf[i] = &#39;b&#39;
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			i--
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			buf[i] = &#39;0&#39;
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		case 8:
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			if buf[i] != &#39;0&#39; {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				i--
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>				buf[i] = &#39;0&#39;
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		case 16:
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			<span class="comment">// Add a leading 0x or 0X.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			i--
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			buf[i] = digits[16]
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			i--
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			buf[i] = &#39;0&#39;
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if verb == &#39;O&#39; {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		i--
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		buf[i] = &#39;o&#39;
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		i--
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		buf[i] = &#39;0&#39;
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if negative {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		i--
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		buf[i] = &#39;-&#39;
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	} else if f.plus {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		i--
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		buf[i] = &#39;+&#39;
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	} else if f.space {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		i--
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		buf[i] = &#39; &#39;
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// Left padding with zeros has already been handled like precision earlier</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// or the f.zero flag is ignored due to an explicitly set precision.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	oldZero := f.zero
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	f.zero = false
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	f.pad(buf[i:])
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	f.zero = oldZero
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// truncateString truncates the string s to the specified precision, if present.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>func (f *fmt) truncateString(s string) string {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	if f.precPresent {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		n := f.prec
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		for i := range s {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			n--
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			if n &lt; 0 {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				return s[:i]
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	return s
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// truncate truncates the byte slice b as a string of the specified precision, if present.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>func (f *fmt) truncate(b []byte) []byte {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if f.precPresent {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		n := f.prec
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		for i := 0; i &lt; len(b); {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			n--
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			if n &lt; 0 {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				return b[:i]
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			wid := 1
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			if b[i] &gt;= utf8.RuneSelf {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>				_, wid = utf8.DecodeRune(b[i:])
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>			i += wid
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	return b
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// fmtS formats a string.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func (f *fmt) fmtS(s string) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	s = f.truncateString(s)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	f.padString(s)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span><span class="comment">// fmtBs formats the byte slice b as if it was formatted as string with fmtS.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>func (f *fmt) fmtBs(b []byte) {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	b = f.truncate(b)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	f.pad(b)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// fmtSbx formats a string or byte slice as a hexadecimal encoding of its bytes.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (f *fmt) fmtSbx(s string, b []byte, digits string) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	length := len(b)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	if b == nil {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		<span class="comment">// No byte slice present. Assume string s should be encoded.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		length = len(s)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// Set length to not process more bytes than the precision demands.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	if f.precPresent &amp;&amp; f.prec &lt; length {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		length = f.prec
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// Compute width of the encoding taking into account the f.sharp and f.space flag.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	width := 2 * length
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	if width &gt; 0 {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		if f.space {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			<span class="comment">// Each element encoded by two hexadecimals will get a leading 0x or 0X.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			if f.sharp {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				width *= 2
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>			<span class="comment">// Elements will be separated by a space.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			width += length - 1
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		} else if f.sharp {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			<span class="comment">// Only a leading 0x or 0X will be added for the whole string.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			width += 2
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	} else { <span class="comment">// The byte slice or string that should be encoded is empty.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		if f.widPresent {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			f.writePadding(f.wid)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		return
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	<span class="comment">// Handle padding to the left.</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if f.widPresent &amp;&amp; f.wid &gt; width &amp;&amp; !f.minus {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		f.writePadding(f.wid - width)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">// Write the encoding directly into the output buffer.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	buf := *f.buf
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	if f.sharp {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		<span class="comment">// Add leading 0x or 0X.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		buf = append(buf, &#39;0&#39;, digits[16])
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	var c byte
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	for i := 0; i &lt; length; i++ {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		if f.space &amp;&amp; i &gt; 0 {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			<span class="comment">// Separate elements with a space.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			buf = append(buf, &#39; &#39;)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			if f.sharp {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				<span class="comment">// Add leading 0x or 0X for each element.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>				buf = append(buf, &#39;0&#39;, digits[16])
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		if b != nil {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			c = b[i] <span class="comment">// Take a byte from the input byte slice.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		} else {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			c = s[i] <span class="comment">// Take a byte from the input string.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		<span class="comment">// Encode each byte as two hexadecimal digits.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		buf = append(buf, digits[c&gt;&gt;4], digits[c&amp;0xF])
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	*f.buf = buf
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	<span class="comment">// Handle padding to the right.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	if f.widPresent &amp;&amp; f.wid &gt; width &amp;&amp; f.minus {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		f.writePadding(f.wid - width)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">// fmtSx formats a string as a hexadecimal encoding of its bytes.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>func (f *fmt) fmtSx(s, digits string) {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	f.fmtSbx(s, nil, digits)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">// fmtBx formats a byte slice as a hexadecimal encoding of its bytes.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>func (f *fmt) fmtBx(b []byte, digits string) {
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	f.fmtSbx(&#34;&#34;, b, digits)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">// fmtQ formats a string as a double-quoted, escaped Go string constant.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// If f.sharp is set a raw (backquoted) string may be returned instead</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// if the string does not contain any control characters other than tab.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func (f *fmt) fmtQ(s string) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	s = f.truncateString(s)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if f.sharp &amp;&amp; strconv.CanBackquote(s) {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		f.padString(&#34;`&#34; + s + &#34;`&#34;)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	buf := f.intbuf[:0]
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	if f.plus {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		f.pad(strconv.AppendQuoteToASCII(buf, s))
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	} else {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		f.pad(strconv.AppendQuote(buf, s))
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span><span class="comment">// fmtC formats an integer as a Unicode character.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span><span class="comment">// If the character is not valid Unicode, it will print &#39;\ufffd&#39;.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>func (f *fmt) fmtC(c uint64) {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// Explicitly check whether c exceeds utf8.MaxRune since the conversion</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">// of a uint64 to a rune may lose precision that indicates an overflow.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	r := rune(c)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if c &gt; utf8.MaxRune {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		r = utf8.RuneError
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	buf := f.intbuf[:0]
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	f.pad(utf8.AppendRune(buf, r))
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// fmtQc formats an integer as a single-quoted, escaped Go character constant.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">// If the character is not valid Unicode, it will print &#39;\ufffd&#39;.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>func (f *fmt) fmtQc(c uint64) {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	r := rune(c)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	if c &gt; utf8.MaxRune {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		r = utf8.RuneError
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	buf := f.intbuf[:0]
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	if f.plus {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		f.pad(strconv.AppendQuoteRuneToASCII(buf, r))
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	} else {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		f.pad(strconv.AppendQuoteRune(buf, r))
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">// fmtFloat formats a float64. It assumes that verb is a valid format specifier</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span><span class="comment">// for strconv.AppendFloat and therefore fits into a byte.</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>func (f *fmt) fmtFloat(v float64, size int, verb rune, prec int) {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// Explicit precision in format specifier overrules default precision.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	if f.precPresent {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		prec = f.prec
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">// Format number, reserving space for leading + sign if needed.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	num := strconv.AppendFloat(f.intbuf[:1], v, byte(verb), prec, size)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	if num[1] == &#39;-&#39; || num[1] == &#39;+&#39; {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		num = num[1:]
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	} else {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		num[0] = &#39;+&#39;
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	<span class="comment">// f.space means to add a leading space instead of a &#34;+&#34; sign unless</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	<span class="comment">// the sign is explicitly asked for by f.plus.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if f.space &amp;&amp; num[0] == &#39;+&#39; &amp;&amp; !f.plus {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		num[0] = &#39; &#39;
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">// Special handling for infinities and NaN,</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	<span class="comment">// which don&#39;t look like a number so shouldn&#39;t be padded with zeros.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if num[1] == &#39;I&#39; || num[1] == &#39;N&#39; {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		oldZero := f.zero
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		f.zero = false
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		<span class="comment">// Remove sign before NaN if not asked for.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		if num[1] == &#39;N&#39; &amp;&amp; !f.space &amp;&amp; !f.plus {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			num = num[1:]
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		f.pad(num)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		f.zero = oldZero
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		return
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	<span class="comment">// The sharp flag forces printing a decimal point for non-binary formats</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	<span class="comment">// and retains trailing zeros, which we may need to restore.</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if f.sharp &amp;&amp; verb != &#39;b&#39; {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		digits := 0
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		switch verb {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		case &#39;v&#39;, &#39;g&#39;, &#39;G&#39;, &#39;x&#39;:
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			digits = prec
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			<span class="comment">// If no precision is set explicitly use a precision of 6.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			if digits == -1 {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				digits = 6
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		<span class="comment">// Buffer pre-allocated with enough room for</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		<span class="comment">// exponent notations of the form &#34;e+123&#34; or &#34;p-1023&#34;.</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		var tailBuf [6]byte
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		tail := tailBuf[:0]
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		hasDecimalPoint := false
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		sawNonzeroDigit := false
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		<span class="comment">// Starting from i = 1 to skip sign at num[0].</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		for i := 1; i &lt; len(num); i++ {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>			switch num[i] {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			case &#39;.&#39;:
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				hasDecimalPoint = true
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			case &#39;p&#39;, &#39;P&#39;:
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>				tail = append(tail, num[i:]...)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>				num = num[:i]
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			case &#39;e&#39;, &#39;E&#39;:
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>				if verb != &#39;x&#39; &amp;&amp; verb != &#39;X&#39; {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>					tail = append(tail, num[i:]...)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>					num = num[:i]
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>					break
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				fallthrough
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			default:
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>				if num[i] != &#39;0&#39; {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>					sawNonzeroDigit = true
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>				}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>				<span class="comment">// Count significant digits after the first non-zero digit.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>				if sawNonzeroDigit {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>					digits--
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>				}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		if !hasDecimalPoint {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			<span class="comment">// Leading digit 0 should contribute once to digits.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			if len(num) == 2 &amp;&amp; num[1] == &#39;0&#39; {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				digits--
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>			}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			num = append(num, &#39;.&#39;)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		for digits &gt; 0 {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			num = append(num, &#39;0&#39;)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			digits--
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		num = append(num, tail...)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	<span class="comment">// We want a sign if asked for and if the sign is not positive.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	if f.plus || num[0] != &#39;+&#39; {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		<span class="comment">// If we&#39;re zero padding to the left we want the sign before the leading zeros.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		<span class="comment">// Achieve this by writing the sign out and then padding the unsigned number.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		if f.zero &amp;&amp; f.widPresent &amp;&amp; f.wid &gt; len(num) {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			f.buf.writeByte(num[0])
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>			f.writePadding(f.wid - len(num))
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>			f.buf.write(num[1:])
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			return
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		f.pad(num)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		return
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	<span class="comment">// No sign to show and the number is positive; just print the unsigned number.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	f.pad(num[1:])
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
</pre><p><a href="format.go?m=text">View as plain text</a></p>

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
