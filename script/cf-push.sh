#!/bin/bash

set -ex

name="wildcat-countdown-dev"
if [[ ! -z "$1" ]]; then
    name="$1"
fi

if [[ "$name" == "wildcat-countdown" ]]; then
    manifest="./untracked/manifest.yml"
else
    manifest="./untracked/manifest-dev.yml"
fi

dir=`mktemp -d`
GOOS=linux go build -o "$dir/main" .
cf push -p "$dir" -f $manifest $name
rm -rf "$dir"
