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

import "fmt"

// PrometheusK8S is the main service wrapper
type PrometheusK8S struct {
	// the client for k8s
	client KubeAPI
	// the update events
	updatesCh UpdateEvent
}

// Event represents an update event itself
type Event struct {
	// the source of the event i.e. node or pod
	Type int
	// a reference to the event
	Event interface{}
}

// ShutdownChannel is a channel used to indicate we wish to shut something down
type ShutdownChannel chan bool

// UpdateEvent is a messaging channel used to send events when they have occurred
type UpdateEvent chan *Event

// KubeAPI is the service responsible for communicating with the api and
// producing a stream of events related to node and pods changes
type KubeAPI interface {
	// checks to see if a namespace exists
	NamespaceExists(string) (bool, error)
	// retrieve a list of nodes from kubernetes
	Nodes() ([]*Node, error)
	// retrieve a list of running pods from within a namespace
	Pods(string) ([]*Pod, error)
	// watch for changes in nodes and pods and update
	Watch(UpdateEvent) (ShutdownChannel, error)
}

// Pod is a normalize form of running pod
type Pod struct {
	// the name / id of the pod
	ID string
	// the label name
	Name string
	// the namespace of the pod
	Namespace string
	// the labels associated to the pod
	Labels map[string]string
	// the annotations associated to the pod
	Annotations map[string]string
	// the ip address of the pod
	Address string
}

// Node is the definition of the kubernetes node
type Node struct {
	// the name / ID of the node
	ID string
	// the labels associated to the node
	Labels map[string]string
}

// Targets is the structure of the prometheus file discovery targets
type Targets struct {
	// the array of hosts within this target
	Targets []string `yaml:"targets",json:"targets"`
	// the labels associated to these targets
	Labels map[string]string `yaml:"labels",json:"labels"`
}

// Metrics is the structure used to produce details about the metric endpoints
// being exported by a pod
type Metrics struct {
	// the name of the metric (optional)
	Name string `yaml:"name,omitempty",json:"name,omitempty"`
	// the port of the metric
	Port int `yaml:"port",json:"port"`
	// the endpoint (optional)
	Endpoint string `yaml:"endpoint,omitempty",json:"endpoint,omitempty"`
}

func (r Pod) String() string {
	return fmt.Sprintf(`
ID: %s
Name: %s
Namespace: %s
Labels: %s
Annotations: %s
Address: %s
`, r.ID, r.Name, r.Namespace, r.Labels, r.Annotations, r.Address)
}
