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

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

// Implements the KubeAPI service interface
type kubeAPIImpl struct {
	// the kubernetes api client
	client *unversioned.Client
}

// NewKubeAPI ... creates a new watch service for kubernetes
func NewKubeAPI() (KubeAPI, error) {
	glog.Infof("Creating a new Kube API service, api: %s", getURL())
	service := new(kubeAPIImpl)
	kube, err := service.newAPIClient()
	if err != nil {
		return nil, err
	}
	service.client = kube

	return service, nil
}

// NamespaceExists checks to see if a namespace exists in k8s
func (r kubeAPIImpl) NamespaceExists(namespace string) (bool, error) {
	glog.V(10).Infof("checking for namespace: %s", namespace)
	if namespace == api.NamespaceAll {
		return true, nil
	}

	namespaces, err := r.client.Namespaces().List(labels.Everything(), fields.Everything())
	if err != nil {
		return false, err
	}

	for _, name := range namespaces.Items {
		if name.Name == namespace {
			return true, nil
		}
	}

	return false, nil
}

// Nodes retrieves a list of nodes from Kuberntes, normalize them and give me the list
func (r kubeAPIImpl) Nodes() ([]*Node, error) {
	nodes, err := r.client.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		glog.Errorf("Failed to retrieve a list of nodes from the api, error: %s", err)
		return nil, err
	}
	var list []*Node
	// step: iterate and normalize the node
	for _, x := range nodes.Items {
		glog.V(10).Infof("Adding the node: %s into the list of kubernetes nodes", x.Name)
		node := &Node{
			ID:     x.Name,
			Labels: x.Labels,
		}
		list = append(list, node)
	}
	glog.V(10).Infof("Retrieved the node: %v from api", list)

	return list, nil
}

// Pods retrieves a list of running within the namespace
func (r kubeAPIImpl) Pods(namespace string) ([]*Pod, error) {
	glog.V(10).Infof("Retrieving a list of the running pods")

	// step: get a list of the pods and find the current revision
	var list []*Pod

	pods, err := r.client.Pods(namespace).List(labels.Everything(), fields.Everything())
	if err != nil {
		glog.Errorf("Failed to retrieve a list of pods running, error: %s", err)
		return nil, err
	}

	// step: iterate and normalize the pods
	for _, x := range pods.Items {
		// step: we have to make sure the pod is running, otherwise it probably won't have an IP address
		if x.Status.Phase == api.PodRunning {
			glog.V(10).Infof("Adding the pod: %s, addesss: %s into the running list", x.Name, x.Status.PodIP)
			pod := &Pod{
				ID:          x.Name,
				Name:        x.Labels["name"],
				Namespace:   x.Namespace,
				Labels:      x.Labels,
				Annotations: x.Annotations,
				Address:     x.Status.PodIP,
			}
			list = append(list, pod)
		}
	}

	return list, nil
}

//
// Watch is the main entry-point for the service, we listen out for changes in the
// nodes, pods and the refresh timer
func (r *kubeAPIImpl) Watch(updates UpdateEvent) (ShutdownChannel, error) {
	var err error
	var nodeCh, podsCh watch.Interface

	// step: create the done channel
	shutdownCh := make(ShutdownChannel)
	go func() {
		// step: wait for the cleanup channel
		<-shutdownCh
	}()

	// step: acquire a nodes watch
	if nodeCh, err = r.createNodesWatch(); err != nil {
		return nil, err
	}
	// step: create a channel for pod updates
	if podsCh, err = r.createPodsWatch(); err != nil {
		return nil, err
	}

	// notes: the main loop to the service; we wait for changes in the nodes,
	// the pods or the refresh timer
	go func() {
		glog.V(10).Infof("Starting the event service loop")
		for {
			select {
			case update := <-nodeCh.ResultChan():
				// step: we only care about added or removed nodes, not modified
				if update.Type == watch.Modified {
					continue
				}
				event := newEvent(nodeEvent, update)
				glog.V(5).Infof("Recieved an update to the nodes: %v", event)
				updates <- event
			case update := <-podsCh.ResultChan():
				event := newEvent(podEvent, update)
				glog.V(5).Infof("Recieved an update to the pods: %v", event)
				updates <- event
			}
		}
	}()

	return shutdownCh, nil
}

// createPodsWatch creates a watcher channel for changes on the pods within the configured namespace
func (r kubeAPIImpl) createPodsWatch() (watch.Interface, error) {
	glog.V(10).Infof("Creating a watcher for the kubernetes pods")
	// step: lets retrieve a revision from which to work from
	list, err := r.client.Pods(api.NamespaceAll).List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the list of pods, error: %s", err)
	}

	// step: create a channel for watching the nodes
	ch, err := r.client.Pods(api.NamespaceAll).Watch(labels.Everything(), fields.Everything(),
		api.ListOptions{ResourceVersion: list.ResourceVersion})

	if err != nil {
		return nil, fmt.Errorf("unable to create a watch on pods resources, reason: %s", err)
	}

	return ch, nil
}

// createNodesWatch creates a nodes update interface used to watch changes in nodes
func (r kubeAPIImpl) createNodesWatch() (watch.Interface, error) {
	glog.V(10).Infof("Creating a watcher for the kubernetes nodes")
	// step: lets retrieve a revision from which to work from
	list, err := r.client.Nodes().List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a list of nodes, error: %s", err)
	}

	nodeCh, err := r.client.Nodes().Watch(labels.Everything(), fields.Everything(), api.ListOptions{ResourceVersion: list.ResourceVersion})
	if err != nil {
		return nil, fmt.Errorf("unable to create a watch on node resources, reason: %s", err)
	}

	return nodeCh, nil
}

// newAPIClient creates a new client to speak to the kubernetes api service
func (r *kubeAPIImpl) newAPIClient() (*unversioned.Client, error) {
	// step: create the configuration
	cfg := unversioned.Config{
		Host:     getURL(),
		Insecure: config.HTTPInsecure,
		Version:  config.APIVersion,
	}

	// check: ensure the token file exists
	if config.TokenFile != "" {
		if _, err := os.Stat(config.TokenFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("the token file: %s does not exist", config.TokenFile)
		}

		content, err := ioutil.ReadFile(config.TokenFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read the token file: %s, error: %s", config.TokenFile, err)
		}
		config.Token = string(content)
	}

	// check: are we using a user token to authenticate?
	if config.Token != "" {
		cfg.BearerToken = config.Token
	}

	// check: are we using a cert to authenticate
	if config.CaCertFile != "" {
		cfg.Insecure = false
		cfg.TLSClientConfig = unversioned.TLSClientConfig{
			CAFile: config.CaCertFile,
		}
	}

	// step: initialize the client
	kube, err := unversioned.New(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create a kubernetes api client, reason: %s", err)
	}

	return kube, nil
}

// getURL: generate the url used to communicate with the kubernetes api service
func getURL() string {
	return fmt.Sprintf("%s://%s:%d", config.APIProtocol, config.Host, config.Port)
}
