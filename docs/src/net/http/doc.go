<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/http/doc.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/http">http</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/http">net/http</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package http provides HTTP client and server implementations.
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>[Get], [Head], [Post], and [PostForm] make HTTP (or HTTPS) requests:
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	resp, err := http.Get(&#34;http://example.com/&#34;)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	...
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	resp, err := http.Post(&#34;http://example.com/upload&#34;, &#34;image/jpeg&#34;, &amp;buf)
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	...
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	resp, err := http.PostForm(&#34;http://example.com/form&#34;,
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>		url.Values{&#34;key&#34;: {&#34;Value&#34;}, &#34;id&#34;: {&#34;123&#34;}})
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>The caller must close the response body when finished with it:
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	resp, err := http.Get(&#34;http://example.com/&#34;)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	if err != nil {
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>		// handle error
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	}
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	defer resp.Body.Close()
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	body, err := io.ReadAll(resp.Body)
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	// ...
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span># Clients and Transports
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>For control over HTTP client headers, redirect policy, and other
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>settings, create a [Client]:
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	client := &amp;http.Client{
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		CheckRedirect: redirectPolicyFunc,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	}
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	resp, err := client.Get(&#34;http://example.com&#34;)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	// ...
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	req, err := http.NewRequest(&#34;GET&#34;, &#34;http://example.com&#34;, nil)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	// ...
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	req.Header.Add(&#34;If-None-Match&#34;, `W/&#34;wyzzy&#34;`)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	resp, err := client.Do(req)
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	// ...
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>For control over proxies, TLS configuration, keep-alives,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>compression, and other settings, create a [Transport]:
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	tr := &amp;http.Transport{
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		MaxIdleConns:       10,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		IdleConnTimeout:    30 * time.Second,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		DisableCompression: true,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	client := &amp;http.Client{Transport: tr}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	resp, err := client.Get(&#34;https://example.com&#34;)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>Clients and Transports are safe for concurrent use by multiple
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>goroutines and for efficiency should only be created once and re-used.
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span># Servers
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>ListenAndServe starts an HTTP server with a given address and handler.
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>The handler is usually nil, which means to use [DefaultServeMux].
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>[Handle] and [HandleFunc] add handlers to [DefaultServeMux]:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	http.Handle(&#34;/foo&#34;, fooHandler)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	http.HandleFunc(&#34;/bar&#34;, func(w http.ResponseWriter, r *http.Request) {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		fmt.Fprintf(w, &#34;Hello, %q&#34;, html.EscapeString(r.URL.Path))
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	})
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	log.Fatal(http.ListenAndServe(&#34;:8080&#34;, nil))
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>More control over the server&#39;s behavior is available by creating a
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>custom Server:
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	s := &amp;http.Server{
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		Addr:           &#34;:8080&#34;,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		Handler:        myHandler,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		ReadTimeout:    10 * time.Second,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		WriteTimeout:   10 * time.Second,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		MaxHeaderBytes: 1 &lt;&lt; 20,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	log.Fatal(s.ListenAndServe())
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span># HTTP/2
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>Starting with Go 1.6, the http package has transparent support for the
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>HTTP/2 protocol when using HTTPS. Programs that must disable HTTP/2
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>can do so by setting [Transport.TLSNextProto] (for clients) or
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>[Server.TLSNextProto] (for servers) to a non-nil, empty
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>map. Alternatively, the following GODEBUG settings are
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>currently supported:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	GODEBUG=http2client=0  # disable HTTP/2 client support
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	GODEBUG=http2server=0  # disable HTTP/2 server support
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	GODEBUG=http2debug=1   # enable verbose HTTP/2 debug logs
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	GODEBUG=http2debug=2   # ... even more verbose, with frame dumps
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>Please report any issues before disabling HTTP/2 support: https://golang.org/s/http2bug
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>The http package&#39;s [Transport] and [Server] both automatically enable
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>HTTP/2 support for simple configurations. To enable HTTP/2 for more
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>complex configurations, to use lower-level HTTP/2 features, or to use
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>a newer version of Go&#39;s http2 package, import &#34;golang.org/x/net/http2&#34;
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>directly and use its ConfigureTransport and/or ConfigureServer
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>functions. Manually configuring HTTP/2 via the golang.org/x/net/http2
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>package takes precedence over the net/http package&#39;s built-in HTTP/2
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>support.
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>*/</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>package http
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>
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
