<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/encoding/csv/writer.go - Go Documentation Server</title>

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
<a href="writer.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/encoding">encoding</a>/<a href="http://localhost:8080/src/encoding/csv">csv</a>/<span class="text-muted">writer.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/encoding/csv">encoding/csv</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package csv
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// A Writer writes records using CSV encoding.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// As returned by [NewWriter], a Writer writes records terminated by a</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// newline and uses &#39;,&#39; as the field delimiter. The exported fields can be</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// changed to customize the details before</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// the first call to [Writer.Write] or [Writer.WriteAll].</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// [Writer.Comma] is the field delimiter.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// If [Writer.UseCRLF] is true,</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// the Writer ends each output line with \r\n instead of \n.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// The writes of individual records are buffered.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// After all data has been written, the client should call the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// [Writer.Flush] method to guarantee all data has been forwarded to</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// the underlying [io.Writer].  Any errors that occurred should</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// be checked by calling the [Writer.Error] method.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>type Writer struct {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	Comma   rune <span class="comment">// Field delimiter (set to &#39;,&#39; by NewWriter)</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	UseCRLF bool <span class="comment">// True to use \r\n as the line terminator</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	w       *bufio.Writer
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// NewWriter returns a new Writer that writes to w.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func NewWriter(w io.Writer) *Writer {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	return &amp;Writer{
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		Comma: &#39;,&#39;,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		w:     bufio.NewWriter(w),
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// Write writes a single CSV record to w along with any necessary quoting.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// A record is a slice of strings with each string being one field.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// Writes are buffered, so [Writer.Flush] must eventually be called to ensure</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// that the record is written to the underlying [io.Writer].</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>func (w *Writer) Write(record []string) error {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	if !validDelim(w.Comma) {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		return errInvalidDelim
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	for n, field := range record {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		if n &gt; 0 {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			if _, err := w.w.WriteRune(w.Comma); err != nil {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>				return err
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		<span class="comment">// If we don&#39;t have to have a quoted field then just</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		<span class="comment">// write out the field and continue to the next field.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if !w.fieldNeedsQuotes(field) {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			if _, err := w.w.WriteString(field); err != nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>				return err
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			continue
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if err := w.w.WriteByte(&#39;&#34;&#39;); err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return err
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		for len(field) &gt; 0 {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			<span class="comment">// Search for special characters.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			i := strings.IndexAny(field, &#34;\&#34;\r\n&#34;)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			if i &lt; 0 {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>				i = len(field)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			<span class="comment">// Copy verbatim everything before the special character.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			if _, err := w.w.WriteString(field[:i]); err != nil {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				return err
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			field = field[i:]
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			<span class="comment">// Encode the special character.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			if len(field) &gt; 0 {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>				var err error
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				switch field[0] {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				case &#39;&#34;&#39;:
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>					_, err = w.w.WriteString(`&#34;&#34;`)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>				case &#39;\r&#39;:
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>					if !w.UseCRLF {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>						err = w.w.WriteByte(&#39;\r&#39;)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>					}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>				case &#39;\n&#39;:
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>					if w.UseCRLF {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>						_, err = w.w.WriteString(&#34;\r\n&#34;)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>					} else {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>						err = w.w.WriteByte(&#39;\n&#39;)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>					}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>				}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				field = field[1:]
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>				if err != nil {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>					return err
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		if err := w.w.WriteByte(&#39;&#34;&#39;); err != nil {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			return err
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	var err error
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if w.UseCRLF {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		_, err = w.w.WriteString(&#34;\r\n&#34;)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	} else {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		err = w.w.WriteByte(&#39;\n&#39;)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	return err
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// Flush writes any buffered data to the underlying [io.Writer].</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// To check if an error occurred during Flush, call [Writer.Error].</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>func (w *Writer) Flush() {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	w.w.Flush()
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// Error reports any error that has occurred during</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// a previous [Writer.Write] or [Writer.Flush].</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (w *Writer) Error() error {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	_, err := w.w.Write(nil)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	return err
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// WriteAll writes multiple CSV records to w using [Writer.Write] and</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// then calls [Writer.Flush], returning any error from the Flush.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>func (w *Writer) WriteAll(records [][]string) error {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	for _, record := range records {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		err := w.Write(record)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		if err != nil {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			return err
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return w.w.Flush()
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// fieldNeedsQuotes reports whether our field must be enclosed in quotes.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// Fields with a Comma, fields with a quote or newline, and</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">// fields which start with a space must be enclosed in quotes.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span><span class="comment">// We used to quote empty strings, but we do not anymore (as of Go 1.4).</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span><span class="comment">// The two representations should be equivalent, but Postgres distinguishes</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// quoted vs non-quoted empty string during database imports, and it has</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// an option to force the quoted behavior for non-quoted CSV but it has</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// no option to force the non-quoted behavior for quoted CSV, making</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// CSV with quoted empty strings strictly less useful.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// Not quoting the empty string also makes this package match the behavior</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// of Microsoft Excel and Google Drive.</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// For Postgres, quote the data terminating string `\.`.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>func (w *Writer) fieldNeedsQuotes(field string) bool {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if field == &#34;&#34; {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		return false
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	if field == `\.` {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		return true
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if w.Comma &lt; utf8.RuneSelf {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		for i := 0; i &lt; len(field); i++ {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			c := field[i]
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			if c == &#39;\n&#39; || c == &#39;\r&#39; || c == &#39;&#34;&#39; || c == byte(w.Comma) {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				return true
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	} else {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		if strings.ContainsRune(field, w.Comma) || strings.ContainsAny(field, &#34;\&#34;\r\n&#34;) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			return true
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	r1, _ := utf8.DecodeRuneInString(field)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	return unicode.IsSpace(r1)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
</pre><p><a href="writer.go?m=text">View as plain text</a></p>

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
