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
	"os"
	"strconv"
)

// getEnvString get the value from the environment or use the default
func getEnvString(key, value string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}

	return value
}

// getEnvInt retrieve the value from the environment or use the default
func getEnvInt(key string, defaultValue int) int {
	if os.Getenv(key) != "" {
		value, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			return defaultValue
		}
		return value
	}

	return defaultValue
}

// writeConfigFile write the contents to stdout of a file
func writeConfigFile(content []byte, directory, filename string, dryRun bool) (err error) {
	var file *os.File

	if dryRun {
		file = os.Stdout
	}
	if !dryRun {
		file, err = os.Open(fmt.Sprintf("%s/%s", directory, filename))
		if err != nil {
			return
		}
	}

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}
