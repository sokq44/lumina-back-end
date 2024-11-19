<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/profbuf.go - Go Documentation Server</title>

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
<a href="profbuf.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">profbuf.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// A profBuf is a lock-free buffer for profiling events,</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// safe for concurrent use by one reader and one writer.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// The writer may be a signal handler running without a user g.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The reader is assumed to be a user g.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Each logged event corresponds to a fixed size header, a list of</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// uintptrs (typically a stack), and exactly one unsafe.Pointer tag.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// The header and uintptrs are stored in the circular buffer data and the</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// tag is stored in a circular buffer tags, running in parallel.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// In the circular buffer data, each event takes 2+hdrsize+len(stk)</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// words: the value 2+hdrsize+len(stk), then the time of the event, then</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// hdrsize words giving the fixed-size header, and then len(stk) words</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// for the stack.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The current effective offsets into the tags and data circular buffers</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// for reading and writing are stored in the high 30 and low 32 bits of r and w.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// The bottom bits of the high 32 are additional flag bits in w, unused in r.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// &#34;Effective&#34; offsets means the total number of reads or writes, mod 2^length.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// The offset in the buffer is the effective offset mod the length of the buffer.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// To make wraparound mod 2^length match wraparound mod length of the buffer,</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// the length of the buffer must be a power of two.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// If the reader catches up to the writer, a flag passed to read controls</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// whether the read blocks until more data is available. A read returns a</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// pointer to the buffer data itself; the caller is assumed to be done with</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// that data at the next read. The read offset rNext tracks the next offset to</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// be returned by read. By definition, r ≤ rNext ≤ w (before wraparound),</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// and rNext is only used by the reader, so it can be accessed without atomics.</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// If the writer gets ahead of the reader, so that the buffer fills,</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// future writes are discarded and replaced in the output stream by an</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// overflow entry, which has size 2+hdrsize+1, time set to the time of</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// the first discarded write, a header of all zeroed words, and a &#34;stack&#34;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// containing one word, the number of discarded writes.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// Between the time the buffer fills and the buffer becomes empty enough</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// to hold more data, the overflow entry is stored as a pending overflow</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// entry in the fields overflow and overflowTime. The pending overflow</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// entry can be turned into a real record by either the writer or the</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// reader. If the writer is called to write a new record and finds that</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// the output buffer has room for both the pending overflow entry and the</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// new record, the writer emits the pending overflow entry and the new</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// record into the buffer. If the reader is called to read data and finds</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// that the output buffer is empty but that there is a pending overflow</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// entry, the reader will return a synthesized record for the pending</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// overflow entry.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// Only the writer can create or add to a pending overflow entry, but</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// either the reader or the writer can clear the pending overflow entry.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// A pending overflow entry is indicated by the low 32 bits of &#39;overflow&#39;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// holding the number of discarded writes, and overflowTime holding the</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// time of the first discarded write. The high 32 bits of &#39;overflow&#39;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// increment each time the low 32 bits transition from zero to non-zero</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// or vice versa. This sequence number avoids ABA problems in the use of</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// compare-and-swap to coordinate between reader and writer.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// The overflowTime is only written when the low 32 bits of overflow are</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// zero, that is, only when there is no pending overflow entry, in</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// preparation for creating a new one. The reader can therefore fetch and</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// clear the entry atomically using</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//	for {</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//		overflow = load(&amp;b.overflow)</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//		if uint32(overflow) == 0 {</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//			// no pending entry</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//			break</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//		}</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//		time = load(&amp;b.overflowTime)</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//		if cas(&amp;b.overflow, overflow, ((overflow&gt;&gt;32)+1)&lt;&lt;32) {</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//			// pending entry cleared</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//			break</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//		}</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//	if uint32(overflow) &gt; 0 {</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//		emit entry for uint32(overflow), time</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>type profBuf struct {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// accessed atomically</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	r, w         profAtomic
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	overflow     atomic.Uint64
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	overflowTime atomic.Uint64
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	eof          atomic.Uint32
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// immutable (excluding slice content)</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	hdrsize uintptr
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	data    []uint64
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	tags    []unsafe.Pointer
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// owned by reader</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	rNext       profIndex
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	overflowBuf []uint64 <span class="comment">// for use by reader to return overflow record</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	wait        note
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// A profAtomic is the atomically-accessed word holding a profIndex.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>type profAtomic uint64
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// A profIndex is the packet tag and data counts and flags bits, described above.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>type profIndex uint64
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>const (
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	profReaderSleeping profIndex = 1 &lt;&lt; 32 <span class="comment">// reader is sleeping and must be woken up</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	profWriteExtra     profIndex = 1 &lt;&lt; 33 <span class="comment">// overflow or eof waiting</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>func (x *profAtomic) load() profIndex {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	return profIndex(atomic.Load64((*uint64)(x)))
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>func (x *profAtomic) store(new profIndex) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	atomic.Store64((*uint64)(x), uint64(new))
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>func (x *profAtomic) cas(old, new profIndex) bool {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	return atomic.Cas64((*uint64)(x), uint64(old), uint64(new))
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func (x profIndex) dataCount() uint32 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	return uint32(x)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>func (x profIndex) tagCount() uint32 {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	return uint32(x &gt;&gt; 34)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// countSub subtracts two counts obtained from profIndex.dataCount or profIndex.tagCount,</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// assuming that they are no more than 2^29 apart (guaranteed since they are never more than</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// len(data) or len(tags) apart, respectively).</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// tagCount wraps at 2^30, while dataCount wraps at 2^32.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// This function works for both.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>func countSub(x, y uint32) int {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// x-y is 32-bit signed or 30-bit signed; sign-extend to 32 bits and convert to int.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	return int(int32(x-y) &lt;&lt; 2 &gt;&gt; 2)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// addCountsAndClearFlags returns the packed form of &#34;x + (data, tag) - all flags&#34;.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>func (x profIndex) addCountsAndClearFlags(data, tag int) profIndex {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	return profIndex((uint64(x)&gt;&gt;34+uint64(uint32(tag)&lt;&lt;2&gt;&gt;2))&lt;&lt;34 | uint64(uint32(x)+uint32(data)))
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// hasOverflow reports whether b has any overflow records pending.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func (b *profBuf) hasOverflow() bool {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	return uint32(b.overflow.Load()) &gt; 0
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// takeOverflow consumes the pending overflow records, returning the overflow count</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// and the time of the first overflow.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// When called by the reader, it is racing against incrementOverflow.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>func (b *profBuf) takeOverflow() (count uint32, time uint64) {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	overflow := b.overflow.Load()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	time = b.overflowTime.Load()
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	for {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		count = uint32(overflow)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		if count == 0 {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			time = 0
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			break
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// Increment generation, clear overflow count in low bits.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if b.overflow.CompareAndSwap(overflow, ((overflow&gt;&gt;32)+1)&lt;&lt;32) {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			break
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		overflow = b.overflow.Load()
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		time = b.overflowTime.Load()
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	return uint32(overflow), time
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span><span class="comment">// incrementOverflow records a single overflow at time now.</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span><span class="comment">// It is racing against a possible takeOverflow in the reader.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>func (b *profBuf) incrementOverflow(now int64) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	for {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		overflow := b.overflow.Load()
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// Once we see b.overflow reach 0, it&#39;s stable: no one else is changing it underfoot.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// We need to set overflowTime if we&#39;re incrementing b.overflow from 0.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		if uint32(overflow) == 0 {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// Store overflowTime first so it&#39;s always available when overflow != 0.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			b.overflowTime.Store(uint64(now))
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			b.overflow.Store((((overflow &gt;&gt; 32) + 1) &lt;&lt; 32) + 1)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			break
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		<span class="comment">// Otherwise we&#39;re racing to increment against reader</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		<span class="comment">// who wants to set b.overflow to 0.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		<span class="comment">// Out of paranoia, leave 2³²-1 a sticky overflow value,</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		<span class="comment">// to avoid wrapping around. Extremely unlikely.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		if int32(overflow) == -1 {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			break
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if b.overflow.CompareAndSwap(overflow, overflow+1) {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			break
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// newProfBuf returns a new profiling buffer with room for</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// a header of hdrsize words and a buffer of at least bufwords words.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>func newProfBuf(hdrsize, bufwords, tags int) *profBuf {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if min := 2 + hdrsize + 1; bufwords &lt; min {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		bufwords = min
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Buffer sizes must be power of two, so that we don&#39;t have to</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// worry about uint32 wraparound changing the effective position</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// within the buffers. We store 30 bits of count; limiting to 28</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// gives us some room for intermediate calculations.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if bufwords &gt;= 1&lt;&lt;28 || tags &gt;= 1&lt;&lt;28 {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		throw(&#34;newProfBuf: buffer too large&#34;)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	var i int
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	for i = 1; i &lt; bufwords; i &lt;&lt;= 1 {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	bufwords = i
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	for i = 1; i &lt; tags; i &lt;&lt;= 1 {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	tags = i
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	b := new(profBuf)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	b.hdrsize = uintptr(hdrsize)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	b.data = make([]uint64, bufwords)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	b.tags = make([]unsafe.Pointer, tags)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	b.overflowBuf = make([]uint64, 2+b.hdrsize+1)
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	return b
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span><span class="comment">// canWriteRecord reports whether the buffer has room</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span><span class="comment">// for a single contiguous record with a stack of length nstk.</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>func (b *profBuf) canWriteRecord(nstk int) bool {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	br := b.r.load()
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	bw := b.w.load()
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// room for tag?</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if countSub(br.tagCount(), bw.tagCount())+len(b.tags) &lt; 1 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		return false
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// room for data?</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	nd := countSub(br.dataCount(), bw.dataCount()) + len(b.data)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	want := 2 + int(b.hdrsize) + nstk
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	i := int(bw.dataCount() % uint32(len(b.data)))
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	if i+want &gt; len(b.data) {
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t fit in trailing fragment of slice.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		<span class="comment">// Skip over that and start over at beginning of slice.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		nd -= len(b.data) - i
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	return nd &gt;= want
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span><span class="comment">// canWriteTwoRecords reports whether the buffer has room</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// for two records with stack lengths nstk1, nstk2, in that order.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// Each record must be contiguous on its own, but the two</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// records need not be contiguous (one can be at the end of the buffer</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span><span class="comment">// and the other can wrap around and start at the beginning of the buffer).</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>func (b *profBuf) canWriteTwoRecords(nstk1, nstk2 int) bool {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	br := b.r.load()
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	bw := b.w.load()
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	<span class="comment">// room for tag?</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	if countSub(br.tagCount(), bw.tagCount())+len(b.tags) &lt; 2 {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		return false
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	<span class="comment">// room for data?</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	nd := countSub(br.dataCount(), bw.dataCount()) + len(b.data)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// first record</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	want := 2 + int(b.hdrsize) + nstk1
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	i := int(bw.dataCount() % uint32(len(b.data)))
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	if i+want &gt; len(b.data) {
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t fit in trailing fragment of slice.</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		<span class="comment">// Skip over that and start over at beginning of slice.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		nd -= len(b.data) - i
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		i = 0
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	i += want
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	nd -= want
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// second record</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	want = 2 + int(b.hdrsize) + nstk2
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if i+want &gt; len(b.data) {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		<span class="comment">// Can&#39;t fit in trailing fragment of slice.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		<span class="comment">// Skip over that and start over at beginning of slice.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		nd -= len(b.data) - i
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		i = 0
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	return nd &gt;= want
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span><span class="comment">// write writes an entry to the profiling buffer b.</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">// The entry begins with a fixed hdr, which must have</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// length b.hdrsize, followed by a variable-sized stack</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// and a single tag pointer *tagPtr (or nil if tagPtr is nil).</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// No write barriers allowed because this might be called from a signal handler.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>func (b *profBuf) write(tagPtr *unsafe.Pointer, now int64, hdr []uint64, stk []uintptr) {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if b == nil {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		return
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if len(hdr) &gt; int(b.hdrsize) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		throw(&#34;misuse of profBuf.write&#34;)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if hasOverflow := b.hasOverflow(); hasOverflow &amp;&amp; b.canWriteTwoRecords(1, len(stk)) {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		<span class="comment">// Room for both an overflow record and the one being written.</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		<span class="comment">// Write the overflow record if the reader hasn&#39;t gotten to it yet.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		<span class="comment">// Only racing against reader, not other writers.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		count, time := b.takeOverflow()
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if count &gt; 0 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			var stk [1]uintptr
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			stk[0] = uintptr(count)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>			b.write(nil, int64(time), nil, stk[:])
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	} else if hasOverflow || !b.canWriteRecord(len(stk)) {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		<span class="comment">// Pending overflow without room to write overflow and new records</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// or no overflow but also no room for new record.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		b.incrementOverflow(now)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		b.wakeupExtra()
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		return
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// There&#39;s room: write the record.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	br := b.r.load()
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	bw := b.w.load()
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// Profiling tag</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// The tag is a pointer, but we can&#39;t run a write barrier here.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">// We have interrupted the OS-level execution of gp, but the</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// runtime still sees gp as executing. In effect, we are running</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// in place of the real gp. Since gp is the only goroutine that</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// can overwrite gp.labels, the value of gp.labels is stable during</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	<span class="comment">// this signal handler: it will still be reachable from gp when</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// we finish executing. If a GC is in progress right now, it must</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// keep gp.labels alive, because gp.labels is reachable from gp.</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// If gp were to overwrite gp.labels, the deletion barrier would</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	<span class="comment">// still shade that pointer, which would preserve it for the</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">// in-progress GC, so all is well. Any future GC will see the</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">// value we copied when scanning b.tags (heap-allocated).</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	<span class="comment">// We arrange that the store here is always overwriting a nil,</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	<span class="comment">// so there is no need for a deletion barrier on b.tags[wt].</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	wt := int(bw.tagCount() % uint32(len(b.tags)))
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	if tagPtr != nil {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		*(*uintptr)(unsafe.Pointer(&amp;b.tags[wt])) = uintptr(*tagPtr)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// Main record.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// It has to fit in a contiguous section of the slice, so if it doesn&#39;t fit at the end,</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">// leave a rewind marker (0) and start over at the beginning of the slice.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	wd := int(bw.dataCount() % uint32(len(b.data)))
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	nd := countSub(br.dataCount(), bw.dataCount()) + len(b.data)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	skip := 0
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	if wd+2+int(b.hdrsize)+len(stk) &gt; len(b.data) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		b.data[wd] = 0
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		skip = len(b.data) - wd
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		nd -= skip
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		wd = 0
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	data := b.data[wd:]
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	data[0] = uint64(2 + b.hdrsize + uintptr(len(stk))) <span class="comment">// length</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	data[1] = uint64(now)                               <span class="comment">// time stamp</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// header, zero-padded</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	i := uintptr(copy(data[2:2+b.hdrsize], hdr))
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	for ; i &lt; b.hdrsize; i++ {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		data[2+i] = 0
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	for i, pc := range stk {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		data[2+b.hdrsize+uintptr(i)] = uint64(pc)
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	for {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// Commit write.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// Racing with reader setting flag bits in b.w, to avoid lost wakeups.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		old := b.w.load()
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		new := old.addCountsAndClearFlags(skip+2+len(stk)+int(b.hdrsize), 1)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		if !b.w.cas(old, new) {
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>			continue
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		<span class="comment">// If there was a reader, wake it up.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		if old&amp;profReaderSleeping != 0 {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>			notewakeup(&amp;b.wait)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		break
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// close signals that there will be no more writes on the buffer.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// Once all the data has been read from the buffer, reads will return eof=true.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func (b *profBuf) close() {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	if b.eof.Load() &gt; 0 {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		throw(&#34;runtime: profBuf already closed&#34;)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	b.eof.Store(1)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	b.wakeupExtra()
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">// wakeupExtra must be called after setting one of the &#34;extra&#34;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">// atomic fields b.overflow or b.eof.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">// It records the change in b.w and wakes up the reader if needed.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>func (b *profBuf) wakeupExtra() {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	for {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		old := b.w.load()
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		new := old | profWriteExtra
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		if !b.w.cas(old, new) {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			continue
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		if old&amp;profReaderSleeping != 0 {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			notewakeup(&amp;b.wait)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		break
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	}
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// profBufReadMode specifies whether to block when no data is available to read.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>type profBufReadMode int
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>const (
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	profBufBlocking profBufReadMode = iota
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	profBufNonBlocking
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>var overflowTag [1]unsafe.Pointer <span class="comment">// always nil</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>func (b *profBuf) read(mode profBufReadMode) (data []uint64, tags []unsafe.Pointer, eof bool) {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	if b == nil {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		return nil, nil, true
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	br := b.rNext
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// Commit previous read, returning that part of the ring to the writer.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">// First clear tags that have now been read, both to avoid holding</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// up the memory they point at for longer than necessary</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// and so that b.write can assume it is always overwriting</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	<span class="comment">// nil tag entries (see comment in b.write).</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	rPrev := b.r.load()
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	if rPrev != br {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		ntag := countSub(br.tagCount(), rPrev.tagCount())
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		ti := int(rPrev.tagCount() % uint32(len(b.tags)))
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		for i := 0; i &lt; ntag; i++ {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			b.tags[ti] = nil
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			if ti++; ti == len(b.tags) {
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				ti = 0
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		b.r.store(br)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>Read:
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	bw := b.w.load()
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	numData := countSub(bw.dataCount(), br.dataCount())
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	if numData == 0 {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		if b.hasOverflow() {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			<span class="comment">// No data to read, but there is overflow to report.</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			<span class="comment">// Racing with writer flushing b.overflow into a real record.</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			count, time := b.takeOverflow()
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			if count == 0 {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				<span class="comment">// Lost the race, go around again.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				goto Read
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>			<span class="comment">// Won the race, report overflow.</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			dst := b.overflowBuf
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			dst[0] = uint64(2 + b.hdrsize + 1)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			dst[1] = time
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; b.hdrsize; i++ {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>				dst[2+i] = 0
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			}
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			dst[2+b.hdrsize] = uint64(count)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			return dst[:2+b.hdrsize+1], overflowTag[:1], false
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		if b.eof.Load() &gt; 0 {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			<span class="comment">// No data, no overflow, EOF set: done.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			return nil, nil, true
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		if bw&amp;profWriteExtra != 0 {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			<span class="comment">// Writer claims to have published extra information (overflow or eof).</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// Attempt to clear notification and then check again.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			<span class="comment">// If we fail to clear the notification it means b.w changed,</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			<span class="comment">// so we still need to check again.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			b.w.cas(bw, bw&amp;^profWriteExtra)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			goto Read
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// Nothing to read right now.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">// Return or sleep according to mode.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		if mode == profBufNonBlocking {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			<span class="comment">// Necessary on Darwin, notetsleepg below does not work in signal handler, root cause of #61768.</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			return nil, nil, false
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		if !b.w.cas(bw, bw|profReaderSleeping) {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			goto Read
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		<span class="comment">// Committed to sleeping.</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		notetsleepg(&amp;b.wait, -1)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		noteclear(&amp;b.wait)
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		goto Read
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	data = b.data[br.dataCount()%uint32(len(b.data)):]
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	if len(data) &gt; numData {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		data = data[:numData]
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	} else {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		numData -= len(data) <span class="comment">// available in case of wraparound</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	skip := 0
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	if data[0] == 0 {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		<span class="comment">// Wraparound record. Go back to the beginning of the ring.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		skip = len(data)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		data = b.data
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		if len(data) &gt; numData {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			data = data[:numData]
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		}
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	ntag := countSub(bw.tagCount(), br.tagCount())
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	if ntag == 0 {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		throw(&#34;runtime: malformed profBuf buffer - tag and data out of sync&#34;)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	tags = b.tags[br.tagCount()%uint32(len(b.tags)):]
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if len(tags) &gt; ntag {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		tags = tags[:ntag]
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	}
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// Count out whole data records until either data or tags is done.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	<span class="comment">// They are always in sync in the buffer, but due to an end-of-slice</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	<span class="comment">// wraparound we might need to stop early and return the rest</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	<span class="comment">// in the next call.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	di := 0
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	ti := 0
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	for di &lt; len(data) &amp;&amp; data[di] != 0 &amp;&amp; ti &lt; len(tags) {
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		if uintptr(di)+uintptr(data[di]) &gt; uintptr(len(data)) {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			throw(&#34;runtime: malformed profBuf buffer - invalid size&#34;)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		di += int(data[di])
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		ti++
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// Remember how much we returned, to commit read on next call.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	b.rNext = br.addCountsAndClearFlags(skip+di, ti)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	if raceenabled {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		<span class="comment">// Match racereleasemerge in runtime_setProfLabel,</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		<span class="comment">// so that the setting of the labels in runtime_setProfLabel</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		<span class="comment">// is treated as happening before any use of the labels</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		<span class="comment">// by our caller. The synchronization on labelSync itself is a fiction</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		<span class="comment">// for the race detector. The actual synchronization is handled</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		<span class="comment">// by the fact that the signal handler only reads from the current</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		<span class="comment">// goroutine and uses atomics to write the updated queue indices,</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		<span class="comment">// and then the read-out from the signal handler buffer uses</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		<span class="comment">// atomics to read those queue indices.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		raceacquire(unsafe.Pointer(&amp;labelSync))
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	return data[:di], tags[:ti], false
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
</pre><p><a href="profbuf.go?m=text">View as plain text</a></p>

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
