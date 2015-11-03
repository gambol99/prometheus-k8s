/*
Copyright 2014 Rohith All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func main() {
	// step: parse the command line configuration
	if err := parseConfig(); err != nil {
		glog.Errorf("invalid configuration, error: %s", err)
		os.Exit(1)
	}

	// step: create the service
	service, err := NewPrometheusK8S()
	if err != nil {
		glog.Errorf("failed to create prometheus-k8s service, error: %s", err)
		os.Exit(1)
	}

	// step: create a exit channel
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// step: generate the onetime config
	err = service.GenerateConfiguration()
	if err != nil {
		glog.Errorf("failed to generate the initial configuration, error: %s", err)
	}

	// step: start the service processor
	go func() {
		err := service.StartServiceProcessor()
		if err != nil {
			glog.Fatalf("failed to start the service, error: %s", err)
		}
	}()

	// step: wait for a exit signal
	<-signalChannel
}
