package signalhandler

import (
	"os"
	"os/signal"
)

type HandlerFunc func()

type signalHandler struct {
	channel    chan os.Signal
	handler    HandlerFunc
	oldHandler HandlerFunc
}

func New(handler HandlerFunc, sig ...os.Signal) *signalHandler {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, sig...)
	h := &signalHandler{
		channel: sigchan,
		handler: handler,
	}
	h.listen()
	return h
}

func (s *signalHandler) listen() {
	go func() {
		for {
			<-s.channel
			s.handler()
		}
	}()
}

func (s *signalHandler) WithSignalBlocked(signalFreeFunc func()) {
	s.oldHandler = s.handler
	signalRaised := false
	s.handler = func() {
		signalRaised = true
	}
	signalFreeFunc()
	s.handler = s.oldHandler
	if signalRaised {
		s.handler()
	}
}

func (s *signalHandler) SetHandler(handler HandlerFunc) {
	s.handler = handler
}
