package signalhandler

import (
	"os"
	"os/signal"
)

type HandlerFunc func()

type SignalHandler struct {
	channel    chan os.Signal
	handler    HandlerFunc
	oldHandler HandlerFunc
	_          struct{}
}

func New(handler HandlerFunc, sig ...os.Signal) *SignalHandler {
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, sig...)
	h := &SignalHandler{
		channel: sigchan,
		handler: handler,
	}
	h.listen()
	return h
}

func (s *SignalHandler) listen() {
	go func() {
		for {
			<-s.channel
			s.handler()
		}
	}()
}

func (s *SignalHandler) WithSignalBlocked(signalFreeFunc func() error) error {
	return s.WithSignalBlockedAndSignalMessageFunc(signalFreeFunc, nil)
}

func (s *SignalHandler) WithSignalBlockedAndSignalMessageFunc(signalFreeFunc func() error, messageFunc func()) error {
	s.oldHandler = s.handler
	signalRaised := false
	s.handler = func() {
		signalRaised = true
		if messageFunc != nil {
			messageFunc()
		}
	}
	err = signalFreeFunc()
	s.handler = s.oldHandler
	if signalRaised {
		s.handler()
	}
	return err
}

func (s *SignalHandler) SetHandler(handler HandlerFunc) {
	s.handler = handler
}
