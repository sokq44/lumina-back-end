<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/bytes/buffer.go - Go Documentation Server</title>

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
<a href="buffer.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/bytes">bytes</a>/<span class="text-muted">buffer.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/bytes">bytes</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package bytes
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Simple byte buffer for marshaling data.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// smallBufferSize is an initial allocation minimal capacity.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>const smallBufferSize = 64
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// A Buffer is a variable-sized buffer of bytes with [Buffer.Read] and [Buffer.Write] methods.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The zero value for Buffer is an empty buffer ready to use.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>type Buffer struct {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	buf      []byte <span class="comment">// contents are the bytes buf[off : len(buf)]</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	off      int    <span class="comment">// read at &amp;buf[off], write at &amp;buf[len(buf)]</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	lastRead readOp <span class="comment">// last read operation, so that Unread* can work correctly.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The readOp constants describe the last action performed on</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// the buffer, so that UnreadRune and UnreadByte can check for</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// invalid usage. opReadRuneX constants are chosen such that</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// converted to int they correspond to the rune size that was read.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>type readOp int8
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// Don&#39;t use iota for these, as the values need to correspond with the</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// names and comments, which is easier to see when being explicit.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>const (
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	opRead      readOp = -1 <span class="comment">// Any other read operation.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	opInvalid   readOp = 0  <span class="comment">// Non-read operation.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	opReadRune1 readOp = 1  <span class="comment">// Read rune of size 1.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	opReadRune2 readOp = 2  <span class="comment">// Read rune of size 2.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	opReadRune3 readOp = 3  <span class="comment">// Read rune of size 3.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	opReadRune4 readOp = 4  <span class="comment">// Read rune of size 4.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>var ErrTooLarge = errors.New(&#34;bytes.Buffer: too large&#34;)
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>var errNegativeRead = errors.New(&#34;bytes.Buffer: reader returned negative count from Read&#34;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>const maxInt = int(^uint(0) &gt;&gt; 1)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// The slice is valid for use only until the next buffer modification (that is,</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// only until the next call to a method like [Buffer.Read], [Buffer.Write], [Buffer.Reset], or [Buffer.Truncate]).</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// The slice aliases the buffer content at least until the next buffer modification,</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// so immediate changes to the slice will affect the result of future reads.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// AvailableBuffer returns an empty buffer with b.Available() capacity.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// This buffer is intended to be appended to and</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// passed to an immediately succeeding [Buffer.Write] call.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// The buffer is only valid until the next write operation on b.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func (b *Buffer) AvailableBuffer() []byte { return b.buf[len(b.buf):] }
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// String returns the contents of the unread portion of the buffer</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// as a string. If the [Buffer] is a nil pointer, it returns &#34;&lt;nil&gt;&#34;.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// To build strings more efficiently, see the strings.Builder type.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func (b *Buffer) String() string {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if b == nil {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		<span class="comment">// Special case, useful in debugging.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return &#34;&lt;nil&gt;&#34;
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return string(b.buf[b.off:])
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// empty reports whether the unread portion of the buffer is empty.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>func (b *Buffer) empty() bool { return len(b.buf) &lt;= b.off }
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// Len returns the number of bytes of the unread portion of the buffer;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// b.Len() == len(b.Bytes()).</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>func (b *Buffer) Len() int { return len(b.buf) - b.off }
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// Cap returns the capacity of the buffer&#39;s underlying byte slice, that is, the</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// total space allocated for the buffer&#39;s data.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (b *Buffer) Cap() int { return cap(b.buf) }
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// Available returns how many bytes are unused in the buffer.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (b *Buffer) Available() int { return cap(b.buf) - len(b.buf) }
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// Truncate discards all but the first n unread bytes from the buffer</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// but continues to use the same allocated storage.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// It panics if n is negative or greater than the length of the buffer.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (b *Buffer) Truncate(n int) {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	if n == 0 {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		b.Reset()
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if n &lt; 0 || n &gt; b.Len() {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		panic(&#34;bytes.Buffer: truncation out of range&#34;)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	b.buf = b.buf[:b.off+n]
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// Reset resets the buffer to be empty,</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// but it retains the underlying storage for use by future writes.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// Reset is the same as [Buffer.Truncate](0).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>func (b *Buffer) Reset() {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	b.buf = b.buf[:0]
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	b.off = 0
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// tryGrowByReslice is an inlineable version of grow for the fast-case where the</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// internal buffer only needs to be resliced.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// It returns the index where bytes should be written and whether it succeeded.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if l := len(b.buf); n &lt;= cap(b.buf)-l {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		b.buf = b.buf[:l+n]
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return l, true
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	return 0, false
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// grow grows the buffer to guarantee space for n more bytes.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// It returns the index where bytes should be written.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// If the buffer can&#39;t grow it will panic with ErrTooLarge.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func (b *Buffer) grow(n int) int {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	m := b.Len()
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// If buffer is empty, reset to recover space.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if m == 0 &amp;&amp; b.off != 0 {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		b.Reset()
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Try to grow by means of a reslice.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if i, ok := b.tryGrowByReslice(n); ok {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		return i
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	if b.buf == nil &amp;&amp; n &lt;= smallBufferSize {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		b.buf = make([]byte, n, smallBufferSize)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		return 0
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	c := cap(b.buf)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if n &lt;= c/2-m {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// We can slide things down instead of allocating a new</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// slice. We only need m+n &lt;= c to slide, but</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// we instead let capacity get twice as large so we</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// don&#39;t spend all our time copying.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		copy(b.buf, b.buf[b.off:])
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	} else if c &gt; maxInt-c-n {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		panic(ErrTooLarge)
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	} else {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// Add b.off to account for b.buf[:b.off] being sliced off the front.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		b.buf = growSlice(b.buf[b.off:], b.off+n)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// Restore b.off and len(b.buf).</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	b.off = 0
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	b.buf = b.buf[:m+n]
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	return m
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// Grow grows the buffer&#39;s capacity, if necessary, to guarantee space for</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// another n bytes. After Grow(n), at least n bytes can be written to the</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// buffer without another allocation.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// If n is negative, Grow will panic.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// If the buffer can&#39;t grow it will panic with [ErrTooLarge].</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>func (b *Buffer) Grow(n int) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		panic(&#34;bytes.Buffer.Grow: negative count&#34;)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	m := b.grow(n)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	b.buf = b.buf[:m]
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// Write appends the contents of p to the buffer, growing the buffer as</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// needed. The return value n is the length of p; err is always nil. If the</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// buffer becomes too large, Write will panic with [ErrTooLarge].</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func (b *Buffer) Write(p []byte) (n int, err error) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	m, ok := b.tryGrowByReslice(len(p))
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if !ok {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		m = b.grow(len(p))
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	return copy(b.buf[m:], p), nil
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// WriteString appends the contents of s to the buffer, growing the buffer as</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span><span class="comment">// needed. The return value n is the length of s; err is always nil. If the</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// buffer becomes too large, WriteString will panic with [ErrTooLarge].</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>func (b *Buffer) WriteString(s string) (n int, err error) {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	m, ok := b.tryGrowByReslice(len(s))
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	if !ok {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		m = b.grow(len(s))
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return copy(b.buf[m:], s), nil
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span><span class="comment">// MinRead is the minimum slice size passed to a Read call by</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span><span class="comment">// [Buffer.ReadFrom]. As long as the [Buffer] has at least MinRead bytes beyond</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// what is required to hold the contents of r, ReadFrom will not grow the</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// underlying buffer.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>const MinRead = 512
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// ReadFrom reads data from r until EOF and appends it to the buffer, growing</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// the buffer as needed. The return value n is the number of bytes read. Any</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// error except io.EOF encountered during the read is also returned. If the</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	for {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		i := b.grow(MinRead)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		b.buf = b.buf[:i]
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		m, e := r.Read(b.buf[i:cap(b.buf)])
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		if m &lt; 0 {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			panic(errNegativeRead)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		b.buf = b.buf[:i+m]
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		n += int64(m)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		if e == io.EOF {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			return n, nil <span class="comment">// e is EOF, so return nil explicitly</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		if e != nil {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			return n, e
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// growSlice grows b by n, preserving the original content of b.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span><span class="comment">// If the allocation fails, it panics with ErrTooLarge.</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>func growSlice(b []byte, n int) []byte {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	defer func() {
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		if recover() != nil {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			panic(ErrTooLarge)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}()
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// TODO(http://golang.org/issue/51462): We should rely on the append-make</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// pattern so that the compiler can call runtime.growslice. For example:</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">//	return append(b, make([]byte, n)...)</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// This avoids unnecessary zero-ing of the first len(b) bytes of the</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// allocated slice, but this pattern causes b to escape onto the heap.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// Instead use the append-make pattern with a nil slice to ensure that</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">// we allocate buffers rounded up to the closest size class.</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	c := len(b) + n <span class="comment">// ensure enough space for n elements</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if c &lt; 2*cap(b) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		<span class="comment">// The growth rate has historically always been 2x. In the future,</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		<span class="comment">// we could rely purely on append to determine the growth rate.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		c = 2 * cap(b)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	b2 := append([]byte(nil), make([]byte, c)...)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	copy(b2, b)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	return b2[:len(b)]
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span><span class="comment">// WriteTo writes data to w until the buffer is drained or an error occurs.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span><span class="comment">// The return value n is the number of bytes written; it always fits into an</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// int, but it is int64 to match the io.WriterTo interface. Any error</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span><span class="comment">// encountered during the write is also returned.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	if nBytes := b.Len(); nBytes &gt; 0 {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		m, e := w.Write(b.buf[b.off:])
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		if m &gt; nBytes {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			panic(&#34;bytes.Buffer.WriteTo: invalid Write count&#34;)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		b.off += m
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		n = int64(m)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		if e != nil {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			return n, e
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		<span class="comment">// all bytes should have been written, by definition of</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		<span class="comment">// Write method in io.Writer</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		if m != nBytes {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			return n, io.ErrShortWrite
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">// Buffer is now empty; reset.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	b.Reset()
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return n, nil
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// WriteByte appends the byte c to the buffer, growing the buffer as needed.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// The returned error is always nil, but is included to match [bufio.Writer]&#39;s</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// WriteByte. If the buffer becomes too large, WriteByte will panic with</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// [ErrTooLarge].</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>func (b *Buffer) WriteByte(c byte) error {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	m, ok := b.tryGrowByReslice(1)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	if !ok {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		m = b.grow(1)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	b.buf[m] = c
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	return nil
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">// WriteRune appends the UTF-8 encoding of Unicode code point r to the</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span><span class="comment">// buffer, returning its length and an error, which is always nil but is</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">// included to match [bufio.Writer]&#39;s WriteRune. The buffer is grown as needed;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span><span class="comment">// if it becomes too large, WriteRune will panic with [ErrTooLarge].</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>func (b *Buffer) WriteRune(r rune) (n int, err error) {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// Compare as uint32 to correctly handle negative runes.</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	if uint32(r) &lt; utf8.RuneSelf {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		b.WriteByte(byte(r))
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		return 1, nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	m, ok := b.tryGrowByReslice(utf8.UTFMax)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if !ok {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		m = b.grow(utf8.UTFMax)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	b.buf = utf8.AppendRune(b.buf[:m], r)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	return len(b.buf) - m, nil
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// Read reads the next len(p) bytes from the buffer or until the buffer</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// is drained. The return value n is the number of bytes read. If the</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// buffer has no data to return, err is io.EOF (unless len(p) is zero);</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span><span class="comment">// otherwise it is nil.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>func (b *Buffer) Read(p []byte) (n int, err error) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	if b.empty() {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		<span class="comment">// Buffer is empty, reset to recover space.</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		b.Reset()
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		if len(p) == 0 {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			return 0, nil
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	n = copy(p, b.buf[b.off:])
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	b.off += n
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		b.lastRead = opRead
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	return n, nil
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// Next returns a slice containing the next n bytes from the buffer,</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// advancing the buffer as if the bytes had been returned by [Buffer.Read].</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// If there are fewer than n bytes in the buffer, Next returns the entire buffer.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// The slice is only valid until the next call to a read or write method.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>func (b *Buffer) Next(n int) []byte {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	m := b.Len()
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if n &gt; m {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		n = m
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	data := b.buf[b.off : b.off+n]
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	b.off += n
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	if n &gt; 0 {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		b.lastRead = opRead
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	return data
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// ReadByte reads and returns the next byte from the buffer.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// If no byte is available, it returns error io.EOF.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>func (b *Buffer) ReadByte() (byte, error) {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	if b.empty() {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		<span class="comment">// Buffer is empty, reset to recover space.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		b.Reset()
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	c := b.buf[b.off]
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	b.off++
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	b.lastRead = opRead
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	return c, nil
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// ReadRune reads and returns the next UTF-8-encoded</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// Unicode code point from the buffer.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// If no bytes are available, the error returned is io.EOF.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// If the bytes are an erroneous UTF-8 encoding, it</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// consumes one byte and returns U+FFFD, 1.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func (b *Buffer) ReadRune() (r rune, size int, err error) {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	if b.empty() {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		<span class="comment">// Buffer is empty, reset to recover space.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		b.Reset()
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		return 0, 0, io.EOF
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	c := b.buf[b.off]
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	if c &lt; utf8.RuneSelf {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		b.off++
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		b.lastRead = opReadRune1
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		return rune(c), 1, nil
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	r, n := utf8.DecodeRune(b.buf[b.off:])
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	b.off += n
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	b.lastRead = readOp(n)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	return r, n, nil
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// UnreadRune unreads the last rune returned by [Buffer.ReadRune].</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// If the most recent read or write operation on the buffer was</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// not a successful [Buffer.ReadRune], UnreadRune returns an error.  (In this regard</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// it is stricter than [Buffer.UnreadByte], which will unread the last byte</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// from any read operation.)</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func (b *Buffer) UnreadRune() error {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	if b.lastRead &lt;= opInvalid {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		return errors.New(&#34;bytes.Buffer: UnreadRune: previous operation was not a successful ReadRune&#34;)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	if b.off &gt;= int(b.lastRead) {
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		b.off -= int(b.lastRead)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	return nil
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>var errUnreadByte = errors.New(&#34;bytes.Buffer: UnreadByte: previous operation was not a successful read&#34;)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">// UnreadByte unreads the last byte returned by the most recent successful</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">// read operation that read at least one byte. If a write has happened since</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">// the last read, if the last read returned an error, or if the read read zero</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">// bytes, UnreadByte returns an error.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>func (b *Buffer) UnreadByte() error {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	if b.lastRead == opInvalid {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		return errUnreadByte
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	b.lastRead = opInvalid
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	if b.off &gt; 0 {
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		b.off--
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	return nil
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span><span class="comment">// ReadBytes reads until the first occurrence of delim in the input,</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span><span class="comment">// returning a slice containing the data up to and including the delimiter.</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span><span class="comment">// If ReadBytes encounters an error before finding a delimiter,</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span><span class="comment">// it returns the data read before the error and the error itself (often io.EOF).</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span><span class="comment">// ReadBytes returns err != nil if and only if the returned data does not end in</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span><span class="comment">// delim.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	slice, err := b.readSlice(delim)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	<span class="comment">// return a copy of slice. The buffer&#39;s backing array may</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	<span class="comment">// be overwritten by later calls.</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	line = append(line, slice...)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	return line, err
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// readSlice is like ReadBytes but returns a reference to internal buffer data.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>func (b *Buffer) readSlice(delim byte) (line []byte, err error) {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	i := IndexByte(b.buf[b.off:], delim)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	end := b.off + i + 1
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	if i &lt; 0 {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		end = len(b.buf)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		err = io.EOF
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	line = b.buf[b.off:end]
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	b.off = end
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	b.lastRead = opRead
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	return line, err
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// ReadString reads until the first occurrence of delim in the input,</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span><span class="comment">// returning a string containing the data up to and including the delimiter.</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// If ReadString encounters an error before finding a delimiter,</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// it returns the data read before the error and the error itself (often io.EOF).</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span><span class="comment">// ReadString returns err != nil if and only if the returned data does not end</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span><span class="comment">// in delim.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>func (b *Buffer) ReadString(delim byte) (line string, err error) {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	slice, err := b.readSlice(delim)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	return string(slice), err
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// NewBuffer creates and initializes a new [Buffer] using buf as its</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// initial contents. The new [Buffer] takes ownership of buf, and the</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">// caller should not use buf after this call. NewBuffer is intended to</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span><span class="comment">// prepare a [Buffer] to read existing data. It can also be used to set</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">// the initial size of the internal buffer for writing. To do that,</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">// buf should have the desired capacity but a length of zero.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// In most cases, new([Buffer]) (or just declaring a [Buffer] variable) is</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// sufficient to initialize a [Buffer].</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>func NewBuffer(buf []byte) *Buffer { return &amp;Buffer{buf: buf} }
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// NewBufferString creates and initializes a new [Buffer] using string s as its</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">// initial contents. It is intended to prepare a buffer to read an existing</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">// string.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">// In most cases, new([Buffer]) (or just declaring a [Buffer] variable) is</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">// sufficient to initialize a [Buffer].</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>func NewBufferString(s string) *Buffer {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	return &amp;Buffer{buf: []byte(s)}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
</pre><p><a href="buffer.go?m=text">View as plain text</a></p>

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
