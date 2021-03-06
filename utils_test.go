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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvString(t *testing.T) {
	assert.NotEmpty(t, getEnvString("USER", ""))
	assert.Equal(t, getEnvString("NOTHING", "TEST"), "TEST")
}

func TestGetEnvInt(t *testing.T) {
	value := getEnvInt("TEST_NUMBER", 4)
	assert.Equal(t, value, 4)
	os.Setenv("TEST_NUMBER", "10")
	value = getEnvInt("TEST_NUMBER", 4)
	assert.Equal(t, value, 10)
	os.Setenv("TEST_NUMBER", "10")
}
