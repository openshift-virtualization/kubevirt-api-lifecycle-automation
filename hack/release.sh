#!/bin/bash

source $(dirname "$0")/build/common.sh
source $(dirname "$0")/build/config.sh

set -e

if [ -z "${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI}" ]; then
    echo >&2 "no working container runtime found. Neither docker nor podman seems to work."
    exit 1
fi

if [ -z ${GPG_PRIVATE_KEY_FILE} ]; then
    echo "GPG_PRIVATE_KEY_FILE env var must be set"
    exit 1
elif [ -z ${GPG_PASSPHRASE_FILE} ]; then
    echo "GPG_PASSPHRASE_FILE env var must be set"
    exit 1
elif [ -z ${GITHUB_API_TOKEN_FILE} ]; then
    echo "GITHUB_API_TOKEN_FILE env var must be set"
    exit 1
fi

GIT_USER=${GIT_USER:-$(git config user.name)}
GIT_EMAIL=${GIT_EMAIL:-$(git config user.email)}

echo "git user:  $GIT_USER"
echo "git email: $GIT_EMAIL"

${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} pull quay.io/kubevirtci/release-tool:latest

echo "${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} run -it --rm \
-v ${GPG_PRIVATE_KEY_FILE}:/home/releaser/gpg-private \
-v ${GPG_PASSPHRASE_FILE}:/home/releaser/gpg-passphrase \
-v ${GITHUB_API_TOKEN_FILE}:/home/releaser/github-api-token \
-v /home/.gitconfig:/home/releaser/.gitconfig \
quay.io/kubevirtci/release-tool:latest \
--org=kubevirt \
--repo=kubevirt-api-lifecycle-automation \
--git-email \"${GIT_EMAIL}\" \
--git-user \"${GIT_USER}\"
\"$@\""

${KUBEVIRT_API_LIFECYCLE_AUTOMATION_CRI} run -it --rm \
    -v ${GPG_PRIVATE_KEY_FILE}:/home/releaser/gpg-private \
    -v ${GPG_PASSPHRASE_FILE}:/home/releaser/gpg-passphrase \
    -v ${GITHUB_API_TOKEN_FILE}:/home/releaser/github-api-token \
    -v /home/.gitconfig:/home/.gitconfig \
    quay.io/kubevirtci/release-tool:latest \
    --org=kubevirt \
    --repo=kubevirt-api-lifecycle-automation \
    --git-email "${GIT_EMAIL}" \
    --git-user "${GIT_USER}" \
    "$@"