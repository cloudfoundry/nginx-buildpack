#!/usr/bin/env bash
set -euo pipefail

export ROOT=`dirname $(readlink -f ${BASH_SOURCE%/*})`
$ROOT/scripts/install_tools.sh

cd $ROOT/src/nginx/
ginkgo -r -skipPackage=brats,integration
