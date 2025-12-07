#!/usr/bin/env sh
set -e

echo "Removing goreleaser metadata from dist..."
rm -f dist/config.yaml
rm -f dist/metadata.json
rm -f dist/artifacts.json

echo "Moving plugin specific configurations to dist..."
cp -f plugin.yaml "dist/"
cp -f vendor.complete "dist/"

echo "Packaging helm chart with 'helm plugin package'..."
echo "$HELM_KEY_PASSPHRASE" | helm plugin package \
    --sign \
    --key "$SIGNING_KEY_EMAIL" \
    --keyring secring.gpg \
    --passphrase-file - \
    -d "." \
    "dist/"

echo "Moving generated OCI artifacts into dist..."
mv -f *.tgz dist/
mv -f *.tgz.prov dist/
ls "dist"
