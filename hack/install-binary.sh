#!/usr/bin/env sh
set -e

echo "Installing helm vendor plugin"

# Configuration
PROJECT_NAME="helm-vendor-plugin"
PROJECT_ORG="${PROJECT_ORG:-Shikachuu}"
BINARY_NAME="vendor-plugin"

# Load Helm environment variables
eval $(${HELM_BIN} env)

# Get plugin directory
if [ -z "$HELM_PLUGIN_DIR" ]; then
  HELM_PLUGIN_DIR="$HELM_PLUGINS/$PROJECT_NAME"
fi

# Skip if requested
if [ "$SKIP_BIN_INSTALL" = "1" ]; then
  echo "Skipping binary install"
  exit 0
fi

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  arm64|aarch64) ARCH="arm64" ;;
  x86_64) ARCH="amd64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')

# Check for curl or wget
if ! type curl >/dev/null 2>&1 && ! type wget >/dev/null 2>&1; then
  echo "Either curl or wget is required"
  exit 1
fi

# Get version from plugin.yaml or use latest
VERSION=$(grep "^version:" "$HELM_PLUGIN_DIR/plugin.yaml" | awk '{print $2}')
if [ -n "$VERSION" ]; then
  DOWNLOAD_URL="https://github.com/$PROJECT_ORG/$PROJECT_NAME/releases/download/v$VERSION/$BINARY_NAME-$OS-$ARCH"
else
  API_URL="https://api.github.com/repos/$PROJECT_ORG/$PROJECT_NAME/releases/latest"
  if type curl >/dev/null 2>&1; then
    DOWNLOAD_URL=$(curl -s "$API_URL" | grep "browser_download_url.*$BINARY_NAME-$OS-$ARCH\"" | cut -d '"' -f 4)
  else
    DOWNLOAD_URL=$(wget -q -O - "$API_URL" | grep "browser_download_url.*$BINARY_NAME-$OS-$ARCH\"" | cut -d '"' -f 4)
  fi
fi

# Download binary
BINDIR="$HELM_PLUGIN_DIR/bin"
mkdir -p "$BINDIR"
LOCAL_BINARY="$BINDIR/$BINARY_NAME"

echo "Downloading $DOWNLOAD_URL"
if type curl >/dev/null 2>&1; then
  HTTP_CODE=$(curl -sL --write-out "%{http_code}" "$DOWNLOAD_URL" --output "$LOCAL_BINARY")
  if [ "$HTTP_CODE" != "200" ]; then
    echo "Failed to download binary (HTTP $HTTP_CODE)"
    exit 1
  fi
else
  if ! wget -q -O "$LOCAL_BINARY" "$DOWNLOAD_URL"; then
    echo "Failed to download binary"
    exit 1
  fi
fi

chmod +x "$LOCAL_BINARY"

# Test installation
echo "$BINARY_NAME installed into $LOCAL_BINARY"
"$LOCAL_BINARY" version || true
