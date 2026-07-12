#!/bin/bash

SCRIPT_DIR="$(dirname "$(realpath $0)")"
BUILD_DIR="$SCRIPT_DIR/build"
FRONTEND_DIR="$SCRIPT_DIR/content"
BACKEND_DIR="$SCRIPT_DIR"

if [ -d $BUILD_DIR ]; then
    rm -r $BUILD_DIR && echo "Deleting directory $BUILD_DIR"
fi
mkdir $BUILD_DIR && echo "Created directory $BUILD_DIR"

# Copy CSS and JS files to build dirbuild dir
cp -r "$FRONTEND_DIR/static" "$BUILD_DIR/static"

# Build server
go build -o "$BUILD_DIR" .
