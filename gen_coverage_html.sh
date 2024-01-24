#!/usr/bin/bash

set -o errexit   # abort on nonzero exitstatus
set -o nounset   # abort on unbound variable
set -o pipefail  # don't hide errors within pipes

export COVDATA="$(pwd)/covdata"

mkdir -p ${COVDATA}
mkdir -p $(pwd)/coverage

go test -v \
  -cover \
  -covermode=atomic \
  -coverpkg=go.l0nax.org/typact \
  ./... -args -test.gocoverdir="${COVDATA}"

(
  cd ./testing/option
  go test -v \
    -cover \
    -covermode=atomic \
    -coverpkg=go.l0nax.org/typact \
    ./... -args -test.gocoverdir="${COVDATA}"
)

###### Generating report
mkdir -p ./coverage
go tool covdata merge -i="${COVDATA}" -o ./coverage

go tool covdata percent -i ./coverage
go tool cover -html=coverage.txt

rm -rf ${COVDATA}
rm -rf ./coverage
rm -rf ./unit-tests.xml
rm -rf ./coverage.txt
rm -rf ./coverage.xml
