#!/bin/bash

set -e

name="wildcat-countdown"
if [[ ! -z "$1" ]]; then
    name="$1"
fi

dir=`mktemp -d`
GOOS=linux go build -o "$dir/main" .
cf push -b binary_buildpack -c './main' -p "$dir" "$name"
rm -rf "$dir"