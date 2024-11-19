<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/vdso_linux.go - Go Documentation Server</title>

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
<a href="vdso_linux.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">vdso_linux.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build linux &amp;&amp; (386 || amd64 || arm || arm64 || loong64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x)</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package runtime
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Look up symbols in the Linux vDSO.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// This code was originally based on the sample Linux vDSO parser at</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/tools/testing/selftests/vDSO/parse_vdso.c</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// This implements the ELF dynamic linking spec at</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// http://sco.com/developers/gabi/latest/ch5.dynamic.html</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The version section is documented at</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// https://refspecs.linuxfoundation.org/LSB_3.2.0/LSB-Core-generic/LSB-Core-generic/symversion.html</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const (
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	_AT_SYSINFO_EHDR = 33
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	_PT_LOAD    = 1 <span class="comment">/* Loadable program segment */</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	_PT_DYNAMIC = 2 <span class="comment">/* Dynamic linking information */</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	_DT_NULL     = 0          <span class="comment">/* Marks end of dynamic section */</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	_DT_HASH     = 4          <span class="comment">/* Dynamic symbol hash table */</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	_DT_STRTAB   = 5          <span class="comment">/* Address of string table */</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	_DT_SYMTAB   = 6          <span class="comment">/* Address of symbol table */</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	_DT_GNU_HASH = 0x6ffffef5 <span class="comment">/* GNU-style dynamic symbol hash table */</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	_DT_VERSYM   = 0x6ffffff0
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	_DT_VERDEF   = 0x6ffffffc
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	_VER_FLG_BASE = 0x1 <span class="comment">/* Version definition of file itself */</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	_SHN_UNDEF = 0 <span class="comment">/* Undefined section */</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	_SHT_DYNSYM = 11 <span class="comment">/* Dynamic linker symbol table */</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	_STT_FUNC = 2 <span class="comment">/* Symbol is a code object */</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	_STT_NOTYPE = 0 <span class="comment">/* Symbol type is not specified */</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	_STB_GLOBAL = 1 <span class="comment">/* Global symbol */</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	_STB_WEAK   = 2 <span class="comment">/* Weak symbol */</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	_EI_NIDENT = 16
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// Maximum indices for the array types used when traversing the vDSO ELF structures.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Computed from architecture-specific max provided by vdso_linux_*.go</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	vdsoSymTabSize     = vdsoArrayMax / unsafe.Sizeof(elfSym{})
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	vdsoDynSize        = vdsoArrayMax / unsafe.Sizeof(elfDyn{})
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	vdsoSymStringsSize = vdsoArrayMax     <span class="comment">// byte</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	vdsoVerSymSize     = vdsoArrayMax / 2 <span class="comment">// uint16</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	vdsoHashSize       = vdsoArrayMax / 4 <span class="comment">// uint32</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// vdsoBloomSizeScale is a scaling factor for gnuhash tables which are uint32 indexed,</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">// but contain uintptrs</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	vdsoBloomSizeScale = unsafe.Sizeof(uintptr(0)) / 4 <span class="comment">// uint32</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">/* How to extract and insert information held in the st_info field.  */</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func _ELF_ST_BIND(val byte) byte { return val &gt;&gt; 4 }
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func _ELF_ST_TYPE(val byte) byte { return val &amp; 0xf }
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>type vdsoSymbolKey struct {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	name    string
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	symHash uint32
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	gnuHash uint32
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	ptr     *uintptr
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>type vdsoVersionKey struct {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	version string
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	verHash uint32
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>type vdsoInfo struct {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	valid bool
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">/* Load information */</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	loadAddr   uintptr
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	loadOffset uintptr <span class="comment">/* loadAddr - recorded vaddr */</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">/* Symbol table */</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	symtab     *[vdsoSymTabSize]elfSym
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	symstrings *[vdsoSymStringsSize]byte
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	chain      []uint32
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	bucket     []uint32
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	symOff     uint32
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	isGNUHash  bool
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">/* Version table */</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	versym *[vdsoVerSymSize]uint16
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	verdef *elfVerdef
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// see vdso_linux_*.go for vdsoSymbolKeys[] and vdso*Sym vars</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func vdsoInitFromSysinfoEhdr(info *vdsoInfo, hdr *elfEhdr) {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	info.valid = false
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	info.loadAddr = uintptr(unsafe.Pointer(hdr))
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	pt := unsafe.Pointer(info.loadAddr + uintptr(hdr.e_phoff))
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	<span class="comment">// We need two things from the segment table: the load offset</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// and the dynamic table.</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	var foundVaddr bool
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	var dyn *[vdsoDynSize]elfDyn
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	for i := uint16(0); i &lt; hdr.e_phnum; i++ {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		pt := (*elfPhdr)(add(pt, uintptr(i)*unsafe.Sizeof(elfPhdr{})))
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		switch pt.p_type {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		case _PT_LOAD:
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			if !foundVaddr {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>				foundVaddr = true
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>				info.loadOffset = info.loadAddr + uintptr(pt.p_offset-pt.p_vaddr)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		case _PT_DYNAMIC:
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			dyn = (*[vdsoDynSize]elfDyn)(unsafe.Pointer(info.loadAddr + uintptr(pt.p_offset)))
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if !foundVaddr || dyn == nil {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return <span class="comment">// Failed</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">// Fish out the useful bits of the dynamic table.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	var hash, gnuhash *[vdsoHashSize]uint32
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	info.symstrings = nil
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	info.symtab = nil
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	info.versym = nil
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	info.verdef = nil
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	for i := 0; dyn[i].d_tag != _DT_NULL; i++ {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		dt := &amp;dyn[i]
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		p := info.loadOffset + uintptr(dt.d_val)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		switch dt.d_tag {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		case _DT_STRTAB:
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			info.symstrings = (*[vdsoSymStringsSize]byte)(unsafe.Pointer(p))
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		case _DT_SYMTAB:
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			info.symtab = (*[vdsoSymTabSize]elfSym)(unsafe.Pointer(p))
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		case _DT_HASH:
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			hash = (*[vdsoHashSize]uint32)(unsafe.Pointer(p))
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		case _DT_GNU_HASH:
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			gnuhash = (*[vdsoHashSize]uint32)(unsafe.Pointer(p))
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		case _DT_VERSYM:
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			info.versym = (*[vdsoVerSymSize]uint16)(unsafe.Pointer(p))
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		case _DT_VERDEF:
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			info.verdef = (*elfVerdef)(unsafe.Pointer(p))
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if info.symstrings == nil || info.symtab == nil || (hash == nil &amp;&amp; gnuhash == nil) {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return <span class="comment">// Failed</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if info.verdef == nil {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		info.versym = nil
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if gnuhash != nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// Parse the GNU hash table header.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		nbucket := gnuhash[0]
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		info.symOff = gnuhash[1]
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		bloomSize := gnuhash[2]
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		info.bucket = gnuhash[4+bloomSize*uint32(vdsoBloomSizeScale):][:nbucket]
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		info.chain = gnuhash[4+bloomSize*uint32(vdsoBloomSizeScale)+nbucket:]
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		info.isGNUHash = true
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	} else {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// Parse the hash table header.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		nbucket := hash[0]
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		nchain := hash[1]
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		info.bucket = hash[2 : 2+nbucket]
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		info.chain = hash[2+nbucket : 2+nbucket+nchain]
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// That&#39;s all we need.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	info.valid = true
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func vdsoFindVersion(info *vdsoInfo, ver *vdsoVersionKey) int32 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if !info.valid {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		return 0
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	def := info.verdef
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	for {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		if def.vd_flags&amp;_VER_FLG_BASE == 0 {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			aux := (*elfVerdaux)(add(unsafe.Pointer(def), uintptr(def.vd_aux)))
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			if def.vd_hash == ver.verHash &amp;&amp; ver.version == gostringnocopy(&amp;info.symstrings[aux.vda_name]) {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>				return int32(def.vd_ndx &amp; 0x7fff)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if def.vd_next == 0 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			break
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		def = (*elfVerdef)(add(unsafe.Pointer(def), uintptr(def.vd_next)))
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return -1 <span class="comment">// cannot match any version</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>func vdsoParseSymbols(info *vdsoInfo, version int32) {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if !info.valid {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		return
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	apply := func(symIndex uint32, k vdsoSymbolKey) bool {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		sym := &amp;info.symtab[symIndex]
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		typ := _ELF_ST_TYPE(sym.st_info)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		bind := _ELF_ST_BIND(sym.st_info)
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		<span class="comment">// On ppc64x, VDSO functions are of type _STT_NOTYPE.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if typ != _STT_FUNC &amp;&amp; typ != _STT_NOTYPE || bind != _STB_GLOBAL &amp;&amp; bind != _STB_WEAK || sym.st_shndx == _SHN_UNDEF {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			return false
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		if k.name != gostringnocopy(&amp;info.symstrings[sym.st_name]) {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			return false
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// Check symbol version.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if info.versym != nil &amp;&amp; version != 0 &amp;&amp; int32(info.versym[symIndex]&amp;0x7fff) != version {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			return false
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		*k.ptr = info.loadOffset + uintptr(sym.st_value)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		return true
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if !info.isGNUHash {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		<span class="comment">// Old-style DT_HASH table.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		for _, k := range vdsoSymbolKeys {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			if len(info.bucket) &gt; 0 {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>				for chain := info.bucket[k.symHash%uint32(len(info.bucket))]; chain != 0; chain = info.chain[chain] {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>					if apply(chain, k) {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>						break
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>					}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>				}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// New-style DT_GNU_HASH table.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	for _, k := range vdsoSymbolKeys {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		symIndex := info.bucket[k.gnuHash%uint32(len(info.bucket))]
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		if symIndex &lt; info.symOff {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			continue
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		for ; ; symIndex++ {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			hash := info.chain[symIndex-info.symOff]
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			if hash|1 == k.gnuHash|1 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				<span class="comment">// Found a hash match.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>				if apply(symIndex, k) {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>					break
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			if hash&amp;1 != 0 {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				<span class="comment">// End of chain.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				break
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>func vdsoauxv(tag, val uintptr) {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	switch tag {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	case _AT_SYSINFO_EHDR:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		if val == 0 {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			<span class="comment">// Something went wrong</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			return
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		var info vdsoInfo
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		<span class="comment">// TODO(rsc): I don&#39;t understand why the compiler thinks info escapes</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		<span class="comment">// when passed to the three functions below.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		info1 := (*vdsoInfo)(noescape(unsafe.Pointer(&amp;info)))
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		vdsoInitFromSysinfoEhdr(info1, (*elfEhdr)(unsafe.Pointer(val)))
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		vdsoParseSymbols(info1, vdsoFindVersion(info1, &amp;vdsoLinuxVersion))
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// vdsoMarker reports whether PC is on the VDSO page.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">//go:nosplit</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>func inVDSOPage(pc uintptr) bool {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	for _, k := range vdsoSymbolKeys {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		if *k.ptr != 0 {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			page := *k.ptr &amp;^ (physPageSize - 1)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			return pc &gt;= page &amp;&amp; pc &lt; page+physPageSize
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	return false
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
</pre><p><a href="vdso_linux.go?m=text">View as plain text</a></p>

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
