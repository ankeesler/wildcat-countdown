#!/bin/bash

set -e

`dirname $0`/cf-push.sh wildcat-countdown-test
curl wildcat-countdown-test.cfapps.io
