package golib

/*

A library of useful functions to work with regular expressions

TODO: Apply this to the Server.Module.Controller's prioritizing of more specific patterns over less specific ones.

*/

import (
	"bytes"
)
// ref: http://savage.net.au/Ron/html/graphviz2.marpa/Lexing.and.Parsing.Overview.html
// ref: https://swtch.com/~rsc/regexp/regexp3.html

type LexicalToken struct {
	Value		string
	IsComplete	bool
	IsEncapsulation	bool
}

type TokenType int

const (
        TOKEN_TYPE_LITERAL		TokenType = iota
        TOKEN_TYPE_ENCAPSULATION
	TOKEN_TYPE_MATCHER
	TOKEN_TYPE_ALTERATION
	TOKEN_TYPE_COUNT
	TOKEN_TYPE_BEGIN
	TOKEN_TYPE_END
)

// Ordered Tree data structure
// TODO: Can/should we generalize this to support other types of ordered trees such as DOM/XML/JSON/CSV?
type PatternTokenSet []PatternToken

type PatternToken struct {
	Token		string
	Type		TokenType
	Valid		bool
	Children	[]PatternTokenSet
}

// Score a pattern for "specificity": a count of literal (non-matcher/wildcard) characters
// ref: https://cs.stackexchange.com/questions/10786/how-to-find-specificity-of-a-regex-match
// ref: https://golang.org/pkg/regexp/#Compile
// ref: https://stackoverflow.com/questions/3978438/dfa-vs-nfa-engines-what-is-the-difference-in-their-capabilities-and-limitations
// ref: https://github.com/timtadh/lexmachine
func TokenizePattern(pattern string) *[]PatternToken {

	tokens := make([]PatternToken, 0)

	// Special Characters which affect the token boundary
	encapsulationOpeners := "{[(<"
	encapsulationClosers := ">)]}"
	escape := '\\'
	greed := "+*"
	matchers := ".^$|"
	escapedMatchers := "dDsSp"
	// ref: https://golang.org/pkg/regexp/syntax/

	// Tokenize the pattern
	escaped := false
	encapsulationDepth := 0
	var tokenBuffer bytes.Buffer
	for (i := 0; i < len(pattern); i++) {
		if escaped {

		} else {
			if '\\' == pattern[i] {
				escaped = true
			} else {
				unescaped.WriteString(pattern[i])
			}
		}
	}
}

// Check whether the supplied string is a valid regexp pattern
func PatternIsValid(pattern string) bool {
	re, err := regexp.Compile(pattern)
	return nil == err
}

func PatternSpecificity(pattern string) int {
}

/*
	score := 0

	// Score one point for each escaped literal character
	var unescaped bytes.Buffer
	for (i := 0; i < len(pattern); i++) {
		if escaped {
			escaped = false
			// TODO: If the escaped character was not a special matcher, then we score it
			if ! strings.Contains(escapedMatchers, pattern[i]) {
				score++
			}
		} else {
			if escape == pattern[i] {
				escaped = true
			} else {
				unescaped.WriteString(pattern[i])
			}
		}
	}

	// TODO: Now scrub out all the special matcher/wildcard chars
	unescapedStr := string(unescaped)
*/

