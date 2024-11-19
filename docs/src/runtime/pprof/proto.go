<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/pprof/proto.go - Go Documentation Server</title>

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
<a href="proto.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/pprof">pprof</a>/<span class="text-muted">proto.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/pprof">runtime/pprof</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package pprof
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;compress/gzip&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// lostProfileEvent is the function to which lost profiling</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// events are attributed.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// (The name shows up in the pprof graphs.)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>func lostProfileEvent() { lostProfileEvent() }
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// A profileBuilder writes a profile incrementally from a</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// stream of profile samples delivered by the runtime.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>type profileBuilder struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	start      time.Time
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	end        time.Time
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	havePeriod bool
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	period     int64
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	m          profMap
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// encoding state</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	w         io.Writer
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	zw        *gzip.Writer
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	pb        protobuf
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	strings   []string
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	stringMap map[string]int
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	locs      map[uintptr]locInfo <span class="comment">// list of locInfo starting with the given PC.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	funcs     map[string]int      <span class="comment">// Package path-qualified function name to Function.ID</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	mem       []memMap
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	deck      pcDeck
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>type memMap struct {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// initialized as reading mapping</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	start   uintptr <span class="comment">// Address at which the binary (or DLL) is loaded into memory.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	end     uintptr <span class="comment">// The limit of the address range occupied by this mapping.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	offset  uint64  <span class="comment">// Offset in the binary that corresponds to the first mapped address.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	file    string  <span class="comment">// The object this entry is loaded from.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	buildID string  <span class="comment">// A string that uniquely identifies a particular program version with high probability.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	funcs symbolizeFlag
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	fake  bool <span class="comment">// map entry was faked; /proc/self/maps wasn&#39;t available</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// symbolizeFlag keeps track of symbolization result.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//	0                  : no symbol lookup was performed</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//	1&lt;&lt;0 (lookupTried) : symbol lookup was performed</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//	1&lt;&lt;1 (lookupFailed): symbol lookup was performed but failed</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>type symbolizeFlag uint8
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>const (
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	lookupTried  symbolizeFlag = 1 &lt;&lt; iota
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	lookupFailed symbolizeFlag = 1 &lt;&lt; iota
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>const (
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// message Profile</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	tagProfile_SampleType        = 1  <span class="comment">// repeated ValueType</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	tagProfile_Sample            = 2  <span class="comment">// repeated Sample</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	tagProfile_Mapping           = 3  <span class="comment">// repeated Mapping</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	tagProfile_Location          = 4  <span class="comment">// repeated Location</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	tagProfile_Function          = 5  <span class="comment">// repeated Function</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	tagProfile_StringTable       = 6  <span class="comment">// repeated string</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	tagProfile_DropFrames        = 7  <span class="comment">// int64 (string table index)</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	tagProfile_KeepFrames        = 8  <span class="comment">// int64 (string table index)</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	tagProfile_TimeNanos         = 9  <span class="comment">// int64</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	tagProfile_DurationNanos     = 10 <span class="comment">// int64</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	tagProfile_PeriodType        = 11 <span class="comment">// ValueType (really optional string???)</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	tagProfile_Period            = 12 <span class="comment">// int64</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	tagProfile_Comment           = 13 <span class="comment">// repeated int64</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	tagProfile_DefaultSampleType = 14 <span class="comment">// int64</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// message ValueType</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	tagValueType_Type = 1 <span class="comment">// int64 (string table index)</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	tagValueType_Unit = 2 <span class="comment">// int64 (string table index)</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// message Sample</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	tagSample_Location = 1 <span class="comment">// repeated uint64</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	tagSample_Value    = 2 <span class="comment">// repeated int64</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	tagSample_Label    = 3 <span class="comment">// repeated Label</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// message Label</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	tagLabel_Key = 1 <span class="comment">// int64 (string table index)</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	tagLabel_Str = 2 <span class="comment">// int64 (string table index)</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	tagLabel_Num = 3 <span class="comment">// int64</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// message Mapping</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	tagMapping_ID              = 1  <span class="comment">// uint64</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	tagMapping_Start           = 2  <span class="comment">// uint64</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	tagMapping_Limit           = 3  <span class="comment">// uint64</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	tagMapping_Offset          = 4  <span class="comment">// uint64</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	tagMapping_Filename        = 5  <span class="comment">// int64 (string table index)</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	tagMapping_BuildID         = 6  <span class="comment">// int64 (string table index)</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	tagMapping_HasFunctions    = 7  <span class="comment">// bool</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	tagMapping_HasFilenames    = 8  <span class="comment">// bool</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	tagMapping_HasLineNumbers  = 9  <span class="comment">// bool</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	tagMapping_HasInlineFrames = 10 <span class="comment">// bool</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// message Location</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	tagLocation_ID        = 1 <span class="comment">// uint64</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	tagLocation_MappingID = 2 <span class="comment">// uint64</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	tagLocation_Address   = 3 <span class="comment">// uint64</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	tagLocation_Line      = 4 <span class="comment">// repeated Line</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// message Line</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	tagLine_FunctionID = 1 <span class="comment">// uint64</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	tagLine_Line       = 2 <span class="comment">// int64</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// message Function</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	tagFunction_ID         = 1 <span class="comment">// uint64</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	tagFunction_Name       = 2 <span class="comment">// int64 (string table index)</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	tagFunction_SystemName = 3 <span class="comment">// int64 (string table index)</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	tagFunction_Filename   = 4 <span class="comment">// int64 (string table index)</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	tagFunction_StartLine  = 5 <span class="comment">// int64</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">// stringIndex adds s to the string table if not already present</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// and returns the index of s in the string table.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func (b *profileBuilder) stringIndex(s string) int64 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	id, ok := b.stringMap[s]
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if !ok {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		id = len(b.strings)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		b.strings = append(b.strings, s)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		b.stringMap[s] = id
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	return int64(id)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func (b *profileBuilder) flush() {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	const dataFlush = 4096
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	if b.pb.nest == 0 &amp;&amp; len(b.pb.data) &gt; dataFlush {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		b.zw.Write(b.pb.data)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		b.pb.data = b.pb.data[:0]
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// pbValueType encodes a ValueType message to b.pb.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (b *profileBuilder) pbValueType(tag int, typ, unit string) {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	b.pb.int64(tagValueType_Type, b.stringIndex(typ))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	b.pb.int64(tagValueType_Unit, b.stringIndex(unit))
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	b.pb.endMessage(tag, start)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// pbSample encodes a Sample message to b.pb.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (b *profileBuilder) pbSample(values []int64, locs []uint64, labels func()) {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	b.pb.int64s(tagSample_Value, values)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	b.pb.uint64s(tagSample_Location, locs)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if labels != nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		labels()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	b.pb.endMessage(tagProfile_Sample, start)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	b.flush()
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// pbLabel encodes a Label message to b.pb.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>func (b *profileBuilder) pbLabel(tag int, key, str string, num int64) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	b.pb.int64Opt(tagLabel_Key, b.stringIndex(key))
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	b.pb.int64Opt(tagLabel_Str, b.stringIndex(str))
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	b.pb.int64Opt(tagLabel_Num, num)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	b.pb.endMessage(tag, start)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span><span class="comment">// pbLine encodes a Line message to b.pb.</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>func (b *profileBuilder) pbLine(tag int, funcID uint64, line int64) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagLine_FunctionID, funcID)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	b.pb.int64Opt(tagLine_Line, line)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	b.pb.endMessage(tag, start)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// pbMapping encodes a Mapping message to b.pb.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func (b *profileBuilder) pbMapping(tag int, id, base, limit, offset uint64, file, buildID string, hasFuncs bool) {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagMapping_ID, id)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagMapping_Start, base)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagMapping_Limit, limit)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagMapping_Offset, offset)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	b.pb.int64Opt(tagMapping_Filename, b.stringIndex(file))
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	b.pb.int64Opt(tagMapping_BuildID, b.stringIndex(buildID))
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// TODO: we set HasFunctions if all symbols from samples were symbolized (hasFuncs).</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// Decide what to do about HasInlineFrames and HasLineNumbers.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// Also, another approach to handle the mapping entry with</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// incomplete symbolization results is to duplicate the mapping</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// entry (but with different Has* fields values) and use</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	<span class="comment">// different entries for symbolized locations and unsymbolized locations.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	if hasFuncs {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		b.pb.bool(tagMapping_HasFunctions, true)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	b.pb.endMessage(tag, start)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func allFrames(addr uintptr) ([]runtime.Frame, symbolizeFlag) {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// Expand this one address using CallersFrames so we can cache</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">// each expansion. In general, CallersFrames takes a whole</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// stack, but in this case we know there will be no skips in</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// the stack and we have return PCs anyway.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	frames := runtime.CallersFrames([]uintptr{addr})
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	frame, more := frames.Next()
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if frame.Function == &#34;runtime.goexit&#34; {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// Short-circuit if we see runtime.goexit so the loop</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// below doesn&#39;t allocate a useless empty location.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		return nil, 0
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	symbolizeResult := lookupTried
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if frame.PC == 0 || frame.Function == &#34;&#34; || frame.File == &#34;&#34; || frame.Line == 0 {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		symbolizeResult |= lookupFailed
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	if frame.PC == 0 {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// If we failed to resolve the frame, at least make up</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// a reasonable call PC. This mostly happens in tests.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		frame.PC = addr - 1
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	ret := []runtime.Frame{frame}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	for frame.Function != &#34;runtime.goexit&#34; &amp;&amp; more {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		frame, more = frames.Next()
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		ret = append(ret, frame)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	return ret, symbolizeResult
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>type locInfo struct {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// location id assigned by the profileBuilder</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	id uint64
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// sequence of PCs, including the fake PCs returned by the traceback</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">// to represent inlined functions</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// https://github.com/golang/go/blob/d6f2f833c93a41ec1c68e49804b8387a06b131c5/src/runtime/traceback.go#L347-L368</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	pcs []uintptr
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// firstPCFrames and firstPCSymbolizeResult hold the results of the</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// allFrames call for the first (leaf-most) PC this locInfo represents</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	firstPCFrames          []runtime.Frame
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	firstPCSymbolizeResult symbolizeFlag
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// newProfileBuilder returns a new profileBuilder.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// CPU profiling data obtained from the runtime can be added</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// by calling b.addCPUData, and then the eventual profile</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// can be obtained by calling b.finish.</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>func newProfileBuilder(w io.Writer) *profileBuilder {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	zw, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	b := &amp;profileBuilder{
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		w:         w,
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		zw:        zw,
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		start:     time.Now(),
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		strings:   []string{&#34;&#34;},
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		stringMap: map[string]int{&#34;&#34;: 0},
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		locs:      map[uintptr]locInfo{},
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		funcs:     map[string]int{},
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	b.readMapping()
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	return b
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// addCPUData adds the CPU profiling data to the profile.</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// The data must be a whole number of records, as delivered by the runtime.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// len(tags) must be equal to the number of records in data.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>func (b *profileBuilder) addCPUData(data []uint64, tags []unsafe.Pointer) error {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if !b.havePeriod {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// first record is period</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		if len(data) &lt; 3 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;truncated profile&#34;)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		if data[0] != 3 || data[2] == 0 {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;malformed profile&#34;)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// data[2] is sampling rate in Hz. Convert to sampling</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		<span class="comment">// period in nanoseconds.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		b.period = 1e9 / int64(data[2])
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		b.havePeriod = true
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		data = data[3:]
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		<span class="comment">// Consume tag slot. Note that there isn&#39;t a meaningful tag</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		<span class="comment">// value for this record.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		tags = tags[1:]
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// Parse CPU samples from the profile.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// Each sample is 3+n uint64s:</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">//	data[0] = 3+n</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">//	data[1] = time stamp (ignored)</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">//	data[2] = count</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">//	data[3:3+n] = stack</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// If the count is 0 and the stack has length 1,</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">// that&#39;s an overflow record inserted by the runtime</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// to indicate that stack[0] samples were lost.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// Otherwise the count is usually 1,</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// but in a few special cases like lost non-Go samples</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// there can be larger counts.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// Because many samples with the same stack arrive,</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// we want to deduplicate immediately, which we do</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// using the b.m profMap.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	for len(data) &gt; 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		if len(data) &lt; 3 || data[0] &gt; uint64(len(data)) {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;truncated profile&#34;)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if data[0] &lt; 3 || tags != nil &amp;&amp; len(tags) &lt; 1 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;malformed profile&#34;)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		if len(tags) &lt; 1 {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;mismatched profile records and tags&#34;)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		count := data[2]
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		stk := data[3:data[0]]
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		data = data[data[0]:]
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		tag := tags[0]
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		tags = tags[1:]
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		if count == 0 &amp;&amp; len(stk) == 1 {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			<span class="comment">// overflow record</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			count = uint64(stk[0])
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			stk = []uint64{
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				<span class="comment">// gentraceback guarantees that PCs in the</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				<span class="comment">// stack can be unconditionally decremented and</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				<span class="comment">// still be valid, so we must do the same.</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				uint64(abi.FuncPCABIInternal(lostProfileEvent) + 1),
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			}
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		b.m.lookup(stk, tag).count += int64(count)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if len(tags) != 0 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;mismatched profile records and tags&#34;)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	return nil
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// build completes and returns the constructed profile.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>func (b *profileBuilder) build() {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	b.end = time.Now()
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	b.pb.int64Opt(tagProfile_TimeNanos, b.start.UnixNano())
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	if b.havePeriod { <span class="comment">// must be CPU profile</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		b.pbValueType(tagProfile_SampleType, &#34;samples&#34;, &#34;count&#34;)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		b.pbValueType(tagProfile_SampleType, &#34;cpu&#34;, &#34;nanoseconds&#34;)
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		b.pb.int64Opt(tagProfile_DurationNanos, b.end.Sub(b.start).Nanoseconds())
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		b.pbValueType(tagProfile_PeriodType, &#34;cpu&#34;, &#34;nanoseconds&#34;)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		b.pb.int64Opt(tagProfile_Period, b.period)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	values := []int64{0, 0}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	var locs []uint64
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	for e := b.m.all; e != nil; e = e.nextAll {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		values[0] = e.count
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		values[1] = e.count * b.period
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		var labels func()
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		if e.tag != nil {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			labels = func() {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				for k, v := range *(*labelMap)(e.tag) {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>					b.pbLabel(tagSample_Label, k, v, 0)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>				}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>			}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		locs = b.appendLocsForStack(locs[:0], e.stk)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		b.pbSample(values, locs, labels)
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	for i, m := range b.mem {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		hasFunctions := m.funcs == lookupTried <span class="comment">// lookupTried but not lookupFailed</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		b.pbMapping(tagProfile_Mapping, uint64(i+1), uint64(m.start), uint64(m.end), m.offset, m.file, m.buildID, hasFunctions)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// TODO: Anything for tagProfile_DropFrames?</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// TODO: Anything for tagProfile_KeepFrames?</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	b.pb.strings(tagProfile_StringTable, b.strings)
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	b.zw.Write(b.pb.data)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	b.zw.Close()
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// appendLocsForStack appends the location IDs for the given stack trace to the given</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// location ID slice, locs. The addresses in the stack are return PCs or 1 + the PC of</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// an inline marker as the runtime traceback function returns.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">// It may return an empty slice even if locs is non-empty, for example if locs consists</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span><span class="comment">// solely of runtime.goexit. We still count these empty stacks in profiles in order to</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span><span class="comment">// get the right cumulative sample count.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// It may emit to b.pb, so there must be no message encoding in progress.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>func (b *profileBuilder) appendLocsForStack(locs []uint64, stk []uintptr) (newLocs []uint64) {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	b.deck.reset()
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	<span class="comment">// The last frame might be truncated. Recover lost inline frames.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	stk = runtime_expandFinalInlineFrame(stk)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	for len(stk) &gt; 0 {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		addr := stk[0]
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		if l, ok := b.locs[addr]; ok {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			<span class="comment">// When generating code for an inlined function, the compiler adds</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			<span class="comment">// NOP instructions to the outermost function as a placeholder for</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>			<span class="comment">// each layer of inlining. When the runtime generates tracebacks for</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			<span class="comment">// stacks that include inlined functions, it uses the addresses of</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>			<span class="comment">// those NOPs as &#34;fake&#34; PCs on the stack as if they were regular</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			<span class="comment">// function call sites. But if a profiling signal arrives while the</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			<span class="comment">// CPU is executing one of those NOPs, its PC will show up as a leaf</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			<span class="comment">// in the profile with its own Location entry. So, always check</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			<span class="comment">// whether addr is a &#34;fake&#34; PC in the context of the current call</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			<span class="comment">// stack by trying to add it to the inlining deck before assuming</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			<span class="comment">// that the deck is complete.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			if len(b.deck.pcs) &gt; 0 {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>				if added := b.deck.tryAdd(addr, l.firstPCFrames, l.firstPCSymbolizeResult); added {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>					stk = stk[1:]
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>					continue
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>				}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			<span class="comment">// first record the location if there is any pending accumulated info.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			if id := b.emitLocation(); id &gt; 0 {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>				locs = append(locs, id)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			<span class="comment">// then, record the cached location.</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			locs = append(locs, l.id)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			<span class="comment">// Skip the matching pcs.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			<span class="comment">// Even if stk was truncated due to the stack depth</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			<span class="comment">// limit, expandFinalInlineFrame above has already</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			<span class="comment">// fixed the truncation, ensuring it is long enough.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			stk = stk[len(l.pcs):]
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			continue
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		frames, symbolizeResult := allFrames(addr)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		if len(frames) == 0 { <span class="comment">// runtime.goexit.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			if id := b.emitLocation(); id &gt; 0 {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				locs = append(locs, id)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			stk = stk[1:]
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			continue
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		if added := b.deck.tryAdd(addr, frames, symbolizeResult); added {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			stk = stk[1:]
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			continue
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		<span class="comment">// add failed because this addr is not inlined with the</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		<span class="comment">// existing PCs in the deck. Flush the deck and retry handling</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		<span class="comment">// this pc.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		if id := b.emitLocation(); id &gt; 0 {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			locs = append(locs, id)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		<span class="comment">// check cache again - previous emitLocation added a new entry</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		if l, ok := b.locs[addr]; ok {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			locs = append(locs, l.id)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			stk = stk[len(l.pcs):] <span class="comment">// skip the matching pcs.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		} else {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			b.deck.tryAdd(addr, frames, symbolizeResult) <span class="comment">// must succeed.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			stk = stk[1:]
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	if id := b.emitLocation(); id &gt; 0 { <span class="comment">// emit remaining location.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		locs = append(locs, id)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	return locs
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">// Here&#39;s an example of how Go 1.17 writes out inlined functions, compiled for</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">// linux/amd64. The disassembly of main.main shows two levels of inlining: main</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// calls b, b calls a, a does some work.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span><span class="comment">//   inline.go:9   0x4553ec  90              NOPL                 // func main()    { b(v) }</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span><span class="comment">//   inline.go:6   0x4553ed  90              NOPL                 // func b(v *int) { a(v) }</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span><span class="comment">//   inline.go:5   0x4553ee  48c7002a000000  MOVQ $0x2a, 0(AX)    // func a(v *int) { *v = 42 }</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span><span class="comment">// If a profiling signal arrives while executing the MOVQ at 0x4553ee (for line</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span><span class="comment">// 5), the runtime will report the stack as the MOVQ frame being called by the</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span><span class="comment">// NOPL at 0x4553ed (for line 6) being called by the NOPL at 0x4553ec (for line</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span><span class="comment">// 9).</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span><span class="comment">// The role of pcDeck is to collapse those three frames back into a single</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span><span class="comment">// location at 0x4553ee, with file/line/function symbolization info representing</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span><span class="comment">// the three layers of calls. It does that via sequential calls to pcDeck.tryAdd</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span><span class="comment">// starting with the leaf-most address. The fourth call to pcDeck.tryAdd will be</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span><span class="comment">// for the caller of main.main. Because main.main was not inlined in its caller,</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span><span class="comment">// the deck will reject the addition, and the fourth PC on the stack will get</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span><span class="comment">// its own location.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// pcDeck is a helper to detect a sequence of inlined functions from</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">// a stack trace returned by the runtime.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span><span class="comment">// The stack traces returned by runtime&#39;s trackback functions are fully</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span><span class="comment">// expanded (at least for Go functions) and include the fake pcs representing</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span><span class="comment">// inlined functions. The profile proto expects the inlined functions to be</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span><span class="comment">// encoded in one Location message.</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span><span class="comment">// https://github.com/google/pprof/blob/5e965273ee43930341d897407202dd5e10e952cb/proto/profile.proto#L177-L184</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// Runtime does not directly expose whether a frame is for an inlined function</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// and looking up debug info is not ideal, so we use a heuristic to filter</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">// the fake pcs and restore the inlined and entry functions. Inlined functions</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// have the following properties:</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">//	Frame&#39;s Func is nil (note: also true for non-Go functions), and</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//	Frame&#39;s Entry matches its entry function frame&#39;s Entry (note: could also be true for recursive calls and non-Go functions), and</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">//	Frame&#39;s Name does not match its entry function frame&#39;s name (note: inlined functions cannot be directly recursive).</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span><span class="comment">// As reading and processing the pcs in a stack trace one by one (from leaf to the root),</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// we use pcDeck to temporarily hold the observed pcs and their expanded frames</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// until we observe the entry function frame.</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>type pcDeck struct {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	pcs             []uintptr
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	frames          []runtime.Frame
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	symbolizeResult symbolizeFlag
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	<span class="comment">// firstPCFrames indicates the number of frames associated with the first</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// (leaf-most) PC in the deck</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	firstPCFrames int
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	<span class="comment">// firstPCSymbolizeResult holds the results of the allFrames call for the</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	<span class="comment">// first (leaf-most) PC in the deck</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	firstPCSymbolizeResult symbolizeFlag
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>func (d *pcDeck) reset() {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	d.pcs = d.pcs[:0]
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	d.frames = d.frames[:0]
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	d.symbolizeResult = 0
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	d.firstPCFrames = 0
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	d.firstPCSymbolizeResult = 0
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span><span class="comment">// tryAdd tries to add the pc and Frames expanded from it (most likely one,</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span><span class="comment">// since the stack trace is already fully expanded) and the symbolizeResult</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">// to the deck. If it fails the caller needs to flush the deck and retry.</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>func (d *pcDeck) tryAdd(pc uintptr, frames []runtime.Frame, symbolizeResult symbolizeFlag) (success bool) {
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	if existing := len(d.frames); existing &gt; 0 {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// &#39;d.frames&#39; are all expanded from one &#39;pc&#39; and represent all</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// inlined functions so we check only the last one.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		newFrame := frames[0]
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		last := d.frames[existing-1]
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		if last.Func != nil { <span class="comment">// the last frame can&#39;t be inlined. Flush.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			return false
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		if last.Entry == 0 || newFrame.Entry == 0 { <span class="comment">// Possibly not a Go function. Don&#39;t try to merge.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			return false
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		if last.Entry != newFrame.Entry { <span class="comment">// newFrame is for a different function.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			return false
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>		}
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		if runtime_FrameSymbolName(&amp;last) == runtime_FrameSymbolName(&amp;newFrame) { <span class="comment">// maybe recursion.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			return false
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	d.pcs = append(d.pcs, pc)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	d.frames = append(d.frames, frames...)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	d.symbolizeResult |= symbolizeResult
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	if len(d.pcs) == 1 {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		d.firstPCFrames = len(d.frames)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		d.firstPCSymbolizeResult = symbolizeResult
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	return true
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">// emitLocation emits the new location and function information recorded in the deck</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span><span class="comment">// and returns the location ID encoded in the profile protobuf.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span><span class="comment">// It emits to b.pb, so there must be no message encoding in progress.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span><span class="comment">// It resets the deck.</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>func (b *profileBuilder) emitLocation() uint64 {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if len(b.deck.pcs) == 0 {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return 0
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	defer b.deck.reset()
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	addr := b.deck.pcs[0]
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	firstFrame := b.deck.frames[0]
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	<span class="comment">// We can&#39;t write out functions while in the middle of the</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	<span class="comment">// Location message, so record new functions we encounter and</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	<span class="comment">// write them out after the Location.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	type newFunc struct {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		id         uint64
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		name, file string
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		startLine  int64
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	newFuncs := make([]newFunc, 0, 8)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	id := uint64(len(b.locs)) + 1
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	b.locs[addr] = locInfo{
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		id:                     id,
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		pcs:                    append([]uintptr{}, b.deck.pcs...),
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		firstPCSymbolizeResult: b.deck.firstPCSymbolizeResult,
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		firstPCFrames:          append([]runtime.Frame{}, b.deck.frames[:b.deck.firstPCFrames]...),
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	start := b.pb.startMessage()
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagLocation_ID, id)
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	b.pb.uint64Opt(tagLocation_Address, uint64(firstFrame.PC))
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	for _, frame := range b.deck.frames {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		<span class="comment">// Write out each line in frame expansion.</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		funcName := runtime_FrameSymbolName(&amp;frame)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		funcID := uint64(b.funcs[funcName])
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		if funcID == 0 {
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			funcID = uint64(len(b.funcs)) + 1
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>			b.funcs[funcName] = int(funcID)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			newFuncs = append(newFuncs, newFunc{
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>				id:        funcID,
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>				name:      funcName,
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>				file:      frame.File,
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				startLine: int64(runtime_FrameStartLine(&amp;frame)),
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>			})
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		b.pbLine(tagLocation_Line, funcID, int64(frame.Line))
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	for i := range b.mem {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>		if b.mem[i].start &lt;= addr &amp;&amp; addr &lt; b.mem[i].end || b.mem[i].fake {
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			b.pb.uint64Opt(tagLocation_MappingID, uint64(i+1))
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			m := b.mem[i]
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			m.funcs |= b.deck.symbolizeResult
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			b.mem[i] = m
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			break
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	b.pb.endMessage(tagProfile_Location, start)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	<span class="comment">// Write out functions we found during frame expansion.</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	for _, fn := range newFuncs {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		start := b.pb.startMessage()
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		b.pb.uint64Opt(tagFunction_ID, fn.id)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		b.pb.int64Opt(tagFunction_Name, b.stringIndex(fn.name))
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		b.pb.int64Opt(tagFunction_SystemName, b.stringIndex(fn.name))
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		b.pb.int64Opt(tagFunction_Filename, b.stringIndex(fn.file))
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		b.pb.int64Opt(tagFunction_StartLine, fn.startLine)
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		b.pb.endMessage(tagProfile_Function, start)
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	b.flush()
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	return id
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>var space = []byte(&#34; &#34;)
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>var newline = []byte(&#34;\n&#34;)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>func parseProcSelfMaps(data []byte, addMapping func(lo, hi, offset uint64, file, buildID string)) {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// $ cat /proc/self/maps</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">// 00400000-0040b000 r-xp 00000000 fc:01 787766                             /bin/cat</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// 0060a000-0060b000 r--p 0000a000 fc:01 787766                             /bin/cat</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// 0060b000-0060c000 rw-p 0000b000 fc:01 787766                             /bin/cat</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	<span class="comment">// 014ab000-014cc000 rw-p 00000000 00:00 0                                  [heap]</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	<span class="comment">// 7f7d76af8000-7f7d7797c000 r--p 00000000 fc:01 1318064                    /usr/lib/locale/locale-archive</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// 7f7d7797c000-7f7d77b36000 r-xp 00000000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77b36000-7f7d77d36000 ---p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77d36000-7f7d77d3a000 r--p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77d3a000-7f7d77d3c000 rw-p 001be000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77d3c000-7f7d77d41000 rw-p 00000000 00:00 0</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77d41000-7f7d77d64000 r-xp 00000000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77f3f000-7f7d77f42000 rw-p 00000000 00:00 0</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77f61000-7f7d77f63000 rw-p 00000000 00:00 0</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77f63000-7f7d77f64000 r--p 00022000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77f64000-7f7d77f65000 rw-p 00023000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// 7f7d77f65000-7f7d77f66000 rw-p 00000000 00:00 0</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// 7ffc342a2000-7ffc342c3000 rw-p 00000000 00:00 0                          [stack]</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// 7ffc34343000-7ffc34345000 r-xp 00000000 00:00 0                          [vdso]</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">// ffffffffff600000-ffffffffff601000 r-xp 00000000 00:00 0                  [vsyscall]</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	var line []byte
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	<span class="comment">// next removes and returns the next field in the line.</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	<span class="comment">// It also removes from line any spaces following the field.</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	next := func() []byte {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		var f []byte
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		f, line, _ = bytes.Cut(line, space)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		line = bytes.TrimLeft(line, &#34; &#34;)
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		return f
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	for len(data) &gt; 0 {
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		line, data, _ = bytes.Cut(data, newline)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		addr := next()
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		loStr, hiStr, ok := strings.Cut(string(addr), &#34;-&#34;)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		if !ok {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			continue
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		lo, err := strconv.ParseUint(loStr, 16, 64)
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		if err != nil {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			continue
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		hi, err := strconv.ParseUint(hiStr, 16, 64)
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		if err != nil {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>			continue
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>		perm := next()
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		if len(perm) &lt; 4 || perm[2] != &#39;x&#39; {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			<span class="comment">// Only interested in executable mappings.</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			continue
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		offset, err := strconv.ParseUint(string(next()), 16, 64)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>		if err != nil {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			continue
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		next()          <span class="comment">// dev</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		inode := next() <span class="comment">// inode</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		if line == nil {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>			continue
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		file := string(line)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		<span class="comment">// Trim deleted file marker.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		deletedStr := &#34; (deleted)&#34;
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		deletedLen := len(deletedStr)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		if len(file) &gt;= deletedLen &amp;&amp; file[len(file)-deletedLen:] == deletedStr {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			file = file[:len(file)-deletedLen]
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		if len(inode) == 1 &amp;&amp; inode[0] == &#39;0&#39; &amp;&amp; file == &#34;&#34; {
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			<span class="comment">// Huge-page text mappings list the initial fragment of</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			<span class="comment">// mapped but unpopulated memory as being inode 0.</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t report that part.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			<span class="comment">// But [vdso] and [vsyscall] are inode 0, so let non-empty file names through.</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			continue
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>		<span class="comment">// TODO: pprof&#39;s remapMappingIDs makes one adjustment:</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		<span class="comment">// 1. If there is an /anon_hugepage mapping first and it is</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		<span class="comment">// consecutive to a next mapping, drop the /anon_hugepage.</span>
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s no indication why this is needed.</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		<span class="comment">// Let&#39;s try not doing this and see what breaks.</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		<span class="comment">// If we do need it, it would go here, before we</span>
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		<span class="comment">// enter the mappings into b.mem in the first place.</span>
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		buildID, _ := elfBuildID(file)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		addMapping(lo, hi, offset, file, buildID)
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	}
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>func (b *profileBuilder) addMapping(lo, hi, offset uint64, file, buildID string) {
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	b.addMappingEntry(lo, hi, offset, file, buildID, false)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>func (b *profileBuilder) addMappingEntry(lo, hi, offset uint64, file, buildID string, fake bool) {
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>	b.mem = append(b.mem, memMap{
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		start:   uintptr(lo),
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		end:     uintptr(hi),
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		offset:  offset,
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		file:    file,
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		buildID: buildID,
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		fake:    fake,
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	})
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>
</pre><p><a href="proto.go?m=text">View as plain text</a></p>

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
