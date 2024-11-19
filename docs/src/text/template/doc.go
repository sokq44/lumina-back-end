<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/text/template/doc.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/text">text</a>/<a href="http://localhost:8080/src/text/template">template</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/text/template">text/template</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2011 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package template implements data-driven templates for generating textual output.
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>To generate HTML output, see [html/template], which has the same interface
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>as this package but automatically secures HTML output against certain attacks.
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>Templates are executed by applying them to a data structure. Annotations in the
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>template refer to elements of the data structure (typically a field of a struct
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>or a key in a map) to control execution and derive values to be displayed.
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>Execution of the template walks the structure and sets the cursor, represented
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>by a period &#39;.&#39; and called &#34;dot&#34;, to the value at the current location in the
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>structure as execution proceeds.
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>The input text for a template is UTF-8-encoded text in any format.
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>&#34;Actions&#34;--data evaluations or control structures--are delimited by
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>&#34;{{&#34; and &#34;}}&#34;; all text outside actions is copied to the output unchanged.
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>Once parsed, a template may be executed safely in parallel, although if parallel
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>executions share a Writer the output may be interleaved.
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>Here is a trivial example that prints &#34;17 items are made of wool&#34;.
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	type Inventory struct {
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		Material string
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		Count    uint
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	}
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	sweaters := Inventory{&#34;wool&#34;, 17}
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	tmpl, err := template.New(&#34;test&#34;).Parse(&#34;{{.Count}} items are made of {{.Material}}&#34;)
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	if err != nil { panic(err) }
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	err = tmpl.Execute(os.Stdout, sweaters)
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	if err != nil { panic(err) }
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>More intricate examples appear below.
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>Text and spaces
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>By default, all text between actions is copied verbatim when the template is
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>executed. For example, the string &#34; items are made of &#34; in the example above
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>appears on standard output when the program is run.
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>However, to aid in formatting template source code, if an action&#39;s left
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>delimiter (by default &#34;{{&#34;) is followed immediately by a minus sign and white
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>space, all trailing white space is trimmed from the immediately preceding text.
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>Similarly, if the right delimiter (&#34;}}&#34;) is preceded by white space and a minus
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>sign, all leading white space is trimmed from the immediately following text.
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>In these trim markers, the white space must be present:
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>&#34;{{- 3}}&#34; is like &#34;{{3}}&#34; but trims the immediately preceding text, while
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>&#34;{{-3}}&#34; parses as an action containing the number -3.
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>For instance, when executing the template whose source is
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	&#34;{{23 -}} &lt; {{- 45}}&#34;
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>the generated output would be
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	&#34;23&lt;45&#34;
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>For this trimming, the definition of white space characters is the same as in Go:
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>space, horizontal tab, carriage return, and newline.
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>Actions
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>Here is the list of actions. &#34;Arguments&#34; and &#34;pipelines&#34; are evaluations of
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>data, defined in detail in the corresponding sections that follow.
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>*/</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">//	{{/* a comment */}}</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span><span class="comment">//	{{- /* a comment with white space trimmed from preceding and following text */ -}}</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span><span class="comment">//		A comment; discarded. May contain newlines.</span>
<span id="L74" class="ln">    74&nbsp;&nbsp;</span><span class="comment">//		Comments do not nest and must start and end at the</span>
<span id="L75" class="ln">    75&nbsp;&nbsp;</span><span class="comment">//		delimiters, as shown here.</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span><span class="comment">/*
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	{{pipeline}}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		The default textual representation (the same as would be
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		printed by fmt.Print) of the value of the pipeline is copied
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		to the output.
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	{{if pipeline}} T1 {{end}}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		If the value of the pipeline is empty, no output is generated;
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		otherwise, T1 is executed. The empty values are false, 0, any
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		nil pointer or interface value, and any array, slice, map, or
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		string of length zero.
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		Dot is unaffected.
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	{{if pipeline}} T1 {{else}} T0 {{end}}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		If the value of the pipeline is empty, T0 is executed;
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		otherwise, T1 is executed. Dot is unaffected.
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		To simplify the appearance of if-else chains, the else action
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		of an if may include another if directly; the effect is exactly
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		the same as writing
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>			{{if pipeline}} T1 {{else}}{{if pipeline}} T0 {{end}}{{end}}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	{{range pipeline}} T1 {{end}}
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		The value of the pipeline must be an array, slice, map, or channel.
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		If the value of the pipeline has length zero, nothing is output;
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		otherwise, dot is set to the successive elements of the array,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		slice, or map and T1 is executed. If the value is a map and the
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		keys are of basic type with a defined order, the elements will be
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		visited in sorted key order.
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	{{range pipeline}} T1 {{else}} T0 {{end}}
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		The value of the pipeline must be an array, slice, map, or channel.
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		If the value of the pipeline has length zero, dot is unaffected and
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		T0 is executed; otherwise, dot is set to the successive elements
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		of the array, slice, or map and T1 is executed.
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	{{break}}
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		The innermost {{range pipeline}} loop is ended early, stopping the
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		current iteration and bypassing all remaining iterations.
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	{{continue}}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		The current iteration of the innermost {{range pipeline}} loop is
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		stopped, and the loop starts the next iteration.
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	{{template &#34;name&#34;}}
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		The template with the specified name is executed with nil data.
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	{{template &#34;name&#34; pipeline}}
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		The template with the specified name is executed with dot set
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		to the value of the pipeline.
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	{{block &#34;name&#34; pipeline}} T1 {{end}}
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		A block is shorthand for defining a template
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>			{{define &#34;name&#34;}} T1 {{end}}
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		and then executing it in place
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>			{{template &#34;name&#34; pipeline}}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		The typical use is to define a set of root templates that are
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		then customized by redefining the block templates within.
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	{{with pipeline}} T1 {{end}}
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		If the value of the pipeline is empty, no output is generated;
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		otherwise, dot is set to the value of the pipeline and T1 is
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		executed.
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	{{with pipeline}} T1 {{else}} T0 {{end}}
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		If the value of the pipeline is empty, dot is unaffected and T0
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		is executed; otherwise, dot is set to the value of the pipeline
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		and T1 is executed.
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>Arguments
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>An argument is a simple value, denoted by one of the following.
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	- A boolean, string, character, integer, floating-point, imaginary
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	  or complex constant in Go syntax. These behave like Go&#39;s untyped
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	  constants. Note that, as in Go, whether a large integer constant
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	  overflows when assigned or passed to a function can depend on whether
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	  the host machine&#39;s ints are 32 or 64 bits.
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	- The keyword nil, representing an untyped Go nil.
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	- The character &#39;.&#39; (period):
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		.
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	  The result is the value of dot.
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	- A variable name, which is a (possibly empty) alphanumeric string
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	  preceded by a dollar sign, such as
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		$piOver2
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	  or
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		$
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	  The result is the value of the variable.
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	  Variables are described below.
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	- The name of a field of the data, which must be a struct, preceded
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	  by a period, such as
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		.Field
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	  The result is the value of the field. Field invocations may be
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	  chained:
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	    .Field1.Field2
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	  Fields can also be evaluated on variables, including chaining:
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	    $x.Field1.Field2
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	- The name of a key of the data, which must be a map, preceded
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	  by a period, such as
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>		.Key
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	  The result is the map element value indexed by the key.
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	  Key invocations may be chained and combined with fields to any
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	  depth:
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	    .Field1.Key1.Field2.Key2
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	  Although the key must be an alphanumeric identifier, unlike with
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	  field names they do not need to start with an upper case letter.
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	  Keys can also be evaluated on variables, including chaining:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	    $x.key1.key2
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	- The name of a niladic method of the data, preceded by a period,
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	  such as
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>		.Method
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	  The result is the value of invoking the method with dot as the
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	  receiver, dot.Method(). Such a method must have one return value (of
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	  any type) or two return values, the second of which is an error.
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	  If it has two and the returned error is non-nil, execution terminates
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	  and an error is returned to the caller as the value of Execute.
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	  Method invocations may be chained and combined with fields and keys
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	  to any depth:
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	    .Field1.Key1.Method1.Field2.Key2.Method2
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	  Methods can also be evaluated on variables, including chaining:
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	    $x.Method1.Field
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	- The name of a niladic function, such as
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>		fun
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	  The result is the value of invoking the function, fun(). The return
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	  types and values behave as in methods. Functions and function
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	  names are described below.
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	- A parenthesized instance of one the above, for grouping. The result
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	  may be accessed by a field or map key invocation.
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>		print (.F1 arg1) (.F2 arg2)
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>		(.StructValuedMethod &#34;arg&#34;).Field
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>Arguments may evaluate to any type; if they are pointers the implementation
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>automatically indirects to the base type when required.
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>If an evaluation yields a function value, such as a function-valued
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>field of a struct, the function is not invoked automatically, but it
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>can be used as a truth value for an if action and the like. To invoke
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>it, use the call function, defined below.
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>Pipelines
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>A pipeline is a possibly chained sequence of &#34;commands&#34;. A command is a simple
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>value (argument) or a function or method call, possibly with multiple arguments:
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	Argument
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>		The result is the value of evaluating the argument.
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	.Method [Argument...]
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		The method can be alone or the last element of a chain but,
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>		unlike methods in the middle of a chain, it can take arguments.
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>		The result is the value of calling the method with the
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>		arguments:
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			dot.Method(Argument1, etc.)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	functionName [Argument...]
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>		The result is the value of calling the function associated
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>		with the name:
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			function(Argument1, etc.)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>		Functions and function names are described below.
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>A pipeline may be &#34;chained&#34; by separating a sequence of commands with pipeline
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>characters &#39;|&#39;. In a chained pipeline, the result of each command is
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>passed as the last argument of the following command. The output of the final
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>command in the pipeline is the value of the pipeline.
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>The output of a command will be either one value or two values, the second of
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>which has type error. If that second value is present and evaluates to
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>non-nil, execution terminates and the error is returned to the caller of
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>Execute.
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>Variables
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>A pipeline inside an action may initialize a variable to capture the result.
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>The initialization has syntax
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	$variable := pipeline
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>where $variable is the name of the variable. An action that declares a
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>variable produces no output.
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>Variables previously declared can also be assigned, using the syntax
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	$variable = pipeline
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>If a &#34;range&#34; action initializes a variable, the variable is set to the
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>successive elements of the iteration. Also, a &#34;range&#34; may declare two
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>variables, separated by a comma:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	range $index, $element := pipeline
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>in which case $index and $element are set to the successive values of the
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>array/slice index or map key and element, respectively. Note that if there is
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>only one variable, it is assigned the element; this is opposite to the
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>convention in Go range clauses.
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>A variable&#39;s scope extends to the &#34;end&#34; action of the control structure (&#34;if&#34;,
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>&#34;with&#34;, or &#34;range&#34;) in which it is declared, or to the end of the template if
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>there is no such control structure. A template invocation does not inherit
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>variables from the point of its invocation.
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>When execution begins, $ is set to the data argument passed to Execute, that is,
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>to the starting value of dot.
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>Examples
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>Here are some example one-line templates demonstrating pipelines and variables.
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>All produce the quoted word &#34;output&#34;:
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>	{{&#34;\&#34;output\&#34;&#34;}}
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		A string constant.
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>	{{`&#34;output&#34;`}}
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		A raw string constant.
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	{{printf &#34;%q&#34; &#34;output&#34;}}
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>		A function call.
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	{{&#34;output&#34; | printf &#34;%q&#34;}}
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>		A function call whose final argument comes from the previous
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>		command.
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>	{{printf &#34;%q&#34; (print &#34;out&#34; &#34;put&#34;)}}
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>		A parenthesized argument.
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>	{{&#34;put&#34; | printf &#34;%s%s&#34; &#34;out&#34; | printf &#34;%q&#34;}}
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>		A more elaborate call.
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>	{{&#34;output&#34; | printf &#34;%s&#34; | printf &#34;%q&#34;}}
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>		A longer chain.
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>	{{with &#34;output&#34;}}{{printf &#34;%q&#34; .}}{{end}}
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>		A with action using dot.
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>	{{with $x := &#34;output&#34; | printf &#34;%q&#34;}}{{$x}}{{end}}
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>		A with action that creates and uses a variable.
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>	{{with $x := &#34;output&#34;}}{{printf &#34;%q&#34; $x}}{{end}}
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>		A with action that uses the variable in another action.
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>	{{with $x := &#34;output&#34;}}{{$x | printf &#34;%q&#34;}}{{end}}
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>		The same, but pipelined.
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>Functions
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>During execution functions are found in two function maps: first in the
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>template, then in the global function map. By default, no functions are defined
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>in the template but the Funcs method can be used to add them.
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>Predefined global functions are named as follows.
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>	and
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		Returns the boolean AND of its arguments by returning the
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>		first empty argument or the last argument. That is,
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>		&#34;and x y&#34; behaves as &#34;if x then y else x.&#34;
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>		Evaluation proceeds through the arguments left to right
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>		and returns when the result is determined.
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>	call
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>		Returns the result of calling the first argument, which
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>		must be a function, with the remaining arguments as parameters.
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>		Thus &#34;call .X.Y 1 2&#34; is, in Go notation, dot.X.Y(1, 2) where
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		Y is a func-valued field, map entry, or the like.
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		The first argument must be the result of an evaluation
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		that yields a value of function type (as distinct from
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		a predefined function such as print). The function must
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		return either one or two result values, the second of which
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>		is of type error. If the arguments don&#39;t match the function
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>		or the returned error value is non-nil, execution stops.
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	html
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		Returns the escaped HTML equivalent of the textual
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		representation of its arguments. This function is unavailable
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		in html/template, with a few exceptions.
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>	index
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		Returns the result of indexing its first argument by the
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>		following arguments. Thus &#34;index x 1 2 3&#34; is, in Go syntax,
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>		x[1][2][3]. Each indexed item must be a map, slice, or array.
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	slice
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>		slice returns the result of slicing its first argument by the
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>		remaining arguments. Thus &#34;slice x 1 2&#34; is, in Go syntax, x[1:2],
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>		while &#34;slice x&#34; is x[:], &#34;slice x 1&#34; is x[1:], and &#34;slice x 1 2 3&#34;
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>		is x[1:2:3]. The first argument must be a string, slice, or array.
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	js
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>		Returns the escaped JavaScript equivalent of the textual
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>		representation of its arguments.
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>	len
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>		Returns the integer length of its argument.
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	not
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>		Returns the boolean negation of its single argument.
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>	or
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>		Returns the boolean OR of its arguments by returning the
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>		first non-empty argument or the last argument, that is,
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>		&#34;or x y&#34; behaves as &#34;if x then x else y&#34;.
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>		Evaluation proceeds through the arguments left to right
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>		and returns when the result is determined.
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>	print
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>		An alias for fmt.Sprint
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>	printf
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>		An alias for fmt.Sprintf
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	println
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		An alias for fmt.Sprintln
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	urlquery
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>		Returns the escaped value of the textual representation of
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>		its arguments in a form suitable for embedding in a URL query.
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>		This function is unavailable in html/template, with a few
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		exceptions.
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>The boolean functions take any zero value to be false and a non-zero
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>value to be true.
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>There is also a set of binary comparison operators defined as
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>functions:
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>	eq
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 == arg2
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>	ne
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 != arg2
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	lt
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 &lt; arg2
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	le
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 &lt;= arg2
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>	gt
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 &gt; arg2
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	ge
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		Returns the boolean truth of arg1 &gt;= arg2
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>For simpler multi-way equality tests, eq (only) accepts two or more
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>arguments and compares the second and subsequent to the first,
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>returning in effect
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	arg1==arg2 || arg1==arg3 || arg1==arg4 ...
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>(Unlike with || in Go, however, eq is a function call and all the
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>arguments will be evaluated.)
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>The comparison functions work on any values whose type Go defines as
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>comparable. For basic types such as integers, the rules are relaxed:
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>size and exact type are ignored, so any integer value, signed or unsigned,
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>may be compared with any other integer value. (The arithmetic value is compared,
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>not the bit pattern, so all negative integers are less than all unsigned integers.)
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>However, as usual, one may not compare an int with a float32 and so on.
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>Associated templates
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>Each template is named by a string specified when it is created. Also, each
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>template is associated with zero or more other templates that it may invoke by
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>name; such associations are transitive and form a name space of templates.
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>A template may use a template invocation to instantiate another associated
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>template; see the explanation of the &#34;template&#34; action above. The name must be
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>that of a template associated with the template that contains the invocation.
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>Nested template definitions
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>When parsing a template, another template may be defined and associated with the
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>template being parsed. Template definitions must appear at the top level of the
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>template, much like global variables in a Go program.
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>The syntax of such definitions is to surround each template declaration with a
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>&#34;define&#34; and &#34;end&#34; action.
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>The define action names the template being created by providing a string
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>constant. Here is a simple example:
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>	{{define &#34;T1&#34;}}ONE{{end}}
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	{{define &#34;T2&#34;}}TWO{{end}}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>	{{define &#34;T3&#34;}}{{template &#34;T1&#34;}} {{template &#34;T2&#34;}}{{end}}
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	{{template &#34;T3&#34;}}
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>This defines two templates, T1 and T2, and a third T3 that invokes the other two
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>when it is executed. Finally it invokes T3. If executed this template will
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>produce the text
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>	ONE TWO
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>By construction, a template may reside in only one association. If it&#39;s
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>necessary to have a template addressable from multiple associations, the
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>template definition must be parsed multiple times to create distinct *Template
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>values, or must be copied with [Template.Clone] or [Template.AddParseTree].
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>Parse may be called multiple times to assemble the various associated templates;
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>see [ParseFiles], [ParseGlob], [Template.ParseFiles] and [Template.ParseGlob]
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>for simple ways to parse related templates stored in files.
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>A template may be executed directly or through [Template.ExecuteTemplate], which executes
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>an associated template identified by name. To invoke our example above, we
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>might write,
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>	err := tmpl.Execute(os.Stdout, &#34;no data needed&#34;)
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>	if err != nil {
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		log.Fatalf(&#34;execution failed: %s&#34;, err)
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>or to invoke a particular template explicitly by name,
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	err := tmpl.ExecuteTemplate(os.Stdout, &#34;T2&#34;, &#34;no data needed&#34;)
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>	if err != nil {
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>		log.Fatalf(&#34;execution failed: %s&#34;, err)
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>	}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>*/</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>package template
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>
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
