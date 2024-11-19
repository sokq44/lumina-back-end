<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/time/format.go - Go Documentation Server</title>

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
<a href="format.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/time">time</a>/<span class="text-muted">format.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/time">time</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package time
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;errors&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// These are predefined layouts for use in Time.Format and time.Parse.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// The reference time used in these layouts is the specific time stamp:</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//	01/02 03:04:05PM &#39;06 -0700</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// (January 2, 15:04:05, 2006, in time zone seven hours west of GMT).</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// That value is recorded as the constant named Layout, listed below. As a Unix</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// time, this is 1136239445. Since MST is GMT-0700, the reference would be</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// printed by the Unix date command as:</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//	Mon Jan 2 15:04:05 MST 2006</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// It is a regrettable historic error that the date uses the American convention</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// of putting the numerical month before the day.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// The example for Time.Format demonstrates the working of the layout string</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// in detail and is a good reference.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Note that the RFC822, RFC850, and RFC1123 formats should be applied</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// only to local times. Applying them to UTC times will use &#34;UTC&#34; as the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// time zone abbreviation, while strictly speaking those RFCs require the</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// use of &#34;GMT&#34; in that case.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// In general RFC1123Z should be used instead of RFC1123 for servers</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// that insist on that format, and RFC3339 should be preferred for new protocols.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// RFC3339, RFC822, RFC822Z, RFC1123, and RFC1123Z are useful for formatting;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// when used with time.Parse they do not accept all the time formats</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// permitted by the RFCs and they do accept time formats not formally defined.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// The RFC3339Nano format removes trailing zeros from the seconds field</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// and thus may not sort correctly once formatted.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// Most programs can use one of the defined constants as the layout passed to</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// Format or Parse. The rest of this comment can be ignored unless you are</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// creating a custom layout string.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// To define your own format, write down what the reference time would look like</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// formatted your way; see the values of constants like ANSIC, StampMicro or</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Kitchen for examples. The model is to demonstrate what the reference time</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// looks like so that the Format and Parse methods can apply the same</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// transformation to a general time value.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// Here is a summary of the components of a layout string. Each element shows by</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// example the formatting of an element of the reference time. Only these values</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// are recognized. Text in the layout string that is not recognized as part of</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// the reference time is echoed verbatim during Format and expected to appear</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// verbatim in the input to Parse.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//	Year: &#34;2006&#34; &#34;06&#34;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//	Month: &#34;Jan&#34; &#34;January&#34; &#34;01&#34; &#34;1&#34;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//	Day of the week: &#34;Mon&#34; &#34;Monday&#34;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//	Day of the month: &#34;2&#34; &#34;_2&#34; &#34;02&#34;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">//	Day of the year: &#34;__2&#34; &#34;002&#34;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//	Hour: &#34;15&#34; &#34;3&#34; &#34;03&#34; (PM or AM)</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">//	Minute: &#34;4&#34; &#34;04&#34;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">//	Second: &#34;5&#34; &#34;05&#34;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//	AM/PM mark: &#34;PM&#34;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// Numeric time zone offsets format as follows:</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">//	&#34;-0700&#34;     ±hhmm</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//	&#34;-07:00&#34;    ±hh:mm</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//	&#34;-07&#34;       ±hh</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//	&#34;-070000&#34;   ±hhmmss</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//	&#34;-07:00:00&#34; ±hh:mm:ss</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">// Replacing the sign in the format with a Z triggers</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// the ISO 8601 behavior of printing Z instead of an</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">// offset for the UTC zone. Thus:</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">//	&#34;Z0700&#34;      Z or ±hhmm</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">//	&#34;Z07:00&#34;     Z or ±hh:mm</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//	&#34;Z07&#34;        Z or ±hh</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//	&#34;Z070000&#34;    Z or ±hhmmss</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//	&#34;Z07:00:00&#34;  Z or ±hh:mm:ss</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// Within the format string, the underscores in &#34;_2&#34; and &#34;__2&#34; represent spaces</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// that may be replaced by digits if the following number has multiple digits,</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// for compatibility with fixed-width Unix time formats. A leading zero represents</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// a zero-padded value.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// The formats __2 and 002 are space-padded and zero-padded</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// three-character day of year; there is no unpadded day of year format.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// A comma or decimal point followed by one or more zeros represents</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// a fractional second, printed to the given number of decimal places.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// A comma or decimal point followed by one or more nines represents</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// a fractional second, printed to the given number of decimal places, with</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// trailing zeros removed.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// For example &#34;15:04:05,000&#34; or &#34;15:04:05.000&#34; formats or parses with</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// millisecond precision.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Some valid layouts are invalid time values for time.Parse, due to formats</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// such as _ for space padding and Z for zone information.</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>const (
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	Layout      = &#34;01/02 03:04:05PM &#39;06 -0700&#34; <span class="comment">// The reference time, in numerical order.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	ANSIC       = &#34;Mon Jan _2 15:04:05 2006&#34;
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	UnixDate    = &#34;Mon Jan _2 15:04:05 MST 2006&#34;
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	RubyDate    = &#34;Mon Jan 02 15:04:05 -0700 2006&#34;
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	RFC822      = &#34;02 Jan 06 15:04 MST&#34;
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	RFC822Z     = &#34;02 Jan 06 15:04 -0700&#34; <span class="comment">// RFC822 with numeric zone</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	RFC850      = &#34;Monday, 02-Jan-06 15:04:05 MST&#34;
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	RFC1123     = &#34;Mon, 02 Jan 2006 15:04:05 MST&#34;
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	RFC1123Z    = &#34;Mon, 02 Jan 2006 15:04:05 -0700&#34; <span class="comment">// RFC1123 with numeric zone</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	RFC3339     = &#34;2006-01-02T15:04:05Z07:00&#34;
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	RFC3339Nano = &#34;2006-01-02T15:04:05.999999999Z07:00&#34;
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	Kitchen     = &#34;3:04PM&#34;
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// Handy time stamps.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	Stamp      = &#34;Jan _2 15:04:05&#34;
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	StampMilli = &#34;Jan _2 15:04:05.000&#34;
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	StampMicro = &#34;Jan _2 15:04:05.000000&#34;
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	StampNano  = &#34;Jan _2 15:04:05.000000000&#34;
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	DateTime   = &#34;2006-01-02 15:04:05&#34;
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	DateOnly   = &#34;2006-01-02&#34;
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	TimeOnly   = &#34;15:04:05&#34;
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>const (
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	_                        = iota
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	stdLongMonth             = iota + stdNeedDate  <span class="comment">// &#34;January&#34;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	stdMonth                                       <span class="comment">// &#34;Jan&#34;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	stdNumMonth                                    <span class="comment">// &#34;1&#34;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	stdZeroMonth                                   <span class="comment">// &#34;01&#34;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	stdLongWeekDay                                 <span class="comment">// &#34;Monday&#34;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	stdWeekDay                                     <span class="comment">// &#34;Mon&#34;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	stdDay                                         <span class="comment">// &#34;2&#34;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	stdUnderDay                                    <span class="comment">// &#34;_2&#34;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	stdZeroDay                                     <span class="comment">// &#34;02&#34;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	stdUnderYearDay                                <span class="comment">// &#34;__2&#34;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	stdZeroYearDay                                 <span class="comment">// &#34;002&#34;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	stdHour                  = iota + stdNeedClock <span class="comment">// &#34;15&#34;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	stdHour12                                      <span class="comment">// &#34;3&#34;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	stdZeroHour12                                  <span class="comment">// &#34;03&#34;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	stdMinute                                      <span class="comment">// &#34;4&#34;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	stdZeroMinute                                  <span class="comment">// &#34;04&#34;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	stdSecond                                      <span class="comment">// &#34;5&#34;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	stdZeroSecond                                  <span class="comment">// &#34;05&#34;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	stdLongYear              = iota + stdNeedDate  <span class="comment">// &#34;2006&#34;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	stdYear                                        <span class="comment">// &#34;06&#34;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	stdPM                    = iota + stdNeedClock <span class="comment">// &#34;PM&#34;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	stdpm                                          <span class="comment">// &#34;pm&#34;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	stdTZ                    = iota                <span class="comment">// &#34;MST&#34;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	stdISO8601TZ                                   <span class="comment">// &#34;Z0700&#34;  // prints Z for UTC</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	stdISO8601SecondsTZ                            <span class="comment">// &#34;Z070000&#34;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	stdISO8601ShortTZ                              <span class="comment">// &#34;Z07&#34;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	stdISO8601ColonTZ                              <span class="comment">// &#34;Z07:00&#34; // prints Z for UTC</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	stdISO8601ColonSecondsTZ                       <span class="comment">// &#34;Z07:00:00&#34;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	stdNumTZ                                       <span class="comment">// &#34;-0700&#34;  // always numeric</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	stdNumSecondsTz                                <span class="comment">// &#34;-070000&#34;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	stdNumShortTZ                                  <span class="comment">// &#34;-07&#34;    // always numeric</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	stdNumColonTZ                                  <span class="comment">// &#34;-07:00&#34; // always numeric</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	stdNumColonSecondsTZ                           <span class="comment">// &#34;-07:00:00&#34;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	stdFracSecond0                                 <span class="comment">// &#34;.0&#34;, &#34;.00&#34;, ... , trailing zeros included</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	stdFracSecond9                                 <span class="comment">// &#34;.9&#34;, &#34;.99&#34;, ..., trailing zeros omitted</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	stdNeedDate       = 1 &lt;&lt; 8             <span class="comment">// need month, day, year</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	stdNeedClock      = 2 &lt;&lt; 8             <span class="comment">// need hour, minute, second</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	stdArgShift       = 16                 <span class="comment">// extra argument in high bits, above low stdArgShift</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	stdSeparatorShift = 28                 <span class="comment">// extra argument in high 4 bits for fractional second separators</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	stdMask           = 1&lt;&lt;stdArgShift - 1 <span class="comment">// mask out argument</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span><span class="comment">// std0x records the std values for &#34;01&#34;, &#34;02&#34;, ..., &#34;06&#34;.</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>var std0x = [...]int{stdZeroMonth, stdZeroDay, stdZeroHour12, stdZeroMinute, stdZeroSecond, stdYear}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// startsWithLowerCase reports whether the string has a lower-case letter at the beginning.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// Its purpose is to prevent matching strings like &#34;Month&#34; when looking for &#34;Mon&#34;.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func startsWithLowerCase(str string) bool {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if len(str) == 0 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		return false
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	c := str[0]
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	return &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39;
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span><span class="comment">// nextStdChunk finds the first occurrence of a std string in</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span><span class="comment">// layout and returns the text before, the std string, and the text after.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>func nextStdChunk(layout string) (prefix string, std int, suffix string) {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	for i := 0; i &lt; len(layout); i++ {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		switch c := int(layout[i]); c {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		case &#39;J&#39;: <span class="comment">// January, Jan</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 &amp;&amp; layout[i:i+3] == &#34;Jan&#34; {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				if len(layout) &gt;= i+7 &amp;&amp; layout[i:i+7] == &#34;January&#34; {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>					return layout[0:i], stdLongMonth, layout[i+7:]
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>				if !startsWithLowerCase(layout[i+3:]) {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>					return layout[0:i], stdMonth, layout[i+3:]
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>				}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		case &#39;M&#39;: <span class="comment">// Monday, Mon, MST</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				if layout[i:i+3] == &#34;Mon&#34; {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>					if len(layout) &gt;= i+6 &amp;&amp; layout[i:i+6] == &#34;Monday&#34; {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>						return layout[0:i], stdLongWeekDay, layout[i+6:]
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>					}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>					if !startsWithLowerCase(layout[i+3:]) {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>						return layout[0:i], stdWeekDay, layout[i+3:]
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>					}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				if layout[i:i+3] == &#34;MST&#34; {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>					return layout[0:i], stdTZ, layout[i+3:]
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		case &#39;0&#39;: <span class="comment">// 01, 02, 03, 04, 05, 06, 002</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			if len(layout) &gt;= i+2 &amp;&amp; &#39;1&#39; &lt;= layout[i+1] &amp;&amp; layout[i+1] &lt;= &#39;6&#39; {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				return layout[0:i], std0x[layout[i+1]-&#39;1&#39;], layout[i+2:]
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 &amp;&amp; layout[i+1] == &#39;0&#39; &amp;&amp; layout[i+2] == &#39;2&#39; {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>				return layout[0:i], stdZeroYearDay, layout[i+3:]
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		case &#39;1&#39;: <span class="comment">// 15, 1</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			if len(layout) &gt;= i+2 &amp;&amp; layout[i+1] == &#39;5&#39; {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>				return layout[0:i], stdHour, layout[i+2:]
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			return layout[0:i], stdNumMonth, layout[i+1:]
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		case &#39;2&#39;: <span class="comment">// 2006, 2</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			if len(layout) &gt;= i+4 &amp;&amp; layout[i:i+4] == &#34;2006&#34; {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				return layout[0:i], stdLongYear, layout[i+4:]
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			return layout[0:i], stdDay, layout[i+1:]
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		case &#39;_&#39;: <span class="comment">// _2, _2006, __2</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			if len(layout) &gt;= i+2 &amp;&amp; layout[i+1] == &#39;2&#39; {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				<span class="comment">//_2006 is really a literal _, followed by stdLongYear</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				if len(layout) &gt;= i+5 &amp;&amp; layout[i+1:i+5] == &#34;2006&#34; {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>					return layout[0 : i+1], stdLongYear, layout[i+5:]
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>				}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				return layout[0:i], stdUnderDay, layout[i+2:]
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 &amp;&amp; layout[i+1] == &#39;_&#39; &amp;&amp; layout[i+2] == &#39;2&#39; {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>				return layout[0:i], stdUnderYearDay, layout[i+3:]
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		case &#39;3&#39;:
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			return layout[0:i], stdHour12, layout[i+1:]
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		case &#39;4&#39;:
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			return layout[0:i], stdMinute, layout[i+1:]
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		case &#39;5&#39;:
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>			return layout[0:i], stdSecond, layout[i+1:]
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		case &#39;P&#39;: <span class="comment">// PM</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			if len(layout) &gt;= i+2 &amp;&amp; layout[i+1] == &#39;M&#39; {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				return layout[0:i], stdPM, layout[i+2:]
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		case &#39;p&#39;: <span class="comment">// pm</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			if len(layout) &gt;= i+2 &amp;&amp; layout[i+1] == &#39;m&#39; {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>				return layout[0:i], stdpm, layout[i+2:]
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		case &#39;-&#39;: <span class="comment">// -070000, -07:00:00, -0700, -07:00, -07</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			if len(layout) &gt;= i+7 &amp;&amp; layout[i:i+7] == &#34;-070000&#34; {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				return layout[0:i], stdNumSecondsTz, layout[i+7:]
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			if len(layout) &gt;= i+9 &amp;&amp; layout[i:i+9] == &#34;-07:00:00&#34; {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>				return layout[0:i], stdNumColonSecondsTZ, layout[i+9:]
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			if len(layout) &gt;= i+5 &amp;&amp; layout[i:i+5] == &#34;-0700&#34; {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>				return layout[0:i], stdNumTZ, layout[i+5:]
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			if len(layout) &gt;= i+6 &amp;&amp; layout[i:i+6] == &#34;-07:00&#34; {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>				return layout[0:i], stdNumColonTZ, layout[i+6:]
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 &amp;&amp; layout[i:i+3] == &#34;-07&#34; {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>				return layout[0:i], stdNumShortTZ, layout[i+3:]
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		case &#39;Z&#39;: <span class="comment">// Z070000, Z07:00:00, Z0700, Z07:00,</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			if len(layout) &gt;= i+7 &amp;&amp; layout[i:i+7] == &#34;Z070000&#34; {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>				return layout[0:i], stdISO8601SecondsTZ, layout[i+7:]
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>			if len(layout) &gt;= i+9 &amp;&amp; layout[i:i+9] == &#34;Z07:00:00&#34; {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				return layout[0:i], stdISO8601ColonSecondsTZ, layout[i+9:]
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>			}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			if len(layout) &gt;= i+5 &amp;&amp; layout[i:i+5] == &#34;Z0700&#34; {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>				return layout[0:i], stdISO8601TZ, layout[i+5:]
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			if len(layout) &gt;= i+6 &amp;&amp; layout[i:i+6] == &#34;Z07:00&#34; {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				return layout[0:i], stdISO8601ColonTZ, layout[i+6:]
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>			if len(layout) &gt;= i+3 &amp;&amp; layout[i:i+3] == &#34;Z07&#34; {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				return layout[0:i], stdISO8601ShortTZ, layout[i+3:]
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		case &#39;.&#39;, &#39;,&#39;: <span class="comment">// ,000, or .000, or ,999, or .999 - repeated digits for fractional seconds.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			if i+1 &lt; len(layout) &amp;&amp; (layout[i+1] == &#39;0&#39; || layout[i+1] == &#39;9&#39;) {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>				ch := layout[i+1]
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				j := i + 1
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				for j &lt; len(layout) &amp;&amp; layout[j] == ch {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>					j++
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>				}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				<span class="comment">// String of digits must end here - only fractional second is all digits.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>				if !isDigit(layout, j) {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>					code := stdFracSecond0
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>					if layout[i+1] == &#39;9&#39; {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>						code = stdFracSecond9
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>					}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					std := stdFracSecond(code, j-(i+1), c)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>					return layout[0:i], std, layout[j:]
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	return layout, 0, &#34;&#34;
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>var longDayNames = []string{
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	&#34;Sunday&#34;,
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	&#34;Monday&#34;,
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	&#34;Tuesday&#34;,
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	&#34;Wednesday&#34;,
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	&#34;Thursday&#34;,
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	&#34;Friday&#34;,
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	&#34;Saturday&#34;,
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>var shortDayNames = []string{
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	&#34;Sun&#34;,
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	&#34;Mon&#34;,
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	&#34;Tue&#34;,
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	&#34;Wed&#34;,
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	&#34;Thu&#34;,
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	&#34;Fri&#34;,
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	&#34;Sat&#34;,
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>var shortMonthNames = []string{
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	&#34;Jan&#34;,
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	&#34;Feb&#34;,
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	&#34;Mar&#34;,
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	&#34;Apr&#34;,
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	&#34;May&#34;,
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	&#34;Jun&#34;,
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	&#34;Jul&#34;,
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	&#34;Aug&#34;,
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	&#34;Sep&#34;,
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	&#34;Oct&#34;,
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	&#34;Nov&#34;,
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	&#34;Dec&#34;,
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>var longMonthNames = []string{
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	&#34;January&#34;,
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	&#34;February&#34;,
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	&#34;March&#34;,
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	&#34;April&#34;,
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	&#34;May&#34;,
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	&#34;June&#34;,
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	&#34;July&#34;,
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	&#34;August&#34;,
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	&#34;September&#34;,
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	&#34;October&#34;,
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	&#34;November&#34;,
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	&#34;December&#34;,
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span><span class="comment">// match reports whether s1 and s2 match ignoring case.</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// It is assumed s1 and s2 are the same length.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>func match(s1, s2 string) bool {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s1); i++ {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		c1 := s1[i]
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		c2 := s2[i]
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		if c1 != c2 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			<span class="comment">// Switch to lower-case; &#39;a&#39;-&#39;A&#39; is known to be a single bit.</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			c1 |= &#39;a&#39; - &#39;A&#39;
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>			c2 |= &#39;a&#39; - &#39;A&#39;
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			if c1 != c2 || c1 &lt; &#39;a&#39; || c1 &gt; &#39;z&#39; {
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>				return false
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>			}
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	return true
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>func lookup(tab []string, val string) (int, string, error) {
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	for i, v := range tab {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		if len(val) &gt;= len(v) &amp;&amp; match(val[0:len(v)], v) {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>			return i, val[len(v):], nil
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	}
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	return -1, val, errBad
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// appendInt appends the decimal form of x to b and returns the result.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// If the decimal form (excluding sign) is shorter than width, the result is padded with leading 0&#39;s.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">// Duplicates functionality in strconv, but avoids dependency.</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>func appendInt(b []byte, x int, width int) []byte {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	u := uint(x)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		b = append(b, &#39;-&#39;)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		u = uint(-x)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	<span class="comment">// 2-digit and 4-digit fields are the most common in time formats.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	utod := func(u uint) byte { return &#39;0&#39; + byte(u) }
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	switch {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	case width == 2 &amp;&amp; u &lt; 1e2:
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		return append(b, utod(u/1e1), utod(u%1e1))
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	case width == 4 &amp;&amp; u &lt; 1e4:
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		return append(b, utod(u/1e3), utod(u/1e2%1e1), utod(u/1e1%1e1), utod(u%1e1))
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	<span class="comment">// Compute the number of decimal digits.</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	var n int
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	if u == 0 {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		n = 1
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	}
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	for u2 := u; u2 &gt; 0; u2 /= 10 {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		n++
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// Add 0-padding.</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	for pad := width - n; pad &gt; 0; pad-- {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		b = append(b, &#39;0&#39;)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	<span class="comment">// Ensure capacity.</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	if len(b)+n &lt;= cap(b) {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		b = b[:len(b)+n]
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	} else {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		b = append(b, make([]byte, n)...)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// Assemble decimal in reverse order.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	i := len(b) - 1
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	for u &gt;= 10 &amp;&amp; i &gt; 0 {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		q := u / 10
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		b[i] = utod(u - q*10)
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		u = q
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		i--
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	}
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	b[i] = utod(u)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	return b
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span><span class="comment">// Never printed, just needs to be non-nil for return by atoi.</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>var errAtoi = errors.New(&#34;time: invalid number&#34;)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span><span class="comment">// Duplicates functionality in strconv, but avoids dependency.</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>func atoi[bytes []byte | string](s bytes) (x int, err error) {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	neg := false
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	if len(s) &gt; 0 &amp;&amp; (s[0] == &#39;-&#39; || s[0] == &#39;+&#39;) {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		neg = s[0] == &#39;-&#39;
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		s = s[1:]
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	q, rem, err := leadingInt(s)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	x = int(q)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	if err != nil || len(rem) &gt; 0 {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		return 0, errAtoi
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	if neg {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		x = -x
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	return x, nil
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// The &#34;std&#34; value passed to appendNano contains two packed fields: the number of</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// digits after the decimal and the separator character (period or comma).</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// These functions pack and unpack that variable.</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>func stdFracSecond(code, n, c int) int {
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	<span class="comment">// Use 0xfff to make the failure case even more absurd.</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	if c == &#39;.&#39; {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		return code | ((n &amp; 0xfff) &lt;&lt; stdArgShift)
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	return code | ((n &amp; 0xfff) &lt;&lt; stdArgShift) | 1&lt;&lt;stdSeparatorShift
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>}
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>func digitsLen(std int) int {
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	return (std &gt;&gt; stdArgShift) &amp; 0xfff
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>func separator(std int) byte {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	if (std &gt;&gt; stdSeparatorShift) == 0 {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		return &#39;.&#39;
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	return &#39;,&#39;
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span><span class="comment">// appendNano appends a fractional second, as nanoseconds, to b</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span><span class="comment">// and returns the result. The nanosec must be within [0, 999999999].</span>
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>func appendNano(b []byte, nanosec int, std int) []byte {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	trim := std&amp;stdMask == stdFracSecond9
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	n := digitsLen(std)
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	if trim &amp;&amp; (n == 0 || nanosec == 0) {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		return b
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	dot := separator(std)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	b = append(b, dot)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	b = appendInt(b, nanosec, 9)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	if n &lt; 9 {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		b = b[:len(b)-9+n]
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if trim {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		for len(b) &gt; 0 &amp;&amp; b[len(b)-1] == &#39;0&#39; {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			b = b[:len(b)-1]
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>		}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>		if len(b) &gt; 0 &amp;&amp; b[len(b)-1] == dot {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			b = b[:len(b)-1]
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	return b
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// String returns the time formatted using the format string</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span><span class="comment">//	&#34;2006-01-02 15:04:05.999999999 -0700 MST&#34;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">// If the time has a monotonic clock reading, the returned string</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span><span class="comment">// includes a final field &#34;m=±&lt;value&gt;&#34;, where value is the monotonic</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// clock reading formatted as a decimal number of seconds.</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// The returned string is meant for debugging; for a stable serialized</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">// representation, use t.MarshalText, t.MarshalBinary, or t.Format</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">// with an explicit format string.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>func (t Time) String() string {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	s := t.Format(&#34;2006-01-02 15:04:05.999999999 -0700 MST&#34;)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// Format monotonic clock reading as m=±ddd.nnnnnnnnn.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if t.wall&amp;hasMonotonic != 0 {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		m2 := uint64(t.ext)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		sign := byte(&#39;+&#39;)
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if t.ext &lt; 0 {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			sign = &#39;-&#39;
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			m2 = -m2
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>		m1, m2 := m2/1e9, m2%1e9
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		m0, m1 := m1/1e9, m1%1e9
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		buf := make([]byte, 0, 24)
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		buf = append(buf, &#34; m=&#34;...)
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		buf = append(buf, sign)
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		wid := 0
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if m0 != 0 {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			buf = appendInt(buf, int(m0), 0)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			wid = 9
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		buf = appendInt(buf, int(m1), wid)
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		buf = append(buf, &#39;.&#39;)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		buf = appendInt(buf, int(m2), 9)
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		s += string(buf)
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>	}
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	return s
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// GoString implements fmt.GoStringer and formats t to be printed in Go source</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span><span class="comment">// code.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>func (t Time) GoString() string {
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	abs := t.abs()
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	year, month, day, _ := absDate(abs, true)
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	hour, minute, second := absClock(abs)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	buf := make([]byte, 0, len(&#34;time.Date(9999, time.September, 31, 23, 59, 59, 999999999, time.Local)&#34;))
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	buf = append(buf, &#34;time.Date(&#34;...)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	buf = appendInt(buf, year, 0)
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	if January &lt;= month &amp;&amp; month &lt;= December {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		buf = append(buf, &#34;, time.&#34;...)
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		buf = append(buf, longMonthNames[month-1]...)
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	} else {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// It&#39;s difficult to construct a time.Time with a date outside the</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// standard range but we might as well try to handle the case.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		buf = appendInt(buf, int(month), 0)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	buf = appendInt(buf, day, 0)
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	buf = appendInt(buf, hour, 0)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	buf = appendInt(buf, minute, 0)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	buf = appendInt(buf, second, 0)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	buf = appendInt(buf, t.Nanosecond(), 0)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	buf = append(buf, &#34;, &#34;...)
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	switch loc := t.Location(); loc {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	case UTC, nil:
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		buf = append(buf, &#34;time.UTC&#34;...)
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	case Local:
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		buf = append(buf, &#34;time.Local&#34;...)
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	default:
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>		<span class="comment">// there are several options for how we could display this, none of</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>		<span class="comment">// which are great:</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>		<span class="comment">// - use Location(loc.name), which is not technically valid syntax</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>		<span class="comment">// - use LoadLocation(loc.name), which will cause a syntax error when</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		<span class="comment">// embedded and also would require us to escape the string without</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		<span class="comment">// importing fmt or strconv</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		<span class="comment">// - try to use FixedZone, which would also require escaping the name</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		<span class="comment">// and would represent e.g. &#34;America/Los_Angeles&#34; daylight saving time</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		<span class="comment">// shifts inaccurately</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		<span class="comment">// - use the pointer format, which is no worse than you&#39;d get with the</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		<span class="comment">// old fmt.Sprintf(&#34;%#v&#34;, t) format.</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		<span class="comment">// Of these, Location(loc.name) is the least disruptive. This is an edge</span>
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>		<span class="comment">// case we hope not to hit too often.</span>
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		buf = append(buf, `time.Location(`...)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		buf = append(buf, quote(loc.name)...)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		buf = append(buf, &#39;)&#39;)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	buf = append(buf, &#39;)&#39;)
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	return string(buf)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// Format returns a textual representation of the time value formatted according</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">// to the layout defined by the argument. See the documentation for the</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">// constant called Layout to see how to represent the layout format.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L618" class="ln">   618&nbsp;&nbsp;</span><span class="comment">// The executable example for Time.Format demonstrates the working</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">// of the layout string in detail and is a good reference.</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>func (t Time) Format(layout string) string {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	const bufSize = 64
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	var b []byte
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	max := len(layout) + 10
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	if max &lt; bufSize {
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>		var buf [bufSize]byte
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		b = buf[:0]
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>	} else {
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>		b = make([]byte, 0, max)
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	b = t.AppendFormat(b, layout)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>	return string(b)
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// AppendFormat is like Format but appends the textual</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span><span class="comment">// representation to b and returns the extended buffer.</span>
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>func (t Time) AppendFormat(b []byte, layout string) []byte {
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	<span class="comment">// Optimize for RFC3339 as it accounts for over half of all representations.</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	switch layout {
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	case RFC3339:
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		return t.appendFormatRFC3339(b, false)
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	case RFC3339Nano:
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>		return t.appendFormatRFC3339(b, true)
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	default:
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>		return t.appendFormat(b, layout)
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	}
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>}
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>func (t Time) appendFormat(b []byte, layout string) []byte {
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	var (
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		name, offset, abs = t.locabs()
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>		year  int = -1
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		month Month
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		day   int
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		yday  int
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		hour  int = -1
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		min   int
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		sec   int
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// Each iteration generates one std value.</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	for layout != &#34;&#34; {
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>		prefix, std, suffix := nextStdChunk(layout)
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		if prefix != &#34;&#34; {
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			b = append(b, prefix...)
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>		}
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		if std == 0 {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>			break
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		}
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>		layout = suffix
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		<span class="comment">// Compute year, month, day if needed.</span>
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		if year &lt; 0 &amp;&amp; std&amp;stdNeedDate != 0 {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>			year, month, day, yday = absDate(abs, true)
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			yday++
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		<span class="comment">// Compute hour, minute, second if needed.</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		if hour &lt; 0 &amp;&amp; std&amp;stdNeedClock != 0 {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			hour, min, sec = absClock(abs)
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		}
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		switch std &amp; stdMask {
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		case stdYear:
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			y := year
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>			if y &lt; 0 {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>				y = -y
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>			}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			b = appendInt(b, y%100, 2)
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		case stdLongYear:
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>			b = appendInt(b, year, 4)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		case stdMonth:
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			b = append(b, month.String()[:3]...)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		case stdLongMonth:
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			m := month.String()
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>			b = append(b, m...)
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		case stdNumMonth:
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			b = appendInt(b, int(month), 0)
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		case stdZeroMonth:
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>			b = appendInt(b, int(month), 2)
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>		case stdWeekDay:
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			b = append(b, absWeekday(abs).String()[:3]...)
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		case stdLongWeekDay:
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			s := absWeekday(abs).String()
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			b = append(b, s...)
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		case stdDay:
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>			b = appendInt(b, day, 0)
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>		case stdUnderDay:
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			if day &lt; 10 {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				b = append(b, &#39; &#39;)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			b = appendInt(b, day, 0)
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		case stdZeroDay:
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>			b = appendInt(b, day, 2)
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		case stdUnderYearDay:
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>			if yday &lt; 100 {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>				b = append(b, &#39; &#39;)
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>				if yday &lt; 10 {
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>					b = append(b, &#39; &#39;)
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>				}
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>			}
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>			b = appendInt(b, yday, 0)
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		case stdZeroYearDay:
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			b = appendInt(b, yday, 3)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		case stdHour:
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			b = appendInt(b, hour, 2)
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		case stdHour12:
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>			<span class="comment">// Noon is 12PM, midnight is 12AM.</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>			hr := hour % 12
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>			if hr == 0 {
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>				hr = 12
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			}
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			b = appendInt(b, hr, 0)
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		case stdZeroHour12:
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>			<span class="comment">// Noon is 12PM, midnight is 12AM.</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>			hr := hour % 12
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>			if hr == 0 {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>				hr = 12
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>			}
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>			b = appendInt(b, hr, 2)
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		case stdMinute:
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>			b = appendInt(b, min, 0)
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		case stdZeroMinute:
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>			b = appendInt(b, min, 2)
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		case stdSecond:
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>			b = appendInt(b, sec, 0)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		case stdZeroSecond:
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>			b = appendInt(b, sec, 2)
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		case stdPM:
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>			if hour &gt;= 12 {
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>				b = append(b, &#34;PM&#34;...)
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>			} else {
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>				b = append(b, &#34;AM&#34;...)
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>			}
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>		case stdpm:
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>			if hour &gt;= 12 {
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>				b = append(b, &#34;pm&#34;...)
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>			} else {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>				b = append(b, &#34;am&#34;...)
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>			}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		case stdISO8601TZ, stdISO8601ColonTZ, stdISO8601SecondsTZ, stdISO8601ShortTZ, stdISO8601ColonSecondsTZ, stdNumTZ, stdNumColonTZ, stdNumSecondsTz, stdNumShortTZ, stdNumColonSecondsTZ:
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>			<span class="comment">// Ugly special case. We cheat and take the &#34;Z&#34; variants</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>			<span class="comment">// to mean &#34;the time zone as formatted for ISO 8601&#34;.</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>			if offset == 0 &amp;&amp; (std == stdISO8601TZ || std == stdISO8601ColonTZ || std == stdISO8601SecondsTZ || std == stdISO8601ShortTZ || std == stdISO8601ColonSecondsTZ) {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>				b = append(b, &#39;Z&#39;)
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>				break
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>			}
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			zone := offset / 60 <span class="comment">// convert to minutes</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			absoffset := offset
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>			if zone &lt; 0 {
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>				b = append(b, &#39;-&#39;)
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>				zone = -zone
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				absoffset = -absoffset
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>			} else {
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>				b = append(b, &#39;+&#39;)
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>			b = appendInt(b, zone/60, 2)
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>			if std == stdISO8601ColonTZ || std == stdNumColonTZ || std == stdISO8601ColonSecondsTZ || std == stdNumColonSecondsTZ {
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>				b = append(b, &#39;:&#39;)
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			}
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>			if std != stdNumShortTZ &amp;&amp; std != stdISO8601ShortTZ {
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>				b = appendInt(b, zone%60, 2)
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>			}
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>			<span class="comment">// append seconds if appropriate</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>			if std == stdISO8601SecondsTZ || std == stdNumSecondsTz || std == stdNumColonSecondsTZ || std == stdISO8601ColonSecondsTZ {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>				if std == stdNumColonSecondsTZ || std == stdISO8601ColonSecondsTZ {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>					b = append(b, &#39;:&#39;)
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>				}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>				b = appendInt(b, absoffset%60, 2)
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>			}
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		case stdTZ:
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>			if name != &#34;&#34; {
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>				b = append(b, name...)
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>				break
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>			}
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>			<span class="comment">// No time zone known for this time, but we must print one.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>			<span class="comment">// Use the -0700 format.</span>
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>			zone := offset / 60 <span class="comment">// convert to minutes</span>
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>			if zone &lt; 0 {
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>				b = append(b, &#39;-&#39;)
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>				zone = -zone
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>			} else {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>				b = append(b, &#39;+&#39;)
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>			}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>			b = appendInt(b, zone/60, 2)
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>			b = appendInt(b, zone%60, 2)
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		case stdFracSecond0, stdFracSecond9:
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>			b = appendNano(b, t.Nanosecond(), std)
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>	}
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>	return b
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>}
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>var errBad = errors.New(&#34;bad value for field&#34;) <span class="comment">// placeholder not passed to user</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span><span class="comment">// ParseError describes a problem parsing a time string.</span>
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>type ParseError struct {
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>	Layout     string
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>	Value      string
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>	LayoutElem string
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	ValueElem  string
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	Message    string
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>}
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span><span class="comment">// newParseError creates a new ParseError.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span><span class="comment">// The provided value and valueElem are cloned to avoid escaping their values.</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>func newParseError(layout, value, layoutElem, valueElem, message string) *ParseError {
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>	valueCopy := cloneString(value)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>	valueElemCopy := cloneString(valueElem)
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>	return &amp;ParseError{layout, valueCopy, layoutElem, valueElemCopy, message}
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>}
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span><span class="comment">// cloneString returns a string copy of s.</span>
<span id="L836" class="ln">   836&nbsp;&nbsp;</span><span class="comment">// Do not use strings.Clone to avoid dependency on strings package.</span>
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>func cloneString(s string) string {
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	return string([]byte(s))
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>}
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>
<span id="L841" class="ln">   841&nbsp;&nbsp;</span><span class="comment">// These are borrowed from unicode/utf8 and strconv and replicate behavior in</span>
<span id="L842" class="ln">   842&nbsp;&nbsp;</span><span class="comment">// that package, since we can&#39;t take a dependency on either.</span>
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>const (
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>	lowerhex  = &#34;0123456789abcdef&#34;
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>	runeSelf  = 0x80
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>	runeError = &#39;\uFFFD&#39;
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>)
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>func quote(s string) string {
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	buf := make([]byte, 1, len(s)+2) <span class="comment">// slice will be at least len(s) + quotes</span>
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	buf[0] = &#39;&#34;&#39;
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>	for i, c := range s {
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		if c &gt;= runeSelf || c &lt; &#39; &#39; {
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>			<span class="comment">// This means you are asking us to parse a time.Duration or</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			<span class="comment">// time.Location with unprintable or non-ASCII characters in it.</span>
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>			<span class="comment">// We don&#39;t expect to hit this case very often. We could try to</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>			<span class="comment">// reproduce strconv.Quote&#39;s behavior with full fidelity but</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>			<span class="comment">// given how rarely we expect to hit these edge cases, speed and</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>			<span class="comment">// conciseness are better.</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>			var width int
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>			if c == runeError {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>				width = 1
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>				if i+2 &lt; len(s) &amp;&amp; s[i:i+3] == string(runeError) {
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>					width = 3
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>				}
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>			} else {
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>				width = len(string(c))
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			}
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			for j := 0; j &lt; width; j++ {
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>				buf = append(buf, `\x`...)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>				buf = append(buf, lowerhex[s[i+j]&gt;&gt;4])
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>				buf = append(buf, lowerhex[s[i+j]&amp;0xF])
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			}
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		} else {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>			if c == &#39;&#34;&#39; || c == &#39;\\&#39; {
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>				buf = append(buf, &#39;\\&#39;)
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			}
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			buf = append(buf, string(c)...)
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		}
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	buf = append(buf, &#39;&#34;&#39;)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>	return string(buf)
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>}
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>
<span id="L885" class="ln">   885&nbsp;&nbsp;</span><span class="comment">// Error returns the string representation of a ParseError.</span>
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>func (e *ParseError) Error() string {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	if e.Message == &#34;&#34; {
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		return &#34;parsing time &#34; +
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			quote(e.Value) + &#34; as &#34; +
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>			quote(e.Layout) + &#34;: cannot parse &#34; +
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>			quote(e.ValueElem) + &#34; as &#34; +
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>			quote(e.LayoutElem)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>	return &#34;parsing time &#34; +
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		quote(e.Value) + e.Message
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>}
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span><span class="comment">// isDigit reports whether s[i] is in range and is a decimal digit.</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>func isDigit[bytes []byte | string](s bytes, i int) bool {
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>	if len(s) &lt;= i {
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		return false
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>	}
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>	c := s[i]
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>	return &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39;
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>}
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>
<span id="L907" class="ln">   907&nbsp;&nbsp;</span><span class="comment">// getnum parses s[0:1] or s[0:2] (fixed forces s[0:2])</span>
<span id="L908" class="ln">   908&nbsp;&nbsp;</span><span class="comment">// as a decimal integer and returns the integer and the</span>
<span id="L909" class="ln">   909&nbsp;&nbsp;</span><span class="comment">// remainder of the string.</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span>func getnum(s string, fixed bool) (int, string, error) {
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>	if !isDigit(s, 0) {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>		return 0, s, errBad
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	}
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>	if !isDigit(s, 1) {
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>		if fixed {
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>			return 0, s, errBad
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>		}
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>		return int(s[0] - &#39;0&#39;), s[1:], nil
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	}
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>	return int(s[0]-&#39;0&#39;)*10 + int(s[1]-&#39;0&#39;), s[2:], nil
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>}
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span><span class="comment">// getnum3 parses s[0:1], s[0:2], or s[0:3] (fixed forces s[0:3])</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span><span class="comment">// as a decimal integer and returns the integer and the remainder</span>
<span id="L925" class="ln">   925&nbsp;&nbsp;</span><span class="comment">// of the string.</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>func getnum3(s string, fixed bool) (int, string, error) {
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>	var n, i int
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	for i = 0; i &lt; 3 &amp;&amp; isDigit(s, i); i++ {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		n = n*10 + int(s[i]-&#39;0&#39;)
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>	}
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>	if i == 0 || fixed &amp;&amp; i != 3 {
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>		return 0, s, errBad
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>	}
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>	return n, s[i:], nil
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>}
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>func cutspace(s string) string {
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>	for len(s) &gt; 0 &amp;&amp; s[0] == &#39; &#39; {
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		s = s[1:]
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>	}
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>	return s
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>}
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span><span class="comment">// skip removes the given prefix from value,</span>
<span id="L945" class="ln">   945&nbsp;&nbsp;</span><span class="comment">// treating runs of space characters as equivalent.</span>
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>func skip(value, prefix string) (string, error) {
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>	for len(prefix) &gt; 0 {
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		if prefix[0] == &#39; &#39; {
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>			if len(value) &gt; 0 &amp;&amp; value[0] != &#39; &#39; {
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>				return value, errBad
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			}
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>			prefix = cutspace(prefix)
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>			value = cutspace(value)
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>			continue
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>		}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>		if len(value) == 0 || value[0] != prefix[0] {
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>			return value, errBad
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>		}
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>		prefix = prefix[1:]
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>		value = value[1:]
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>	}
<span id="L962" class="ln">   962&nbsp;&nbsp;</span>	return value, nil
<span id="L963" class="ln">   963&nbsp;&nbsp;</span>}
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>
<span id="L965" class="ln">   965&nbsp;&nbsp;</span><span class="comment">// Parse parses a formatted string and returns the time value it represents.</span>
<span id="L966" class="ln">   966&nbsp;&nbsp;</span><span class="comment">// See the documentation for the constant called Layout to see how to</span>
<span id="L967" class="ln">   967&nbsp;&nbsp;</span><span class="comment">// represent the format. The second argument must be parseable using</span>
<span id="L968" class="ln">   968&nbsp;&nbsp;</span><span class="comment">// the format string (layout) provided as the first argument.</span>
<span id="L969" class="ln">   969&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L970" class="ln">   970&nbsp;&nbsp;</span><span class="comment">// The example for Time.Format demonstrates the working of the layout string</span>
<span id="L971" class="ln">   971&nbsp;&nbsp;</span><span class="comment">// in detail and is a good reference.</span>
<span id="L972" class="ln">   972&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span><span class="comment">// When parsing (only), the input may contain a fractional second</span>
<span id="L974" class="ln">   974&nbsp;&nbsp;</span><span class="comment">// field immediately after the seconds field, even if the layout does not</span>
<span id="L975" class="ln">   975&nbsp;&nbsp;</span><span class="comment">// signify its presence. In that case either a comma or a decimal point</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span><span class="comment">// followed by a maximal series of digits is parsed as a fractional second.</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span><span class="comment">// Fractional seconds are truncated to nanosecond precision.</span>
<span id="L978" class="ln">   978&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L979" class="ln">   979&nbsp;&nbsp;</span><span class="comment">// Elements omitted from the layout are assumed to be zero or, when</span>
<span id="L980" class="ln">   980&nbsp;&nbsp;</span><span class="comment">// zero is impossible, one, so parsing &#34;3:04pm&#34; returns the time</span>
<span id="L981" class="ln">   981&nbsp;&nbsp;</span><span class="comment">// corresponding to Jan 1, year 0, 15:04:00 UTC (note that because the year is</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span><span class="comment">// 0, this time is before the zero Time).</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span><span class="comment">// Years must be in the range 0000..9999. The day of the week is checked</span>
<span id="L984" class="ln">   984&nbsp;&nbsp;</span><span class="comment">// for syntax but it is otherwise ignored.</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L986" class="ln">   986&nbsp;&nbsp;</span><span class="comment">// For layouts specifying the two-digit year 06, a value NN &gt;= 69 will be treated</span>
<span id="L987" class="ln">   987&nbsp;&nbsp;</span><span class="comment">// as 19NN and a value NN &lt; 69 will be treated as 20NN.</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L989" class="ln">   989&nbsp;&nbsp;</span><span class="comment">// The remainder of this comment describes the handling of time zones.</span>
<span id="L990" class="ln">   990&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L991" class="ln">   991&nbsp;&nbsp;</span><span class="comment">// In the absence of a time zone indicator, Parse returns a time in UTC.</span>
<span id="L992" class="ln">   992&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span><span class="comment">// When parsing a time with a zone offset like -0700, if the offset corresponds</span>
<span id="L994" class="ln">   994&nbsp;&nbsp;</span><span class="comment">// to a time zone used by the current location (Local), then Parse uses that</span>
<span id="L995" class="ln">   995&nbsp;&nbsp;</span><span class="comment">// location and zone in the returned time. Otherwise it records the time as</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span><span class="comment">// being in a fabricated location with time fixed at the given zone offset.</span>
<span id="L997" class="ln">   997&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L998" class="ln">   998&nbsp;&nbsp;</span><span class="comment">// When parsing a time with a zone abbreviation like MST, if the zone abbreviation</span>
<span id="L999" class="ln">   999&nbsp;&nbsp;</span><span class="comment">// has a defined offset in the current location, then that offset is used.</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span><span class="comment">// The zone abbreviation &#34;UTC&#34; is recognized as UTC regardless of location.</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span><span class="comment">// If the zone abbreviation is unknown, Parse records the time as being</span>
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span><span class="comment">// in a fabricated location with the given zone abbreviation and a zero offset.</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span><span class="comment">// This choice means that such a time can be parsed and reformatted with the</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span><span class="comment">// same layout losslessly, but the exact instant used in the representation will</span>
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span><span class="comment">// differ by the actual zone offset. To avoid such problems, prefer time layouts</span>
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span><span class="comment">// that use a numeric zone offset, or use ParseInLocation.</span>
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>func Parse(layout, value string) (Time, error) {
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>	<span class="comment">// Optimize for RFC3339 as it accounts for over half of all representations.</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	if layout == RFC3339 || layout == RFC3339Nano {
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>		if t, ok := parseRFC3339(value, Local); ok {
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>			return t, nil
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>		}
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	}
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	return parse(layout, value, UTC, Local)
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>}
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span><span class="comment">// ParseInLocation is like Parse but differs in two important ways.</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span><span class="comment">// First, in the absence of time zone information, Parse interprets a time as UTC;</span>
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span><span class="comment">// ParseInLocation interprets the time as in the given location.</span>
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span><span class="comment">// Second, when given a zone offset or abbreviation, Parse tries to match it</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span><span class="comment">// against the Local location; ParseInLocation uses the given location.</span>
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>func ParseInLocation(layout, value string, loc *Location) (Time, error) {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>	<span class="comment">// Optimize for RFC3339 as it accounts for over half of all representations.</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>	if layout == RFC3339 || layout == RFC3339Nano {
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>		if t, ok := parseRFC3339(value, loc); ok {
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>			return t, nil
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>		}
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>	}
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>	return parse(layout, value, loc, loc)
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>}
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>func parse(layout, value string, defaultLocation, local *Location) (Time, error) {
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>	alayout, avalue := layout, value
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>	rangeErrString := &#34;&#34; <span class="comment">// set if a value is out of range</span>
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>	amSet := false       <span class="comment">// do we need to subtract 12 from the hour for midnight?</span>
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>	pmSet := false       <span class="comment">// do we need to add 12 to the hour?</span>
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>	<span class="comment">// Time being constructed.</span>
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>	var (
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>		year       int
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		month      int = -1
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		day        int = -1
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		yday       int = -1
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>		hour       int
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		min        int
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		sec        int
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		nsec       int
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>		z          *Location
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		zoneOffset int = -1
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		zoneName   string
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>	)
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>	<span class="comment">// Each iteration processes one std value.</span>
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>	for {
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>		var err error
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>		prefix, std, suffix := nextStdChunk(layout)
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>		stdstr := layout[len(prefix) : len(layout)-len(suffix)]
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>		value, err = skip(value, prefix)
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>		if err != nil {
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, prefix, value, &#34;&#34;)
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>		}
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>		if std == 0 {
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>			if len(value) != 0 {
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>				return Time{}, newParseError(alayout, avalue, &#34;&#34;, value, &#34;: extra text: &#34;+quote(value))
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			}
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>			break
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>		}
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>		layout = suffix
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>		var p string
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>		hold := value
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>		switch std &amp; stdMask {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>		case stdYear:
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>			if len(value) &lt; 2 {
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>				err = errBad
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>				break
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>			}
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>			p, value = value[0:2], value[2:]
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>			year, err = atoi(p)
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>			if err != nil {
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>				break
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>			}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>			if year &gt;= 69 { <span class="comment">// Unix time starts Dec 31 1969 in some time zones</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>				year += 1900
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>			} else {
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>				year += 2000
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>			}
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>		case stdLongYear:
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>			if len(value) &lt; 4 || !isDigit(value, 0) {
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>				err = errBad
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>				break
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>			}
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>			p, value = value[0:4], value[4:]
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>			year, err = atoi(p)
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>		case stdMonth:
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>			month, value, err = lookup(shortMonthNames, value)
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>			month++
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>		case stdLongMonth:
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>			month, value, err = lookup(longMonthNames, value)
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>			month++
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>		case stdNumMonth, stdZeroMonth:
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>			month, value, err = getnum(value, std == stdZeroMonth)
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>			if err == nil &amp;&amp; (month &lt;= 0 || 12 &lt; month) {
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>				rangeErrString = &#34;month&#34;
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>			}
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>		case stdWeekDay:
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>			<span class="comment">// Ignore weekday except for error checking.</span>
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>			_, value, err = lookup(shortDayNames, value)
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>		case stdLongWeekDay:
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>			_, value, err = lookup(longDayNames, value)
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>		case stdDay, stdUnderDay, stdZeroDay:
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>			if std == stdUnderDay &amp;&amp; len(value) &gt; 0 &amp;&amp; value[0] == &#39; &#39; {
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>				value = value[1:]
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>			}
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>			day, value, err = getnum(value, std == stdZeroDay)
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			<span class="comment">// Note that we allow any one- or two-digit day here.</span>
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>			<span class="comment">// The month, day, year combination is validated after we&#39;ve completed parsing.</span>
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>		case stdUnderYearDay, stdZeroYearDay:
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>			for i := 0; i &lt; 2; i++ {
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>				if std == stdUnderYearDay &amp;&amp; len(value) &gt; 0 &amp;&amp; value[0] == &#39; &#39; {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>					value = value[1:]
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>				}
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>			}
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>			yday, value, err = getnum3(value, std == stdZeroYearDay)
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>			<span class="comment">// Note that we allow any one-, two-, or three-digit year-day here.</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			<span class="comment">// The year-day, year combination is validated after we&#39;ve completed parsing.</span>
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>		case stdHour:
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>			hour, value, err = getnum(value, false)
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>			if hour &lt; 0 || 24 &lt;= hour {
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>				rangeErrString = &#34;hour&#34;
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>			}
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>		case stdHour12, stdZeroHour12:
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>			hour, value, err = getnum(value, std == stdZeroHour12)
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>			if hour &lt; 0 || 12 &lt; hour {
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>				rangeErrString = &#34;hour&#34;
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>			}
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>		case stdMinute, stdZeroMinute:
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>			min, value, err = getnum(value, std == stdZeroMinute)
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>			if min &lt; 0 || 60 &lt;= min {
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>				rangeErrString = &#34;minute&#34;
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>			}
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>		case stdSecond, stdZeroSecond:
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>			sec, value, err = getnum(value, std == stdZeroSecond)
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>			if err != nil {
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>				break
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>			}
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>			if sec &lt; 0 || 60 &lt;= sec {
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>				rangeErrString = &#34;second&#34;
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>				break
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>			}
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>			<span class="comment">// Special case: do we have a fractional second but no</span>
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>			<span class="comment">// fractional second in the format?</span>
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>			if len(value) &gt;= 2 &amp;&amp; commaOrPeriod(value[0]) &amp;&amp; isDigit(value, 1) {
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>				_, std, _ = nextStdChunk(layout)
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>				std &amp;= stdMask
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>				if std == stdFracSecond0 || std == stdFracSecond9 {
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>					<span class="comment">// Fractional second in the layout; proceed normally</span>
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>					break
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>				}
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>				<span class="comment">// No fractional second in the layout but we have one in the input.</span>
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>				n := 2
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>				for ; n &lt; len(value) &amp;&amp; isDigit(value, n); n++ {
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>				}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>				nsec, rangeErrString, err = parseNanoseconds(value, n)
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>				value = value[n:]
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>			}
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>		case stdPM:
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>			if len(value) &lt; 2 {
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>				err = errBad
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>				break
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>			}
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>			p, value = value[0:2], value[2:]
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>			switch p {
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>			case &#34;PM&#34;:
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>				pmSet = true
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>			case &#34;AM&#34;:
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>				amSet = true
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>			default:
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>				err = errBad
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>			}
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>		case stdpm:
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>			if len(value) &lt; 2 {
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>				err = errBad
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>				break
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>			}
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			p, value = value[0:2], value[2:]
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>			switch p {
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			case &#34;pm&#34;:
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>				pmSet = true
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>			case &#34;am&#34;:
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>				amSet = true
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>			default:
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>				err = errBad
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>			}
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>		case stdISO8601TZ, stdISO8601ColonTZ, stdISO8601SecondsTZ, stdISO8601ShortTZ, stdISO8601ColonSecondsTZ, stdNumTZ, stdNumShortTZ, stdNumColonTZ, stdNumSecondsTz, stdNumColonSecondsTZ:
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>			if (std == stdISO8601TZ || std == stdISO8601ShortTZ || std == stdISO8601ColonTZ) &amp;&amp; len(value) &gt;= 1 &amp;&amp; value[0] == &#39;Z&#39; {
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>				value = value[1:]
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>				z = UTC
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>				break
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>			}
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>			var sign, hour, min, seconds string
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>			if std == stdISO8601ColonTZ || std == stdNumColonTZ {
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>				if len(value) &lt; 6 {
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>					err = errBad
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>					break
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>				}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>				if value[3] != &#39;:&#39; {
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span>					err = errBad
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span>					break
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>				}
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>				sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], &#34;00&#34;, value[6:]
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>			} else if std == stdNumShortTZ || std == stdISO8601ShortTZ {
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>				if len(value) &lt; 3 {
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>					err = errBad
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>					break
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>				}
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>				sign, hour, min, seconds, value = value[0:1], value[1:3], &#34;00&#34;, &#34;00&#34;, value[3:]
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>			} else if std == stdISO8601ColonSecondsTZ || std == stdNumColonSecondsTZ {
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>				if len(value) &lt; 9 {
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>					err = errBad
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>					break
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>				}
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>				if value[3] != &#39;:&#39; || value[6] != &#39;:&#39; {
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>					err = errBad
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>					break
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>				}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>				sign, hour, min, seconds, value = value[0:1], value[1:3], value[4:6], value[7:9], value[9:]
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>			} else if std == stdISO8601SecondsTZ || std == stdNumSecondsTz {
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>				if len(value) &lt; 7 {
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span>					err = errBad
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>					break
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>				}
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>				sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], value[5:7], value[7:]
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>			} else {
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>				if len(value) &lt; 5 {
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>					err = errBad
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>					break
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>				}
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>				sign, hour, min, seconds, value = value[0:1], value[1:3], value[3:5], &#34;00&#34;, value[5:]
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>			}
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>			var hr, mm, ss int
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>			hr, _, err = getnum(hour, true)
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>			if err == nil {
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>				mm, _, err = getnum(min, true)
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>			}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>			if err == nil {
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>				ss, _, err = getnum(seconds, true)
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>			}
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>			zoneOffset = (hr*60+mm)*60 + ss <span class="comment">// offset is in seconds</span>
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>			switch sign[0] {
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>			case &#39;+&#39;:
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>			case &#39;-&#39;:
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>				zoneOffset = -zoneOffset
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>			default:
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>				err = errBad
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>			}
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>		case stdTZ:
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>			<span class="comment">// Does it look like a time zone?</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>			if len(value) &gt;= 3 &amp;&amp; value[0:3] == &#34;UTC&#34; {
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>				z = UTC
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>				value = value[3:]
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>				break
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>			}
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>			n, ok := parseTimeZone(value)
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>			if !ok {
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>				err = errBad
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>				break
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>			}
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>			zoneName, value = value[:n], value[n:]
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>		case stdFracSecond0:
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>			<span class="comment">// stdFracSecond0 requires the exact number of digits as specified in</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>			<span class="comment">// the layout.</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>			ndigit := 1 + digitsLen(std)
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>			if len(value) &lt; ndigit {
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>				err = errBad
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>				break
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>			}
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>			nsec, rangeErrString, err = parseNanoseconds(value, ndigit)
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>			value = value[ndigit:]
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>		case stdFracSecond9:
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>			if len(value) &lt; 2 || !commaOrPeriod(value[0]) || value[1] &lt; &#39;0&#39; || &#39;9&#39; &lt; value[1] {
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>				<span class="comment">// Fractional second omitted.</span>
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>				break
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>			}
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>			<span class="comment">// Take any number of digits, even more than asked for,</span>
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>			<span class="comment">// because it is what the stdSecond case would do.</span>
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>			i := 0
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>			for i+1 &lt; len(value) &amp;&amp; &#39;0&#39; &lt;= value[i+1] &amp;&amp; value[i+1] &lt;= &#39;9&#39; {
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>				i++
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>			}
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>			nsec, rangeErrString, err = parseNanoseconds(value, 1+i)
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>			value = value[1+i:]
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>		}
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>		if rangeErrString != &#34;&#34; {
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, stdstr, value, &#34;: &#34;+rangeErrString+&#34; out of range&#34;)
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>		}
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>		if err != nil {
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, stdstr, hold, &#34;&#34;)
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		}
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>	}
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>	if pmSet &amp;&amp; hour &lt; 12 {
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>		hour += 12
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>	} else if amSet &amp;&amp; hour == 12 {
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>		hour = 0
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>	}
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>	<span class="comment">// Convert yday to day, month.</span>
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>	if yday &gt;= 0 {
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>		var d int
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>		var m int
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>		if isLeap(year) {
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>			if yday == 31+29 {
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>				m = int(February)
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>				d = 29
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>			} else if yday &gt; 31+29 {
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>				yday--
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>			}
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		}
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>		if yday &lt; 1 || yday &gt; 365 {
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, &#34;&#34;, value, &#34;: day-of-year out of range&#34;)
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>		}
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>		if m == 0 {
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>			m = (yday-1)/31 + 1
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>			if int(daysBefore[m]) &lt; yday {
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>				m++
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>			}
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>			d = yday - int(daysBefore[m-1])
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>		}
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>		<span class="comment">// If month, day already seen, yday&#39;s m, d must match.</span>
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, set them from m, d.</span>
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>		if month &gt;= 0 &amp;&amp; month != m {
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, &#34;&#34;, value, &#34;: day-of-year does not match month&#34;)
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>		}
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>		month = m
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>		if day &gt;= 0 &amp;&amp; day != d {
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>			return Time{}, newParseError(alayout, avalue, &#34;&#34;, value, &#34;: day-of-year does not match day&#34;)
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>		}
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>		day = d
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>	} else {
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>		if month &lt; 0 {
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>			month = int(January)
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		}
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>		if day &lt; 0 {
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>			day = 1
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>		}
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>	}
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>	<span class="comment">// Validate the day of the month.</span>
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>	if day &lt; 1 || day &gt; daysIn(Month(month), year) {
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>		return Time{}, newParseError(alayout, avalue, &#34;&#34;, value, &#34;: day out of range&#34;)
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>	}
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>	if z != nil {
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>		return Date(year, Month(month), day, hour, min, sec, nsec, z), nil
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>	}
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>	if zoneOffset != -1 {
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>		t := Date(year, Month(month), day, hour, min, sec, nsec, UTC)
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>		t.addSec(-int64(zoneOffset))
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>		<span class="comment">// Look for local zone with the given offset.</span>
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>		<span class="comment">// If that zone was in effect at the given time, use it.</span>
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>		name, offset, _, _, _ := local.lookup(t.unixSec())
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>		if offset == zoneOffset &amp;&amp; (zoneName == &#34;&#34; || name == zoneName) {
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>			t.setLoc(local)
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>			return t, nil
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>		}
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>		<span class="comment">// Otherwise create fake zone to record offset.</span>
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>		zoneNameCopy := cloneString(zoneName) <span class="comment">// avoid leaking the input value</span>
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>		t.setLoc(FixedZone(zoneNameCopy, zoneOffset))
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>		return t, nil
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>	}
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>	if zoneName != &#34;&#34; {
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>		t := Date(year, Month(month), day, hour, min, sec, nsec, UTC)
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>		<span class="comment">// Look for local zone with the given offset.</span>
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>		<span class="comment">// If that zone was in effect at the given time, use it.</span>
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>		offset, ok := local.lookupName(zoneName, t.unixSec())
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>		if ok {
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>			t.addSec(-int64(offset))
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>			t.setLoc(local)
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>			return t, nil
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>		}
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, create fake zone with unknown offset.</span>
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>		if len(zoneName) &gt; 3 &amp;&amp; zoneName[:3] == &#34;GMT&#34; {
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>			offset, _ = atoi(zoneName[3:]) <span class="comment">// Guaranteed OK by parseGMT.</span>
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>			offset *= 3600
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>		}
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>		zoneNameCopy := cloneString(zoneName) <span class="comment">// avoid leaking the input value</span>
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>		t.setLoc(FixedZone(zoneNameCopy, offset))
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>		return t, nil
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>	}
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>	<span class="comment">// Otherwise, fall back to default.</span>
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>	return Date(year, Month(month), day, hour, min, sec, nsec, defaultLocation), nil
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>}
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span><span class="comment">// parseTimeZone parses a time zone string and returns its length. Time zones</span>
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span><span class="comment">// are human-generated and unpredictable. We can&#39;t do precise error checking.</span>
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span><span class="comment">// On the other hand, for a correct parse there must be a time zone at the</span>
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span><span class="comment">// beginning of the string, so it&#39;s almost always true that there&#39;s one</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span><span class="comment">// there. We look at the beginning of the string for a run of upper-case letters.</span>
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span><span class="comment">// If there are more than 5, it&#39;s an error.</span>
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span><span class="comment">// If there are 4 or 5 and the last is a T, it&#39;s a time zone.</span>
<span id="L1408" class="ln">  1408&nbsp;&nbsp;</span><span class="comment">// If there are 3, it&#39;s a time zone.</span>
<span id="L1409" class="ln">  1409&nbsp;&nbsp;</span><span class="comment">// Otherwise, other than special cases, it&#39;s not a time zone.</span>
<span id="L1410" class="ln">  1410&nbsp;&nbsp;</span><span class="comment">// GMT is special because it can have an hour offset.</span>
<span id="L1411" class="ln">  1411&nbsp;&nbsp;</span>func parseTimeZone(value string) (length int, ok bool) {
<span id="L1412" class="ln">  1412&nbsp;&nbsp;</span>	if len(value) &lt; 3 {
<span id="L1413" class="ln">  1413&nbsp;&nbsp;</span>		return 0, false
<span id="L1414" class="ln">  1414&nbsp;&nbsp;</span>	}
<span id="L1415" class="ln">  1415&nbsp;&nbsp;</span>	<span class="comment">// Special case 1: ChST and MeST are the only zones with a lower-case letter.</span>
<span id="L1416" class="ln">  1416&nbsp;&nbsp;</span>	if len(value) &gt;= 4 &amp;&amp; (value[:4] == &#34;ChST&#34; || value[:4] == &#34;MeST&#34;) {
<span id="L1417" class="ln">  1417&nbsp;&nbsp;</span>		return 4, true
<span id="L1418" class="ln">  1418&nbsp;&nbsp;</span>	}
<span id="L1419" class="ln">  1419&nbsp;&nbsp;</span>	<span class="comment">// Special case 2: GMT may have an hour offset; treat it specially.</span>
<span id="L1420" class="ln">  1420&nbsp;&nbsp;</span>	if value[:3] == &#34;GMT&#34; {
<span id="L1421" class="ln">  1421&nbsp;&nbsp;</span>		length = parseGMT(value)
<span id="L1422" class="ln">  1422&nbsp;&nbsp;</span>		return length, true
<span id="L1423" class="ln">  1423&nbsp;&nbsp;</span>	}
<span id="L1424" class="ln">  1424&nbsp;&nbsp;</span>	<span class="comment">// Special Case 3: Some time zones are not named, but have +/-00 format</span>
<span id="L1425" class="ln">  1425&nbsp;&nbsp;</span>	if value[0] == &#39;+&#39; || value[0] == &#39;-&#39; {
<span id="L1426" class="ln">  1426&nbsp;&nbsp;</span>		length = parseSignedOffset(value)
<span id="L1427" class="ln">  1427&nbsp;&nbsp;</span>		ok := length &gt; 0 <span class="comment">// parseSignedOffset returns 0 in case of bad input</span>
<span id="L1428" class="ln">  1428&nbsp;&nbsp;</span>		return length, ok
<span id="L1429" class="ln">  1429&nbsp;&nbsp;</span>	}
<span id="L1430" class="ln">  1430&nbsp;&nbsp;</span>	<span class="comment">// How many upper-case letters are there? Need at least three, at most five.</span>
<span id="L1431" class="ln">  1431&nbsp;&nbsp;</span>	var nUpper int
<span id="L1432" class="ln">  1432&nbsp;&nbsp;</span>	for nUpper = 0; nUpper &lt; 6; nUpper++ {
<span id="L1433" class="ln">  1433&nbsp;&nbsp;</span>		if nUpper &gt;= len(value) {
<span id="L1434" class="ln">  1434&nbsp;&nbsp;</span>			break
<span id="L1435" class="ln">  1435&nbsp;&nbsp;</span>		}
<span id="L1436" class="ln">  1436&nbsp;&nbsp;</span>		if c := value[nUpper]; c &lt; &#39;A&#39; || &#39;Z&#39; &lt; c {
<span id="L1437" class="ln">  1437&nbsp;&nbsp;</span>			break
<span id="L1438" class="ln">  1438&nbsp;&nbsp;</span>		}
<span id="L1439" class="ln">  1439&nbsp;&nbsp;</span>	}
<span id="L1440" class="ln">  1440&nbsp;&nbsp;</span>	switch nUpper {
<span id="L1441" class="ln">  1441&nbsp;&nbsp;</span>	case 0, 1, 2, 6:
<span id="L1442" class="ln">  1442&nbsp;&nbsp;</span>		return 0, false
<span id="L1443" class="ln">  1443&nbsp;&nbsp;</span>	case 5: <span class="comment">// Must end in T to match.</span>
<span id="L1444" class="ln">  1444&nbsp;&nbsp;</span>		if value[4] == &#39;T&#39; {
<span id="L1445" class="ln">  1445&nbsp;&nbsp;</span>			return 5, true
<span id="L1446" class="ln">  1446&nbsp;&nbsp;</span>		}
<span id="L1447" class="ln">  1447&nbsp;&nbsp;</span>	case 4:
<span id="L1448" class="ln">  1448&nbsp;&nbsp;</span>		<span class="comment">// Must end in T, except one special case.</span>
<span id="L1449" class="ln">  1449&nbsp;&nbsp;</span>		if value[3] == &#39;T&#39; || value[:4] == &#34;WITA&#34; {
<span id="L1450" class="ln">  1450&nbsp;&nbsp;</span>			return 4, true
<span id="L1451" class="ln">  1451&nbsp;&nbsp;</span>		}
<span id="L1452" class="ln">  1452&nbsp;&nbsp;</span>	case 3:
<span id="L1453" class="ln">  1453&nbsp;&nbsp;</span>		return 3, true
<span id="L1454" class="ln">  1454&nbsp;&nbsp;</span>	}
<span id="L1455" class="ln">  1455&nbsp;&nbsp;</span>	return 0, false
<span id="L1456" class="ln">  1456&nbsp;&nbsp;</span>}
<span id="L1457" class="ln">  1457&nbsp;&nbsp;</span>
<span id="L1458" class="ln">  1458&nbsp;&nbsp;</span><span class="comment">// parseGMT parses a GMT time zone. The input string is known to start &#34;GMT&#34;.</span>
<span id="L1459" class="ln">  1459&nbsp;&nbsp;</span><span class="comment">// The function checks whether that is followed by a sign and a number in the</span>
<span id="L1460" class="ln">  1460&nbsp;&nbsp;</span><span class="comment">// range -23 through +23 excluding zero.</span>
<span id="L1461" class="ln">  1461&nbsp;&nbsp;</span>func parseGMT(value string) int {
<span id="L1462" class="ln">  1462&nbsp;&nbsp;</span>	value = value[3:]
<span id="L1463" class="ln">  1463&nbsp;&nbsp;</span>	if len(value) == 0 {
<span id="L1464" class="ln">  1464&nbsp;&nbsp;</span>		return 3
<span id="L1465" class="ln">  1465&nbsp;&nbsp;</span>	}
<span id="L1466" class="ln">  1466&nbsp;&nbsp;</span>
<span id="L1467" class="ln">  1467&nbsp;&nbsp;</span>	return 3 + parseSignedOffset(value)
<span id="L1468" class="ln">  1468&nbsp;&nbsp;</span>}
<span id="L1469" class="ln">  1469&nbsp;&nbsp;</span>
<span id="L1470" class="ln">  1470&nbsp;&nbsp;</span><span class="comment">// parseSignedOffset parses a signed timezone offset (e.g. &#34;+03&#34; or &#34;-04&#34;).</span>
<span id="L1471" class="ln">  1471&nbsp;&nbsp;</span><span class="comment">// The function checks for a signed number in the range -23 through +23 excluding zero.</span>
<span id="L1472" class="ln">  1472&nbsp;&nbsp;</span><span class="comment">// Returns length of the found offset string or 0 otherwise.</span>
<span id="L1473" class="ln">  1473&nbsp;&nbsp;</span>func parseSignedOffset(value string) int {
<span id="L1474" class="ln">  1474&nbsp;&nbsp;</span>	sign := value[0]
<span id="L1475" class="ln">  1475&nbsp;&nbsp;</span>	if sign != &#39;-&#39; &amp;&amp; sign != &#39;+&#39; {
<span id="L1476" class="ln">  1476&nbsp;&nbsp;</span>		return 0
<span id="L1477" class="ln">  1477&nbsp;&nbsp;</span>	}
<span id="L1478" class="ln">  1478&nbsp;&nbsp;</span>	x, rem, err := leadingInt(value[1:])
<span id="L1479" class="ln">  1479&nbsp;&nbsp;</span>
<span id="L1480" class="ln">  1480&nbsp;&nbsp;</span>	<span class="comment">// fail if nothing consumed by leadingInt</span>
<span id="L1481" class="ln">  1481&nbsp;&nbsp;</span>	if err != nil || value[1:] == rem {
<span id="L1482" class="ln">  1482&nbsp;&nbsp;</span>		return 0
<span id="L1483" class="ln">  1483&nbsp;&nbsp;</span>	}
<span id="L1484" class="ln">  1484&nbsp;&nbsp;</span>	if x &gt; 23 {
<span id="L1485" class="ln">  1485&nbsp;&nbsp;</span>		return 0
<span id="L1486" class="ln">  1486&nbsp;&nbsp;</span>	}
<span id="L1487" class="ln">  1487&nbsp;&nbsp;</span>	return len(value) - len(rem)
<span id="L1488" class="ln">  1488&nbsp;&nbsp;</span>}
<span id="L1489" class="ln">  1489&nbsp;&nbsp;</span>
<span id="L1490" class="ln">  1490&nbsp;&nbsp;</span>func commaOrPeriod(b byte) bool {
<span id="L1491" class="ln">  1491&nbsp;&nbsp;</span>	return b == &#39;.&#39; || b == &#39;,&#39;
<span id="L1492" class="ln">  1492&nbsp;&nbsp;</span>}
<span id="L1493" class="ln">  1493&nbsp;&nbsp;</span>
<span id="L1494" class="ln">  1494&nbsp;&nbsp;</span>func parseNanoseconds[bytes []byte | string](value bytes, nbytes int) (ns int, rangeErrString string, err error) {
<span id="L1495" class="ln">  1495&nbsp;&nbsp;</span>	if !commaOrPeriod(value[0]) {
<span id="L1496" class="ln">  1496&nbsp;&nbsp;</span>		err = errBad
<span id="L1497" class="ln">  1497&nbsp;&nbsp;</span>		return
<span id="L1498" class="ln">  1498&nbsp;&nbsp;</span>	}
<span id="L1499" class="ln">  1499&nbsp;&nbsp;</span>	if nbytes &gt; 10 {
<span id="L1500" class="ln">  1500&nbsp;&nbsp;</span>		value = value[:10]
<span id="L1501" class="ln">  1501&nbsp;&nbsp;</span>		nbytes = 10
<span id="L1502" class="ln">  1502&nbsp;&nbsp;</span>	}
<span id="L1503" class="ln">  1503&nbsp;&nbsp;</span>	if ns, err = atoi(value[1:nbytes]); err != nil {
<span id="L1504" class="ln">  1504&nbsp;&nbsp;</span>		return
<span id="L1505" class="ln">  1505&nbsp;&nbsp;</span>	}
<span id="L1506" class="ln">  1506&nbsp;&nbsp;</span>	if ns &lt; 0 {
<span id="L1507" class="ln">  1507&nbsp;&nbsp;</span>		rangeErrString = &#34;fractional second&#34;
<span id="L1508" class="ln">  1508&nbsp;&nbsp;</span>		return
<span id="L1509" class="ln">  1509&nbsp;&nbsp;</span>	}
<span id="L1510" class="ln">  1510&nbsp;&nbsp;</span>	<span class="comment">// We need nanoseconds, which means scaling by the number</span>
<span id="L1511" class="ln">  1511&nbsp;&nbsp;</span>	<span class="comment">// of missing digits in the format, maximum length 10.</span>
<span id="L1512" class="ln">  1512&nbsp;&nbsp;</span>	scaleDigits := 10 - nbytes
<span id="L1513" class="ln">  1513&nbsp;&nbsp;</span>	for i := 0; i &lt; scaleDigits; i++ {
<span id="L1514" class="ln">  1514&nbsp;&nbsp;</span>		ns *= 10
<span id="L1515" class="ln">  1515&nbsp;&nbsp;</span>	}
<span id="L1516" class="ln">  1516&nbsp;&nbsp;</span>	return
<span id="L1517" class="ln">  1517&nbsp;&nbsp;</span>}
<span id="L1518" class="ln">  1518&nbsp;&nbsp;</span>
<span id="L1519" class="ln">  1519&nbsp;&nbsp;</span>var errLeadingInt = errors.New(&#34;time: bad [0-9]*&#34;) <span class="comment">// never printed</span>
<span id="L1520" class="ln">  1520&nbsp;&nbsp;</span>
<span id="L1521" class="ln">  1521&nbsp;&nbsp;</span><span class="comment">// leadingInt consumes the leading [0-9]* from s.</span>
<span id="L1522" class="ln">  1522&nbsp;&nbsp;</span>func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
<span id="L1523" class="ln">  1523&nbsp;&nbsp;</span>	i := 0
<span id="L1524" class="ln">  1524&nbsp;&nbsp;</span>	for ; i &lt; len(s); i++ {
<span id="L1525" class="ln">  1525&nbsp;&nbsp;</span>		c := s[i]
<span id="L1526" class="ln">  1526&nbsp;&nbsp;</span>		if c &lt; &#39;0&#39; || c &gt; &#39;9&#39; {
<span id="L1527" class="ln">  1527&nbsp;&nbsp;</span>			break
<span id="L1528" class="ln">  1528&nbsp;&nbsp;</span>		}
<span id="L1529" class="ln">  1529&nbsp;&nbsp;</span>		if x &gt; 1&lt;&lt;63/10 {
<span id="L1530" class="ln">  1530&nbsp;&nbsp;</span>			<span class="comment">// overflow</span>
<span id="L1531" class="ln">  1531&nbsp;&nbsp;</span>			return 0, rem, errLeadingInt
<span id="L1532" class="ln">  1532&nbsp;&nbsp;</span>		}
<span id="L1533" class="ln">  1533&nbsp;&nbsp;</span>		x = x*10 + uint64(c) - &#39;0&#39;
<span id="L1534" class="ln">  1534&nbsp;&nbsp;</span>		if x &gt; 1&lt;&lt;63 {
<span id="L1535" class="ln">  1535&nbsp;&nbsp;</span>			<span class="comment">// overflow</span>
<span id="L1536" class="ln">  1536&nbsp;&nbsp;</span>			return 0, rem, errLeadingInt
<span id="L1537" class="ln">  1537&nbsp;&nbsp;</span>		}
<span id="L1538" class="ln">  1538&nbsp;&nbsp;</span>	}
<span id="L1539" class="ln">  1539&nbsp;&nbsp;</span>	return x, s[i:], nil
<span id="L1540" class="ln">  1540&nbsp;&nbsp;</span>}
<span id="L1541" class="ln">  1541&nbsp;&nbsp;</span>
<span id="L1542" class="ln">  1542&nbsp;&nbsp;</span><span class="comment">// leadingFraction consumes the leading [0-9]* from s.</span>
<span id="L1543" class="ln">  1543&nbsp;&nbsp;</span><span class="comment">// It is used only for fractions, so does not return an error on overflow,</span>
<span id="L1544" class="ln">  1544&nbsp;&nbsp;</span><span class="comment">// it just stops accumulating precision.</span>
<span id="L1545" class="ln">  1545&nbsp;&nbsp;</span>func leadingFraction(s string) (x uint64, scale float64, rem string) {
<span id="L1546" class="ln">  1546&nbsp;&nbsp;</span>	i := 0
<span id="L1547" class="ln">  1547&nbsp;&nbsp;</span>	scale = 1
<span id="L1548" class="ln">  1548&nbsp;&nbsp;</span>	overflow := false
<span id="L1549" class="ln">  1549&nbsp;&nbsp;</span>	for ; i &lt; len(s); i++ {
<span id="L1550" class="ln">  1550&nbsp;&nbsp;</span>		c := s[i]
<span id="L1551" class="ln">  1551&nbsp;&nbsp;</span>		if c &lt; &#39;0&#39; || c &gt; &#39;9&#39; {
<span id="L1552" class="ln">  1552&nbsp;&nbsp;</span>			break
<span id="L1553" class="ln">  1553&nbsp;&nbsp;</span>		}
<span id="L1554" class="ln">  1554&nbsp;&nbsp;</span>		if overflow {
<span id="L1555" class="ln">  1555&nbsp;&nbsp;</span>			continue
<span id="L1556" class="ln">  1556&nbsp;&nbsp;</span>		}
<span id="L1557" class="ln">  1557&nbsp;&nbsp;</span>		if x &gt; (1&lt;&lt;63-1)/10 {
<span id="L1558" class="ln">  1558&nbsp;&nbsp;</span>			<span class="comment">// It&#39;s possible for overflow to give a positive number, so take care.</span>
<span id="L1559" class="ln">  1559&nbsp;&nbsp;</span>			overflow = true
<span id="L1560" class="ln">  1560&nbsp;&nbsp;</span>			continue
<span id="L1561" class="ln">  1561&nbsp;&nbsp;</span>		}
<span id="L1562" class="ln">  1562&nbsp;&nbsp;</span>		y := x*10 + uint64(c) - &#39;0&#39;
<span id="L1563" class="ln">  1563&nbsp;&nbsp;</span>		if y &gt; 1&lt;&lt;63 {
<span id="L1564" class="ln">  1564&nbsp;&nbsp;</span>			overflow = true
<span id="L1565" class="ln">  1565&nbsp;&nbsp;</span>			continue
<span id="L1566" class="ln">  1566&nbsp;&nbsp;</span>		}
<span id="L1567" class="ln">  1567&nbsp;&nbsp;</span>		x = y
<span id="L1568" class="ln">  1568&nbsp;&nbsp;</span>		scale *= 10
<span id="L1569" class="ln">  1569&nbsp;&nbsp;</span>	}
<span id="L1570" class="ln">  1570&nbsp;&nbsp;</span>	return x, scale, s[i:]
<span id="L1571" class="ln">  1571&nbsp;&nbsp;</span>}
<span id="L1572" class="ln">  1572&nbsp;&nbsp;</span>
<span id="L1573" class="ln">  1573&nbsp;&nbsp;</span>var unitMap = map[string]uint64{
<span id="L1574" class="ln">  1574&nbsp;&nbsp;</span>	&#34;ns&#34;: uint64(Nanosecond),
<span id="L1575" class="ln">  1575&nbsp;&nbsp;</span>	&#34;us&#34;: uint64(Microsecond),
<span id="L1576" class="ln">  1576&nbsp;&nbsp;</span>	&#34;µs&#34;: uint64(Microsecond), <span class="comment">// U+00B5 = micro symbol</span>
<span id="L1577" class="ln">  1577&nbsp;&nbsp;</span>	&#34;μs&#34;: uint64(Microsecond), <span class="comment">// U+03BC = Greek letter mu</span>
<span id="L1578" class="ln">  1578&nbsp;&nbsp;</span>	&#34;ms&#34;: uint64(Millisecond),
<span id="L1579" class="ln">  1579&nbsp;&nbsp;</span>	&#34;s&#34;:  uint64(Second),
<span id="L1580" class="ln">  1580&nbsp;&nbsp;</span>	&#34;m&#34;:  uint64(Minute),
<span id="L1581" class="ln">  1581&nbsp;&nbsp;</span>	&#34;h&#34;:  uint64(Hour),
<span id="L1582" class="ln">  1582&nbsp;&nbsp;</span>}
<span id="L1583" class="ln">  1583&nbsp;&nbsp;</span>
<span id="L1584" class="ln">  1584&nbsp;&nbsp;</span><span class="comment">// ParseDuration parses a duration string.</span>
<span id="L1585" class="ln">  1585&nbsp;&nbsp;</span><span class="comment">// A duration string is a possibly signed sequence of</span>
<span id="L1586" class="ln">  1586&nbsp;&nbsp;</span><span class="comment">// decimal numbers, each with optional fraction and a unit suffix,</span>
<span id="L1587" class="ln">  1587&nbsp;&nbsp;</span><span class="comment">// such as &#34;300ms&#34;, &#34;-1.5h&#34; or &#34;2h45m&#34;.</span>
<span id="L1588" class="ln">  1588&nbsp;&nbsp;</span><span class="comment">// Valid time units are &#34;ns&#34;, &#34;us&#34; (or &#34;µs&#34;), &#34;ms&#34;, &#34;s&#34;, &#34;m&#34;, &#34;h&#34;.</span>
<span id="L1589" class="ln">  1589&nbsp;&nbsp;</span>func ParseDuration(s string) (Duration, error) {
<span id="L1590" class="ln">  1590&nbsp;&nbsp;</span>	<span class="comment">// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+</span>
<span id="L1591" class="ln">  1591&nbsp;&nbsp;</span>	orig := s
<span id="L1592" class="ln">  1592&nbsp;&nbsp;</span>	var d uint64
<span id="L1593" class="ln">  1593&nbsp;&nbsp;</span>	neg := false
<span id="L1594" class="ln">  1594&nbsp;&nbsp;</span>
<span id="L1595" class="ln">  1595&nbsp;&nbsp;</span>	<span class="comment">// Consume [-+]?</span>
<span id="L1596" class="ln">  1596&nbsp;&nbsp;</span>	if s != &#34;&#34; {
<span id="L1597" class="ln">  1597&nbsp;&nbsp;</span>		c := s[0]
<span id="L1598" class="ln">  1598&nbsp;&nbsp;</span>		if c == &#39;-&#39; || c == &#39;+&#39; {
<span id="L1599" class="ln">  1599&nbsp;&nbsp;</span>			neg = c == &#39;-&#39;
<span id="L1600" class="ln">  1600&nbsp;&nbsp;</span>			s = s[1:]
<span id="L1601" class="ln">  1601&nbsp;&nbsp;</span>		}
<span id="L1602" class="ln">  1602&nbsp;&nbsp;</span>	}
<span id="L1603" class="ln">  1603&nbsp;&nbsp;</span>	<span class="comment">// Special case: if all that is left is &#34;0&#34;, this is zero.</span>
<span id="L1604" class="ln">  1604&nbsp;&nbsp;</span>	if s == &#34;0&#34; {
<span id="L1605" class="ln">  1605&nbsp;&nbsp;</span>		return 0, nil
<span id="L1606" class="ln">  1606&nbsp;&nbsp;</span>	}
<span id="L1607" class="ln">  1607&nbsp;&nbsp;</span>	if s == &#34;&#34; {
<span id="L1608" class="ln">  1608&nbsp;&nbsp;</span>		return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1609" class="ln">  1609&nbsp;&nbsp;</span>	}
<span id="L1610" class="ln">  1610&nbsp;&nbsp;</span>	for s != &#34;&#34; {
<span id="L1611" class="ln">  1611&nbsp;&nbsp;</span>		var (
<span id="L1612" class="ln">  1612&nbsp;&nbsp;</span>			v, f  uint64      <span class="comment">// integers before, after decimal point</span>
<span id="L1613" class="ln">  1613&nbsp;&nbsp;</span>			scale float64 = 1 <span class="comment">// value = v + f/scale</span>
<span id="L1614" class="ln">  1614&nbsp;&nbsp;</span>		)
<span id="L1615" class="ln">  1615&nbsp;&nbsp;</span>
<span id="L1616" class="ln">  1616&nbsp;&nbsp;</span>		var err error
<span id="L1617" class="ln">  1617&nbsp;&nbsp;</span>
<span id="L1618" class="ln">  1618&nbsp;&nbsp;</span>		<span class="comment">// The next character must be [0-9.]</span>
<span id="L1619" class="ln">  1619&nbsp;&nbsp;</span>		if !(s[0] == &#39;.&#39; || &#39;0&#39; &lt;= s[0] &amp;&amp; s[0] &lt;= &#39;9&#39;) {
<span id="L1620" class="ln">  1620&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1621" class="ln">  1621&nbsp;&nbsp;</span>		}
<span id="L1622" class="ln">  1622&nbsp;&nbsp;</span>		<span class="comment">// Consume [0-9]*</span>
<span id="L1623" class="ln">  1623&nbsp;&nbsp;</span>		pl := len(s)
<span id="L1624" class="ln">  1624&nbsp;&nbsp;</span>		v, s, err = leadingInt(s)
<span id="L1625" class="ln">  1625&nbsp;&nbsp;</span>		if err != nil {
<span id="L1626" class="ln">  1626&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1627" class="ln">  1627&nbsp;&nbsp;</span>		}
<span id="L1628" class="ln">  1628&nbsp;&nbsp;</span>		pre := pl != len(s) <span class="comment">// whether we consumed anything before a period</span>
<span id="L1629" class="ln">  1629&nbsp;&nbsp;</span>
<span id="L1630" class="ln">  1630&nbsp;&nbsp;</span>		<span class="comment">// Consume (\.[0-9]*)?</span>
<span id="L1631" class="ln">  1631&nbsp;&nbsp;</span>		post := false
<span id="L1632" class="ln">  1632&nbsp;&nbsp;</span>		if s != &#34;&#34; &amp;&amp; s[0] == &#39;.&#39; {
<span id="L1633" class="ln">  1633&nbsp;&nbsp;</span>			s = s[1:]
<span id="L1634" class="ln">  1634&nbsp;&nbsp;</span>			pl := len(s)
<span id="L1635" class="ln">  1635&nbsp;&nbsp;</span>			f, scale, s = leadingFraction(s)
<span id="L1636" class="ln">  1636&nbsp;&nbsp;</span>			post = pl != len(s)
<span id="L1637" class="ln">  1637&nbsp;&nbsp;</span>		}
<span id="L1638" class="ln">  1638&nbsp;&nbsp;</span>		if !pre &amp;&amp; !post {
<span id="L1639" class="ln">  1639&nbsp;&nbsp;</span>			<span class="comment">// no digits (e.g. &#34;.s&#34; or &#34;-.s&#34;)</span>
<span id="L1640" class="ln">  1640&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1641" class="ln">  1641&nbsp;&nbsp;</span>		}
<span id="L1642" class="ln">  1642&nbsp;&nbsp;</span>
<span id="L1643" class="ln">  1643&nbsp;&nbsp;</span>		<span class="comment">// Consume unit.</span>
<span id="L1644" class="ln">  1644&nbsp;&nbsp;</span>		i := 0
<span id="L1645" class="ln">  1645&nbsp;&nbsp;</span>		for ; i &lt; len(s); i++ {
<span id="L1646" class="ln">  1646&nbsp;&nbsp;</span>			c := s[i]
<span id="L1647" class="ln">  1647&nbsp;&nbsp;</span>			if c == &#39;.&#39; || &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; {
<span id="L1648" class="ln">  1648&nbsp;&nbsp;</span>				break
<span id="L1649" class="ln">  1649&nbsp;&nbsp;</span>			}
<span id="L1650" class="ln">  1650&nbsp;&nbsp;</span>		}
<span id="L1651" class="ln">  1651&nbsp;&nbsp;</span>		if i == 0 {
<span id="L1652" class="ln">  1652&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: missing unit in duration &#34; + quote(orig))
<span id="L1653" class="ln">  1653&nbsp;&nbsp;</span>		}
<span id="L1654" class="ln">  1654&nbsp;&nbsp;</span>		u := s[:i]
<span id="L1655" class="ln">  1655&nbsp;&nbsp;</span>		s = s[i:]
<span id="L1656" class="ln">  1656&nbsp;&nbsp;</span>		unit, ok := unitMap[u]
<span id="L1657" class="ln">  1657&nbsp;&nbsp;</span>		if !ok {
<span id="L1658" class="ln">  1658&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: unknown unit &#34; + quote(u) + &#34; in duration &#34; + quote(orig))
<span id="L1659" class="ln">  1659&nbsp;&nbsp;</span>		}
<span id="L1660" class="ln">  1660&nbsp;&nbsp;</span>		if v &gt; 1&lt;&lt;63/unit {
<span id="L1661" class="ln">  1661&nbsp;&nbsp;</span>			<span class="comment">// overflow</span>
<span id="L1662" class="ln">  1662&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1663" class="ln">  1663&nbsp;&nbsp;</span>		}
<span id="L1664" class="ln">  1664&nbsp;&nbsp;</span>		v *= unit
<span id="L1665" class="ln">  1665&nbsp;&nbsp;</span>		if f &gt; 0 {
<span id="L1666" class="ln">  1666&nbsp;&nbsp;</span>			<span class="comment">// float64 is needed to be nanosecond accurate for fractions of hours.</span>
<span id="L1667" class="ln">  1667&nbsp;&nbsp;</span>			<span class="comment">// v &gt;= 0 &amp;&amp; (f*unit/scale) &lt;= 3.6e+12 (ns/h, h is the largest unit)</span>
<span id="L1668" class="ln">  1668&nbsp;&nbsp;</span>			v += uint64(float64(f) * (float64(unit) / scale))
<span id="L1669" class="ln">  1669&nbsp;&nbsp;</span>			if v &gt; 1&lt;&lt;63 {
<span id="L1670" class="ln">  1670&nbsp;&nbsp;</span>				<span class="comment">// overflow</span>
<span id="L1671" class="ln">  1671&nbsp;&nbsp;</span>				return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1672" class="ln">  1672&nbsp;&nbsp;</span>			}
<span id="L1673" class="ln">  1673&nbsp;&nbsp;</span>		}
<span id="L1674" class="ln">  1674&nbsp;&nbsp;</span>		d += v
<span id="L1675" class="ln">  1675&nbsp;&nbsp;</span>		if d &gt; 1&lt;&lt;63 {
<span id="L1676" class="ln">  1676&nbsp;&nbsp;</span>			return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1677" class="ln">  1677&nbsp;&nbsp;</span>		}
<span id="L1678" class="ln">  1678&nbsp;&nbsp;</span>	}
<span id="L1679" class="ln">  1679&nbsp;&nbsp;</span>	if neg {
<span id="L1680" class="ln">  1680&nbsp;&nbsp;</span>		return -Duration(d), nil
<span id="L1681" class="ln">  1681&nbsp;&nbsp;</span>	}
<span id="L1682" class="ln">  1682&nbsp;&nbsp;</span>	if d &gt; 1&lt;&lt;63-1 {
<span id="L1683" class="ln">  1683&nbsp;&nbsp;</span>		return 0, errors.New(&#34;time: invalid duration &#34; + quote(orig))
<span id="L1684" class="ln">  1684&nbsp;&nbsp;</span>	}
<span id="L1685" class="ln">  1685&nbsp;&nbsp;</span>	return Duration(d), nil
<span id="L1686" class="ln">  1686&nbsp;&nbsp;</span>}
<span id="L1687" class="ln">  1687&nbsp;&nbsp;</span>
</pre><p><a href="format.go?m=text">View as plain text</a></p>

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
