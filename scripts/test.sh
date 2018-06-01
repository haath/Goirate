#!/bin/bash

#############################################
#       Code coverage generation
#
# https://gitlab.com/pantomath-io/demo-tools
#############################################

COVERAGE_DIR="build/coverage"
OUTPUT="build/test_coverage.txt"
OUTPUT_HTML="build/test_coverage.html"
PKG_LIST=$(go list ./... | grep -v /vendor/)

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR" ;

# Create a coverage file for each package
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" ;
done ;

# Merge the coverage profile files
echo 'mode: count' > "${OUTPUT}" ;
tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> "${OUTPUT}" ;

# Remove temporary files
rm -rf $COVERAGE_DIR ;

# Generate the HTML coverage report
go tool cover -html="${OUTPUT}" -o "${OUTPUT_HTML}" ;

# Generate the text coverage report
go tool cover -func="${OUTPUT}" -o ${OUTPUT} ;

# Print out the coverage report
cat ${OUTPUT} ;