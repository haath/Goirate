#!/bin/bash

PUBLISH_URL="https://api.bintray.com/content/gmantaos/Goirate/$1/$CI_COMMIT_TAG/goirate-$1?publish=1&override=1"

echo $PUBLISH_URL

curl -T build/goirate -ugmantaos:$BINTRAY_API_KEY $PUBLISH_URL

echo ""