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

# Keep the pipefail here, or the unit tests won't return an error when needed.
#!/usr/bin/env bash
set -e
set -o pipefail

source /etc/profile.d/gimme.sh
export GOPATH="/root/go"
eval "$@"
