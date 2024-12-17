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
 * Copyright The Kubevirt Authors
 *
 */

package kubevirt_api_lifecycle_automation

import (
	"context"
	"fmt"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/client"
	envmanager "github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/kubevirt-api-lifecycle-automation/env-manager"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/patch"
	"os"
	"path"
	"strconv"

	"github.com/kubevirt/kubevirt-api-lifecycle-automation/pkg/log"
	"k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	virtv1 "kubevirt.io/api/core/v1"
)

const (
	machineTypeEnvName     = "MACHINE_TYPE_GLOB"
	namespaceEnvName       = "NAMESPACE"
	restartRequiredEnvName = "RESTART_REQUIRED"
	labelSelectorEnvName   = "LABEL_SELECTOR"
)

type MachineTypeUpdater struct {
	virtClient      client.KubevirtApiLifecycleAutomationClient
	machineTypeGlob string
	namespace       string
	labelSelector   labels.Selector
	restartRequired bool
}

var EnvVarManager envmanager.EnvVarManager = envmanager.EnvVarManagerImpl{}

func Execute() {
	log.InitializeLogging("machine-type-updater")

	virtCli, err := client.GetKubevirtApiLifecycleAutomationClient()
	if err != nil {
		log.Log.Errorf("Error retrieving virt client: %v", err)
		os.Exit(1)
	}

	app, err := NewMachineTypeUpdater(virtCli)
	if err != nil {
		os.Exit(1)
	}
	app.Run()

}

func NewMachineTypeUpdater(virtCli client.KubevirtApiLifecycleAutomationClient) (*MachineTypeUpdater, error) {
	updater := MachineTypeUpdater{
		virtClient: virtCli,
	}
	err := updater.initVariables()
	if err != nil {
		log.Log.Errorf("Error initializing variables: %v", err)
		return nil, err
	}

	return &updater, nil
}

func (c *MachineTypeUpdater) getNamespaces() ([]string, error) {
	if c.namespace != metav1.NamespaceAll {
		return []string{c.namespace}, nil
	}

	namespacesList, err := c.virtClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, len(namespacesList.Items))
	for i, ns := range namespacesList.Items {
		namespaces[i] = ns.Name
	}
	return namespaces, nil
}

func (c *MachineTypeUpdater) Run() {
	defer utilruntime.HandleCrash()

	log.Log.Info("Starting machine-type-updater")
	defer log.Log.Info("Shutting down machine-type-updater")

	namespaces, err := c.getNamespaces()
	if err != nil {
		log.Log.Errorf("Error while listing namespaces: %v", err)
		os.Exit(1)
	}

	for _, ns := range namespaces {
		vmList, err := c.virtClient.KubevirtClient().KubevirtV1().VirtualMachines(c.namespace).List(context.Background(), metav1.ListOptions{LabelSelector: c.labelSelector.String()})
		if err != nil {
			log.Log.Errorf("Error getting vm list: %s from namespace: %v", err.Error(), ns)
			continue
		}
		for _, vm := range vmList.Items {
			c.execute(&vm)
		}
	}
}

func (c *MachineTypeUpdater) execute(vm *virtv1.VirtualMachine) error {
	shouldUpdate, err := shouldUpdateMachineType(vm.Spec.Template.Spec.Domain.Machine, c.machineTypeGlob)
	if err != nil {
		// The only possible error is the bad pattern.
		// This should never happen since the pattern was
		// already checked before starting the controller.
		// In case something is changed, panic!
		panic(err)
	}

	if !shouldUpdate {
		return nil
	}

	err = c.patchMachineType(vm)
	if err != nil {
		log.Log.Errorf("Error patching vm %s/%s: %v\nSkipping...", vm.Namespace, vm.Name, err)
		return nil
	}

	runStrategy, err := vm.RunStrategy()
	if err != nil {
		log.Log.Errorf("Error getting RunStrategy from vm %s/%s: %v\nSkipping...", vm.Namespace, vm.Name, err)
		return nil
	}

	if runStrategy == virtv1.RunStrategyHalted || vm.Status.PrintableStatus != virtv1.VirtualMachineStatusRunning {
		return nil
	}

	if runStrategy == virtv1.RunStrategyOnce {
		return c.virtClient.KubevirtClient().KubevirtV1().VirtualMachines(vm.Namespace).Stop(context.Background(), vm.Name, &virtv1.StopOptions{})

	}

	return c.virtClient.KubevirtClient().KubevirtV1().VirtualMachines(vm.Namespace).Restart(context.Background(), vm.Name, &virtv1.RestartOptions{})
}

func (c *MachineTypeUpdater) patchMachineType(vm *virtv1.VirtualMachine) error {
	// removing the machine type field from the VM spec reverts it to
	// the default machine type of the VM's arch
	patches := []patch.PatchOperation{
		{
			Op:   patch.PatchRemoveOp,
			Path: "/spec/template/spec/domain/machine",
		},
	}

	payload, err := patch.GeneratePatchPayload(patches...)
	if err != nil {
		// This is a programmer's error and should not happen
		return fmt.Errorf("failed to generate patch payload: %v", err)
	}

	_, err = c.virtClient.KubevirtClient().KubevirtV1().VirtualMachines(vm.Namespace).Patch(context.Background(), vm.Name, types.JSONPatchType, payload, metav1.PatchOptions{})
	return err
}

func (c *MachineTypeUpdater) initVariables() error {
	c.labelSelector = labels.Everything()

	machineTypeEnvValue, exists := EnvVarManager.LookupEnv(machineTypeEnvName)
	if !exists {
		return fmt.Errorf("no machine type was specified")
	}

	// Execute a match with empty string to check if the pattern is correct
	_, err := path.Match(machineTypeEnvValue, "")
	if err != nil {
		return fmt.Errorf("syntax error in pattern of %s environment variable, value \"%s\"", machineTypeEnvName, machineTypeEnvValue)
	}

	c.machineTypeGlob = machineTypeEnvValue

	namespaceEnv, exists := EnvVarManager.LookupEnv(namespaceEnvName)
	if exists {
		errs := validation.ValidateNamespaceName(namespaceEnv, false)
		if errs != nil {
			return fmt.Errorf("syntax error in %s environment variable, value \"%s\": %v", namespaceEnvName, namespaceEnv, errs)
		}
		c.namespace = namespaceEnv
	}

	restartEnv, exists := EnvVarManager.LookupEnv(restartRequiredEnvName)
	if exists {
		restartRequired, err := strconv.ParseBool(restartEnv)
		if err != nil {
			return fmt.Errorf("error parsing %s environment variable, value \"%s\": %v", restartRequiredEnvName, restartEnv, err)
		}

		c.restartRequired = restartRequired
	}

	selectorEnv, exists := EnvVarManager.LookupEnv(labelSelectorEnvName)
	if exists {
		labelSelector, err := labels.Parse(selectorEnv)
		if err != nil {
			return fmt.Errorf("error parsing %s environment variable, value \"%s\": %v", labelSelectorEnvName, selectorEnv, err)
		}

		c.labelSelector = labelSelector
	}

	return nil
}

func shouldUpdateMachineType(currentMachine *virtv1.Machine, machineTypeGlob string) (bool, error) {
	if currentMachine == nil {
		return false, nil
	}

	return path.Match(machineTypeGlob, currentMachine.Type)
}
