#!/usr/bin/bash

set -o errexit  # abort on nonzero exitstatus
set -o nounset  # abort on unbound variable
set -o pipefail # don't hide errors within pipes

export COVDATA="$(pwd)/covdata"
export COVERPKG="go.l0nax.org/typact,go.l0nax.org/typact/std,go.l0nax.org/typact/std/option,go.l0nax.org/typact/std/exp,go.l0nax.org/typact/std/exp/cmpop,go.l0nax.org/typact/std/exp/pred,go.l0nax.org/typact/std/exp/xslices,go.l0nax.org/typact/std/xhash"

function clean_data {
  rm -rf ${COVDATA}
  rm -rf ./coverage
  rm -rf ./unit-tests.xml
  rm -rf ./coverage.txt
  rm -rf ./coverage.xml
}

clean_data
mkdir -p ${COVDATA}
mkdir -p $(pwd)/coverage

go test -v \
  -cover \
  -covermode=atomic \
  -coverpkg=${COVERPKG} \
  ./ -args -test.gocoverdir="${COVDATA}"

go test -v \
  -cover \
  -covermode=atomic \
  -coverpkg=${COVERPKG} \
  ./std/... -args -test.gocoverdir="${COVDATA}"

(
  cd ./testing/option
  go test -v \
    -cover \
    -covermode=atomic \
    -coverpkg=${COVERPKG} \
    ./... -args -test.gocoverdir="${COVDATA}"
)

(
  cd ./testing/std/exp
  go test -v \
    -cover \
    -covermode=atomic \
    -coverpkg=${COVERPKG} \
    ./... -args -test.gocoverdir="${COVDATA}"
)

(
  cd ./std/option
  go test -v \
    -cover \
    -covermode=atomic \
    -coverpkg=${COVERPKG} \
    ./... -args -test.gocoverdir="${COVDATA}"
)

(
  cd ./std
  go test -v \
    -cover \
    -covermode=atomic \
    -coverpkg=${COVERPKG} \
    ./... -args -test.gocoverdir="${COVDATA}"
)

###### Generating report
mkdir -p ./coverage
go tool covdata merge -i="${COVDATA}" -o ./coverage

go tool covdata percent -i ./coverage
go tool covdata textfmt -i ./coverage -o ./coverage.txt
go tool cover -html=coverage.txt

clean_data
