package tests

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/kubevirt-api-lifecycle-automation/resources/operator"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/libvmi"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/util"
	libvirtxml "github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/capabilities"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/framework"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/libnode"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/libvmifact"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	virt1 "kubevirt.io/api/core/v1"
	"sync"
	"time"
)

var _ = Describe("Kubevirt Api Lifecycle Automation tests", func() {
	f := framework.NewFramework("job-test")
	var kubevirtNS string
	var once sync.Once
	Context("Kubevirt Api Lifecycle Automation Job", func() {
		BeforeEach(func() {
			once.Do(func() {
				c, err := f.GetVirtClient()
				Expect(err).ToNot(HaveOccurred())
				kv, err := c.KubevirtV1().KubeVirts(v1.NamespaceAll).List(context.Background(), v1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(kv.Items).To(HaveLen(1))
				kubevirtNS = kv.Items[0].Namespace
			})
		})

		It("Simple job test", func(ctx context.Context) {
			c, err := f.GetVirtClient()
			Expect(err).ToNot(HaveOccurred())
			nodeList := getAllSchedulableNodes(f.K8sClient)
			Expect(len(nodeList.Items)).Should(BeNumerically(">", 0))
			targetNode := nodeList.Items[0]
			defaultMachineType := getDefaultMachineType(f, targetNode.Name)
			supportedMachineTypes, err := getSupportedMachineTypes(f, targetNode.Name, kubevirtNS)
			Expect(err).ToNot(HaveOccurred())
			var targetMachineType string
			for _, machine := range supportedMachineTypes {
				if machine.Deprecated != "yes" && machine.Name != "" && machine.Name != defaultMachineType {
					targetMachineType = machine.Name
					break
				}
			}
			vmi := libvmifact.NewGuestless(
				libvmi.WithNamespace(f.Namespace.Name),
				libvmi.WithNodeSelectorFor(targetNode.Name),
				libvmi.WithMachineType(targetMachineType),
			)
			vm := libvmi.NewVirtualMachine(
				vmi,
				libvmi.WithRunStrategy(virt1.RunStrategyAlways),
			)
			vm, err = c.KubevirtV1().VirtualMachines(f.Namespace.Name).Create(context.Background(), vm, v1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Eventually(thisVM(f, vm)).WithTimeout(120 * time.Second).WithPolling(time.Second).Should(beReady())
			activateKubevirtApiLifecycleAutomationJob(f, "*rhel*", "true")
			Eventually(thisVM(f, vm)).WithTimeout(20 * time.Second).WithPolling(time.Second).Should(HaveMachineType(defaultMachineType))
		})
	})
})

// thisVM fetches the latest state of the VirtualMachine. If the object does not exist, nil is returned.
func thisVM(f *framework.Framework, vm *virt1.VirtualMachine) func() (*virt1.VirtualMachine, error) {
	return thisVMWith(f, vm.Namespace, vm.Name)
}

// thisVMWith fetches the latest state of the VirtualMachine based on namespace and name. If the object does not exist, nil is returned.
func thisVMWith(f *framework.Framework, namespace string, name string) func() (*virt1.VirtualMachine, error) {
	return func() (p *virt1.VirtualMachine, err error) {
		virtClient, err := f.GetVirtClient()
		Expect(err).ToNot(HaveOccurred())
		p, err = virtClient.KubevirtV1().VirtualMachines(namespace).Get(context.Background(), name, v1.GetOptions{})
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return
	}
}

func beReady() gomegatypes.GomegaMatcher {
	return gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Status": gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Ready": BeTrue(),
		}),
	}))
}

func HaveMachineType(machineType string) gomegatypes.GomegaMatcher {
	return gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Spec": gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Template": gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Spec": gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Domain": gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Machine": gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
							"Type": Equal(machineType),
						})),
					}),
				}),
			})),
		}),
	}))
}

// getAllSchedulableNodes returns list of Nodes which are "KubeVirt" schedulable.
func getAllSchedulableNodes(virtClient *kubernetes.Clientset) *k8sv1.NodeList {
	nodeList, err := virtClient.CoreV1().Nodes().List(context.Background(), v1.ListOptions{
		LabelSelector: virt1.NodeSchedulable + "=" + "true",
	})
	Expect(err).ToNot(HaveOccurred(), "Should list compute nodeList")
	return nodeList
}

// getSupportedMachineTypes returns list of supported machine types of a node
func getSupportedMachineTypes(f *framework.Framework, targetNodeName string, kubevirtNS string) ([]libvirtxml.CapsGuestMachine, error) {
	var machines []libvirtxml.CapsGuestMachine
	var capabilities *libvirtxml.Caps
	cmd := []string{"bash", "-c", fmt.Sprintf("cat /var/lib/kubevirt-node-labeller/capabilities.xml")}
	stdout, stderr, err := libnode.ExecuteCommandOnNodeThroughVirtHandler(f, targetNodeName, kubevirtNS, cmd)
	Expect(err).ToNot(HaveOccurred())
	Expect(stderr).To(BeEmpty())
	Expect(stdout).ToNot(BeEmpty())
	err = xml.Unmarshal([]byte(stdout), &capabilities)
	if err != nil {
		return nil, err
	}
	for _, guest := range capabilities.Guests {
		machines = append(machines, guest.Arch.Machines...)
	}
	return machines, nil
}

// getDefaultMachineType the default machinetype for a node
func getDefaultMachineType(f *framework.Framework, targetNodeName string) string {
	c, err := f.GetVirtClient()
	Expect(err).ToNot(HaveOccurred())
	vmi := libvmifact.NewGuestless(
		libvmi.WithNamespace(f.Namespace.Name),
		libvmi.WithNodeSelectorFor(targetNodeName),
	)
	vm := libvmi.NewVirtualMachine(
		vmi,
		libvmi.WithRunStrategy(virt1.RunStrategyAlways),
	)
	vm, err = c.KubevirtV1().VirtualMachines(f.Namespace.Name).Create(context.Background(), vm, v1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())
	err = c.KubevirtV1().VirtualMachines(f.Namespace.Name).Delete(context.Background(), vm.Name, v1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())

	return vm.Spec.Template.Spec.Domain.Machine.Type
}

func activateKubevirtApiLifecycleAutomationJob(f *framework.Framework, machineTypeGlob, restartRequired string) {
	newJob := operator.CreateKubevirtApiLifecycleAutomation(f.KubevirtApiLifecycleAutomationNamespace, machineTypeGlob, v1.NamespaceAll, restartRequired, "", "", f.KubevirtApiLifecycleAutomationImage, "Always")
	newJob.ObjectMeta.Name = "" // Clear the name field as we are using GenerateName
	newJob.ObjectMeta.GenerateName = util.KubevirtApiLifecycleAutomationResourceName + "-"
	newJob.Spec.Suspend = nil
	newJob, err := f.K8sClient.BatchV1().Jobs(f.KubevirtApiLifecycleAutomationNamespace).Create(context.Background(), newJob, v1.CreateOptions{})
	defer func() {
		err = f.K8sClient.BatchV1().Jobs(f.KubevirtApiLifecycleAutomationNamespace).Delete(context.Background(), newJob.Name, v1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred())
	}()
	Expect(err).ToNot(HaveOccurred())
	By("Waiting for the job to complete...")
	Eventually(func() int32 {
		newJob, err = f.K8sClient.BatchV1().Jobs(f.KubevirtApiLifecycleAutomationNamespace).Get(context.Background(), newJob.Name, v1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		return newJob.Status.Succeeded
	}, 180*time.Second, 5*time.Second).Should(BeNumerically(">", 0))
}
