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
 * Copyright 2020 Red Hat, Inc.
 *
 */

package libvmifact

import (
	kvirtv1 "kubevirt.io/api/core/v1"

	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/libvmi"
)

func NewGuestless(opts ...libvmi.Option) *kvirtv1.VirtualMachineInstance {
	opts = append(
		[]libvmi.Option{libvmi.WithResourceMemory(qemuMinimumMemory())},
		opts...)
	return libvmi.New(opts...)
}

func qemuMinimumMemory() string {
	if isARM64() {
		// required to start qemu on ARM with UEFI firmware
		// https://github.com/kubevirt/kubevirt/pull/11366#issuecomment-1970247448
		const armMinimalBootableMemory = "128Mi"
		return armMinimalBootableMemory
	}
	return "1Mi"
}
