<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/version.go - Go Documentation Server</title>

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
<a href="version.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">version.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2021 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package types
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;go/version&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/goversion&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// A goVersion is a Go language version string of the form &#34;go1.%d&#34;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// where d is the minor version number. goVersion strings don&#39;t</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// contain release numbers (&#34;go1.20.1&#34; is not a valid goVersion).</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type goVersion string
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// asGoVersion returns v as a goVersion (e.g., &#34;go1.20.1&#34; becomes &#34;go1.20&#34;).</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// If v is not a valid Go version, the result is the empty string.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>func asGoVersion(v string) goVersion {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	return goVersion(version.Lang(v))
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// isValid reports whether v is a valid Go version.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>func (v goVersion) isValid() bool {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	return v != &#34;&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// cmp returns -1, 0, or +1 depending on whether x &lt; y, x == y, or x &gt; y,</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// interpreted as Go versions.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func (x goVersion) cmp(y goVersion) int {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	return version.Compare(string(x), string(y))
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>var (
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// Go versions that introduced language changes</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	go1_9  = asGoVersion(&#34;go1.9&#34;)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	go1_13 = asGoVersion(&#34;go1.13&#34;)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	go1_14 = asGoVersion(&#34;go1.14&#34;)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	go1_17 = asGoVersion(&#34;go1.17&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	go1_18 = asGoVersion(&#34;go1.18&#34;)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	go1_20 = asGoVersion(&#34;go1.20&#34;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	go1_21 = asGoVersion(&#34;go1.21&#34;)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	go1_22 = asGoVersion(&#34;go1.22&#34;)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// current (deployed) Go version</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	go_current = asGoVersion(fmt.Sprintf(&#34;go1.%d&#34;, goversion.Version))
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// langCompat reports an error if the representation of a numeric</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// literal is not compatible with the current language version.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func (check *Checker) langCompat(lit *ast.BasicLit) {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	s := lit.Value
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if len(s) &lt;= 2 || check.allowVersion(check.pkg, lit, go1_13) {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		return
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// len(s) &gt; 2</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if strings.Contains(s, &#34;_&#34;) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		check.versionErrorf(lit, go1_13, &#34;underscores in numeric literals&#34;)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if s[0] != &#39;0&#39; {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	radix := s[1]
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if radix == &#39;b&#39; || radix == &#39;B&#39; {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		check.versionErrorf(lit, go1_13, &#34;binary literals&#34;)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		return
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if radix == &#39;o&#39; || radix == &#39;O&#39; {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		check.versionErrorf(lit, go1_13, &#34;0o/0O-style octal literals&#34;)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		return
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	if lit.Kind != token.INT &amp;&amp; (radix == &#39;x&#39; || radix == &#39;X&#39;) {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		check.versionErrorf(lit, go1_13, &#34;hexadecimal floating-point literals&#34;)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// allowVersion reports whether the given package is allowed to use version v.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (check *Checker) allowVersion(pkg *Package, at positioner, v goVersion) bool {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// We assume that imported packages have all been checked,</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// so we only have to check for the local package.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if pkg != check.pkg {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return true
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// If no explicit file version is specified,</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// fileVersion corresponds to the module version.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	var fileVersion goVersion
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if pos := at.Pos(); pos.IsValid() {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		<span class="comment">// We need version.Lang below because file versions</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// can be (unaltered) Config.GoVersion strings that</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		<span class="comment">// may contain dot-release information.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		fileVersion = asGoVersion(check.versions[check.fileFor(pos)])
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return !fileVersion.isValid() || fileVersion.cmp(v) &gt;= 0
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// verifyVersionf is like allowVersion but also accepts a format string and arguments</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// which are used to report a version error if allowVersion returns false. It uses the</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// current package.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (check *Checker) verifyVersionf(at positioner, v goVersion, format string, args ...interface{}) bool {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if !check.allowVersion(check.pkg, at, v) {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		check.versionErrorf(at, v, format, args...)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return false
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	return true
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// TODO(gri) Consider a more direct (position-independent) mechanism</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">//           to identify which file we&#39;re in so that version checks</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//           work correctly in the absence of correct position info.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">// fileFor returns the *ast.File which contains the position pos.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// If there are no files, the result is nil.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// The position must be valid.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>func (check *Checker) fileFor(pos token.Pos) *ast.File {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	assert(pos.IsValid())
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// Eval and CheckExpr tests may not have any source files.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if len(check.files) == 0 {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return nil
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	for _, file := range check.files {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		if file.FileStart &lt;= pos &amp;&amp; pos &lt; file.FileEnd {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			return file
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	panic(check.sprintf(&#34;file not found for pos = %d (%s)&#34;, int(pos), check.fset.Position(pos)))
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
</pre><p><a href="version.go?m=text">View as plain text</a></p>

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
