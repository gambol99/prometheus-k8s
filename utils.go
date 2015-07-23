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
	"net/url"
	"strconv"
	"os"
)

func validateURL(location string) (err error) {
	_, err = url.Parse(location)
	return
}

func getEnvString(key, value string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return value
}

func getEnvInt(key string, default_value int) int {
	if os.Getenv(key) != "" {
		value, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			return default_value
		}
		return value
	}
	return default_value
}
