<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/internal/trace/mud.go - Go Documentation Server</title>

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
<a href="mud.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/internal">internal</a>/<a href="http://localhost:8080/src/internal/trace">trace</a>/<span class="text-muted">mud.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/internal/trace">internal/trace</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2017 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package trace
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;math&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;sort&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>)
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// mud is an updatable mutator utilization distribution.</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// This is a continuous distribution of duration over mutator</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// utilization. For example, the integral from mutator utilization a</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// to b is the total duration during which the mutator utilization was</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// in the range [a, b].</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// This distribution is *not* normalized (it is not a probability</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">// distribution). This makes it easier to work with as it&#39;s being</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// updated.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// It is represented as the sum of scaled uniform distribution</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">// functions and Dirac delta functions (which are treated as</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">// degenerate uniform distributions).</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>type mud struct {
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	sorted, unsorted []edge
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	<span class="comment">// trackMass is the inverse cumulative sum to track as the</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	<span class="comment">// distribution is updated.</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	trackMass float64
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	<span class="comment">// trackBucket is the bucket in which trackMass falls. If the</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	<span class="comment">// total mass of the distribution is &lt; trackMass, this is</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	<span class="comment">// len(hist).</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	trackBucket int
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	<span class="comment">// trackSum is the cumulative sum of hist[:trackBucket]. Once</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	<span class="comment">// trackSum &gt;= trackMass, trackBucket must be recomputed.</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	trackSum float64
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	<span class="comment">// hist is a hierarchical histogram of distribution mass.</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	hist [mudDegree]float64
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>}
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>const (
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	<span class="comment">// mudDegree is the number of buckets in the MUD summary</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	<span class="comment">// histogram.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	mudDegree = 1024
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>)
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>type edge struct {
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	<span class="comment">// At x, the function increases by y.</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	x, delta float64
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	<span class="comment">// Additionally at x is a Dirac delta function with area dirac.</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	dirac float64
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>}
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span><span class="comment">// add adds a uniform function over [l, r] scaled so the total weight</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span><span class="comment">// of the uniform is area. If l==r, this adds a Dirac delta function.</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>func (d *mud) add(l, r, area float64) {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	if area == 0 {
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		return
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	}
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	if r &lt; l {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		l, r = r, l
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	<span class="comment">// Add the edges.</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	if l == r {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		d.unsorted = append(d.unsorted, edge{l, 0, area})
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	} else {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		delta := area / (r - l)
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		d.unsorted = append(d.unsorted, edge{l, delta, 0}, edge{r, -delta, 0})
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	<span class="comment">// Update the histogram.</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	h := &amp;d.hist
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	lbFloat, lf := math.Modf(l * mudDegree)
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	lb := int(lbFloat)
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	if lb &gt;= mudDegree {
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		lb, lf = mudDegree-1, 1
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	if l == r {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		h[lb] += area
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	} else {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		rbFloat, rf := math.Modf(r * mudDegree)
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		rb := int(rbFloat)
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		if rb &gt;= mudDegree {
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>			rb, rf = mudDegree-1, 1
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		if lb == rb {
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			h[lb] += area
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		} else {
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			perBucket := area / (r - l) / mudDegree
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>			h[lb] += perBucket * (1 - lf)
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			h[rb] += perBucket * rf
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			for i := lb + 1; i &lt; rb; i++ {
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>				h[i] += perBucket
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>			}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	}
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// Update mass tracking.</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	if thresh := float64(d.trackBucket) / mudDegree; l &lt; thresh {
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		if r &lt; thresh {
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>			d.trackSum += area
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		} else {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			d.trackSum += area * (thresh - l) / (r - l)
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		if d.trackSum &gt;= d.trackMass {
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>			<span class="comment">// The tracked mass now falls in a different</span>
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>			<span class="comment">// bucket. Recompute the inverse cumulative sum.</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			d.setTrackMass(d.trackMass)
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	}
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span><span class="comment">// setTrackMass sets the mass to track the inverse cumulative sum for.</span>
<span id="L119" class="ln">   119&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span><span class="comment">// Specifically, mass is a cumulative duration, and the mutator</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span><span class="comment">// utilization bounds for this duration can be queried using</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span><span class="comment">// approxInvCumulativeSum.</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>func (d *mud) setTrackMass(mass float64) {
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	d.trackMass = mass
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">// Find the bucket currently containing trackMass by computing</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">// the cumulative sum.</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	sum := 0.0
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	for i, val := range d.hist[:] {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		newSum := sum + val
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		if newSum &gt; mass {
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			<span class="comment">// mass falls in bucket i.</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			d.trackBucket = i
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			d.trackSum = sum
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>			return
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		sum = newSum
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	}
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	d.trackBucket = len(d.hist)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	d.trackSum = sum
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>}
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">// approxInvCumulativeSum is like invCumulativeSum, but specifically</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">// operates on the tracked mass and returns an upper and lower bound</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">// approximation of the inverse cumulative sum.</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">// The true inverse cumulative sum will be in the range [lower, upper).</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>func (d *mud) approxInvCumulativeSum() (float64, float64, bool) {
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	if d.trackBucket == len(d.hist) {
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		return math.NaN(), math.NaN(), false
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	}
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	return float64(d.trackBucket) / mudDegree, float64(d.trackBucket+1) / mudDegree, true
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span><span class="comment">// invCumulativeSum returns x such that the integral of d from -∞ to x</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span><span class="comment">// is y. If the total weight of d is less than y, it returns the</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span><span class="comment">// maximum of the distribution and false.</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span><span class="comment">// Specifically, y is a cumulative duration, and invCumulativeSum</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span><span class="comment">// returns the mutator utilization x such that at least y time has</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span><span class="comment">// been spent with mutator utilization &lt;= x.</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>func (d *mud) invCumulativeSum(y float64) (float64, bool) {
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	if len(d.sorted) == 0 &amp;&amp; len(d.unsorted) == 0 {
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		return math.NaN(), false
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	}
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// Sort edges.</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	edges := d.unsorted
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	sort.Slice(edges, func(i, j int) bool {
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		return edges[i].x &lt; edges[j].x
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	})
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">// Merge with sorted edges.</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	d.unsorted = nil
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if d.sorted == nil {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		d.sorted = edges
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	} else {
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		oldSorted := d.sorted
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>		newSorted := make([]edge, len(oldSorted)+len(edges))
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>		i, j := 0, 0
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		for o := range newSorted {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			if i &gt;= len(oldSorted) {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>				copy(newSorted[o:], edges[j:])
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>				break
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			} else if j &gt;= len(edges) {
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				copy(newSorted[o:], oldSorted[i:])
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				break
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>			} else if oldSorted[i].x &lt; edges[j].x {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				newSorted[o] = oldSorted[i]
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				i++
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>			} else {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>				newSorted[o] = edges[j]
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>				j++
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>			}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>		}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>		d.sorted = newSorted
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	<span class="comment">// Traverse edges in order computing a cumulative sum.</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	csum, rate, prevX := 0.0, 0.0, 0.0
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	for _, e := range d.sorted {
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>		newCsum := csum + (e.x-prevX)*rate
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>		if newCsum &gt;= y {
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>			<span class="comment">// y was exceeded between the previous edge</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>			<span class="comment">// and this one.</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>			if rate == 0 {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>				<span class="comment">// Anywhere between prevX and</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>				<span class="comment">// e.x will do. We return e.x</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>				<span class="comment">// because that takes care of</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				<span class="comment">// the y==0 case naturally.</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				return e.x, true
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>			}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>			return (y-csum)/rate + prevX, true
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>		newCsum += e.dirac
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>		if newCsum &gt;= y {
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>			<span class="comment">// y was exceeded by the Dirac delta at e.x.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>			return e.x, true
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		}
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		csum, prevX = newCsum, e.x
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>		rate += e.delta
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	}
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	return prevX, false
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>}
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
</pre><p><a href="mud.go?m=text">View as plain text</a></p>

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
