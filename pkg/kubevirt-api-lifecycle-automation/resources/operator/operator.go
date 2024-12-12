package operator

import (
	utils2 "github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/util"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/controller-lifecycle-operator-sdk/pkg/sdk/resources"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	roleName        = "kubevirt-api-lifecycle-automation"
	clusterRoleName = roleName + "-cluster"
)

func getClusterPolicyRules() []rbacv1.PolicyRule {
	rules := []rbacv1.PolicyRule{
		{
			APIGroups: []string{
				"kubevirt.io",
			},
			Resources: []string{
				"virtualmachines",
			},
			Verbs: []string{
				"list",
				"patch",
			},
		},
		{
			APIGroups: []string{
				"",
			},
			Resources: []string{
				"namespaces",
			},
			Verbs: []string{
				"list",
			},
		},
		{
			APIGroups: []string{
				virtv1.SubresourceGroupName,
			},
			Resources: []string{
				"virtualmachines/stop",
				"virtualmachines/restart",
			},
			Verbs: []string{
				"update",
			},
		},
	}
	rules = append(rules)
	return rules
}

func createClusterRole() *rbacv1.ClusterRole {
	return utils2.ResourceBuilder.CreateOperatorClusterRole(clusterRoleName, getClusterPolicyRules())
}

func createClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
	return utils2.ResourceBuilder.CreateOperatorClusterRoleBinding(utils2.KubevirtApiLifecycleAutomationResourceName, clusterRoleName, utils2.KubevirtApiLifecycleAutomationResourceName, namespace)
}

func createClusterRBAC(args *FactoryArgs) []client.Object {
	return []client.Object{
		createClusterRole(),
		createClusterRoleBinding(args.NamespacedArgs.Namespace),
	}
}
func createNamespacedRBAC(args *FactoryArgs) []client.Object {
	return []client.Object{
		createServiceAccount(args.NamespacedArgs.Namespace),
	}
}
func createServiceAccount(namespace string) *corev1.ServiceAccount {
	return utils2.ResourceBuilder.CreateOperatorServiceAccount(utils2.KubevirtApiLifecycleAutomationResourceName, namespace)
}

func createJob(args *FactoryArgs) []client.Object {
	return []client.Object{
		CreateKubevirtApiLifecycleAutomation(args.NamespacedArgs.Namespace,
			args.NamespacedArgs.MachineTypeGlob,
			args.NamespacedArgs.TargetNamespace,
			args.NamespacedArgs.RestartRequired,
			args.NamespacedArgs.LabelSelector,
			args.NamespacedArgs.Verbosity,
			args.Image,
			args.NamespacedArgs.PullPolicy),
	}
}

func createJobEnvVar(machineTypeGlob, targetNamespace, restartRequired, labelSelector, verbosity string) []corev1.EnvVar {
	envVar := []corev1.EnvVar{
		{
			Name:  "MACHINE_TYPE_GLOB",
			Value: machineTypeGlob,
		},
		{
			Name:  "RESTART_REQUIRED",
			Value: restartRequired,
		},
		{
			Name:  "VERBOSITY",
			Value: verbosity,
		},
	}
	if labelSelector != "" {
		envVar = append(envVar, corev1.EnvVar{Name: "LABEL_SELECTOR", Value: labelSelector})
	}
	if targetNamespace != "" {
		envVar = append(envVar, corev1.EnvVar{Name: "NAMESPACE", Value: targetNamespace})
	}
	return envVar
}

func CreateKubevirtApiLifecycleAutomation(namespace, machineTypeGlob, targetNamespace, restartRequired, labelSelector, verbosity, kubevirtApiLifecycleAutomationImage, pullPolicy string) *v1.Job {
	container := corev1.Container{
		Name:            "kubevirt-api-lifecycle-automation",
		Image:           kubevirtApiLifecycleAutomationImage,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("100m"),
				corev1.ResourceMemory: resource.MustParse("50M"),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Privileged:               ptr.To(false),
			AllowPrivilegeEscalation: ptr.To(false),
			RunAsNonRoot:             ptr.To(true),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			},
		},
	}
	container.Env = createJobEnvVar(machineTypeGlob, targetNamespace, restartRequired, labelSelector, verbosity)

	labels := resources.WithLabels(map[string]string{"name": "kubevirt-api-lifecycle-automation"}, utils2.JobLabels)
	cj := &v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubevirt-api-lifecycle-automation",
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.JobSpec{
			Suspend: ptr.To(true),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName:            "kubevirt-api-lifecycle-automation",
					HostPID:                       true,
					HostUsers:                     ptr.To(true),
					TerminationGracePeriodSeconds: ptr.To(int64(5)),
					Containers:                    []corev1.Container{container},
					PriorityClassName:             "system-node-critical",
					RestartPolicy:                 corev1.RestartPolicyNever,
				},
			},
		},
		Status: v1.JobStatus{},
	}

	return cj
}
