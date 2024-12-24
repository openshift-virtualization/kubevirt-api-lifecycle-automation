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

determine_kubevirt_api_lifecycle_automation_bin() {
    if [ "${KUBEVIRTCI_RUNTIME-}" = "podman" ]; then
        echo podman
    elif [ "${KUBEVIRTCI_RUNTIME-}" = "docker" ]; then
        echo docker
    else
        if docker ps >/dev/null 2>&1; then
            echo docker
        elif podman ps >/dev/null 2>&1; then
            echo podman
        else
            echo ""
        fi
    fi
}


KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR="$(cd $(dirname $0)/../../ && pwd -P)"


BIN_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/bin
OUT_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/_out
CMD_OUT_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/cmd
TESTS_OUT_DIR=${OUT_DIR}/tests
BUILD_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/hack/build
MANIFEST_TEMPLATE_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/manifests/templates
MANIFEST_GENERATED_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/manifests/generated
CACHE_DIR=${OUT_DIR}/gocache
VENDOR_DIR=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/vendor
ARCHITECTURE="${BUILD_ARCH:-$(uname -m)}"
HOST_ARCHITECTURE="$(uname -m)"
KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI="$(determine_kubevirt_api_lifecycle_automation_bin)"


