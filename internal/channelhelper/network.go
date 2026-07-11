package channelhelper

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/StorkKershaw/hound/internal/windows/networking/connectivity"
	"github.com/go-ole/go-ole"
)

const (
	RO_INIT_MULTITHREADED = 0x01
)

func MonitorNetwork(in <-chan struct{}, out chan<- struct{}) {
	runtime.LockOSThread()

	if err := ole.RoInitialize(RO_INIT_MULTITHREADED); err != nil {
		log.Printf("could not initialize Windows Runtime: %v", err)
		return
	}

	handler := connectivity.NewNetworkStatusChangedEventHandler(
		ole.NewGUID(connectivity.GUIDNetworkStatusChangedEventHandler),
		func(instance *connectivity.NetworkStatusChangedEventHandler, sender unsafe.Pointer) {
			out <- struct{}{}
		},
	)

	token, err := connectivity.NetworkInformationAddNetworkStatusChanged(handler)
	if err != nil {
		log.Printf("Error adding network status changed handler: %v", err)
		return
	}

	out <- struct{}{}

	<-in

	if err := connectivity.NetworkInformationRemoveNetworkStatusChanged(token); err != nil {
		log.Printf("Error removing network status changed handler: %v", err)
		return
	}
}
