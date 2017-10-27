#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

REPO_NAME="github.com/r2d4/crd"
CRD_NAMESPACE="r2d4.com"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ${GOPATH}/src/k8s.io/code-generator)}

${CODEGEN_PKG}/generate-groups.sh all \
  ${REPO_NAME}/pkg/client ${REPO_NAME}/pkg/apis \
  ${CRD_NAMESPACE}:v1 \
