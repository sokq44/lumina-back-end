<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/http/servemux121.go - Go Documentation Server</title>

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
<a href="servemux121.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/http">http</a>/<span class="text-muted">servemux121.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/http">net/http</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package http
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// This file implements ServeMux behavior as in Go 1.21.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">// The behavior is controlled by a GODEBUG setting.</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// Most of this code is derived from commit 08e35cc334.</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// Changes are minimal: aside from the different receiver type,</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// they mostly involve renaming functions, usually by unexporting them.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>import (
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;internal/godebug&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;net/url&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>)
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>var httpmuxgo121 = godebug.New(&#34;httpmuxgo121&#34;)
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>var use121 bool
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Read httpmuxgo121 once at startup, since dealing with changes to it during</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// program execution is too complex and error-prone.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>func init() {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	if httpmuxgo121.Value() == &#34;1&#34; {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		use121 = true
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		httpmuxgo121.IncNonDefault()
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>}
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// serveMux121 holds the state of a ServeMux needed for Go 1.21 behavior.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>type serveMux121 struct {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	mu    sync.RWMutex
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	m     map[string]muxEntry
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	es    []muxEntry <span class="comment">// slice of entries sorted from longest to shortest.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	hosts bool       <span class="comment">// whether any patterns contain hostnames</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>type muxEntry struct {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	h       Handler
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	pattern string
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// Formerly ServeMux.Handle.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>func (mux *serveMux121) handle(pattern string, handler Handler) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	mux.mu.Lock()
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	defer mux.mu.Unlock()
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	if pattern == &#34;&#34; {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		panic(&#34;http: invalid pattern&#34;)
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if handler == nil {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		panic(&#34;http: nil handler&#34;)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	if _, exist := mux.m[pattern]; exist {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		panic(&#34;http: multiple registrations for &#34; + pattern)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	if mux.m == nil {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		mux.m = make(map[string]muxEntry)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	e := muxEntry{h: handler, pattern: pattern}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	mux.m[pattern] = e
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	if pattern[len(pattern)-1] == &#39;/&#39; {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		mux.es = appendSorted(mux.es, e)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if pattern[0] != &#39;/&#39; {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		mux.hosts = true
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	n := len(es)
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	i := sort.Search(n, func(i int) bool {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		return len(es[i].pattern) &lt; len(e.pattern)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	})
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if i == n {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return append(es, e)
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// we now know that i points at where we want to insert</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	es = append(es, muxEntry{}) <span class="comment">// try to grow the slice in place, any entry works.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	copy(es[i+1:], es[i:])      <span class="comment">// Move shorter entries down</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	es[i] = e
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return es
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span><span class="comment">// Formerly ServeMux.HandleFunc.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>func (mux *serveMux121) handleFunc(pattern string, handler func(ResponseWriter, *Request)) {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if handler == nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		panic(&#34;http: nil handler&#34;)
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	mux.handle(pattern, HandlerFunc(handler))
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span><span class="comment">// Formerly ServeMux.Handler.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>func (mux *serveMux121) findHandler(r *Request) (h Handler, pattern string) {
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	<span class="comment">// CONNECT requests are not canonicalized.</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	if r.Method == &#34;CONNECT&#34; {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		<span class="comment">// If r.URL.Path is /tree and its handler is not registered,</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		<span class="comment">// the /tree -&gt; /tree/ redirect applies to CONNECT requests</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		<span class="comment">// but the path canonicalization does not.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if u, ok := mux.redirectToPathSlash(r.URL.Host, r.URL.Path, r.URL); ok {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		return mux.handler(r.Host, r.URL.Path)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	<span class="comment">// All other requests have any port stripped and path cleaned</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// before passing to mux.handler.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	host := stripHostPort(r.Host)
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	path := cleanPath(r.URL.Path)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// If the given path is /tree and its handler is not registered,</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	<span class="comment">// redirect for /tree/.</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if u, ok := mux.redirectToPathSlash(host, path, r.URL); ok {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	if path != r.URL.Path {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		_, pattern = mux.handler(host, path)
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		u := &amp;url.URL{Path: path, RawQuery: r.URL.RawQuery}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		return RedirectHandler(u.String(), StatusMovedPermanently), pattern
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	return mux.handler(host, r.URL.Path)
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">// handler is the main implementation of findHandler.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">// The path is known to be in canonical form, except for CONNECT methods.</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>func (mux *serveMux121) handler(host, path string) (h Handler, pattern string) {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	mux.mu.RLock()
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	defer mux.mu.RUnlock()
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">// Host-specific pattern takes precedence over generic ones</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if mux.hosts {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		h, pattern = mux.match(host + path)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if h == nil {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		h, pattern = mux.match(path)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	}
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if h == nil {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		h, pattern = NotFoundHandler(), &#34;&#34;
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	return
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// Find a handler on a handler map given a path string.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span><span class="comment">// Most-specific (longest) pattern wins.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func (mux *serveMux121) match(path string) (h Handler, pattern string) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Check for exact match first.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	v, ok := mux.m[path]
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	if ok {
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		return v.h, v.pattern
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	}
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// Check for longest valid match.  mux.es contains all patterns</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// that end in / sorted from longest to shortest.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	for _, e := range mux.es {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		if strings.HasPrefix(path, e.pattern) {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			return e.h, e.pattern
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	return nil, &#34;&#34;
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// redirectToPathSlash determines if the given path needs appending &#34;/&#34; to it.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// This occurs when a handler for path + &#34;/&#34; was already registered, but</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span><span class="comment">// not for path itself. If the path needs appending to, it creates a new</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// URL, setting the path to u.Path + &#34;/&#34; and returning true to indicate so.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>func (mux *serveMux121) redirectToPathSlash(host, path string, u *url.URL) (*url.URL, bool) {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	mux.mu.RLock()
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	shouldRedirect := mux.shouldRedirectRLocked(host, path)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	mux.mu.RUnlock()
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if !shouldRedirect {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		return u, false
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	path = path + &#34;/&#34;
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	u = &amp;url.URL{Path: path, RawQuery: u.RawQuery}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	return u, true
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// shouldRedirectRLocked reports whether the given path and host should be redirected to</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// path+&#34;/&#34;. This should happen if a handler is registered for path+&#34;/&#34; but</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span><span class="comment">// not path -- see comments at ServeMux.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>func (mux *serveMux121) shouldRedirectRLocked(host, path string) bool {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	p := []string{path, host + path}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	for _, c := range p {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		if _, exist := mux.m[c]; exist {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>			return false
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		}
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	n := len(path)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if n == 0 {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return false
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	for _, c := range p {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		if _, exist := mux.m[c+&#34;/&#34;]; exist {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>			return path[n-1] != &#39;/&#39;
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		}
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	return false
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
</pre><p><a href="servemux121.go?m=text">View as plain text</a></p>

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
