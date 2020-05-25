package easylexer

import "fmt"

// Ошибка, которая возвращается при обнаружении неизвестного токена.
type UnknownTokenError struct {
	Literal  string
	Position Position
}

// Получить сообщение об ошибке в виде строки.
func (se UnknownTokenError) Error() string {
	return fmt.Sprintf("%d:%d:UnknownTokenError: %#v", se.Position.Line+1, se.Position.Column+1, se.Literal)
}
