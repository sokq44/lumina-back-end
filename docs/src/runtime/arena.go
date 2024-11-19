<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/arena.go - Go Documentation Server</title>

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
<a href="arena.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<span class="text-muted">arena.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime">runtime</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Implementation of (safe) user arenas.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file contains the implementation of user arenas wherein Go values can</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// be manually allocated and freed in bulk. The act of manually freeing memory,</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// potentially before a GC cycle, means that a garbage collection cycle can be</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// delayed, improving efficiency by reducing GC cycle frequency. There are other</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// potential efficiency benefits, such as improved locality and access to a more</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// efficient allocation strategy.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// What makes the arenas here safe is that once they are freed, accessing the</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// arena&#39;s memory will cause an explicit program fault, and the arena&#39;s address</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// space will not be reused until no more pointers into it are found. There&#39;s one</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// exception to this: if an arena allocated memory that isn&#39;t exhausted, it&#39;s placed</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// back into a pool for reuse. This means that a crash is not always guaranteed.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// While this may seem unsafe, it still prevents memory corruption, and is in fact</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// necessary in order to make new(T) a valid implementation of arenas. Such a property</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// is desirable to allow for a trivial implementation. (It also avoids complexities</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// that arise from synchronization with the GC when trying to set the arena chunks to</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// fault while the GC is active.)</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The implementation works in layers. At the bottom, arenas are managed in chunks.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Each chunk must be a multiple of the heap arena size, or the heap arena size must</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// be divisible by the arena chunks. The address space for each chunk, and each</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// corresponding heapArena for that address space, are eternally reserved for use as</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// arena chunks. That is, they can never be used for the general heap. Each chunk</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// is also represented by a single mspan, and is modeled as a single large heap</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// allocation. It must be, because each chunk contains ordinary Go values that may</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// point into the heap, so it must be scanned just like any other object. Any</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// pointer into a chunk will therefore always cause the whole chunk to be scanned</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// while its corresponding arena is still live.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// Chunks may be allocated either from new memory mapped by the OS on our behalf,</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// or by reusing old freed chunks. When chunks are freed, their underlying memory</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// is returned to the OS, set to fault on access, and may not be reused until the</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// program doesn&#39;t point into the chunk anymore (the code refers to this state as</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// &#34;quarantined&#34;), a property checked by the GC.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// The sweeper handles moving chunks out of this quarantine state to be ready for</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// reuse. When the chunk is placed into the quarantine state, its corresponding</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// span is marked as noscan so that the GC doesn&#39;t try to scan memory that would</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// cause a fault.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// At the next layer are the user arenas themselves. They consist of a single</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// active chunk which new Go values are bump-allocated into and a list of chunks</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// that were exhausted when allocating into the arena. Once the arena is freed,</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// it frees all full chunks it references, and places the active one onto a reuse</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// list for a future arena to use. Each arena keeps its list of referenced chunks</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// explicitly live until it is freed. Each user arena also maps to an object which</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// has a finalizer attached that ensures the arena&#39;s chunks are all freed even if</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// the arena itself is never explicitly freed.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// Pointer-ful memory is bump-allocated from low addresses to high addresses in each</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// chunk, while pointer-free memory is bump-allocated from high address to low</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// addresses. The reason for this is to take advantage of a GC optimization wherein</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">// the GC will stop scanning an object when there are no more pointers in it, which</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// also allows us to elide clearing the heap bitmap for pointer-free Go values</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// allocated into arenas.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// Note that arenas are not safe to use concurrently.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// In summary, there are 2 resources: arenas, and arena chunks. They exist in the</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// following lifecycle:</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// (1) A new arena is created via newArena.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// (2) Chunks are allocated to hold memory allocated into the arena with new or slice.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//    (a) Chunks are first allocated from the reuse list of partially-used chunks.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//    (b) If there are no such chunks, then chunks on the ready list are taken.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//    (c) Failing all the above, memory for a new chunk is mapped.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// (3) The arena is freed, or all references to it are dropped, triggering its finalizer.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//    (a) If the GC is not active, exhausted chunks are set to fault and placed on a</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//        quarantine list.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//    (b) If the GC is active, exhausted chunks are placed on a fault list and will</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//        go through step (a) at a later point in time.</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//    (c) Any remaining partially-used chunk is placed on a reuse list.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// (4) Once no more pointers are found into quarantined arena chunks, the sweeper</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//     takes these chunks out of quarantine and places them on the ready list.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>package runtime
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>import (
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	&#34;internal/goarch&#34;
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	&#34;internal/goexperiment&#34;
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	&#34;runtime/internal/math&#34;
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>)
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// Functions starting with arena_ are meant to be exported to downstream users</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// of arenas. They should wrap these functions in a higher-lever API.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// The underlying arena and its resources are managed through an opaque unsafe.Pointer.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// arena_newArena is a wrapper around newUserArena.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">//go:linkname arena_newArena arena.runtime_arena_newArena</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func arena_newArena() unsafe.Pointer {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	return unsafe.Pointer(newUserArena())
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// arena_arena_New is a wrapper around (*userArena).new, except that typ</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// is an any (must be a *_type, still) and typ must be a type descriptor</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// for a pointer to the type to actually be allocated, i.e. pass a *T</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// to allocate a T. This is necessary because this function returns a *T.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">//go:linkname arena_arena_New arena.runtime_arena_arena_New</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func arena_arena_New(arena unsafe.Pointer, typ any) any {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	t := (*_type)(efaceOf(&amp;typ).data)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if t.Kind_&amp;kindMask != kindPtr {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		throw(&#34;arena_New: non-pointer type&#34;)
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	te := (*ptrtype)(unsafe.Pointer(t)).Elem
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	x := ((*userArena)(arena)).new(te)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	var result any
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	e := efaceOf(&amp;result)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	e._type = t
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	e.data = x
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	return result
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// arena_arena_Slice is a wrapper around (*userArena).slice.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">//go:linkname arena_arena_Slice arena.runtime_arena_arena_Slice</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func arena_arena_Slice(arena unsafe.Pointer, slice any, cap int) {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	((*userArena)(arena)).slice(slice, cap)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// arena_arena_Free is a wrapper around (*userArena).free.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">//go:linkname arena_arena_Free arena.runtime_arena_arena_Free</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>func arena_arena_Free(arena unsafe.Pointer) {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	((*userArena)(arena)).free()
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// arena_heapify takes a value that lives in an arena and makes a copy</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// of it on the heap. Values that don&#39;t live in an arena are returned unmodified.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//go:linkname arena_heapify arena.runtime_arena_heapify</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func arena_heapify(s any) any {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	var v unsafe.Pointer
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	e := efaceOf(&amp;s)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	t := e._type
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	switch t.Kind_ &amp; kindMask {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	case kindString:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		v = stringStructOf((*string)(e.data)).str
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	case kindSlice:
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		v = (*slice)(e.data).array
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	case kindPtr:
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		v = e.data
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	default:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		panic(&#34;arena: Clone only supports pointers, slices, and strings&#34;)
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	span := spanOf(uintptr(v))
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if span == nil || !span.isUserArenaChunk {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// Not stored in a user arena chunk.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		return s
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Heap-allocate storage for a copy.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	var x any
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	switch t.Kind_ &amp; kindMask {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	case kindString:
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		s1 := s.(string)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		s2, b := rawstring(len(s1))
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		copy(b, s1)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		x = s2
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	case kindSlice:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		len := (*slice)(e.data).len
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		et := (*slicetype)(unsafe.Pointer(t)).Elem
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		sl := new(slice)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		*sl = slice{makeslicecopy(et, len, len, (*slice)(e.data).array), len, len}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		xe := efaceOf(&amp;x)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		xe._type = t
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		xe.data = unsafe.Pointer(sl)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	case kindPtr:
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		et := (*ptrtype)(unsafe.Pointer(t)).Elem
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		e2 := newobject(et)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		typedmemmove(et, e2, e.data)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		xe := efaceOf(&amp;x)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		xe._type = t
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		xe.data = e2
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	return x
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>const (
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// userArenaChunkBytes is the size of a user arena chunk.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	userArenaChunkBytesMax = 8 &lt;&lt; 20
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	userArenaChunkBytes    = uintptr(int64(userArenaChunkBytesMax-heapArenaBytes)&amp;(int64(userArenaChunkBytesMax-heapArenaBytes)&gt;&gt;63) + heapArenaBytes) <span class="comment">// min(userArenaChunkBytesMax, heapArenaBytes)</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// userArenaChunkPages is the number of pages a user arena chunk uses.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	userArenaChunkPages = userArenaChunkBytes / pageSize
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// userArenaChunkMaxAllocBytes is the maximum size of an object that can</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// be allocated from an arena. This number is chosen to cap worst-case</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// fragmentation of user arenas to 25%. Larger allocations are redirected</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// to the heap.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	userArenaChunkMaxAllocBytes = userArenaChunkBytes / 4
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func init() {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if userArenaChunkPages*pageSize != userArenaChunkBytes {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		throw(&#34;user arena chunk size is not a multiple of the page size&#34;)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	if userArenaChunkBytes%physPageSize != 0 {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		throw(&#34;user arena chunk size is not a multiple of the physical page size&#34;)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if userArenaChunkBytes &lt; heapArenaBytes {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		if heapArenaBytes%userArenaChunkBytes != 0 {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			throw(&#34;user arena chunk size is smaller than a heap arena, but doesn&#39;t divide it&#34;)
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	} else {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		if userArenaChunkBytes%heapArenaBytes != 0 {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			throw(&#34;user arena chunks size is larger than a heap arena, but not a multiple&#34;)
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	lockInit(&amp;userArenaState.lock, lockRankUserArenaState)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span><span class="comment">// userArenaChunkReserveBytes returns the amount of additional bytes to reserve for</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span><span class="comment">// heap metadata.</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>func userArenaChunkReserveBytes() uintptr {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		<span class="comment">// In the allocation headers experiment, we reserve the end of the chunk for</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// a pointer/scalar bitmap. We also reserve space for a dummy _type that</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// refers to the bitmap. The PtrBytes field of the dummy _type indicates how</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// many of those bits are valid.</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		return userArenaChunkBytes/goarch.PtrSize/8 + unsafe.Sizeof(_type{})
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	return 0
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>type userArena struct {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// full is a list of full chunks that have not enough free memory left, and</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// that we&#39;ll free once this user arena is freed.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">// Can&#39;t use mSpanList here because it&#39;s not-in-heap.</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	fullList *mspan
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// active is the user arena chunk we&#39;re currently allocating into.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	active *mspan
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">// refs is a set of references to the arena chunks so that they&#39;re kept alive.</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// The last reference in the list always refers to active, while the rest of</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// them correspond to fullList. Specifically, the head of fullList is the</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// second-to-last one, fullList.next is the third-to-last, and so on.</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// In other words, every time a new chunk becomes active, its appended to this</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// list.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	refs []unsafe.Pointer
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	<span class="comment">// defunct is true if free has been called on this arena.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// This is just a best-effort way to discover a concurrent allocation</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// and free. Also used to detect a double-free.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	defunct atomic.Bool
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// newUserArena creates a new userArena ready to be used.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>func newUserArena() *userArena {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	a := new(userArena)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	SetFinalizer(a, func(a *userArena) {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		<span class="comment">// If arena handle is dropped without being freed, then call</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		<span class="comment">// free on the arena, so the arena chunks are never reclaimed</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// by the garbage collector.</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		a.free()
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	})
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	a.refill()
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	return a
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// new allocates a new object of the provided type into the arena, and returns</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// its pointer.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// This operation is not safe to call concurrently with other operations on the</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// same arena.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>func (a *userArena) new(typ *_type) unsafe.Pointer {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	return a.alloc(typ, -1)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// slice allocates a new slice backing store. slice must be a pointer to a slice</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// (i.e. *[]T), because userArenaSlice will update the slice directly.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// cap determines the capacity of the slice backing store and must be non-negative.</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// This operation is not safe to call concurrently with other operations on the</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// same arena.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>func (a *userArena) slice(sl any, cap int) {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if cap &lt; 0 {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		panic(&#34;userArena.slice: negative cap&#34;)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	i := efaceOf(&amp;sl)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	typ := i._type
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindMask != kindPtr {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		panic(&#34;slice result of non-ptr type&#34;)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	typ = (*ptrtype)(unsafe.Pointer(typ)).Elem
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	if typ.Kind_&amp;kindMask != kindSlice {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		panic(&#34;slice of non-ptr-to-slice type&#34;)
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	typ = (*slicetype)(unsafe.Pointer(typ)).Elem
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// t is now the element type of the slice we want to allocate.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	*((*slice)(i.data)) = slice{a.alloc(typ, cap), cap, cap}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span><span class="comment">// free returns the userArena&#39;s chunks back to mheap and marks it as defunct.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span><span class="comment">// Must be called at most once for any given arena.</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// This operation is not safe to call concurrently with other operations on the</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span><span class="comment">// same arena.</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>func (a *userArena) free() {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	<span class="comment">// Check for a double-free.</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if a.defunct.Load() {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		panic(&#34;arena double free&#34;)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	<span class="comment">// Mark ourselves as defunct.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	a.defunct.Store(true)
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	SetFinalizer(a, nil)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// Free all the full arenas.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// The refs on this list are in reverse order from the second-to-last.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	s := a.fullList
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	i := len(a.refs) - 2
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	for s != nil {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		a.fullList = s.next
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		s.next = nil
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		freeUserArenaChunk(s, a.refs[i])
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		s = a.fullList
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		i--
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	}
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if a.fullList != nil || i &gt;= 0 {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		<span class="comment">// There&#39;s still something left on the full list, or we</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		<span class="comment">// failed to actually iterate over the entire refs list.</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		throw(&#34;full list doesn&#39;t match refs list in length&#34;)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	<span class="comment">// Put the active chunk onto the reuse list.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	<span class="comment">// Note that active&#39;s reference is always the last reference in refs.</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	s = a.active
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	if s != nil {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		if raceenabled || msanenabled || asanenabled {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			<span class="comment">// Don&#39;t reuse arenas with sanitizers enabled. We want to catch</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			<span class="comment">// any use-after-free errors aggressively.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>			freeUserArenaChunk(s, a.refs[len(a.refs)-1])
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		} else {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			lock(&amp;userArenaState.lock)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			userArenaState.reuse = append(userArenaState.reuse, liveUserArenaChunk{s, a.refs[len(a.refs)-1]})
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			unlock(&amp;userArenaState.lock)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	<span class="comment">// nil out a.active so that a race with freeing will more likely cause a crash.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	a.active = nil
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	a.refs = nil
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span><span class="comment">// alloc reserves space in the current chunk or calls refill and reserves space</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span><span class="comment">// in a new chunk. If cap is negative, the type will be taken literally, otherwise</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// it will be considered as an element type for a slice backing store with capacity</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// cap.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (a *userArena) alloc(typ *_type, cap int) unsafe.Pointer {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	s := a.active
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	var x unsafe.Pointer
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	for {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		x = s.userArenaNextFree(typ, cap)
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		if x != nil {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>			break
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		s = a.refill()
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	return x
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// refill inserts the current arena chunk onto the full list and obtains a new</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">// one, either from the partial list or allocating a new one, both from mheap.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>func (a *userArena) refill() *mspan {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	<span class="comment">// If there&#39;s an active chunk, assume it&#39;s full.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	s := a.active
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	if s != nil {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		if s.userArenaChunkFree.size() &gt; userArenaChunkMaxAllocBytes {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s difficult to tell when we&#39;re actually out of memory</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			<span class="comment">// in a chunk because the allocation that failed may still leave</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			<span class="comment">// some free space available. However, that amount of free space</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			<span class="comment">// should never exceed the maximum allocation size.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			throw(&#34;wasted too much memory in an arena chunk&#34;)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		s.next = a.fullList
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		a.fullList = s
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		a.active = nil
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		s = nil
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	var x unsafe.Pointer
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">// Check the partially-used list.</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	lock(&amp;userArenaState.lock)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	if len(userArenaState.reuse) &gt; 0 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		<span class="comment">// Pick off the last arena chunk from the list.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>		n := len(userArenaState.reuse) - 1
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		x = userArenaState.reuse[n].x
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		s = userArenaState.reuse[n].mspan
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		userArenaState.reuse[n].x = nil
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		userArenaState.reuse[n].mspan = nil
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		userArenaState.reuse = userArenaState.reuse[:n]
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	unlock(&amp;userArenaState.lock)
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	if s == nil {
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// Allocate a new one.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		x, s = newUserArenaChunk()
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		if s == nil {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			throw(&#34;out of memory&#34;)
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	}
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	a.refs = append(a.refs, x)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	a.active = s
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	return s
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>}
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>type liveUserArenaChunk struct {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	*mspan <span class="comment">// Must represent a user arena chunk.</span>
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// Reference to mspan.base() to keep the chunk alive.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	x unsafe.Pointer
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>}
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>var userArenaState struct {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	lock mutex
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// reuse contains a list of partially-used and already-live</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// user arena chunks that can be quickly reused for another</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// arena.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// Protected by lock.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	reuse []liveUserArenaChunk
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">// fault contains full user arena chunks that need to be faulted.</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	<span class="comment">// Protected by lock.</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	fault []liveUserArenaChunk
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span><span class="comment">// userArenaNextFree reserves space in the user arena for an item of the specified</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span><span class="comment">// type. If cap is not -1, this is for an array of cap elements of type t.</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>func (s *mspan) userArenaNextFree(typ *_type, cap int) unsafe.Pointer {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	size := typ.Size_
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	if cap &gt; 0 {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		if size &gt; ^uintptr(0)/uintptr(cap) {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			<span class="comment">// Overflow.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>			throw(&#34;out of memory&#34;)
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		size *= uintptr(cap)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	if size == 0 || cap == 0 {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		return unsafe.Pointer(&amp;zerobase)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	if size &gt; userArenaChunkMaxAllocBytes {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		<span class="comment">// Redirect allocations that don&#39;t fit into a chunk well directly</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		<span class="comment">// from the heap.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		if cap &gt;= 0 {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			return newarray(typ, cap)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		return newobject(typ)
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>	}
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	<span class="comment">// Prevent preemption as we set up the space for a new object.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	<span class="comment">// Act like we&#39;re allocating.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	if mp.mallocing != 0 {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		throw(&#34;malloc deadlock&#34;)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	if mp.gsignal == getg() {
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		throw(&#34;malloc during signal&#34;)
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	mp.mallocing = 1
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	var ptr unsafe.Pointer
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	if typ.PtrBytes == 0 {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		<span class="comment">// Allocate pointer-less objects from the tail end of the chunk.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		v, ok := s.userArenaChunkFree.takeFromBack(size, typ.Align_)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		if ok {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			ptr = unsafe.Pointer(v)
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	} else {
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		v, ok := s.userArenaChunkFree.takeFromFront(size, typ.Align_)
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		if ok {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			ptr = unsafe.Pointer(v)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	if ptr == nil {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		<span class="comment">// Failed to allocate.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		mp.mallocing = 0
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		releasem(mp)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		return nil
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>	if s.needzero != 0 {
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		throw(&#34;arena chunk needs zeroing, but should already be zeroed&#34;)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	<span class="comment">// Set up heap bitmap and do extra accounting.</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	if typ.PtrBytes != 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		if cap &gt;= 0 {
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			userArenaHeapBitsSetSliceType(typ, cap, ptr, s)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>		} else {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			userArenaHeapBitsSetType(typ, ptr, s)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>		c := getMCache(mp)
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		if c == nil {
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			throw(&#34;mallocgc called without a P or outside bootstrapping&#34;)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		if cap &gt; 0 {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			c.scanAlloc += size - (typ.Size_ - typ.PtrBytes)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		} else {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			c.scanAlloc += typ.PtrBytes
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	<span class="comment">// Ensure that the stores above that initialize x to</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	<span class="comment">// type-safe memory and set the heap bits occur before</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">// the caller can make ptr observable to the garbage</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">// collector. Otherwise, on weakly ordered machines,</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	<span class="comment">// the garbage collector could follow a pointer to x,</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	<span class="comment">// but see uninitialized memory or stale heap bits.</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	publicationBarrier()
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	mp.mallocing = 0
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	releasem(mp)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	return ptr
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">// userArenaHeapBitsSetSliceType is the equivalent of heapBitsSetType but for</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// Go slice backing store values allocated in a user arena chunk. It sets up the</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// heap bitmap for n consecutive values with type typ allocated at address ptr.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>func userArenaHeapBitsSetSliceType(typ *_type, n int, ptr unsafe.Pointer, s *mspan) {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	mem, overflow := math.MulUintptr(typ.Size_, uintptr(n))
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if overflow || n &lt; 0 || mem &gt; maxAlloc {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		panic(plainError(&#34;runtime: allocation size out of range&#34;))
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	for i := 0; i &lt; n; i++ {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		userArenaHeapBitsSetType(typ, add(ptr, uintptr(i)*typ.Size_), s)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// newUserArenaChunk allocates a user arena chunk, which maps to a single</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// heap arena and single span. Returns a pointer to the base of the chunk</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// (this is really important: we need to keep the chunk alive) and the span.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func newUserArenaChunk() (unsafe.Pointer, *mspan) {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	if gcphase == _GCmarktermination {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		throw(&#34;newUserArenaChunk called with gcphase == _GCmarktermination&#34;)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// Deduct assist credit. Because user arena chunks are modeled as one</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	<span class="comment">// giant heap object which counts toward heapLive, we&#39;re obligated to</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	<span class="comment">// assist the GC proportionally (and it&#39;s worth noting that the arena</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	<span class="comment">// does represent additional work for the GC, but we also have no idea</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	<span class="comment">// what that looks like until we actually allocate things into the</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	<span class="comment">// arena).</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	deductAssistCredit(userArenaChunkBytes)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	<span class="comment">// Set mp.mallocing to keep from being preempted by GC.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	if mp.mallocing != 0 {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		throw(&#34;malloc deadlock&#34;)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	}
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	if mp.gsignal == getg() {
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		throw(&#34;malloc during signal&#34;)
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	}
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	mp.mallocing = 1
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	<span class="comment">// Allocate a new user arena.</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	var span *mspan
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		span = mheap_.allocUserArenaChunk()
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	})
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	if span == nil {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		throw(&#34;out of memory&#34;)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	x := unsafe.Pointer(span.base())
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// Allocate black during GC.</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	<span class="comment">// All slots hold nil so no scanning is needed.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	<span class="comment">// This may be racing with GC so do it atomically if there can be</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	<span class="comment">// a race marking the bit.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	if gcphase != _GCoff {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		gcmarknewobject(span, span.base())
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	if raceenabled {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Track individual objects.</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		racemalloc(unsafe.Pointer(span.base()), span.elemsize)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	if msanenabled {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Track individual objects.</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		msanmalloc(unsafe.Pointer(span.base()), span.elemsize)
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	}
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	if asanenabled {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		<span class="comment">// TODO(mknyszek): Track individual objects.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		rzSize := computeRZlog(span.elemsize)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		span.elemsize -= rzSize
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		if goexperiment.AllocHeaders {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			span.largeType.Size_ = span.elemsize
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		rzStart := span.base() + span.elemsize
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		span.userArenaChunkFree = makeAddrRange(span.base(), rzStart)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		asanpoison(unsafe.Pointer(rzStart), span.limit-rzStart)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		asanunpoison(unsafe.Pointer(span.base()), span.elemsize)
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>	if rate := MemProfileRate; rate &gt; 0 {
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		c := getMCache(mp)
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		if c == nil {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			throw(&#34;newUserArenaChunk called without a P or outside bootstrapping&#34;)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>		<span class="comment">// Note cache c only valid while m acquired; see #47302</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		if rate != 1 &amp;&amp; userArenaChunkBytes &lt; c.nextSample {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>			c.nextSample -= userArenaChunkBytes
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		} else {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			profilealloc(mp, unsafe.Pointer(span.base()), userArenaChunkBytes)
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		}
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	mp.mallocing = 0
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	releasem(mp)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>	<span class="comment">// Again, because this chunk counts toward heapLive, potentially trigger a GC.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>	if t := (gcTrigger{kind: gcTriggerHeap}); t.test() {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>		gcStart(t)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	}
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	if debug.malloc {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>		if debug.allocfreetrace != 0 {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			tracealloc(unsafe.Pointer(span.base()), userArenaChunkBytes, nil)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		if inittrace.active &amp;&amp; inittrace.id == getg().goid {
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>			<span class="comment">// Init functions are executed sequentially in a single goroutine.</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			inittrace.bytes += uint64(userArenaChunkBytes)
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>		}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	<span class="comment">// Double-check it&#39;s aligned to the physical page size. Based on the current</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>	<span class="comment">// implementation this is trivially true, but it need not be in the future.</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	<span class="comment">// However, if it&#39;s not aligned to the physical page size then we can&#39;t properly</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	<span class="comment">// set it to fault later.</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	if uintptr(x)%physPageSize != 0 {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		throw(&#34;user arena chunk is not aligned to the physical page size&#34;)
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	return x, span
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span><span class="comment">// isUnusedUserArenaChunk indicates that the arena chunk has been set to fault</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span><span class="comment">// and doesn&#39;t contain any scannable memory anymore. However, it might still be</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span><span class="comment">// mSpanInUse as it sits on the quarantine list, since it needs to be swept.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span><span class="comment">// This is not safe to execute unless the caller has ownership of the mspan or</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span><span class="comment">// the world is stopped (preemption is prevented while the relevant state changes).</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// This is really only meant to be used by accounting tests in the runtime to</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// distinguish when a span shouldn&#39;t be counted (since mSpanInUse might not be</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span><span class="comment">// enough).</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>func (s *mspan) isUnusedUserArenaChunk() bool {
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	return s.isUserArenaChunk &amp;&amp; s.spanclass == makeSpanClass(0, true)
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>}
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span><span class="comment">// setUserArenaChunkToFault sets the address space for the user arena chunk to fault</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span><span class="comment">// and releases any underlying memory resources.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span><span class="comment">// Must be in a non-preemptible state to ensure the consistency of statistics</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span><span class="comment">// exported to MemStats.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>func (s *mspan) setUserArenaChunkToFault() {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	if !s.isUserArenaChunk {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		throw(&#34;invalid span in heapArena for user arena&#34;)
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	if s.npages*pageSize != userArenaChunkBytes {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		throw(&#34;span on userArena.faultList has invalid size&#34;)
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	}
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	<span class="comment">// Update the span class to be noscan. What we want to happen is that</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">// any pointer into the span keeps it from getting recycled, so we want</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// the mark bit to get set, but we&#39;re about to set the address space to fault,</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// so we have to prevent the GC from scanning this memory.</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// It&#39;s OK to set it here because (1) a GC isn&#39;t in progress, so the scanning code</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	<span class="comment">// won&#39;t make a bad decision, (2) we&#39;re currently non-preemptible and in the runtime,</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">// so a GC is blocked from starting. We might race with sweeping, which could</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">// put it on the &#34;wrong&#34; sweep list, but really don&#39;t care because the chunk is</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	<span class="comment">// treated as a large object span and there&#39;s no meaningful difference between scan</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	<span class="comment">// and noscan large objects in the sweeper. The STW at the start of the GC acts as a</span>
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	<span class="comment">// barrier for this update.</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>	s.spanclass = makeSpanClass(0, true)
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	<span class="comment">// Actually set the arena chunk to fault, so we&#39;ll get dangling pointer errors.</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>	<span class="comment">// sysFault currently uses a method on each OS that forces it to evacuate all</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>	<span class="comment">// memory backing the chunk.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>	sysFault(unsafe.Pointer(s.base()), s.npages*pageSize)
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	<span class="comment">// Everything on the list is counted as in-use, however sysFault transitions to</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	<span class="comment">// Reserved, not Prepared, so we skip updating heapFree or heapReleased and just</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	<span class="comment">// remove the memory from the total altogether; it&#39;s just address space now.</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>	gcController.heapInUse.add(-int64(s.npages * pageSize))
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	<span class="comment">// Count this as a free of an object right now as opposed to when</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	<span class="comment">// the span gets off the quarantine list. The main reason is so that the</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	<span class="comment">// amount of bytes allocated doesn&#39;t exceed how much is counted as</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>	<span class="comment">// &#34;mapped ready,&#34; which could cause a deadlock in the pacer.</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	gcController.totalFree.Add(int64(s.elemsize))
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	<span class="comment">// Update consistent stats to match.</span>
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	<span class="comment">// We&#39;re non-preemptible, so it&#39;s safe to update consistent stats (our P</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>	<span class="comment">// won&#39;t change out from under us).</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.committed, -int64(s.npages*pageSize))
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.inHeap, -int64(s.npages*pageSize))
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeFreeCount, 1)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeFree, int64(s.elemsize))
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	<span class="comment">// This counts as a free, so update heapLive.</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	gcController.update(-int64(s.elemsize), 0)
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// Mark it as free for the race detector.</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	if raceenabled {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		racefree(unsafe.Pointer(s.base()), s.elemsize)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	systemstack(func() {
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		<span class="comment">// Add the user arena to the quarantine list.</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		lock(&amp;mheap_.lock)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		mheap_.userArena.quarantineList.insert(s)
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		unlock(&amp;mheap_.lock)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	})
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>}
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span><span class="comment">// inUserArenaChunk returns true if p points to a user arena chunk.</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>func inUserArenaChunk(p uintptr) bool {
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	s := spanOf(p)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	if s == nil {
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		return false
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>	}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	return s.isUserArenaChunk
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// freeUserArenaChunk releases the user arena represented by s back to the runtime.</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span><span class="comment">// x must be a live pointer within s.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span><span class="comment">// The runtime will set the user arena to fault once it&#39;s safe (the GC is no longer running)</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span><span class="comment">// and then once the user arena is no longer referenced by the application, will allow it to</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span><span class="comment">// be reused.</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>func freeUserArenaChunk(s *mspan, x unsafe.Pointer) {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	if !s.isUserArenaChunk {
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		throw(&#34;span is not for a user arena&#34;)
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	if s.npages*pageSize != userArenaChunkBytes {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		throw(&#34;invalid user arena span size&#34;)
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	}
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	<span class="comment">// Mark the region as free to various santizers immediately instead</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	<span class="comment">// of handling them at sweep time.</span>
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	if raceenabled {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>		racefree(unsafe.Pointer(s.base()), s.elemsize)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	}
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	if msanenabled {
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>		msanfree(unsafe.Pointer(s.base()), s.elemsize)
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>	}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	if asanenabled {
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		asanpoison(unsafe.Pointer(s.base()), s.elemsize)
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	}
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	<span class="comment">// Make ourselves non-preemptible as we manipulate state and statistics.</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	<span class="comment">// Also required by setUserArenaChunksToFault.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	mp := acquirem()
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	<span class="comment">// We can only set user arenas to fault if we&#39;re in the _GCoff phase.</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	if gcphase == _GCoff {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		lock(&amp;userArenaState.lock)
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		faultList := userArenaState.fault
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>		userArenaState.fault = nil
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		unlock(&amp;userArenaState.lock)
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		s.setUserArenaChunkToFault()
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		for _, lc := range faultList {
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			lc.mspan.setUserArenaChunkToFault()
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>		}
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>		<span class="comment">// Until the chunks are set to fault, keep them alive via the fault list.</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		KeepAlive(x)
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		KeepAlive(faultList)
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	} else {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		<span class="comment">// Put the user arena on the fault list.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		lock(&amp;userArenaState.lock)
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		userArenaState.fault = append(userArenaState.fault, liveUserArenaChunk{s, x})
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		unlock(&amp;userArenaState.lock)
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	}
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	releasem(mp)
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>}
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span><span class="comment">// allocUserArenaChunk attempts to reuse a free user arena chunk represented</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span><span class="comment">// as a span.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span><span class="comment">// Must be in a non-preemptible state to ensure the consistency of statistics</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span><span class="comment">// exported to MemStats.</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span><span class="comment">// Acquires the heap lock. Must run on the system stack for that reason.</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span><span class="comment">//go:systemstack</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>func (h *mheap) allocUserArenaChunk() *mspan {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	var s *mspan
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	var base uintptr
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	<span class="comment">// First check the free list.</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	lock(&amp;h.lock)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	if !h.userArena.readyList.isEmpty() {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		s = h.userArena.readyList.first
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		h.userArena.readyList.remove(s)
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		base = s.base()
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>	} else {
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// Free list was empty, so allocate a new arena.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		hintList := &amp;h.userArena.arenaHints
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		if raceenabled {
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>			<span class="comment">// In race mode just use the regular heap hints. We might fragment</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			<span class="comment">// the address space, but the race detector requires that the heap</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>			<span class="comment">// is mapped contiguously.</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>			hintList = &amp;h.arenaHints
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		}
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		v, size := h.sysAlloc(userArenaChunkBytes, hintList, false)
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		if size%userArenaChunkBytes != 0 {
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>			throw(&#34;sysAlloc size is not divisible by userArenaChunkBytes&#34;)
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>		if size &gt; userArenaChunkBytes {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			<span class="comment">// We got more than we asked for. This can happen if</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			<span class="comment">// heapArenaSize &gt; userArenaChunkSize, or if sysAlloc just returns</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			<span class="comment">// some extra as a result of trying to find an aligned region.</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>			<span class="comment">// Divide it up and put it on the ready list.</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>			for i := userArenaChunkBytes; i &lt; size; i += userArenaChunkBytes {
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>				s := h.allocMSpanLocked()
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>				s.init(uintptr(v)+i, userArenaChunkPages)
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>				h.userArena.readyList.insertBack(s)
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			}
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>			size = userArenaChunkBytes
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		base = uintptr(v)
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>		if base == 0 {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			<span class="comment">// Out of memory.</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			unlock(&amp;h.lock)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>			return nil
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>		}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>		s = h.allocMSpanLocked()
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	}
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	unlock(&amp;h.lock)
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	<span class="comment">// sysAlloc returns Reserved address space, and any span we&#39;re</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	<span class="comment">// reusing is set to fault (so, also Reserved), so transition</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	<span class="comment">// it to Prepared and then Ready.</span>
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	<span class="comment">// Unlike (*mheap).grow, just map in everything that we</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	<span class="comment">// asked for. We&#39;re likely going to use it all.</span>
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	sysMap(unsafe.Pointer(base), userArenaChunkBytes, &amp;gcController.heapReleased)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	sysUsed(unsafe.Pointer(base), userArenaChunkBytes, userArenaChunkBytes)
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	<span class="comment">// Model the user arena as a heap span for a large object.</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	spc := makeSpanClass(0, false)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	h.initSpan(s, spanAllocHeap, spc, base, userArenaChunkPages)
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>	s.isUserArenaChunk = true
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	s.elemsize -= userArenaChunkReserveBytes()
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	s.limit = s.base() + s.elemsize
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>	s.freeindex = 1
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	s.allocCount = 1
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	<span class="comment">// Account for this new arena chunk memory.</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	gcController.heapInUse.add(int64(userArenaChunkBytes))
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	gcController.heapReleased.add(-int64(userArenaChunkBytes))
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	stats := memstats.heapStats.acquire()
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.inHeap, int64(userArenaChunkBytes))
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	atomic.Xaddint64(&amp;stats.committed, int64(userArenaChunkBytes))
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	<span class="comment">// Model the arena as a single large malloc.</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeAlloc, int64(s.elemsize))
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>	atomic.Xadd64(&amp;stats.largeAllocCount, 1)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	memstats.heapStats.release()
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	<span class="comment">// Count the alloc in inconsistent, internal stats.</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	gcController.totalAlloc.Add(int64(s.elemsize))
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	<span class="comment">// Update heapLive.</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>	gcController.update(int64(s.elemsize), 0)
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>	<span class="comment">// This must clear the entire heap bitmap so that it&#39;s safe</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	<span class="comment">// to allocate noscan data without writing anything out.</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	s.initHeapBits(true)
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>	<span class="comment">// Clear the span preemptively. It&#39;s an arena chunk, so let&#39;s assume</span>
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	<span class="comment">// everything is going to be used.</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>	<span class="comment">// This also seems to make a massive difference as to whether or</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>	<span class="comment">// not Linux decides to back this memory with transparent huge</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>	<span class="comment">// pages. There&#39;s latency involved in this zeroing, but the hugepage</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	<span class="comment">// gains are almost always worth it. Note: it&#39;s important that we</span>
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	<span class="comment">// clear even if it&#39;s freshly mapped and we know there&#39;s no point</span>
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	<span class="comment">// to zeroing as *that* is the critical signal to use huge pages.</span>
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	memclrNoHeapPointers(unsafe.Pointer(s.base()), s.elemsize)
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	s.needzero = 0
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	s.freeIndexForScan = 1
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	<span class="comment">// Set up the range for allocation.</span>
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	s.userArenaChunkFree = makeAddrRange(base, base+s.elemsize)
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// Put the large span in the mcentral swept list so that it&#39;s</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	<span class="comment">// visible to the background sweeper.</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	h.central[spc].mcentral.fullSwept(h.sweepgen).push(s)
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	if goexperiment.AllocHeaders {
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>		<span class="comment">// Set up an allocation header. Avoid write barriers here because this type</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>		<span class="comment">// is not a real type, and it exists in an invalid location.</span>
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		*(*uintptr)(unsafe.Pointer(&amp;s.largeType)) = uintptr(unsafe.Pointer(s.limit))
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>		*(*uintptr)(unsafe.Pointer(&amp;s.largeType.GCData)) = s.limit + unsafe.Sizeof(_type{})
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		s.largeType.PtrBytes = 0
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		s.largeType.Size_ = s.elemsize
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	return s
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
</pre><p><a href="arena.go?m=text">View as plain text</a></p>

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
