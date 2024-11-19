<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/cmplx/sqrt.go - Go Documentation Server</title>

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
<a href="sqrt.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/cmplx">cmplx</a>/<span class="text-muted">sqrt.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/cmplx">math/cmplx</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package cmplx
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;math&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// The original C code, the long comment, and the constants</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// below are from http://netlib.sandia.gov/cephes/c9x-complex/clog.c.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// The go code is a simplified version of the original C.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Cephes Math Library Release 2.8:  June, 2000</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// The readme file at http://netlib.sandia.gov/cephes/ says:</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//    Some software in this archive may be from the book _Methods and</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// Programs for Mathematical Functions_ (Prentice-Hall or Simon &amp; Schuster</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// International, 1989) or from the Cephes Mathematical Library, a</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// commercial product. In either event, it is copyrighted by the author.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// What you see here may be used freely but it comes with no support or</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// guarantee.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//   The two known misprints in the book are repaired here in the</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// source listings for the gamma function and the incomplete beta</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// integral.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//   Stephen L. Moshier</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//   moshier@na-net.ornl.gov</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// Complex square root</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// DESCRIPTION:</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// If z = x + iy,  r = |z|, then</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//                       1/2</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// Re w  =  [ (r + x)/2 ]   ,</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//                       1/2</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// Im w  =  [ (r - x)/2 ]   .</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// Cancellation error in r-x or r+x is avoided by using the</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// identity  2 Re w Im w  =  y.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// Note that -w is also a square root of z. The root chosen</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// is always in the right half plane and Im w has the same sign as y.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// ACCURACY:</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//                      Relative error:</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// arithmetic   domain     # trials      peak         rms</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//    DEC       -10,+10     25000       3.2e-17     9.6e-18</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//    IEEE      -10,+10   1,000,000     2.9e-16     6.1e-17</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// Sqrt returns the square root of x.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// The result r is chosen so that real(r) ≥ 0 and imag(r) has the same sign as imag(x).</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>func Sqrt(x complex128) complex128 {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if imag(x) == 0 {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		<span class="comment">// Ensure that imag(r) has the same sign as imag(x) for imag(x) == signed zero.</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if real(x) == 0 {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			return complex(0, imag(x))
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if real(x) &lt; 0 {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			return complex(0, math.Copysign(math.Sqrt(-real(x)), imag(x)))
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		return complex(math.Sqrt(real(x)), imag(x))
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	} else if math.IsInf(imag(x), 0) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return complex(math.Inf(1.0), imag(x))
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if real(x) == 0 {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		if imag(x) &lt; 0 {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			r := math.Sqrt(-0.5 * imag(x))
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>			return complex(r, -r)
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		r := math.Sqrt(0.5 * imag(x))
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		return complex(r, r)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	a := real(x)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	b := imag(x)
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	var scale float64
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Rescale to avoid internal overflow or underflow.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if math.Abs(a) &gt; 4 || math.Abs(b) &gt; 4 {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		a *= 0.25
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		b *= 0.25
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		scale = 2
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	} else {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		a *= 1.8014398509481984e16 <span class="comment">// 2**54</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		b *= 1.8014398509481984e16
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		scale = 7.450580596923828125e-9 <span class="comment">// 2**-27</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	r := math.Hypot(a, b)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	var t float64
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	if a &gt; 0 {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		t = math.Sqrt(0.5*r + 0.5*a)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		r = scale * math.Abs((0.5*b)/t)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		t *= scale
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	} else {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		r = math.Sqrt(0.5*r - 0.5*a)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		t = scale * math.Abs((0.5*b)/r)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		r *= scale
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	}
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if b &lt; 0 {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		return complex(t, -r)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	return complex(t, r)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
</pre><p><a href="sqrt.go?m=text">View as plain text</a></p>

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
