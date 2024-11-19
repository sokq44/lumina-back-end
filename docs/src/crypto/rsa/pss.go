<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/crypto/rsa/pss.go - Go Documentation Server</title>

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
<a href="pss.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/crypto">crypto</a>/<a href="http://localhost:8080/src/crypto/rsa">rsa</a>/<span class="text-muted">pss.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/crypto/rsa">crypto/rsa</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package rsa
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file implements the RSASSA-PSS signature scheme according to RFC 8017.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;crypto/internal/boring&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;hash&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// Per RFC 8017, Section 9.1</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//     EM = MGF1 xor DB || H( 8*0x00 || mHash || salt ) || 0xbc</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// where</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//     DB = PS || 0x01 || salt</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// and PS can be empty so</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//     emLen = dbLen + hLen + 1 = psLen + sLen + hLen + 2</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func emsaPSSEncode(mHash []byte, emBits int, salt []byte, hash hash.Hash) ([]byte, error) {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// See RFC 8017, Section 9.1.1.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	hLen := hash.Size()
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	sLen := len(salt)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	emLen := (emBits + 7) / 8
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	<span class="comment">// 1.  If the length of M is greater than the input limitation for the</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">//     hash function (2^61 - 1 octets for SHA-1), output &#34;message too</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">//     long&#34; and stop.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	<span class="comment">// 2.  Let mHash = Hash(M), an octet string of length hLen.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	if len(mHash) != hLen {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		return nil, errors.New(&#34;crypto/rsa: input must be hashed with given hash&#34;)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// 3.  If emLen &lt; hLen + sLen + 2, output &#34;encoding error&#34; and stop.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	if emLen &lt; hLen+sLen+2 {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		return nil, ErrMessageTooLong
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	em := make([]byte, emLen)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	psLen := emLen - sLen - hLen - 2
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	db := em[:psLen+1+sLen]
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	h := em[psLen+1+sLen : emLen-1]
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// 4.  Generate a random octet string salt of length sLen; if sLen = 0,</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	<span class="comment">//     then salt is the empty string.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// 5.  Let</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	<span class="comment">//       M&#39; = (0x)00 00 00 00 00 00 00 00 || mHash || salt;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">//     M&#39; is an octet string of length 8 + hLen + sLen with eight</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	<span class="comment">//     initial zero octets.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// 6.  Let H = Hash(M&#39;), an octet string of length hLen.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	var prefix [8]byte
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	hash.Write(prefix[:])
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	hash.Write(mHash)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	hash.Write(salt)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	h = hash.Sum(h[:0])
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	hash.Reset()
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// 7.  Generate an octet string PS consisting of emLen - sLen - hLen - 2</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	<span class="comment">//     zero octets. The length of PS may be 0.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// 8.  Let DB = PS || 0x01 || salt; DB is an octet string of length</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">//     emLen - hLen - 1.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	db[psLen] = 0x01
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	copy(db[psLen+1:], salt)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// 9.  Let dbMask = MGF(H, emLen - hLen - 1).</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// 10. Let maskedDB = DB \xor dbMask.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	mgf1XOR(db, hash, h)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// 11. Set the leftmost 8 * emLen - emBits bits of the leftmost octet in</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">//     maskedDB to zero.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	db[0] &amp;= 0xff &gt;&gt; (8*emLen - emBits)
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// 12. Let EM = maskedDB || H || 0xbc.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	em[emLen-1] = 0xbc
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// 13. Output EM.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	return em, nil
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>func emsaPSSVerify(mHash, em []byte, emBits, sLen int, hash hash.Hash) error {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	<span class="comment">// See RFC 8017, Section 9.1.2.</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	hLen := hash.Size()
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if sLen == PSSSaltLengthEqualsHash {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		sLen = hLen
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	emLen := (emBits + 7) / 8
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	if emLen != len(em) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		return errors.New(&#34;rsa: internal error: inconsistent length&#34;)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// 1.  If the length of M is greater than the input limitation for the</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">//     hash function (2^61 - 1 octets for SHA-1), output &#34;inconsistent&#34;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">//     and stop.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// 2.  Let mHash = Hash(M), an octet string of length hLen.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if hLen != len(mHash) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return ErrVerification
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// 3.  If emLen &lt; hLen + sLen + 2, output &#34;inconsistent&#34; and stop.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	if emLen &lt; hLen+sLen+2 {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		return ErrVerification
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// 4.  If the rightmost octet of EM does not have hexadecimal value</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//     0xbc, output &#34;inconsistent&#34; and stop.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	if em[emLen-1] != 0xbc {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return ErrVerification
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// 5.  Let maskedDB be the leftmost emLen - hLen - 1 octets of EM, and</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">//     let H be the next hLen octets.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	db := em[:emLen-hLen-1]
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	h := em[emLen-hLen-1 : emLen-1]
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">// 6.  If the leftmost 8 * emLen - emBits bits of the leftmost octet in</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">//     maskedDB are not all equal to zero, output &#34;inconsistent&#34; and</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">//     stop.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	var bitMask byte = 0xff &gt;&gt; (8*emLen - emBits)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if em[0] &amp; ^bitMask != 0 {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		return ErrVerification
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// 7.  Let dbMask = MGF(H, emLen - hLen - 1).</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">// 8.  Let DB = maskedDB \xor dbMask.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	mgf1XOR(db, hash, h)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// 9.  Set the leftmost 8 * emLen - emBits bits of the leftmost octet in DB</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">//     to zero.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	db[0] &amp;= bitMask
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// If we don&#39;t know the salt length, look for the 0x01 delimiter.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	if sLen == PSSSaltLengthAuto {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		psLen := bytes.IndexByte(db, 0x01)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		if psLen &lt; 0 {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>			return ErrVerification
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		sLen = len(db) - psLen - 1
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">// 10. If the emLen - hLen - sLen - 2 leftmost octets of DB are not zero</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">//     or if the octet at position emLen - hLen - sLen - 1 (the leftmost</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">//     position is &#34;position 1&#34;) does not have hexadecimal value 0x01,</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">//     output &#34;inconsistent&#34; and stop.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	psLen := emLen - hLen - sLen - 2
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	for _, e := range db[:psLen] {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		if e != 0x00 {
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			return ErrVerification
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	if db[psLen] != 0x01 {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return ErrVerification
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// 11.  Let salt be the last sLen octets of DB.</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	salt := db[len(db)-sLen:]
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">// 12.  Let</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//          M&#39; = (0x)00 00 00 00 00 00 00 00 || mHash || salt ;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//     M&#39; is an octet string of length 8 + hLen + sLen with eight</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">//     initial zero octets.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// 13. Let H&#39; = Hash(M&#39;), an octet string of length hLen.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	var prefix [8]byte
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	hash.Write(prefix[:])
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	hash.Write(mHash)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	hash.Write(salt)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	h0 := hash.Sum(nil)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	<span class="comment">// 14. If H = H&#39;, output &#34;consistent.&#34; Otherwise, output &#34;inconsistent.&#34;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	if !bytes.Equal(h0, h) { <span class="comment">// TODO: constant time?</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		return ErrVerification
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	return nil
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">// signPSSWithSalt calculates the signature of hashed using PSS with specified salt.</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// Note that hashed must be the result of hashing the input message using the</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// given hash function. salt is a random sequence of bytes whose length will be</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// later used to verify the signature.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>func signPSSWithSalt(priv *PrivateKey, hash crypto.Hash, hashed, salt []byte) ([]byte, error) {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	emBits := priv.N.BitLen() - 1
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	em, err := emsaPSSEncode(hashed, emBits, salt, hash.New())
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if err != nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return nil, err
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	if boring.Enabled {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		bkey, err := boringPrivateKey(priv)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if err != nil {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			return nil, err
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		<span class="comment">// Note: BoringCrypto always does decrypt &#34;withCheck&#34;.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		<span class="comment">// (It&#39;s not just decrypt.)</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		s, err := boring.DecryptRSANoPadding(bkey, em)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		if err != nil {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			return nil, err
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		return s, nil
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// RFC 8017: &#34;Note that the octet length of EM will be one less than k if</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">// modBits - 1 is divisible by 8 and equal to k otherwise, where k is the</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// length in octets of the RSA modulus n.&#34; 🙄</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// This is extremely annoying, as all other encrypt and decrypt inputs are</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// always the exact same size as the modulus. Since it only happens for</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// weird modulus sizes, fix it by padding inefficiently.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	if emLen, k := len(em), priv.Size(); emLen &lt; k {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		emNew := make([]byte, k)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		copy(emNew[k-emLen:], em)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		em = emNew
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	return decrypt(priv, em, withCheck)
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>}
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>const (
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// PSSSaltLengthAuto causes the salt in a PSS signature to be as large</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	<span class="comment">// as possible when signing, and to be auto-detected when verifying.</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	PSSSaltLengthAuto = 0
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// PSSSaltLengthEqualsHash causes the salt length to equal the length</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	<span class="comment">// of the hash used in the signature.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	PSSSaltLengthEqualsHash = -1
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>)
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span><span class="comment">// PSSOptions contains options for creating and verifying PSS signatures.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>type PSSOptions struct {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// SaltLength controls the length of the salt used in the PSS signature. It</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// can either be a positive number of bytes, or one of the special</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// PSSSaltLength constants.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	SaltLength int
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// Hash is the hash function used to generate the message digest. If not</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// zero, it overrides the hash function passed to SignPSS. It&#39;s required</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// when using PrivateKey.Sign.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	Hash crypto.Hash
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span><span class="comment">// HashFunc returns opts.Hash so that [PSSOptions] implements [crypto.SignerOpts].</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>func (opts *PSSOptions) HashFunc() crypto.Hash {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	return opts.Hash
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>func (opts *PSSOptions) saltLength() int {
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	if opts == nil {
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		return PSSSaltLengthAuto
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return opts.SaltLength
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>var invalidSaltLenErr = errors.New(&#34;crypto/rsa: PSSOptions.SaltLength cannot be negative&#34;)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// SignPSS calculates the signature of digest using PSS.</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// digest must be the result of hashing the input message using the given hash</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// function. The opts argument may be nil, in which case sensible defaults are</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span><span class="comment">// used. If opts.Hash is set, it overrides hash.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span><span class="comment">// The signature is randomized depending on the message, key, and salt size,</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span><span class="comment">// using bytes from rand. Most applications should use [crypto/rand.Reader] as</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// rand.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>func SignPSS(rand io.Reader, priv *PrivateKey, hash crypto.Hash, digest []byte, opts *PSSOptions) ([]byte, error) {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	<span class="comment">// Note that while we don&#39;t commit to deterministic execution with respect</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// to the rand stream, we also don&#39;t apply MaybeReadByte, so per Hyrum&#39;s Law</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	<span class="comment">// it&#39;s probably relied upon by some. It&#39;s a tolerable promise because a</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	<span class="comment">// well-specified number of random bytes is included in the signature, in a</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// well-specified way.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	if boring.Enabled &amp;&amp; rand == boring.RandReader {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		bkey, err := boringPrivateKey(priv)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		if err != nil {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			return nil, err
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>		return boring.SignRSAPSS(bkey, hash, digest, opts.saltLength())
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	}
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	boring.UnreachableExceptTests()
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	if opts != nil &amp;&amp; opts.Hash != 0 {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		hash = opts.Hash
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	}
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	saltLength := opts.saltLength()
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	switch saltLength {
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	case PSSSaltLengthAuto:
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		saltLength = (priv.N.BitLen()-1+7)/8 - 2 - hash.Size()
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		if saltLength &lt; 0 {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>			return nil, ErrMessageTooLong
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	case PSSSaltLengthEqualsHash:
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		saltLength = hash.Size()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	default:
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		<span class="comment">// If we get here saltLength is either &gt; 0 or &lt; -1, in the</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		<span class="comment">// latter case we fail out.</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		if saltLength &lt;= 0 {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>			return nil, invalidSaltLenErr
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	}
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	salt := make([]byte, saltLength)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>	if _, err := io.ReadFull(rand, salt); err != nil {
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		return nil, err
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	}
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	return signPSSWithSalt(priv, hash, digest, salt)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>}
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>
<span id="L335" class="ln">   335&nbsp;&nbsp;</span><span class="comment">// VerifyPSS verifies a PSS signature.</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span><span class="comment">// A valid signature is indicated by returning a nil error. digest must be the</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span><span class="comment">// result of hashing the input message using the given hash function. The opts</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// argument may be nil, in which case sensible defaults are used. opts.Hash is</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span><span class="comment">// ignored.</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>func VerifyPSS(pub *PublicKey, hash crypto.Hash, digest []byte, sig []byte, opts *PSSOptions) error {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	if boring.Enabled {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		bkey, err := boringPublicKey(pub)
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		if err != nil {
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>			return err
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		if err := boring.VerifyRSAPSS(bkey, hash, digest, sig, opts.saltLength()); err != nil {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>			return ErrVerification
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>		return nil
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>	}
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	if len(sig) != pub.Size() {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		return ErrVerification
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// Salt length must be either one of the special constants (-1 or 0)</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">// or otherwise positive. If it is &lt; PSSSaltLengthEqualsHash (-1)</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	<span class="comment">// we return an error.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	if opts.saltLength() &lt; PSSSaltLengthEqualsHash {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		return invalidSaltLenErr
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	}
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	emBits := pub.N.BitLen() - 1
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	emLen := (emBits + 7) / 8
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	em, err := encrypt(pub, sig)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	if err != nil {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		return ErrVerification
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	<span class="comment">// Like in signPSSWithSalt, deal with mismatches between emLen and the size</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// of the modulus. The spec would have us wire emLen into the encoding</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	<span class="comment">// function, but we&#39;d rather always encode to the size of the modulus and</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	<span class="comment">// then strip leading zeroes if necessary. This only happens for weird</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	<span class="comment">// modulus sizes anyway.</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	for len(em) &gt; emLen &amp;&amp; len(em) &gt; 0 {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		if em[0] != 0 {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>			return ErrVerification
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		em = em[1:]
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	}
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	return emsaPSSVerify(digest, em, emBits, opts.saltLength(), hash.New())
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>
</pre><p><a href="pss.go?m=text">View as plain text</a></p>

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
