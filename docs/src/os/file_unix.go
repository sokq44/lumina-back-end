<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/os/file_unix.go - Go Documentation Server</title>

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
<a href="file_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/os">os</a>/<span class="text-muted">file_unix.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || (js &amp;&amp; wasm) || wasip1</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package os
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/poll&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/syscall/unix&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;io/fs&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	_ &#34;unsafe&#34; <span class="comment">// for go:linkname</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const _UTIME_OMIT = unix.UTIME_OMIT
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// fixLongPath is a noop on non-Windows platforms.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func fixLongPath(path string) string {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	return path
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func rename(oldname, newname string) error {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	fi, err := Lstat(newname)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if err == nil &amp;&amp; fi.IsDir() {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		<span class="comment">// There are two independent errors this function can return:</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		<span class="comment">// one for a bad oldname, and one for a bad newname.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		<span class="comment">// At this point we&#39;ve determined the newname is bad.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		<span class="comment">// But just in case oldname is also bad, prioritize returning</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		<span class="comment">// the oldname error because that&#39;s what we did historically.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		<span class="comment">// However, if the old name and new name are not the same, yet</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		<span class="comment">// they refer to the same file, it implies a case-only</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		<span class="comment">// rename on a case-insensitive filesystem, which is ok.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		if ofi, err := Lstat(oldname); err != nil {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>			if pe, ok := err.(*PathError); ok {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>				err = pe.Err
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>			}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			return &amp;LinkError{&#34;rename&#34;, oldname, newname, err}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		} else if newname == oldname || !SameFile(fi, ofi) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>			return &amp;LinkError{&#34;rename&#34;, oldname, newname, syscall.EEXIST}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	err = ignoringEINTR(func() error {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		return syscall.Rename(oldname, newname)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	})
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if err != nil {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		return &amp;LinkError{&#34;rename&#34;, oldname, newname, err}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	return nil
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// file is the real representation of *File.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// The extra level of indirection ensures that no clients of os</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// can overwrite this data, which could cause the finalizer</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// to close the wrong file descriptor.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>type file struct {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	pfd         poll.FD
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	name        string
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	dirinfo     *dirInfo <span class="comment">// nil unless directory being read</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	nonblock    bool     <span class="comment">// whether we set nonblocking mode</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	stdoutOrErr bool     <span class="comment">// whether this is stdout or stderr</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	appendMode  bool     <span class="comment">// whether file is opened for appending</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// Fd returns the integer Unix file descriptor referencing the open file.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// If f is closed, the file descriptor becomes invalid.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// If f is garbage collected, a finalizer may close the file descriptor,</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// making it invalid; see runtime.SetFinalizer for more information on when</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// a finalizer might be run. On Unix systems this will cause the SetDeadline</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">// methods to stop working.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// Because file descriptors can be reused, the returned file descriptor may</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// only be closed through the Close method of f, or by its finalizer during</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// garbage collection. Otherwise, during garbage collection the finalizer</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// may close an unrelated file descriptor with the same (reused) number.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// As an alternative, see the f.SyscallConn method.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (f *File) Fd() uintptr {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if f == nil {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return ^(uintptr(0))
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// If we put the file descriptor into nonblocking mode,</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// then set it to blocking mode before we return it,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// because historically we have always returned a descriptor</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// opened in blocking mode. The File will continue to work,</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// but any blocking operation will tie up a thread.</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if f.nonblock {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		f.pfd.SetBlocking()
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return uintptr(f.pfd.Sysfd)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// NewFile returns a new File with the given file descriptor and</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// name. The returned value will be nil if fd is not a valid file</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// descriptor. On Unix systems, if the file descriptor is in</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// non-blocking mode, NewFile will attempt to return a pollable File</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// (one for which the SetDeadline methods work).</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// After passing it to NewFile, fd may become invalid under the same</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// conditions described in the comments of the Fd method, and the same</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// constraints apply.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func NewFile(fd uintptr, name string) *File {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	fdi := int(fd)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if fdi &lt; 0 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return nil
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	kind := kindNewFile
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	appendMode := false
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if flags, err := unix.Fcntl(fdi, syscall.F_GETFL, 0); err == nil {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		if unix.HasNonblockFlag(flags) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			kind = kindNonBlock
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		appendMode = flags&amp;syscall.O_APPEND != 0
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	f := newFile(fdi, name, kind)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	f.appendMode = appendMode
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return f
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// net_newUnixFile is a hidden entry point called by net.conn.File.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// This is used so that a nonblocking network connection will become</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// blocking if code calls the Fd method. We don&#39;t want that for direct</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// calls to NewFile: passing a nonblocking descriptor to NewFile should</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// remain nonblocking if you get it back using Fd. But for net.conn.File</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// the call to NewFile is hidden from the user. Historically in that case</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// the Fd method has returned a blocking descriptor, and we want to</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// retain that behavior because existing code expects it and depends on it.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">//go:linkname net_newUnixFile net.newUnixFile</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>func net_newUnixFile(fd int, name string) *File {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if fd &lt; 0 {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		panic(&#34;invalid FD&#34;)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	f := newFile(fd, name, kindNonBlock)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	f.nonblock = true <span class="comment">// tell Fd to return blocking descriptor</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	return f
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// newFileKind describes the kind of file to newFile.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>type newFileKind int
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>const (
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// kindNewFile means that the descriptor was passed to us via NewFile.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	kindNewFile newFileKind = iota
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">// kindOpenFile means that the descriptor was opened using</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Open, Create, or OpenFile (without O_NONBLOCK).</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	kindOpenFile
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// kindPipe means that the descriptor was opened using Pipe.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	kindPipe
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// kindNonBlock means that the descriptor is already in</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// non-blocking mode.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	kindNonBlock
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// kindNoPoll means that we should not put the descriptor into</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// non-blocking mode, because we know it is not a pipe or FIFO.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// Used by openFdAt for directories.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	kindNoPoll
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// newFile is like NewFile, but if called from OpenFile or Pipe</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// (as passed in the kind parameter) it tries to add the file to</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// the runtime poller.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func newFile(fd int, name string, kind newFileKind) *File {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	f := &amp;File{&amp;file{
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		pfd: poll.FD{
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			Sysfd:         fd,
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			IsStream:      true,
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			ZeroReadIsEOF: true,
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		},
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		name:        name,
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		stdoutOrErr: fd == 1 || fd == 2,
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	}}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	pollable := kind == kindOpenFile || kind == kindPipe || kind == kindNonBlock
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// If the caller passed a non-blocking filedes (kindNonBlock),</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// we assume they know what they are doing so we allow it to be</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// used with kqueue.</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if kind == kindOpenFile {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		switch runtime.GOOS {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		case &#34;darwin&#34;, &#34;ios&#34;, &#34;dragonfly&#34;, &#34;freebsd&#34;, &#34;netbsd&#34;, &#34;openbsd&#34;:
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			var st syscall.Stat_t
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			err := ignoringEINTR(func() error {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				return syscall.Fstat(fd, &amp;st)
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			})
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			typ := st.Mode &amp; syscall.S_IFMT
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t try to use kqueue with regular files on *BSDs.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			<span class="comment">// On FreeBSD a regular file is always</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			<span class="comment">// reported as ready for writing.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			<span class="comment">// On Dragonfly, NetBSD and OpenBSD the fd is signaled</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			<span class="comment">// only once as ready (both read and write).</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			<span class="comment">// Issue 19093.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			<span class="comment">// Also don&#39;t add directories to the netpoller.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			if err == nil &amp;&amp; (typ == syscall.S_IFREG || typ == syscall.S_IFDIR) {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				pollable = false
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			<span class="comment">// In addition to the behavior described above for regular files,</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			<span class="comment">// on Darwin, kqueue does not work properly with fifos:</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			<span class="comment">// closing the last writer does not cause a kqueue event</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			<span class="comment">// for any readers. See issue #24164.</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			if (runtime.GOOS == &#34;darwin&#34; || runtime.GOOS == &#34;ios&#34;) &amp;&amp; typ == syscall.S_IFIFO {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				pollable = false
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	clearNonBlock := false
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if pollable {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		if kind == kindNonBlock {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			<span class="comment">// The descriptor is already in non-blocking mode.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// We only set f.nonblock if we put the file into</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// non-blocking mode.</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		} else if err := syscall.SetNonblock(fd, true); err == nil {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			f.nonblock = true
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			clearNonBlock = true
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		} else {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			pollable = false
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// An error here indicates a failure to register</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// with the netpoll system. That can happen for</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// a file descriptor that is not supported by</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// epoll/kqueue; for example, disk files on</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// Linux systems. We assume that any real error</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// will show up in later I/O.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// We do restore the blocking behavior if it was set by us.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	if pollErr := f.pfd.Init(&#34;file&#34;, pollable); pollErr != nil &amp;&amp; clearNonBlock {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		if err := syscall.SetNonblock(fd, false); err == nil {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			f.nonblock = false
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	runtime.SetFinalizer(f.file, (*file).close)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	return f
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>func sigpipe() <span class="comment">// implemented in package runtime</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// epipecheck raises SIGPIPE if we get an EPIPE error on standard</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// output or standard error. See the SIGPIPE docs in os/signal, and</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// issue 11845.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>func epipecheck(file *File, e error) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	if e == syscall.EPIPE &amp;&amp; file.stdoutOrErr {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		sigpipe()
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// DevNull is the name of the operating system&#39;s “null device.”</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// On Unix-like systems, it is &#34;/dev/null&#34;; on Windows, &#34;NUL&#34;.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>const DevNull = &#34;/dev/null&#34;
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// openFileNolog is the Unix implementation of OpenFile.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// Changes here should be reflected in openFdAt, if relevant.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>func openFileNolog(name string, flag int, perm FileMode) (*File, error) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	setSticky := false
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if !supportsCreateWithStickyBit &amp;&amp; flag&amp;O_CREATE != 0 &amp;&amp; perm&amp;ModeSticky != 0 {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		if _, err := Stat(name); IsNotExist(err) {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			setSticky = true
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	var r int
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	var s poll.SysFile
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	for {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		var e error
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		r, s, e = open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		if e == nil {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			break
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">// We have to check EINTR here, per issues 11180 and 39237.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		if e == syscall.EINTR {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			continue
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return nil, &amp;PathError{Op: &#34;open&#34;, Path: name, Err: e}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	<span class="comment">// open(2) itself won&#39;t handle the sticky bit on *BSD and Solaris</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	if setSticky {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		setStickyBit(name)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// There&#39;s a race here with fork/exec, which we are</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// content to live with. See ../syscall/exec_unix.go.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	if !supportsCloseOnExec {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		syscall.CloseOnExec(r)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	kind := kindOpenFile
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if unix.HasNonblockFlag(flag) {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		kind = kindNonBlock
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	f := newFile(r, name, kind)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	f.pfd.SysFile = s
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	return f, nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (file *file) close() error {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if file == nil {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		return syscall.EINVAL
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	if file.dirinfo != nil {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		file.dirinfo.close()
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		file.dirinfo = nil
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	var err error
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	if e := file.pfd.Close(); e != nil {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if e == poll.ErrFileClosing {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			e = ErrClosed
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		err = &amp;PathError{Op: &#34;close&#34;, Path: file.name, Err: e}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// no need for a finalizer anymore</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	runtime.SetFinalizer(file, nil)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	return err
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// seek sets the offset for the next Read or Write on file to offset, interpreted</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// according to whence: 0 means relative to the origin of the file, 1 means</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// relative to the current offset, and 2 means relative to the end.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span><span class="comment">// It returns the new offset and an error, if any.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>func (f *File) seek(offset int64, whence int) (ret int64, err error) {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	if f.dirinfo != nil {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		<span class="comment">// Free cached dirinfo, so we allocate a new one if we</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		<span class="comment">// access this file as a directory again. See #35767 and #37161.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		f.dirinfo.close()
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		f.dirinfo = nil
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	ret, err = f.pfd.Seek(offset, whence)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	runtime.KeepAlive(f)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	return ret, err
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// Truncate changes the size of the named file.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// If the file is a symbolic link, it changes the size of the link&#39;s target.</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// If there is an error, it will be of type *PathError.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func Truncate(name string, size int64) error {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	e := ignoringEINTR(func() error {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		return syscall.Truncate(name, size)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	})
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	if e != nil {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		return &amp;PathError{Op: &#34;truncate&#34;, Path: name, Err: e}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	return nil
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// Remove removes the named file or (empty) directory.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// If there is an error, it will be of type *PathError.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func Remove(name string) error {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	<span class="comment">// System call interface forces us to know</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// whether name is a file or directory.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	<span class="comment">// Try both: it is cheaper on average than</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// doing a Stat plus the right one.</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	e := ignoringEINTR(func() error {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		return syscall.Unlink(name)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	})
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	if e == nil {
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		return nil
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	e1 := ignoringEINTR(func() error {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		return syscall.Rmdir(name)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	})
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	if e1 == nil {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		return nil
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// Both failed: figure out which error to return.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">// OS X and Linux differ on whether unlink(dir)</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// returns EISDIR, so can&#39;t use that. However,</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// both agree that rmdir(file) returns ENOTDIR,</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// so we can use that to decide which error is real.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	<span class="comment">// Rmdir might also return ENOTDIR if given a bad</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	<span class="comment">// file path, like /etc/passwd/foo, but in that case,</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// both errors will be ENOTDIR, so it&#39;s okay to</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// use the error from unlink.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	if e1 != syscall.ENOTDIR {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		e = e1
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	return &amp;PathError{Op: &#34;remove&#34;, Path: name, Err: e}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>func tempDir() string {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	dir := Getenv(&#34;TMPDIR&#34;)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	if dir == &#34;&#34; {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		if runtime.GOOS == &#34;android&#34; {
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			dir = &#34;/data/local/tmp&#34;
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		} else {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			dir = &#34;/tmp&#34;
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	return dir
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">// Link creates newname as a hard link to the oldname file.</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// If there is an error, it will be of type *LinkError.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>func Link(oldname, newname string) error {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	e := ignoringEINTR(func() error {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		return syscall.Link(oldname, newname)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	})
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	if e != nil {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		return &amp;LinkError{&#34;link&#34;, oldname, newname, e}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	return nil
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// Symlink creates newname as a symbolic link to oldname.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">// On Windows, a symlink to a non-existent oldname creates a file symlink;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">// if oldname is later created as a directory the symlink will not work.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">// If there is an error, it will be of type *LinkError.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>func Symlink(oldname, newname string) error {
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	e := ignoringEINTR(func() error {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		return syscall.Symlink(oldname, newname)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	})
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	if e != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		return &amp;LinkError{&#34;symlink&#34;, oldname, newname, e}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	return nil
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func readlink(name string) (string, error) {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	for len := 128; ; len *= 2 {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		b := make([]byte, len)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		var (
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			n int
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			e error
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		)
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		for {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			n, e = fixCount(syscall.Readlink(name, b))
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			if e != syscall.EINTR {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				break
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		<span class="comment">// buffer too small</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		if (runtime.GOOS == &#34;aix&#34; || runtime.GOOS == &#34;wasip1&#34;) &amp;&amp; e == syscall.ERANGE {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			continue
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		if e != nil {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			return &#34;&#34;, &amp;PathError{Op: &#34;readlink&#34;, Path: name, Err: e}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		if n &lt; len {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			return string(b[0:n]), nil
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>type unixDirent struct {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	parent string
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	name   string
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	typ    FileMode
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	info   FileInfo
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>func (d *unixDirent) Name() string   { return d.name }
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>func (d *unixDirent) IsDir() bool    { return d.typ.IsDir() }
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>func (d *unixDirent) Type() FileMode { return d.typ }
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>func (d *unixDirent) Info() (FileInfo, error) {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if d.info != nil {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		return d.info, nil
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	return lstat(d.parent + &#34;/&#34; + d.name)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>func (d *unixDirent) String() string {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	return fs.FormatDirEntry(d)
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>func newUnixDirent(parent, name string, typ FileMode) (DirEntry, error) {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	ude := &amp;unixDirent{
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		parent: parent,
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		name:   name,
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		typ:    typ,
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	if typ != ^FileMode(0) &amp;&amp; !testingForceReadDirLstat {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		return ude, nil
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	info, err := lstat(parent + &#34;/&#34; + name)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	if err != nil {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		return nil, err
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	ude.typ = info.Mode().Type()
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	ude.info = info
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	return ude, nil
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>
</pre><p><a href="file_unix.go?m=text">View as plain text</a></p>

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
