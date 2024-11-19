<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/http/httptest/server.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="server.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<a href="http://localhost:8080/src/net/http">http</a>/<a href="http://localhost:8080/src/net/http/httptest">httptest</a>/<span class="text-muted">server.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net/http/httptest">net/http/httptest</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Implementation of Server</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>package httptest
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>import (
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;crypto/tls&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;crypto/x509&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;flag&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;log&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;net&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;net/http&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;net/http/internal/testcert&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// A Server is an HTTP server listening on a system-chosen port on the</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// local loopback interface, for use in end-to-end HTTP tests.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type Server struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	URL      string <span class="comment">// base URL of form http://ipaddr:port with no trailing slash</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	Listener net.Listener
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// EnableHTTP2 controls whether HTTP/2 is enabled</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	<span class="comment">// on the server. It must be set between calling</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// NewUnstartedServer and calling Server.StartTLS.</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	EnableHTTP2 bool
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	<span class="comment">// TLS is the optional TLS configuration, populated with a new config</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// after TLS is started. If set on an unstarted server before StartTLS</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// is called, existing fields are copied into the new config.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	TLS *tls.Config
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// Config may be changed after calling NewUnstartedServer and</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// before Start or StartTLS.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	Config *http.Server
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	<span class="comment">// certificate is a parsed version of the TLS config certificate, if present.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	certificate *x509.Certificate
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	<span class="comment">// wg counts the number of outstanding HTTP requests on this server.</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	<span class="comment">// Close blocks until all requests are finished.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	wg sync.WaitGroup
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	mu     sync.Mutex <span class="comment">// guards closed and conns</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	closed bool
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	conns  map[net.Conn]http.ConnState <span class="comment">// except terminal states</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	<span class="comment">// client is configured for use with the server.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Its transport is automatically closed when Close is called.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	client *http.Client
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>func newLocalListener() net.Listener {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if serveFlag != &#34;&#34; {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		l, err := net.Listen(&#34;tcp&#34;, serveFlag)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		if err != nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>			panic(fmt.Sprintf(&#34;httptest: failed to listen on %v: %v&#34;, serveFlag, err))
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		return l
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	l, err := net.Listen(&#34;tcp&#34;, &#34;127.0.0.1:0&#34;)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if err != nil {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		if l, err = net.Listen(&#34;tcp6&#34;, &#34;[::1]:0&#34;); err != nil {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>			panic(fmt.Sprintf(&#34;httptest: failed to listen on a port: %v&#34;, err))
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	return l
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span><span class="comment">// When debugging a particular http server-based test,</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// this flag lets you run</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span><span class="comment">//	go test -run=&#39;^BrokenTest$&#39; -httptest.serve=127.0.0.1:8000</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// to start the broken server so you can interact with it manually.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span><span class="comment">// We only register this flag if it looks like the caller knows about it</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span><span class="comment">// and is trying to use it as we don&#39;t want to pollute flags and this</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span><span class="comment">// isn&#39;t really part of our API. Don&#39;t depend on this.</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>var serveFlag string
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>func init() {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if strSliceContainsPrefix(os.Args, &#34;-httptest.serve=&#34;) || strSliceContainsPrefix(os.Args, &#34;--httptest.serve=&#34;) {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		flag.StringVar(&amp;serveFlag, &#34;httptest.serve&#34;, &#34;&#34;, &#34;if non-empty, httptest.NewServer serves on this address and blocks.&#34;)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>func strSliceContainsPrefix(v []string, pre string) bool {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	for _, s := range v {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if strings.HasPrefix(s, pre) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return true
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	return false
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span><span class="comment">// NewServer starts and returns a new [Server].</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span><span class="comment">// The caller should call Close when finished, to shut it down.</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>func NewServer(handler http.Handler) *Server {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	ts := NewUnstartedServer(handler)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	ts.Start()
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	return ts
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span><span class="comment">// NewUnstartedServer returns a new [Server] but doesn&#39;t start it.</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span><span class="comment">// After changing its configuration, the caller should call Start or</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span><span class="comment">// StartTLS.</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span><span class="comment">// The caller should call Close when finished, to shut it down.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>func NewUnstartedServer(handler http.Handler) *Server {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	return &amp;Server{
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		Listener: newLocalListener(),
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		Config:   &amp;http.Server{Handler: handler},
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span><span class="comment">// Start starts a server from NewUnstartedServer.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>func (s *Server) Start() {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	if s.URL != &#34;&#34; {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		panic(&#34;Server already started&#34;)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	}
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	if s.client == nil {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		s.client = &amp;http.Client{Transport: &amp;http.Transport{}}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	s.URL = &#34;http://&#34; + s.Listener.Addr().String()
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	s.wrap()
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	s.goServe()
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if serveFlag != &#34;&#34; {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		fmt.Fprintln(os.Stderr, &#34;httptest: serving on&#34;, s.URL)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		select {}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>}
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">// StartTLS starts TLS on a server from NewUnstartedServer.</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>func (s *Server) StartTLS() {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if s.URL != &#34;&#34; {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		panic(&#34;Server already started&#34;)
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	if s.client == nil {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		s.client = &amp;http.Client{}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	cert, err := tls.X509KeyPair(testcert.LocalhostCert, testcert.LocalhostKey)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	if err != nil {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;httptest: NewTLSServer: %v&#34;, err))
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	existingConfig := s.TLS
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	if existingConfig != nil {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		s.TLS = existingConfig.Clone()
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	} else {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		s.TLS = new(tls.Config)
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	if s.TLS.NextProtos == nil {
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		nextProtos := []string{&#34;http/1.1&#34;}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		if s.EnableHTTP2 {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			nextProtos = []string{&#34;h2&#34;}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		s.TLS.NextProtos = nextProtos
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	if len(s.TLS.Certificates) == 0 {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		s.TLS.Certificates = []tls.Certificate{cert}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	s.certificate, err = x509.ParseCertificate(s.TLS.Certificates[0].Certificate[0])
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	if err != nil {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		panic(fmt.Sprintf(&#34;httptest: NewTLSServer: %v&#34;, err))
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	certpool := x509.NewCertPool()
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	certpool.AddCert(s.certificate)
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	s.client.Transport = &amp;http.Transport{
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		TLSClientConfig: &amp;tls.Config{
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			RootCAs: certpool,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		},
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		ForceAttemptHTTP2: s.EnableHTTP2,
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	s.Listener = tls.NewListener(s.Listener, s.TLS)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	s.URL = &#34;https://&#34; + s.Listener.Addr().String()
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	s.wrap()
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	s.goServe()
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span><span class="comment">// NewTLSServer starts and returns a new [Server] using TLS.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span><span class="comment">// The caller should call Close when finished, to shut it down.</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>func NewTLSServer(handler http.Handler) *Server {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	ts := NewUnstartedServer(handler)
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	ts.StartTLS()
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	return ts
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>type closeIdleTransport interface {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	CloseIdleConnections()
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">// Close shuts down the server and blocks until all outstanding</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">// requests on this server have completed.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>func (s *Server) Close() {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	s.mu.Lock()
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	if !s.closed {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		s.closed = true
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		s.Listener.Close()
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		s.Config.SetKeepAlivesEnabled(false)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		for c, st := range s.conns {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			<span class="comment">// Force-close any idle connections (those between</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>			<span class="comment">// requests) and new connections (those which connected</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			<span class="comment">// but never sent a request). StateNew connections are</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			<span class="comment">// super rare and have only been seen (in</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>			<span class="comment">// previously-flaky tests) in the case of</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			<span class="comment">// socket-late-binding races from the http Client</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			<span class="comment">// dialing this server and then getting an idle</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// connection before the dial completed. There is thus</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			<span class="comment">// a connected connection in StateNew with no</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			<span class="comment">// associated Request. We only close StateIdle and</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>			<span class="comment">// StateNew because they&#39;re not doing anything. It&#39;s</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			<span class="comment">// possible StateNew is about to do something in a few</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			<span class="comment">// milliseconds, but a previous CL to check again in a</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			<span class="comment">// few milliseconds wasn&#39;t liked (early versions of</span>
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			<span class="comment">// https://golang.org/cl/15151) so now we just</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			<span class="comment">// forcefully close StateNew. The docs for Server.Close say</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			<span class="comment">// we wait for &#34;outstanding requests&#34;, so we don&#39;t close things</span>
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			<span class="comment">// in StateActive.</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			if st == http.StateIdle || st == http.StateNew {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				s.closeConn(c)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		<span class="comment">// If this server doesn&#39;t shut down in 5 seconds, tell the user why.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		t := time.AfterFunc(5*time.Second, s.logCloseHangDebugInfo)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		defer t.Stop()
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	s.mu.Unlock()
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Not part of httptest.Server&#39;s correctness, but assume most</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	<span class="comment">// users of httptest.Server will be using the standard</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	<span class="comment">// transport, so help them out and close any idle connections for them.</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if t, ok := http.DefaultTransport.(closeIdleTransport); ok {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		t.CloseIdleConnections()
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	<span class="comment">// Also close the client idle connections.</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if s.client != nil {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if t, ok := s.client.Transport.(closeIdleTransport); ok {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			t.CloseIdleConnections()
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		}
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	s.wg.Wait()
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>func (s *Server) logCloseHangDebugInfo() {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	s.mu.Lock()
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	defer s.mu.Unlock()
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	var buf strings.Builder
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	buf.WriteString(&#34;httptest.Server blocked in Close after 5 seconds, waiting for connections:\n&#34;)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	for c, st := range s.conns {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		fmt.Fprintf(&amp;buf, &#34;  %T %p %v in state %v\n&#34;, c, c, c.RemoteAddr(), st)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	log.Print(buf.String())
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>}
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span><span class="comment">// CloseClientConnections closes any open HTTP connections to the test Server.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>func (s *Server) CloseClientConnections() {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	s.mu.Lock()
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	nconn := len(s.conns)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	ch := make(chan struct{}, nconn)
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	for c := range s.conns {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		go s.closeConnChan(c, ch)
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	s.mu.Unlock()
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	<span class="comment">// Wait for outstanding closes to finish.</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	<span class="comment">// Out of paranoia for making a late change in Go 1.6, we</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	<span class="comment">// bound how long this can wait, since golang.org/issue/14291</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	<span class="comment">// isn&#39;t fully understood yet. At least this should only be used</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	<span class="comment">// in tests.</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	timer := time.NewTimer(5 * time.Second)
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>	defer timer.Stop()
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	for i := 0; i &lt; nconn; i++ {
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		select {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		case &lt;-ch:
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		case &lt;-timer.C:
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			<span class="comment">// Too slow. Give up.</span>
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>			return
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>	}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">// Certificate returns the certificate used by the server, or nil if</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">// the server doesn&#39;t use TLS.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>func (s *Server) Certificate() *x509.Certificate {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	return s.certificate
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">// Client returns an HTTP client configured for making requests to the server.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">// It is configured to trust the server&#39;s TLS test certificate and will</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">// close its idle connections on [Server.Close].</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>func (s *Server) Client() *http.Client {
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	return s.client
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>func (s *Server) goServe() {
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	s.wg.Add(1)
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>	go func() {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		defer s.wg.Done()
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		s.Config.Serve(s.Listener)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	}()
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span><span class="comment">// wrap installs the connection state-tracking hook to know which</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span><span class="comment">// connections are idle.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>func (s *Server) wrap() {
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	oldHook := s.Config.ConnState
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	s.Config.ConnState = func(c net.Conn, cs http.ConnState) {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		s.mu.Lock()
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		defer s.mu.Unlock()
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		switch cs {
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		case http.StateNew:
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			if _, exists := s.conns[c]; exists {
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>				panic(&#34;invalid state transition&#34;)
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>			}
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>			if s.conns == nil {
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>				s.conns = make(map[net.Conn]http.ConnState)
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>			}
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			<span class="comment">// Add c to the set of tracked conns and increment it to the</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			<span class="comment">// waitgroup.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>			s.wg.Add(1)
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>			s.conns[c] = cs
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>			if s.closed {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>				<span class="comment">// Probably just a socket-late-binding dial from</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>				<span class="comment">// the default transport that lost the race (and</span>
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>				<span class="comment">// thus this connection is now idle and will</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				<span class="comment">// never be used).</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>				s.closeConn(c)
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>			}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		case http.StateActive:
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>			if oldState, ok := s.conns[c]; ok {
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>				if oldState != http.StateNew &amp;&amp; oldState != http.StateIdle {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>					panic(&#34;invalid state transition&#34;)
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>				}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>				s.conns[c] = cs
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>			}
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		case http.StateIdle:
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>			if oldState, ok := s.conns[c]; ok {
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>				if oldState != http.StateActive {
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>					panic(&#34;invalid state transition&#34;)
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>				}
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>				s.conns[c] = cs
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>			}
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			if s.closed {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>				s.closeConn(c)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>			}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		case http.StateHijacked, http.StateClosed:
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>			<span class="comment">// Remove c from the set of tracked conns and decrement it from the</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>			<span class="comment">// waitgroup, unless it was previously removed.</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>			if _, ok := s.conns[c]; ok {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>				delete(s.conns, c)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>				<span class="comment">// Keep Close from returning until the user&#39;s ConnState hook</span>
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>				<span class="comment">// (if any) finishes.</span>
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>				defer s.wg.Done()
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		}
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		if oldHook != nil {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>			oldHook(c, cs)
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		}
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span><span class="comment">// closeConn closes c.</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span><span class="comment">// s.mu must be held.</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>func (s *Server) closeConn(c net.Conn) { s.closeConnChan(c, nil) }
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>
<span id="L378" class="ln">   378&nbsp;&nbsp;</span><span class="comment">// closeConnChan is like closeConn, but takes an optional channel to receive a value</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span><span class="comment">// when the goroutine closing c is done.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>func (s *Server) closeConnChan(c net.Conn, done chan&lt;- struct{}) {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	c.Close()
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	if done != nil {
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		done &lt;- struct{}{}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>
</pre><p><a href="server.go?m=text">View as plain text</a></p>

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
