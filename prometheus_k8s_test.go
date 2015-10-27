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

func newTestPrometheusK8S(t *testing.T) *PrometheusK8S {
	fakeAPI := newFakeKubeAPI(t)
	return &PrometheusK8S{
		client:    fakeAPI,
		updatesCh: make(UpdateEvent, 10),
	}
}

func TestGeneratePodsConfiguration(t *testing.T) {
	ks8 := newTestPrometheusK8S(t)
	content, err := ks8.generatePodsConfiguration()
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
	t.Logf("pod config:\n%s", content)
}

func TestGenerateNodesConfiguration(t *testing.T) {
	ks8 := newTestPrometheusK8S(t)
	content, err := ks8.generateNodesConfiguration()
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
	t.Logf("node config:\n%s", content)
}
