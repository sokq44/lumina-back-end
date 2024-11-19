<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/attr.go - Go Documentation Server</title>

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
<a href="attr.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">attr.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/html/template">html/template</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package template
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;strings&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>)
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">// attrTypeMap[n] describes the value of the given attribute.</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span><span class="comment">// If an attribute affects (or can mask) the encoding or interpretation of</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// other content, or affects the contents, idempotency, or credentials of a</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// network message, then the value in this map is contentTypeUnsafe.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// This map is derived from HTML5, specifically</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">// https://www.w3.org/TR/html5/Overview.html#attributes-1</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// as well as &#34;%URI&#34;-typed attributes from</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">// https://www.w3.org/TR/html4/index/attributes.html</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>var attrTypeMap = map[string]contentType{
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;accept&#34;:          contentTypePlain,
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;accept-charset&#34;:  contentTypeUnsafe,
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;action&#34;:          contentTypeURL,
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;alt&#34;:             contentTypePlain,
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;archive&#34;:         contentTypeURL,
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;async&#34;:           contentTypeUnsafe,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;autocomplete&#34;:    contentTypePlain,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;autofocus&#34;:       contentTypePlain,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;autoplay&#34;:        contentTypePlain,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;background&#34;:      contentTypeURL,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	&#34;border&#34;:          contentTypePlain,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	&#34;checked&#34;:         contentTypePlain,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	&#34;cite&#34;:            contentTypeURL,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	&#34;challenge&#34;:       contentTypeUnsafe,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	&#34;charset&#34;:         contentTypeUnsafe,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	&#34;class&#34;:           contentTypePlain,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	&#34;classid&#34;:         contentTypeURL,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	&#34;codebase&#34;:        contentTypeURL,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	&#34;cols&#34;:            contentTypePlain,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	&#34;colspan&#34;:         contentTypePlain,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	&#34;content&#34;:         contentTypeUnsafe,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	&#34;contenteditable&#34;: contentTypePlain,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	&#34;contextmenu&#34;:     contentTypePlain,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	&#34;controls&#34;:        contentTypePlain,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	&#34;coords&#34;:          contentTypePlain,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	&#34;crossorigin&#34;:     contentTypeUnsafe,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	&#34;data&#34;:            contentTypeURL,
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	&#34;datetime&#34;:        contentTypePlain,
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	&#34;default&#34;:         contentTypePlain,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	&#34;defer&#34;:           contentTypeUnsafe,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	&#34;dir&#34;:             contentTypePlain,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	&#34;dirname&#34;:         contentTypePlain,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	&#34;disabled&#34;:        contentTypePlain,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	&#34;draggable&#34;:       contentTypePlain,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	&#34;dropzone&#34;:        contentTypePlain,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	&#34;enctype&#34;:         contentTypeUnsafe,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	&#34;for&#34;:             contentTypePlain,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	&#34;form&#34;:            contentTypeUnsafe,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	&#34;formaction&#34;:      contentTypeURL,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	&#34;formenctype&#34;:     contentTypeUnsafe,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	&#34;formmethod&#34;:      contentTypeUnsafe,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	&#34;formnovalidate&#34;:  contentTypeUnsafe,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	&#34;formtarget&#34;:      contentTypePlain,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	&#34;headers&#34;:         contentTypePlain,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	&#34;height&#34;:          contentTypePlain,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	&#34;hidden&#34;:          contentTypePlain,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	&#34;high&#34;:            contentTypePlain,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	&#34;href&#34;:            contentTypeURL,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	&#34;hreflang&#34;:        contentTypePlain,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	&#34;http-equiv&#34;:      contentTypeUnsafe,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	&#34;icon&#34;:            contentTypeURL,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	&#34;id&#34;:              contentTypePlain,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	&#34;ismap&#34;:           contentTypePlain,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	&#34;keytype&#34;:         contentTypeUnsafe,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	&#34;kind&#34;:            contentTypePlain,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	&#34;label&#34;:           contentTypePlain,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	&#34;lang&#34;:            contentTypePlain,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	&#34;language&#34;:        contentTypeUnsafe,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	&#34;list&#34;:            contentTypePlain,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	&#34;longdesc&#34;:        contentTypeURL,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	&#34;loop&#34;:            contentTypePlain,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	&#34;low&#34;:             contentTypePlain,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	&#34;manifest&#34;:        contentTypeURL,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	&#34;max&#34;:             contentTypePlain,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	&#34;maxlength&#34;:       contentTypePlain,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	&#34;media&#34;:           contentTypePlain,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	&#34;mediagroup&#34;:      contentTypePlain,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	&#34;method&#34;:          contentTypeUnsafe,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	&#34;min&#34;:             contentTypePlain,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	&#34;multiple&#34;:        contentTypePlain,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&#34;name&#34;:            contentTypePlain,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	&#34;novalidate&#34;:      contentTypeUnsafe,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	<span class="comment">// Skip handler names from</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/TR/html5/webappapis.html#event-handlers-on-elements,-document-objects,-and-window-objects</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	<span class="comment">// since we have special handling in attrType.</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	&#34;open&#34;:        contentTypePlain,
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	&#34;optimum&#34;:     contentTypePlain,
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	&#34;pattern&#34;:     contentTypeUnsafe,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	&#34;placeholder&#34;: contentTypePlain,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	&#34;poster&#34;:      contentTypeURL,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	&#34;profile&#34;:     contentTypeURL,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	&#34;preload&#34;:     contentTypePlain,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	&#34;pubdate&#34;:     contentTypePlain,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	&#34;radiogroup&#34;:  contentTypePlain,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	&#34;readonly&#34;:    contentTypePlain,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	&#34;rel&#34;:         contentTypeUnsafe,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	&#34;required&#34;:    contentTypePlain,
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	&#34;reversed&#34;:    contentTypePlain,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	&#34;rows&#34;:        contentTypePlain,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	&#34;rowspan&#34;:     contentTypePlain,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	&#34;sandbox&#34;:     contentTypeUnsafe,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	&#34;spellcheck&#34;:  contentTypePlain,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	&#34;scope&#34;:       contentTypePlain,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	&#34;scoped&#34;:      contentTypePlain,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	&#34;seamless&#34;:    contentTypePlain,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	&#34;selected&#34;:    contentTypePlain,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	&#34;shape&#34;:       contentTypePlain,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	&#34;size&#34;:        contentTypePlain,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	&#34;sizes&#34;:       contentTypePlain,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	&#34;span&#34;:        contentTypePlain,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	&#34;src&#34;:         contentTypeURL,
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	&#34;srcdoc&#34;:      contentTypeHTML,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	&#34;srclang&#34;:     contentTypePlain,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	&#34;srcset&#34;:      contentTypeSrcset,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	&#34;start&#34;:       contentTypePlain,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	&#34;step&#34;:        contentTypePlain,
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	&#34;style&#34;:       contentTypeCSS,
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	&#34;tabindex&#34;:    contentTypePlain,
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	&#34;target&#34;:      contentTypePlain,
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	&#34;title&#34;:       contentTypePlain,
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	&#34;type&#34;:        contentTypeUnsafe,
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	&#34;usemap&#34;:      contentTypeURL,
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	&#34;value&#34;:       contentTypeUnsafe,
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	&#34;width&#34;:       contentTypePlain,
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	&#34;wrap&#34;:        contentTypePlain,
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	&#34;xmlns&#34;:       contentTypeURL,
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>}
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">// attrType returns a conservative (upper-bound on authority) guess at the</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">// type of the lowercase named attribute.</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>func attrType(name string) contentType {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	if strings.HasPrefix(name, &#34;data-&#34;) {
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		<span class="comment">// Strip data- so that custom attribute heuristics below are</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		<span class="comment">// widely applied.</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		<span class="comment">// Treat data-action as URL below.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		name = name[5:]
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	} else if prefix, short, ok := strings.Cut(name, &#34;:&#34;); ok {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		if prefix == &#34;xmlns&#34; {
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			return contentTypeURL
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		<span class="comment">// Treat svg:href and xlink:href as href below.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		name = short
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	if t, ok := attrTypeMap[name]; ok {
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		return t
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	}
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">// Treat partial event handler names as script.</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if strings.HasPrefix(name, &#34;on&#34;) {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return contentTypeJS
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	}
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">// Heuristics to prevent &#34;javascript:...&#34; injection in custom</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">// data attributes and custom attributes like g:tweetUrl.</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">// https://www.w3.org/TR/html5/dom.html#embedding-custom-non-visible-data-with-the-data-*-attributes</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">// &#34;Custom data attributes are intended to store custom data</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">//  private to the page or application, for which there are no</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">//  more appropriate attributes or elements.&#34;</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">// Developers seem to store URL content in data URLs that start</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">// or end with &#34;URI&#34; or &#34;URL&#34;.</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	if strings.Contains(name, &#34;src&#34;) ||
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		strings.Contains(name, &#34;uri&#34;) ||
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		strings.Contains(name, &#34;url&#34;) {
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		return contentTypeURL
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	return contentTypePlain
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>}
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>
</pre><p><a href="attr.go?m=text">View as plain text</a></p>

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
