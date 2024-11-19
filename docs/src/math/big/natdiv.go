<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/math/big/natdiv.go - Go Documentation Server</title>

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
<a href="natdiv.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/math">math</a>/<a href="http://localhost:8080/src/math/big">big</a>/<span class="text-muted">natdiv.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/math/big">math/big</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2009 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span><span class="comment">/*
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>Multi-precision division. Here be dragons.
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>Given u and v, where u is n+m digits, and v is n digits (with no leading zeros),
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>the goal is to return quo, rem such that u = quo*v + rem, where 0 ≤ rem &lt; v.
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>That is, quo = ⌊u/v⌋ where ⌊x⌋ denotes the floor (truncation to integer) of x,
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>and rem = u - quo·v.
<span id="L13" class="ln">    13&nbsp;&nbsp;</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>Long Division
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>Division in a computer proceeds the same as long division in elementary school,
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>but computers are not as good as schoolchildren at following vague directions,
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>so we have to be much more precise about the actual steps and what can happen.
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>We work from most to least significant digit of the quotient, doing:
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span> • Guess a digit q, the number of v to subtract from the current
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>   section of u to zero out the topmost digit.
<span id="L25" class="ln">    25&nbsp;&nbsp;</span> • Check the guess by multiplying q·v and comparing it against
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>   the current section of u, adjusting the guess as needed.
<span id="L27" class="ln">    27&nbsp;&nbsp;</span> • Subtract q·v from the current section of u.
<span id="L28" class="ln">    28&nbsp;&nbsp;</span> • Add q to the corresponding section of the result quo.
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>When all digits have been processed, the final remainder is left in u
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>and returned as rem.
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>For example, here is a sketch of dividing 5 digits by 3 digits (n=3, m=2).
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>	                 q₂ q₁ q₀
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>	         _________________
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	v₂ v₁ v₀ ) u₄ u₃ u₂ u₁ u₀
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>	           ↓  ↓  ↓  |  |
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	          [u₄ u₃ u₂]|  |
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	        - [  q₂·v  ]|  |
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>	        ----------- ↓  |
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>	          [  rem  | u₁]|
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>	        - [    q₁·v   ]|
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>	           ----------- ↓
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	             [  rem  | u₀]
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>	           - [    q₀·v   ]
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>	              ------------
<span id="L48" class="ln">    48&nbsp;&nbsp;</span>	                [  rem   ]
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>Instead of creating new storage for the remainders and copying digits from u
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>as indicated by the arrows, we use u&#39;s storage directly as both the source
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>and destination of the subtractions, so that the remainders overwrite
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>successive overlapping sections of u as the division proceeds, using a slice
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>of u to identify the current section. This avoids all the copying as well as
<span id="L55" class="ln">    55&nbsp;&nbsp;</span>shifting of remainders.
<span id="L56" class="ln">    56&nbsp;&nbsp;</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>Division of u with n+m digits by v with n digits (in base B) can in general
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>produce at most m+1 digits, because:
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>  • u &lt; B^(n+m)               [B^(n+m) has n+m+1 digits]
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>  • v ≥ B^(n-1)               [B^(n-1) is the smallest n-digit number]
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>  • u/v &lt; B^(n+m) / B^(n-1)   [divide bounds for u, v]
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>  • u/v &lt; B^(m+1)             [simplify]
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>The first step is special: it takes the top n digits of u and divides them by
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>the n digits of v, producing the first quotient digit and an n-digit remainder.
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>In the example, q₂ = ⌊u₄u₃u₂ / v⌋.
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>The first step divides n digits by n digits to ensure that it produces only a
<span id="L70" class="ln">    70&nbsp;&nbsp;</span>single digit.
<span id="L71" class="ln">    71&nbsp;&nbsp;</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>Each subsequent step appends the next digit from u to the remainder and divides
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>those n+1 digits by the n digits of v, producing another quotient digit and a
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>new n-digit remainder.
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>Subsequent steps divide n+1 digits by n digits, an operation that in general
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>might produce two digits. However, as used in the algorithm, that division is
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>guaranteed to produce only a single digit. The dividend is of the form
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>rem·B + d, where rem is a remainder from the previous step and d is a single
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>digit, so:
<span id="L81" class="ln">    81&nbsp;&nbsp;</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span> • rem ≤ v - 1                 [rem is a remainder from dividing by v]
<span id="L83" class="ln">    83&nbsp;&nbsp;</span> • rem·B ≤ v·B - B             [multiply by B]
<span id="L84" class="ln">    84&nbsp;&nbsp;</span> • d ≤ B - 1                   [d is a single digit]
<span id="L85" class="ln">    85&nbsp;&nbsp;</span> • rem·B + d ≤ v·B - 1         [add]
<span id="L86" class="ln">    86&nbsp;&nbsp;</span> • rem·B + d &lt; v·B             [change ≤ to &lt;]
<span id="L87" class="ln">    87&nbsp;&nbsp;</span> • (rem·B + d)/v &lt; B           [divide by v]
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>Guess and Check
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span>At each step we need to divide n+1 digits by n digits, but this is for the
<span id="L93" class="ln">    93&nbsp;&nbsp;</span>implementation of division by n digits, so we can&#39;t just invoke a division
<span id="L94" class="ln">    94&nbsp;&nbsp;</span>routine: we _are_ the division routine. Instead, we guess at the answer and
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>then check it using multiplication. If the guess is wrong, we correct it.
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>How can this guessing possibly be efficient? It turns out that the following
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>statement (let&#39;s call it the Good Guess Guarantee) is true.
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>If
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>
<span id="L102" class="ln">   102&nbsp;&nbsp;</span> • q = ⌊u/v⌋ where u is n+1 digits and v is n digits,
<span id="L103" class="ln">   103&nbsp;&nbsp;</span> • q &lt; B, and
<span id="L104" class="ln">   104&nbsp;&nbsp;</span> • the topmost digit of v = vₙ₋₁ ≥ B/2,
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>then q̂ = ⌊uₙuₙ₋₁ / vₙ₋₁⌋ satisfies q ≤ q̂ ≤ q+2. (Proof below.)
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>That is, if we know the answer has only a single digit and we guess an answer
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>by ignoring the bottom n-1 digits of u and v, using a 2-by-1-digit division,
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>then that guess is at least as large as the correct answer. It is also not
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>too much larger: it is off by at most two from the correct answer.
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>Note that in the first step of the overall division, which is an n-by-n-digit
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>division, the 2-by-1 guess uses an implicit uₙ = 0.
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>Note that using a 2-by-1-digit division here does not mean calling ourselves
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>recursively. Instead, we use an efficient direct hardware implementation of
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>that operation.
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>Note that because q is u/v rounded down, q·v must not exceed u: u ≥ q·v.
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>If a guess q̂ is too big, it will not satisfy this test. Viewed a different way,
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>the remainder r̂ for a given q̂ is u - q̂·v, which must be positive. If it is
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>negative, then the guess q̂ is too big.
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>This gives us a way to compute q. First compute q̂ with 2-by-1-digit division.
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>Then, while u &lt; q̂·v, decrement q̂; this loop executes at most twice, because
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>q̂ ≤ q+2.
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>Scaling Inputs
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>The Good Guess Guarantee requires that the top digit of v (vₙ₋₁) be at least B/2.
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>For example in base 10, ⌊172/19⌋ = 9, but ⌊18/1⌋ = 18: the guess is wildly off
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>because the first digit 1 is smaller than B/2 = 5.
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>We can ensure that v has a large top digit by multiplying both u and v by the
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>right amount. Continuing the example, if we multiply both 172 and 19 by 3, we
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>now have ⌊516/57⌋, the leading digit of v is now ≥ 5, and sure enough
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>⌊51/5⌋ = 10 is much closer to the correct answer 9. It would be easier here
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>to multiply by 4, because that can be done with a shift. Specifically, we can
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>always count the number of leading zeros i in the first digit of v and then
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>shift both u and v left by i bits.
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>Having scaled u and v, the value ⌊u/v⌋ is unchanged, but the remainder will
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>be scaled: 172 mod 19 is 1, but 516 mod 57 is 3. We have to divide the remainder
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>by the scaling factor (shifting right i bits) when we finish.
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>Note that these shifts happen before and after the entire division algorithm,
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>not at each step in the per-digit iteration.
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>Note the effect of scaling inputs on the size of the possible quotient.
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>In the scaled u/v, u can gain a digit from scaling; v never does, because we
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>pick the scaling factor to make v&#39;s top digit larger but without overflowing.
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>If u and v have n+m and n digits after scaling, then:
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>  • u &lt; B^(n+m)               [B^(n+m) has n+m+1 digits]
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>  • v ≥ B^n / 2               [vₙ₋₁ ≥ B/2, so vₙ₋₁·B^(n-1) ≥ B^n/2]
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>  • u/v &lt; B^(n+m) / (B^n / 2) [divide bounds for u, v]
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>  • u/v &lt; 2 B^m               [simplify]
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>The quotient can still have m+1 significant digits, but if so the top digit
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>must be a 1. This provides a different way to handle the first digit of the
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>result: compare the top n digits of u against v and fill in either a 0 or a 1.
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>Refining Guesses
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>Before we check whether u &lt; q̂·v, we can adjust our guess to change it from
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>q̂ = ⌊uₙuₙ₋₁ / vₙ₋₁⌋ into the refined guess ⌊uₙuₙ₋₁uₙ₋₂ / vₙ₋₁vₙ₋₂⌋.
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>Although not mentioned above, the Good Guess Guarantee also promises that this
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>3-by-2-digit division guess is more precise and at most one away from the real
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>answer q. The improvement from the 2-by-1 to the 3-by-2 guess can also be done
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>without n-digit math.
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>If we have a guess q̂ = ⌊uₙuₙ₋₁ / vₙ₋₁⌋ and we want to see if it also equal to
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>⌊uₙuₙ₋₁uₙ₋₂ / vₙ₋₁vₙ₋₂⌋, we can use the same check we would for the full division:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>if uₙuₙ₋₁uₙ₋₂ &lt; q̂·vₙ₋₁vₙ₋₂, then the guess is too large and should be reduced.
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>Checking uₙuₙ₋₁uₙ₋₂ &lt; q̂·vₙ₋₁vₙ₋₂ is the same as uₙuₙ₋₁uₙ₋₂ - q̂·vₙ₋₁vₙ₋₂ &lt; 0,
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>and
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>	uₙuₙ₋₁uₙ₋₂ - q̂·vₙ₋₁vₙ₋₂ = (uₙuₙ₋₁·B + uₙ₋₂) - q̂·(vₙ₋₁·B + vₙ₋₂)
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>	                          [splitting off the bottom digit]
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>	                      = (uₙuₙ₋₁ - q̂·vₙ₋₁)·B + uₙ₋₂ - q̂·vₙ₋₂
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>	                          [regrouping]
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>The expression (uₙuₙ₋₁ - q̂·vₙ₋₁) is the remainder of uₙuₙ₋₁ / vₙ₋₁.
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>If the initial guess returns both q̂ and its remainder r̂, then checking
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>whether uₙuₙ₋₁uₙ₋₂ &lt; q̂·vₙ₋₁vₙ₋₂ is the same as checking r̂·B + uₙ₋₂ &lt; q̂·vₙ₋₂.
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>If we find that r̂·B + uₙ₋₂ &lt; q̂·vₙ₋₂, then we can adjust the guess by
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>decrementing q̂ and adding vₙ₋₁ to r̂. We repeat until r̂·B + uₙ₋₂ ≥ q̂·vₙ₋₂.
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>(As before, this fixup is only needed at most twice.)
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>Now that q̂ = ⌊uₙuₙ₋₁uₙ₋₂ / vₙ₋₁vₙ₋₂⌋, as mentioned above it is at most one
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>away from the correct q, and we&#39;ve avoided doing any n-digit math.
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>(If we need the new remainder, it can be computed as r̂·B + uₙ₋₂ - q̂·vₙ₋₂.)
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>The final check u &lt; q̂·v and the possible fixup must be done at full precision.
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>For random inputs, a fixup at this step is exceedingly rare: the 3-by-2 guess
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>is not often wrong at all. But still we must do the check. Note that since the
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>3-by-2 guess is off by at most 1, it can be convenient to perform the final
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>u &lt; q̂·v as part of the computation of the remainder r = u - q̂·v. If the
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>subtraction underflows, decremeting q̂ and adding one v back to r is enough to
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>arrive at the final q, r.
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>That&#39;s the entirety of long division: scale the inputs, and then loop over
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>each output position, guessing, checking, and correcting the next output digit.
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>For a 2n-digit number divided by an n-digit number (the worst size-n case for
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>division complexity), this algorithm uses n+1 iterations, each of which must do
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>at least the 1-by-n-digit multiplication q̂·v. That&#39;s O(n) iterations of
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>O(n) time each, so O(n²) time overall.
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>Recursive Division
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>For very large inputs, it is possible to improve on the O(n²) algorithm.
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>Let&#39;s call a group of n/2 real digits a (very) “wide digit”. We can run the
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>standard long division algorithm explained above over the wide digits instead of
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>the actual digits. This will result in many fewer steps, but the math involved in
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>each step is more work.
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>Where basic long division uses a 2-by-1-digit division to guess the initial q̂,
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>the new algorithm must use a 2-by-1-wide-digit division, which is of course
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>really an n-by-n/2-digit division. That&#39;s OK: if we implement n-digit division
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>in terms of n/2-digit division, the recursion will terminate when the divisor
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>becomes small enough to handle with standard long division or even with the
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>2-by-1 hardware instruction.
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>For example, here is a sketch of dividing 10 digits by 4, proceeding with
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>wide digits corresponding to two regular digits. The first step, still special,
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>must leave off a (regular) digit, dividing 5 by 4 and producing a 4-digit
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>remainder less than v. The middle steps divide 6 digits by 4, guaranteed to
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>produce two output digits each (one wide digit) with 4-digit remainders.
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>The final step must use what it has: the 4-digit remainder plus one more,
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>5 digits to divide by 4.
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>	                       q₆ q₅ q₄ q₃ q₂ q₁ q₀
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>	            _______________________________
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>	v₃ v₂ v₁ v₀ ) u₉ u₈ u₇ u₆ u₅ u₄ u₃ u₂ u₁ u₀
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>	              ↓  ↓  ↓  ↓  ↓  |  |  |  |  |
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>	             [u₉ u₈ u₇ u₆ u₅]|  |  |  |  |
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>	           - [    q₆q₅·v    ]|  |  |  |  |
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>	           ----------------- ↓  ↓  |  |  |
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>	                [    rem    |u₄ u₃]|  |  |
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>	              - [     q₄q₃·v      ]|  |  |
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>	              -------------------- ↓  ↓  |
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>	                      [    rem    |u₂ u₁]|
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>	                    - [     q₂q₁·v      ]|
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>	                    -------------------- ↓
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>	                            [    rem    |u₀]
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>	                          - [     q₀·v     ]
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>	                          ------------------
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>	                               [    rem    ]
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>An alternative would be to look ahead to how well n/2 divides into n+m and
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>adjust the first step to use fewer digits as needed, making the first step
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>more special to make the last step not special at all. For example, using the
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>same input, we could choose to use only 4 digits in the first step, leaving
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>a full wide digit for the last step:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>	                       q₆ q₅ q₄ q₃ q₂ q₁ q₀
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>	            _______________________________
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>	v₃ v₂ v₁ v₀ ) u₉ u₈ u₇ u₆ u₅ u₄ u₃ u₂ u₁ u₀
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>	              ↓  ↓  ↓  ↓  |  |  |  |  |  |
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	             [u₉ u₈ u₇ u₆]|  |  |  |  |  |
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>	           - [    q₆·v   ]|  |  |  |  |  |
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	           -------------- ↓  ↓  |  |  |  |
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>	             [    rem    |u₅ u₄]|  |  |  |
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	           - [     q₅q₄·v      ]|  |  |  |
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>	           -------------------- ↓  ↓  |  |
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	                   [    rem    |u₃ u₂]|  |
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>	                 - [     q₃q₂·v      ]|  |
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>	                 -------------------- ↓  ↓
<span id="L276" class="ln">   276&nbsp;&nbsp;</span>	                         [    rem    |u₁ u₀]
<span id="L277" class="ln">   277&nbsp;&nbsp;</span>	                       - [     q₁q₀·v      ]
<span id="L278" class="ln">   278&nbsp;&nbsp;</span>	                       ---------------------
<span id="L279" class="ln">   279&nbsp;&nbsp;</span>	                               [    rem    ]
<span id="L280" class="ln">   280&nbsp;&nbsp;</span>
<span id="L281" class="ln">   281&nbsp;&nbsp;</span>Today, the code in divRecursiveStep works like the first example. Perhaps in
<span id="L282" class="ln">   282&nbsp;&nbsp;</span>the future we will make it work like the alternative, to avoid a special case
<span id="L283" class="ln">   283&nbsp;&nbsp;</span>in the final iteration.
<span id="L284" class="ln">   284&nbsp;&nbsp;</span>
<span id="L285" class="ln">   285&nbsp;&nbsp;</span>Either way, each step is a 3-by-2-wide-digit division approximated first by
<span id="L286" class="ln">   286&nbsp;&nbsp;</span>a 2-by-1-wide-digit division, just as we did for regular digits in long division.
<span id="L287" class="ln">   287&nbsp;&nbsp;</span>Because the actual answer we want is a 3-by-2-wide-digit division, instead of
<span id="L288" class="ln">   288&nbsp;&nbsp;</span>multiplying q̂·v directly during the fixup, we can use the quick refinement
<span id="L289" class="ln">   289&nbsp;&nbsp;</span>from long division (an n/2-by-n/2 multiply) to correct q to its actual value
<span id="L290" class="ln">   290&nbsp;&nbsp;</span>and also compute the remainder (as mentioned above), and then stop after that,
<span id="L291" class="ln">   291&nbsp;&nbsp;</span>never doing a full n-by-n multiply.
<span id="L292" class="ln">   292&nbsp;&nbsp;</span>
<span id="L293" class="ln">   293&nbsp;&nbsp;</span>Instead of using an n-by-n/2-digit division to produce n/2 digits, we can add
<span id="L294" class="ln">   294&nbsp;&nbsp;</span>(not discard) one more real digit, doing an (n+1)-by-(n/2+1)-digit division that
<span id="L295" class="ln">   295&nbsp;&nbsp;</span>produces n/2+1 digits. That single extra digit tightens the Good Guess Guarantee
<span id="L296" class="ln">   296&nbsp;&nbsp;</span>to q ≤ q̂ ≤ q+1 and lets us drop long division&#39;s special treatment of the first
<span id="L297" class="ln">   297&nbsp;&nbsp;</span>digit. These benefits are discussed more after the Good Guess Guarantee proof
<span id="L298" class="ln">   298&nbsp;&nbsp;</span>below.
<span id="L299" class="ln">   299&nbsp;&nbsp;</span>
<span id="L300" class="ln">   300&nbsp;&nbsp;</span>
<span id="L301" class="ln">   301&nbsp;&nbsp;</span>How Fast is Recursive Division?
<span id="L302" class="ln">   302&nbsp;&nbsp;</span>
<span id="L303" class="ln">   303&nbsp;&nbsp;</span>For a 2n-by-n-digit division, this algorithm runs a 4-by-2 long division over
<span id="L304" class="ln">   304&nbsp;&nbsp;</span>wide digits, producing two wide digits plus a possible leading regular digit 1,
<span id="L305" class="ln">   305&nbsp;&nbsp;</span>which can be handled without a recursive call. That is, the algorithm uses two
<span id="L306" class="ln">   306&nbsp;&nbsp;</span>full iterations, each using an n-by-n/2-digit division and an n/2-by-n/2-digit
<span id="L307" class="ln">   307&nbsp;&nbsp;</span>multiplication, along with a few n-digit additions and subtractions. The standard
<span id="L308" class="ln">   308&nbsp;&nbsp;</span>n-by-n-digit multiplication algorithm requires O(n²) time, making the overall
<span id="L309" class="ln">   309&nbsp;&nbsp;</span>algorithm require time T(n) where
<span id="L310" class="ln">   310&nbsp;&nbsp;</span>
<span id="L311" class="ln">   311&nbsp;&nbsp;</span>	T(n) = 2T(n/2) + O(n) + O(n²)
<span id="L312" class="ln">   312&nbsp;&nbsp;</span>
<span id="L313" class="ln">   313&nbsp;&nbsp;</span>which, by the Bentley-Haken-Saxe theorem, ends up reducing to T(n) = O(n²).
<span id="L314" class="ln">   314&nbsp;&nbsp;</span>This is not an improvement over regular long division.
<span id="L315" class="ln">   315&nbsp;&nbsp;</span>
<span id="L316" class="ln">   316&nbsp;&nbsp;</span>When the number of digits n becomes large enough, Karatsuba&#39;s algorithm for
<span id="L317" class="ln">   317&nbsp;&nbsp;</span>multiplication can be used instead, which takes O(n^log₂3) = O(n^1.6) time.
<span id="L318" class="ln">   318&nbsp;&nbsp;</span>(Karatsuba multiplication is implemented in func karatsuba in nat.go.)
<span id="L319" class="ln">   319&nbsp;&nbsp;</span>That makes the overall recursive division algorithm take O(n^1.6) time as well,
<span id="L320" class="ln">   320&nbsp;&nbsp;</span>which is an improvement, but again only for large enough numbers.
<span id="L321" class="ln">   321&nbsp;&nbsp;</span>
<span id="L322" class="ln">   322&nbsp;&nbsp;</span>It is not critical to make sure that every recursion does only two recursive
<span id="L323" class="ln">   323&nbsp;&nbsp;</span>calls. While in general the number of recursive calls can change the time
<span id="L324" class="ln">   324&nbsp;&nbsp;</span>analysis, in this case doing three calls does not change the analysis:
<span id="L325" class="ln">   325&nbsp;&nbsp;</span>
<span id="L326" class="ln">   326&nbsp;&nbsp;</span>	T(n) = 3T(n/2) + O(n) + O(n^log₂3)
<span id="L327" class="ln">   327&nbsp;&nbsp;</span>
<span id="L328" class="ln">   328&nbsp;&nbsp;</span>ends up being T(n) = O(n^log₂3). Because the Karatsuba multiplication taking
<span id="L329" class="ln">   329&nbsp;&nbsp;</span>time O(n^log₂3) is itself doing 3 half-sized recursions, doing three for the
<span id="L330" class="ln">   330&nbsp;&nbsp;</span>division does not hurt the asymptotic performance. Of course, it is likely
<span id="L331" class="ln">   331&nbsp;&nbsp;</span>still faster in practice to do two.
<span id="L332" class="ln">   332&nbsp;&nbsp;</span>
<span id="L333" class="ln">   333&nbsp;&nbsp;</span>
<span id="L334" class="ln">   334&nbsp;&nbsp;</span>Proof of the Good Guess Guarantee
<span id="L335" class="ln">   335&nbsp;&nbsp;</span>
<span id="L336" class="ln">   336&nbsp;&nbsp;</span>Given numbers x, y, let us break them into the quotients and remainders when
<span id="L337" class="ln">   337&nbsp;&nbsp;</span>divided by some scaling factor S, with the added constraints that the quotient
<span id="L338" class="ln">   338&nbsp;&nbsp;</span>x/y and the high part of y are both less than some limit T, and that the high
<span id="L339" class="ln">   339&nbsp;&nbsp;</span>part of y is at least half as big as T.
<span id="L340" class="ln">   340&nbsp;&nbsp;</span>
<span id="L341" class="ln">   341&nbsp;&nbsp;</span>	x₁ = ⌊x/S⌋        y₁ = ⌊y/S⌋
<span id="L342" class="ln">   342&nbsp;&nbsp;</span>	x₀ = x mod S      y₀ = y mod S
<span id="L343" class="ln">   343&nbsp;&nbsp;</span>
<span id="L344" class="ln">   344&nbsp;&nbsp;</span>	x  = x₁·S + x₀    0 ≤ x₀ &lt; S    x/y &lt; T
<span id="L345" class="ln">   345&nbsp;&nbsp;</span>	y  = y₁·S + y₀    0 ≤ y₀ &lt; S    T/2 ≤ y₁ &lt; T
<span id="L346" class="ln">   346&nbsp;&nbsp;</span>
<span id="L347" class="ln">   347&nbsp;&nbsp;</span>And consider the two truncated quotients:
<span id="L348" class="ln">   348&nbsp;&nbsp;</span>
<span id="L349" class="ln">   349&nbsp;&nbsp;</span>	q = ⌊x/y⌋
<span id="L350" class="ln">   350&nbsp;&nbsp;</span>	q̂ = ⌊x₁/y₁⌋
<span id="L351" class="ln">   351&nbsp;&nbsp;</span>
<span id="L352" class="ln">   352&nbsp;&nbsp;</span>We will prove that q ≤ q̂ ≤ q+2.
<span id="L353" class="ln">   353&nbsp;&nbsp;</span>
<span id="L354" class="ln">   354&nbsp;&nbsp;</span>The guarantee makes no real demands on the scaling factor S: it is simply the
<span id="L355" class="ln">   355&nbsp;&nbsp;</span>magnitude of the digits cut from both x and y to produce x₁ and y₁.
<span id="L356" class="ln">   356&nbsp;&nbsp;</span>The guarantee makes only limited demands on T: it must be large enough to hold
<span id="L357" class="ln">   357&nbsp;&nbsp;</span>the quotient x/y, and y₁ must have roughly the same size.
<span id="L358" class="ln">   358&nbsp;&nbsp;</span>
<span id="L359" class="ln">   359&nbsp;&nbsp;</span>To apply to the earlier discussion of 2-by-1 guesses in long division,
<span id="L360" class="ln">   360&nbsp;&nbsp;</span>we would choose:
<span id="L361" class="ln">   361&nbsp;&nbsp;</span>
<span id="L362" class="ln">   362&nbsp;&nbsp;</span>	S  = Bⁿ⁻¹
<span id="L363" class="ln">   363&nbsp;&nbsp;</span>	T  = B
<span id="L364" class="ln">   364&nbsp;&nbsp;</span>	x  = u
<span id="L365" class="ln">   365&nbsp;&nbsp;</span>	x₁ = uₙuₙ₋₁
<span id="L366" class="ln">   366&nbsp;&nbsp;</span>	x₀ = uₙ₋₂...u₀
<span id="L367" class="ln">   367&nbsp;&nbsp;</span>	y  = v
<span id="L368" class="ln">   368&nbsp;&nbsp;</span>	y₁ = vₙ₋₁
<span id="L369" class="ln">   369&nbsp;&nbsp;</span>	y₀ = vₙ₋₂...u₀
<span id="L370" class="ln">   370&nbsp;&nbsp;</span>
<span id="L371" class="ln">   371&nbsp;&nbsp;</span>These simpler variables avoid repeating those longer expressions in the proof.
<span id="L372" class="ln">   372&nbsp;&nbsp;</span>
<span id="L373" class="ln">   373&nbsp;&nbsp;</span>Note also that, by definition, truncating division ⌊x/y⌋ satisfies
<span id="L374" class="ln">   374&nbsp;&nbsp;</span>
<span id="L375" class="ln">   375&nbsp;&nbsp;</span>	x/y - 1 &lt; ⌊x/y⌋ ≤ x/y.
<span id="L376" class="ln">   376&nbsp;&nbsp;</span>
<span id="L377" class="ln">   377&nbsp;&nbsp;</span>This fact will be used a few times in the proofs.
<span id="L378" class="ln">   378&nbsp;&nbsp;</span>
<span id="L379" class="ln">   379&nbsp;&nbsp;</span>Proof that q ≤ q̂:
<span id="L380" class="ln">   380&nbsp;&nbsp;</span>
<span id="L381" class="ln">   381&nbsp;&nbsp;</span>	q̂·y₁ = ⌊x₁/y₁⌋·y₁                      [by definition, q̂ = ⌊x₁/y₁⌋]
<span id="L382" class="ln">   382&nbsp;&nbsp;</span>	     &gt; (x₁/y₁ - 1)·y₁                  [x₁/y₁ - 1 &lt; ⌊x₁/y₁⌋]
<span id="L383" class="ln">   383&nbsp;&nbsp;</span>	     = x₁ - y₁                         [distribute y₁]
<span id="L384" class="ln">   384&nbsp;&nbsp;</span>
<span id="L385" class="ln">   385&nbsp;&nbsp;</span>	So q̂·y₁ &gt; x₁ - y₁.
<span id="L386" class="ln">   386&nbsp;&nbsp;</span>	Since q̂·y₁ is an integer, q̂·y₁ ≥ x₁ - y₁ + 1.
<span id="L387" class="ln">   387&nbsp;&nbsp;</span>
<span id="L388" class="ln">   388&nbsp;&nbsp;</span>	q̂ - q = q̂ - ⌊x/y⌋                      [by definition, q = ⌊x/y⌋]
<span id="L389" class="ln">   389&nbsp;&nbsp;</span>	      ≥ q̂ - x/y                        [⌊x/y⌋ &lt; x/y]
<span id="L390" class="ln">   390&nbsp;&nbsp;</span>	      = (1/y)·(q̂·y - x)                [factor out 1/y]
<span id="L391" class="ln">   391&nbsp;&nbsp;</span>	      ≥ (1/y)·(q̂·y₁·S - x)             [y = y₁·S + y₀ ≥ y₁·S]
<span id="L392" class="ln">   392&nbsp;&nbsp;</span>	      ≥ (1/y)·((x₁ - y₁ + 1)·S - x)    [above: q̂·y₁ ≥ x₁ - y₁ + 1]
<span id="L393" class="ln">   393&nbsp;&nbsp;</span>	      = (1/y)·(x₁·S - y₁·S + S - x)    [distribute S]
<span id="L394" class="ln">   394&nbsp;&nbsp;</span>	      = (1/y)·(S - x₀ - y₁·S)          [-x = -x₁·S - x₀]
<span id="L395" class="ln">   395&nbsp;&nbsp;</span>	      &gt; -y₁·S / y                      [x₀ &lt; S, so S - x₀ &lt; 0; drop it]
<span id="L396" class="ln">   396&nbsp;&nbsp;</span>	      ≥ -1                             [y₁·S ≤ y]
<span id="L397" class="ln">   397&nbsp;&nbsp;</span>
<span id="L398" class="ln">   398&nbsp;&nbsp;</span>	So q̂ - q &gt; -1.
<span id="L399" class="ln">   399&nbsp;&nbsp;</span>	Since q̂ - q is an integer, q̂ - q ≥ 0, or equivalently q ≤ q̂.
<span id="L400" class="ln">   400&nbsp;&nbsp;</span>
<span id="L401" class="ln">   401&nbsp;&nbsp;</span>Proof that q̂ ≤ q+2:
<span id="L402" class="ln">   402&nbsp;&nbsp;</span>
<span id="L403" class="ln">   403&nbsp;&nbsp;</span>	x₁/y₁ - x/y = x₁·S/y₁·S - x/y          [multiply left term by S/S]
<span id="L404" class="ln">   404&nbsp;&nbsp;</span>	            ≤ x/y₁·S - x/y             [x₁S ≤ x]
<span id="L405" class="ln">   405&nbsp;&nbsp;</span>	            = (x/y)·(y/y₁·S - 1)       [factor out x/y]
<span id="L406" class="ln">   406&nbsp;&nbsp;</span>	            = (x/y)·((y - y₁·S)/y₁·S)  [move -1 into y/y₁·S fraction]
<span id="L407" class="ln">   407&nbsp;&nbsp;</span>	            = (x/y)·(y₀/y₁·S)          [y - y₁·S = y₀]
<span id="L408" class="ln">   408&nbsp;&nbsp;</span>	            = (x/y)·(1/y₁)·(y₀/S)      [factor out 1/y₁]
<span id="L409" class="ln">   409&nbsp;&nbsp;</span>	            &lt; (x/y)·(1/y₁)             [y₀ &lt; S, so y₀/S &lt; 1]
<span id="L410" class="ln">   410&nbsp;&nbsp;</span>	            ≤ (x/y)·(2/T)              [y₁ ≥ T/2, so 1/y₁ ≤ 2/T]
<span id="L411" class="ln">   411&nbsp;&nbsp;</span>	            &lt; T·(2/T)                  [x/y &lt; T]
<span id="L412" class="ln">   412&nbsp;&nbsp;</span>	            = 2                        [T·(2/T) = 2]
<span id="L413" class="ln">   413&nbsp;&nbsp;</span>
<span id="L414" class="ln">   414&nbsp;&nbsp;</span>	So x₁/y₁ - x/y &lt; 2.
<span id="L415" class="ln">   415&nbsp;&nbsp;</span>
<span id="L416" class="ln">   416&nbsp;&nbsp;</span>	q̂ - q = ⌊x₁/y₁⌋ - q                    [by definition, q̂ = ⌊x₁/y₁⌋]
<span id="L417" class="ln">   417&nbsp;&nbsp;</span>	      = ⌊x₁/y₁⌋ - ⌊x/y⌋                [by definition, q = ⌊x/y⌋]
<span id="L418" class="ln">   418&nbsp;&nbsp;</span>	      ≤ x₁/y₁ - ⌊x/y⌋                  [⌊x₁/y₁⌋ ≤ x₁/y₁]
<span id="L419" class="ln">   419&nbsp;&nbsp;</span>	      &lt; x₁/y₁ - (x/y - 1)              [⌊x/y⌋ &gt; x/y - 1]
<span id="L420" class="ln">   420&nbsp;&nbsp;</span>	      = (x₁/y₁ - x/y) + 1              [regrouping]
<span id="L421" class="ln">   421&nbsp;&nbsp;</span>	      &lt; 2 + 1                          [above: x₁/y₁ - x/y &lt; 2]
<span id="L422" class="ln">   422&nbsp;&nbsp;</span>	      = 3
<span id="L423" class="ln">   423&nbsp;&nbsp;</span>
<span id="L424" class="ln">   424&nbsp;&nbsp;</span>	So q̂ - q &lt; 3.
<span id="L425" class="ln">   425&nbsp;&nbsp;</span>	Since q̂ - q is an integer, q̂ - q ≤ 2.
<span id="L426" class="ln">   426&nbsp;&nbsp;</span>
<span id="L427" class="ln">   427&nbsp;&nbsp;</span>Note that when x/y &lt; T/2, the bounds tighten to x₁/y₁ - x/y &lt; 1 and therefore
<span id="L428" class="ln">   428&nbsp;&nbsp;</span>q̂ - q ≤ 1.
<span id="L429" class="ln">   429&nbsp;&nbsp;</span>
<span id="L430" class="ln">   430&nbsp;&nbsp;</span>Note also that in the general case 2n-by-n division where we don&#39;t know that
<span id="L431" class="ln">   431&nbsp;&nbsp;</span>x/y &lt; T, we do know that x/y &lt; 2T, yielding the bound q̂ - q ≤ 4. So we could
<span id="L432" class="ln">   432&nbsp;&nbsp;</span>remove the special case first step of long division as long as we allow the
<span id="L433" class="ln">   433&nbsp;&nbsp;</span>first fixup loop to run up to four times. (Using a simple comparison to decide
<span id="L434" class="ln">   434&nbsp;&nbsp;</span>whether the first digit is 0 or 1 is still more efficient, though.)
<span id="L435" class="ln">   435&nbsp;&nbsp;</span>
<span id="L436" class="ln">   436&nbsp;&nbsp;</span>Finally, note that when dividing three leading base-B digits by two (scaled),
<span id="L437" class="ln">   437&nbsp;&nbsp;</span>we have T = B² and x/y &lt; B = T/B, a much tighter bound than x/y &lt; T.
<span id="L438" class="ln">   438&nbsp;&nbsp;</span>This in turn yields the much tighter bound x₁/y₁ - x/y &lt; 2/B. This means that
<span id="L439" class="ln">   439&nbsp;&nbsp;</span>⌊x₁/y₁⌋ and ⌊x/y⌋ can only differ when x/y is less than 2/B greater than an
<span id="L440" class="ln">   440&nbsp;&nbsp;</span>integer. For random x and y, the chance of this is 2/B, or, for large B,
<span id="L441" class="ln">   441&nbsp;&nbsp;</span>approximately zero. This means that after we produce the 3-by-2 guess in the
<span id="L442" class="ln">   442&nbsp;&nbsp;</span>long division algorithm, the fixup loop essentially never runs.
<span id="L443" class="ln">   443&nbsp;&nbsp;</span>
<span id="L444" class="ln">   444&nbsp;&nbsp;</span>In the recursive algorithm, the extra digit in (2·⌊n/2⌋+1)-by-(⌊n/2⌋+1)-digit
<span id="L445" class="ln">   445&nbsp;&nbsp;</span>division has exactly the same effect: the probability of needing a fixup is the
<span id="L446" class="ln">   446&nbsp;&nbsp;</span>same 2/B. Even better, we can allow the general case x/y &lt; 2T and the fixup
<span id="L447" class="ln">   447&nbsp;&nbsp;</span>probability only grows to 4/B, still essentially zero.
<span id="L448" class="ln">   448&nbsp;&nbsp;</span>
<span id="L449" class="ln">   449&nbsp;&nbsp;</span>
<span id="L450" class="ln">   450&nbsp;&nbsp;</span>References
<span id="L451" class="ln">   451&nbsp;&nbsp;</span>
<span id="L452" class="ln">   452&nbsp;&nbsp;</span>There are no great references for implementing long division; thus this comment.
<span id="L453" class="ln">   453&nbsp;&nbsp;</span>Here are some notes about what to expect from the obvious references.
<span id="L454" class="ln">   454&nbsp;&nbsp;</span>
<span id="L455" class="ln">   455&nbsp;&nbsp;</span>Knuth Volume 2 (Seminumerical Algorithms) section 4.3.1 is the usual canonical
<span id="L456" class="ln">   456&nbsp;&nbsp;</span>reference for long division, but that entire series is highly compressed, never
<span id="L457" class="ln">   457&nbsp;&nbsp;</span>repeating a necessary fact and leaving important insights to the exercises.
<span id="L458" class="ln">   458&nbsp;&nbsp;</span>For example, no rationale whatsoever is given for the calculation that extends
<span id="L459" class="ln">   459&nbsp;&nbsp;</span>q̂ from a 2-by-1 to a 3-by-2 guess, nor why it reduces the error bound.
<span id="L460" class="ln">   460&nbsp;&nbsp;</span>The proof that the calculation even has the desired effect is left to exercises.
<span id="L461" class="ln">   461&nbsp;&nbsp;</span>The solutions to those exercises provided at the back of the book are entirely
<span id="L462" class="ln">   462&nbsp;&nbsp;</span>calculations, still with no explanation as to what is going on or how you would
<span id="L463" class="ln">   463&nbsp;&nbsp;</span>arrive at the idea of doing those exact calculations. Nowhere is it mentioned
<span id="L464" class="ln">   464&nbsp;&nbsp;</span>that this test extends the 2-by-1 guess into a 3-by-2 guess. The proof of the
<span id="L465" class="ln">   465&nbsp;&nbsp;</span>Good Guess Guarantee is only for the 2-by-1 guess and argues by contradiction,
<span id="L466" class="ln">   466&nbsp;&nbsp;</span>making it difficult to understand how modifications like adding another digit
<span id="L467" class="ln">   467&nbsp;&nbsp;</span>or adjusting the quotient range affects the overall bound.
<span id="L468" class="ln">   468&nbsp;&nbsp;</span>
<span id="L469" class="ln">   469&nbsp;&nbsp;</span>All that said, Knuth remains the canonical reference. It is dense but packed
<span id="L470" class="ln">   470&nbsp;&nbsp;</span>full of information and references, and the proofs are simpler than many other
<span id="L471" class="ln">   471&nbsp;&nbsp;</span>presentations. The proofs above are reworkings of Knuth&#39;s to remove the
<span id="L472" class="ln">   472&nbsp;&nbsp;</span>arguments by contradiction and add explanations or steps that Knuth omitted.
<span id="L473" class="ln">   473&nbsp;&nbsp;</span>But beware of errors in older printings. Take the published errata with you.
<span id="L474" class="ln">   474&nbsp;&nbsp;</span>
<span id="L475" class="ln">   475&nbsp;&nbsp;</span>Brinch Hansen&#39;s “Multiple-length Division Revisited: a Tour of the Minefield”
<span id="L476" class="ln">   476&nbsp;&nbsp;</span>starts with a blunt critique of Knuth&#39;s presentation (among others) and then
<span id="L477" class="ln">   477&nbsp;&nbsp;</span>presents a more detailed and easier to follow treatment of long division,
<span id="L478" class="ln">   478&nbsp;&nbsp;</span>including an implementation in Pascal. But the algorithm and implementation
<span id="L479" class="ln">   479&nbsp;&nbsp;</span>work entirely in terms of 3-by-2 division, which is much less useful on modern
<span id="L480" class="ln">   480&nbsp;&nbsp;</span>hardware than an algorithm using 2-by-1 division. The proofs are a bit too
<span id="L481" class="ln">   481&nbsp;&nbsp;</span>focused on digit counting and seem needlessly complex, especially compared to
<span id="L482" class="ln">   482&nbsp;&nbsp;</span>the ones given above.
<span id="L483" class="ln">   483&nbsp;&nbsp;</span>
<span id="L484" class="ln">   484&nbsp;&nbsp;</span>Burnikel and Ziegler&#39;s “Fast Recursive Division” introduced the key insight of
<span id="L485" class="ln">   485&nbsp;&nbsp;</span>implementing division by an n-digit divisor using recursive calls to division
<span id="L486" class="ln">   486&nbsp;&nbsp;</span>by an n/2-digit divisor, relying on Karatsuba multiplication to yield a
<span id="L487" class="ln">   487&nbsp;&nbsp;</span>sub-quadratic run time. However, the presentation decisions are made almost
<span id="L488" class="ln">   488&nbsp;&nbsp;</span>entirely for the purpose of simplifying the run-time analysis, rather than
<span id="L489" class="ln">   489&nbsp;&nbsp;</span>simplifying the presentation. Instead of a single algorithm that loops over
<span id="L490" class="ln">   490&nbsp;&nbsp;</span>quotient digits, the paper presents two mutually-recursive algorithms, for
<span id="L491" class="ln">   491&nbsp;&nbsp;</span>2n-by-n and 3n-by-2n. The paper also does not present any general (n+m)-by-n
<span id="L492" class="ln">   492&nbsp;&nbsp;</span>algorithm.
<span id="L493" class="ln">   493&nbsp;&nbsp;</span>
<span id="L494" class="ln">   494&nbsp;&nbsp;</span>The proofs in the paper are remarkably complex, especially considering that
<span id="L495" class="ln">   495&nbsp;&nbsp;</span>the algorithm is at its core just long division on wide digits, so that the
<span id="L496" class="ln">   496&nbsp;&nbsp;</span>usual long division proofs apply essentially unaltered.
<span id="L497" class="ln">   497&nbsp;&nbsp;</span>*/</span>
<span id="L498" class="ln">   498&nbsp;&nbsp;</span>
<span id="L499" class="ln">   499&nbsp;&nbsp;</span>package big
<span id="L500" class="ln">   500&nbsp;&nbsp;</span>
<span id="L501" class="ln">   501&nbsp;&nbsp;</span>import &#34;math/bits&#34;
<span id="L502" class="ln">   502&nbsp;&nbsp;</span>
<span id="L503" class="ln">   503&nbsp;&nbsp;</span><span class="comment">// rem returns r such that r = u%v.</span>
<span id="L504" class="ln">   504&nbsp;&nbsp;</span><span class="comment">// It uses z as the storage for r.</span>
<span id="L505" class="ln">   505&nbsp;&nbsp;</span>func (z nat) rem(u, v nat) (r nat) {
<span id="L506" class="ln">   506&nbsp;&nbsp;</span>	if alias(z, u) {
<span id="L507" class="ln">   507&nbsp;&nbsp;</span>		z = nil
<span id="L508" class="ln">   508&nbsp;&nbsp;</span>	}
<span id="L509" class="ln">   509&nbsp;&nbsp;</span>	qp := getNat(0)
<span id="L510" class="ln">   510&nbsp;&nbsp;</span>	q, r := qp.div(z, u, v)
<span id="L511" class="ln">   511&nbsp;&nbsp;</span>	*qp = q
<span id="L512" class="ln">   512&nbsp;&nbsp;</span>	putNat(qp)
<span id="L513" class="ln">   513&nbsp;&nbsp;</span>	return r
<span id="L514" class="ln">   514&nbsp;&nbsp;</span>}
<span id="L515" class="ln">   515&nbsp;&nbsp;</span>
<span id="L516" class="ln">   516&nbsp;&nbsp;</span><span class="comment">// div returns q, r such that q = ⌊u/v⌋ and r = u%v = u - q·v.</span>
<span id="L517" class="ln">   517&nbsp;&nbsp;</span><span class="comment">// It uses z and z2 as the storage for q and r.</span>
<span id="L518" class="ln">   518&nbsp;&nbsp;</span>func (z nat) div(z2, u, v nat) (q, r nat) {
<span id="L519" class="ln">   519&nbsp;&nbsp;</span>	if len(v) == 0 {
<span id="L520" class="ln">   520&nbsp;&nbsp;</span>		panic(&#34;division by zero&#34;)
<span id="L521" class="ln">   521&nbsp;&nbsp;</span>	}
<span id="L522" class="ln">   522&nbsp;&nbsp;</span>
<span id="L523" class="ln">   523&nbsp;&nbsp;</span>	if u.cmp(v) &lt; 0 {
<span id="L524" class="ln">   524&nbsp;&nbsp;</span>		q = z[:0]
<span id="L525" class="ln">   525&nbsp;&nbsp;</span>		r = z2.set(u)
<span id="L526" class="ln">   526&nbsp;&nbsp;</span>		return
<span id="L527" class="ln">   527&nbsp;&nbsp;</span>	}
<span id="L528" class="ln">   528&nbsp;&nbsp;</span>
<span id="L529" class="ln">   529&nbsp;&nbsp;</span>	if len(v) == 1 {
<span id="L530" class="ln">   530&nbsp;&nbsp;</span>		<span class="comment">// Short division: long optimized for a single-word divisor.</span>
<span id="L531" class="ln">   531&nbsp;&nbsp;</span>		<span class="comment">// In that case, the 2-by-1 guess is all we need at each step.</span>
<span id="L532" class="ln">   532&nbsp;&nbsp;</span>		var r2 Word
<span id="L533" class="ln">   533&nbsp;&nbsp;</span>		q, r2 = z.divW(u, v[0])
<span id="L534" class="ln">   534&nbsp;&nbsp;</span>		r = z2.setWord(r2)
<span id="L535" class="ln">   535&nbsp;&nbsp;</span>		return
<span id="L536" class="ln">   536&nbsp;&nbsp;</span>	}
<span id="L537" class="ln">   537&nbsp;&nbsp;</span>
<span id="L538" class="ln">   538&nbsp;&nbsp;</span>	q, r = z.divLarge(z2, u, v)
<span id="L539" class="ln">   539&nbsp;&nbsp;</span>	return
<span id="L540" class="ln">   540&nbsp;&nbsp;</span>}
<span id="L541" class="ln">   541&nbsp;&nbsp;</span>
<span id="L542" class="ln">   542&nbsp;&nbsp;</span><span class="comment">// divW returns q, r such that q = ⌊x/y⌋ and r = x%y = x - q·y.</span>
<span id="L543" class="ln">   543&nbsp;&nbsp;</span><span class="comment">// It uses z as the storage for q.</span>
<span id="L544" class="ln">   544&nbsp;&nbsp;</span><span class="comment">// Note that y is a single digit (Word), not a big number.</span>
<span id="L545" class="ln">   545&nbsp;&nbsp;</span>func (z nat) divW(x nat, y Word) (q nat, r Word) {
<span id="L546" class="ln">   546&nbsp;&nbsp;</span>	m := len(x)
<span id="L547" class="ln">   547&nbsp;&nbsp;</span>	switch {
<span id="L548" class="ln">   548&nbsp;&nbsp;</span>	case y == 0:
<span id="L549" class="ln">   549&nbsp;&nbsp;</span>		panic(&#34;division by zero&#34;)
<span id="L550" class="ln">   550&nbsp;&nbsp;</span>	case y == 1:
<span id="L551" class="ln">   551&nbsp;&nbsp;</span>		q = z.set(x) <span class="comment">// result is x</span>
<span id="L552" class="ln">   552&nbsp;&nbsp;</span>		return
<span id="L553" class="ln">   553&nbsp;&nbsp;</span>	case m == 0:
<span id="L554" class="ln">   554&nbsp;&nbsp;</span>		q = z[:0] <span class="comment">// result is 0</span>
<span id="L555" class="ln">   555&nbsp;&nbsp;</span>		return
<span id="L556" class="ln">   556&nbsp;&nbsp;</span>	}
<span id="L557" class="ln">   557&nbsp;&nbsp;</span>	<span class="comment">// m &gt; 0</span>
<span id="L558" class="ln">   558&nbsp;&nbsp;</span>	z = z.make(m)
<span id="L559" class="ln">   559&nbsp;&nbsp;</span>	r = divWVW(z, 0, x, y)
<span id="L560" class="ln">   560&nbsp;&nbsp;</span>	q = z.norm()
<span id="L561" class="ln">   561&nbsp;&nbsp;</span>	return
<span id="L562" class="ln">   562&nbsp;&nbsp;</span>}
<span id="L563" class="ln">   563&nbsp;&nbsp;</span>
<span id="L564" class="ln">   564&nbsp;&nbsp;</span><span class="comment">// modW returns x % d.</span>
<span id="L565" class="ln">   565&nbsp;&nbsp;</span>func (x nat) modW(d Word) (r Word) {
<span id="L566" class="ln">   566&nbsp;&nbsp;</span>	<span class="comment">// TODO(agl): we don&#39;t actually need to store the q value.</span>
<span id="L567" class="ln">   567&nbsp;&nbsp;</span>	var q nat
<span id="L568" class="ln">   568&nbsp;&nbsp;</span>	q = q.make(len(x))
<span id="L569" class="ln">   569&nbsp;&nbsp;</span>	return divWVW(q, 0, x, d)
<span id="L570" class="ln">   570&nbsp;&nbsp;</span>}
<span id="L571" class="ln">   571&nbsp;&nbsp;</span>
<span id="L572" class="ln">   572&nbsp;&nbsp;</span><span class="comment">// divWVW overwrites z with ⌊x/y⌋, returning the remainder r.</span>
<span id="L573" class="ln">   573&nbsp;&nbsp;</span><span class="comment">// The caller must ensure that len(z) = len(x).</span>
<span id="L574" class="ln">   574&nbsp;&nbsp;</span>func divWVW(z []Word, xn Word, x []Word, y Word) (r Word) {
<span id="L575" class="ln">   575&nbsp;&nbsp;</span>	r = xn
<span id="L576" class="ln">   576&nbsp;&nbsp;</span>	if len(x) == 1 {
<span id="L577" class="ln">   577&nbsp;&nbsp;</span>		qq, rr := bits.Div(uint(r), uint(x[0]), uint(y))
<span id="L578" class="ln">   578&nbsp;&nbsp;</span>		z[0] = Word(qq)
<span id="L579" class="ln">   579&nbsp;&nbsp;</span>		return Word(rr)
<span id="L580" class="ln">   580&nbsp;&nbsp;</span>	}
<span id="L581" class="ln">   581&nbsp;&nbsp;</span>	rec := reciprocalWord(y)
<span id="L582" class="ln">   582&nbsp;&nbsp;</span>	for i := len(z) - 1; i &gt;= 0; i-- {
<span id="L583" class="ln">   583&nbsp;&nbsp;</span>		z[i], r = divWW(r, x[i], y, rec)
<span id="L584" class="ln">   584&nbsp;&nbsp;</span>	}
<span id="L585" class="ln">   585&nbsp;&nbsp;</span>	return r
<span id="L586" class="ln">   586&nbsp;&nbsp;</span>}
<span id="L587" class="ln">   587&nbsp;&nbsp;</span>
<span id="L588" class="ln">   588&nbsp;&nbsp;</span><span class="comment">// div returns q, r such that q = ⌊uIn/vIn⌋ and r = uIn%vIn = uIn - q·vIn.</span>
<span id="L589" class="ln">   589&nbsp;&nbsp;</span><span class="comment">// It uses z and u as the storage for q and r.</span>
<span id="L590" class="ln">   590&nbsp;&nbsp;</span><span class="comment">// The caller must ensure that len(vIn) ≥ 2 (use divW otherwise)</span>
<span id="L591" class="ln">   591&nbsp;&nbsp;</span><span class="comment">// and that len(uIn) ≥ len(vIn) (the answer is 0, uIn otherwise).</span>
<span id="L592" class="ln">   592&nbsp;&nbsp;</span>func (z nat) divLarge(u, uIn, vIn nat) (q, r nat) {
<span id="L593" class="ln">   593&nbsp;&nbsp;</span>	n := len(vIn)
<span id="L594" class="ln">   594&nbsp;&nbsp;</span>	m := len(uIn) - n
<span id="L595" class="ln">   595&nbsp;&nbsp;</span>
<span id="L596" class="ln">   596&nbsp;&nbsp;</span>	<span class="comment">// Scale the inputs so vIn&#39;s top bit is 1 (see “Scaling Inputs” above).</span>
<span id="L597" class="ln">   597&nbsp;&nbsp;</span>	<span class="comment">// vIn is treated as a read-only input (it may be in use by another</span>
<span id="L598" class="ln">   598&nbsp;&nbsp;</span>	<span class="comment">// goroutine), so we must make a copy.</span>
<span id="L599" class="ln">   599&nbsp;&nbsp;</span>	<span class="comment">// uIn is copied to u.</span>
<span id="L600" class="ln">   600&nbsp;&nbsp;</span>	shift := nlz(vIn[n-1])
<span id="L601" class="ln">   601&nbsp;&nbsp;</span>	vp := getNat(n)
<span id="L602" class="ln">   602&nbsp;&nbsp;</span>	v := *vp
<span id="L603" class="ln">   603&nbsp;&nbsp;</span>	shlVU(v, vIn, shift)
<span id="L604" class="ln">   604&nbsp;&nbsp;</span>	u = u.make(len(uIn) + 1)
<span id="L605" class="ln">   605&nbsp;&nbsp;</span>	u[len(uIn)] = shlVU(u[0:len(uIn)], uIn, shift)
<span id="L606" class="ln">   606&nbsp;&nbsp;</span>
<span id="L607" class="ln">   607&nbsp;&nbsp;</span>	<span class="comment">// The caller should not pass aliased z and u, since those are</span>
<span id="L608" class="ln">   608&nbsp;&nbsp;</span>	<span class="comment">// the two different outputs, but correct just in case.</span>
<span id="L609" class="ln">   609&nbsp;&nbsp;</span>	if alias(z, u) {
<span id="L610" class="ln">   610&nbsp;&nbsp;</span>		z = nil
<span id="L611" class="ln">   611&nbsp;&nbsp;</span>	}
<span id="L612" class="ln">   612&nbsp;&nbsp;</span>	q = z.make(m + 1)
<span id="L613" class="ln">   613&nbsp;&nbsp;</span>
<span id="L614" class="ln">   614&nbsp;&nbsp;</span>	<span class="comment">// Use basic or recursive long division depending on size.</span>
<span id="L615" class="ln">   615&nbsp;&nbsp;</span>	if n &lt; divRecursiveThreshold {
<span id="L616" class="ln">   616&nbsp;&nbsp;</span>		q.divBasic(u, v)
<span id="L617" class="ln">   617&nbsp;&nbsp;</span>	} else {
<span id="L618" class="ln">   618&nbsp;&nbsp;</span>		q.divRecursive(u, v)
<span id="L619" class="ln">   619&nbsp;&nbsp;</span>	}
<span id="L620" class="ln">   620&nbsp;&nbsp;</span>	putNat(vp)
<span id="L621" class="ln">   621&nbsp;&nbsp;</span>
<span id="L622" class="ln">   622&nbsp;&nbsp;</span>	q = q.norm()
<span id="L623" class="ln">   623&nbsp;&nbsp;</span>
<span id="L624" class="ln">   624&nbsp;&nbsp;</span>	<span class="comment">// Undo scaling of remainder.</span>
<span id="L625" class="ln">   625&nbsp;&nbsp;</span>	shrVU(u, u, shift)
<span id="L626" class="ln">   626&nbsp;&nbsp;</span>	r = u.norm()
<span id="L627" class="ln">   627&nbsp;&nbsp;</span>
<span id="L628" class="ln">   628&nbsp;&nbsp;</span>	return q, r
<span id="L629" class="ln">   629&nbsp;&nbsp;</span>}
<span id="L630" class="ln">   630&nbsp;&nbsp;</span>
<span id="L631" class="ln">   631&nbsp;&nbsp;</span><span class="comment">// divBasic implements long division as described above.</span>
<span id="L632" class="ln">   632&nbsp;&nbsp;</span><span class="comment">// It overwrites q with ⌊u/v⌋ and overwrites u with the remainder r.</span>
<span id="L633" class="ln">   633&nbsp;&nbsp;</span><span class="comment">// q must be large enough to hold ⌊u/v⌋.</span>
<span id="L634" class="ln">   634&nbsp;&nbsp;</span>func (q nat) divBasic(u, v nat) {
<span id="L635" class="ln">   635&nbsp;&nbsp;</span>	n := len(v)
<span id="L636" class="ln">   636&nbsp;&nbsp;</span>	m := len(u) - n
<span id="L637" class="ln">   637&nbsp;&nbsp;</span>
<span id="L638" class="ln">   638&nbsp;&nbsp;</span>	qhatvp := getNat(n + 1)
<span id="L639" class="ln">   639&nbsp;&nbsp;</span>	qhatv := *qhatvp
<span id="L640" class="ln">   640&nbsp;&nbsp;</span>
<span id="L641" class="ln">   641&nbsp;&nbsp;</span>	<span class="comment">// Set up for divWW below, precomputing reciprocal argument.</span>
<span id="L642" class="ln">   642&nbsp;&nbsp;</span>	vn1 := v[n-1]
<span id="L643" class="ln">   643&nbsp;&nbsp;</span>	rec := reciprocalWord(vn1)
<span id="L644" class="ln">   644&nbsp;&nbsp;</span>
<span id="L645" class="ln">   645&nbsp;&nbsp;</span>	<span class="comment">// Compute each digit of quotient.</span>
<span id="L646" class="ln">   646&nbsp;&nbsp;</span>	for j := m; j &gt;= 0; j-- {
<span id="L647" class="ln">   647&nbsp;&nbsp;</span>		<span class="comment">// Compute the 2-by-1 guess q̂.</span>
<span id="L648" class="ln">   648&nbsp;&nbsp;</span>		<span class="comment">// The first iteration must invent a leading 0 for u.</span>
<span id="L649" class="ln">   649&nbsp;&nbsp;</span>		qhat := Word(_M)
<span id="L650" class="ln">   650&nbsp;&nbsp;</span>		var ujn Word
<span id="L651" class="ln">   651&nbsp;&nbsp;</span>		if j+n &lt; len(u) {
<span id="L652" class="ln">   652&nbsp;&nbsp;</span>			ujn = u[j+n]
<span id="L653" class="ln">   653&nbsp;&nbsp;</span>		}
<span id="L654" class="ln">   654&nbsp;&nbsp;</span>
<span id="L655" class="ln">   655&nbsp;&nbsp;</span>		<span class="comment">// ujn ≤ vn1, or else q̂ would be more than one digit.</span>
<span id="L656" class="ln">   656&nbsp;&nbsp;</span>		<span class="comment">// For ujn == vn1, we set q̂ to the max digit M above.</span>
<span id="L657" class="ln">   657&nbsp;&nbsp;</span>		<span class="comment">// Otherwise, we compute the 2-by-1 guess.</span>
<span id="L658" class="ln">   658&nbsp;&nbsp;</span>		if ujn != vn1 {
<span id="L659" class="ln">   659&nbsp;&nbsp;</span>			var rhat Word
<span id="L660" class="ln">   660&nbsp;&nbsp;</span>			qhat, rhat = divWW(ujn, u[j+n-1], vn1, rec)
<span id="L661" class="ln">   661&nbsp;&nbsp;</span>
<span id="L662" class="ln">   662&nbsp;&nbsp;</span>			<span class="comment">// Refine q̂ to a 3-by-2 guess. See “Refining Guesses” above.</span>
<span id="L663" class="ln">   663&nbsp;&nbsp;</span>			vn2 := v[n-2]
<span id="L664" class="ln">   664&nbsp;&nbsp;</span>			x1, x2 := mulWW(qhat, vn2)
<span id="L665" class="ln">   665&nbsp;&nbsp;</span>			ujn2 := u[j+n-2]
<span id="L666" class="ln">   666&nbsp;&nbsp;</span>			for greaterThan(x1, x2, rhat, ujn2) { <span class="comment">// x1x2 &gt; r̂ u[j+n-2]</span>
<span id="L667" class="ln">   667&nbsp;&nbsp;</span>				qhat--
<span id="L668" class="ln">   668&nbsp;&nbsp;</span>				prevRhat := rhat
<span id="L669" class="ln">   669&nbsp;&nbsp;</span>				rhat += vn1
<span id="L670" class="ln">   670&nbsp;&nbsp;</span>				<span class="comment">// If r̂  overflows, then</span>
<span id="L671" class="ln">   671&nbsp;&nbsp;</span>				<span class="comment">// r̂ u[j+n-2]v[n-1] is now definitely &gt; x1 x2.</span>
<span id="L672" class="ln">   672&nbsp;&nbsp;</span>				if rhat &lt; prevRhat {
<span id="L673" class="ln">   673&nbsp;&nbsp;</span>					break
<span id="L674" class="ln">   674&nbsp;&nbsp;</span>				}
<span id="L675" class="ln">   675&nbsp;&nbsp;</span>				<span class="comment">// TODO(rsc): No need for a full mulWW.</span>
<span id="L676" class="ln">   676&nbsp;&nbsp;</span>				<span class="comment">// x2 += vn2; if x2 overflows, x1++</span>
<span id="L677" class="ln">   677&nbsp;&nbsp;</span>				x1, x2 = mulWW(qhat, vn2)
<span id="L678" class="ln">   678&nbsp;&nbsp;</span>			}
<span id="L679" class="ln">   679&nbsp;&nbsp;</span>		}
<span id="L680" class="ln">   680&nbsp;&nbsp;</span>
<span id="L681" class="ln">   681&nbsp;&nbsp;</span>		<span class="comment">// Compute q̂·v.</span>
<span id="L682" class="ln">   682&nbsp;&nbsp;</span>		qhatv[n] = mulAddVWW(qhatv[0:n], v, qhat, 0)
<span id="L683" class="ln">   683&nbsp;&nbsp;</span>		qhl := len(qhatv)
<span id="L684" class="ln">   684&nbsp;&nbsp;</span>		if j+qhl &gt; len(u) &amp;&amp; qhatv[n] == 0 {
<span id="L685" class="ln">   685&nbsp;&nbsp;</span>			qhl--
<span id="L686" class="ln">   686&nbsp;&nbsp;</span>		}
<span id="L687" class="ln">   687&nbsp;&nbsp;</span>
<span id="L688" class="ln">   688&nbsp;&nbsp;</span>		<span class="comment">// Subtract q̂·v from the current section of u.</span>
<span id="L689" class="ln">   689&nbsp;&nbsp;</span>		<span class="comment">// If it underflows, q̂·v &gt; u, which we fix up</span>
<span id="L690" class="ln">   690&nbsp;&nbsp;</span>		<span class="comment">// by decrementing q̂ and adding v back.</span>
<span id="L691" class="ln">   691&nbsp;&nbsp;</span>		c := subVV(u[j:j+qhl], u[j:], qhatv)
<span id="L692" class="ln">   692&nbsp;&nbsp;</span>		if c != 0 {
<span id="L693" class="ln">   693&nbsp;&nbsp;</span>			c := addVV(u[j:j+n], u[j:], v)
<span id="L694" class="ln">   694&nbsp;&nbsp;</span>			<span class="comment">// If n == qhl, the carry from subVV and the carry from addVV</span>
<span id="L695" class="ln">   695&nbsp;&nbsp;</span>			<span class="comment">// cancel out and don&#39;t affect u[j+n].</span>
<span id="L696" class="ln">   696&nbsp;&nbsp;</span>			if n &lt; qhl {
<span id="L697" class="ln">   697&nbsp;&nbsp;</span>				u[j+n] += c
<span id="L698" class="ln">   698&nbsp;&nbsp;</span>			}
<span id="L699" class="ln">   699&nbsp;&nbsp;</span>			qhat--
<span id="L700" class="ln">   700&nbsp;&nbsp;</span>		}
<span id="L701" class="ln">   701&nbsp;&nbsp;</span>
<span id="L702" class="ln">   702&nbsp;&nbsp;</span>		<span class="comment">// Save quotient digit.</span>
<span id="L703" class="ln">   703&nbsp;&nbsp;</span>		<span class="comment">// Caller may know the top digit is zero and not leave room for it.</span>
<span id="L704" class="ln">   704&nbsp;&nbsp;</span>		if j == m &amp;&amp; m == len(q) &amp;&amp; qhat == 0 {
<span id="L705" class="ln">   705&nbsp;&nbsp;</span>			continue
<span id="L706" class="ln">   706&nbsp;&nbsp;</span>		}
<span id="L707" class="ln">   707&nbsp;&nbsp;</span>		q[j] = qhat
<span id="L708" class="ln">   708&nbsp;&nbsp;</span>	}
<span id="L709" class="ln">   709&nbsp;&nbsp;</span>
<span id="L710" class="ln">   710&nbsp;&nbsp;</span>	putNat(qhatvp)
<span id="L711" class="ln">   711&nbsp;&nbsp;</span>}
<span id="L712" class="ln">   712&nbsp;&nbsp;</span>
<span id="L713" class="ln">   713&nbsp;&nbsp;</span><span class="comment">// greaterThan reports whether the two digit numbers x1 x2 &gt; y1 y2.</span>
<span id="L714" class="ln">   714&nbsp;&nbsp;</span><span class="comment">// TODO(rsc): In contradiction to most of this file, x1 is the high</span>
<span id="L715" class="ln">   715&nbsp;&nbsp;</span><span class="comment">// digit and x2 is the low digit. This should be fixed.</span>
<span id="L716" class="ln">   716&nbsp;&nbsp;</span>func greaterThan(x1, x2, y1, y2 Word) bool {
<span id="L717" class="ln">   717&nbsp;&nbsp;</span>	return x1 &gt; y1 || x1 == y1 &amp;&amp; x2 &gt; y2
<span id="L718" class="ln">   718&nbsp;&nbsp;</span>}
<span id="L719" class="ln">   719&nbsp;&nbsp;</span>
<span id="L720" class="ln">   720&nbsp;&nbsp;</span><span class="comment">// divRecursiveThreshold is the number of divisor digits</span>
<span id="L721" class="ln">   721&nbsp;&nbsp;</span><span class="comment">// at which point divRecursive is faster than divBasic.</span>
<span id="L722" class="ln">   722&nbsp;&nbsp;</span>const divRecursiveThreshold = 100
<span id="L723" class="ln">   723&nbsp;&nbsp;</span>
<span id="L724" class="ln">   724&nbsp;&nbsp;</span><span class="comment">// divRecursive implements recursive division as described above.</span>
<span id="L725" class="ln">   725&nbsp;&nbsp;</span><span class="comment">// It overwrites z with ⌊u/v⌋ and overwrites u with the remainder r.</span>
<span id="L726" class="ln">   726&nbsp;&nbsp;</span><span class="comment">// z must be large enough to hold ⌊u/v⌋.</span>
<span id="L727" class="ln">   727&nbsp;&nbsp;</span><span class="comment">// This function is just for allocating and freeing temporaries</span>
<span id="L728" class="ln">   728&nbsp;&nbsp;</span><span class="comment">// around divRecursiveStep, the real implementation.</span>
<span id="L729" class="ln">   729&nbsp;&nbsp;</span>func (z nat) divRecursive(u, v nat) {
<span id="L730" class="ln">   730&nbsp;&nbsp;</span>	<span class="comment">// Recursion depth is (much) less than 2 log₂(len(v)).</span>
<span id="L731" class="ln">   731&nbsp;&nbsp;</span>	<span class="comment">// Allocate a slice of temporaries to be reused across recursion,</span>
<span id="L732" class="ln">   732&nbsp;&nbsp;</span>	<span class="comment">// plus one extra temporary not live across the recursion.</span>
<span id="L733" class="ln">   733&nbsp;&nbsp;</span>	recDepth := 2 * bits.Len(uint(len(v)))
<span id="L734" class="ln">   734&nbsp;&nbsp;</span>	tmp := getNat(3 * len(v))
<span id="L735" class="ln">   735&nbsp;&nbsp;</span>	temps := make([]*nat, recDepth)
<span id="L736" class="ln">   736&nbsp;&nbsp;</span>
<span id="L737" class="ln">   737&nbsp;&nbsp;</span>	z.clear()
<span id="L738" class="ln">   738&nbsp;&nbsp;</span>	z.divRecursiveStep(u, v, 0, tmp, temps)
<span id="L739" class="ln">   739&nbsp;&nbsp;</span>
<span id="L740" class="ln">   740&nbsp;&nbsp;</span>	<span class="comment">// Free temporaries.</span>
<span id="L741" class="ln">   741&nbsp;&nbsp;</span>	for _, n := range temps {
<span id="L742" class="ln">   742&nbsp;&nbsp;</span>		if n != nil {
<span id="L743" class="ln">   743&nbsp;&nbsp;</span>			putNat(n)
<span id="L744" class="ln">   744&nbsp;&nbsp;</span>		}
<span id="L745" class="ln">   745&nbsp;&nbsp;</span>	}
<span id="L746" class="ln">   746&nbsp;&nbsp;</span>	putNat(tmp)
<span id="L747" class="ln">   747&nbsp;&nbsp;</span>}
<span id="L748" class="ln">   748&nbsp;&nbsp;</span>
<span id="L749" class="ln">   749&nbsp;&nbsp;</span><span class="comment">// divRecursiveStep is the actual implementation of recursive division.</span>
<span id="L750" class="ln">   750&nbsp;&nbsp;</span><span class="comment">// It adds ⌊u/v⌋ to z and overwrites u with the remainder r.</span>
<span id="L751" class="ln">   751&nbsp;&nbsp;</span><span class="comment">// z must be large enough to hold ⌊u/v⌋.</span>
<span id="L752" class="ln">   752&nbsp;&nbsp;</span><span class="comment">// It uses temps[depth] (allocating if needed) as a temporary live across</span>
<span id="L753" class="ln">   753&nbsp;&nbsp;</span><span class="comment">// the recursive call. It also uses tmp, but not live across the recursion.</span>
<span id="L754" class="ln">   754&nbsp;&nbsp;</span>func (z nat) divRecursiveStep(u, v nat, depth int, tmp *nat, temps []*nat) {
<span id="L755" class="ln">   755&nbsp;&nbsp;</span>	<span class="comment">// u is a subsection of the original and may have leading zeros.</span>
<span id="L756" class="ln">   756&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc): The v = v.norm() is useless and should be removed.</span>
<span id="L757" class="ln">   757&nbsp;&nbsp;</span>	<span class="comment">// We know (and require) that v&#39;s top digit is ≥ B/2.</span>
<span id="L758" class="ln">   758&nbsp;&nbsp;</span>	u = u.norm()
<span id="L759" class="ln">   759&nbsp;&nbsp;</span>	v = v.norm()
<span id="L760" class="ln">   760&nbsp;&nbsp;</span>	if len(u) == 0 {
<span id="L761" class="ln">   761&nbsp;&nbsp;</span>		z.clear()
<span id="L762" class="ln">   762&nbsp;&nbsp;</span>		return
<span id="L763" class="ln">   763&nbsp;&nbsp;</span>	}
<span id="L764" class="ln">   764&nbsp;&nbsp;</span>
<span id="L765" class="ln">   765&nbsp;&nbsp;</span>	<span class="comment">// Fall back to basic division if the problem is now small enough.</span>
<span id="L766" class="ln">   766&nbsp;&nbsp;</span>	n := len(v)
<span id="L767" class="ln">   767&nbsp;&nbsp;</span>	if n &lt; divRecursiveThreshold {
<span id="L768" class="ln">   768&nbsp;&nbsp;</span>		z.divBasic(u, v)
<span id="L769" class="ln">   769&nbsp;&nbsp;</span>		return
<span id="L770" class="ln">   770&nbsp;&nbsp;</span>	}
<span id="L771" class="ln">   771&nbsp;&nbsp;</span>
<span id="L772" class="ln">   772&nbsp;&nbsp;</span>	<span class="comment">// Nothing to do if u is shorter than v (implies u &lt; v).</span>
<span id="L773" class="ln">   773&nbsp;&nbsp;</span>	m := len(u) - n
<span id="L774" class="ln">   774&nbsp;&nbsp;</span>	if m &lt; 0 {
<span id="L775" class="ln">   775&nbsp;&nbsp;</span>		return
<span id="L776" class="ln">   776&nbsp;&nbsp;</span>	}
<span id="L777" class="ln">   777&nbsp;&nbsp;</span>
<span id="L778" class="ln">   778&nbsp;&nbsp;</span>	<span class="comment">// We consider B digits in a row as a single wide digit.</span>
<span id="L779" class="ln">   779&nbsp;&nbsp;</span>	<span class="comment">// (See “Recursive Division” above.)</span>
<span id="L780" class="ln">   780&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L781" class="ln">   781&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc): rename B to Wide, to avoid confusion with _B,</span>
<span id="L782" class="ln">   782&nbsp;&nbsp;</span>	<span class="comment">// which is something entirely different.</span>
<span id="L783" class="ln">   783&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc): Look into whether using ⌈n/2⌉ is better than ⌊n/2⌋.</span>
<span id="L784" class="ln">   784&nbsp;&nbsp;</span>	B := n / 2
<span id="L785" class="ln">   785&nbsp;&nbsp;</span>
<span id="L786" class="ln">   786&nbsp;&nbsp;</span>	<span class="comment">// Allocate a nat for qhat below.</span>
<span id="L787" class="ln">   787&nbsp;&nbsp;</span>	if temps[depth] == nil {
<span id="L788" class="ln">   788&nbsp;&nbsp;</span>		temps[depth] = getNat(n) <span class="comment">// TODO(rsc): Can be just B+1.</span>
<span id="L789" class="ln">   789&nbsp;&nbsp;</span>	} else {
<span id="L790" class="ln">   790&nbsp;&nbsp;</span>		*temps[depth] = temps[depth].make(B + 1)
<span id="L791" class="ln">   791&nbsp;&nbsp;</span>	}
<span id="L792" class="ln">   792&nbsp;&nbsp;</span>
<span id="L793" class="ln">   793&nbsp;&nbsp;</span>	<span class="comment">// Compute each wide digit of the quotient.</span>
<span id="L794" class="ln">   794&nbsp;&nbsp;</span>	<span class="comment">//</span>
<span id="L795" class="ln">   795&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc): Change the loop to be</span>
<span id="L796" class="ln">   796&nbsp;&nbsp;</span>	<span class="comment">//	for j := (m+B-1)/B*B; j &gt; 0; j -= B {</span>
<span id="L797" class="ln">   797&nbsp;&nbsp;</span>	<span class="comment">// which will make the final step a regular step, letting us</span>
<span id="L798" class="ln">   798&nbsp;&nbsp;</span>	<span class="comment">// delete what amounts to an extra copy of the loop body below.</span>
<span id="L799" class="ln">   799&nbsp;&nbsp;</span>	j := m
<span id="L800" class="ln">   800&nbsp;&nbsp;</span>	for j &gt; B {
<span id="L801" class="ln">   801&nbsp;&nbsp;</span>		<span class="comment">// Divide u[j-B:j+n] (3 wide digits) by v (2 wide digits).</span>
<span id="L802" class="ln">   802&nbsp;&nbsp;</span>		<span class="comment">// First make the 2-by-1-wide-digit guess using a recursive call.</span>
<span id="L803" class="ln">   803&nbsp;&nbsp;</span>		<span class="comment">// Then extend the guess to the full 3-by-2 (see “Refining Guesses”).</span>
<span id="L804" class="ln">   804&nbsp;&nbsp;</span>		<span class="comment">//</span>
<span id="L805" class="ln">   805&nbsp;&nbsp;</span>		<span class="comment">// For the 2-by-1-wide-digit guess, instead of doing 2B-by-B-digit,</span>
<span id="L806" class="ln">   806&nbsp;&nbsp;</span>		<span class="comment">// we use a (2B+1)-by-(B+1) digit, which handles the possibility that</span>
<span id="L807" class="ln">   807&nbsp;&nbsp;</span>		<span class="comment">// the result has an extra leading 1 digit as well as guaranteeing</span>
<span id="L808" class="ln">   808&nbsp;&nbsp;</span>		<span class="comment">// that the computed q̂ will be off by at most 1 instead of 2.</span>
<span id="L809" class="ln">   809&nbsp;&nbsp;</span>
<span id="L810" class="ln">   810&nbsp;&nbsp;</span>		<span class="comment">// s is the number of digits to drop from the 3B- and 2B-digit chunks.</span>
<span id="L811" class="ln">   811&nbsp;&nbsp;</span>		<span class="comment">// We drop B-1 to be left with 2B+1 and B+1.</span>
<span id="L812" class="ln">   812&nbsp;&nbsp;</span>		s := (B - 1)
<span id="L813" class="ln">   813&nbsp;&nbsp;</span>
<span id="L814" class="ln">   814&nbsp;&nbsp;</span>		<span class="comment">// uu is the up-to-3B-digit section of u we are working on.</span>
<span id="L815" class="ln">   815&nbsp;&nbsp;</span>		uu := u[j-B:]
<span id="L816" class="ln">   816&nbsp;&nbsp;</span>
<span id="L817" class="ln">   817&nbsp;&nbsp;</span>		<span class="comment">// Compute the 2-by-1 guess q̂, leaving r̂ in uu[s:B+n].</span>
<span id="L818" class="ln">   818&nbsp;&nbsp;</span>		qhat := *temps[depth]
<span id="L819" class="ln">   819&nbsp;&nbsp;</span>		qhat.clear()
<span id="L820" class="ln">   820&nbsp;&nbsp;</span>		qhat.divRecursiveStep(uu[s:B+n], v[s:], depth+1, tmp, temps)
<span id="L821" class="ln">   821&nbsp;&nbsp;</span>		qhat = qhat.norm()
<span id="L822" class="ln">   822&nbsp;&nbsp;</span>
<span id="L823" class="ln">   823&nbsp;&nbsp;</span>		<span class="comment">// Extend to a 3-by-2 quotient and remainder.</span>
<span id="L824" class="ln">   824&nbsp;&nbsp;</span>		<span class="comment">// Because divRecursiveStep overwrote the top part of uu with</span>
<span id="L825" class="ln">   825&nbsp;&nbsp;</span>		<span class="comment">// the remainder r̂, the full uu already contains the equivalent</span>
<span id="L826" class="ln">   826&nbsp;&nbsp;</span>		<span class="comment">// of r̂·B + uₙ₋₂ from the “Refining Guesses” discussion.</span>
<span id="L827" class="ln">   827&nbsp;&nbsp;</span>		<span class="comment">// Subtracting q̂·vₙ₋₂ from it will compute the full-length remainder.</span>
<span id="L828" class="ln">   828&nbsp;&nbsp;</span>		<span class="comment">// If that subtraction underflows, q̂·v &gt; u, which we fix up</span>
<span id="L829" class="ln">   829&nbsp;&nbsp;</span>		<span class="comment">// by decrementing q̂ and adding v back, same as in long division.</span>
<span id="L830" class="ln">   830&nbsp;&nbsp;</span>
<span id="L831" class="ln">   831&nbsp;&nbsp;</span>		<span class="comment">// TODO(rsc): Instead of subtract and fix-up, this code is computing</span>
<span id="L832" class="ln">   832&nbsp;&nbsp;</span>		<span class="comment">// q̂·vₙ₋₂ and decrementing q̂ until that product is ≤ u.</span>
<span id="L833" class="ln">   833&nbsp;&nbsp;</span>		<span class="comment">// But we can do the subtraction directly, as in the comment above</span>
<span id="L834" class="ln">   834&nbsp;&nbsp;</span>		<span class="comment">// and in long division, because we know that q̂ is wrong by at most one.</span>
<span id="L835" class="ln">   835&nbsp;&nbsp;</span>		qhatv := tmp.make(3 * n)
<span id="L836" class="ln">   836&nbsp;&nbsp;</span>		qhatv.clear()
<span id="L837" class="ln">   837&nbsp;&nbsp;</span>		qhatv = qhatv.mul(qhat, v[:s])
<span id="L838" class="ln">   838&nbsp;&nbsp;</span>		for i := 0; i &lt; 2; i++ {
<span id="L839" class="ln">   839&nbsp;&nbsp;</span>			e := qhatv.cmp(uu.norm())
<span id="L840" class="ln">   840&nbsp;&nbsp;</span>			if e &lt;= 0 {
<span id="L841" class="ln">   841&nbsp;&nbsp;</span>				break
<span id="L842" class="ln">   842&nbsp;&nbsp;</span>			}
<span id="L843" class="ln">   843&nbsp;&nbsp;</span>			subVW(qhat, qhat, 1)
<span id="L844" class="ln">   844&nbsp;&nbsp;</span>			c := subVV(qhatv[:s], qhatv[:s], v[:s])
<span id="L845" class="ln">   845&nbsp;&nbsp;</span>			if len(qhatv) &gt; s {
<span id="L846" class="ln">   846&nbsp;&nbsp;</span>				subVW(qhatv[s:], qhatv[s:], c)
<span id="L847" class="ln">   847&nbsp;&nbsp;</span>			}
<span id="L848" class="ln">   848&nbsp;&nbsp;</span>			addAt(uu[s:], v[s:], 0)
<span id="L849" class="ln">   849&nbsp;&nbsp;</span>		}
<span id="L850" class="ln">   850&nbsp;&nbsp;</span>		if qhatv.cmp(uu.norm()) &gt; 0 {
<span id="L851" class="ln">   851&nbsp;&nbsp;</span>			panic(&#34;impossible&#34;)
<span id="L852" class="ln">   852&nbsp;&nbsp;</span>		}
<span id="L853" class="ln">   853&nbsp;&nbsp;</span>		c := subVV(uu[:len(qhatv)], uu[:len(qhatv)], qhatv)
<span id="L854" class="ln">   854&nbsp;&nbsp;</span>		if c &gt; 0 {
<span id="L855" class="ln">   855&nbsp;&nbsp;</span>			subVW(uu[len(qhatv):], uu[len(qhatv):], c)
<span id="L856" class="ln">   856&nbsp;&nbsp;</span>		}
<span id="L857" class="ln">   857&nbsp;&nbsp;</span>		addAt(z, qhat, j-B)
<span id="L858" class="ln">   858&nbsp;&nbsp;</span>		j -= B
<span id="L859" class="ln">   859&nbsp;&nbsp;</span>	}
<span id="L860" class="ln">   860&nbsp;&nbsp;</span>
<span id="L861" class="ln">   861&nbsp;&nbsp;</span>	<span class="comment">// TODO(rsc): Rewrite loop as described above and delete all this code.</span>
<span id="L862" class="ln">   862&nbsp;&nbsp;</span>
<span id="L863" class="ln">   863&nbsp;&nbsp;</span>	<span class="comment">// Now u &lt; (v&lt;&lt;B), compute lower bits in the same way.</span>
<span id="L864" class="ln">   864&nbsp;&nbsp;</span>	<span class="comment">// Choose shift = B-1 again.</span>
<span id="L865" class="ln">   865&nbsp;&nbsp;</span>	s := B - 1
<span id="L866" class="ln">   866&nbsp;&nbsp;</span>	qhat := *temps[depth]
<span id="L867" class="ln">   867&nbsp;&nbsp;</span>	qhat.clear()
<span id="L868" class="ln">   868&nbsp;&nbsp;</span>	qhat.divRecursiveStep(u[s:].norm(), v[s:], depth+1, tmp, temps)
<span id="L869" class="ln">   869&nbsp;&nbsp;</span>	qhat = qhat.norm()
<span id="L870" class="ln">   870&nbsp;&nbsp;</span>	qhatv := tmp.make(3 * n)
<span id="L871" class="ln">   871&nbsp;&nbsp;</span>	qhatv.clear()
<span id="L872" class="ln">   872&nbsp;&nbsp;</span>	qhatv = qhatv.mul(qhat, v[:s])
<span id="L873" class="ln">   873&nbsp;&nbsp;</span>	<span class="comment">// Set the correct remainder as before.</span>
<span id="L874" class="ln">   874&nbsp;&nbsp;</span>	for i := 0; i &lt; 2; i++ {
<span id="L875" class="ln">   875&nbsp;&nbsp;</span>		if e := qhatv.cmp(u.norm()); e &gt; 0 {
<span id="L876" class="ln">   876&nbsp;&nbsp;</span>			subVW(qhat, qhat, 1)
<span id="L877" class="ln">   877&nbsp;&nbsp;</span>			c := subVV(qhatv[:s], qhatv[:s], v[:s])
<span id="L878" class="ln">   878&nbsp;&nbsp;</span>			if len(qhatv) &gt; s {
<span id="L879" class="ln">   879&nbsp;&nbsp;</span>				subVW(qhatv[s:], qhatv[s:], c)
<span id="L880" class="ln">   880&nbsp;&nbsp;</span>			}
<span id="L881" class="ln">   881&nbsp;&nbsp;</span>			addAt(u[s:], v[s:], 0)
<span id="L882" class="ln">   882&nbsp;&nbsp;</span>		}
<span id="L883" class="ln">   883&nbsp;&nbsp;</span>	}
<span id="L884" class="ln">   884&nbsp;&nbsp;</span>	if qhatv.cmp(u.norm()) &gt; 0 {
<span id="L885" class="ln">   885&nbsp;&nbsp;</span>		panic(&#34;impossible&#34;)
<span id="L886" class="ln">   886&nbsp;&nbsp;</span>	}
<span id="L887" class="ln">   887&nbsp;&nbsp;</span>	c := subVV(u[0:len(qhatv)], u[0:len(qhatv)], qhatv)
<span id="L888" class="ln">   888&nbsp;&nbsp;</span>	if c &gt; 0 {
<span id="L889" class="ln">   889&nbsp;&nbsp;</span>		c = subVW(u[len(qhatv):], u[len(qhatv):], c)
<span id="L890" class="ln">   890&nbsp;&nbsp;</span>	}
<span id="L891" class="ln">   891&nbsp;&nbsp;</span>	if c &gt; 0 {
<span id="L892" class="ln">   892&nbsp;&nbsp;</span>		panic(&#34;impossible&#34;)
<span id="L893" class="ln">   893&nbsp;&nbsp;</span>	}
<span id="L894" class="ln">   894&nbsp;&nbsp;</span>
<span id="L895" class="ln">   895&nbsp;&nbsp;</span>	<span class="comment">// Done!</span>
<span id="L896" class="ln">   896&nbsp;&nbsp;</span>	addAt(z, qhat.norm(), 0)
<span id="L897" class="ln">   897&nbsp;&nbsp;</span>}
<span id="L898" class="ln">   898&nbsp;&nbsp;</span>
</pre><p><a href="natdiv.go?m=text">View as plain text</a></p>

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
