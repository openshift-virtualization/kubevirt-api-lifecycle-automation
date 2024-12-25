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

set -e
script_dir="$(cd "$(dirname "$0")" && pwd -P)"
source "${script_dir}"/common.sh
source "${script_dir}"/config.sh

mkdir -p "${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/_out"

# update this whenever new builder tag is created
BUILDER_IMAGE=${BUILDER_IMAGE:-quay.io/bmordeha/kubevirt-api-lifecycle-automation-builder:2412110934-d596255}

BUILDER_VOLUME="kubevirt-api-lifecycle-automation-volume"
DOCKER_CA_CERT_FILE="${DOCKER_CA_CERT_FILE:-}"
DOCKERIZED_CUSTOM_CA_PATH="/etc/pki/ca-trust/source/anchors/custom-ca.crt"

DISABLE_SECCOMP=${DISABLE_SECCOMP:-}

SYNC_OUT=${SYNC_OUT:-true}
SYNC_VENDOR=${SYNC_VENDOR:-false}

# Create the persistent docker volume
if [ -z "$(${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} volume list | grep ${BUILDER_VOLUME})" ]; then
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} volume create ${BUILDER_VOLUME}
fi

# Make sure that the output directory exists
echo "Making sure output directory exists..."
${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} run -v "${BUILDER_VOLUME}:/root:rw,z" --security-opt label=disable $DISABLE_SECCOMP --rm --entrypoint "/entrypoint.sh" ${BUILDER_IMAGE} mkdir -p /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation/_out

${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} run -v "${BUILDER_VOLUME}:/root:rw,z" --security-opt label=disable $DISABLE_SECCOMP --rm --entrypoint "/entrypoint.sh" ${BUILDER_IMAGE} git config --global --add safe.directory /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation
echo "Starting rsyncd"
# Start an rsyncd instance and make sure it gets stopped after the script exits
RSYNC_CID_KUBEVIRT_API_LIFECYCLE_AUTOMATION=$(${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} run -d -v "${BUILDER_VOLUME}:/root:rw,z" --security-opt label=disable $DISABLE_SECCOMP --cap-add SYS_CHROOT --expose 873 -P --entrypoint "/entrypoint.sh" ${BUILDER_IMAGE} /usr/bin/rsync --no-detach --daemon --verbose)

function finish() {
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} stop --time 1 ${RSYNC_CID_KUBEVIRT_API_LIFECYCLE_AUTOMATION} >/dev/null 2>&1
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} rm -f ${RSYNC_CID_KUBEVIRT_API_LIFECYCLE_AUTOMATION} >/dev/null 2>&1
}
trap finish EXIT

RSYNCD_PORT=$(${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} port $RSYNC_CID_KUBEVIRT_API_LIFECYCLE_AUTOMATION | cut -d':' -f2)

rsynch_fail_count=0

while ! rsync ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/ "rsync://root@127.0.0.1:${RSYNCD_PORT}/build/" &>/dev/null; do
    if [[ "$rsynch_fail_count" -eq 0 ]]; then
        printf "Waiting for rsyncd to be ready"
        sleep .1
    elif [[ "$rsynch_fail_count" -lt 30 ]]; then
        printf "."
        sleep 1
    else
        printf "failed"
        break
    fi
    rsynch_fail_count=$((rsynch_fail_count + 1))
done

printf "\n"

rsynch_fail_count=0

_rsync() {
    rsync -al "$@"
}

echo "Rsyncing ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR} to container"

# Copy kubevirt-api-lifecycle-automation into the persistent docker volume
_rsync \
    --delete \
    --exclude 'cluster-up/cluster/**/.kubectl' \
    --exclude 'cluster-up/cluster/**/.oc' \
    --exclude 'cluster-up/cluster/**/.kubeconfig' \
    --exclude ".vagrant" \
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/ \
    "rsync://root@127.0.0.1:${RSYNCD_PORT}/build"

# Run the command
test -t 1 && USE_TTY="-it"
if ! $KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI exec -w /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation ${USE_TTY} ${RSYNC_CID_KUBEVIRT_API_LIFECYCLE_AUTOMATION} /entrypoint.sh "$@"; then
    # Copy the build output out of the container, make sure that _out exactly matches the build result
    if [ "$SYNC_OUT" = "true" ]; then
        _rsync --delete "rsync://root@127.0.0.1:${RSYNCD_PORT}/out" ${OUT_DIR}
    fi
    exit 1
fi

# Copy the whole kubevirt-api-lifecycle-automation data out to get generated sources and formatting changes
_rsync \
    --exclude 'cluster-up/cluster/**/.kubectl' \
    --exclude 'cluster-up/cluster/**/.oc' \
    --exclude 'cluster-up/cluster/**/.kubeconfig' \
    --exclude 'pkg/client-go/kubevirt/clientset/versioned/typed/core/v1/generated_expansion.go' \
    --exclude "_out" \
    --exclude "bin" \
    --exclude "vendor" \
    --exclude ".vagrant" \
    --exclude ".git" \
    "rsync://root@127.0.0.1:${RSYNCD_PORT}/build" \
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/

if [ "$SYNC_VENDOR" = "true" ] && [ -n $VENDOR_DIR ]; then
    _rsync --delete "rsync://root@127.0.0.1:${RSYNCD_PORT}/vendor" "${VENDOR_DIR}/"
fi


# Copy the build output out of the container, make sure that _out exactly matches the build result
if [ "$SYNC_OUT" = "true" ]; then
    _rsync --delete "rsync://root@127.0.0.1:${RSYNCD_PORT}/out" ${OUT_DIR}
fi