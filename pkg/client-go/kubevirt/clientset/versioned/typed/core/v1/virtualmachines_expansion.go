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

package v1

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "kubevirt.io/api/core/v1"
)

const (
	cannotMarshalJSONErrFmt = "cannot Marshal to json: %s"
	vmSubresourceURLFmt     = "/apis/subresources.kubevirt.io/%s"
)

type KubeVirtExpansion interface{}

type VirtualMachineInstanceExpansion interface{}

type VirtualMachineInstanceMigrationExpansion interface{}

type VirtualMachineInstancePresetExpansion interface{}

type VirtualMachineInstanceReplicaSetExpansion interface{}

type VirtualMachineExpansion interface {
	Restart(ctx context.Context, name string, restartOptions *v1.RestartOptions) error
	Start(ctx context.Context, name string, startOptions *v1.StartOptions) error
	Stop(ctx context.Context, name string, stopOptions *v1.StopOptions) error
}

func (c *virtualMachines) Restart(ctx context.Context, name string, restartOptions *v1.RestartOptions) error {
	body, err := json.Marshal(restartOptions)
	if err != nil {
		return fmt.Errorf(cannotMarshalJSONErrFmt, err)
	}
	return c.client.Put().
		AbsPath(fmt.Sprintf(vmSubresourceURLFmt, v1.ApiStorageVersion)).
		Namespace(c.ns).
		Resource("virtualmachines").
		Name(name).
		SubResource("restart").
		Body(body).
		Do(ctx).
		Error()
}

func (c *virtualMachines) Start(ctx context.Context, name string, startOptions *v1.StartOptions) error {
	optsJson, err := json.Marshal(startOptions)
	if err != nil {
		return err
	}
	return c.client.Put().
		AbsPath(fmt.Sprintf(vmSubresourceURLFmt, v1.ApiStorageVersion)).
		Namespace(c.ns).
		Resource("virtualmachines").
		Name(name).
		SubResource("start").
		Body(optsJson).
		Do(ctx).
		Error()
}

func (c *virtualMachines) Stop(ctx context.Context, name string, stopOptions *v1.StopOptions) error {
	optsJson, err := json.Marshal(stopOptions)
	if err != nil {
		return err
	}
	return c.client.Put().
		AbsPath(fmt.Sprintf(vmSubresourceURLFmt, v1.ApiStorageVersion)).
		Namespace(c.ns).
		Resource("virtualmachines").
		Name(name).
		SubResource("stop").
		Body(optsJson).
		Do(ctx).
		Error()
}
