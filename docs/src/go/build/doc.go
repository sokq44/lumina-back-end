<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/build/doc.go - Go Documentation Server</title>

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
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/build">build</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/build">go/build</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package build gathers information about Go packages.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// # Go Path</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// The Go path is a list of directory trees containing Go source code.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// It is consulted to resolve imports that cannot be found in the standard</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Go tree. The default path is the value of the GOPATH environment</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// variable, interpreted as a path list appropriate to the operating system</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// (on Unix, the variable is a colon-separated string;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// on Windows, a semicolon-separated string;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// on Plan 9, a list).</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Each directory listed in the Go path must have a prescribed structure:</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The src/ directory holds source code. The path below &#39;src&#39; determines</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// the import path or executable name.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// The pkg/ directory holds installed package objects.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// As in the Go tree, each target operating system and</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// architecture pair has its own subdirectory of pkg</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// (pkg/GOOS_GOARCH).</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// If DIR is a directory listed in the Go path, a package with</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// source in DIR/src/foo/bar can be imported as &#34;foo/bar&#34; and</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// has its compiled form installed to &#34;DIR/pkg/GOOS_GOARCH/foo/bar.a&#34;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// (or, for gccgo, &#34;DIR/pkg/gccgo/foo/libbar.a&#34;).</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// The bin/ directory holds compiled commands.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// Each command is named for its source directory, but only</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// using the final element, not the entire path. That is, the</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// command with source in DIR/src/foo/quux is installed into</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// DIR/bin/quux, not DIR/bin/foo/quux. The foo/ is stripped</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// so that you can add DIR/bin to your PATH to get at the</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// installed commands.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// Here&#39;s an example directory layout:</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//	GOPATH=/home/user/gocode</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//	/home/user/gocode/</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//	    src/</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//	        foo/</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//	            bar/               (go code in package bar)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//	                x.go</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//	            quux/              (go code in package main)</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//	                y.go</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//	    bin/</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//	        quux                   (installed command)</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//	    pkg/</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//	        linux_amd64/</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//	            foo/</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	                bar.a          (installed package object)</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// # Build Constraints</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// A build constraint, also known as a build tag, is a condition under which a</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// file should be included in the package. Build constraints are given by a</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// line comment that begins</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//	//go:build</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// Build constraints may also be part of a file&#39;s name</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// (for example, source_windows.go will only be included if the target</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// operating system is windows).</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// See &#39;go help buildconstraint&#39;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// (https://golang.org/cmd/go/#hdr-Build_constraints) for details.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// # Binary-Only Packages</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// In Go 1.12 and earlier, it was possible to distribute packages in binary</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// form without including the source code used for compiling the package.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// The package was distributed with a source file not excluded by build</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// constraints and containing a &#34;//go:binary-only-package&#34; comment. Like a</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// build constraint, this comment appeared at the top of a file, preceded</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// only by blank lines and other line comments and with a blank line</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// following the comment, to separate it from the package documentation.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// Unlike build constraints, this comment is only recognized in non-test</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// Go source files.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// The minimal source code for a binary-only package was therefore:</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//	//go:binary-only-package</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//	package mypkg</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// The source code could include additional Go code. That code was never</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// compiled but would be processed by tools like godoc and might be useful</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// as end-user documentation.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// &#34;go build&#34; and other commands no longer support binary-only-packages.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// [Import] and [ImportDir] will still set the BinaryOnly flag in packages</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// containing these comments for use in tools and error messages.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>package build
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
