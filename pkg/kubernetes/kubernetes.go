package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	clientSet *kubernetes.Clientset
}

func NewKubernetesClient(k8sContext, configPath string) (*KubernetesClient, error) {
	config, err := buildConfigFromFlags(k8sContext, configPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not build k8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "could not build k8s clientset from config")
	}

	return &KubernetesClient{
		clientSet: clientset,
	}, nil
}

func (k *KubernetesClient) GetPods() (map[string]string, error) {
	pods, err := k.clientSet.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: "version",
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get pod list from Kubernetes")
	}

	services := cleanupList(pods.Items)
	return services, nil
}

// Returns a de-duped, cleaned up & sorted list of services
func cleanupList(list []types.Pod) map[string]string {
	pods := make(map[string]string)
	for _, pod := range list {
		app := pod.ObjectMeta.Labels["app"]
		if _, ok := pods[app]; ok {
			continue
		}
		// TODO - can refactor this
		if (strings.HasPrefix(app, "api") || strings.HasPrefix(app, "svc")) && !strings.HasSuffix(app, "docs-site") {
			pods[app] = pod.ObjectMeta.Labels["version"]
		}
	}

	return pods
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
