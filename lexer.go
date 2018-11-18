package golib

import (

)

type encapsulation struct {
	opener	string
	closer	string
}

type Token struct {
	name	string
	value	string
}

type TokenSequence []*Token

type Lexer struct {
	encapsulations		[]*encapsulation
}

func NewLexer() *Lexer {
	lexer := Lexer{
		encapsulations:	make([]*encapsualtion, 0)
	}
	return &lexer
}

func (lex *Lexer) AddEncapsulation(opener string, closer string) {
	lex.encapsulations = append(
		lex.encapsulations,
		newEncapsulation(opener, closer),
	)
}

func newEncapsualtion(opener string, closer string) *encapsulation {
	return &encapsulation{
		opener:	opener,
		closer: closer,
	}
}

func (lex *Lexer) Tokenize(content string) *TokenSequence {
	tokenSequence := make([]*Token, 0)
	// TODO: Do some parsing here!
	return &tokenSequence
}

