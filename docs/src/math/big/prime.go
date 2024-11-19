<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/prime.go - Go Documentation Server</title>

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
<a href="prime.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">prime.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2016 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package big
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;math/rand&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// ProbablyPrime reports whether x is probably prime,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// applying the Miller-Rabin test with n pseudorandomly chosen bases</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// as well as a Baillie-PSW test.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// If x is prime, ProbablyPrime returns true.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// If x is chosen randomly and not prime, ProbablyPrime probably returns false.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// The probability of returning true for a randomly chosen non-prime is at most ¼ⁿ.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// ProbablyPrime is 100% accurate for inputs less than 2⁶⁴.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// See Menezes et al., Handbook of Applied Cryptography, 1997, pp. 145-149,</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// and FIPS 186-4 Appendix F for further discussion of the error probabilities.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// ProbablyPrime is not suitable for judging primes that an adversary may</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// have crafted to fool the test.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// As of Go 1.8, ProbablyPrime(0) is allowed and applies only a Baillie-PSW test.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Before Go 1.8, ProbablyPrime applied only the Miller-Rabin tests, and ProbablyPrime(0) panicked.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>func (x *Int) ProbablyPrime(n int) bool {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// Note regarding the doc comment above:</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// It would be more precise to say that the Baillie-PSW test uses the</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// extra strong Lucas test as its Lucas test, but since no one knows</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// how to tell any of the Lucas tests apart inside a Baillie-PSW test</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// (they all work equally well empirically), that detail need not be</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// documented or implicitly guaranteed.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// The comment does avoid saying &#34;the&#34; Baillie-PSW test</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// because of this general ambiguity.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	if n &lt; 0 {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		panic(&#34;negative n for ProbablyPrime&#34;)
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	}
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	if x.neg || len(x.abs) == 0 {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		return false
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	<span class="comment">// primeBitMask records the primes &lt; 64.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	const primeBitMask uint64 = 1&lt;&lt;2 | 1&lt;&lt;3 | 1&lt;&lt;5 | 1&lt;&lt;7 |
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		1&lt;&lt;11 | 1&lt;&lt;13 | 1&lt;&lt;17 | 1&lt;&lt;19 | 1&lt;&lt;23 | 1&lt;&lt;29 | 1&lt;&lt;31 |
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		1&lt;&lt;37 | 1&lt;&lt;41 | 1&lt;&lt;43 | 1&lt;&lt;47 | 1&lt;&lt;53 | 1&lt;&lt;59 | 1&lt;&lt;61
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	w := x.abs[0]
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	if len(x.abs) == 1 &amp;&amp; w &lt; 64 {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		return primeBitMask&amp;(1&lt;&lt;w) != 0
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if w&amp;1 == 0 {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return false <span class="comment">// x is even</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	const primesA = 3 * 5 * 7 * 11 * 13 * 17 * 19 * 23 * 37
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	const primesB = 29 * 31 * 41 * 43 * 47 * 53
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	var rA, rB uint32
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	switch _W {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	case 32:
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		rA = uint32(x.abs.modW(primesA))
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		rB = uint32(x.abs.modW(primesB))
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	case 64:
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		r := x.abs.modW((primesA * primesB) &amp; _M)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		rA = uint32(r % primesA)
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		rB = uint32(r % primesB)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	default:
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		panic(&#34;math/big: invalid word size&#34;)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	if rA%3 == 0 || rA%5 == 0 || rA%7 == 0 || rA%11 == 0 || rA%13 == 0 || rA%17 == 0 || rA%19 == 0 || rA%23 == 0 || rA%37 == 0 ||
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		rB%29 == 0 || rB%31 == 0 || rB%41 == 0 || rB%43 == 0 || rB%47 == 0 || rB%53 == 0 {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		return false
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	return x.abs.probablyPrimeMillerRabin(n+1, true) &amp;&amp; x.abs.probablyPrimeLucas()
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// probablyPrimeMillerRabin reports whether n passes reps rounds of the</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// Miller-Rabin primality test, using pseudo-randomly chosen bases.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// If force2 is true, one of the rounds is forced to use base 2.</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// See Handbook of Applied Cryptography, p. 139, Algorithm 4.24.</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// The number n is known to be non-zero.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>func (n nat) probablyPrimeMillerRabin(reps int, force2 bool) bool {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	nm1 := nat(nil).sub(n, natOne)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// determine q, k such that nm1 = q &lt;&lt; k</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	k := nm1.trailingZeroBits()
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	q := nat(nil).shr(nm1, k)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	nm3 := nat(nil).sub(nm1, natTwo)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	rand := rand.New(rand.NewSource(int64(n[0])))
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	var x, y, quotient nat
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	nm3Len := nm3.bitLen()
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>NextRandom:
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	for i := 0; i &lt; reps; i++ {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		if i == reps-1 &amp;&amp; force2 {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			x = x.set(natTwo)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		} else {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			x = x.random(rand, nm3, nm3Len)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			x = x.add(x, natTwo)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		y = y.expNN(x, q, n, false)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if y.cmp(natOne) == 0 || y.cmp(nm1) == 0 {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			continue
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		for j := uint(1); j &lt; k; j++ {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			y = y.sqr(y)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			quotient, y = quotient.div(y, y, n)
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			if y.cmp(nm1) == 0 {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				continue NextRandom
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			if y.cmp(natOne) == 0 {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>				return false
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		return false
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	return true
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// probablyPrimeLucas reports whether n passes the &#34;almost extra strong&#34; Lucas probable prime test,</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// using Baillie-OEIS parameter selection. This corresponds to &#34;AESLPSP&#34; on Jacobsen&#39;s tables (link below).</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// The combination of this test and a Miller-Rabin/Fermat test with base 2 gives a Baillie-PSW test.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// References:</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// Baillie and Wagstaff, &#34;Lucas Pseudoprimes&#34;, Mathematics of Computation 35(152),</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// October 1980, pp. 1391-1417, especially page 1401.</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// https://www.ams.org/journals/mcom/1980-35-152/S0025-5718-1980-0583518-6/S0025-5718-1980-0583518-6.pdf</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// Grantham, &#34;Frobenius Pseudoprimes&#34;, Mathematics of Computation 70(234),</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">// March 2000, pp. 873-891.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// https://www.ams.org/journals/mcom/2001-70-234/S0025-5718-00-01197-2/S0025-5718-00-01197-2.pdf</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">// Baillie, &#34;Extra strong Lucas pseudoprimes&#34;, OEIS A217719, https://oeis.org/A217719.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">// Jacobsen, &#34;Pseudoprime Statistics, Tables, and Data&#34;, http://ntheory.org/pseudoprimes.html.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// Nicely, &#34;The Baillie-PSW Primality Test&#34;, https://web.archive.org/web/20191121062007/http://www.trnicely.net/misc/bpsw.html.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// (Note that Nicely&#39;s definition of the &#34;extra strong&#34; test gives the wrong Jacobi condition,</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// as pointed out by Jacobsen.)</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">// Crandall and Pomerance, Prime Numbers: A Computational Perspective, 2nd ed.</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">// Springer, 2005.</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>func (n nat) probablyPrimeLucas() bool {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">// Discard 0, 1.</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if len(n) == 0 || n.cmp(natOne) == 0 {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return false
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">// Two is the only even prime.</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Already checked by caller, but here to allow testing in isolation.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if n[0]&amp;1 == 0 {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return n.cmp(natTwo) == 0
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Baillie-OEIS &#34;method C&#34; for choosing D, P, Q,</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// as in https://oeis.org/A217719/a217719.txt:</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// try increasing P ≥ 3 such that D = P² - 4 (so Q = 1)</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// until Jacobi(D, n) = -1.</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">// The search is expected to succeed for non-square n after just a few trials.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// After more than expected failures, check whether n is square</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// (which would cause Jacobi(D, n) = 1 for all D not dividing n).</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	p := Word(3)
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	d := nat{1}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	t1 := nat(nil) <span class="comment">// temp</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	intD := &amp;Int{abs: d}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	intN := &amp;Int{abs: n}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	for ; ; p++ {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		if p &gt; 10000 {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			<span class="comment">// This is widely believed to be impossible.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			<span class="comment">// If we get a report, we&#39;ll want the exact number n.</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			panic(&#34;math/big: internal error: cannot find (D/n) = -1 for &#34; + intN.String())
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		}
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		d[0] = p*p - 4
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		j := Jacobi(intD, intN)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		if j == -1 {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			break
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		if j == 0 {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			<span class="comment">// d = p²-4 = (p-2)(p+2).</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			<span class="comment">// If (d/n) == 0 then d shares a prime factor with n.</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			<span class="comment">// Since the loop proceeds in increasing p and starts with p-2==1,</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>			<span class="comment">// the shared prime factor must be p+2.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			<span class="comment">// If p+2 == n, then n is prime; otherwise p+2 is a proper factor of n.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			return len(n) == 1 &amp;&amp; n[0] == p+2
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		if p == 40 {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			<span class="comment">// We&#39;ll never find (d/n) = -1 if n is a square.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			<span class="comment">// If n is a non-square we expect to find a d in just a few attempts on average.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			<span class="comment">// After 40 attempts, take a moment to check if n is indeed a square.</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			t1 = t1.sqrt(n)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>			t1 = t1.sqr(t1)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>			if t1.cmp(n) == 0 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>				return false
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	}
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Grantham definition of &#34;extra strong Lucas pseudoprime&#34;, after Thm 2.3 on p. 876</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	<span class="comment">// (D, P, Q above have become Δ, b, 1):</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// Let U_n = U_n(b, 1), V_n = V_n(b, 1), and Δ = b²-4.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	<span class="comment">// An extra strong Lucas pseudoprime to base b is a composite n = 2^r s + Jacobi(Δ, n),</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	<span class="comment">// where s is odd and gcd(n, 2*Δ) = 1, such that either (i) U_s ≡ 0 mod n and V_s ≡ ±2 mod n,</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	<span class="comment">// or (ii) V_{2^t s} ≡ 0 mod n for some 0 ≤ t &lt; r-1.</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// We know gcd(n, Δ) = 1 or else we&#39;d have found Jacobi(d, n) == 0 above.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	<span class="comment">// We know gcd(n, 2) = 1 because n is odd.</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	<span class="comment">// Arrange s = (n - Jacobi(Δ, n)) / 2^r = (n+1) / 2^r.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	s := nat(nil).add(n, natOne)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	r := int(s.trailingZeroBits())
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	s = s.shr(s, uint(r))
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	nm2 := nat(nil).sub(n, natTwo) <span class="comment">// n-2</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	<span class="comment">// We apply the &#34;almost extra strong&#34; test, which checks the above conditions</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	<span class="comment">// except for U_s ≡ 0 mod n, which allows us to avoid computing any U_k values.</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	<span class="comment">// Jacobsen points out that maybe we should just do the full extra strong test:</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	<span class="comment">// &#34;It is also possible to recover U_n using Crandall and Pomerance equation 3.13:</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	<span class="comment">// U_n = D^-1 (2V_{n+1} - PV_n) allowing us to run the full extra-strong test</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	<span class="comment">// at the cost of a single modular inversion. This computation is easy and fast in GMP,</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	<span class="comment">// so we can get the full extra-strong test at essentially the same performance as the</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	<span class="comment">// almost extra strong test.&#34;</span>
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	<span class="comment">// Compute Lucas sequence V_s(b, 1), where:</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	<span class="comment">//	V(0) = 2</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	<span class="comment">//	V(1) = P</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	<span class="comment">//	V(k) = P V(k-1) - Q V(k-2).</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	<span class="comment">// (Remember that due to method C above, P = b, Q = 1.)</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// In general V(k) = α^k + β^k, where α and β are roots of x² - Px + Q.</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// Crandall and Pomerance (p.147) observe that for 0 ≤ j ≤ k,</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	<span class="comment">//	V(j+k) = V(j)V(k) - V(k-j).</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// So in particular, to quickly double the subscript:</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	<span class="comment">//	V(2k) = V(k)² - 2</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	<span class="comment">//	V(2k+1) = V(k) V(k+1) - P</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	<span class="comment">// We can therefore start with k=0 and build up to k=s in log₂(s) steps.</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	natP := nat(nil).setWord(p)
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	vk := nat(nil).setWord(2)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	vk1 := nat(nil).setWord(p)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	t2 := nat(nil) <span class="comment">// temp</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	for i := int(s.bitLen()); i &gt;= 0; i-- {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		if s.bit(uint(i)) != 0 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			<span class="comment">// k&#39; = 2k+1</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			<span class="comment">// V(k&#39;) = V(2k+1) = V(k) V(k+1) - P.</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>			t1 = t1.mul(vk, vk1)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>			t1 = t1.add(t1, n)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			t1 = t1.sub(t1, natP)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>			t2, vk = t2.div(vk, t1, n)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>			<span class="comment">// V(k&#39;+1) = V(2k+2) = V(k+1)² - 2.</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			t1 = t1.sqr(vk1)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>			t1 = t1.add(t1, nm2)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			t2, vk1 = t2.div(vk1, t1, n)
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		} else {
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>			<span class="comment">// k&#39; = 2k</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>			<span class="comment">// V(k&#39;+1) = V(2k+1) = V(k) V(k+1) - P.</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>			t1 = t1.mul(vk, vk1)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			t1 = t1.add(t1, n)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>			t1 = t1.sub(t1, natP)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			t2, vk1 = t2.div(vk1, t1, n)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>			<span class="comment">// V(k&#39;) = V(2k) = V(k)² - 2</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			t1 = t1.sqr(vk)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			t1 = t1.add(t1, nm2)
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			t2, vk = t2.div(vk, t1, n)
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	}
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// Now k=s, so vk = V(s). Check V(s) ≡ ±2 (mod n).</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	if vk.cmp(natTwo) == 0 || vk.cmp(nm2) == 0 {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>		<span class="comment">// Check U(s) ≡ 0.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		<span class="comment">// As suggested by Jacobsen, apply Crandall and Pomerance equation 3.13:</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		<span class="comment">//	U(k) = D⁻¹ (2 V(k+1) - P V(k))</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		<span class="comment">// Since we are checking for U(k) == 0 it suffices to check 2 V(k+1) == P V(k) mod n,</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>		<span class="comment">// or P V(k) - 2 V(k+1) == 0 mod n.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		t1 := t1.mul(vk, natP)
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		t2 := t2.shl(vk1, 1)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		if t1.cmp(t2) &lt; 0 {
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			t1, t2 = t2, t1
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>		}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		t1 = t1.sub(t1, t2)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		t3 := vk1 <span class="comment">// steal vk1, no longer needed below</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		vk1 = nil
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		_ = vk1
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		t2, t3 = t2.div(t3, t1, n)
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		if len(t3) == 0 {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			return true
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	<span class="comment">// Check V(2^t s) ≡ 0 mod n for some 0 ≤ t &lt; r-1.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	for t := 0; t &lt; r-1; t++ {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		if len(vk) == 0 { <span class="comment">// vk == 0</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>			return true
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>		}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		<span class="comment">// Optimization: V(k) = 2 is a fixed point for V(k&#39;) = V(k)² - 2,</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		<span class="comment">// so if V(k) = 2, we can stop: we will never find a future V(k) == 0.</span>
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		if len(vk) == 1 &amp;&amp; vk[0] == 2 { <span class="comment">// vk == 2</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			return false
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>		}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		<span class="comment">// k&#39; = 2k</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		<span class="comment">// V(k&#39;) = V(2k) = V(k)² - 2</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>		t1 = t1.sqr(vk)
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		t1 = t1.sub(t1, natTwo)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		t2, vk = t2.div(vk, t1, n)
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	return false
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
</pre><p><a href="prime.go?m=text">View as plain text</a></p>

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
