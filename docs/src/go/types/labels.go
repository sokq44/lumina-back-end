<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>src/go/types/labels.go - Go Documentation Server</title>

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
<a href="labels.go#" id="menu-button"><span id="menu-button-arrow">&#9661;</span></a>
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
    <a href="http://localhost:8080/src">src</a>/<a href="http://localhost:8080/src/go">go</a>/<a href="http://localhost:8080/src/go/types">types</a>/<span class="text-muted">labels.go</span>
  </h1>





  <h2>
    Documentation: <a href="http://localhost:8080/pkg/go/types">go/types</a>
  </h2>



<div id="nav"></div>


<script type='text/javascript'>document.ANALYSIS_DATA = null;</script>
<pre><span id="L1" class="ln">     1&nbsp;&nbsp;</span><span class="comment">// Copyright 2013 The Go Authors. All rights reserved.</span>
<span id="L2" class="ln">     2&nbsp;&nbsp;</span><span class="comment">// Use of this source code is governed by a BSD-style</span>
<span id="L3" class="ln">     3&nbsp;&nbsp;</span><span class="comment">// license that can be found in the LICENSE file.</span>
<span id="L4" class="ln">     4&nbsp;&nbsp;</span>
<span id="L5" class="ln">     5&nbsp;&nbsp;</span>package types
<span id="L6" class="ln">     6&nbsp;&nbsp;</span>
<span id="L7" class="ln">     7&nbsp;&nbsp;</span>import (
<span id="L8" class="ln">     8&nbsp;&nbsp;</span>	&#34;go/ast&#34;
<span id="L9" class="ln">     9&nbsp;&nbsp;</span>	&#34;go/token&#34;
<span id="L10" class="ln">    10&nbsp;&nbsp;</span>	. &#34;internal/types/errors&#34;
<span id="L11" class="ln">    11&nbsp;&nbsp;</span>)
<span id="L12" class="ln">    12&nbsp;&nbsp;</span>
<span id="L13" class="ln">    13&nbsp;&nbsp;</span><span class="comment">// labels checks correct label use in body.</span>
<span id="L14" class="ln">    14&nbsp;&nbsp;</span>func (check *Checker) labels(body *ast.BlockStmt) {
<span id="L15" class="ln">    15&nbsp;&nbsp;</span>	<span class="comment">// set of all labels in this body</span>
<span id="L16" class="ln">    16&nbsp;&nbsp;</span>	all := NewScope(nil, body.Pos(), body.End(), &#34;label&#34;)
<span id="L17" class="ln">    17&nbsp;&nbsp;</span>
<span id="L18" class="ln">    18&nbsp;&nbsp;</span>	fwdJumps := check.blockBranches(all, nil, nil, body.List)
<span id="L19" class="ln">    19&nbsp;&nbsp;</span>
<span id="L20" class="ln">    20&nbsp;&nbsp;</span>	<span class="comment">// If there are any forward jumps left, no label was found for</span>
<span id="L21" class="ln">    21&nbsp;&nbsp;</span>	<span class="comment">// the corresponding goto statements. Either those labels were</span>
<span id="L22" class="ln">    22&nbsp;&nbsp;</span>	<span class="comment">// never defined, or they are inside blocks and not reachable</span>
<span id="L23" class="ln">    23&nbsp;&nbsp;</span>	<span class="comment">// for the respective gotos.</span>
<span id="L24" class="ln">    24&nbsp;&nbsp;</span>	for _, jmp := range fwdJumps {
<span id="L25" class="ln">    25&nbsp;&nbsp;</span>		var msg string
<span id="L26" class="ln">    26&nbsp;&nbsp;</span>		var code Code
<span id="L27" class="ln">    27&nbsp;&nbsp;</span>		name := jmp.Label.Name
<span id="L28" class="ln">    28&nbsp;&nbsp;</span>		if alt := all.Lookup(name); alt != nil {
<span id="L29" class="ln">    29&nbsp;&nbsp;</span>			msg = &#34;goto %s jumps into block&#34;
<span id="L30" class="ln">    30&nbsp;&nbsp;</span>			alt.(*Label).used = true <span class="comment">// avoid another error</span>
<span id="L31" class="ln">    31&nbsp;&nbsp;</span>			code = JumpIntoBlock
<span id="L32" class="ln">    32&nbsp;&nbsp;</span>		} else {
<span id="L33" class="ln">    33&nbsp;&nbsp;</span>			msg = &#34;label %s not declared&#34;
<span id="L34" class="ln">    34&nbsp;&nbsp;</span>			code = UndeclaredLabel
<span id="L35" class="ln">    35&nbsp;&nbsp;</span>		}
<span id="L36" class="ln">    36&nbsp;&nbsp;</span>		check.errorf(jmp.Label, code, msg, name)
<span id="L37" class="ln">    37&nbsp;&nbsp;</span>	}
<span id="L38" class="ln">    38&nbsp;&nbsp;</span>
<span id="L39" class="ln">    39&nbsp;&nbsp;</span>	<span class="comment">// spec: &#34;It is illegal to define a label that is never used.&#34;</span>
<span id="L40" class="ln">    40&nbsp;&nbsp;</span>	for name, obj := range all.elems {
<span id="L41" class="ln">    41&nbsp;&nbsp;</span>		obj = resolve(name, obj)
<span id="L42" class="ln">    42&nbsp;&nbsp;</span>		if lbl := obj.(*Label); !lbl.used {
<span id="L43" class="ln">    43&nbsp;&nbsp;</span>			check.softErrorf(lbl, UnusedLabel, &#34;label %s declared and not used&#34;, lbl.name)
<span id="L44" class="ln">    44&nbsp;&nbsp;</span>		}
<span id="L45" class="ln">    45&nbsp;&nbsp;</span>	}
<span id="L46" class="ln">    46&nbsp;&nbsp;</span>}
<span id="L47" class="ln">    47&nbsp;&nbsp;</span>
<span id="L48" class="ln">    48&nbsp;&nbsp;</span><span class="comment">// A block tracks label declarations in a block and its enclosing blocks.</span>
<span id="L49" class="ln">    49&nbsp;&nbsp;</span>type block struct {
<span id="L50" class="ln">    50&nbsp;&nbsp;</span>	parent *block                      <span class="comment">// enclosing block</span>
<span id="L51" class="ln">    51&nbsp;&nbsp;</span>	lstmt  *ast.LabeledStmt            <span class="comment">// labeled statement to which this block belongs, or nil</span>
<span id="L52" class="ln">    52&nbsp;&nbsp;</span>	labels map[string]*ast.LabeledStmt <span class="comment">// allocated lazily</span>
<span id="L53" class="ln">    53&nbsp;&nbsp;</span>}
<span id="L54" class="ln">    54&nbsp;&nbsp;</span>
<span id="L55" class="ln">    55&nbsp;&nbsp;</span><span class="comment">// insert records a new label declaration for the current block.</span>
<span id="L56" class="ln">    56&nbsp;&nbsp;</span><span class="comment">// The label must not have been declared before in any block.</span>
<span id="L57" class="ln">    57&nbsp;&nbsp;</span>func (b *block) insert(s *ast.LabeledStmt) {
<span id="L58" class="ln">    58&nbsp;&nbsp;</span>	name := s.Label.Name
<span id="L59" class="ln">    59&nbsp;&nbsp;</span>	if debug {
<span id="L60" class="ln">    60&nbsp;&nbsp;</span>		assert(b.gotoTarget(name) == nil)
<span id="L61" class="ln">    61&nbsp;&nbsp;</span>	}
<span id="L62" class="ln">    62&nbsp;&nbsp;</span>	labels := b.labels
<span id="L63" class="ln">    63&nbsp;&nbsp;</span>	if labels == nil {
<span id="L64" class="ln">    64&nbsp;&nbsp;</span>		labels = make(map[string]*ast.LabeledStmt)
<span id="L65" class="ln">    65&nbsp;&nbsp;</span>		b.labels = labels
<span id="L66" class="ln">    66&nbsp;&nbsp;</span>	}
<span id="L67" class="ln">    67&nbsp;&nbsp;</span>	labels[name] = s
<span id="L68" class="ln">    68&nbsp;&nbsp;</span>}
<span id="L69" class="ln">    69&nbsp;&nbsp;</span>
<span id="L70" class="ln">    70&nbsp;&nbsp;</span><span class="comment">// gotoTarget returns the labeled statement in the current</span>
<span id="L71" class="ln">    71&nbsp;&nbsp;</span><span class="comment">// or an enclosing block with the given label name, or nil.</span>
<span id="L72" class="ln">    72&nbsp;&nbsp;</span>func (b *block) gotoTarget(name string) *ast.LabeledStmt {
<span id="L73" class="ln">    73&nbsp;&nbsp;</span>	for s := b; s != nil; s = s.parent {
<span id="L74" class="ln">    74&nbsp;&nbsp;</span>		if t := s.labels[name]; t != nil {
<span id="L75" class="ln">    75&nbsp;&nbsp;</span>			return t
<span id="L76" class="ln">    76&nbsp;&nbsp;</span>		}
<span id="L77" class="ln">    77&nbsp;&nbsp;</span>	}
<span id="L78" class="ln">    78&nbsp;&nbsp;</span>	return nil
<span id="L79" class="ln">    79&nbsp;&nbsp;</span>}
<span id="L80" class="ln">    80&nbsp;&nbsp;</span>
<span id="L81" class="ln">    81&nbsp;&nbsp;</span><span class="comment">// enclosingTarget returns the innermost enclosing labeled</span>
<span id="L82" class="ln">    82&nbsp;&nbsp;</span><span class="comment">// statement with the given label name, or nil.</span>
<span id="L83" class="ln">    83&nbsp;&nbsp;</span>func (b *block) enclosingTarget(name string) *ast.LabeledStmt {
<span id="L84" class="ln">    84&nbsp;&nbsp;</span>	for s := b; s != nil; s = s.parent {
<span id="L85" class="ln">    85&nbsp;&nbsp;</span>		if t := s.lstmt; t != nil &amp;&amp; t.Label.Name == name {
<span id="L86" class="ln">    86&nbsp;&nbsp;</span>			return t
<span id="L87" class="ln">    87&nbsp;&nbsp;</span>		}
<span id="L88" class="ln">    88&nbsp;&nbsp;</span>	}
<span id="L89" class="ln">    89&nbsp;&nbsp;</span>	return nil
<span id="L90" class="ln">    90&nbsp;&nbsp;</span>}
<span id="L91" class="ln">    91&nbsp;&nbsp;</span>
<span id="L92" class="ln">    92&nbsp;&nbsp;</span><span class="comment">// blockBranches processes a block&#39;s statement list and returns the set of outgoing forward jumps.</span>
<span id="L93" class="ln">    93&nbsp;&nbsp;</span><span class="comment">// all is the scope of all declared labels, parent the set of labels declared in the immediately</span>
<span id="L94" class="ln">    94&nbsp;&nbsp;</span><span class="comment">// enclosing block, and lstmt is the labeled statement this block is associated with (or nil).</span>
<span id="L95" class="ln">    95&nbsp;&nbsp;</span>func (check *Checker) blockBranches(all *Scope, parent *block, lstmt *ast.LabeledStmt, list []ast.Stmt) []*ast.BranchStmt {
<span id="L96" class="ln">    96&nbsp;&nbsp;</span>	b := &amp;block{parent: parent, lstmt: lstmt}
<span id="L97" class="ln">    97&nbsp;&nbsp;</span>
<span id="L98" class="ln">    98&nbsp;&nbsp;</span>	var (
<span id="L99" class="ln">    99&nbsp;&nbsp;</span>		varDeclPos         token.Pos
<span id="L100" class="ln">   100&nbsp;&nbsp;</span>		fwdJumps, badJumps []*ast.BranchStmt
<span id="L101" class="ln">   101&nbsp;&nbsp;</span>	)
<span id="L102" class="ln">   102&nbsp;&nbsp;</span>
<span id="L103" class="ln">   103&nbsp;&nbsp;</span>	<span class="comment">// All forward jumps jumping over a variable declaration are possibly</span>
<span id="L104" class="ln">   104&nbsp;&nbsp;</span>	<span class="comment">// invalid (they may still jump out of the block and be ok).</span>
<span id="L105" class="ln">   105&nbsp;&nbsp;</span>	<span class="comment">// recordVarDecl records them for the given position.</span>
<span id="L106" class="ln">   106&nbsp;&nbsp;</span>	recordVarDecl := func(pos token.Pos) {
<span id="L107" class="ln">   107&nbsp;&nbsp;</span>		varDeclPos = pos
<span id="L108" class="ln">   108&nbsp;&nbsp;</span>		badJumps = append(badJumps[:0], fwdJumps...) <span class="comment">// copy fwdJumps to badJumps</span>
<span id="L109" class="ln">   109&nbsp;&nbsp;</span>	}
<span id="L110" class="ln">   110&nbsp;&nbsp;</span>
<span id="L111" class="ln">   111&nbsp;&nbsp;</span>	jumpsOverVarDecl := func(jmp *ast.BranchStmt) bool {
<span id="L112" class="ln">   112&nbsp;&nbsp;</span>		if varDeclPos.IsValid() {
<span id="L113" class="ln">   113&nbsp;&nbsp;</span>			for _, bad := range badJumps {
<span id="L114" class="ln">   114&nbsp;&nbsp;</span>				if jmp == bad {
<span id="L115" class="ln">   115&nbsp;&nbsp;</span>					return true
<span id="L116" class="ln">   116&nbsp;&nbsp;</span>				}
<span id="L117" class="ln">   117&nbsp;&nbsp;</span>			}
<span id="L118" class="ln">   118&nbsp;&nbsp;</span>		}
<span id="L119" class="ln">   119&nbsp;&nbsp;</span>		return false
<span id="L120" class="ln">   120&nbsp;&nbsp;</span>	}
<span id="L121" class="ln">   121&nbsp;&nbsp;</span>
<span id="L122" class="ln">   122&nbsp;&nbsp;</span>	blockBranches := func(lstmt *ast.LabeledStmt, list []ast.Stmt) {
<span id="L123" class="ln">   123&nbsp;&nbsp;</span>		<span class="comment">// Unresolved forward jumps inside the nested block</span>
<span id="L124" class="ln">   124&nbsp;&nbsp;</span>		<span class="comment">// become forward jumps in the current block.</span>
<span id="L125" class="ln">   125&nbsp;&nbsp;</span>		fwdJumps = append(fwdJumps, check.blockBranches(all, b, lstmt, list)...)
<span id="L126" class="ln">   126&nbsp;&nbsp;</span>	}
<span id="L127" class="ln">   127&nbsp;&nbsp;</span>
<span id="L128" class="ln">   128&nbsp;&nbsp;</span>	var stmtBranches func(ast.Stmt)
<span id="L129" class="ln">   129&nbsp;&nbsp;</span>	stmtBranches = func(s ast.Stmt) {
<span id="L130" class="ln">   130&nbsp;&nbsp;</span>		switch s := s.(type) {
<span id="L131" class="ln">   131&nbsp;&nbsp;</span>		case *ast.DeclStmt:
<span id="L132" class="ln">   132&nbsp;&nbsp;</span>			if d, _ := s.Decl.(*ast.GenDecl); d != nil &amp;&amp; d.Tok == token.VAR {
<span id="L133" class="ln">   133&nbsp;&nbsp;</span>				recordVarDecl(d.Pos())
<span id="L134" class="ln">   134&nbsp;&nbsp;</span>			}
<span id="L135" class="ln">   135&nbsp;&nbsp;</span>
<span id="L136" class="ln">   136&nbsp;&nbsp;</span>		case *ast.LabeledStmt:
<span id="L137" class="ln">   137&nbsp;&nbsp;</span>			<span class="comment">// declare non-blank label</span>
<span id="L138" class="ln">   138&nbsp;&nbsp;</span>			if name := s.Label.Name; name != &#34;_&#34; {
<span id="L139" class="ln">   139&nbsp;&nbsp;</span>				lbl := NewLabel(s.Label.Pos(), check.pkg, name)
<span id="L140" class="ln">   140&nbsp;&nbsp;</span>				if alt := all.Insert(lbl); alt != nil {
<span id="L141" class="ln">   141&nbsp;&nbsp;</span>					check.softErrorf(lbl, DuplicateLabel, &#34;label %s already declared&#34;, name)
<span id="L142" class="ln">   142&nbsp;&nbsp;</span>					check.reportAltDecl(alt)
<span id="L143" class="ln">   143&nbsp;&nbsp;</span>					<span class="comment">// ok to continue</span>
<span id="L144" class="ln">   144&nbsp;&nbsp;</span>				} else {
<span id="L145" class="ln">   145&nbsp;&nbsp;</span>					b.insert(s)
<span id="L146" class="ln">   146&nbsp;&nbsp;</span>					check.recordDef(s.Label, lbl)
<span id="L147" class="ln">   147&nbsp;&nbsp;</span>				}
<span id="L148" class="ln">   148&nbsp;&nbsp;</span>				<span class="comment">// resolve matching forward jumps and remove them from fwdJumps</span>
<span id="L149" class="ln">   149&nbsp;&nbsp;</span>				i := 0
<span id="L150" class="ln">   150&nbsp;&nbsp;</span>				for _, jmp := range fwdJumps {
<span id="L151" class="ln">   151&nbsp;&nbsp;</span>					if jmp.Label.Name == name {
<span id="L152" class="ln">   152&nbsp;&nbsp;</span>						<span class="comment">// match</span>
<span id="L153" class="ln">   153&nbsp;&nbsp;</span>						lbl.used = true
<span id="L154" class="ln">   154&nbsp;&nbsp;</span>						check.recordUse(jmp.Label, lbl)
<span id="L155" class="ln">   155&nbsp;&nbsp;</span>						if jumpsOverVarDecl(jmp) {
<span id="L156" class="ln">   156&nbsp;&nbsp;</span>							check.softErrorf(
<span id="L157" class="ln">   157&nbsp;&nbsp;</span>								jmp.Label,
<span id="L158" class="ln">   158&nbsp;&nbsp;</span>								JumpOverDecl,
<span id="L159" class="ln">   159&nbsp;&nbsp;</span>								&#34;goto %s jumps over variable declaration at line %d&#34;,
<span id="L160" class="ln">   160&nbsp;&nbsp;</span>								name,
<span id="L161" class="ln">   161&nbsp;&nbsp;</span>								check.fset.Position(varDeclPos).Line,
<span id="L162" class="ln">   162&nbsp;&nbsp;</span>							)
<span id="L163" class="ln">   163&nbsp;&nbsp;</span>							<span class="comment">// ok to continue</span>
<span id="L164" class="ln">   164&nbsp;&nbsp;</span>						}
<span id="L165" class="ln">   165&nbsp;&nbsp;</span>					} else {
<span id="L166" class="ln">   166&nbsp;&nbsp;</span>						<span class="comment">// no match - record new forward jump</span>
<span id="L167" class="ln">   167&nbsp;&nbsp;</span>						fwdJumps[i] = jmp
<span id="L168" class="ln">   168&nbsp;&nbsp;</span>						i++
<span id="L169" class="ln">   169&nbsp;&nbsp;</span>					}
<span id="L170" class="ln">   170&nbsp;&nbsp;</span>				}
<span id="L171" class="ln">   171&nbsp;&nbsp;</span>				fwdJumps = fwdJumps[:i]
<span id="L172" class="ln">   172&nbsp;&nbsp;</span>				lstmt = s
<span id="L173" class="ln">   173&nbsp;&nbsp;</span>			}
<span id="L174" class="ln">   174&nbsp;&nbsp;</span>			stmtBranches(s.Stmt)
<span id="L175" class="ln">   175&nbsp;&nbsp;</span>
<span id="L176" class="ln">   176&nbsp;&nbsp;</span>		case *ast.BranchStmt:
<span id="L177" class="ln">   177&nbsp;&nbsp;</span>			if s.Label == nil {
<span id="L178" class="ln">   178&nbsp;&nbsp;</span>				return <span class="comment">// checked in 1st pass (check.stmt)</span>
<span id="L179" class="ln">   179&nbsp;&nbsp;</span>			}
<span id="L180" class="ln">   180&nbsp;&nbsp;</span>
<span id="L181" class="ln">   181&nbsp;&nbsp;</span>			<span class="comment">// determine and validate target</span>
<span id="L182" class="ln">   182&nbsp;&nbsp;</span>			name := s.Label.Name
<span id="L183" class="ln">   183&nbsp;&nbsp;</span>			switch s.Tok {
<span id="L184" class="ln">   184&nbsp;&nbsp;</span>			case token.BREAK:
<span id="L185" class="ln">   185&nbsp;&nbsp;</span>				<span class="comment">// spec: &#34;If there is a label, it must be that of an enclosing</span>
<span id="L186" class="ln">   186&nbsp;&nbsp;</span>				<span class="comment">// &#34;for&#34;, &#34;switch&#34;, or &#34;select&#34; statement, and that is the one</span>
<span id="L187" class="ln">   187&nbsp;&nbsp;</span>				<span class="comment">// whose execution terminates.&#34;</span>
<span id="L188" class="ln">   188&nbsp;&nbsp;</span>				valid := false
<span id="L189" class="ln">   189&nbsp;&nbsp;</span>				if t := b.enclosingTarget(name); t != nil {
<span id="L190" class="ln">   190&nbsp;&nbsp;</span>					switch t.Stmt.(type) {
<span id="L191" class="ln">   191&nbsp;&nbsp;</span>					case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.ForStmt, *ast.RangeStmt:
<span id="L192" class="ln">   192&nbsp;&nbsp;</span>						valid = true
<span id="L193" class="ln">   193&nbsp;&nbsp;</span>					}
<span id="L194" class="ln">   194&nbsp;&nbsp;</span>				}
<span id="L195" class="ln">   195&nbsp;&nbsp;</span>				if !valid {
<span id="L196" class="ln">   196&nbsp;&nbsp;</span>					check.errorf(s.Label, MisplacedLabel, &#34;invalid break label %s&#34;, name)
<span id="L197" class="ln">   197&nbsp;&nbsp;</span>					return
<span id="L198" class="ln">   198&nbsp;&nbsp;</span>				}
<span id="L199" class="ln">   199&nbsp;&nbsp;</span>
<span id="L200" class="ln">   200&nbsp;&nbsp;</span>			case token.CONTINUE:
<span id="L201" class="ln">   201&nbsp;&nbsp;</span>				<span class="comment">// spec: &#34;If there is a label, it must be that of an enclosing</span>
<span id="L202" class="ln">   202&nbsp;&nbsp;</span>				<span class="comment">// &#34;for&#34; statement, and that is the one whose execution advances.&#34;</span>
<span id="L203" class="ln">   203&nbsp;&nbsp;</span>				valid := false
<span id="L204" class="ln">   204&nbsp;&nbsp;</span>				if t := b.enclosingTarget(name); t != nil {
<span id="L205" class="ln">   205&nbsp;&nbsp;</span>					switch t.Stmt.(type) {
<span id="L206" class="ln">   206&nbsp;&nbsp;</span>					case *ast.ForStmt, *ast.RangeStmt:
<span id="L207" class="ln">   207&nbsp;&nbsp;</span>						valid = true
<span id="L208" class="ln">   208&nbsp;&nbsp;</span>					}
<span id="L209" class="ln">   209&nbsp;&nbsp;</span>				}
<span id="L210" class="ln">   210&nbsp;&nbsp;</span>				if !valid {
<span id="L211" class="ln">   211&nbsp;&nbsp;</span>					check.errorf(s.Label, MisplacedLabel, &#34;invalid continue label %s&#34;, name)
<span id="L212" class="ln">   212&nbsp;&nbsp;</span>					return
<span id="L213" class="ln">   213&nbsp;&nbsp;</span>				}
<span id="L214" class="ln">   214&nbsp;&nbsp;</span>
<span id="L215" class="ln">   215&nbsp;&nbsp;</span>			case token.GOTO:
<span id="L216" class="ln">   216&nbsp;&nbsp;</span>				if b.gotoTarget(name) == nil {
<span id="L217" class="ln">   217&nbsp;&nbsp;</span>					<span class="comment">// label may be declared later - add branch to forward jumps</span>
<span id="L218" class="ln">   218&nbsp;&nbsp;</span>					fwdJumps = append(fwdJumps, s)
<span id="L219" class="ln">   219&nbsp;&nbsp;</span>					return
<span id="L220" class="ln">   220&nbsp;&nbsp;</span>				}
<span id="L221" class="ln">   221&nbsp;&nbsp;</span>
<span id="L222" class="ln">   222&nbsp;&nbsp;</span>			default:
<span id="L223" class="ln">   223&nbsp;&nbsp;</span>				check.errorf(s, InvalidSyntaxTree, &#34;branch statement: %s %s&#34;, s.Tok, name)
<span id="L224" class="ln">   224&nbsp;&nbsp;</span>				return
<span id="L225" class="ln">   225&nbsp;&nbsp;</span>			}
<span id="L226" class="ln">   226&nbsp;&nbsp;</span>
<span id="L227" class="ln">   227&nbsp;&nbsp;</span>			<span class="comment">// record label use</span>
<span id="L228" class="ln">   228&nbsp;&nbsp;</span>			obj := all.Lookup(name)
<span id="L229" class="ln">   229&nbsp;&nbsp;</span>			obj.(*Label).used = true
<span id="L230" class="ln">   230&nbsp;&nbsp;</span>			check.recordUse(s.Label, obj)
<span id="L231" class="ln">   231&nbsp;&nbsp;</span>
<span id="L232" class="ln">   232&nbsp;&nbsp;</span>		case *ast.AssignStmt:
<span id="L233" class="ln">   233&nbsp;&nbsp;</span>			if s.Tok == token.DEFINE {
<span id="L234" class="ln">   234&nbsp;&nbsp;</span>				recordVarDecl(s.Pos())
<span id="L235" class="ln">   235&nbsp;&nbsp;</span>			}
<span id="L236" class="ln">   236&nbsp;&nbsp;</span>
<span id="L237" class="ln">   237&nbsp;&nbsp;</span>		case *ast.BlockStmt:
<span id="L238" class="ln">   238&nbsp;&nbsp;</span>			blockBranches(lstmt, s.List)
<span id="L239" class="ln">   239&nbsp;&nbsp;</span>
<span id="L240" class="ln">   240&nbsp;&nbsp;</span>		case *ast.IfStmt:
<span id="L241" class="ln">   241&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L242" class="ln">   242&nbsp;&nbsp;</span>			if s.Else != nil {
<span id="L243" class="ln">   243&nbsp;&nbsp;</span>				stmtBranches(s.Else)
<span id="L244" class="ln">   244&nbsp;&nbsp;</span>			}
<span id="L245" class="ln">   245&nbsp;&nbsp;</span>
<span id="L246" class="ln">   246&nbsp;&nbsp;</span>		case *ast.CaseClause:
<span id="L247" class="ln">   247&nbsp;&nbsp;</span>			blockBranches(nil, s.Body)
<span id="L248" class="ln">   248&nbsp;&nbsp;</span>
<span id="L249" class="ln">   249&nbsp;&nbsp;</span>		case *ast.SwitchStmt:
<span id="L250" class="ln">   250&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L251" class="ln">   251&nbsp;&nbsp;</span>
<span id="L252" class="ln">   252&nbsp;&nbsp;</span>		case *ast.TypeSwitchStmt:
<span id="L253" class="ln">   253&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L254" class="ln">   254&nbsp;&nbsp;</span>
<span id="L255" class="ln">   255&nbsp;&nbsp;</span>		case *ast.CommClause:
<span id="L256" class="ln">   256&nbsp;&nbsp;</span>			blockBranches(nil, s.Body)
<span id="L257" class="ln">   257&nbsp;&nbsp;</span>
<span id="L258" class="ln">   258&nbsp;&nbsp;</span>		case *ast.SelectStmt:
<span id="L259" class="ln">   259&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L260" class="ln">   260&nbsp;&nbsp;</span>
<span id="L261" class="ln">   261&nbsp;&nbsp;</span>		case *ast.ForStmt:
<span id="L262" class="ln">   262&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L263" class="ln">   263&nbsp;&nbsp;</span>
<span id="L264" class="ln">   264&nbsp;&nbsp;</span>		case *ast.RangeStmt:
<span id="L265" class="ln">   265&nbsp;&nbsp;</span>			stmtBranches(s.Body)
<span id="L266" class="ln">   266&nbsp;&nbsp;</span>		}
<span id="L267" class="ln">   267&nbsp;&nbsp;</span>	}
<span id="L268" class="ln">   268&nbsp;&nbsp;</span>
<span id="L269" class="ln">   269&nbsp;&nbsp;</span>	for _, s := range list {
<span id="L270" class="ln">   270&nbsp;&nbsp;</span>		stmtBranches(s)
<span id="L271" class="ln">   271&nbsp;&nbsp;</span>	}
<span id="L272" class="ln">   272&nbsp;&nbsp;</span>
<span id="L273" class="ln">   273&nbsp;&nbsp;</span>	return fwdJumps
<span id="L274" class="ln">   274&nbsp;&nbsp;</span>}
<span id="L275" class="ln">   275&nbsp;&nbsp;</span>
</pre><p><a href="labels.go?m=text">View as plain text</a></p>

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
