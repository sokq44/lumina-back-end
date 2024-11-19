<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/archive/tar/common.go - Go Documentation Server</title>

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
<a href="common.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/archive">archive</a>/<a href="http://localhost:8080/src/archive/tar">tar</a>/<span class="text-muted">common.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/archive/tar">archive/tar</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package tar implements access to tar archives.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Tape archives (tar) are a file format for storing a sequence of files that</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// can be read and written in a streaming manner.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// This package aims to cover most variations of the format,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// including those produced by GNU and BSD tar tools.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>package tar
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>import (
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;internal/godebug&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;io/fs&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;path&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// BUG: Use of the Uid and Gid fields in Header could overflow on 32-bit</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// architectures. If a large value is encountered when decoding, the result</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// stored in Header will be the truncated version.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>var tarinsecurepath = godebug.New(&#34;tarinsecurepath&#34;)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>var (
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	ErrHeader          = errors.New(&#34;archive/tar: invalid tar header&#34;)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	ErrWriteTooLong    = errors.New(&#34;archive/tar: write too long&#34;)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	ErrFieldTooLong    = errors.New(&#34;archive/tar: header field too long&#34;)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	ErrWriteAfterClose = errors.New(&#34;archive/tar: write after close&#34;)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	ErrInsecurePath    = errors.New(&#34;archive/tar: insecure file path&#34;)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	errMissData        = errors.New(&#34;archive/tar: sparse file references non-existent data&#34;)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	errUnrefData       = errors.New(&#34;archive/tar: sparse file contains unreferenced data&#34;)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	errWriteHole       = errors.New(&#34;archive/tar: write non-NUL byte in sparse hole&#34;)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>type headerError []string
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (he headerError) Error() string {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	const prefix = &#34;archive/tar: cannot encode header&#34;
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	var ss []string
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	for _, s := range he {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		if s != &#34;&#34; {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			ss = append(ss, s)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if len(ss) == 0 {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return prefix
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%s: %v&#34;, prefix, strings.Join(ss, &#34;; and &#34;))
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// Type flags for Header.Typeflag.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>const (
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;0&#39; indicates a regular file.</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	TypeReg = &#39;0&#39;
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// Deprecated: Use TypeReg instead.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	TypeRegA = &#39;\x00&#39;
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;1&#39; to &#39;6&#39; are header-only flags and may not have a data body.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	TypeLink    = &#39;1&#39; <span class="comment">// Hard link</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	TypeSymlink = &#39;2&#39; <span class="comment">// Symbolic link</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	TypeChar    = &#39;3&#39; <span class="comment">// Character device node</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	TypeBlock   = &#39;4&#39; <span class="comment">// Block device node</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	TypeDir     = &#39;5&#39; <span class="comment">// Directory</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	TypeFifo    = &#39;6&#39; <span class="comment">// FIFO node</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;7&#39; is reserved.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	TypeCont = &#39;7&#39;
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;x&#39; is used by the PAX format to store key-value records that</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// are only relevant to the next file.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// This package transparently handles these types.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	TypeXHeader = &#39;x&#39;
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;g&#39; is used by the PAX format to store key-value records that</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// are relevant to all subsequent files.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// This package only supports parsing and composing such headers,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// but does not currently support persisting the global state across files.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	TypeXGlobalHeader = &#39;g&#39;
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// Type &#39;S&#39; indicates a sparse file in the GNU format.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	TypeGNUSparse = &#39;S&#39;
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// Types &#39;L&#39; and &#39;K&#39; are used by the GNU format for a meta file</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// used to store the path or link name for the next file.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// This package transparently handles these types.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	TypeGNULongName = &#39;L&#39;
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	TypeGNULongLink = &#39;K&#39;
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Keywords for PAX extended header records.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>const (
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	paxNone     = &#34;&#34; <span class="comment">// Indicates that no PAX key is suitable</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	paxPath     = &#34;path&#34;
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	paxLinkpath = &#34;linkpath&#34;
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	paxSize     = &#34;size&#34;
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	paxUid      = &#34;uid&#34;
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	paxGid      = &#34;gid&#34;
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	paxUname    = &#34;uname&#34;
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	paxGname    = &#34;gname&#34;
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	paxMtime    = &#34;mtime&#34;
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	paxAtime    = &#34;atime&#34;
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	paxCtime    = &#34;ctime&#34;   <span class="comment">// Removed from later revision of PAX spec, but was valid</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	paxCharset  = &#34;charset&#34; <span class="comment">// Currently unused</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	paxComment  = &#34;comment&#34; <span class="comment">// Currently unused</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	paxSchilyXattr = &#34;SCHILY.xattr.&#34;
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// Keywords for GNU sparse files in a PAX extended header.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	paxGNUSparse          = &#34;GNU.sparse.&#34;
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	paxGNUSparseNumBlocks = &#34;GNU.sparse.numblocks&#34;
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	paxGNUSparseOffset    = &#34;GNU.sparse.offset&#34;
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	paxGNUSparseNumBytes  = &#34;GNU.sparse.numbytes&#34;
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	paxGNUSparseMap       = &#34;GNU.sparse.map&#34;
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	paxGNUSparseName      = &#34;GNU.sparse.name&#34;
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	paxGNUSparseMajor     = &#34;GNU.sparse.major&#34;
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	paxGNUSparseMinor     = &#34;GNU.sparse.minor&#34;
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	paxGNUSparseSize      = &#34;GNU.sparse.size&#34;
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	paxGNUSparseRealSize  = &#34;GNU.sparse.realsize&#34;
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// basicKeys is a set of the PAX keys for which we have built-in support.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// This does not contain &#34;charset&#34; or &#34;comment&#34;, which are both PAX-specific,</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// so adding them as first-class features of Header is unlikely.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// Users can use the PAXRecords field to set it themselves.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>var basicKeys = map[string]bool{
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	paxPath: true, paxLinkpath: true, paxSize: true, paxUid: true, paxGid: true,
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	paxUname: true, paxGname: true, paxMtime: true, paxAtime: true, paxCtime: true,
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// A Header represents a single header in a tar archive.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// Some fields may not be populated.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// For forward compatibility, users that retrieve a Header from Reader.Next,</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// mutate it in some ways, and then pass it back to Writer.WriteHeader</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// should do so by creating a new Header and copying the fields</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// that they are interested in preserving.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>type Header struct {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// Typeflag is the type of header entry.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// The zero value is automatically promoted to either TypeReg or TypeDir</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// depending on the presence of a trailing slash in Name.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	Typeflag byte
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	Name     string <span class="comment">// Name of file entry</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	Linkname string <span class="comment">// Target name of link (valid for TypeLink or TypeSymlink)</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	Size  int64  <span class="comment">// Logical file size in bytes</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	Mode  int64  <span class="comment">// Permission and mode bits</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	Uid   int    <span class="comment">// User ID of owner</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	Gid   int    <span class="comment">// Group ID of owner</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	Uname string <span class="comment">// User name of owner</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	Gname string <span class="comment">// Group name of owner</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// If the Format is unspecified, then Writer.WriteHeader rounds ModTime</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// to the nearest second and ignores the AccessTime and ChangeTime fields.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// To use AccessTime or ChangeTime, specify the Format as PAX or GNU.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// To use sub-second resolution, specify the Format as PAX.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	ModTime    time.Time <span class="comment">// Modification time</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	AccessTime time.Time <span class="comment">// Access time (requires either PAX or GNU support)</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	ChangeTime time.Time <span class="comment">// Change time (requires either PAX or GNU support)</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	Devmajor int64 <span class="comment">// Major device number (valid for TypeChar or TypeBlock)</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	Devminor int64 <span class="comment">// Minor device number (valid for TypeChar or TypeBlock)</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Xattrs stores extended attributes as PAX records under the</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// &#34;SCHILY.xattr.&#34; namespace.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// The following are semantically equivalent:</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">//  h.Xattrs[key] = value</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">//  h.PAXRecords[&#34;SCHILY.xattr.&#34;+key] = value</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// When Writer.WriteHeader is called, the contents of Xattrs will take</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// precedence over those in PAXRecords.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// Deprecated: Use PAXRecords instead.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	Xattrs map[string]string
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// PAXRecords is a map of PAX extended header records.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">// User-defined records should have keys of the following form:</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//	VENDOR.keyword</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Where VENDOR is some namespace in all uppercase, and keyword may</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// not contain the &#39;=&#39; character (e.g., &#34;GOLANG.pkg.version&#34;).</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// The key and value should be non-empty UTF-8 strings.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// When Writer.WriteHeader is called, PAX records derived from the</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// other fields in Header take precedence over PAXRecords.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	PAXRecords map[string]string
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// Format specifies the format of the tar header.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// This is set by Reader.Next as a best-effort guess at the format.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// Since the Reader liberally reads some non-compliant files,</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// it is possible for this to be FormatUnknown.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// If the format is unspecified when Writer.WriteHeader is called,</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">// then it uses the first format (in the order of USTAR, PAX, GNU)</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// capable of encoding this Header (see Format).</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	Format Format
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// sparseEntry represents a Length-sized fragment at Offset in the file.</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>type sparseEntry struct{ Offset, Length int64 }
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>func (s sparseEntry) endOffset() int64 { return s.Offset + s.Length }
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// A sparse file can be represented as either a sparseDatas or a sparseHoles.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// As long as the total size is known, they are equivalent and one can be</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// converted to the other form and back. The various tar formats with sparse</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// file support represent sparse files in the sparseDatas form. That is, they</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// specify the fragments in the file that has data, and treat everything else as</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// having zero bytes. As such, the encoding and decoding logic in this package</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// deals with sparseDatas.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// However, the external API uses sparseHoles instead of sparseDatas because the</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// zero value of sparseHoles logically represents a normal file (i.e., there are</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// no holes in it). On the other hand, the zero value of sparseDatas implies</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// that the file has no data in it, which is rather odd.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span><span class="comment">// As an example, if the underlying raw file contains the 10-byte data:</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">//	var compactFile = &#34;abcdefgh&#34;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span><span class="comment">// And the sparse map has the following entries:</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">//	var spd sparseDatas = []sparseEntry{</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">//		{Offset: 2,  Length: 5},  // Data fragment for 2..6</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">//		{Offset: 18, Length: 3},  // Data fragment for 18..20</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">//	var sph sparseHoles = []sparseEntry{</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span><span class="comment">//		{Offset: 0,  Length: 2},  // Hole fragment for 0..1</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">//		{Offset: 7,  Length: 11}, // Hole fragment for 7..17</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">//		{Offset: 21, Length: 4},  // Hole fragment for 21..24</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// Then the content of the resulting sparse file with a Header.Size of 25 is:</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">//	var sparseFile = &#34;\x00&#34;*2 + &#34;abcde&#34; + &#34;\x00&#34;*11 + &#34;fgh&#34; + &#34;\x00&#34;*4</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>type (
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	sparseDatas []sparseEntry
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	sparseHoles []sparseEntry
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// validateSparseEntries reports whether sp is a valid sparse map.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// It does not matter whether sp represents data fragments or hole fragments.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>func validateSparseEntries(sp []sparseEntry, size int64) bool {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// Validate all sparse entries. These are the same checks as performed by</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">// the BSD tar utility.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	if size &lt; 0 {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		return false
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	var pre sparseEntry
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	for _, cur := range sp {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		switch {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		case cur.Offset &lt; 0 || cur.Length &lt; 0:
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			return false <span class="comment">// Negative values are never okay</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		case cur.Offset &gt; math.MaxInt64-cur.Length:
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			return false <span class="comment">// Integer overflow with large length</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		case cur.endOffset() &gt; size:
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			return false <span class="comment">// Region extends beyond the actual size</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		case pre.endOffset() &gt; cur.Offset:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			return false <span class="comment">// Regions cannot overlap and must be in order</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		pre = cur
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	return true
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// alignSparseEntries mutates src and returns dst where each fragment&#39;s</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// starting offset is aligned up to the nearest block edge, and each</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// ending offset is aligned down to the nearest block edge.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// Even though the Go tar Reader and the BSD tar utility can handle entries</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// with arbitrary offsets and lengths, the GNU tar utility can only handle</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// offsets and lengths that are multiples of blockSize.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>func alignSparseEntries(src []sparseEntry, size int64) []sparseEntry {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	dst := src[:0]
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	for _, s := range src {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		pos, end := s.Offset, s.endOffset()
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		pos += blockPadding(+pos) <span class="comment">// Round-up to nearest blockSize</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		if end != size {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			end -= blockPadding(-end) <span class="comment">// Round-down to nearest blockSize</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		if pos &lt; end {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			dst = append(dst, sparseEntry{Offset: pos, Length: end - pos})
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	return dst
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// invertSparseEntries converts a sparse map from one form to the other.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// If the input is sparseHoles, then it will output sparseDatas and vice-versa.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// The input must have been already validated.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">// This function mutates src and returns a normalized map where:</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">//   - adjacent fragments are coalesced together</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">//   - only the last fragment may be empty</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">//   - the endOffset of the last fragment is the total size</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>func invertSparseEntries(src []sparseEntry, size int64) []sparseEntry {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	dst := src[:0]
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	var pre sparseEntry
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	for _, cur := range src {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		if cur.Length == 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			continue <span class="comment">// Skip empty fragments</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		pre.Length = cur.Offset - pre.Offset
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if pre.Length &gt; 0 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			dst = append(dst, pre) <span class="comment">// Only add non-empty fragments</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		pre.Offset = cur.endOffset()
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	pre.Length = size - pre.Offset <span class="comment">// Possibly the only empty fragment</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	return append(dst, pre)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">// fileState tracks the number of logical (includes sparse holes) and physical</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// (actual in tar archive) bytes remaining for the current file.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// Invariant: logicalRemaining &gt;= physicalRemaining</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>type fileState interface {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	logicalRemaining() int64
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	physicalRemaining() int64
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// allowedFormats determines which formats can be used.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// The value returned is the logical OR of multiple possible formats.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// If the value is FormatUnknown, then the input Header cannot be encoded</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// and an error is returned explaining why.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// As a by-product of checking the fields, this function returns paxHdrs, which</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// contain all fields that could not be directly encoded.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// A value receiver ensures that this method does not mutate the source Header.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>func (h Header) allowedFormats() (format Format, paxHdrs map[string]string, err error) {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	format = FormatUSTAR | FormatPAX | FormatGNU
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	paxHdrs = make(map[string]string)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	var whyNoUSTAR, whyNoPAX, whyNoGNU string
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	var preferPAX bool <span class="comment">// Prefer PAX over USTAR</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	verifyString := func(s string, size int, name, paxKey string) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		<span class="comment">// NUL-terminator is optional for path and linkpath.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// Technically, it is required for uname and gname,</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		<span class="comment">// but neither GNU nor BSD tar checks for it.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		tooLong := len(s) &gt; size
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		allowLongGNU := paxKey == paxPath || paxKey == paxLinkpath
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		if hasNUL(s) || (tooLong &amp;&amp; !allowLongGNU) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			whyNoGNU = fmt.Sprintf(&#34;GNU cannot encode %s=%q&#34;, name, s)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			format.mustNotBe(FormatGNU)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		if !isASCII(s) || tooLong {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			canSplitUSTAR := paxKey == paxPath
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			if _, _, ok := splitUSTARPath(s); !canSplitUSTAR || !ok {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>				whyNoUSTAR = fmt.Sprintf(&#34;USTAR cannot encode %s=%q&#34;, name, s)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>				format.mustNotBe(FormatUSTAR)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>			}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			if paxKey == paxNone {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				whyNoPAX = fmt.Sprintf(&#34;PAX cannot encode %s=%q&#34;, name, s)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>				format.mustNotBe(FormatPAX)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			} else {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>				paxHdrs[paxKey] = s
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			}
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		if v, ok := h.PAXRecords[paxKey]; ok &amp;&amp; v == s {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			paxHdrs[paxKey] = v
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	verifyNumeric := func(n int64, size int, name, paxKey string) {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if !fitsInBase256(size, n) {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			whyNoGNU = fmt.Sprintf(&#34;GNU cannot encode %s=%d&#34;, name, n)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			format.mustNotBe(FormatGNU)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		if !fitsInOctal(size, n) {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			whyNoUSTAR = fmt.Sprintf(&#34;USTAR cannot encode %s=%d&#34;, name, n)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			format.mustNotBe(FormatUSTAR)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			if paxKey == paxNone {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				whyNoPAX = fmt.Sprintf(&#34;PAX cannot encode %s=%d&#34;, name, n)
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				format.mustNotBe(FormatPAX)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>			} else {
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>				paxHdrs[paxKey] = strconv.FormatInt(n, 10)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		if v, ok := h.PAXRecords[paxKey]; ok &amp;&amp; v == strconv.FormatInt(n, 10) {
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			paxHdrs[paxKey] = v
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	}
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	verifyTime := func(ts time.Time, size int, name, paxKey string) {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		if ts.IsZero() {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			return <span class="comment">// Always okay</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		}
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		if !fitsInBase256(size, ts.Unix()) {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			whyNoGNU = fmt.Sprintf(&#34;GNU cannot encode %s=%v&#34;, name, ts)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			format.mustNotBe(FormatGNU)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		isMtime := paxKey == paxMtime
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		fitsOctal := fitsInOctal(size, ts.Unix())
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if (isMtime &amp;&amp; !fitsOctal) || !isMtime {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			whyNoUSTAR = fmt.Sprintf(&#34;USTAR cannot encode %s=%v&#34;, name, ts)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			format.mustNotBe(FormatUSTAR)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		needsNano := ts.Nanosecond() != 0
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		if !isMtime || !fitsOctal || needsNano {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			preferPAX = true <span class="comment">// USTAR may truncate sub-second measurements</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			if paxKey == paxNone {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>				whyNoPAX = fmt.Sprintf(&#34;PAX cannot encode %s=%v&#34;, name, ts)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				format.mustNotBe(FormatPAX)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			} else {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>				paxHdrs[paxKey] = formatPAXTime(ts)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		if v, ok := h.PAXRecords[paxKey]; ok &amp;&amp; v == formatPAXTime(ts) {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			paxHdrs[paxKey] = v
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	<span class="comment">// Check basic fields.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	var blk block
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	v7 := blk.toV7()
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	ustar := blk.toUSTAR()
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	gnu := blk.toGNU()
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	verifyString(h.Name, len(v7.name()), &#34;Name&#34;, paxPath)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	verifyString(h.Linkname, len(v7.linkName()), &#34;Linkname&#34;, paxLinkpath)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	verifyString(h.Uname, len(ustar.userName()), &#34;Uname&#34;, paxUname)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	verifyString(h.Gname, len(ustar.groupName()), &#34;Gname&#34;, paxGname)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	verifyNumeric(h.Mode, len(v7.mode()), &#34;Mode&#34;, paxNone)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	verifyNumeric(int64(h.Uid), len(v7.uid()), &#34;Uid&#34;, paxUid)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	verifyNumeric(int64(h.Gid), len(v7.gid()), &#34;Gid&#34;, paxGid)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	verifyNumeric(h.Size, len(v7.size()), &#34;Size&#34;, paxSize)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	verifyNumeric(h.Devmajor, len(ustar.devMajor()), &#34;Devmajor&#34;, paxNone)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	verifyNumeric(h.Devminor, len(ustar.devMinor()), &#34;Devminor&#34;, paxNone)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	verifyTime(h.ModTime, len(v7.modTime()), &#34;ModTime&#34;, paxMtime)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	verifyTime(h.AccessTime, len(gnu.accessTime()), &#34;AccessTime&#34;, paxAtime)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	verifyTime(h.ChangeTime, len(gnu.changeTime()), &#34;ChangeTime&#34;, paxCtime)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// Check for header-only types.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	var whyOnlyPAX, whyOnlyGNU string
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	switch h.Typeflag {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	case TypeReg, TypeChar, TypeBlock, TypeFifo, TypeGNUSparse:
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		<span class="comment">// Exclude TypeLink and TypeSymlink, since they may reference directories.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		if strings.HasSuffix(h.Name, &#34;/&#34;) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			return FormatUnknown, nil, headerError{&#34;filename may not have trailing slash&#34;}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	case TypeXHeader, TypeGNULongName, TypeGNULongLink:
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		return FormatUnknown, nil, headerError{&#34;cannot manually encode TypeXHeader, TypeGNULongName, or TypeGNULongLink headers&#34;}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	case TypeXGlobalHeader:
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		h2 := Header{Name: h.Name, Typeflag: h.Typeflag, Xattrs: h.Xattrs, PAXRecords: h.PAXRecords, Format: h.Format}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		if !reflect.DeepEqual(h, h2) {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			return FormatUnknown, nil, headerError{&#34;only PAXRecords should be set for TypeXGlobalHeader&#34;}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		whyOnlyPAX = &#34;only PAX supports TypeXGlobalHeader&#34;
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		format.mayOnlyBe(FormatPAX)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	if !isHeaderOnlyType(h.Typeflag) &amp;&amp; h.Size &lt; 0 {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		return FormatUnknown, nil, headerError{&#34;negative size on header-only type&#34;}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// Check PAX records.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	if len(h.Xattrs) &gt; 0 {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		for k, v := range h.Xattrs {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			paxHdrs[paxSchilyXattr+k] = v
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		whyOnlyPAX = &#34;only PAX supports Xattrs&#34;
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		format.mayOnlyBe(FormatPAX)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	if len(h.PAXRecords) &gt; 0 {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		for k, v := range h.PAXRecords {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			switch _, exists := paxHdrs[k]; {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			case exists:
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>				continue <span class="comment">// Do not overwrite existing records</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			case h.Typeflag == TypeXGlobalHeader:
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>				paxHdrs[k] = v <span class="comment">// Copy all records</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			case !basicKeys[k] &amp;&amp; !strings.HasPrefix(k, paxGNUSparse):
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>				paxHdrs[k] = v <span class="comment">// Ignore local records that may conflict</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		whyOnlyPAX = &#34;only PAX supports PAXRecords&#34;
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		format.mayOnlyBe(FormatPAX)
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	for k, v := range paxHdrs {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		if !validPAXRecord(k, v) {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			return FormatUnknown, nil, headerError{fmt.Sprintf(&#34;invalid PAX record: %q&#34;, k+&#34; = &#34;+v)}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	<span class="comment">// TODO(dsnet): Re-enable this when adding sparse support.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	<span class="comment">// See https://golang.org/issue/22735</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	<span class="comment">/*
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		// Check sparse files.
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		if len(h.SparseHoles) &gt; 0 || h.Typeflag == TypeGNUSparse {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			if isHeaderOnlyType(h.Typeflag) {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>				return FormatUnknown, nil, headerError{&#34;header-only type cannot be sparse&#34;}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			if !validateSparseEntries(h.SparseHoles, h.Size) {
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>				return FormatUnknown, nil, headerError{&#34;invalid sparse holes&#34;}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			if h.Typeflag == TypeGNUSparse {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>				whyOnlyGNU = &#34;only GNU supports TypeGNUSparse&#34;
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>				format.mayOnlyBe(FormatGNU)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			} else {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				whyNoGNU = &#34;GNU supports sparse files only with TypeGNUSparse&#34;
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>				format.mustNotBe(FormatGNU)
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			whyNoUSTAR = &#34;USTAR does not support sparse files&#34;
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			format.mustNotBe(FormatUSTAR)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	*/</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	<span class="comment">// Check desired format.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	if wantFormat := h.Format; wantFormat != FormatUnknown {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		if wantFormat.has(FormatPAX) &amp;&amp; !preferPAX {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			wantFormat.mayBe(FormatUSTAR) <span class="comment">// PAX implies USTAR allowed too</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		format.mayOnlyBe(wantFormat) <span class="comment">// Set union of formats allowed and format wanted</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	if format == FormatUnknown {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		switch h.Format {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		case FormatUSTAR:
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			err = headerError{&#34;Format specifies USTAR&#34;, whyNoUSTAR, whyOnlyPAX, whyOnlyGNU}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		case FormatPAX:
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			err = headerError{&#34;Format specifies PAX&#34;, whyNoPAX, whyOnlyGNU}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>		case FormatGNU:
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			err = headerError{&#34;Format specifies GNU&#34;, whyNoGNU, whyOnlyPAX}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		default:
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>			err = headerError{whyNoUSTAR, whyNoPAX, whyNoGNU, whyOnlyPAX, whyOnlyGNU}
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	return format, paxHdrs, err
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// FileInfo returns an fs.FileInfo for the Header.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>func (h *Header) FileInfo() fs.FileInfo {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	return headerFileInfo{h}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// headerFileInfo implements fs.FileInfo.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>type headerFileInfo struct {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	h *Header
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>func (fi headerFileInfo) Size() int64        { return fi.h.Size }
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>func (fi headerFileInfo) IsDir() bool        { return fi.Mode().IsDir() }
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>func (fi headerFileInfo) ModTime() time.Time { return fi.h.ModTime }
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>func (fi headerFileInfo) Sys() any           { return fi.h }
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// Name returns the base name of the file.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>func (fi headerFileInfo) Name() string {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	if fi.IsDir() {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		return path.Base(path.Clean(fi.h.Name))
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	return path.Base(fi.h.Name)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">// Mode returns the permission and mode bits for the headerFileInfo.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>func (fi headerFileInfo) Mode() (mode fs.FileMode) {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	<span class="comment">// Set file permission bits.</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	mode = fs.FileMode(fi.h.Mode).Perm()
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	<span class="comment">// Set setuid, setgid and sticky bits.</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	if fi.h.Mode&amp;c_ISUID != 0 {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		mode |= fs.ModeSetuid
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if fi.h.Mode&amp;c_ISGID != 0 {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		mode |= fs.ModeSetgid
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if fi.h.Mode&amp;c_ISVTX != 0 {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		mode |= fs.ModeSticky
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	<span class="comment">// Set file mode bits; clear perm, setuid, setgid, and sticky bits.</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	switch m := fs.FileMode(fi.h.Mode) &amp;^ 07777; m {
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	case c_ISDIR:
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		mode |= fs.ModeDir
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	case c_ISFIFO:
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		mode |= fs.ModeNamedPipe
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	case c_ISLNK:
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		mode |= fs.ModeSymlink
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	case c_ISBLK:
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		mode |= fs.ModeDevice
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	case c_ISCHR:
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		mode |= fs.ModeDevice
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		mode |= fs.ModeCharDevice
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	case c_ISSOCK:
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		mode |= fs.ModeSocket
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	switch fi.h.Typeflag {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	case TypeSymlink:
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		mode |= fs.ModeSymlink
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	case TypeChar:
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		mode |= fs.ModeDevice
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		mode |= fs.ModeCharDevice
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	case TypeBlock:
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		mode |= fs.ModeDevice
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	case TypeDir:
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		mode |= fs.ModeDir
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	case TypeFifo:
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		mode |= fs.ModeNamedPipe
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	return mode
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>func (fi headerFileInfo) String() string {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	return fs.FormatFileInfo(fi)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// sysStat, if non-nil, populates h from system-dependent fields of fi.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>var sysStat func(fi fs.FileInfo, h *Header) error
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>const (
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	<span class="comment">// Mode constants from the USTAR spec:</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	<span class="comment">// See http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	c_ISUID = 04000 <span class="comment">// Set uid</span>
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	c_ISGID = 02000 <span class="comment">// Set gid</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	c_ISVTX = 01000 <span class="comment">// Save text (sticky bit)</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	<span class="comment">// Common Unix mode constants; these are not defined in any common tar standard.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	<span class="comment">// Header.FileInfo understands these, but FileInfoHeader will never produce these.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	c_ISDIR  = 040000  <span class="comment">// Directory</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	c_ISFIFO = 010000  <span class="comment">// FIFO</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	c_ISREG  = 0100000 <span class="comment">// Regular file</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	c_ISLNK  = 0120000 <span class="comment">// Symbolic link</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	c_ISBLK  = 060000  <span class="comment">// Block special file</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	c_ISCHR  = 020000  <span class="comment">// Character special file</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	c_ISSOCK = 0140000 <span class="comment">// Socket</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>)
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// FileInfoHeader creates a partially-populated [Header] from fi.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">// If fi describes a symlink, FileInfoHeader records link as the link target.</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">// If fi describes a directory, a slash is appended to the name.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span><span class="comment">// Since fs.FileInfo&#39;s Name method only returns the base name of</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span><span class="comment">// the file it describes, it may be necessary to modify Header.Name</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span><span class="comment">// to provide the full path name of the file.</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>func FileInfoHeader(fi fs.FileInfo, link string) (*Header, error) {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	if fi == nil {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		return nil, errors.New(&#34;archive/tar: FileInfo is nil&#34;)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	fm := fi.Mode()
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	h := &amp;Header{
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		Name:    fi.Name(),
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		ModTime: fi.ModTime(),
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		Mode:    int64(fm.Perm()), <span class="comment">// or&#39;d with c_IS* constants later</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	switch {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	case fm.IsRegular():
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		h.Typeflag = TypeReg
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		h.Size = fi.Size()
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	case fi.IsDir():
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		h.Typeflag = TypeDir
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		h.Name += &#34;/&#34;
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	case fm&amp;fs.ModeSymlink != 0:
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		h.Typeflag = TypeSymlink
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		h.Linkname = link
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	case fm&amp;fs.ModeDevice != 0:
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		if fm&amp;fs.ModeCharDevice != 0 {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			h.Typeflag = TypeChar
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		} else {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			h.Typeflag = TypeBlock
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		}
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	case fm&amp;fs.ModeNamedPipe != 0:
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		h.Typeflag = TypeFifo
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	case fm&amp;fs.ModeSocket != 0:
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;archive/tar: sockets not supported&#34;)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	default:
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;archive/tar: unknown file mode %v&#34;, fm)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	if fm&amp;fs.ModeSetuid != 0 {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		h.Mode |= c_ISUID
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	}
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	if fm&amp;fs.ModeSetgid != 0 {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		h.Mode |= c_ISGID
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if fm&amp;fs.ModeSticky != 0 {
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		h.Mode |= c_ISVTX
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	<span class="comment">// If possible, populate additional fields from OS-specific</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	<span class="comment">// FileInfo fields.</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	if sys, ok := fi.Sys().(*Header); ok {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		<span class="comment">// This FileInfo came from a Header (not the OS). Use the</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		<span class="comment">// original Header to populate all remaining fields.</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		h.Uid = sys.Uid
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		h.Gid = sys.Gid
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		h.Uname = sys.Uname
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		h.Gname = sys.Gname
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		h.AccessTime = sys.AccessTime
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		h.ChangeTime = sys.ChangeTime
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>		if sys.Xattrs != nil {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>			h.Xattrs = make(map[string]string)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>			for k, v := range sys.Xattrs {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>				h.Xattrs[k] = v
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			}
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		if sys.Typeflag == TypeLink {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			<span class="comment">// hard link</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			h.Typeflag = TypeLink
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			h.Size = 0
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			h.Linkname = sys.Linkname
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		if sys.PAXRecords != nil {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			h.PAXRecords = make(map[string]string)
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			for k, v := range sys.PAXRecords {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				h.PAXRecords[k] = v
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	if sysStat != nil {
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		return h, sysStat(fi, h)
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	}
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	return h, nil
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span><span class="comment">// isHeaderOnlyType checks if the given type flag is of the type that has no</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span><span class="comment">// data section even if a size is specified.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>func isHeaderOnlyType(flag byte) bool {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	switch flag {
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	case TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, TypeFifo:
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		return true
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	default:
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		return false
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>}
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
</pre><p><a href="common.go?m=text">View as plain text</a></p>

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
