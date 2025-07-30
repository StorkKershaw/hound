package main

import (
	"log"

	"github.com/StorkKershaw/hound/internal/channel"
)

func main() {
	log.SetFlags(0)

	termination := channel.ProduceBy(terminate)
	statusChange := channel.TransformBy(termination, monitorNetworkStatus)
	connectionChange := channel.TransformBy(statusChange, getNetworkConnection)
	errors := channel.TransformBy(connectionChange, authenticate)
	<-channel.TransformBy(errors, notify)
}
