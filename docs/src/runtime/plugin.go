<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/plugin.go - Go Documentation Server</title>

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
<a href="plugin.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">plugin.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//go:linkname plugin_lastmoduleinit plugin.lastmoduleinit</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>func plugin_lastmoduleinit() (path string, syms map[string]any, initTasks []*initTask, errstr string) {
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	var md *moduledata
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	for pmd := firstmoduledata.next; pmd != nil; pmd = pmd.next {
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>		if pmd.bad {
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>			md = nil <span class="comment">// we only want the last module</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>			continue
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		}
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		md = pmd
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	}
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	if md == nil {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>		throw(&#34;runtime: no plugin module data&#34;)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	if md.pluginpath == &#34;&#34; {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		throw(&#34;runtime: plugin has empty pluginpath&#34;)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	if md.typemap != nil {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		return &#34;&#34;, nil, nil, &#34;plugin already loaded&#34;
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	}
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	for _, pmd := range activeModules() {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		if pmd.pluginpath == md.pluginpath {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			md.bad = true
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>			return &#34;&#34;, nil, nil, &#34;plugin already loaded&#34;
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		if inRange(pmd.text, pmd.etext, md.text, md.etext) ||
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>			inRange(pmd.bss, pmd.ebss, md.bss, md.ebss) ||
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>			inRange(pmd.data, pmd.edata, md.data, md.edata) ||
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			inRange(pmd.types, pmd.etypes, md.types, md.etypes) {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>			println(&#34;plugin: new module data overlaps with previous moduledata&#34;)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>			println(&#34;\tpmd.text-etext=&#34;, hex(pmd.text), &#34;-&#34;, hex(pmd.etext))
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>			println(&#34;\tpmd.bss-ebss=&#34;, hex(pmd.bss), &#34;-&#34;, hex(pmd.ebss))
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>			println(&#34;\tpmd.data-edata=&#34;, hex(pmd.data), &#34;-&#34;, hex(pmd.edata))
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			println(&#34;\tpmd.types-etypes=&#34;, hex(pmd.types), &#34;-&#34;, hex(pmd.etypes))
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			println(&#34;\tmd.text-etext=&#34;, hex(md.text), &#34;-&#34;, hex(md.etext))
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			println(&#34;\tmd.bss-ebss=&#34;, hex(md.bss), &#34;-&#34;, hex(md.ebss))
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			println(&#34;\tmd.data-edata=&#34;, hex(md.data), &#34;-&#34;, hex(md.edata))
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			println(&#34;\tmd.types-etypes=&#34;, hex(md.types), &#34;-&#34;, hex(md.etypes))
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			throw(&#34;plugin: new module data overlaps with previous moduledata&#34;)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	for _, pkghash := range md.pkghashes {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		if pkghash.linktimehash != *pkghash.runtimehash {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			md.bad = true
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			return &#34;&#34;, nil, nil, &#34;plugin was built with a different version of package &#34; + pkghash.modulename
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// Initialize the freshly loaded module.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	modulesinit()
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	typelinksinit()
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	pluginftabverify(md)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	moduledataverify1(md)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	lock(&amp;itabLock)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	for _, i := range md.itablinks {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		itabAdd(i)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	unlock(&amp;itabLock)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// Build a map of symbol names to symbols. Here in the runtime</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// we fill out the first word of the interface, the type. We</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// pass these zero value interfaces to the plugin package,</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// where the symbol value is filled in (usually via cgo).</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// Because functions are handled specially in the plugin package,</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// function symbol names are prefixed here with &#39;.&#39; to avoid</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// a dependency on the reflect package.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	syms = make(map[string]any, len(md.ptab))
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	for _, ptab := range md.ptab {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		symName := resolveNameOff(unsafe.Pointer(md.types), ptab.name)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		t := toRType((*_type)(unsafe.Pointer(md.types))).typeOff(ptab.typ) <span class="comment">// TODO can this stack of conversions be simpler?</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		var val any
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		valp := (*[2]unsafe.Pointer)(unsafe.Pointer(&amp;val))
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		(*valp)[0] = unsafe.Pointer(t)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		name := symName.Name()
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		if t.Kind_&amp;kindMask == kindFunc {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			name = &#34;.&#34; + name
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		syms[name] = val
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return md.pluginpath, syms, md.inittasks, &#34;&#34;
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func pluginftabverify(md *moduledata) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	badtable := false
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	for i := 0; i &lt; len(md.ftab); i++ {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		entry := md.textAddr(md.ftab[i].entryoff)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		if md.minpc &lt;= entry &amp;&amp; entry &lt;= md.maxpc {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			continue
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		f := funcInfo{(*_func)(unsafe.Pointer(&amp;md.pclntable[md.ftab[i].funcoff])), md}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		name := funcname(f)
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		<span class="comment">// A common bug is f.entry has a relocation to a duplicate</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		<span class="comment">// function symbol, meaning if we search for its PC we get</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		<span class="comment">// a valid entry with a name that is useful for debugging.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		name2 := &#34;none&#34;
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		entry2 := uintptr(0)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		f2 := findfunc(entry)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		if f2.valid() {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			name2 = funcname(f2)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			entry2 = f2.entry()
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		badtable = true
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		println(&#34;ftab entry&#34;, hex(entry), &#34;/&#34;, hex(entry2), &#34;: &#34;,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			name, &#34;/&#34;, name2, &#34;outside pc range:[&#34;, hex(md.minpc), &#34;,&#34;, hex(md.maxpc), &#34;], modulename=&#34;, md.modulename, &#34;, pluginpath=&#34;, md.pluginpath)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if badtable {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		throw(&#34;runtime: plugin has bad symbol table&#34;)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// inRange reports whether v0 or v1 are in the range [r0, r1].</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>func inRange(r0, r1, v0, v1 uintptr) bool {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	return (v0 &gt;= r0 &amp;&amp; v0 &lt;= r1) || (v1 &gt;= r0 &amp;&amp; v1 &lt;= r1)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// A ptabEntry is generated by the compiler for each exported function</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// and global variable in the main package of a plugin. It is used to</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// initialize the plugin module&#39;s symbol map.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>type ptabEntry struct {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	name nameOff
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	typ  typeOff
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
</pre><p><a href="plugin.go?m=text">View as plain text</a></p>

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
