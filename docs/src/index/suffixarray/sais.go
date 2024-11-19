<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/index/suffixarray/sais.go - Go Documentation Server</title>

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
<a href="sais.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/index">index</a>/<a href="http://localhost:8080/src/index/suffixarray">suffixarray</a>/<span class="text-muted">sais.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/index/suffixarray">index/suffixarray</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2019 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Suffix array construction by induced sorting (SAIS).</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// See Ge Nong, Sen Zhang, and Wai Hong Chen,</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// &#34;Two Efficient Algorithms for Linear Time Suffix Array Construction&#34;,</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// especially section 3 (https://ieeexplore.ieee.org/document/5582081).</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// See also http://zork.net/~st/jottings/sais.html.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// With optimizations inspired by Yuta Mori&#39;s sais-lite</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// (https://sites.google.com/site/yuta256/sais).</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// And with other new optimizations.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Many of these functions are parameterized by the sizes of</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// the types they operate on. The generator gen.go makes</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// copies of these functions for use with other sizes.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Specifically:</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// - A function with a name ending in _8_32 takes []byte and []int32 arguments</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//   and is duplicated into _32_32, _8_64, and _64_64 forms.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//   The _32_32 and _64_64_ suffixes are shortened to plain _32 and _64.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//   Any lines in the function body that contain the text &#34;byte-only&#34; or &#34;256&#34;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//   are stripped when creating _32_32 and _64_64 forms.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//   (Those lines are typically 8-bit-specific optimizations.)</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// - A function with a name ending only in _32 operates on []int32</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//   and is duplicated into a _64 form. (Note that it may still take a []byte,</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//   but there is no need for a version of the function in which the []byte</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//   is widened to a full integer array.)</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// The overall runtime of this code is linear in the input size:</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// it runs a sequence of linear passes to reduce the problem to</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// a subproblem at most half as big, invokes itself recursively,</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// and then runs a sequence of linear passes to turn the answer</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// for the subproblem into the answer for the original problem.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// This gives T(N) = O(N) + T(N/2) = O(N) + O(N/2) + O(N/4) + ... = O(N).</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// The outline of the code, with the forward and backward scans</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// through O(N)-sized arrays called out, is:</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// sais_I_N</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//	placeLMS_I_B</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//		bucketMax_I_B</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text&gt; (1)</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (2)</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -text, random bucket&gt; (3)</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//	induceSubL_I_B</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//		bucketMin_I_B</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text, often optimized away&gt; (4)</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (5)</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//		&lt;scan +sa, random text, random bucket&gt; (6)</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	induceSubS_I_B</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//		bucketMax_I_B</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text, often optimized away&gt; (7)</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (8)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -sa, random text, random bucket&gt; (9)</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//	assignID_I_B</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//		&lt;scan +sa, random text substrings&gt; (10)</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//	map_B</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -sa&gt; (11)</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//	recurse_B</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//		(recursive call to sais_B_B for a subproblem of size at most 1/2 input, often much smaller)</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//	unmap_I_B</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -text&gt; (12)</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//		&lt;scan +sa&gt; (13)</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//	expand_I_B</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//		bucketMax_I_B</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text, often optimized away&gt; (14)</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (15)</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -sa, random text, random bucket&gt; (16)</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//	induceL_I_B</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//		bucketMin_I_B</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text, often optimized away&gt; (17)</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (18)</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//		&lt;scan +sa, random text, random bucket&gt; (19)</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//	induceS_I_B</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">//		bucketMax_I_B</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">//			freq_I_B</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//				&lt;scan +text, often optimized away&gt; (20)</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//			&lt;scan +freq&gt; (21)</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//		&lt;scan -sa, random text, random bucket&gt; (22)</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// Here, _B indicates the suffix array size (_32 or _64) and _I the input size (_8 or _B).</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// The outline shows there are in general 22 scans through</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// O(N)-sized arrays for a given level of the recursion.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// In the top level, operating on 8-bit input text,</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// the six freq scans are fixed size (256) instead of potentially</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// input-sized. Also, the frequency is counted once and cached</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// whenever there is room to do so (there is nearly always room in general,</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// and always room at the top level), which eliminates all but</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// the first freq_I_B text scans (that is, 5 of the 6).</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// So the top level of the recursion only does 22 - 6 - 5 = 11</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// input-sized scans and a typical level does 16 scans.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// The linear scans do not cost anywhere near as much as</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// the random accesses to the text made during a few of</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// the scans (specifically #6, #9, #16, #19, #22 marked above).</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// In real texts, there is not much but some locality to</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// the accesses, due to the repetitive structure of the text</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// (the same reason Burrows-Wheeler compression is so effective).</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// For random inputs, there is no locality, which makes those</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// accesses even more expensive, especially once the text</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// no longer fits in cache.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// For example, running on 50 MB of Go source code, induceSubL_8_32</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// (which runs only once, at the top level of the recursion)</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// takes 0.44s, while on 50 MB of random input, it takes 2.55s.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// Nearly all the relative slowdown is explained by the text access:</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span><span class="comment">//		c0, c1 := text[k-1], text[k]</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// That line runs for 0.23s on the Go text and 2.02s on random text.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//go:generate go run gen.go</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>package suffixarray
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// text_32 returns the suffix array for the input text.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// It requires that len(text) fit in an int32</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// and that the caller zero sa.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>func text_32(text []byte, sa []int32) {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if int(int32(len(text))) != len(text) || len(text) != len(sa) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		panic(&#34;suffixarray: misuse of text_32&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	sais_8_32(text, 256, sa, make([]int32, 2*256))
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// sais_8_32 computes the suffix array of text.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// The text must contain only values in [0, textMax).</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// The suffix array is stored in sa, which the caller</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// must ensure is already zeroed.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// The caller must also provide temporary space tmp</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// with len(tmp) ≥ textMax. If len(tmp) ≥ 2*textMax</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// then the algorithm runs a little faster.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// If sais_8_32 modifies tmp, it sets tmp[0] = -1 on return.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func sais_8_32(text []byte, textMax int, sa, tmp []int32) {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if len(sa) != len(text) || len(tmp) &lt; textMax {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		panic(&#34;suffixarray: misuse of sais_8_32&#34;)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">// Trivial base cases. Sorting 0 or 1 things is easy.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if len(text) == 0 {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		return
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if len(text) == 1 {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		sa[0] = 0
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">// Establish slices indexed by text character</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">// holding character frequency and bucket-sort offsets.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">// If there&#39;s only enough tmp for one slice,</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// we make it the bucket offsets and recompute</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// the character frequency each time we need it.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	var freq, bucket []int32
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if len(tmp) &gt;= 2*textMax {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		freq, bucket = tmp[:textMax], tmp[textMax:2*textMax]
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		freq[0] = -1 <span class="comment">// mark as uninitialized</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	} else {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		freq, bucket = nil, tmp[:textMax]
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">// The SAIS algorithm.</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">// Each of these calls makes one scan through sa.</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// See the individual functions for documentation</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">// about each&#39;s role in the algorithm.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	numLMS := placeLMS_8_32(text, sa, freq, bucket)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if numLMS &lt;= 1 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		<span class="comment">// 0 or 1 items are already sorted. Do nothing.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	} else {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		induceSubL_8_32(text, sa, freq, bucket)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		induceSubS_8_32(text, sa, freq, bucket)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		length_8_32(text, sa, numLMS)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		maxID := assignID_8_32(text, sa, numLMS)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if maxID &lt; numLMS {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			map_32(sa, numLMS)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			recurse_32(sa, tmp, numLMS, maxID)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			unmap_8_32(text, sa, numLMS)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		} else {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// If maxID == numLMS, then each LMS-substring</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// is unique, so the relative ordering of two LMS-suffixes</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">// is determined by just the leading LMS-substring.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			<span class="comment">// That is, the LMS-suffix sort order matches the</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>			<span class="comment">// (simpler) LMS-substring sort order.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			<span class="comment">// Copy the original LMS-substring order into the</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			<span class="comment">// suffix array destination.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			copy(sa, sa[len(sa)-numLMS:])
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		expand_8_32(text, freq, bucket, sa, numLMS)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	induceL_8_32(text, sa, freq, bucket)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	induceS_8_32(text, sa, freq, bucket)
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	<span class="comment">// Mark for caller that we overwrote tmp.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	tmp[0] = -1
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// freq_8_32 returns the character frequencies</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// for text, as a slice indexed by character value.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// If freq is nil, freq_8_32 uses and returns bucket.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// If freq is non-nil, freq_8_32 assumes that freq[0] &gt;= 0</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// means the frequencies are already computed.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">// If the frequency data is overwritten or uninitialized,</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// the caller must set freq[0] = -1 to force recomputation</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// the next time it is needed.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>func freq_8_32(text []byte, freq, bucket []int32) []int32 {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if freq != nil &amp;&amp; freq[0] &gt;= 0 {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		return freq <span class="comment">// already computed</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if freq == nil {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		freq = bucket
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	freq = freq[:256] <span class="comment">// eliminate bounds check for freq[c] below</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	for i := range freq {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		freq[i] = 0
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	for _, c := range text {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		freq[c]++
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	return freq
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span><span class="comment">// bucketMin_8_32 stores into bucket[c] the minimum index</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// in the bucket for character c in a bucket-sort of text.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func bucketMin_8_32(text []byte, freq, bucket []int32) {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	freq = freq_8_32(text, freq, bucket)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	freq = freq[:256]     <span class="comment">// establish len(freq) = 256, so 0 ≤ i &lt; 256 below</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[i] below</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	total := int32(0)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	for i, n := range freq {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		bucket[i] = total
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		total += n
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// bucketMax_8_32 stores into bucket[c] the maximum index</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// in the bucket for character c in a bucket-sort of text.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span><span class="comment">// The bucket indexes for c are [min, max).</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span><span class="comment">// That is, max is one past the final index in that bucket.</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>func bucketMax_8_32(text []byte, freq, bucket []int32) {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	freq = freq_8_32(text, freq, bucket)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	freq = freq[:256]     <span class="comment">// establish len(freq) = 256, so 0 ≤ i &lt; 256 below</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[i] below</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	total := int32(0)
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	for i, n := range freq {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		total += n
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		bucket[i] = total
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span><span class="comment">// The SAIS algorithm proceeds in a sequence of scans through sa.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// Each of the following functions implements one scan,</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span><span class="comment">// and the functions appear here in the order they execute in the algorithm.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span><span class="comment">// placeLMS_8_32 places into sa the indexes of the</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span><span class="comment">// final characters of the LMS substrings of text,</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// sorted into the rightmost ends of their correct buckets</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span><span class="comment">// in the suffix array.</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span><span class="comment">// The imaginary sentinel character at the end of the text</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// is the final character of the final LMS substring, but there</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span><span class="comment">// is no bucket for the imaginary sentinel character,</span>
<span id="L271" class="ln">   271&nbsp;&nbsp;</span><span class="comment">// which has a smaller value than any real character.</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span><span class="comment">// The caller must therefore pretend that sa[-1] == len(text).</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span><span class="comment">// The text indexes of LMS-substring characters are always ≥ 1</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span><span class="comment">// (the first LMS-substring must be preceded by one or more L-type</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// characters that are not part of any LMS-substring),</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// so using 0 as a “not present” suffix array entry is safe,</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// both in this function and in most later functions</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// (until induceL_8_32 below).</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>func placeLMS_8_32(text []byte, sa, freq, bucket []int32) int {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	bucketMax_8_32(text, freq, bucket)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	numLMS := 0
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	lastB := int32(-1)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[c1] below</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	<span class="comment">// The next stanza of code (until the blank line) loop backward</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	<span class="comment">// over text, stopping to execute a code body at each position i</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// such that text[i] is an L-character and text[i+1] is an S-character.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	<span class="comment">// That is, i+1 is the position of the start of an LMS-substring.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	<span class="comment">// These could be hoisted out into a function with a callback,</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// but at a significant speed cost. Instead, we just write these</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// seven lines a few times in this source file. The copies below</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// refer back to the pattern established by this original as the</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// &#34;LMS-substring iterator&#34;.</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// In every scan through the text, c0, c1 are successive characters of text.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	<span class="comment">// In this backward scan, c0 == text[i] and c1 == text[i+1].</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	<span class="comment">// By scanning backward, we can keep track of whether the current</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	<span class="comment">// position is type-S or type-L according to the usual definition:</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	<span class="comment">//	- position len(text) is type S with text[len(text)] == -1 (the sentinel)</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">//	- position i is type S if text[i] &lt; text[i+1], or if text[i] == text[i+1] &amp;&amp; i+1 is type S.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	<span class="comment">//	- position i is type L if text[i] &gt; text[i+1], or if text[i] == text[i+1] &amp;&amp; i+1 is type L.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	<span class="comment">// The backward scan lets us maintain the current type,</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	<span class="comment">// update it when we see c0 != c1, and otherwise leave it alone.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	<span class="comment">// We want to identify all S positions with a preceding L.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// Position len(text) is one such position by definition, but we have</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	<span class="comment">// nowhere to write it down, so we eliminate it by untruthfully</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	<span class="comment">// setting isTypeS = false at the start of the loop.</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	c0, c1, isTypeS := byte(0), byte(0), false
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	for i := len(text) - 1; i &gt;= 0; i-- {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		c0, c1 = text[i], c0
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		if c0 &lt; c1 {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			isTypeS = true
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		} else if c0 &gt; c1 &amp;&amp; isTypeS {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>			isTypeS = false
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			<span class="comment">// Bucket the index i+1 for the start of an LMS-substring.</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			b := bucket[c1] - 1
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			bucket[c1] = b
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			sa[b] = int32(i + 1)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			lastB = b
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			numLMS++
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	<span class="comment">// We recorded the LMS-substring starts but really want the ends.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// Luckily, with two differences, the start indexes and the end indexes are the same.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// The first difference is that the rightmost LMS-substring&#39;s end index is len(text),</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// so the caller must pretend that sa[-1] == len(text), as noted above.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	<span class="comment">// The second difference is that the first leftmost LMS-substring start index</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	<span class="comment">// does not end an earlier LMS-substring, so as an optimization we can omit</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// that leftmost LMS-substring start index (the last one we wrote).</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	<span class="comment">// Exception: if numLMS &lt;= 1, the caller is not going to bother with</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// the recursion at all and will treat the result as containing LMS-substring starts.</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	<span class="comment">// In that case, we don&#39;t remove the final entry.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	if numLMS &gt; 1 {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		sa[lastB] = 0
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	return numLMS
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">// induceSubL_8_32 inserts the L-type text indexes of LMS-substrings</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// into sa, assuming that the final characters of the LMS-substrings</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// are already inserted into sa, sorted by final character, and at the</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// right (not left) end of the corresponding character bucket.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// Each LMS-substring has the form (as a regexp) /S+L+S/:</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// one or more S-type, one or more L-type, final S-type.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">// induceSubL_8_32 leaves behind only the leftmost L-type text</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">// index for each LMS-substring. That is, it removes the final S-type</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">// indexes that are present on entry, and it inserts but then removes</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">// the interior L-type indexes too.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">// (Only the leftmost L-type index is needed by induceSubS_8_32.)</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>func induceSubL_8_32(text []byte, sa, freq, bucket []int32) {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// Initialize positions for left side of character buckets.</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	bucketMin_8_32(text, freq, bucket)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[cB] below</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	<span class="comment">// As we scan the array left-to-right, each sa[i] = j &gt; 0 is a correctly</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	<span class="comment">// sorted suffix array entry (for text[j:]) for which we know that j-1 is type L.</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	<span class="comment">// Because j-1 is type L, inserting it into sa now will sort it correctly.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// But we want to distinguish a j-1 with j-2 of type L from type S.</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	<span class="comment">// We can process the former but want to leave the latter for the caller.</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	<span class="comment">// We record the difference by negating j-1 if it is preceded by type S.</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	<span class="comment">// Either way, the insertion (into the text[j-1] bucket) is guaranteed to</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// happen at sa[i´] for some i´ &gt; i, that is, in the portion of sa we have</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// yet to scan. A single pass therefore sees indexes j, j-1, j-2, j-3,</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	<span class="comment">// and so on, in sorted but not necessarily adjacent order, until it finds</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// one preceded by an index of type S, at which point it must stop.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	<span class="comment">// As we scan through the array, we clear the worked entries (sa[i] &gt; 0) to zero,</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	<span class="comment">// and we flip sa[i] &lt; 0 to -sa[i], so that the loop finishes with sa containing</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	<span class="comment">// only the indexes of the leftmost L-type indexes for each LMS-substring.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	<span class="comment">// The suffix array sa therefore serves simultaneously as input, output,</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// and a miraculously well-tailored work queue.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	<span class="comment">// placeLMS_8_32 left out the implicit entry sa[-1] == len(text),</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	<span class="comment">// corresponding to the identified type-L index len(text)-1.</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	<span class="comment">// Process it before the left-to-right scan of sa proper.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// See body in loop for commentary.</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	k := len(text) - 1
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	c0, c1 := text[k-1], text[k]
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	if c0 &lt; c1 {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>		k = -k
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	<span class="comment">// Cache recently used bucket index:</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	<span class="comment">// we&#39;re processing suffixes in sorted order</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	<span class="comment">// and accessing buckets indexed by the</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	<span class="comment">// byte before the sorted order, which still</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	<span class="comment">// has very good locality.</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	<span class="comment">// Invariant: b is cached, possibly dirty copy of bucket[cB].</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	cB := c1
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	b := bucket[cB]
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	sa[b] = int32(k)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	b++
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	for i := 0; i &lt; len(sa); i++ {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		j := int(sa[i])
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		if j == 0 {
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			<span class="comment">// Skip empty entry.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			continue
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		if j &lt; 0 {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			<span class="comment">// Leave discovered type-S index for caller.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			sa[i] = int32(-j)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			continue
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		}
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		sa[i] = 0
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		<span class="comment">// Index j was on work queue, meaning k := j-1 is L-type,</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		<span class="comment">// so we can now place k correctly into sa.</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is L-type, queue k for processing later in this loop.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is S-type (text[k-1] &lt; text[k]), queue -k to save for the caller.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		k := j - 1
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		c0, c1 := text[k-1], text[k]
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		if c0 &lt; c1 {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			k = -k
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		if cB != c1 {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			bucket[cB] = b
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			cB = c1
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			b = bucket[cB]
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		sa[b] = int32(k)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		b++
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>}
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span><span class="comment">// induceSubS_8_32 inserts the S-type text indexes of LMS-substrings</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span><span class="comment">// into sa, assuming that the leftmost L-type text indexes are already</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span><span class="comment">// inserted into sa, sorted by LMS-substring suffix, and at the</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span><span class="comment">// left end of the corresponding character bucket.</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">// Each LMS-substring has the form (as a regexp) /S+L+S/:</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">// one or more S-type, one or more L-type, final S-type.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// induceSubS_8_32 leaves behind only the leftmost S-type text</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span><span class="comment">// index for each LMS-substring, in sorted order, at the right end of sa.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span><span class="comment">// That is, it removes the L-type indexes that are present on entry,</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span><span class="comment">// and it inserts but then removes the interior S-type indexes too,</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span><span class="comment">// leaving the LMS-substring start indexes packed into sa[len(sa)-numLMS:].</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// (Only the LMS-substring start indexes are processed by the recursion.)</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func induceSubS_8_32(text []byte, sa, freq, bucket []int32) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// Initialize positions for right side of character buckets.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	bucketMax_8_32(text, freq, bucket)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[cB] below</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	<span class="comment">// Analogous to induceSubL_8_32 above,</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	<span class="comment">// as we scan the array right-to-left, each sa[i] = j &gt; 0 is a correctly</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	<span class="comment">// sorted suffix array entry (for text[j:]) for which we know that j-1 is type S.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	<span class="comment">// Because j-1 is type S, inserting it into sa now will sort it correctly.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	<span class="comment">// But we want to distinguish a j-1 with j-2 of type S from type L.</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	<span class="comment">// We can process the former but want to leave the latter for the caller.</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// We record the difference by negating j-1 if it is preceded by type L.</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	<span class="comment">// Either way, the insertion (into the text[j-1] bucket) is guaranteed to</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	<span class="comment">// happen at sa[i´] for some i´ &lt; i, that is, in the portion of sa we have</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// yet to scan. A single pass therefore sees indexes j, j-1, j-2, j-3,</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	<span class="comment">// and so on, in sorted but not necessarily adjacent order, until it finds</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	<span class="comment">// one preceded by an index of type L, at which point it must stop.</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	<span class="comment">// That index (preceded by one of type L) is an LMS-substring start.</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	<span class="comment">// As we scan through the array, we clear the worked entries (sa[i] &gt; 0) to zero,</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// and we flip sa[i] &lt; 0 to -sa[i] and compact into the top of sa,</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	<span class="comment">// so that the loop finishes with the top of sa containing exactly</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	<span class="comment">// the LMS-substring start indexes, sorted by LMS-substring.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	<span class="comment">// Cache recently used bucket index:</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>	cB := byte(0)
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	b := bucket[cB]
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>	top := len(sa)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	for i := len(sa) - 1; i &gt;= 0; i-- {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>		j := int(sa[i])
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		if j == 0 {
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			<span class="comment">// Skip empty entry.</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			continue
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		}
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		sa[i] = 0
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		if j &lt; 0 {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>			<span class="comment">// Leave discovered LMS-substring start index for caller.</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			top--
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			sa[top] = int32(-j)
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			continue
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		<span class="comment">// Index j was on work queue, meaning k := j-1 is S-type,</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		<span class="comment">// so we can now place k correctly into sa.</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is S-type, queue k for processing later in this loop.</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is L-type (text[k-1] &gt; text[k]), queue -k to save for the caller.</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		k := j - 1
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		c1 := text[k]
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		c0 := text[k-1]
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		if c0 &gt; c1 {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			k = -k
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		if cB != c1 {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			bucket[cB] = b
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			cB = c1
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			b = bucket[cB]
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		b--
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		sa[b] = int32(k)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>}
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span><span class="comment">// length_8_32 computes and records the length of each LMS-substring in text.</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span><span class="comment">// The length of the LMS-substring at index j is stored at sa[j/2],</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span><span class="comment">// avoiding the LMS-substring indexes already stored in the top half of sa.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span><span class="comment">// (If index j is an LMS-substring start, then index j-1 is type L and cannot be.)</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span><span class="comment">// There are two exceptions, made for optimizations in name_8_32 below.</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// First, the final LMS-substring is recorded as having length 0, which is otherwise</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">// impossible, instead of giving it a length that includes the implicit sentinel.</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">// This ensures the final LMS-substring has length unequal to all others</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">// and therefore can be detected as different without text comparison</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span><span class="comment">// (it is unequal because it is the only one that ends in the implicit sentinel,</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// and the text comparison would be problematic since the implicit sentinel</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// is not actually present at text[len(text)]).</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">// Second, to avoid text comparison entirely, if an LMS-substring is very short,</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">// sa[j/2] records its actual text instead of its length, so that if two such</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span><span class="comment">// substrings have matching “length,” the text need not be read at all.</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">// The definition of “very short” is that the text bytes must pack into a uint32,</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// and the unsigned encoding e must be ≥ len(text), so that it can be</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span><span class="comment">// distinguished from a valid length.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>func length_8_32(text []byte, sa []int32, numLMS int) {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	end := 0 <span class="comment">// index of current LMS-substring end (0 indicates final LMS-substring)</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	<span class="comment">// The encoding of N text bytes into a “length” word</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>	<span class="comment">// adds 1 to each byte, packs them into the bottom</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// N*8 bits of a word, and then bitwise inverts the result.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// That is, the text sequence A B C (hex 41 42 43)</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	<span class="comment">// encodes as ^uint32(0x42_43_44).</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	<span class="comment">// LMS-substrings can never start or end with 0xFF.</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>	<span class="comment">// Adding 1 ensures the encoded byte sequence never</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>	<span class="comment">// starts or ends with 0x00, so that present bytes can be</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>	<span class="comment">// distinguished from zero-padding in the top bits,</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	<span class="comment">// so the length need not be separately encoded.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>	<span class="comment">// Inverting the bytes increases the chance that a</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	<span class="comment">// 4-byte encoding will still be ≥ len(text).</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	<span class="comment">// In particular, if the first byte is ASCII (&lt;= 0x7E, so +1 &lt;= 0x7F)</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	<span class="comment">// then the high bit of the inversion will be set,</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	<span class="comment">// making it clearly not a valid length (it would be a negative one).</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	<span class="comment">// cx holds the pre-inverted encoding (the packed incremented bytes).</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	cx := uint32(0) <span class="comment">// byte-only</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	<span class="comment">// This stanza (until the blank line) is the &#34;LMS-substring iterator&#34;,</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// described in placeLMS_8_32 above, with one line added to maintain cx.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	c0, c1, isTypeS := byte(0), byte(0), false
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	for i := len(text) - 1; i &gt;= 0; i-- {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		c0, c1 = text[i], c0
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		cx = cx&lt;&lt;8 | uint32(c1+1) <span class="comment">// byte-only</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		if c0 &lt; c1 {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>			isTypeS = true
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		} else if c0 &gt; c1 &amp;&amp; isTypeS {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			isTypeS = false
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			<span class="comment">// Index j = i+1 is the start of an LMS-substring.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			<span class="comment">// Compute length or encoded text to store in sa[j/2].</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>			j := i + 1
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>			var code int32
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			if end == 0 {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>				code = 0
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>			} else {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>				code = int32(end - j)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>				if code &lt;= 32/8 &amp;&amp; ^cx &gt;= uint32(len(text)) { <span class="comment">// byte-only</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>					code = int32(^cx) <span class="comment">// byte-only</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>				} <span class="comment">// byte-only</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			sa[j&gt;&gt;1] = code
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>			end = j + 1
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			cx = uint32(c1 + 1) <span class="comment">// byte-only</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>}
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span><span class="comment">// assignID_8_32 assigns a dense ID numbering to the</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span><span class="comment">// set of LMS-substrings respecting string ordering and equality,</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span><span class="comment">// returning the maximum assigned ID.</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span><span class="comment">// For example given the input &#34;ababab&#34;, the LMS-substrings</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span><span class="comment">// are &#34;aba&#34;, &#34;aba&#34;, and &#34;ab&#34;, renumbered as 2 2 1.</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// sa[len(sa)-numLMS:] holds the LMS-substring indexes</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// sorted in string order, so to assign numbers we can</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// consider each in turn, removing adjacent duplicates.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// The new ID for the LMS-substring at index j is written to sa[j/2],</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span><span class="comment">// overwriting the length previously stored there (by length_8_32 above).</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>func assignID_8_32(text []byte, sa []int32, numLMS int) int {
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	id := 0
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	lastLen := int32(-1) <span class="comment">// impossible</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	lastPos := int32(0)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	for _, j := range sa[len(sa)-numLMS:] {
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		<span class="comment">// Is the LMS-substring at index j new, or is it the same as the last one we saw?</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		n := sa[j/2]
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		if n != lastLen {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			goto New
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		if uint32(n) &gt;= uint32(len(text)) {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			<span class="comment">// “Length” is really encoded full text, and they match.</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			goto Same
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		}
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		{
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>			<span class="comment">// Compare actual texts.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>			n := int(n)
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>			this := text[j:][:n]
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>			last := text[lastPos:][:n]
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			for i := 0; i &lt; n; i++ {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>				if this[i] != last[i] {
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>					goto New
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>				}
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			}
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>			goto Same
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	New:
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		id++
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>		lastPos = j
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		lastLen = n
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	Same:
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>		sa[j/2] = int32(id)
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	}
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	return id
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>}
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span><span class="comment">// map_32 maps the LMS-substrings in text to their new IDs,</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span><span class="comment">// producing the subproblem for the recursion.</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// The mapping itself was mostly applied by assignID_8_32:</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">// sa[i] is either 0, the ID for the LMS-substring at index 2*i,</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// or the ID for the LMS-substring at index 2*i+1.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// To produce the subproblem we need only remove the zeros</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// and change ID into ID-1 (our IDs start at 1, but text chars start at 0).</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L637" class="ln">   637&nbsp;&nbsp;</span><span class="comment">// map_32 packs the result, which is the input to the recursion,</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span><span class="comment">// into the top of sa, so that the recursion result can be stored</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span><span class="comment">// in the bottom of sa, which sets up for expand_8_32 well.</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>func map_32(sa []int32, numLMS int) {
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	w := len(sa)
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	for i := len(sa) / 2; i &gt;= 0; i-- {
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>		j := sa[i]
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		if j &gt; 0 {
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>			w--
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>			sa[w] = j - 1
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		}
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	}
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>}
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span><span class="comment">// recurse_32 calls sais_32 recursively to solve the subproblem we&#39;ve built.</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span><span class="comment">// The subproblem is at the right end of sa, the suffix array result will be</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span><span class="comment">// written at the left end of sa, and the middle of sa is available for use as</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span><span class="comment">// temporary frequency and bucket storage.</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>func recurse_32(sa, oldTmp []int32, numLMS, maxID int) {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	dst, saTmp, text := sa[:numLMS], sa[numLMS:len(sa)-numLMS], sa[len(sa)-numLMS:]
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// Set up temporary space for recursive call.</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// We must pass sais_32 a tmp buffer with at least maxID entries.</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// The subproblem is guaranteed to have length at most len(sa)/2,</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// so that sa can hold both the subproblem and its suffix array.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	<span class="comment">// Nearly all the time, however, the subproblem has length &lt; len(sa)/3,</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	<span class="comment">// in which case there is a subproblem-sized middle of sa that</span>
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>	<span class="comment">// we can reuse for temporary space (saTmp).</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>	<span class="comment">// When recurse_32 is called from sais_8_32, oldTmp is length 512</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>	<span class="comment">// (from text_32), and saTmp will typically be much larger, so we&#39;ll use saTmp.</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	<span class="comment">// When deeper recursions come back to recurse_32, now oldTmp is</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	<span class="comment">// the saTmp from the top-most recursion, it is typically larger than</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	<span class="comment">// the current saTmp (because the current sa gets smaller and smaller</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>	<span class="comment">// as the recursion gets deeper), and we keep reusing that top-most</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	<span class="comment">// large saTmp instead of the offered smaller ones.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>	<span class="comment">// Why is the subproblem length so often just under len(sa)/3?</span>
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	<span class="comment">// See Nong, Zhang, and Chen, section 3.6 for a plausible explanation.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	<span class="comment">// In brief, the len(sa)/2 case would correspond to an SLSLSLSLSLSL pattern</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>	<span class="comment">// in the input, perfect alternation of larger and smaller input bytes.</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>	<span class="comment">// Real text doesn&#39;t do that. If each L-type index is randomly followed</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	<span class="comment">// by either an L-type or S-type index, then half the substrings will</span>
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	<span class="comment">// be of the form SLS, but the other half will be longer. Of that half,</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	<span class="comment">// half (a quarter overall) will be SLLS; an eighth will be SLLLS, and so on.</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>	<span class="comment">// Not counting the final S in each (which overlaps the first S in the next),</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	<span class="comment">// This works out to an average length 2×½ + 3×¼ + 4×⅛ + ... = 3.</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>	<span class="comment">// The space we need is further reduced by the fact that many of the</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	<span class="comment">// short patterns like SLS will often be the same character sequences</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>	<span class="comment">// repeated throughout the text, reducing maxID relative to numLMS.</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>	<span class="comment">// For short inputs, the averages may not run in our favor, but then we</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>	<span class="comment">// can often fall back to using the length-512 tmp available in the</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>	<span class="comment">// top-most call. (Also a short allocation would not be a big deal.)</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>	<span class="comment">// For pathological inputs, we fall back to allocating a new tmp of length</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>	<span class="comment">// max(maxID, numLMS/2). This level of the recursion needs maxID,</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>	<span class="comment">// and all deeper levels of the recursion will need no more than numLMS/2,</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>	<span class="comment">// so this one allocation is guaranteed to suffice for the entire stack</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	<span class="comment">// of recursive calls.</span>
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>	tmp := oldTmp
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>	if len(tmp) &lt; len(saTmp) {
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		tmp = saTmp
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>	}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	if len(tmp) &lt; numLMS {
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		<span class="comment">// TestSAIS/forcealloc reaches this code.</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		n := maxID
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		if n &lt; numLMS/2 {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			n = numLMS / 2
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		tmp = make([]int32, n)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	<span class="comment">// sais_32 requires that the caller arrange to clear dst,</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>	<span class="comment">// because in general the caller may know dst is</span>
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	<span class="comment">// freshly-allocated and already cleared. But this one is not.</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	for i := range dst {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		dst[i] = 0
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>	sais_32(text, maxID, dst, tmp)
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>}
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>
<span id="L719" class="ln">   719&nbsp;&nbsp;</span><span class="comment">// unmap_8_32 unmaps the subproblem back to the original.</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span><span class="comment">// sa[:numLMS] is the LMS-substring numbers, which don&#39;t matter much anymore.</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span><span class="comment">// sa[len(sa)-numLMS:] is the sorted list of those LMS-substring numbers.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span><span class="comment">// The key part is that if the list says K that means the K&#39;th substring.</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span><span class="comment">// We can replace sa[:numLMS] with the indexes of the LMS-substrings.</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span><span class="comment">// Then if the list says K it really means sa[K].</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span><span class="comment">// Having mapped the list back to LMS-substring indexes,</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span><span class="comment">// we can place those into the right buckets.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>func unmap_8_32(text []byte, sa []int32, numLMS int) {
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>	unmap := sa[len(sa)-numLMS:]
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>	j := len(unmap)
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// &#34;LMS-substring iterator&#34; (see placeLMS_8_32 above).</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	c0, c1, isTypeS := byte(0), byte(0), false
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	for i := len(text) - 1; i &gt;= 0; i-- {
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		c0, c1 = text[i], c0
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		if c0 &lt; c1 {
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			isTypeS = true
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		} else if c0 &gt; c1 &amp;&amp; isTypeS {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>			isTypeS = false
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			<span class="comment">// Populate inverse map.</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>			j--
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			unmap[j] = int32(i + 1)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		}
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	<span class="comment">// Apply inverse map to subproblem suffix array.</span>
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>	sa = sa[:numLMS]
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	for i := 0; i &lt; len(sa); i++ {
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		sa[i] = unmap[sa[i]]
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>	}
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>}
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// expand_8_32 distributes the compacted, sorted LMS-suffix indexes</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span><span class="comment">// from sa[:numLMS] into the tops of the appropriate buckets in sa,</span>
<span id="L755" class="ln">   755&nbsp;&nbsp;</span><span class="comment">// preserving the sorted order and making room for the L-type indexes</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span><span class="comment">// to be slotted into the sorted sequence by induceL_8_32.</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>func expand_8_32(text []byte, freq, bucket, sa []int32, numLMS int) {
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	bucketMax_8_32(text, freq, bucket)
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bound check for bucket[c] below</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>	<span class="comment">// Loop backward through sa, always tracking</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	<span class="comment">// the next index to populate from sa[:numLMS].</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	<span class="comment">// When we get to one, populate it.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	<span class="comment">// Zero the rest of the slots; they have dead values in them.</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	x := numLMS - 1
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	saX := sa[x]
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	c := text[saX]
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>	b := bucket[c] - 1
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>	bucket[c] = b
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>	for i := len(sa) - 1; i &gt;= 0; i-- {
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>		if i != int(b) {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>			sa[i] = 0
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			continue
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>		sa[i] = saX
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		<span class="comment">// Load next entry to put down (if any).</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		if x &gt; 0 {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			x--
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			saX = sa[x] <span class="comment">// TODO bounds check</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>			c = text[saX]
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			b = bucket[c] - 1
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>			bucket[c] = b
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		}
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	}
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>}
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span><span class="comment">// induceL_8_32 inserts L-type text indexes into sa,</span>
<span id="L790" class="ln">   790&nbsp;&nbsp;</span><span class="comment">// assuming that the leftmost S-type indexes are inserted</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span><span class="comment">// into sa, in sorted order, in the right bucket halves.</span>
<span id="L792" class="ln">   792&nbsp;&nbsp;</span><span class="comment">// It leaves all the L-type indexes in sa, but the</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span><span class="comment">// leftmost L-type indexes are negated, to mark them</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span><span class="comment">// for processing by induceS_8_32.</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>func induceL_8_32(text []byte, sa, freq, bucket []int32) {
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">// Initialize positions for left side of character buckets.</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	bucketMin_8_32(text, freq, bucket)
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[cB] below</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	<span class="comment">// This scan is similar to the one in induceSubL_8_32 above.</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	<span class="comment">// That one arranges to clear all but the leftmost L-type indexes.</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>	<span class="comment">// This scan leaves all the L-type indexes and the original S-type</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	<span class="comment">// indexes, but it negates the positive leftmost L-type indexes</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>	<span class="comment">// (the ones that induceS_8_32 needs to process).</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>	<span class="comment">// expand_8_32 left out the implicit entry sa[-1] == len(text),</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>	<span class="comment">// corresponding to the identified type-L index len(text)-1.</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	<span class="comment">// Process it before the left-to-right scan of sa proper.</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>	<span class="comment">// See body in loop for commentary.</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>	k := len(text) - 1
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>	c0, c1 := text[k-1], text[k]
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	if c0 &lt; c1 {
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		k = -k
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	<span class="comment">// Cache recently used bucket index.</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	cB := c1
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>	b := bucket[cB]
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>	sa[b] = int32(k)
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	b++
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	for i := 0; i &lt; len(sa); i++ {
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		j := int(sa[i])
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		if j &lt;= 0 {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			<span class="comment">// Skip empty or negated entry (including negated zero).</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>			continue
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		}
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// Index j was on work queue, meaning k := j-1 is L-type,</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>		<span class="comment">// so we can now place k correctly into sa.</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is L-type, queue k for processing later in this loop.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is S-type (text[k-1] &lt; text[k]), queue -k to save for the caller.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		<span class="comment">// If k is zero, k-1 doesn&#39;t exist, so we only need to leave it</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">// for the caller. The caller can&#39;t tell the difference between</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		<span class="comment">// an empty slot and a non-empty zero, but there&#39;s no need</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		<span class="comment">// to distinguish them anyway: the final suffix array will end up</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		<span class="comment">// with one zero somewhere, and that will be a real zero.</span>
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		k := j - 1
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		c1 := text[k]
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>		if k &gt; 0 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			if c0 := text[k-1]; c0 &lt; c1 {
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>				k = -k
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			}
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		}
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>		if cB != c1 {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>			bucket[cB] = b
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			cB = c1
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>			b = bucket[cB]
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>		sa[b] = int32(k)
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		b++
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>}
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>func induceS_8_32(text []byte, sa, freq, bucket []int32) {
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">// Initialize positions for right side of character buckets.</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	bucketMax_8_32(text, freq, bucket)
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	bucket = bucket[:256] <span class="comment">// eliminate bounds check for bucket[cB] below</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	cB := byte(0)
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>	b := bucket[cB]
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	for i := len(sa) - 1; i &gt;= 0; i-- {
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		j := int(sa[i])
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		if j &gt;= 0 {
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>			<span class="comment">// Skip non-flagged entry.</span>
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			<span class="comment">// (This loop can&#39;t see an empty entry; 0 means the real zero index.)</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			continue
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>		}
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>		<span class="comment">// Negative j is a work queue entry; rewrite to positive j for final suffix array.</span>
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>		j = -j
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		sa[i] = int32(j)
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>		<span class="comment">// Index j was on work queue (encoded as -j but now decoded),</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>		<span class="comment">// meaning k := j-1 is L-type,</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>		<span class="comment">// so we can now place k correctly into sa.</span>
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is S-type, queue -k for processing later in this loop.</span>
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>		<span class="comment">// If k-1 is L-type (text[k-1] &gt; text[k]), queue k to save for the caller.</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>		<span class="comment">// If k is zero, k-1 doesn&#39;t exist, so we only need to leave it</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>		<span class="comment">// for the caller.</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>		k := j - 1
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>		c1 := text[k]
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		if k &gt; 0 {
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>			if c0 := text[k-1]; c0 &lt;= c1 {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>				k = -k
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>			}
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		}
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>		if cB != c1 {
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			bucket[cB] = b
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>			cB = c1
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>			b = bucket[cB]
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		}
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		b--
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>		sa[b] = int32(k)
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>	}
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>}
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
</pre><p><a href="sais.go?m=text">View as plain text</a></p>

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
