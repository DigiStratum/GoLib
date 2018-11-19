package golib

/*
A general purpose lexer

This is a learning experiment. Maybe it will turn into something useful? The ones out there seem to
be more complex than we need for our purposes.

ref: https://github.com/AdamColton/parlex/blob/master/lexer/simplelexer/lexer.go
ref: http://matt.might.net/articles/grammars-bnf-ebnf/
ref: https://github.com/goccmack/gocc

ref: http://matt.might.net/articles/implementation-of-regular-expression-matching-in-scheme-with-derivatives/

ref: https://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html

-----------

BNF Grammar of Regular Expressions
Following the precedence rules given previously, a BNF grammar for Perl-style regular expressions can be constructed as follows.
ref: https://stackoverflow.com/questions/265457/regex-grammar
ref: http://www.cs.sfu.ca/~cameron/Teaching/384/99-3/regexp-plg.html
ref: https://web.archive.org/web/20090129224504/http://faqts.com/knowledge_base/view.phtml/aid/25718/fid/200

<RE>	::=	<union> | <simple-RE>
<union>	::=	<RE> "|" <simple-RE>
<simple-RE>	::=	<concatenation> | <basic-RE>
<concatenation>	::=	<simple-RE> <basic-RE>
<basic-RE>	::=	<star> | <plus> | <elementary-RE>
<star>	::=	<elementary-RE> "*"
<plus>	::=	<elementary-RE> "+"
<elementary-RE>	::=	<group> | <any> | <eos> | <char> | <set>
<group>	::=	"(" <RE> ")"
<any>	::=	"."
<eos>	::=	"$"
<char>	::=	any non metacharacter | "\" metacharacter
<set>	::=	<positive-set> | <negative-set>
<positive-set>	::=	"[" <set-items> "]"
<negative-set>	::=	"[^" <set-items> "]"
<set-items>	::=	<set-item> | <set-item> <set-items>
<set-items>	::=	<range> | <char>
<range>	::=	<char> "-" <char>

*/




import (

)

type Lexeme struct {
	name	string
	value	string
}

type LexemeSpec struct {
	name	string
	pattern	string
}

type Lexemes []*Lexeme
type LexemeSpecs []*LexemeSpec

type Lexer struct {
	lexemeSpecs	[]LexemeSpec
	processed	bool
	lexemeSequence	*Lexemes
	content		*string
}

func newLexemes() *Lexemes {
	lexemeSequence := make([]*Lexeme, 0)
	return &lexemeSequence
}

func NewLexer(content *string) *Lexer {
	lexer := Lexer{
		parser_encapsulations:	make([]*encapsualtion, 0),
		lexemeSequence:		newLexemes(),
		content:		content,
	}
	return &lexer
}

// Add a LexemeSpec for this Lexer
func (lex *Lexer) AddLexemeSpec(name string, pattern string) {
	lex.parser_lexemeSpecs = append(
		lex.parser_lexemeSpecs,
		&LexemeSpec{ name: name, pattern: pattern },
	)
}

func newLexemeSpec(name string, pattern string) *LexemeSpec {
	return &LexemeSpec{
		name:		name,
		pattern:	pattern,
	}
}

func (lex *Lexer) Lexemeize() {
	if lex.processed { return }
	// TODO: Do some parsing here!
	lex.processed := true
}

