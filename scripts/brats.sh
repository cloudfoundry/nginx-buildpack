#!/usr/bin/env bash
set -euo pipefail

export ROOT="$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc
$ROOT/scripts/install_tools.sh

GINKGO_NODES=${GINKGO_NODES:-3}
GINKGO_ATTEMPTS=${GINKGO_ATTEMPTS:-1}
export CF_STACK=${CF_STACK:-cflinuxfs2}

cd $ROOT/src/nginx/brats
ginkgo -mod vendor -r --flakeAttempts=$GINKGO_ATTEMPTS -nodes $GINKGO_NODES
