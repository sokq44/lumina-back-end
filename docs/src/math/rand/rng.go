<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/rand/rng.go - Go Documentation Server</title>

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
<a href="rng.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/rand">rand</a>/<span class="text-muted">rng.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/rand">math/rand</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package rand
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">/*
<span id="L8" class="ln">     8&nbsp;&nbsp;</span> * Uniform distribution
<span id="L9" class="ln">     9&nbsp;&nbsp;</span> *
<span id="L10" class="ln">    10&nbsp;&nbsp;</span> * algorithm by
<span id="L11" class="ln">    11&nbsp;&nbsp;</span> * DP Mitchell and JA Reeds
<span id="L12" class="ln">    12&nbsp;&nbsp;</span> */</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>const (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	rngLen   = 607
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	rngTap   = 273
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	rngMax   = 1 &lt;&lt; 63
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	rngMask  = rngMax - 1
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	int32max = (1 &lt;&lt; 31) - 1
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>)
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>var (
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// rngCooked used for seeding. See gen_cooked.go for details.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	rngCooked [rngLen]int64 = [...]int64{
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		-4181792142133755926, -4576982950128230565, 1395769623340756751, 5333664234075297259,
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		-6347679516498800754, 9033628115061424579, 7143218595135194537, 4812947590706362721,
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		7937252194349799378, 5307299880338848416, 8209348851763925077, -7107630437535961764,
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		4593015457530856296, 8140875735541888011, -5903942795589686782, -603556388664454774,
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>		-7496297993371156308, 113108499721038619, 4569519971459345583, -4160538177779461077,
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>		-6835753265595711384, -6507240692498089696, 6559392774825876886, 7650093201692370310,
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>		7684323884043752161, -8965504200858744418, -2629915517445760644, 271327514973697897,
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		-6433985589514657524, 1065192797246149621, 3344507881999356393, -4763574095074709175,
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>		7465081662728599889, 1014950805555097187, -4773931307508785033, -5742262670416273165,
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>		2418672789110888383, 5796562887576294778, 4484266064449540171, 3738982361971787048,
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		-4699774852342421385, 10530508058128498, -589538253572429690, -6598062107225984180,
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		8660405965245884302, 10162832508971942, -2682657355892958417, 7031802312784620857,
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>		6240911277345944669, 831864355460801054, -1218937899312622917, 2116287251661052151,
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>		2202309800992166967, 9161020366945053561, 4069299552407763864, 4936383537992622449,
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>		457351505131524928, -8881176990926596454, -6375600354038175299, -7155351920868399290,
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>		4368649989588021065, 887231587095185257, -3659780529968199312, -2407146836602825512,
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		5616972787034086048, -751562733459939242, 1686575021641186857, -5177887698780513806,
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		-4979215821652996885, -1375154703071198421, 5632136521049761902, -8390088894796940536,
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>		-193645528485698615, -5979788902190688516, -4907000935050298721, -285522056888777828,
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		-2776431630044341707, 1679342092332374735, 6050638460742422078, -2229851317345194226,
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		-1582494184340482199, 5881353426285907985, 812786550756860885, 4541845584483343330,
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>		-6497901820577766722, 4980675660146853729, -4012602956251539747, -329088717864244987,
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		-2896929232104691526, 1495812843684243920, -2153620458055647789, 7370257291860230865,
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		-2466442761497833547, 4706794511633873654, -1398851569026877145, 8549875090542453214,
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>		-9189721207376179652, -7894453601103453165, 7297902601803624459, 1011190183918857495,
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		-6985347000036920864, 5147159997473910359, -8326859945294252826, 2659470849286379941,
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>		6097729358393448602, -7491646050550022124, -5117116194870963097, -896216826133240300,
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		-745860416168701406, 5803876044675762232, -787954255994554146, -3234519180203704564,
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>		-4507534739750823898, -1657200065590290694, 505808562678895611, -4153273856159712438,
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>		-8381261370078904295, 572156825025677802, 1791881013492340891, 3393267094866038768,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>		-5444650186382539299, 2352769483186201278, -7930912453007408350, -325464993179687389,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>		-3441562999710612272, -6489413242825283295, 5092019688680754699, -227247482082248967,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>		4234737173186232084, 5027558287275472836, 4635198586344772304, -536033143587636457,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		5907508150730407386, -8438615781380831356, 972392927514829904, -3801314342046600696,
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		-4064951393885491917, -174840358296132583, 2407211146698877100, -1640089820333676239,
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		3940796514530962282, -5882197405809569433, 3095313889586102949, -1818050141166537098,
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		5832080132947175283, 7890064875145919662, 8184139210799583195, -8073512175445549678,
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>		-7758774793014564506, -4581724029666783935, 3516491885471466898, -8267083515063118116,
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		6657089965014657519, 5220884358887979358, 1796677326474620641, 5340761970648932916,
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		1147977171614181568, 5066037465548252321, 2574765911837859848, 1085848279845204775,
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		-5873264506986385449, 6116438694366558490, 2107701075971293812, -7420077970933506541,
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		2469478054175558874, -1855128755834809824, -5431463669011098282, -9038325065738319171,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>		-6966276280341336160, 7217693971077460129, -8314322083775271549, 7196649268545224266,
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		-3585711691453906209, -5267827091426810625, 8057528650917418961, -5084103596553648165,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		-2601445448341207749, -7850010900052094367, 6527366231383600011, 3507654575162700890,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		9202058512774729859, 1954818376891585542, -2582991129724600103, 8299563319178235687,
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		-5321504681635821435, 7046310742295574065, -2376176645520785576, -7650733936335907755,
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>		8850422670118399721, 3631909142291992901, 5158881091950831288, -6340413719511654215,
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		4763258931815816403, 6280052734341785344, -4979582628649810958, 2043464728020827976,
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		-2678071570832690343, 4562580375758598164, 5495451168795427352, -7485059175264624713,
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>		553004618757816492, 6895160632757959823, -989748114590090637, 7139506338801360852,
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		-672480814466784139, 5535668688139305547, 2430933853350256242, -3821430778991574732,
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		-1063731997747047009, -3065878205254005442, 7632066283658143750, 6308328381617103346,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		3681878764086140361, 3289686137190109749, 6587997200611086848, 244714774258135476,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		-5143583659437639708, 8090302575944624335, 2945117363431356361, -8359047641006034763,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>		3009039260312620700, -793344576772241777, 401084700045993341, -1968749590416080887,
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>		4707864159563588614, -3583123505891281857, -3240864324164777915, -5908273794572565703,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		-3719524458082857382, -5281400669679581926, 8118566580304798074, 3839261274019871296,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>		7062410411742090847, -8481991033874568140, 6027994129690250817, -6725542042704711878,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>		-2971981702428546974, -7854441788951256975, 8809096399316380241, 6492004350391900708,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		2462145737463489636, -8818543617934476634, -5070345602623085213, -8961586321599299868,
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>		-3758656652254704451, -8630661632476012791, 6764129236657751224, -709716318315418359,
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		-3403028373052861600, -8838073512170985897, -3999237033416576341, -2920240395515973663,
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>		-2073249475545404416, 368107899140673753, -6108185202296464250, -6307735683270494757,
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>		4782583894627718279, 6718292300699989587, 8387085186914375220, 3387513132024756289,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		4654329375432538231, -292704475491394206, -3848998599978456535, 7623042350483453954,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>		7725442901813263321, 9186225467561587250, -5132344747257272453, -6865740430362196008,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>		2530936820058611833, 1636551876240043639, -3658707362519810009, 1452244145334316253,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>		-7161729655835084979, -7943791770359481772, 9108481583171221009, -3200093350120725999,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>		5007630032676973346, 2153168792952589781, 6720334534964750538, -3181825545719981703,
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		3433922409283786309, 2285479922797300912, 3110614940896576130, -2856812446131932915,
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		-3804580617188639299, 7163298419643543757, 4891138053923696990, 580618510277907015,
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>		1684034065251686769, 4429514767357295841, -8893025458299325803, -8103734041042601133,
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		7177515271653460134, 4589042248470800257, -1530083407795771245, 143607045258444228,
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		246994305896273627, -8356954712051676521, 6473547110565816071, 3092379936208876896,
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		2058427839513754051, -4089587328327907870, 8785882556301281247, -3074039370013608197,
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		-637529855400303673, 6137678347805511274, -7152924852417805802, 5708223427705576541,
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		-3223714144396531304, 4358391411789012426, 325123008708389849, 6837621693887290924,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		4843721905315627004, -3212720814705499393, -3825019837890901156, 4602025990114250980,
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		1044646352569048800, 9106614159853161675, -8394115921626182539, -4304087667751778808,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>		2681532557646850893, 3681559472488511871, -3915372517896561773, -2889241648411946534,
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		-6564663803938238204, -8060058171802589521, 581945337509520675, 3648778920718647903,
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		-4799698790548231394, -7602572252857820065, 220828013409515943, -1072987336855386047,
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		4287360518296753003, -4633371852008891965, 5513660857261085186, -2258542936462001533,
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		-8744380348503999773, 8746140185685648781, 228500091334420247, 1356187007457302238,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>		3019253992034194581, 3152601605678500003, -8793219284148773595, 5559581553696971176,
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		4916432985369275664, -8559797105120221417, -5802598197927043732, 2868348622579915573,
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		-7224052902810357288, -5894682518218493085, 2587672709781371173, -7706116723325376475,
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		3092343956317362483, -5561119517847711700, 972445599196498113, -1558506600978816441,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		1708913533482282562, -2305554874185907314, -6005743014309462908, -6653329009633068701,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		-483583197311151195, 2488075924621352812, -4529369641467339140, -4663743555056261452,
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>		2997203966153298104, 1282559373026354493, 240113143146674385, 8665713329246516443,
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>		628141331766346752, -4651421219668005332, -7750560848702540400, 7596648026010355826,
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		-3132152619100351065, 7834161864828164065, 7103445518877254909, 4390861237357459201,
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		-4780718172614204074, -319889632007444440, 622261699494173647, -3186110786557562560,
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>		-8718967088789066690, -1948156510637662747, -8212195255998774408, -7028621931231314745,
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>		2623071828615234808, -4066058308780939700, -5484966924888173764, -6683604512778046238,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		-6756087640505506466, 5256026990536851868, 7841086888628396109, 6640857538655893162,
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		-8021284697816458310, -7109857044414059830, -1689021141511844405, -4298087301956291063,
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		-4077748265377282003, -998231156719803476, 2719520354384050532, 9132346697815513771,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		4332154495710163773, -2085582442760428892, 6994721091344268833, -2556143461985726874,
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>		-8567931991128098309, 59934747298466858, -3098398008776739403, -265597256199410390,
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>		2332206071942466437, -7522315324568406181, 3154897383618636503, -7585605855467168281,
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>		-6762850759087199275, 197309393502684135, -8579694182469508493, 2543179307861934850,
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>		4350769010207485119, -4468719947444108136, -7207776534213261296, -1224312577878317200,
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		4287946071480840813, 8362686366770308971, 6486469209321732151, -5605644191012979782,
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		-1669018511020473564, 4450022655153542367, -7618176296641240059, -3896357471549267421,
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		-4596796223304447488, -6531150016257070659, -8982326463137525940, -4125325062227681798,
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>		-1306489741394045544, -8338554946557245229, 5329160409530630596, 7790979528857726136,
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>		4955070238059373407, -4304834761432101506, -6215295852904371179, 3007769226071157901,
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>		-6753025801236972788, 8928702772696731736, 7856187920214445904, -4748497451462800923,
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		7900176660600710914, -7082800908938549136, -6797926979589575837, -6737316883512927978,
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>		4186670094382025798, 1883939007446035042, -414705992779907823, 3734134241178479257,
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>		4065968871360089196, 6953124200385847784, -7917685222115876751, -7585632937840318161,
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>		-5567246375906782599, -5256612402221608788, 3106378204088556331, -2894472214076325998,
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>		4565385105440252958, 1979884289539493806, -6891578849933910383, 3783206694208922581,
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>		8464961209802336085, 2843963751609577687, 3030678195484896323, -4429654462759003204,
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>		4459239494808162889, 402587895800087237, 8057891408711167515, 4541888170938985079,
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>		1042662272908816815, -3666068979732206850, 2647678726283249984, 2144477441549833761,
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>		-3417019821499388721, -2105601033380872185, 5916597177708541638, -8760774321402454447,
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>		8833658097025758785, 5970273481425315300, 563813119381731307, -6455022486202078793,
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		1598828206250873866, -4016978389451217698, -2988328551145513985, -6071154634840136312,
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>		8469693267274066490, 125672920241807416, -3912292412830714870, -2559617104544284221,
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>		-486523741806024092, -4735332261862713930, 5923302823487327109, -9082480245771672572,
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		-1808429243461201518, 7990420780896957397, 4317817392807076702, 3625184369705367340,
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>		-6482649271566653105, -3480272027152017464, -3225473396345736649, -368878695502291645,
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		-3981164001421868007, -8522033136963788610, 7609280429197514109, 3020985755112334161,
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>		-2572049329799262942, 2635195723621160615, 5144520864246028816, -8188285521126945980,
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		1567242097116389047, 8172389260191636581, -2885551685425483535, -7060359469858316883,
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>		-6480181133964513127, -7317004403633452381, 6011544915663598137, 5932255307352610768,
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>		2241128460406315459, -8327867140638080220, 3094483003111372717, 4583857460292963101,
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>		9079887171656594975, -384082854924064405, -3460631649611717935, 4225072055348026230,
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>		-7385151438465742745, 3801620336801580414, -399845416774701952, -7446754431269675473,
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		7899055018877642622, 5421679761463003041, 5521102963086275121, -4975092593295409910,
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>		8735487530905098534, -7462844945281082830, -2080886987197029914, -1000715163927557685,
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>		-4253840471931071485, -5828896094657903328, 6424174453260338141, 359248545074932887,
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>		-5949720754023045210, -2426265837057637212, 3030918217665093212, -9077771202237461772,
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>		-3186796180789149575, 740416251634527158, -2142944401404840226, 6951781370868335478,
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		399922722363687927, -8928469722407522623, -1378421100515597285, -8343051178220066766,
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		-3030716356046100229, -8811767350470065420, 9026808440365124461, 6440783557497587732,
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		4615674634722404292, 539897290441580544, 2096238225866883852, 8751955639408182687,
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		-7316147128802486205, 7381039757301768559, 6157238513393239656, -1473377804940618233,
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>		8629571604380892756, 5280433031239081479, 7101611890139813254, 2479018537985767835,
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>		7169176924412769570, -1281305539061572506, -7865612307799218120, 2278447439451174845,
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>		3625338785743880657, 6477479539006708521, 8976185375579272206, -3712000482142939688,
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		1326024180520890843, 7537449876596048829, 5464680203499696154, 3189671183162196045,
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		6346751753565857109, -8982212049534145501, -6127578587196093755, -245039190118465649,
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>		-6320577374581628592, 7208698530190629697, 7276901792339343736, -7490986807540332668,
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>		4133292154170828382, 2918308698224194548, -7703910638917631350, -3929437324238184044,
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>		-4300543082831323144, -6344160503358350167, 5896236396443472108, -758328221503023383,
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		-1894351639983151068, -307900319840287220, -6278469401177312761, -2171292963361310674,
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		8382142935188824023, 9103922860780351547, 4152330101494654406,
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>	}
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>)
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>type rngSource struct {
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	tap  int           <span class="comment">// index into vec</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	feed int           <span class="comment">// index into vec</span>
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	vec  [rngLen]int64 <span class="comment">// current feedback register</span>
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>}
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span><span class="comment">// seed rng x[n+1] = 48271 * x[n] mod (2**31 - 1)</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>func seedrand(x int32) int32 {
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	const (
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>		A = 48271
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>		Q = 44488
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>		R = 3399
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>	)
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>	hi := x / Q
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	lo := x % Q
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>	x = A*lo - R*hi
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>	if x &lt; 0 {
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>		x += int32max
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>	}
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>	return x
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>}
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">// Seed uses the provided seed value to initialize the generator to a deterministic state.</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>func (rng *rngSource) Seed(seed int64) {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>	rng.tap = 0
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>	rng.feed = rngLen - rngTap
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>	seed = seed % int32max
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>	if seed &lt; 0 {
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>		seed += int32max
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>	}
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>	if seed == 0 {
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>		seed = 89482311
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>	}
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>	x := int32(seed)
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>	for i := -20; i &lt; rngLen; i++ {
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>		x = seedrand(x)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>		if i &gt;= 0 {
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>			var u int64
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>			u = int64(x) &lt;&lt; 40
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			x = seedrand(x)
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>			u ^= int64(x) &lt;&lt; 20
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>			x = seedrand(x)
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			u ^= int64(x)
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>			u ^= rngCooked[i]
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			rng.vec[i] = u
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>		}
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>	}
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>}
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span><span class="comment">// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.</span>
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>func (rng *rngSource) Int63() int64 {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>	return int64(rng.Uint64() &amp; rngMask)
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span><span class="comment">// Uint64 returns a non-negative pseudo-random 64-bit integer as a uint64.</span>
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>func (rng *rngSource) Uint64() uint64 {
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	rng.tap--
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	if rng.tap &lt; 0 {
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>		rng.tap += rngLen
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	}
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	rng.feed--
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	if rng.feed &lt; 0 {
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		rng.feed += rngLen
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	}
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	x := rng.vec[rng.feed] + rng.vec[rng.tap]
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	rng.vec[rng.feed] = x
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	return uint64(x)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>}
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>
</pre><p><a href="rng.go?m=text">View as plain text</a></p>

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
