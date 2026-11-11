#!/usr/bin/env bash

# Enable xtrace if the DEBUG environment variable is set
if [[ ${DEBUG-} =~ ^(1|yes|true)$ ]]; then
  set -o xtrace # Trace the execution of the script (debug)
fi

# No errexit: each case expects some runs to fail, and the harness must report every case.
set -o nounset  # Disallow expansion of unset variables
set -o pipefail # Use last non-zero exit code in a pipeline

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly SCRIPT_DIR
readonly CHECK="${SCRIPT_DIR}/check-node-pins.sh"

failures=0
cases=0
root=
workdir=

# Starts a case in its own empty directory under root.
new_workdir() {
  cases=$((cases + 1))
  workdir="${root}/${cases}"
  mkdir -p "${workdir}/.github/workflows"
}

node_from() {
  printf 'FROM --platform=%s node:%s@sha256:8996bb58 AS client\n' "\$BUILDPLATFORM" "$1"
}

# Writes Dockerfile.ci and ci.yaml into a new workdir; an empty argument omits that pin.
fixture() {
  local image=$1 ci=$2

  new_workdir
  {
    printf 'FROM alpine:3.24@sha256:be26560 AS base\n'
    [[ -z $image ]] || node_from "$image"
  } > "${workdir}/Dockerfile.ci"
  {
    printf '      - uses: actions/setup-node@820762786026740c76f36085b0efc47a31fe5020\n'
    printf '        env:\n'
    [[ -z $ci ]] || printf '          NODE_VERSION: %s\n' "$ci"
  } > "${workdir}/.github/workflows/ci.yaml"
}

# Runs the check in workdir and compares its exit status and stderr with what the case expects.
expect() {
  local name=$1 want_status=$2 want_err=$3
  local err status

  err=$(cd "$workdir" && "$CHECK" 2>&1 > /dev/null)
  status=$?
  if [[ $status -ne $want_status || $err != *"$want_err"* ]]; then
    printf 'FAIL %s: status %s (want %s), stderr: %s\n' "$name" "$status" "$want_status" "$err"
    failures=$((failures + 1))
  else
    printf 'ok   %s\n' "$name"
  fi
}

test_cases() {
  fixture 26.13.0-alpine 26.13.0
  expect "matching pins pass" 0 ""

  fixture 26.13.0-alpine '"26.13.0"'
  expect "quoted NODE_VERSION passes" 0 ""

  fixture 26.13.0-alpine 24.22.2
  expect "image ahead of CI fails" 1 \
    "node:26.13.0 but .github/workflows/ci.yaml tests on NODE_VERSION 24.22.2"

  fixture 26-alpine 26.13.0
  expect "floating image tag fails" 1 "node:26 but"

  fixture 26.13.0-alpine ""
  expect "missing NODE_VERSION fails" 1 "no NODE_VERSION found"

  fixture "" 26.13.0
  expect "missing node image fails" 1 "no node: image tag found"

  fixture 26.13.0-alpine 26.13.0
  node_from 24.22.2-alpine >> "${workdir}/Dockerfile.ci"
  expect "two node images fail" 1 "more than one node: image tag in Dockerfile.ci: 26.13.0, 24.22.2"

  new_workdir
  expect "missing files fail" 1 "cannot read Dockerfile.ci"

  new_workdir
  cp "${SCRIPT_DIR}/../Dockerfile.ci" "$workdir"
  cp "${SCRIPT_DIR}/../.github/workflows/ci.yaml" "${workdir}/.github/workflows"
  expect "this repository's pins pass" 0 ""
}

main() {
  root=$(mktemp -d)
  test_cases
  if ((failures > 0)); then
    printf '%d of %d case(s) failed\n' "$failures" "$cases"
    exit 1
  fi
}

# Invoke main with args if not sourced (https://stackoverflow.com/a/28776166/8787985)
if ! (return 0 2> /dev/null); then
  main "$@"
fi
