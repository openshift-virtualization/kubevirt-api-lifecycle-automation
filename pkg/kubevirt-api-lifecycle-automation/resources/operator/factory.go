package operator

import (
	"fmt"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/kubevirt-api-lifecycle-automation/resources/args"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utils "kubevirt.io/controller-lifecycle-operator-sdk/pkg/sdk/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FactoryArgs contains the required parameters to generate all cluster-scoped resources
type FactoryArgs struct {
	NamespacedArgs args.FactoryArgs
	Image          string
}

type factoryFunc func(*FactoryArgs) []client.Object

func aggregateFactoryFunc(funcs ...factoryFunc) factoryFunc {
	return func(args *FactoryArgs) []client.Object {
		var result []client.Object
		for _, f := range funcs {
			result = append(result, f(args)...)
		}
		return result
	}
}

// CreateOperatorResourceGroup creates all cluster resources from a specific group/component
func CreateOperatorResourceGroup(group string, args *FactoryArgs) ([]client.Object, error) {
	f, ok := kubevirtApiLifecycleAutomationFactoryFunctions[group]
	if !ok {
		return nil, fmt.Errorf("group %s does not exist", group)
	}

	resources := f(args)
	for _, r := range resources {
		utils.ValidateGVKs([]runtime.Object{r})
	}
	return resources, nil
}

var kubevirtApiLifecycleAutomationFactoryFunctions = map[string]factoryFunc{
	"kubevirt-api-lifecycle-automation-cluster-rbac": createClusterRBAC,
	"kubevirt-api-lifecycle-automation-rbac":         createNamespacedRBAC,
	"kubevirt-api-lifecycle-automation":              createJob,
	"everything":                aggregateFactoryFunc(createClusterRBAC, createNamespacedRBAC, createJob),
}

// ClusterServiceVersionData - Data arguments used to create kubevirt api lifecycle automation job's CSV manifest
type ClusterServiceVersionData struct {
	CsvVersion         string
	ReplacesCsvVersion string
	Namespace          string
	ImagePullPolicy    string
	ImagePullSecrets   []corev1.LocalObjectReference
	IconBase64         string
	Verbosity          string
	OperatorVersion    string
	ControllerImage    string
	WebhookServerImage string
	OperatorImage      string
}
