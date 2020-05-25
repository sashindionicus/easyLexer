package easylexer_test

import (
	"testing"

	"github.com/sashindionicus/easyLexer"
)

func TestUnknownTokenError(t *testing.T) {
	err := simplexer.UnknownTokenError{Literal: "test", Position: simplexer.Position{Line: 0, Column: 1}}
	except := "1:2:UnknownTokenError: \"test\""

	if err.Error() != except {
		t.Errorf("excepted %#v but got %s", except, err.Error())
	}
}
