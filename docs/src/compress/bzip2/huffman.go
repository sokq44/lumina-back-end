<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/compress/bzip2/huffman.go - Go Documentation Server</title>

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
<a href="huffman.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/compress">compress</a>/<a href="http://localhost:8080/src/compress/bzip2">bzip2</a>/<span class="text-muted">huffman.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/compress/bzip2">compress/bzip2</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package bzip2
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import &#34;sort&#34;
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">// A huffmanTree is a binary tree which is navigated, bit-by-bit to reach a</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">// symbol.</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>type huffmanTree struct {
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	<span class="comment">// nodes contains all the non-leaf nodes in the tree. nodes[0] is the</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	<span class="comment">// root of the tree and nextNode contains the index of the next element</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	<span class="comment">// of nodes to use when the tree is being constructed.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	nodes    []huffmanNode
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	nextNode int
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>}
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// A huffmanNode is a node in the tree. left and right contain indexes into the</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// nodes slice of the tree. If left or right is invalidNodeValue then the child</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// is a left node and its value is in leftValue/rightValue.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// The symbols are uint16s because bzip2 encodes not only MTF indexes in the</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// tree, but also two magic values for run-length encoding and an EOF symbol.</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// Thus there are more than 256 possible symbols.</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type huffmanNode struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	left, right           uint16
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	leftValue, rightValue uint16
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>}
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// invalidNodeValue is an invalid index which marks a leaf node in the tree.</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>const invalidNodeValue = 0xffff
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">// Decode reads bits from the given bitReader and navigates the tree until a</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">// symbol is found.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>func (t *huffmanTree) Decode(br *bitReader) (v uint16) {
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	nodeIndex := uint16(0) <span class="comment">// node 0 is the root of the tree.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	for {
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		node := &amp;t.nodes[nodeIndex]
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		var bit uint16
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		if br.bits &gt; 0 {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>			<span class="comment">// Get next bit - fast path.</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>			br.bits--
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>			bit = uint16(br.n&gt;&gt;(br.bits&amp;63)) &amp; 1
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		} else {
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>			<span class="comment">// Get next bit - slow path.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			<span class="comment">// Use ReadBits to retrieve a single bit</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>			<span class="comment">// from the underling io.ByteReader.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>			bit = uint16(br.ReadBits(1))
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		}
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		<span class="comment">// Trick a compiler into generating conditional move instead of branch,</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		<span class="comment">// by making both loads unconditional.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		l, r := node.left, node.right
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		if bit == 1 {
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			nodeIndex = l
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		} else {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>			nodeIndex = r
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if nodeIndex == invalidNodeValue {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			<span class="comment">// We found a leaf. Use the value of bit to decide</span>
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			<span class="comment">// whether is a left or a right value.</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			l, r := node.leftValue, node.rightValue
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			if bit == 1 {
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>				v = l
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			} else {
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>				v = r
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			}
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>			return
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	}
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span><span class="comment">// newHuffmanTree builds a Huffman tree from a slice containing the code</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span><span class="comment">// lengths of each symbol. The maximum code length is 32 bits.</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>func newHuffmanTree(lengths []uint8) (huffmanTree, error) {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	<span class="comment">// There are many possible trees that assign the same code length to</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	<span class="comment">// each symbol (consider reflecting a tree down the middle, for</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">// example). Since the code length assignments determine the</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">// efficiency of the tree, each of these trees is equally good. In</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">// order to minimize the amount of information needed to build a tree</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	<span class="comment">// bzip2 uses a canonical tree so that it can be reconstructed given</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	<span class="comment">// only the code length assignments.</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if len(lengths) &lt; 2 {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		panic(&#34;newHuffmanTree: too few symbols&#34;)
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	}
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	var t huffmanTree
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	<span class="comment">// First we sort the code length assignments by ascending code length,</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	<span class="comment">// using the symbol value to break ties.</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	pairs := make([]huffmanSymbolLengthPair, len(lengths))
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	for i, length := range lengths {
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		pairs[i].value = uint16(i)
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		pairs[i].length = length
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	sort.Slice(pairs, func(i, j int) bool {
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		if pairs[i].length &lt; pairs[j].length {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>			return true
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		}
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		if pairs[i].length &gt; pairs[j].length {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			return false
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		if pairs[i].value &lt; pairs[j].value {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			return true
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		}
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		return false
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	})
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	<span class="comment">// Now we assign codes to the symbols, starting with the longest code.</span>
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	<span class="comment">// We keep the codes packed into a uint32, at the most-significant end.</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	<span class="comment">// So branches are taken from the MSB downwards. This makes it easy to</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	<span class="comment">// sort them later.</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	code := uint32(0)
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	length := uint8(32)
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	codes := make([]huffmanCode, len(lengths))
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	for i := len(pairs) - 1; i &gt;= 0; i-- {
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		if length &gt; pairs[i].length {
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>			length = pairs[i].length
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		}
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		codes[i].code = code
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		codes[i].codeLen = length
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		codes[i].value = pairs[i].value
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		<span class="comment">// We need to &#39;increment&#39; the code, which means treating |code|</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		<span class="comment">// like a |length| bit number.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		code += 1 &lt;&lt; (32 - length)
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">// Now we can sort by the code so that the left half of each branch are</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">// grouped together, recursively.</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	sort.Slice(codes, func(i, j int) bool {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		return codes[i].code &lt; codes[j].code
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	})
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	t.nodes = make([]huffmanNode, len(codes))
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	_, err := buildHuffmanNode(&amp;t, codes, 0)
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	return t, err
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>}
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// huffmanSymbolLengthPair contains a symbol and its code length.</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>type huffmanSymbolLengthPair struct {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	value  uint16
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	length uint8
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span><span class="comment">// huffmanCode contains a symbol, its code and code length.</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>type huffmanCode struct {
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	code    uint32
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	codeLen uint8
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	value   uint16
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>}
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// buildHuffmanNode takes a slice of sorted huffmanCodes and builds a node in</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// the Huffman tree at the given level. It returns the index of the newly</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span><span class="comment">// constructed node.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>func buildHuffmanNode(t *huffmanTree, codes []huffmanCode, level uint32) (nodeIndex uint16, err error) {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	test := uint32(1) &lt;&lt; (31 - level)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">// We have to search the list of codes to find the divide between the left and right sides.</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	firstRightIndex := len(codes)
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	for i, code := range codes {
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		if code.code&amp;test != 0 {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>			firstRightIndex = i
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>			break
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	left := codes[:firstRightIndex]
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	right := codes[firstRightIndex:]
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	if len(left) == 0 || len(right) == 0 {
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		<span class="comment">// There is a superfluous level in the Huffman tree indicating</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		<span class="comment">// a bug in the encoder. However, this bug has been observed in</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		<span class="comment">// the wild so we handle it.</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>		<span class="comment">// If this function was called recursively then we know that</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>		<span class="comment">// len(codes) &gt;= 2 because, otherwise, we would have hit the</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>		<span class="comment">// &#34;leaf node&#34; case, below, and not recurred.</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		<span class="comment">// However, for the initial call it&#39;s possible that len(codes)</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		<span class="comment">// is zero or one. Both cases are invalid because a zero length</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		<span class="comment">// tree cannot encode anything and a length-1 tree can only</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		<span class="comment">// encode EOF and so is superfluous. We reject both.</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		if len(codes) &lt; 2 {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>			return 0, StructuralError(&#34;empty Huffman tree&#34;)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>		}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		<span class="comment">// In this case the recursion doesn&#39;t always reduce the length</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>		<span class="comment">// of codes so we need to ensure termination via another</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>		<span class="comment">// mechanism.</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		if level == 31 {
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>			<span class="comment">// Since len(codes) &gt;= 2 the only way that the values</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			<span class="comment">// can match at all 32 bits is if they are equal, which</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>			<span class="comment">// is invalid. This ensures that we never enter</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>			<span class="comment">// infinite recursion.</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			return 0, StructuralError(&#34;equal symbols in Huffman tree&#34;)
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>		}
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		if len(left) == 0 {
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>			return buildHuffmanNode(t, right, level+1)
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>		}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>		return buildHuffmanNode(t, left, level+1)
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	}
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	nodeIndex = uint16(t.nextNode)
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	node := &amp;t.nodes[t.nextNode]
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	t.nextNode++
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	if len(left) == 1 {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>		<span class="comment">// leaf node</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		node.left = invalidNodeValue
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		node.leftValue = left[0].value
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	} else {
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>		node.left, err = buildHuffmanNode(t, left, level+1)
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	}
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	if err != nil {
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		return
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	}
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if len(right) == 1 {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		<span class="comment">// leaf node</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		node.right = invalidNodeValue
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		node.rightValue = right[0].value
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	} else {
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		node.right, err = buildHuffmanNode(t, right, level+1)
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	return
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
</pre><p><a href="huffman.go?m=text">View as plain text</a></p>

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
