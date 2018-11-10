#!/bin/bash

set -e

name="wildcat-countdown"
if [[ ! -z "$1" ]]; then
    name="$1"
fi

dir=`mktemp -d`
GOOS=linux go build -o "$dir/main" .
cf push -p "$dir" -f ./untracked/manifest.yml
rm -rf "$dir"
