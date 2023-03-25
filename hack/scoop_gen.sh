#!/usr/bin/env bash

set -o nounset
set -o pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

GITHUB_REF=${GITHUB_REF:-""}

if [[ -z "$GITHUB_REF" ]]; then
  echo "GITHUB_REF is empty"
  exit -1
fi

SCOOP_VERSION=${GITHUB_REF#refs/tags/}
SCOOP_VERSION=${SCOOP_VERSION#v}

SCOOP_HASH="$(curl -s https://github.com/go-musicfox/go-musicfox/releases/download/v${SCOOP_VERSION}/go-musicfox_${SCOOP_VERSION}_windows_amd64.zip | sha256sum)"

eval "cat <<EOF
$(< "$ROOT"/deploy/scoop/go-musicfox.json.tpl)
EOF
"  > "$ROOT"/deploy/scoop/go-musicfox.json