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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	metrics := &Metrics{
		Name:     "test",
		Port:     9090,
		Endpoint: "/metrics",
	}
	content, err := encode(metrics)
	assert.Nil(t, err)
	assert.NotEmpty(t, content)

}

func TestDecode(t *testing.T) {
	metrics := &Metrics{
		Name:     "test",
		Port:     9090,
		Endpoint: "/metrics",
	}
	content, err := encode(metrics)
	assert.Nil(t, err)
	assert.NotEmpty(t, content)

	var decoded Metrics
	err = decode(content, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, decoded.Name, "test")
	assert.Equal(t, decoded.Port, 9090)
	assert.Equal(t, decoded.Endpoint, "/metrics")
}
