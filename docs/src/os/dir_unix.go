<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/os/dir_unix.go - Go Documentation Server</title>

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
<a href="dir_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/os">os</a>/<span class="text-muted">dir_unix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/os">os</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build aix || dragonfly || freebsd || (js &amp;&amp; wasm) || wasip1 || linux || netbsd || openbsd || solaris</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package os
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Auxiliary information if the File describes a directory</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type dirInfo struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	buf  *[]byte <span class="comment">// buffer for directory I/O</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	nbuf int     <span class="comment">// length of buf; return value from Getdirentries</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	bufp int     <span class="comment">// location of next record in buf.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>const (
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// More than 5760 to work around https://golang.org/issue/24015.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	blockSize = 8192
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>var dirBufPool = sync.Pool{
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	New: func() any {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		<span class="comment">// The buffer must be at least a block long.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		buf := make([]byte, blockSize)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		return &amp;buf
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	},
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func (d *dirInfo) close() {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	if d.buf != nil {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		dirBufPool.Put(d.buf)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		d.buf = nil
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func (f *File) readdir(n int, mode readdirMode) (names []string, dirents []DirEntry, infos []FileInfo, err error) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// If this file has no dirinfo, create one.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	if f.dirinfo == nil {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		f.dirinfo = new(dirInfo)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		f.dirinfo.buf = dirBufPool.Get().(*[]byte)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	d := f.dirinfo
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Change the meaning of n for the implementation below.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// The n above was for the public interface of &#34;if n &lt;= 0,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// Readdir returns all the FileInfo from the directory in a</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// single slice&#34;.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// But below, we use only negative to mean looping until the</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// end and positive to mean bounded, with positive</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// terminating at 0.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if n == 0 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		n = -1
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	for n != 0 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		<span class="comment">// Refill the buffer if necessary</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		if d.bufp &gt;= d.nbuf {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			d.bufp = 0
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			var errno error
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			d.nbuf, errno = f.pfd.ReadDirent(*d.buf)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			runtime.KeepAlive(f)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			if errno != nil {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>				return names, dirents, infos, &amp;PathError{Op: &#34;readdirent&#34;, Path: f.name, Err: errno}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			if d.nbuf &lt;= 0 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>				break <span class="comment">// EOF</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// Drain the buffer</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		buf := (*d.buf)[d.bufp:d.nbuf]
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		reclen, ok := direntReclen(buf)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		if !ok || reclen &gt; uint64(len(buf)) {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			break
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		rec := buf[:reclen]
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		d.bufp += int(reclen)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		ino, ok := direntIno(rec)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		if !ok {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			break
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		<span class="comment">// When building to wasip1, the host runtime might be running on Windows</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		<span class="comment">// or might expose a remote file system which does not have the concept</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		<span class="comment">// of inodes. Therefore, we cannot make the assumption that it is safe</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		<span class="comment">// to skip entries with zero inodes.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if ino == 0 &amp;&amp; runtime.GOOS != &#34;wasip1&#34; {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			continue
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		const namoff = uint64(unsafe.Offsetof(syscall.Dirent{}.Name))
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		namlen, ok := direntNamlen(rec)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if !ok || namoff+namlen &gt; uint64(len(rec)) {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			break
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		name := rec[namoff : namoff+namlen]
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		for i, c := range name {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			if c == 0 {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				name = name[:i]
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>				break
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		<span class="comment">// Check for useless names before allocating a string.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if string(name) == &#34;.&#34; || string(name) == &#34;..&#34; {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			continue
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if n &gt; 0 { <span class="comment">// see &#39;n == 0&#39; comment above</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			n--
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		if mode == readdirName {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			names = append(names, string(name))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		} else if mode == readdirDirEntry {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			de, err := newUnixDirent(f.name, string(name), direntType(rec))
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			if IsNotExist(err) {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>				<span class="comment">// File disappeared between readdir and stat.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				<span class="comment">// Treat as if it didn&#39;t exist.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				continue
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			if err != nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				return nil, dirents, nil, err
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			dirents = append(dirents, de)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		} else {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			info, err := lstat(f.name + &#34;/&#34; + string(name))
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			if IsNotExist(err) {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>				<span class="comment">// File disappeared between readdir + stat.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				<span class="comment">// Treat as if it didn&#39;t exist.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>				continue
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			if err != nil {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				return nil, nil, infos, err
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			infos = append(infos, info)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	if n &gt; 0 &amp;&amp; len(names)+len(dirents)+len(infos) == 0 {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		return nil, nil, nil, io.EOF
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	return names, dirents, infos, nil
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// readInt returns the size-bytes unsigned integer in native byte order at offset off.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func readInt(b []byte, off, size uintptr) (u uint64, ok bool) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if len(b) &lt; int(off+size) {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return 0, false
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if isBigEndian {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return readIntBE(b[off:], size), true
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	return readIntLE(b[off:], size), true
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>func readIntBE(b []byte, size uintptr) uint64 {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	switch size {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	case 1:
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		return uint64(b[0])
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	case 2:
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		_ = b[1] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return uint64(b[1]) | uint64(b[0])&lt;&lt;8
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	case 4:
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		_ = b[3] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return uint64(b[3]) | uint64(b[2])&lt;&lt;8 | uint64(b[1])&lt;&lt;16 | uint64(b[0])&lt;&lt;24
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	case 8:
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		_ = b[7] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		return uint64(b[7]) | uint64(b[6])&lt;&lt;8 | uint64(b[5])&lt;&lt;16 | uint64(b[4])&lt;&lt;24 |
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			uint64(b[3])&lt;&lt;32 | uint64(b[2])&lt;&lt;40 | uint64(b[1])&lt;&lt;48 | uint64(b[0])&lt;&lt;56
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	default:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		panic(&#34;syscall: readInt with unsupported size&#34;)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func readIntLE(b []byte, size uintptr) uint64 {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	switch size {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	case 1:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return uint64(b[0])
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	case 2:
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		_ = b[1] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		return uint64(b[0]) | uint64(b[1])&lt;&lt;8
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	case 4:
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		_ = b[3] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		return uint64(b[0]) | uint64(b[1])&lt;&lt;8 | uint64(b[2])&lt;&lt;16 | uint64(b[3])&lt;&lt;24
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	case 8:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		_ = b[7] <span class="comment">// bounds check hint to compiler; see golang.org/issue/14808</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		return uint64(b[0]) | uint64(b[1])&lt;&lt;8 | uint64(b[2])&lt;&lt;16 | uint64(b[3])&lt;&lt;24 |
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			uint64(b[4])&lt;&lt;32 | uint64(b[5])&lt;&lt;40 | uint64(b[6])&lt;&lt;48 | uint64(b[7])&lt;&lt;56
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	default:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		panic(&#34;syscall: readInt with unsupported size&#34;)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
</pre><p><a href="dir_unix.go?m=text">View as plain text</a></p>

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
