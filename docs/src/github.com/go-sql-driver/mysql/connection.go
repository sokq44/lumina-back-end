<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/go-sql-driver/mysql/connection.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="connection.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver">go-sql-driver</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver/mysql">mysql</a>/<span class="text-muted">connection.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/github.com/go-sql-driver/mysql">github.com/go-sql-driver/mysql</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Go MySQL Driver - A MySQL-Driver for Go&#39;s database/sql package</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This Source Code Form is subject to the terms of the Mozilla Public</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// License, v. 2.0. If a copy of the MPL was not distributed with this file,</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// You can obtain one at http://mozilla.org/MPL/2.0/.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package mysql
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;context&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;database/sql&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;database/sql/driver&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;net&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>type mysqlConn struct {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	buf              buffer
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	netConn          net.Conn
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	rawConn          net.Conn    <span class="comment">// underlying connection when netConn is TLS connection.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	result           mysqlResult <span class="comment">// managed by clearResult() and handleOkPacket().</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	cfg              *Config
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	connector        *connector
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	maxAllowedPacket int
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	maxWriteSize     int
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	writeTimeout     time.Duration
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	flags            clientFlag
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	status           statusFlag
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	sequence         uint8
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	parseTime        bool
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// for context support (Go 1.8+)</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	watching bool
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	watcher  chan&lt;- context.Context
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	closech  chan struct{}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	finished chan&lt;- struct{}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	canceled atomicError <span class="comment">// set non-nil if conn is canceled</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	closed   atomicBool  <span class="comment">// set when conn is closed, before closech is closed</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// Helper function to call per-connection logger.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func (mc *mysqlConn) log(v ...any) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	mc.cfg.Logger.Print(v...)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// Handles parameters set in DSN after the connection is established</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>func (mc *mysqlConn) handleParams() (err error) {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	var cmdSet strings.Builder
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	for param, val := range mc.cfg.Params {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		switch param {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		<span class="comment">// Charset: character_set_connection, character_set_client, character_set_results</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		case &#34;charset&#34;:
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			charsets := strings.Split(val, &#34;,&#34;)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			for _, cs := range charsets {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>				<span class="comment">// ignore errors here - a charset may not exist</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>				if mc.cfg.Collation != &#34;&#34; {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>					err = mc.exec(&#34;SET NAMES &#34; + cs + &#34; COLLATE &#34; + mc.cfg.Collation)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>				} else {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>					err = mc.exec(&#34;SET NAMES &#34; + cs)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>				}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>				if err == nil {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>					break
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>				}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			if err != nil {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>				return
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		<span class="comment">// Other system vars accumulated in a single SET command</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		default:
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			if cmdSet.Len() == 0 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>				<span class="comment">// Heuristic: 29 chars for each other key=value to reduce reallocations</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>				cmdSet.Grow(4 + len(param) + 3 + len(val) + 30*(len(mc.cfg.Params)-1))
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>				cmdSet.WriteString(&#34;SET &#34;)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			} else {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				cmdSet.WriteString(&#34;, &#34;)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			cmdSet.WriteString(param)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			cmdSet.WriteString(&#34; = &#34;)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>			cmdSet.WriteString(val)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	if cmdSet.Len() &gt; 0 {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		err = mc.exec(cmdSet.String())
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if err != nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			return
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	return
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>func (mc *mysqlConn) markBadConn(err error) error {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if mc == nil {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		return err
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if err != errBadConnNoWrite {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		return err
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	return driver.ErrBadConn
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>func (mc *mysqlConn) Begin() (driver.Tx, error) {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	return mc.begin(false)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>func (mc *mysqlConn) begin(readOnly bool) (driver.Tx, error) {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		mc.log(ErrInvalidConn)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	var q string
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if readOnly {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		q = &#34;START TRANSACTION READ ONLY&#34;
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	} else {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		q = &#34;START TRANSACTION&#34;
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	err := mc.exec(q)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if err == nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return &amp;mysqlTx{mc}, err
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	return nil, mc.markBadConn(err)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>func (mc *mysqlConn) Close() (err error) {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">// Makes Close idempotent</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if !mc.closed.Load() {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		err = mc.writeCommandPacket(comQuit)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	mc.cleanup()
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	mc.clearResult()
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	return
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// Closes the network connection and unsets internal variables. Do not call this</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// function after successfully authentication, call Close instead. This function</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// is called before auth or on auth failure because MySQL will have already</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// closed the network connection.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func (mc *mysqlConn) cleanup() {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if mc.closed.Swap(true) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		return
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// Makes cleanup idempotent</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	close(mc.closech)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	conn := mc.rawConn
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if conn == nil {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if err := conn.Close(); err != nil {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		mc.log(err)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// This function can be called from multiple goroutines.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// So we can not mc.clearResult() here.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// Caller should do it if they are in safe goroutine.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>func (mc *mysqlConn) error() error {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if err := mc.canceled.Value(); err != nil {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			return err
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		return ErrInvalidConn
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	return nil
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>func (mc *mysqlConn) Prepare(query string) (driver.Stmt, error) {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		mc.log(ErrInvalidConn)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	err := mc.writeCommandPacketStr(comStmtPrepare, query)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	if err != nil {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// STMT_PREPARE is safe to retry.  So we can return ErrBadConn here.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		mc.log(err)
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	stmt := &amp;mysqlStmt{
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		mc: mc,
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// Read Result</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	columnCount, err := stmt.readPrepareResultPacket()
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if err == nil {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		if stmt.paramCount &gt; 0 {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			if err = mc.readUntilEOF(); err != nil {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				return nil, err
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		if columnCount &gt; 0 {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			err = mc.readUntilEOF()
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	return stmt, err
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>func (mc *mysqlConn) interpolateParams(query string, args []driver.Value) (string, error) {
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Number of ? should be same to len(args)</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if strings.Count(query, &#34;?&#34;) != len(args) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return &#34;&#34;, driver.ErrSkip
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	buf, err := mc.buf.takeCompleteBuffer()
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	if err != nil {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// can not take the buffer. Something must be wrong with the connection</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		mc.log(err)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		return &#34;&#34;, ErrInvalidConn
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	buf = buf[:0]
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	argPos := 0
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	for i := 0; i &lt; len(query); i++ {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		q := strings.IndexByte(query[i:], &#39;?&#39;)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		if q == -1 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			buf = append(buf, query[i:]...)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			break
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		buf = append(buf, query[i:i+q]...)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		i += q
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		arg := args[argPos]
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		argPos++
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if arg == nil {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			buf = append(buf, &#34;NULL&#34;...)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			continue
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		switch v := arg.(type) {
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		case int64:
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>			buf = strconv.AppendInt(buf, v, 10)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		case uint64:
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			<span class="comment">// Handle uint64 explicitly because our custom ConvertValue emits unsigned values</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			buf = strconv.AppendUint(buf, v, 10)
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		case float64:
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			buf = strconv.AppendFloat(buf, v, &#39;g&#39;, -1, 64)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		case bool:
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			if v {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				buf = append(buf, &#39;1&#39;)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			} else {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>				buf = append(buf, &#39;0&#39;)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		case time.Time:
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			if v.IsZero() {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				buf = append(buf, &#34;&#39;0000-00-00&#39;&#34;...)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			} else {
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>				buf = append(buf, &#39;\&#39;&#39;)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>				buf, err = appendDateTime(buf, v.In(mc.cfg.Loc), mc.cfg.timeTruncate)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>				if err != nil {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>					return &#34;&#34;, err
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>				buf = append(buf, &#39;\&#39;&#39;)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		case json.RawMessage:
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			buf = append(buf, &#39;\&#39;&#39;)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			if mc.status&amp;statusNoBackslashEscapes == 0 {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>				buf = escapeBytesBackslash(buf, v)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			} else {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				buf = escapeBytesQuotes(buf, v)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			buf = append(buf, &#39;\&#39;&#39;)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		case []byte:
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			if v == nil {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>				buf = append(buf, &#34;NULL&#34;...)
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			} else {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				buf = append(buf, &#34;_binary&#39;&#34;...)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>				if mc.status&amp;statusNoBackslashEscapes == 0 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>					buf = escapeBytesBackslash(buf, v)
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>				} else {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>					buf = escapeBytesQuotes(buf, v)
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>				}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				buf = append(buf, &#39;\&#39;&#39;)
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		case string:
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			buf = append(buf, &#39;\&#39;&#39;)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>			if mc.status&amp;statusNoBackslashEscapes == 0 {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>				buf = escapeStringBackslash(buf, v)
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>			} else {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>				buf = escapeStringQuotes(buf, v)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>			buf = append(buf, &#39;\&#39;&#39;)
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		default:
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>			return &#34;&#34;, driver.ErrSkip
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		if len(buf)+4 &gt; mc.maxAllowedPacket {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			return &#34;&#34;, driver.ErrSkip
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	if argPos != len(args) {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		return &#34;&#34;, driver.ErrSkip
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	}
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	return string(buf), nil
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>}
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>func (mc *mysqlConn) Exec(query string, args []driver.Value) (driver.Result, error) {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		mc.log(ErrInvalidConn)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	if len(args) != 0 {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if !mc.cfg.InterpolateParams {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			return nil, driver.ErrSkip
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		<span class="comment">// try to interpolate the parameters to save extra roundtrips for preparing and closing a statement</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		prepared, err := mc.interpolateParams(query, args)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		if err != nil {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>			return nil, err
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		query = prepared
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	err := mc.exec(query)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if err == nil {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		copied := mc.result
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return &amp;copied, err
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	return nil, mc.markBadConn(err)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// Internal function to execute commands</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>func (mc *mysqlConn) exec(query string) error {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	handleOk := mc.clearResult()
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	if err := mc.writeCommandPacketStr(comQuery, query); err != nil {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		return mc.markBadConn(err)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// Read Result</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	resLen, err := handleOk.readResultSetHeaderPacket()
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	if err != nil {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		return err
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	if resLen &gt; 0 {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		<span class="comment">// columns</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		if err := mc.readUntilEOF(); err != nil {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			return err
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		<span class="comment">// rows</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		if err := mc.readUntilEOF(); err != nil {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			return err
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		}
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	return handleOk.discardResults()
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>func (mc *mysqlConn) Query(query string, args []driver.Value) (driver.Rows, error) {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	return mc.query(query, args)
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>func (mc *mysqlConn) query(query string, args []driver.Value) (*textRows, error) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	handleOk := mc.clearResult()
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		mc.log(ErrInvalidConn)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if len(args) != 0 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		if !mc.cfg.InterpolateParams {
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			return nil, driver.ErrSkip
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		}
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		<span class="comment">// try client-side prepare to reduce roundtrip</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		prepared, err := mc.interpolateParams(query, args)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		if err != nil {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			return nil, err
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		query = prepared
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	err := mc.writeCommandPacketStr(comQuery, query)
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	if err == nil {
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>		<span class="comment">// Read Result</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		var resLen int
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>		resLen, err = handleOk.readResultSetHeaderPacket()
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		if err == nil {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			rows := new(textRows)
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			rows.mc = mc
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>			if resLen == 0 {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>				rows.rs.done = true
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>				switch err := rows.NextResultSet(); err {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>				case nil, io.EOF:
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>					return rows, nil
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>				default:
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>					return nil, err
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>				}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			}
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			<span class="comment">// Columns</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			rows.rs.columns, err = mc.readColumns(resLen)
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>			return rows, err
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>		}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	}
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	return nil, mc.markBadConn(err)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>}
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">// Gets the value of the given MySQL System Variable</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">// The returned byte slice is only valid until the next read</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>func (mc *mysqlConn) getSystemVar(name string) ([]byte, error) {
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	handleOk := mc.clearResult()
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	if err := mc.writeCommandPacketStr(comQuery, &#34;SELECT @@&#34;+name); err != nil {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		return nil, err
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	<span class="comment">// Read Result</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	resLen, err := handleOk.readResultSetHeaderPacket()
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	if err == nil {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		rows := new(textRows)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		rows.mc = mc
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		rows.rs.columns = []mysqlField{{fieldType: fieldTypeVarChar}}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		if resLen &gt; 0 {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			<span class="comment">// Columns</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			if err := mc.readUntilEOF(); err != nil {
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>				return nil, err
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>		}
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		dest := make([]driver.Value, resLen)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		if err = rows.readRow(dest); err == nil {
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			return dest[0].([]byte), mc.readUntilEOF()
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>		}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	}
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	return nil, err
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span><span class="comment">// finish is called when the query has canceled.</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func (mc *mysqlConn) cancel(err error) {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	mc.canceled.Set(err)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	mc.cleanup()
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span><span class="comment">// finish is called when the query has succeeded.</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>func (mc *mysqlConn) finish() {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	if !mc.watching || mc.finished == nil {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		return
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	select {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	case mc.finished &lt;- struct{}{}:
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		mc.watching = false
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	case &lt;-mc.closech:
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">// Ping implements driver.Pinger interface</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>func (mc *mysqlConn) Ping(ctx context.Context) (err error) {
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>		mc.log(ErrInvalidConn)
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		return driver.ErrBadConn
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	if err = mc.watchCancel(ctx); err != nil {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>		return
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	}
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	defer mc.finish()
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	handleOk := mc.clearResult()
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	if err = mc.writeCommandPacket(comPing); err != nil {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		return mc.markBadConn(err)
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	return handleOk.readResultOK()
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>}
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">// BeginTx implements driver.ConnBeginTx interface</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>func (mc *mysqlConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if err := mc.watchCancel(ctx); err != nil {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		return nil, err
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	defer mc.finish()
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	if sql.IsolationLevel(opts.Isolation) != sql.LevelDefault {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		level, err := mapIsolationLevel(opts.Isolation)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		if err != nil {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			return nil, err
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		err = mc.exec(&#34;SET TRANSACTION ISOLATION LEVEL &#34; + level)
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		if err != nil {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>			return nil, err
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	}
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	return mc.begin(opts.ReadOnly)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>func (mc *mysqlConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	dargs, err := namedValueToValue(args)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	if err != nil {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		return nil, err
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	if err := mc.watchCancel(ctx); err != nil {
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		return nil, err
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	rows, err := mc.query(query, dargs)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	if err != nil {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		mc.finish()
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>		return nil, err
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	}
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	rows.finish = mc.finish
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>	return rows, err
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>func (mc *mysqlConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	dargs, err := namedValueToValue(args)
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	if err != nil {
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		return nil, err
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>	}
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	if err := mc.watchCancel(ctx); err != nil {
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		return nil, err
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	defer mc.finish()
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	return mc.Exec(query, dargs)
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>func (mc *mysqlConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>	if err := mc.watchCancel(ctx); err != nil {
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		return nil, err
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	}
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	stmt, err := mc.Prepare(query)
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	mc.finish()
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	if err != nil {
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		return nil, err
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	select {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	default:
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	case &lt;-ctx.Done():
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		stmt.Close()
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		return nil, ctx.Err()
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	}
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	return stmt, nil
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>func (stmt *mysqlStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	dargs, err := namedValueToValue(args)
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	if err != nil {
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>		return nil, err
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	if err := stmt.mc.watchCancel(ctx); err != nil {
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		return nil, err
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	rows, err := stmt.query(dargs)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	if err != nil {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		stmt.mc.finish()
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		return nil, err
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	rows.finish = stmt.mc.finish
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	return rows, err
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>func (stmt *mysqlStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	dargs, err := namedValueToValue(args)
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if err != nil {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return nil, err
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	if err := stmt.mc.watchCancel(ctx); err != nil {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		return nil, err
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	defer stmt.mc.finish()
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	return stmt.Exec(dargs)
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>}
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>func (mc *mysqlConn) watchCancel(ctx context.Context) error {
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	if mc.watching {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		<span class="comment">// Reach here if canceled,</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>		<span class="comment">// so the connection is already invalid</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		mc.cleanup()
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		return nil
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	<span class="comment">// When ctx is already cancelled, don&#39;t watch it.</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	if err := ctx.Err(); err != nil {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>		return err
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	<span class="comment">// When ctx is not cancellable, don&#39;t watch it.</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	if ctx.Done() == nil {
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		return nil
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	<span class="comment">// When watcher is not alive, can&#39;t watch it.</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	if mc.watcher == nil {
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>		return nil
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>	}
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	mc.watching = true
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>	mc.watcher &lt;- ctx
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	return nil
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>}
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>func (mc *mysqlConn) startWatcher() {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	watcher := make(chan context.Context, 1)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	mc.watcher = watcher
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	finished := make(chan struct{})
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	mc.finished = finished
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	go func() {
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>		for {
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>			var ctx context.Context
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>			select {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>			case ctx = &lt;-watcher:
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>			case &lt;-mc.closech:
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>				return
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>			}
<span id="L633" class="ln">   633&nbsp;&nbsp;</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>			select {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>			case &lt;-ctx.Done():
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>				mc.cancel(ctx.Err())
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>			case &lt;-finished:
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>			case &lt;-mc.closech:
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>				return
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>			}
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>		}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	}()
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>func (mc *mysqlConn) CheckNamedValue(nv *driver.NamedValue) (err error) {
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	nv.Value, err = converter{}.ConvertValue(nv.Value)
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	return
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>}
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span><span class="comment">// ResetSession implements driver.SessionResetter.</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span><span class="comment">// (From Go 1.10)</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>func (mc *mysqlConn) ResetSession(ctx context.Context) error {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>	if mc.closed.Load() {
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>		return driver.ErrBadConn
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>	}
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	<span class="comment">// Perform a stale connection check. We only perform this check for</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	<span class="comment">// the first query on a connection that has been checked out of the</span>
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	<span class="comment">// connection pool: a fresh connection from the pool is more likely</span>
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	<span class="comment">// to be stale, and it has not performed any previous writes that</span>
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>	<span class="comment">// could cause data corruption, so it&#39;s safe to return ErrBadConn</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// if the check fails.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	if mc.cfg.CheckConnLiveness {
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>		conn := mc.netConn
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>		if mc.rawConn != nil {
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			conn = mc.rawConn
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>		}
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>		var err error
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>		if mc.cfg.ReadTimeout != 0 {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>			err = conn.SetReadDeadline(time.Now().Add(mc.cfg.ReadTimeout))
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>		}
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>		if err == nil {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>			err = connCheck(conn)
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>		if err != nil {
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>			mc.log(&#34;closing bad idle connection: &#34;, err)
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>			return driver.ErrBadConn
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	return nil
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>
<span id="L684" class="ln">   684&nbsp;&nbsp;</span><span class="comment">// IsValid implements driver.Validator interface</span>
<span id="L685" class="ln">   685&nbsp;&nbsp;</span><span class="comment">// (From Go 1.15)</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>func (mc *mysqlConn) IsValid() bool {
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	return !mc.closed.Load()
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
</pre><p><a href="connection.go?m=text">View as plain text</a></p>

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
