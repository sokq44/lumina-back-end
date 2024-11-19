<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/os/user/cgo_lookup_cgo.go - Go Documentation Server</title>

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
<a href="cgo_lookup_cgo.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/os">os</a>/<a href="http://localhost:8080/src/os/user">user</a>/<span class="text-muted">cgo_lookup_cgo.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/os/user">os/user</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build cgo &amp;&amp; !osusergo &amp;&amp; unix &amp;&amp; !android &amp;&amp; !darwin</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package user
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;syscall&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">/*
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>#cgo solaris CFLAGS: -D_POSIX_PTHREAD_SEMANTICS
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>#cgo CFLAGS: -fno-stack-protector
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>#include &lt;unistd.h&gt;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>#include &lt;sys/types.h&gt;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>#include &lt;pwd.h&gt;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>#include &lt;grp.h&gt;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>#include &lt;stdlib.h&gt;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>#include &lt;string.h&gt;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>static struct passwd mygetpwuid_r(int uid, char *buf, size_t buflen, int *found, int *perr) {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	struct passwd pwd;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	struct passwd *result;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	memset (&amp;pwd, 0, sizeof(pwd));
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	*perr = getpwuid_r(uid, &amp;pwd, buf, buflen, &amp;result);
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	*found = result != NULL;
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	return pwd;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>static struct passwd mygetpwnam_r(const char *name, char *buf, size_t buflen, int *found, int *perr) {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	struct passwd pwd;
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	struct passwd *result;
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	memset(&amp;pwd, 0, sizeof(pwd));
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	*perr = getpwnam_r(name, &amp;pwd, buf, buflen, &amp;result);
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	*found = result != NULL;
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	return pwd;
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>static struct group mygetgrgid_r(int gid, char *buf, size_t buflen, int *found, int *perr) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	struct group grp;
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	struct group *result;
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	memset(&amp;grp, 0, sizeof(grp));
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	*perr = getgrgid_r(gid, &amp;grp, buf, buflen, &amp;result);
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	*found = result != NULL;
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	return grp;
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>static struct group mygetgrnam_r(const char *name, char *buf, size_t buflen, int *found, int *perr) {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	struct group grp;
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	struct group *result;
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	memset(&amp;grp, 0, sizeof(grp));
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	*perr = getgrnam_r(name, &amp;grp, buf, buflen, &amp;result);
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	*found = result != NULL;
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	return grp;
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>*/</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>import &#34;C&#34;
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>type _C_char = C.char
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>type _C_int = C.int
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>type _C_gid_t = C.gid_t
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>type _C_uid_t = C.uid_t
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>type _C_size_t = C.size_t
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>type _C_struct_group = C.struct_group
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>type _C_struct_passwd = C.struct_passwd
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>type _C_long = C.long
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func _C_pw_uid(p *_C_struct_passwd) _C_uid_t   { return p.pw_uid }
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func _C_pw_uidp(p *_C_struct_passwd) *_C_uid_t { return &amp;p.pw_uid }
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func _C_pw_gid(p *_C_struct_passwd) _C_gid_t   { return p.pw_gid }
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>func _C_pw_gidp(p *_C_struct_passwd) *_C_gid_t { return &amp;p.pw_gid }
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func _C_pw_name(p *_C_struct_passwd) *_C_char  { return p.pw_name }
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func _C_pw_gecos(p *_C_struct_passwd) *_C_char { return p.pw_gecos }
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func _C_pw_dir(p *_C_struct_passwd) *_C_char   { return p.pw_dir }
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func _C_gr_gid(g *_C_struct_group) _C_gid_t  { return g.gr_gid }
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func _C_gr_name(g *_C_struct_group) *_C_char { return g.gr_name }
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>func _C_GoString(p *_C_char) string { return C.GoString(p) }
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func _C_getpwnam_r(name *_C_char, buf *_C_char, size _C_size_t) (pwd _C_struct_passwd, found bool, errno syscall.Errno) {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	var f, e _C_int
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	pwd = C.mygetpwnam_r(name, buf, size, &amp;f, &amp;e)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	return pwd, f != 0, syscall.Errno(e)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>func _C_getpwuid_r(uid _C_uid_t, buf *_C_char, size _C_size_t) (pwd _C_struct_passwd, found bool, errno syscall.Errno) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	var f, e _C_int
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	pwd = C.mygetpwuid_r(_C_int(uid), buf, size, &amp;f, &amp;e)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	return pwd, f != 0, syscall.Errno(e)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func _C_getgrnam_r(name *_C_char, buf *_C_char, size _C_size_t) (grp _C_struct_group, found bool, errno syscall.Errno) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	var f, e _C_int
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	grp = C.mygetgrnam_r(name, buf, size, &amp;f, &amp;e)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	return grp, f != 0, syscall.Errno(e)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func _C_getgrgid_r(gid _C_gid_t, buf *_C_char, size _C_size_t) (grp _C_struct_group, found bool, errno syscall.Errno) {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	var f, e _C_int
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	grp = C.mygetgrgid_r(_C_int(gid), buf, size, &amp;f, &amp;e)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	return grp, f != 0, syscall.Errno(e)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>const (
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	_C__SC_GETPW_R_SIZE_MAX = C._SC_GETPW_R_SIZE_MAX
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	_C__SC_GETGR_R_SIZE_MAX = C._SC_GETGR_R_SIZE_MAX
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func _C_sysconf(key _C_int) _C_long { return C.sysconf(key) }
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
</pre><p><a href="cgo_lookup_cgo.go?m=text">View as plain text</a></p>

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
