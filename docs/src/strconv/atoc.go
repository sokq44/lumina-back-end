<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/strconv/atoc.go - Go Documentation Server</title>

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
<a href="atoc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/strconv">strconv</a>/<span class="text-muted">atoc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/strconv">strconv</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package strconv
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>const fnParseComplex = &#34;ParseComplex&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// convErr splits an error returned by parseFloatPrefix</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// into a syntax or range error for ParseComplex.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>func convErr(err error, s string) (syntax, range_ error) {
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	if x, ok := err.(*NumError); ok {
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>		x.Func = fnParseComplex
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>		x.Num = cloneString(s)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>		if x.Err == ErrRange {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>			return nil, x
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	}
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	return err, nil
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>}
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// ParseComplex converts the string s to a complex number</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// with the precision specified by bitSize: 64 for complex64, or 128 for complex128.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// When bitSize=64, the result still has type complex128, but it will be</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// convertible to complex64 without changing its value.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// The number represented by s must be of the form N, Ni, or N±Ni, where N stands</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// for a floating-point number as recognized by ParseFloat, and i is the imaginary</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// component. If the second N is unsigned, a + sign is required between the two components</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// as indicated by the ±. If the second N is NaN, only a + sign is accepted.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// The form may be parenthesized and cannot contain any spaces.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// The resulting complex number consists of the two components converted by ParseFloat.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// The errors that ParseComplex returns have concrete type *NumError</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// and include err.Num = s.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// If s is not syntactically well-formed, ParseComplex returns err.Err = ErrSyntax.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// If s is syntactically well-formed but either component is more than 1/2 ULP</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// away from the largest floating point number of the given component&#39;s size,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// ParseComplex returns err.Err = ErrRange and c = ±Inf for the respective component.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func ParseComplex(s string, bitSize int) (complex128, error) {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	size := 64
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if bitSize == 64 {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		size = 32 <span class="comment">// complex64 uses float32 parts</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	orig := s
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	<span class="comment">// Remove parentheses, if any.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if len(s) &gt;= 2 &amp;&amp; s[0] == &#39;(&#39; &amp;&amp; s[len(s)-1] == &#39;)&#39; {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		s = s[1 : len(s)-1]
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	var pending error <span class="comment">// pending range error, or nil</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Read real part (possibly imaginary part if followed by &#39;i&#39;).</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	re, n, err := parseFloatPrefix(s, size)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if err != nil {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		err, pending = convErr(err, orig)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if err != nil {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			return 0, err
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	s = s[n:]
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// If we have nothing left, we&#39;re done.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if len(s) == 0 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return complex(re, 0), pending
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// Otherwise, look at the next character.</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	switch s[0] {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	case &#39;+&#39;:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		<span class="comment">// Consume the &#39;+&#39; to avoid an error if we have &#34;+NaNi&#34;, but</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// do this only if we don&#39;t have a &#34;++&#34; (don&#39;t hide that error).</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		if len(s) &gt; 1 &amp;&amp; s[1] != &#39;+&#39; {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			s = s[1:]
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	case &#39;-&#39;:
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		<span class="comment">// ok</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	case &#39;i&#39;:
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		<span class="comment">// If &#39;i&#39; is the last character, we only have an imaginary part.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if len(s) == 1 {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			return complex(0, re), pending
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		fallthrough
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	default:
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		return 0, syntaxError(fnParseComplex, orig)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// Read imaginary part.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	im, n, err := parseFloatPrefix(s, size)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if err != nil {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		err, pending = convErr(err, orig)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if err != nil {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return 0, err
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	s = s[n:]
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if s != &#34;i&#34; {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		return 0, syntaxError(fnParseComplex, orig)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	return complex(re, im), pending
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
</pre><p><a href="atoc.go?m=text">View as plain text</a></p>

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
