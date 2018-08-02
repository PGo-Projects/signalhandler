package signalhandler

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func sampleHandler() {
	fmt.Println("This is a sample handler")
}

func TestSignalHandler(t *testing.T) {
	handler := New(sampleHandler, os.Interrupt, syscall.SIGTERM)
	handler.WithSignalBlocked(func() error {
		fmt.Println("Inside Ctrl-C free environment")
		time.Sleep(5 * time.Second)
		fmt.Println("Exiting Ctrl-C free environment")
		return nil
	})
}
