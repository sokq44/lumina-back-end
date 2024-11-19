<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/net/nss.go - Go Documentation Server</title>

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
<a href="nss.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/net">net</a>/<span class="text-muted">nss.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/net">net</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2015 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package net
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;internal/bytealg&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>)
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>const (
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	nssConfigPath = &#34;/etc/nsswitch.conf&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>)
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>var nssConfig nsswitchConfig
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>type nsswitchConfig struct {
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	initOnce sync.Once <span class="comment">// guards init of nsswitchConfig</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	<span class="comment">// ch is used as a semaphore that only allows one lookup at a</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// time to recheck nsswitch.conf</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	ch          chan struct{} <span class="comment">// guards lastChecked and modTime</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	lastChecked time.Time     <span class="comment">// last time nsswitch.conf was checked</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	mu      sync.Mutex <span class="comment">// protects nssConf</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	nssConf *nssConf
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>func getSystemNSS() *nssConf {
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	nssConfig.tryUpdate()
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	nssConfig.mu.Lock()
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	conf := nssConfig.nssConf
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	nssConfig.mu.Unlock()
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	return conf
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// init initializes conf and is only called via conf.initOnce.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func (conf *nsswitchConfig) init() {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	conf.nssConf = parseNSSConfFile(&#34;/etc/nsswitch.conf&#34;)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	conf.lastChecked = time.Now()
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	conf.ch = make(chan struct{}, 1)
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// tryUpdate tries to update conf.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>func (conf *nsswitchConfig) tryUpdate() {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	conf.initOnce.Do(conf.init)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	<span class="comment">// Ensure only one update at a time checks nsswitch.conf</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	if !conf.tryAcquireSema() {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		return
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	defer conf.releaseSema()
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	now := time.Now()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if conf.lastChecked.After(now.Add(-5 * time.Second)) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		return
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	conf.lastChecked = now
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	var mtime time.Time
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	if fi, err := os.Stat(nssConfigPath); err == nil {
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		mtime = fi.ModTime()
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if mtime.Equal(conf.nssConf.mtime) {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		return
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	nssConf := parseNSSConfFile(nssConfigPath)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	conf.mu.Lock()
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	conf.nssConf = nssConf
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	conf.mu.Unlock()
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>func (conf *nsswitchConfig) acquireSema() {
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	conf.ch &lt;- struct{}{}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>func (conf *nsswitchConfig) tryAcquireSema() bool {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	select {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	case conf.ch &lt;- struct{}{}:
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		return true
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	default:
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return false
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func (conf *nsswitchConfig) releaseSema() {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	&lt;-conf.ch
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>}
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span><span class="comment">// nssConf represents the state of the machine&#39;s /etc/nsswitch.conf file.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>type nssConf struct {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	mtime   time.Time              <span class="comment">// time of nsswitch.conf modification</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	err     error                  <span class="comment">// any error encountered opening or parsing the file</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	sources map[string][]nssSource <span class="comment">// keyed by database (e.g. &#34;hosts&#34;)</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>type nssSource struct {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	source   string <span class="comment">// e.g. &#34;compat&#34;, &#34;files&#34;, &#34;mdns4_minimal&#34;</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	criteria []nssCriterion
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span><span class="comment">// standardCriteria reports all specified criteria have the default</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span><span class="comment">// status actions.</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>func (s nssSource) standardCriteria() bool {
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	for i, crit := range s.criteria {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		if !crit.standardStatusAction(i == len(s.criteria)-1) {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			return false
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	return true
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// nssCriterion is the parsed structure of one of the criteria in brackets</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// after an NSS source name.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>type nssCriterion struct {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	negate bool   <span class="comment">// if &#34;!&#34; was present</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	status string <span class="comment">// e.g. &#34;success&#34;, &#34;unavail&#34; (lowercase)</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	action string <span class="comment">// e.g. &#34;return&#34;, &#34;continue&#34; (lowercase)</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span><span class="comment">// standardStatusAction reports whether c is equivalent to not</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span><span class="comment">// specifying the criterion at all. last is whether this criteria is the</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// last in the list.</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>func (c nssCriterion) standardStatusAction(last bool) bool {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	if c.negate {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		return false
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	var def string
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	switch c.status {
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	case &#34;success&#34;:
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		def = &#34;return&#34;
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	case &#34;notfound&#34;, &#34;unavail&#34;, &#34;tryagain&#34;:
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		def = &#34;continue&#34;
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	default:
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		<span class="comment">// Unknown status</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return false
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	if last &amp;&amp; c.action == &#34;return&#34; {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		return true
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	return c.action == def
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>func parseNSSConfFile(file string) *nssConf {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	f, err := open(file)
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	if err != nil {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		return &amp;nssConf{err: err}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	defer f.close()
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	mtime, _, err := f.stat()
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if err != nil {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return &amp;nssConf{err: err}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	conf := parseNSSConf(f)
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	conf.mtime = mtime
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	return conf
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>func parseNSSConf(f *file) *nssConf {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	conf := new(nssConf)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	for line, ok := f.readLine(); ok; line, ok = f.readLine() {
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		line = trimSpace(removeComment(line))
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if len(line) == 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			continue
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		colon := bytealg.IndexByteString(line, &#39;:&#39;)
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		if colon == -1 {
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			conf.err = errors.New(&#34;no colon on line&#34;)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>			return conf
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		}
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		db := trimSpace(line[:colon])
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		srcs := line[colon+1:]
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		for {
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			srcs = trimSpace(srcs)
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			if len(srcs) == 0 {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				break
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			sp := bytealg.IndexByteString(srcs, &#39; &#39;)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			var src string
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			if sp == -1 {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>				src = srcs
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				srcs = &#34;&#34; <span class="comment">// done</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>			} else {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>				src = srcs[:sp]
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				srcs = trimSpace(srcs[sp+1:])
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			var criteria []nssCriterion
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			<span class="comment">// See if there&#39;s a criteria block in brackets.</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			if len(srcs) &gt; 0 &amp;&amp; srcs[0] == &#39;[&#39; {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>				bclose := bytealg.IndexByteString(srcs, &#39;]&#39;)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>				if bclose == -1 {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>					conf.err = errors.New(&#34;unclosed criterion bracket&#34;)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>					return conf
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>				}
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				var err error
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				criteria, err = parseCriteria(srcs[1:bclose])
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>				if err != nil {
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>					conf.err = errors.New(&#34;invalid criteria: &#34; + srcs[1:bclose])
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>					return conf
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				}
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				srcs = srcs[bclose+1:]
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>			}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			if conf.sources == nil {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				conf.sources = make(map[string][]nssSource)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			conf.sources[db] = append(conf.sources[db], nssSource{
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				source:   src,
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>				criteria: criteria,
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			})
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		}
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	return conf
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span><span class="comment">// parses &#34;foo=bar !foo=bar&#34;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>func parseCriteria(x string) (c []nssCriterion, err error) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	err = foreachField(x, func(f string) error {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		not := false
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		if len(f) &gt; 0 &amp;&amp; f[0] == &#39;!&#39; {
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			not = true
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			f = f[1:]
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		if len(f) &lt; 3 {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			return errors.New(&#34;criterion too short&#34;)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		eq := bytealg.IndexByteString(f, &#39;=&#39;)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		if eq == -1 {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			return errors.New(&#34;criterion lacks equal sign&#34;)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		if hasUpperCase(f) {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			lower := []byte(f)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			lowerASCIIBytes(lower)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			f = string(lower)
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		c = append(c, nssCriterion{
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			negate: not,
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>			status: f[:eq],
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			action: f[eq+1:],
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		})
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		return nil
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	})
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	return
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
</pre><p><a href="nss.go?m=text">View as plain text</a></p>

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
