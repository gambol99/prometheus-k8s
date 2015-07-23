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
	"flag"
	"fmt"
)

//
// Config ... Configuration for the service
type Config struct {
	// the url to the api proxy
	Host string
	// the port the api proxy is running
	Port int
	// the kubernetes namespace to listen in
	Namespace string
	// the kubernetes token file if any
	TokenFile string
	// the cert used to verify to kubernetes
	CaCertFile string
	// the metrics annotation used
	MetricAnnotation string
	// the filename of the nodes yaml
	NodesConfigFilename string
	// the filename of the pods yaml
	PodsConfigFilename string
	// the directory to save the configuration
	ConfigDirectory string
	// the refresh interval
	RefreshInterval int
	// the api version
	APIVersion string
	// the protocol to use when connecting to the api
	APIProtocol string
	// toggle to indicate if we shoud add all the kubernetes nodes as targets
	WithNodes bool
	// a toggle to produce the endpoints for pods
	WithPods bool
	// a dry run - i.e. only display to screen
	DryRun bool
	// Insure https
	HttpInsecure bool
}

var (
	config Config
)

func init() {
	flag.StringVar(&config.Host, "api", getEnvString("KUBERNETES_SERVICE_HOST", "127.0.0.1"), "the host / ip address the kubectl proxy is running")
	flag.StringVar(&config.NodesConfigFilename, "node-file", "nodes.yml", "the filename of the nodes yaml file")
	flag.StringVar(&config.PodsConfigFilename, "pod-file", "pods.yml", "the filename of of the pods yaml")
	flag.StringVar(&config.APIVersion, "api-version", "v1", "the protocol to use when connecting to the api")
	flag.StringVar(&config.APIProtocol, "api-protocol", "http", "the kubernetes api version to use")
	flag.StringVar(&config.ConfigDirectory, "config", ".", "the directory save the genrated files into")
	flag.StringVar(&config.MetricAnnotation, "metrics", METRICS_ANNOTATION, "the tag used in the pods annotations")
	flag.StringVar(&config.TokenFile, "bearer-token-file", "", "The file containing the bearer token.")
	flag.StringVar(&config.CaCertFile, "ca-cert-file", "", "The file containing the CA certificate.")
	flag.BoolVar(&config.HttpInsecure, "insecure", true, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure.")
	flag.IntVar(&config.Port, "port", getEnvInt("KUBERNETES_SERVICE_PORT", 8001), "the port the api proxy is running on")
	flag.IntVar(&config.RefreshInterval, "interval", 300, "the refresh interval in seconds that we perform a forced refresh")
	flag.BoolVar(&config.WithNodes, "nodes", true, "generate the metric endpoints for all kubernetes nodes in the cluster")
	flag.BoolVar(&config.WithPods, "pods", true, "generate the metric endpoints for pods which container prometheus endpoints")
	flag.BoolVar(&config.DryRun, "dry-run", false, "perform a dry run, display output to screen only")
}

func parseConfig() error {
	// step: parse the command line arguments
	flag.Parse()
	// step: validate we have everything we need to proceed
	location := fmt.Sprintf("%s://%s:%d", config.APIProtocol, config.Host, config.Port)
	// check: ensure the location is valid
	if err := validateURL(location); err != nil {
		return fmt.Errorf("invalid URL specified, please check the url and port, error: %s", err)
	}
	return nil
}
