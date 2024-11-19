<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/map.go - Go Documentation Server</title>

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
<a href="map.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">map.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package runtime
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file contains the implementation of Go&#39;s map type.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// A map is just a hash table. The data is arranged</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// into an array of buckets. Each bucket contains up to</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// 8 key/elem pairs. The low-order bits of the hash are</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// used to select a bucket. Each bucket contains a few</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// high-order bits of each hash to distinguish the entries</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// within a single bucket.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// If more than 8 keys hash to a bucket, we chain on</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// extra buckets.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// When the hashtable grows, we allocate a new array</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// of buckets twice as big. Buckets are incrementally</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// copied from the old bucket array to the new bucket array.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// Map iterators walk through the array of buckets and</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// return the keys in walk order (bucket #, then overflow</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// chain order, then bucket index).  To maintain iteration</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// semantics, we never move keys within their bucket (if</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// we did, keys might be returned 0 or 2 times).  When</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// growing the table, iterators remain iterating through the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// old table and must check the new table if the bucket</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// they are iterating through has been moved (&#34;evacuated&#34;)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// to the new table.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// Picking loadFactor: too large and we have lots of overflow</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// buckets, too small and we waste a lot of space. I wrote</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// a simple program to check some stats for different loads:</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// (64-bit, 8 byte keys and elems)</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//  loadFactor    %overflow  bytes/entry     hitprobe    missprobe</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//        4.00         2.13        20.77         3.00         4.00</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//        4.50         4.05        17.30         3.25         4.50</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//        5.00         6.85        14.77         3.50         5.00</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//        5.50        10.55        12.94         3.75         5.50</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//        6.00        15.27        11.67         4.00         6.00</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//        6.50        20.90        10.79         4.25         6.50</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//        7.00        27.14        10.15         4.50         7.00</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//        7.50        34.03         9.73         4.75         7.50</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//        8.00        41.10         9.40         5.00         8.00</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// %overflow   = percentage of buckets which have an overflow bucket</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// bytes/entry = overhead bytes used per key/elem pair</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// hitprobe    = # of entries to check when looking up a present key</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// missprobe   = # of entries to check when looking up an absent key</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// Keep in mind this data is for maximally loaded tables, i.e. just</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// before the table grows. Typical tables will be somewhat less loaded.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>import (
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	&#34;internal/abi&#34;
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	&#34;runtime/internal/math&#34;
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>const (
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Maximum number of key/elem pairs a bucket can hold.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	bucketCntBits = abi.MapBucketCountBits
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	bucketCnt     = abi.MapBucketCount
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// Maximum average load of a bucket that triggers growth is bucketCnt*13/16 (about 80% full)</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// Because of minimum alignment rules, bucketCnt is known to be at least 8.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	<span class="comment">// Represent as loadFactorNum/loadFactorDen, to allow integer math.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	loadFactorDen = 2
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	loadFactorNum = loadFactorDen * bucketCnt * 13 / 16
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// Maximum key or elem size to keep inline (instead of mallocing per element).</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// Must fit in a uint8.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// Fast versions cannot handle big elems - the cutoff size for</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// fast versions in cmd/compile/internal/gc/walk.go must be at most this elem.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	maxKeySize  = abi.MapMaxKeyBytes
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	maxElemSize = abi.MapMaxElemBytes
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// data offset should be the size of the bmap struct, but needs to be</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// aligned correctly. For amd64p32 this means 64-bit alignment</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// even though pointers are 32 bit.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	dataOffset = unsafe.Offsetof(struct {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		b bmap
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		v int64
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}{}.v)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// Possible tophash values. We reserve a few possibilities for special marks.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Each bucket (including its overflow buckets, if any) will have either all or none of its</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// entries in the evacuated* states (except during the evacuate() method, which only happens</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// during map writes and thus no one else can observe the map during that time).</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	emptyRest      = 0 <span class="comment">// this cell is empty, and there are no more non-empty cells at higher indexes or overflows.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	emptyOne       = 1 <span class="comment">// this cell is empty</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	evacuatedX     = 2 <span class="comment">// key/elem is valid.  Entry has been evacuated to first half of larger table.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	evacuatedY     = 3 <span class="comment">// same as above, but evacuated to second half of larger table.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	evacuatedEmpty = 4 <span class="comment">// cell is empty, bucket is evacuated.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	minTopHash     = 5 <span class="comment">// minimum tophash for a normal filled cell.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// flags</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	iterator     = 1 <span class="comment">// there may be an iterator using buckets</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	oldIterator  = 2 <span class="comment">// there may be an iterator using oldbuckets</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	hashWriting  = 4 <span class="comment">// a goroutine is writing to the map</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	sameSizeGrow = 8 <span class="comment">// the current map growth is to a new map of the same size</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// sentinel bucket ID for iterator checks</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	noCheck = 1&lt;&lt;(8*goarch.PtrSize) - 1
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// isEmpty reports whether the given tophash array entry represents an empty bucket entry.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>func isEmpty(x uint8) bool {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	return x &lt;= emptyOne
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// A header for a Go map.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>type hmap struct {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// Make sure this stays in sync with the compiler&#39;s definition.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	count     int <span class="comment">// # live cells == size of map.  Must be first (used by len() builtin)</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	flags     uint8
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	B         uint8  <span class="comment">// log_2 of # of buckets (can hold up to loadFactor * 2^B items)</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	noverflow uint16 <span class="comment">// approximate number of overflow buckets; see incrnoverflow for details</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	hash0     uint32 <span class="comment">// hash seed</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	buckets    unsafe.Pointer <span class="comment">// array of 2^B Buckets. may be nil if count==0.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	oldbuckets unsafe.Pointer <span class="comment">// previous bucket array of half the size, non-nil only when growing</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	nevacuate  uintptr        <span class="comment">// progress counter for evacuation (buckets less than this have been evacuated)</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	extra *mapextra <span class="comment">// optional fields</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// mapextra holds fields that are not present on all maps.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>type mapextra struct {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// If both key and elem do not contain pointers and are inline, then we mark bucket</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// type as containing no pointers. This avoids scanning such maps.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// However, bmap.overflow is a pointer. In order to keep overflow buckets</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// overflow and oldoverflow are only used if key and elem do not contain pointers.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// overflow contains overflow buckets for hmap.buckets.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// oldoverflow contains overflow buckets for hmap.oldbuckets.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">// The indirection allows to store a pointer to the slice in hiter.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	overflow    *[]*bmap
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	oldoverflow *[]*bmap
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// nextOverflow holds a pointer to a free overflow bucket.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	nextOverflow *bmap
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// A bucket for a Go map.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>type bmap struct {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// tophash generally contains the top byte of the hash value</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// for each key in this bucket. If tophash[0] &lt; minTopHash,</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// tophash[0] is a bucket evacuation state instead.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	tophash [bucketCnt]uint8
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Followed by bucketCnt keys and then bucketCnt elems.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// NOTE: packing all the keys together and then all the elems together makes the</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// code a bit more complicated than alternating key/elem/key/elem/... but it allows</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// us to eliminate padding which would be needed for, e.g., map[int64]int8.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// Followed by an overflow pointer.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span><span class="comment">// A hash iteration structure.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span><span class="comment">// If you modify hiter, also change cmd/compile/internal/reflectdata/reflect.go</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span><span class="comment">// and reflect/value.go to match the layout of this structure.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>type hiter struct {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	key         unsafe.Pointer <span class="comment">// Must be in first position.  Write nil to indicate iteration end (see cmd/compile/internal/walk/range.go).</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	elem        unsafe.Pointer <span class="comment">// Must be in second position (see cmd/compile/internal/walk/range.go).</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	t           *maptype
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	h           *hmap
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	buckets     unsafe.Pointer <span class="comment">// bucket ptr at hash_iter initialization time</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	bptr        *bmap          <span class="comment">// current bucket</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	overflow    *[]*bmap       <span class="comment">// keeps overflow buckets of hmap.buckets alive</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	oldoverflow *[]*bmap       <span class="comment">// keeps overflow buckets of hmap.oldbuckets alive</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	startBucket uintptr        <span class="comment">// bucket iteration started at</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	offset      uint8          <span class="comment">// intra-bucket offset to start from during iteration (should be big enough to hold bucketCnt-1)</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	wrapped     bool           <span class="comment">// already wrapped around from end of bucket array to beginning</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	B           uint8
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	i           uint8
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	bucket      uintptr
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	checkBucket uintptr
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span><span class="comment">// bucketShift returns 1&lt;&lt;b, optimized for code generation.</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>func bucketShift(b uint8) uintptr {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// Masking the shift amount allows overflow checks to be elided.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	return uintptr(1) &lt;&lt; (b &amp; (goarch.PtrSize*8 - 1))
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// bucketMask returns 1&lt;&lt;b - 1, optimized for code generation.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>func bucketMask(b uint8) uintptr {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	return bucketShift(b) - 1
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span><span class="comment">// tophash calculates the tophash value for hash.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>func tophash(hash uintptr) uint8 {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	top := uint8(hash &gt;&gt; (goarch.PtrSize*8 - 8))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if top &lt; minTopHash {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		top += minTopHash
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	return top
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func evacuated(b *bmap) bool {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	h := b.tophash[0]
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	return h &gt; emptyOne &amp;&amp; h &lt; minTopHash
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>func (b *bmap) overflow(t *maptype) *bmap {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	return *(**bmap)(add(unsafe.Pointer(b), uintptr(t.BucketSize)-goarch.PtrSize))
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func (b *bmap) setoverflow(t *maptype, ovf *bmap) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	*(**bmap)(add(unsafe.Pointer(b), uintptr(t.BucketSize)-goarch.PtrSize)) = ovf
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>func (b *bmap) keys() unsafe.Pointer {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	return add(unsafe.Pointer(b), dataOffset)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// incrnoverflow increments h.noverflow.</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// noverflow counts the number of overflow buckets.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// This is used to trigger same-size map growth.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// See also tooManyOverflowBuckets.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span><span class="comment">// To keep hmap small, noverflow is a uint16.</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// When there are few buckets, noverflow is an exact count.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// When there are many buckets, noverflow is an approximate count.</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>func (h *hmap) incrnoverflow() {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	<span class="comment">// We trigger same-size map growth if there are</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// as many overflow buckets as buckets.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// We need to be able to count to 1&lt;&lt;h.B.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	if h.B &lt; 16 {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		h.noverflow++
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		return
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// Increment with probability 1/(1&lt;&lt;(h.B-15)).</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// When we reach 1&lt;&lt;15 - 1, we will have approximately</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// as many overflow buckets as buckets.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	mask := uint32(1)&lt;&lt;(h.B-15) - 1
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">// Example: if h.B == 18, then mask == 7,</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">// and rand() &amp; 7 == 0 with probability 1/8.</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if uint32(rand())&amp;mask == 0 {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		h.noverflow++
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>func (h *hmap) newoverflow(t *maptype, b *bmap) *bmap {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	var ovf *bmap
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	if h.extra != nil &amp;&amp; h.extra.nextOverflow != nil {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		<span class="comment">// We have preallocated overflow buckets available.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		<span class="comment">// See makeBucketArray for more details.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		ovf = h.extra.nextOverflow
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		if ovf.overflow(t) == nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			<span class="comment">// We&#39;re not at the end of the preallocated overflow buckets. Bump the pointer.</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			h.extra.nextOverflow = (*bmap)(add(unsafe.Pointer(ovf), uintptr(t.BucketSize)))
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		} else {
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			<span class="comment">// This is the last preallocated overflow bucket.</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			<span class="comment">// Reset the overflow pointer on this bucket,</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			<span class="comment">// which was set to a non-nil sentinel value.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			ovf.setoverflow(t, nil)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			h.extra.nextOverflow = nil
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	} else {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		ovf = (*bmap)(newobject(t.Bucket))
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	h.incrnoverflow()
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	if t.Bucket.PtrBytes == 0 {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		h.createOverflow()
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		*h.extra.overflow = append(*h.extra.overflow, ovf)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	b.setoverflow(t, ovf)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	return ovf
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>func (h *hmap) createOverflow() {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	if h.extra == nil {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		h.extra = new(mapextra)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	if h.extra.overflow == nil {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		h.extra.overflow = new([]*bmap)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>}
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>func makemap64(t *maptype, hint int64, h *hmap) *hmap {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	if int64(int(hint)) != hint {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		hint = 0
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	return makemap(t, int(hint), h)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// makemap_small implements Go map creation for make(map[k]v) and</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// make(map[k]v, hint) when hint is known to be at most bucketCnt</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">// at compile time and the map needs to be allocated on the heap.</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>func makemap_small() *hmap {
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	h := new(hmap)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	h.hash0 = uint32(rand())
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	return h
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// makemap implements Go map creation for make(map[k]v, hint).</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// If the compiler has determined that the map or the first bucket</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">// can be created on the stack, h and/or bucket may be non-nil.</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">// If h != nil, the map can be created directly in h.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">// If h.buckets != nil, bucket pointed to can be used as the first bucket.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>func makemap(t *maptype, hint int, h *hmap) *hmap {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	mem, overflow := math.MulUintptr(uintptr(hint), t.Bucket.Size_)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if overflow || mem &gt; maxAlloc {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		hint = 0
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// initialize Hmap</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	if h == nil {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		h = new(hmap)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	h.hash0 = uint32(rand())
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// Find the size parameter B which will hold the requested # of elements.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// For hint &lt; 0 overLoadFactor returns false since hint &lt; bucketCnt.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	B := uint8(0)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	for overLoadFactor(hint, B) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		B++
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	h.B = B
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// allocate initial hash table</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// if B == 0, the buckets field is allocated lazily later (in mapassign)</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// If hint is large zeroing this memory could take a while.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if h.B != 0 {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		var nextOverflow *bmap
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		h.buckets, nextOverflow = makeBucketArray(t, h.B, nil)
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		if nextOverflow != nil {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			h.extra = new(mapextra)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			h.extra.nextOverflow = nextOverflow
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	return h
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// makeBucketArray initializes a backing array for map buckets.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">// 1&lt;&lt;b is the minimum number of buckets to allocate.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// dirtyalloc should either be nil or a bucket array previously</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// allocated by makeBucketArray with the same t and b parameters.</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">// If dirtyalloc is nil a new backing array will be alloced and</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">// otherwise dirtyalloc will be cleared and reused as backing array.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>func makeBucketArray(t *maptype, b uint8, dirtyalloc unsafe.Pointer) (buckets unsafe.Pointer, nextOverflow *bmap) {
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	base := bucketShift(b)
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	nbuckets := base
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	<span class="comment">// For small b, overflow buckets are unlikely.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	<span class="comment">// Avoid the overhead of the calculation.</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	if b &gt;= 4 {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		<span class="comment">// Add on the estimated number of overflow buckets</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		<span class="comment">// required to insert the median number of elements</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		<span class="comment">// used with this value of b.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		nbuckets += bucketShift(b - 4)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		sz := t.Bucket.Size_ * nbuckets
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		up := roundupsize(sz, t.Bucket.PtrBytes == 0)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		if up != sz {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			nbuckets = up / t.Bucket.Size_
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	if dirtyalloc == nil {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		buckets = newarray(t.Bucket, int(nbuckets))
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	} else {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		<span class="comment">// dirtyalloc was previously generated by</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		<span class="comment">// the above newarray(t.Bucket, int(nbuckets))</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		<span class="comment">// but may not be empty.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		buckets = dirtyalloc
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		size := t.Bucket.Size_ * nbuckets
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		if t.Bucket.PtrBytes != 0 {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			memclrHasPointers(buckets, size)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		} else {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			memclrNoHeapPointers(buckets, size)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	if base != nbuckets {
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// We preallocated some overflow buckets.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		<span class="comment">// To keep the overhead of tracking these overflow buckets to a minimum,</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		<span class="comment">// we use the convention that if a preallocated overflow bucket&#39;s overflow</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		<span class="comment">// pointer is nil, then there are more available by bumping the pointer.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		<span class="comment">// We need a safe non-nil pointer for the last overflow bucket; just use buckets.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		nextOverflow = (*bmap)(add(buckets, base*uintptr(t.BucketSize)))
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		last := (*bmap)(add(buckets, (nbuckets-1)*uintptr(t.BucketSize)))
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		last.setoverflow(t, (*bmap)(buckets))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	return buckets, nextOverflow
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// mapaccess1 returns a pointer to h[key].  Never returns nil, instead</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// it will return a reference to the zero object for the elem type if</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// the key is not in the map.</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// NOTE: The returned pointer may keep the whole map live, so don&#39;t</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// hold onto it for very long.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; h != nil {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(mapaccess1)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, pc)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		raceReadObjectPC(t.Key, key, callerpc, pc)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	if msanenabled &amp;&amp; h != nil {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		msanread(key, t.Key.Size_)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if asanenabled &amp;&amp; h != nil {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		asanread(key, t.Key.Size_)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		if err := mapKeyError(t, key); err != nil {
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			panic(err) <span class="comment">// see issue 23734</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		return unsafe.Pointer(&amp;zeroVal[0])
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		fatal(&#34;concurrent map read and map write&#34;)
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	hash := t.Hasher(key, uintptr(h.hash0))
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	m := bucketMask(h.B)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	b := (*bmap)(add(h.buckets, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	if c := h.oldbuckets; c != nil {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		if !h.sameSizeGrow() {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			<span class="comment">// There used to be half as many buckets; mask down one more power of two.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			m &gt;&gt;= 1
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		oldb := (*bmap)(add(c, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		if !evacuated(oldb) {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			b = oldb
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	top := tophash(hash)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>bucketloop:
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	for ; b != nil; b = b.overflow(t) {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			if b.tophash[i] != top {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>				if b.tophash[i] == emptyRest {
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>					break bucketloop
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>				}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>				continue
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>				k = *((*unsafe.Pointer)(k))
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			if t.Key.Equal(key, k) {
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>				e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>				if t.IndirectElem() {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>					e = *((*unsafe.Pointer)(e))
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>				}
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>				return e
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	return unsafe.Pointer(&amp;zeroVal[0])
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>func mapaccess2(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, bool) {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; h != nil {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(mapaccess2)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, pc)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		raceReadObjectPC(t.Key, key, callerpc, pc)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	}
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	if msanenabled &amp;&amp; h != nil {
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		msanread(key, t.Key.Size_)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	if asanenabled &amp;&amp; h != nil {
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		asanread(key, t.Key.Size_)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		if err := mapKeyError(t, key); err != nil {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			panic(err) <span class="comment">// see issue 23734</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>		return unsafe.Pointer(&amp;zeroVal[0]), false
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	}
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		fatal(&#34;concurrent map read and map write&#34;)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	hash := t.Hasher(key, uintptr(h.hash0))
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	m := bucketMask(h.B)
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	b := (*bmap)(add(h.buckets, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	if c := h.oldbuckets; c != nil {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		if !h.sameSizeGrow() {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// There used to be half as many buckets; mask down one more power of two.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			m &gt;&gt;= 1
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		oldb := (*bmap)(add(c, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		if !evacuated(oldb) {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			b = oldb
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	top := tophash(hash)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>bucketloop:
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	for ; b != nil; b = b.overflow(t) {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			if b.tophash[i] != top {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				if b.tophash[i] == emptyRest {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>					break bucketloop
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>				}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>				continue
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>				k = *((*unsafe.Pointer)(k))
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			if t.Key.Equal(key, k) {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>				if t.IndirectElem() {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>					e = *((*unsafe.Pointer)(e))
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>				}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>				return e, true
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	return unsafe.Pointer(&amp;zeroVal[0]), false
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">// returns both key and elem. Used by map iterator.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>func mapaccessK(t *maptype, h *hmap, key unsafe.Pointer) (unsafe.Pointer, unsafe.Pointer) {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		return nil, nil
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	hash := t.Hasher(key, uintptr(h.hash0))
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	m := bucketMask(h.B)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	b := (*bmap)(add(h.buckets, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	if c := h.oldbuckets; c != nil {
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>		if !h.sameSizeGrow() {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>			<span class="comment">// There used to be half as many buckets; mask down one more power of two.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>			m &gt;&gt;= 1
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		}
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		oldb := (*bmap)(add(c, (hash&amp;m)*uintptr(t.BucketSize)))
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		if !evacuated(oldb) {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>			b = oldb
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	}
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	top := tophash(hash)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>bucketloop:
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	for ; b != nil; b = b.overflow(t) {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			if b.tophash[i] != top {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>				if b.tophash[i] == emptyRest {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>					break bucketloop
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>				}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>				continue
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>				k = *((*unsafe.Pointer)(k))
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>			}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>			if t.Key.Equal(key, k) {
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>				e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>				if t.IndirectElem() {
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>					e = *((*unsafe.Pointer)(e))
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>				}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>				return k, e
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>			}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	return nil, nil
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>func mapaccess1_fat(t *maptype, h *hmap, key, zero unsafe.Pointer) unsafe.Pointer {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	e := mapaccess1(t, h, key)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	if e == unsafe.Pointer(&amp;zeroVal[0]) {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		return zero
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	return e
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>func mapaccess2_fat(t *maptype, h *hmap, key, zero unsafe.Pointer) (unsafe.Pointer, bool) {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	e := mapaccess1(t, h, key)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if e == unsafe.Pointer(&amp;zeroVal[0]) {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		return zero, false
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	return e, true
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">// Like mapaccess, but allocates a slot for the key if it is not present in the map.</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	if h == nil {
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		panic(plainError(&#34;assignment to entry in nil map&#34;))
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	}
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if raceenabled {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(mapassign)
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		racewritepc(unsafe.Pointer(h), callerpc, pc)
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		raceReadObjectPC(t.Key, key, callerpc, pc)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	if msanenabled {
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>		msanread(key, t.Key.Size_)
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	if asanenabled {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		asanread(key, t.Key.Size_)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	hash := t.Hasher(key, uintptr(h.hash0))
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// Set hashWriting after calling t.hasher, since t.hasher may panic,</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	<span class="comment">// in which case we have not actually done a write.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	h.flags ^= hashWriting
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	if h.buckets == nil {
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		h.buckets = newobject(t.Bucket) <span class="comment">// newarray(t.Bucket, 1)</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>again:
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	bucket := hash &amp; bucketMask(h.B)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	if h.growing() {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		growWork(t, h, bucket)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	b := (*bmap)(add(h.buckets, bucket*uintptr(t.BucketSize)))
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	top := tophash(hash)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	var inserti *uint8
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	var insertk unsafe.Pointer
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	var elem unsafe.Pointer
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>bucketloop:
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	for {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>			if b.tophash[i] != top {
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>				if isEmpty(b.tophash[i]) &amp;&amp; inserti == nil {
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>					inserti = &amp;b.tophash[i]
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>					insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>					elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>				}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>				if b.tophash[i] == emptyRest {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>					break bucketloop
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>				}
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				continue
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>				k = *((*unsafe.Pointer)(k))
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>			}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			if !t.Key.Equal(key, k) {
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>				continue
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>			}
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			<span class="comment">// already have a mapping for key. Update it.</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>			if t.NeedKeyUpdate() {
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>				typedmemmove(t.Key, k, key)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>			}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			goto done
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		ovf := b.overflow(t)
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		if ovf == nil {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>			break
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		}
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		b = ovf
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	}
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	<span class="comment">// Did not find mapping for key. Allocate new cell &amp; add entry.</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// If we hit the max load factor or we have too many overflow buckets,</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	<span class="comment">// and we&#39;re not already in the middle of growing, start growing.</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	if !h.growing() &amp;&amp; (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		hashGrow(t, h)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>		goto again <span class="comment">// Growing the table invalidates everything, so try again</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	}
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	if inserti == nil {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		<span class="comment">// The current bucket and all the overflow buckets connected to it are full, allocate a new one.</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		newb := h.newoverflow(t, b)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		inserti = &amp;newb.tophash[0]
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		insertk = add(unsafe.Pointer(newb), dataOffset)
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		elem = add(insertk, bucketCnt*uintptr(t.KeySize))
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// store new key/elem at insert position</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	if t.IndirectKey() {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		kmem := newobject(t.Key)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		*(*unsafe.Pointer)(insertk) = kmem
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		insertk = kmem
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	if t.IndirectElem() {
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		vmem := newobject(t.Elem)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		*(*unsafe.Pointer)(elem) = vmem
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	}
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	typedmemmove(t.Key, insertk, key)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	*inserti = top
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	h.count++
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>done:
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting == 0 {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	h.flags &amp;^= hashWriting
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	if t.IndirectElem() {
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		elem = *((*unsafe.Pointer)(elem))
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	return elem
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>}
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>func mapdelete(t *maptype, h *hmap, key unsafe.Pointer) {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; h != nil {
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(mapdelete)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		racewritepc(unsafe.Pointer(h), callerpc, pc)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		raceReadObjectPC(t.Key, key, callerpc, pc)
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	if msanenabled &amp;&amp; h != nil {
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		msanread(key, t.Key.Size_)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>	}
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	if asanenabled &amp;&amp; h != nil {
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		asanread(key, t.Key.Size_)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>		if err := mapKeyError(t, key); err != nil {
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			panic(err) <span class="comment">// see issue 23734</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>		}
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		return
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	}
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	hash := t.Hasher(key, uintptr(h.hash0))
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	<span class="comment">// Set hashWriting after calling t.hasher, since t.hasher may panic,</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	<span class="comment">// in which case we have not actually done a write (delete).</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	h.flags ^= hashWriting
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	bucket := hash &amp; bucketMask(h.B)
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	if h.growing() {
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		growWork(t, h, bucket)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	b := (*bmap)(add(h.buckets, bucket*uintptr(t.BucketSize)))
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	bOrig := b
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	top := tophash(hash)
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>search:
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	for ; b != nil; b = b.overflow(t) {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			if b.tophash[i] != top {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>				if b.tophash[i] == emptyRest {
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>					break search
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>				}
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>				continue
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			}
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.KeySize))
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			k2 := k
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>				k2 = *((*unsafe.Pointer)(k2))
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>			}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			if !t.Key.Equal(key, k2) {
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>				continue
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>			<span class="comment">// Only clear key if there are pointers in it.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>				*(*unsafe.Pointer)(k) = nil
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			} else if t.Key.PtrBytes != 0 {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>				memclrHasPointers(k, t.Key.Size_)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>			}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			if t.IndirectElem() {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>				*(*unsafe.Pointer)(e) = nil
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			} else if t.Elem.PtrBytes != 0 {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>				memclrHasPointers(e, t.Elem.Size_)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			} else {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>				memclrNoHeapPointers(e, t.Elem.Size_)
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			}
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			b.tophash[i] = emptyOne
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			<span class="comment">// If the bucket now ends in a bunch of emptyOne states,</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>			<span class="comment">// change those to emptyRest states.</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>			<span class="comment">// It would be nice to make this a separate function, but</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			<span class="comment">// for loops are not currently inlineable.</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			if i == bucketCnt-1 {
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>				if b.overflow(t) != nil &amp;&amp; b.overflow(t).tophash[0] != emptyRest {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>					goto notLast
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>				}
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			} else {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				if b.tophash[i+1] != emptyRest {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>					goto notLast
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>				}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			for {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>				b.tophash[i] = emptyRest
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				if i == 0 {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>					if b == bOrig {
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>						break <span class="comment">// beginning of initial bucket, we&#39;re done.</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>					}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>					<span class="comment">// Find previous bucket, continue at its last entry.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>					c := b
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>					for b = bOrig; b.overflow(t) != c; b = b.overflow(t) {
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>					}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>					i = bucketCnt - 1
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>				} else {
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>					i--
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>				}
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>				if b.tophash[i] != emptyOne {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>					break
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>				}
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		notLast:
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>			h.count--
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			<span class="comment">// Reset the hash seed to make it more difficult for attackers to</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			<span class="comment">// repeatedly trigger hash collisions. See issue 25237.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			if h.count == 0 {
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>				h.hash0 = uint32(rand())
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			}
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>			break search
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		}
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	}
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting == 0 {
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	}
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	h.flags &amp;^= hashWriting
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>}
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span><span class="comment">// mapiterinit initializes the hiter struct used for ranging over maps.</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span><span class="comment">// The hiter struct pointed to by &#39;it&#39; is allocated on the stack</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span><span class="comment">// by the compilers order pass or on the heap by reflect_mapiterinit.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span><span class="comment">// Both need to have zeroed hiter since the struct contains pointers.</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>func mapiterinit(t *maptype, h *hmap, it *hiter) {
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; h != nil {
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(mapiterinit))
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	}
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	it.t = t
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		return
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>	}
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	if unsafe.Sizeof(hiter{})/goarch.PtrSize != 12 {
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		throw(&#34;hash_iter size incorrect&#34;) <span class="comment">// see cmd/compile/internal/reflectdata/reflect.go</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>	}
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	it.h = h
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	<span class="comment">// grab snapshot of bucket state</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	it.B = h.B
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>	it.buckets = h.buckets
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>	if t.Bucket.PtrBytes == 0 {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		<span class="comment">// Allocate the current slice and remember pointers to both current and old.</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		<span class="comment">// This preserves all relevant overflow buckets alive even if</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		<span class="comment">// the table grows and/or overflow buckets are added to the table</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		<span class="comment">// while we are iterating.</span>
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		h.createOverflow()
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		it.overflow = h.extra.overflow
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>		it.oldoverflow = h.extra.oldoverflow
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>	}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	<span class="comment">// decide where to start</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	r := uintptr(rand())
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	it.startBucket = r &amp; bucketMask(h.B)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	it.offset = uint8(r &gt;&gt; h.B &amp; (bucketCnt - 1))
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	<span class="comment">// iterator state</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	it.bucket = it.startBucket
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	<span class="comment">// Remember we have an iterator.</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>	<span class="comment">// Can run concurrently with another mapiterinit().</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	if old := h.flags; old&amp;(iterator|oldIterator) != iterator|oldIterator {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		atomic.Or8(&amp;h.flags, iterator|oldIterator)
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	}
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	mapiternext(it)
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>}
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>func mapiternext(it *hiter) {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	h := it.h
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	if raceenabled {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(mapiternext))
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	}
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		fatal(&#34;concurrent map iteration and map write&#34;)
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	t := it.t
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	bucket := it.bucket
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	b := it.bptr
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	i := it.i
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	checkBucket := it.checkBucket
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>next:
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	if b == nil {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		if bucket == it.startBucket &amp;&amp; it.wrapped {
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			<span class="comment">// end of iteration</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			it.key = nil
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>			it.elem = nil
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>			return
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		}
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		if h.growing() &amp;&amp; it.B == h.B {
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			<span class="comment">// Iterator was started in the middle of a grow, and the grow isn&#39;t done yet.</span>
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>			<span class="comment">// If the bucket we&#39;re looking at hasn&#39;t been filled in yet (i.e. the old</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			<span class="comment">// bucket hasn&#39;t been evacuated) then we need to iterate through the old</span>
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			<span class="comment">// bucket and only return the ones that will be migrated to this bucket.</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			oldbucket := bucket &amp; it.h.oldbucketmask()
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			b = (*bmap)(add(h.oldbuckets, oldbucket*uintptr(t.BucketSize)))
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			if !evacuated(b) {
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>				checkBucket = bucket
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>			} else {
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>				b = (*bmap)(add(it.buckets, bucket*uintptr(t.BucketSize)))
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>				checkBucket = noCheck
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>			}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		} else {
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>			b = (*bmap)(add(it.buckets, bucket*uintptr(t.BucketSize)))
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>			checkBucket = noCheck
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		}
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>		bucket++
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		if bucket == bucketShift(it.B) {
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>			bucket = 0
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>			it.wrapped = true
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>		}
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>		i = 0
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	for ; i &lt; bucketCnt; i++ {
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>		offi := (i + it.offset) &amp; (bucketCnt - 1)
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>		if isEmpty(b.tophash[offi]) || b.tophash[offi] == evacuatedEmpty {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>			<span class="comment">// TODO: emptyRest is hard to use here, as we start iterating</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>			<span class="comment">// in the middle of a bucket. It&#39;s feasible, just tricky.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>			continue
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		}
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>		k := add(unsafe.Pointer(b), dataOffset+uintptr(offi)*uintptr(t.KeySize))
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		if t.IndirectKey() {
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>			k = *((*unsafe.Pointer)(k))
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>		}
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>		e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+uintptr(offi)*uintptr(t.ValueSize))
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>		if checkBucket != noCheck &amp;&amp; !h.sameSizeGrow() {
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>			<span class="comment">// Special case: iterator was started during a grow to a larger size</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>			<span class="comment">// and the grow is not done yet. We&#39;re working on a bucket whose</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>			<span class="comment">// oldbucket has not been evacuated yet. Or at least, it wasn&#39;t</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>			<span class="comment">// evacuated when we started the bucket. So we&#39;re iterating</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>			<span class="comment">// through the oldbucket, skipping any keys that will go</span>
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>			<span class="comment">// to the other new bucket (each oldbucket expands to two</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>			<span class="comment">// buckets during a grow).</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>			if t.ReflexiveKey() || t.Key.Equal(k, k) {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>				<span class="comment">// If the item in the oldbucket is not destined for</span>
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>				<span class="comment">// the current new bucket in the iteration, skip it.</span>
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>				hash := t.Hasher(k, uintptr(h.hash0))
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>				if hash&amp;bucketMask(it.B) != checkBucket {
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>					continue
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>				}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>			} else {
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>				<span class="comment">// Hash isn&#39;t repeatable if k != k (NaNs).  We need a</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>				<span class="comment">// repeatable and randomish choice of which direction</span>
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>				<span class="comment">// to send NaNs during evacuation. We&#39;ll use the low</span>
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>				<span class="comment">// bit of tophash to decide which way NaNs go.</span>
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>				<span class="comment">// NOTE: this case is why we need two evacuate tophash</span>
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>				<span class="comment">// values, evacuatedX and evacuatedY, that differ in</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>				<span class="comment">// their low bit.</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>				if checkBucket&gt;&gt;(it.B-1) != uintptr(b.tophash[offi]&amp;1) {
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>					continue
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>				}
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>			}
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		}
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		if (b.tophash[offi] != evacuatedX &amp;&amp; b.tophash[offi] != evacuatedY) ||
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			!(t.ReflexiveKey() || t.Key.Equal(k, k)) {
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			<span class="comment">// This is the golden data, we can return it.</span>
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			<span class="comment">// OR</span>
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			<span class="comment">// key!=key, so the entry can&#39;t be deleted or updated, so we can just return it.</span>
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			<span class="comment">// That&#39;s lucky for us because when key!=key we can&#39;t look it up successfully.</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>			it.key = k
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>			if t.IndirectElem() {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>				e = *((*unsafe.Pointer)(e))
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>			}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>			it.elem = e
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		} else {
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>			<span class="comment">// The hash table has grown since the iterator was started.</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>			<span class="comment">// The golden data for this key is now somewhere else.</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>			<span class="comment">// Check the current hash table for the data.</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>			<span class="comment">// This code handles the case where the key</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>			<span class="comment">// has been deleted, updated, or deleted and reinserted.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>			<span class="comment">// NOTE: we need to regrab the key as it has potentially been</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>			<span class="comment">// updated to an equal() but not identical key (e.g. +0.0 vs -0.0).</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			rk, re := mapaccessK(t, h, k)
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			if rk == nil {
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>				continue <span class="comment">// key has been deleted</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>			}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>			it.key = rk
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>			it.elem = re
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>		}
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>		it.bucket = bucket
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>		if it.bptr != b { <span class="comment">// avoid unnecessary write barrier; see issue 14921</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>			it.bptr = b
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>		}
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		it.i = i + 1
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>		it.checkBucket = checkBucket
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>		return
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	}
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	b = b.overflow(t)
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>	i = 0
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	goto next
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>}
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span><span class="comment">// mapclear deletes all keys from a map.</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>func mapclear(t *maptype, h *hmap) {
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	if raceenabled &amp;&amp; h != nil {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		pc := abi.FuncPCABIInternal(mapclear)
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>		racewritepc(unsafe.Pointer(h), callerpc, pc)
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	}
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		return
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting != 0 {
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>	}
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	h.flags ^= hashWriting
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	<span class="comment">// Mark buckets empty, so existing iterators can be terminated, see issue #59411.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	markBucketsEmpty := func(bucket unsafe.Pointer, mask uintptr) {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt;= mask; i++ {
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>			b := (*bmap)(add(bucket, i*uintptr(t.BucketSize)))
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>			for ; b != nil; b = b.overflow(t) {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>				for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>					b.tophash[i] = emptyRest
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>				}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>			}
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>		}
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	}
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>	markBucketsEmpty(h.buckets, bucketMask(h.B))
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	if oldBuckets := h.oldbuckets; oldBuckets != nil {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		markBucketsEmpty(oldBuckets, h.oldbucketmask())
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>	}
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>	h.flags &amp;^= sameSizeGrow
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	h.oldbuckets = nil
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	h.nevacuate = 0
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>	h.noverflow = 0
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>	h.count = 0
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	<span class="comment">// Reset the hash seed to make it more difficult for attackers to</span>
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	<span class="comment">// repeatedly trigger hash collisions. See issue 25237.</span>
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>	h.hash0 = uint32(rand())
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>	<span class="comment">// Keep the mapextra allocation but clear any extra information.</span>
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	if h.extra != nil {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>		*h.extra = mapextra{}
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	}
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>	<span class="comment">// makeBucketArray clears the memory pointed to by h.buckets</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	<span class="comment">// and recovers any overflow buckets by generating them</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	<span class="comment">// as if h.buckets was newly alloced.</span>
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>	_, nextOverflow := makeBucketArray(t, h.B, h.buckets)
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>	if nextOverflow != nil {
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		<span class="comment">// If overflow buckets are created then h.extra</span>
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		<span class="comment">// will have been allocated during initial bucket creation.</span>
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		h.extra.nextOverflow = nextOverflow
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>	}
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>	if h.flags&amp;hashWriting == 0 {
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		fatal(&#34;concurrent map writes&#34;)
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>	}
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>	h.flags &amp;^= hashWriting
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>}
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>func hashGrow(t *maptype, h *hmap) {
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	<span class="comment">// If we&#39;ve hit the load factor, get bigger.</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>	<span class="comment">// Otherwise, there are too many overflow buckets,</span>
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>	<span class="comment">// so keep the same number of buckets and &#34;grow&#34; laterally.</span>
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>	bigger := uint8(1)
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>	if !overLoadFactor(h.count+1, h.B) {
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>		bigger = 0
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>		h.flags |= sameSizeGrow
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>	}
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>	oldbuckets := h.buckets
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>	newbuckets, nextOverflow := makeBucketArray(t, h.B+bigger, nil)
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>	flags := h.flags &amp;^ (iterator | oldIterator)
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>	if h.flags&amp;iterator != 0 {
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		flags |= oldIterator
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>	}
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>	<span class="comment">// commit the grow (atomic wrt gc)</span>
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>	h.B += bigger
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>	h.flags = flags
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>	h.oldbuckets = oldbuckets
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>	h.buckets = newbuckets
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>	h.nevacuate = 0
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>	h.noverflow = 0
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>	if h.extra != nil &amp;&amp; h.extra.overflow != nil {
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>		<span class="comment">// Promote current overflow buckets to the old generation.</span>
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>		if h.extra.oldoverflow != nil {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>			throw(&#34;oldoverflow is not nil&#34;)
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>		}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>		h.extra.oldoverflow = h.extra.overflow
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>		h.extra.overflow = nil
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>	}
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>	if nextOverflow != nil {
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>		if h.extra == nil {
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>			h.extra = new(mapextra)
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>		}
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>		h.extra.nextOverflow = nextOverflow
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>	}
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>	<span class="comment">// the actual copying of the hash table data is done incrementally</span>
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>	<span class="comment">// by growWork() and evacuate().</span>
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>}
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span><span class="comment">// overLoadFactor reports whether count items placed in 1&lt;&lt;B buckets is over loadFactor.</span>
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>func overLoadFactor(count int, B uint8) bool {
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>	return count &gt; bucketCnt &amp;&amp; uintptr(count) &gt; loadFactorNum*(bucketShift(B)/loadFactorDen)
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>}
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span><span class="comment">// tooManyOverflowBuckets reports whether noverflow buckets is too many for a map with 1&lt;&lt;B buckets.</span>
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span><span class="comment">// Note that most of these overflow buckets must be in sparse use;</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span><span class="comment">// if use was dense, then we&#39;d have already triggered regular map growth.</span>
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>func tooManyOverflowBuckets(noverflow uint16, B uint8) bool {
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>	<span class="comment">// If the threshold is too low, we do extraneous work.</span>
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>	<span class="comment">// If the threshold is too high, maps that grow and shrink can hold on to lots of unused memory.</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>	<span class="comment">// &#34;too many&#34; means (approximately) as many overflow buckets as regular buckets.</span>
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>	<span class="comment">// See incrnoverflow for more details.</span>
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>	if B &gt; 15 {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		B = 15
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>	}
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>	<span class="comment">// The compiler doesn&#39;t see here that B &lt; 16; mask B to generate shorter shift code.</span>
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>	return noverflow &gt;= uint16(1)&lt;&lt;(B&amp;15)
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>}
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span><span class="comment">// growing reports whether h is growing. The growth may be to the same size or bigger.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>func (h *hmap) growing() bool {
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>	return h.oldbuckets != nil
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>}
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span><span class="comment">// sameSizeGrow reports whether the current growth is to a map of the same size.</span>
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>func (h *hmap) sameSizeGrow() bool {
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>	return h.flags&amp;sameSizeGrow != 0
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>}
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span><span class="comment">// noldbuckets calculates the number of buckets prior to the current map growth.</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>func (h *hmap) noldbuckets() uintptr {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>	oldB := h.B
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>	if !h.sameSizeGrow() {
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>		oldB--
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>	}
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>	return bucketShift(oldB)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>}
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span><span class="comment">// oldbucketmask provides a mask that can be applied to calculate n % noldbuckets().</span>
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>func (h *hmap) oldbucketmask() uintptr {
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>	return h.noldbuckets() - 1
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>}
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>func growWork(t *maptype, h *hmap, bucket uintptr) {
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>	<span class="comment">// make sure we evacuate the oldbucket corresponding</span>
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>	<span class="comment">// to the bucket we&#39;re about to use</span>
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>	evacuate(t, h, bucket&amp;h.oldbucketmask())
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>	<span class="comment">// evacuate one more oldbucket to make progress on growing</span>
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>	if h.growing() {
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>		evacuate(t, h, h.nevacuate)
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>	}
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>}
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>func bucketEvacuated(t *maptype, h *hmap, bucket uintptr) bool {
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>	b := (*bmap)(add(h.oldbuckets, bucket*uintptr(t.BucketSize)))
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>	return evacuated(b)
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>}
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span><span class="comment">// evacDst is an evacuation destination.</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>type evacDst struct {
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>	b *bmap          <span class="comment">// current destination bucket</span>
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>	i int            <span class="comment">// key/elem index into b</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>	k unsafe.Pointer <span class="comment">// pointer to current key storage</span>
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>	e unsafe.Pointer <span class="comment">// pointer to current elem storage</span>
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>func evacuate(t *maptype, h *hmap, oldbucket uintptr) {
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>	b := (*bmap)(add(h.oldbuckets, oldbucket*uintptr(t.BucketSize)))
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>	newbit := h.noldbuckets()
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>	if !evacuated(b) {
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>		<span class="comment">// TODO: reuse overflow buckets instead of using new ones, if there</span>
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>		<span class="comment">// is no iterator using the old buckets.  (If !oldIterator.)</span>
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>		<span class="comment">// xy contains the x and y (low and high) evacuation destinations.</span>
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>		var xy [2]evacDst
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>		x := &amp;xy[0]
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>		x.b = (*bmap)(add(h.buckets, oldbucket*uintptr(t.BucketSize)))
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>		x.k = add(unsafe.Pointer(x.b), dataOffset)
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>		x.e = add(x.k, bucketCnt*uintptr(t.KeySize))
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>		if !h.sameSizeGrow() {
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>			<span class="comment">// Only calculate y pointers if we&#39;re growing bigger.</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>			<span class="comment">// Otherwise GC can see bad pointers.</span>
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>			y := &amp;xy[1]
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>			y.b = (*bmap)(add(h.buckets, (oldbucket+newbit)*uintptr(t.BucketSize)))
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>			y.k = add(unsafe.Pointer(y.b), dataOffset)
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>			y.e = add(y.k, bucketCnt*uintptr(t.KeySize))
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>		}
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>		for ; b != nil; b = b.overflow(t) {
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset)
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>			e := add(k, bucketCnt*uintptr(t.KeySize))
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>			for i := 0; i &lt; bucketCnt; i, k, e = i+1, add(k, uintptr(t.KeySize)), add(e, uintptr(t.ValueSize)) {
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>				top := b.tophash[i]
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>				if isEmpty(top) {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>					b.tophash[i] = evacuatedEmpty
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>					continue
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>				}
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>				if top &lt; minTopHash {
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>					throw(&#34;bad map state&#34;)
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>				}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>				k2 := k
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>				if t.IndirectKey() {
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>					k2 = *((*unsafe.Pointer)(k2))
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>				}
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>				var useY uint8
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>				if !h.sameSizeGrow() {
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>					<span class="comment">// Compute hash to make our evacuation decision (whether we need</span>
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>					<span class="comment">// to send this key/elem to bucket x or bucket y).</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>					hash := t.Hasher(k2, uintptr(h.hash0))
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>					if h.flags&amp;iterator != 0 &amp;&amp; !t.ReflexiveKey() &amp;&amp; !t.Key.Equal(k2, k2) {
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>						<span class="comment">// If key != key (NaNs), then the hash could be (and probably</span>
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>						<span class="comment">// will be) entirely different from the old hash. Moreover,</span>
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>						<span class="comment">// it isn&#39;t reproducible. Reproducibility is required in the</span>
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>						<span class="comment">// presence of iterators, as our evacuation decision must</span>
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>						<span class="comment">// match whatever decision the iterator made.</span>
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>						<span class="comment">// Fortunately, we have the freedom to send these keys either</span>
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>						<span class="comment">// way. Also, tophash is meaningless for these kinds of keys.</span>
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>						<span class="comment">// We let the low bit of tophash drive the evacuation decision.</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>						<span class="comment">// We recompute a new random tophash for the next level so</span>
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>						<span class="comment">// these keys will get evenly distributed across all buckets</span>
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>						<span class="comment">// after multiple grows.</span>
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>						useY = top &amp; 1
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>						top = tophash(hash)
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>					} else {
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>						if hash&amp;newbit != 0 {
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>							useY = 1
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>						}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>					}
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>				}
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>				if evacuatedX+1 != evacuatedY || evacuatedX^1 != evacuatedY {
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>					throw(&#34;bad evacuatedN&#34;)
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>				}
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>				b.tophash[i] = evacuatedX + useY <span class="comment">// evacuatedX + 1 == evacuatedY</span>
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				dst := &amp;xy[useY]                 <span class="comment">// evacuation destination</span>
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>				if dst.i == bucketCnt {
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>					dst.b = h.newoverflow(t, dst.b)
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>					dst.i = 0
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>					dst.k = add(unsafe.Pointer(dst.b), dataOffset)
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>					dst.e = add(dst.k, bucketCnt*uintptr(t.KeySize))
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>				}
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>				dst.b.tophash[dst.i&amp;(bucketCnt-1)] = top <span class="comment">// mask dst.i as an optimization, to avoid a bounds check</span>
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>				if t.IndirectKey() {
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>					*(*unsafe.Pointer)(dst.k) = k2 <span class="comment">// copy pointer</span>
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>				} else {
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>					typedmemmove(t.Key, dst.k, k) <span class="comment">// copy elem</span>
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>				}
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>				if t.IndirectElem() {
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>					*(*unsafe.Pointer)(dst.e) = *(*unsafe.Pointer)(e)
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>				} else {
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>					typedmemmove(t.Elem, dst.e, e)
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>				}
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>				dst.i++
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>				<span class="comment">// These updates might push these pointers past the end of the</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>				<span class="comment">// key or elem arrays.  That&#39;s ok, as we have the overflow pointer</span>
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>				<span class="comment">// at the end of the bucket to protect against pointing past the</span>
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>				<span class="comment">// end of the bucket.</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>				dst.k = add(dst.k, uintptr(t.KeySize))
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>				dst.e = add(dst.e, uintptr(t.ValueSize))
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>			}
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>		}
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>		<span class="comment">// Unlink the overflow buckets &amp; clear key/elem to help GC.</span>
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>		if h.flags&amp;oldIterator == 0 &amp;&amp; t.Bucket.PtrBytes != 0 {
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>			b := add(h.oldbuckets, oldbucket*uintptr(t.BucketSize))
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>			<span class="comment">// Preserve b.tophash because the evacuation</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>			<span class="comment">// state is maintained there.</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>			ptr := add(b, dataOffset)
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>			n := uintptr(t.BucketSize) - dataOffset
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>			memclrHasPointers(ptr, n)
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>		}
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>	}
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>	if oldbucket == h.nevacuate {
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>		advanceEvacuationMark(h, t, newbit)
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>	}
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>}
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>func advanceEvacuationMark(h *hmap, t *maptype, newbit uintptr) {
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>	h.nevacuate++
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>	<span class="comment">// Experiments suggest that 1024 is overkill by at least an order of magnitude.</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>	<span class="comment">// Put it in there as a safeguard anyway, to ensure O(1) behavior.</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>	stop := h.nevacuate + 1024
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>	if stop &gt; newbit {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>		stop = newbit
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>	}
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>	for h.nevacuate != stop &amp;&amp; bucketEvacuated(t, h, h.nevacuate) {
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>		h.nevacuate++
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>	}
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>	if h.nevacuate == newbit { <span class="comment">// newbit == # of oldbuckets</span>
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>		<span class="comment">// Growing is all done. Free old main bucket array.</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		h.oldbuckets = nil
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>		<span class="comment">// Can discard old overflow buckets as well.</span>
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>		<span class="comment">// If they are still referenced by an iterator,</span>
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		<span class="comment">// then the iterator holds a pointers to the slice.</span>
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>		if h.extra != nil {
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>			h.extra.oldoverflow = nil
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>		}
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>		h.flags &amp;^= sameSizeGrow
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>	}
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>}
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span><span class="comment">// Reflect stubs. Called from ../reflect/asm_*.s</span>
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_makemap reflect.makemap</span>
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>func reflect_makemap(t *maptype, cap int) *hmap {
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	<span class="comment">// Check invariants and reflects math.</span>
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>	if t.Key.Equal == nil {
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>		throw(&#34;runtime.reflect_makemap: unsupported map key type&#34;)
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>	}
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>	if t.Key.Size_ &gt; maxKeySize &amp;&amp; (!t.IndirectKey() || t.KeySize != uint8(goarch.PtrSize)) ||
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		t.Key.Size_ &lt;= maxKeySize &amp;&amp; (t.IndirectKey() || t.KeySize != uint8(t.Key.Size_)) {
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>		throw(&#34;key size wrong&#34;)
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>	}
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>	if t.Elem.Size_ &gt; maxElemSize &amp;&amp; (!t.IndirectElem() || t.ValueSize != uint8(goarch.PtrSize)) ||
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>		t.Elem.Size_ &lt;= maxElemSize &amp;&amp; (t.IndirectElem() || t.ValueSize != uint8(t.Elem.Size_)) {
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>		throw(&#34;elem size wrong&#34;)
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>	}
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>	if t.Key.Align_ &gt; bucketCnt {
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		throw(&#34;key align too big&#34;)
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>	}
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>	if t.Elem.Align_ &gt; bucketCnt {
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		throw(&#34;elem align too big&#34;)
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>	}
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>	if t.Key.Size_%uintptr(t.Key.Align_) != 0 {
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		throw(&#34;key size not a multiple of key align&#34;)
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>	}
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>	if t.Elem.Size_%uintptr(t.Elem.Align_) != 0 {
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>		throw(&#34;elem size not a multiple of elem align&#34;)
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>	}
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>	if bucketCnt &lt; 8 {
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		throw(&#34;bucketsize too small for proper alignment&#34;)
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>	}
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>	if dataOffset%uintptr(t.Key.Align_) != 0 {
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		throw(&#34;need padding in bucket (key)&#34;)
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>	}
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>	if dataOffset%uintptr(t.Elem.Align_) != 0 {
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>		throw(&#34;need padding in bucket (elem)&#34;)
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>	}
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	return makemap(t, cap, nil)
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>}
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapaccess reflect.mapaccess</span>
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>func reflect_mapaccess(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>	elem, ok := mapaccess2(t, h, key)
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>	if !ok {
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>		<span class="comment">// reflect wants nil for a missing element</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>		elem = nil
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>	}
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	return elem
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>}
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapaccess_faststr reflect.mapaccess_faststr</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>func reflect_mapaccess_faststr(t *maptype, h *hmap, key string) unsafe.Pointer {
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>	elem, ok := mapaccess2_faststr(t, h, key)
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	if !ok {
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>		<span class="comment">// reflect wants nil for a missing element</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>		elem = nil
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>	}
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>	return elem
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>}
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapassign reflect.mapassign0</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>func reflect_mapassign(t *maptype, h *hmap, key unsafe.Pointer, elem unsafe.Pointer) {
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>	p := mapassign(t, h, key)
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>	typedmemmove(t.Elem, p, elem)
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>}
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapassign_faststr reflect.mapassign_faststr0</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>func reflect_mapassign_faststr(t *maptype, h *hmap, key string, elem unsafe.Pointer) {
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>	p := mapassign_faststr(t, h, key)
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>	typedmemmove(t.Elem, p, elem)
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>}
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapdelete reflect.mapdelete</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>func reflect_mapdelete(t *maptype, h *hmap, key unsafe.Pointer) {
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>	mapdelete(t, h, key)
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>}
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapdelete_faststr reflect.mapdelete_faststr</span>
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>func reflect_mapdelete_faststr(t *maptype, h *hmap, key string) {
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>	mapdelete_faststr(t, h, key)
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>}
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapiterinit reflect.mapiterinit</span>
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>func reflect_mapiterinit(t *maptype, h *hmap, it *hiter) {
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>	mapiterinit(t, h, it)
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>}
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapiternext reflect.mapiternext</span>
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>func reflect_mapiternext(it *hiter) {
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>	mapiternext(it)
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>}
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapiterkey reflect.mapiterkey</span>
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>func reflect_mapiterkey(it *hiter) unsafe.Pointer {
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>	return it.key
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>}
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapiterelem reflect.mapiterelem</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>func reflect_mapiterelem(it *hiter) unsafe.Pointer {
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>	return it.elem
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>}
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_maplen reflect.maplen</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>func reflect_maplen(h *hmap) int {
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>	if h == nil {
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span>		return 0
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span>	}
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span>	if raceenabled {
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(reflect_maplen))
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>	}
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>	return h.count
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>}
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span><span class="comment">//go:linkname reflect_mapclear reflect.mapclear</span>
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>func reflect_mapclear(t *maptype, h *hmap) {
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	mapclear(t, h)
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>}
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span><span class="comment">//go:linkname reflectlite_maplen internal/reflectlite.maplen</span>
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>func reflectlite_maplen(h *hmap) int {
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>	if h == nil {
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>		return 0
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>	}
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>	if raceenabled {
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>		callerpc := getcallerpc()
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>		racereadpc(unsafe.Pointer(h), callerpc, abi.FuncPCABIInternal(reflect_maplen))
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>	}
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>	return h.count
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>}
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>var zeroVal [abi.ZeroValSize]byte
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span><span class="comment">// mapinitnoop is a no-op function known the Go linker; if a given global</span>
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span><span class="comment">// map (of the right size) is determined to be dead, the linker will</span>
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span><span class="comment">// rewrite the relocation (from the package init func) from the outlined</span>
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span><span class="comment">// map init function to this symbol. Defined in assembly so as to avoid</span>
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span><span class="comment">// complications with instrumentation (coverage, etc).</span>
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>func mapinitnoop()
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span><span class="comment">// mapclone for implementing maps.Clone</span>
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span><span class="comment">//go:linkname mapclone maps.clone</span>
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>func mapclone(m any) any {
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>	e := efaceOf(&amp;m)
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>	e.data = unsafe.Pointer(mapclone2((*maptype)(unsafe.Pointer(e._type)), (*hmap)(e.data)))
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>	return m
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>}
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span><span class="comment">// moveToBmap moves a bucket from src to dst. It returns the destination bucket or new destination bucket if it overflows</span>
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span><span class="comment">// and the pos that the next key/value will be written, if pos == bucketCnt means needs to written in overflow bucket.</span>
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>func moveToBmap(t *maptype, h *hmap, dst *bmap, pos int, src *bmap) (*bmap, int) {
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>	for i := 0; i &lt; bucketCnt; i++ {
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>		if isEmpty(src.tophash[i]) {
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>			continue
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span>		}
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span>
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span>		for ; pos &lt; bucketCnt; pos++ {
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>			if isEmpty(dst.tophash[pos]) {
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>				break
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>			}
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>		}
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>		if pos == bucketCnt {
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>			dst = h.newoverflow(t, dst)
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>			pos = 0
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>		}
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span>		srcK := add(unsafe.Pointer(src), dataOffset+uintptr(i)*uintptr(t.KeySize))
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span>		srcEle := add(unsafe.Pointer(src), dataOffset+bucketCnt*uintptr(t.KeySize)+uintptr(i)*uintptr(t.ValueSize))
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>		dstK := add(unsafe.Pointer(dst), dataOffset+uintptr(pos)*uintptr(t.KeySize))
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>		dstEle := add(unsafe.Pointer(dst), dataOffset+bucketCnt*uintptr(t.KeySize)+uintptr(pos)*uintptr(t.ValueSize))
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		dst.tophash[pos] = src.tophash[i]
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>		if t.IndirectKey() {
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>			srcK = *(*unsafe.Pointer)(srcK)
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>			if t.NeedKeyUpdate() {
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>				kStore := newobject(t.Key)
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>				typedmemmove(t.Key, kStore, srcK)
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>				srcK = kStore
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>			}
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>			<span class="comment">// Note: if NeedKeyUpdate is false, then the memory</span>
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>			<span class="comment">// used to store the key is immutable, so we can share</span>
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>			<span class="comment">// it between the original map and its clone.</span>
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>			*(*unsafe.Pointer)(dstK) = srcK
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>		} else {
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>			typedmemmove(t.Key, dstK, srcK)
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>		}
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>		if t.IndirectElem() {
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>			srcEle = *(*unsafe.Pointer)(srcEle)
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>			eStore := newobject(t.Elem)
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>			typedmemmove(t.Elem, eStore, srcEle)
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>			*(*unsafe.Pointer)(dstEle) = eStore
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>		} else {
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>			typedmemmove(t.Elem, dstEle, srcEle)
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>		}
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>		pos++
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>		h.count++
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>	}
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>	return dst, pos
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>}
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>func mapclone2(t *maptype, src *hmap) *hmap {
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>	dst := makemap(t, src.count, nil)
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>	dst.hash0 = src.hash0
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>	dst.nevacuate = 0
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>	<span class="comment">//flags do not need to be copied here, just like a new map has no flags.</span>
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>	if src.count == 0 {
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>		return dst
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>	}
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	if src.flags&amp;hashWriting != 0 {
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>		fatal(&#34;concurrent map clone and map write&#34;)
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>	}
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>	if src.B == 0 &amp;&amp; !(t.IndirectKey() &amp;&amp; t.NeedKeyUpdate()) &amp;&amp; !t.IndirectElem() {
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>		<span class="comment">// Quick copy for small maps.</span>
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span>		dst.buckets = newobject(t.Bucket)
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>		dst.count = src.count
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>		typedmemmove(t.Bucket, dst.buckets, src.buckets)
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>		return dst
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>	}
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>	if dst.B == 0 {
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>		dst.buckets = newobject(t.Bucket)
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>	}
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>	dstArraySize := int(bucketShift(dst.B))
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>	srcArraySize := int(bucketShift(src.B))
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>	for i := 0; i &lt; dstArraySize; i++ {
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>		dstBmap := (*bmap)(add(dst.buckets, uintptr(i*int(t.BucketSize))))
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>		pos := 0
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>		for j := 0; j &lt; srcArraySize; j += dstArraySize {
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>			srcBmap := (*bmap)(add(src.buckets, uintptr((i+j)*int(t.BucketSize))))
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>			for srcBmap != nil {
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>				dstBmap, pos = moveToBmap(t, dst, dstBmap, pos, srcBmap)
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>				srcBmap = srcBmap.overflow(t)
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>			}
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>		}
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span>	}
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span>	if src.oldbuckets == nil {
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>		return dst
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	}
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>	oldB := src.B
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>	srcOldbuckets := src.oldbuckets
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>	if !src.sameSizeGrow() {
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>		oldB--
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>	}
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>	oldSrcArraySize := int(bucketShift(oldB))
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>	for i := 0; i &lt; oldSrcArraySize; i++ {
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>		srcBmap := (*bmap)(add(srcOldbuckets, uintptr(i*int(t.BucketSize))))
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>		if evacuated(srcBmap) {
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>			continue
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>		}
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>		if oldB &gt;= dst.B { <span class="comment">// main bucket bits in dst is less than oldB bits in src</span>
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>			dstBmap := (*bmap)(add(dst.buckets, (uintptr(i)&amp;bucketMask(dst.B))*uintptr(t.BucketSize)))
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>			for dstBmap.overflow(t) != nil {
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>				dstBmap = dstBmap.overflow(t)
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>			}
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>			pos := 0
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>			for srcBmap != nil {
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>				dstBmap, pos = moveToBmap(t, dst, dstBmap, pos, srcBmap)
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>				srcBmap = srcBmap.overflow(t)
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>			}
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>			continue
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>		}
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>		<span class="comment">// oldB &lt; dst.B, so a single source bucket may go to multiple destination buckets.</span>
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>		<span class="comment">// Process entries one at a time.</span>
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>		for srcBmap != nil {
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>			<span class="comment">// move from oldBlucket to new bucket</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>			for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>				if isEmpty(srcBmap.tophash[i]) {
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>					continue
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>				}
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>				if src.flags&amp;hashWriting != 0 {
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span>					fatal(&#34;concurrent map clone and map write&#34;)
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span>				}
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span>
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span>				srcK := add(unsafe.Pointer(srcBmap), dataOffset+i*uintptr(t.KeySize))
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span>				if t.IndirectKey() {
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>					srcK = *((*unsafe.Pointer)(srcK))
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>				}
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>				srcEle := add(unsafe.Pointer(srcBmap), dataOffset+bucketCnt*uintptr(t.KeySize)+i*uintptr(t.ValueSize))
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>				if t.IndirectElem() {
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>					srcEle = *((*unsafe.Pointer)(srcEle))
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>				}
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>				dstEle := mapassign(t, dst, srcK)
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>				typedmemmove(t.Elem, dstEle, srcEle)
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>			}
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>			srcBmap = srcBmap.overflow(t)
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>		}
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>	}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>	return dst
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>}
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span><span class="comment">// keys for implementing maps.keys</span>
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span><span class="comment">//go:linkname keys maps.keys</span>
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>func keys(m any, p unsafe.Pointer) {
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>	e := efaceOf(&amp;m)
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	t := (*maptype)(unsafe.Pointer(e._type))
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>	h := (*hmap)(e.data)
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		return
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>	}
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>	s := (*slice)(p)
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>	r := int(rand())
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>	offset := uint8(r &gt;&gt; h.B &amp; (bucketCnt - 1))
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>	if h.B == 0 {
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>		copyKeys(t, h, (*bmap)(h.buckets), s, offset)
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>		return
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>	}
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>	arraySize := int(bucketShift(h.B))
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>	buckets := h.buckets
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>	for i := 0; i &lt; arraySize; i++ {
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>		bucket := (i + r) &amp; (arraySize - 1)
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>		b := (*bmap)(add(buckets, uintptr(bucket)*uintptr(t.BucketSize)))
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>		copyKeys(t, h, b, s, offset)
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>	}
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>	if h.growing() {
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>		oldArraySize := int(h.noldbuckets())
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>		for i := 0; i &lt; oldArraySize; i++ {
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>			bucket := (i + r) &amp; (oldArraySize - 1)
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>			b := (*bmap)(add(h.oldbuckets, uintptr(bucket)*uintptr(t.BucketSize)))
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>			if evacuated(b) {
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>				continue
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>			}
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>			copyKeys(t, h, b, s, offset)
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>		}
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>	}
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>	return
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>}
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>func copyKeys(t *maptype, h *hmap, b *bmap, s *slice, offset uint8) {
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>	for b != nil {
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>			offi := (i + uintptr(offset)) &amp; (bucketCnt - 1)
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>			if isEmpty(b.tophash[offi]) {
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>				continue
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>			}
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>			if h.flags&amp;hashWriting != 0 {
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>				fatal(&#34;concurrent map read and map write&#34;)
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>			}
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>			k := add(unsafe.Pointer(b), dataOffset+offi*uintptr(t.KeySize))
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>			if t.IndirectKey() {
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>				k = *((*unsafe.Pointer)(k))
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>			}
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>			if s.len &gt;= s.cap {
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>				fatal(&#34;concurrent map read and map write&#34;)
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>			}
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>			typedmemmove(t.Key, add(s.array, uintptr(s.len)*uintptr(t.Key.Size())), k)
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>			s.len++
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>		}
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>		b = b.overflow(t)
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>	}
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>}
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span><span class="comment">// values for implementing maps.values</span>
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span><span class="comment">//go:linkname values maps.values</span>
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>func values(m any, p unsafe.Pointer) {
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>	e := efaceOf(&amp;m)
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>	t := (*maptype)(unsafe.Pointer(e._type))
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>	h := (*hmap)(e.data)
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>	if h == nil || h.count == 0 {
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>		return
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>	}
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	s := (*slice)(p)
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>	r := int(rand())
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	offset := uint8(r &gt;&gt; h.B &amp; (bucketCnt - 1))
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>	if h.B == 0 {
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>		copyValues(t, h, (*bmap)(h.buckets), s, offset)
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>		return
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	}
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>	arraySize := int(bucketShift(h.B))
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>	buckets := h.buckets
<span id="L1688" class="ln">  1688&nbsp;&nbsp;</span>	for i := 0; i &lt; arraySize; i++ {
<span id="L1689" class="ln">  1689&nbsp;&nbsp;</span>		bucket := (i + r) &amp; (arraySize - 1)
<span id="L1690" class="ln">  1690&nbsp;&nbsp;</span>		b := (*bmap)(add(buckets, uintptr(bucket)*uintptr(t.BucketSize)))
<span id="L1691" class="ln">  1691&nbsp;&nbsp;</span>		copyValues(t, h, b, s, offset)
<span id="L1692" class="ln">  1692&nbsp;&nbsp;</span>	}
<span id="L1693" class="ln">  1693&nbsp;&nbsp;</span>
<span id="L1694" class="ln">  1694&nbsp;&nbsp;</span>	if h.growing() {
<span id="L1695" class="ln">  1695&nbsp;&nbsp;</span>		oldArraySize := int(h.noldbuckets())
<span id="L1696" class="ln">  1696&nbsp;&nbsp;</span>		for i := 0; i &lt; oldArraySize; i++ {
<span id="L1697" class="ln">  1697&nbsp;&nbsp;</span>			bucket := (i + r) &amp; (oldArraySize - 1)
<span id="L1698" class="ln">  1698&nbsp;&nbsp;</span>			b := (*bmap)(add(h.oldbuckets, uintptr(bucket)*uintptr(t.BucketSize)))
<span id="L1699" class="ln">  1699&nbsp;&nbsp;</span>			if evacuated(b) {
<span id="L1700" class="ln">  1700&nbsp;&nbsp;</span>				continue
<span id="L1701" class="ln">  1701&nbsp;&nbsp;</span>			}
<span id="L1702" class="ln">  1702&nbsp;&nbsp;</span>			copyValues(t, h, b, s, offset)
<span id="L1703" class="ln">  1703&nbsp;&nbsp;</span>		}
<span id="L1704" class="ln">  1704&nbsp;&nbsp;</span>	}
<span id="L1705" class="ln">  1705&nbsp;&nbsp;</span>	return
<span id="L1706" class="ln">  1706&nbsp;&nbsp;</span>}
<span id="L1707" class="ln">  1707&nbsp;&nbsp;</span>
<span id="L1708" class="ln">  1708&nbsp;&nbsp;</span>func copyValues(t *maptype, h *hmap, b *bmap, s *slice, offset uint8) {
<span id="L1709" class="ln">  1709&nbsp;&nbsp;</span>	for b != nil {
<span id="L1710" class="ln">  1710&nbsp;&nbsp;</span>		for i := uintptr(0); i &lt; bucketCnt; i++ {
<span id="L1711" class="ln">  1711&nbsp;&nbsp;</span>			offi := (i + uintptr(offset)) &amp; (bucketCnt - 1)
<span id="L1712" class="ln">  1712&nbsp;&nbsp;</span>			if isEmpty(b.tophash[offi]) {
<span id="L1713" class="ln">  1713&nbsp;&nbsp;</span>				continue
<span id="L1714" class="ln">  1714&nbsp;&nbsp;</span>			}
<span id="L1715" class="ln">  1715&nbsp;&nbsp;</span>
<span id="L1716" class="ln">  1716&nbsp;&nbsp;</span>			if h.flags&amp;hashWriting != 0 {
<span id="L1717" class="ln">  1717&nbsp;&nbsp;</span>				fatal(&#34;concurrent map read and map write&#34;)
<span id="L1718" class="ln">  1718&nbsp;&nbsp;</span>			}
<span id="L1719" class="ln">  1719&nbsp;&nbsp;</span>
<span id="L1720" class="ln">  1720&nbsp;&nbsp;</span>			ele := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.KeySize)+offi*uintptr(t.ValueSize))
<span id="L1721" class="ln">  1721&nbsp;&nbsp;</span>			if t.IndirectElem() {
<span id="L1722" class="ln">  1722&nbsp;&nbsp;</span>				ele = *((*unsafe.Pointer)(ele))
<span id="L1723" class="ln">  1723&nbsp;&nbsp;</span>			}
<span id="L1724" class="ln">  1724&nbsp;&nbsp;</span>			if s.len &gt;= s.cap {
<span id="L1725" class="ln">  1725&nbsp;&nbsp;</span>				fatal(&#34;concurrent map read and map write&#34;)
<span id="L1726" class="ln">  1726&nbsp;&nbsp;</span>			}
<span id="L1727" class="ln">  1727&nbsp;&nbsp;</span>			typedmemmove(t.Elem, add(s.array, uintptr(s.len)*uintptr(t.Elem.Size())), ele)
<span id="L1728" class="ln">  1728&nbsp;&nbsp;</span>			s.len++
<span id="L1729" class="ln">  1729&nbsp;&nbsp;</span>		}
<span id="L1730" class="ln">  1730&nbsp;&nbsp;</span>		b = b.overflow(t)
<span id="L1731" class="ln">  1731&nbsp;&nbsp;</span>	}
<span id="L1732" class="ln">  1732&nbsp;&nbsp;</span>}
<span id="L1733" class="ln">  1733&nbsp;&nbsp;</span>
</pre><p><a href="map.go?m=text">View as plain text</a></p>

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
