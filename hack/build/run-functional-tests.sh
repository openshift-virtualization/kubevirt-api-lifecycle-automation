#!/usr/bin/env bash

#Copyright 2025 The KubevirtApiLifecycleAutomation Authors.
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

set -eo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd -P)"
source hack/build/config.sh
source hack/build/common.sh
source cluster-up/hack/common.sh

KUBEVIRTCI_CONFIG_PATH="$(
    cd "$(dirname "$BASH_SOURCE[0]")/../../"
    echo "$(pwd)/_ci-configs"
)"

# functional testing
BASE_PATH=${KUBEVIRTCI_CONFIG_PATH:-$PWD}
KUBECONFIG=${KUBECONFIG:-$BASE_PATH/$KUBEVIRT_PROVIDER/.kubeconfig}
GOCLI=${GOCLI:-${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/cluster-up/cli.sh}
KUBE_URL=${KUBE_URL:-""}
KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE:-kubevirt-api-lifecycle-automation}
KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE:-registry:5000/kubevirt-api-lifecycle-automation:latest}

OPERATOR_CONTAINER_IMAGE=$(./cluster-up/kubectl.sh get job -n $KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE kubevirt-api-lifecycle-automation -o'custom-columns=spec:spec.template.spec.containers[0].image' --no-headers)
DOCKER_PREFIX=${OPERATOR_CONTAINER_IMAGE%/*}
DOCKER_TAG=${OPERATOR_CONTAINER_IMAGE##*:}

if [ -z "${KUBECTL+x}" ]; then
    kubevirtci_kubectl="${BASE_PATH}/${KUBEVIRT_PROVIDER}/.kubectl"
    if [ -e ${kubevirtci_kubectl} ]; then
        KUBECTL=${kubevirtci_kubectl}
    else
        KUBECTL=$(which kubectl)
    fi
fi

# parsetTestOpts sets 'pkgs' and test_args
parseTestOpts "${@}"

arg_kubeurl="${KUBE_URL:+-kubeurl=$KUBE_URL}"
arg_namespace="${KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE:+-kubevirt-api-lifecycle-automation-namespace=$KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE}"
arg_kubevirt_api_lifecycle_automation_image="${KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE:+-kubevirt-api-lifecycle-automation-image=$KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE}"
arg_kubeconfig_kubevirt_api_lifecycle_automation="${KUBECONFIG:+-kubeconfig-kubevirt-api-lifecycle-automation=$KUBECONFIG}"
arg_kubeconfig="${KUBECONFIG:+-kubeconfig=$KUBECONFIG}"
arg_kubectl="${KUBECTL:+-kubectl-path-kubevirt-api-lifecycle-automation=$KUBECTL}"
arg_oc="${KUBECTL:+-oc-path-kubevirt-api-lifecycle-automation=$KUBECTL}"
arg_gocli="${GOCLI:+-gocli-path-kubevirt-api-lifecycle-automation=$GOCLI}"
arg_docker_prefix="${DOCKER_PREFIX:+-docker-prefix=$DOCKER_PREFIX}"
arg_docker_tag="${DOCKER_TAG:+-docker-tag=$DOCKER_TAG}"

test_args="${test_args}  -ginkgo.v  ${arg_kubeurl} ${arg_namespace} ${arg_kubevirt_api_lifecycle_automation_image} ${arg_kubeconfig} ${arg_kubeconfig_kubevirt_api_lifecycle_automation} ${arg_kubectl} ${arg_oc} ${arg_gocli} ${arg_docker_prefix} ${arg_docker_tag}"

(
    export TESTS_WORKDIR=${AAQ_DIR}/tests
    ginkgo_args="--trace --timeout=8h --v"
    ${TESTS_OUT_DIR}/ginkgo ${ginkgo_args} ${TESTS_OUT_DIR}/tests.test -- ${test_args}
)