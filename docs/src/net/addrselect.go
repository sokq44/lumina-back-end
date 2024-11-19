<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/addrselect.go - Go Documentation Server</title>

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
<a href="addrselect.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">addrselect.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net">net</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Minimal RFC 6724 address selection.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package net
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;net/netip&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func sortByRFC6724(addrs []IPAddr) {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	if len(addrs) &lt; 2 {
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>		return
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	sortByRFC6724withSrcs(addrs, srcAddrs(addrs))
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>func sortByRFC6724withSrcs(addrs []IPAddr, srcs []netip.Addr) {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	if len(addrs) != len(srcs) {
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>		panic(&#34;internal error&#34;)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	}
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	addrAttr := make([]ipAttr, len(addrs))
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	srcAttr := make([]ipAttr, len(srcs))
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	for i, v := range addrs {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		addrAttrIP, _ := netip.AddrFromSlice(v.IP)
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		addrAttr[i] = ipAttrOf(addrAttrIP)
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		srcAttr[i] = ipAttrOf(srcs[i])
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	sort.Stable(&amp;byRFC6724{
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		addrs:    addrs,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		addrAttr: addrAttr,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		srcs:     srcs,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		srcAttr:  srcAttr,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	})
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// srcAddrs tries to UDP-connect to each address to see if it has a</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// route. (This doesn&#39;t send any packets). The destination port</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// number is irrelevant.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>func srcAddrs(addrs []IPAddr) []netip.Addr {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	srcs := make([]netip.Addr, len(addrs))
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	dst := UDPAddr{Port: 9}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	for i := range addrs {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		dst.IP = addrs[i].IP
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		dst.Zone = addrs[i].Zone
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		c, err := DialUDP(&#34;udp&#34;, nil, &amp;dst)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		if err == nil {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			if src, ok := c.LocalAddr().(*UDPAddr); ok {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>				srcs[i], _ = netip.AddrFromSlice(src.IP)
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>			c.Close()
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	return srcs
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>type ipAttr struct {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	Scope      scope
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	Precedence uint8
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	Label      uint8
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func ipAttrOf(ip netip.Addr) ipAttr {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if !ip.IsValid() {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		return ipAttr{}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	match := rfc6724policyTable.Classify(ip)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	return ipAttr{
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		Scope:      classifyScope(ip),
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		Precedence: match.Precedence,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		Label:      match.Label,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>type byRFC6724 struct {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	addrs    []IPAddr <span class="comment">// addrs to sort</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	addrAttr []ipAttr
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	srcs     []netip.Addr <span class="comment">// or not valid addr if unreachable</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	srcAttr  []ipAttr
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>func (s *byRFC6724) Len() int { return len(s.addrs) }
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>func (s *byRFC6724) Swap(i, j int) {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	s.addrs[i], s.addrs[j] = s.addrs[j], s.addrs[i]
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	s.srcs[i], s.srcs[j] = s.srcs[j], s.srcs[i]
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	s.addrAttr[i], s.addrAttr[j] = s.addrAttr[j], s.addrAttr[i]
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	s.srcAttr[i], s.srcAttr[j] = s.srcAttr[j], s.srcAttr[i]
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Less reports whether i is a better destination address for this</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// host than j.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span><span class="comment">// The algorithm and variable names comes from RFC 6724 section 6.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>func (s *byRFC6724) Less(i, j int) bool {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	DA := s.addrs[i].IP
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	DB := s.addrs[j].IP
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	SourceDA := s.srcs[i]
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	SourceDB := s.srcs[j]
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	attrDA := &amp;s.addrAttr[i]
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	attrDB := &amp;s.addrAttr[j]
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	attrSourceDA := &amp;s.srcAttr[i]
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	attrSourceDB := &amp;s.srcAttr[j]
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	const preferDA = true
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	const preferDB = false
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	<span class="comment">// Rule 1: Avoid unusable destinations.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	<span class="comment">// If DB is known to be unreachable or if Source(DB) is undefined, then</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">// prefer DA.  Similarly, if DA is known to be unreachable or if</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// Source(DA) is undefined, then prefer DB.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if !SourceDA.IsValid() &amp;&amp; !SourceDB.IsValid() {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		return false <span class="comment">// &#34;equal&#34;</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	if !SourceDB.IsValid() {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		return preferDA
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if !SourceDA.IsValid() {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return preferDB
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// Rule 2: Prefer matching scope.</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// If Scope(DA) = Scope(Source(DA)) and Scope(DB) &lt;&gt; Scope(Source(DB)),</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// then prefer DA.  Similarly, if Scope(DA) &lt;&gt; Scope(Source(DA)) and</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// Scope(DB) = Scope(Source(DB)), then prefer DB.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if attrDA.Scope == attrSourceDA.Scope &amp;&amp; attrDB.Scope != attrSourceDB.Scope {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		return preferDA
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	if attrDA.Scope != attrSourceDA.Scope &amp;&amp; attrDB.Scope == attrSourceDB.Scope {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		return preferDB
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Rule 3: Avoid deprecated addresses.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// If Source(DA) is deprecated and Source(DB) is not, then prefer DB.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// Similarly, if Source(DA) is not deprecated and Source(DB) is</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">// deprecated, then prefer DA.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">// TODO(bradfitz): implement? low priority for now.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// Rule 4: Prefer home addresses.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">// If Source(DA) is simultaneously a home address and care-of address</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// and Source(DB) is not, then prefer DA.  Similarly, if Source(DB) is</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">// simultaneously a home address and care-of address and Source(DA) is</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">// not, then prefer DB.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">// TODO(bradfitz): implement? low priority for now.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Rule 5: Prefer matching label.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">// If Label(Source(DA)) = Label(DA) and Label(Source(DB)) &lt;&gt; Label(DB),</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// then prefer DA.  Similarly, if Label(Source(DA)) &lt;&gt; Label(DA) and</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// Label(Source(DB)) = Label(DB), then prefer DB.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if attrSourceDA.Label == attrDA.Label &amp;&amp;
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		attrSourceDB.Label != attrDB.Label {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return preferDA
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	if attrSourceDA.Label != attrDA.Label &amp;&amp;
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		attrSourceDB.Label == attrDB.Label {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		return preferDB
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	}
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// Rule 6: Prefer higher precedence.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// If Precedence(DA) &gt; Precedence(DB), then prefer DA.  Similarly, if</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// Precedence(DA) &lt; Precedence(DB), then prefer DB.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if attrDA.Precedence &gt; attrDB.Precedence {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		return preferDA
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if attrDA.Precedence &lt; attrDB.Precedence {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return preferDB
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">// Rule 7: Prefer native transport.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// If DA is reached via an encapsulating transition mechanism (e.g.,</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">// IPv6 in IPv4) and DB is not, then prefer DB.  Similarly, if DB is</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">// reached via encapsulation and DA is not, then prefer DA.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// TODO(bradfitz): implement? low priority for now.</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// Rule 8: Prefer smaller scope.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// If Scope(DA) &lt; Scope(DB), then prefer DA.  Similarly, if Scope(DA) &gt;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// Scope(DB), then prefer DB.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	if attrDA.Scope &lt; attrDB.Scope {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		return preferDA
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	if attrDA.Scope &gt; attrDB.Scope {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		return preferDB
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Rule 9: Use the longest matching prefix.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// When DA and DB belong to the same address family (both are IPv6 or</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// both are IPv4 [but see below]): If CommonPrefixLen(Source(DA), DA) &gt;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">// CommonPrefixLen(Source(DB), DB), then prefer DA.  Similarly, if</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">// CommonPrefixLen(Source(DA), DA) &lt; CommonPrefixLen(Source(DB), DB),</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// then prefer DB.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// However, applying this rule to IPv4 addresses causes</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// problems (see issues 13283 and 18518), so limit to IPv6.</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if DA.To4() == nil &amp;&amp; DB.To4() == nil {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		commonA := commonPrefixLen(SourceDA, DA)
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		commonB := commonPrefixLen(SourceDB, DB)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		if commonA &gt; commonB {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			return preferDA
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		if commonA &lt; commonB {
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			return preferDB
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Rule 10: Otherwise, leave the order unchanged.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// If DA preceded DB in the original list, prefer DA.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">// Otherwise, prefer DB.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	return false <span class="comment">// &#34;equal&#34;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>type policyTableEntry struct {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	Prefix     netip.Prefix
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	Precedence uint8
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	Label      uint8
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>type policyTable []policyTableEntry
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span><span class="comment">// RFC 6724 section 2.1.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span><span class="comment">// Items are sorted by the size of their Prefix.Mask.Size,</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>var rfc6724policyTable = policyTable{
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	{
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		<span class="comment">// &#34;::1/128&#34;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}), 128),
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		Precedence: 50,
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		Label:      0,
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	},
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	{
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		<span class="comment">// &#34;::ffff:0:0/96&#34;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		<span class="comment">// IPv4-compatible, etc.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}), 96),
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		Precedence: 35,
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		Label:      4,
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	},
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	{
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		<span class="comment">// &#34;::/96&#34;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{}), 96),
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		Precedence: 1,
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		Label:      3,
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	},
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	{
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		<span class="comment">// &#34;2001::/32&#34;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		<span class="comment">// Teredo</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0x20, 0x01}), 32),
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		Precedence: 5,
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		Label:      5,
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	},
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	{
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		<span class="comment">// &#34;2002::/16&#34;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		<span class="comment">// 6to4</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0x20, 0x02}), 16),
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		Precedence: 30,
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		Label:      2,
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	},
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	{
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		<span class="comment">// &#34;3ffe::/16&#34;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0x3f, 0xfe}), 16),
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		Precedence: 1,
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		Label:      12,
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	},
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	{
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>		<span class="comment">// &#34;fec0::/10&#34;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0xfe, 0xc0}), 10),
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		Precedence: 1,
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		Label:      11,
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	},
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	{
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		<span class="comment">// &#34;fc00::/7&#34;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{0xfc}), 7),
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		Precedence: 3,
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		Label:      13,
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	},
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	{
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// &#34;::/0&#34;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		Prefix:     netip.PrefixFrom(netip.AddrFrom16([16]byte{}), 0),
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		Precedence: 40,
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		Label:      1,
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	},
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>}
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">// Classify returns the policyTableEntry of the entry with the longest</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// matching prefix that contains ip.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// The table t must be sorted from largest mask size to smallest.</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>func (t policyTable) Classify(ip netip.Addr) policyTableEntry {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// Prefix.Contains() will not match an IPv6 prefix for an IPv4 address.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if ip.Is4() {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		ip = netip.AddrFrom16(ip.As16())
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	for _, ent := range t {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		if ent.Prefix.Contains(ip) {
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>			return ent
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	return policyTableEntry{}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">// RFC 6724 section 3.1.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>type scope uint8
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>const (
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	scopeInterfaceLocal scope = 0x1
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	scopeLinkLocal      scope = 0x2
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	scopeAdminLocal     scope = 0x4
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	scopeSiteLocal      scope = 0x5
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	scopeOrgLocal       scope = 0x8
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	scopeGlobal         scope = 0xe
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>func classifyScope(ip netip.Addr) scope {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		return scopeLinkLocal
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	ipv6 := ip.Is6() &amp;&amp; !ip.Is4In6()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	ipv6AsBytes := ip.As16()
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	if ipv6 &amp;&amp; ip.IsMulticast() {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		return scope(ipv6AsBytes[1] &amp; 0xf)
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	<span class="comment">// Site-local addresses are defined in RFC 3513 section 2.5.6</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// (and deprecated in RFC 3879).</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	if ipv6 &amp;&amp; ipv6AsBytes[0] == 0xfe &amp;&amp; ipv6AsBytes[1]&amp;0xc0 == 0xc0 {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		return scopeSiteLocal
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	return scopeGlobal
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span><span class="comment">// commonPrefixLen reports the length of the longest prefix (looking</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span><span class="comment">// at the most significant, or leftmost, bits) that the</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// two addresses have in common, up to the length of a&#39;s prefix (i.e.,</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">// the portion of the address not including the interface ID).</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// If a or b is an IPv4 address as an IPv6 address, the IPv4 addresses</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// are compared (with max common prefix length of 32).</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// If a and b are different IP versions, 0 is returned.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L342" class="ln">   342&nbsp;&nbsp;</span><span class="comment">// See https://tools.ietf.org/html/rfc6724#section-2.2</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>func commonPrefixLen(a netip.Addr, b IP) (cpl int) {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	if b4 := b.To4(); b4 != nil {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		b = b4
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	aAsSlice := a.AsSlice()
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	if len(aAsSlice) != len(b) {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		return 0
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	<span class="comment">// If IPv6, only up to the prefix (first 64 bits)</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	if len(aAsSlice) &gt; 8 {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		aAsSlice = aAsSlice[:8]
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		b = b[:8]
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	}
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	for len(aAsSlice) &gt; 0 {
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		if aAsSlice[0] == b[0] {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>			cpl += 8
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			aAsSlice = aAsSlice[1:]
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			b = b[1:]
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			continue
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		bits := 8
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		ab, bb := aAsSlice[0], b[0]
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		for {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			ab &gt;&gt;= 1
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			bb &gt;&gt;= 1
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>			bits--
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			if ab == bb {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>				cpl += bits
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>				return
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>			}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		}
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	return
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>}
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
</pre><p><a href="addrselect.go?m=text">View as plain text</a></p>

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
