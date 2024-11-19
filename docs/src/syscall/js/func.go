<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/syscall/js/func.go - Go Documentation Server</title>

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
<a href="func.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/syscall">syscall</a>/<a href="http://localhost:8080/src/syscall/js">js</a>/<span class="text-muted">func.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/syscall/js">syscall/js</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2018 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build js &amp;&amp; wasm</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package js
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import &#34;sync&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>var (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	funcsMu    sync.Mutex
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	funcs             = make(map[uint32]func(Value, []Value) any)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	nextFuncID uint32 = 1
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Func is a wrapped Go function to be called by JavaScript.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>type Func struct {
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	Value <span class="comment">// the JavaScript function that invokes the Go function</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	id    uint32
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>}
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// FuncOf returns a function to be used by JavaScript.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// The Go function fn is called with the value of JavaScript&#39;s &#34;this&#34; keyword and the</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// arguments of the invocation. The return value of the invocation is</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// the result of the Go function mapped back to JavaScript according to ValueOf.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// Invoking the wrapped Go function from JavaScript will</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// pause the event loop and spawn a new goroutine.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// Other wrapped functions which are triggered during a call from Go to JavaScript</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">// get executed on the same goroutine.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// As a consequence, if one wrapped function blocks, JavaScript&#39;s event loop</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// is blocked until that function returns. Hence, calling any async JavaScript</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// API, which requires the event loop, like fetch (http.Client), will cause an</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// immediate deadlock. Therefore a blocking function should explicitly start a</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// new goroutine.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// Func.Release must be called to free up resources when the function will not be invoked any more.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func FuncOf(fn func(this Value, args []Value) any) Func {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	funcsMu.Lock()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	id := nextFuncID
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	nextFuncID++
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	funcs[id] = fn
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	funcsMu.Unlock()
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	return Func{
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		id:    id,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		Value: jsGo.Call(&#34;_makeFuncWrapper&#34;, id),
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// Release frees up resources allocated for the function.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// The function must not be invoked after calling Release.</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// It is allowed to call Release while the function is still running.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>func (c Func) Release() {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	funcsMu.Lock()
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	delete(funcs, c.id)
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	funcsMu.Unlock()
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// setEventHandler is defined in the runtime package.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>func setEventHandler(fn func() bool)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>func init() {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	setEventHandler(handleEvent)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">// handleEvent retrieves the pending event (window._pendingEvent) and calls the js.Func on it.</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// It returns true if an event was handled.</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>func handleEvent() bool {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	<span class="comment">// Retrieve the event from js</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	cb := jsGo.Get(&#34;_pendingEvent&#34;)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if cb.IsNull() {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		return false
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	jsGo.Set(&#34;_pendingEvent&#34;, Null())
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	id := uint32(cb.Get(&#34;id&#34;).Int())
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if id == 0 { <span class="comment">// zero indicates deadlock</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		select {}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// Retrieve the associated js.Func</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	funcsMu.Lock()
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	f, ok := funcs[id]
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	funcsMu.Unlock()
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	if !ok {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		Global().Get(&#34;console&#34;).Call(&#34;error&#34;, &#34;call to released function&#34;)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		return true
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// Call the js.Func with arguments</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	this := cb.Get(&#34;this&#34;)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	argsObj := cb.Get(&#34;args&#34;)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	args := make([]Value, argsObj.Length())
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	for i := range args {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		args[i] = argsObj.Index(i)
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	result := f(this, args)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// Return the result to js</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	cb.Set(&#34;result&#34;, result)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	return true
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
</pre><p><a href="func.go?m=text">View as plain text</a></p>

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
