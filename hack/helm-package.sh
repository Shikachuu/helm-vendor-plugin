#!/usr/bin/env sh
set -e

PLUGIN_VERSION=$(yq eval '.version' plugin.yaml)

mkdir -p "tmp-$OS-$ARCH/bin"
mkdir -p "tmp-$OS-$ARCH/hack"
cp "$BIN_PATH" "tmp-$OS-$ARCH/bin/$BIN_NAME"
cp plugin.yaml "tmp-$OS-$ARCH/"
cp vendor.complete "tmp-$OS-$ARCH/"
cp hack/install-binary.sh "tmp-$OS-$ARCH/hack/install-binary.sh"

mkdir -p "dist-$OS-$ARCH"

echo "$HELM_KEY_PASSPHRASE" | helm plugin package \
    --sign \
    --key "$SIGNING_KEY_EMAIL" \
    --keyring secring.gpg \
    --passphrase-file - \
    -d "dist-$OS-$ARCH" \
    "tmp-$OS-$ARCH/"

mkdir -p dist
mv "dist-$OS-$ARCH/vendor-$PLUGIN_VERSION.tgz" "dist/vendor-$VERSION-$OS-$ARCH.tgz"
mv "dist-$OS-$ARCH/vendor-$PLUGIN_VERSION.tgz.prov" "dist/vendor-$VERSION-$OS-$ARCH.tgz.prov"

rm -rf "tmp-$OS-$ARCH" "dist-$OS-$ARCH"

