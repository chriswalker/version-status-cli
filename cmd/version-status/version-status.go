package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chriswalker/version-status-cli/internal/app"
)

var (
	first     = flag.String("first", "", "(required) first k8s context")
	second    = flag.String("second", "", "(required) second k8s context")
	diffsonly = flag.Bool("diffsonly", false, "(optional) whether to report only services with version differences, or all services")
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if *first == "" || *second == "" {
		fmt.Printf("%s %s\n", *first, *second)
		flag.Usage()
		os.Exit(1)
	}
	app := app.NewApp(*kubeconfig)
	app.GetVersionStatus([]string{*first, *second}, *diffsonly)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
