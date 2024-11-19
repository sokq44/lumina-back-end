<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/tagptr_64bit.go - Go Documentation Server</title>

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
<a href="tagptr_64bit.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">tagptr_64bit.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build amd64 || arm64 || loong64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x || wasm</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/goos&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>const (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// addrBits is the number of bits needed to represent a virtual address.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// See heapAddrBits for a table of address space sizes on</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// various architectures. 48 bits is enough for all</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// architectures except s390x.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// On AMD64, virtual addresses are 48-bit (or 57-bit) numbers sign extended to 64.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// We shift the address left 16 to eliminate the sign extended part and make</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// room in the bottom for the count.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// On s390x, virtual addresses are 64-bit. There&#39;s not much we</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// can do about this, so we just hope that the kernel doesn&#39;t</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// get to really high addresses and panic if it does.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	addrBits = 48
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// In addition to the 16 bits taken from the top, we can take 3 from the</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// bottom, because node must be pointer-aligned, giving a total of 19 bits</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// of count.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	tagBits = 64 - addrBits + 3
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// On AIX, 64-bit addresses are split into 36-bit segment number and 28-bit</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// offset in segment.  Segment numbers in the range 0x0A0000000-0x0AFFFFFFF(LSA)</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// are available for mmap.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// We assume all tagged addresses are from memory allocated with mmap.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// We use one bit to distinguish between the two ranges.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	aixAddrBits = 57
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	aixTagBits  = 64 - aixAddrBits + 3
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// riscv64 SV57 mode gives 56 bits of userspace VA.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// tagged pointer code supports it,</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// but broader support for SV57 mode is incomplete,</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// and there may be other issues (see #54104).</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	riscv64AddrBits = 56
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	riscv64TagBits  = 64 - riscv64AddrBits + 3
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// The number of bits stored in the numeric tag of a taggedPointer</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>const taggedPointerBits = (goos.IsAix * aixTagBits) + (goarch.IsRiscv64 * riscv64TagBits) + ((1 - goos.IsAix) * (1 - goarch.IsRiscv64) * tagBits)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// taggedPointerPack created a taggedPointer from a pointer and a tag.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Tag bits that don&#39;t fit in the result are discarded.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>func taggedPointerPack(ptr unsafe.Pointer, tag uintptr) taggedPointer {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if GOOS == &#34;aix&#34; {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		if GOARCH != &#34;ppc64&#34; {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			throw(&#34;check this code for aix on non-ppc64&#34;)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return taggedPointer(uint64(uintptr(ptr))&lt;&lt;(64-aixAddrBits) | uint64(tag&amp;(1&lt;&lt;aixTagBits-1)))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	if GOARCH == &#34;riscv64&#34; {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return taggedPointer(uint64(uintptr(ptr))&lt;&lt;(64-riscv64AddrBits) | uint64(tag&amp;(1&lt;&lt;riscv64TagBits-1)))
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return taggedPointer(uint64(uintptr(ptr))&lt;&lt;(64-addrBits) | uint64(tag&amp;(1&lt;&lt;tagBits-1)))
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// Pointer returns the pointer from a taggedPointer.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func (tp taggedPointer) pointer() unsafe.Pointer {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if GOARCH == &#34;amd64&#34; {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		<span class="comment">// amd64 systems can place the stack above the VA hole, so we need to sign extend</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		<span class="comment">// val before unpacking.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		return unsafe.Pointer(uintptr(int64(tp) &gt;&gt; tagBits &lt;&lt; 3))
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if GOOS == &#34;aix&#34; {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		return unsafe.Pointer(uintptr((tp &gt;&gt; aixTagBits &lt;&lt; 3) | 0xa&lt;&lt;56))
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if GOARCH == &#34;riscv64&#34; {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return unsafe.Pointer(uintptr(tp &gt;&gt; riscv64TagBits &lt;&lt; 3))
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	return unsafe.Pointer(uintptr(tp &gt;&gt; tagBits &lt;&lt; 3))
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// Tag returns the tag from a taggedPointer.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func (tp taggedPointer) tag() uintptr {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return uintptr(tp &amp; (1&lt;&lt;taggedPointerBits - 1))
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
</pre><p><a href="tagptr_64bit.go?m=text">View as plain text</a></p>

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
