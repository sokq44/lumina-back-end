<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/github.com/go-sql-driver/mysql/collations.go - Go Documentation Server</title>

<link type="text/css" rel="stylesheet" href="../../../../lib/godoc/style.css">

<script>window.initFuncs = [];</script>
<script src="../../../../lib/godoc/jquery.js" defer></script>



<script>var goVersion = "go1.22.2";</script>
<script src="../../../../lib/godoc/godocs.js" defer></script>
</head>
<body>

<div id='lowframe' style="position: fixed; bottom: 0; left: 0; height: 0; width: 100%; border-top: thin solid grey; background-color: white; overflow: auto;">
...
</div><!-- #lowframe -->

<div id="topbar" class="wide"><div class="container">
<div class="top-heading" id="heading-wide"><a href="../../../../index.html">Go Documentation Server</a></div>
<div class="top-heading" id="heading-narrow"><a href="../../../../index.html">GoDoc</a></div>
<a href="collations.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/github.com">github.com</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver">go-sql-driver</a>/<a href="http://localhost:8080/src/github.com/go-sql-driver/mysql">mysql</a>/<span class="text-muted">collations.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/github.com/go-sql-driver/mysql">github.com/go-sql-driver/mysql</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Go MySQL Driver - A MySQL-Driver for Go&#39;s database/sql package</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// Copyright 2014 The Go-MySQL-Driver Authors. All rights reserved.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">// This Source Code Form is subject to the terms of the Mozilla Public</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">// License, v. 2.0. If a copy of the MPL was not distributed with this file,</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">// You can obtain one at http://mozilla.org/MPL/2.0/.</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>package mysql
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>const defaultCollation = &#34;utf8mb4_general_ci&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>const binaryCollationID = 63
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span><span class="comment">// A list of available collations mapped to the internal ID.</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span><span class="comment">// To update this map use the following MySQL query:</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">//	SELECT COLLATION_NAME, ID FROM information_schema.COLLATIONS WHERE ID&lt;256 ORDER BY ID</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">// Handshake packet have only 1 byte for collation_id.  So we can&#39;t use collations with ID &gt; 255.</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">// ucs2, utf16, and utf32 can&#39;t be used for connection charset.</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">// https://dev.mysql.com/doc/refman/5.7/en/charset-connection.html#charset-connection-impermissible-client-charset</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">// They are commented out to reduce this map.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>var collations = map[string]byte{
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;big5_chinese_ci&#34;:      1,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>	&#34;latin2_czech_cs&#34;:      2,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>	&#34;dec8_swedish_ci&#34;:      3,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>	&#34;cp850_general_ci&#34;:     4,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>	&#34;latin1_german1_ci&#34;:    5,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>	&#34;hp8_english_ci&#34;:       6,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>	&#34;koi8r_general_ci&#34;:     7,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>	&#34;latin1_swedish_ci&#34;:    8,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>	&#34;latin2_general_ci&#34;:    9,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>	&#34;swe7_swedish_ci&#34;:      10,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	&#34;ascii_general_ci&#34;:     11,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	&#34;ujis_japanese_ci&#34;:     12,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	&#34;sjis_japanese_ci&#34;:     13,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	&#34;cp1251_bulgarian_ci&#34;:  14,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	&#34;latin1_danish_ci&#34;:     15,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	&#34;hebrew_general_ci&#34;:    16,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	&#34;tis620_thai_ci&#34;:       18,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	&#34;euckr_korean_ci&#34;:      19,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	&#34;latin7_estonian_cs&#34;:   20,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	&#34;latin2_hungarian_ci&#34;:  21,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	&#34;koi8u_general_ci&#34;:     22,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	&#34;cp1251_ukrainian_ci&#34;:  23,
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	&#34;gb2312_chinese_ci&#34;:    24,
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	&#34;greek_general_ci&#34;:     25,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>	&#34;cp1250_general_ci&#34;:    26,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	&#34;latin2_croatian_ci&#34;:   27,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	&#34;gbk_chinese_ci&#34;:       28,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	&#34;cp1257_lithuanian_ci&#34;: 29,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>	&#34;latin5_turkish_ci&#34;:    30,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>	&#34;latin1_german2_ci&#34;:    31,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>	&#34;armscii8_general_ci&#34;:  32,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>	&#34;utf8_general_ci&#34;:      33,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	&#34;cp1250_czech_cs&#34;:      34,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_general_ci&#34;:          35,</span>
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	&#34;cp866_general_ci&#34;:    36,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>	&#34;keybcs2_general_ci&#34;:  37,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	&#34;macce_general_ci&#34;:    38,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	&#34;macroman_general_ci&#34;: 39,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	&#34;cp852_general_ci&#34;:    40,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	&#34;latin7_general_ci&#34;:   41,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>	&#34;latin7_general_cs&#34;:   42,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	&#34;macce_bin&#34;:           43,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	&#34;cp1250_croatian_ci&#34;:  44,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>	&#34;utf8mb4_general_ci&#34;:  45,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>	&#34;utf8mb4_bin&#34;:         46,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>	&#34;latin1_bin&#34;:          47,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>	&#34;latin1_general_ci&#34;:   48,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>	&#34;latin1_general_cs&#34;:   49,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	&#34;cp1251_bin&#34;:          50,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	&#34;cp1251_general_ci&#34;:   51,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>	&#34;cp1251_general_cs&#34;:   52,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	&#34;macroman_bin&#34;:        53,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_general_ci&#34;:         54,</span>
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_bin&#34;:                55,</span>
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16le_general_ci&#34;:       56,</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	&#34;cp1256_general_ci&#34;: 57,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	&#34;cp1257_bin&#34;:        58,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>	&#34;cp1257_general_ci&#34;: 59,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_general_ci&#34;:         60,</span>
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_bin&#34;:                61,</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16le_bin&#34;:              62,</span>
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	&#34;binary&#34;:          63,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>	&#34;armscii8_bin&#34;:    64,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	&#34;ascii_bin&#34;:       65,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	&#34;cp1250_bin&#34;:      66,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>	&#34;cp1256_bin&#34;:      67,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>	&#34;cp866_bin&#34;:       68,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>	&#34;dec8_bin&#34;:        69,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>	&#34;greek_bin&#34;:       70,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>	&#34;hebrew_bin&#34;:      71,
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>	&#34;hp8_bin&#34;:         72,
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	&#34;keybcs2_bin&#34;:     73,
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>	&#34;koi8r_bin&#34;:       74,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	&#34;koi8u_bin&#34;:       75,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	&#34;utf8_tolower_ci&#34;: 76,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>	&#34;latin2_bin&#34;:      77,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	&#34;latin5_bin&#34;:      78,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	&#34;latin7_bin&#34;:      79,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	&#34;cp850_bin&#34;:       80,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	&#34;cp852_bin&#34;:       81,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	&#34;swe7_bin&#34;:        82,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	&#34;utf8_bin&#34;:        83,
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	&#34;big5_bin&#34;:        84,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>	&#34;euckr_bin&#34;:       85,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	&#34;gb2312_bin&#34;:      86,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>	&#34;gbk_bin&#34;:         87,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	&#34;sjis_bin&#34;:        88,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	&#34;tis620_bin&#34;:      89,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_bin&#34;:                 90,</span>
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>	&#34;ujis_bin&#34;:            91,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>	&#34;geostd8_general_ci&#34;:  92,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	&#34;geostd8_bin&#34;:         93,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	&#34;latin1_spanish_ci&#34;:   94,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>	&#34;cp932_japanese_ci&#34;:   95,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	&#34;cp932_bin&#34;:           96,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	&#34;eucjpms_japanese_ci&#34;: 97,
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	&#34;eucjpms_bin&#34;:         98,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	&#34;cp1250_polish_ci&#34;:    99,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_unicode_ci&#34;:         101,</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_icelandic_ci&#34;:       102,</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_latvian_ci&#34;:         103,</span>
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_romanian_ci&#34;:        104,</span>
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_slovenian_ci&#34;:       105,</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_polish_ci&#34;:          106,</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_estonian_ci&#34;:        107,</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_spanish_ci&#34;:         108,</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_swedish_ci&#34;:         109,</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_turkish_ci&#34;:         110,</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_czech_ci&#34;:           111,</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_danish_ci&#34;:          112,</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_lithuanian_ci&#34;:      113,</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_slovak_ci&#34;:          114,</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_spanish2_ci&#34;:        115,</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_roman_ci&#34;:           116,</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_persian_ci&#34;:         117,</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_esperanto_ci&#34;:       118,</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_hungarian_ci&#34;:       119,</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_sinhala_ci&#34;:         120,</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_german2_ci&#34;:         121,</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_croatian_ci&#34;:        122,</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_unicode_520_ci&#34;:     123,</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf16_vietnamese_ci&#34;:      124,</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_unicode_ci&#34;:          128,</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_icelandic_ci&#34;:        129,</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_latvian_ci&#34;:          130,</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_romanian_ci&#34;:         131,</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_slovenian_ci&#34;:        132,</span>
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_polish_ci&#34;:           133,</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_estonian_ci&#34;:         134,</span>
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_spanish_ci&#34;:          135,</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_swedish_ci&#34;:          136,</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_turkish_ci&#34;:          137,</span>
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_czech_ci&#34;:            138,</span>
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_danish_ci&#34;:           139,</span>
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_lithuanian_ci&#34;:       140,</span>
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_slovak_ci&#34;:           141,</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_spanish2_ci&#34;:         142,</span>
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_roman_ci&#34;:            143,</span>
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_persian_ci&#34;:          144,</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_esperanto_ci&#34;:        145,</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_hungarian_ci&#34;:        146,</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_sinhala_ci&#34;:          147,</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_german2_ci&#34;:          148,</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_croatian_ci&#34;:         149,</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_unicode_520_ci&#34;:      150,</span>
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_vietnamese_ci&#34;:       151,</span>
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	<span class="comment">//&#34;ucs2_general_mysql500_ci&#34;: 159,</span>
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_unicode_ci&#34;:         160,</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_icelandic_ci&#34;:       161,</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_latvian_ci&#34;:         162,</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_romanian_ci&#34;:        163,</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_slovenian_ci&#34;:       164,</span>
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_polish_ci&#34;:          165,</span>
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_estonian_ci&#34;:        166,</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_spanish_ci&#34;:         167,</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_swedish_ci&#34;:         168,</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_turkish_ci&#34;:         169,</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_czech_ci&#34;:           170,</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_danish_ci&#34;:          171,</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_lithuanian_ci&#34;:      172,</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_slovak_ci&#34;:          173,</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_spanish2_ci&#34;:        174,</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_roman_ci&#34;:           175,</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_persian_ci&#34;:         176,</span>
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_esperanto_ci&#34;:       177,</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_hungarian_ci&#34;:       178,</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_sinhala_ci&#34;:         179,</span>
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_german2_ci&#34;:         180,</span>
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_croatian_ci&#34;:        181,</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_unicode_520_ci&#34;:     182,</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	<span class="comment">//&#34;utf32_vietnamese_ci&#34;:      183,</span>
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	&#34;utf8_unicode_ci&#34;:          192,
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	&#34;utf8_icelandic_ci&#34;:        193,
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>	&#34;utf8_latvian_ci&#34;:          194,
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	&#34;utf8_romanian_ci&#34;:         195,
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	&#34;utf8_slovenian_ci&#34;:        196,
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>	&#34;utf8_polish_ci&#34;:           197,
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>	&#34;utf8_estonian_ci&#34;:         198,
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>	&#34;utf8_spanish_ci&#34;:          199,
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>	&#34;utf8_swedish_ci&#34;:          200,
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	&#34;utf8_turkish_ci&#34;:          201,
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	&#34;utf8_czech_ci&#34;:            202,
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>	&#34;utf8_danish_ci&#34;:           203,
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	&#34;utf8_lithuanian_ci&#34;:       204,
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	&#34;utf8_slovak_ci&#34;:           205,
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>	&#34;utf8_spanish2_ci&#34;:         206,
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	&#34;utf8_roman_ci&#34;:            207,
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	&#34;utf8_persian_ci&#34;:          208,
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>	&#34;utf8_esperanto_ci&#34;:        209,
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	&#34;utf8_hungarian_ci&#34;:        210,
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>	&#34;utf8_sinhala_ci&#34;:          211,
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	&#34;utf8_german2_ci&#34;:          212,
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	&#34;utf8_croatian_ci&#34;:         213,
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>	&#34;utf8_unicode_520_ci&#34;:      214,
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>	&#34;utf8_vietnamese_ci&#34;:       215,
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>	&#34;utf8_general_mysql500_ci&#34;: 223,
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>	&#34;utf8mb4_unicode_ci&#34;:       224,
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	&#34;utf8mb4_icelandic_ci&#34;:     225,
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>	&#34;utf8mb4_latvian_ci&#34;:       226,
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>	&#34;utf8mb4_romanian_ci&#34;:      227,
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	&#34;utf8mb4_slovenian_ci&#34;:     228,
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>	&#34;utf8mb4_polish_ci&#34;:        229,
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	&#34;utf8mb4_estonian_ci&#34;:      230,
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	&#34;utf8mb4_spanish_ci&#34;:       231,
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	&#34;utf8mb4_swedish_ci&#34;:       232,
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>	&#34;utf8mb4_turkish_ci&#34;:       233,
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>	&#34;utf8mb4_czech_ci&#34;:         234,
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>	&#34;utf8mb4_danish_ci&#34;:        235,
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>	&#34;utf8mb4_lithuanian_ci&#34;:    236,
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	&#34;utf8mb4_slovak_ci&#34;:        237,
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>	&#34;utf8mb4_spanish2_ci&#34;:      238,
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>	&#34;utf8mb4_roman_ci&#34;:         239,
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>	&#34;utf8mb4_persian_ci&#34;:       240,
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	&#34;utf8mb4_esperanto_ci&#34;:     241,
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	&#34;utf8mb4_hungarian_ci&#34;:     242,
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	&#34;utf8mb4_sinhala_ci&#34;:       243,
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	&#34;utf8mb4_german2_ci&#34;:       244,
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	&#34;utf8mb4_croatian_ci&#34;:      245,
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	&#34;utf8mb4_unicode_520_ci&#34;:   246,
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	&#34;utf8mb4_vietnamese_ci&#34;:    247,
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	&#34;gb18030_chinese_ci&#34;:       248,
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	&#34;gb18030_bin&#34;:              249,
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	&#34;gb18030_unicode_520_ci&#34;:   250,
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	&#34;utf8mb4_0900_ai_ci&#34;:       255,
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span><span class="comment">// A denylist of collations which is unsafe to interpolate parameters.</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span><span class="comment">// These multibyte encodings may contains 0x5c (`\`) in their trailing bytes.</span>
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>var unsafeCollations = map[string]bool{
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	&#34;big5_chinese_ci&#34;:        true,
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	&#34;sjis_japanese_ci&#34;:       true,
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>	&#34;gbk_chinese_ci&#34;:         true,
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	&#34;big5_bin&#34;:               true,
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	&#34;gb2312_bin&#34;:             true,
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	&#34;gbk_bin&#34;:                true,
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>	&#34;sjis_bin&#34;:               true,
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	&#34;cp932_japanese_ci&#34;:      true,
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>	&#34;cp932_bin&#34;:              true,
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	&#34;gb18030_chinese_ci&#34;:     true,
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	&#34;gb18030_bin&#34;:            true,
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	&#34;gb18030_unicode_520_ci&#34;: true,
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
</pre><p><a href="collations.go?m=text">View as plain text</a></p>

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
