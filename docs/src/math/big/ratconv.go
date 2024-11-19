<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/ratconv.go - Go Documentation Server</title>

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
<a href="ratconv.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">ratconv.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This file implements rat-to-string conversion functions.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package big
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>func ratTok(ch rune) bool {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	return strings.ContainsRune(&#34;+-/0123456789.eE&#34;, ch)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>}
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>var ratZero Rat
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>var _ fmt.Scanner = &amp;ratZero <span class="comment">// *Rat must implement fmt.Scanner</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// Scan is a support routine for fmt.Scanner. It accepts the formats</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// &#39;e&#39;, &#39;E&#39;, &#39;f&#39;, &#39;F&#39;, &#39;g&#39;, &#39;G&#39;, and &#39;v&#39;. All formats are equivalent.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func (z *Rat) Scan(s fmt.ScanState, ch rune) error {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	tok, err := s.Token(true, ratTok)
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	if err != nil {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		return err
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	if !strings.ContainsRune(&#34;efgEFGv&#34;, ch) {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		return errors.New(&#34;Rat.Scan: invalid verb&#34;)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	}
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	if _, ok := z.SetString(string(tok)); !ok {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		return errors.New(&#34;Rat.Scan: invalid syntax&#34;)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	return nil
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// SetString sets z to the value of s and returns z and a boolean indicating</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// success. s can be given as a (possibly signed) fraction &#34;a/b&#34;, or as a</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// floating-point number optionally followed by an exponent.</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// If a fraction is provided, both the dividend and the divisor may be a</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// decimal integer or independently use a prefix of “0b”, “0” or “0o”,</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// or “0x” (or their upper-case variants) to denote a binary, octal, or</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// hexadecimal integer, respectively. The divisor may not be signed.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// If a floating-point number is provided, it may be in decimal form or</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// use any of the same prefixes as above but for “0” to denote a non-decimal</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// mantissa. A leading “0” is considered a decimal leading 0; it does not</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// indicate octal representation in this case.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// An optional base-10 “e” or base-2 “p” (or their upper-case variants)</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// exponent may be provided as well, except for hexadecimal floats which</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// only accept an (optional) “p” exponent (because an “e” or “E” cannot</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// be distinguished from a mantissa digit). If the exponent&#39;s absolute value</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// is too large, the operation may fail.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// The entire string, not just a prefix, must be valid for success. If the</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// operation failed, the value of z is undefined but the returned value is nil.</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func (z *Rat) SetString(s string) (*Rat, bool) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if len(s) == 0 {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return nil, false
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	<span class="comment">// len(s) &gt; 0</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	<span class="comment">// parse fraction a/b, if any</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if sep := strings.Index(s, &#34;/&#34;); sep &gt;= 0 {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		if _, ok := z.a.SetString(s[:sep], 0); !ok {
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			return nil, false
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		r := strings.NewReader(s[sep+1:])
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		var err error
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if z.b.abs, _, _, err = z.b.abs.scan(r, 0, false); err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return nil, false
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		<span class="comment">// entire string must have been consumed</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		if _, err = r.ReadByte(); err != io.EOF {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			return nil, false
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		if len(z.b.abs) == 0 {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			return nil, false
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		return z.norm(), true
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// parse floating-point number</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	r := strings.NewReader(s)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// sign</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	neg, err := scanSign(r)
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if err != nil {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		return nil, false
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// mantissa</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	var base int
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	var fcount int <span class="comment">// fractional digit count; valid if &lt;= 0</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	z.a.abs, base, fcount, err = z.a.abs.scan(r, 0, true)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	if err != nil {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		return nil, false
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// exponent</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	var exp int64
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	var ebase int
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	exp, ebase, err = scanExponent(r, true, true)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	if err != nil {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		return nil, false
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	<span class="comment">// there should be no unread characters left</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if _, err = r.ReadByte(); err != io.EOF {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return nil, false
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// special-case 0 (see also issue #16176)</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if len(z.a.abs) == 0 {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		return z.norm(), true
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// len(z.a.abs) &gt; 0</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// The mantissa may have a radix point (fcount &lt;= 0) and there</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	<span class="comment">// may be a nonzero exponent exp. The radix point amounts to a</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	<span class="comment">// division by base**(-fcount), which equals a multiplication by</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">// base**fcount. An exponent means multiplication by ebase**exp.</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">// Multiplications are commutative, so we can apply them in any</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">// order. We only have powers of 2 and 10, and we split powers</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// of 10 into the product of the same powers of 2 and 5. This</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// may reduce the size of shift/multiplication factors or</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">// divisors required to create the final fraction, depending</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">// on the actual floating-point value.</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">// determine binary or decimal exponent contribution of radix point</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	var exp2, exp5 int64
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if fcount &lt; 0 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		<span class="comment">// The mantissa has a radix point ddd.dddd; and</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		<span class="comment">// -fcount is the number of digits to the right</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		<span class="comment">// of &#39;.&#39;. Adjust relevant exponent accordingly.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		d := int64(fcount)
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		switch base {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		case 10:
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			exp5 = d
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			fallthrough <span class="comment">// 10**e == 5**e * 2**e</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		case 2:
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>			exp2 = d
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		case 8:
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			exp2 = d * 3 <span class="comment">// octal digits are 3 bits each</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		case 16:
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			exp2 = d * 4 <span class="comment">// hexadecimal digits are 4 bits each</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		default:
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>			panic(&#34;unexpected mantissa base&#34;)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		<span class="comment">// fcount consumed - not needed anymore</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">// take actual exponent into account</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	switch ebase {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	case 10:
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		exp5 += exp
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		fallthrough <span class="comment">// see fallthrough above</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	case 2:
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		exp2 += exp
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	default:
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		panic(&#34;unexpected exponent base&#34;)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// exp consumed - not needed anymore</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// apply exp5 contributions</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// (start with exp5 so the numbers to multiply are smaller)</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	if exp5 != 0 {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		n := exp5
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		if n &lt; 0 {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			n = -n
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>			if n &lt; 0 {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>				<span class="comment">// This can occur if -n overflows. -(-1 &lt;&lt; 63) would become</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>				<span class="comment">// -1 &lt;&lt; 63, which is still negative.</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>				return nil, false
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		if n &gt; 1e6 {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			return nil, false <span class="comment">// avoid excessively large exponents</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		pow5 := z.b.abs.expNN(natFive, nat(nil).setWord(Word(n)), nil, false) <span class="comment">// use underlying array of z.b.abs</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		if exp5 &gt; 0 {
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			z.a.abs = z.a.abs.mul(z.a.abs, pow5)
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			z.b.abs = z.b.abs.setWord(1)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		} else {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			z.b.abs = pow5
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	} else {
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		z.b.abs = z.b.abs.setWord(1)
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// apply exp2 contributions</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	if exp2 &lt; -1e7 || exp2 &gt; 1e7 {
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		return nil, false <span class="comment">// avoid excessively large exponents</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	}
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if exp2 &gt; 0 {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		z.a.abs = z.a.abs.shl(z.a.abs, uint(exp2))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	} else if exp2 &lt; 0 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		z.b.abs = z.b.abs.shl(z.b.abs, uint(-exp2))
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	z.a.neg = neg &amp;&amp; len(z.a.abs) &gt; 0 <span class="comment">// 0 has no sign</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	return z.norm(), true
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">// scanExponent scans the longest possible prefix of r representing a base 10</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">// (“e”, “E”) or a base 2 (“p”, “P”) exponent, if any. It returns the</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">// exponent, the exponent base (10 or 2), or a read or syntax error, if any.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">// If sepOk is set, an underscore character “_” may appear between successive</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">// exponent digits; such underscores do not change the value of the exponent.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">// Incorrect placement of underscores is reported as an error if there are no</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">// other errors. If sepOk is not set, underscores are not recognized and thus</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// terminate scanning like any other character that is not a valid digit.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">//	exponent = ( &#34;e&#34; | &#34;E&#34; | &#34;p&#34; | &#34;P&#34; ) [ sign ] digits .</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">//	sign     = &#34;+&#34; | &#34;-&#34; .</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">//	digits   = digit { [ &#39;_&#39; ] digit } .</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">//	digit    = &#34;0&#34; ... &#34;9&#34; .</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span><span class="comment">// A base 2 exponent is only permitted if base2ok is set.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>func scanExponent(r io.ByteScanner, base2ok, sepOk bool) (exp int64, base int, err error) {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// one char look-ahead</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	ch, err := r.ReadByte()
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	if err != nil {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		if err == io.EOF {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			err = nil
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		return 0, 10, err
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">// exponent char</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	switch ch {
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	case &#39;e&#39;, &#39;E&#39;:
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		base = 10
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	case &#39;p&#39;, &#39;P&#39;:
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>		if base2ok {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			base = 2
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>			break <span class="comment">// ok</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		fallthrough <span class="comment">// binary exponent not permitted</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	default:
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>		r.UnreadByte() <span class="comment">// ch does not belong to exponent anymore</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		return 0, 10, nil
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	}
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// sign</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	var digits []byte
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	ch, err = r.ReadByte()
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if err == nil &amp;&amp; (ch == &#39;+&#39; || ch == &#39;-&#39;) {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		if ch == &#39;-&#39; {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			digits = append(digits, &#39;-&#39;)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		}
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		ch, err = r.ReadByte()
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	<span class="comment">// prev encodes the previously seen char: it is one</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// of &#39;_&#39;, &#39;0&#39; (a digit), or &#39;.&#39; (anything else). A</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	<span class="comment">// valid separator &#39;_&#39; may only occur after a digit.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	prev := &#39;.&#39;
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	invalSep := false
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// exponent value</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	hasDigits := false
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	for err == nil {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		if &#39;0&#39; &lt;= ch &amp;&amp; ch &lt;= &#39;9&#39; {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			digits = append(digits, ch)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			prev = &#39;0&#39;
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			hasDigits = true
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		} else if ch == &#39;_&#39; &amp;&amp; sepOk {
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			if prev != &#39;0&#39; {
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>				invalSep = true
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			prev = &#39;_&#39;
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		} else {
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			r.UnreadByte() <span class="comment">// ch does not belong to number anymore</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			break
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>		ch, err = r.ReadByte()
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	}
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	if err == io.EOF {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		err = nil
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>	if err == nil &amp;&amp; !hasDigits {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		err = errNoDigits
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	if err == nil {
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		exp, err = strconv.ParseInt(string(digits), 10, 64)
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	<span class="comment">// other errors take precedence over invalid separators</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if err == nil &amp;&amp; (invalSep || prev == &#39;_&#39;) {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		err = errInvalSep
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	return
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// String returns a string representation of x in the form &#34;a/b&#34; (even if b == 1).</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>func (x *Rat) String() string {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	return string(x.marshal())
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span><span class="comment">// marshal implements String returning a slice of bytes</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (x *Rat) marshal() []byte {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	var buf []byte
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	buf = x.a.Append(buf, 10)
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	buf = append(buf, &#39;/&#39;)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	if len(x.b.abs) != 0 {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		buf = x.b.Append(buf, 10)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	} else {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		buf = append(buf, &#39;1&#39;)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	return buf
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span><span class="comment">// RatString returns a string representation of x in the form &#34;a/b&#34; if b != 1,</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span><span class="comment">// and in the form &#34;a&#34; if b == 1.</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>func (x *Rat) RatString() string {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	if x.IsInt() {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		return x.a.String()
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	return x.String()
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span><span class="comment">// FloatString returns a string representation of x in decimal form with prec</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span><span class="comment">// digits of precision after the radix point. The last digit is rounded to</span>
<span id="L329" class="ln">   329&nbsp;&nbsp;</span><span class="comment">// nearest, with halves rounded away from zero.</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>func (x *Rat) FloatString(prec int) string {
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	var buf []byte
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	if x.IsInt() {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		buf = x.a.Append(buf, 10)
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		if prec &gt; 0 {
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			buf = append(buf, &#39;.&#39;)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			for i := prec; i &gt; 0; i-- {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				buf = append(buf, &#39;0&#39;)
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>		}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		return string(buf)
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	<span class="comment">// x.b.abs != 0</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	q, r := nat(nil).div(nat(nil), x.a.abs, x.b.abs)
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	p := natOne
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	if prec &gt; 0 {
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		p = nat(nil).expNN(natTen, nat(nil).setUint64(uint64(prec)), nil, false)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	}
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	r = r.mul(r, p)
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	r, r2 := r.div(nat(nil), r, x.b.abs)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	<span class="comment">// see if we need to round up</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	r2 = r2.add(r2, r2)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	if x.b.abs.cmp(r2) &lt;= 0 {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		r = r.add(r, natOne)
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		if r.cmp(p) &gt;= 0 {
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			q = nat(nil).add(q, natOne)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			r = nat(nil).sub(r, p)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		}
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	if x.a.neg {
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		buf = append(buf, &#39;-&#39;)
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	buf = append(buf, q.utoa(10)...) <span class="comment">// itoa ignores sign if q == 0</span>
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	if prec &gt; 0 {
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		buf = append(buf, &#39;.&#39;)
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>		rs := r.utoa(10)
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		for i := prec - len(rs); i &gt; 0; i-- {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>			buf = append(buf, &#39;0&#39;)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		}
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		buf = append(buf, rs...)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	return string(buf)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>}
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>
<span id="L382" class="ln">   382&nbsp;&nbsp;</span><span class="comment">// Note: FloatPrec (below) is in this file rather than rat.go because</span>
<span id="L383" class="ln">   383&nbsp;&nbsp;</span><span class="comment">//       its results are relevant for decimal representation/printing.</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span><span class="comment">// FloatPrec returns the number n of non-repeating digits immediately</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// following the decimal point of the decimal representation of x.</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// The boolean result indicates whether a decimal representation of x</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// with that many fractional digits is exact or rounded.</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// Examples:</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">//	x      n    exact    decimal representation n fractional digits</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">//	0      0    true     0</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">//	1      0    true     1</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">//	1/2    1    true     0.5</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">//	1/3    0    false    0       (0.333... rounded)</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span><span class="comment">//	1/4    2    true     0.25</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span><span class="comment">//	1/6    1    false    0.2     (0.166... rounded)</span>
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>func (x *Rat) FloatPrec() (n int, exact bool) {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>	<span class="comment">// Determine q and largest p2, p5 such that d = q·2^p2·5^p5.</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	<span class="comment">// The results n, exact are:</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	<span class="comment">//     n = max(p2, p5)</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	<span class="comment">//     exact = q == 1</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	<span class="comment">// For details see:</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	<span class="comment">// https://en.wikipedia.org/wiki/Repeating_decimal#Reciprocals_of_integers_not_coprime_to_10</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	d := x.Denom().abs <span class="comment">// d &gt;= 1</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	<span class="comment">// Determine p2 by counting factors of 2.</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	<span class="comment">// p2 corresponds to the trailing zero bits in d.</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	<span class="comment">// Do this first to reduce q as much as possible.</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>	var q nat
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	p2 := d.trailingZeroBits()
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	q = q.shr(d, p2)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	<span class="comment">// Determine p5 by counting factors of 5.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	<span class="comment">// Build a table starting with an initial power of 5,</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	<span class="comment">// and use repeated squaring until the factor doesn&#39;t</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	<span class="comment">// divide q anymore. Then use the table to determine</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	<span class="comment">// the power of 5 in q.</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	const fp = 13        <span class="comment">// f == 5^fp</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	var tab []nat        <span class="comment">// tab[i] == (5^fp)^(2^i) == 5^(fp·2^i)</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	f := nat{1220703125} <span class="comment">// == 5^fp (must fit into a uint32 Word)</span>
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	var t, r nat         <span class="comment">// temporaries</span>
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	for {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		if _, r = t.div(r, q, f); len(r) != 0 {
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			break <span class="comment">// f doesn&#39;t divide q evenly</span>
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>		}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>		tab = append(tab, f)
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>		f = nat(nil).sqr(f) <span class="comment">// nat(nil) to ensure a new f for each table entry</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	}
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	<span class="comment">// Factor q using the table entries, if any.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	<span class="comment">// We start with the largest factor f = tab[len(tab)-1]</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	<span class="comment">// that evenly divides q. It does so at most once because</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	<span class="comment">// otherwise f·f would also divide q. That can&#39;t be true</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>	<span class="comment">// because f·f is the next higher table entry, contradicting</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	<span class="comment">// how f was chosen in the first place.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>	<span class="comment">// The same reasoning applies to the subsequent factors.</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	var p5 uint
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>	for i := len(tab) - 1; i &gt;= 0; i-- {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		if t, r = t.div(r, q, tab[i]); len(r) == 0 {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			p5 += fp * (1 &lt;&lt; i) <span class="comment">// tab[i] == 5^(fp·2^i)</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			q = q.set(t)
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	}
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	<span class="comment">// If fp != 1, we may still have multiples of 5 left.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	for {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		if t, r = t.div(r, q, natFive); len(r) != 0 {
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			break
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>		p5++
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>		q = q.set(t)
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>	}
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	return int(max(p2, p5)), q.cmp(natOne) == 0
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>}
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>
</pre><p><a href="ratconv.go?m=text">View as plain text</a></p>

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
