#!/usr/bin/env bash

# Enable xtrace if the DEBUG environment variable is set
if [[ ${DEBUG-} =~ ^(1|yes|true)$ ]]; then
  set -o xtrace # Trace the execution of the script (debug)
fi

set -o errexit  # Exit on most errors (see the manual)
set -o nounset  # Disallow expansion of unset variables
set -o pipefail # Use last non-zero exit code in a pipeline

readonly DOCKERFILE=Dockerfile.ci
readonly WORKFLOW=.github/workflows/ci.yaml
readonly IMAGE_PATTERN='^FROM[[:space:]].*[[:space:]]node:([^-@[:space:]]+)[-@[:space:]].*$'
readonly CI_PATTERN='^[[:space:]]+NODE_VERSION:[[:space:]]+"?([^"[:space:]#]+)"?.*$'

die() {
  printf 'check-node-pins: %s\n' "$*" >&2
  exit 1
}

# Prints the one value that pattern's first group captures in file; dies on none or several.
extract() {
  local file=$1 pattern=$2 what=$3
  local values

  [[ -r $file ]] || die "cannot read ${file}; run from the repository root"
  values=$(sed -nE "s/${pattern}/\1/p" "$file")
  [[ -n $values ]] || die "no ${what} found in ${file}"
  [[ $values != *$'\n'* ]] || die "more than one ${what} in ${file}: ${values//$'\n'/, }"
  printf '%s\n' "$values"
}

main() {
  local image ci

  image=$(extract "$DOCKERFILE" "$IMAGE_PATTERN" "node: image tag")
  ci=$(extract "$WORKFLOW" "$CI_PATTERN" "NODE_VERSION")

  if [[ $image != "$ci" ]]; then
    die "${DOCKERFILE} builds on node:${image} but ${WORKFLOW} tests on NODE_VERSION ${ci}." \
      "Set both to the same exact version; Renovate's node group can split them on a major" \
      "(knowledge/dependency-updates/knowledge.md)."
  fi
}

# Invoke main with args if not sourced (https://stackoverflow.com/a/28776166/8787985)
if ! (return 0 2> /dev/null); then
  main "$@"
fi
