package tests_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/flags"
	"github.com/onsi/ginkgo/v2"
	ginkgo_reporters "github.com/onsi/ginkgo/v2/reporters"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	qe_reporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/framework"
	corev1 "k8s.io/api/core/v1"
)

const (
	pollInterval     = 2 * time.Second
	nsDeletedTimeout = 270 * time.Second
)

var (
	afterSuiteReporters = []Reporter{}
)

var (
	kubectlPath          = flag.String("kubectl-path-kubevirt-api-lifecycle-automation", "kubectl", "The path to the kubectl binary")
	ocPath               = flag.String("oc-path-kubevirt-api-lifecycle-automation", "oc", "The path to the oc binary")
	kubevirtApiLifecycleAutomationInstallNs = flag.String("kubevirt-api-lifecycle-automation-namespace", "kubevirt-api-lifecycle-automation", "The namespace of the Kubevirt Api Lifecycle Automation")
	kubevirtApiLifecycleAutomationImage     = flag.String("kubevirt-api-lifecycle-automation-image", "kubevirt-api-lifecycle-automation", "The image of Kubevirt Api Lifecycle Automation")
	kubeConfig           = flag.String("kubeconfig-kubevirt-api-lifecycle-automation", "/var/run/kubernetes/admin.kubeconfig", "The absolute path to the kubeconfig file")
	kubeURL              = flag.String("kubeurl", "", "kube URL url:port")
	goCLIPath            = flag.String("gocli-path-kubevirt-api-lifecycle-automation", "cli.sh", "The path to cli script")
	dockerPrefix         = flag.String("docker-prefix", "", "The docker host:port")
	dockerTag            = flag.String("docker-tag", "", "The docker tag")
)

// kubevirtApiLifecycleAutomationFailHandler call ginkgo.Fail with printing the additional information
func kubevirtApiLifecycleAutomationFailHandler(message string, callerSkip ...int) {
	if len(callerSkip) > 0 {
		callerSkip[0]++
	}
	ginkgo.Fail(message, callerSkip...)
}

func TestTests(t *testing.T) {
	defer GinkgoRecover()
	RegisterFailHandler(kubevirtApiLifecycleAutomationFailHandler)
	BuildTestSuite()
	RunSpecs(t, "Tests Suite")
}

// To understand the order in which things are run, read http://onsi.github.io/ginkgo/#understanding-ginkgos-lifecycle
// flag parsing happens AFTER ginkgo has constructed the entire testing tree. So anything that uses information from flags
// cannot work when called during test tree construction.
func BuildTestSuite() {
	BeforeSuite(func() {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Reading parameters\n")
		flags.NormalizeFlags()
		// Read flags, and configure client instances
		framework.ClientsInstance.KubectlPath = *kubectlPath
		framework.ClientsInstance.OcPath = *ocPath
		framework.ClientsInstance.KubevirtApiLifecycleAutomationNamespace = *kubevirtApiLifecycleAutomationInstallNs
		framework.ClientsInstance.KubeConfig = *kubeConfig
		framework.ClientsInstance.KubeURL = *kubeURL
		framework.ClientsInstance.GoCLIPath = *goCLIPath
		framework.ClientsInstance.DockerPrefix = *dockerPrefix
		framework.ClientsInstance.DockerTag = *dockerTag
		framework.ClientsInstance.KubevirtApiLifecycleAutomationImage = *kubevirtApiLifecycleAutomationImage

		fmt.Fprintf(ginkgo.GinkgoWriter, "Kubectl path: %s\n", framework.ClientsInstance.KubectlPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "OC path: %s\n", framework.ClientsInstance.OcPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "KubevirtApiLifecycleAutomation install NS: %s\n", framework.ClientsInstance.KubevirtApiLifecycleAutomationNamespace)
		fmt.Fprintf(ginkgo.GinkgoWriter, "Kubeconfig: %s\n", framework.ClientsInstance.KubeConfig)
		fmt.Fprintf(ginkgo.GinkgoWriter, "KubeURL: %s\n", framework.ClientsInstance.KubeURL)
		fmt.Fprintf(ginkgo.GinkgoWriter, "GO CLI path: %s\n", framework.ClientsInstance.GoCLIPath)
		fmt.Fprintf(ginkgo.GinkgoWriter, "DockerPrefix: %s\n", framework.ClientsInstance.DockerPrefix)
		fmt.Fprintf(ginkgo.GinkgoWriter, "DockerTag: %s\n", framework.ClientsInstance.DockerTag)

		restConfig, err := framework.ClientsInstance.LoadConfig()
		if err != nil {
			// Can't use Expect here due this being called outside of an It block, and Expect
			// requires any calls to it to be inside an It block.
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to load RestConfig err:%v", err))
		}

		framework.ClientsInstance.RestConfig = restConfig
		// clients
		kcs, err := framework.ClientsInstance.GetKubeClient()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to create K8SClient: %v", err))
		}
		framework.ClientsInstance.K8sClient = kcs

		crClient, err := framework.ClientsInstance.GetCrClient()
		if err != nil {
			Fail(fmt.Sprintf("ERROR, unable to create CrClient: %v", err))
		}
		framework.ClientsInstance.CrClient = crClient

		dyn, err := framework.ClientsInstance.GetDynamicClient()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("ERROR, unable to create DynamicClient: %v", err))
		}
		framework.ClientsInstance.DynamicClient = dyn

		if qe_reporters.Polarion.Run {
			afterSuiteReporters = append(afterSuiteReporters, &qe_reporters.Polarion)
		}
		if qe_reporters.JunitOutput != "" {
			afterSuiteReporters = append(afterSuiteReporters, qe_reporters.NewJunitReporter())
		}

	})

	AfterSuite(func() {
		k8sClient := framework.ClientsInstance.K8sClient
		Eventually(func() []corev1.Namespace {
			nsList, _ := k8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: framework.NsPrefixLabel})
			return nsList.Items
		}, nsDeletedTimeout, pollInterval).Should(BeEmpty())
	})

	var _ = ReportAfterSuite("TestTests", func(report Report) {
		for _, reporter := range afterSuiteReporters {
			ginkgo_reporters.ReportViaDeprecatedReporter(reporter, report)
		}
	})

}
