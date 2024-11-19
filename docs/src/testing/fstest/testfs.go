<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/testing/fstest/testfs.go - Go Documentation Server</title>

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
<a href="testfs.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/testing">testing</a>/<a href="http://localhost:8080/src/testing/fstest">fstest</a>/<span class="text-muted">testfs.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/testing/fstest">testing/fstest</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2020 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// Package fstest implements support for testing implementations and users of file systems.</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>package fstest
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>import (
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;errors&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;io/fs&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;path&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;reflect&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;testing/iotest&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// TestFS tests a file system implementation.</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// It walks the entire tree of files in fsys,</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// opening and checking that each file behaves correctly.</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// It also checks that the file system contains at least the expected files.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// As a special case, if no expected files are listed, fsys must be empty.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Otherwise, fsys must contain at least the listed files; it can also contain others.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">// The contents of fsys must not change concurrently with TestFS.</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">// If TestFS finds any misbehaviors, it returns an error reporting all of them.</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">// The error text spans multiple lines, one per detected misbehavior.</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// Typical usage inside a test is:</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//	if err := fstest.TestFS(myFS, &#34;file/that/should/be/present&#34;); err != nil {</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//		t.Fatal(err)</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func TestFS(fsys fs.FS, expected ...string) error {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	if err := testFS(fsys, expected...); err != nil {
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		return err
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	}
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	for _, name := range expected {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		if i := strings.Index(name, &#34;/&#34;); i &gt;= 0 {
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>			dir, dirSlash := name[:i], name[:i+1]
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			var subExpected []string
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			for _, name := range expected {
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>				if strings.HasPrefix(name, dirSlash) {
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>					subExpected = append(subExpected, name[len(dirSlash):])
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>				}
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			}
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			sub, err := fs.Sub(fsys, dir)
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			if err != nil {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>				return err
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>			}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			if err := testFS(sub, subExpected...); err != nil {
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>				return fmt.Errorf(&#34;testing fs.Sub(fsys, %s): %v&#34;, dir, err)
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>			}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>			break <span class="comment">// one sub-test is enough</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	return nil
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>}
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>func testFS(fsys fs.FS, expected ...string) error {
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	t := fsTester{fsys: fsys}
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	t.checkDir(&#34;.&#34;)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	t.checkOpen(&#34;.&#34;)
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	found := make(map[string]bool)
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	for _, dir := range t.dirs {
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		found[dir] = true
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	}
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	for _, file := range t.files {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		found[file] = true
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	delete(found, &#34;.&#34;)
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	if len(expected) == 0 &amp;&amp; len(found) &gt; 0 {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		var list []string
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		for k := range found {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>			if k != &#34;.&#34; {
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>				list = append(list, k)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>			}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		sort.Strings(list)
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		if len(list) &gt; 15 {
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>			list = append(list[:10], &#34;...&#34;)
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		}
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		t.errorf(&#34;expected empty file system but found files:\n%s&#34;, strings.Join(list, &#34;\n&#34;))
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	for _, name := range expected {
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		if !found[name] {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			t.errorf(&#34;expected but not found: %s&#34;, name)
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	if len(t.errText) == 0 {
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		return nil
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	return errors.New(&#34;TestFS found errors:\n&#34; + string(t.errText))
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span><span class="comment">// An fsTester holds state for running the test.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>type fsTester struct {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	fsys    fs.FS
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	errText []byte
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	dirs    []string
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	files   []string
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span><span class="comment">// errorf adds an error line to errText.</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>func (t *fsTester) errorf(format string, args ...any) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	if len(t.errText) &gt; 0 {
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		t.errText = append(t.errText, &#39;\n&#39;)
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	}
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	t.errText = append(t.errText, fmt.Sprintf(format, args...)...)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>func (t *fsTester) openDir(dir string) fs.ReadDirFile {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	f, err := t.fsys.Open(dir)
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	if err != nil {
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Open: %v&#34;, dir, err)
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return nil
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	d, ok := f.(fs.ReadDirFile)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if !ok {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		f.Close()
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Open returned File type %T, not a fs.ReadDirFile&#34;, dir, f)
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		return nil
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	return d
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">// checkDir checks the directory dir, which is expected to exist</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">// (it is either the root or was found in a directory listing with IsDir true).</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>func (t *fsTester) checkDir(dir string) {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">// Read entire directory.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	t.dirs = append(t.dirs, dir)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	d := t.openDir(dir)
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	if d == nil {
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		return
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	list, err := d.ReadDir(-1)
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	if err != nil {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		d.Close()
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		t.errorf(&#34;%s: ReadDir(-1): %v&#34;, dir, err)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		return
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	}
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">// Check all children.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	var prefix string
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	if dir == &#34;.&#34; {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		prefix = &#34;&#34;
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	} else {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		prefix = dir + &#34;/&#34;
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	for _, info := range list {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		name := info.Name()
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		switch {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		case name == &#34;.&#34;, name == &#34;..&#34;, name == &#34;&#34;:
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>			t.errorf(&#34;%s: ReadDir: child has invalid name: %#q&#34;, dir, name)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			continue
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		case strings.Contains(name, &#34;/&#34;):
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			t.errorf(&#34;%s: ReadDir: child name contains slash: %#q&#34;, dir, name)
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			continue
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		case strings.Contains(name, `\`):
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			t.errorf(&#34;%s: ReadDir: child name contains backslash: %#q&#34;, dir, name)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			continue
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		path := prefix + name
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		t.checkStat(path, info)
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		t.checkOpen(path)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		if info.IsDir() {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			t.checkDir(path)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		} else {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			t.checkFile(path)
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">// Check ReadDir(-1) at EOF.</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	list2, err := d.ReadDir(-1)
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	if len(list2) &gt; 0 || err != nil {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		d.Close()
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		t.errorf(&#34;%s: ReadDir(-1) at EOF = %d entries, %v, wanted 0 entries, nil&#34;, dir, len(list2), err)
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		return
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	}
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">// Check ReadDir(1) at EOF (different results).</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	list2, err = d.ReadDir(1)
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	if len(list2) &gt; 0 || err != io.EOF {
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		d.Close()
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		t.errorf(&#34;%s: ReadDir(1) at EOF = %d entries, %v, wanted 0 entries, EOF&#34;, dir, len(list2), err)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		return
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	}
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">// Check that close does not report an error.</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	if err := d.Close(); err != nil {
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Close: %v&#34;, dir, err)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	<span class="comment">// Check that closing twice doesn&#39;t crash.</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	<span class="comment">// The return value doesn&#39;t matter.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	d.Close()
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	<span class="comment">// Reopen directory, read a second time, make sure contents match.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	if d = t.openDir(dir); d == nil {
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		return
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	}
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	defer d.Close()
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	list2, err = d.ReadDir(-1)
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	if err != nil {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		t.errorf(&#34;%s: second Open+ReadDir(-1): %v&#34;, dir, err)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		return
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	t.checkDirList(dir, &#34;first Open+ReadDir(-1) vs second Open+ReadDir(-1)&#34;, list, list2)
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	<span class="comment">// Reopen directory, read a third time in pieces, make sure contents match.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	if d = t.openDir(dir); d == nil {
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		return
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	}
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	defer d.Close()
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	list2 = nil
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	for {
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		n := 1
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		if len(list2) &gt; 0 {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			n = 2
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		frag, err := d.ReadDir(n)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		if len(frag) &gt; n {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			t.errorf(&#34;%s: third Open: ReadDir(%d) after %d: %d entries (too many)&#34;, dir, n, len(list2), len(frag))
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			return
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		}
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		list2 = append(list2, frag...)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		if err == io.EOF {
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			break
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		}
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		if err != nil {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			t.errorf(&#34;%s: third Open: ReadDir(%d) after %d: %v&#34;, dir, n, len(list2), err)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>			return
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>		if n == 0 {
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>			t.errorf(&#34;%s: third Open: ReadDir(%d) after %d: 0 entries but nil error&#34;, dir, n, len(list2))
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			return
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		}
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	t.checkDirList(dir, &#34;first Open+ReadDir(-1) vs third Open+ReadDir(1,2) loop&#34;, list, list2)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	<span class="comment">// If fsys has ReadDir, check that it matches and is sorted.</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	if fsys, ok := t.fsys.(fs.ReadDirFS); ok {
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>		list2, err := fsys.ReadDir(dir)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		if err != nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			t.errorf(&#34;%s: fsys.ReadDir: %v&#34;, dir, err)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			return
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		t.checkDirList(dir, &#34;first Open+ReadDir(-1) vs fsys.ReadDir&#34;, list, list2)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		for i := 0; i+1 &lt; len(list2); i++ {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			if list2[i].Name() &gt;= list2[i+1].Name() {
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>				t.errorf(&#34;%s: fsys.ReadDir: list not sorted: %s before %s&#34;, dir, list2[i].Name(), list2[i+1].Name())
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>			}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>		}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	}
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	<span class="comment">// Check fs.ReadDir as well.</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	list2, err = fs.ReadDir(t.fsys, dir)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	if err != nil {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>		t.errorf(&#34;%s: fs.ReadDir: %v&#34;, dir, err)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		return
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	}
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	t.checkDirList(dir, &#34;first Open+ReadDir(-1) vs fs.ReadDir&#34;, list, list2)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	for i := 0; i+1 &lt; len(list2); i++ {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		if list2[i].Name() &gt;= list2[i+1].Name() {
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			t.errorf(&#34;%s: fs.ReadDir: list not sorted: %s before %s&#34;, dir, list2[i].Name(), list2[i+1].Name())
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		}
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	t.checkGlob(dir, list2)
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
<span id="L276" class="ln">   276&nbsp;&nbsp;</span><span class="comment">// formatEntry formats an fs.DirEntry into a string for error messages and comparison.</span>
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>func formatEntry(entry fs.DirEntry) string {
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%s IsDir=%v Type=%v&#34;, entry.Name(), entry.IsDir(), entry.Type())
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>}
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span><span class="comment">// formatInfoEntry formats an fs.FileInfo into a string like the result of formatEntry, for error messages and comparison.</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>func formatInfoEntry(info fs.FileInfo) string {
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%s IsDir=%v Type=%v&#34;, info.Name(), info.IsDir(), info.Mode().Type())
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>}
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>
<span id="L286" class="ln">   286&nbsp;&nbsp;</span><span class="comment">// formatInfo formats an fs.FileInfo into a string for error messages and comparison.</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>func formatInfo(info fs.FileInfo) string {
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>	return fmt.Sprintf(&#34;%s IsDir=%v Mode=%v Size=%d ModTime=%v&#34;, info.Name(), info.IsDir(), info.Mode(), info.Size(), info.ModTime())
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>
<span id="L291" class="ln">   291&nbsp;&nbsp;</span><span class="comment">// checkGlob checks that various glob patterns work if the file system implements GlobFS.</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>func (t *fsTester) checkGlob(dir string, list []fs.DirEntry) {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>	if _, ok := t.fsys.(fs.GlobFS); !ok {
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>		return
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>	}
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>	<span class="comment">// Make a complex glob pattern prefix that only matches dir.</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	var glob string
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>	if dir != &#34;.&#34; {
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		elem := strings.Split(dir, &#34;/&#34;)
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		for i, e := range elem {
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>			var pattern []rune
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>			for j, r := range e {
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>				if r == &#39;*&#39; || r == &#39;?&#39; || r == &#39;\\&#39; || r == &#39;[&#39; || r == &#39;-&#39; {
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>					pattern = append(pattern, &#39;\\&#39;, r)
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>					continue
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>				}
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>				switch (i + j) % 5 {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>				case 0:
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>					pattern = append(pattern, r)
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>				case 1:
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>					pattern = append(pattern, &#39;[&#39;, r, &#39;]&#39;)
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>				case 2:
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>					pattern = append(pattern, &#39;[&#39;, r, &#39;-&#39;, r, &#39;]&#39;)
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>				case 3:
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>					pattern = append(pattern, &#39;[&#39;, &#39;\\&#39;, r, &#39;]&#39;)
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>				case 4:
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>					pattern = append(pattern, &#39;[&#39;, &#39;\\&#39;, r, &#39;-&#39;, &#39;\\&#39;, r, &#39;]&#39;)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>				}
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>			}
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>			elem[i] = string(pattern)
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		glob = strings.Join(elem, &#34;/&#34;) + &#34;/&#34;
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	}
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	<span class="comment">// Test that malformed patterns are detected.</span>
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>	<span class="comment">// The error is likely path.ErrBadPattern but need not be.</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>	if _, err := t.fsys.(fs.GlobFS).Glob(glob + &#34;nonexist/[]&#34;); err == nil {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Glob(%#q): bad pattern not detected&#34;, dir, glob+&#34;nonexist/[]&#34;)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	}
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	<span class="comment">// Try to find a letter that appears in only some of the final names.</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>	c := rune(&#39;a&#39;)
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>	for ; c &lt;= &#39;z&#39;; c++ {
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		have, haveNot := false, false
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		for _, d := range list {
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>			if strings.ContainsRune(d.Name(), c) {
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>				have = true
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>			} else {
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>				haveNot = true
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>			}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		}
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		if have &amp;&amp; haveNot {
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>			break
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>		}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>	}
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>	if c &gt; &#39;z&#39; {
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>		c = &#39;a&#39;
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	}
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	glob += &#34;*&#34; + string(c) + &#34;*&#34;
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	var want []string
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	for _, d := range list {
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		if strings.ContainsRune(d.Name(), c) {
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>			want = append(want, path.Join(dir, d.Name()))
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		}
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	}
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	names, err := t.fsys.(fs.GlobFS).Glob(glob)
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	if err != nil {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Glob(%#q): %v&#34;, dir, glob, err)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>		return
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	if reflect.DeepEqual(want, names) {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		return
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	if !sort.StringsAreSorted(names) {
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Glob(%#q): unsorted output:\n%s&#34;, dir, glob, strings.Join(names, &#34;\n&#34;))
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>		sort.Strings(names)
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>	var problems []string
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	for len(want) &gt; 0 || len(names) &gt; 0 {
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>		switch {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		case len(want) &gt; 0 &amp;&amp; len(names) &gt; 0 &amp;&amp; want[0] == names[0]:
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>			want, names = want[1:], names[1:]
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>		case len(want) &gt; 0 &amp;&amp; (len(names) == 0 || want[0] &lt; names[0]):
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>			problems = append(problems, &#34;missing: &#34;+want[0])
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>			want = want[1:]
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		default:
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>			problems = append(problems, &#34;extra: &#34;+names[0])
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>			names = names[1:]
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>		}
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	}
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	t.errorf(&#34;%s: Glob(%#q): wrong output:\n%s&#34;, dir, glob, strings.Join(problems, &#34;\n&#34;))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>}
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span><span class="comment">// checkStat checks that a direct stat of path matches entry,</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span><span class="comment">// which was found in the parent&#39;s directory listing.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>func (t *fsTester) checkStat(path string, entry fs.DirEntry) {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	file, err := t.fsys.Open(path)
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	if err != nil {
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Open: %v&#34;, path, err)
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		return
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	}
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>	info, err := file.Stat()
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	file.Close()
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	if err != nil {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Stat: %v&#34;, path, err)
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>		return
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>	}
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	fentry := formatEntry(entry)
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	fientry := formatInfoEntry(info)
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	<span class="comment">// Note: mismatch here is OK for symlink, because Open dereferences symlink.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	if fentry != fientry &amp;&amp; entry.Type()&amp;fs.ModeSymlink == 0 {
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>		t.errorf(&#34;%s: mismatch:\n\tentry = %s\n\tfile.Stat() = %s&#34;, path, fentry, fientry)
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	}
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	einfo, err := entry.Info()
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	if err != nil {
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>		t.errorf(&#34;%s: entry.Info: %v&#34;, path, err)
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>		return
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	}
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>	finfo := formatInfo(info)
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	if entry.Type()&amp;fs.ModeSymlink != 0 {
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>		<span class="comment">// For symlink, just check that entry.Info matches entry on common fields.</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>		<span class="comment">// Open deferences symlink, so info itself may differ.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>		feentry := formatInfoEntry(einfo)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>		if fentry != feentry {
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>			t.errorf(&#34;%s: mismatch\n\tentry = %s\n\tentry.Info() = %s\n&#34;, path, fentry, feentry)
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>		}
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>	} else {
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>		feinfo := formatInfo(einfo)
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>		if feinfo != finfo {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>			t.errorf(&#34;%s: mismatch:\n\tentry.Info() = %s\n\tfile.Stat() = %s\n&#34;, path, feinfo, finfo)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	<span class="comment">// Stat should be the same as Open+Stat, even for symlinks.</span>
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	info2, err := fs.Stat(t.fsys, path)
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>	if err != nil {
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>		t.errorf(&#34;%s: fs.Stat: %v&#34;, path, err)
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>		return
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>	}
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	finfo2 := formatInfo(info2)
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>	if finfo2 != finfo {
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		t.errorf(&#34;%s: fs.Stat(...) = %s\n\twant %s&#34;, path, finfo2, finfo)
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>	}
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	if fsys, ok := t.fsys.(fs.StatFS); ok {
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>		info2, err := fsys.Stat(path)
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>		if err != nil {
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			t.errorf(&#34;%s: fsys.Stat: %v&#34;, path, err)
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			return
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>		}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>		finfo2 := formatInfo(info2)
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>		if finfo2 != finfo {
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			t.errorf(&#34;%s: fsys.Stat(...) = %s\n\twant %s&#34;, path, finfo2, finfo)
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		}
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	}
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>}
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>
<span id="L454" class="ln">   454&nbsp;&nbsp;</span><span class="comment">// checkDirList checks that two directory lists contain the same files and file info.</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span><span class="comment">// The order of the lists need not match.</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>func (t *fsTester) checkDirList(dir, desc string, list1, list2 []fs.DirEntry) {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	old := make(map[string]fs.DirEntry)
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	checkMode := func(entry fs.DirEntry) {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		if entry.IsDir() != (entry.Type()&amp;fs.ModeDir != 0) {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			if entry.IsDir() {
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>				t.errorf(&#34;%s: ReadDir returned %s with IsDir() = true, Type() &amp; ModeDir = 0&#34;, dir, entry.Name())
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>			} else {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>				t.errorf(&#34;%s: ReadDir returned %s with IsDir() = false, Type() &amp; ModeDir = ModeDir&#34;, dir, entry.Name())
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>			}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>		}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	}
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>	for _, entry1 := range list1 {
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>		old[entry1.Name()] = entry1
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>		checkMode(entry1)
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>	}
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>	var diffs []string
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>	for _, entry2 := range list2 {
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>		entry1 := old[entry2.Name()]
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>		if entry1 == nil {
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>			checkMode(entry2)
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>			diffs = append(diffs, &#34;+ &#34;+formatEntry(entry2))
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>			continue
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>		}
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>		if formatEntry(entry1) != formatEntry(entry2) {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>			diffs = append(diffs, &#34;- &#34;+formatEntry(entry1), &#34;+ &#34;+formatEntry(entry2))
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>		}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>		delete(old, entry2.Name())
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>	}
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	for _, entry1 := range old {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>		diffs = append(diffs, &#34;- &#34;+formatEntry(entry1))
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	}
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	if len(diffs) == 0 {
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>		return
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	}
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>	sort.Slice(diffs, func(i, j int) bool {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>		fi := strings.Fields(diffs[i])
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>		fj := strings.Fields(diffs[j])
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>		<span class="comment">// sort by name (i &lt; j) and then +/- (j &lt; i, because + &lt; -)</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		return fi[1]+&#34; &#34;+fj[0] &lt; fj[1]+&#34; &#34;+fi[0]
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>	})
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>	t.errorf(&#34;%s: diff %s:\n\t%s&#34;, dir, desc, strings.Join(diffs, &#34;\n\t&#34;))
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">// checkFile checks that basic file reading works correctly.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>func (t *fsTester) checkFile(file string) {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	t.files = append(t.files, file)
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	<span class="comment">// Read entire file.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	f, err := t.fsys.Open(file)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if err != nil {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Open: %v&#34;, file, err)
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>		return
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	}
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	data, err := io.ReadAll(f)
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	if err != nil {
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		f.Close()
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Open+ReadAll: %v&#34;, file, err)
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>		return
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	}
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	if err := f.Close(); err != nil {
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>		t.errorf(&#34;%s: Close: %v&#34;, file, err)
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	}
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	<span class="comment">// Check that closing twice doesn&#39;t crash.</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	<span class="comment">// The return value doesn&#39;t matter.</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	f.Close()
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	<span class="comment">// Check that ReadFile works if present.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>	if fsys, ok := t.fsys.(fs.ReadFileFS); ok {
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		data2, err := fsys.ReadFile(file)
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		if err != nil {
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>			t.errorf(&#34;%s: fsys.ReadFile: %v&#34;, file, err)
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>			return
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>		}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>		t.checkFileRead(file, &#34;ReadAll vs fsys.ReadFile&#34;, data, data2)
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		<span class="comment">// Modify the data and check it again. Modifying the</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>		<span class="comment">// returned byte slice should not affect the next call.</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>		for i := range data2 {
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			data2[i]++
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>		}
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>		data2, err = fsys.ReadFile(file)
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>		if err != nil {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>			t.errorf(&#34;%s: second call to fsys.ReadFile: %v&#34;, file, err)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>			return
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>		}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		t.checkFileRead(file, &#34;Readall vs second fsys.ReadFile&#34;, data, data2)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		t.checkBadPath(file, &#34;ReadFile&#34;,
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>			func(name string) error { _, err := fsys.ReadFile(name); return err })
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	}
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	<span class="comment">// Check that fs.ReadFile works with t.fsys.</span>
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	data2, err := fs.ReadFile(t.fsys, file)
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	if err != nil {
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		t.errorf(&#34;%s: fs.ReadFile: %v&#34;, file, err)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>		return
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	}
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	t.checkFileRead(file, &#34;ReadAll vs fs.ReadFile&#34;, data, data2)
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	<span class="comment">// Use iotest.TestReader to check small reads, Seek, ReadAt.</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>	f, err = t.fsys.Open(file)
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	if err != nil {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		t.errorf(&#34;%s: second Open: %v&#34;, file, err)
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		return
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	}
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	defer f.Close()
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	if err := iotest.TestReader(f, data); err != nil {
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>		t.errorf(&#34;%s: failed TestReader:\n\t%s&#34;, file, strings.ReplaceAll(err.Error(), &#34;\n&#34;, &#34;\n\t&#34;))
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	}
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>}
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>func (t *fsTester) checkFileRead(file, desc string, data1, data2 []byte) {
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	if string(data1) != string(data2) {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		t.errorf(&#34;%s: %s: different data returned\n\t%q\n\t%q&#34;, file, desc, data1, data2)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		return
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>	}
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span><span class="comment">// checkBadPath checks that various invalid forms of file&#39;s name cannot be opened using t.fsys.Open.</span>
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>func (t *fsTester) checkOpen(file string) {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	t.checkBadPath(file, &#34;Open&#34;, func(file string) error {
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>		f, err := t.fsys.Open(file)
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>		if err == nil {
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>			f.Close()
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		}
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>		return err
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>	})
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>}
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span><span class="comment">// checkBadPath checks that various invalid forms of file&#39;s name cannot be opened using open.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>func (t *fsTester) checkBadPath(file string, desc string, open func(string) error) {
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	bad := []string{
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>		&#34;/&#34; + file,
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		file + &#34;/.&#34;,
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	if file == &#34;.&#34; {
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>		bad = append(bad, &#34;/&#34;)
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	}
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	if i := strings.Index(file, &#34;/&#34;); i &gt;= 0 {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		bad = append(bad,
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>			file[:i]+&#34;//&#34;+file[i+1:],
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>			file[:i]+&#34;/./&#34;+file[i+1:],
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>			file[:i]+`\`+file[i+1:],
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>			file[:i]+&#34;/../&#34;+file,
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>		)
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	}
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	if i := strings.LastIndex(file, &#34;/&#34;); i &gt;= 0 {
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>		bad = append(bad,
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>			file[:i]+&#34;//&#34;+file[i+1:],
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>			file[:i]+&#34;/./&#34;+file[i+1:],
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>			file[:i]+`\`+file[i+1:],
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>			file+&#34;/../&#34;+file[i+1:],
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	for _, b := range bad {
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>		if err := open(b); err == nil {
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>			t.errorf(&#34;%s: %s(%s) succeeded, want error&#34;, file, desc, b)
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>		}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>	}
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>}
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>
</pre><p><a href="testfs.go?m=text">View as plain text</a></p>

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
