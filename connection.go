package main

import (
	"log"

	"github.com/StorkKershaw/hound/internal/windows/networking/connectivity"
)

func getNetworkConnection(in <-chan struct{}, out chan<- struct{}) {
	for range in {
		func() { // closure to execute `defer` after each iteration
			profile, err := connectivity.NetworkInformationGetInternetConnectionProfile()
			if profile == nil || err != nil {
				log.Printf("Error getting internet connection profile: %v", err)
				return
			}
			defer profile.Release()

			level, err := profile.GetNetworkConnectivityLevel()
			if err != nil {
				log.Printf("Error getting network connectivity level: %v", err)
				return
			}
			if level != connectivity.NetworkConnectivityLevelConstrainedInternetAccess {
				log.Printf("No authentication is required for network connectivity level: %v", level)
				return
			}

			out <- struct{}{}
		}()
	}
}
