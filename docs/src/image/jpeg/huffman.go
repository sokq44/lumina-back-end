<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/image/jpeg/huffman.go - Go Documentation Server</title>

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
<a href="huffman.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/image">image</a>/<a href="http://localhost:8080/src/image/jpeg">jpeg</a>/<span class="text-muted">huffman.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/image/jpeg">image/jpeg</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package jpeg
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// maxCodeLength is the maximum (inclusive) number of bits in a Huffman code.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>const maxCodeLength = 16
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// maxNCodes is the maximum (inclusive) number of codes in a Huffman tree.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>const maxNCodes = 256
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// lutSize is the log-2 size of the Huffman decoder&#39;s look-up table.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const lutSize = 8
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// huffman is a Huffman decoder, specified in section C.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type huffman struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// length is the number of codes in the tree.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	nCodes int32
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// lut is the look-up table for the next lutSize bits in the bit-stream.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// The high 8 bits of the uint16 are the encoded value. The low 8 bits</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// are 1 plus the code length, or 0 if the value is too large to fit in</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// lutSize bits.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	lut [1 &lt;&lt; lutSize]uint16
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// vals are the decoded values, sorted by their encoding.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	vals [maxNCodes]uint8
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// minCodes[i] is the minimum code of length i, or -1 if there are no</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// codes of that length.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	minCodes [maxCodeLength]int32
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// maxCodes[i] is the maximum code of length i, or -1 if there are no</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// codes of that length.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	maxCodes [maxCodeLength]int32
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// valsIndices[i] is the index into vals of minCodes[i].</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	valsIndices [maxCodeLength]int32
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// errShortHuffmanData means that an unexpected EOF occurred while decoding</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// Huffman data.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>var errShortHuffmanData = FormatError(&#34;short Huffman data&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// ensureNBits reads bytes from the byte buffer to ensure that d.bits.n is at</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// least n. For best performance (avoiding function calls inside hot loops),</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// the caller is the one responsible for first checking that d.bits.n &lt; n.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func (d *decoder) ensureNBits(n int32) error {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	for {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		c, err := d.readByteStuffedByte()
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		if err != nil {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			if err == io.ErrUnexpectedEOF {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>				return errShortHuffmanData
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			return err
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		d.bits.a = d.bits.a&lt;&lt;8 | uint32(c)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		d.bits.n += 8
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		if d.bits.m == 0 {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			d.bits.m = 1 &lt;&lt; 7
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		} else {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			d.bits.m &lt;&lt;= 8
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if d.bits.n &gt;= n {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			break
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return nil
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// receiveExtend is the composition of RECEIVE and EXTEND, specified in section</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// F.2.2.1.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func (d *decoder) receiveExtend(t uint8) (int32, error) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if d.bits.n &lt; int32(t) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		if err := d.ensureNBits(int32(t)); err != nil {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			return 0, err
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	d.bits.n -= int32(t)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	d.bits.m &gt;&gt;= t
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	s := int32(1) &lt;&lt; t
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	x := int32(d.bits.a&gt;&gt;uint8(d.bits.n)) &amp; (s - 1)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if x &lt; s&gt;&gt;1 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		x += ((-1) &lt;&lt; t) + 1
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	return x, nil
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// processDHT processes a Define Huffman Table marker, and initializes a huffman</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// struct from its contents. Specified in section B.2.4.2.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (d *decoder) processDHT(n int) error {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	for n &gt; 0 {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if n &lt; 17 {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			return FormatError(&#34;DHT has wrong length&#34;)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if err := d.readFull(d.tmp[:17]); err != nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return err
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		tc := d.tmp[0] &gt;&gt; 4
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		if tc &gt; maxTc {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			return FormatError(&#34;bad Tc value&#34;)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		th := d.tmp[0] &amp; 0x0f
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// The baseline th &lt;= 1 restriction is specified in table B.5.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if th &gt; maxTh || (d.baseline &amp;&amp; th &gt; 1) {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			return FormatError(&#34;bad Th value&#34;)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		h := &amp;d.huff[tc][th]
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		<span class="comment">// Read nCodes and h.vals (and derive h.nCodes).</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// nCodes[i] is the number of codes with code length i.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		<span class="comment">// h.nCodes is the total number of codes.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		h.nCodes = 0
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		var nCodes [maxCodeLength]int32
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		for i := range nCodes {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			nCodes[i] = int32(d.tmp[i+1])
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			h.nCodes += nCodes[i]
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		if h.nCodes == 0 {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			return FormatError(&#34;Huffman table has zero length&#34;)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		if h.nCodes &gt; maxNCodes {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			return FormatError(&#34;Huffman table has excessive length&#34;)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		n -= int(h.nCodes) + 17
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		if n &lt; 0 {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			return FormatError(&#34;DHT has wrong length&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		if err := d.readFull(h.vals[:h.nCodes]); err != nil {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			return err
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		<span class="comment">// Derive the look-up table.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		for i := range h.lut {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			h.lut[i] = 0
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		var x, code uint32
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		for i := uint32(0); i &lt; lutSize; i++ {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			code &lt;&lt;= 1
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			for j := int32(0); j &lt; nCodes[i]; j++ {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>				<span class="comment">// The codeLength is 1+i, so shift code by 8-(1+i) to</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>				<span class="comment">// calculate the high bits for every 8-bit sequence</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				<span class="comment">// whose codeLength&#39;s high bits matches code.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>				<span class="comment">// The high 8 bits of lutValue are the encoded value.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>				<span class="comment">// The low 8 bits are 1 plus the codeLength.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				base := uint8(code &lt;&lt; (7 - i))
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				lutValue := uint16(h.vals[x])&lt;&lt;8 | uint16(2+i)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				for k := uint8(0); k &lt; 1&lt;&lt;(7-i); k++ {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>					h.lut[base|k] = lutValue
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				code++
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				x++
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		<span class="comment">// Derive minCodes, maxCodes, and valsIndices.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		var c, index int32
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		for i, n := range nCodes {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			if n == 0 {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>				h.minCodes[i] = -1
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>				h.maxCodes[i] = -1
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>				h.valsIndices[i] = -1
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			} else {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				h.minCodes[i] = c
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>				h.maxCodes[i] = c + n - 1
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>				h.valsIndices[i] = index
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				c += n
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				index += n
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			c &lt;&lt;= 1
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	return nil
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// decodeHuffman returns the next Huffman-coded value from the bit-stream,</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// decoded according to h.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>func (d *decoder) decodeHuffman(h *huffman) (uint8, error) {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if h.nCodes == 0 {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return 0, FormatError(&#34;uninitialized Huffman table&#34;)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if d.bits.n &lt; 8 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		if err := d.ensureNBits(8); err != nil {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			if err != errMissingFF00 &amp;&amp; err != errShortHuffmanData {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				return 0, err
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// There are no more bytes of data in this segment, but we may still</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">// be able to read the next symbol out of the previously read bits.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			<span class="comment">// First, undo the readByte that the ensureNBits call made.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			if d.bytes.nUnreadable != 0 {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>				d.unreadByteStuffedByte()
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			goto slowPath
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if v := h.lut[(d.bits.a&gt;&gt;uint32(d.bits.n-lutSize))&amp;0xff]; v != 0 {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		n := (v &amp; 0xff) - 1
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		d.bits.n -= int32(n)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		d.bits.m &gt;&gt;= n
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return uint8(v &gt;&gt; 8), nil
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>slowPath:
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	for i, code := 0, int32(0); i &lt; maxCodeLength; i++ {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if d.bits.n == 0 {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			if err := d.ensureNBits(1); err != nil {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				return 0, err
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		if d.bits.a&amp;d.bits.m != 0 {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			code |= 1
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		d.bits.n--
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		d.bits.m &gt;&gt;= 1
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		if code &lt;= h.maxCodes[i] {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			return h.vals[h.valsIndices[i]+code-h.minCodes[i]], nil
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		code &lt;&lt;= 1
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	return 0, FormatError(&#34;bad Huffman code&#34;)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>func (d *decoder) decodeBit() (bool, error) {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	if d.bits.n == 0 {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if err := d.ensureNBits(1); err != nil {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			return false, err
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	ret := d.bits.a&amp;d.bits.m != 0
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	d.bits.n--
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	d.bits.m &gt;&gt;= 1
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	return ret, nil
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>func (d *decoder) decodeBits(n int32) (uint32, error) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if d.bits.n &lt; n {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if err := d.ensureNBits(n); err != nil {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			return 0, err
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	ret := d.bits.a &gt;&gt; uint32(d.bits.n-n)
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	ret &amp;= (1 &lt;&lt; uint32(n)) - 1
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	d.bits.n -= n
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	d.bits.m &gt;&gt;= uint32(n)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	return ret, nil
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
</pre><p><a href="huffman.go?m=text">View as plain text</a></p>

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
