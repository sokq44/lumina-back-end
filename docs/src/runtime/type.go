<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/type.go - Go Documentation Server</title>

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
<a href="type.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">type.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Runtime type representation.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>type nameOff = abi.NameOff
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>type typeOff = abi.TypeOff
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>type textOff = abi.TextOff
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type _type = abi.Type
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// rtype is a wrapper that allows us to define additional methods.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type rtype struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	*abi.Type <span class="comment">// embedding is okay here (unlike reflect) because none of this is public</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func (t rtype) string() string {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	s := t.nameOff(t.Str).Name()
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	if t.TFlag&amp;abi.TFlagExtraStar != 0 {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		return s[1:]
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	return s
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func (t rtype) uncommon() *uncommontype {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	return t.Uncommon()
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>func (t rtype) name() string {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	if t.TFlag&amp;abi.TFlagNamed == 0 {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	s := t.string()
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	i := len(s) - 1
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	sqBrackets := 0
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	for i &gt;= 0 &amp;&amp; (s[i] != &#39;.&#39; || sqBrackets != 0) {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		switch s[i] {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		case &#39;]&#39;:
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			sqBrackets++
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		case &#39;[&#39;:
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			sqBrackets--
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		i--
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	return s[i+1:]
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// pkgpath returns the path of the package where t was defined, if</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// available. This is not the same as the reflect package&#39;s PkgPath</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// method, in that it returns the package path for struct and interface</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// types, not just named types.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func (t rtype) pkgpath() string {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if u := t.uncommon(); u != nil {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return t.nameOff(u.PkgPath).Name()
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	switch t.Kind_ &amp; kindMask {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	case kindStruct:
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		st := (*structtype)(unsafe.Pointer(t.Type))
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return st.PkgPath.Name()
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	case kindInterface:
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		it := (*interfacetype)(unsafe.Pointer(t.Type))
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		return it.PkgPath.Name()
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	return &#34;&#34;
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// reflectOffs holds type offsets defined at run time by the reflect package.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// When a type is defined at run time, its *rtype data lives on the heap.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// There are a wide range of possible addresses the heap may use, that</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// may not be representable as a 32-bit offset. Moreover the GC may</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// one day start moving heap memory, in which case there is no stable</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// offset that can be defined.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// To provide stable offsets, we add pin *rtype objects in a global map</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// and treat the offset as an identifier. We use negative offsets that</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// do not overlap with any compile-time module offsets.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// Entries are created by reflect.addReflectOff.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>var reflectOffs struct {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	lock mutex
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	next int32
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	m    map[int32]unsafe.Pointer
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	minv map[unsafe.Pointer]int32
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func reflectOffsLock() {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	lock(&amp;reflectOffs.lock)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if raceenabled {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;reflectOffs.lock))
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func reflectOffsUnlock() {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if raceenabled {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		racerelease(unsafe.Pointer(&amp;reflectOffs.lock))
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	unlock(&amp;reflectOffs.lock)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func resolveNameOff(ptrInModule unsafe.Pointer, off nameOff) name {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if off == 0 {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return name{}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	base := uintptr(ptrInModule)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	for md := &amp;firstmoduledata; md != nil; md = md.next {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		if base &gt;= md.types &amp;&amp; base &lt; md.etypes {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			res := md.types + uintptr(off)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			if res &gt; md.etypes {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>				println(&#34;runtime: nameOff&#34;, hex(off), &#34;out of range&#34;, hex(md.types), &#34;-&#34;, hex(md.etypes))
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				throw(&#34;runtime: name offset out of range&#34;)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			return name{Bytes: (*byte)(unsafe.Pointer(res))}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// No module found. see if it is a run time name.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	reflectOffsLock()
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	res, found := reflectOffs.m[int32(off)]
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	reflectOffsUnlock()
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if !found {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		println(&#34;runtime: nameOff&#34;, hex(off), &#34;base&#34;, hex(base), &#34;not in ranges:&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		for next := &amp;firstmoduledata; next != nil; next = next.next {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			println(&#34;\ttypes&#34;, hex(next.types), &#34;etypes&#34;, hex(next.etypes))
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		throw(&#34;runtime: name offset base pointer out of range&#34;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	return name{Bytes: (*byte)(res)}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func (t rtype) nameOff(off nameOff) name {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	return resolveNameOff(unsafe.Pointer(t.Type), off)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func resolveTypeOff(ptrInModule unsafe.Pointer, off typeOff) *_type {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if off == 0 || off == -1 {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// -1 is the sentinel value for unreachable code.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		<span class="comment">// See cmd/link/internal/ld/data.go:relocsym.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		return nil
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	base := uintptr(ptrInModule)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	var md *moduledata
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	for next := &amp;firstmoduledata; next != nil; next = next.next {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		if base &gt;= next.types &amp;&amp; base &lt; next.etypes {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			md = next
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			break
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if md == nil {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		reflectOffsLock()
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		res := reflectOffs.m[int32(off)]
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		reflectOffsUnlock()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		if res == nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			println(&#34;runtime: typeOff&#34;, hex(off), &#34;base&#34;, hex(base), &#34;not in ranges:&#34;)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			for next := &amp;firstmoduledata; next != nil; next = next.next {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				println(&#34;\ttypes&#34;, hex(next.types), &#34;etypes&#34;, hex(next.etypes))
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			throw(&#34;runtime: type offset base pointer out of range&#34;)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return (*_type)(res)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if t := md.typemap[off]; t != nil {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return t
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	res := md.types + uintptr(off)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if res &gt; md.etypes {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		println(&#34;runtime: typeOff&#34;, hex(off), &#34;out of range&#34;, hex(md.types), &#34;-&#34;, hex(md.etypes))
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		throw(&#34;runtime: type offset out of range&#34;)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	return (*_type)(unsafe.Pointer(res))
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func (t rtype) typeOff(off typeOff) *_type {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	return resolveTypeOff(unsafe.Pointer(t.Type), off)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>func (t rtype) textOff(off textOff) unsafe.Pointer {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if off == -1 {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// -1 is the sentinel value for unreachable code.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// See cmd/link/internal/ld/data.go:relocsym.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		return unsafe.Pointer(abi.FuncPCABIInternal(unreachableMethod))
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	base := uintptr(unsafe.Pointer(t.Type))
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	var md *moduledata
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	for next := &amp;firstmoduledata; next != nil; next = next.next {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		if base &gt;= next.types &amp;&amp; base &lt; next.etypes {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			md = next
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			break
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	if md == nil {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		reflectOffsLock()
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		res := reflectOffs.m[int32(off)]
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		reflectOffsUnlock()
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		if res == nil {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			println(&#34;runtime: textOff&#34;, hex(off), &#34;base&#34;, hex(base), &#34;not in ranges:&#34;)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			for next := &amp;firstmoduledata; next != nil; next = next.next {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				println(&#34;\ttypes&#34;, hex(next.types), &#34;etypes&#34;, hex(next.etypes))
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			throw(&#34;runtime: text offset base pointer out of range&#34;)
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		return res
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	res := md.textAddr(uint32(off))
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	return unsafe.Pointer(res)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>type uncommontype = abi.UncommonType
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>type interfacetype = abi.InterfaceType
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>type maptype = abi.MapType
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>type arraytype = abi.ArrayType
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>type chantype = abi.ChanType
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>type slicetype = abi.SliceType
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>type functype = abi.FuncType
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>type ptrtype = abi.PtrType
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>type name = abi.Name
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>type structtype = abi.StructType
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>func pkgPath(n name) string {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	if n.Bytes == nil || *n.Data(0)&amp;(1&lt;&lt;2) == 0 {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	i, l := n.ReadVarint(1)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	off := 1 + i + l
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if *n.Data(0)&amp;(1&lt;&lt;1) != 0 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		i2, l2 := n.ReadVarint(off)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		off += i2 + l2
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	var nameOff nameOff
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	copy((*[4]byte)(unsafe.Pointer(&amp;nameOff))[:], (*[4]byte)(unsafe.Pointer(n.Data(off)))[:])
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	pkgPathName := resolveNameOff(unsafe.Pointer(n.Bytes), nameOff)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	return pkgPathName.Name()
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">// typelinksinit scans the types from extra modules and builds the</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span><span class="comment">// moduledata typemap used to de-duplicate type pointers.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>func typelinksinit() {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	if firstmoduledata.next == nil {
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		return
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	typehash := make(map[uint32][]*_type, len(firstmoduledata.typelinks))
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	modules := activeModules()
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	prev := modules[0]
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	for _, md := range modules[1:] {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		<span class="comment">// Collect types from the previous module into typehash.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	collect:
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		for _, tl := range prev.typelinks {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			var t *_type
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			if prev.typemap == nil {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				t = (*_type)(unsafe.Pointer(prev.types + uintptr(tl)))
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			} else {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>				t = prev.typemap[typeOff(tl)]
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			<span class="comment">// Add to typehash if not seen before.</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			tlist := typehash[t.Hash]
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			for _, tcur := range tlist {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>				if tcur == t {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>					continue collect
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>				}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			typehash[t.Hash] = append(tlist, t)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		}
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		if md.typemap == nil {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			<span class="comment">// If any of this module&#39;s typelinks match a type from a</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			<span class="comment">// prior module, prefer that prior type by adding the offset</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			<span class="comment">// to this module&#39;s typemap.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			tm := make(map[typeOff]*_type, len(md.typelinks))
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			pinnedTypemaps = append(pinnedTypemaps, tm)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			md.typemap = tm
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			for _, tl := range md.typelinks {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>				t := (*_type)(unsafe.Pointer(md.types + uintptr(tl)))
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				for _, candidate := range typehash[t.Hash] {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>					seen := map[_typePair]struct{}{}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>					if typesEqual(t, candidate, seen) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>						t = candidate
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>						break
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>					}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>				}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				md.typemap[typeOff(tl)] = t
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		prev = md
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>type _typePair struct {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	t1 *_type
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	t2 *_type
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>func toRType(t *abi.Type) rtype {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	return rtype{t}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// typesEqual reports whether two types are equal.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// Everywhere in the runtime and reflect packages, it is assumed that</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// there is exactly one *_type per Go type, so that pointer equality</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// can be used to test if types are equal. There is one place that</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// breaks this assumption: buildmode=shared. In this case a type can</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span><span class="comment">// appear as two different pieces of memory. This is hidden from the</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span><span class="comment">// runtime and reflect package by the per-module typemap built in</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// typelinksinit. It uses typesEqual to map types from later modules</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// back into earlier ones.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span><span class="comment">// Only typelinksinit needs this function.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>func typesEqual(t, v *_type, seen map[_typePair]struct{}) bool {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	tp := _typePair{t, v}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	if _, ok := seen[tp]; ok {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return true
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// mark these types as seen, and thus equivalent which prevents an infinite loop if</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// the two types are identical, but recursively defined and loaded from</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// different modules</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	seen[tp] = struct{}{}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	if t == v {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		return true
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	kind := t.Kind_ &amp; kindMask
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	if kind != v.Kind_&amp;kindMask {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		return false
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	rt, rv := toRType(t), toRType(v)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	if rt.string() != rv.string() {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		return false
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	ut := t.Uncommon()
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	uv := v.Uncommon()
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	if ut != nil || uv != nil {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		if ut == nil || uv == nil {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			return false
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		pkgpatht := rt.nameOff(ut.PkgPath).Name()
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		pkgpathv := rv.nameOff(uv.PkgPath).Name()
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		if pkgpatht != pkgpathv {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			return false
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	if kindBool &lt;= kind &amp;&amp; kind &lt;= kindComplex128 {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		return true
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	switch kind {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	case kindString, kindUnsafePointer:
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		return true
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	case kindArray:
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		at := (*arraytype)(unsafe.Pointer(t))
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		av := (*arraytype)(unsafe.Pointer(v))
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		return typesEqual(at.Elem, av.Elem, seen) &amp;&amp; at.Len == av.Len
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	case kindChan:
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		ct := (*chantype)(unsafe.Pointer(t))
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		cv := (*chantype)(unsafe.Pointer(v))
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		return ct.Dir == cv.Dir &amp;&amp; typesEqual(ct.Elem, cv.Elem, seen)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	case kindFunc:
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		ft := (*functype)(unsafe.Pointer(t))
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		fv := (*functype)(unsafe.Pointer(v))
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		if ft.OutCount != fv.OutCount || ft.InCount != fv.InCount {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			return false
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		tin, vin := ft.InSlice(), fv.InSlice()
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		for i := 0; i &lt; len(tin); i++ {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			if !typesEqual(tin[i], vin[i], seen) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				return false
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>			}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		tout, vout := ft.OutSlice(), fv.OutSlice()
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		for i := 0; i &lt; len(tout); i++ {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			if !typesEqual(tout[i], vout[i], seen) {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>				return false
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		return true
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	case kindInterface:
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		it := (*interfacetype)(unsafe.Pointer(t))
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		iv := (*interfacetype)(unsafe.Pointer(v))
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		if it.PkgPath.Name() != iv.PkgPath.Name() {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>			return false
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		if len(it.Methods) != len(iv.Methods) {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			return false
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		for i := range it.Methods {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			tm := &amp;it.Methods[i]
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			vm := &amp;iv.Methods[i]
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// Note the mhdr array can be relocated from</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// another module. See #17724.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			tname := resolveNameOff(unsafe.Pointer(tm), tm.Name)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			vname := resolveNameOff(unsafe.Pointer(vm), vm.Name)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			if tname.Name() != vname.Name() {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>				return false
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			if pkgPath(tname) != pkgPath(vname) {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				return false
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			tityp := resolveTypeOff(unsafe.Pointer(tm), tm.Typ)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			vityp := resolveTypeOff(unsafe.Pointer(vm), vm.Typ)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			if !typesEqual(tityp, vityp, seen) {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>				return false
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		}
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		return true
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	case kindMap:
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		mt := (*maptype)(unsafe.Pointer(t))
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		mv := (*maptype)(unsafe.Pointer(v))
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		return typesEqual(mt.Key, mv.Key, seen) &amp;&amp; typesEqual(mt.Elem, mv.Elem, seen)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	case kindPtr:
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		pt := (*ptrtype)(unsafe.Pointer(t))
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		pv := (*ptrtype)(unsafe.Pointer(v))
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		return typesEqual(pt.Elem, pv.Elem, seen)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	case kindSlice:
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		st := (*slicetype)(unsafe.Pointer(t))
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		sv := (*slicetype)(unsafe.Pointer(v))
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return typesEqual(st.Elem, sv.Elem, seen)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	case kindStruct:
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		st := (*structtype)(unsafe.Pointer(t))
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		sv := (*structtype)(unsafe.Pointer(v))
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		if len(st.Fields) != len(sv.Fields) {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			return false
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		if st.PkgPath.Name() != sv.PkgPath.Name() {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			return false
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		for i := range st.Fields {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			tf := &amp;st.Fields[i]
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			vf := &amp;sv.Fields[i]
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			if tf.Name.Name() != vf.Name.Name() {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>				return false
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			if !typesEqual(tf.Typ, vf.Typ, seen) {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>				return false
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			if tf.Name.Tag() != vf.Name.Tag() {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>				return false
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			if tf.Offset != vf.Offset {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>				return false
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			if tf.Name.IsEmbedded() != vf.Name.IsEmbedded() {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>				return false
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		return true
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	default:
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		println(&#34;runtime: impossible type kind&#34;, kind)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		throw(&#34;runtime: impossible type kind&#34;)
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		return false
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
</pre><p><a href="type.go?m=text">View as plain text</a></p>

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
