<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/backend/handlers/user.go - Go Documentation Server</title>

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
<a href="user.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/backend">backend</a>/<a href="http://localhost:8080/src/backend/handlers">handlers</a>/<span class="text-muted">user.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/backend/handlers">backend/handlers</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// This package provides all the handlers for all the endpoints the API has.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// It includes handlers for:</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span><span class="comment">//   - user registration,</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">//   - user verification (through e-mail),</span>
<span id="L6" class="ln">     6&nbsp;&nbsp;</span><span class="comment">//   - logging into one&#39;s account,</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span><span class="comment">//   - logging out from one&#39;s account,</span>
<span id="L8" class="ln">     8&nbsp;&nbsp;</span><span class="comment">//   - getting user&#39;s data,</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span><span class="comment">//   - modifying user&#39;s data,</span>
<span id="L10" class="ln">    10&nbsp;&nbsp;</span><span class="comment">//   - initializing password change,</span>
<span id="L11" class="ln">    11&nbsp;&nbsp;</span><span class="comment">//   - changing user&#39;s password</span>
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>package handlers
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>import (
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	&#34;backend/config&#34;
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	&#34;backend/models&#34;
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>	&#34;backend/utils/crypt&#34;
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	&#34;backend/utils/database&#34;
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>	&#34;backend/utils/emails&#34;
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	&#34;backend/utils/errhandle&#34;
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	&#34;backend/utils/jwt&#34;
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	&#34;encoding/json&#34;
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	&#34;net/http&#34;
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>)
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>var db *database.Database = database.GetDb()
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>var em *emails.SmtpClient = emails.GetEmails()
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">// # RegisterUser</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//	This handler allows the registratin process.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//	This means creating an unverified user, after validation, in the database and sending a verification e-mail message.</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//	Methods: POST</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//	Request Body:</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//		username: &#34;...&#34;,</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//		email: &#34;...&#34;,</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span><span class="comment">//		password: &#34;...&#34;</span>
<span id="L43" class="ln">    43&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L44" class="ln">    44&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L45" class="ln">    45&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L46" class="ln">    46&nbsp;&nbsp;</span><span class="comment">//		201 (Created): A new unverified user record was created in the database.</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Error while decing the request body or validating the sent data.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span><span class="comment">//		409 (Conflict): A user with the provided credentials already exists.</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while checking whether a user with the provided credentials already exists.</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span><span class="comment">//		Could also occur when there&#39;s been a problem while creating a new user record in the database. Can also be caused</span>
<span id="L54" class="ln">    54&nbsp;&nbsp;</span><span class="comment">//		by a problem while the generation and storing of a new email verification token. The last possible cause for this</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">//		response to happen could be an error while sending the verification e-mail. Refer to the logs for more information.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>var RegisterUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>		Username string `json:&#34;username&#34;`
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>		Email    string `json:&#34;email&#34;`
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		Password string `json:&#34;password&#34;`
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	var body RequestBody
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while decoding the request body: %v&#34;, err),
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>		}
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>			return
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		}
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>	}
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>	u := models.User{
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>		Username: body.Username,
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>		Email:    body.Email,
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>		Password: body.Password,
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>	}
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>	if u.Validate(false).Handle(w, r) {
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>		return
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>	}
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>	exists, e := db.UserExists(u)
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		return
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	if exists {
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>			Type:          errhandle.DatabaseError,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>			ServerMessage: &#34;user already exists&#34;,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>			ClientMessage: &#34;A user with these credentials already exists.&#34;,
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			Status:        http.StatusConflict,
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>		}
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>	}
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	u.Password = crypt.Sha256(body.Password)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>	if db.CreateUser(u).Handle(w, r) {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>		return
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	token, e := crypt.RandomString(128)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		return
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	duration := time.Duration(config.EmailVerTime)
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>	verification := models.EmailVerification{
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>		Token:   token,
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>		UserId:  u.Id,
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>		Expires: time.Now().Add(duration),
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>	}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>	if db.CreateEmailVerification(verification).Handle(w, r) {
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		return
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>	}
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>	if em.SendVerificationEmail(u.Email, token).Handle(w, r) {
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>		return
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>	}
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusCreated)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span><span class="comment">// # VerifyEmail</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span><span class="comment">//	This handler allows the user to verify himself with an email verification token generated earlier.</span>
<span id="L131" class="ln">   131&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span><span class="comment">//	Methods: PATCH</span>
<span id="L133" class="ln">   133&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L134" class="ln">   134&nbsp;&nbsp;</span><span class="comment">//	RequestBody:</span>
<span id="L135" class="ln">   135&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">//		token: &#34;...&#34;</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//		204 (No Content): User with the given token has been verified.</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Problem while decoding the request body.</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): No such e-mail verification token or unverified user was found in the database.</span>
<span id="L145" class="ln">   145&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L146" class="ln">   146&nbsp;&nbsp;</span><span class="comment">//		410 (Gone): The provided e-mail verification token has expired.</span>
<span id="L147" class="ln">   147&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the provided e-mail verification token. Could also be caused</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span><span class="comment">//		by an error while deleting the e-mail verification token or when verifying a user. Refer to the logs for more</span>
<span id="L150" class="ln">   150&nbsp;&nbsp;</span><span class="comment">//		information.</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>var VerifyEmail http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>		Token string `json:&#34;token&#34;`
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>	}
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	var body RequestBody
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while retrieving the access_token cookie: %v&#34;, err),
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>			return
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>		}
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>	}
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>	emailValidation, e := db.GetEmailVerificationByToken(body.Token)
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>		return
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>	}
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if emailValidation.Expires.Before(time.Now()) {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			Type:          errhandle.DatabaseError,
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			ServerMessage: &#34;email validation token has expired&#34;,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			ClientMessage: &#34;The verification link is invalid or has expired.&#34;,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			Status:        http.StatusGone,
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>		e.Handle(w, r)
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		return
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	}
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	e = db.DeleteEmailVerificationById(emailValidation.Id)
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		return
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	e = db.VerifyUser(emailValidation.UserId)
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>		return
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>	}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusNoContent)
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>}
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>
<span id="L198" class="ln">   198&nbsp;&nbsp;</span><span class="comment">// # LoginUser</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span><span class="comment">//	This handler allows the user to log into his account.</span>
<span id="L201" class="ln">   201&nbsp;&nbsp;</span><span class="comment">//	This means generating access and refresh token which are passed in HTTP-ONLY Cookies.</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span><span class="comment">//	Methods: POST</span>
<span id="L204" class="ln">   204&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L205" class="ln">   205&nbsp;&nbsp;</span><span class="comment">//	Request Body:</span>
<span id="L206" class="ln">   206&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span><span class="comment">//		email: &#34;...&#34;,</span>
<span id="L208" class="ln">   208&nbsp;&nbsp;</span><span class="comment">//		password: &#34;...&#34;</span>
<span id="L209" class="ln">   209&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L211" class="ln">   211&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L212" class="ln">   212&nbsp;&nbsp;</span><span class="comment">//		200 (OK): User has been logged in and his &#39;session&#39; started.</span>
<span id="L213" class="ln">   213&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L214" class="ln">   214&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Problem while decoding the request body.</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): User or his refresh token couldn&#39;t be found.</span>
<span id="L217" class="ln">   217&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the user or his refresh token from the database.</span>
<span id="L219" class="ln">   219&nbsp;&nbsp;</span><span class="comment">//		Could also occur when there&#39;s been a problem with access or refresh token generation. Another reason</span>
<span id="L220" class="ln">   220&nbsp;&nbsp;</span><span class="comment">//		for its occurance could be an error while storing the refresh toke in the database.</span>
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>var LoginUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>		Email    string `json:&#34;email&#34;`
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>		Password string `json:&#34;password&#34;`
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>	}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>	var body RequestBody
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while decoding the request body: %v&#34;, err),
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>		}
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>			return
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		}
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>	}
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	user, e := db.GetUserByEmail(body.Email)
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>		return
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	}
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	refreshToken, _ := db.GetRefreshTokenByUserId(user.Id)
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	if refreshToken != nil {
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>		w.WriteHeader(http.StatusOK)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>		return
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	}
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	hashedPasswd := crypt.Sha256(body.Password)
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	if !user.Verified || hashedPasswd != user.Password {
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>		w.WriteHeader(http.StatusForbidden)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>		return
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	}
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>	now := time.Now()
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>	access, e := jwt.GenerateAccessToken(user.Id, now)
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>		return
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>	}
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	refresh, e := jwt.GenerateRefreshToken(user.Id, now)
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>		return
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	e = db.CreateRefreshToken(refresh)
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		return
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	http.SetCookie(w, &amp;http.Cookie{
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>		Name:     &#34;access_token&#34;,
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>		Value:    access,
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>		HttpOnly: true,
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>		Path:     &#34;/&#34;,
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>		Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	})
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>	http.SetCookie(w, &amp;http.Cookie{
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>		Name:     &#34;refresh_token&#34;,
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>		Value:    refresh.Token,
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>		HttpOnly: true,
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>		Path:     &#34;/&#34;,
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>		Expires:  now.Add(time.Duration(config.JwtRefExpTime)),
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>	})
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusOK)
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>}
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>
<span id="L292" class="ln">   292&nbsp;&nbsp;</span><span class="comment">// # LogoutUser</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L294" class="ln">   294&nbsp;&nbsp;</span><span class="comment">//	This handler allows the user to log out from his account.</span>
<span id="L295" class="ln">   295&nbsp;&nbsp;</span><span class="comment">//	This means destroying his session (deleting the access and refresh tokens Cookies).</span>
<span id="L296" class="ln">   296&nbsp;&nbsp;</span><span class="comment">//	In order to access this endpoint the user must be logged in.</span>
<span id="L297" class="ln">   297&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L298" class="ln">   298&nbsp;&nbsp;</span><span class="comment">//	Methods: DELETE</span>
<span id="L299" class="ln">   299&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span><span class="comment">//		200 (OK): User has been logged out which means his &#39;session&#39; has been destroyed.</span>
<span id="L302" class="ln">   302&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the refresh token or while deleting it from the database.</span>
<span id="L304" class="ln">   304&nbsp;&nbsp;</span><span class="comment">//		Refer to the logs for more information.</span>
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>var LogoutUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>	refreshCookie, err := r.Cookie(&#34;refresh_token&#34;)
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>	if err != nil {
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while retrieving the refresh_token cookie: %v&#34;, err),
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>			Status:        http.StatusInternalServerError,
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>		}
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>			return
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>		}
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>	}
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>	e := db.DeleteRefreshTokenByToken(refreshCookie.Value)
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>		return
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>	}
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>	http.SetCookie(w, &amp;http.Cookie{
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>		Name:     &#34;access_token&#34;,
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>		Value:    &#34;&#34;,
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>		HttpOnly: true,
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>		Path:     &#34;/&#34;,
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>		Expires:  time.Unix(0, 0),
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>	})
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>	http.SetCookie(w, &amp;http.Cookie{
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>		Name:     &#34;refresh_token&#34;,
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>		Value:    &#34;&#34;,
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>		HttpOnly: true,
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>		Path:     &#34;/&#34;,
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>		Expires:  time.Unix(0, 0),
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>	})
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusOK)
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>}
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>
<span id="L343" class="ln">   343&nbsp;&nbsp;</span><span class="comment">// # GetUser</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L345" class="ln">   345&nbsp;&nbsp;</span><span class="comment">//	This handler allows the user to get his data.</span>
<span id="L346" class="ln">   346&nbsp;&nbsp;</span><span class="comment">//	In order to access this endpoint the user must be logged in.</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L348" class="ln">   348&nbsp;&nbsp;</span><span class="comment">//	Methods: GET</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L350" class="ln">   350&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L351" class="ln">   351&nbsp;&nbsp;</span><span class="comment">//		200 (OK): There&#39;s been no error while getting user&#39;s information.</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L353" class="ln">   353&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): When the token isn&#39;t properly formed (e.g. doesn&#39;t have three parts).</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L355" class="ln">   355&nbsp;&nbsp;</span><span class="comment">//		404 (Not found): Couldn&#39;t find the user in the the database.</span>
<span id="L356" class="ln">   356&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L357" class="ln">   357&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the access token or while unmarshaling it&#39;s payload.</span>
<span id="L358" class="ln">   358&nbsp;&nbsp;</span><span class="comment">//		Could also be caused by an error while retrieving the user from the database. Will also accur when there&#39;s</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span><span class="comment">//		an error while encoding data to the response. Refer to the logs for more information.</span>
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>var GetUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>	access, err := r.Cookie(&#34;access_token&#34;)
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	if err != nil {
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while retrieving the access_token cookie: %v&#34;, err),
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>			Status:        http.StatusInternalServerError,
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>		}
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>			return
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>		}
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>	}
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>	claims, e := jwt.DecodePayload(access.Value)
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>		return
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>	}
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>	userId := claims[&#34;user&#34;].(string)
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>	user, e := db.GetUserById(userId)
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>		return
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	}
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	userData := map[string]string{
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>		&#34;username&#34;: user.Username,
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>		&#34;email&#34;:    user.Email,
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	}
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	if err := json.NewEncoder(w).Encode(userData); err != nil {
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while encoding json data to the response: %v&#34;, err),
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>			Status:        http.StatusInternalServerError,
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>		}
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>			return
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>		}
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	}
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>}
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>
<span id="L402" class="ln">   402&nbsp;&nbsp;</span><span class="comment">// # ModifyUser</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L404" class="ln">   404&nbsp;&nbsp;</span><span class="comment">//	This handler allows the user to modify his data.</span>
<span id="L405" class="ln">   405&nbsp;&nbsp;</span><span class="comment">//	In order to access this endpoint the user must be logged in.</span>
<span id="L406" class="ln">   406&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L407" class="ln">   407&nbsp;&nbsp;</span><span class="comment">//	Methods: PATCH</span>
<span id="L408" class="ln">   408&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L409" class="ln">   409&nbsp;&nbsp;</span><span class="comment">//	Request Body:</span>
<span id="L410" class="ln">   410&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L411" class="ln">   411&nbsp;&nbsp;</span><span class="comment">//		username: &#34;...&#34;,</span>
<span id="L412" class="ln">   412&nbsp;&nbsp;</span><span class="comment">//		email: &#34;...&#34;</span>
<span id="L413" class="ln">   413&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L415" class="ln">   415&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span><span class="comment">//		200 (OK): There&#39;s been no error while modifying the user.</span>
<span id="L417" class="ln">   417&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L418" class="ln">   418&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Couldn&#39;t decode the body or the new data validation failed.</span>
<span id="L419" class="ln">   419&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L420" class="ln">   420&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): Couldn&#39;t find the user in the the database.</span>
<span id="L421" class="ln">   421&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L422" class="ln">   422&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the user from the database or while updating his data.</span>
<span id="L423" class="ln">   423&nbsp;&nbsp;</span><span class="comment">//		Refer to the logs for more information.</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>var ModifyUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>		Username string `json:&#34;username&#34;`
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>		Email    string `json:&#34;email&#34;`
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>	}
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>	var body RequestBody
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while decoding the request body: %v&#34;, err),
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>		}
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>			return
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>		}
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>	}
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>	user, e := db.GetUserByEmail(body.Email)
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>		return
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>	}
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>	var newUser models.User = models.User{
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>		Id:       user.Id,
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>		Username: body.Username,
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>		Email:    body.Email,
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>		Password: user.Password,
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>		Verified: user.Verified,
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>	}
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>	if newUser.Validate(true).Handle(w, r) {
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>		return
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>	}
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>	if db.UpdateUser(newUser).Handle(w, r) {
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>		return
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>	}
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>}
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>
<span id="L463" class="ln">   463&nbsp;&nbsp;</span><span class="comment">// # PasswordChangeInit</span>
<span id="L464" class="ln">   464&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L465" class="ln">   465&nbsp;&nbsp;</span><span class="comment">//	This handler is responsible for initiating the password change procedure for a user.</span>
<span id="L466" class="ln">   466&nbsp;&nbsp;</span><span class="comment">//	This means generating a password change token, storing it in the database and providing it</span>
<span id="L467" class="ln">   467&nbsp;&nbsp;</span><span class="comment">//	for a user through an e-mail message.</span>
<span id="L468" class="ln">   468&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span><span class="comment">//	Methods: POST</span>
<span id="L470" class="ln">   470&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L471" class="ln">   471&nbsp;&nbsp;</span><span class="comment">//	Request Body:</span>
<span id="L472" class="ln">   472&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L473" class="ln">   473&nbsp;&nbsp;</span><span class="comment">//		email: &#34;...&#34;</span>
<span id="L474" class="ln">   474&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L476" class="ln">   476&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L477" class="ln">   477&nbsp;&nbsp;</span><span class="comment">//		201 (Created): The password change procedure has been initialized, which means that a new record of password change</span>
<span id="L478" class="ln">   478&nbsp;&nbsp;</span><span class="comment">//		token has been created in the database.</span>
<span id="L479" class="ln">   479&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L480" class="ln">   480&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Problem while decoding the request body.</span>
<span id="L481" class="ln">   481&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L482" class="ln">   482&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): Couldn&#39;t find any user with the provided e-mail address.</span>
<span id="L483" class="ln">   483&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving a user with the specified e-mail. Could also be caused by</span>
<span id="L485" class="ln">   485&nbsp;&nbsp;</span><span class="comment">//		an error while generating, storing or sending the new password change token Refer to the logs for more information.</span>
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>var PasswordChangeInit http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>		Email string `json:&#34;email&#34;`
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>	}
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>	var body RequestBody
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while decoding the request body: %v&#34;, err),
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>		}
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>			return
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>		}
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>	}
<span id="L503" class="ln">   503&nbsp;&nbsp;</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span>	u, e := db.GetUserByEmail(body.Email)
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>		return
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>	}
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	token, e := crypt.RandomString(128)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>		return
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	}
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>	duration := time.Duration(config.PasswdChangeTime)
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>	passwdChange := models.PasswordChange{
<span id="L516" class="ln">   516&nbsp;&nbsp;</span>		Token:   token,
<span id="L517" class="ln">   517&nbsp;&nbsp;</span>		UserId:  u.Id,
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>		Expires: time.Now().Add(duration),
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	}
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>	if db.CreatePasswordChange(passwdChange).Handle(w, r) {
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>		return
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>	}
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>	if em.SendPasswordChangeEmail(body.Email, token).Handle(w, r) {
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		return
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>	}
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusCreated)
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>}
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span><span class="comment">// # ChangePassword</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L533" class="ln">   533&nbsp;&nbsp;</span><span class="comment">//	This handler is responsible for the password change.</span>
<span id="L534" class="ln">   534&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L535" class="ln">   535&nbsp;&nbsp;</span><span class="comment">//	Methods: PATCH</span>
<span id="L536" class="ln">   536&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L537" class="ln">   537&nbsp;&nbsp;</span><span class="comment">//	Request Body:</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span><span class="comment">//	{</span>
<span id="L539" class="ln">   539&nbsp;&nbsp;</span><span class="comment">//		password: &#34;...&#34;,</span>
<span id="L540" class="ln">   540&nbsp;&nbsp;</span><span class="comment">//		token: &#34;...&#34;</span>
<span id="L541" class="ln">   541&nbsp;&nbsp;</span><span class="comment">//	}</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">//		200 (OK): User&#39;s password has been changed.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L546" class="ln">   546&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): Problem while decoding the request body or while validating the new password.</span>
<span id="L547" class="ln">   547&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L548" class="ln">   548&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): Couldn&#39;t find the specified password change token or any user assigned to it in the database.</span>
<span id="L549" class="ln">   549&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L550" class="ln">   550&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): Problem while retrieving the desired password change token or the user assigned</span>
<span id="L551" class="ln">   551&nbsp;&nbsp;</span><span class="comment">//		to it from the database. Could also be caused by an error while deleting the password change token or updating</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span><span class="comment">//		the user assigned to it. Refer to the logs for more information.</span>
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>var ChangePassword http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>	type RequestBody struct {
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		Password string `json:&#34;password&#34;`
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>		Token    string `json:&#34;token&#34;`
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	}
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	var body RequestBody
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>		e := errhandle.Error{
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>			Type:          errhandle.HandlerError,
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;error while decoding the request body: %v&#34;, err),
<span id="L564" class="ln">   564&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>			Status:        http.StatusBadRequest,
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>		}
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>			return
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>		}
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>	}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span>	passwordChange, e := db.GetPasswordChangeByToken(body.Token)
<span id="L573" class="ln">   573&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>		return
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	}
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>	user, e := db.GetUserById(passwordChange.UserId)
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>	if e.Handle(w, r) {
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		return
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	user.Password = body.Password
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>	if user.Validate(false).Handle(w, r) {
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>		return
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	}
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>	if db.DeletePasswordChangeById(passwordChange.Id).Handle(w, r) {
<span id="L588" class="ln">   588&nbsp;&nbsp;</span>		return
<span id="L589" class="ln">   589&nbsp;&nbsp;</span>	}
<span id="L590" class="ln">   590&nbsp;&nbsp;</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span>	user.Password = crypt.Sha256(body.Password)
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>	if db.UpdateUser(*user).Handle(w, r) {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>		return
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	}
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	w.WriteHeader(http.StatusOK)
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>}
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>
</pre><p><a href="user.go?m=text">View as plain text</a></p>

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
