#!/bin/bash

set -e

dir=`mktemp -d`
GOOS=linux go build -o "$dir/main" .
cf push -b binary_buildpack -c './main' -p "$dir" wildcat-countdown
rm -rf "$dir"
