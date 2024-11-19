<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/go-sql-driver/mysql/statement.go - Go Documentation Server</title>

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
<a href="statement.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver">go-sql-driver</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver/mysql">mysql</a>/<span class="text-muted">statement.go</span>
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
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;database/sql/driver&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type mysqlStmt struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	mc         *mysqlConn
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	id         uint32
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	paramCount int
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>func (stmt *mysqlStmt) Close() error {
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	if stmt.mc == nil || stmt.mc.closed.Load() {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		<span class="comment">// driver.Stmt.Close can be called more than once, thus this function</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		<span class="comment">// has to be idempotent.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		<span class="comment">// See also Issue #450 and golang/go#16019.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		<span class="comment">//errLog.Print(ErrInvalidConn)</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		return driver.ErrBadConn
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	err := stmt.mc.writeCommandPacketUint32(comStmtClose, stmt.id)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	stmt.mc = nil
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	return err
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func (stmt *mysqlStmt) NumInput() int {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	return stmt.paramCount
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func (stmt *mysqlStmt) ColumnConverter(idx int) driver.ValueConverter {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	return converter{}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>func (stmt *mysqlStmt) CheckNamedValue(nv *driver.NamedValue) (err error) {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	nv.Value, err = converter{}.ConvertValue(nv.Value)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	return
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>func (stmt *mysqlStmt) Exec(args []driver.Value) (driver.Result, error) {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if stmt.mc.closed.Load() {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		stmt.mc.log(ErrInvalidConn)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	err := stmt.writeExecutePacket(args)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if err != nil {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return nil, stmt.mc.markBadConn(err)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	mc := stmt.mc
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	handleOk := stmt.mc.clearResult()
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">// Read Result</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	resLen, err := handleOk.readResultSetHeaderPacket()
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if err != nil {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return nil, err
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	if resLen &gt; 0 {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		<span class="comment">// Columns</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if err = mc.readUntilEOF(); err != nil {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return nil, err
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		<span class="comment">// Rows</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		if err := mc.readUntilEOF(); err != nil {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			return nil, err
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	if err := handleOk.discardResults(); err != nil {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return nil, err
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	copied := mc.result
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return &amp;copied, nil
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (stmt *mysqlStmt) Query(args []driver.Value) (driver.Rows, error) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return stmt.query(args)
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (stmt *mysqlStmt) query(args []driver.Value) (*binaryRows, error) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if stmt.mc.closed.Load() {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		stmt.mc.log(ErrInvalidConn)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		return nil, driver.ErrBadConn
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// Send command</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	err := stmt.writeExecutePacket(args)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if err != nil {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		return nil, stmt.mc.markBadConn(err)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	mc := stmt.mc
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// Read Result</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	handleOk := stmt.mc.clearResult()
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	resLen, err := handleOk.readResultSetHeaderPacket()
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	if err != nil {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		return nil, err
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	rows := new(binaryRows)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if resLen &gt; 0 {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		rows.mc = mc
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		rows.rs.columns, err = mc.readColumns(resLen)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	} else {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		rows.rs.done = true
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		switch err := rows.NextResultSet(); err {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		case nil, io.EOF:
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			return rows, nil
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		default:
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			return nil, err
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	return rows, err
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>var jsonType = reflect.TypeOf(json.RawMessage{})
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>type converter struct{}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// ConvertValue mirrors the reference/default converter in database/sql/driver</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// with _one_ exception.  We support uint64 with their high bit and the default</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// implementation does not.  This function should be kept in sync with</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// database/sql/driver defaultConverter.ConvertValue() except for that</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// deliberate difference.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func (c converter) ConvertValue(v any) (driver.Value, error) {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	if driver.IsValue(v) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		return v, nil
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if vr, ok := v.(driver.Valuer); ok {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		sv, err := callValuerValue(vr)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		if err != nil {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>			return nil, err
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		if driver.IsValue(sv) {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			return sv, nil
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		}
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// A value returned from the Valuer interface can be &#34;a type handled by</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// a database driver&#39;s NamedValueChecker interface&#34; so we should accept</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// uint64 here as well.</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		if u, ok := sv.(uint64); ok {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			return u, nil
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;non-Value type %T returned from Value&#34;, sv)
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	rv := reflect.ValueOf(v)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	switch rv.Kind() {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	case reflect.Ptr:
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// indirect pointers</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if rv.IsNil() {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			return nil, nil
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		} else {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			return c.ConvertValue(rv.Elem().Interface())
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		return rv.Int(), nil
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		return rv.Uint(), nil
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	case reflect.Float32, reflect.Float64:
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		return rv.Float(), nil
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		return rv.Bool(), nil
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	case reflect.Slice:
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		switch t := rv.Type(); {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		case t == jsonType:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			return v, nil
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		case t.Elem().Kind() == reflect.Uint8:
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			return rv.Bytes(), nil
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		default:
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;unsupported type %T, a slice of %s&#34;, v, t.Elem().Kind())
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	case reflect.String:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return rv.String(), nil
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	return nil, fmt.Errorf(&#34;unsupported type %T, a %s&#34;, v, rv.Kind())
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>var valuerReflectType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">// callValuerValue returns vr.Value(), with one exception:</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// If vr.Value is an auto-generated method on a pointer type and the</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// pointer is nil, it would panic at runtime in the panicwrap</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">// method. Treat it like nil instead.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">// This is so people can implement driver.Value on value types and</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">// still use nil pointers to those types to mean nil/NULL, just like</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// string/*string.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// This is an exact copy of the same-named unexported function from the</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// database/sql package.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Ptr &amp;&amp;
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		rv.IsNil() &amp;&amp;
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		rv.Type().Elem().Implements(valuerReflectType) {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return nil, nil
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	return vr.Value()
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
</pre><p><a href="statement.go?m=text">View as plain text</a></p>

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
