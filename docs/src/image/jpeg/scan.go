<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/image/jpeg/scan.go - Go Documentation Server</title>

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
<a href="scan.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/image">image</a>/<a href="http://localhost:8080/src/image/jpeg">jpeg</a>/<span class="text-muted">scan.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/image/jpeg">image/jpeg</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package jpeg
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;image&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// makeImg allocates and initializes the destination image.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>func (d *decoder) makeImg(mxx, myy int) {
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	if d.nComp == 1 {
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>		m := image.NewGray(image.Rect(0, 0, 8*mxx, 8*myy))
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>		d.img1 = m.SubImage(image.Rect(0, 0, d.width, d.height)).(*image.Gray)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		return
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	h0 := d.comp[0].h
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	v0 := d.comp[0].v
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	hRatio := h0 / d.comp[1].h
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	vRatio := v0 / d.comp[1].v
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	var subsampleRatio image.YCbCrSubsampleRatio
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	switch hRatio&lt;&lt;4 | vRatio {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	case 0x11:
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio444
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	case 0x12:
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio440
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	case 0x21:
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio422
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	case 0x22:
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio420
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	case 0x41:
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio411
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	case 0x42:
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		subsampleRatio = image.YCbCrSubsampleRatio410
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	default:
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		panic(&#34;unreachable&#34;)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	m := image.NewYCbCr(image.Rect(0, 0, 8*h0*mxx, 8*v0*myy), subsampleRatio)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	d.img3 = m.SubImage(image.Rect(0, 0, d.width, d.height)).(*image.YCbCr)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if d.nComp == 4 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		h3, v3 := d.comp[3].h, d.comp[3].v
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		d.blackPix = make([]byte, 8*h3*mxx*8*v3*myy)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		d.blackStride = 8 * h3 * mxx
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// Specified in section B.2.3.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>func (d *decoder) processSOS(n int) error {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	if d.nComp == 0 {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		return FormatError(&#34;missing SOF marker&#34;)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if n &lt; 6 || 4+2*d.nComp &lt; n || n%2 != 0 {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		return FormatError(&#34;SOS has wrong length&#34;)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if err := d.readFull(d.tmp[:n]); err != nil {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		return err
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	nComp := int(d.tmp[0])
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if n != 4+2*nComp {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return FormatError(&#34;SOS length inconsistent with number of components&#34;)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	var scan [maxComponents]struct {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		compIndex uint8
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		td        uint8 <span class="comment">// DC table selector.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		ta        uint8 <span class="comment">// AC table selector.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	totalHV := 0
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	for i := 0; i &lt; nComp; i++ {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		cs := d.tmp[1+2*i] <span class="comment">// Component selector.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		compIndex := -1
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		for j, comp := range d.comp[:d.nComp] {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			if cs == comp.c {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>				compIndex = j
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		if compIndex &lt; 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			return FormatError(&#34;unknown component selector&#34;)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		scan[i].compIndex = uint8(compIndex)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">// Section B.2.3 states that &#34;the value of Cs_j shall be different from</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		<span class="comment">// the values of Cs_1 through Cs_(j-1)&#34;. Since we have previously</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		<span class="comment">// verified that a frame&#39;s component identifiers (C_i values in section</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		<span class="comment">// B.2.2) are unique, it suffices to check that the implicit indexes</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		<span class="comment">// into d.comp are unique.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		for j := 0; j &lt; i; j++ {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			if scan[i].compIndex == scan[j].compIndex {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				return FormatError(&#34;repeated component selector&#34;)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		totalHV += d.comp[compIndex].h * d.comp[compIndex].v
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// The baseline t &lt;= 1 restriction is specified in table B.3.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		scan[i].td = d.tmp[2+2*i] &gt;&gt; 4
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		if t := scan[i].td; t &gt; maxTh || (d.baseline &amp;&amp; t &gt; 1) {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			return FormatError(&#34;bad Td value&#34;)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		scan[i].ta = d.tmp[2+2*i] &amp; 0x0f
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if t := scan[i].ta; t &gt; maxTh || (d.baseline &amp;&amp; t &gt; 1) {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			return FormatError(&#34;bad Ta value&#34;)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// Section B.2.3 states that if there is more than one component then the</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	<span class="comment">// total H*V values in a scan must be &lt;= 10.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if d.nComp &gt; 1 &amp;&amp; totalHV &gt; 10 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return FormatError(&#34;total sampling factors too large&#34;)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// zigStart and zigEnd are the spectral selection bounds.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// ah and al are the successive approximation high and low values.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// The spec calls these values Ss, Se, Ah and Al.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// For progressive JPEGs, these are the two more-or-less independent</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// aspects of progression. Spectral selection progression is when not</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// all of a block&#39;s 64 DCT coefficients are transmitted in one pass.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// For example, three passes could transmit coefficient 0 (the DC</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// component), coefficients 1-5, and coefficients 6-63, in zig-zag</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// order. Successive approximation is when not all of the bits of a</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// band of coefficients are transmitted in one pass. For example,</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// three passes could transmit the 6 most significant bits, followed</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// by the second-least significant bit, followed by the least</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// significant bit.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// For sequential JPEGs, these parameters are hard-coded to 0/63/0/0, as</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// per table B.3.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	zigStart, zigEnd, ah, al := int32(0), int32(blockSize-1), uint32(0), uint32(0)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if d.progressive {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		zigStart = int32(d.tmp[1+2*nComp])
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		zigEnd = int32(d.tmp[2+2*nComp])
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		ah = uint32(d.tmp[3+2*nComp] &gt;&gt; 4)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		al = uint32(d.tmp[3+2*nComp] &amp; 0x0f)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		if (zigStart == 0 &amp;&amp; zigEnd != 0) || zigStart &gt; zigEnd || blockSize &lt;= zigEnd {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			return FormatError(&#34;bad spectral selection bounds&#34;)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		if zigStart != 0 &amp;&amp; nComp != 1 {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			return FormatError(&#34;progressive AC coefficients for more than one component&#34;)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if ah != 0 &amp;&amp; ah != al+1 {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			return FormatError(&#34;bad successive approximation values&#34;)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// mxx and myy are the number of MCUs (Minimum Coded Units) in the image.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	h0, v0 := d.comp[0].h, d.comp[0].v <span class="comment">// The h and v values from the Y components.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	mxx := (d.width + 8*h0 - 1) / (8 * h0)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	myy := (d.height + 8*v0 - 1) / (8 * v0)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if d.img1 == nil &amp;&amp; d.img3 == nil {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		d.makeImg(mxx, myy)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if d.progressive {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		for i := 0; i &lt; nComp; i++ {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			compIndex := scan[i].compIndex
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			if d.progCoeffs[compIndex] == nil {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				d.progCoeffs[compIndex] = make([]block, mxx*myy*d.comp[compIndex].h*d.comp[compIndex].v)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	d.bits = bits{}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	mcu, expectedRST := 0, uint8(rst0Marker)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	var (
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// b is the decoded coefficients, in natural (not zig-zag) order.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		b  block
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		dc [maxComponents]int32
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// bx and by are the location of the current block, in units of 8x8</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// blocks: the third block in the first row has (bx, by) = (2, 0).</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		bx, by     int
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		blockCount int
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	for my := 0; my &lt; myy; my++ {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		for mx := 0; mx &lt; mxx; mx++ {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			for i := 0; i &lt; nComp; i++ {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				compIndex := scan[i].compIndex
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>				hi := d.comp[compIndex].h
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>				vi := d.comp[compIndex].v
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				for j := 0; j &lt; hi*vi; j++ {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>					<span class="comment">// The blocks are traversed one MCU at a time. For 4:2:0 chroma</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>					<span class="comment">// subsampling, there are four Y 8x8 blocks in every 16x16 MCU.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>					<span class="comment">//</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>					<span class="comment">// For a sequential 32x16 pixel image, the Y blocks visiting order is:</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>					<span class="comment">//	0 1 4 5</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>					<span class="comment">//	2 3 6 7</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>					<span class="comment">//</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>					<span class="comment">// For progressive images, the interleaved scans (those with nComp &gt; 1)</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>					<span class="comment">// are traversed as above, but non-interleaved scans are traversed left</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>					<span class="comment">// to right, top to bottom:</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>					<span class="comment">//	0 1 2 3</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>					<span class="comment">//	4 5 6 7</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>					<span class="comment">// Only DC scans (zigStart == 0) can be interleaved. AC scans must have</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>					<span class="comment">// only one component.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>					<span class="comment">//</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>					<span class="comment">// To further complicate matters, for non-interleaved scans, there is no</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>					<span class="comment">// data for any blocks that are inside the image at the MCU level but</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>					<span class="comment">// outside the image at the pixel level. For example, a 24x16 pixel 4:2:0</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>					<span class="comment">// progressive image consists of two 16x16 MCUs. The interleaved scans</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					<span class="comment">// will process 8 Y blocks:</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>					<span class="comment">//	0 1 4 5</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>					<span class="comment">//	2 3 6 7</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>					<span class="comment">// The non-interleaved scans will process only 6 Y blocks:</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>					<span class="comment">//	0 1 2</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>					<span class="comment">//	3 4 5</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>					if nComp != 1 {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>						bx = hi*mx + j%hi
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>						by = vi*my + j/hi
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>					} else {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>						q := mxx * hi
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>						bx = blockCount % q
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>						by = blockCount / q
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>						blockCount++
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>						if bx*8 &gt;= d.width || by*8 &gt;= d.height {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>							continue
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>						}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>					}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					<span class="comment">// Load the previous partially decoded coefficients, if applicable.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>					if d.progressive {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>						b = d.progCoeffs[compIndex][by*mxx*hi+bx]
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>					} else {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>						b = block{}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>					}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>					if ah != 0 {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>						if err := d.refine(&amp;b, &amp;d.huff[acTable][scan[i].ta], zigStart, zigEnd, 1&lt;&lt;al); err != nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>							return err
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>						}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>					} else {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>						zig := zigStart
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>						if zig == 0 {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>							zig++
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>							<span class="comment">// Decode the DC coefficient, as specified in section F.2.2.1.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>							value, err := d.decodeHuffman(&amp;d.huff[dcTable][scan[i].td])
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>							if err != nil {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>								return err
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>							}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>							if value &gt; 16 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>								return UnsupportedError(&#34;excessive DC component&#34;)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>							}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>							dcDelta, err := d.receiveExtend(value)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>							if err != nil {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>								return err
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>							}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>							dc[compIndex] += dcDelta
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>							b[0] = dc[compIndex] &lt;&lt; al
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>						}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>						if zig &lt;= zigEnd &amp;&amp; d.eobRun &gt; 0 {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>							d.eobRun--
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>						} else {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>							<span class="comment">// Decode the AC coefficients, as specified in section F.2.2.2.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>							huff := &amp;d.huff[acTable][scan[i].ta]
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>							for ; zig &lt;= zigEnd; zig++ {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>								value, err := d.decodeHuffman(huff)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>								if err != nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>									return err
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>								}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>								val0 := value &gt;&gt; 4
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>								val1 := value &amp; 0x0f
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>								if val1 != 0 {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>									zig += int32(val0)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>									if zig &gt; zigEnd {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>										break
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>									}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>									ac, err := d.receiveExtend(val1)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>									if err != nil {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>										return err
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>									}
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>									b[unzig[zig]] = ac &lt;&lt; al
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>								} else {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>									if val0 != 0x0f {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>										d.eobRun = uint16(1 &lt;&lt; val0)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>										if val0 != 0 {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>											bits, err := d.decodeBits(int32(val0))
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>											if err != nil {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>												return err
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>											}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>											d.eobRun |= uint16(bits)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>										}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>										d.eobRun--
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>										break
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>									}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>									zig += 0x0f
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>								}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>							}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>						}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>					}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>					if d.progressive {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>						<span class="comment">// Save the coefficients.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>						d.progCoeffs[compIndex][by*mxx*hi+bx] = b
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>						<span class="comment">// At this point, we could call reconstructBlock to dequantize and perform the</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>						<span class="comment">// inverse DCT, to save early stages of a progressive image to the *image.YCbCr</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>						<span class="comment">// buffers (the whole point of progressive encoding), but in Go, the jpeg.Decode</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>						<span class="comment">// function does not return until the entire image is decoded, so we &#34;continue&#34;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>						<span class="comment">// here to avoid wasted computation. Instead, reconstructBlock is called on each</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>						<span class="comment">// accumulated block by the reconstructProgressiveImage method after all of the</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>						<span class="comment">// SOS markers are processed.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>						continue
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>					}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>					if err := d.reconstructBlock(&amp;b, bx, by, int(compIndex)); err != nil {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>						return err
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>					}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				} <span class="comment">// for j</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			} <span class="comment">// for i</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			mcu++
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			if d.ri &gt; 0 &amp;&amp; mcu%d.ri == 0 &amp;&amp; mcu &lt; mxx*myy {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				<span class="comment">// A more sophisticated decoder could use RST[0-7] markers to resynchronize from corrupt input,</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				<span class="comment">// but this one assumes well-formed input, and hence the restart marker follows immediately.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>				if err := d.readFull(d.tmp[:2]); err != nil {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>					return err
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>				<span class="comment">// Section F.1.2.3 says that &#34;Byte alignment of markers is</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>				<span class="comment">// achieved by padding incomplete bytes with 1-bits. If padding</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>				<span class="comment">// with 1-bits creates a X’FF’ value, a zero byte is stuffed</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>				<span class="comment">// before adding the marker.&#34;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>				<span class="comment">// Seeing &#34;\xff\x00&#34; here is not spec compliant, as we are not</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>				<span class="comment">// expecting an *incomplete* byte (that needed padding). Still,</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>				<span class="comment">// some real world encoders (see golang.org/issue/28717) insert</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				<span class="comment">// it, so we accept it and re-try the 2 byte read.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>				<span class="comment">//</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>				<span class="comment">// libjpeg issues a warning (but not an error) for this:</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				<span class="comment">// https://github.com/LuaDist/libjpeg/blob/6c0fcb8ddee365e7abc4d332662b06900612e923/jdmarker.c#L1041-L1046</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>				if d.tmp[0] == 0xff &amp;&amp; d.tmp[1] == 0x00 {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>					if err := d.readFull(d.tmp[:2]); err != nil {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>						return err
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>					}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				if d.tmp[0] != 0xff || d.tmp[1] != expectedRST {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>					return FormatError(&#34;bad RST marker&#34;)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				expectedRST++
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>				if expectedRST == rst7Marker+1 {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>					expectedRST = rst0Marker
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				<span class="comment">// Reset the Huffman decoder.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>				d.bits = bits{}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>				<span class="comment">// Reset the DC components, as per section F.2.1.3.1.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>				dc = [maxComponents]int32{}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				<span class="comment">// Reset the progressive decoder state, as per section G.1.2.2.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>				d.eobRun = 0
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		} <span class="comment">// for mx</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	} <span class="comment">// for my</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	return nil
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// refine decodes a successive approximation refinement block, as specified in</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// section G.1.2.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>func (d *decoder) refine(b *block, h *huffman, zigStart, zigEnd, delta int32) error {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// Refining a DC component is trivial.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	if zigStart == 0 {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		if zigEnd != 0 {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			panic(&#34;unreachable&#34;)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		bit, err := d.decodeBit()
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		if err != nil {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			return err
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		if bit {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			b[0] |= delta
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		return nil
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// Refining AC components is more complicated; see sections G.1.2.2 and G.1.2.3.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	zig := zigStart
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	if d.eobRun == 0 {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	loop:
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		for ; zig &lt;= zigEnd; zig++ {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			z := int32(0)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			value, err := d.decodeHuffman(h)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			if err != nil {
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>				return err
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			val0 := value &gt;&gt; 4
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			val1 := value &amp; 0x0f
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			switch val1 {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			case 0:
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				if val0 != 0x0f {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>					d.eobRun = uint16(1 &lt;&lt; val0)
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>					if val0 != 0 {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>						bits, err := d.decodeBits(int32(val0))
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>						if err != nil {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>							return err
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>						}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>						d.eobRun |= uint16(bits)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>					}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>					break loop
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>				}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			case 1:
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>				z = delta
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				bit, err := d.decodeBit()
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				if err != nil {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>					return err
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>				}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>				if !bit {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>					z = -z
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			default:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>				return FormatError(&#34;unexpected Huffman code&#34;)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			zig, err = d.refineNonZeroes(b, zig, zigEnd, int32(val0), delta)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			if err != nil {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				return err
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			if zig &gt; zigEnd {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				return FormatError(&#34;too many coefficients&#34;)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			if z != 0 {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>				b[unzig[zig]] = z
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if d.eobRun &gt; 0 {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		d.eobRun--
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		if _, err := d.refineNonZeroes(b, zig, zigEnd, -1, delta); err != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			return err
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	return nil
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span><span class="comment">// refineNonZeroes refines non-zero entries of b in zig-zag order. If nz &gt;= 0,</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span><span class="comment">// the first nz zero entries are skipped over.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>func (d *decoder) refineNonZeroes(b *block, zig, zigEnd, nz, delta int32) (int32, error) {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	for ; zig &lt;= zigEnd; zig++ {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		u := unzig[zig]
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		if b[u] == 0 {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			if nz == 0 {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>				break
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			nz--
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			continue
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		bit, err := d.decodeBit()
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		if err != nil {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			return 0, err
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		if !bit {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			continue
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		if b[u] &gt;= 0 {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			b[u] += delta
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		} else {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			b[u] -= delta
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	return zig, nil
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>func (d *decoder) reconstructProgressiveImage() error {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// The h0, mxx, by and bx variables have the same meaning as in the</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// processSOS method.</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	h0 := d.comp[0].h
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	mxx := (d.width + 8*h0 - 1) / (8 * h0)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	for i := 0; i &lt; d.nComp; i++ {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		if d.progCoeffs[i] == nil {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>			continue
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		v := 8 * d.comp[0].v / d.comp[i].v
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		h := 8 * d.comp[0].h / d.comp[i].h
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		stride := mxx * d.comp[i].h
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		for by := 0; by*v &lt; d.height; by++ {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			for bx := 0; bx*h &lt; d.width; bx++ {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>				if err := d.reconstructBlock(&amp;d.progCoeffs[i][by*stride+bx], bx, by, i); err != nil {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>					return err
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>				}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	return nil
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">// reconstructBlock dequantizes, performs the inverse DCT and stores the block</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// to the image.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>func (d *decoder) reconstructBlock(b *block, bx, by, compIndex int) error {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	qt := &amp;d.quant[d.comp[compIndex].tq]
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	for zig := 0; zig &lt; blockSize; zig++ {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		b[unzig[zig]] *= qt[zig]
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	idct(b)
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	dst, stride := []byte(nil), 0
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if d.nComp == 1 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		dst, stride = d.img1.Pix[8*(by*d.img1.Stride+bx):], d.img1.Stride
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	} else {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		switch compIndex {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		case 0:
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			dst, stride = d.img3.Y[8*(by*d.img3.YStride+bx):], d.img3.YStride
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		case 1:
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			dst, stride = d.img3.Cb[8*(by*d.img3.CStride+bx):], d.img3.CStride
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		case 2:
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			dst, stride = d.img3.Cr[8*(by*d.img3.CStride+bx):], d.img3.CStride
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		case 3:
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			dst, stride = d.blackPix[8*(by*d.blackStride+bx):], d.blackStride
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		default:
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			return UnsupportedError(&#34;too many components&#34;)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	<span class="comment">// Level shift by +128, clip to [0, 255], and write to dst.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	for y := 0; y &lt; 8; y++ {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		y8 := y * 8
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		yStride := y * stride
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		for x := 0; x &lt; 8; x++ {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			c := b[y8+x]
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			if c &lt; -128 {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>				c = 0
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			} else if c &gt; 127 {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>				c = 255
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			} else {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>				c += 128
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			dst[yStride+x] = uint8(c)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	return nil
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>
</pre><p><a href="scan.go?m=text">View as plain text</a></p>

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
