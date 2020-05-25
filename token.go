package easylexer

import (
	"regexp"
	"strconv"
	"strings"
)

// TokenID является идентификатором для TokenType.
type TokenID int

// Идентификаторы токенов по умолчанию.
const (
	OTHER TokenID = -(iota + 1)
	IDENT
	NUMBER
	STRING
)

/*
Преобразование в читаемую строку.
Добавленные пользователем идентификаторы токенов преобразуются в НЕИЗВЕСТНЫЕ.
*/
func (id TokenID) String() string {
	switch id {
	case OTHER:
		return "OTHER"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	default:
		return "UNKNOWN(" + strconv.Itoa(int(id)) + ")"
	}
}

/*
TokenType - это правило для создания токена.
GetID возвращает TokenID этого TokenType.
TokenID может делиться с другим TokenType.
FindToken возвращает новый токен, если заголовок первого аргумента был сопоставлен с шаблоном этого TokenType.
Второй аргумент - это позиция токена в буфере. В почти реализации Позиция будет напрямую переходить в токен результата.
*/
type TokenType interface {
	GetID() TokenID
	FindToken(string, Position) *Token
}

/*
RegexpTokenType - это реализация TokenType с регулярным выражением.
Идентификатором является TokenID для этого типа токена.
Re - регулярное выражение токена. Это должно начинаться с "^".
*/
type RegexpTokenType struct {
	ID TokenID
	Re *regexp.Regexp
}

/*
Сделайте новый RegexpTokenType.
id - это TokenID нового RegexpTokenType.
Это регулярное выражение токена.
*/
func NewRegexpTokenType(id TokenID, re string) *RegexpTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &RegexpTokenType{
		ID: id,
		Re: regexp.MustCompile(re),
	}
}

// Получить читаемую строку TokenID.
func (rtt *RegexpTokenType) String() string {
	return rtt.ID.String()
}

// GetID возвращает идентификатор этого типа токена.
func (rtt *RegexpTokenType) GetID() TokenID {
	return rtt.ID
}

// FindToken возвращает новый токен, если s начинается с этого токена.
func (rtt *RegexpTokenType) FindToken(s string, p Position) *Token {
	m := rtt.Re.FindStringSubmatch(s)
	if len(m) > 0 {
		return &Token{
			Type:       rtt,
			Literal:    m[0],
			Submatches: m[1:],
			Position:   p,
		}
	}
	return nil
}

/*
PatternTokenType - это тип токена словаря.
PatternTokenType имеет несколько строк и находит токен, который идеально подходит им.
*/
type PatternTokenType struct {
	ID       TokenID
	Patterns []string
}

/*
Создайте новый тип токена.
id - это идентификатор токена нового шаблона TokenType.
шаблоны это массив шаблонов.
*/
func NewPatternTokenType(id TokenID, patterns []string) *PatternTokenType {
	return &PatternTokenType{
		ID:       id,
		Patterns: patterns,
	}
}

// Получить читаемую строку TokenID.
func (ptt *PatternTokenType) String() string {
	return ptt.ID.String()
}

// GetID возвращает идентификатор типа токена.
func (ptt *PatternTokenType) GetID() TokenID {
	return ptt.ID
}

// FindToken возвращает новый токен, если s начинается с этого токена.
func (ptt *PatternTokenType) FindToken(s string, p Position) *Token {
	for _, x := range ptt.Patterns {
		if strings.HasPrefix(s, x) {
			return &Token{
				Type:     ptt,
				Literal:  x,
				Position: p,
			}
		}
	}
	return nil
}

// Данные найденного токена.
type Token struct {
	Type       TokenType
	Literal    string   // Строка соответствует.
	Submatches []string // Подчеркивания регулярного выражения.
	Position   Position // Положение токена.
}
