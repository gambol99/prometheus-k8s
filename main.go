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
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

var (
	// the api for the kubernetes service
	kubeapi KubeAPI
)

const (
	METRICS_ANNOTATION = "metrics"
)

func main() {
	var err error

	// step: parse the command line configuration
	if err := parseConfig(); err != nil {
		glog.Errorf("Invalid configuration, error: %s", err)
		os.Exit(1)
	}
	// step: create a new watcher
	if kubeapi, err = NewKubeAPI(); err != nil {
		glog.Errorf("Failed to create a client for Kubernetes API, error: %s", err)
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	glog.Infof("Starting the Prometheus Watcher Service, version: %s", VERSION)

	// step: we create a channel to receive updates from the api
	serviceUpdatesCh := make(UpdateEvent, 10)
	// step: we start watching out for events from the api
	if _, err := kubeapi.Watch(serviceUpdatesCh); err != nil {
		glog.Errorf("Failed to start watching out for events from kubernetes, error: %s", err)
		os.Exit(1)
	}

	// step: perform a initial render
	generateConfiguration()

	var event_loop = true
	// step: lets create a ticker to enforce refreshing
	ticker := time.NewTimer(time.Second * time.Duration(config.RefreshInterval))

	// step: we loop forever, or until a signal destroy
	for event_loop {
		// step: we need to start watching for changes in the api
		select {
		case <-ticker.C:
			glog.V(4).Infof("We have received a refresh interval, regenerating the config")
			generateConfiguration()
		case event := <-serviceUpdatesCh:
			glog.V(4).Infof("We have received an update event from the watcher service, event: %s", event)
			// step: generate the content and write
			generateConfiguration()
		case <-signalChannel:
			// we need to shutdown the service
			glog.Infof("Shutting down the service, we have received a kill signal")
			event_loop = false
		}
	}
}

// generateConfiguration: render the configuration to file/s
func generateConfiguration() error {
	glog.Infof("Updating the configuration of the prometheus endpoints / targets")
	// step: are we generating the nodes?
	if config.WithNodes {
		content, err := generateNodesConfiguration()
		if err != nil {
			glog.Errorf("Unable to retrieve the list of nodes: error: %s", err)
			return err
		}
		writeConfiguration(content, fmt.Sprintf("%s/%s", config.ConfigDirectory, config.NodesConfigFilename))
	}
	if config.WithNodes {
		content, err := generatePodsConfiguration()
		if err != nil {
			glog.Errorf("Unable to retrieve the list of pods: error: %s", err)
			return err
		}
		writeConfiguration(content, fmt.Sprintf("%s/%s", config.ConfigDirectory, config.PodsConfigFilename))
	}
	return nil
}

// writeConfiguration: write the configuration to a file or to screen
func writeConfiguration(content []byte, filename string) error {
	if config.DryRun {
		fmt.Println("----")
		fmt.Printf("filename: '%s'\n", filename)
		fmt.Printf("content: \n%s", content)
		fmt.Println("----")
	} else {
		// step: open the file for
		err := ioutil.WriteFile(filename, content, os.FileMode(0664))
		if err != nil {
			glog.Errorf("Failed to write to file: '%s', error: %s", filename, err)
			return err
		}
	}
	return nil
}

// renderNodes: write the node config file to the disk
func generateNodesConfiguration() ([]byte, error) {
	glog.V(4).Infof("Rendering the nodes to configuration")
	// step: get the current list of nodes from the kubernetes
	nodes, err := kubeapi.Nodes()
	if err != nil {
		return nil, err
	}
	// step: create the targets group
	tgroups := make([]*Targets, 0)
	tgroups = append(tgroups, NewTarget())
	for _, node := range nodes {
		tgroups[0].Targets = append(tgroups[0].Targets, fmt.Sprintf("%s:%d", node.ID, config.NodePort))
	}
	tgroups[0].Labels["role"] = "kubernetes_nodes"
	// step: marshall the config
	output, err := encode(tgroups)
	if err != nil {
		glog.Errorf("Failed to marshall the data, error: %s", err)
	}
	return output, nil
}

// renderPods: write the pod config to disk
func generatePodsConfiguration() ([]byte, error) {
	glog.V(4).Infof("Generating the pod services configuration")

	// step: get the current listing of pods
	pods, err := kubeapi.Pods()
	if err != nil {
		glog.Errorf("Unable to retrieve the list of pods, error: %s", err)
	}
	glog.V(5).Infof("Retrieved %d pods from kubernetes: %s", len(pods), pods)

	// step: we iterate around and find all pods of the same 'Name' - effectively we are
	// grouping by the spec.labels['name'] for target groups, we also filter out any pods
	// which do not have a metrics annotation
	service_groups := make(map[string][]*Metrics, 0)
	for _, pod := range pods {
		if _, found := service_groups[pod.Name]; !found {
			// check: does the pod have a metrics annotation?
			if _, found := pod.Annotations[METRICS_ANNOTATION]; found {
				// check: decode the metrics
				metrics, err := decodeMetrics(pod.Annotations[METRICS_ANNOTATION])
				if err != nil {
					glog.Errorf("Skipping pod: '%s', name: '%s' as the metrics config is invalid, error: %s",
						pod.ID, pod.Name, err)
					continue
				}
				service_groups[pod.Name] = metrics
			}
		}
	}

	// check: if we have nothing to render, lets return an empty string
	if len(service_groups) <= 0 {
		return []byte(""), nil
	}

	groups := make([]*Targets, 0)

	// step: now we iterate the pods again, group by the service_names and produce
	// the target groups per service name
	for service_name, metrics := range service_groups {
		target := NewTarget()
		target.Labels["pod"] = service_name

		for _, pod := range pods {
			if pod.Name == service_name {
				// step: copy in the rest of the pod labels
				for k, v := range pod.Labels {
					if k == "pod" {
						continue
					}
					target.Labels[k] = v
				}
				// step: we produce a endpoint for each metrics listed
				for _, metric := range metrics {
					target.Targets = append(target.Targets, fmt.Sprintf("%s:%d", pod.Address, metric.Port))
				}
			}
		}
		// step: append the group to the groups
		groups = append(groups, target)
	}

	// step: marshall the config into format
	format, err := encode(groups)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall the target into format, error: %s", err)
	}
	return format, nil
}

