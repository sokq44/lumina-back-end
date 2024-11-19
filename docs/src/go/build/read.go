<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/build/read.go - Go Documentation Server</title>

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
<a href="read.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/build">build</a>/<span class="text-muted">read.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/build">go/build</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package build
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;go/parser&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;go/scanner&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>type importReader struct {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	b    *bufio.Reader
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	buf  []byte
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	peek byte
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	err  error
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	eof  bool
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	nerr int
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	pos  token.Position
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>var bom = []byte{0xef, 0xbb, 0xbf}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func newImportReader(name string, r io.Reader) *importReader {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	b := bufio.NewReader(r)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// Remove leading UTF-8 BOM.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// Per https://golang.org/ref/spec#Source_code_representation:</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// a compiler may ignore a UTF-8-encoded byte order mark (U+FEFF)</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// if it is the first Unicode code point in the source text.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	if leadingBytes, err := b.Peek(3); err == nil &amp;&amp; bytes.Equal(leadingBytes, bom) {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		b.Discard(3)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	}
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	return &amp;importReader{
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		b: b,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		pos: token.Position{
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			Filename: name,
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			Line:     1,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			Column:   1,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		},
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func isIdent(c byte) bool {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return &#39;A&#39; &lt;= c &amp;&amp; c &lt;= &#39;Z&#39; || &#39;a&#39; &lt;= c &amp;&amp; c &lt;= &#39;z&#39; || &#39;0&#39; &lt;= c &amp;&amp; c &lt;= &#39;9&#39; || c == &#39;_&#39; || c &gt;= utf8.RuneSelf
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>var (
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	errSyntax = errors.New(&#34;syntax error&#34;)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	errNUL    = errors.New(&#34;unexpected NUL in input&#34;)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// syntaxError records a syntax error, but only if an I/O error has not already been recorded.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>func (r *importReader) syntaxError() {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if r.err == nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		r.err = errSyntax
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// readByte reads the next byte from the input, saves it in buf, and returns it.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// If an error occurs, readByte records the error in r.err and returns 0.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (r *importReader) readByte() byte {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	c, err := r.b.ReadByte()
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if err == nil {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		r.buf = append(r.buf, c)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		if c == 0 {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			err = errNUL
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if err != nil {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		if err == io.EOF {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			r.eof = true
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		} else if r.err == nil {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			r.err = err
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		}
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		c = 0
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return c
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// readByteNoBuf is like readByte but doesn&#39;t buffer the byte.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// It exhausts r.buf before reading from r.b.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>func (r *importReader) readByteNoBuf() byte {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	var c byte
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	var err error
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	if len(r.buf) &gt; 0 {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		c = r.buf[0]
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		r.buf = r.buf[1:]
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	} else {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		c, err = r.b.ReadByte()
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		if err == nil &amp;&amp; c == 0 {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			err = errNUL
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if err != nil {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if err == io.EOF {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			r.eof = true
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		} else if r.err == nil {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			r.err = err
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		}
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		return 0
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	r.pos.Offset++
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if c == &#39;\n&#39; {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		r.pos.Line++
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		r.pos.Column = 1
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	} else {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		r.pos.Column++
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	return c
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// peekByte returns the next byte from the input reader but does not advance beyond it.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// If skipSpace is set, peekByte skips leading spaces and comments.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func (r *importReader) peekByte(skipSpace bool) byte {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if r.err != nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		if r.nerr++; r.nerr &gt; 10000 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			panic(&#34;go/build: import reader looping&#34;)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		return 0
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// Use r.peek as first input byte.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">// Don&#39;t just return r.peek here: it might have been left by peekByte(false)</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// and this might be peekByte(true).</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	c := r.peek
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	if c == 0 {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		c = r.readByte()
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	for r.err == nil &amp;&amp; !r.eof {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		if skipSpace {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			<span class="comment">// For the purposes of this reader, semicolons are never necessary to</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			<span class="comment">// understand the input and are treated as spaces.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			switch c {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			case &#39; &#39;, &#39;\f&#39;, &#39;\t&#39;, &#39;\r&#39;, &#39;\n&#39;, &#39;;&#39;:
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				c = r.readByte()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				continue
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			case &#39;/&#39;:
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>				c = r.readByte()
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>				if c == &#39;/&#39; {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>					for c != &#39;\n&#39; &amp;&amp; r.err == nil &amp;&amp; !r.eof {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>						c = r.readByte()
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>					}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				} else if c == &#39;*&#39; {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>					var c1 byte
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>					for (c != &#39;*&#39; || c1 != &#39;/&#39;) &amp;&amp; r.err == nil {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>						if r.eof {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>							r.syntaxError()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>						}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>						c, c1 = c1, r.readByte()
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>					}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>				} else {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>					r.syntaxError()
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>				}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				c = r.readByte()
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>				continue
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		break
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	r.peek = c
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	return r.peek
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// nextByte is like peekByte but advances beyond the returned byte.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>func (r *importReader) nextByte(skipSpace bool) byte {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	c := r.peekByte(skipSpace)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	r.peek = 0
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	return c
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>var goEmbed = []byte(&#34;go:embed&#34;)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// findEmbed advances the input reader to the next //go:embed comment.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span><span class="comment">// It reports whether it found a comment.</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// (Otherwise it found an error or EOF.)</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>func (r *importReader) findEmbed(first bool) bool {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// The import block scan stopped after a non-space character,</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// so the reader is not at the start of a line on the first call.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// After that, each //go:embed extraction leaves the reader</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// at the end of a line.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	startLine := !first
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	var c byte
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	for r.err == nil &amp;&amp; !r.eof {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		c = r.readByteNoBuf()
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	Reswitch:
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		switch c {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		default:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			startLine = false
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		case &#39;\n&#39;:
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			startLine = true
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		case &#39; &#39;, &#39;\t&#39;:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			<span class="comment">// leave startLine alone</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		case &#39;&#34;&#39;:
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			startLine = false
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			for r.err == nil {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>				if r.eof {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>					r.syntaxError()
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>				c = r.readByteNoBuf()
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>				if c == &#39;\\&#39; {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					r.readByteNoBuf()
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>					if r.err != nil {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>						r.syntaxError()
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>						return false
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>					}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>					continue
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				if c == &#39;&#34;&#39; {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>					c = r.readByteNoBuf()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>					goto Reswitch
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>				}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			goto Reswitch
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		case &#39;`&#39;:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			startLine = false
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			for r.err == nil {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				if r.eof {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>					r.syntaxError()
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>				}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>				c = r.readByteNoBuf()
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				if c == &#39;`&#39; {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>					c = r.readByteNoBuf()
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>					goto Reswitch
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>				}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		case &#39;\&#39;&#39;:
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			startLine = false
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			for r.err == nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				if r.eof {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>					r.syntaxError()
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>				}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>				c = r.readByteNoBuf()
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>				if c == &#39;\\&#39; {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>					r.readByteNoBuf()
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>					if r.err != nil {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>						r.syntaxError()
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>						return false
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>					}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>					continue
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>				}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				if c == &#39;\&#39;&#39; {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>					c = r.readByteNoBuf()
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>					goto Reswitch
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				}
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		case &#39;/&#39;:
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			c = r.readByteNoBuf()
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			switch c {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			default:
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>				startLine = false
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>				goto Reswitch
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			case &#39;*&#39;:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				var c1 byte
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>				for (c != &#39;*&#39; || c1 != &#39;/&#39;) &amp;&amp; r.err == nil {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>					if r.eof {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>						r.syntaxError()
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>					}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>					c, c1 = c1, r.readByteNoBuf()
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>				}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				startLine = false
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			case &#39;/&#39;:
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				if startLine {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>					<span class="comment">// Try to read this as a //go:embed comment.</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>					for i := range goEmbed {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>						c = r.readByteNoBuf()
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>						if c != goEmbed[i] {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>							goto SkipSlashSlash
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>						}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>					}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>					c = r.readByteNoBuf()
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>					if c == &#39; &#39; || c == &#39;\t&#39; {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>						<span class="comment">// Found one!</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>						return true
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>					}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>				}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			SkipSlashSlash:
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>				for c != &#39;\n&#39; &amp;&amp; r.err == nil &amp;&amp; !r.eof {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>					c = r.readByteNoBuf()
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				startLine = true
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	return false
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span><span class="comment">// readKeyword reads the given keyword from the input.</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span><span class="comment">// If the keyword is not present, readKeyword records a syntax error.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>func (r *importReader) readKeyword(kw string) {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	r.peekByte(true)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	for i := 0; i &lt; len(kw); i++ {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		if r.nextByte(false) != kw[i] {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			r.syntaxError()
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			return
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if isIdent(r.peekByte(false)) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		r.syntaxError()
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span><span class="comment">// readIdent reads an identifier from the input.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span><span class="comment">// If an identifier is not present, readIdent records a syntax error.</span>
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>func (r *importReader) readIdent() {
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	c := r.peekByte(true)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if !isIdent(c) {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		r.syntaxError()
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	for isIdent(r.peekByte(false)) {
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		r.peek = 0
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// readString reads a quoted string literal from the input.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// If an identifier is not present, readString records a syntax error.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>func (r *importReader) readString() {
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	switch r.nextByte(true) {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	case &#39;`&#39;:
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		for r.err == nil {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			if r.nextByte(false) == &#39;`&#39; {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				break
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			}
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			if r.eof {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				r.syntaxError()
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		}
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	case &#39;&#34;&#39;:
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		for r.err == nil {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>			c := r.nextByte(false)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			if c == &#39;&#34;&#39; {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>				break
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			if r.eof || c == &#39;\n&#39; {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				r.syntaxError()
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			if c == &#39;\\&#39; {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>				r.nextByte(false)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	default:
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		r.syntaxError()
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	}
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span><span class="comment">// readImport reads an import clause - optional identifier followed by quoted string -</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span><span class="comment">// from the input.</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>func (r *importReader) readImport() {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	c := r.peekByte(true)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	if c == &#39;.&#39; {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		r.peek = 0
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	} else if isIdent(c) {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		r.readIdent()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	r.readString()
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// readComments is like io.ReadAll, except that it only reads the leading</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span><span class="comment">// block of comments in the file.</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>func readComments(f io.Reader) ([]byte, error) {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	r := newImportReader(&#34;&#34;, f)
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	r.peekByte(true)
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	if r.err == nil &amp;&amp; !r.eof {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		<span class="comment">// Didn&#39;t reach EOF, so must have found a non-space byte. Remove it.</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		r.buf = r.buf[:len(r.buf)-1]
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	return r.buf, r.err
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>}
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// readGoInfo expects a Go file as input and reads the file up to and including the import section.</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// It records what it learned in *info.</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// If info.fset is non-nil, readGoInfo parses the file and sets info.parsed, info.parseErr,</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// info.imports and info.embeds.</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// It only returns an error if there are problems reading the file,</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">// not for syntax errors in the file itself.</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>func readGoInfo(f io.Reader, info *fileInfo) error {
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	r := newImportReader(info.name, f)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	r.readKeyword(&#34;package&#34;)
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	r.readIdent()
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	for r.peekByte(true) == &#39;i&#39; {
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>		r.readKeyword(&#34;import&#34;)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		if r.peekByte(true) == &#39;(&#39; {
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			r.nextByte(false)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			for r.peekByte(true) != &#39;)&#39; &amp;&amp; r.err == nil {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>				r.readImport()
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			r.nextByte(false)
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>		} else {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			r.readImport()
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	info.header = r.buf
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// If we stopped successfully before EOF, we read a byte that told us we were done.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// Return all but that last byte, which would cause a syntax error if we let it through.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if r.err == nil &amp;&amp; !r.eof {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		info.header = r.buf[:len(r.buf)-1]
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// If we stopped for a syntax error, consume the whole file so that</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	<span class="comment">// we are sure we don&#39;t change the errors that go/parser returns.</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	if r.err == errSyntax {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		r.err = nil
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		for r.err == nil &amp;&amp; !r.eof {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>			r.readByte()
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		info.header = r.buf
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>	if r.err != nil {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return r.err
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	if info.fset == nil {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		return nil
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	<span class="comment">// Parse file header &amp; record imports.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	info.parsed, info.parseErr = parser.ParseFile(info.fset, info.name, info.header, parser.ImportsOnly|parser.ParseComments)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	if info.parseErr != nil {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		return nil
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>	}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	hasEmbed := false
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	for _, decl := range info.parsed.Decls {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		d, ok := decl.(*ast.GenDecl)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		if !ok {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			continue
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		for _, dspec := range d.Specs {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>			spec, ok := dspec.(*ast.ImportSpec)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			if !ok {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				continue
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>			}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>			quoted := spec.Path.Value
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>			path, err := strconv.Unquote(quoted)
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			if err != nil {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;parser returned invalid quoted string: &lt;%s&gt;&#34;, quoted)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			if !isValidImport(path) {
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>				<span class="comment">// The parser used to return a parse error for invalid import paths, but</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>				<span class="comment">// no longer does, so check for and create the error here instead.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>				info.parseErr = scanner.Error{Pos: info.fset.Position(spec.Pos()), Msg: &#34;invalid import path: &#34; + path}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>				info.imports = nil
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>				return nil
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>			}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>			if path == &#34;embed&#34; {
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>				hasEmbed = true
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			}
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			doc := spec.Doc
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			if doc == nil &amp;&amp; len(d.Specs) == 1 {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>				doc = d.Doc
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			info.imports = append(info.imports, fileImport{path, spec.Pos(), doc})
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>		}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	<span class="comment">// Extract directives.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	for _, group := range info.parsed.Comments {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		if group.Pos() &gt;= info.parsed.Package {
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			break
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		}
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		for _, c := range group.List {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			if strings.HasPrefix(c.Text, &#34;//go:&#34;) {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>				info.directives = append(info.directives, Directive{c.Text, info.fset.Position(c.Slash)})
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>			}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		}
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	<span class="comment">// If the file imports &#34;embed&#34;,</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	<span class="comment">// we have to look for //go:embed comments</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	<span class="comment">// in the remainder of the file.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	<span class="comment">// The compiler will enforce the mapping of comments to</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	<span class="comment">// declared variables. We just need to know the patterns.</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	<span class="comment">// If there were //go:embed comments earlier in the file</span>
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	<span class="comment">// (near the package statement or imports), the compiler</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	<span class="comment">// will reject them. They can be (and have already been) ignored.</span>
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	if hasEmbed {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		var line []byte
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>		for first := true; r.findEmbed(first); first = false {
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>			line = line[:0]
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>			pos := r.pos
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>			for {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>				c := r.readByteNoBuf()
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>				if c == &#39;\n&#39; || r.err != nil || r.eof {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>					break
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>				}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>				line = append(line, c)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			<span class="comment">// Add args if line is well-formed.</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			<span class="comment">// Ignore badly-formed lines - the compiler will report them when it finds them,</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			<span class="comment">// and we can pretend they are not there to help go list succeed with what it knows.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			embs, err := parseGoEmbed(string(line), pos)
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>			if err == nil {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>				info.embeds = append(info.embeds, embs...)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	return nil
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>}
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span><span class="comment">// isValidImport checks if the import is a valid import using the more strict</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">// checks allowed by the implementation restriction in https://go.dev/ref/spec#Import_declarations.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span><span class="comment">// It was ported from the function of the same name that was removed from the</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span><span class="comment">// parser in CL 424855, when the parser stopped doing these checks.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>func isValidImport(s string) bool {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	const illegalChars = `!&#34;#$%&amp;&#39;()*,:;&lt;=&gt;?[\]^{|}` + &#34;`\uFFFD&#34;
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>	for _, r := range s {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			return false
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	}
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	return s != &#34;&#34;
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// parseGoEmbed parses the text following &#34;//go:embed&#34; to extract the glob patterns.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// It accepts unquoted space-separated patterns as well as double-quoted and back-quoted Go strings.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// This is based on a similar function in cmd/compile/internal/gc/noder.go;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">// this version calculates position information as well.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>func parseGoEmbed(args string, pos token.Position) ([]fileEmbed, error) {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	trimBytes := func(n int) {
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		pos.Offset += n
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		pos.Column += utf8.RuneCountInString(args[:n])
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		args = args[n:]
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	trimSpace := func() {
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		trim := strings.TrimLeftFunc(args, unicode.IsSpace)
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>		trimBytes(len(args) - len(trim))
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	}
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	var list []fileEmbed
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	for trimSpace(); args != &#34;&#34;; trimSpace() {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		var path string
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		pathPos := pos
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	Switch:
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		switch args[0] {
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>		default:
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			i := len(args)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			for j, c := range args {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>				if unicode.IsSpace(c) {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>					i = j
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>					break
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>				}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>			}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>			path = args[:i]
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>			trimBytes(i)
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		case &#39;`&#39;:
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>			var ok bool
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>			path, _, ok = strings.Cut(args[1:], &#34;`&#34;)
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>			if !ok {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>				return nil, fmt.Errorf(&#34;invalid quoted string in //go:embed: %s&#34;, args)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>			}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>			trimBytes(1 + len(path) + 1)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		case &#39;&#34;&#39;:
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			i := 1
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>			for ; i &lt; len(args); i++ {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>				if args[i] == &#39;\\&#39; {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>					i++
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>					continue
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>				}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>				if args[i] == &#39;&#34;&#39; {
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>					q, err := strconv.Unquote(args[:i+1])
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>					if err != nil {
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>						return nil, fmt.Errorf(&#34;invalid quoted string in //go:embed: %s&#34;, args[:i+1])
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>					}
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>					path = q
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>					trimBytes(i + 1)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>					break Switch
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>				}
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>			}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>			if i &gt;= len(args) {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>				return nil, fmt.Errorf(&#34;invalid quoted string in //go:embed: %s&#34;, args)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>			}
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>		}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		if args != &#34;&#34; {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			r, _ := utf8.DecodeRuneInString(args)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			if !unicode.IsSpace(r) {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>				return nil, fmt.Errorf(&#34;invalid quoted string in //go:embed: %s&#34;, args)
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>		list = append(list, fileEmbed{path, pathPos})
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	}
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	return list, nil
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>}
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
</pre><p><a href="read.go?m=text">View as plain text</a></p>

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
