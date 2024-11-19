<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/go-sql-driver/mysql/rows.go - Go Documentation Server</title>

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
<a href="rows.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver">go-sql-driver</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver/mysql">mysql</a>/<span class="text-muted">rows.go</span>
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
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type resultSet struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	columns     []mysqlField
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	columnNames []string
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	done        bool
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>type mysqlRows struct {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	mc     *mysqlConn
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	rs     resultSet
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	finish func()
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>type binaryRows struct {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	mysqlRows
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>type textRows struct {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	mysqlRows
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>func (rows *mysqlRows) Columns() []string {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if rows.rs.columnNames != nil {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		return rows.rs.columnNames
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	columns := make([]string, len(rows.rs.columns))
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if rows.mc != nil &amp;&amp; rows.mc.cfg.ColumnsWithAlias {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		for i := range columns {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			if tableName := rows.rs.columns[i].tableName; len(tableName) &gt; 0 {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>				columns[i] = tableName + &#34;.&#34; + rows.rs.columns[i].name
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			} else {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				columns[i] = rows.rs.columns[i].name
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	} else {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		for i := range columns {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			columns[i] = rows.rs.columns[i].name
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	rows.rs.columnNames = columns
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return columns
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func (rows *mysqlRows) ColumnTypeDatabaseTypeName(i int) string {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	return rows.rs.columns[i].typeDatabaseName()
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// func (rows *mysqlRows) ColumnTypeLength(i int) (length int64, ok bool) {</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span><span class="comment">// 	return int64(rows.rs.columns[i].length), true</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">// }</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (rows *mysqlRows) ColumnTypeNullable(i int) (nullable, ok bool) {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return rows.rs.columns[i].flags&amp;flagNotNULL == 0, true
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>func (rows *mysqlRows) ColumnTypePrecisionScale(i int) (int64, int64, bool) {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	column := rows.rs.columns[i]
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	decimals := int64(column.decimals)
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	switch column.fieldType {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	case fieldTypeDecimal, fieldTypeNewDecimal:
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if decimals &gt; 0 {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			return int64(column.length) - 2, decimals, true
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		return int64(column.length) - 1, decimals, true
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case fieldTypeTimestamp, fieldTypeDateTime, fieldTypeTime:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return decimals, decimals, true
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	case fieldTypeFloat, fieldTypeDouble:
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		if decimals == 0x1f {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			return math.MaxInt64, math.MaxInt64, true
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		return math.MaxInt64, decimals, true
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	return 0, 0, false
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func (rows *mysqlRows) ColumnTypeScanType(i int) reflect.Type {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	return rows.rs.columns[i].scanType()
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (rows *mysqlRows) Close() (err error) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if f := rows.finish; f != nil {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		f()
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		rows.finish = nil
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	mc := rows.mc
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if mc == nil {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return nil
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if err := mc.error(); err != nil {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return err
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// flip the buffer for this connection if we need to drain it.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// note that for a successful query (i.e. one where rows.next()</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// has been called until it returns false), `rows.mc` will be nil</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// by the time the user calls `(*Rows).Close`, so we won&#39;t reach this</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// see: https://github.com/golang/go/commit/651ddbdb5056ded455f47f9c494c67b389622a47</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	mc.buf.flip()
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// Remove unread packets from stream</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	if !rows.rs.done {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		err = mc.readUntilEOF()
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if err == nil {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		handleOk := mc.clearResult()
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		if err = handleOk.discardResults(); err != nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			return err
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	rows.mc = nil
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	return err
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (rows *mysqlRows) HasNextResultSet() (b bool) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	if rows.mc == nil {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		return false
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	return rows.mc.status&amp;statusMoreResultsExists != 0
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>func (rows *mysqlRows) nextResultSet() (int, error) {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if rows.mc == nil {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if err := rows.mc.error(); err != nil {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		return 0, err
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Remove unread packets from stream</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if !rows.rs.done {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		if err := rows.mc.readUntilEOF(); err != nil {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			return 0, err
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		rows.rs.done = true
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if !rows.HasNextResultSet() {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		rows.mc = nil
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return 0, io.EOF
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	rows.rs = resultSet{}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// rows.mc.affectedRows and rows.mc.insertIds accumulate on each call to</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// nextResultSet.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	resLen, err := rows.mc.resultUnchanged().readResultSetHeaderPacket()
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if err != nil {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		<span class="comment">// Clean up about multi-results flag</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		rows.rs.done = true
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		rows.mc.status = rows.mc.status &amp; (^statusMoreResultsExists)
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	return resLen, err
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>func (rows *mysqlRows) nextNotEmptyResultSet() (int, error) {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	for {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		resLen, err := rows.nextResultSet()
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if err != nil {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			return 0, err
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if resLen &gt; 0 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			return resLen, nil
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		rows.rs.done = true
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>func (rows *binaryRows) NextResultSet() error {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	resLen, err := rows.nextNotEmptyResultSet()
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if err != nil {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		return err
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	rows.rs.columns, err = rows.mc.readColumns(resLen)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	return err
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>func (rows *binaryRows) Next(dest []driver.Value) error {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if mc := rows.mc; mc != nil {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if err := mc.error(); err != nil {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			return err
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		<span class="comment">// Fetch next row from stream</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		return rows.readRow(dest)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	return io.EOF
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>func (rows *textRows) NextResultSet() (err error) {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	resLen, err := rows.nextNotEmptyResultSet()
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	if err != nil {
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		return err
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	rows.rs.columns, err = rows.mc.readColumns(resLen)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	return err
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func (rows *textRows) Next(dest []driver.Value) error {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	if mc := rows.mc; mc != nil {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if err := mc.error(); err != nil {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			return err
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		<span class="comment">// Fetch next row from stream</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		return rows.readRow(dest)
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	return io.EOF
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
</pre><p><a href="rows.go?m=text">View as plain text</a></p>

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
