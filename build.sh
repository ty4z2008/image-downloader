#!/bin/bash

if [ $# -lt 1 ]; then
	echo "go build script need entry directory"
	exit 1
fi

# Save the pwd before we run anything
PRE_PWD=`pwd`

# Determine the build script's actual directory, following symlinks
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
BUILD_DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"

# Derive the project name from the directory
PROJECT="$(basename $BUILD_DIR)"

# Build the project
if [ ! -d "bin" ]; then
	mkdir -p bin
fi
GOOS=darwin GOARCH=amd64 go build  -o "bin/${PROJECT}-osx" $1
GOOS=linux GOARCH=amd64 go build  -o "bin/${PROJECT}-linux" $1
GOOS=windows GOARCH=amd64 go build  -o "bin/${PROJECT}-windows" $1

EXIT_STATUS=$?

if [ $EXIT_STATUS == 0 ]; then
  echo "Build succeeded"
else
  echo "Build failed"
fi

# Change back to where we were
cd $PRE_PWD

exit $EXIT_STATUS