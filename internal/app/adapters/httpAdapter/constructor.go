package httpAdapter

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/config"
	"github.com/Rostislaved/ethereum-parser/internal/app/parser"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout       = 5 * time.Second
	_defaultWriteTimeout      = 50 * time.Second
	_defaultShutdownTimeout   = 3 * time.Second
	_defaultReadHeaderTimeout = 5 * time.Second
)

type HttpAdapter struct {
	config config.Server
	server *http.Server
	parser *parser.Parser
	notify chan error
}

func New(config config.Server, parser *parser.Parser) *HttpAdapter {
	adapter := &HttpAdapter{
		config: config,
		parser: parser,
		notify: make(chan error, 1),
	}

	httpServer := &http.Server{
		Handler:           adapter.getMux(),
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
		ReadHeaderTimeout: _defaultReadHeaderTimeout,
		Addr:              config.Addr,
	}

	adapter.server = httpServer

	return adapter
}
