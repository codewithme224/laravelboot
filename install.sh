#!/bin/bash

# install.sh for laravelboot
set -e

REPO="codewithme224/laravelboot"
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

echo "🔍 Detecting latest version..."
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "❌ Could not detect latest version."
    exit 1
fi

echo "📦 Downloading laravelboot $LATEST_TAG for $PLATFORM/$ARCH..."
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/laravelboot_${LATEST_TAG#v}_${PLATFORM}_${ARCH}.tar.gz"

curl -L -o laravelboot.tar.gz "$URL"

echo "📂 Extracting..."
tar -xzf laravelboot.tar.gz

echo "🚀 Installing to /usr/local/bin (may require sudo)..."
sudo mv laravelboot /usr/local/bin/
chmod +x /usr/local/bin/laravelboot

rm laravelboot.tar.gz
if [ -f "LICENSE" ]; then rm LICENSE; fi
if [ -f "README.md" ]; then rm README.md; fi

echo "✨ laravelboot installed successfully!"
laravelboot --help || true
