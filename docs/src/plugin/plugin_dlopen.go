<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/plugin/plugin_dlopen.go - Go Documentation Server</title>

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
<a href="plugin_dlopen.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/plugin">plugin</a>/<span class="text-muted">plugin_dlopen.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/plugin">plugin</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build (linux &amp;&amp; cgo) || (darwin &amp;&amp; cgo) || (freebsd &amp;&amp; cgo)</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package plugin
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">/*
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>#cgo linux LDFLAGS: -ldl
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>#include &lt;dlfcn.h&gt;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>#include &lt;limits.h&gt;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>#include &lt;stdlib.h&gt;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>#include &lt;stdint.h&gt;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>#include &lt;stdio.h&gt;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>static uintptr_t pluginOpen(const char* path, char** err) {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	if (h == NULL) {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		*err = (char*)dlerror();
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	return (uintptr_t)h;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>static void* pluginLookup(uintptr_t h, const char* name, char** err) {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	void* r = dlsym((void*)h, name);
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	if (r == NULL) {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		*err = (char*)dlerror();
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	return r;
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>*/</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>import &#34;C&#34;
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>import (
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func open(name string) (*Plugin, error) {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	cPath := make([]byte, C.PATH_MAX+1)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	cRelName := make([]byte, len(name)+1)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	copy(cRelName, name)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	if C.realpath(
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		(*C.char)(unsafe.Pointer(&amp;cRelName[0])),
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		(*C.char)(unsafe.Pointer(&amp;cPath[0]))) == nil {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		return nil, errors.New(`plugin.Open(&#34;` + name + `&#34;): realpath failed`)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	filepath := C.GoString((*C.char)(unsafe.Pointer(&amp;cPath[0])))
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	pluginsMu.Lock()
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if p := plugins[filepath]; p != nil {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		pluginsMu.Unlock()
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		if p.err != &#34;&#34; {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			return nil, errors.New(`plugin.Open(&#34;` + name + `&#34;): ` + p.err + ` (previous failure)`)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		&lt;-p.loaded
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return p, nil
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	var cErr *C.char
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	h := C.pluginOpen((*C.char)(unsafe.Pointer(&amp;cPath[0])), &amp;cErr)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if h == 0 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		pluginsMu.Unlock()
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return nil, errors.New(`plugin.Open(&#34;` + name + `&#34;): ` + C.GoString(cErr))
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// TODO(crawshaw): look for plugin note, confirm it is a Go plugin</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// and it was built with the correct toolchain.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if len(name) &gt; 3 &amp;&amp; name[len(name)-3:] == &#34;.so&#34; {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		name = name[:len(name)-3]
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if plugins == nil {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		plugins = make(map[string]*Plugin)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	pluginpath, syms, initTasks, errstr := lastmoduleinit()
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	if errstr != &#34;&#34; {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		plugins[filepath] = &amp;Plugin{
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			pluginpath: pluginpath,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			err:        errstr,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		pluginsMu.Unlock()
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		return nil, errors.New(`plugin.Open(&#34;` + name + `&#34;): ` + errstr)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// This function can be called from the init function of a plugin.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// Drop a placeholder in the map so subsequent opens can wait on it.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	p := &amp;Plugin{
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		pluginpath: pluginpath,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		loaded:     make(chan struct{}),
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	plugins[filepath] = p
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	pluginsMu.Unlock()
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	doInit(initTasks)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// Fill out the value of each plugin symbol.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	updatedSyms := map[string]any{}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	for symName, sym := range syms {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		isFunc := symName[0] == &#39;.&#39;
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if isFunc {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			delete(syms, symName)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			symName = symName[1:]
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		fullName := pluginpath + &#34;.&#34; + symName
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		cname := make([]byte, len(fullName)+1)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		copy(cname, fullName)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		p := C.pluginLookup(h, (*C.char)(unsafe.Pointer(&amp;cname[0])), &amp;cErr)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		if p == nil {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			return nil, errors.New(`plugin.Open(&#34;` + name + `&#34;): could not find symbol ` + symName + `: ` + C.GoString(cErr))
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		valp := (*[2]unsafe.Pointer)(unsafe.Pointer(&amp;sym))
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if isFunc {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			(*valp)[1] = unsafe.Pointer(&amp;p)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		} else {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			(*valp)[1] = p
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		<span class="comment">// we can&#39;t add to syms during iteration as we&#39;ll end up processing</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		<span class="comment">// some symbols twice with the inability to tell if the symbol is a function</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		updatedSyms[symName] = sym
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	p.syms = updatedSyms
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	close(p.loaded)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	return p, nil
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>func lookup(p *Plugin, symName string) (Symbol, error) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if s := p.syms[symName]; s != nil {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		return s, nil
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	return nil, errors.New(&#34;plugin: symbol &#34; + symName + &#34; not found in plugin &#34; + p.pluginpath)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>var (
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	pluginsMu sync.Mutex
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	plugins   map[string]*Plugin
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// lastmoduleinit is defined in package runtime.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func lastmoduleinit() (pluginpath string, syms map[string]any, inittasks []*initTask, errstr string)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// doInit is defined in package runtime.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//go:linkname doInit runtime.doInit</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func doInit(t []*initTask)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>type initTask struct {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// fields defined in runtime.initTask. We only handle pointers to an initTask</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// in this package, so the contents are irrelevant.</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
</pre><p><a href="plugin_dlopen.go?m=text">View as plain text</a></p>

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
