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
mv "dist-$OS-$ARCH/vendor-$PLUGIN_VERSION.tgz" "dist/$OS-$ARCH-$VERSION.tgz"
mv "dist-$OS-$ARCH/vendor-$PLUGIN_VERSION.tgz.prov" "dist/$OS-$ARCH-$VERSION.tgz.prov"

# Update the filename reference in the provenance file
# Use different sed syntax for macOS vs Linux
if [ "$(uname)" = "Darwin" ]; then
    sed -i '' "s/vendor-$PLUGIN_VERSION\.tgz/$OS-$ARCH-$VERSION.tgz/g" "dist/$OS-$ARCH-$VERSION.tgz.prov"
else
    sed -i "s/vendor-$PLUGIN_VERSION\.tgz/$OS-$ARCH-$VERSION.tgz/g" "dist/$OS-$ARCH-$VERSION.tgz.prov"
fi

rm -rf "tmp-$OS-$ARCH" "dist-$OS-$ARCH"

