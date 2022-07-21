#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

"$(dirname "${BASH_SOURCE[0]}")"/generate-groups.sh all \
  github.com/eddiezane/that-conference-k8s-controller/pkg/generated \
  github.com/eddiezane/that-conference-k8s-controller/pkg/apis \
  "pictures:v1" \
  --output-base=. \
  --trim-path-prefix "github.com/eddiezane/that-conference-k8s-controller"
