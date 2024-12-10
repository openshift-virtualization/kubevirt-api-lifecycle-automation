#!/bin/bash -e

echo kubevirt-api-lifecycle-automation

source ./hack/build/config.sh
source ./hack/build/common.sh
source ./cluster-up/hack/common.sh
source ./cluster-up/cluster/${KUBEVIRT_PROVIDER}/provider.sh

if [ "${KUBEVIRT_PROVIDER}" = "external" ]; then
   KUBEVIRT_API_LIFECYCLE_AUTOMATION_SYNC_PROVIDER="external"
else
   KUBEVIRT_API_LIFECYCLE_AUTOMATION_SYNC_PROVIDER="kubevirtci"
fi
source ./cluster-sync/${KUBEVIRT_API_LIFECYCLE_AUTOMATION_SYNC_PROVIDER}/provider.sh


KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE:-kubevirt-api-lifecycle-automation}
KUBEVIRT_API_LIFECYCLE_AUTOMATION_INSTALL_TIMEOUT=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_INSTALL_TIMEOUT:-120}
KUBEVIRT_API_LIFECYCLE_AUTOMATION_AVAILABLE_TIMEOUT=${KUBEVIRT_API_LIFECYCLE_AUTOMATION_AVAILABLE_TIMEOUT:-600}

# Set controller verbosity to 3 for functional tests.
export VERBOSITY=3

PULL_POLICY=${PULL_POLICY:-IfNotPresent}
# The default DOCKER_PREFIX is set to kubevirt and used for builds, however we don't use that for cluster-sync
# instead we use a local registry; so here we'll check for anything != "external"
# wel also confuse this by swapping the setting of the DOCKER_PREFIX variable around based on it's context, for
# build and push it's localhost, but for manifests, we sneak in a change to point a registry container on the
# kubernetes cluster.  So, we introduced this MANIFEST_REGISTRY variable specifically to deal with that and not
# have to refactor/rewrite any of the code that works currently.
MANIFEST_REGISTRY=$DOCKER_PREFIX

if [ "${KUBEVIRT_PROVIDER}" != "external" ]; then
  registry=${IMAGE_REGISTRY:-localhost:$(_port registry)}
  DOCKER_PREFIX=${registry}
  MANIFEST_REGISTRY="registry:5000"
fi

if [ "${KUBEVIRT_PROVIDER}" == "external" ]; then
  # No kubevirtci local registry, likely using something external
  if [[ $(${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} login --help | grep authfile) ]]; then
    registry_provider=$(echo "$DOCKER_PREFIX" | cut -d '/' -f 1)
    echo "Please log in to "${registry_provider}", bazel push expects external registry creds to be in ~/.docker/config.json"
    ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} login --authfile "${HOME}/.docker/config.json" $registry_provider
  fi
fi

# Need to set the DOCKER_PREFIX appropriately in the call to `make docker push`, otherwise make will just pass in the default `kubevirt`

DOCKER_PREFIX=$MANIFEST_REGISTRY PULL_POLICY=$PULL_POLICY make manifests
DOCKER_PREFIX=$DOCKER_PREFIX make push

function check_kubevirt_api_lifecycle_automation_exists() {
  # Check if the kubevirt-api-lifecycle-automation Job exists in the specified namespace
  kubectl get job kubevirt-api-lifecycle-automation -n "$KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE" &> /dev/null

  if [ $? -eq 0 ]; then
    echo "Job kubevirt-api-lifecycle-automation exists in namespace $KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE."
    return 0
  else
    echo "Job kubevirt-api-lifecycle-automation does not exist in namespace $KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE."
    return 1
  fi
}

function wait_kubevirt_api_lifecycle_automation_available {
  retry_count="${KUBEVIRT_API_LIFECYCLE_AUTOMATION_INSTALL_TIMEOUT}"
  echo "Waiting for kubevirt-api-lifecycle-automation Job in namespace '$KUBEVIRT_API_LIFECYCLE_AUTOMATION_NAMESPACE' to be ready..."

  # Loop for the specified number of retries
  for ((i = 0; i < retry_count; i++)); do
    # Check if DaemonSet pods are ready
    if check_kubevirt_api_lifecycle_automation_exists ; then
      echo "Job kubevirt-api-lifecycle-automation is available."
      exit 0
    fi

    # Wait for 1 second before retrying
    sleep 1
  done
    echo "Warning: kubevirt-api-lifecycle-automation doesn't exist!"
}

mkdir -p ./_out/tests

# Install KUBEVIRT API LIFECYCLE AUTOMATION JOB
install_kubevirt_api_lifecycle_automation

wait_kubevirt_api_lifecycle_automation_available
