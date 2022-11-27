package signalListener

import (
	"os"
	"os/signal"
	"syscall"
)

type SignalListener struct {
	notify chan os.Signal
}

func New() *SignalListener {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	return &SignalListener{
		notify: signalChan,
	}
}

func (s *SignalListener) Notify() <-chan os.Signal {
	return s.notify
}
