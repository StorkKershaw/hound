package main

import (
	"log"

	"github.com/StorkKershaw/hound/internal/channel"
	helper "github.com/StorkKershaw/hound/internal/channelhelper"
	"github.com/StorkKershaw/hound/internal/toast"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(toast.Writer{})

	termination := channel.Produce(helper.Terminate)
	statusChange := channel.Transform(termination, helper.MonitorNetwork)
	connectionChange := channel.Transform(statusChange, helper.WatchConnectivity)
	<-channel.Transform(connectionChange, helper.Authenticate)
}
