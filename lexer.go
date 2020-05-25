package easylexer

import (
	"io"
	"strings"
)

// Определены значения по умолчанию для свойств Lexer в качестве значения пакета.
var (
	DefaultWhitespace = NewPatternTokenType(-1, []string{" ", "\t", "\r", "\n"})

	DefaultTokenTypes = []TokenType{
		NewRegexpTokenType(IDENT, `[a-zA-Z_][a-zA-Z0-9_]*`),
		NewRegexpTokenType(NUMBER, `[0-9]+(?:\.[0-9]+)?`),
		NewRegexpTokenType(STRING, `\"([^"]*)\"`),
		NewRegexpTokenType(OTHER, `.`),
	}
)

/*
Лексический анализатор.
Пробел - это TokenType для пропуска символов, таких как пробелы.
Значением по умолчанию является симплексер. DefaultWhitespace.
Не пропускает никаких символов, если пробел равен нулю.
TokenTypes - это массив TokenType.
Lexer будет последовательно проверять TokenTypes и возвращать первый соответствующий токен.
По умолчанию используется симплексер. DefaultTokenTypes.
Lexer никогда не будет использовать его, даже если добавит TokenType после OTHER.
Потому что ДРУГОЙ примет любой отдельный символ.
*/

type Lexer struct {
	reader     io.Reader
	buf        string
	loadedLine string
	nextPos    Position
	Whitespace TokenType
	TokenTypes []TokenType
}

// Создание нового Lexer.
func NewLexer(reader io.Reader) *Lexer {
	l := new(Lexer)
	l.reader = reader

	l.Whitespace = DefaultWhitespace
	l.TokenTypes = DefaultTokenTypes

	return l
}

func (l *Lexer) readBufIfNeed() {
	if len(l.buf) < 1024 {
		buf := make([]byte, 2048)
		l.reader.Read(buf)
		l.buf += strings.TrimRight(string(buf), "\x00")
	}
}

func (l *Lexer) consumeBuffer(t *Token) {
	if t == nil {
		return
	}

	l.buf = l.buf[len(t.Literal):]

	l.nextPos = shiftPos(l.nextPos, t.Literal)

	if idx := strings.LastIndex(t.Literal, "\n"); idx >= 0 {
		l.loadedLine = t.Literal[idx+1:]
	} else {
		l.loadedLine += t.Literal
	}
}

func (l *Lexer) skipWhitespace() {
	if l.Whitespace == nil {
		return
	}

	for true {
		l.readBufIfNeed()

		if t := l.Whitespace.FindToken(l.buf, l.nextPos); t != nil {
			l.consumeBuffer(t)
		} else {
			break
		}
	}
}

func (l *Lexer) makeError() error {
	for shift, _ := range l.buf {
		if l.Whitespace != nil && l.Whitespace.FindToken(l.buf[shift:], l.nextPos) != nil {
			return UnknownTokenError{
				Literal:  l.buf[:shift],
				Position: l.nextPos,
			}
		}

		for _, tokenType := range l.TokenTypes {
			if tokenType.FindToken(l.buf[shift:], l.nextPos) != nil {
				return UnknownTokenError{
					Literal:  l.buf[:shift],
					Position: l.nextPos,
				}
			}
		}
	}

	return UnknownTokenError{
		Literal:  l.buf,
		Position: l.nextPos,
	}
}

/*
Посмотрите первый токен в буфере.
Возвращает nil как * Token, если буфер пуст.
*/
func (l *Lexer) Peek() (*Token, error) {
	for _, tokenType := range l.TokenTypes {
		l.skipWhitespace()

		l.readBufIfNeed()
		if t := tokenType.FindToken(l.buf, l.nextPos); t != nil {
			return t, nil
		}
	}

	if len(l.buf) > 0 {
		return nil, l.makeError()
	}

	return nil, nil
}

/*
Сканирование получит первый токен в буфере и удалит его из буфера.
*/
func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()

	l.consumeBuffer(t)

	return t, e
}

/*
GetCurrentLine возвращает строку последнего отсканированного токена.
*/
func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()

	if idx := strings.Index(l.buf, "\n"); idx >= 0 {
		return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
	} else {
		return l.loadedLine + l.buf
	}
}
