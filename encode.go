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
	"github.com/golang/glog"
)

// encode ... convert / marshall the data structure into the specified format
func encode(data interface{}) (output []byte, err error) {
	output, err = yaml.Marshal(data)
	if err != nil {
		glog.Errorf("Failed to marshall the structure to yaml, %s, error: %s", data, err)
		return nil, fmt.Errorf("marshalling failure, data: %V, error: %s", data, err)
	}
	return
}

// decode ... decodes the string into a actual structure
func decode(input []byte, output interface{}) error {
	err := yaml.Unmarshal(input, output)
	if err != nil {
		glog.Errorf("Failed to decode the yaml: %s, error: %s", input, err)
		return err
	}
	return nil
}

// decodeMetrics ... decodes the annotated metrics back into the correct structure
func decodeMetrics(cfg string) ([]*Metrics, error) {
	var metrics []*Metrics
	if err := decode([]byte(cfg), &metrics); err != nil {
		return nil, fmt.Errorf("invalid metric config, error: %s", err)
	}
	return metrics, nil
}
