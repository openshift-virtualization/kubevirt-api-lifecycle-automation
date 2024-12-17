package util

import (
	utils "kubevirt.io/controller-lifecycle-operator-sdk/pkg/sdk/resources"
)

const (
	// KubevirtApiLifecycleAutomationLabel is the labe applied to all non operator resources
	KubevirtApiLifecycleAutomationLabel = "kubevirt-api-lifecycle-automation.io"
	// AppKubernetesManagedByLabel is the Kubernetes recommended managed-by label
	AppKubernetesManagedByLabel = "app.kubernetes.io/managed-by"
	// AppKubernetesComponentLabel is the Kubernetes recommended component label
	AppKubernetesComponentLabel = "app.kubernetes.io/component"
	KubevirtApiLifecycleAutomationResourceName     = "kubevirt-api-lifecycle-automation"
)

var commonLabels = map[string]string{
	KubevirtApiLifecycleAutomationLabel:            "",
	AppKubernetesManagedByLabel: "kubevirt-api-lifecycle-automation",
	AppKubernetesComponentLabel: "virtualization",
}

var JobLabels = map[string]string{
	"kubevirt-api-lifecycle-automation.io": "",
	"tier":            "node",
}

// ResourceBuilder helps in creating k8s resources
var ResourceBuilder = utils.NewResourceBuilder(commonLabels, JobLabels)
