<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/go-sql-driver/mysql/packets.go - Go Documentation Server</title>

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
<a href="packets.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver">go-sql-driver</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver/mysql">mysql</a>/<span class="text-muted">packets.go</span>
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
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;crypto/tls&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;database/sql/driver&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Packets documentation:</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/client-server-protocol.html</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Read packet to buffer &#39;data&#39;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>func (mc *mysqlConn) readPacket() ([]byte, error) {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	var prevData []byte
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	for {
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		<span class="comment">// read packet header</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		data, err := mc.buf.readNext(4)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		if err != nil {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			if cerr := mc.canceled.Value(); cerr != nil {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>				return nil, cerr
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>			}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>			mc.log(err)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			mc.Close()
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>			return nil, ErrInvalidConn
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		<span class="comment">// packet length [24 bit]</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		pktLen := int(uint32(data[0]) | uint32(data[1])&lt;&lt;8 | uint32(data[2])&lt;&lt;16)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// check packet sync [8 bit]</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		if data[3] != mc.sequence {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>			mc.Close()
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			if data[3] &gt; mc.sequence {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				return nil, ErrPktSyncMul
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			return nil, ErrPktSync
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		mc.sequence++
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		<span class="comment">// packets with length 0 terminate a previous packet which is a</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		<span class="comment">// multiple of (2^24)-1 bytes long</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		if pktLen == 0 {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			<span class="comment">// there was no previous packet</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			if prevData == nil {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>				mc.log(ErrMalformPkt)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>				mc.Close()
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>				return nil, ErrInvalidConn
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			return prevData, nil
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		<span class="comment">// read packet body [pktLen bytes]</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		data, err = mc.buf.readNext(pktLen)
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		if err != nil {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			if cerr := mc.canceled.Value(); cerr != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>				return nil, cerr
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			mc.log(err)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			mc.Close()
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			return nil, ErrInvalidConn
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		<span class="comment">// return data if this was the last packet</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		if pktLen &lt; maxPacketSize {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>			<span class="comment">// zero allocations for non-split packets</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			if prevData == nil {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				return data, nil
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>			}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			return append(prevData, data...), nil
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		prevData = append(prevData, data...)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// Write packet buffer &#39;data&#39;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func (mc *mysqlConn) writePacket(data []byte) error {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	pktLen := len(data) - 4
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if pktLen &gt; mc.maxAllowedPacket {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return ErrPktTooLarge
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	for {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		var size int
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		if pktLen &gt;= maxPacketSize {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			data[0] = 0xff
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			data[1] = 0xff
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			data[2] = 0xff
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			size = maxPacketSize
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		} else {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>			data[0] = byte(pktLen)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>			data[1] = byte(pktLen &gt;&gt; 8)
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			data[2] = byte(pktLen &gt;&gt; 16)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			size = pktLen
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		data[3] = mc.sequence
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		<span class="comment">// Write packet</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		if mc.writeTimeout &gt; 0 {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			if err := mc.netConn.SetWriteDeadline(time.Now().Add(mc.writeTimeout)); err != nil {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>				return err
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		n, err := mc.netConn.Write(data[:4+size])
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		if err == nil &amp;&amp; n == 4+size {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			mc.sequence++
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			if size != maxPacketSize {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				return nil
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			pktLen -= size
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			data = data[size:]
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			continue
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// Handle error</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		if err == nil { <span class="comment">// n != len(data)</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			mc.cleanup()
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			mc.log(ErrMalformPkt)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		} else {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if cerr := mc.canceled.Value(); cerr != nil {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				return cerr
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			if n == 0 &amp;&amp; pktLen == len(data)-4 {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				<span class="comment">// only for the first loop iteration when nothing was written yet</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>				return errBadConnNoWrite
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			mc.cleanup()
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			mc.log(err)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		return ErrInvalidConn
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">/******************************************************************************
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>*                           Initialization Process                            *
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>******************************************************************************/</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// Handshake Initialization Packet</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>func (mc *mysqlConn) readHandshakePacket() (data []byte, plugin string, err error) {
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	data, err = mc.readPacket()
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if err != nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		<span class="comment">// for init we can rewrite this to ErrBadConn for sql.Driver to retry, since</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		<span class="comment">// in connection initialization we don&#39;t risk retrying non-idempotent actions.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		if err == ErrInvalidConn {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			return nil, &#34;&#34;, driver.ErrBadConn
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		return
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if data[0] == iERR {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return nil, &#34;&#34;, mc.handleErrorPacket(data)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// protocol version [1 byte]</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	if data[0] &lt; minProtocolVersion {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		return nil, &#34;&#34;, fmt.Errorf(
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			&#34;unsupported protocol version %d. Version %d or higher is required&#34;,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			data[0],
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			minProtocolVersion,
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// server version [null terminated string]</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">// connection id [4 bytes]</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">// first part of the password cipher [8 bytes]</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	authData := data[pos : pos+8]
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">// (filler) always 0x00 [1 byte]</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	pos += 8 + 1
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// capability flags (lower 2 bytes) [2 bytes]</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	mc.flags = clientFlag(binary.LittleEndian.Uint16(data[pos : pos+2]))
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	if mc.flags&amp;clientProtocol41 == 0 {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		return nil, &#34;&#34;, ErrOldProtocol
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	if mc.flags&amp;clientSSL == 0 &amp;&amp; mc.cfg.TLS != nil {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		if mc.cfg.AllowFallbackToPlaintext {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			mc.cfg.TLS = nil
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		} else {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			return nil, &#34;&#34;, ErrNoTLS
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>		}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	pos += 2
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	if len(data) &gt; pos {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		<span class="comment">// character set [1 byte]</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		<span class="comment">// status flags [2 bytes]</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		<span class="comment">// capability flags (upper 2 bytes) [2 bytes]</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		<span class="comment">// length of auth-plugin-data [1 byte]</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		<span class="comment">// reserved (all [00]) [10 bytes]</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		pos += 1 + 2 + 2 + 1 + 10
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		<span class="comment">// second part of the password cipher [minimum 13 bytes],</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		<span class="comment">// where len=MAX(13, length of auth-plugin-data - 8)</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		<span class="comment">// The web documentation is ambiguous about the length. However,</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		<span class="comment">// according to mysql-5.7/sql/auth/sql_authentication.cc line 538,</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		<span class="comment">// the 13th byte is &#34;\0 byte, terminating the second part of</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		<span class="comment">// a scramble&#34;. So the second part of the password cipher is</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// a NULL terminated string that&#39;s at least 13 bytes with the</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// last byte being NULL.</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		<span class="comment">// The official Python library uses the fixed length 12</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		<span class="comment">// which seems to work but technically could have a hidden bug.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		authData = append(authData, data[pos:pos+12]...)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		pos += 13
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// EOF if version (&gt;= 5.5.7 and &lt; 5.5.10) or (&gt;= 5.6.0 and &lt; 5.6.2)</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// \NUL otherwise</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		if end := bytes.IndexByte(data[pos:], 0x00); end != -1 {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			plugin = string(data[pos : pos+end])
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		} else {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			plugin = string(data[pos:])
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		}
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		<span class="comment">// make a memory safe copy of the cipher slice</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		var b [20]byte
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		copy(b[:], authData)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		return b[:], plugin, nil
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// make a memory safe copy of the cipher slice</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	var b [8]byte
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	copy(b[:], authData)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	return b[:], plugin, nil
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span><span class="comment">// Client Authentication Packet</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>func (mc *mysqlConn) writeHandshakeResponsePacket(authResp []byte, plugin string) error {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// Adjust client flags based on server support</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	clientFlags := clientProtocol41 |
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		clientSecureConn |
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		clientLongPassword |
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		clientTransactions |
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		clientLocalFiles |
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		clientPluginAuth |
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		clientMultiResults |
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		clientConnectAttrs |
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		mc.flags&amp;clientLongFlag
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if mc.cfg.ClientFoundRows {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		clientFlags |= clientFoundRows
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// To enable TLS / SSL</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	if mc.cfg.TLS != nil {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		clientFlags |= clientSSL
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	if mc.cfg.MultiStatements {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		clientFlags |= clientMultiStatements
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// encode length of the auth plugin data</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	var authRespLEIBuf [9]byte
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	authRespLen := len(authResp)
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	authRespLEI := appendLengthEncodedInteger(authRespLEIBuf[:0], uint64(authRespLen))
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if len(authRespLEI) &gt; 1 {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		<span class="comment">// if the length can not be written in 1 byte, it must be written as a</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">// length encoded integer</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		clientFlags |= clientPluginAuthLenEncClientData
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	pktLen := 4 + 4 + 1 + 23 + len(mc.cfg.User) + 1 + len(authRespLEI) + len(authResp) + 21 + 1
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	<span class="comment">// To specify a db name</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	if n := len(mc.cfg.DBName); n &gt; 0 {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		clientFlags |= clientConnectWithDB
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		pktLen += n + 1
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	}
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// encode length of the connection attributes</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	var connAttrsLEIBuf [9]byte
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	connAttrsLen := len(mc.connector.encodedAttributes)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	connAttrsLEI := appendLengthEncodedInteger(connAttrsLEIBuf[:0], uint64(connAttrsLen))
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	pktLen += len(connAttrsLEI) + len(mc.connector.encodedAttributes)
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	<span class="comment">// Calculate packet length and get buffer with that size</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	data, err := mc.buf.takeBuffer(pktLen + 4)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	if err != nil {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		mc.log(err)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	<span class="comment">// ClientFlags [32 bit]</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	data[4] = byte(clientFlags)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	data[5] = byte(clientFlags &gt;&gt; 8)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	data[6] = byte(clientFlags &gt;&gt; 16)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	data[7] = byte(clientFlags &gt;&gt; 24)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// MaxPacketSize [32 bit] (none)</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	data[8] = 0x00
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	data[9] = 0x00
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	data[10] = 0x00
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	data[11] = 0x00
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	<span class="comment">// Collation ID [1 byte]</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	cname := mc.cfg.Collation
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	if cname == &#34;&#34; {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		cname = defaultCollation
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	var found bool
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	data[12], found = collations[cname]
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if !found {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		<span class="comment">// Note possibility for false negatives:</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		<span class="comment">// could be triggered  although the collation is valid if the</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		<span class="comment">// collations map does not contain entries the server supports.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;unknown collation: %q&#34;, cname)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>	<span class="comment">// Filler [23 bytes] (all 0x00)</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	pos := 13
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>	for ; pos &lt; 13+23; pos++ {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		data[pos] = 0
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>	}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	<span class="comment">// SSL Connection Request Packet</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	<span class="comment">// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	if mc.cfg.TLS != nil {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		<span class="comment">// Send TLS / SSL request packet</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		if err := mc.writePacket(data[:(4+4+1+23)+4]); err != nil {
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>			return err
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		<span class="comment">// Switch to TLS</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		tlsConn := tls.Client(mc.netConn, mc.cfg.TLS)
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		if err := tlsConn.Handshake(); err != nil {
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>			return err
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		}
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		mc.netConn = tlsConn
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		mc.buf.nc = tlsConn
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	<span class="comment">// User [null terminated string]</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	if len(mc.cfg.User) &gt; 0 {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		pos += copy(data[pos:], mc.cfg.User)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	}
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	data[pos] = 0x00
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	pos++
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	<span class="comment">// Auth Data [length encoded integer]</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	pos += copy(data[pos:], authRespLEI)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	pos += copy(data[pos:], authResp)
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// Databasename [null terminated string]</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	if len(mc.cfg.DBName) &gt; 0 {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		pos += copy(data[pos:], mc.cfg.DBName)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		data[pos] = 0x00
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		pos++
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	pos += copy(data[pos:], plugin)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	data[pos] = 0x00
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	pos++
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	<span class="comment">// Connection Attributes</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	pos += copy(data[pos:], connAttrsLEI)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	pos += copy(data[pos:], []byte(mc.connector.encodedAttributes))
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	<span class="comment">// Send Auth packet</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	return mc.writePacket(data[:pos])
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>}
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::AuthSwitchResponse</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>func (mc *mysqlConn) writeAuthSwitchPacket(authData []byte) error {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	pktLen := 4 + len(authData)
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	data, err := mc.buf.takeSmallBuffer(pktLen)
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	if err != nil {
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		mc.log(err)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	<span class="comment">// Add the auth data [EOF]</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	copy(data[4:], authData)
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	return mc.writePacket(data)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">/******************************************************************************
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>*                             Command Packets                                 *
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>******************************************************************************/</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>func (mc *mysqlConn) writeCommandPacket(command byte) error {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	<span class="comment">// Reset Packet Sequence</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	mc.sequence = 0
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	data, err := mc.buf.takeSmallBuffer(4 + 1)
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	if err != nil {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		mc.log(err)
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	}
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// Add command byte</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	data[4] = command
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	<span class="comment">// Send CMD packet</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	return mc.writePacket(data)
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>}
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>func (mc *mysqlConn) writeCommandPacketStr(command byte, arg string) error {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	<span class="comment">// Reset Packet Sequence</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	mc.sequence = 0
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	pktLen := 1 + len(arg)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	data, err := mc.buf.takeBuffer(pktLen + 4)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	if err != nil {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		mc.log(err)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// Add command byte</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	data[4] = command
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// Add arg</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	copy(data[5:], arg)
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	<span class="comment">// Send CMD packet</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	return mc.writePacket(data)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>}
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>func (mc *mysqlConn) writeCommandPacketUint32(command byte, arg uint32) error {
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	<span class="comment">// Reset Packet Sequence</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	mc.sequence = 0
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	data, err := mc.buf.takeSmallBuffer(4 + 1 + 4)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	if err != nil {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		mc.log(err)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	<span class="comment">// Add command byte</span>
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	data[4] = command
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	<span class="comment">// Add arg [32 bit]</span>
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>	data[5] = byte(arg)
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>	data[6] = byte(arg &gt;&gt; 8)
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>	data[7] = byte(arg &gt;&gt; 16)
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	data[8] = byte(arg &gt;&gt; 24)
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>	<span class="comment">// Send CMD packet</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	return mc.writePacket(data)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>}
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">/******************************************************************************
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>*                              Result Packets                                 *
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>******************************************************************************/</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>func (mc *mysqlConn) readAuthResult() ([]byte, string, error) {
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>	data, err := mc.readPacket()
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>	if err != nil {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>		return nil, &#34;&#34;, err
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	}
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	<span class="comment">// packet indicator</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>	switch data[0] {
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	case iOK:
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>		<span class="comment">// resultUnchanged, since auth happens before any queries or</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>		<span class="comment">// commands have been executed.</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		return nil, &#34;&#34;, mc.resultUnchanged().handleOkPacket(data)
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	case iAuthMoreData:
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>		return data[1:], &#34;&#34;, err
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	case iEOF:
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		if len(data) == 1 {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			<span class="comment">// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::OldAuthSwitchRequest</span>
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			return nil, &#34;mysql_old_password&#34;, nil
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		}
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		pluginEndIndex := bytes.IndexByte(data, 0x00)
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		if pluginEndIndex &lt; 0 {
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			return nil, &#34;&#34;, ErrMalformPkt
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>		}
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		plugin := string(data[1:pluginEndIndex])
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>		authData := data[pluginEndIndex+1:]
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>		return authData, plugin, nil
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	default: <span class="comment">// Error otherwise</span>
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		return nil, &#34;&#34;, mc.handleErrorPacket(data)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span><span class="comment">// Returns error if Packet is not a &#39;Result OK&#39;-Packet</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>func (mc *okHandler) readResultOK() error {
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	data, err := mc.conn().readPacket()
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	if err != nil {
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		return err
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	}
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>	if data[0] == iOK {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		return mc.handleOkPacket(data)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	return mc.conn().handleErrorPacket(data)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// Result Set Header Packet</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::Resultset</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>func (mc *okHandler) readResultSetHeaderPacket() (int, error) {
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">// handleOkPacket replaces both values; other cases leave the values unchanged.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	mc.result.affectedRows = append(mc.result.affectedRows, 0)
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	mc.result.insertIds = append(mc.result.insertIds, 0)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	data, err := mc.conn().readPacket()
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if err == nil {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		switch data[0] {
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		case iOK:
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			return 0, mc.handleOkPacket(data)
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		case iERR:
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>			return 0, mc.conn().handleErrorPacket(data)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		case iLocalInFile:
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			return 0, mc.handleInFileRequest(string(data[1:]))
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>		}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		<span class="comment">// column count</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		num, _, _ := readLengthEncodedInteger(data)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		<span class="comment">// ignore remaining data in the packet. see #1478.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>		return int(num), nil
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>	return 0, err
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>}
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// Error Packet</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/generic-response-packets.html#packet-ERR_Packet</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>func (mc *mysqlConn) handleErrorPacket(data []byte) error {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	if data[0] != iERR {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		return ErrMalformPkt
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	<span class="comment">// 0xff [1 byte]</span>
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	<span class="comment">// Error Number [16 bit uint]</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	errno := binary.LittleEndian.Uint16(data[1:3])
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	<span class="comment">// 1792: ER_CANT_EXECUTE_IN_READ_ONLY_TRANSACTION</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	<span class="comment">// 1290: ER_OPTION_PREVENTS_STATEMENT (returned by Aurora during failover)</span>
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	if (errno == 1792 || errno == 1290) &amp;&amp; mc.cfg.RejectReadOnly {
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		<span class="comment">// Oops; we are connected to a read-only connection, and won&#39;t be able</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		<span class="comment">// to issue any write statements. Since RejectReadOnly is configured,</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		<span class="comment">// we throw away this connection hoping this one would have write</span>
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>		<span class="comment">// permission. This is specifically for a possible race condition</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		<span class="comment">// during failover (e.g. on AWS Aurora). See README.md for more.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>		<span class="comment">// We explicitly close the connection before returning</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		<span class="comment">// driver.ErrBadConn to ensure that `database/sql` purges this</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		<span class="comment">// connection and initiates a new one for next statement next time.</span>
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>		mc.Close()
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		return driver.ErrBadConn
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	}
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	me := &amp;MySQLError{Number: errno}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	pos := 3
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	<span class="comment">// SQL State [optional: # + 5bytes string]</span>
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	if data[3] == 0x23 {
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		copy(me.SQLState[:], data[4:4+5])
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>		pos = 9
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	<span class="comment">// Error Message [string]</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	me.Message = string(data[pos:])
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	return me
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>func readStatus(b []byte) statusFlag {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	return statusFlag(b[0]) | statusFlag(b[1])&lt;&lt;8
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span><span class="comment">// Returns an instance of okHandler for codepaths where mysqlConn.result doesn&#39;t</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span><span class="comment">// need to be cleared first (e.g. during authentication, or while additional</span>
<span id="L602" class="ln">   602&nbsp;&nbsp;</span><span class="comment">// resultsets are being fetched.)</span>
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>func (mc *mysqlConn) resultUnchanged() *okHandler {
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	return (*okHandler)(mc)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>}
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span><span class="comment">// okHandler represents the state of the connection when mysqlConn.result has</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span><span class="comment">// been prepared for processing of OK packets.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span><span class="comment">// To correctly populate mysqlConn.result (updated by handleOkPacket()), all</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span><span class="comment">// callpaths must either:</span>
<span id="L612" class="ln">   612&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L613" class="ln">   613&nbsp;&nbsp;</span><span class="comment">// 1. first clear it using clearResult(), or</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span><span class="comment">// 2. confirm that they don&#39;t need to (by calling resultUnchanged()).</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L616" class="ln">   616&nbsp;&nbsp;</span><span class="comment">// Both return an instance of type *okHandler.</span>
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>type okHandler mysqlConn
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span><span class="comment">// Exposes the underlying type&#39;s methods.</span>
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>func (mc *okHandler) conn() *mysqlConn {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	return (*mysqlConn)(mc)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span><span class="comment">// clearResult clears the connection&#39;s stored affectedRows and insertIds</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// fields.</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span><span class="comment">// It returns a handler that can process OK responses.</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>func (mc *mysqlConn) clearResult() *okHandler {
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>	mc.result = mysqlResult{}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>	return (*okHandler)(mc)
<span id="L631" class="ln">   631&nbsp;&nbsp;</span>}
<span id="L632" class="ln">   632&nbsp;&nbsp;</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// Ok Packet</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/generic-response-packets.html#packet-OK_Packet</span>
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>func (mc *okHandler) handleOkPacket(data []byte) error {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	var n, m int
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	var affectedRows, insertId uint64
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	<span class="comment">// 0x00 [1 byte]</span>
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// Affected rows [Length Coded Binary]</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	affectedRows, _, n = readLengthEncodedInteger(data[1:])
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>	<span class="comment">// Insert id [Length Coded Binary]</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	insertId, _, m = readLengthEncodedInteger(data[1+n:])
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>	<span class="comment">// Update for the current statement result (only used by</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>	<span class="comment">// readResultSetHeaderPacket).</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>	if len(mc.result.affectedRows) &gt; 0 {
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		mc.result.affectedRows[len(mc.result.affectedRows)-1] = int64(affectedRows)
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>	}
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>	if len(mc.result.insertIds) &gt; 0 {
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		mc.result.insertIds[len(mc.result.insertIds)-1] = int64(insertId)
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>	}
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	<span class="comment">// server_status [2 bytes]</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>	mc.status = readStatus(data[1+n+m : 1+n+m+2])
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	if mc.status&amp;statusMoreResultsExists != 0 {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>		return nil
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	}
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	<span class="comment">// warning count [2 bytes]</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>	return nil
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>}
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span><span class="comment">// Read Packets as Field Packets until EOF-Packet or an Error appears</span>
<span id="L668" class="ln">   668&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-Protocol::ColumnDefinition41</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>func (mc *mysqlConn) readColumns(count int) ([]mysqlField, error) {
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>	columns := make([]mysqlField, count)
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>	for i := 0; ; i++ {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>		data, err := mc.readPacket()
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		if err != nil {
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>			return nil, err
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>		}
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		<span class="comment">// EOF Packet</span>
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		if data[0] == iEOF &amp;&amp; (len(data) == 5 || len(data) == 1) {
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>			if i == count {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>				return columns, nil
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>			}
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;column count mismatch n:%d len:%d&#34;, count, len(columns))
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		}
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		<span class="comment">// Catalog</span>
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>		pos, err := skipLengthEncodedString(data)
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		if err != nil {
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>			return nil, err
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		}
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		<span class="comment">// Database [len coded string]</span>
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>		n, err := skipLengthEncodedString(data[pos:])
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>		if err != nil {
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			return nil, err
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>		}
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>		pos += n
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>		<span class="comment">// Table [len coded string]</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		if mc.cfg.ColumnsWithAlias {
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>			tableName, _, n, err := readLengthEncodedString(data[pos:])
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>			if err != nil {
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>				return nil, err
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>			}
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			pos += n
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>			columns[i].tableName = string(tableName)
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		} else {
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>			n, err = skipLengthEncodedString(data[pos:])
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>			if err != nil {
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>				return nil, err
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>			}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>			pos += n
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>		}
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>		<span class="comment">// Original table [len coded string]</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>		n, err = skipLengthEncodedString(data[pos:])
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>		if err != nil {
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>			return nil, err
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>		}
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		pos += n
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		<span class="comment">// Name [len coded string]</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>		name, _, n, err := readLengthEncodedString(data[pos:])
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>		if err != nil {
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>			return nil, err
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>		}
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>		columns[i].name = string(name)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		pos += n
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		<span class="comment">// Original name [len coded string]</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		n, err = skipLengthEncodedString(data[pos:])
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>		if err != nil {
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>			return nil, err
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>		}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>		pos += n
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>		<span class="comment">// Filler [uint8]</span>
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>		pos++
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>		<span class="comment">// Charset [charset, collation uint8]</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		columns[i].charSet = data[pos]
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		pos += 2
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		<span class="comment">// Length [uint32]</span>
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		columns[i].length = binary.LittleEndian.Uint32(data[pos : pos+4])
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>		pos += 4
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>		<span class="comment">// Field type [uint8]</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>		columns[i].fieldType = fieldType(data[pos])
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>		pos++
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>		<span class="comment">// Flags [uint16]</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>		columns[i].flags = fieldFlag(binary.LittleEndian.Uint16(data[pos : pos+2]))
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		pos += 2
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		<span class="comment">// Decimals [uint8]</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>		columns[i].decimals = data[pos]
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>		<span class="comment">//pos++</span>
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>		<span class="comment">// Default value [len coded binary]</span>
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		<span class="comment">//if pos &lt; len(data) {</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		<span class="comment">//	defaultVal, _, err = bytesToLengthCodedBinary(data[pos:])</span>
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>		<span class="comment">//}</span>
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	}
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>}
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>
<span id="L767" class="ln">   767&nbsp;&nbsp;</span><span class="comment">// Read Packets as Field Packets until EOF-Packet or an Error appears</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-query-response.html#packet-ProtocolText::ResultsetRow</span>
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>func (rows *textRows) readRow(dest []driver.Value) error {
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	mc := rows.mc
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	if rows.rs.done {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>		return io.EOF
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	}
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	data, err := mc.readPacket()
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>	if err != nil {
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		return err
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	}
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">// EOF Packet</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	if data[0] == iEOF &amp;&amp; len(data) == 5 {
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>		<span class="comment">// server_status [2 bytes]</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>		rows.mc.status = readStatus(data[3:])
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>		rows.rs.done = true
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>		if !rows.HasNextResultSet() {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>			rows.mc = nil
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		}
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>		return io.EOF
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>	}
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	if data[0] == iERR {
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>		rows.mc = nil
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>		return mc.handleErrorPacket(data)
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	}
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">// RowSet Packet</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	var (
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		n      int
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>		isNull bool
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		pos    int = 0
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	)
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>	for i := range dest {
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		<span class="comment">// Read bytes and convert to string</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		var buf []byte
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		buf, isNull, n, err = readLengthEncodedString(data[pos:])
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		pos += n
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>		if err != nil {
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>			return err
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>		if isNull {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>			dest[i] = nil
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>			continue
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>		}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		switch rows.rs.columns[i].fieldType {
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		case fieldTypeTimestamp,
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>			fieldTypeDateTime,
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>			fieldTypeDate,
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>			fieldTypeNewDate:
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>			if mc.parseTime {
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>				dest[i], err = parseDateTime(buf, mc.cfg.Loc)
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>			} else {
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>				dest[i] = buf
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>			}
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		case fieldTypeTiny, fieldTypeShort, fieldTypeInt24, fieldTypeYear, fieldTypeLong:
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>			dest[i], err = strconv.ParseInt(string(buf), 10, 64)
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		case fieldTypeLongLong:
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>			if rows.rs.columns[i].flags&amp;flagUnsigned != 0 {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>				dest[i], err = strconv.ParseUint(string(buf), 10, 64)
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>			} else {
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>				dest[i], err = strconv.ParseInt(string(buf), 10, 64)
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>			}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		case fieldTypeFloat:
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			var d float64
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>			d, err = strconv.ParseFloat(string(buf), 32)
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			dest[i] = float32(d)
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>		case fieldTypeDouble:
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			dest[i], err = strconv.ParseFloat(string(buf), 64)
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>		default:
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			dest[i] = buf
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		if err != nil {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			return err
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>	}
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>	return nil
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span><span class="comment">// Reads Packets until EOF-Packet or an Error appears. Returns count of Packets read</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>func (mc *mysqlConn) readUntilEOF() error {
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	for {
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>		data, err := mc.readPacket()
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		if err != nil {
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>			return err
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		}
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		switch data[0] {
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		case iERR:
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			return mc.handleErrorPacket(data)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>		case iEOF:
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			if len(data) == 5 {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>				mc.status = readStatus(data[3:])
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>			}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			return nil
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	}
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>}
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>
<span id="L878" class="ln">   878&nbsp;&nbsp;</span><span class="comment">/******************************************************************************
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>*                           Prepared Statements                               *
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>******************************************************************************/</span>
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>
<span id="L882" class="ln">   882&nbsp;&nbsp;</span><span class="comment">// Prepare Result Packets</span>
<span id="L883" class="ln">   883&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-stmt-prepare-response.html</span>
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>func (stmt *mysqlStmt) readPrepareResultPacket() (uint16, error) {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>	data, err := stmt.mc.readPacket()
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	if err == nil {
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>		<span class="comment">// packet indicator [1 byte]</span>
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>		if data[0] != iOK {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>			return 0, stmt.mc.handleErrorPacket(data)
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>		}
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		<span class="comment">// statement id [4 bytes]</span>
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>		stmt.id = binary.LittleEndian.Uint32(data[1:5])
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>		<span class="comment">// Column count [16 bit uint]</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>		columnCount := binary.LittleEndian.Uint16(data[5:7])
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>		<span class="comment">// Param count [16 bit uint]</span>
<span id="L899" class="ln">   899&nbsp;&nbsp;</span>		stmt.paramCount = int(binary.LittleEndian.Uint16(data[7:9]))
<span id="L900" class="ln">   900&nbsp;&nbsp;</span>
<span id="L901" class="ln">   901&nbsp;&nbsp;</span>		<span class="comment">// Reserved [8 bit]</span>
<span id="L902" class="ln">   902&nbsp;&nbsp;</span>
<span id="L903" class="ln">   903&nbsp;&nbsp;</span>		<span class="comment">// Warning count [16 bit uint]</span>
<span id="L904" class="ln">   904&nbsp;&nbsp;</span>
<span id="L905" class="ln">   905&nbsp;&nbsp;</span>		return columnCount, nil
<span id="L906" class="ln">   906&nbsp;&nbsp;</span>	}
<span id="L907" class="ln">   907&nbsp;&nbsp;</span>	return 0, err
<span id="L908" class="ln">   908&nbsp;&nbsp;</span>}
<span id="L909" class="ln">   909&nbsp;&nbsp;</span>
<span id="L910" class="ln">   910&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-stmt-send-long-data.html</span>
<span id="L911" class="ln">   911&nbsp;&nbsp;</span>func (stmt *mysqlStmt) writeCommandLongData(paramID int, arg []byte) error {
<span id="L912" class="ln">   912&nbsp;&nbsp;</span>	maxLen := stmt.mc.maxAllowedPacket - 1
<span id="L913" class="ln">   913&nbsp;&nbsp;</span>	pktLen := maxLen
<span id="L914" class="ln">   914&nbsp;&nbsp;</span>
<span id="L915" class="ln">   915&nbsp;&nbsp;</span>	<span class="comment">// After the header (bytes 0-3) follows before the data:</span>
<span id="L916" class="ln">   916&nbsp;&nbsp;</span>	<span class="comment">// 1 byte command</span>
<span id="L917" class="ln">   917&nbsp;&nbsp;</span>	<span class="comment">// 4 bytes stmtID</span>
<span id="L918" class="ln">   918&nbsp;&nbsp;</span>	<span class="comment">// 2 bytes paramID</span>
<span id="L919" class="ln">   919&nbsp;&nbsp;</span>	const dataOffset = 1 + 4 + 2
<span id="L920" class="ln">   920&nbsp;&nbsp;</span>
<span id="L921" class="ln">   921&nbsp;&nbsp;</span>	<span class="comment">// Cannot use the write buffer since</span>
<span id="L922" class="ln">   922&nbsp;&nbsp;</span>	<span class="comment">// a) the buffer is too small</span>
<span id="L923" class="ln">   923&nbsp;&nbsp;</span>	<span class="comment">// b) it is in use</span>
<span id="L924" class="ln">   924&nbsp;&nbsp;</span>	data := make([]byte, 4+1+4+2+len(arg))
<span id="L925" class="ln">   925&nbsp;&nbsp;</span>
<span id="L926" class="ln">   926&nbsp;&nbsp;</span>	copy(data[4+dataOffset:], arg)
<span id="L927" class="ln">   927&nbsp;&nbsp;</span>
<span id="L928" class="ln">   928&nbsp;&nbsp;</span>	for argLen := len(arg); argLen &gt; 0; argLen -= pktLen - dataOffset {
<span id="L929" class="ln">   929&nbsp;&nbsp;</span>		if dataOffset+argLen &lt; maxLen {
<span id="L930" class="ln">   930&nbsp;&nbsp;</span>			pktLen = dataOffset + argLen
<span id="L931" class="ln">   931&nbsp;&nbsp;</span>		}
<span id="L932" class="ln">   932&nbsp;&nbsp;</span>
<span id="L933" class="ln">   933&nbsp;&nbsp;</span>		stmt.mc.sequence = 0
<span id="L934" class="ln">   934&nbsp;&nbsp;</span>		<span class="comment">// Add command byte [1 byte]</span>
<span id="L935" class="ln">   935&nbsp;&nbsp;</span>		data[4] = comStmtSendLongData
<span id="L936" class="ln">   936&nbsp;&nbsp;</span>
<span id="L937" class="ln">   937&nbsp;&nbsp;</span>		<span class="comment">// Add stmtID [32 bit]</span>
<span id="L938" class="ln">   938&nbsp;&nbsp;</span>		data[5] = byte(stmt.id)
<span id="L939" class="ln">   939&nbsp;&nbsp;</span>		data[6] = byte(stmt.id &gt;&gt; 8)
<span id="L940" class="ln">   940&nbsp;&nbsp;</span>		data[7] = byte(stmt.id &gt;&gt; 16)
<span id="L941" class="ln">   941&nbsp;&nbsp;</span>		data[8] = byte(stmt.id &gt;&gt; 24)
<span id="L942" class="ln">   942&nbsp;&nbsp;</span>
<span id="L943" class="ln">   943&nbsp;&nbsp;</span>		<span class="comment">// Add paramID [16 bit]</span>
<span id="L944" class="ln">   944&nbsp;&nbsp;</span>		data[9] = byte(paramID)
<span id="L945" class="ln">   945&nbsp;&nbsp;</span>		data[10] = byte(paramID &gt;&gt; 8)
<span id="L946" class="ln">   946&nbsp;&nbsp;</span>
<span id="L947" class="ln">   947&nbsp;&nbsp;</span>		<span class="comment">// Send CMD packet</span>
<span id="L948" class="ln">   948&nbsp;&nbsp;</span>		err := stmt.mc.writePacket(data[:4+pktLen])
<span id="L949" class="ln">   949&nbsp;&nbsp;</span>		if err == nil {
<span id="L950" class="ln">   950&nbsp;&nbsp;</span>			data = data[pktLen-dataOffset:]
<span id="L951" class="ln">   951&nbsp;&nbsp;</span>			continue
<span id="L952" class="ln">   952&nbsp;&nbsp;</span>		}
<span id="L953" class="ln">   953&nbsp;&nbsp;</span>		return err
<span id="L954" class="ln">   954&nbsp;&nbsp;</span>
<span id="L955" class="ln">   955&nbsp;&nbsp;</span>	}
<span id="L956" class="ln">   956&nbsp;&nbsp;</span>
<span id="L957" class="ln">   957&nbsp;&nbsp;</span>	<span class="comment">// Reset Packet Sequence</span>
<span id="L958" class="ln">   958&nbsp;&nbsp;</span>	stmt.mc.sequence = 0
<span id="L959" class="ln">   959&nbsp;&nbsp;</span>	return nil
<span id="L960" class="ln">   960&nbsp;&nbsp;</span>}
<span id="L961" class="ln">   961&nbsp;&nbsp;</span>
<span id="L962" class="ln">   962&nbsp;&nbsp;</span><span class="comment">// Execute Prepared Statement</span>
<span id="L963" class="ln">   963&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/com-stmt-execute.html</span>
<span id="L964" class="ln">   964&nbsp;&nbsp;</span>func (stmt *mysqlStmt) writeExecutePacket(args []driver.Value) error {
<span id="L965" class="ln">   965&nbsp;&nbsp;</span>	if len(args) != stmt.paramCount {
<span id="L966" class="ln">   966&nbsp;&nbsp;</span>		return fmt.Errorf(
<span id="L967" class="ln">   967&nbsp;&nbsp;</span>			&#34;argument count mismatch (got: %d; has: %d)&#34;,
<span id="L968" class="ln">   968&nbsp;&nbsp;</span>			len(args),
<span id="L969" class="ln">   969&nbsp;&nbsp;</span>			stmt.paramCount,
<span id="L970" class="ln">   970&nbsp;&nbsp;</span>		)
<span id="L971" class="ln">   971&nbsp;&nbsp;</span>	}
<span id="L972" class="ln">   972&nbsp;&nbsp;</span>
<span id="L973" class="ln">   973&nbsp;&nbsp;</span>	const minPktLen = 4 + 1 + 4 + 1 + 4
<span id="L974" class="ln">   974&nbsp;&nbsp;</span>	mc := stmt.mc
<span id="L975" class="ln">   975&nbsp;&nbsp;</span>
<span id="L976" class="ln">   976&nbsp;&nbsp;</span>	<span class="comment">// Determine threshold dynamically to avoid packet size shortage.</span>
<span id="L977" class="ln">   977&nbsp;&nbsp;</span>	longDataSize := mc.maxAllowedPacket / (stmt.paramCount + 1)
<span id="L978" class="ln">   978&nbsp;&nbsp;</span>	if longDataSize &lt; 64 {
<span id="L979" class="ln">   979&nbsp;&nbsp;</span>		longDataSize = 64
<span id="L980" class="ln">   980&nbsp;&nbsp;</span>	}
<span id="L981" class="ln">   981&nbsp;&nbsp;</span>
<span id="L982" class="ln">   982&nbsp;&nbsp;</span>	<span class="comment">// Reset packet-sequence</span>
<span id="L983" class="ln">   983&nbsp;&nbsp;</span>	mc.sequence = 0
<span id="L984" class="ln">   984&nbsp;&nbsp;</span>
<span id="L985" class="ln">   985&nbsp;&nbsp;</span>	var data []byte
<span id="L986" class="ln">   986&nbsp;&nbsp;</span>	var err error
<span id="L987" class="ln">   987&nbsp;&nbsp;</span>
<span id="L988" class="ln">   988&nbsp;&nbsp;</span>	if len(args) == 0 {
<span id="L989" class="ln">   989&nbsp;&nbsp;</span>		data, err = mc.buf.takeBuffer(minPktLen)
<span id="L990" class="ln">   990&nbsp;&nbsp;</span>	} else {
<span id="L991" class="ln">   991&nbsp;&nbsp;</span>		data, err = mc.buf.takeCompleteBuffer()
<span id="L992" class="ln">   992&nbsp;&nbsp;</span>		<span class="comment">// In this case the len(data) == cap(data) which is used to optimise the flow below.</span>
<span id="L993" class="ln">   993&nbsp;&nbsp;</span>	}
<span id="L994" class="ln">   994&nbsp;&nbsp;</span>	if err != nil {
<span id="L995" class="ln">   995&nbsp;&nbsp;</span>		<span class="comment">// cannot take the buffer. Something must be wrong with the connection</span>
<span id="L996" class="ln">   996&nbsp;&nbsp;</span>		mc.log(err)
<span id="L997" class="ln">   997&nbsp;&nbsp;</span>		return errBadConnNoWrite
<span id="L998" class="ln">   998&nbsp;&nbsp;</span>	}
<span id="L999" class="ln">   999&nbsp;&nbsp;</span>
<span id="L1000" class="ln">  1000&nbsp;&nbsp;</span>	<span class="comment">// command [1 byte]</span>
<span id="L1001" class="ln">  1001&nbsp;&nbsp;</span>	data[4] = comStmtExecute
<span id="L1002" class="ln">  1002&nbsp;&nbsp;</span>
<span id="L1003" class="ln">  1003&nbsp;&nbsp;</span>	<span class="comment">// statement_id [4 bytes]</span>
<span id="L1004" class="ln">  1004&nbsp;&nbsp;</span>	data[5] = byte(stmt.id)
<span id="L1005" class="ln">  1005&nbsp;&nbsp;</span>	data[6] = byte(stmt.id &gt;&gt; 8)
<span id="L1006" class="ln">  1006&nbsp;&nbsp;</span>	data[7] = byte(stmt.id &gt;&gt; 16)
<span id="L1007" class="ln">  1007&nbsp;&nbsp;</span>	data[8] = byte(stmt.id &gt;&gt; 24)
<span id="L1008" class="ln">  1008&nbsp;&nbsp;</span>
<span id="L1009" class="ln">  1009&nbsp;&nbsp;</span>	<span class="comment">// flags (0: CURSOR_TYPE_NO_CURSOR) [1 byte]</span>
<span id="L1010" class="ln">  1010&nbsp;&nbsp;</span>	data[9] = 0x00
<span id="L1011" class="ln">  1011&nbsp;&nbsp;</span>
<span id="L1012" class="ln">  1012&nbsp;&nbsp;</span>	<span class="comment">// iteration_count (uint32(1)) [4 bytes]</span>
<span id="L1013" class="ln">  1013&nbsp;&nbsp;</span>	data[10] = 0x01
<span id="L1014" class="ln">  1014&nbsp;&nbsp;</span>	data[11] = 0x00
<span id="L1015" class="ln">  1015&nbsp;&nbsp;</span>	data[12] = 0x00
<span id="L1016" class="ln">  1016&nbsp;&nbsp;</span>	data[13] = 0x00
<span id="L1017" class="ln">  1017&nbsp;&nbsp;</span>
<span id="L1018" class="ln">  1018&nbsp;&nbsp;</span>	if len(args) &gt; 0 {
<span id="L1019" class="ln">  1019&nbsp;&nbsp;</span>		pos := minPktLen
<span id="L1020" class="ln">  1020&nbsp;&nbsp;</span>
<span id="L1021" class="ln">  1021&nbsp;&nbsp;</span>		var nullMask []byte
<span id="L1022" class="ln">  1022&nbsp;&nbsp;</span>		if maskLen, typesLen := (len(args)+7)/8, 1+2*len(args); pos+maskLen+typesLen &gt;= cap(data) {
<span id="L1023" class="ln">  1023&nbsp;&nbsp;</span>			<span class="comment">// buffer has to be extended but we don&#39;t know by how much so</span>
<span id="L1024" class="ln">  1024&nbsp;&nbsp;</span>			<span class="comment">// we depend on append after all data with known sizes fit.</span>
<span id="L1025" class="ln">  1025&nbsp;&nbsp;</span>			<span class="comment">// We stop at that because we deal with a lot of columns here</span>
<span id="L1026" class="ln">  1026&nbsp;&nbsp;</span>			<span class="comment">// which makes the required allocation size hard to guess.</span>
<span id="L1027" class="ln">  1027&nbsp;&nbsp;</span>			tmp := make([]byte, pos+maskLen+typesLen)
<span id="L1028" class="ln">  1028&nbsp;&nbsp;</span>			copy(tmp[:pos], data[:pos])
<span id="L1029" class="ln">  1029&nbsp;&nbsp;</span>			data = tmp
<span id="L1030" class="ln">  1030&nbsp;&nbsp;</span>			nullMask = data[pos : pos+maskLen]
<span id="L1031" class="ln">  1031&nbsp;&nbsp;</span>			<span class="comment">// No need to clean nullMask as make ensures that.</span>
<span id="L1032" class="ln">  1032&nbsp;&nbsp;</span>			pos += maskLen
<span id="L1033" class="ln">  1033&nbsp;&nbsp;</span>		} else {
<span id="L1034" class="ln">  1034&nbsp;&nbsp;</span>			nullMask = data[pos : pos+maskLen]
<span id="L1035" class="ln">  1035&nbsp;&nbsp;</span>			for i := range nullMask {
<span id="L1036" class="ln">  1036&nbsp;&nbsp;</span>				nullMask[i] = 0
<span id="L1037" class="ln">  1037&nbsp;&nbsp;</span>			}
<span id="L1038" class="ln">  1038&nbsp;&nbsp;</span>			pos += maskLen
<span id="L1039" class="ln">  1039&nbsp;&nbsp;</span>		}
<span id="L1040" class="ln">  1040&nbsp;&nbsp;</span>
<span id="L1041" class="ln">  1041&nbsp;&nbsp;</span>		<span class="comment">// newParameterBoundFlag 1 [1 byte]</span>
<span id="L1042" class="ln">  1042&nbsp;&nbsp;</span>		data[pos] = 0x01
<span id="L1043" class="ln">  1043&nbsp;&nbsp;</span>		pos++
<span id="L1044" class="ln">  1044&nbsp;&nbsp;</span>
<span id="L1045" class="ln">  1045&nbsp;&nbsp;</span>		<span class="comment">// type of each parameter [len(args)*2 bytes]</span>
<span id="L1046" class="ln">  1046&nbsp;&nbsp;</span>		paramTypes := data[pos:]
<span id="L1047" class="ln">  1047&nbsp;&nbsp;</span>		pos += len(args) * 2
<span id="L1048" class="ln">  1048&nbsp;&nbsp;</span>
<span id="L1049" class="ln">  1049&nbsp;&nbsp;</span>		<span class="comment">// value of each parameter [n bytes]</span>
<span id="L1050" class="ln">  1050&nbsp;&nbsp;</span>		paramValues := data[pos:pos]
<span id="L1051" class="ln">  1051&nbsp;&nbsp;</span>		valuesCap := cap(paramValues)
<span id="L1052" class="ln">  1052&nbsp;&nbsp;</span>
<span id="L1053" class="ln">  1053&nbsp;&nbsp;</span>		for i, arg := range args {
<span id="L1054" class="ln">  1054&nbsp;&nbsp;</span>			<span class="comment">// build NULL-bitmap</span>
<span id="L1055" class="ln">  1055&nbsp;&nbsp;</span>			if arg == nil {
<span id="L1056" class="ln">  1056&nbsp;&nbsp;</span>				nullMask[i/8] |= 1 &lt;&lt; (uint(i) &amp; 7)
<span id="L1057" class="ln">  1057&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeNULL)
<span id="L1058" class="ln">  1058&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1059" class="ln">  1059&nbsp;&nbsp;</span>				continue
<span id="L1060" class="ln">  1060&nbsp;&nbsp;</span>			}
<span id="L1061" class="ln">  1061&nbsp;&nbsp;</span>
<span id="L1062" class="ln">  1062&nbsp;&nbsp;</span>			if v, ok := arg.(json.RawMessage); ok {
<span id="L1063" class="ln">  1063&nbsp;&nbsp;</span>				arg = []byte(v)
<span id="L1064" class="ln">  1064&nbsp;&nbsp;</span>			}
<span id="L1065" class="ln">  1065&nbsp;&nbsp;</span>			<span class="comment">// cache types and values</span>
<span id="L1066" class="ln">  1066&nbsp;&nbsp;</span>			switch v := arg.(type) {
<span id="L1067" class="ln">  1067&nbsp;&nbsp;</span>			case int64:
<span id="L1068" class="ln">  1068&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeLongLong)
<span id="L1069" class="ln">  1069&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1070" class="ln">  1070&nbsp;&nbsp;</span>
<span id="L1071" class="ln">  1071&nbsp;&nbsp;</span>				if cap(paramValues)-len(paramValues)-8 &gt;= 0 {
<span id="L1072" class="ln">  1072&nbsp;&nbsp;</span>					paramValues = paramValues[:len(paramValues)+8]
<span id="L1073" class="ln">  1073&nbsp;&nbsp;</span>					binary.LittleEndian.PutUint64(
<span id="L1074" class="ln">  1074&nbsp;&nbsp;</span>						paramValues[len(paramValues)-8:],
<span id="L1075" class="ln">  1075&nbsp;&nbsp;</span>						uint64(v),
<span id="L1076" class="ln">  1076&nbsp;&nbsp;</span>					)
<span id="L1077" class="ln">  1077&nbsp;&nbsp;</span>				} else {
<span id="L1078" class="ln">  1078&nbsp;&nbsp;</span>					paramValues = append(paramValues,
<span id="L1079" class="ln">  1079&nbsp;&nbsp;</span>						uint64ToBytes(uint64(v))...,
<span id="L1080" class="ln">  1080&nbsp;&nbsp;</span>					)
<span id="L1081" class="ln">  1081&nbsp;&nbsp;</span>				}
<span id="L1082" class="ln">  1082&nbsp;&nbsp;</span>
<span id="L1083" class="ln">  1083&nbsp;&nbsp;</span>			case uint64:
<span id="L1084" class="ln">  1084&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeLongLong)
<span id="L1085" class="ln">  1085&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x80 <span class="comment">// type is unsigned</span>
<span id="L1086" class="ln">  1086&nbsp;&nbsp;</span>
<span id="L1087" class="ln">  1087&nbsp;&nbsp;</span>				if cap(paramValues)-len(paramValues)-8 &gt;= 0 {
<span id="L1088" class="ln">  1088&nbsp;&nbsp;</span>					paramValues = paramValues[:len(paramValues)+8]
<span id="L1089" class="ln">  1089&nbsp;&nbsp;</span>					binary.LittleEndian.PutUint64(
<span id="L1090" class="ln">  1090&nbsp;&nbsp;</span>						paramValues[len(paramValues)-8:],
<span id="L1091" class="ln">  1091&nbsp;&nbsp;</span>						uint64(v),
<span id="L1092" class="ln">  1092&nbsp;&nbsp;</span>					)
<span id="L1093" class="ln">  1093&nbsp;&nbsp;</span>				} else {
<span id="L1094" class="ln">  1094&nbsp;&nbsp;</span>					paramValues = append(paramValues,
<span id="L1095" class="ln">  1095&nbsp;&nbsp;</span>						uint64ToBytes(uint64(v))...,
<span id="L1096" class="ln">  1096&nbsp;&nbsp;</span>					)
<span id="L1097" class="ln">  1097&nbsp;&nbsp;</span>				}
<span id="L1098" class="ln">  1098&nbsp;&nbsp;</span>
<span id="L1099" class="ln">  1099&nbsp;&nbsp;</span>			case float64:
<span id="L1100" class="ln">  1100&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeDouble)
<span id="L1101" class="ln">  1101&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1102" class="ln">  1102&nbsp;&nbsp;</span>
<span id="L1103" class="ln">  1103&nbsp;&nbsp;</span>				if cap(paramValues)-len(paramValues)-8 &gt;= 0 {
<span id="L1104" class="ln">  1104&nbsp;&nbsp;</span>					paramValues = paramValues[:len(paramValues)+8]
<span id="L1105" class="ln">  1105&nbsp;&nbsp;</span>					binary.LittleEndian.PutUint64(
<span id="L1106" class="ln">  1106&nbsp;&nbsp;</span>						paramValues[len(paramValues)-8:],
<span id="L1107" class="ln">  1107&nbsp;&nbsp;</span>						math.Float64bits(v),
<span id="L1108" class="ln">  1108&nbsp;&nbsp;</span>					)
<span id="L1109" class="ln">  1109&nbsp;&nbsp;</span>				} else {
<span id="L1110" class="ln">  1110&nbsp;&nbsp;</span>					paramValues = append(paramValues,
<span id="L1111" class="ln">  1111&nbsp;&nbsp;</span>						uint64ToBytes(math.Float64bits(v))...,
<span id="L1112" class="ln">  1112&nbsp;&nbsp;</span>					)
<span id="L1113" class="ln">  1113&nbsp;&nbsp;</span>				}
<span id="L1114" class="ln">  1114&nbsp;&nbsp;</span>
<span id="L1115" class="ln">  1115&nbsp;&nbsp;</span>			case bool:
<span id="L1116" class="ln">  1116&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeTiny)
<span id="L1117" class="ln">  1117&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1118" class="ln">  1118&nbsp;&nbsp;</span>
<span id="L1119" class="ln">  1119&nbsp;&nbsp;</span>				if v {
<span id="L1120" class="ln">  1120&nbsp;&nbsp;</span>					paramValues = append(paramValues, 0x01)
<span id="L1121" class="ln">  1121&nbsp;&nbsp;</span>				} else {
<span id="L1122" class="ln">  1122&nbsp;&nbsp;</span>					paramValues = append(paramValues, 0x00)
<span id="L1123" class="ln">  1123&nbsp;&nbsp;</span>				}
<span id="L1124" class="ln">  1124&nbsp;&nbsp;</span>
<span id="L1125" class="ln">  1125&nbsp;&nbsp;</span>			case []byte:
<span id="L1126" class="ln">  1126&nbsp;&nbsp;</span>				<span class="comment">// Common case (non-nil value) first</span>
<span id="L1127" class="ln">  1127&nbsp;&nbsp;</span>				if v != nil {
<span id="L1128" class="ln">  1128&nbsp;&nbsp;</span>					paramTypes[i+i] = byte(fieldTypeString)
<span id="L1129" class="ln">  1129&nbsp;&nbsp;</span>					paramTypes[i+i+1] = 0x00
<span id="L1130" class="ln">  1130&nbsp;&nbsp;</span>
<span id="L1131" class="ln">  1131&nbsp;&nbsp;</span>					if len(v) &lt; longDataSize {
<span id="L1132" class="ln">  1132&nbsp;&nbsp;</span>						paramValues = appendLengthEncodedInteger(paramValues,
<span id="L1133" class="ln">  1133&nbsp;&nbsp;</span>							uint64(len(v)),
<span id="L1134" class="ln">  1134&nbsp;&nbsp;</span>						)
<span id="L1135" class="ln">  1135&nbsp;&nbsp;</span>						paramValues = append(paramValues, v...)
<span id="L1136" class="ln">  1136&nbsp;&nbsp;</span>					} else {
<span id="L1137" class="ln">  1137&nbsp;&nbsp;</span>						if err := stmt.writeCommandLongData(i, v); err != nil {
<span id="L1138" class="ln">  1138&nbsp;&nbsp;</span>							return err
<span id="L1139" class="ln">  1139&nbsp;&nbsp;</span>						}
<span id="L1140" class="ln">  1140&nbsp;&nbsp;</span>					}
<span id="L1141" class="ln">  1141&nbsp;&nbsp;</span>					continue
<span id="L1142" class="ln">  1142&nbsp;&nbsp;</span>				}
<span id="L1143" class="ln">  1143&nbsp;&nbsp;</span>
<span id="L1144" class="ln">  1144&nbsp;&nbsp;</span>				<span class="comment">// Handle []byte(nil) as a NULL value</span>
<span id="L1145" class="ln">  1145&nbsp;&nbsp;</span>				nullMask[i/8] |= 1 &lt;&lt; (uint(i) &amp; 7)
<span id="L1146" class="ln">  1146&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeNULL)
<span id="L1147" class="ln">  1147&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1148" class="ln">  1148&nbsp;&nbsp;</span>
<span id="L1149" class="ln">  1149&nbsp;&nbsp;</span>			case string:
<span id="L1150" class="ln">  1150&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeString)
<span id="L1151" class="ln">  1151&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1152" class="ln">  1152&nbsp;&nbsp;</span>
<span id="L1153" class="ln">  1153&nbsp;&nbsp;</span>				if len(v) &lt; longDataSize {
<span id="L1154" class="ln">  1154&nbsp;&nbsp;</span>					paramValues = appendLengthEncodedInteger(paramValues,
<span id="L1155" class="ln">  1155&nbsp;&nbsp;</span>						uint64(len(v)),
<span id="L1156" class="ln">  1156&nbsp;&nbsp;</span>					)
<span id="L1157" class="ln">  1157&nbsp;&nbsp;</span>					paramValues = append(paramValues, v...)
<span id="L1158" class="ln">  1158&nbsp;&nbsp;</span>				} else {
<span id="L1159" class="ln">  1159&nbsp;&nbsp;</span>					if err := stmt.writeCommandLongData(i, []byte(v)); err != nil {
<span id="L1160" class="ln">  1160&nbsp;&nbsp;</span>						return err
<span id="L1161" class="ln">  1161&nbsp;&nbsp;</span>					}
<span id="L1162" class="ln">  1162&nbsp;&nbsp;</span>				}
<span id="L1163" class="ln">  1163&nbsp;&nbsp;</span>
<span id="L1164" class="ln">  1164&nbsp;&nbsp;</span>			case time.Time:
<span id="L1165" class="ln">  1165&nbsp;&nbsp;</span>				paramTypes[i+i] = byte(fieldTypeString)
<span id="L1166" class="ln">  1166&nbsp;&nbsp;</span>				paramTypes[i+i+1] = 0x00
<span id="L1167" class="ln">  1167&nbsp;&nbsp;</span>
<span id="L1168" class="ln">  1168&nbsp;&nbsp;</span>				var a [64]byte
<span id="L1169" class="ln">  1169&nbsp;&nbsp;</span>				var b = a[:0]
<span id="L1170" class="ln">  1170&nbsp;&nbsp;</span>
<span id="L1171" class="ln">  1171&nbsp;&nbsp;</span>				if v.IsZero() {
<span id="L1172" class="ln">  1172&nbsp;&nbsp;</span>					b = append(b, &#34;0000-00-00&#34;...)
<span id="L1173" class="ln">  1173&nbsp;&nbsp;</span>				} else {
<span id="L1174" class="ln">  1174&nbsp;&nbsp;</span>					b, err = appendDateTime(b, v.In(mc.cfg.Loc), mc.cfg.timeTruncate)
<span id="L1175" class="ln">  1175&nbsp;&nbsp;</span>					if err != nil {
<span id="L1176" class="ln">  1176&nbsp;&nbsp;</span>						return err
<span id="L1177" class="ln">  1177&nbsp;&nbsp;</span>					}
<span id="L1178" class="ln">  1178&nbsp;&nbsp;</span>				}
<span id="L1179" class="ln">  1179&nbsp;&nbsp;</span>
<span id="L1180" class="ln">  1180&nbsp;&nbsp;</span>				paramValues = appendLengthEncodedInteger(paramValues,
<span id="L1181" class="ln">  1181&nbsp;&nbsp;</span>					uint64(len(b)),
<span id="L1182" class="ln">  1182&nbsp;&nbsp;</span>				)
<span id="L1183" class="ln">  1183&nbsp;&nbsp;</span>				paramValues = append(paramValues, b...)
<span id="L1184" class="ln">  1184&nbsp;&nbsp;</span>
<span id="L1185" class="ln">  1185&nbsp;&nbsp;</span>			default:
<span id="L1186" class="ln">  1186&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;cannot convert type: %T&#34;, arg)
<span id="L1187" class="ln">  1187&nbsp;&nbsp;</span>			}
<span id="L1188" class="ln">  1188&nbsp;&nbsp;</span>		}
<span id="L1189" class="ln">  1189&nbsp;&nbsp;</span>
<span id="L1190" class="ln">  1190&nbsp;&nbsp;</span>		<span class="comment">// Check if param values exceeded the available buffer</span>
<span id="L1191" class="ln">  1191&nbsp;&nbsp;</span>		<span class="comment">// In that case we must build the data packet with the new values buffer</span>
<span id="L1192" class="ln">  1192&nbsp;&nbsp;</span>		if valuesCap != cap(paramValues) {
<span id="L1193" class="ln">  1193&nbsp;&nbsp;</span>			data = append(data[:pos], paramValues...)
<span id="L1194" class="ln">  1194&nbsp;&nbsp;</span>			if err = mc.buf.store(data); err != nil {
<span id="L1195" class="ln">  1195&nbsp;&nbsp;</span>				mc.log(err)
<span id="L1196" class="ln">  1196&nbsp;&nbsp;</span>				return errBadConnNoWrite
<span id="L1197" class="ln">  1197&nbsp;&nbsp;</span>			}
<span id="L1198" class="ln">  1198&nbsp;&nbsp;</span>		}
<span id="L1199" class="ln">  1199&nbsp;&nbsp;</span>
<span id="L1200" class="ln">  1200&nbsp;&nbsp;</span>		pos += len(paramValues)
<span id="L1201" class="ln">  1201&nbsp;&nbsp;</span>		data = data[:pos]
<span id="L1202" class="ln">  1202&nbsp;&nbsp;</span>	}
<span id="L1203" class="ln">  1203&nbsp;&nbsp;</span>
<span id="L1204" class="ln">  1204&nbsp;&nbsp;</span>	return mc.writePacket(data)
<span id="L1205" class="ln">  1205&nbsp;&nbsp;</span>}
<span id="L1206" class="ln">  1206&nbsp;&nbsp;</span>
<span id="L1207" class="ln">  1207&nbsp;&nbsp;</span><span class="comment">// For each remaining resultset in the stream, discards its rows and updates</span>
<span id="L1208" class="ln">  1208&nbsp;&nbsp;</span><span class="comment">// mc.affectedRows and mc.insertIds.</span>
<span id="L1209" class="ln">  1209&nbsp;&nbsp;</span>func (mc *okHandler) discardResults() error {
<span id="L1210" class="ln">  1210&nbsp;&nbsp;</span>	for mc.status&amp;statusMoreResultsExists != 0 {
<span id="L1211" class="ln">  1211&nbsp;&nbsp;</span>		resLen, err := mc.readResultSetHeaderPacket()
<span id="L1212" class="ln">  1212&nbsp;&nbsp;</span>		if err != nil {
<span id="L1213" class="ln">  1213&nbsp;&nbsp;</span>			return err
<span id="L1214" class="ln">  1214&nbsp;&nbsp;</span>		}
<span id="L1215" class="ln">  1215&nbsp;&nbsp;</span>		if resLen &gt; 0 {
<span id="L1216" class="ln">  1216&nbsp;&nbsp;</span>			<span class="comment">// columns</span>
<span id="L1217" class="ln">  1217&nbsp;&nbsp;</span>			if err := mc.conn().readUntilEOF(); err != nil {
<span id="L1218" class="ln">  1218&nbsp;&nbsp;</span>				return err
<span id="L1219" class="ln">  1219&nbsp;&nbsp;</span>			}
<span id="L1220" class="ln">  1220&nbsp;&nbsp;</span>			<span class="comment">// rows</span>
<span id="L1221" class="ln">  1221&nbsp;&nbsp;</span>			if err := mc.conn().readUntilEOF(); err != nil {
<span id="L1222" class="ln">  1222&nbsp;&nbsp;</span>				return err
<span id="L1223" class="ln">  1223&nbsp;&nbsp;</span>			}
<span id="L1224" class="ln">  1224&nbsp;&nbsp;</span>		}
<span id="L1225" class="ln">  1225&nbsp;&nbsp;</span>	}
<span id="L1226" class="ln">  1226&nbsp;&nbsp;</span>	return nil
<span id="L1227" class="ln">  1227&nbsp;&nbsp;</span>}
<span id="L1228" class="ln">  1228&nbsp;&nbsp;</span>
<span id="L1229" class="ln">  1229&nbsp;&nbsp;</span><span class="comment">// http://dev.mysql.com/doc/internals/en/binary-protocol-resultset-row.html</span>
<span id="L1230" class="ln">  1230&nbsp;&nbsp;</span>func (rows *binaryRows) readRow(dest []driver.Value) error {
<span id="L1231" class="ln">  1231&nbsp;&nbsp;</span>	data, err := rows.mc.readPacket()
<span id="L1232" class="ln">  1232&nbsp;&nbsp;</span>	if err != nil {
<span id="L1233" class="ln">  1233&nbsp;&nbsp;</span>		return err
<span id="L1234" class="ln">  1234&nbsp;&nbsp;</span>	}
<span id="L1235" class="ln">  1235&nbsp;&nbsp;</span>
<span id="L1236" class="ln">  1236&nbsp;&nbsp;</span>	<span class="comment">// packet indicator [1 byte]</span>
<span id="L1237" class="ln">  1237&nbsp;&nbsp;</span>	if data[0] != iOK {
<span id="L1238" class="ln">  1238&nbsp;&nbsp;</span>		<span class="comment">// EOF Packet</span>
<span id="L1239" class="ln">  1239&nbsp;&nbsp;</span>		if data[0] == iEOF &amp;&amp; len(data) == 5 {
<span id="L1240" class="ln">  1240&nbsp;&nbsp;</span>			rows.mc.status = readStatus(data[3:])
<span id="L1241" class="ln">  1241&nbsp;&nbsp;</span>			rows.rs.done = true
<span id="L1242" class="ln">  1242&nbsp;&nbsp;</span>			if !rows.HasNextResultSet() {
<span id="L1243" class="ln">  1243&nbsp;&nbsp;</span>				rows.mc = nil
<span id="L1244" class="ln">  1244&nbsp;&nbsp;</span>			}
<span id="L1245" class="ln">  1245&nbsp;&nbsp;</span>			return io.EOF
<span id="L1246" class="ln">  1246&nbsp;&nbsp;</span>		}
<span id="L1247" class="ln">  1247&nbsp;&nbsp;</span>		mc := rows.mc
<span id="L1248" class="ln">  1248&nbsp;&nbsp;</span>		rows.mc = nil
<span id="L1249" class="ln">  1249&nbsp;&nbsp;</span>
<span id="L1250" class="ln">  1250&nbsp;&nbsp;</span>		<span class="comment">// Error otherwise</span>
<span id="L1251" class="ln">  1251&nbsp;&nbsp;</span>		return mc.handleErrorPacket(data)
<span id="L1252" class="ln">  1252&nbsp;&nbsp;</span>	}
<span id="L1253" class="ln">  1253&nbsp;&nbsp;</span>
<span id="L1254" class="ln">  1254&nbsp;&nbsp;</span>	<span class="comment">// NULL-bitmap,  [(column-count + 7 + 2) / 8 bytes]</span>
<span id="L1255" class="ln">  1255&nbsp;&nbsp;</span>	pos := 1 + (len(dest)+7+2)&gt;&gt;3
<span id="L1256" class="ln">  1256&nbsp;&nbsp;</span>	nullMask := data[1:pos]
<span id="L1257" class="ln">  1257&nbsp;&nbsp;</span>
<span id="L1258" class="ln">  1258&nbsp;&nbsp;</span>	for i := range dest {
<span id="L1259" class="ln">  1259&nbsp;&nbsp;</span>		<span class="comment">// Field is NULL</span>
<span id="L1260" class="ln">  1260&nbsp;&nbsp;</span>		<span class="comment">// (byte &gt;&gt; bit-pos) % 2 == 1</span>
<span id="L1261" class="ln">  1261&nbsp;&nbsp;</span>		if ((nullMask[(i+2)&gt;&gt;3] &gt;&gt; uint((i+2)&amp;7)) &amp; 1) == 1 {
<span id="L1262" class="ln">  1262&nbsp;&nbsp;</span>			dest[i] = nil
<span id="L1263" class="ln">  1263&nbsp;&nbsp;</span>			continue
<span id="L1264" class="ln">  1264&nbsp;&nbsp;</span>		}
<span id="L1265" class="ln">  1265&nbsp;&nbsp;</span>
<span id="L1266" class="ln">  1266&nbsp;&nbsp;</span>		<span class="comment">// Convert to byte-coded string</span>
<span id="L1267" class="ln">  1267&nbsp;&nbsp;</span>		switch rows.rs.columns[i].fieldType {
<span id="L1268" class="ln">  1268&nbsp;&nbsp;</span>		case fieldTypeNULL:
<span id="L1269" class="ln">  1269&nbsp;&nbsp;</span>			dest[i] = nil
<span id="L1270" class="ln">  1270&nbsp;&nbsp;</span>			continue
<span id="L1271" class="ln">  1271&nbsp;&nbsp;</span>
<span id="L1272" class="ln">  1272&nbsp;&nbsp;</span>		<span class="comment">// Numeric Types</span>
<span id="L1273" class="ln">  1273&nbsp;&nbsp;</span>		case fieldTypeTiny:
<span id="L1274" class="ln">  1274&nbsp;&nbsp;</span>			if rows.rs.columns[i].flags&amp;flagUnsigned != 0 {
<span id="L1275" class="ln">  1275&nbsp;&nbsp;</span>				dest[i] = int64(data[pos])
<span id="L1276" class="ln">  1276&nbsp;&nbsp;</span>			} else {
<span id="L1277" class="ln">  1277&nbsp;&nbsp;</span>				dest[i] = int64(int8(data[pos]))
<span id="L1278" class="ln">  1278&nbsp;&nbsp;</span>			}
<span id="L1279" class="ln">  1279&nbsp;&nbsp;</span>			pos++
<span id="L1280" class="ln">  1280&nbsp;&nbsp;</span>			continue
<span id="L1281" class="ln">  1281&nbsp;&nbsp;</span>
<span id="L1282" class="ln">  1282&nbsp;&nbsp;</span>		case fieldTypeShort, fieldTypeYear:
<span id="L1283" class="ln">  1283&nbsp;&nbsp;</span>			if rows.rs.columns[i].flags&amp;flagUnsigned != 0 {
<span id="L1284" class="ln">  1284&nbsp;&nbsp;</span>				dest[i] = int64(binary.LittleEndian.Uint16(data[pos : pos+2]))
<span id="L1285" class="ln">  1285&nbsp;&nbsp;</span>			} else {
<span id="L1286" class="ln">  1286&nbsp;&nbsp;</span>				dest[i] = int64(int16(binary.LittleEndian.Uint16(data[pos : pos+2])))
<span id="L1287" class="ln">  1287&nbsp;&nbsp;</span>			}
<span id="L1288" class="ln">  1288&nbsp;&nbsp;</span>			pos += 2
<span id="L1289" class="ln">  1289&nbsp;&nbsp;</span>			continue
<span id="L1290" class="ln">  1290&nbsp;&nbsp;</span>
<span id="L1291" class="ln">  1291&nbsp;&nbsp;</span>		case fieldTypeInt24, fieldTypeLong:
<span id="L1292" class="ln">  1292&nbsp;&nbsp;</span>			if rows.rs.columns[i].flags&amp;flagUnsigned != 0 {
<span id="L1293" class="ln">  1293&nbsp;&nbsp;</span>				dest[i] = int64(binary.LittleEndian.Uint32(data[pos : pos+4]))
<span id="L1294" class="ln">  1294&nbsp;&nbsp;</span>			} else {
<span id="L1295" class="ln">  1295&nbsp;&nbsp;</span>				dest[i] = int64(int32(binary.LittleEndian.Uint32(data[pos : pos+4])))
<span id="L1296" class="ln">  1296&nbsp;&nbsp;</span>			}
<span id="L1297" class="ln">  1297&nbsp;&nbsp;</span>			pos += 4
<span id="L1298" class="ln">  1298&nbsp;&nbsp;</span>			continue
<span id="L1299" class="ln">  1299&nbsp;&nbsp;</span>
<span id="L1300" class="ln">  1300&nbsp;&nbsp;</span>		case fieldTypeLongLong:
<span id="L1301" class="ln">  1301&nbsp;&nbsp;</span>			if rows.rs.columns[i].flags&amp;flagUnsigned != 0 {
<span id="L1302" class="ln">  1302&nbsp;&nbsp;</span>				val := binary.LittleEndian.Uint64(data[pos : pos+8])
<span id="L1303" class="ln">  1303&nbsp;&nbsp;</span>				if val &gt; math.MaxInt64 {
<span id="L1304" class="ln">  1304&nbsp;&nbsp;</span>					dest[i] = uint64ToString(val)
<span id="L1305" class="ln">  1305&nbsp;&nbsp;</span>				} else {
<span id="L1306" class="ln">  1306&nbsp;&nbsp;</span>					dest[i] = int64(val)
<span id="L1307" class="ln">  1307&nbsp;&nbsp;</span>				}
<span id="L1308" class="ln">  1308&nbsp;&nbsp;</span>			} else {
<span id="L1309" class="ln">  1309&nbsp;&nbsp;</span>				dest[i] = int64(binary.LittleEndian.Uint64(data[pos : pos+8]))
<span id="L1310" class="ln">  1310&nbsp;&nbsp;</span>			}
<span id="L1311" class="ln">  1311&nbsp;&nbsp;</span>			pos += 8
<span id="L1312" class="ln">  1312&nbsp;&nbsp;</span>			continue
<span id="L1313" class="ln">  1313&nbsp;&nbsp;</span>
<span id="L1314" class="ln">  1314&nbsp;&nbsp;</span>		case fieldTypeFloat:
<span id="L1315" class="ln">  1315&nbsp;&nbsp;</span>			dest[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[pos : pos+4]))
<span id="L1316" class="ln">  1316&nbsp;&nbsp;</span>			pos += 4
<span id="L1317" class="ln">  1317&nbsp;&nbsp;</span>			continue
<span id="L1318" class="ln">  1318&nbsp;&nbsp;</span>
<span id="L1319" class="ln">  1319&nbsp;&nbsp;</span>		case fieldTypeDouble:
<span id="L1320" class="ln">  1320&nbsp;&nbsp;</span>			dest[i] = math.Float64frombits(binary.LittleEndian.Uint64(data[pos : pos+8]))
<span id="L1321" class="ln">  1321&nbsp;&nbsp;</span>			pos += 8
<span id="L1322" class="ln">  1322&nbsp;&nbsp;</span>			continue
<span id="L1323" class="ln">  1323&nbsp;&nbsp;</span>
<span id="L1324" class="ln">  1324&nbsp;&nbsp;</span>		<span class="comment">// Length coded Binary Strings</span>
<span id="L1325" class="ln">  1325&nbsp;&nbsp;</span>		case fieldTypeDecimal, fieldTypeNewDecimal, fieldTypeVarChar,
<span id="L1326" class="ln">  1326&nbsp;&nbsp;</span>			fieldTypeBit, fieldTypeEnum, fieldTypeSet, fieldTypeTinyBLOB,
<span id="L1327" class="ln">  1327&nbsp;&nbsp;</span>			fieldTypeMediumBLOB, fieldTypeLongBLOB, fieldTypeBLOB,
<span id="L1328" class="ln">  1328&nbsp;&nbsp;</span>			fieldTypeVarString, fieldTypeString, fieldTypeGeometry, fieldTypeJSON:
<span id="L1329" class="ln">  1329&nbsp;&nbsp;</span>			var isNull bool
<span id="L1330" class="ln">  1330&nbsp;&nbsp;</span>			var n int
<span id="L1331" class="ln">  1331&nbsp;&nbsp;</span>			dest[i], isNull, n, err = readLengthEncodedString(data[pos:])
<span id="L1332" class="ln">  1332&nbsp;&nbsp;</span>			pos += n
<span id="L1333" class="ln">  1333&nbsp;&nbsp;</span>			if err == nil {
<span id="L1334" class="ln">  1334&nbsp;&nbsp;</span>				if !isNull {
<span id="L1335" class="ln">  1335&nbsp;&nbsp;</span>					continue
<span id="L1336" class="ln">  1336&nbsp;&nbsp;</span>				} else {
<span id="L1337" class="ln">  1337&nbsp;&nbsp;</span>					dest[i] = nil
<span id="L1338" class="ln">  1338&nbsp;&nbsp;</span>					continue
<span id="L1339" class="ln">  1339&nbsp;&nbsp;</span>				}
<span id="L1340" class="ln">  1340&nbsp;&nbsp;</span>			}
<span id="L1341" class="ln">  1341&nbsp;&nbsp;</span>			return err
<span id="L1342" class="ln">  1342&nbsp;&nbsp;</span>
<span id="L1343" class="ln">  1343&nbsp;&nbsp;</span>		case
<span id="L1344" class="ln">  1344&nbsp;&nbsp;</span>			fieldTypeDate, fieldTypeNewDate, <span class="comment">// Date YYYY-MM-DD</span>
<span id="L1345" class="ln">  1345&nbsp;&nbsp;</span>			fieldTypeTime,                         <span class="comment">// Time [-][H]HH:MM:SS[.fractal]</span>
<span id="L1346" class="ln">  1346&nbsp;&nbsp;</span>			fieldTypeTimestamp, fieldTypeDateTime: <span class="comment">// Timestamp YYYY-MM-DD HH:MM:SS[.fractal]</span>
<span id="L1347" class="ln">  1347&nbsp;&nbsp;</span>
<span id="L1348" class="ln">  1348&nbsp;&nbsp;</span>			num, isNull, n := readLengthEncodedInteger(data[pos:])
<span id="L1349" class="ln">  1349&nbsp;&nbsp;</span>			pos += n
<span id="L1350" class="ln">  1350&nbsp;&nbsp;</span>
<span id="L1351" class="ln">  1351&nbsp;&nbsp;</span>			switch {
<span id="L1352" class="ln">  1352&nbsp;&nbsp;</span>			case isNull:
<span id="L1353" class="ln">  1353&nbsp;&nbsp;</span>				dest[i] = nil
<span id="L1354" class="ln">  1354&nbsp;&nbsp;</span>				continue
<span id="L1355" class="ln">  1355&nbsp;&nbsp;</span>			case rows.rs.columns[i].fieldType == fieldTypeTime:
<span id="L1356" class="ln">  1356&nbsp;&nbsp;</span>				<span class="comment">// database/sql does not support an equivalent to TIME, return a string</span>
<span id="L1357" class="ln">  1357&nbsp;&nbsp;</span>				var dstlen uint8
<span id="L1358" class="ln">  1358&nbsp;&nbsp;</span>				switch decimals := rows.rs.columns[i].decimals; decimals {
<span id="L1359" class="ln">  1359&nbsp;&nbsp;</span>				case 0x00, 0x1f:
<span id="L1360" class="ln">  1360&nbsp;&nbsp;</span>					dstlen = 8
<span id="L1361" class="ln">  1361&nbsp;&nbsp;</span>				case 1, 2, 3, 4, 5, 6:
<span id="L1362" class="ln">  1362&nbsp;&nbsp;</span>					dstlen = 8 + 1 + decimals
<span id="L1363" class="ln">  1363&nbsp;&nbsp;</span>				default:
<span id="L1364" class="ln">  1364&nbsp;&nbsp;</span>					return fmt.Errorf(
<span id="L1365" class="ln">  1365&nbsp;&nbsp;</span>						&#34;protocol error, illegal decimals value %d&#34;,
<span id="L1366" class="ln">  1366&nbsp;&nbsp;</span>						rows.rs.columns[i].decimals,
<span id="L1367" class="ln">  1367&nbsp;&nbsp;</span>					)
<span id="L1368" class="ln">  1368&nbsp;&nbsp;</span>				}
<span id="L1369" class="ln">  1369&nbsp;&nbsp;</span>				dest[i], err = formatBinaryTime(data[pos:pos+int(num)], dstlen)
<span id="L1370" class="ln">  1370&nbsp;&nbsp;</span>			case rows.mc.parseTime:
<span id="L1371" class="ln">  1371&nbsp;&nbsp;</span>				dest[i], err = parseBinaryDateTime(num, data[pos:], rows.mc.cfg.Loc)
<span id="L1372" class="ln">  1372&nbsp;&nbsp;</span>			default:
<span id="L1373" class="ln">  1373&nbsp;&nbsp;</span>				var dstlen uint8
<span id="L1374" class="ln">  1374&nbsp;&nbsp;</span>				if rows.rs.columns[i].fieldType == fieldTypeDate {
<span id="L1375" class="ln">  1375&nbsp;&nbsp;</span>					dstlen = 10
<span id="L1376" class="ln">  1376&nbsp;&nbsp;</span>				} else {
<span id="L1377" class="ln">  1377&nbsp;&nbsp;</span>					switch decimals := rows.rs.columns[i].decimals; decimals {
<span id="L1378" class="ln">  1378&nbsp;&nbsp;</span>					case 0x00, 0x1f:
<span id="L1379" class="ln">  1379&nbsp;&nbsp;</span>						dstlen = 19
<span id="L1380" class="ln">  1380&nbsp;&nbsp;</span>					case 1, 2, 3, 4, 5, 6:
<span id="L1381" class="ln">  1381&nbsp;&nbsp;</span>						dstlen = 19 + 1 + decimals
<span id="L1382" class="ln">  1382&nbsp;&nbsp;</span>					default:
<span id="L1383" class="ln">  1383&nbsp;&nbsp;</span>						return fmt.Errorf(
<span id="L1384" class="ln">  1384&nbsp;&nbsp;</span>							&#34;protocol error, illegal decimals value %d&#34;,
<span id="L1385" class="ln">  1385&nbsp;&nbsp;</span>							rows.rs.columns[i].decimals,
<span id="L1386" class="ln">  1386&nbsp;&nbsp;</span>						)
<span id="L1387" class="ln">  1387&nbsp;&nbsp;</span>					}
<span id="L1388" class="ln">  1388&nbsp;&nbsp;</span>				}
<span id="L1389" class="ln">  1389&nbsp;&nbsp;</span>				dest[i], err = formatBinaryDateTime(data[pos:pos+int(num)], dstlen)
<span id="L1390" class="ln">  1390&nbsp;&nbsp;</span>			}
<span id="L1391" class="ln">  1391&nbsp;&nbsp;</span>
<span id="L1392" class="ln">  1392&nbsp;&nbsp;</span>			if err == nil {
<span id="L1393" class="ln">  1393&nbsp;&nbsp;</span>				pos += int(num)
<span id="L1394" class="ln">  1394&nbsp;&nbsp;</span>				continue
<span id="L1395" class="ln">  1395&nbsp;&nbsp;</span>			} else {
<span id="L1396" class="ln">  1396&nbsp;&nbsp;</span>				return err
<span id="L1397" class="ln">  1397&nbsp;&nbsp;</span>			}
<span id="L1398" class="ln">  1398&nbsp;&nbsp;</span>
<span id="L1399" class="ln">  1399&nbsp;&nbsp;</span>		<span class="comment">// Please report if this happens!</span>
<span id="L1400" class="ln">  1400&nbsp;&nbsp;</span>		default:
<span id="L1401" class="ln">  1401&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;unknown field type %d&#34;, rows.rs.columns[i].fieldType)
<span id="L1402" class="ln">  1402&nbsp;&nbsp;</span>		}
<span id="L1403" class="ln">  1403&nbsp;&nbsp;</span>	}
<span id="L1404" class="ln">  1404&nbsp;&nbsp;</span>
<span id="L1405" class="ln">  1405&nbsp;&nbsp;</span>	return nil
<span id="L1406" class="ln">  1406&nbsp;&nbsp;</span>}
<span id="L1407" class="ln">  1407&nbsp;&nbsp;</span>
</pre><p><a href="packets.go?m=text">View as plain text</a></p>

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
