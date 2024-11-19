<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/archive/tar/reader.go - Go Documentation Server</title>

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
<a href="reader.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/archive">archive</a>/<a href="http://localhost:8080/src/archive/tar">tar</a>/<span class="text-muted">reader.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/archive/tar">archive/tar</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package tar
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;bytes&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;io&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;path/filepath&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;strconv&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>)
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// Reader provides sequential access to the contents of a tar archive.</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// Reader.Next advances to the next file in the archive (including the first),</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// and then Reader can be treated as an io.Reader to access the file&#39;s data.</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>type Reader struct {
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	r    io.Reader
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	pad  int64      <span class="comment">// Amount of padding (ignored) after current file entry</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	curr fileReader <span class="comment">// Reader for current file entry</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	blk  block      <span class="comment">// Buffer to use as temporary local storage</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	<span class="comment">// err is a persistent error.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	<span class="comment">// It is only the responsibility of every exported method of Reader to</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	<span class="comment">// ensure that this error is sticky.</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	err error
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>type fileReader interface {
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	io.Reader
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	fileState
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	WriteTo(io.Writer) (int64, error)
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>}
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">// NewReader creates a new [Reader] reading from r.</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>func NewReader(r io.Reader) *Reader {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	return &amp;Reader{r: r, curr: &amp;regFileReader{r, 0}}
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>}
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">// Next advances to the next entry in the tar archive.</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">// The Header.Size determines how many bytes can be read for the next file.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">// Any remaining data in the current file is automatically discarded.</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">// At the end of the archive, Next returns the error io.EOF.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// If Next encounters a non-local name (as defined by [filepath.IsLocal])</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">// and the GODEBUG environment variable contains `tarinsecurepath=0`,</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">// Next returns the header with an [ErrInsecurePath] error.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">// A future version of Go may introduce this behavior by default.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">// Programs that want to accept non-local names can ignore</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">// the [ErrInsecurePath] error and use the returned header.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>func (tr *Reader) Next() (*Header, error) {
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	if tr.err != nil {
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		return nil, tr.err
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	}
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	hdr, err := tr.next()
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	tr.err = err
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if err == nil &amp;&amp; !filepath.IsLocal(hdr.Name) {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		if tarinsecurepath.Value() == &#34;0&#34; {
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>			tarinsecurepath.IncNonDefault()
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>			err = ErrInsecurePath
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		}
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	}
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	return hdr, err
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>}
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>func (tr *Reader) next() (*Header, error) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	var paxHdrs map[string]string
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	var gnuLongName, gnuLongLink string
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	<span class="comment">// Externally, Next iterates through the tar archive as if it is a series of</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	<span class="comment">// files. Internally, the tar format often uses fake &#34;files&#34; to add meta</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	<span class="comment">// data that describes the next file. These meta data &#34;files&#34; should not</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// normally be visible to the outside. As such, this loop iterates through</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">// one or more &#34;header files&#34; until it finds a &#34;normal file&#34;.</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	format := FormatUSTAR | FormatPAX | FormatGNU
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	for {
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		<span class="comment">// Discard the remainder of the file and any padding.</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		if err := discard(tr.r, tr.curr.physicalRemaining()); err != nil {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>			return nil, err
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		if _, err := tryReadFull(tr.r, tr.blk[:tr.pad]); err != nil {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>			return nil, err
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		}
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		tr.pad = 0
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		hdr, rawHdr, err := tr.readHeader()
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		if err != nil {
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			return nil, err
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		}
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		if err := tr.handleRegularFile(hdr); err != nil {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			return nil, err
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		format.mayOnlyBe(hdr.Format)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		<span class="comment">// Check for PAX/GNU special headers and files.</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		switch hdr.Typeflag {
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		case TypeXHeader, TypeXGlobalHeader:
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>			format.mayOnlyBe(FormatPAX)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>			paxHdrs, err = parsePAX(tr)
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			if err != nil {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>				return nil, err
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			}
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			if hdr.Typeflag == TypeXGlobalHeader {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>				mergePAX(hdr, paxHdrs)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>				return &amp;Header{
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>					Name:       hdr.Name,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>					Typeflag:   hdr.Typeflag,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>					Xattrs:     hdr.Xattrs,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>					PAXRecords: hdr.PAXRecords,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>					Format:     format,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				}, nil
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>			}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			continue <span class="comment">// This is a meta header affecting the next header</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		case TypeGNULongName, TypeGNULongLink:
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			format.mayOnlyBe(FormatGNU)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			realname, err := readSpecialFile(tr)
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>			if err != nil {
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>				return nil, err
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>			}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>			var p parser
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>			switch hdr.Typeflag {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			case TypeGNULongName:
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				gnuLongName = p.parseString(realname)
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>			case TypeGNULongLink:
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>				gnuLongLink = p.parseString(realname)
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>			}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			continue <span class="comment">// This is a meta header affecting the next header</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		default:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			<span class="comment">// The old GNU sparse format is handled here since it is technically</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			<span class="comment">// just a regular file with additional attributes.</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>			if err := mergePAX(hdr, paxHdrs); err != nil {
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>				return nil, err
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>			if gnuLongName != &#34;&#34; {
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				hdr.Name = gnuLongName
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>			}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>			if gnuLongLink != &#34;&#34; {
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>				hdr.Linkname = gnuLongLink
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>			}
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>			if hdr.Typeflag == TypeRegA {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>				if strings.HasSuffix(hdr.Name, &#34;/&#34;) {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>					hdr.Typeflag = TypeDir <span class="comment">// Legacy archives use trailing slash for directories</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				} else {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>					hdr.Typeflag = TypeReg
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				}
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>			}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>			<span class="comment">// The extended headers may have updated the size.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>			<span class="comment">// Thus, setup the regFileReader again after merging PAX headers.</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>			if err := tr.handleRegularFile(hdr); err != nil {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>				return nil, err
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>			}
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			<span class="comment">// Sparse formats rely on being able to read from the logical data</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			<span class="comment">// section; there must be a preceding call to handleRegularFile.</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			if err := tr.handleSparseFile(hdr, rawHdr); err != nil {
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>				return nil, err
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>			}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			<span class="comment">// Set the final guess at the format.</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			if format.has(FormatUSTAR) &amp;&amp; format.has(FormatPAX) {
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>				format.mayOnlyBe(FormatUSTAR)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			}
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			hdr.Format = format
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			return hdr, nil <span class="comment">// This is a file, so stop</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span><span class="comment">// handleRegularFile sets up the current file reader and padding such that it</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span><span class="comment">// can only read the following logical data section. It will properly handle</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span><span class="comment">// special headers that contain no data section.</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>func (tr *Reader) handleRegularFile(hdr *Header) error {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	nb := hdr.Size
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	if isHeaderOnlyType(hdr.Typeflag) {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		nb = 0
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	}
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	if nb &lt; 0 {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		return ErrHeader
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	}
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	tr.pad = blockPadding(nb)
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	tr.curr = &amp;regFileReader{r: tr.r, nb: nb}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	return nil
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>}
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span><span class="comment">// handleSparseFile checks if the current file is a sparse format of any type</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span><span class="comment">// and sets the curr reader appropriately.</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>func (tr *Reader) handleSparseFile(hdr *Header, rawHdr *block) error {
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	var spd sparseDatas
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	var err error
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if hdr.Typeflag == TypeGNUSparse {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		spd, err = tr.readOldGNUSparseMap(hdr, rawHdr)
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	} else {
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		spd, err = tr.readGNUSparsePAXHeaders(hdr)
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	<span class="comment">// If sp is non-nil, then this is a sparse file.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	<span class="comment">// Note that it is possible for len(sp) == 0.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	if err == nil &amp;&amp; spd != nil {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if isHeaderOnlyType(hdr.Typeflag) || !validateSparseEntries(spd, hdr.Size) {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			return ErrHeader
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		sph := invertSparseEntries(spd, hdr.Size)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		tr.curr = &amp;sparseFileReader{tr.curr, sph, 0}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	return err
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">// readGNUSparsePAXHeaders checks the PAX headers for GNU sparse headers.</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">// If they are found, then this function reads the sparse map and returns it.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">// This assumes that 0.0 headers have already been converted to 0.1 headers</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">// by the PAX header parsing logic.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>func (tr *Reader) readGNUSparsePAXHeaders(hdr *Header) (sparseDatas, error) {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	<span class="comment">// Identify the version of GNU headers.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	var is1x0 bool
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	major, minor := hdr.PAXRecords[paxGNUSparseMajor], hdr.PAXRecords[paxGNUSparseMinor]
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	switch {
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	case major == &#34;0&#34; &amp;&amp; (minor == &#34;0&#34; || minor == &#34;1&#34;):
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		is1x0 = false
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	case major == &#34;1&#34; &amp;&amp; minor == &#34;0&#34;:
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		is1x0 = true
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	case major != &#34;&#34; || minor != &#34;&#34;:
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		return nil, nil <span class="comment">// Unknown GNU sparse PAX version</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	case hdr.PAXRecords[paxGNUSparseMap] != &#34;&#34;:
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		is1x0 = false <span class="comment">// 0.0 and 0.1 did not have explicit version records, so guess</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	default:
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		return nil, nil <span class="comment">// Not a PAX format GNU sparse file.</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	hdr.Format.mayOnlyBe(FormatPAX)
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	<span class="comment">// Update hdr from GNU sparse PAX headers.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	if name := hdr.PAXRecords[paxGNUSparseName]; name != &#34;&#34; {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>		hdr.Name = name
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	}
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	size := hdr.PAXRecords[paxGNUSparseSize]
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	if size == &#34;&#34; {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>		size = hdr.PAXRecords[paxGNUSparseRealSize]
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if size != &#34;&#34; {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		n, err := strconv.ParseInt(size, 10, 64)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		if err != nil {
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>			return nil, ErrHeader
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>		hdr.Size = n
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	}
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	<span class="comment">// Read the sparse map according to the appropriate format.</span>
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	if is1x0 {
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		return readGNUSparseMap1x0(tr.curr)
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	}
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	return readGNUSparseMap0x1(hdr.PAXRecords)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>}
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>
<span id="L260" class="ln">   260&nbsp;&nbsp;</span><span class="comment">// mergePAX merges paxHdrs into hdr for all relevant fields of Header.</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>func mergePAX(hdr *Header, paxHdrs map[string]string) (err error) {
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	for k, v := range paxHdrs {
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		if v == &#34;&#34; {
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>			continue <span class="comment">// Keep the original USTAR value</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		}
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		var id64 int64
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		switch k {
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>		case paxPath:
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>			hdr.Name = v
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		case paxLinkpath:
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>			hdr.Linkname = v
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>		case paxUname:
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>			hdr.Uname = v
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		case paxGname:
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>			hdr.Gname = v
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		case paxUid:
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>			id64, err = strconv.ParseInt(v, 10, 64)
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>			hdr.Uid = int(id64) <span class="comment">// Integer overflow possible</span>
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>		case paxGid:
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>			id64, err = strconv.ParseInt(v, 10, 64)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>			hdr.Gid = int(id64) <span class="comment">// Integer overflow possible</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		case paxAtime:
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>			hdr.AccessTime, err = parsePAXTime(v)
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		case paxMtime:
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>			hdr.ModTime, err = parsePAXTime(v)
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		case paxCtime:
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>			hdr.ChangeTime, err = parsePAXTime(v)
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		case paxSize:
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>			hdr.Size, err = strconv.ParseInt(v, 10, 64)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		default:
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>			if strings.HasPrefix(k, paxSchilyXattr) {
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>				if hdr.Xattrs == nil {
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>					hdr.Xattrs = make(map[string]string)
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>				}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>				hdr.Xattrs[k[len(paxSchilyXattr):]] = v
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>			}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		}
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>		if err != nil {
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>			return ErrHeader
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>		}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>	}
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	hdr.PAXRecords = paxHdrs
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>	return nil
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>
<span id="L306" class="ln">   306&nbsp;&nbsp;</span><span class="comment">// parsePAX parses PAX headers.</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span><span class="comment">// If an extended header (type &#39;x&#39;) is invalid, ErrHeader is returned.</span>
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>func parsePAX(r io.Reader) (map[string]string, error) {
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>	buf, err := readSpecialFile(r)
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>	if err != nil {
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>		return nil, err
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>	}
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>	sbuf := string(buf)
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	<span class="comment">// For GNU PAX sparse format 0.0 support.</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>	<span class="comment">// This function transforms the sparse format 0.0 headers into format 0.1</span>
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	<span class="comment">// headers since 0.0 headers were not PAX compliant.</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>	var sparseMap []string
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	paxHdrs := make(map[string]string)
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	for len(sbuf) &gt; 0 {
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		key, value, residual, err := parsePAXRecord(sbuf)
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		if err != nil {
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>			return nil, ErrHeader
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		}
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		sbuf = residual
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		switch key {
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		case paxGNUSparseOffset, paxGNUSparseNumBytes:
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>			<span class="comment">// Validate sparse header order and value.</span>
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>			if (len(sparseMap)%2 == 0 &amp;&amp; key != paxGNUSparseOffset) ||
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>				(len(sparseMap)%2 == 1 &amp;&amp; key != paxGNUSparseNumBytes) ||
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>				strings.Contains(value, &#34;,&#34;) {
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>				return nil, ErrHeader
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>			}
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>			sparseMap = append(sparseMap, value)
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		default:
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>			paxHdrs[key] = value
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		}
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	}
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	if len(sparseMap) &gt; 0 {
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		paxHdrs[paxGNUSparseMap] = strings.Join(sparseMap, &#34;,&#34;)
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>	}
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	return paxHdrs, nil
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>}
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">// readHeader reads the next block header and assumes that the underlying reader</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">// is already aligned to a block boundary. It returns the raw block of the</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">// header in case further processing is required.</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">// The err will be set to io.EOF only when one of the following occurs:</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//   - Exactly 0 bytes are read and EOF is hit.</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">//   - Exactly 1 block of zeros is read and EOF is hit.</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">//   - At least 2 blocks of zeros are read.</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>func (tr *Reader) readHeader() (*Header, *block, error) {
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>	<span class="comment">// Two blocks of zero bytes marks the end of the archive.</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>	if _, err := io.ReadFull(tr.r, tr.blk[:]); err != nil {
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>		return nil, nil, err <span class="comment">// EOF is okay here; exactly 0 bytes read</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>	}
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	if bytes.Equal(tr.blk[:], zeroBlock[:]) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		if _, err := io.ReadFull(tr.r, tr.blk[:]); err != nil {
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>			return nil, nil, err <span class="comment">// EOF is okay here; exactly 1 block of zeros read</span>
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		}
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>		if bytes.Equal(tr.blk[:], zeroBlock[:]) {
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			return nil, nil, io.EOF <span class="comment">// normal EOF; exactly 2 block of zeros read</span>
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		}
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		return nil, nil, ErrHeader <span class="comment">// Zero block and then non-zero block</span>
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>	<span class="comment">// Verify the header matches a known format.</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>	format := tr.blk.getFormat()
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	if format == FormatUnknown {
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>		return nil, nil, ErrHeader
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	}
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	var p parser
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	hdr := new(Header)
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	<span class="comment">// Unpack the V7 header.</span>
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	v7 := tr.blk.toV7()
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	hdr.Typeflag = v7.typeFlag()[0]
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	hdr.Name = p.parseString(v7.name())
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	hdr.Linkname = p.parseString(v7.linkName())
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	hdr.Size = p.parseNumeric(v7.size())
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	hdr.Mode = p.parseNumeric(v7.mode())
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	hdr.Uid = int(p.parseNumeric(v7.uid()))
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>	hdr.Gid = int(p.parseNumeric(v7.gid()))
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	hdr.ModTime = time.Unix(p.parseNumeric(v7.modTime()), 0)
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	<span class="comment">// Unpack format specific fields.</span>
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	if format &gt; formatV7 {
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>		ustar := tr.blk.toUSTAR()
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>		hdr.Uname = p.parseString(ustar.userName())
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>		hdr.Gname = p.parseString(ustar.groupName())
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		hdr.Devmajor = p.parseNumeric(ustar.devMajor())
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		hdr.Devminor = p.parseNumeric(ustar.devMinor())
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		var prefix string
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>		switch {
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>		case format.has(FormatUSTAR | FormatPAX):
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>			hdr.Format = format
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>			ustar := tr.blk.toUSTAR()
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>			prefix = p.parseString(ustar.prefix())
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>			<span class="comment">// For Format detection, check if block is properly formatted since</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>			<span class="comment">// the parser is more liberal than what USTAR actually permits.</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>			notASCII := func(r rune) bool { return r &gt;= 0x80 }
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>			if bytes.IndexFunc(tr.blk[:], notASCII) &gt;= 0 {
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>				hdr.Format = FormatUnknown <span class="comment">// Non-ASCII characters in block.</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>			}
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>			nul := func(b []byte) bool { return int(b[len(b)-1]) == 0 }
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>			if !(nul(v7.size()) &amp;&amp; nul(v7.mode()) &amp;&amp; nul(v7.uid()) &amp;&amp; nul(v7.gid()) &amp;&amp;
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>				nul(v7.modTime()) &amp;&amp; nul(ustar.devMajor()) &amp;&amp; nul(ustar.devMinor())) {
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>				hdr.Format = FormatUnknown <span class="comment">// Numeric fields must end in NUL</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>			}
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>		case format.has(formatSTAR):
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>			star := tr.blk.toSTAR()
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>			prefix = p.parseString(star.prefix())
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>			hdr.AccessTime = time.Unix(p.parseNumeric(star.accessTime()), 0)
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>			hdr.ChangeTime = time.Unix(p.parseNumeric(star.changeTime()), 0)
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>		case format.has(FormatGNU):
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>			hdr.Format = format
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>			var p2 parser
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>			gnu := tr.blk.toGNU()
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>			if b := gnu.accessTime(); b[0] != 0 {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>				hdr.AccessTime = time.Unix(p2.parseNumeric(b), 0)
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>			}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>			if b := gnu.changeTime(); b[0] != 0 {
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>				hdr.ChangeTime = time.Unix(p2.parseNumeric(b), 0)
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>			}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>			<span class="comment">// Prior to Go1.8, the Writer had a bug where it would output</span>
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			<span class="comment">// an invalid tar file in certain rare situations because the logic</span>
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			<span class="comment">// incorrectly believed that the old GNU format had a prefix field.</span>
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			<span class="comment">// This is wrong and leads to an output file that mangles the</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			<span class="comment">// atime and ctime fields, which are often left unused.</span>
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>			<span class="comment">// In order to continue reading tar files created by former, buggy</span>
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			<span class="comment">// versions of Go, we skeptically parse the atime and ctime fields.</span>
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>			<span class="comment">// If we are unable to parse them and the prefix field looks like</span>
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>			<span class="comment">// an ASCII string, then we fallback on the pre-Go1.8 behavior</span>
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>			<span class="comment">// of treating these fields as the USTAR prefix field.</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>			<span class="comment">// Note that this will not use the fallback logic for all possible</span>
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>			<span class="comment">// files generated by a pre-Go1.8 toolchain. If the generated file</span>
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>			<span class="comment">// happened to have a prefix field that parses as valid</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>			<span class="comment">// atime and ctime fields (e.g., when they are valid octal strings),</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>			<span class="comment">// then it is impossible to distinguish between a valid GNU file</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>			<span class="comment">// and an invalid pre-Go1.8 file.</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>			<span class="comment">//</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>			<span class="comment">// See https://golang.org/issues/12594</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>			<span class="comment">// See https://golang.org/issues/21005</span>
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>			if p2.err != nil {
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>				hdr.AccessTime, hdr.ChangeTime = time.Time{}, time.Time{}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>				ustar := tr.blk.toUSTAR()
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>				if s := p.parseString(ustar.prefix()); isASCII(s) {
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>					prefix = s
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>				}
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>				hdr.Format = FormatUnknown <span class="comment">// Buggy file is not GNU</span>
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>			}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>		}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>		if len(prefix) &gt; 0 {
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>			hdr.Name = prefix + &#34;/&#34; + hdr.Name
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>		}
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>	}
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>	return hdr, &amp;tr.blk, p.err
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>}
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">// readOldGNUSparseMap reads the sparse map from the old GNU sparse format.</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">// The sparse map is stored in the tar header if it&#39;s small enough.</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">// If it&#39;s larger than four entries, then one or more extension headers are used</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">// to store the rest of the sparse map.</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">// The Header.Size does not reflect the size of any extended headers used.</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">// Thus, this function will read from the raw io.Reader to fetch extra headers.</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">// This method mutates blk in the process.</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>func (tr *Reader) readOldGNUSparseMap(hdr *Header, blk *block) (sparseDatas, error) {
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>	<span class="comment">// Make sure that the input format is GNU.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>	<span class="comment">// Unfortunately, the STAR format also has a sparse header format that uses</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>	<span class="comment">// the same type flag but has a completely different layout.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>	if blk.getFormat() != FormatGNU {
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>		return nil, ErrHeader
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>	}
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>	hdr.Format.mayOnlyBe(FormatGNU)
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>	var p parser
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	hdr.Size = p.parseNumeric(blk.toGNU().realSize())
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>	if p.err != nil {
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>		return nil, p.err
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>	}
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	s := blk.toGNU().sparse()
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	spd := make(sparseDatas, 0, s.maxEntries())
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>	for {
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>		for i := 0; i &lt; s.maxEntries(); i++ {
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			<span class="comment">// This termination condition is identical to GNU and BSD tar.</span>
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			if s.entry(i).offset()[0] == 0x00 {
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>				break <span class="comment">// Don&#39;t return, need to process extended headers (even if empty)</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>			}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>			offset := p.parseNumeric(s.entry(i).offset())
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			length := p.parseNumeric(s.entry(i).length())
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>			if p.err != nil {
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>				return nil, p.err
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>			}
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>			spd = append(spd, sparseEntry{Offset: offset, Length: length})
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>		}
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		if s.isExtended()[0] &gt; 0 {
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>			<span class="comment">// There are more entries. Read an extension header and parse its entries.</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>			if _, err := mustReadFull(tr.r, blk[:]); err != nil {
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>				return nil, err
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>			}
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>			s = blk.toSparse()
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>			continue
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>		}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>		return spd, nil <span class="comment">// Done</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>	}
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>}
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>
<span id="L519" class="ln">   519&nbsp;&nbsp;</span><span class="comment">// readGNUSparseMap1x0 reads the sparse map as stored in GNU&#39;s PAX sparse format</span>
<span id="L520" class="ln">   520&nbsp;&nbsp;</span><span class="comment">// version 1.0. The format of the sparse map consists of a series of</span>
<span id="L521" class="ln">   521&nbsp;&nbsp;</span><span class="comment">// newline-terminated numeric fields. The first field is the number of entries</span>
<span id="L522" class="ln">   522&nbsp;&nbsp;</span><span class="comment">// and is always present. Following this are the entries, consisting of two</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span><span class="comment">// fields (offset, length). This function must stop reading at the end</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span><span class="comment">// boundary of the block containing the last newline.</span>
<span id="L525" class="ln">   525&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L526" class="ln">   526&nbsp;&nbsp;</span><span class="comment">// Note that the GNU manual says that numeric values should be encoded in octal</span>
<span id="L527" class="ln">   527&nbsp;&nbsp;</span><span class="comment">// format. However, the GNU tar utility itself outputs these values in decimal.</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span><span class="comment">// As such, this library treats values as being encoded in decimal.</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>func readGNUSparseMap1x0(r io.Reader) (sparseDatas, error) {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>	var (
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		cntNewline int64
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		buf        bytes.Buffer
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		blk        block
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>	)
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	<span class="comment">// feedTokens copies data in blocks from r into buf until there are</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>	<span class="comment">// at least cnt newlines in buf. It will not read more blocks than needed.</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	feedTokens := func(n int64) error {
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>		for cntNewline &lt; n {
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>			if _, err := mustReadFull(r, blk[:]); err != nil {
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>				return err
<span id="L542" class="ln">   542&nbsp;&nbsp;</span>			}
<span id="L543" class="ln">   543&nbsp;&nbsp;</span>			buf.Write(blk[:])
<span id="L544" class="ln">   544&nbsp;&nbsp;</span>			for _, c := range blk {
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>				if c == &#39;\n&#39; {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>					cntNewline++
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>				}
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>			}
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		}
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>		return nil
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>	}
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	<span class="comment">// nextToken gets the next token delimited by a newline. This assumes that</span>
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	<span class="comment">// at least one newline exists in the buffer.</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>	nextToken := func() string {
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		cntNewline--
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>		tok, _ := buf.ReadString(&#39;\n&#39;)
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>		return strings.TrimRight(tok, &#34;\n&#34;)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	}
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	<span class="comment">// Parse for the number of entries.</span>
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>	<span class="comment">// Use integer overflow resistant math to check this.</span>
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>	if err := feedTokens(1); err != nil {
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>		return nil, err
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>	}
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	numEntries, err := strconv.ParseInt(nextToken(), 10, 0) <span class="comment">// Intentionally parse as native int</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	if err != nil || numEntries &lt; 0 || int(2*numEntries) &lt; int(numEntries) {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>		return nil, ErrHeader
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>	<span class="comment">// Parse for all member entries.</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	<span class="comment">// numEntries is trusted after this since a potential attacker must have</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	<span class="comment">// committed resources proportional to what this library used.</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>	if err := feedTokens(2 * numEntries); err != nil {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>		return nil, err
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	}
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	spd := make(sparseDatas, 0, numEntries)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	for i := int64(0); i &lt; numEntries; i++ {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		offset, err1 := strconv.ParseInt(nextToken(), 10, 64)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>		length, err2 := strconv.ParseInt(nextToken(), 10, 64)
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>		if err1 != nil || err2 != nil {
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>			return nil, ErrHeader
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		}
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		spd = append(spd, sparseEntry{Offset: offset, Length: length})
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>	return spd, nil
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>}
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// readGNUSparseMap0x1 reads the sparse map as stored in GNU&#39;s PAX sparse format</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// version 0.1. The sparse map is stored in the PAX headers.</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>func readGNUSparseMap0x1(paxHdrs map[string]string) (sparseDatas, error) {
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	<span class="comment">// Get number of entries.</span>
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	<span class="comment">// Use integer overflow resistant math to check this.</span>
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	numEntriesStr := paxHdrs[paxGNUSparseNumBlocks]
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>	numEntries, err := strconv.ParseInt(numEntriesStr, 10, 0) <span class="comment">// Intentionally parse as native int</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	if err != nil || numEntries &lt; 0 || int(2*numEntries) &lt; int(numEntries) {
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>		return nil, ErrHeader
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	}
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	<span class="comment">// There should be two numbers in sparseMap for each entry.</span>
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	sparseMap := strings.Split(paxHdrs[paxGNUSparseMap], &#34;,&#34;)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	if len(sparseMap) == 1 &amp;&amp; sparseMap[0] == &#34;&#34; {
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>		sparseMap = sparseMap[:0]
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	}
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	if int64(len(sparseMap)) != 2*numEntries {
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>		return nil, ErrHeader
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	}
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	<span class="comment">// Loop through the entries in the sparse map.</span>
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>	<span class="comment">// numEntries is trusted now.</span>
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	spd := make(sparseDatas, 0, numEntries)
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	for len(sparseMap) &gt;= 2 {
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>		offset, err1 := strconv.ParseInt(sparseMap[0], 10, 64)
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>		length, err2 := strconv.ParseInt(sparseMap[1], 10, 64)
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>		if err1 != nil || err2 != nil {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>			return nil, ErrHeader
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>		}
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		spd = append(spd, sparseEntry{Offset: offset, Length: length})
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>		sparseMap = sparseMap[2:]
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	}
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>	return spd, nil
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>}
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span><span class="comment">// Read reads from the current file in the tar archive.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span><span class="comment">// It returns (0, io.EOF) when it reaches the end of that file,</span>
<span id="L626" class="ln">   626&nbsp;&nbsp;</span><span class="comment">// until [Next] is called to advance to the next file.</span>
<span id="L627" class="ln">   627&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span><span class="comment">// If the current file is sparse, then the regions marked as a hole</span>
<span id="L629" class="ln">   629&nbsp;&nbsp;</span><span class="comment">// are read back as NUL-bytes.</span>
<span id="L630" class="ln">   630&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// Calling Read on special types like [TypeLink], [TypeSymlink], [TypeChar],</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">// [TypeBlock], [TypeDir], and [TypeFifo] returns (0, [io.EOF]) regardless of what</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// the [Header.Size] claims.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>func (tr *Reader) Read(b []byte) (int, error) {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	if tr.err != nil {
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>		return 0, tr.err
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>	}
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	n, err := tr.curr.Read(b)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	if err != nil &amp;&amp; err != io.EOF {
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>		tr.err = err
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	}
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	return n, err
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>}
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span><span class="comment">// writeTo writes the content of the current file to w.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span><span class="comment">// The bytes written matches the number of remaining bytes in the current file.</span>
<span id="L647" class="ln">   647&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span><span class="comment">// If the current file is sparse and w is an io.WriteSeeker,</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span><span class="comment">// then writeTo uses Seek to skip past holes defined in Header.SparseHoles,</span>
<span id="L650" class="ln">   650&nbsp;&nbsp;</span><span class="comment">// assuming that skipped regions are filled with NULs.</span>
<span id="L651" class="ln">   651&nbsp;&nbsp;</span><span class="comment">// This always writes the last byte to ensure w is the right size.</span>
<span id="L652" class="ln">   652&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L653" class="ln">   653&nbsp;&nbsp;</span><span class="comment">// TODO(dsnet): Re-export this when adding sparse file support.</span>
<span id="L654" class="ln">   654&nbsp;&nbsp;</span><span class="comment">// See https://golang.org/issue/22735</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>func (tr *Reader) writeTo(w io.Writer) (int64, error) {
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>	if tr.err != nil {
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		return 0, tr.err
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>	}
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>	n, err := tr.curr.WriteTo(w)
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>	if err != nil {
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>		tr.err = err
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>	}
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>	return n, err
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>}
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>
<span id="L666" class="ln">   666&nbsp;&nbsp;</span><span class="comment">// regFileReader is a fileReader for reading data from a regular file entry.</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>type regFileReader struct {
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>	r  io.Reader <span class="comment">// Underlying Reader</span>
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>	nb int64     <span class="comment">// Number of remaining bytes to read</span>
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>}
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>func (fr *regFileReader) Read(b []byte) (n int, err error) {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>	if int64(len(b)) &gt; fr.nb {
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>		b = b[:fr.nb]
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>	}
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>	if len(b) &gt; 0 {
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>		n, err = fr.r.Read(b)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>		fr.nb -= int64(n)
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>	}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>	switch {
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>	case err == io.EOF &amp;&amp; fr.nb &gt; 0:
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		return n, io.ErrUnexpectedEOF
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>	case err == nil &amp;&amp; fr.nb == 0:
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		return n, io.EOF
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>	default:
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		return n, err
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>	}
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>}
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>func (fr *regFileReader) WriteTo(w io.Writer) (int64, error) {
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>	return io.Copy(w, struct{ io.Reader }{fr})
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>}
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>
<span id="L694" class="ln">   694&nbsp;&nbsp;</span><span class="comment">// logicalRemaining implements fileState.logicalRemaining.</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>func (fr regFileReader) logicalRemaining() int64 {
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>	return fr.nb
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>}
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>
<span id="L699" class="ln">   699&nbsp;&nbsp;</span><span class="comment">// physicalRemaining implements fileState.physicalRemaining.</span>
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>func (fr regFileReader) physicalRemaining() int64 {
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>	return fr.nb
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>}
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span><span class="comment">// sparseFileReader is a fileReader for reading data from a sparse file entry.</span>
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>type sparseFileReader struct {
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>	fr  fileReader  <span class="comment">// Underlying fileReader</span>
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>	sp  sparseHoles <span class="comment">// Normalized list of sparse holes</span>
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	pos int64       <span class="comment">// Current position in sparse file</span>
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>}
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>func (sr *sparseFileReader) Read(b []byte) (n int, err error) {
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>	finished := int64(len(b)) &gt;= sr.logicalRemaining()
<span id="L713" class="ln">   713&nbsp;&nbsp;</span>	if finished {
<span id="L714" class="ln">   714&nbsp;&nbsp;</span>		b = b[:sr.logicalRemaining()]
<span id="L715" class="ln">   715&nbsp;&nbsp;</span>	}
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	b0 := b
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>	endPos := sr.pos + int64(len(b))
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>	for endPos &gt; sr.pos &amp;&amp; err == nil {
<span id="L720" class="ln">   720&nbsp;&nbsp;</span>		var nf int <span class="comment">// Bytes read in fragment</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span>		holeStart, holeEnd := sr.sp[0].Offset, sr.sp[0].endOffset()
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>		if sr.pos &lt; holeStart { <span class="comment">// In a data fragment</span>
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>			bf := b[:min(int64(len(b)), holeStart-sr.pos)]
<span id="L724" class="ln">   724&nbsp;&nbsp;</span>			nf, err = tryReadFull(sr.fr, bf)
<span id="L725" class="ln">   725&nbsp;&nbsp;</span>		} else { <span class="comment">// In a hole fragment</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span>			bf := b[:min(int64(len(b)), holeEnd-sr.pos)]
<span id="L727" class="ln">   727&nbsp;&nbsp;</span>			nf, err = tryReadFull(zeroReader{}, bf)
<span id="L728" class="ln">   728&nbsp;&nbsp;</span>		}
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>		b = b[nf:]
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>		sr.pos += int64(nf)
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>		if sr.pos &gt;= holeEnd &amp;&amp; len(sr.sp) &gt; 1 {
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>			sr.sp = sr.sp[1:] <span class="comment">// Ensure last fragment always remains</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>		}
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	}
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>	n = len(b0) - len(b)
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	switch {
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	case err == io.EOF:
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>		return n, errMissData <span class="comment">// Less data in dense file than sparse file</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	case err != nil:
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>		return n, err
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>	case sr.logicalRemaining() == 0 &amp;&amp; sr.physicalRemaining() &gt; 0:
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>		return n, errUnrefData <span class="comment">// More data in dense file than sparse file</span>
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>	case finished:
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>		return n, io.EOF
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	default:
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>		return n, nil
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>	}
<span id="L749" class="ln">   749&nbsp;&nbsp;</span>}
<span id="L750" class="ln">   750&nbsp;&nbsp;</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span>func (sr *sparseFileReader) WriteTo(w io.Writer) (n int64, err error) {
<span id="L752" class="ln">   752&nbsp;&nbsp;</span>	ws, ok := w.(io.WriteSeeker)
<span id="L753" class="ln">   753&nbsp;&nbsp;</span>	if ok {
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>		if _, err := ws.Seek(0, io.SeekCurrent); err != nil {
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>			ok = false <span class="comment">// Not all io.Seeker can really seek</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>		}
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	}
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	if !ok {
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>		return io.Copy(w, struct{ io.Reader }{sr})
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	}
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>	var writeLastByte bool
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	pos0 := sr.pos
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>	for sr.logicalRemaining() &gt; 0 &amp;&amp; !writeLastByte &amp;&amp; err == nil {
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>		var nf int64 <span class="comment">// Size of fragment</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>		holeStart, holeEnd := sr.sp[0].Offset, sr.sp[0].endOffset()
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>		if sr.pos &lt; holeStart { <span class="comment">// In a data fragment</span>
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>			nf = holeStart - sr.pos
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>			nf, err = io.CopyN(ws, sr.fr, nf)
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>		} else { <span class="comment">// In a hole fragment</span>
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>			nf = holeEnd - sr.pos
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>			if sr.physicalRemaining() == 0 {
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>				writeLastByte = true
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>				nf--
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>			}
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>			_, err = ws.Seek(nf, io.SeekCurrent)
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>		}
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>		sr.pos += nf
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>		if sr.pos &gt;= holeEnd &amp;&amp; len(sr.sp) &gt; 1 {
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>			sr.sp = sr.sp[1:] <span class="comment">// Ensure last fragment always remains</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>		}
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	}
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	<span class="comment">// If the last fragment is a hole, then seek to 1-byte before EOF, and</span>
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>	<span class="comment">// write a single byte to ensure the file is the right size.</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	if writeLastByte &amp;&amp; err == nil {
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>		_, err = ws.Write([]byte{0})
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		sr.pos++
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	}
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	n = sr.pos - pos0
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>	switch {
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	case err == io.EOF:
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>		return n, errMissData <span class="comment">// Less data in dense file than sparse file</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	case err != nil:
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>		return n, err
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	case sr.logicalRemaining() == 0 &amp;&amp; sr.physicalRemaining() &gt; 0:
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>		return n, errUnrefData <span class="comment">// More data in dense file than sparse file</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	default:
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>		return n, nil
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>	}
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>}
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>func (sr sparseFileReader) logicalRemaining() int64 {
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>	return sr.sp[len(sr.sp)-1].endOffset() - sr.pos
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>}
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>func (sr sparseFileReader) physicalRemaining() int64 {
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>	return sr.fr.physicalRemaining()
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>}
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>type zeroReader struct{}
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>func (zeroReader) Read(b []byte) (int, error) {
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>	for i := range b {
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		b[i] = 0
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>	}
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>	return len(b), nil
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>}
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>
<span id="L820" class="ln">   820&nbsp;&nbsp;</span><span class="comment">// mustReadFull is like io.ReadFull except it returns</span>
<span id="L821" class="ln">   821&nbsp;&nbsp;</span><span class="comment">// io.ErrUnexpectedEOF when io.EOF is hit before len(b) bytes are read.</span>
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>func mustReadFull(r io.Reader, b []byte) (int, error) {
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>	n, err := tryReadFull(r, b)
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>	if err == io.EOF {
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		err = io.ErrUnexpectedEOF
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>	}
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>	return n, err
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>}
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span><span class="comment">// tryReadFull is like io.ReadFull except it returns</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span><span class="comment">// io.EOF when it is hit before len(b) bytes are read.</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>func tryReadFull(r io.Reader, b []byte) (n int, err error) {
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>	for len(b) &gt; n &amp;&amp; err == nil {
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		var nn int
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		nn, err = r.Read(b[n:])
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		n += nn
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>	}
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>	if len(b) == n &amp;&amp; err == io.EOF {
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>		err = nil
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>	}
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>	return n, err
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>
<span id="L844" class="ln">   844&nbsp;&nbsp;</span><span class="comment">// readSpecialFile is like io.ReadAll except it returns</span>
<span id="L845" class="ln">   845&nbsp;&nbsp;</span><span class="comment">// ErrFieldTooLong if more than maxSpecialFileSize is read.</span>
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>func readSpecialFile(r io.Reader) ([]byte, error) {
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>	buf, err := io.ReadAll(io.LimitReader(r, maxSpecialFileSize+1))
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>	if len(buf) &gt; maxSpecialFileSize {
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		return nil, ErrFieldTooLong
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>	}
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>	return buf, err
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>
<span id="L854" class="ln">   854&nbsp;&nbsp;</span><span class="comment">// discard skips n bytes in r, reporting an error if unable to do so.</span>
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>func discard(r io.Reader, n int64) error {
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>	<span class="comment">// If possible, Seek to the last byte before the end of the data section.</span>
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>	<span class="comment">// Do this because Seek is often lazy about reporting errors; this will mask</span>
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>	<span class="comment">// the fact that the stream may be truncated. We can rely on the</span>
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	<span class="comment">// io.CopyN done shortly afterwards to trigger any IO errors.</span>
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>	var seekSkipped int64 <span class="comment">// Number of bytes skipped via Seek</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	if sr, ok := r.(io.Seeker); ok &amp;&amp; n &gt; 1 {
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>		<span class="comment">// Not all io.Seeker can actually Seek. For example, os.Stdin implements</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>		<span class="comment">// io.Seeker, but calling Seek always returns an error and performs</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>		<span class="comment">// no action. Thus, we try an innocent seek to the current position</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>		<span class="comment">// to see if Seek is really supported.</span>
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>		pos1, err := sr.Seek(0, io.SeekCurrent)
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>		if pos1 &gt;= 0 &amp;&amp; err == nil {
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>			<span class="comment">// Seek seems supported, so perform the real Seek.</span>
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>			pos2, err := sr.Seek(n-1, io.SeekCurrent)
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>			if pos2 &lt; 0 || err != nil {
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>				return err
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>			}
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>			seekSkipped = pos2 - pos1
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>		}
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>	}
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>	copySkipped, err := io.CopyN(io.Discard, r, n-seekSkipped)
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>	if err == io.EOF &amp;&amp; seekSkipped+copySkipped &lt; n {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>		err = io.ErrUnexpectedEOF
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>	}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>	return err
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>
</pre><p><a href="reader.go?m=text">View as plain text</a></p>

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
