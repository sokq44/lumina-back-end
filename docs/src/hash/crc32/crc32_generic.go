<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/hash/crc32/crc32_generic.go - Go Documentation Server</title>

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
<a href="crc32_generic.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/hash">hash</a>/<a href="http://localhost:8080/src/hash/crc32">crc32</a>/<span class="text-muted">crc32_generic.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/hash/crc32">hash/crc32</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file contains CRC32 algorithms that are not specific to any architecture</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// and don&#39;t use hardware acceleration.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// The simple (and slow) CRC32 implementation only uses a 256*4 bytes table.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// The slicing-by-8 algorithm is a faster implementation that uses a bigger</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// table (8*256*4 bytes).</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>package crc32
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// simpleMakeTable allocates and constructs a Table for the specified</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// polynomial. The table is suitable for use with the simple algorithm</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// (simpleUpdate).</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>func simpleMakeTable(poly uint32) *Table {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	t := new(Table)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	simplePopulateTable(poly, t)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	return t
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// simplePopulateTable constructs a Table for the specified polynomial, suitable</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// for use with simpleUpdate.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func simplePopulateTable(poly uint32, t *Table) {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	for i := 0; i &lt; 256; i++ {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		crc := uint32(i)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		for j := 0; j &lt; 8; j++ {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			if crc&amp;1 == 1 {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>				crc = (crc &gt;&gt; 1) ^ poly
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			} else {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>				crc &gt;&gt;= 1
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		t[i] = crc
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// simpleUpdate uses the simple algorithm to update the CRC, given a table that</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// was previously computed using simpleMakeTable.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func simpleUpdate(crc uint32, tab *Table, p []byte) uint32 {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	crc = ^crc
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	for _, v := range p {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		crc = tab[byte(crc)^v] ^ (crc &gt;&gt; 8)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	return ^crc
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// Use slicing-by-8 when payload &gt;= this value.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>const slicing8Cutoff = 16
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// slicing8Table is array of 8 Tables, used by the slicing-by-8 algorithm.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>type slicing8Table [8]Table
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// slicingMakeTable constructs a slicing8Table for the specified polynomial. The</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// table is suitable for use with the slicing-by-8 algorithm (slicingUpdate).</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func slicingMakeTable(poly uint32) *slicing8Table {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	t := new(slicing8Table)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	simplePopulateTable(poly, &amp;t[0])
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	for i := 0; i &lt; 256; i++ {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		crc := t[0][i]
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		for j := 1; j &lt; 8; j++ {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			crc = t[0][crc&amp;0xFF] ^ (crc &gt;&gt; 8)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			t[j][i] = crc
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return t
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// slicingUpdate uses the slicing-by-8 algorithm to update the CRC, given a</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// table that was previously computed using slicingMakeTable.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func slicingUpdate(crc uint32, tab *slicing8Table, p []byte) uint32 {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if len(p) &gt;= slicing8Cutoff {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		crc = ^crc
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		for len(p) &gt; 8 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			crc ^= uint32(p[0]) | uint32(p[1])&lt;&lt;8 | uint32(p[2])&lt;&lt;16 | uint32(p[3])&lt;&lt;24
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			crc = tab[0][p[7]] ^ tab[1][p[6]] ^ tab[2][p[5]] ^ tab[3][p[4]] ^
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>				tab[4][crc&gt;&gt;24] ^ tab[5][(crc&gt;&gt;16)&amp;0xFF] ^
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>				tab[6][(crc&gt;&gt;8)&amp;0xFF] ^ tab[7][crc&amp;0xFF]
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			p = p[8:]
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		crc = ^crc
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	if len(p) == 0 {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		return crc
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return simpleUpdate(crc, &amp;tab[0], p)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
</pre><p><a href="crc32_generic.go?m=text">View as plain text</a></p>

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
