<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/gob/doc.go - Go Documentation Server</title>

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
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/gob">gob</a>/<span class="text-muted">doc.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package gob manages streams of gobs - binary values exchanged between an
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>[Encoder] (transmitter) and a [Decoder] (receiver). A typical use is transporting
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>arguments and results of remote procedure calls (RPCs) such as those provided by
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>[net/rpc].
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>The implementation compiles a custom codec for each data type in the stream and
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>is most efficient when a single [Encoder] is used to transmit a stream of values,
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>amortizing the cost of compilation.
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span># Basics
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>A stream of gobs is self-describing. Each data item in the stream is preceded by
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>a specification of its type, expressed in terms of a small set of predefined
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>types. Pointers are not transmitted, but the things they point to are
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>transmitted; that is, the values are flattened. Nil pointers are not permitted,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>as they have no value. Recursive types work fine, but
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>recursive values (data with cycles) are problematic. This may change.
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>To use gobs, create an [Encoder] and present it with a series of data items as
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>values or addresses that can be dereferenced to values. The [Encoder] makes sure
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>all type information is sent before it is needed. At the receive side, a
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>[Decoder] retrieves values from the encoded stream and unpacks them into local
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>variables.
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span># Types and Values
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>The source and destination values/types need not correspond exactly. For structs,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>fields (identified by name) that are in the source but absent from the receiving
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>variable will be ignored. Fields that are in the receiving variable but missing
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>from the transmitted type or value will be ignored in the destination. If a field
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>with the same name is present in both, their types must be compatible. Both the
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>receiver and transmitter will do all necessary indirection and dereferencing to
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>convert between gobs and actual Go values. For instance, a gob type that is
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>schematically,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	struct { A, B int }
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>can be sent from or received into any of these Go types:
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	struct { A, B int }	// the same
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	*struct { A, B int }	// extra indirection of the struct
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	struct { *A, **B int }	// extra indirection of the fields
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	struct { A, B int64 }	// different concrete value type; see below
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>It may also be received into any of these:
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	struct { A, B int }	// the same
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	struct { B, A int }	// ordering doesn&#39;t matter; matching is by name
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	struct { A, B, C int }	// extra field (C) ignored
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	struct { B int }	// missing field (A) ignored; data will be dropped
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	struct { B, C int }	// missing field (A) ignored; extra field (C) ignored.
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>Attempting to receive into these types will draw a decode error:
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	struct { A int; B uint }	// change of signedness for B
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	struct { A int; B float }	// change of type for B
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	struct { }			// no field names in common
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	struct { C, D int }		// no field names in common
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>Integers are transmitted two ways: arbitrary precision signed integers or
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>arbitrary precision unsigned integers. There is no int8, int16 etc.
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>discrimination in the gob format; there are only signed and unsigned integers. As
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>described below, the transmitter sends the value in a variable-length encoding;
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>the receiver accepts the value and stores it in the destination variable.
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>Floating-point numbers are always sent using IEEE-754 64-bit precision (see
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>below).
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>Signed integers may be received into any signed integer variable: int, int16, etc.;
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>unsigned integers may be received into any unsigned integer variable; and floating
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>point values may be received into any floating point variable. However,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>the destination variable must be able to represent the value or the decode
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>operation will fail.
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>Structs, arrays and slices are also supported. Structs encode and decode only
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>exported fields. Strings and arrays of bytes are supported with a special,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>efficient representation (see below). When a slice is decoded, if the existing
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>slice has capacity the slice will be extended in place; if not, a new array is
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>allocated. Regardless, the length of the resulting slice reports the number of
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>elements decoded.
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>In general, if allocation is required, the decoder will allocate memory. If not,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>it will update the destination variables with values read from the stream. It does
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>not initialize them first, so if the destination is a compound value such as a
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>map, struct, or slice, the decoded values will be merged elementwise into the
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>existing variables.
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>Functions and channels will not be sent in a gob. Attempting to encode such a value
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>at the top level will fail. A struct field of chan or func type is treated exactly
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>like an unexported field and is ignored.
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>Gob can encode a value of any type implementing the [GobEncoder] or
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>[encoding.BinaryMarshaler] interfaces by calling the corresponding method,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>in that order of preference.
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>Gob can decode a value of any type implementing the [GobDecoder] or
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>[encoding.BinaryUnmarshaler] interfaces by calling the corresponding method,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>again in that order of preference.
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span># Encoding Details
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>This section documents the encoding, details that are not important for most
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>users. Details are presented bottom-up.
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>An unsigned integer is sent one of two ways. If it is less than 128, it is sent
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>as a byte with that value. Otherwise it is sent as a minimal-length big-endian
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>(high byte first) byte stream holding the value, preceded by one byte holding the
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>byte count, negated. Thus 0 is transmitted as (00), 7 is transmitted as (07) and
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>256 is transmitted as (FE 01 00).
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>A boolean is encoded within an unsigned integer: 0 for false, 1 for true.
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>A signed integer, i, is encoded within an unsigned integer, u. Within u, bits 1
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>upward contain the value; bit 0 says whether they should be complemented upon
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>receipt. The encode algorithm looks like this:
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	var u uint
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		u = (^uint(i) &lt;&lt; 1) | 1 // complement i, bit 0 is 1
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	} else {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		u = (uint(i) &lt;&lt; 1) // do not complement i, bit 0 is 0
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	encodeUnsigned(u)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>The low bit is therefore analogous to a sign bit, but making it the complement bit
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>instead guarantees that the largest negative integer is not a special case. For
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>example, -129=^128=(^256&gt;&gt;1) encodes as (FE 01 01).
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>Floating-point numbers are always sent as a representation of a float64 value.
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>That value is converted to a uint64 using [math.Float64bits]. The uint64 is then
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>byte-reversed and sent as a regular unsigned integer. The byte-reversal means the
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>exponent and high-precision part of the mantissa go first. Since the low bits are
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>often zero, this can save encoding bytes. For instance, 17.0 is encoded in only
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>three bytes (FE 31 40).
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>Strings and slices of bytes are sent as an unsigned count followed by that many
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>uninterpreted bytes of the value.
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>All other slices and arrays are sent as an unsigned count followed by that many
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>elements using the standard gob encoding for their type, recursively.
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>Maps are sent as an unsigned count followed by that many key, element
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>pairs. Empty but non-nil maps are sent, so if the receiver has not allocated
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>one already, one will always be allocated on receipt unless the transmitted map
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>is nil and not at the top level.
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>In slices and arrays, as well as maps, all elements, even zero-valued elements,
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>are transmitted, even if all the elements are zero.
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>Structs are sent as a sequence of (field number, field value) pairs. The field
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>value is sent using the standard gob encoding for its type, recursively. If a
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>field has the zero value for its type (except for arrays; see above), it is omitted
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>from the transmission. The field number is defined by the type of the encoded
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>struct: the first field of the encoded type is field 0, the second is field 1,
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>etc. When encoding a value, the field numbers are delta encoded for efficiency
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>and the fields are always sent in order of increasing field number; the deltas are
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>therefore unsigned. The initialization for the delta encoding sets the field
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>number to -1, so an unsigned integer field 0 with value 7 is transmitted as unsigned
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>delta = 1, unsigned value = 7 or (01 07). Finally, after all the fields have been
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>sent a terminating mark denotes the end of the struct. That mark is a delta=0
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>value, which has representation (00).
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>Interface types are not checked for compatibility; all interface types are
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>treated, for transmission, as members of a single &#34;interface&#34; type, analogous to
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>int or []byte - in effect they&#39;re all treated as interface{}. Interface values
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>are transmitted as a string identifying the concrete type being sent (a name
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>that must be pre-defined by calling [Register]), followed by a byte count of the
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>length of the following data (so the value can be skipped if it cannot be
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>stored), followed by the usual encoding of concrete (dynamic) value stored in
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>the interface value. (A nil interface value is identified by the empty string
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>and transmits no value.) Upon receipt, the decoder verifies that the unpacked
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>concrete item satisfies the interface of the receiving variable.
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>If a value is passed to [Encoder.Encode] and the type is not a struct (or pointer to struct,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>etc.), for simplicity of processing it is represented as a struct of one field.
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>The only visible effect of this is to encode a zero byte after the value, just as
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>after the last field of an encoded struct, so that the decode algorithm knows when
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>the top-level value is complete.
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>The representation of types is described below. When a type is defined on a given
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>connection between an [Encoder] and [Decoder], it is assigned a signed integer type
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>id. When [Encoder.Encode](v) is called, it makes sure there is an id assigned for
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>the type of v and all its elements and then it sends the pair (typeid, encoded-v)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>where typeid is the type id of the encoded type of v and encoded-v is the gob
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>encoding of the value v.
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>To define a type, the encoder chooses an unused, positive type id and sends the
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>pair (-type id, encoded-type) where encoded-type is the gob encoding of a wireType
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>description, constructed from these types:
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	type wireType struct {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		ArrayT           *ArrayType
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		SliceT           *SliceType
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		StructT          *StructType
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		MapT             *MapType
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		GobEncoderT      *gobEncoderType
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		BinaryMarshalerT *gobEncoderType
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		TextMarshalerT   *gobEncoderType
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	type arrayType struct {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		CommonType
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		Elem typeId
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		Len  int
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	type CommonType struct {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		Name string // the name of the struct type
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		Id  int    // the id of the type, repeated so it&#39;s inside the type
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	type sliceType struct {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		CommonType
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		Elem typeId
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	type structType struct {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		CommonType
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		Field []*fieldType // the fields of the struct.
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	type fieldType struct {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		Name string // the name of the field.
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		Id   int    // the type id of the field, which must be already defined
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	type mapType struct {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		CommonType
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		Key  typeId
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		Elem typeId
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	type gobEncoderType struct {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		CommonType
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>If there are nested type ids, the types for all inner type ids must be defined
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>before the top-level type id is used to describe an encoded-v.
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>For simplicity in setup, the connection is defined to understand these types a
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>priori, as well as the basic gob types int, uint, etc. Their ids are:
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	bool        1
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	int         2
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	uint        3
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	float       4
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	[]byte      5
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	string      6
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	complex     7
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	interface   8
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	// gap for reserved ids.
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	WireType    16
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	ArrayType   17
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	CommonType  18
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	SliceType   19
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	StructType  20
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	FieldType   21
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	// 22 is slice of fieldType.
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	MapType     23
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>Finally, each message created by a call to Encode is preceded by an encoded
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>unsigned integer count of the number of bytes remaining in the message. After
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>the initial type name, interface values are wrapped the same way; in effect, the
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>interface value acts like a recursive invocation of Encode.
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>In summary, a gob stream looks like
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	(byteCount (-type id, encoding of a wireType)* (type id, encoding of a value))*
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>where * signifies zero or more repetitions and the type id of a value must
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>be predefined or be defined before the value in the stream.
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>Compatibility: Any future changes to the package will endeavor to maintain
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>compatibility with streams encoded using previous versions. That is, any released
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>version of this package should be able to decode data written with any previously
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>released version, subject to issues such as security fixes. See the Go compatibility
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>document for background: https://golang.org/doc/go1compat
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>See &#34;Gobs of data&#34; for a design discussion of the gob wire format:
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>https://blog.golang.org/gobs-of-data
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span># Security
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>This package is not designed to be hardened against adversarial inputs, and is
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>outside the scope of https://go.dev/security/policy. In particular, the [Decoder]
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>does only basic sanity checking on decoded input sizes, and its limits are not
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>configurable. Care should be taken when decoding gob data from untrusted
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>sources, which may consume significant resources.
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>*/</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>package gob
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">/*
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>Grammar:
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>Tokens starting with a lower case letter are terminals; int(n)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>and uint(n) represent the signed/unsigned encodings of the value n.
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>GobStream:
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	DelimitedMessage*
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>DelimitedMessage:
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	uint(lengthOfMessage) Message
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>Message:
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	TypeSequence TypedValue
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>TypeSequence
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	(TypeDefinition DelimitedTypeDefinition*)?
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>DelimitedTypeDefinition:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	uint(lengthOfTypeDefinition) TypeDefinition
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>TypedValue:
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	int(typeId) Value
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>TypeDefinition:
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	int(-typeId) encodingOfWireType
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>Value:
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	SingletonValue | StructValue
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>SingletonValue:
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	uint(0) FieldValue
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>FieldValue:
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	builtinValue | ArrayValue | MapValue | SliceValue | StructValue | InterfaceValue
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>InterfaceValue:
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	NilInterfaceValue | NonNilInterfaceValue
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>NilInterfaceValue:
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	uint(0)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>NonNilInterfaceValue:
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	ConcreteTypeName TypeSequence InterfaceContents
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>ConcreteTypeName:
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	uint(lengthOfName) [already read=n] name
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>InterfaceContents:
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	int(concreteTypeId) DelimitedValue
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>DelimitedValue:
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	uint(length) Value
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>ArrayValue:
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	uint(n) FieldValue*n [n elements]
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>MapValue:
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	uint(n) (FieldValue FieldValue)*n  [n (key, value) pairs]
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>SliceValue:
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	uint(n) FieldValue*n [n elements]
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>StructValue:
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	(uint(fieldDelta) FieldValue)*
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>*/</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">/*
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>For implementers and the curious, here is an encoded example. Given
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	type Point struct {X, Y int}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>and the value
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	p := Point{22, 33}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>the bytes transmitted that encode p will be:
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	1f ff 81 03 01 01 05 50 6f 69 6e 74 01 ff 82 00
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	01 02 01 01 58 01 04 00 01 01 59 01 04 00 00 00
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	07 ff 82 01 2c 01 42 00
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>They are determined as follows.
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>Since this is the first transmission of type Point, the type descriptor
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>for Point itself must be sent before the value. This is the first type
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>we&#39;ve sent on this Encoder, so it has type id 65 (0 through 64 are
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>reserved).
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	1f	// This item (a type descriptor) is 31 bytes long.
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	ff 81	// The negative of the id for the type we&#39;re defining, -65.
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		// This is one byte (indicated by FF = -1) followed by
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		// ^-65&lt;&lt;1 | 1. The low 1 bit signals to complement the
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		// rest upon receipt.
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	// Now we send a type descriptor, which is itself a struct (wireType).
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	// The type of wireType itself is known (it&#39;s built in, as is the type of
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	// all its components), so we just need to send a *value* of type wireType
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	// that represents type &#34;Point&#34;.
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	// Here starts the encoding of that value.
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	// Set the field number implicitly to -1; this is done at the beginning
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	// of every struct, including nested structs.
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	03	// Add 3 to field number; now 2 (wireType.structType; this is a struct).
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		// structType starts with an embedded CommonType, which appears
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		// as a regular structure here too.
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	01	// add 1 to field number (now 0); start of embedded CommonType.
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	01	// add 1 to field number (now 0, the name of the type)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	05	// string is (unsigned) 5 bytes long
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	50 6f 69 6e 74	// wireType.structType.CommonType.name = &#34;Point&#34;
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	01	// add 1 to field number (now 1, the id of the type)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	ff 82	// wireType.structType.CommonType._id = 65
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	00	// end of embedded wiretype.structType.CommonType struct
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	01	// add 1 to field number (now 1, the field array in wireType.structType)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	02	// There are two fields in the type (len(structType.field))
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	01	// Start of first field structure; add 1 to get field number 0: field[0].name
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	01	// 1 byte
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	58	// structType.field[0].name = &#34;X&#34;
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	01	// Add 1 to get field number 1: field[0].id
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	04	// structType.field[0].typeId is 2 (signed int).
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	00	// End of structType.field[0]; start structType.field[1]; set field number to -1.
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	01	// Add 1 to get field number 0: field[1].name
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	01	// 1 byte
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	59	// structType.field[1].name = &#34;Y&#34;
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	01	// Add 1 to get field number 1: field[1].id
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	04	// struct.Type.field[1].typeId is 2 (signed int).
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	00	// End of structType.field[1]; end of structType.field.
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	00	// end of wireType.structType structure
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	00	// end of wireType structure
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>Now we can send the Point value. Again the field number resets to -1:
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	07	// this value is 7 bytes long
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	ff 82	// the type number, 65 (1 byte (-FF) followed by 65&lt;&lt;1)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	01	// add one to field number, yielding field 0
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	2c	// encoding of signed &#34;22&#34; (0x2c = 44 = 22&lt;&lt;1); Point.x = 22
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	01	// add one to field number, yielding field 1
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	42	// encoding of signed &#34;33&#34; (0x42 = 66 = 33&lt;&lt;1); Point.y = 33
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	00	// end of structure
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>The type encoding is long and fairly intricate but we send it only once.
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>If p is transmitted a second time, the type is already known so the
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>output will be just:
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	07 ff 82 01 2c 01 42 00
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>A single non-struct value at top level is transmitted like a field with
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>delta tag 0. For instance, a signed integer with value 3 presented as
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>the argument to Encode will emit:
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	03 04 00 06
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>Which represents:
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	03	// this value is 3 bytes long
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	04	// the type number, 2, represents an integer
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	00	// tag delta 0
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	06	// value 3
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>*/</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
