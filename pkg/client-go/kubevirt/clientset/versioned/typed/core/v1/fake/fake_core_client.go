/*
Copyright 2025 The KubevirtApiLifecycleAutomation Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/client-go/kubevirt/clientset/versioned/typed/core/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKubevirtV1 struct {
	*testing.Fake
}

func (c *FakeKubevirtV1) KubeVirts(namespace string) v1.KubeVirtInterface {
	return &FakeKubeVirts{c, namespace}
}

func (c *FakeKubevirtV1) VirtualMachines(namespace string) v1.VirtualMachineInterface {
	return &FakeVirtualMachines{c, namespace}
}

func (c *FakeKubevirtV1) VirtualMachineInstances(namespace string) v1.VirtualMachineInstanceInterface {
	return &FakeVirtualMachineInstances{c, namespace}
}

func (c *FakeKubevirtV1) VirtualMachineInstanceMigrations(namespace string) v1.VirtualMachineInstanceMigrationInterface {
	return &FakeVirtualMachineInstanceMigrations{c, namespace}
}

func (c *FakeKubevirtV1) VirtualMachineInstancePresets(namespace string) v1.VirtualMachineInstancePresetInterface {
	return &FakeVirtualMachineInstancePresets{c, namespace}
}

func (c *FakeKubevirtV1) VirtualMachineInstanceReplicaSets(namespace string) v1.VirtualMachineInstanceReplicaSetInterface {
	return &FakeVirtualMachineInstanceReplicaSets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKubevirtV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
