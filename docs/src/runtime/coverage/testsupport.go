<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/runtime/coverage/testsupport.go - Go Documentation Server</title>

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
<a href="testsupport.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/runtime">runtime</a>/<a href="http://localhost:8080/src/runtime/coverage">coverage</a>/<span class="text-muted">testsupport.go</span>
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
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;internal/coverage&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;internal/coverage/calloc&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;internal/coverage/cformat&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;internal/coverage/cmerge&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;internal/coverage/decodecounter&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;internal/coverage/decodemeta&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;internal/coverage/pods&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;os&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;path/filepath&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;runtime/internal/atomic&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;unsafe&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>)
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// processCoverTestDir is called (via a linknamed reference) from</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// testmain code when &#34;go test -cover&#34; is in effect. It is not</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">// intended to be used other than internally by the Go command&#39;s</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// generated code.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>func processCoverTestDir(dir string, cfile string, cm string, cpkg string) error {
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	return processCoverTestDirInternal(dir, cfile, cm, cpkg, os.Stdout)
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">// processCoverTestDirInternal is an io.Writer version of processCoverTestDir,</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// exposed for unit testing.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>func processCoverTestDirInternal(dir string, cfile string, cm string, cpkg string, w io.Writer) error {
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	cmode := coverage.ParseCounterMode(cm)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if cmode == coverage.CtrModeInvalid {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;invalid counter mode %q&#34;, cm)
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	<span class="comment">// Emit meta-data and counter data.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	ml := getCovMetaList()
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	if len(ml) == 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		<span class="comment">// This corresponds to the case where we have a package that</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		<span class="comment">// contains test code but no functions (which is fine). In this</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		<span class="comment">// case there is no need to emit anything.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	} else {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		if err := emitMetaDataToDirectory(dir, ml); err != nil {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			return err
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		if err := emitCounterDataToDirectory(dir); err != nil {
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			return err
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	}
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	<span class="comment">// Collect pods from test run. For the majority of cases we would</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	<span class="comment">// expect to see a single pod here, but allow for multiple pods in</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">// case the test harness is doing extra work to collect data files</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	<span class="comment">// from builds that it kicks off as part of the testing.</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	podlist, err := pods.CollectPods([]string{dir}, false)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	if err != nil {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;reading from %s: %v&#34;, dir, err)
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	<span class="comment">// Open text output file if appropriate.</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	var tf *os.File
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	var tfClosed bool
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	if cfile != &#34;&#34; {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		var err error
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		tf, err = os.Create(cfile)
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if err != nil {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;internal error: opening coverage data output file %q: %v&#34;, cfile, err)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		defer func() {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			if !tfClosed {
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>				tfClosed = true
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>				tf.Close()
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		}()
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// Read/process the pods.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	ts := &amp;tstate{
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		cm:    &amp;cmerge.Merger{},
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		cf:    cformat.NewFormatter(cmode),
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		cmode: cmode,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	<span class="comment">// Generate the expected hash string based on the final meta-data</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	<span class="comment">// hash for this test, then look only for pods that refer to that</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	<span class="comment">// hash (just in case there are multiple instrumented executables</span>
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	<span class="comment">// in play). See issue #57924 for more on this.</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	hashstring := fmt.Sprintf(&#34;%x&#34;, finalHash)
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	importpaths := make(map[string]struct{})
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	for _, p := range podlist {
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		if !strings.Contains(p.MetaFile, hashstring) {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			continue
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		}
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		if err := ts.processPod(p, importpaths); err != nil {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			return err
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	metafilespath := filepath.Join(dir, coverage.MetaFilesFileName)
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if _, err := os.Stat(metafilespath); err == nil {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if err := ts.readAuxMetaFiles(metafilespath, importpaths); err != nil {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			return err
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		}
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	<span class="comment">// Emit percent.</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	if err := ts.cf.EmitPercent(w, cpkg, true, true); err != nil {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		return err
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	}
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	<span class="comment">// Emit text output.</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if tf != nil {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		if err := ts.cf.EmitTextual(tf); err != nil {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			return err
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		tfClosed = true
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		if err := tf.Close(); err != nil {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;closing %s: %v&#34;, cfile, err)
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	}
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return nil
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>type tstate struct {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	calloc.BatchCounterAlloc
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	cm    *cmerge.Merger
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	cf    *cformat.Formatter
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	cmode coverage.CounterMode
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// processPod reads coverage counter data for a specific pod.</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>func (ts *tstate) processPod(p pods.Pod, importpaths map[string]struct{}) error {
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">// Open meta-data file</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	f, err := os.Open(p.MetaFile)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	if err != nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;unable to open meta-data file %s: %v&#34;, p.MetaFile, err)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	defer func() {
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		f.Close()
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	}()
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	var mfr *decodemeta.CoverageMetaFileReader
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	mfr, err = decodemeta.NewCoverageMetaFileReader(f, nil)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	if err != nil {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;error reading meta-data file %s: %v&#34;, p.MetaFile, err)
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	newmode := mfr.CounterMode()
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	if newmode != ts.cmode {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;internal error: counter mode clash: %q from test harness, %q from data file %s&#34;, ts.cmode.String(), newmode.String(), p.MetaFile)
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	newgran := mfr.CounterGranularity()
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	if err := ts.cm.SetModeAndGranularity(p.MetaFile, cmode, newgran); err != nil {
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		return err
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">// A map to store counter data, indexed by pkgid/fnid tuple.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	pmm := make(map[pkfunc][]uint32)
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// Helper to read a single counter data file.</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	readcdf := func(cdf string) error {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		cf, err := os.Open(cdf)
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		if err != nil {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;opening counter data file %s: %s&#34;, cdf, err)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		defer cf.Close()
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		var cdr *decodecounter.CounterDataReader
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		cdr, err = decodecounter.NewCounterDataReader(cdf, cf)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		if err != nil {
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;reading counter data file %s: %s&#34;, cdf, err)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		}
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		var data decodecounter.FuncPayload
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		for {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			ok, err := cdr.NextFunc(&amp;data)
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			if err != nil {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;reading counter data file %s: %v&#34;, cdf, err)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>			}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			if !ok {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				break
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			<span class="comment">// NB: sanity check on pkg and func IDs?</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			key := pkfunc{pk: data.PkgIdx, fcn: data.FuncIdx}
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			if prev, found := pmm[key]; found {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				<span class="comment">// Note: no overflow reporting here.</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				if err, _ := ts.cm.MergeCounters(data.Counters, prev); err != nil {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>					return fmt.Errorf(&#34;processing counter data file %s: %v&#34;, cdf, err)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			}
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			c := ts.AllocateCounters(len(data.Counters))
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>			copy(c, data.Counters)
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>			pmm[key] = c
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		return nil
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// Read counter data files.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	for _, cdf := range p.CounterDataFiles {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if err := readcdf(cdf); err != nil {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			return err
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	}
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	<span class="comment">// Visit meta-data file.</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	np := uint32(mfr.NumPackages())
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	payload := []byte{}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	for pkIdx := uint32(0); pkIdx &lt; np; pkIdx++ {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>		var pd *decodemeta.CoverageMetaDataDecoder
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>		pd, payload, err = mfr.GetPackageDecoder(pkIdx, payload)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		if err != nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>			return fmt.Errorf(&#34;reading pkg %d from meta-file %s: %s&#34;, pkIdx, p.MetaFile, err)
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>		ts.cf.SetPackage(pd.PackagePath())
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		importpaths[pd.PackagePath()] = struct{}{}
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		var fd coverage.FuncDesc
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		nf := pd.NumFuncs()
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		for fnIdx := uint32(0); fnIdx &lt; nf; fnIdx++ {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			if err := pd.ReadFunc(fnIdx, &amp;fd); err != nil {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;reading meta-data file %s: %v&#34;,
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>					p.MetaFile, err)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			}
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			key := pkfunc{pk: pkIdx, fcn: fnIdx}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			counters, haveCounters := pmm[key]
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			for i := 0; i &lt; len(fd.Units); i++ {
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>				u := fd.Units[i]
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>				<span class="comment">// Skip units with non-zero parent (no way to represent</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>				<span class="comment">// these in the existing format).</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>				if u.Parent != 0 {
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>					continue
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>				}
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				count := uint32(0)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>				if haveCounters {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>					count = counters[i]
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>				}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>				ts.cf.AddUnit(fd.Srcfile, fd.Funcname, fd.Lit, u, count)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>			}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	}
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	return nil
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>type pkfunc struct {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	pk, fcn uint32
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>func (ts *tstate) readAuxMetaFiles(metafiles string, importpaths map[string]struct{}) error {
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	<span class="comment">// Unmarshall the information on available aux metafiles into</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	<span class="comment">// a MetaFileCollection struct.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	var mfc coverage.MetaFileCollection
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	data, err := os.ReadFile(metafiles)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if err != nil {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;error reading auxmetafiles file %q: %v&#34;, metafiles, err)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	if err := json.Unmarshal(data, &amp;mfc); err != nil {
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		return fmt.Errorf(&#34;error reading auxmetafiles file %q: %v&#34;, metafiles, err)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	}
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	<span class="comment">// Walk through each available aux meta-file. If we&#39;ve already</span>
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	<span class="comment">// seen the package path in question during the walk of the</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	<span class="comment">// &#34;regular&#34; meta-data file, then we can skip the package,</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	<span class="comment">// otherwise construct a dummy pod with the single meta-data file</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	<span class="comment">// (no counters) and invoke processPod on it.</span>
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	for i := range mfc.ImportPaths {
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		p := mfc.ImportPaths[i]
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		if _, ok := importpaths[p]; ok {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			continue
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>		var pod pods.Pod
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		pod.MetaFile = mfc.MetaFileFragments[i]
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>		if err := ts.processPod(pod, importpaths); err != nil {
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>			return err
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		}
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	}
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	return nil
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>}
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span><span class="comment">// snapshot returns a snapshot of coverage percentage at a moment of</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// time within a running test, so as to support the testing.Coverage()</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span><span class="comment">// function. This version doesn&#39;t examine coverage meta-data, so the</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span><span class="comment">// result it returns will be less accurate (more &#34;slop&#34;) due to the</span>
<span id="L284" class="ln">   284&nbsp;&nbsp;</span><span class="comment">// fact that we don&#39;t look at the meta data to see how many statements</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span><span class="comment">// are associated with each counter.</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>func snapshot() float64 {
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	cl := getCovCounterList()
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	if len(cl) == 0 {
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>		<span class="comment">// no work to do here.</span>
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		return 0.0
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>	}
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	tot := uint64(0)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	totExec := uint64(0)
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	for _, c := range cl {
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>		sd := unsafe.Slice((*atomic.Uint32)(unsafe.Pointer(c.Counters)), c.Len)
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		tot += uint64(len(sd))
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		for i := 0; i &lt; len(sd); i++ {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			<span class="comment">// Skip ahead until the next non-zero value.</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>			if sd[i].Load() == 0 {
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>				continue
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			<span class="comment">// We found a function that was executed.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>			nCtrs := sd[i+coverage.NumCtrsOffset].Load()
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>			cst := i + coverage.FirstCtrOffset
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>			if cst+int(nCtrs) &gt; len(sd) {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				break
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>			}
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			counters := sd[cst : cst+int(nCtrs)]
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			for i := range counters {
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>				if counters[i].Load() != 0 {
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>					totExec++
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>				}
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			}
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>			i += coverage.FirstCtrOffset + int(nCtrs) - 1
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	}
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	if tot == 0 {
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		return 0.0
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	}
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	return float64(totExec) / float64(tot)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>}
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>
</pre><p><a href="testsupport.go?m=text">View as plain text</a></p>

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
