<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/database/sql/convert.go - Go Documentation Server</title>

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
<a href="convert.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/database">database</a>/<a href="http://localhost:8080/src/database/sql">sql</a>/<span class="text-muted">convert.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/database/sql">database/sql</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Type conversions for Scan.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package sql
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;database/sql/driver&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;unicode&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;unicode/utf8&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>var errNilPtr = errors.New(&#34;destination pointer is nil&#34;) <span class="comment">// embedded in descriptive error</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>func describeNamedValue(nv *driver.NamedValue) string {
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	if len(nv.Name) == 0 {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		return fmt.Sprintf(&#34;$%d&#34;, nv.Ordinal)
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	}
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;with name %q&#34;, nv.Name)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>}
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>func validateNamedValueName(name string) error {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	if len(name) == 0 {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		return nil
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	r, _ := utf8.DecodeRuneInString(name)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	if unicode.IsLetter(r) {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		return nil
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	return fmt.Errorf(&#34;name %q does not begin with a letter&#34;, name)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// ccChecker wraps the driver.ColumnConverter and allows it to be used</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// as if it were a NamedValueChecker. If the driver ColumnConverter</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// is not present then the NamedValueChecker will return driver.ErrSkip.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>type ccChecker struct {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	cci  driver.ColumnConverter
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	want int
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (c ccChecker) CheckNamedValue(nv *driver.NamedValue) error {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if c.cci == nil {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		return driver.ErrSkip
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// The column converter shouldn&#39;t be called on any index</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	<span class="comment">// it isn&#39;t expecting. The final error will be thrown</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// in the argument converter loop.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	index := nv.Ordinal - 1
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	if c.want &lt;= index {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		return nil
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	}
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">// First, see if the value itself knows how to convert</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// itself to a driver type. For example, a NullString</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">// struct changing into a string or nil.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	if vr, ok := nv.Value.(driver.Valuer); ok {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		sv, err := callValuerValue(vr)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		if err != nil {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			return err
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if !driver.IsValue(sv) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;non-subset type %T returned from Value&#34;, sv)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		nv.Value = sv
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// Second, ask the column to sanity check itself. For</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// example, drivers might use this to make sure that</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// an int64 values being inserted into a 16-bit</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">// integer field is in range (before getting</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// truncated), or that a nil can&#39;t go into a NOT NULL</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">// column before going across the network to get the</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// same error.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	var err error
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	arg := nv.Value
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	nv.Value, err = c.cci.ColumnConverter(index).ConvertValue(arg)
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	if err != nil {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		return err
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if !driver.IsValue(nv.Value) {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;driver ColumnConverter error converted %T to unsupported type %T&#34;, arg, nv.Value)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	return nil
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// defaultCheckNamedValue wraps the default ColumnConverter to have the same</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// function signature as the CheckNamedValue in the driver.NamedValueChecker</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// interface.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func defaultCheckNamedValue(nv *driver.NamedValue) (err error) {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	nv.Value, err = driver.DefaultParameterConverter.ConvertValue(nv.Value)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	return err
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// driverArgsConnLocked converts arguments from callers of Stmt.Exec and</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// Stmt.Query into driver Values.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">// The statement ds may be nil, if no statement is available.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// ci must be locked.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>func driverArgsConnLocked(ci driver.Conn, ds *driverStmt, args []any) ([]driver.NamedValue, error) {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	nvargs := make([]driver.NamedValue, len(args))
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// -1 means the driver doesn&#39;t know how to count the number of</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// placeholders, so we won&#39;t sanity check input here and instead let the</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// driver deal with errors.</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	want := -1
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	var si driver.Stmt
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	var cc ccChecker
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if ds != nil {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		si = ds.si
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		want = ds.si.NumInput()
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		cc.want = want
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// Check all types of interfaces from the start.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// Drivers may opt to use the NamedValueChecker for special</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// argument types, then return driver.ErrSkip to pass it along</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// to the column converter.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	nvc, ok := si.(driver.NamedValueChecker)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if !ok {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		nvc, ok = ci.(driver.NamedValueChecker)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	cci, ok := si.(driver.ColumnConverter)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if ok {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		cc.cci = cci
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// Loop through all the arguments, checking each one.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// If no error is returned simply increment the index</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// and continue. However if driver.ErrRemoveArgument</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// is returned the argument is not included in the query</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// argument list.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	var err error
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	var n int
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	for _, arg := range args {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		nv := &amp;nvargs[n]
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if np, ok := arg.(NamedArg); ok {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			if err = validateNamedValueName(np.Name); err != nil {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				return nil, err
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>			arg = np.Value
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			nv.Name = np.Name
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		nv.Ordinal = n + 1
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		nv.Value = arg
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		<span class="comment">// Checking sequence has four routes:</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		<span class="comment">// A: 1. Default</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		<span class="comment">// B: 1. NamedValueChecker 2. Column Converter 3. Default</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		<span class="comment">// C: 1. NamedValueChecker 3. Default</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		<span class="comment">// D: 1. Column Converter 2. Default</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// The only time a Column Converter is called is first</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// or after NamedValueConverter. If first it is handled before</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		<span class="comment">// the nextCheck label. Thus for repeats tries only when the</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		<span class="comment">// NamedValueConverter is selected should the Column Converter</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		<span class="comment">// be used in the retry.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		checker := defaultCheckNamedValue
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		nextCC := false
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		switch {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		case nvc != nil:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			nextCC = cci != nil
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			checker = nvc.CheckNamedValue
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		case cci != nil:
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			checker = cc.CheckNamedValue
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	nextCheck:
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		err = checker(nv)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		switch err {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		case nil:
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			n++
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			continue
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		case driver.ErrRemoveArgument:
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			nvargs = nvargs[:len(nvargs)-1]
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			continue
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		case driver.ErrSkip:
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			if nextCC {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				nextCC = false
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				checker = cc.CheckNamedValue
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			} else {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				checker = defaultCheckNamedValue
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			goto nextCheck
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		default:
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;sql: converting argument %s type: %v&#34;, describeNamedValue(nv), err)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// Check the length of arguments after conversion to allow for omitted</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// arguments.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if want != -1 &amp;&amp; len(nvargs) != want {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return nil, fmt.Errorf(&#34;sql: expected %d arguments, got %d&#34;, want, len(nvargs))
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	return nvargs, nil
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// convertAssign is the same as convertAssignRows, but without the optional</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// rows argument.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func convertAssign(dest, src any) error {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	return convertAssignRows(dest, src, nil)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// convertAssignRows copies to dest the value in src, converting it if possible.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// An error is returned if the copy would result in loss of information.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// dest should be a pointer type. If rows is passed in, the rows will</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// be used as the parent for any cursor values converted from a</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// driver.Rows to a *Rows.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func convertAssignRows(dest, src any, rows *Rows) error {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// Common cases, without reflect.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	switch s := src.(type) {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	case string:
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		case *string:
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			if d == nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>				return errNilPtr
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			*d = s
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			return nil
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		case *[]byte:
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			if d == nil {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>				return errNilPtr
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			*d = []byte(s)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			return nil
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		case *RawBytes:
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			if d == nil {
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				return errNilPtr
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			*d = append((*d)[:0], s...)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			return nil
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	case []byte:
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		case *string:
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>			if d == nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>				return errNilPtr
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>			*d = string(s)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			return nil
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		case *any:
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>			if d == nil {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>				return errNilPtr
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>			}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			*d = bytes.Clone(s)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			return nil
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		case *[]byte:
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			if d == nil {
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>				return errNilPtr
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			}
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			*d = bytes.Clone(s)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			return nil
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		case *RawBytes:
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			if d == nil {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>				return errNilPtr
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			*d = s
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			return nil
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		}
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	case time.Time:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		case *time.Time:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			*d = s
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			return nil
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		case *string:
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>			*d = s.Format(time.RFC3339Nano)
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			return nil
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		case *[]byte:
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>			if d == nil {
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>				return errNilPtr
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>			*d = []byte(s.Format(time.RFC3339Nano))
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			return nil
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		case *RawBytes:
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			if d == nil {
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>				return errNilPtr
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			*d = s.AppendFormat((*d)[:0], time.RFC3339Nano)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			return nil
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	case decimalDecompose:
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		case decimalCompose:
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>			return d.Compose(s.Decompose(nil))
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	case nil:
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		case *any:
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			if d == nil {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>				return errNilPtr
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>			}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			*d = nil
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			return nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		case *[]byte:
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			if d == nil {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>				return errNilPtr
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>			*d = nil
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>			return nil
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		case *RawBytes:
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			if d == nil {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				return errNilPtr
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>			}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>			*d = nil
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			return nil
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// The driver is returning a cursor the client may iterate over.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	case driver.Rows:
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		switch d := dest.(type) {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		case *Rows:
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			if d == nil {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>				return errNilPtr
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>			}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			if rows == nil {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				return errors.New(&#34;invalid context to convert cursor rows, missing parent *Rows&#34;)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			rows.closemu.Lock()
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>			*d = Rows{
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>				dc:          rows.dc,
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>				releaseConn: func(error) {},
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>				rowsi:       s,
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			}
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			<span class="comment">// Chain the cancel function.</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			parentCancel := rows.cancel
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			rows.cancel = func() {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>				<span class="comment">// When Rows.cancel is called, the closemu will be locked as well.</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>				<span class="comment">// So we can access rs.lasterr.</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				d.close(rows.lasterr)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				if parentCancel != nil {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>					parentCancel()
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>				}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>			rows.closemu.Unlock()
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			return nil
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	var sv reflect.Value
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	switch d := dest.(type) {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	case *string:
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>		sv = reflect.ValueOf(src)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		switch sv.Kind() {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		case reflect.Bool,
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			reflect.Float32, reflect.Float64:
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			*d = asString(src)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			return nil
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	case *[]byte:
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		sv = reflect.ValueOf(src)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		if b, ok := asBytes(nil, sv); ok {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			*d = b
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			return nil
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	case *RawBytes:
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		sv = reflect.ValueOf(src)
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			*d = RawBytes(b)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>			return nil
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	case *bool:
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>		bv, err := driver.Bool.ConvertValue(src)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		if err == nil {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			*d = bv.(bool)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		return err
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	case *any:
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>		*d = src
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		return nil
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	if scanner, ok := dest.(Scanner); ok {
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		return scanner.Scan(src)
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	dpv := reflect.ValueOf(dest)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	if dpv.Kind() != reflect.Pointer {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		return errors.New(&#34;destination not a pointer&#34;)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	}
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	if dpv.IsNil() {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		return errNilPtr
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	}
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	if !sv.IsValid() {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>		sv = reflect.ValueOf(src)
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	dv := reflect.Indirect(dpv)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	if sv.IsValid() &amp;&amp; sv.Type().AssignableTo(dv.Type()) {
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>		switch b := src.(type) {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		case []byte:
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>			dv.Set(reflect.ValueOf(bytes.Clone(b)))
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>		default:
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			dv.Set(sv)
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		}
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>		return nil
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if dv.Kind() == sv.Kind() &amp;&amp; sv.Type().ConvertibleTo(dv.Type()) {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		dv.Set(sv.Convert(dv.Type()))
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		return nil
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	<span class="comment">// The following conversions use a string value as an intermediate representation</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// to convert between various numeric types.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// This also allows scanning into user defined types such as &#34;type Int int64&#34;.</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// For symmetry, also check for string destination types.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	switch dv.Kind() {
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	case reflect.Pointer:
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>		if src == nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			dv.SetZero()
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			return nil
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		}
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		dv.Set(reflect.New(dv.Type().Elem()))
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>		return convertAssignRows(dv.Interface(), src, rows)
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		if src == nil {
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting NULL to %s is unsupported&#34;, dv.Kind())
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		s := asString(src)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>		if err != nil {
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			err = strconvErr(err)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting driver.Value type %T (%q) to a %s: %v&#34;, src, s, dv.Kind(), err)
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		}
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>		dv.SetInt(i64)
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		return nil
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		if src == nil {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting NULL to %s is unsupported&#34;, dv.Kind())
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>		}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		s := asString(src)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		if err != nil {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			err = strconvErr(err)
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting driver.Value type %T (%q) to a %s: %v&#34;, src, s, dv.Kind(), err)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		dv.SetUint(u64)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		return nil
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	case reflect.Float32, reflect.Float64:
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		if src == nil {
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting NULL to %s is unsupported&#34;, dv.Kind())
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>		s := asString(src)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		if err != nil {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			err = strconvErr(err)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting driver.Value type %T (%q) to a %s: %v&#34;, src, s, dv.Kind(), err)
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		}
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>		dv.SetFloat(f64)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		return nil
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	case reflect.String:
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>		if src == nil {
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;converting NULL to %s is unsupported&#34;, dv.Kind())
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>		}
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		switch v := src.(type) {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		case string:
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>			dv.SetString(v)
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			return nil
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>		case []byte:
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			dv.SetString(string(v))
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>			return nil
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	}
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	return fmt.Errorf(&#34;unsupported Scan, storing driver.Value type %T into type %T&#34;, src, dest)
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>func strconvErr(err error) error {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	if ne, ok := err.(*strconv.NumError); ok {
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		return ne.Err
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	return err
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>}
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>func asString(src any) string {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	switch v := src.(type) {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	case string:
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>		return v
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	case []byte:
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		return string(v)
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>	}
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>	rv := reflect.ValueOf(src)
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>	switch rv.Kind() {
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		return strconv.FormatInt(rv.Int(), 10)
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		return strconv.FormatUint(rv.Uint(), 10)
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	case reflect.Float64:
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		return strconv.FormatFloat(rv.Float(), &#39;g&#39;, -1, 64)
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	case reflect.Float32:
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		return strconv.FormatFloat(rv.Float(), &#39;g&#39;, -1, 32)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		return strconv.FormatBool(rv.Bool())
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%v&#34;, src)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>}
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	switch rv.Kind() {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		return strconv.AppendInt(buf, rv.Int(), 10), true
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		return strconv.AppendUint(buf, rv.Uint(), 10), true
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>	case reflect.Float32:
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		return strconv.AppendFloat(buf, rv.Float(), &#39;g&#39;, -1, 32), true
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	case reflect.Float64:
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		return strconv.AppendFloat(buf, rv.Float(), &#39;g&#39;, -1, 64), true
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	case reflect.Bool:
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		return strconv.AppendBool(buf, rv.Bool()), true
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	case reflect.String:
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		s := rv.String()
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		return append(buf, s...), true
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	return
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>var valuerReflectType = reflect.TypeFor[driver.Valuer]()
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">// callValuerValue returns vr.Value(), with one exception:</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">// If vr.Value is an auto-generated method on a pointer type and the</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">// pointer is nil, it would panic at runtime in the panicwrap</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">// method. Treat it like nil instead.</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">// Issue 8415.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">// This is so people can implement driver.Value on value types and</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">// still use nil pointers to those types to mean nil/NULL, just like</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">// string/*string.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// This function is mirrored in the database/sql/driver package.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Pointer &amp;&amp;
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		rv.IsNil() &amp;&amp;
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		rv.Type().Elem().Implements(valuerReflectType) {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		return nil, nil
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	return vr.Value()
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// decimal composes or decomposes a decimal value to and from individual parts.</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">// There are four parts: a boolean negative flag, a form byte with three possible states</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span><span class="comment">// (finite=0, infinite=1, NaN=2), a base-2 big-endian integer</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span><span class="comment">// coefficient (also known as a significand) as a []byte, and an int32 exponent.</span>
<span id="L557" class="ln">   557&nbsp;&nbsp;</span><span class="comment">// These are composed into a final value as &#34;decimal = (neg) (form=finite) coefficient * 10 ^ exponent&#34;.</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span><span class="comment">// A zero length coefficient is a zero value.</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span><span class="comment">// The big-endian integer coefficient stores the most significant byte first (at coefficient[0]).</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span><span class="comment">// If the form is not finite the coefficient and exponent should be ignored.</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span><span class="comment">// The negative parameter may be set to true for any form, although implementations are not required</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span><span class="comment">// to respect the negative parameter in the non-finite form.</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span><span class="comment">// Implementations may choose to set the negative parameter to true on a zero or NaN value,</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span><span class="comment">// but implementations that do not differentiate between negative and positive</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span><span class="comment">// zero or NaN values should ignore the negative parameter without error.</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span><span class="comment">// If an implementation does not support Infinity it may be converted into a NaN without error.</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span><span class="comment">// If a value is set that is larger than what is supported by an implementation,</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span><span class="comment">// an error must be returned.</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span><span class="comment">// Implementations must return an error if a NaN or Infinity is attempted to be set while neither</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span><span class="comment">// are supported.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// NOTE(kardianos): This is an experimental interface. See https://golang.org/issue/30870</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>type decimal interface {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	decimalDecompose
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	decimalCompose
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>}
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>type decimalDecompose interface {
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	<span class="comment">// Decompose returns the internal decimal state in parts.</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	<span class="comment">// If the provided buf has sufficient capacity, buf may be returned as the coefficient with</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	<span class="comment">// the value set and length set as appropriate.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>type decimalCompose interface {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	<span class="comment">// Compose sets the internal decimal value from parts. If the value cannot be</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	<span class="comment">// represented then an error should be returned.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	Compose(form byte, negative bool, coefficient []byte, exponent int32) error
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>}
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>
</pre><p><a href="convert.go?m=text">View as plain text</a></p>

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
