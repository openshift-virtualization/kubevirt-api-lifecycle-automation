#!/usr/bin/env bash

set -e

source hack/build/common.sh
go build -C vendor/github.com/onsi/ginkgo/v2/ginkgo -o /usr/bin/ginkgo

# Find every folder containing tests
for dir in $(find ${KUBEVIRT_API_LIFECYCLE_AUTOMATION_DIR}/pkg/ -type f -name '*_test.go' -printf '%h\n' | sort -u); do
    # If there is no file ending with _suite_test.go, bootstrap ginkgo
    SUITE_FILE=$(find $dir -maxdepth 1 -type f -name '*_suite_test.go')
    if [ -z "$SUITE_FILE" ]; then
        (cd $dir && /usr/bin/ginkgo bootstrap || :)
    fi
done
