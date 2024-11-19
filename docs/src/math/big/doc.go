<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/doc.go - Go Documentation Server</title>

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
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package big implements arbitrary-precision arithmetic (big numbers).
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>The following numeric types are supported:
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	Int    signed integers
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	Rat    rational numbers
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	Float  floating-point numbers
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>The zero value for an [Int], [Rat], or [Float] correspond to 0. Thus, new
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>values can be declared in the usual ways and denote 0 without further
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>initialization:
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	var x Int        // &amp;x is an *Int of value 0
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	var r = &amp;Rat{}   // r is a *Rat of value 0
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	y := new(Float)  // y is a *Float of value 0
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>Alternatively, new values can be allocated and initialized with factory
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>functions of the form:
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	func NewT(v V) *T
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>For instance, [NewInt](x) returns an *[Int] set to the value of the int64
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>argument x, [NewRat](a, b) returns a *[Rat] set to the fraction a/b where
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>a and b are int64 values, and [NewFloat](f) returns a *[Float] initialized
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>to the float64 argument f. More flexibility is provided with explicit
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>setters, for instance:
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	var z1 Int
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	z1.SetUint64(123)                 // z1 := 123
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	z2 := new(Rat).SetFloat64(1.25)   // z2 := 5/4
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	z3 := new(Float).SetInt(z1)       // z3 := 123.0
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>Setters, numeric operations and predicates are represented as methods of
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>the form:
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	func (z *T) SetV(v V) *T          // z = v
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	func (z *T) Unary(x *T) *T        // z = unary x
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	func (z *T) Binary(x, y *T) *T    // z = x binary y
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	func (x *T) Pred() P              // p = pred(x)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>with T one of [Int], [Rat], or [Float]. For unary and binary operations, the
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>result is the receiver (usually named z in that case; see below); if it
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>is one of the operands x or y it may be safely overwritten (and its memory
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>reused).
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>Arithmetic expressions are typically written as a sequence of individual
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>method calls, with each call corresponding to an operation. The receiver
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>denotes the result and the method arguments are the operation&#39;s operands.
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>For instance, given three *Int values a, b and c, the invocation
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	c.Add(a, b)
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>computes the sum a + b and stores the result in c, overwriting whatever
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>value was held in c before. Unless specified otherwise, operations permit
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>aliasing of parameters, so it is perfectly ok to write
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	sum.Add(sum, x)
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>to accumulate values x in a sum.
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>(By always passing in a result value via the receiver, memory use can be
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>much better controlled. Instead of having to allocate new memory for each
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>result, an operation can reuse the space allocated for the result value,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>and overwrite that value with the new result in the process.)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>Notational convention: Incoming method parameters (including the receiver)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>are named consistently in the API to clarify their use. Incoming operands
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>are usually named x, y, a, b, and so on, but never z. A parameter specifying
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>the result is named z (typically the receiver).
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>For instance, the arguments for (*Int).Add are named x and y, and because
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>the receiver specifies the result destination, it is called z:
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	func (z *Int) Add(x, y *Int) *Int
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>Methods of this form typically return the incoming receiver as well, to
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>enable simple call chaining.
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>Methods which don&#39;t require a result value to be passed in (for instance,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>[Int.Sign]), simply return the result. In this case, the receiver is typically
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>the first operand, named x:
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	func (x *Int) Sign() int
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>Various methods support conversions between strings and corresponding
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>numeric values, and vice versa: *[Int], *[Rat], and *[Float] values implement
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>the Stringer interface for a (default) string representation of the value,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>but also provide SetString methods to initialize a value from a string in
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>a variety of supported formats (see the respective SetString documentation).
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>Finally, *[Int], *[Rat], and *[Float] satisfy [fmt.Scanner] for scanning
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>and (except for *[Rat]) the Formatter interface for formatted printing.
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>*/</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>package big
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
