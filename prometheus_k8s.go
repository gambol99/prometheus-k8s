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
	"strings"
	"time"

	"github.com/golang/glog"
)

// NewPrometheusK8S creates a new prometheus service
func NewPrometheusK8S() (*PrometheusK8S, error) {
	glog.Infof("starting the Prometheus Watcher Service, version: %s", Version)

	client, err := NewKubeAPI()
	if err != nil {
		return nil, err
	}
	updatesCh := make(UpdateEvent, 10)

	return &PrometheusK8S{
		client:    client,
		updatesCh: updatesCh,
	}, nil
}

// StartServiceProcessor starts the service processor
func (r *PrometheusK8S) StartServiceProcessor() error {
	// step: we start watching out for events from the api
	if _, err := r.client.Watch(r.updatesCh); err != nil {
		glog.Errorf("failed to start watching out for events from kubernetes, error: %s", err)
		return err
	}

	// step: lets create a ticker to enforce refreshing
	ticker := time.NewTimer(time.Second * time.Duration(config.RefreshInterval))

	for {
		select {
		case <-ticker.C:
			glog.V(5).Infof("we have received a refresh interval, regenerating the config")
			r.GenerateConfiguration()
		case event := <-r.updatesCh:
			glog.V(4).Infof("we have received an update event from the watcher service, event: %s", event)
			// step: generate the content and write
			r.GenerateConfiguration()
		}
	}
}

// GenerateConfiguration render the configuration to file/s
func (r *PrometheusK8S) GenerateConfiguration() error {
	glog.Infof("generating the configuration of the prometheus nodes and services")

	// step: are we generating the nodes?
	if config.WithNodes {
		content, err := r.generateNodesConfiguration()
		if err != nil {
			glog.Errorf("Unable to retrieve the list of nodes: error: %s", err)
			return err
		}

		err = writeConfigFile(content, config.ConfigDirectory, config.NodesConfigFilename, config.DryRun)
		if err != nil {
			glog.Errorf("failed to write the node configuration, error: %s", err)
		}
	}

	if config.WithPods {
		content, err := r.generatePodsConfiguration()
		if err != nil {
			glog.Errorf("gnable to retrieve the list of pods: error: %s", err)
			return err
		}

		err = writeConfigFile(content, config.ConfigDirectory, config.PodsConfigFilename, config.DryRun)
		if err != nil {
			glog.Errorf("failed to write the pods configuration, error: %s", err)
		}
	}

	return nil
}

// generateNodesConfiguration generates the node config
func (r *PrometheusK8S) generateNodesConfiguration() ([]byte, error) {
	glog.V(4).Infof("generating the nodes configuration")
	// step: get the current list of nodes from the kubernetes
	nodes, err := r.client.Nodes()
	if err != nil {
		return nil, err
	}

	// step: create the targets group
	var targets []*Targets

	targets = append(targets, newTarget())
	for _, node := range nodes {
		targets[0].Targets = append(targets[0].Targets, fmt.Sprintf("%s:%d", node.ID, config.NodePort))
	}
	targets[0].Labels["role"] = "kubernetes_node"

	// step: marshall the config
	output, err := encode(targets)
	if err != nil {
		glog.Errorf("Failed to marshall the data, error: %s", err)
	}

	return output, nil
}

// renderPods: write the pod config to disk
func (r *PrometheusK8S) generatePodsConfiguration() ([]byte, error) {
	glog.V(4).Infof("generating the pod services configuration, namespaces: %s", config.Namespaces)

	var content []byte
	var targets []*Targets

	// step: get the current listing of pods
	namespaces := strings.Split(config.Namespaces, ",")

	for _, namespace := range namespaces {
		// step: check the namespace exists and if not, just skip
		found, err := r.client.NamespaceExists(namespace)
		if err != nil {
			glog.Errorf("unable to determine if the namespace: %s exists, error: %s", namespace, err)
			return content, err
		} else if !found {
			glog.Warningf("the namespace: %s does not exist, skipping retrieveing config", namespace)
			continue
		}

		// step: grab the pods within the specified namespace
		pods, err := r.client.Pods(namespace)
		if err != nil {
			glog.Errorf("unable to retrieve the list of pods with namespace: %s, error: %s", namespace, err)
			return content, err
		}

		glog.V(5).Infof("retrieved %d pods from namespace: %s, pods: #%v", len(pods), namespace, pods)

		// step: we iterate around and find all pods of the same 'Name' - effectively we are
		// grouping by the spec.labels['name'] for target groups, we also filter out any pods
		// which do not have a metrics annotation
		serviceGroups := make(map[string][]*Metrics, 0)
		for _, pod := range pods {
			// step: check if this pod name has already been found
			if _, found := serviceGroups[pod.Name]; found {
				continue
			}
			// step: check of the pod has annotations
			if _, found = pod.Annotations[config.MetricAnnotation]; !found {
				continue
			}

			// check: decode the metrics annotations
			metrics, err := decodeMetrics(pod.Annotations[config.MetricAnnotation])
			if err != nil {
				glog.Errorf("skipping pod: '%s', name: '%s' as the metrics config is invalid, error: %s", pod.ID, pod.Name, err)
				continue
			}

			serviceGroups[pod.Name] = metrics
		}

		// check: if we have nothing to render, lets continue
		if len(serviceGroups) <= 0 {
			continue
		}

		// step: now we iterate the pods again, group by the service_names and produce
		// the target groups per service name
		for serviceName, metrics := range serviceGroups {
			target := newTarget()
			target.Labels["pod"] = serviceName

			for _, pod := range pods {
				if pod.Name == serviceName {
					// step: copy in the rest of the pod labels
					target.Labels["namespace"] = pod.Namespace
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
			targets = append(targets, target)
		}
	}

	// step: marshall the config into format
	content, err := encode(targets)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshall the target into format, error: %s", err)
	}

	return content, nil
}
