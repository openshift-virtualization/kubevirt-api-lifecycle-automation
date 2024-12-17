package libnode

import (
	"context"
	"fmt"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/exec"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/framework"
	k8sv1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	virtv1 "kubevirt.io/api/core/v1"
)

func ExecuteCommandOnNodeThroughVirtHandler(f *framework.Framework, nodeName string, namespace string, command []string) (stdout, stderr string, err error) {
	virtHandlerPod, found, err := getVirtHandler(f, nodeName, namespace)
	if !found && err == nil {
		err = fmt.Errorf(fmt.Sprintf("Virt handler was not found on node:%v in namespace: %v", nodeName, namespace))
	}
	if err != nil {
		return "", "", err
	}
	return exec.ExecuteCommandOnPodWithResults(f, virtHandlerPod, "virt-handler", command)
}

func getVirtHandler(f *framework.Framework, nodeName string, namespace string) (*k8sv1.Pod, bool, error) {
	handlerNodeSelector := fields.ParseSelectorOrDie("spec.nodeName=" + nodeName)
	labelSelector, err := labels.Parse(virtv1.AppLabel + " in (virt-handler)")
	if err != nil {
		return nil, false, err
	}
	cli, err := f.GetKubeClient()
	if err != nil {
		return nil, false, err
	}
	pods, err := cli.CoreV1().Pods(namespace).List(context.Background(),
		k8smetav1.ListOptions{
			FieldSelector: handlerNodeSelector.String(),
			LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, false, err
	}
	if len(pods.Items) > 1 {
		return nil, false, fmt.Errorf("Expected to find one Pod, found %d Pods", len(pods.Items))
	}

	if len(pods.Items) == 0 {
		return nil, false, nil
	}
	return &pods.Items[0], true, nil
}
