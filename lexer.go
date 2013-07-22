package dateparser

import (
	"bufio"
	"io"
	"unicode"
)

const (
	_ST_NONE = iota
	_ST_LETTER
	_ST_DIGIT
	_ST_LETTER_PERIOD
	_ST_DIGIT_PERIOD
)

type lexer struct {
	stream     *bufio.Reader
	charStack  []rune
	tokenStack [][]rune
	eof        bool
}

func newLexer(stream io.Reader) (lex *lexer) {
	return &lexer{
		stream:     bufio.NewReader(stream),
		charStack:  nil,
		tokenStack: nil,
		eof:        false,
	}
}

func (lex *lexer) lex() (tokenStr string, err error) {
	var token []rune
	
	if len(lex.tokenStack) > 0 {
		token = lex.tokenStack[0]
		lex.tokenStack = lex.tokenStack[1:]
		return string(encode(token)), nil
	}

	state := _ST_NONE
	seenLetters := false

loop:
	for !lex.eof {
		var nextChar rune = '\x00'

		if len(lex.charStack) > 0 {
			nextChar = lex.charStack[0]
			lex.charStack = lex.charStack[1:]

		} else {
			for nextChar == '\x00' {
				nextChar, _, err = lex.stream.ReadRune()

				if err == io.EOF {
					lex.eof = true
					break loop
				}

				if err != nil {
					return "", err
				}
			}
		}

		switch state {
		case _ST_NONE:
			token = []rune{nextChar}

			switch {
			case unicode.IsLetter(nextChar):
				state = _ST_LETTER
			case unicode.IsDigit(nextChar):
				state = _ST_DIGIT
			case unicode.IsSpace(nextChar):
				token = []rune{' '}
				break loop
			default:
				break loop
			}

		case _ST_LETTER:
			seenLetters = true

			switch {
			case nextChar == '.':
				state = _ST_LETTER_PERIOD
				fallthrough
			case unicode.IsLetter(nextChar):
				token = append(token, nextChar)
			default:
				lex.charStack = append(lex.charStack, nextChar)
				break loop
			}

		case _ST_DIGIT:
			switch {
			case nextChar == '.':
				state = _ST_DIGIT_PERIOD
				fallthrough
			case unicode.IsDigit(nextChar):
				token = append(token, nextChar)
			default:
				lex.charStack = append(lex.charStack, nextChar)
				break loop
			}

		case _ST_LETTER_PERIOD:
			seenLetters = true

			switch {
			case unicode.IsDigit(nextChar) && token[len(token)-1] == '.':
				state = _ST_DIGIT_PERIOD
				fallthrough
			case nextChar == '.' || unicode.IsLetter(nextChar):
				token = append(token, nextChar)
			default:
				lex.charStack = append(lex.charStack, nextChar)
				break loop
			}

		case _ST_DIGIT_PERIOD:
			switch {
			case unicode.IsLetter(nextChar) && token[len(token)-1] == '.':
				state = _ST_LETTER_PERIOD
				fallthrough
			case nextChar == '.' || unicode.IsDigit(nextChar):
				token = append(token, nextChar)
			default:
				lex.charStack = append(lex.charStack, nextChar)
				break loop
			}
		}
	}

	if (state == _ST_LETTER_PERIOD || state == _ST_DIGIT_PERIOD) &&
			(seenLetters || has2Periods(token) || token[len(token)-1] == '.') {

		lastIndex := 0
		var newToken []rune
		
		for i, char := range token {
			if char == '.' {
				part := token[lastIndex:i]
				
				if lastIndex == 0 {
					newToken = part
				} else {
					lex.tokenStack = append(lex.tokenStack, []rune{'.'})
					if len(part) > 0 {
						lex.tokenStack = append(lex.tokenStack, part)
					}
				}
				
				lastIndex = i + 1
			}
		}
		
		if lastIndex < len(token) {
			lex.tokenStack = append(lex.tokenStack, []rune{'.'})
			lex.tokenStack = append(lex.tokenStack, token[lastIndex:])
		}
		
		token = newToken
	}
	
	return string(encode(token)), nil
}

func (lex *lexer) lexAll() (tokens []string, err error) {
	for {
		token, err := lex.lex()
		if err != nil {
			return nil, err
		}
		
		if token == "" {
			break
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens, nil
}
