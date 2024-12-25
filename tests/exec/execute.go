/*
 * This file is part of the kubevirt project
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
 * Copyright 2022 Red Hat, Inc.
 *
 */

package exec

import (
	"bytes"
	"github.com/kubevirt/kubevirt-api-lifecycle-automation/tests/framework"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecuteCommandOnPodWithResults(f *framework.Framework, pod *k8sv1.Pod, containerName string, command []string) (stdout, stderr string, err error) {
	var (
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)
	options := remotecommand.StreamOptions{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
		Tty:    false,
	}
	err = ExecuteCommandOnPodWithOptions(f, pod, containerName, command, options)
	return stdoutBuf.String(), stderrBuf.String(), err
}

func ExecuteCommandOnPodWithOptions(f *framework.Framework, pod *k8sv1.Pod, containerName string, command []string, options remotecommand.StreamOptions) error {
	cli, _ := f.GetKubeClient()
	req := cli.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Param("container", containerName)

	req.VersionedParams(&k8sv1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	virtConfig := f.RestConfig

	executor, err := remotecommand.NewSPDYExecutor(virtConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	return executor.Stream(options)
}
