<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/errors/wrap.go - Go Documentation Server</title>

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
<a href="wrap.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/errors">errors</a>/<span class="text-muted">wrap.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/errors">errors</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2018 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package errors
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;internal/reflectlite&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// Unwrap returns the result of calling the Unwrap method on err, if err&#39;s</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// type contains an Unwrap method returning error.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// Otherwise, Unwrap returns nil.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// Unwrap only calls a method of the form &#34;Unwrap() error&#34;.</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// In particular Unwrap does not unwrap errors returned by [Join].</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>func Unwrap(err error) error {
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	u, ok := err.(interface {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>		Unwrap() error
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	})
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	if !ok {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>		return nil
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	}
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	return u.Unwrap()
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>}
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// Is reports whether any error in err&#39;s tree matches target.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// The tree consists of err itself, followed by the errors obtained by repeatedly</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// calling its Unwrap() error or Unwrap() []error method. When err wraps multiple</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// errors, Is examines err followed by a depth-first traversal of its children.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// An error is considered to match a target if it is equal to that target or if</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// it implements a method Is(error) bool such that Is(target) returns true.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// An error type might provide an Is method so it can be treated as equivalent</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// to an existing error. For example, if MyError defines</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// then Is(MyError{}, fs.ErrExist) returns true. See [syscall.Errno.Is] for</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">// an example in the standard library. An Is method should only shallowly</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// compare err and the target and not call [Unwrap] on either.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>func Is(err, target error) bool {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	if target == nil {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		return err == target
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	isComparable := reflectlite.TypeOf(target).Comparable()
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	return is(err, target, isComparable)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>func is(err, target error, targetComparable bool) bool {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	for {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		if targetComparable &amp;&amp; err == target {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			return true
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		if x, ok := err.(interface{ Is(error) bool }); ok &amp;&amp; x.Is(target) {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			return true
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		switch x := err.(type) {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		case interface{ Unwrap() error }:
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			err = x.Unwrap()
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			if err == nil {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>				return false
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		case interface{ Unwrap() []error }:
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			for _, err := range x.Unwrap() {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>				if is(err, target, targetComparable) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>					return true
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>				}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			return false
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		default:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return false
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">// As finds the first error in err&#39;s tree that matches target, and if one is found, sets</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// target to that error value and returns true. Otherwise, it returns false.</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// The tree consists of err itself, followed by the errors obtained by repeatedly</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// calling its Unwrap() error or Unwrap() []error method. When err wraps multiple</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// errors, As examines err followed by a depth-first traversal of its children.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span><span class="comment">// An error matches target if the error&#39;s concrete value is assignable to the value</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">// pointed to by target, or if the error has a method As(interface{}) bool such that</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">// As(target) returns true. In the latter case, the As method is responsible for</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">// setting target.</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// An error type might provide an As method so it can be treated as if it were a</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// different error type.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// As panics if target is not a non-nil pointer to either a type that implements</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span><span class="comment">// error, or to any interface type.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>func As(err error, target any) bool {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	if err == nil {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		return false
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	if target == nil {
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		panic(&#34;errors: target cannot be nil&#34;)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	}
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	val := reflectlite.ValueOf(target)
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	typ := val.Type()
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	if typ.Kind() != reflectlite.Ptr || val.IsNil() {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		panic(&#34;errors: target must be a non-nil pointer&#34;)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	targetType := typ.Elem()
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	if targetType.Kind() != reflectlite.Interface &amp;&amp; !targetType.Implements(errorType) {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		panic(&#34;errors: *target must be interface or implement error&#34;)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	return as(err, target, val, targetType)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>func as(err error, target any, targetVal reflectlite.Value, targetType reflectlite.Type) bool {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	for {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		if reflectlite.TypeOf(err).AssignableTo(targetType) {
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			targetVal.Elem().Set(reflectlite.ValueOf(err))
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			return true
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		if x, ok := err.(interface{ As(any) bool }); ok &amp;&amp; x.As(target) {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			return true
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		switch x := err.(type) {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		case interface{ Unwrap() error }:
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>			err = x.Unwrap()
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			if err == nil {
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				return false
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		case interface{ Unwrap() []error }:
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			for _, err := range x.Unwrap() {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				if err == nil {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>					continue
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>				}
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>				if as(err, target, targetVal, targetType) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>					return true
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>				}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>			return false
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		default:
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			return false
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>var errorType = reflectlite.TypeOf((*error)(nil)).Elem()
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
</pre><p><a href="wrap.go?m=text">View as plain text</a></p>

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
