package client

//go:generate mockgen -source $GOFILE -package=$GOPACKAGE -destination=generated_mock_$GOFILE

/*
 ATTENTION: Rerun code generators when interface signatures are modified.
*/

import (
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/client-go/kubevirt/clientset/versioned"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubevirtApiLifecycleAutomationClient interface {
	RestClient() *rest.RESTClient
	kubernetes.Interface
	KubevirtClient() versioned.Interface
	DiscoveryClient() discovery.DiscoveryInterface
	Config() *rest.Config
}

type kubevirtCli struct {
	master          string
	kubeconfig      string
	restClient      *rest.RESTClient
	config          *rest.Config
	kubevirtClient  *versioned.Clientset
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   dynamic.Interface
	*kubernetes.Clientset
}

func (k kubevirtCli) KubevirtClient() versioned.Interface {
	return k.kubevirtClient
}

func (k kubevirtCli) Config() *rest.Config {
	return k.config
}

func (k kubevirtCli) RestClient() *rest.RESTClient {
	return k.restClient
}
func (k kubevirtCli) DiscoveryClient() discovery.DiscoveryInterface {
	return k.discoveryClient
}
