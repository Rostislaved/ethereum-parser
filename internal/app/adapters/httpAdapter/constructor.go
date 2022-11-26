package httpAdapter

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/parser"
)

type HttpAdapter struct {
	parser *parser.Parser
}

func New(parser *parser.Parser) *HttpAdapter {
	return &HttpAdapter{
		parser: parser,
	}
}
