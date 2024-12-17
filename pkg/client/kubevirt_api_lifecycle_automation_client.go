package client

import (
	"flag"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/client-go/kubevirt/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

var (
	kubeconfig string
	master     string
)

var (
	SchemeBuilder  runtime.SchemeBuilder
	Scheme         *runtime.Scheme
	Codecs         serializer.CodecFactory
	ParameterCodec runtime.ParameterCodec
)

func init() {
	// This allows consumers of the KubeVirt client go package to
	// customize what version the client uses. Without specifying a
	// version, all versions are registered. While this techincally
	// file to register all versions, so k8s ecosystem libraries
	// do not work well with this. By explicitly setting the env var,
	// consumers of our client go can avoid these scenarios by only
	// registering a single version
	Scheme = runtime.NewScheme()
	AddToScheme := SchemeBuilder.AddToScheme
	Codecs = serializer.NewCodecFactory(Scheme)
	ParameterCodec = runtime.NewParameterCodec(Scheme)
	AddToScheme(Scheme)
	AddToScheme(scheme.Scheme)
}

type RestConfigHookFunc func(*rest.Config)

var restConfigHooks []RestConfigHookFunc
var restConfigHooksLock sync.Mutex

var kubevirtApiLifecycleAutomationclient KubevirtApiLifecycleAutomationClient
var once sync.Once

// Init adds the default `kubeconfig` and `master` flags. It is not added by default to allow integration into
// the different controller generators which normally add these flags too.
func Init() {
	if flag.CommandLine.Lookup("kubeconfig") == nil {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	if flag.CommandLine.Lookup("master") == nil {
		flag.StringVar(&master, "master", "", "master url")
	}
}

func executeRestConfigHooks(config *rest.Config) {
	restConfigHooksLock.Lock()
	defer restConfigHooksLock.Unlock()

	for _, hookFn := range restConfigHooks {
		hookFn(config)
	}
}

func GetKubevirtApiLifecycleAutomationClientFromRESTConfig(config *rest.Config) (KubevirtApiLifecycleAutomationClient, error) {
	shallowCopy := *config
	shallowCopy.GroupVersion = &v1.SchemeGroupVersion
	shallowCopy.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: Codecs}
	shallowCopy.APIPath = "/apis"
	shallowCopy.ContentType = runtime.ContentTypeJSON
	if config.UserAgent == "" {
		config.UserAgent = restclient.DefaultKubernetesUserAgent()
	}

	executeRestConfigHooks(&shallowCopy)

	restClient, err := rest.RESTClientFor(&shallowCopy)
	if err != nil {
		return nil, err
	}

	coreClient, err := kubernetes.NewForConfig(&shallowCopy)
	if err != nil {
		return nil, err
	}

	kubevirtClient, err := versioned.NewForConfig(&shallowCopy)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(&shallowCopy)
	if err != nil {
		return nil, err
	}

	return &kubevirtCli{
		master,
		kubeconfig,
		restClient,
		&shallowCopy,
		kubevirtClient,
		discoveryClient,
		dynamicClient,
		coreClient,
	}, nil
}

func GetKubevirtApiLifecycleAutomationClientFromFlags(master string, kubeconfig string) (KubevirtApiLifecycleAutomationClient, error) {
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, err
	}
	return GetKubevirtApiLifecycleAutomationClientFromRESTConfig(config)
}

func GetKubevirtApiLifecycleAutomationClient() (KubevirtApiLifecycleAutomationClient, error) {
	var err error
	once.Do(func() {
		kubevirtApiLifecycleAutomationclient, err = GetKubevirtApiLifecycleAutomationClientFromFlags(master, kubeconfig)
	})
	return kubevirtApiLifecycleAutomationclient, err
}
