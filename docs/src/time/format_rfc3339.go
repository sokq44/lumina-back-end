<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/time/format_rfc3339.go - Go Documentation Server</title>

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
<a href="format_rfc3339.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/time">time</a>/<span class="text-muted">format_rfc3339.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/time">time</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package time
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;errors&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// RFC 3339 is the most commonly used format.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// It is implicitly used by the Time.(Marshal|Unmarshal)(Text|JSON) methods.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// Also, according to analysis on https://go.dev/issue/52746,</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// RFC 3339 accounts for 57% of all explicitly specified time formats,</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// with the second most popular format only being used 8% of the time.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The overwhelming use of RFC 3339 compared to all other formats justifies</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// the addition of logic to optimize formatting and parsing.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>func (t Time) appendFormatRFC3339(b []byte, nanos bool) []byte {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	_, offset, abs := t.locabs()
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// Format date.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	year, month, day, _ := absDate(abs, true)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	b = appendInt(b, year, 4)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	b = append(b, &#39;-&#39;)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	b = appendInt(b, int(month), 2)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	b = append(b, &#39;-&#39;)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	b = appendInt(b, day, 2)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	b = append(b, &#39;T&#39;)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// Format time.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	hour, min, sec := absClock(abs)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	b = appendInt(b, hour, 2)
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	b = append(b, &#39;:&#39;)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	b = appendInt(b, min, 2)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	b = append(b, &#39;:&#39;)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	b = appendInt(b, sec, 2)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if nanos {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		std := stdFracSecond(stdFracSecond9, 9, &#39;.&#39;)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		b = appendNano(b, t.Nanosecond(), std)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if offset == 0 {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		return append(b, &#39;Z&#39;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// Format zone.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	zone := offset / 60 <span class="comment">// convert to minutes</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if zone &lt; 0 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		b = append(b, &#39;-&#39;)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		zone = -zone
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	} else {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		b = append(b, &#39;+&#39;)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	b = appendInt(b, zone/60, 2)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	b = append(b, &#39;:&#39;)
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	b = appendInt(b, zone%60, 2)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return b
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func (t Time) appendStrictRFC3339(b []byte) ([]byte, error) {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	n0 := len(b)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	b = t.appendFormatRFC3339(b, true)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// Not all valid Go timestamps can be serialized as valid RFC 3339.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">// Explicitly check for these edge cases.</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// See https://go.dev/issue/4556 and https://go.dev/issue/54580.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	num2 := func(b []byte) byte { return 10*(b[0]-&#39;0&#39;) + (b[1] - &#39;0&#39;) }
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	switch {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	case b[n0+len(&#34;9999&#34;)] != &#39;-&#39;: <span class="comment">// year must be exactly 4 digits wide</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return b, errors.New(&#34;year outside of range [0,9999]&#34;)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	case b[len(b)-1] != &#39;Z&#39;:
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		c := b[len(b)-len(&#34;Z07:00&#34;)]
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		if (&#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;) || num2(b[len(b)-len(&#34;07:00&#34;):]) &gt;= 24 {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			return b, errors.New(&#34;timezone hour outside of range [0,23]&#34;)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return b, nil
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>func parseRFC3339[bytes []byte | string](s bytes, local *Location) (Time, bool) {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// parseUint parses s as an unsigned decimal integer and</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// verifies that it is within some range.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// If it is invalid or out-of-range,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// it sets ok to false and returns the min value.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	ok := true
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	parseUint := func(s bytes, min, max int) (x int) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		for _, c := range []byte(s) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			if c &lt; &#39;0&#39; || &#39;9&#39; &lt; c {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				ok = false
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>				return min
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			x = x*10 + int(c) - &#39;0&#39;
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if x &lt; min || max &lt; x {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			ok = false
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			return min
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		return x
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Parse the date and time.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if len(s) &lt; len(&#34;2006-01-02T15:04:05&#34;) {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		return Time{}, false
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	year := parseUint(s[0:4], 0, 9999)                       <span class="comment">// e.g., 2006</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	month := parseUint(s[5:7], 1, 12)                        <span class="comment">// e.g., 01</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	day := parseUint(s[8:10], 1, daysIn(Month(month), year)) <span class="comment">// e.g., 02</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	hour := parseUint(s[11:13], 0, 23)                       <span class="comment">// e.g., 15</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	min := parseUint(s[14:16], 0, 59)                        <span class="comment">// e.g., 04</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	sec := parseUint(s[17:19], 0, 59)                        <span class="comment">// e.g., 05</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	if !ok || !(s[4] == &#39;-&#39; &amp;&amp; s[7] == &#39;-&#39; &amp;&amp; s[10] == &#39;T&#39; &amp;&amp; s[13] == &#39;:&#39; &amp;&amp; s[16] == &#39;:&#39;) {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		return Time{}, false
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	s = s[19:]
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// Parse the fractional second.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	var nsec int
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	if len(s) &gt;= 2 &amp;&amp; s[0] == &#39;.&#39; &amp;&amp; isDigit(s, 1) {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		n := 2
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		for ; n &lt; len(s) &amp;&amp; isDigit(s, n); n++ {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		nsec, _, _ = parseNanoseconds(s, n)
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		s = s[n:]
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Parse the time zone.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	t := Date(year, Month(month), day, hour, min, sec, nsec, UTC)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	if len(s) != 1 || s[0] != &#39;Z&#39; {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		if len(s) != len(&#34;-07:00&#34;) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			return Time{}, false
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		hr := parseUint(s[1:3], 0, 23) <span class="comment">// e.g., 07</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		mm := parseUint(s[4:6], 0, 59) <span class="comment">// e.g., 00</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		if !ok || !((s[0] == &#39;-&#39; || s[0] == &#39;+&#39;) &amp;&amp; s[3] == &#39;:&#39;) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			return Time{}, false
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		zoneOffset := (hr*60 + mm) * 60
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		if s[0] == &#39;-&#39; {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			zoneOffset *= -1
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		t.addSec(-int64(zoneOffset))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		<span class="comment">// Use local zone with the given offset if possible.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if _, offset, _, _, _ := local.lookup(t.unixSec()); offset == zoneOffset {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			t.setLoc(local)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		} else {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			t.setLoc(FixedZone(&#34;&#34;, zoneOffset))
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	return t, true
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func parseStrictRFC3339(b []byte) (Time, error) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	t, ok := parseRFC3339(b, Local)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if !ok {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		t, err := Parse(RFC3339, string(b))
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		if err != nil {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			return Time{}, err
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// The parse template syntax cannot correctly validate RFC 3339.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// Explicitly check for cases that Parse is unable to validate for.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// See https://go.dev/issue/54580.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		num2 := func(b []byte) byte { return 10*(b[0]-&#39;0&#39;) + (b[1] - &#39;0&#39;) }
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		switch {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// TODO(https://go.dev/issue/54580): Strict parsing is disabled for now.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		<span class="comment">// Enable this again with a GODEBUG opt-out.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		case true:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			return t, nil
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		case b[len(&#34;2006-01-02T&#34;)+1] == &#39;:&#39;: <span class="comment">// hour must be two digits</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return Time{}, &amp;ParseError{RFC3339, string(b), &#34;15&#34;, string(b[len(&#34;2006-01-02T&#34;):][:1]), &#34;&#34;}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		case b[len(&#34;2006-01-02T15:04:05&#34;)] == &#39;,&#39;: <span class="comment">// sub-second separator must be a period</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			return Time{}, &amp;ParseError{RFC3339, string(b), &#34;.&#34;, &#34;,&#34;, &#34;&#34;}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		case b[len(b)-1] != &#39;Z&#39;:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			switch {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			case num2(b[len(b)-len(&#34;07:00&#34;):]) &gt;= 24: <span class="comment">// timezone hour must be in range</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				return Time{}, &amp;ParseError{RFC3339, string(b), &#34;Z07:00&#34;, string(b[len(b)-len(&#34;Z07:00&#34;):]), &#34;: timezone hour out of range&#34;}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			case num2(b[len(b)-len(&#34;00&#34;):]) &gt;= 60: <span class="comment">// timezone minute must be in range</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>				return Time{}, &amp;ParseError{RFC3339, string(b), &#34;Z07:00&#34;, string(b[len(b)-len(&#34;Z07:00&#34;):]), &#34;: timezone minute out of range&#34;}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		default: <span class="comment">// unknown error; should not occur</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			return Time{}, &amp;ParseError{RFC3339, string(b), RFC3339, string(b), &#34;&#34;}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	return t, nil
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
</pre><p><a href="format_rfc3339.go?m=text">View as plain text</a></p>

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
