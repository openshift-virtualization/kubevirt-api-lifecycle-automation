/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors
 *
 */

package fake

import (
	"context"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/client-go/testing"

	v1 "kubevirt.io/api/core/v1"
)

func (c *FakeVirtualMachines) Restart(ctx context.Context, name string, restartOptions *v1.RestartOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewPutSubresourceAction(virtualmachinesResource, c.ns, "restart", name, restartOptions), nil)

	return err
}

func (c *FakeVirtualMachines) Start(ctx context.Context, name string, startOptions *v1.StartOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewPutSubresourceAction(virtualmachinesResource, c.ns, "start", name, startOptions), nil)

	return err
}

func (c *FakeVirtualMachines) Stop(ctx context.Context, name string, stopOptions *v1.StopOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewPutSubresourceAction(virtualmachinesResource, c.ns, "stop", name, stopOptions), nil)

	return err
}
