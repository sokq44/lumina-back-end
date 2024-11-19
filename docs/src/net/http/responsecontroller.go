<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/http/responsecontroller.go - Go Documentation Server</title>

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
<a href="responsecontroller.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/http">http</a>/<span class="text-muted">responsecontroller.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/http">net/http</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package http
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bufio&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;net&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// A ResponseController is used by an HTTP handler to control the response.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// A ResponseController may not be used after the [Handler.ServeHTTP] method has returned.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>type ResponseController struct {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	rw ResponseWriter
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// NewResponseController creates a [ResponseController] for a request.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// The ResponseWriter should be the original value passed to the [Handler.ServeHTTP] method,</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// or have an Unwrap method returning the original ResponseWriter.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// If the ResponseWriter implements any of the following methods, the ResponseController</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// will call them as appropriate:</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//	Flush()</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//	FlushError() error // alternative Flush returning an error</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//	Hijack() (net.Conn, *bufio.ReadWriter, error)</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//	SetReadDeadline(deadline time.Time) error</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//	SetWriteDeadline(deadline time.Time) error</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//	EnableFullDuplex() error</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// If the ResponseWriter does not support a method, ResponseController returns</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// an error matching [ErrNotSupported].</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func NewResponseController(rw ResponseWriter) *ResponseController {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	return &amp;ResponseController{rw}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>type rwUnwrapper interface {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	Unwrap() ResponseWriter
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// Flush flushes buffered data to the client.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func (c *ResponseController) Flush() error {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	rw := c.rw
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	for {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		switch t := rw.(type) {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		case interface{ FlushError() error }:
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			return t.FlushError()
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		case Flusher:
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			t.Flush()
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			return nil
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		case rwUnwrapper:
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>			rw = t.Unwrap()
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		default:
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			return errNotSupported()
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// Hijack lets the caller take over the connection.</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// See the Hijacker interface for details.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func (c *ResponseController) Hijack() (net.Conn, *bufio.ReadWriter, error) {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	rw := c.rw
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	for {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		switch t := rw.(type) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		case Hijacker:
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			return t.Hijack()
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		case rwUnwrapper:
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			rw = t.Unwrap()
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		default:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return nil, nil, errNotSupported()
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// SetReadDeadline sets the deadline for reading the entire request, including the body.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// Reads from the request body after the deadline has been exceeded will return an error.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// A zero value means no deadline.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// Setting the read deadline after it has been exceeded will not extend it.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func (c *ResponseController) SetReadDeadline(deadline time.Time) error {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	rw := c.rw
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	for {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		switch t := rw.(type) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		case interface{ SetReadDeadline(time.Time) error }:
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>			return t.SetReadDeadline(deadline)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		case rwUnwrapper:
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			rw = t.Unwrap()
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		default:
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			return errNotSupported()
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// SetWriteDeadline sets the deadline for writing the response.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span><span class="comment">// Writes to the response body after the deadline has been exceeded will not block,</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// but may succeed if the data has been buffered.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// A zero value means no deadline.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// Setting the write deadline after it has been exceeded will not extend it.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func (c *ResponseController) SetWriteDeadline(deadline time.Time) error {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	rw := c.rw
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	for {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		switch t := rw.(type) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		case interface{ SetWriteDeadline(time.Time) error }:
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			return t.SetWriteDeadline(deadline)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		case rwUnwrapper:
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			rw = t.Unwrap()
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		default:
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			return errNotSupported()
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// EnableFullDuplex indicates that the request handler will interleave reads from [Request.Body]</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// with writes to the [ResponseWriter].</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// For HTTP/1 requests, the Go HTTP server by default consumes any unread portion of</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span><span class="comment">// the request body before beginning to write the response, preventing handlers from</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// concurrently reading from the request and writing the response.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span><span class="comment">// Calling EnableFullDuplex disables this behavior and permits handlers to continue to read</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// from the request while concurrently writing the response.</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// For HTTP/2 requests, the Go HTTP server always permits concurrent reads and responses.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (c *ResponseController) EnableFullDuplex() error {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	rw := c.rw
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	for {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		switch t := rw.(type) {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		case interface{ EnableFullDuplex() error }:
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			return t.EnableFullDuplex()
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		case rwUnwrapper:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			rw = t.Unwrap()
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		default:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			return errNotSupported()
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	}
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// errNotSupported returns an error that Is ErrNotSupported,</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// but is not == to it.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>func errNotSupported() error {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return fmt.Errorf(&#34;%w&#34;, ErrNotSupported)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
</pre><p><a href="responsecontroller.go?m=text">View as plain text</a></p>

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
