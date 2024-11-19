<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/syscall/env_unix.go - Go Documentation Server</title>

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
<a href="env_unix.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/syscall">syscall</a>/<span class="text-muted">env_unix.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/syscall">syscall</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2010 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//go:build unix || (js &amp;&amp; wasm) || plan9 || wasip1</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// Unix environment variables.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package syscall
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>import (
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;sync&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>var (
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	<span class="comment">// envOnce guards initialization by copyenv, which populates env.</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	envOnce sync.Once
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// envLock guards env and envs.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	envLock sync.RWMutex
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// env maps from an environment variable to its first occurrence in envs.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	env map[string]int
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// envs is provided by the runtime. elements are expected to</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// be of the form &#34;key=value&#34;. An empty string means deleted</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	<span class="comment">// (or a duplicate to be ignored).</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	envs []string = runtime_envs()
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>func runtime_envs() []string <span class="comment">// in package runtime</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>func copyenv() {
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	env = make(map[string]int)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	for i, s := range envs {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		for j := 0; j &lt; len(s); j++ {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>			if s[j] == &#39;=&#39; {
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>				key := s[:j]
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>				if _, ok := env[key]; !ok {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>					env[key] = i <span class="comment">// first mention of key</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>				} else {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>					<span class="comment">// Clear duplicate keys. This permits Unsetenv to</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>					<span class="comment">// safely delete only the first item without</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>					<span class="comment">// worrying about unshadowing a later one,</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>					<span class="comment">// which might be a security problem.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>					envs[i] = &#34;&#34;
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>				}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>				break
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		}
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>func Unsetenv(key string) error {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	envOnce.Do(copyenv)
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	envLock.Lock()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	defer envLock.Unlock()
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if i, ok := env[key]; ok {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		envs[i] = &#34;&#34;
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		delete(env, key)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	runtimeUnsetenv(key)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return nil
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func Getenv(key string) (value string, found bool) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	envOnce.Do(copyenv)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	if len(key) == 0 {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	envLock.RLock()
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	defer envLock.RUnlock()
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	i, ok := env[key]
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	if !ok {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		return &#34;&#34;, false
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	}
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	s := envs[i]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	for i := 0; i &lt; len(s); i++ {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if s[i] == &#39;=&#39; {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			return s[i+1:], true
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	return &#34;&#34;, false
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>}
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>func Setenv(key, value string) error {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	envOnce.Do(copyenv)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	if len(key) == 0 {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		return EINVAL
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	for i := 0; i &lt; len(key); i++ {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		if key[i] == &#39;=&#39; || key[i] == 0 {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			return EINVAL
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// On Plan 9, null is used as a separator, eg in $path.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if runtime.GOOS != &#34;plan9&#34; {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		for i := 0; i &lt; len(value); i++ {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>			if value[i] == 0 {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>				return EINVAL
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	envLock.Lock()
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	defer envLock.Unlock()
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	i, ok := env[key]
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	kv := key + &#34;=&#34; + value
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	if ok {
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		envs[i] = kv
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	} else {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		i = len(envs)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		envs = append(envs, kv)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	env[key] = i
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	runtimeSetenv(key, value)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	return nil
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>func Clearenv() {
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	envOnce.Do(copyenv) <span class="comment">// prevent copyenv in Getenv/Setenv</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	envLock.Lock()
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	defer envLock.Unlock()
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	for k := range env {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		runtimeUnsetenv(k)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	env = make(map[string]int)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	envs = []string{}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>func Environ() []string {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	envOnce.Do(copyenv)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	envLock.RLock()
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	defer envLock.RUnlock()
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	a := make([]string, 0, len(envs))
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	for _, env := range envs {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		if env != &#34;&#34; {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>			a = append(a, env)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	return a
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
</pre><p><a href="env_unix.go?m=text">View as plain text</a></p>

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
