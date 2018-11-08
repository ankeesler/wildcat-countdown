#!/bin/bash

set -e
set -o pipefail

`dirname $0`/cf-push.sh wildcat-countdown-test
curl wildcat-countdown-test.cfapps.io > /dev/null

cf logs wildcat-countdown-test --recent | grep "hello, tuna" > /dev/null
