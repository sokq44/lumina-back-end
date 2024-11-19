<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/html/template/doc.go - Go Documentation Server</title>

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
<a href="doc.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/html">html</a>/<a href="http://localhost:8080/src/html/template">template</a>/<span class="text-muted">doc.go</span>
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
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package template (html/template) implements data-driven templates for
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>generating HTML output safe against code injection. It provides the
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>same interface as [text/template] and should be used instead of
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>[text/template] whenever the output is HTML.
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>The documentation here focuses on the security features of the package.
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>For information about how to program the templates themselves, see the
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>documentation for [text/template].
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span># Introduction
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>This package wraps [text/template] so you can share its template API
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>to parse and execute HTML templates safely.
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	tmpl, err := template.New(&#34;name&#34;).Parse(...)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	// Error checking elided
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	err = tmpl.Execute(out, data)
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>If successful, tmpl will now be injection-safe. Otherwise, err is an error
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>defined in the docs for ErrorCode.
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>HTML templates treat data values as plain text which should be encoded so they
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>can be safely embedded in an HTML document. The escaping is contextual, so
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>actions can appear within JavaScript, CSS, and URI contexts.
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>The security model used by this package assumes that template authors are
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>trusted, while Execute&#39;s data parameter is not. More details are
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>provided below.
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>Example
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	import &#34;text/template&#34;
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	...
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	t, err := template.New(&#34;foo&#34;).Parse(`{{define &#34;T&#34;}}Hello, {{.}}!{{end}}`)
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	err = t.ExecuteTemplate(out, &#34;T&#34;, &#34;&lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;&#34;)
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>produces
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>but the contextual autoescaping in html/template
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	import &#34;html/template&#34;
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	...
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	t, err := template.New(&#34;foo&#34;).Parse(`{{define &#34;T&#34;}}Hello, {{.}}!{{end}}`)
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	err = t.ExecuteTemplate(out, &#34;T&#34;, &#34;&lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;&#34;)
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>produces safe, escaped HTML output
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	Hello, &amp;lt;script&amp;gt;alert(&amp;#39;you have been pwned&amp;#39;)&amp;lt;/script&amp;gt;!
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span># Contexts
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>This package understands HTML, CSS, JavaScript, and URIs. It adds sanitizing
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>functions to each simple action pipeline, so given the excerpt
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	&lt;a href=&#34;/search?q={{.}}&#34;&gt;{{.}}&lt;/a&gt;
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>At parse time each {{.}} is overwritten to add escaping functions as necessary.
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>In this case it becomes
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	&lt;a href=&#34;/search?q={{. | urlescaper | attrescaper}}&#34;&gt;{{. | htmlescaper}}&lt;/a&gt;
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>where urlescaper, attrescaper, and htmlescaper are aliases for internal escaping
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>functions.
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>For these internal escaping functions, if an action pipeline evaluates to
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>a nil interface value, it is treated as though it were an empty string.
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span># Namespaced and data- attributes
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>Attributes with a namespace are treated as if they had no namespace.
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>Given the excerpt
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	&lt;a my:href=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>At parse time the attribute will be treated as if it were just &#34;href&#34;.
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>So at parse time the template becomes:
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	&lt;a my:href=&#34;{{. | urlescaper | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>Similarly to attributes with namespaces, attributes with a &#34;data-&#34; prefix are
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>treated as if they had no &#34;data-&#34; prefix. So given
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&lt;a data-href=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>At parse time this becomes
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	&lt;a data-href=&#34;{{. | urlescaper | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>If an attribute has both a namespace and a &#34;data-&#34; prefix, only the namespace
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>will be removed when determining the context. For example
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	&lt;a my:data-href=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>This is handled as if &#34;my:data-href&#34; was just &#34;data-href&#34; and not &#34;href&#34; as
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>it would be if the &#34;data-&#34; prefix were to be ignored too. Thus at parse
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>time this becomes just
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	&lt;a my:data-href=&#34;{{. | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>As a special case, attributes with the namespace &#34;xmlns&#34; are always treated
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>as containing URLs. Given the excerpts
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	&lt;a xmlns:title=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	&lt;a xmlns:href=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	&lt;a xmlns:onclick=&#34;{{.}}&#34;&gt;&lt;/a&gt;
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>At parse time they become:
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	&lt;a xmlns:title=&#34;{{. | urlescaper | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	&lt;a xmlns:href=&#34;{{. | urlescaper | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	&lt;a xmlns:onclick=&#34;{{. | urlescaper | attrescaper}}&#34;&gt;&lt;/a&gt;
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span># Errors
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>See the documentation of ErrorCode for details.
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span># A fuller picture
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>The rest of this package comment may be skipped on first reading; it includes
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>details necessary to understand escaping contexts and error messages. Most users
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>will not need to understand these details.
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span># Contexts
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>Assuming {{.}} is `O&#39;Reilly: How are &lt;i&gt;you&lt;/i&gt;?`, the table below shows
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>how {{.}} appears when used in the context to the left.
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	Context                          {{.}} After
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	{{.}}                            O&#39;Reilly: How are &amp;lt;i&amp;gt;you&amp;lt;/i&amp;gt;?
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	&lt;a title=&#39;{{.}}&#39;&gt;                O&amp;#39;Reilly: How are you?
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	&lt;a href=&#34;/{{.}}&#34;&gt;                O&amp;#39;Reilly: How are %3ci%3eyou%3c/i%3e?
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	&lt;a href=&#34;?q={{.}}&#34;&gt;              O&amp;#39;Reilly%3a%20How%20are%3ci%3e...%3f
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	&lt;a onx=&#39;f(&#34;{{.}}&#34;)&#39;&gt;             O\x27Reilly: How are \x3ci\x3eyou...?
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	&lt;a onx=&#39;f({{.}})&#39;&gt;               &#34;O\x27Reilly: How are \x3ci\x3eyou...?&#34;
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	&lt;a onx=&#39;pattern = /{{.}}/;&#39;&gt;     O\x27Reilly: How are \x3ci\x3eyou...\x3f
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>If used in an unsafe context, then the value might be filtered out:
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	Context                          {{.}} After
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	&lt;a href=&#34;{{.}}&#34;&gt;                 #ZgotmplZ
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>since &#34;O&#39;Reilly:&#34; is not an allowed protocol like &#34;http:&#34;.
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>If {{.}} is the innocuous word, `left`, then it can appear more widely,
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	Context                              {{.}} After
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	{{.}}                                left
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	&lt;a title=&#39;{{.}}&#39;&gt;                    left
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	&lt;a href=&#39;{{.}}&#39;&gt;                     left
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	&lt;a href=&#39;/{{.}}&#39;&gt;                    left
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	&lt;a href=&#39;?dir={{.}}&#39;&gt;                left
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	&lt;a style=&#34;border-{{.}}: 4px&#34;&gt;        left
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	&lt;a style=&#34;align: {{.}}&#34;&gt;             left
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	&lt;a style=&#34;background: &#39;{{.}}&#39;&gt;       left
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	&lt;a style=&#34;background: url(&#39;{{.}}&#39;)&gt;  left
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	&lt;style&gt;p.{{.}} {color:red}&lt;/style&gt;   left
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>Non-string values can be used in JavaScript contexts.
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>If {{.}} is
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	struct{A,B string}{ &#34;foo&#34;, &#34;bar&#34; }
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>in the escaped template
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	&lt;script&gt;var pair = {{.}};&lt;/script&gt;
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>then the template output is
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	&lt;script&gt;var pair = {&#34;A&#34;: &#34;foo&#34;, &#34;B&#34;: &#34;bar&#34;};&lt;/script&gt;
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>See package json to understand how non-string content is marshaled for
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>embedding in JavaScript contexts.
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span># Typed Strings
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>By default, this package assumes that all pipelines produce a plain text string.
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>It adds escaping pipeline stages necessary to correctly and safely embed that
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>plain text string in the appropriate context.
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>When a data value is not plain text, you can make sure it is not over-escaped
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>by marking it with its type.
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>Types HTML, JS, URL, and others from content.go can carry safe content that is
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>exempted from escaping.
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>The template
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	Hello, {{.}}!
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>can be invoked with
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	tmpl.Execute(out, template.HTML(`&lt;b&gt;World&lt;/b&gt;`))
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>to produce
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	Hello, &lt;b&gt;World&lt;/b&gt;!
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>instead of the
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	Hello, &amp;lt;b&amp;gt;World&amp;lt;b&amp;gt;!
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>that would have been produced if {{.}} was a regular string.
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span># Security Model
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>https://rawgit.com/mikesamuel/sanitized-jquery-templates/trunk/safetemplate.html#problem_definition defines &#34;safe&#34; as used by this package.
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>This package assumes that template authors are trusted, that Execute&#39;s data
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>parameter is not, and seeks to preserve the properties below in the face
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>of untrusted data:
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>Structure Preservation Property:
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>&#34;... when a template author writes an HTML tag in a safe templating language,
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>the browser will interpret the corresponding portion of the output as a tag
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>regardless of the values of untrusted data, and similarly for other structures
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>such as attribute boundaries and JS and CSS string boundaries.&#34;
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>Code Effect Property:
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>&#34;... only code specified by the template author should run as a result of
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>injecting the template output into a page and all code specified by the
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>template author should run as a result of the same.&#34;
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>Least Surprise Property:
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>&#34;A developer (or code reviewer) familiar with HTML, CSS, and JavaScript, who
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>knows that contextual autoescaping happens should be able to look at a {{.}}
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>and correctly infer what sanitization happens.&#34;
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>As a consequence of the Least Surprise Property, template actions within an
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>ECMAScript 6 template literal are disabled by default.
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>Handling string interpolation within these literals is rather complex resulting
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>in no clear safe way to support it.
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>To re-enable template actions within ECMAScript 6 template literals, use the
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>GODEBUG=jstmpllitinterp=1 environment variable.
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>*/</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>package template
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
</pre><p><a href="doc.go?m=text">View as plain text</a></p>

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
