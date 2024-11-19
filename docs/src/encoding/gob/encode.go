<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/gob/encode.go - Go Documentation Server</title>

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
<a href="encode.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/gob">gob</a>/<span class="text-muted">encode.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:generate go run encgen.go -output enc_helpers.go</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package gob
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;encoding&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;math/bits&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>const uint64Size = 8
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>type encHelper func(state *encoderState, v reflect.Value) bool
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// encoderState is the global execution state of an instance of the encoder.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Field numbers are delta encoded and always increase. The field</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// number is initialized to -1 so 0 comes out as delta(1). A delta of</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// 0 terminates the structure.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type encoderState struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	enc      *Encoder
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	b        *encBuffer
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	sendZero bool                 <span class="comment">// encoding an array element or map key/value pair; send zero values</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	fieldnum int                  <span class="comment">// the last field number written.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	buf      [1 + uint64Size]byte <span class="comment">// buffer used by the encoder; here to avoid allocation.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	next     *encoderState        <span class="comment">// for free list</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// encBuffer is an extremely simple, fast implementation of a write-only byte buffer.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// It never returns a non-nil error, but Write returns an error value so it matches io.Writer.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>type encBuffer struct {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	data    []byte
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	scratch [64]byte
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>var encBufferPool = sync.Pool{
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	New: func() any {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		e := new(encBuffer)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		e.data = e.scratch[0:0]
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		return e
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	},
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func (e *encBuffer) writeByte(c byte) {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	e.data = append(e.data, c)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (e *encBuffer) Write(p []byte) (int, error) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	e.data = append(e.data, p...)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	return len(p), nil
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (e *encBuffer) WriteString(s string) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	e.data = append(e.data, s...)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>func (e *encBuffer) Len() int {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	return len(e.data)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>func (e *encBuffer) Bytes() []byte {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	return e.data
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func (e *encBuffer) Reset() {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if len(e.data) &gt;= tooBig {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		e.data = e.scratch[0:0]
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	} else {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		e.data = e.data[0:0]
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (enc *Encoder) newEncoderState(b *encBuffer) *encoderState {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	e := enc.freeList
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if e == nil {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		e = new(encoderState)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		e.enc = enc
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	} else {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		enc.freeList = e.next
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	e.sendZero = false
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	e.fieldnum = 0
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	e.b = b
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	if len(b.data) == 0 {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		b.data = b.scratch[0:0]
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return e
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (enc *Encoder) freeEncoderState(e *encoderState) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	e.next = enc.freeList
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	enc.freeList = e
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// Unsigned integers have a two-state encoding. If the number is less</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// than 128 (0 through 0x7F), its value is written directly.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// Otherwise the value is written in big-endian byte order preceded</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// by the byte length, negated.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// encodeUint writes an encoded unsigned integer to state.b.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func (state *encoderState) encodeUint(x uint64) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if x &lt;= 0x7F {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		state.b.writeByte(uint8(x))
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		return
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	binary.BigEndian.PutUint64(state.buf[1:], x)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	bc := bits.LeadingZeros64(x) &gt;&gt; 3      <span class="comment">// 8 - bytelen(x)</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	state.buf[bc] = uint8(bc - uint64Size) <span class="comment">// and then we subtract 8 to get -bytelen(x)</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	state.b.Write(state.buf[bc : uint64Size+1])
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// encodeInt writes an encoded signed integer to state.w.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// The low bit of the encoding says whether to bit complement the (other bits of the)</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// uint to recover the int.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func (state *encoderState) encodeInt(i int64) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	var x uint64
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		x = uint64(^i&lt;&lt;1) | 1
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	} else {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		x = uint64(i &lt;&lt; 1)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	state.encodeUint(x)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// encOp is the signature of an encoding operator for a given type.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>type encOp func(i *encInstr, state *encoderState, v reflect.Value)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// The &#39;instructions&#39; of the encoding machine</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>type encInstr struct {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	op    encOp
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	field int   <span class="comment">// field number in input</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	index []int <span class="comment">// struct index</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	indir int   <span class="comment">// how many pointer indirections to reach the value in the struct</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// update emits a field number and updates the state to record its value for delta encoding.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// If the instruction pointer is nil, it does nothing</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func (state *encoderState) update(instr *encInstr) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if instr != nil {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		state.encodeUint(uint64(instr.field - state.fieldnum))
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		state.fieldnum = instr.field
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// Each encoder for a composite is responsible for handling any</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// indirections associated with the elements of the data structure.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// If any pointer so reached is nil, no bytes are written. If the</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// data item is zero, no bytes are written. Single values - ints,</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// strings etc. - are indirected before calling their encoders.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// Otherwise, the output (for a scalar) is the field number, as an</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// encoded integer, followed by the field data in its appropriate</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// format.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// encIndirect dereferences pv indir times and returns the result.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>func encIndirect(pv reflect.Value, indir int) reflect.Value {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	for ; indir &gt; 0; indir-- {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		if pv.IsNil() {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			break
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		pv = pv.Elem()
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	return pv
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// encBool encodes the bool referenced by v as an unsigned 0 or 1.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func encBool(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	b := v.Bool()
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	if b || state.sendZero {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		state.update(i)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if b {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			state.encodeUint(1)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		} else {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			state.encodeUint(0)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// encInt encodes the signed integer (int int8 int16 int32 int64) referenced by v.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>func encInt(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	value := v.Int()
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	if value != 0 || state.sendZero {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		state.update(i)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		state.encodeInt(value)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">// encUint encodes the unsigned integer (uint uint8 uint16 uint32 uint64 uintptr) referenced by v.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func encUint(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	value := v.Uint()
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if value != 0 || state.sendZero {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		state.update(i)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		state.encodeUint(value)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// floatBits returns a uint64 holding the bits of a floating-point number.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// Floating-point numbers are transmitted as uint64s holding the bits</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// of the underlying representation. They are sent byte-reversed, with</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// the exponent end coming out first, so integer floating point numbers</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// (for example) transmit more compactly. This routine does the</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// swizzling.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func floatBits(f float64) uint64 {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	u := math.Float64bits(f)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	return bits.ReverseBytes64(u)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// encFloat encodes the floating point value (float32 float64) referenced by v.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>func encFloat(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	f := v.Float()
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	if f != 0 || state.sendZero {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		bits := floatBits(f)
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		state.update(i)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		state.encodeUint(bits)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// encComplex encodes the complex value (complex64 complex128) referenced by v.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// Complex numbers are just a pair of floating-point numbers, real part first.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>func encComplex(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	c := v.Complex()
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	if c != 0+0i || state.sendZero {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		rpart := floatBits(real(c))
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		ipart := floatBits(imag(c))
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		state.update(i)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		state.encodeUint(rpart)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		state.encodeUint(ipart)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// encUint8Array encodes the byte array referenced by v.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span><span class="comment">// Byte arrays are encoded as an unsigned count followed by the raw bytes.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>func encUint8Array(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	b := v.Bytes()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if len(b) &gt; 0 || state.sendZero {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		state.update(i)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		state.encodeUint(uint64(len(b)))
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		state.b.Write(b)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span><span class="comment">// encString encodes the string referenced by v.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// Strings are encoded as an unsigned count followed by the raw bytes.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>func encString(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	s := v.String()
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	if len(s) &gt; 0 || state.sendZero {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		state.update(i)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		state.encodeUint(uint64(len(s)))
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		state.b.WriteString(s)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// encStructTerminator encodes the end of an encoded struct</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// as delta field number of 0.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>func encStructTerminator(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	state.encodeUint(0)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// Execution engine</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// encEngine an array of instructions indexed by field number of the encoding</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// data, typically a struct. It is executed top to bottom, walking the struct.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>type encEngine struct {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	instr []encInstr
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>const singletonField = 0
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// valid reports whether the value is valid and a non-nil pointer.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// (Slices, maps, and chans take care of themselves.)</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>func valid(v reflect.Value) bool {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	switch v.Kind() {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	case reflect.Invalid:
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		return false
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	case reflect.Pointer:
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		return !v.IsNil()
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	return true
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// encodeSingle encodes a single top-level non-struct value.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>func (enc *Encoder) encodeSingle(b *encBuffer, engine *encEngine, value reflect.Value) {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	defer enc.freeEncoderState(state)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	state.fieldnum = singletonField
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// There is no surrounding struct to frame the transmission, so we must</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// generate data even if the item is zero. To do this, set sendZero.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	state.sendZero = true
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	instr := &amp;engine.instr[singletonField]
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	if instr.indir &gt; 0 {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		value = encIndirect(value, instr.indir)
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	if valid(value) {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		instr.op(instr, state, value)
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// encodeStruct encodes a single struct value.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (enc *Encoder) encodeStruct(b *encBuffer, engine *encEngine, value reflect.Value) {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if !valid(value) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		return
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	defer enc.freeEncoderState(state)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	for i := 0; i &lt; len(engine.instr); i++ {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		instr := &amp;engine.instr[i]
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		if i &gt;= value.NumField() {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			<span class="comment">// encStructTerminator</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			instr.op(instr, state, reflect.Value{})
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			break
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		field := value.FieldByIndex(instr.index)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		if instr.indir &gt; 0 {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			field = encIndirect(field, instr.indir)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			<span class="comment">// TODO: Is field guaranteed valid? If so we could avoid this check.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			if !valid(field) {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				continue
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		instr.op(instr, state, field)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span><span class="comment">// encodeArray encodes an array.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>func (enc *Encoder) encodeArray(b *encBuffer, value reflect.Value, op encOp, elemIndir int, length int, helper encHelper) {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	defer enc.freeEncoderState(state)
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	state.sendZero = true
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	state.encodeUint(uint64(length))
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if helper != nil &amp;&amp; helper(state, value) {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		return
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	for i := 0; i &lt; length; i++ {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		elem := value.Index(i)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		if elemIndir &gt; 0 {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			elem = encIndirect(elem, elemIndir)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			<span class="comment">// TODO: Is elem guaranteed valid? If so we could avoid this check.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			if !valid(elem) {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>				errorf(&#34;encodeArray: nil element&#34;)
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		op(nil, state, elem)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// encodeReflectValue is a helper for maps. It encodes the value v.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>func encodeReflectValue(state *encoderState, v reflect.Value, op encOp, indir int) {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	for i := 0; i &lt; indir &amp;&amp; v.IsValid(); i++ {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		v = reflect.Indirect(v)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	if !v.IsValid() {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		errorf(&#34;encodeReflectValue: nil element&#34;)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	op(nil, state, v)
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// encodeMap encodes a map as unsigned count followed by key:value pairs.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>func (enc *Encoder) encodeMap(b *encBuffer, mv reflect.Value, keyOp, elemOp encOp, keyIndir, elemIndir int) {
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	state.sendZero = true
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	state.encodeUint(uint64(mv.Len()))
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	mi := mv.MapRange()
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	for mi.Next() {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		encodeReflectValue(state, mi.Key(), keyOp, keyIndir)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		encodeReflectValue(state, mi.Value(), elemOp, elemIndir)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	enc.freeEncoderState(state)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// encodeInterface encodes the interface value iv.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span><span class="comment">// To send an interface, we send a string identifying the concrete type, followed</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// by the type identifier (which might require defining that type right now), followed</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// by the concrete value. A nil value gets sent as the empty string for the name,</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span><span class="comment">// followed by no value.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>func (enc *Encoder) encodeInterface(b *encBuffer, iv reflect.Value) {
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// Gobs can encode nil interface values but not typed interface</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	<span class="comment">// values holding nil pointers, since nil pointers point to no value.</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	elem := iv.Elem()
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	if elem.Kind() == reflect.Pointer &amp;&amp; elem.IsNil() {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		errorf(&#34;gob: cannot encode nil pointer of type %s inside interface&#34;, iv.Elem().Type())
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	state.sendZero = true
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	if iv.IsNil() {
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		state.encodeUint(0)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		return
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	ut := userType(iv.Elem().Type())
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	namei, ok := concreteTypeToName.Load(ut.base)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	if !ok {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		errorf(&#34;type not registered for interface: %s&#34;, ut.base)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	name := namei.(string)
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	<span class="comment">// Send the name.</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	state.encodeUint(uint64(len(name)))
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	state.b.WriteString(name)
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	<span class="comment">// Define the type id if necessary.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	enc.sendTypeDescriptor(enc.writer(), state, ut)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// Send the type id.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	enc.sendTypeId(state, ut)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	<span class="comment">// Encode the value into a new buffer. Any nested type definitions</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	<span class="comment">// should be written to b, before the encoded value.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	enc.pushWriter(b)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	data := encBufferPool.Get().(*encBuffer)
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	data.Write(spaceForLength)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	enc.encode(data, elem, ut)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if enc.err != nil {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		error_(enc.err)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	enc.popWriter()
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	enc.writeMessage(b, data)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	data.Reset()
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	encBufferPool.Put(data)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	if enc.err != nil {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		error_(enc.err)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	enc.freeEncoderState(state)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span><span class="comment">// encodeGobEncoder encodes a value that implements the GobEncoder interface.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span><span class="comment">// The data is sent as a byte array.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>func (enc *Encoder) encodeGobEncoder(b *encBuffer, ut *userTypeInfo, v reflect.Value) {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// TODO: should we catch panics from the called method?</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	var data []byte
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	var err error
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// We know it&#39;s one of these.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	switch ut.externalEnc {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	case xGob:
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		data, err = v.Interface().(GobEncoder).GobEncode()
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	case xBinary:
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		data, err = v.Interface().(encoding.BinaryMarshaler).MarshalBinary()
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	case xText:
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		data, err = v.Interface().(encoding.TextMarshaler).MarshalText()
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	}
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	if err != nil {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		error_(err)
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	state := enc.newEncoderState(b)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	state.fieldnum = -1
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	state.encodeUint(uint64(len(data)))
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	state.b.Write(data)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	enc.freeEncoderState(state)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>var encOpTable = [...]encOp{
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	reflect.Bool:       encBool,
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	reflect.Int:        encInt,
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	reflect.Int8:       encInt,
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	reflect.Int16:      encInt,
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	reflect.Int32:      encInt,
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	reflect.Int64:      encInt,
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	reflect.Uint:       encUint,
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	reflect.Uint8:      encUint,
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	reflect.Uint16:     encUint,
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	reflect.Uint32:     encUint,
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	reflect.Uint64:     encUint,
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	reflect.Uintptr:    encUint,
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	reflect.Float32:    encFloat,
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	reflect.Float64:    encFloat,
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	reflect.Complex64:  encComplex,
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	reflect.Complex128: encComplex,
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	reflect.String:     encString,
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// encOpFor returns (a pointer to) the encoding op for the base type under rt and</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">// the indirection count to reach it.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>func encOpFor(rt reflect.Type, inProgress map[reflect.Type]*encOp, building map[*typeInfo]bool) (*encOp, int) {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	ut := userType(rt)
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	<span class="comment">// If the type implements GobEncoder, we handle it without further processing.</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	if ut.externalEnc != 0 {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		return gobEncodeOpFor(ut)
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	<span class="comment">// If this type is already in progress, it&#39;s a recursive type (e.g. map[string]*T).</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	<span class="comment">// Return the pointer to the op we&#39;re already building.</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	if opPtr := inProgress[rt]; opPtr != nil {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		return opPtr, ut.indir
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	typ := ut.base
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	indir := ut.indir
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	k := typ.Kind()
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	var op encOp
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	if int(k) &lt; len(encOpTable) {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		op = encOpTable[k]
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	if op == nil {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		inProgress[rt] = &amp;op
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		<span class="comment">// Special cases</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		switch t := typ; t.Kind() {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		case reflect.Slice:
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			if t.Elem().Kind() == reflect.Uint8 {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>				op = encUint8Array
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>				break
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			<span class="comment">// Slices have a header; we decode it to find the underlying array.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			elemOp, elemIndir := encOpFor(t.Elem(), inProgress, building)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			helper := encSliceHelper[t.Elem().Kind()]
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			op = func(i *encInstr, state *encoderState, slice reflect.Value) {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>				if !state.sendZero &amp;&amp; slice.Len() == 0 {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>					return
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>				}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>				state.update(i)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>				state.enc.encodeArray(state.b, slice, *elemOp, elemIndir, slice.Len(), helper)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		case reflect.Array:
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>			<span class="comment">// True arrays have size in the type.</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			elemOp, elemIndir := encOpFor(t.Elem(), inProgress, building)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>			helper := encArrayHelper[t.Elem().Kind()]
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>			op = func(i *encInstr, state *encoderState, array reflect.Value) {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>				state.update(i)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>				state.enc.encodeArray(state.b, array, *elemOp, elemIndir, array.Len(), helper)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>			}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		case reflect.Map:
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			keyOp, keyIndir := encOpFor(t.Key(), inProgress, building)
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			elemOp, elemIndir := encOpFor(t.Elem(), inProgress, building)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			op = func(i *encInstr, state *encoderState, mv reflect.Value) {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				<span class="comment">// We send zero-length (but non-nil) maps because the</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				<span class="comment">// receiver might want to use the map.  (Maps don&#39;t use append.)</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				if !state.sendZero &amp;&amp; mv.IsNil() {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>					return
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>				}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>				state.update(i)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>				state.enc.encodeMap(state.b, mv, *keyOp, *elemOp, keyIndir, elemIndir)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		case reflect.Struct:
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			<span class="comment">// Generate a closure that calls out to the engine for the nested type.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			getEncEngine(userType(typ), building)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			info := mustGetTypeInfo(typ)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			op = func(i *encInstr, state *encoderState, sv reflect.Value) {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>				state.update(i)
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>				<span class="comment">// indirect through info to delay evaluation for recursive structs</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				enc := info.encoder.Load()
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>				state.enc.encodeStruct(state.b, enc, sv)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		case reflect.Interface:
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			op = func(i *encInstr, state *encoderState, iv reflect.Value) {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>				if !state.sendZero &amp;&amp; (!iv.IsValid() || iv.IsNil()) {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>					return
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>				}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>				state.update(i)
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				state.enc.encodeInterface(state.b, iv)
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>			}
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	if op == nil {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		errorf(&#34;can&#39;t happen: encode type %s&#34;, rt)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	return &amp;op, indir
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span><span class="comment">// gobEncodeOpFor returns the op for a type that is known to implement GobEncoder.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>func gobEncodeOpFor(ut *userTypeInfo) (*encOp, int) {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	rt := ut.user
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	if ut.encIndir == -1 {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		rt = reflect.PointerTo(rt)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	} else if ut.encIndir &gt; 0 {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		for i := int8(0); i &lt; ut.encIndir; i++ {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			rt = rt.Elem()
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	var op encOp
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	op = func(i *encInstr, state *encoderState, v reflect.Value) {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		if ut.encIndir == -1 {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			<span class="comment">// Need to climb up one level to turn value into pointer.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			if !v.CanAddr() {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>				errorf(&#34;unaddressable value of type %s&#34;, rt)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>			}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>			v = v.Addr()
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>		}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		if !state.sendZero &amp;&amp; v.IsZero() {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>			return
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		state.update(i)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		state.enc.encodeGobEncoder(state.b, ut, v)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	return &amp;op, int(ut.encIndir) <span class="comment">// encIndir: op will get called with p == address of receiver.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span><span class="comment">// compileEnc returns the engine to compile the type.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>func compileEnc(ut *userTypeInfo, building map[*typeInfo]bool) *encEngine {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	srt := ut.base
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	engine := new(encEngine)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	seen := make(map[reflect.Type]*encOp)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	rt := ut.base
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	if ut.externalEnc != 0 {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		rt = ut.user
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	if ut.externalEnc == 0 &amp;&amp; srt.Kind() == reflect.Struct {
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		for fieldNum, wireFieldNum := 0, 0; fieldNum &lt; srt.NumField(); fieldNum++ {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>			f := srt.Field(fieldNum)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			if !isSent(srt, &amp;f) {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>				continue
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			op, indir := encOpFor(f.Type, seen, building)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			engine.instr = append(engine.instr, encInstr{*op, wireFieldNum, f.Index, indir})
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>			wireFieldNum++
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		}
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		if srt.NumField() &gt; 0 &amp;&amp; len(engine.instr) == 0 {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			errorf(&#34;type %s has no exported fields&#34;, rt)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		engine.instr = append(engine.instr, encInstr{encStructTerminator, 0, nil, 0})
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	} else {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		engine.instr = make([]encInstr, 1)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		op, indir := encOpFor(rt, seen, building)
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		engine.instr[0] = encInstr{*op, singletonField, nil, indir}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	return engine
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>}
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>
<span id="L623" class="ln">   623&nbsp;&nbsp;</span><span class="comment">// getEncEngine returns the engine to compile the type.</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>func getEncEngine(ut *userTypeInfo, building map[*typeInfo]bool) *encEngine {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	info, err := getTypeInfo(ut)
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	if err != nil {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>		error_(err)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	enc := info.encoder.Load()
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	if enc == nil {
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>		enc = buildEncEngine(info, ut, building)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>	}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	return enc
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>}
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>func buildEncEngine(info *typeInfo, ut *userTypeInfo, building map[*typeInfo]bool) *encEngine {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	<span class="comment">// Check for recursive types.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	if building != nil &amp;&amp; building[info] {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		return nil
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>	}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	info.encInit.Lock()
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	defer info.encInit.Unlock()
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	enc := info.encoder.Load()
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	if enc == nil {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>		if building == nil {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			building = make(map[*typeInfo]bool)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		building[info] = true
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		enc = compileEnc(ut, building)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		info.encoder.Store(enc)
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	return enc
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>func (enc *Encoder) encode(b *encBuffer, value reflect.Value, ut *userTypeInfo) {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	defer catchError(&amp;enc.err)
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	engine := getEncEngine(ut, nil)
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	indir := ut.indir
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	if ut.externalEnc != 0 {
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		indir = int(ut.encIndir)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	for i := 0; i &lt; indir; i++ {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		value = reflect.Indirect(value)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	if ut.externalEnc == 0 &amp;&amp; value.Type().Kind() == reflect.Struct {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		enc.encodeStruct(b, engine, value)
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	} else {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		enc.encodeSingle(b, engine, value)
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
</pre><p><a href="encode.go?m=text">View as plain text</a></p>

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
