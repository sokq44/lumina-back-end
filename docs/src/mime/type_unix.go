<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/mime/type_unix.go - Go Documentation Server</title>

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
<a href="type_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/mime">mime</a>/<span class="text-muted">type_unix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/mime">mime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || (js &amp;&amp; wasm) || wasip1</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package mime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>func init() {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	osInitMime = initMimeUnix
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// See https://specifications.freedesktop.org/shared-mime-info-spec/shared-mime-info-spec-0.21.html</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// for the FreeDesktop Shared MIME-info Database specification.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>var mimeGlobs = []string{
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;/usr/local/share/mime/globs2&#34;,
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;/usr/share/mime/globs2&#34;,
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// Common locations for mime.types files on unix.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>var typeFiles = []string{
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;/etc/mime.types&#34;,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;/etc/apache2/mime.types&#34;,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	&#34;/etc/apache/mime.types&#34;,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	&#34;/etc/httpd/conf/mime.types&#34;,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func loadMimeGlobsFile(filename string) error {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	f, err := os.Open(filename)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if err != nil {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return err
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	defer f.Close()
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	scanner := bufio.NewScanner(f)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	for scanner.Scan() {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		<span class="comment">// Each line should be of format: weight:mimetype:glob[:morefields...]</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		fields := strings.Split(scanner.Text(), &#34;:&#34;)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		if len(fields) &lt; 3 || len(fields[0]) &lt; 1 || len(fields[2]) &lt; 3 {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			continue
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		} else if fields[0][0] == &#39;#&#39; || fields[2][0] != &#39;*&#39; || fields[2][1] != &#39;.&#39; {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			continue
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		extension := fields[2][1:]
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		if strings.ContainsAny(extension, &#34;?*[&#34;) {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			<span class="comment">// Not a bare extension, but a glob. Ignore for now:</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			<span class="comment">// - we do not have an implementation for this glob</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			<span class="comment">//   syntax (translation to path/filepath.Match could</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			<span class="comment">//   be possible)</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			<span class="comment">// - support for globs with weight ordering would have</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			<span class="comment">//   performance impact to all lookups to support the</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			<span class="comment">//   rarely seen glob entries</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			<span class="comment">// - trying to match glob metacharacters literally is</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			<span class="comment">//   not useful</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			continue
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if _, ok := mimeTypes.Load(extension); ok {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ve already seen this extension.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			<span class="comment">// The file is in weight order, so we keep</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			<span class="comment">// the first entry that we see.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			continue
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		setExtensionType(extension, fields[1])
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if err := scanner.Err(); err != nil {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		panic(err)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	return nil
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func loadMimeFile(filename string) {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	f, err := os.Open(filename)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if err != nil {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	defer f.Close()
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	scanner := bufio.NewScanner(f)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	for scanner.Scan() {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		fields := strings.Fields(scanner.Text())
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		if len(fields) &lt;= 1 || fields[0][0] == &#39;#&#39; {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			continue
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		mimeType := fields[0]
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		for _, ext := range fields[1:] {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			if ext[0] == &#39;#&#39; {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>				break
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			setExtensionType(&#34;.&#34;+ext, mimeType)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	if err := scanner.Err(); err != nil {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		panic(err)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func initMimeUnix() {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	for _, filename := range mimeGlobs {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if err := loadMimeGlobsFile(filename); err == nil {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			return <span class="comment">// Stop checking more files if mimetype database is found.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// Fallback if no system-generated mimetype database exists.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	for _, filename := range typeFiles {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		loadMimeFile(filename)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>func initMimeForTests() map[string]string {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	mimeGlobs = []string{&#34;&#34;}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	typeFiles = []string{&#34;testdata/test.types&#34;}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return map[string]string{
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		&#34;.T1&#34;:  &#34;application/test&#34;,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		&#34;.t2&#34;:  &#34;text/test; charset=utf-8&#34;,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		&#34;.png&#34;: &#34;image/png&#34;,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
</pre><p><a href="type_unix.go?m=text">View as plain text</a></p>

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
