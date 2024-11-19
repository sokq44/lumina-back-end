<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/backend/middleware/middleware.go - Go Documentation Server</title>

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
<a href="middleware.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/backend">backend</a>/<a href="http://localhost:8080/src/backend/middleware">middleware</a>/<span class="text-muted">middleware.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/backend/middleware">backend/middleware</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// This package is responsible for all kinds of middleware through which a request goes before the desired handler function.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Right now middlewares available are:</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">//   - authentication middleware,</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span><span class="comment">//   - HTTP Method miedleware,</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package middleware
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;backend/config&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;backend/utils/database&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	&#34;backend/utils/errhandle&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>	&#34;backend/utils/jwt&#34;
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>	&#34;fmt&#34;
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>	&#34;net/http&#34;
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>	&#34;time&#34;
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>)
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span><span class="comment">// -</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L19" class="ln">    19&nbsp;&nbsp;</span><span class="comment">//	This middleware is responsible for securing any sensitive endpoint. It operates on the access and refresh tokens which</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span><span class="comment">//	are generated upon logging into the user&#39;s account and store in HTTP-Only Cookies. In short what it does is:</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span><span class="comment">//	  - checking how to tokens were signed (whether they were generated on the server),</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span><span class="comment">//	  - checking whether the refresh token has expired (if so, it replies with status 401)</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span><span class="comment">//	  - checking whether the refresh token contains valid data (whether it&#39;s assigned to the same person as stored in the payload),</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span><span class="comment">//	  - checking whether the access token has expired and if so, generating another one (we&#39;ve already checked whether the person</span>
<span id="L25" class="ln">    25&nbsp;&nbsp;</span><span class="comment">//	    is who he says he is in all of the points above)</span>
<span id="L26" class="ln">    26&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L27" class="ln">    27&nbsp;&nbsp;</span><span class="comment">//	Methods: (All methods are accepted)</span>
<span id="L28" class="ln">    28&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L29" class="ln">    29&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span><span class="comment">//		400 (Bad Request): At least one of the tokens was generated in a bad way (it doesn&#39;t have exactly 3 parts).</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L32" class="ln">    32&nbsp;&nbsp;</span><span class="comment">//		401 (Unauthorized): There&#39;s no access or refresh token in the HTTP-Only Cookie. At least one token wasn&#39;t generated on</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span><span class="comment">//		this server.</span>
<span id="L34" class="ln">    34&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span><span class="comment">//		404 (Not Found): There&#39;s no such refresh token database record associated with the user given in the refresh token.</span>
<span id="L36" class="ln">    36&nbsp;&nbsp;</span><span class="comment">//		There&#39;s no user with the id given in the access token.</span>
<span id="L37" class="ln">    37&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L38" class="ln">    38&nbsp;&nbsp;</span><span class="comment">//		500 (Internal Server Error): There&#39;s been an error while decoding the access or refresh token cookie. There&#39;s been an</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span><span class="comment">//		error while decoding at least one of the token&#39;s payload. There&#39;s been an error while deleting the refresh token from</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span><span class="comment">//		the database. There&#39;s been an error while trying to identify the user based on his refresh or access token. There&#39;s been</span>
<span id="L41" class="ln">    41&nbsp;&nbsp;</span><span class="comment">//		an error while trying to generate a new access token.</span>
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>func Authenticate(next http.HandlerFunc) http.HandlerFunc {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	return func(w http.ResponseWriter, r *http.Request) {
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		db := database.GetDb()
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>		now := time.Now()
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>		accessToken, refreshToken, e := getRefAccFromRequest(r)
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>			return
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>		}
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>		if !jwt.WasGeneratedWithSecret(refreshToken, config.JwtSecret) || !jwt.WasGeneratedWithSecret(accessToken, config.JwtSecret) {
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>			e := errhandle.Error{
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>				Type:          errhandle.JwtError,
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>				ServerMessage: &#34;one of the tokens or both weren&#39;t created with the server secret&#34;,
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>				ClientMessage: &#34;Your authentication medium wasn&#39;t generated by this server.&#34;,
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>				Status:        http.StatusUnauthorized,
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>			}
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>			e.Handle(w, r)
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>			return
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>		}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>		claimsRefresh, e := jwt.DecodePayload(refreshToken)
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>			return
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>		}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>		claimsAccess, e := jwt.DecodePayload(accessToken)
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>			return
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>		}
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>		expiresRefresh := int64(claimsRefresh[&#34;exp&#34;].(float64))
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if expiresRefresh &lt; now.Unix() {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			e := db.DeleteRefreshTokenByToken(refreshToken)
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>			if e.Handle(w, r) {
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>				return
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>			}
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>			http.SetCookie(w, &amp;http.Cookie{
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>				Name:     &#34;refresh_token&#34;,
<span id="L82" class="ln">    82&nbsp;&nbsp;</span>				Value:    &#34;&#34;,
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>				HttpOnly: true,
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>				Path:     &#34;/&#34;,
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>				Expires:  time.Unix(0, 0),
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			})
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>			http.SetCookie(w, &amp;http.Cookie{
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>				Name:     &#34;access_token&#34;,
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>				Value:    &#34;&#34;,
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>				HttpOnly: true,
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>				Path:     &#34;/&#34;,
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>				Expires:  time.Unix(0, 0),
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>			})
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>			w.WriteHeader(http.StatusUnauthorized)
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>			return
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>		}
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		userId := claimsRefresh[&#34;user&#34;].(string)
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>		tk, e := db.GetRefreshTokenByUserId(userId)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>		if e.Handle(w, r) {
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>			return
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>		}
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>		if tk.UserId != userId {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>			w.WriteHeader(http.StatusUnauthorized)
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>			return
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>		}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>		expiresAccess := int64(claimsAccess[&#34;exp&#34;].(float64))
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if expiresAccess &lt; now.Unix() {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			user, e := db.GetUserById(claimsAccess[&#34;user&#34;].(string))
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>			if e.Handle(w, r) {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>				return
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>			}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>			access, e := jwt.GenerateAccessToken(user.Id, now)
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>			if e.Handle(w, r) {
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>				return
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>			}
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>			http.SetCookie(w, &amp;http.Cookie{
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>				Name:     &#34;access_token&#34;,
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>				Value:    access,
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>				HttpOnly: true,
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>				Path:     &#34;/&#34;,
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>				Expires:  now.Add(time.Duration(config.JwtAccExpTime)),
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>			})
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		}
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>		next(w, r)
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>	}
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span><span class="comment">// -</span>
<span id="L137" class="ln">   137&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span><span class="comment">//	This middleware is responsible for making sure that the request sent to a certain endpoint has a valid HTTP method.</span>
<span id="L139" class="ln">   139&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L140" class="ln">   140&nbsp;&nbsp;</span><span class="comment">//	Methods: (The one specified in the argument)</span>
<span id="L141" class="ln">   141&nbsp;&nbsp;</span><span class="comment">//</span>
<span id="L142" class="ln">   142&nbsp;&nbsp;</span><span class="comment">//	Possible Responses:</span>
<span id="L143" class="ln">   143&nbsp;&nbsp;</span><span class="comment">//		405 (Method Not Allowed): When the request is of an unaccepted method</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>func Method(method string, next http.HandlerFunc) http.HandlerFunc {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>	return func(w http.ResponseWriter, r *http.Request) {
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>		if r.Method != method {
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>			w.WriteHeader(http.StatusMethodNotAllowed)
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>			return
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>		}
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>		next(w, r)
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>	}
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>}
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>func getRefAccFromRequest(r *http.Request) (string, string, *errhandle.Error) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>	access, err := r.Cookie(&#34;access_token&#34;)
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>	if err == http.ErrNoCookie {
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>		return &#34;&#34;, &#34;&#34;, &amp;errhandle.Error{
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>			Type:          errhandle.JwtError,
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>			ServerMessage: &#34;no access_token cookie present&#34;,
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>			ClientMessage: &#34;There was no authentication medium present in the request.&#34;,
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>			Status:        http.StatusUnauthorized,
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>		}
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>	} else if err != nil {
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>		return &#34;&#34;, &#34;&#34;, &amp;errhandle.Error{
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>			Type:          errhandle.JwtError,
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;while trying to retrieve the access_token cookie -&gt; %v&#34;, err),
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>			Status:        http.StatusInternalServerError,
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>		}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>	}
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>	refresh, err := r.Cookie(&#34;refresh_token&#34;)
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>	if err == http.ErrNoCookie {
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>		return &#34;&#34;, &#34;&#34;, &amp;errhandle.Error{
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>			Type:          errhandle.JwtError,
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			ServerMessage: &#34;no refresh_token cookie present&#34;,
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>			ClientMessage: &#34;There was no authentication medium present in the request.&#34;,
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			Status:        http.StatusUnauthorized,
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>		}
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>	} else if err != nil {
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>		return &#34;&#34;, &#34;&#34;, &amp;errhandle.Error{
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			Type:          errhandle.JwtError,
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			ServerMessage: fmt.Sprintf(&#34;while trying to retrieve the refresh_token cookie -&gt; %v&#34;, err),
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>			ClientMessage: &#34;An error has occurred while processing your request.&#34;,
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>			Status:        http.StatusInternalServerError,
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>		}
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>	}
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>	return access.Value, refresh.Value, nil
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>}
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>
</pre><p><a href="middleware.go?m=text">View as plain text</a></p>

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
