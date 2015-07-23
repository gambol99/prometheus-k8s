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
	"gopkg.in/yaml.v2"
)

// Metrics ... is the structure used to produce details about the metric endpoints
// being exported by a pod
type Metrics struct {
	// the name of the metric (optional)
	Name string `yaml:"name,omitempty", json:"name,omitempty"`
	// the port of the metric
	Port int `yaml:"port", json:"port"`
	// the endpoint (optional)
	Endpoint string `yaml:"endpoint,omitempty", json:"endpoint,omitempty"`
}

func DecodeMetrics(cfg string) ([]*Metrics, error) {
	if cfg == "" {
		return nil, fmt.Errorf("invalid metric config, value is empty")
	}
	var metrics []*Metrics
	if err := yaml.Unmarshal([]byte(cfg), &metrics); err != nil {
		return nil, fmt.Errorf("invalid metric config, error: %s", err)
	}
	return metrics, nil
}
