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
	"sync"
	"testing"
)

type fakeKubeAPI struct{}

var (
	apiLock sync.Once
)

func newFakeKubeAPI(t *testing.T) KubeAPI {
	return &fakeKubeAPI{}
}

// checks to see if a namespace exists
func (r fakeKubeAPI) NamespaceExists(namespace string) (bool, error) {
	namespaces := []string{
		"default",
		"platform",
	}
	if namespace == "" {
		return true, nil
	}

	for _, ns := range namespaces {
		if namespace == ns {
			return true, nil
		}
	}

	return false, nil
}

func (r fakeKubeAPI) Nodes() ([]*Node, error) {
	return []*Node{
		{
			ID: "10.50.0.101",
			Labels: map[string]string{
				"kubernetes": "true",
			},
		},
		{
			ID: "10.50.0.102",
			Labels: map[string]string{
				"kubernetes": "true",
			},
		},
		{
			ID: "10.50.0.103",
			Labels: map[string]string{
				"kubernetes": "true",
			},
		},
	}, nil
}

func (r fakeKubeAPI) Pods(namespace string) ([]*Pod, error) {
	pods := map[string][]*Pod{
		"default": {
			{
				ID:        "nginx_8327",
				Name:      "nginx",
				Namespace: "default",
				Labels: map[string]string{
					"name": "nginx",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: collectd-exporter\n  port: 9103\n",
				},
				Address: "10.10.0.100",
			},
			{
				ID:        "nginx_dsd2",
				Name:      "nginx",
				Namespace: "default",
				Labels: map[string]string{
					"name": "nginx",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: collectd-exporter\n  port: 9103\n",
				},
				Address: "10.10.0.101",
			},
			{
				ID:        "nginx_dsdd2",
				Name:      "nginx",
				Namespace: "default",
				Labels: map[string]string{
					"name": "nginx",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: collectd-exporter\n  port: 9103\n",
				},
				Address: "10.10.0.103",
			},
		},
		"platform": {
			{
				ID:        "prometheus",
				Name:      "prometheus",
				Namespace: "platform",
				Labels: map[string]string{
					"name": "prometheus",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: prometheus-exporter\n  port: 1000\n",
				},
				Address: "10.10.2.10",
			},
			{
				ID:        "prometheus",
				Name:      "prometheus",
				Namespace: "platform",
				Labels: map[string]string{
					"name": "nginx",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: prometheus-exporter\n  port: 1000\n",
				},
				Address: "10.10.1.4",
			},
			{
				ID:        "prometheus",
				Name:      "prometheus",
				Namespace: "platform",
				Labels: map[string]string{
					"name": "prometheus",
				},
				Annotations: map[string]string{
					config.MetricAnnotation: "- name: prometheus-exporter\n  port: 1000\n",
				},
				Address: "10.10.0.13",
			},
		},
	}

	if namespace == "" {
		var list []*Pod
		for _, kpod := range pods {
			list = append(list, kpod...)
		}
		return list, nil
	}

	list, _ := pods[namespace]

	return list, nil
}

func (r fakeKubeAPI) Watch(UpdateEvent) (ShutdownChannel, error) {
	return nil, nil
}
