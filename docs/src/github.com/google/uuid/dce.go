<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/google/uuid/dce.go - Go Documentation Server</title>

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
<a href="dce.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/google">google</a>/<a href="http://localhost:8080/src/github.com/google/uuid">uuid</a>/<span class="text-muted">dce.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/github.com/google/uuid">github.com/google/uuid</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 Google Inc.  All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package uuid
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;encoding/binary&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// A Domain represents a Version 2 domain</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>type Domain byte
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Domain constants for DCE Security (Version 2) UUIDs.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>const (
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	Person = Domain(0)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	Group  = Domain(1)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	Org    = Domain(2)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// NewDCESecurity returns a DCE Security (Version 2) UUID.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// The domain should be one of Person, Group or Org.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// On a POSIX system the id should be the users UID for the Person</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// domain and the users GID for the Group.  The meaning of id for</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// the domain Org or on non-POSIX systems is site defined.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// For a given domain/id pair the same token may be returned for up to</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// 7 minutes and 10 seconds.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func NewDCESecurity(domain Domain, id uint32) (UUID, error) {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	uuid, err := NewUUID()
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if err == nil {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		uuid[6] = (uuid[6] &amp; 0x0f) | 0x20 <span class="comment">// Version 2</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		uuid[9] = byte(domain)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		binary.BigEndian.PutUint32(uuid[0:], id)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	return uuid, err
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// NewDCEPerson returns a DCE Security (Version 2) UUID in the person</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// domain with the id returned by os.Getuid.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//  NewDCESecurity(Person, uint32(os.Getuid()))</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func NewDCEPerson() (UUID, error) {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	return NewDCESecurity(Person, uint32(os.Getuid()))
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// NewDCEGroup returns a DCE Security (Version 2) UUID in the group</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// domain with the id returned by os.Getgid.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//  NewDCESecurity(Group, uint32(os.Getgid()))</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func NewDCEGroup() (UUID, error) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	return NewDCESecurity(Group, uint32(os.Getgid()))
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>}
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// Domain returns the domain for a Version 2 UUID.  Domains are only defined</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// for Version 2 UUIDs.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func (uuid UUID) Domain() Domain {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	return Domain(uuid[9])
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// ID returns the id for a Version 2 UUID. IDs are only defined for Version 2</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// UUIDs.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>func (uuid UUID) ID() uint32 {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	return binary.BigEndian.Uint32(uuid[0:4])
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>func (d Domain) String() string {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	switch d {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	case Person:
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		return &#34;Person&#34;
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	case Group:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		return &#34;Group&#34;
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	case Org:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return &#34;Org&#34;
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;Domain%d&#34;, int(d))
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
</pre><p><a href="dce.go?m=text">View as plain text</a></p>

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
