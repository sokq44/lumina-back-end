<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/internal/trace/v2/base.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="base.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/internal">internal</a>/<a href="http://localhost:8080/src/internal/trace">trace</a>/<a href="http://localhost:8080/src/internal/trace/v2">v2</a>/<span class="text-muted">base.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/internal/trace/v2">internal/trace/v2</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2023 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file contains data types that all implementations of the trace format</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// parser need to provide to the rest of the package.</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>package trace
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>import (
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;internal/trace/v2/event&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;internal/trace/v2/event/go122&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;internal/trace/v2/version&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// maxArgs is the maximum number of arguments for &#34;plain&#34; events,</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// i.e. anything that could reasonably be represented as a Base.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>const maxArgs = 5
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// baseEvent is the basic unprocessed event. This serves as a common</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// fundamental data structure across.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type baseEvent struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	typ  event.Type
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	time Time
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	args [maxArgs - 1]uint64
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// extra returns a slice representing extra available space in args</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// that the parser can use to pass data up into Event.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func (e *baseEvent) extra(v version.Version) []uint64 {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	switch v {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	case version.Go122:
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		return e.args[len(go122.Specs()[e.typ].Args)-1:]
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	panic(fmt.Sprintf(&#34;unsupported version: go 1.%d&#34;, v))
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// evTable contains the per-generation data necessary to</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// interpret an individual event.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>type evTable struct {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	freq    frequency
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	strings dataTable[stringID, string]
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	stacks  dataTable[stackID, stack]
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	<span class="comment">// extraStrings are strings that get generated during</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// parsing but haven&#39;t come directly from the trace, so</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// they don&#39;t appear in strings.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	extraStrings   []string
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	extraStringIDs map[string]extraStringID
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	nextExtra      extraStringID
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// addExtraString adds an extra string to the evTable and returns</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// a unique ID for the string in the table.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (t *evTable) addExtraString(s string) extraStringID {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if s == &#34;&#34; {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return 0
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if t.extraStringIDs == nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		t.extraStringIDs = make(map[string]extraStringID)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	if id, ok := t.extraStringIDs[s]; ok {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return id
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	t.nextExtra++
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	id := t.nextExtra
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	t.extraStrings = append(t.extraStrings, s)
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	t.extraStringIDs[s] = id
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	return id
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// getExtraString returns the extra string for the provided ID.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// The ID must have been produced by addExtraString for this evTable.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (t *evTable) getExtraString(id extraStringID) string {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if id == 0 {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return &#34;&#34;
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	return t.extraStrings[id-1]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// dataTable is a mapping from EIs to Es.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>type dataTable[EI ~uint64, E any] struct {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	present []uint8
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	dense   []E
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	sparse  map[EI]E
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// insert tries to add a mapping from id to s.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Returns an error if a mapping for id already exists, regardless</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// of whether or not s is the same in content. This should be used</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// for validation during parsing.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func (d *dataTable[EI, E]) insert(id EI, data E) error {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if d.sparse == nil {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		d.sparse = make(map[EI]E)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if existing, ok := d.get(id); ok {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;multiple %Ts with the same ID: id=%d, new=%v, existing=%v&#34;, data, id, data, existing)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	d.sparse[id] = data
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	return nil
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// compactify attempts to compact sparse into dense.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// This is intended to be called only once after insertions are done.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func (d *dataTable[EI, E]) compactify() {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if d.sparse == nil || len(d.dense) != 0 {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		<span class="comment">// Already compactified.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Find the range of IDs.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	maxID := EI(0)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	minID := ^EI(0)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	for id := range d.sparse {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		if id &gt; maxID {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			maxID = id
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		if id &lt; minID {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			minID = id
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if maxID &gt;= math.MaxInt {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		<span class="comment">// We can&#39;t create a slice big enough to hold maxID elements</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re willing to waste at most 2x memory.</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if int(maxID-minID) &gt; max(len(d.sparse), 2*len(d.sparse)) {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if int(minID) &gt; len(d.sparse) {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	size := int(maxID) + 1
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	d.present = make([]uint8, (size+7)/8)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	d.dense = make([]E, size)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for id, data := range d.sparse {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		d.dense[id] = data
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		d.present[id/8] |= uint8(1) &lt;&lt; (id % 8)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	d.sparse = nil
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// get returns the E for id or false if it doesn&#39;t</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// exist. This should be used for validation during parsing.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>func (d *dataTable[EI, E]) get(id EI) (E, bool) {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if id == 0 {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		return *new(E), true
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	if uint64(id) &lt; uint64(len(d.dense)) {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		if d.present[id/8]&amp;(uint8(1)&lt;&lt;(id%8)) != 0 {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			return d.dense[id], true
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	} else if d.sparse != nil {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		if data, ok := d.sparse[id]; ok {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			return data, true
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	return *new(E), false
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span><span class="comment">// forEach iterates over all ID/value pairs in the data table.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func (d *dataTable[EI, E]) forEach(yield func(EI, E) bool) bool {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	for id, value := range d.dense {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if d.present[id/8]&amp;(uint8(1)&lt;&lt;(id%8)) == 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			continue
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if !yield(EI(id), value) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return false
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if d.sparse == nil {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		return true
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	for id, value := range d.sparse {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		if !yield(id, value) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			return false
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	return true
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// mustGet returns the E for id or panics if it fails.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// This should only be used if id has already been validated.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>func (d *dataTable[EI, E]) mustGet(id EI) E {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	data, ok := d.get(id)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if !ok {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;expected id %d in %T table&#34;, id, data))
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	return data
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// frequency is nanoseconds per timestamp unit.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>type frequency float64
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// mul multiplies an unprocessed to produce a time in nanoseconds.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func (f frequency) mul(t timestamp) Time {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	return Time(float64(t) * float64(f))
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// stringID is an index into the string table for a generation.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>type stringID uint64
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// extraStringID is an index into the extra string table for a generation.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>type extraStringID uint64
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// stackID is an index into the stack table for a generation.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>type stackID uint64
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// cpuSample represents a CPU profiling sample captured by the trace.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>type cpuSample struct {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	schedCtx
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	time  Time
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	stack stackID
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// asEvent produces a complete Event from a cpuSample. It needs</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// the evTable from the generation that created it.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// We don&#39;t just store it as an Event in generation to minimize</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// the amount of pointer data floating around.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func (s cpuSample) asEvent(table *evTable) Event {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// TODO(mknyszek): This is go122-specific, but shouldn&#39;t be.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// Generalize this in the future.</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	e := Event{
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		table: table,
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		ctx:   s.schedCtx,
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		base: baseEvent{
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			typ:  go122.EvCPUSample,
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			time: s.time,
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		},
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	e.base.args[0] = uint64(s.stack)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	return e
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span><span class="comment">// stack represents a goroutine stack sample.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>type stack struct {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	frames []frame
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>func (s stack) String() string {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	var sb strings.Builder
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	for _, frame := range s.frames {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		fmt.Fprintf(&amp;sb, &#34;\t%#v\n&#34;, frame)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	return sb.String()
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// frame represents a single stack frame.</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>type frame struct {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	pc     uint64
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	funcID stringID
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	fileID stringID
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	line   uint64
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
</pre><p><a href="base.go?m=text">View as plain text</a></p>

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
