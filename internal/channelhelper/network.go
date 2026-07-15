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
		log.Printf("could not add network status change handler: %v", err)
		return
	}

	out <- struct{}{}
	log.Printf("watching for network status changes")

	<-in

	if err := connectivity.NetworkInformationRemoveNetworkStatusChanged(token); err != nil {
		log.Printf("could not remove network status change handler: %v", err)
		return
	}
}
