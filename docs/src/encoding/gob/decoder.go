<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/gob/decoder.go - Go Documentation Server</title>

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
<a href="decoder.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/gob">gob</a>/<span class="text-muted">decoder.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package gob
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/saferio&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// tooBig provides a sanity check for sizes; used in several places. Upper limit</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// of is 1GB on 32-bit systems, 8GB on 64-bit, allowing room to grow a little</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// without overflow.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>const tooBig = (1 &lt;&lt; 30) &lt;&lt; (^uint(0) &gt;&gt; 62)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// A Decoder manages the receipt of type and data information read from the</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// remote side of a connection.  It is safe for concurrent use by multiple</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// goroutines.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// The Decoder does only basic sanity checking on decoded input sizes,</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// and its limits are not configurable. Take caution when decoding gob data</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// from untrusted sources.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>type Decoder struct {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	mutex        sync.Mutex                              <span class="comment">// each item must be received atomically</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	r            io.Reader                               <span class="comment">// source of the data</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	buf          decBuffer                               <span class="comment">// buffer for more efficient i/o from r</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	wireType     map[typeId]*wireType                    <span class="comment">// map from remote ID to local description</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	decoderCache map[reflect.Type]map[typeId]**decEngine <span class="comment">// cache of compiled engines</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	ignorerCache map[typeId]**decEngine                  <span class="comment">// ditto for ignored objects</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	freeList     *decoderState                           <span class="comment">// list of free decoderStates; avoids reallocation</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	countBuf     []byte                                  <span class="comment">// used for decoding integers while parsing messages</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	err          error
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// ignoreDepth tracks the depth of recursively parsed ignored fields</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	ignoreDepth int
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// NewDecoder returns a new decoder that reads from the [io.Reader].</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// If r does not also implement [io.ByteReader], it will be wrapped in a</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// [bufio.Reader].</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>func NewDecoder(r io.Reader) *Decoder {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	dec := new(Decoder)
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// We use the ability to read bytes as a plausible surrogate for buffering.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	if _, ok := r.(io.ByteReader); !ok {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		r = bufio.NewReader(r)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	dec.r = r
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	dec.wireType = make(map[typeId]*wireType)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	dec.decoderCache = make(map[reflect.Type]map[typeId]**decEngine)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	dec.ignorerCache = make(map[typeId]**decEngine)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	dec.countBuf = make([]byte, 9) <span class="comment">// counts may be uint64s (unlikely!), require 9 bytes</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	return dec
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// recvType loads the definition of a type.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>func (dec *Decoder) recvType(id typeId) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// Have we already seen this type? That&#39;s an error</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if id &lt; firstUserId || dec.wireType[id] != nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		dec.err = errors.New(&#34;gob: duplicate type received&#34;)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		return
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// Type:</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	wire := new(wireType)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	dec.decodeValue(tWireType, reflect.ValueOf(wire))
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if dec.err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// Remember we&#39;ve seen this type.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	dec.wireType[id] = wire
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>var errBadCount = errors.New(&#34;invalid message length&#34;)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// recvMessage reads the next count-delimited item from the input. It is the converse</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// of Encoder.writeMessage. It returns false on EOF or other error reading the message.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>func (dec *Decoder) recvMessage() bool {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// Read a count.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	nbytes, _, err := decodeUintReader(dec.r, dec.countBuf)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	if err != nil {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		dec.err = err
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return false
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if nbytes &gt;= tooBig {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		dec.err = errBadCount
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		return false
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	dec.readMessage(int(nbytes))
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	return dec.err == nil
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// readMessage reads the next nbytes bytes from the input.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func (dec *Decoder) readMessage(nbytes int) {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	if dec.buf.Len() != 0 {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		<span class="comment">// The buffer should always be empty now.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		panic(&#34;non-empty decoder buffer&#34;)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Read the data</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	var buf []byte
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	buf, dec.err = saferio.ReadData(dec.r, uint64(nbytes))
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	dec.buf.SetBytes(buf)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if dec.err == io.EOF {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		dec.err = io.ErrUnexpectedEOF
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// toInt turns an encoded uint64 into an int, according to the marshaling rules.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>func toInt(x uint64) int64 {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	i := int64(x &gt;&gt; 1)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if x&amp;1 != 0 {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		i = ^i
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return i
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>func (dec *Decoder) nextInt() int64 {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	n, _, err := decodeUintReader(&amp;dec.buf, dec.countBuf)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if err != nil {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		dec.err = err
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return toInt(n)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (dec *Decoder) nextUint() uint64 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	n, _, err := decodeUintReader(&amp;dec.buf, dec.countBuf)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	if err != nil {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		dec.err = err
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	return n
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// decodeTypeSequence parses:</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// TypeSequence</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//	(TypeDefinition DelimitedTypeDefinition*)?</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// and returns the type id of the next value. It returns -1 at</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// EOF.  Upon return, the remainder of dec.buf is the value to be</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// decoded. If this is an interface value, it can be ignored by</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// resetting that buffer.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>func (dec *Decoder) decodeTypeSequence(isInterface bool) typeId {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	firstMessage := true
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	for dec.err == nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		if dec.buf.Len() == 0 {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			if !dec.recvMessage() {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				<span class="comment">// We can only return io.EOF if the input was empty.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				<span class="comment">// If we read one or more type spec messages,</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>				<span class="comment">// require a data item message to follow.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>				<span class="comment">// If we hit an EOF before that, then give ErrUnexpectedEOF.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>				if !firstMessage &amp;&amp; dec.err == io.EOF {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>					dec.err = io.ErrUnexpectedEOF
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>				}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>				break
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		<span class="comment">// Receive a type id.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		id := typeId(dec.nextInt())
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		if id &gt;= 0 {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			<span class="comment">// Value follows.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			return id
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		<span class="comment">// Type definition for (-id) follows.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		dec.recvType(-id)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if dec.err != nil {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			break
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		<span class="comment">// When decoding an interface, after a type there may be a</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		<span class="comment">// DelimitedValue still in the buffer. Skip its count.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		<span class="comment">// (Alternatively, the buffer is empty and the byte count</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		<span class="comment">// will be absorbed by recvMessage.)</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		if dec.buf.Len() &gt; 0 {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			if !isInterface {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				dec.err = errors.New(&#34;extra data in buffer&#34;)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				break
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			dec.nextUint()
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		firstMessage = false
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	return -1
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// Decode reads the next value from the input stream and stores</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// it in the data represented by the empty interface value.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// If e is nil, the value will be discarded. Otherwise,</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span><span class="comment">// the value underlying e must be a pointer to the</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// correct type for the next data item received.</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// If the input is at EOF, Decode returns [io.EOF] and</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span><span class="comment">// does not modify e.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>func (dec *Decoder) Decode(e any) error {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if e == nil {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		return dec.DecodeValue(reflect.Value{})
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	value := reflect.ValueOf(e)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// If e represents a value as opposed to a pointer, the answer won&#39;t</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// get back to the caller. Make sure it&#39;s a pointer.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	if value.Type().Kind() != reflect.Pointer {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		dec.err = errors.New(&#34;gob: attempt to decode into a non-pointer&#34;)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		return dec.err
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return dec.DecodeValue(value)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// DecodeValue reads the next value from the input stream.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// If v is the zero reflect.Value (v.Kind() == Invalid), DecodeValue discards the value.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// Otherwise, it stores the value into v. In that case, v must represent</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// a non-nil pointer to data or be an assignable reflect.Value (v.CanSet())</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// If the input is at EOF, DecodeValue returns [io.EOF] and</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// does not modify v.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>func (dec *Decoder) DecodeValue(v reflect.Value) error {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if v.IsValid() {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if v.Kind() == reflect.Pointer &amp;&amp; !v.IsNil() {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			<span class="comment">// That&#39;s okay, we&#39;ll store through the pointer.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		} else if !v.CanSet() {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			return errors.New(&#34;gob: DecodeValue of unassignable value&#34;)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// Make sure we&#39;re single-threaded through here.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	dec.mutex.Lock()
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	defer dec.mutex.Unlock()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	dec.buf.Reset() <span class="comment">// In case data lingers from previous invocation.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	dec.err = nil
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	id := dec.decodeTypeSequence(false)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	if dec.err == nil {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		dec.decodeValue(id, v)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	return dec.err
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// If debug.go is compiled into the program, debugFunc prints a human-readable</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">// representation of the gob data read from r by calling that file&#39;s Debug function.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span><span class="comment">// Otherwise it is nil.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>var debugFunc func(io.Reader)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
</pre><p><a href="decoder.go?m=text">View as plain text</a></p>

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
