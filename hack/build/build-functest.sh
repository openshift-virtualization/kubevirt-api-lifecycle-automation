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

set -euo pipefail

script_dir="$(cd "$(dirname "$0")" && pwd -P)"
source "${script_dir}"/common.sh

mkdir -p ${TESTS_OUT_DIR}/
# use vendor
export GO111MODULE=${GO111MODULE:-off}
go build -C vendor/github.com/onsi/ginkgo/v2/ginkgo -o /usr/bin/ginkgo
ginkgo build ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/tests/
mv ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/tests/tests.test ${TESTS_OUT_DIR}/
cp -f /usr/bin/ginkgo ${TESTS_OUT_DIR}/
