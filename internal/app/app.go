package app

import (
	"fmt"
	"sync"

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

	c := make(chan map[string]string)
	var wg sync.WaitGroup

	for _, env := range envs {
		wg.Add(1)
		go a.getServices(env, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	results := make([]map[string]string, 0)
	for result := range c {
		results = append(results, result)
	}

	// TODO - resuls may come back in different order; need to tag them
	versions := a.processResults(results[0], results[1])

	outputter := output.NewStdOutputter()
	outputter.Output(versions)

	// Output - TODO: move into output package
	/*
		for _, version := range versions {
			fn := getColourFunc(version)
			fmt.Printf("%s\n", fn("[%s] Staging: %s - Prod: %s", version.ServiceName, version.StagingVersion, version.ProdVersion))
		}
	*/
}

func (a *App) getServices(env string, result chan<- map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

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
