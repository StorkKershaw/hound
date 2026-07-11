package channelhelper

import (
	"os"
	"os/signal"
	"syscall"
)

func Terminate(out chan<- struct{}) {
	termination := make(chan os.Signal, 1)
	defer close(termination)
	signal.Notify(termination, syscall.SIGINT, syscall.SIGTERM)
	<-termination
	out <- struct{}{}
}
