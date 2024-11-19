<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strings/search.go - Go Documentation Server</title>

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
<a href="search.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strings">strings</a>/<span class="text-muted">search.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/strings">strings</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package strings
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// stringFinder efficiently finds strings in a source text. It&#39;s implemented</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// using the Boyer-Moore string search algorithm:</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// https://en.wikipedia.org/wiki/Boyer-Moore_string_search_algorithm</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// https://www.cs.utexas.edu/~moore/publications/fstrpos.pdf (note: this aged</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// document uses 1-based indexing)</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>type stringFinder struct {
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	<span class="comment">// pattern is the string that we are searching for in the text.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	pattern string
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	<span class="comment">// badCharSkip[b] contains the distance between the last byte of pattern</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// and the rightmost occurrence of b in pattern. If b is not in pattern,</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	<span class="comment">// badCharSkip[b] is len(pattern).</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// Whenever a mismatch is found with byte b in the text, we can safely</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// shift the matching frame at least badCharSkip[b] until the next time</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// the matching char could be in alignment.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	badCharSkip [256]int
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// goodSuffixSkip[i] defines how far we can shift the matching frame given</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// that the suffix pattern[i+1:] matches, but the byte pattern[i] does</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// not. There are two cases to consider:</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// 1. The matched suffix occurs elsewhere in pattern (with a different</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// byte preceding it that we might possibly match). In this case, we can</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// shift the matching frame to align with the next suffix chunk. For</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// example, the pattern &#34;mississi&#34; has the suffix &#34;issi&#34; next occurring</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// (in right-to-left order) at index 1, so goodSuffixSkip[3] ==</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// shift+len(suffix) == 3+4 == 7.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// 2. If the matched suffix does not occur elsewhere in pattern, then the</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// matching frame may share part of its prefix with the end of the</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// matching suffix. In this case, goodSuffixSkip[i] will contain how far</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// to shift the frame to align this portion of the prefix to the</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// suffix. For example, in the pattern &#34;abcxxxabc&#34;, when the first</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// mismatch from the back is found to be in position 3, the matching</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// suffix &#34;xxabc&#34; is not found elsewhere in the pattern. However, its</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// rightmost &#34;abc&#34; (at position 6) is a prefix of the whole pattern, so</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// goodSuffixSkip[3] == shift+len(suffix) == 6+5 == 11.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	goodSuffixSkip []int
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func makeStringFinder(pattern string) *stringFinder {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	f := &amp;stringFinder{
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		pattern:        pattern,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		goodSuffixSkip: make([]int, len(pattern)),
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// last is the index of the last character in the pattern.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	last := len(pattern) - 1
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Build bad character table.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Bytes not in the pattern can skip one pattern&#39;s length.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	for i := range f.badCharSkip {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		f.badCharSkip[i] = len(pattern)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// The loop condition is &lt; instead of &lt;= so that the last byte does not</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// have a zero distance to itself. Finding this byte out of place implies</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// that it is not in the last position.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	for i := 0; i &lt; last; i++ {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		f.badCharSkip[pattern[i]] = last - i
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// Build good suffix table.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	<span class="comment">// First pass: set each value to the next index which starts a prefix of</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	<span class="comment">// pattern.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	lastPrefix := last
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	for i := last; i &gt;= 0; i-- {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		if HasPrefix(pattern, pattern[i+1:]) {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			lastPrefix = i + 1
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// lastPrefix is the shift, and (last-i) is len(suffix).</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		f.goodSuffixSkip[i] = lastPrefix + last - i
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// Second pass: find repeats of pattern&#39;s suffix starting from the front.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	for i := 0; i &lt; last; i++ {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			<span class="comment">// (last-i) is the shift, and lenSuffix is len(suffix).</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return f
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func longestCommonSuffix(a, b string) (i int) {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	for ; i &lt; len(a) &amp;&amp; i &lt; len(b); i++ {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if a[len(a)-1-i] != b[len(b)-1-i] {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			break
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	return
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// next returns the index in text of the first occurrence of the pattern. If</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// the pattern is not found, it returns -1.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>func (f *stringFinder) next(text string) int {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	i := len(f.pattern) - 1
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	for i &lt; len(text) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// Compare backwards from the end until the first unmatching character.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		j := len(f.pattern) - 1
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		for j &gt;= 0 &amp;&amp; text[i] == f.pattern[j] {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			i--
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			j--
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		if j &lt; 0 {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			return i + 1 <span class="comment">// match</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	return -1
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
</pre><p><a href="search.go?m=text">View as plain text</a></p>

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
