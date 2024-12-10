#!/usr/bin/env bash

set -e

function install_kubevirt_api_lifecycle_automation {
    _kubectl apply -f "./_out/manifests/release/kubevirt-api-lifecycle-automation.yaml"
}

function delete_kubevirt_api_lifecycle_automation {
  if [ -f "./_out/manifests/release/kubevirt-api-lifecycle-automation.yaml" ]; then
    _kubectl delete --ignore-not-found -f "./_out/manifests/release/kubevirt-api-lifecycle-automation.yaml"
  else
    echo "File ./_out/manifests/release/kubevirt-api-lifecycle-automation.yaml does not exist."
  fi
}