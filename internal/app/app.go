package app

import (
	"fmt"

	"github.com/chriswalker/version-status-cli/internal/output"
	"github.com/chriswalker/version-status-cli/pkg/kubernetes"
)

type App struct {
	configPath string
}

func NewApp(kubeConfigFilepath string) *App {
	return &App{
		configPath: kubeConfigFilepath,
	}
}

func (a *App) GetVersionStatus() {
	envs := []string{"staging", "production"}

	c := make(chan map[string]string, 2)

	for _, env := range envs {
		go a.getServices(env, c)
	}

	results := make([]map[string]string, 0)
	for i := 0; i < cap(c); i++ {
		result := <-c
		results = append(results, result)
	}

	// TODO - resuls may come back in different order; need to tag them
	versions := a.processResults(results[0], results[1])

	outputter := output.NewStdOutputter()
	outputter.Output(versions)
}

func (a *App) getServices(env string, result chan<- map[string]string) {
	client, err := kubernetes.NewKubernetesClient(env, a.configPath)
	if err != nil {
		// TODO
		fmt.Println(err)
	}

	pods, err := client.GetPods()
	if err != nil {
		// TODO
		fmt.Println(err)
	}

	result <- pods
}

func (a *App) processResults(staging, production map[string]string) []output.Version {
	var versions []output.Version

	for service, version := range staging {
		ver := output.Version{
			ServiceName:    service,
			StagingVersion: version,
		}
		if prodVersion, ok := production[service]; ok != false {
			ver.ProdVersion = prodVersion
		}

		versions = append(versions, ver)
	}

	return versions
}
