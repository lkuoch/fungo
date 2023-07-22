package lexer

import "fungo/token"

// Needs to support peeking the next character
type Lexer struct {
	input        string
	position     int  // current char position in input
	readPosition int  // current reading position in input (after current char)
	char         byte // current char
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input, position: 0, readPosition: 0, char: 0}
	lexer.readChar()

	return lexer
}

// Read the next character and advance position in the `input` string
// Only support ASCII characters
func (l *Lexer) readChar() {
	// Assign character if exists
	if l.readPosition >= len(l.input) {
		// In ASCII, the `0th` byte represents null
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func createNewToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhiteSpace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var newToken token.Token

	l.skipWhiteSpace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			prevCh := l.char
			l.readChar()
			newToken = token.Token{
				Type:    token.EQ,
				Literal: string(prevCh) + string(l.char),
			}
		} else {
			newToken = createNewToken(token.ASSIGN, l.char)
		}
	case '!':
		if l.peekChar() == '=' {
			prevCh := l.char
			l.readChar()
			newToken = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(prevCh) + string(l.char),
			}
		} else {
			newToken = createNewToken(token.BANG, l.char)
		}
	case '+':
		newToken = createNewToken(token.PLUS, l.char)
	case '-':
		newToken = createNewToken(token.MINUS, l.char)
	case '/':
		newToken = createNewToken(token.SLASH, l.char)
	case '*':
		newToken = createNewToken(token.ASTERISK, l.char)

	case '<':
		newToken = createNewToken(token.LT, l.char)
	case '>':
		newToken = createNewToken(token.GT, l.char)

	case ',':
		newToken = createNewToken(token.COMMA, l.char)
	case ';':
		newToken = createNewToken(token.SEMICOLON, l.char)

	case '(':
		newToken = createNewToken(token.LPAREN, l.char)
	case ')':
		newToken = createNewToken(token.RPAREN, l.char)
	case '{':
		newToken = createNewToken(token.LBRACE, l.char)
	case '}':
		newToken = createNewToken(token.RBRACE, l.char)

	// ASCII NULL character
	case 0:
		newToken.Literal = ""
		newToken.Type = token.EOF
	default:
		if isLetter(l.char) {
			newToken.Literal = l.readIdentifier()
			newToken.Type = token.LookupIdent(newToken.Literal)
			return newToken
		} else if isDigit(l.char) {
			newToken.Type = token.INT
			newToken.Literal = l.readNumber()
			return newToken
		} else {
			newToken = createNewToken(token.ILLEGAL, l.char)
		}
	}

	// After reading identifier, shift lexer to next place
	l.readChar()

	return newToken
}
