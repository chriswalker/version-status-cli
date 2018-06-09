package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chriswalker/version-status-cli/internal/app"
	//_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	first  = flag.String("first", "", "(required) first k8s context")
	second = flag.String("second", "", "(required) second k8s context")
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
	app.GetVersionStatus([]string{*first, *second})
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
