package backend

import (
	"fmt"
	"log"
	"os"
	"strings"

	"io"

	"k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	coreclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/util/term"
)

type Backend struct {
	kubeConfig string
	config     *restclient.Config
	kubeClient *clientset.Clientset
	podClient  coreclient.CoreInterface
}

func New(kubeConfig string) *Backend {

	var (
		config *restclient.Config
		err    error
	)

	if _, err = os.Stat(kubeConfig); os.IsNotExist(err) { // if file is not exist, will build config in cluster
		if config, err = restclient.InClusterConfig(); err != nil {
			if err != nil {
				panic(fmt.Errorf("Build kubernetes config in cluster error: %v", err))
			}
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			err = fmt.Errorf("Build kubernetes config out of cluster error: %v", err)
			panic(err)
		}
	}

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &Backend{kubeConfig: kubeConfig, config: config, kubeClient: kubeClient, podClient: kubeClient.Core()}
}

func (b *Backend) GetNamespaces() ([]api.Namespace, error) {
	namespaceList, err := b.kubeClient.Core().Namespaces().List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaceList.Items, nil
}

func (b *Backend) GetPods(namespace string) ([]api.Pod, error) {
	if len(strings.TrimSpace(namespace)) == 0 {
		return nil, fmt.Errorf("GetPods error, namespace must be provide")
	}

	podList, err := b.kubeClient.Core().Pods(namespace).List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func (b *Backend) GetContainers(namespace, podName string) ([]api.Container, error) {
	if len(strings.TrimSpace(namespace)) == 0 {
		return nil, fmt.Errorf("GetPods error, namespace must be provide")
	}

	if len(strings.TrimSpace(podName)) == 0 {
		return nil, fmt.Errorf("GetContainers error, podName must be provide")
	}

	pod, err := b.kubeClient.Core().Pods(namespace).Get(podName)
	if err != nil {
		return nil, err
	}

	return pod.Spec.Containers, nil
}

func (b *Backend) ExecPod(namespace, podName, containerName string, commands []string, in io.Reader, out, err io.Writer, size *term.Size) {
	log.Println("=========>>  ExecPod ", namespace, podName, containerName)
	options := &ExecOptions{
		StreamOptions: StreamOptions{
			In: in, Out: out, Err: out,
			TTY: true, Stdin: true,
			PodName:       podName,
			ContainerName: containerName,
			Namespace:     namespace,
			Size:          size,
		},
		PodClient: b.podClient,
		Config:    b.config,
		Executor:  &DefaultRemoteExecutor{},
	}

	options.Complete(commands)
	if err := options.Validate(); err != nil {
		log.Println("[ERROR] ", err)
	}

	if err := options.Run(); err != nil {
		log.Println("[ERROR] ", err.Error(), (err == nil))
	}
}
