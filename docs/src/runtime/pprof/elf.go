<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/pprof/elf.go - Go Documentation Server</title>

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
<a href="elf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/pprof">pprof</a>/<span class="text-muted">elf.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/pprof">runtime/pprof</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package pprof
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>var (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	errBadELF    = errors.New(&#34;malformed ELF binary&#34;)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	errNoBuildID = errors.New(&#34;no NT_GNU_BUILD_ID found in ELF binary&#34;)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// elfBuildID returns the GNU build ID of the named ELF binary,</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// without introducing a dependency on debug/elf and its dependencies.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func elfBuildID(file string) (string, error) {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	buf := make([]byte, 256)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	f, err := os.Open(file)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	if err != nil {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		return &#34;&#34;, err
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	defer f.Close()
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	if _, err := f.ReadAt(buf[:64], 0); err != nil {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		return &#34;&#34;, err
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// ELF file begins with \x7F E L F.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if buf[0] != 0x7F || buf[1] != &#39;E&#39; || buf[2] != &#39;L&#39; || buf[3] != &#39;F&#39; {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return &#34;&#34;, errBadELF
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	var byteOrder binary.ByteOrder
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	switch buf[5] {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	default:
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		return &#34;&#34;, errBadELF
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	case 1: <span class="comment">// little-endian</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		byteOrder = binary.LittleEndian
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	case 2: <span class="comment">// big-endian</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		byteOrder = binary.BigEndian
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	var shnum int
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	var shoff, shentsize int64
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	switch buf[4] {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	default:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return &#34;&#34;, errBadELF
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	case 1: <span class="comment">// 32-bit file header</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		shoff = int64(byteOrder.Uint32(buf[32:]))
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		shentsize = int64(byteOrder.Uint16(buf[46:]))
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if shentsize != 40 {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			return &#34;&#34;, errBadELF
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		shnum = int(byteOrder.Uint16(buf[48:]))
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	case 2: <span class="comment">// 64-bit file header</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		shoff = int64(byteOrder.Uint64(buf[40:]))
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		shentsize = int64(byteOrder.Uint16(buf[58:]))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if shentsize != 64 {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			return &#34;&#34;, errBadELF
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		shnum = int(byteOrder.Uint16(buf[60:]))
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	for i := 0; i &lt; shnum; i++ {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		if _, err := f.ReadAt(buf[:shentsize], shoff+int64(i)*shentsize); err != nil {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			return &#34;&#34;, err
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if typ := byteOrder.Uint32(buf[4:]); typ != 7 { <span class="comment">// SHT_NOTE</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			continue
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		var off, size int64
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		if shentsize == 40 {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			<span class="comment">// 32-bit section header</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			off = int64(byteOrder.Uint32(buf[16:]))
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			size = int64(byteOrder.Uint32(buf[20:]))
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		} else {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			<span class="comment">// 64-bit section header</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			off = int64(byteOrder.Uint64(buf[24:]))
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			size = int64(byteOrder.Uint64(buf[32:]))
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		size += off
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		for off &lt; size {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			if _, err := f.ReadAt(buf[:16], off); err != nil { <span class="comment">// room for header + name GNU\x00</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>				return &#34;&#34;, err
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			nameSize := int(byteOrder.Uint32(buf[0:]))
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			descSize := int(byteOrder.Uint32(buf[4:]))
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			noteType := int(byteOrder.Uint32(buf[8:]))
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			descOff := off + int64(12+(nameSize+3)&amp;^3)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			off = descOff + int64((descSize+3)&amp;^3)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			if nameSize != 4 || noteType != 3 || buf[12] != &#39;G&#39; || buf[13] != &#39;N&#39; || buf[14] != &#39;U&#39; || buf[15] != &#39;\x00&#39; { <span class="comment">// want name GNU\x00 type 3 (NT_GNU_BUILD_ID)</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>				continue
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			if descSize &gt; len(buf) {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>				return &#34;&#34;, errBadELF
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			if _, err := f.ReadAt(buf[:descSize], descOff); err != nil {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>				return &#34;&#34;, err
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			return fmt.Sprintf(&#34;%x&#34;, buf[:descSize]), nil
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	return &#34;&#34;, errNoBuildID
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
</pre><p><a href="elf.go?m=text">View as plain text</a></p>

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
