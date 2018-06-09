package app

import (
	"fmt"

	"github.com/chriswalker/version-status-cli/internal/output"
	"github.com/chriswalker/version-status-cli/pkg/kubernetes"
	"github.com/pkg/errors"
)

type App struct {
	configPath string
}

type services struct {
	context  string
	services map[string]string
	err      error
}

func (s services) Error() error {
	return s.err
}

func NewApp(kubeConfigFilepath string) *App {
	return &App{
		configPath: kubeConfigFilepath,
	}
}

func (a *App) GetVersionStatus(contexts []string) {
	c := make(chan services, 2)

	for _, context := range contexts {
		go a.getServices(context, c)
	}

	results := make(map[string]services, 0)
	for i := 0; i < cap(c); i++ {
		result := <-c
		if result.Error() != nil {
			fmt.Println(result.Error())
			return
		}
		results[result.context] = result
	}

	versions := a.processResults(results[contexts[0]].services, results[contexts[1]].services)

	outputter := output.NewStdOutputter()
	outputter.Output(contexts, versions)
}

func (a *App) getServices(context string, result chan<- services) {
	client, err := kubernetes.NewKubernetesClient(context, a.configPath)
	if err != nil {
		result <- services{
			err: errors.Wrap(err, "could not create k8s client"),
		}
		return
	}

	pods, err := client.GetPods()
	if err != nil {
		result <- services{
			err: errors.Wrapf(err, "could not get list of pods for context '%s'", context),
		}
		return
	}

	result <- services{
		context:  context,
		services: pods,
	}
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
