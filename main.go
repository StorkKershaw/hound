package main

import (
	"log"

	"github.com/StorkKershaw/hound/internal/channel"
	helper "github.com/StorkKershaw/hound/internal/channelhelper"
)

func main() {
	log.SetFlags(0)

	termination := channel.Produce(helper.Terminate)
	statusChange := channel.Transform(termination, helper.MonitorNetwork)
	connectionChange := channel.Transform(statusChange, helper.WatchConnectivity)
	browserErrors := channel.Transform(connectionChange, helper.Authenticate)
	<-channel.Transform(browserErrors, helper.Notify)
}
