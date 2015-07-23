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

//
// Event ... the update event itself
type Event struct {
	// the source of the event i.e. node or pod
	Type int
	// a reference to the event
	Event interface{}
}

// ShutdownChannel ... a channel used to indicate we wish to shut something down
type ShutdownChannel chan bool

// UpdateEvent ... A messaging channel used to send events when they have occurred
type UpdateEvent chan *Event

// KubeAPI ... is the service responsible for commmunicating with the api and
// producing a stream of events related to node and pods changes
type KubeAPI interface {
	// retrieve a list of nodes from kubernetes
	Nodes() ([]*Node, error)
	// retrieve a list of running pods from within a namespace
	Pods() ([]*Pod, error)
	// watch for changes in nodes and pods and update
	Watch(UpdateEvent) (ShutdownChannel, error)
}
