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
KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE_NAME=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_IMAGE_NAME:-kubevirt-api-lifecycle-automation}

DOCKER_PREFIX=${DOCKER_PREFIX:-"quay.io/bmordeha"}
DOCKER_TAG=${DOCKER_TAG:-latest}
VERBOSITY=${VERBOSITY:-1}
PULL_POLICY=${PULL_POLICY:-Always}
KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE:-kubevirt-api-lifecycle-automation}
TARGET_NAMESPACE=${TARGET_NAMESPACE:-""}
RESTART_REQUIRED=${RESTART_REQUIRED:-true}
LABEL_SELECTOR=${LABEL_SELECTOR:-""}
MACHINE_TYPE_GLOBAL=${MACHINE_TYPE_GLOBAL:-"smth-glob9.10.0"}
CR_NAME=${CR_NAME:-kubevirt-api-lifecycle-automation}

function parseTestOpts() {
    pkgs=""
    test_args=""
    while [[ $# -gt 0 ]] && [[ $1 != "" ]]; do
        case "${1}" in
        --test-args=*)
            test_args="${1#*=}"
            shift 1
            ;;
        ./*...)
            pkgs="${pkgs} ${1}"
            shift 1
            ;;
        *)
            echo "ABORT: Unrecognized option \"$1\""
            exit 1
            ;;
        esac
    done
}