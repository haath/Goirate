#!/bin/sh

VERSION_FILE="cmd/goirate/version.go"

GIT_SHA=$(git rev-parse --short HEAD)

VERSION=${CI_COMMIT_TAG:-$GIT_SHA}

sed -i '$ d' $VERSION_FILE

echo "const VERSION = \"${VERSION}\"" >> $VERSION_FILE