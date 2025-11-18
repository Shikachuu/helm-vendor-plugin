#!/usr/bin/env sh
echo "Installing helm vendor plugin"

PROJECT_NAME="helm-vendor-plugin"
PROJECT_ORG="${PROJECT_ORG:-Shikachuu}"
PROJECT_GH="$PROJECT_ORG/$PROJECT_NAME"
BINARY_NAME="vendor-plugin"
export GREP_COLOR="never"

HELM_MAJOR_VERSION=$("${HELM_BIN}" version --client --short | awk -F '.' '{print $1}')

: ${HELM_PLUGIN_DIR:="$("${HELM_BIN}" home --debug=false)/plugins/$PROJECT_NAME"}

# Convert the HELM_PLUGIN_DIR to unix if cygpath is available
if type cygpath >/dev/null 2>&1; then
  HELM_PLUGIN_DIR=$(cygpath -u $HELM_PLUGIN_DIR)
fi

if [ "$SKIP_BIN_INSTALL" = "1" ]; then
  echo "Skipping binary install"
  exit
fi

# initArch discovers the architecture for this system
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    aarch64) ARCH="arm64" ;;
    x86_64) ARCH="amd64" ;;
    *)
      echo "Unsupported architecture: $ARCH"
      echo "Supported architectures: amd64, arm64"
      exit 1
      ;;
  esac
}

# initOS discovers the operating system for this system
initOS() {
  OS=$(uname | tr '[:upper:]' '[:lower:]')

  case "$OS" in
    msys*) OS='windows' ;;
    mingw*) OS='windows' ;;
    darwin) OS='darwin' ;;
  esac
}

# verifySupported checks that the os/arch combination is supported and tools are available
verifySupported() {
  # Check for windows-arm64 which is explicitly not supported
  if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
    echo "No prebuild binary for ${OS}-${ARCH}."
    exit 1
  fi

  if ! type "curl" >/dev/null && ! type "wget" >/dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version
getDownloadURL() {
  # Add .exe extension for Windows binaries
  BINARY_EXT=""
  if [ "$OS" = "windows" ]; then
    BINARY_EXT=".exe"
  fi

  version="$(cat $HELM_PLUGIN_DIR/plugin.yaml | grep "version" | cut -d '"' -f 2)"
  if [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/v$version/bin/$BINARY_NAME-$OS-$ARCH$BINARY_EXT"
  else
    url="https://api.github.com/repos/$PROJECT_GH/releases/latest"
    if type "curl" >/dev/null; then
      DOWNLOAD_URL=$(curl -s $url | grep "bin/$BINARY_NAME-$OS-$ARCH$BINARY_EXT\"" | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
    elif type "wget" >/dev/null; then
      DOWNLOAD_URL=$(wget -q -O - $url | grep "bin/$BINARY_NAME-$OS-$ARCH$BINARY_EXT\"" | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
    fi
  fi
}

# downloadFile downloads the latest binary
downloadFile() {
  BINDIR="$HELM_PLUGIN_DIR/bin"
  rm -rf "$BINDIR"
  mkdir -p "$BINDIR"

  # Local binary name (without .exe extension for consistency)
  LOCAL_BINARY="$BINDIR/$BINARY_NAME"

  echo "Downloading $DOWNLOAD_URL"
  if type "curl" >/dev/null; then
    HTTP_CODE=$(curl -sL --write-out "%{http_code}" "$DOWNLOAD_URL" --output "$LOCAL_BINARY")
    if [ ${HTTP_CODE} -ne 200 ]; then
      echo "Failed to download binary (HTTP ${HTTP_CODE})"
      exit 1
    fi
  elif type "wget" >/dev/null; then
    wget -q -O "$LOCAL_BINARY" "$DOWNLOAD_URL"
  fi

  chmod +x "$LOCAL_BINARY"
}

# fail_trap is executed if an error occurs
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    printf "\tFor support, go to https://github.com/$PROJECT_GH.\n"
  fi
  exit $result
}

# testVersion tests the installed client
testVersion() {
  set +e
  echo "$BINARY_NAME installed into $HELM_PLUGIN_DIR/bin/$BINARY_NAME"
  "${HELM_PLUGIN_DIR}/bin/$BINARY_NAME" version
  set -e
}

# Execution
trap "fail_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
downloadFile
testVersion
