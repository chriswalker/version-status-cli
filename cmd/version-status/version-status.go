/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	"github.com/fatih/color"
	"github.com/pkg/errors"
	types "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	stagingServices, _ := getServices("staging", *kubeconfig)
	fmt.Println("Staging services")
	for service, version := range stagingServices {
		fmt.Printf("[%s] %s\n", color.BlueString(service), version)
	}

	prodServices, _ := getServices("production", *kubeconfig)
	fmt.Println("\n\nProd services")
	for service, version := range prodServices {
		fmt.Printf("[%s] %s\n", color.BlueString(service), version)
	}

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	/*
		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}
	*/
}

func getServices(k8sContext, configPath string) (map[string]string, error) {
	config, err := buildConfigFromFlags(k8sContext, configPath)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: "version",
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not get pods for context '%s'", k8sContext)
	}

	services := cleanupList(pods.Items)
	return services, nil
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Returns a de-duped, cleaned up & sorted list of services
// TODO, better way of doing this & doesn't support int deployments
func cleanupList(list []types.Pod) map[string]string {
	pods := make(map[string]string)
	for _, pod := range list {
		app := pod.ObjectMeta.Labels["app"]
		if _, ok := pods[app]; ok {
			continue
		}
		if (strings.HasPrefix(app, "api") || strings.HasPrefix(app, "svc")) && !strings.HasSuffix(app, "docs-site") {
			pods[app] = pod.ObjectMeta.Labels["version"]
		}
	}

	return pods
}

func getColourFunc(num int) func(format string, a ...interface{}) string {
	if num < 3 {
		return color.New(color.FgGreen).SprintfFunc()
	} else if num < 6 {
		return color.New(color.FgYellow).SprintfFunc()
	}
	return color.New(color.FgRed).SprintfFunc()
}
