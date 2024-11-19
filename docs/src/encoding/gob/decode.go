<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/gob/decode.go - Go Documentation Server</title>

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
<a href="decode.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/gob">gob</a>/<span class="text-muted">decode.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/encoding/gob">encoding/gob</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:generate go run decgen.go -output dec_helpers.go</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package gob
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;encoding&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/saferio&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>var (
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	errBadUint = errors.New(&#34;gob: encoded unsigned integer out of range&#34;)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	errBadType = errors.New(&#34;gob: unknown type id or corrupted data&#34;)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	errRange   = errors.New(&#34;gob: bad data: field numbers out of bounds&#34;)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>type decHelper func(state *decoderState, v reflect.Value, length int, ovfl error) bool
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// decoderState is the execution state of an instance of the decoder. A new state</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// is created for nested objects.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>type decoderState struct {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	dec *Decoder
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// The buffer is stored with an extra indirection because it may be replaced</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// if we load a type during decode (when reading an interface value).</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	b        *decBuffer
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	fieldnum int           <span class="comment">// the last field number read.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	next     *decoderState <span class="comment">// for free list</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// decBuffer is an extremely simple, fast implementation of a read-only byte buffer.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// It is initialized by calling Size and then copying the data into the slice returned by Bytes().</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>type decBuffer struct {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	data   []byte
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	offset int <span class="comment">// Read offset.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func (d *decBuffer) Read(p []byte) (int, error) {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	n := copy(p, d.data[d.offset:])
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	if n == 0 &amp;&amp; len(p) != 0 {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	}
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	d.offset += n
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	return n, nil
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (d *decBuffer) Drop(n int) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if n &gt; d.Len() {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		panic(&#34;drop&#34;)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	d.offset += n
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func (d *decBuffer) ReadByte() (byte, error) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if d.offset &gt;= len(d.data) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	c := d.data[d.offset]
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	d.offset++
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return c, nil
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (d *decBuffer) Len() int {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return len(d.data) - d.offset
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (d *decBuffer) Bytes() []byte {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	return d.data[d.offset:]
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// SetBytes sets the buffer to the bytes, discarding any existing data.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (d *decBuffer) SetBytes(data []byte) {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	d.data = data
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	d.offset = 0
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>func (d *decBuffer) Reset() {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	d.data = d.data[0:0]
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	d.offset = 0
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// We pass the bytes.Buffer separately for easier testing of the infrastructure</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// without requiring a full Decoder.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (dec *Decoder) newDecoderState(buf *decBuffer) *decoderState {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	d := dec.freeList
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if d == nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		d = new(decoderState)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		d.dec = dec
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	} else {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		dec.freeList = d.next
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	d.b = buf
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return d
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>func (dec *Decoder) freeDecoderState(d *decoderState) {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	d.next = dec.freeList
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	dec.freeList = d
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func overflow(name string) error {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	return errors.New(`value for &#34;` + name + `&#34; out of range`)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// decodeUintReader reads an encoded unsigned integer from an io.Reader.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// Used only by the Decoder to read the message length.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func decodeUintReader(r io.Reader, buf []byte) (x uint64, width int, err error) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	width = 1
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	n, err := io.ReadFull(r, buf[0:width])
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	if n == 0 {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	b := buf[0]
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if b &lt;= 0x7f {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return uint64(b), width, nil
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	n = -int(int8(b))
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if n &gt; uint64Size {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		err = errBadUint
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		return
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	width, err = io.ReadFull(r, buf[0:n])
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	if err != nil {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		if err == io.EOF {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			err = io.ErrUnexpectedEOF
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Could check that the high byte is zero but it&#39;s not worth it.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	for _, b := range buf[0:width] {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		x = x&lt;&lt;8 | uint64(b)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	width++ <span class="comment">// +1 for length byte</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	return
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// decodeUint reads an encoded unsigned integer from state.r.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// Does not check for overflow.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func (state *decoderState) decodeUint() (x uint64) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	b, err := state.b.ReadByte()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if err != nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		error_(err)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if b &lt;= 0x7f {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		return uint64(b)
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	n := -int(int8(b))
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if n &gt; uint64Size {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		error_(errBadUint)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	buf := state.b.Bytes()
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if len(buf) &lt; n {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		errorf(&#34;invalid uint data length %d: exceeds input size %d&#34;, n, len(buf))
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t need to check error; it&#39;s safe to loop regardless.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// Could check that the high byte is zero but it&#39;s not worth it.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	for _, b := range buf[0:n] {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		x = x&lt;&lt;8 | uint64(b)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	return x
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// decodeInt reads an encoded signed integer from state.r.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// Does not check for overflow.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>func (state *decoderState) decodeInt() int64 {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	x := state.decodeUint()
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if x&amp;1 != 0 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		return ^int64(x &gt;&gt; 1)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	return int64(x &gt;&gt; 1)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span><span class="comment">// getLength decodes the next uint and makes sure it is a possible</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// size for a data item that follows, which means it must fit in a</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// non-negative int and fit in the buffer.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func (state *decoderState) getLength() (int, bool) {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	n := int(state.decodeUint())
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if n &lt; 0 || state.b.Len() &lt; n || tooBig &lt;= n {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		return 0, false
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	return n, true
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// decOp is the signature of a decoding operator for a given type.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>type decOp func(i *decInstr, state *decoderState, v reflect.Value)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">// The &#39;instructions&#39; of the decoding machine</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>type decInstr struct {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	op    decOp
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	field int   <span class="comment">// field number of the wire type</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	index []int <span class="comment">// field access indices for destination type</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	ovfl  error <span class="comment">// error message for overflow/underflow (for arrays, of the elements)</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// ignoreUint discards a uint value with no destination.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func ignoreUint(i *decInstr, state *decoderState, v reflect.Value) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	state.decodeUint()
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// ignoreTwoUints discards a uint value with no destination. It&#39;s used to skip</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// complex values.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func ignoreTwoUints(i *decInstr, state *decoderState, v reflect.Value) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	state.decodeUint()
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	state.decodeUint()
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// Since the encoder writes no zeros, if we arrive at a decoder we have</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// a value to extract and store. The field number has already been read</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// (it&#39;s how we knew to call this decoder).</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// Each decoder is responsible for handling any indirections associated</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">// with the data structure. If any pointer so reached is nil, allocation must</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">// be done.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// decAlloc takes a value and returns a settable value that can</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// be assigned to. If the value is a pointer, decAlloc guarantees it points to storage.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// The callers to the individual decoders are expected to have used decAlloc.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// The individual decoders don&#39;t need it.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>func decAlloc(v reflect.Value) reflect.Value {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	for v.Kind() == reflect.Pointer {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		if v.IsNil() {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			v.Set(reflect.New(v.Type().Elem()))
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		v = v.Elem()
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	return v
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// decBool decodes a uint and stores it as a boolean in value.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>func decBool(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	value.SetBool(state.decodeUint() != 0)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span><span class="comment">// decInt8 decodes an integer and stores it as an int8 in value.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>func decInt8(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	v := state.decodeInt()
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if v &lt; math.MinInt8 || math.MaxInt8 &lt; v {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	value.SetInt(v)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// decUint8 decodes an unsigned integer and stores it as a uint8 in value.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func decUint8(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	v := state.decodeUint()
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if math.MaxUint8 &lt; v {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	value.SetUint(v)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// decInt16 decodes an integer and stores it as an int16 in value.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>func decInt16(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	v := state.decodeInt()
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	if v &lt; math.MinInt16 || math.MaxInt16 &lt; v {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	value.SetInt(v)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// decUint16 decodes an unsigned integer and stores it as a uint16 in value.</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>func decUint16(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	v := state.decodeUint()
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	if math.MaxUint16 &lt; v {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	value.SetUint(v)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// decInt32 decodes an integer and stores it as an int32 in value.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>func decInt32(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	v := state.decodeInt()
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if v &lt; math.MinInt32 || math.MaxInt32 &lt; v {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	value.SetInt(v)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// decUint32 decodes an unsigned integer and stores it as a uint32 in value.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>func decUint32(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	v := state.decodeUint()
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if math.MaxUint32 &lt; v {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		error_(i.ovfl)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	value.SetUint(v)
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// decInt64 decodes an integer and stores it as an int64 in value.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>func decInt64(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	v := state.decodeInt()
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	value.SetInt(v)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// decUint64 decodes an unsigned integer and stores it as a uint64 in value.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func decUint64(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	v := state.decodeUint()
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	value.SetUint(v)
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// Floating-point numbers are transmitted as uint64s holding the bits</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// of the underlying representation. They are sent byte-reversed, with</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// the exponent end coming out first, so integer floating point numbers</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span><span class="comment">// (for example) transmit more compactly. This routine does the</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// unswizzling.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>func float64FromBits(u uint64) float64 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	v := bits.ReverseBytes64(u)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	return math.Float64frombits(v)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// float32FromBits decodes an unsigned integer, treats it as a 32-bit floating-point</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// number, and returns it. It&#39;s a helper function for float32 and complex64.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// It returns a float64 because that&#39;s what reflection needs, but its return</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span><span class="comment">// value is known to be accurately representable in a float32.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>func float32FromBits(u uint64, ovfl error) float64 {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	v := float64FromBits(u)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	av := v
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	if av &lt; 0 {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		av = -av
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// +Inf is OK in both 32- and 64-bit floats. Underflow is always OK.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if math.MaxFloat32 &lt; av &amp;&amp; av &lt;= math.MaxFloat64 {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		error_(ovfl)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	return v
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// decFloat32 decodes an unsigned integer, treats it as a 32-bit floating-point</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// number, and stores it in value.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func decFloat32(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	value.SetFloat(float32FromBits(state.decodeUint(), i.ovfl))
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// decFloat64 decodes an unsigned integer, treats it as a 64-bit floating-point</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// number, and stores it in value.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>func decFloat64(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	value.SetFloat(float64FromBits(state.decodeUint()))
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// decComplex64 decodes a pair of unsigned integers, treats them as a</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// pair of floating point numbers, and stores them as a complex64 in value.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// The real part comes first.</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>func decComplex64(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	real := float32FromBits(state.decodeUint(), i.ovfl)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	imag := float32FromBits(state.decodeUint(), i.ovfl)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	value.SetComplex(complex(real, imag))
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// decComplex128 decodes a pair of unsigned integers, treats them as a</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// pair of floating point numbers, and stores them as a complex128 in value.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">// The real part comes first.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>func decComplex128(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	real := float64FromBits(state.decodeUint())
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	imag := float64FromBits(state.decodeUint())
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	value.SetComplex(complex(real, imag))
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span><span class="comment">// decUint8Slice decodes a byte slice and stores in value a slice header</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// describing the data.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// uint8 slices are encoded as an unsigned count followed by the raw bytes.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>func decUint8Slice(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	if !ok {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		errorf(&#34;bad %s slice length: %d&#34;, value.Type(), n)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	if value.Cap() &lt; n {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		safe := saferio.SliceCap[byte](uint64(n))
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		if safe &lt; 0 {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			errorf(&#34;%s slice too big: %d elements&#34;, value.Type(), n)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		value.Set(reflect.MakeSlice(value.Type(), safe, safe))
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		ln := safe
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		i := 0
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		for i &lt; n {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			if i &gt;= ln {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>				<span class="comment">// We didn&#39;t allocate the entire slice,</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>				<span class="comment">// due to using saferio.SliceCap.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>				<span class="comment">// Grow the slice for one more element.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>				<span class="comment">// The slice is full, so this should</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>				<span class="comment">// bump up the capacity.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>				value.Grow(1)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			<span class="comment">// Copy into s up to the capacity or n,</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			<span class="comment">// whichever is less.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			ln = value.Cap()
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			if ln &gt; n {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>				ln = n
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>			value.SetLen(ln)
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			sub := value.Slice(i, ln)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			if _, err := state.b.Read(sub.Bytes()); err != nil {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>				errorf(&#34;error decoding []byte at %d: %s&#34;, i, err)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>			}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>			i = ln
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	} else {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		value.SetLen(n)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if _, err := state.b.Read(value.Bytes()); err != nil {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			errorf(&#34;error decoding []byte: %s&#34;, err)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		}
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// decString decodes byte array and stores in value a string header</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// describing the data.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// Strings are encoded as an unsigned count followed by the raw bytes.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>func decString(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if !ok {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		errorf(&#34;bad %s slice length: %d&#34;, value.Type(), n)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// Read the data.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	data := state.b.Bytes()
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if len(data) &lt; n {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		errorf(&#34;invalid string length %d: exceeds input size %d&#34;, n, len(data))
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	s := string(data[:n])
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	value.SetString(s)
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// ignoreUint8Array skips over the data for a byte slice value with no destination.</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>func ignoreUint8Array(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	if !ok {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		errorf(&#34;slice length too large&#34;)
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	bn := state.b.Len()
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	if bn &lt; n {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		errorf(&#34;invalid slice length %d: exceeds input size %d&#34;, n, bn)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// Execution engine</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span><span class="comment">// The encoder engine is an array of instructions indexed by field number of the incoming</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">// decoder. It is executed with random access according to field number.</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>type decEngine struct {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	instr    []decInstr
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	numInstr int <span class="comment">// the number of active instructions</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// decodeSingle decodes a top-level value that is not a struct and stores it in value.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// Such values are preceded by a zero, making them have the memory layout of a</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// struct field (although with an illegal field number).</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>func (dec *Decoder) decodeSingle(engine *decEngine, value reflect.Value) {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	state := dec.newDecoderState(&amp;dec.buf)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	defer dec.freeDecoderState(state)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	state.fieldnum = singletonField
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	if state.decodeUint() != 0 {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		errorf(&#34;decode: corrupted data: non-zero delta for singleton&#34;)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	instr := &amp;engine.instr[singletonField]
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	instr.op(instr, state, value)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// decodeStruct decodes a top-level struct and stores it in value.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">// Indir is for the value, not the type. At the time of the call it may</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span><span class="comment">// differ from ut.indir, which was computed when the engine was built.</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">// This state cannot arise for decodeSingle, which is called directly</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">// from the user&#39;s value, not from the innards of an engine.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>func (dec *Decoder) decodeStruct(engine *decEngine, value reflect.Value) {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	state := dec.newDecoderState(&amp;dec.buf)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	defer dec.freeDecoderState(state)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	for state.b.Len() &gt; 0 {
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		delta := int(state.decodeUint())
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		if delta &lt; 0 {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			errorf(&#34;decode: corrupted data: negative delta&#34;)
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		if delta == 0 { <span class="comment">// struct terminator is zero delta fieldnum</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			break
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		if state.fieldnum &gt;= len(engine.instr)-delta { <span class="comment">// subtract to compare without overflow</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			error_(errRange)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		fieldnum := state.fieldnum + delta
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		instr := &amp;engine.instr[fieldnum]
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		var field reflect.Value
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		if instr.index != nil {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			<span class="comment">// Otherwise the field is unknown to us and instr.op is an ignore op.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			field = value.FieldByIndex(instr.index)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			if field.Kind() == reflect.Pointer {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>				field = decAlloc(field)
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>			}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		instr.op(instr, state, field)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		state.fieldnum = fieldnum
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>}
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>var noValue reflect.Value
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span><span class="comment">// ignoreStruct discards the data for a struct with no destination.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>func (dec *Decoder) ignoreStruct(engine *decEngine) {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	state := dec.newDecoderState(&amp;dec.buf)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	defer dec.freeDecoderState(state)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	for state.b.Len() &gt; 0 {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		delta := int(state.decodeUint())
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		if delta &lt; 0 {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			errorf(&#34;ignore decode: corrupted data: negative delta&#34;)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		if delta == 0 { <span class="comment">// struct terminator is zero delta fieldnum</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			break
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		fieldnum := state.fieldnum + delta
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		if fieldnum &gt;= len(engine.instr) {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			error_(errRange)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		instr := &amp;engine.instr[fieldnum]
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		instr.op(instr, state, noValue)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		state.fieldnum = fieldnum
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// ignoreSingle discards the data for a top-level non-struct value with no</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">// destination. It&#39;s used when calling Decode with a nil value.</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>func (dec *Decoder) ignoreSingle(engine *decEngine) {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	state := dec.newDecoderState(&amp;dec.buf)
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	defer dec.freeDecoderState(state)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	state.fieldnum = singletonField
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	delta := int(state.decodeUint())
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if delta != 0 {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		errorf(&#34;decode: corrupted data: non-zero delta for singleton&#34;)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	}
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	instr := &amp;engine.instr[singletonField]
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	instr.op(instr, state, noValue)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// decodeArrayHelper does the work for decoding arrays and slices.</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>func (dec *Decoder) decodeArrayHelper(state *decoderState, value reflect.Value, elemOp decOp, length int, ovfl error, helper decHelper) {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	if helper != nil &amp;&amp; helper(state, value, length, ovfl) {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		return
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	instr := &amp;decInstr{elemOp, 0, nil, ovfl}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	isPtr := value.Type().Elem().Kind() == reflect.Pointer
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	ln := value.Len()
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	for i := 0; i &lt; length; i++ {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		if state.b.Len() == 0 {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>			errorf(&#34;decoding array or slice: length exceeds input size (%d elements)&#34;, length)
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		if i &gt;= ln {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>			<span class="comment">// This is a slice that we only partially allocated.</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			<span class="comment">// Grow it up to length.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>			value.Grow(1)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>			cp := value.Cap()
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			if cp &gt; length {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>				cp = length
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>			}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>			value.SetLen(cp)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>			ln = cp
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		v := value.Index(i)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		if isPtr {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			v = decAlloc(v)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		}
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		elemOp(instr, state, v)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>}
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span><span class="comment">// decodeArray decodes an array and stores it in value.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// The length is an unsigned integer preceding the elements. Even though the length is redundant</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">// (it&#39;s part of the type), it&#39;s a useful check and is included in the encoding.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>func (dec *Decoder) decodeArray(state *decoderState, value reflect.Value, elemOp decOp, length int, ovfl error, helper decHelper) {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	if n := state.decodeUint(); n != uint64(length) {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		errorf(&#34;length mismatch in decodeArray&#34;)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	dec.decodeArrayHelper(state, value, elemOp, length, ovfl, helper)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span><span class="comment">// decodeIntoValue is a helper for map decoding.</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>func decodeIntoValue(state *decoderState, op decOp, isPtr bool, value reflect.Value, instr *decInstr) reflect.Value {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	v := value
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	if isPtr {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		v = decAlloc(value)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	op(instr, state, v)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	return value
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// decodeMap decodes a map and stores it in value.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// Maps are encoded as a length followed by key:value pairs.</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span><span class="comment">// Because the internals of maps are not visible to us, we must</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span><span class="comment">// use reflection rather than pointer magic.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>func (dec *Decoder) decodeMap(mtyp reflect.Type, state *decoderState, value reflect.Value, keyOp, elemOp decOp, ovfl error) {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	n := int(state.decodeUint())
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	if value.IsNil() {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		value.Set(reflect.MakeMapWithSize(mtyp, n))
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	keyIsPtr := mtyp.Key().Kind() == reflect.Pointer
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	elemIsPtr := mtyp.Elem().Kind() == reflect.Pointer
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	keyInstr := &amp;decInstr{keyOp, 0, nil, ovfl}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	elemInstr := &amp;decInstr{elemOp, 0, nil, ovfl}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	keyP := reflect.New(mtyp.Key())
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	elemP := reflect.New(mtyp.Elem())
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		key := decodeIntoValue(state, keyOp, keyIsPtr, keyP.Elem(), keyInstr)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		elem := decodeIntoValue(state, elemOp, elemIsPtr, elemP.Elem(), elemInstr)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		value.SetMapIndex(key, elem)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		keyP.Elem().SetZero()
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		elemP.Elem().SetZero()
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// ignoreArrayHelper does the work for discarding arrays and slices.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>func (dec *Decoder) ignoreArrayHelper(state *decoderState, elemOp decOp, length int) {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	instr := &amp;decInstr{elemOp, 0, nil, errors.New(&#34;no error&#34;)}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	for i := 0; i &lt; length; i++ {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		if state.b.Len() == 0 {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>			errorf(&#34;decoding array or slice: length exceeds input size (%d elements)&#34;, length)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		elemOp(instr, state, noValue)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// ignoreArray discards the data for an array value with no destination.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>func (dec *Decoder) ignoreArray(state *decoderState, elemOp decOp, length int) {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	if n := state.decodeUint(); n != uint64(length) {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		errorf(&#34;length mismatch in ignoreArray&#34;)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	dec.ignoreArrayHelper(state, elemOp, length)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// ignoreMap discards the data for a map value with no destination.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>func (dec *Decoder) ignoreMap(state *decoderState, keyOp, elemOp decOp) {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	n := int(state.decodeUint())
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	keyInstr := &amp;decInstr{keyOp, 0, nil, errors.New(&#34;no error&#34;)}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	elemInstr := &amp;decInstr{elemOp, 0, nil, errors.New(&#34;no error&#34;)}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		keyOp(keyInstr, state, noValue)
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		elemOp(elemInstr, state, noValue)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>}
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span><span class="comment">// decodeSlice decodes a slice and stores it in value.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">// Slices are encoded as an unsigned length followed by the elements.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>func (dec *Decoder) decodeSlice(state *decoderState, value reflect.Value, elemOp decOp, ovfl error, helper decHelper) {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	u := state.decodeUint()
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	typ := value.Type()
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	size := uint64(typ.Elem().Size())
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	nBytes := u * size
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	n := int(u)
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	<span class="comment">// Take care with overflow in this calculation.</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	if n &lt; 0 || uint64(n) != u || nBytes &gt; tooBig || (size &gt; 0 &amp;&amp; nBytes/size != u) {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		<span class="comment">// We don&#39;t check n against buffer length here because if it&#39;s a slice</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		<span class="comment">// of interfaces, there will be buffer reloads.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		errorf(&#34;%s slice too big: %d elements of %d bytes&#34;, typ.Elem(), u, size)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	}
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	if value.Cap() &lt; n {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		safe := saferio.SliceCapWithSize(size, uint64(n))
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		if safe &lt; 0 {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>			errorf(&#34;%s slice too big: %d elements of %d bytes&#34;, typ.Elem(), u, size)
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>		}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		value.Set(reflect.MakeSlice(typ, safe, safe))
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	} else {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		value.SetLen(n)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	dec.decodeArrayHelper(state, value, elemOp, n, ovfl, helper)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>}
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span><span class="comment">// ignoreSlice skips over the data for a slice value with no destination.</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>func (dec *Decoder) ignoreSlice(state *decoderState, elemOp decOp) {
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	dec.ignoreArrayHelper(state, elemOp, int(state.decodeUint()))
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>}
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// decodeInterface decodes an interface value and stores it in value.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span><span class="comment">// Interfaces are encoded as the name of a concrete type followed by a value.</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// If the name is empty, the value is nil and no value is sent.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>func (dec *Decoder) decodeInterface(ityp reflect.Type, state *decoderState, value reflect.Value) {
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// Read the name of the concrete type.</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	nr := state.decodeUint()
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	if nr &gt; 1&lt;&lt;31 { <span class="comment">// zero is permissible for anonymous types</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		errorf(&#34;invalid type name length %d&#34;, nr)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	}
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	if nr &gt; uint64(state.b.Len()) {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>		errorf(&#34;invalid type name length %d: exceeds input size&#34;, nr)
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	n := int(nr)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	name := state.b.Bytes()[:n]
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// Allocate the destination interface value.</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	if len(name) == 0 {
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		<span class="comment">// Copy the nil interface value to the target.</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		value.SetZero()
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		return
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	}
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	if len(name) &gt; 1024 {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		errorf(&#34;name too long (%d bytes): %.20q...&#34;, len(name), name)
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	<span class="comment">// The concrete type must be registered.</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	typi, ok := nameToConcreteType.Load(string(name))
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	if !ok {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		errorf(&#34;name not registered for interface: %q&#34;, name)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	}
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	typ := typi.(reflect.Type)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	<span class="comment">// Read the type id of the concrete value.</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	concreteId := dec.decodeTypeSequence(true)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	if concreteId &lt; 0 {
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>		error_(dec.err)
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	}
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	<span class="comment">// Byte count of value is next; we don&#39;t care what it is (it&#39;s there</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	<span class="comment">// in case we want to ignore the value by skipping it completely).</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	state.decodeUint()
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	<span class="comment">// Read the concrete value.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	v := allocValue(typ)
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	dec.decodeValue(concreteId, v)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	if dec.err != nil {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>		error_(dec.err)
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	}
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">// Assign the concrete value to the interface.</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// Tread carefully; it might not satisfy the interface.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	if !typ.AssignableTo(ityp) {
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		errorf(&#34;%s is not assignable to type %s&#34;, typ, ityp)
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	}
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	<span class="comment">// Copy the interface value to the target.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	value.Set(v)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>}
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span><span class="comment">// ignoreInterface discards the data for an interface value with no destination.</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>func (dec *Decoder) ignoreInterface(state *decoderState) {
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// Read the name of the concrete type.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	if !ok {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		errorf(&#34;bad interface encoding: name too large for buffer&#34;)
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	}
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	bn := state.b.Len()
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	if bn &lt; n {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		errorf(&#34;invalid interface value length %d: exceeds input size %d&#34;, n, bn)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>	}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	id := dec.decodeTypeSequence(true)
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	if id &lt; 0 {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		error_(dec.err)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	<span class="comment">// At this point, the decoder buffer contains a delimited value. Just toss it.</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	n, ok = state.getLength()
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	if !ok {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		errorf(&#34;bad interface encoding: data length too large for buffer&#34;)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// decodeGobDecoder decodes something implementing the GobDecoder interface.</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// The data is encoded as a byte slice.</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>func (dec *Decoder) decodeGobDecoder(ut *userTypeInfo, state *decoderState, value reflect.Value) {
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// Read the bytes for the value.</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	if !ok {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		errorf(&#34;GobDecoder: length too large for buffer&#34;)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	b := state.b.Bytes()
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	if len(b) &lt; n {
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		errorf(&#34;GobDecoder: invalid data length %d: exceeds input size %d&#34;, n, len(b))
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	b = b[:n]
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	var err error
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	<span class="comment">// We know it&#39;s one of these.</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	switch ut.externalDec {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	case xGob:
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		err = value.Interface().(GobDecoder).GobDecode(b)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	case xBinary:
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		err = value.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b)
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	case xText:
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		err = value.Interface().(encoding.TextUnmarshaler).UnmarshalText(b)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	if err != nil {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		error_(err)
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	}
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span><span class="comment">// ignoreGobDecoder discards the data for a GobDecoder value with no destination.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>func (dec *Decoder) ignoreGobDecoder(state *decoderState) {
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	<span class="comment">// Read the bytes for the value.</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	n, ok := state.getLength()
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	if !ok {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		errorf(&#34;GobDecoder: length too large for buffer&#34;)
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>	}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	bn := state.b.Len()
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	if bn &lt; n {
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>		errorf(&#34;GobDecoder: invalid data length %d: exceeds input size %d&#34;, n, bn)
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	}
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	state.b.Drop(n)
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span><span class="comment">// Index by Go types.</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>var decOpTable = [...]decOp{
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	reflect.Bool:       decBool,
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	reflect.Int8:       decInt8,
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	reflect.Int16:      decInt16,
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	reflect.Int32:      decInt32,
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	reflect.Int64:      decInt64,
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	reflect.Uint8:      decUint8,
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	reflect.Uint16:     decUint16,
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	reflect.Uint32:     decUint32,
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	reflect.Uint64:     decUint64,
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	reflect.Float32:    decFloat32,
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	reflect.Float64:    decFloat64,
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	reflect.Complex64:  decComplex64,
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	reflect.Complex128: decComplex128,
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	reflect.String:     decString,
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span><span class="comment">// Indexed by gob types.  tComplex will be added during type.init().</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>var decIgnoreOpMap = map[typeId]decOp{
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	tBool:    ignoreUint,
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	tInt:     ignoreUint,
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	tUint:    ignoreUint,
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	tFloat:   ignoreUint,
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	tBytes:   ignoreUint8Array,
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	tString:  ignoreUint8Array,
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	tComplex: ignoreTwoUints,
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>}
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span><span class="comment">// decOpFor returns the decoding op for the base type under rt and</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span><span class="comment">// the indirection count to reach it.</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>func (dec *Decoder) decOpFor(wireId typeId, rt reflect.Type, name string, inProgress map[reflect.Type]*decOp) *decOp {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	ut := userType(rt)
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	<span class="comment">// If the type implements GobEncoder, we handle it without further processing.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	if ut.externalDec != 0 {
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		return dec.gobDecodeOpFor(ut)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	<span class="comment">// If this type is already in progress, it&#39;s a recursive type (e.g. map[string]*T).</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	<span class="comment">// Return the pointer to the op we&#39;re already building.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>	if opPtr := inProgress[rt]; opPtr != nil {
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		return opPtr
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	}
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>	typ := ut.base
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	var op decOp
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	k := typ.Kind()
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>	if int(k) &lt; len(decOpTable) {
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>		op = decOpTable[k]
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	}
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	if op == nil {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		inProgress[rt] = &amp;op
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		<span class="comment">// Special cases</span>
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>		switch t := typ; t.Kind() {
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		case reflect.Array:
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>			name = &#34;element of &#34; + name
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			elemId := dec.wireType[wireId].ArrayT.Elem
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>			elemOp := dec.decOpFor(elemId, t.Elem(), name, inProgress)
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>			ovfl := overflow(name)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>			helper := decArrayHelper[t.Elem().Kind()]
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>				state.dec.decodeArray(state, value, *elemOp, t.Len(), ovfl, helper)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			}
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		case reflect.Map:
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			keyId := dec.wireType[wireId].MapT.Key
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			elemId := dec.wireType[wireId].MapT.Elem
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>			keyOp := dec.decOpFor(keyId, t.Key(), &#34;key of &#34;+name, inProgress)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>			elemOp := dec.decOpFor(elemId, t.Elem(), &#34;element of &#34;+name, inProgress)
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>			ovfl := overflow(name)
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>				state.dec.decodeMap(t, state, value, *keyOp, *elemOp, ovfl)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			}
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		case reflect.Slice:
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			name = &#34;element of &#34; + name
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>			if t.Elem().Kind() == reflect.Uint8 {
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>				op = decUint8Slice
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>				break
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>			}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			var elemId typeId
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>			if tt := builtinIdToType(wireId); tt != nil {
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>				elemId = tt.(*sliceType).Elem
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			} else {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>				elemId = dec.wireType[wireId].SliceT.Elem
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			elemOp := dec.decOpFor(elemId, t.Elem(), name, inProgress)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			ovfl := overflow(name)
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			helper := decSliceHelper[t.Elem().Kind()]
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>				state.dec.decodeSlice(state, value, *elemOp, ovfl, helper)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			}
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		case reflect.Struct:
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			<span class="comment">// Generate a closure that calls out to the engine for the nested type.</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			ut := userType(typ)
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			enginePtr, err := dec.getDecEnginePtr(wireId, ut)
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			if err != nil {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>				error_(err)
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>			}
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>				<span class="comment">// indirect through enginePtr to delay evaluation for recursive structs.</span>
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>				dec.decodeStruct(*enginePtr, value)
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>			}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		case reflect.Interface:
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>				state.dec.decodeInterface(t, state, value)
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>			}
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		}
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	}
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	if op == nil {
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		errorf(&#34;decode can&#39;t handle type %s&#34;, rt)
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	}
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	return &amp;op
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>}
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>var maxIgnoreNestingDepth = 10000
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span><span class="comment">// decIgnoreOpFor returns the decoding op for a field that has no destination.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>func (dec *Decoder) decIgnoreOpFor(wireId typeId, inProgress map[typeId]*decOp) *decOp {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	<span class="comment">// Track how deep we&#39;ve recursed trying to skip nested ignored fields.</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	dec.ignoreDepth++
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	defer func() { dec.ignoreDepth-- }()
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	if dec.ignoreDepth &gt; maxIgnoreNestingDepth {
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		error_(errors.New(&#34;invalid nesting depth&#34;))
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	}
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	<span class="comment">// If this type is already in progress, it&#39;s a recursive type (e.g. map[string]*T).</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// Return the pointer to the op we&#39;re already building.</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	if opPtr := inProgress[wireId]; opPtr != nil {
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>		return opPtr
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>	}
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	op, ok := decIgnoreOpMap[wireId]
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	if !ok {
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		inProgress[wireId] = &amp;op
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		if wireId == tInterface {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			<span class="comment">// Special case because it&#39;s a method: the ignored item might</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>			<span class="comment">// define types and we need to record their state in the decoder.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>				state.dec.ignoreInterface(state)
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>			}
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>			return &amp;op
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>		}
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		<span class="comment">// Special cases</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		wire := dec.wireType[wireId]
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		switch {
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		case wire == nil:
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>			errorf(&#34;bad data: undefined type %s&#34;, wireId.string())
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>		case wire.ArrayT != nil:
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>			elemId := wire.ArrayT.Elem
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>			elemOp := dec.decIgnoreOpFor(elemId, inProgress)
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>				state.dec.ignoreArray(state, *elemOp, wire.ArrayT.Len)
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			}
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		case wire.MapT != nil:
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			keyId := dec.wireType[wireId].MapT.Key
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			elemId := dec.wireType[wireId].MapT.Elem
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			keyOp := dec.decIgnoreOpFor(keyId, inProgress)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			elemOp := dec.decIgnoreOpFor(elemId, inProgress)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>				state.dec.ignoreMap(state, *keyOp, *elemOp)
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			}
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		case wire.SliceT != nil:
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>			elemId := wire.SliceT.Elem
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>			elemOp := dec.decIgnoreOpFor(elemId, inProgress)
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>				state.dec.ignoreSlice(state, *elemOp)
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			}
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>		case wire.StructT != nil:
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			<span class="comment">// Generate a closure that calls out to the engine for the nested type.</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>			enginePtr, err := dec.getIgnoreEnginePtr(wireId)
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			if err != nil {
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>				error_(err)
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>			}
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>				<span class="comment">// indirect through enginePtr to delay evaluation for recursive structs</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>				state.dec.ignoreStruct(*enginePtr)
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>			}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		case wire.GobEncoderT != nil, wire.BinaryMarshalerT != nil, wire.TextMarshalerT != nil:
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>				state.dec.ignoreGobDecoder(state)
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>			}
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		}
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>	}
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	if op == nil {
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>		errorf(&#34;bad data: ignore can&#39;t handle type %s&#34;, wireId.string())
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	}
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	return &amp;op
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span><span class="comment">// gobDecodeOpFor returns the op for a type that is known to implement</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span><span class="comment">// GobDecoder.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>func (dec *Decoder) gobDecodeOpFor(ut *userTypeInfo) *decOp {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>	rcvrType := ut.user
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>	if ut.decIndir == -1 {
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		rcvrType = reflect.PointerTo(rcvrType)
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	} else if ut.decIndir &gt; 0 {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		for i := int8(0); i &lt; ut.decIndir; i++ {
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>			rcvrType = rcvrType.Elem()
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		}
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>	var op decOp
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	op = func(i *decInstr, state *decoderState, value reflect.Value) {
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>		<span class="comment">// We now have the base type. We need its address if the receiver is a pointer.</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>		if value.Kind() != reflect.Pointer &amp;&amp; rcvrType.Kind() == reflect.Pointer {
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>			value = value.Addr()
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>		}
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>		state.dec.decodeGobDecoder(ut, state, value)
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	}
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	return &amp;op
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>}
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span><span class="comment">// compatibleType asks: Are these two gob Types compatible?</span>
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span><span class="comment">// Answers the question for basic types, arrays, maps and slices, plus</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span><span class="comment">// GobEncoder/Decoder pairs.</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span><span class="comment">// Structs are considered ok; fields will be checked later.</span>
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>func (dec *Decoder) compatibleType(fr reflect.Type, fw typeId, inProgress map[reflect.Type]typeId) bool {
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	if rhs, ok := inProgress[fr]; ok {
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>		return rhs == fw
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	}
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	inProgress[fr] = fw
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>	ut := userType(fr)
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	wire, ok := dec.wireType[fw]
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>	<span class="comment">// If wire was encoded with an encoding method, fr must have that method.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	<span class="comment">// And if not, it must not.</span>
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">// At most one of the booleans in ut is set.</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	<span class="comment">// We could possibly relax this constraint in the future in order to</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	<span class="comment">// choose the decoding method using the data in the wireType.</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	<span class="comment">// The parentheses look odd but are correct.</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>	if (ut.externalDec == xGob) != (ok &amp;&amp; wire.GobEncoderT != nil) ||
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>		(ut.externalDec == xBinary) != (ok &amp;&amp; wire.BinaryMarshalerT != nil) ||
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>		(ut.externalDec == xText) != (ok &amp;&amp; wire.TextMarshalerT != nil) {
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>		return false
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>	}
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	if ut.externalDec != 0 { <span class="comment">// This test trumps all others.</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		return true
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	}
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	switch t := ut.base; t.Kind() {
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	default:
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>		<span class="comment">// chan, etc: cannot handle.</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>		return false
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		return fw == tBool
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		return fw == tInt
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		return fw == tUint
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	case reflect.Float32, reflect.Float64:
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		return fw == tFloat
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	case reflect.Complex64, reflect.Complex128:
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		return fw == tComplex
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	case reflect.String:
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		return fw == tString
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	case reflect.Interface:
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>		return fw == tInterface
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	case reflect.Array:
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>		if !ok || wire.ArrayT == nil {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			return false
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		}
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>		array := wire.ArrayT
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		return t.Len() == array.Len &amp;&amp; dec.compatibleType(t.Elem(), array.Elem, inProgress)
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>	case reflect.Map:
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		if !ok || wire.MapT == nil {
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>			return false
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		}
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>		MapType := wire.MapT
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>		return dec.compatibleType(t.Key(), MapType.Key, inProgress) &amp;&amp; dec.compatibleType(t.Elem(), MapType.Elem, inProgress)
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>		<span class="comment">// Is it an array of bytes?</span>
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		if t.Elem().Kind() == reflect.Uint8 {
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>			return fw == tBytes
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		}
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>		<span class="comment">// Extract and compare element types.</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		var sw *sliceType
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		if tt := builtinIdToType(fw); tt != nil {
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>			sw, _ = tt.(*sliceType)
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>		} else if wire != nil {
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>			sw = wire.SliceT
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>		}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>		elem := userType(t.Elem()).base
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>		return sw != nil &amp;&amp; dec.compatibleType(elem, sw.Elem, inProgress)
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>	case reflect.Struct:
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>		return true
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>	}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>}
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span><span class="comment">// typeString returns a human-readable description of the type identified by remoteId.</span>
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>func (dec *Decoder) typeString(remoteId typeId) string {
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>	typeLock.Lock()
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>	defer typeLock.Unlock()
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>	if t := idToType(remoteId); t != nil {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		<span class="comment">// globally known type.</span>
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>		return t.string()
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>	}
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	return dec.wireType[remoteId].string()
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>}
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span><span class="comment">// compileSingle compiles the decoder engine for a non-struct top-level value, including</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span><span class="comment">// GobDecoders.</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>func (dec *Decoder) compileSingle(remoteId typeId, ut *userTypeInfo) (engine *decEngine, err error) {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	rt := ut.user
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>	engine = new(decEngine)
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>	engine.instr = make([]decInstr, 1) <span class="comment">// one item</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>	name := rt.String()                <span class="comment">// best we can do</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>	if !dec.compatibleType(rt, remoteId, make(map[reflect.Type]typeId)) {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>		remoteType := dec.typeString(remoteId)
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>		<span class="comment">// Common confusing case: local interface type, remote concrete type.</span>
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>		if ut.base.Kind() == reflect.Interface &amp;&amp; remoteId != tInterface {
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>			return nil, errors.New(&#34;gob: local interface type &#34; + name + &#34; can only be decoded from remote interface type; received concrete type &#34; + remoteType)
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>		}
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		return nil, errors.New(&#34;gob: decoding into local type &#34; + name + &#34;, received remote type &#34; + remoteType)
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	}
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>	op := dec.decOpFor(remoteId, rt, name, make(map[reflect.Type]*decOp))
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	ovfl := errors.New(`value for &#34;` + name + `&#34; out of range`)
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	engine.instr[singletonField] = decInstr{*op, singletonField, nil, ovfl}
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	engine.numInstr = 1
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>	return
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>}
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span><span class="comment">// compileIgnoreSingle compiles the decoder engine for a non-struct top-level value that will be discarded.</span>
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>func (dec *Decoder) compileIgnoreSingle(remoteId typeId) *decEngine {
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>	engine := new(decEngine)
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>	engine.instr = make([]decInstr, 1) <span class="comment">// one item</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>	op := dec.decIgnoreOpFor(remoteId, make(map[typeId]*decOp))
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>	ovfl := overflow(dec.typeString(remoteId))
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	engine.instr[0] = decInstr{*op, 0, nil, ovfl}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>	engine.numInstr = 1
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>	return engine
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>}
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span><span class="comment">// compileDec compiles the decoder engine for a value. If the value is not a struct,</span>
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span><span class="comment">// it calls out to compileSingle.</span>
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>func (dec *Decoder) compileDec(remoteId typeId, ut *userTypeInfo) (engine *decEngine, err error) {
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	defer catchError(&amp;err)
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	rt := ut.base
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>	srt := rt
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>	if srt.Kind() != reflect.Struct || ut.externalDec != 0 {
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>		return dec.compileSingle(remoteId, ut)
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>	}
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	var wireStruct *structType
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>	<span class="comment">// Builtin types can come from global pool; the rest must be defined by the decoder.</span>
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>	<span class="comment">// Also we know we&#39;re decoding a struct now, so the client must have sent one.</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>	if t := builtinIdToType(remoteId); t != nil {
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>		wireStruct, _ = t.(*structType)
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	} else {
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>		wire := dec.wireType[remoteId]
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>		if wire == nil {
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			error_(errBadType)
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>		}
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>		wireStruct = wire.StructT
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	}
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>	if wireStruct == nil {
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>		errorf(&#34;type mismatch in decoder: want struct type %s; got non-struct&#34;, rt)
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>	}
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	engine = new(decEngine)
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	engine.instr = make([]decInstr, len(wireStruct.Field))
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>	seen := make(map[reflect.Type]*decOp)
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>	<span class="comment">// Loop over the fields of the wire type.</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>	for fieldnum := 0; fieldnum &lt; len(wireStruct.Field); fieldnum++ {
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>		wireField := wireStruct.Field[fieldnum]
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>		if wireField.Name == &#34;&#34; {
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>			errorf(&#34;empty name for remote field of type %s&#34;, wireStruct.Name)
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>		}
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>		ovfl := overflow(wireField.Name)
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>		<span class="comment">// Find the field of the local type with the same name.</span>
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>		localField, present := srt.FieldByName(wireField.Name)
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>		<span class="comment">// TODO(r): anonymous names</span>
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>		if !present || !isExported(wireField.Name) {
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>			op := dec.decIgnoreOpFor(wireField.Id, make(map[typeId]*decOp))
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>			engine.instr[fieldnum] = decInstr{*op, fieldnum, nil, ovfl}
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>			continue
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		}
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>		if !dec.compatibleType(localField.Type, wireField.Id, make(map[reflect.Type]typeId)) {
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>			errorf(&#34;wrong type (%s) for received field %s.%s&#34;, localField.Type, wireStruct.Name, wireField.Name)
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>		}
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>		op := dec.decOpFor(wireField.Id, localField.Type, localField.Name, seen)
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>		engine.instr[fieldnum] = decInstr{*op, fieldnum, localField.Index, ovfl}
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>		engine.numInstr++
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>	}
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>	return
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>}
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span><span class="comment">// getDecEnginePtr returns the engine for the specified type.</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>func (dec *Decoder) getDecEnginePtr(remoteId typeId, ut *userTypeInfo) (enginePtr **decEngine, err error) {
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>	rt := ut.user
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>	decoderMap, ok := dec.decoderCache[rt]
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>	if !ok {
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>		decoderMap = make(map[typeId]**decEngine)
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>		dec.decoderCache[rt] = decoderMap
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>	}
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>	if enginePtr, ok = decoderMap[remoteId]; !ok {
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>		<span class="comment">// To handle recursive types, mark this engine as underway before compiling.</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		enginePtr = new(*decEngine)
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		decoderMap[remoteId] = enginePtr
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>		*enginePtr, err = dec.compileDec(remoteId, ut)
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>		if err != nil {
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			delete(decoderMap, remoteId)
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>		}
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>	}
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>	return
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span><span class="comment">// emptyStruct is the type we compile into when ignoring a struct value.</span>
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>type emptyStruct struct{}
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>var emptyStructType = reflect.TypeFor[emptyStruct]()
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span><span class="comment">// getIgnoreEnginePtr returns the engine for the specified type when the value is to be discarded.</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>func (dec *Decoder) getIgnoreEnginePtr(wireId typeId) (enginePtr **decEngine, err error) {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>	var ok bool
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>	if enginePtr, ok = dec.ignorerCache[wireId]; !ok {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>		<span class="comment">// To handle recursive types, mark this engine as underway before compiling.</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>		enginePtr = new(*decEngine)
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		dec.ignorerCache[wireId] = enginePtr
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		wire := dec.wireType[wireId]
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>		if wire != nil &amp;&amp; wire.StructT != nil {
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>			*enginePtr, err = dec.compileDec(wireId, userType(emptyStructType))
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		} else {
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>			*enginePtr = dec.compileIgnoreSingle(wireId)
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>		}
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>		if err != nil {
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			delete(dec.ignorerCache, wireId)
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>		}
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>	}
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>	return
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>}
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span><span class="comment">// decodeValue decodes the data stream representing a value and stores it in value.</span>
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>func (dec *Decoder) decodeValue(wireId typeId, value reflect.Value) {
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>	defer catchError(&amp;dec.err)
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>	<span class="comment">// If the value is nil, it means we should just ignore this item.</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>	if !value.IsValid() {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>		dec.decodeIgnoredValue(wireId)
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>		return
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	}
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>	<span class="comment">// Dereference down to the underlying type.</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	ut := userType(value.Type())
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>	base := ut.base
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	var enginePtr **decEngine
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>	enginePtr, dec.err = dec.getDecEnginePtr(wireId, ut)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>	if dec.err != nil {
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>		return
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>	}
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>	value = decAlloc(value)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>	engine := *enginePtr
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>	if st := base; st.Kind() == reflect.Struct &amp;&amp; ut.externalDec == 0 {
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>		wt := dec.wireType[wireId]
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>		if engine.numInstr == 0 &amp;&amp; st.NumField() &gt; 0 &amp;&amp;
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>			wt != nil &amp;&amp; len(wt.StructT.Field) &gt; 0 {
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>			name := base.Name()
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>			errorf(&#34;type mismatch: no fields matched compiling decoder for %s&#34;, name)
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>		}
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		dec.decodeStruct(engine, value)
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>	} else {
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>		dec.decodeSingle(engine, value)
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>	}
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>}
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span><span class="comment">// decodeIgnoredValue decodes the data stream representing a value of the specified type and discards it.</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>func (dec *Decoder) decodeIgnoredValue(wireId typeId) {
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	var enginePtr **decEngine
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>	enginePtr, dec.err = dec.getIgnoreEnginePtr(wireId)
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>	if dec.err != nil {
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>		return
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>	}
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>	wire := dec.wireType[wireId]
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>	if wire != nil &amp;&amp; wire.StructT != nil {
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>		dec.ignoreStruct(*enginePtr)
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>	} else {
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>		dec.ignoreSingle(*enginePtr)
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>	}
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>}
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>const (
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>	intBits     = 32 &lt;&lt; (^uint(0) &gt;&gt; 63)
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	uintptrBits = 32 &lt;&lt; (^uintptr(0) &gt;&gt; 63)
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>)
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>func init() {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>	var iop, uop decOp
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>	switch intBits {
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	case 32:
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>		iop = decInt32
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		uop = decUint32
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	case 64:
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>		iop = decInt64
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>		uop = decUint64
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	default:
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>		panic(&#34;gob: unknown size of int/uint&#34;)
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>	}
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	decOpTable[reflect.Int] = iop
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	decOpTable[reflect.Uint] = uop
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>	<span class="comment">// Finally uintptr</span>
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>	switch uintptrBits {
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>	case 32:
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		uop = decUint32
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>	case 64:
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>		uop = decUint64
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>	default:
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>		panic(&#34;gob: unknown size of uintptr&#34;)
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	}
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>	decOpTable[reflect.Uintptr] = uop
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>}
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span><span class="comment">// Gob depends on being able to take the address</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span><span class="comment">// of zeroed Values it creates, so use this wrapper instead</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span><span class="comment">// of the standard reflect.Zero.</span>
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span><span class="comment">// Each call allocates once.</span>
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>func allocValue(t reflect.Type) reflect.Value {
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	return reflect.New(t).Elem()
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>
</pre><p><a href="decode.go?m=text">View as plain text</a></p>

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
