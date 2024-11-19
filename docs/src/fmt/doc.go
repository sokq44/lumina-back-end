<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/fmt/doc.go - Go Documentation Server</title>

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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/fmt">fmt</a>/<span class="text-muted">doc.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/fmt">fmt</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>Package fmt implements formatted I/O with functions analogous
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>to C&#39;s printf and scanf.  The format &#39;verbs&#39; are derived from C&#39;s but
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>are simpler.
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span># Printing
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>The verbs:
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>General:
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	%v	the value in a default format
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>		when printing structs, the plus flag (%+v) adds field names
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	%#v	a Go-syntax representation of the value
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	%T	a Go-syntax representation of the type of the value
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	%%	a literal percent sign; consumes no value
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>Boolean:
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	%t	the word true or false
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>Integer:
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	%b	base 2
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	%c	the character represented by the corresponding Unicode code point
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	%d	base 10
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	%o	base 8
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	%O	base 8 with 0o prefix
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	%q	a single-quoted character literal safely escaped with Go syntax.
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	%x	base 16, with lower-case letters for a-f
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	%X	base 16, with upper-case letters for A-F
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	%U	Unicode format: U+1234; same as &#34;U+%04X&#34;
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>Floating-point and complex constituents:
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	%b	decimalless scientific notation with exponent a power of two,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		in the manner of strconv.FormatFloat with the &#39;b&#39; format,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		e.g. -123456p-78
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	%e	scientific notation, e.g. -1.234456e+78
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	%E	scientific notation, e.g. -1.234456E+78
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	%f	decimal point but no exponent, e.g. 123.456
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	%F	synonym for %f
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	%g	%e for large exponents, %f otherwise. Precision is discussed below.
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	%G	%E for large exponents, %F otherwise
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	%x	hexadecimal notation (with decimal power of two exponent), e.g. -0x1.23abcp+20
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	%X	upper-case hexadecimal notation, e.g. -0X1.23ABCP+20
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>String and slice of bytes (treated equivalently with these verbs):
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	%s	the uninterpreted bytes of the string or slice
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	%q	a double-quoted string safely escaped with Go syntax
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	%x	base 16, lower-case, two characters per byte
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	%X	base 16, upper-case, two characters per byte
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>Slice:
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	%p	address of 0th element in base 16 notation, with leading 0x
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>Pointer:
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	%p	base 16 notation, with leading 0x
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	The %b, %d, %o, %x and %X verbs also work with pointers,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	formatting the value exactly as if it were an integer.
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>The default format for %v is:
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	bool:                    %t
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	int, int8 etc.:          %d
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	uint, uint8 etc.:        %d, %#x if printed with %#v
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	float32, complex64, etc: %g
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	string:                  %s
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	chan:                    %p
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	pointer:                 %p
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>For compound objects, the elements are printed using these rules, recursively,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>laid out like this:
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	struct:             {field0 field1 ...}
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	array, slice:       [elem0 elem1 ...]
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	maps:               map[key1:value1 key2:value2 ...]
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	pointer to above:   &amp;{}, &amp;[], &amp;map[]
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>Width is specified by an optional decimal number immediately preceding the verb.
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>If absent, the width is whatever is necessary to represent the value.
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>Precision is specified after the (optional) width by a period followed by a
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>decimal number. If no period is present, a default precision is used.
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>A period with no following number specifies a precision of zero.
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>Examples:
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	%f     default width, default precision
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	%9f    width 9, default precision
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	%.2f   default width, precision 2
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	%9.2f  width 9, precision 2
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	%9.f   width 9, precision 0
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>Width and precision are measured in units of Unicode code points,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>that is, runes. (This differs from C&#39;s printf where the
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>units are always measured in bytes.) Either or both of the flags
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>may be replaced with the character &#39;*&#39;, causing their values to be
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>obtained from the next operand (preceding the one to format),
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>which must be of type int.
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>For most values, width is the minimum number of runes to output,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>padding the formatted form with spaces if necessary.
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>For strings, byte slices and byte arrays, however, precision
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>limits the length of the input to be formatted (not the size of
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>the output), truncating if necessary. Normally it is measured in
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>runes, but for these types when formatted with the %x or %X format
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>it is measured in bytes.
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>For floating-point values, width sets the minimum width of the field and
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>precision sets the number of places after the decimal, if appropriate,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>except that for %g/%G precision sets the maximum number of significant
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>digits (trailing zeros are removed). For example, given 12.345 the format
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>%6.3f prints 12.345 while %.3g prints 12.3. The default precision for %e, %f
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>and %#g is 6; for %g it is the smallest number of digits necessary to identify
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>the value uniquely.
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>For complex numbers, the width and precision apply to the two
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>components independently and the result is parenthesized, so %f applied
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>to 1.2+3.4i produces (1.200000+3.400000i).
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>When formatting a single integer code point or a rune string (type []rune)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>with %q, invalid Unicode code points are changed to the Unicode replacement
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>character, U+FFFD, as in strconv.QuoteRune.
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>Other flags:
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	&#39;+&#39;	always print a sign for numeric values;
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		guarantee ASCII-only output for %q (%+q)
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	&#39;-&#39;	pad with spaces on the right rather than the left (left-justify the field)
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	&#39;#&#39;	alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		0x or 0X for hex (%#x or %#X); suppress 0x for %p (%#p);
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		for %q, print a raw (backquoted) string if strconv.CanBackquote
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		returns true;
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		always print a decimal point for %e, %E, %f, %F, %g and %G;
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		do not remove trailing zeros for %g and %G;
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		write e.g. U+0078 &#39;x&#39; if the character is printable for %U (%#U).
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	&#39; &#39;	(space) leave a space for elided sign in numbers (% d);
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		put spaces between bytes printing strings or slices in hex (% x, % X)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	&#39;0&#39;	pad with leading zeros rather than spaces;
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		for numbers, this moves the padding after the sign;
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		ignored for strings, byte slices and byte arrays
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>Flags are ignored by verbs that do not expect them.
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>For example there is no alternate decimal format, so %#d and %d
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>behave identically.
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>For each Printf-like function, there is also a Print function
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>that takes no format and is equivalent to saying %v for every
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>operand.  Another variant Println inserts blanks between
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>operands and appends a newline.
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>Regardless of the verb, if an operand is an interface value,
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>the internal concrete value is used, not the interface itself.
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>Thus:
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	var i interface{} = 23
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	fmt.Printf(&#34;%v\n&#34;, i)
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>will print 23.
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>Except when printed using the verbs %T and %p, special
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>formatting considerations apply for operands that implement
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>certain interfaces. In order of application:
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>1. If the operand is a reflect.Value, the operand is replaced by the
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>concrete value that it holds, and printing continues with the next rule.
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>2. If an operand implements the Formatter interface, it will
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>be invoked. In this case the interpretation of verbs and flags is
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>controlled by that implementation.
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>3. If the %v verb is used with the # flag (%#v) and the operand
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>implements the GoStringer interface, that will be invoked.
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>If the format (which is implicitly %v for Println etc.) is valid
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>for a string (%s %q %x %X), or is %v but not %#v,
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>the following two rules apply:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>4. If an operand implements the error interface, the Error method
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>will be invoked to convert the object to a string, which will then
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>be formatted as required by the verb (if any).
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>5. If an operand implements method String() string, that method
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>will be invoked to convert the object to a string, which will then
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>be formatted as required by the verb (if any).
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>For compound operands such as slices and structs, the format
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>applies to the elements of each operand, recursively, not to the
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>operand as a whole. Thus %q will quote each element of a slice
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>of strings, and %6.2f will control formatting for each element
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>of a floating-point array.
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>However, when printing a byte slice with a string-like verb
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>(%s %q %x %X), it is treated identically to a string, as a single item.
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>To avoid recursion in cases such as
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	type X string
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	func (x X) String() string { return Sprintf(&#34;&lt;%s&gt;&#34;, x) }
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>convert the value before recurring:
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	func (x X) String() string { return Sprintf(&#34;&lt;%s&gt;&#34;, string(x)) }
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>Infinite recursion can also be triggered by self-referential data
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>structures, such as a slice that contains itself as an element, if
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>that type has a String method. Such pathologies are rare, however,
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>and the package does not protect against them.
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>When printing a struct, fmt cannot and therefore does not invoke
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>formatting methods such as Error or String on unexported fields.
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span># Explicit argument indexes
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>In Printf, Sprintf, and Fprintf, the default behavior is for each
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>formatting verb to format successive arguments passed in the call.
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>However, the notation [n] immediately before the verb indicates that the
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>nth one-indexed argument is to be formatted instead. The same notation
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>before a &#39;*&#39; for a width or precision selects the argument index holding
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>the value. After processing a bracketed expression [n], subsequent verbs
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>will use arguments n+1, n+2, etc. unless otherwise directed.
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>For example,
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	fmt.Sprintf(&#34;%[2]d %[1]d\n&#34;, 11, 22)
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>will yield &#34;22 11&#34;, while
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	fmt.Sprintf(&#34;%[3]*.[2]*[1]f&#34;, 12.0, 2, 6)
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>equivalent to
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	fmt.Sprintf(&#34;%6.2f&#34;, 12.0)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>will yield &#34; 12.00&#34;. Because an explicit index affects subsequent verbs,
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>this notation can be used to print the same values multiple times
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>by resetting the index for the first argument to be repeated:
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	fmt.Sprintf(&#34;%d %d %#[1]x %#x&#34;, 16, 17)
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>will yield &#34;16 17 0x10 0x11&#34;.
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>
<span id="L250" class="ln">   250&nbsp;&nbsp;</span># Format errors
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>If an invalid argument is given for a verb, such as providing
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>a string to %d, the generated string will contain a
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>description of the problem, as in these examples:
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	Wrong type or unknown verb: %!verb(type=value)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>		Printf(&#34;%d&#34;, &#34;hi&#34;):        %!d(string=hi)
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	Too many arguments: %!(EXTRA type=value)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>		Printf(&#34;hi&#34;, &#34;guys&#34;):      hi%!(EXTRA string=guys)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	Too few arguments: %!verb(MISSING)
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		Printf(&#34;hi%d&#34;):            hi%!d(MISSING)
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	Non-int for width or precision: %!(BADWIDTH) or %!(BADPREC)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>		Printf(&#34;%*s&#34;, 4.5, &#34;hi&#34;):  %!(BADWIDTH)hi
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		Printf(&#34;%.*s&#34;, 4.5, &#34;hi&#34;): %!(BADPREC)hi
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	Invalid or invalid use of argument index: %!(BADINDEX)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		Printf(&#34;%*[2]d&#34;, 7):       %!d(BADINDEX)
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>		Printf(&#34;%.[2]d&#34;, 7):       %!d(BADINDEX)
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>All errors begin with the string &#34;%!&#34; followed sometimes
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>by a single character (the verb) and end with a parenthesized
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>description.
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>If an Error or String method triggers a panic when called by a
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>print routine, the fmt package reformats the error message
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>from the panic, decorating it with an indication that it came
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>through the fmt package.  For example, if a String method
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>calls panic(&#34;bad&#34;), the resulting formatted message will look
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>like
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>	%!s(PANIC=bad)
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>The %!s just shows the print verb in use when the failure
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>occurred. If the panic is caused by a nil receiver to an Error
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>or String method, however, the output is the undecorated
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>string, &#34;&lt;nil&gt;&#34;.
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>
<span id="L287" class="ln">   287&nbsp;&nbsp;</span># Scanning
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>An analogous set of functions scans formatted text to yield
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>values.  Scan, Scanf and Scanln read from os.Stdin; Fscan,
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>Fscanf and Fscanln read from a specified io.Reader; Sscan,
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>Sscanf and Sscanln read from an argument string.
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>Scan, Fscan, Sscan treat newlines in the input as spaces.
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>Scanln, Fscanln and Sscanln stop scanning at a newline and
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>require that the items be followed by a newline or EOF.
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>Scanf, Fscanf, and Sscanf parse the arguments according to a
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>format string, analogous to that of Printf. In the text that
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>follows, &#39;space&#39; means any Unicode whitespace character
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>except newline.
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>In the format string, a verb introduced by the % character
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>consumes and parses input; these verbs are described in more
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>detail below. A character other than %, space, or newline in
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>the format consumes exactly that input character, which must
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>be present. A newline with zero or more spaces before it in
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>the format string consumes zero or more spaces in the input
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>followed by a single newline or the end of the input. A space
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>following a newline in the format string consumes zero or more
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>spaces in the input. Otherwise, any run of one or more spaces
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>in the format string consumes as many spaces as possible in
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>the input. Unless the run of spaces in the format string
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>appears adjacent to a newline, the run must consume at least
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>one space from the input or find the end of the input.
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>The handling of spaces and newlines differs from that of C&#39;s
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>scanf family: in C, newlines are treated as any other space,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>and it is never an error when a run of spaces in the format
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>string finds no spaces to consume in the input.
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>The verbs behave analogously to those of Printf.
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>For example, %x will scan an integer as a hexadecimal number,
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>and %v will scan the default representation format for the value.
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>The Printf verbs %p and %T and the flags # and + are not implemented.
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>For floating-point and complex values, all valid formatting verbs
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>(%b %e %E %f %F %g %G %x %X and %v) are equivalent and accept
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>both decimal and hexadecimal notation (for example: &#34;2.3e+7&#34;, &#34;0x4.5p-8&#34;)
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>and digit-separating underscores (for example: &#34;3.14159_26535_89793&#34;).
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>Input processed by verbs is implicitly space-delimited: the
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>implementation of every verb except %c starts by discarding
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>leading spaces from the remaining input, and the %s verb
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>(and %v reading into a string) stops consuming input at the first
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>space or newline character.
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>The familiar base-setting prefixes 0b (binary), 0o and 0 (octal),
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>and 0x (hexadecimal) are accepted when scanning integers
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>without a format or with the %v verb, as are digit-separating
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>underscores.
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>Width is interpreted in the input text but there is no
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>syntax for scanning with a precision (no %5.2f, just %5f).
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>If width is provided, it applies after leading spaces are
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>trimmed and specifies the maximum number of runes to read
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>to satisfy the verb. For example,
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	Sscanf(&#34; 1234567 &#34;, &#34;%5s%d&#34;, &amp;s, &amp;i)
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>will set s to &#34;12345&#34; and i to 67 while
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>	Sscanf(&#34; 12 34 567 &#34;, &#34;%5s%d&#34;, &amp;s, &amp;i)
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>will set s to &#34;12&#34; and i to 34.
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>In all the scanning functions, a carriage return followed
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>immediately by a newline is treated as a plain newline
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>(\r\n means the same as \n).
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>In all the scanning functions, if an operand implements method
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>Scan (that is, it implements the Scanner interface) that
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>method will be used to scan the text for that operand.  Also,
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>if the number of arguments scanned is less than the number of
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>arguments provided, an error is returned.
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>All arguments to be scanned must be either pointers to basic
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>types or implementations of the Scanner interface.
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>Like Scanf and Fscanf, Sscanf need not consume its entire input.
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>There is no way to recover how much of the input string Sscanf used.
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>Note: Fscan etc. can read one character (rune) past the input
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>they return, which means that a loop calling a scan routine
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>may skip some of the input.  This is usually a problem only
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>when there is no space between input values.  If the reader
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>provided to Fscan implements ReadRune, that method will be used
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>to read characters.  If the reader also implements UnreadRune,
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>that method will be used to save the character and successive
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>calls will not lose data.  To attach ReadRune and UnreadRune
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>methods to a reader without that capability, use
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>bufio.NewReader.
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>*/</span>
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>package fmt
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>
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
