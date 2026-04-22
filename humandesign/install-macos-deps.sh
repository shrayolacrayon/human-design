#!/bin/bash
# Installs the Swiss Ephemeris C library (libswe) on macOS for local development.
# Not needed when using Docker.

set -e

echo "Installing Swiss Ephemeris (libswe) for macOS..."

WORK_DIR=$(mktemp -d)
trap "rm -rf $WORK_DIR" EXIT

cd "$WORK_DIR"

# Detect prefix (Apple Silicon Homebrew vs Intel)
if [ -d /opt/homebrew ]; then
    PREFIX=/opt/homebrew
else
    PREFIX=/usr/local
fi

# Download Swiss Ephemeris source (latest stable)
VERSION="2.10.03"
curl -L "https://www.astro.com/ftp/sweph/src/sweph-${VERSION}.tar.gz" -o sweph.tar.gz
tar xzf sweph.tar.gz
cd sweph-*

# Build shared library
gcc -O2 -fPIC -shared -o libswe.dylib \
    sweph.c swephlib.c swejpl.c swedate.c swecl.c swehel.c swevents.c \
    -lm -install_name "${PREFIX}/lib/libswe.dylib"

# Install headers and library
mkdir -p "${PREFIX}/include" "${PREFIX}/lib"
cp swephexp.h sweph.h "${PREFIX}/include/"
cp libswe.dylib "${PREFIX}/lib/"

# Create symlinks
ln -sf "${PREFIX}/lib/libswe.dylib" "${PREFIX}/lib/libswe.1.dylib"

echo ""
echo "Done! libswe installed to ${PREFIX}"
echo ""
echo "Now run from the humandesign/ directory:"
echo "  go run ./cmd/server/"
