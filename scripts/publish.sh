#!/bin/bash

curl -T build/goirate -ugmantaos:$BINTRAY_API_KEY https://api.bintray.com/content/gmantaos/Goirate/$1/$CI_COMMIT_TAG/goirate?publish=1