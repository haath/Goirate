#!/bin/bash

FILE_NAME="goirate-$CI_COMMIT_TAG"

if [[ $1 == *"win"* ]] then
    FILE_NAME="$FILE_NAME.exe"
fi

PUBLISH_URL="https://api.bintray.com/content/gmantaos/Goirate/$1/$CI_COMMIT_TAG/$1/$FILE_NAME?publish=1&override=1"

echo $PUBLISH_URL

curl -T build/goirate -ugmantaos:$BINTRAY_API_KEY $PUBLISH_URL

echo ""