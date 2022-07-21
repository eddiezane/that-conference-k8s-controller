#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

base_dir="$(dirname "${BASH_SOURCE[0]}")/.."

rm -rf "${base_dir}/pkg/generated"
rm -rf "${base_dir}/pkg/apis/pictures/v1/zz_generated.deepcopy.go"
