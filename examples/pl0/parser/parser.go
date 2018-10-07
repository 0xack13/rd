package parser

import (
	"fmt"
	"regexp"

	"github.com/shivamMg/rd"
	. "github.com/shivamMg/rd/examples/pl0/tokens"
)

func Parse(tokens []rd.Token) (*rd.Tree, error) {
	b := rd.NewBuilder(tokens)
	if ok := Program(b); !ok {
		return nil, b.Err()
	}
	return b.Tree(), nil
}

func Program(b *rd.Builder) (ok bool) {
	b.Enter("Program")
	defer b.Exit(&ok)

	return Block(b) && b.Match(Period)
}

func Block(b *rd.Builder) (ok bool) {
	b.Enter("Block")
	defer b.Exit(&ok)

	if b.Match(Const) {
		for {
			if Ident(b) && b.Match(Equal) && Number(b) {
				if b.Match(Comma) {
					continue
				}
				if b.Match(Semicolon) {
					break
				}
			}
			return false
		}
	}
	if b.Match(Var) {
		for {
			if Ident(b) {
				if b.Match(Comma) {
					continue
				}
				if b.Match(Semicolon) {
					break
				}
			}
			return false
		}
	}
	for b.Match(Procedure) {
		if Ident(b) && b.Match(Semicolon) && Block(b) && b.Match(Semicolon) {
			continue
		}
		return false
	}
	return Statement(b)
}

func Statement(b *rd.Builder) (ok bool) {
	b.Enter("Statement")
	defer b.Exit(&ok)

	switch {
	case Ident(b):
		return b.Match(Assignment) && Expression(b)
	case b.Match(Exclam):
		return Expression(b)
	case b.Match(Ques):
		return Ident(b)
	case b.Match(Call):
		return Ident(b)
	case b.Match(Begin):
		for {
			if Statement(b) {
				if b.Match(Semicolon) {
					continue
				}
				break
			}
			return false
		}
		return b.Match(End)
	case b.Match(If):
		return Condition(b) && b.Match(Then) && Statement(b)
	case b.Match(While):
		return Condition(b) && b.Match(Do) && Statement(b)
	}
	return false
}

func Condition(b *rd.Builder) (ok bool) {
	b.Enter("Condition")
	defer b.Exit(&ok)

	switch {
	case b.Match(Odd):
		return Expression(b)
	case Expression(b):
		if b.Match(Equal) || b.Match(Hash) || b.Match(LT) || b.Match(LTE) || b.Match(GT) || b.Match(GTE) {
			return Expression(b)
		}
		return false
	}
	return false
}

func Expression(b *rd.Builder) (ok bool) {
	b.Enter("Expression")
	defer b.Exit(&ok)

	if b.Match(Plus) || b.Match(Minus) {
	}
	for {
		if Term(b) {
			if b.Match(Plus) || b.Match(Minus) {
				continue
			}
			break
		}
		return false
	}
	return true
}

func Term(b *rd.Builder) (ok bool) {
	b.Enter("Term")
	defer b.Exit(&ok)

	for {
		if Factor(b) {
			if b.Match(Mul) || b.Match(Div) {
				continue
			}
			break
		}
		return false
	}
	return true
}

func Factor(b *rd.Builder) (ok bool) {
	b.Enter("Factor")
	defer b.Exit(&ok)

	if Ident(b) {
		return true
	}
	if Number(b) {
		return true
	}
	return b.Match(OpenParen) && Expression(b) && b.Match(CloseParen)
}

func Ident(b *rd.Builder) (ok bool) {
	b.Enter("Ident")
	defer b.Exit(&ok)

	token, ok := b.Next()
	if !ok {
		return false
	}
	if _, ok := token.(Token); ok {
		b.Reset()
		return false
	}
	if ok, _ := regexp.MatchString(`[[:alpha:]]`, fmt.Sprint(token)); ok {
		b.Add(token)
		return true
	}
	b.Reset()
	return false
}

func Number(b *rd.Builder) (ok bool) {
	b.Enter("Number")
	defer b.Exit(&ok)

	token, ok := b.Next()
	if !ok {
		return false
	}
	if _, ok := token.(Token); ok {
		b.Reset()
		return false
	}
	if ok, _ := regexp.MatchString(`[[:digit:]]`, fmt.Sprint(token)); ok {
		b.Add(token)
		return true
	}
	b.Reset()
	return false
}