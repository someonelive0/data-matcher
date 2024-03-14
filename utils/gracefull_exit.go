package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitExit will block until os signal happened
func WaitExit(c chan os.Signal, exit func()) {
	for i := range c {
		switch i {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// fmt.Println("receive exit signal ", i.String(), ",exit...")
			exit()
			os.Exit(0)
		}
	}
}

// NewShutdownSignal new normal Signal channel
func NewShutdownSignal() chan os.Signal {
	c := make(chan os.Signal, 3)
	// SIGHUP: terminal closed
	// SIGINT: Ctrl+C
	// SIGTERM: program exit
	// SIGQUIT: Ctrl+/
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return c
}

/*
	// Add follow codes to main()

	// wait until graceful exit
	signalChan := NewShutdownSignal()
	WaitExit(signalChan, func() {
		// your clean code, such as
		// cancel()
		// close(ch)
		// wg.Wait()
	})
*/
