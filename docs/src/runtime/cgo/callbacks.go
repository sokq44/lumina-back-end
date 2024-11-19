<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/cgo/callbacks.go - Go Documentation Server</title>

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
<a href="callbacks.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/cgo">cgo</a>/<span class="text-muted">callbacks.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/cgo">runtime/cgo</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package cgo
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;unsafe&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// These utility functions are available to be called from code</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// compiled with gcc via crosscall2.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// The declaration of crosscall2 is:</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//   void crosscall2(void (*fn)(void *), void *, int);</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// We need to export the symbol crosscall2 in order to support</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// callbacks from shared libraries. This applies regardless of</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// linking mode.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Compatibility note: SWIG uses crosscall2 in exactly one situation:</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// to call _cgo_panic using the pattern shown below. We need to keep</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// that pattern working. In particular, crosscall2 actually takes four</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// arguments, but it works to call it with three arguments when</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// calling _cgo_panic.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_static crosscall2</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_dynamic crosscall2</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// Panic. The argument is converted into a Go string.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// Call like this in code compiled with gcc:</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//   struct { const char *p; } a;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//   a.p = /* string to pass to panic */;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//   crosscall2(_cgo_panic, &amp;a, sizeof a);</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//   /* The function call will not return.  */</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">// TODO: We should export a regular C function to panic, change SWIG</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">// to use that instead of the above pattern, and then we can drop</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// backwards-compatibility from crosscall2 and stop exporting it.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//go:linkname _runtime_cgo_panic_internal runtime._cgo_panic_internal</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>func _runtime_cgo_panic_internal(p *byte)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_panic _cgo_panic</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_static _cgo_panic</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_dynamic _cgo_panic</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>func _cgo_panic(a *struct{ cstr *byte }) {
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	_runtime_cgo_panic_internal(a.cstr)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_init</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_init x_cgo_init</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_init _cgo_init</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>var x_cgo_init byte
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>var _cgo_init = &amp;x_cgo_init
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_thread_start</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_thread_start x_cgo_thread_start</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_thread_start _cgo_thread_start</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>var x_cgo_thread_start byte
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>var _cgo_thread_start = &amp;x_cgo_thread_start
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// Creates a new system thread without updating any Go state.</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span><span class="comment">// This method is invoked during shared library loading to create a new OS</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span><span class="comment">// thread to perform the runtime initialization. This method is similar to</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span><span class="comment">// _cgo_sys_thread_start except that it doesn&#39;t update any Go state.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_sys_thread_create</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_sys_thread_create x_cgo_sys_thread_create</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_sys_thread_create _cgo_sys_thread_create</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>var x_cgo_sys_thread_create byte
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>var _cgo_sys_thread_create = &amp;x_cgo_sys_thread_create
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">// Indicates whether a dummy thread key has been created or not.</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">// When calling go exported function from C, we register a destructor</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// callback, for a dummy thread key, by using pthread_key_create.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_pthread_key_created</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_pthread_key_created x_cgo_pthread_key_created</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_pthread_key_created _cgo_pthread_key_created</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>var x_cgo_pthread_key_created byte
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>var _cgo_pthread_key_created = &amp;x_cgo_pthread_key_created
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// Export crosscall2 to a c function pointer variable.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span><span class="comment">// Used to dropm in pthread key destructor, while C thread is exiting.</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_crosscall2_ptr</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span><span class="comment">//go:linkname x_crosscall2_ptr x_crosscall2_ptr</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span><span class="comment">//go:linkname _crosscall2_ptr _crosscall2_ptr</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>var x_crosscall2_ptr byte
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>var _crosscall2_ptr = &amp;x_crosscall2_ptr
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// Set the x_crosscall2_ptr C function pointer variable point to crosscall2.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// It&#39;s for the runtime package to call at init time.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>func set_crosscall2()
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">//go:linkname _set_crosscall2 runtime.set_crosscall2</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>var _set_crosscall2 = set_crosscall2
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span><span class="comment">// Store the g into the thread-specific value.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span><span class="comment">// So that pthread_key_destructor will dropm when the thread is exiting.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_bindm</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_bindm x_cgo_bindm</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_bindm _cgo_bindm</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>var x_cgo_bindm byte
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>var _cgo_bindm = &amp;x_cgo_bindm
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// Notifies that the runtime has been initialized.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">// We currently block at every CGO entry point (via _cgo_wait_runtime_init_done)</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// to ensure that the runtime has been initialized before the CGO call is</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// executed. This is necessary for shared libraries where we kickoff runtime</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">// initialization in a separate thread and return without waiting for this</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// thread to complete the init.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_notify_runtime_init_done</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_notify_runtime_init_done x_cgo_notify_runtime_init_done</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_notify_runtime_init_done _cgo_notify_runtime_init_done</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>var x_cgo_notify_runtime_init_done byte
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>var _cgo_notify_runtime_init_done = &amp;x_cgo_notify_runtime_init_done
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Sets the traceback context function. See runtime.SetCgoTraceback.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_set_context_function</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_set_context_function x_cgo_set_context_function</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_set_context_function _cgo_set_context_function</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>var x_cgo_set_context_function byte
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>var _cgo_set_context_function = &amp;x_cgo_set_context_function
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">// Calls a libc function to execute background work injected via libc</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">// interceptors, such as processing pending signals under the thread</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// sanitizer.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// Left as a nil pointer if no libc interceptors are expected.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static _cgo_yield</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_yield _cgo_yield</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>var _cgo_yield unsafe.Pointer
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_static _cgo_topofstack</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//go:cgo_export_dynamic _cgo_topofstack</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// x_cgo_getstackbound gets the thread&#39;s C stack size and</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">// set the G&#39;s stack bound based on the stack size.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">//go:cgo_import_static x_cgo_getstackbound</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">//go:linkname x_cgo_getstackbound x_cgo_getstackbound</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//go:linkname _cgo_getstackbound _cgo_getstackbound</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>var x_cgo_getstackbound byte
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>var _cgo_getstackbound = &amp;x_cgo_getstackbound
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
</pre><p><a href="callbacks.go?m=text">View as plain text</a></p>

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
