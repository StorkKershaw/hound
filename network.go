package main

import (
	"log"
	"unsafe"

	"github.com/StorkKershaw/hound/internal/windows/networking/connectivity"
	"github.com/go-ole/go-ole"
)

const (
	RO_INIT_MULTITHREADED = 0x01
)

func monitorNetworkStatus(in <-chan struct{}, out chan<- struct{}) {
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

	<-in

	if err := connectivity.NetworkInformationRemoveNetworkStatusChanged(token); err != nil {
		log.Printf("Error removing network status changed handler: %v", err)
		return
	}
}
