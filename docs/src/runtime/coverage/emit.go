<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/coverage/emit.go - Go Documentation Server</title>

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
<a href="emit.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/coverage">coverage</a>/<span class="text-muted">emit.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/runtime/coverage">runtime/coverage</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2022 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package coverage
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;crypto/md5&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/coverage&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/coverage/encodecounter&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/coverage/encodemeta&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/coverage/rtcov&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;path/filepath&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;runtime&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;sync/atomic&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// This file contains functions that support the writing of data files</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// emitted at the end of code coverage testing runs, from instrumented</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// executables.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// getCovMetaList returns a list of meta-data blobs registered</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// for the currently executing instrumented program. It is defined in the</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">// runtime.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>func getCovMetaList() []rtcov.CovMetaBlob
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// getCovCounterList returns a list of counter-data blobs registered</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// for the currently executing instrumented program. It is defined in the</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// runtime.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func getCovCounterList() []rtcov.CovCounterBlob
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// getCovPkgMap returns a map storing the remapped package IDs for</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">// hard-coded runtime packages (see internal/coverage/pkgid.go for</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">// more on why hard-coded package IDs are needed). This function</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">// is defined in the runtime.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func getCovPkgMap() map[int]int
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// emitState holds useful state information during the emit process.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// When an instrumented program finishes execution and starts the</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">// process of writing out coverage data, it&#39;s possible that an</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// existing meta-data file already exists in the output directory. In</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// this case openOutputFiles() below will leave the &#39;mf&#39; field below</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// as nil. If a new meta-data file is needed, field &#39;mfname&#39; will be</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// the final desired path of the meta file, &#39;mftmp&#39; will be a</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// temporary file, and &#39;mf&#39; will be an open os.File pointer for</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// &#39;mftmp&#39;. The meta-data file payload will be written to &#39;mf&#39;, the</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">// temp file will be then closed and renamed (from &#39;mftmp&#39; to</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// &#39;mfname&#39;), so as to insure that the meta-data file is created</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// atomically; we want this so that things work smoothly in cases</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// where there are several instances of a given instrumented program</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// all terminating at the same time and trying to create meta-data</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span><span class="comment">// files simultaneously.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span><span class="comment">// For counter data files there is less chance of a collision, hence</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span><span class="comment">// the openOutputFiles() stores the counter data file in &#39;cfname&#39; and</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span><span class="comment">// then places the *io.File into &#39;cf&#39;.</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>type emitState struct {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	mfname string   <span class="comment">// path of final meta-data output file</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	mftmp  string   <span class="comment">// path to meta-data temp file (if needed)</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	mf     *os.File <span class="comment">// open os.File for meta-data temp file</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	cfname string   <span class="comment">// path of final counter data file</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	cftmp  string   <span class="comment">// path to counter data temp file</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	cf     *os.File <span class="comment">// open os.File for counter data file</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	outdir string   <span class="comment">// output directory</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// List of meta-data symbols obtained from the runtime</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	metalist []rtcov.CovMetaBlob
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// List of counter-data symbols obtained from the runtime</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	counterlist []rtcov.CovCounterBlob
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">// Table to use for remapping hard-coded pkg ids.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	pkgmap map[int]int
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// emit debug trace output</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	debug bool
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>var (
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// finalHash is computed at init time from the list of meta-data</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// symbols registered during init. It is used both for writing the</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// meta-data file and counter-data files.</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	finalHash [16]byte
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// Set to true when we&#39;ve computed finalHash + finalMetaLen.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	finalHashComputed bool
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// Total meta-data length.</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	finalMetaLen uint64
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// Records whether we&#39;ve already attempted to write meta-data.</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	metaDataEmitAttempted bool
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	<span class="comment">// Counter mode for this instrumented program run.</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	cmode coverage.CounterMode
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	<span class="comment">// Counter granularity for this instrumented program run.</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	cgran coverage.CounterGranularity
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	<span class="comment">// Cached value of GOCOVERDIR environment variable.</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	goCoverDir string
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Copy of os.Args made at init time, converted into map format.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	capturedOsArgs map[string]string
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// Flag used in tests to signal that coverage data already written.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	covProfileAlreadyEmitted bool
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span><span class="comment">// fileType is used to select between counter-data files and</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span><span class="comment">// meta-data files.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>type fileType int
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>const (
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	noFile = 1 &lt;&lt; iota
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	metaDataFile
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	counterDataFile
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">// emitMetaData emits the meta-data output file for this coverage run.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// This entry point is intended to be invoked by the compiler from</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// an instrumented program&#39;s main package init func.</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>func emitMetaData() {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	if covProfileAlreadyEmitted {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	ml, err := prepareForMetaEmit()
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	if err != nil {
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;error: coverage meta-data prep failed: %v\n&#34;, err)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		if os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34; {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			panic(&#34;meta-data write failure&#34;)
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	}
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	if len(ml) == 0 {
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;program not built with -cover\n&#34;)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		return
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	goCoverDir = os.Getenv(&#34;GOCOVERDIR&#34;)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if goCoverDir == &#34;&#34; {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;warning: GOCOVERDIR not set, no coverage data emitted\n&#34;)
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	if err := emitMetaDataToDirectory(goCoverDir, ml); err != nil {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;error: coverage meta-data emit failed: %v\n&#34;, err)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34; {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			panic(&#34;meta-data write failure&#34;)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		}
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>func modeClash(m coverage.CounterMode) bool {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if m == coverage.CtrModeRegOnly || m == coverage.CtrModeTestMain {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return false
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if cmode == coverage.CtrModeInvalid {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		cmode = m
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return false
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	return cmode != m
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>}
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>func granClash(g coverage.CounterGranularity) bool {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	if cgran == coverage.CtrGranularityInvalid {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		cgran = g
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		return false
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	return cgran != g
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span><span class="comment">// prepareForMetaEmit performs preparatory steps needed prior to</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span><span class="comment">// emitting a meta-data file, notably computing a final hash of</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span><span class="comment">// all meta-data blobs and capturing os args.</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>func prepareForMetaEmit() ([]rtcov.CovMetaBlob, error) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// Ask the runtime for the list of coverage meta-data symbols.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	ml := getCovMetaList()
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">// In the normal case (go build -o prog.exe ... ; ./prog.exe)</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">// len(ml) will always be non-zero, but we check here since at</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">// some point this function will be reachable via user-callable</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">// APIs (for example, to write out coverage data from a server</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">// program that doesn&#39;t ever call os.Exit).</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if len(ml) == 0 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return nil, nil
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	s := &amp;emitState{
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		metalist: ml,
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		debug:    os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34;,
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">// Capture os.Args() now so as to avoid issues if args</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">// are rewritten during program execution.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	capturedOsArgs = captureOsArgs()
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	if s.debug {
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;=+= GOCOVERDIR is %s\n&#34;, os.Getenv(&#34;GOCOVERDIR&#34;))
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;=+= contents of covmetalist:\n&#34;)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>		for k, b := range ml {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			fmt.Fprintf(os.Stderr, &#34;=+= slot: %d path: %s &#34;, k, b.PkgPath)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			if b.PkgID != -1 {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				fmt.Fprintf(os.Stderr, &#34; hcid: %d&#34;, b.PkgID)
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			fmt.Fprintf(os.Stderr, &#34;\n&#34;)
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>		}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		pm := getCovPkgMap()
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;=+= remap table:\n&#34;)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		for from, to := range pm {
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>			fmt.Fprintf(os.Stderr, &#34;=+= from %d to %d\n&#34;,
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				uint32(from), uint32(to))
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	}
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	h := md5.New()
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	tlen := uint64(unsafe.Sizeof(coverage.MetaFileHeader{}))
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	for _, entry := range ml {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		if _, err := h.Write(entry.Hash[:]); err != nil {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>			return nil, err
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		}
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		tlen += uint64(entry.Len)
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		ecm := coverage.CounterMode(entry.CounterMode)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		if modeClash(ecm) {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;coverage counter mode clash: package %s uses mode=%d, but package %s uses mode=%s\n&#34;, ml[0].PkgPath, cmode, entry.PkgPath, ecm)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		ecg := coverage.CounterGranularity(entry.CounterGranularity)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		if granClash(ecg) {
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			return nil, fmt.Errorf(&#34;coverage counter granularity clash: package %s uses gran=%d, but package %s uses gran=%s\n&#34;, ml[0].PkgPath, cgran, entry.PkgPath, ecg)
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	<span class="comment">// Hash mode and granularity as well.</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	h.Write([]byte(cmode.String()))
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	h.Write([]byte(cgran.String()))
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	<span class="comment">// Compute final digest.</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	fh := h.Sum(nil)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	copy(finalHash[:], fh)
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	finalHashComputed = true
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	finalMetaLen = tlen
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	return ml, nil
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span><span class="comment">// emitMetaDataToDirectory emits the meta-data output file to the specified</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span><span class="comment">// directory, returning an error if something went wrong.</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>func emitMetaDataToDirectory(outdir string, ml []rtcov.CovMetaBlob) error {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	ml, err := prepareForMetaEmit()
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	if err != nil {
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		return err
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	}
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	if len(ml) == 0 {
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		return nil
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	}
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	metaDataEmitAttempted = true
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	s := &amp;emitState{
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		metalist: ml,
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		debug:    os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34;,
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		outdir:   outdir,
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// Open output files.</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if err := s.openOutputFiles(finalHash, finalMetaLen, metaDataFile); err != nil {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		return err
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	<span class="comment">// Emit meta-data file only if needed (may already be present).</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	if s.needMetaDataFile() {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		if err := s.emitMetaDataFile(finalHash, finalMetaLen); err != nil {
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			return err
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		}
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	}
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	return nil
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span><span class="comment">// emitCounterData emits the counter data output file for this coverage run.</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span><span class="comment">// This entry point is intended to be invoked by the runtime when an</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span><span class="comment">// instrumented program is terminating or calling os.Exit().</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>func emitCounterData() {
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	if goCoverDir == &#34;&#34; || !finalHashComputed || covProfileAlreadyEmitted {
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		return
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>	if err := emitCounterDataToDirectory(goCoverDir); err != nil {
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		fmt.Fprintf(os.Stderr, &#34;error: coverage counter data emit failed: %v\n&#34;, err)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		if os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34; {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			panic(&#34;counter-data write failure&#34;)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		}
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// emitCounterDataToDirectory emits the counter-data output file for this coverage run.</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>func emitCounterDataToDirectory(outdir string) error {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	<span class="comment">// Ask the runtime for the list of coverage counter symbols.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	cl := getCovCounterList()
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	if len(cl) == 0 {
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		<span class="comment">// no work to do here.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		return nil
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	}
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	if !finalHashComputed {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;error: meta-data not available (binary not built with -cover?)&#34;)
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	}
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>	<span class="comment">// Ask the runtime for the list of coverage counter symbols.</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	pm := getCovPkgMap()
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	s := &amp;emitState{
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		counterlist: cl,
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>		pkgmap:      pm,
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>		outdir:      outdir,
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		debug:       os.Getenv(&#34;GOCOVERDEBUG&#34;) != &#34;&#34;,
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>	<span class="comment">// Open output file.</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	if err := s.openOutputFiles(finalHash, finalMetaLen, counterDataFile); err != nil {
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		return err
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	if s.cf == nil {
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;counter data output file open failed (no additional info&#34;)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	<span class="comment">// Emit counter data file.</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>	if err := s.emitCounterDataFile(finalHash, s.cf); err != nil {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		return err
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>	}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	if err := s.cf.Close(); err != nil {
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;closing counter data file: %v&#34;, err)
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	}
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	<span class="comment">// Counter file has now been closed. Rename the temp to the</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>	<span class="comment">// final desired path.</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	if err := os.Rename(s.cftmp, s.cfname); err != nil {
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;writing %s: rename from %s failed: %v\n&#34;, s.cfname, s.cftmp, err)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	}
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	return nil
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>}
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>
<span id="L339" class="ln">   339&nbsp;&nbsp;</span><span class="comment">// emitCounterDataToWriter emits counter data for this coverage run to an io.Writer.</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>func (s *emitState) emitCounterDataToWriter(w io.Writer) error {
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if err := s.emitCounterDataFile(finalHash, w); err != nil {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		return err
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	return nil
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// openMetaFile determines whether we need to emit a meta-data output</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// file, or whether we can reuse the existing file in the coverage out</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// dir. It updates mfname/mftmp/mf fields in &#39;s&#39;, returning an error</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">// if something went wrong. See the comment on the emitState type</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// definition above for more on how file opening is managed.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>func (s *emitState) openMetaFile(metaHash [16]byte, metaLen uint64) error {
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>	<span class="comment">// Open meta-outfile for reading to see if it exists.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>	fn := fmt.Sprintf(&#34;%s.%x&#34;, coverage.MetaFilePref, metaHash)
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	s.mfname = filepath.Join(s.outdir, fn)
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	fi, err := os.Stat(s.mfname)
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	if err != nil || fi.Size() != int64(metaLen) {
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		<span class="comment">// We need a new meta-file.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>		tname := &#34;tmp.&#34; + fn + strconv.FormatInt(time.Now().UnixNano(), 10)
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		s.mftmp = filepath.Join(s.outdir, tname)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		s.mf, err = os.Create(s.mftmp)
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		if err != nil {
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;creating meta-data file %s: %v&#34;, s.mftmp, err)
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		}
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	return nil
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span><span class="comment">// openCounterFile opens an output file for the counter data portion</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span><span class="comment">// of a test coverage run. If updates the &#39;cfname&#39; and &#39;cf&#39; fields in</span>
<span id="L372" class="ln">   372&nbsp;&nbsp;</span><span class="comment">// &#39;s&#39;, returning an error if something went wrong.</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>func (s *emitState) openCounterFile(metaHash [16]byte) error {
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	processID := os.Getpid()
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	fn := fmt.Sprintf(coverage.CounterFileTempl, coverage.CounterFilePref, metaHash, processID, time.Now().UnixNano())
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	s.cfname = filepath.Join(s.outdir, fn)
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	s.cftmp = filepath.Join(s.outdir, &#34;tmp.&#34;+fn)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	var err error
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	s.cf, err = os.Create(s.cftmp)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	if err != nil {
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;creating counter data file %s: %v&#34;, s.cftmp, err)
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	}
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	return nil
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
<span id="L386" class="ln">   386&nbsp;&nbsp;</span><span class="comment">// openOutputFiles opens output files in preparation for emitting</span>
<span id="L387" class="ln">   387&nbsp;&nbsp;</span><span class="comment">// coverage data. In the case of the meta-data file, openOutputFiles</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span><span class="comment">// may determine that we can reuse an existing meta-data file in the</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">// outdir, in which case it will leave the &#39;mf&#39; field in the state</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// struct as nil. If a new meta-file is needed, the field &#39;mfname&#39;</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span><span class="comment">// will be the final desired path of the meta file, &#39;mftmp&#39; will be a</span>
<span id="L392" class="ln">   392&nbsp;&nbsp;</span><span class="comment">// temporary file, and &#39;mf&#39; will be an open os.File pointer for</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span><span class="comment">// &#39;mftmp&#39;. The idea is that the client/caller will write content into</span>
<span id="L394" class="ln">   394&nbsp;&nbsp;</span><span class="comment">// &#39;mf&#39;, close it, and then rename &#39;mftmp&#39; to &#39;mfname&#39;. This function</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span><span class="comment">// also opens the counter data output file, setting &#39;cf&#39; and &#39;cfname&#39;</span>
<span id="L396" class="ln">   396&nbsp;&nbsp;</span><span class="comment">// in the state struct.</span>
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>func (s *emitState) openOutputFiles(metaHash [16]byte, metaLen uint64, which fileType) error {
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	fi, err := os.Stat(s.outdir)
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	if err != nil {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;output directory %q inaccessible (err: %v); no coverage data written&#34;, s.outdir, err)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>	}
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	if !fi.IsDir() {
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;output directory %q not a directory; no coverage data written&#34;, s.outdir)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	}
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if (which &amp; metaDataFile) != 0 {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		if err := s.openMetaFile(metaHash, metaLen); err != nil {
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			return err
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>		}
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if (which &amp; counterDataFile) != 0 {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		if err := s.openCounterFile(metaHash); err != nil {
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>			return err
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>		}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	return nil
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>}
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">// emitMetaDataFile emits coverage meta-data to a previously opened</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">// temporary file (s.mftmp), then renames the generated file to the</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">// final path (s.mfname).</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>func (s *emitState) emitMetaDataFile(finalHash [16]byte, tlen uint64) error {
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	if err := writeMetaData(s.mf, s.metalist, cmode, cgran, finalHash); err != nil {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;writing %s: %v\n&#34;, s.mftmp, err)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	}
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>	if err := s.mf.Close(); err != nil {
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;closing meta data temp file: %v&#34;, err)
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// Temp file has now been flushed and closed. Rename the temp to the</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	<span class="comment">// final desired path.</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	if err := os.Rename(s.mftmp, s.mfname); err != nil {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;writing %s: rename from %s failed: %v\n&#34;, s.mfname, s.mftmp, err)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>	}
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	return nil
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span><span class="comment">// needMetaDataFile returns TRUE if we need to emit a meta-data file</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span><span class="comment">// for this program run. It should be used only after</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span><span class="comment">// openOutputFiles() has been invoked.</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>func (s *emitState) needMetaDataFile() bool {
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	return s.mf != nil
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>}
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>func writeMetaData(w io.Writer, metalist []rtcov.CovMetaBlob, cmode coverage.CounterMode, gran coverage.CounterGranularity, finalHash [16]byte) error {
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>	mfw := encodemeta.NewCoverageMetaFileWriter(&#34;&lt;io.Writer&gt;&#34;, w)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>	var blobs [][]byte
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>	for _, e := range metalist {
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		sd := unsafe.Slice(e.P, int(e.Len))
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		blobs = append(blobs, sd)
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>	}
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	return mfw.Write(finalHash, blobs, cmode, gran)
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>}
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>func (s *emitState) VisitFuncs(f encodecounter.CounterVisitorFn) error {
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	var tcounters []uint32
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	rdCounters := func(actrs []atomic.Uint32, ctrs []uint32) []uint32 {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		ctrs = ctrs[:0]
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		for i := range actrs {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			ctrs = append(ctrs, actrs[i].Load())
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		return ctrs
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	dpkg := uint32(0)
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>	for _, c := range s.counterlist {
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		sd := unsafe.Slice((*atomic.Uint32)(unsafe.Pointer(c.Counters)), int(c.Len))
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>		for i := 0; i &lt; len(sd); i++ {
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>			<span class="comment">// Skip ahead until the next non-zero value.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>			sdi := sd[i].Load()
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>			if sdi == 0 {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>				continue
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>			}
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			<span class="comment">// We found a function that was executed.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			nCtrs := sd[i+coverage.NumCtrsOffset].Load()
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>			pkgId := sd[i+coverage.PkgIdOffset].Load()
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>			funcId := sd[i+coverage.FuncIdOffset].Load()
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			cst := i + coverage.FirstCtrOffset
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>			counters := sd[cst : cst+int(nCtrs)]
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>			<span class="comment">// Check to make sure that we have at least one live</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>			<span class="comment">// counter. See the implementation note in ClearCoverageCounters</span>
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>			<span class="comment">// for a description of why this is needed.</span>
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>			isLive := false
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>			for i := 0; i &lt; len(counters); i++ {
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>				if counters[i].Load() != 0 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>					isLive = true
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>					break
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>				}
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			}
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			if !isLive {
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>				<span class="comment">// Skip this function.</span>
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				i += coverage.FirstCtrOffset + int(nCtrs) - 1
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>				continue
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			}
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			if s.debug {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>				if pkgId != dpkg {
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>					dpkg = pkgId
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>					fmt.Fprintf(os.Stderr, &#34;\n=+= %d: pk=%d visit live fcn&#34;,
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>						i, pkgId)
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>				}
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>				fmt.Fprintf(os.Stderr, &#34; {i=%d F%d NC%d}&#34;, i, funcId, nCtrs)
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>			<span class="comment">// Vet and/or fix up package ID. A package ID of zero</span>
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			<span class="comment">// indicates that there is some new package X that is a</span>
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			<span class="comment">// runtime dependency, and this package has code that</span>
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			<span class="comment">// executes before its corresponding init package runs.</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>			<span class="comment">// This is a fatal error that we should only see during</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>			<span class="comment">// Go development (e.g. tip).</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>			ipk := int32(pkgId)
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>			if ipk == 0 {
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>				fmt.Fprintf(os.Stderr, &#34;\n&#34;)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>				reportErrorInHardcodedList(int32(i), ipk, funcId, nCtrs)
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>			} else if ipk &lt; 0 {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>				if newId, ok := s.pkgmap[int(ipk)]; ok {
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>					pkgId = uint32(newId)
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>				} else {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>					fmt.Fprintf(os.Stderr, &#34;\n&#34;)
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>					reportErrorInHardcodedList(int32(i), ipk, funcId, nCtrs)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>				}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>			} else {
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>				<span class="comment">// The package ID value stored in the counter array</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>				<span class="comment">// has 1 added to it (so as to preclude the</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>				<span class="comment">// possibility of a zero value ; see</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>				<span class="comment">// runtime.addCovMeta), so subtract off 1 here to form</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>				<span class="comment">// the real package ID.</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>				pkgId--
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			}
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>			tcounters = rdCounters(counters, tcounters)
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>			if err := f(pkgId, funcId, tcounters); err != nil {
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>				return err
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>			}
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>			<span class="comment">// Skip over this function.</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			i += coverage.FirstCtrOffset + int(nCtrs) - 1
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		if s.debug {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>			fmt.Fprintf(os.Stderr, &#34;\n&#34;)
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>		}
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	return nil
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">// captureOsArgs converts os.Args() into the format we use to store</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">// this info in the counter data file (counter data file &#34;args&#34;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span><span class="comment">// section is a generic key-value collection). See the &#39;args&#39; section</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span><span class="comment">// in internal/coverage/defs.go for more info. The args map</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span><span class="comment">// is also used to capture GOOS + GOARCH values as well.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>func captureOsArgs() map[string]string {
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	m := make(map[string]string)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	m[&#34;argc&#34;] = strconv.Itoa(len(os.Args))
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	for k, a := range os.Args {
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>		m[fmt.Sprintf(&#34;argv%d&#34;, k)] = a
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	}
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	m[&#34;GOOS&#34;] = runtime.GOOS
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	m[&#34;GOARCH&#34;] = runtime.GOARCH
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	return m
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span><span class="comment">// emitCounterDataFile emits the counter data portion of a</span>
<span id="L568" class="ln">   568&nbsp;&nbsp;</span><span class="comment">// coverage output file (to the file &#39;s.cf&#39;).</span>
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>func (s *emitState) emitCounterDataFile(finalHash [16]byte, w io.Writer) error {
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	cfw := encodecounter.NewCoverageDataWriter(w, coverage.CtrULeb128)
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	if err := cfw.Write(finalHash, capturedOsArgs, s); err != nil {
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>		return err
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	return nil
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span><span class="comment">// markProfileEmitted signals the runtime/coverage machinery that</span>
<span id="L578" class="ln">   578&nbsp;&nbsp;</span><span class="comment">// coverage data output files have already been written out, and there</span>
<span id="L579" class="ln">   579&nbsp;&nbsp;</span><span class="comment">// is no need to take any additional action at exit time. This</span>
<span id="L580" class="ln">   580&nbsp;&nbsp;</span><span class="comment">// function is called (via linknamed reference) from the</span>
<span id="L581" class="ln">   581&nbsp;&nbsp;</span><span class="comment">// coverage-related boilerplate code in _testmain.go emitted for go</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// unit tests.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>func markProfileEmitted(val bool) {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	covProfileAlreadyEmitted = val
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>func reportErrorInHardcodedList(slot, pkgID int32, fnID, nCtrs uint32) {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>	metaList := getCovMetaList()
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	pkgMap := getCovPkgMap()
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	println(&#34;internal error in coverage meta-data tracking:&#34;)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	println(&#34;encountered bad pkgID:&#34;, pkgID, &#34; at slot:&#34;, slot,
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		&#34; fnID:&#34;, fnID, &#34; numCtrs:&#34;, nCtrs)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	println(&#34;list of hard-coded runtime package IDs needs revising.&#34;)
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	println(&#34;[see the comment on the &#39;rtPkgs&#39; var in &#34;)
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	println(&#34; &lt;goroot&gt;/src/internal/coverage/pkid.go]&#34;)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	println(&#34;registered list:&#34;)
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	for k, b := range metaList {
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>		print(&#34;slot: &#34;, k, &#34; path=&#39;&#34;, b.PkgPath, &#34;&#39; &#34;)
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		if b.PkgID != -1 {
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>			print(&#34; hard-coded id: &#34;, b.PkgID)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>		}
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		println(&#34;&#34;)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	println(&#34;remap table:&#34;)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>	for from, to := range pkgMap {
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>		println(&#34;from &#34;, from, &#34; to &#34;, to)
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	}
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>
</pre><p><a href="emit.go?m=text">View as plain text</a></p>

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
