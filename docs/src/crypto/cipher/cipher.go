<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/cipher/cipher.go - Go Documentation Server</title>

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
<a href="cipher.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/cipher">cipher</a>/<span class="text-muted">cipher.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/cipher">crypto/cipher</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package cipher implements standard block cipher modes that can be wrapped</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// around low-level block cipher implementations.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// See https://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// and NIST Special Publication 800-38A.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package cipher
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// A Block represents an implementation of block cipher</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// using a given key. It provides the capability to encrypt</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// or decrypt individual blocks. The mode implementations</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// extend that capability to streams of blocks.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type Block interface {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// BlockSize returns the cipher&#39;s block size.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	BlockSize() int
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">// Encrypt encrypts the first block in src into dst.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Dst and src must overlap entirely or not at all.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	Encrypt(dst, src []byte)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// Decrypt decrypts the first block in src into dst.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// Dst and src must overlap entirely or not at all.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	Decrypt(dst, src []byte)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// A Stream represents a stream cipher.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>type Stream interface {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// XORKeyStream XORs each byte in the given slice with a byte from the</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// cipher&#39;s key stream. Dst and src must overlap entirely or not at all.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// If len(dst) &lt; len(src), XORKeyStream should panic. It is acceptable</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// to pass a dst bigger than src, and in that case, XORKeyStream will</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// only update dst[:len(src)] and will not touch the rest of dst.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// Multiple calls to XORKeyStream behave as if the concatenation of</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// the src buffers was passed in a single run. That is, Stream</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// maintains state and does not reset at each XORKeyStream call.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	XORKeyStream(dst, src []byte)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// A BlockMode represents a block cipher running in a block-based mode (CBC,</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// ECB etc).</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>type BlockMode interface {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// BlockSize returns the mode&#39;s block size.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	BlockSize() int
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// CryptBlocks encrypts or decrypts a number of blocks. The length of</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// src must be a multiple of the block size. Dst and src must overlap</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// entirely or not at all.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// If len(dst) &lt; len(src), CryptBlocks should panic. It is acceptable</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// to pass a dst bigger than src, and in that case, CryptBlocks will</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// only update dst[:len(src)] and will not touch the rest of dst.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Multiple calls to CryptBlocks behave as if the concatenation of</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// the src buffers was passed in a single run. That is, BlockMode</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// maintains state and does not reset at each CryptBlocks call.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	CryptBlocks(dst, src []byte)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
</pre><p><a href="cipher.go?m=text">View as plain text</a></p>

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
