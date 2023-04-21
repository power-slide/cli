#!/usr/bin/env bash

######## Check for install tools

if ! [ -x "$(command -v sed)" ] ; then
    echo 'Error: Please install sed.' >&2
    exit 1
fi

if ! [ -x "$(command -v curl)" ] ; then
    echo 'Error: Please install curl.' >&2
    exit 1
fi

if ! [ -x "$(command -v jq)" ] ; then
    echo 'Error: Please install jq.' >&2
    exit 1
fi

######## Run Installer

BIN_NAME="pwrsl"
if [ -x "$(command -v ${BIN_NAME})" ] ; then
    echo 'Error: Already installed.' >&2
    exit 1
fi
echo "Installing PowerSlide CLI..."

#### Determine OS, CPU arch, and get download link

ARCH=$(uname -m)
OS=$(uname)
echo "Detected ${OS} on ${ARCH}"

OS=$(echo ${OS} | sed 's/[A-Z]/\L&/g') # downcase
case $ARCH in
    x86_64) ARCH="amd64" ;;
    *) ;;
esac
TARGET="${BIN_NAME}-${OS}-${ARCH}"

LATEST_RELEASE=$(curl -fsSL \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/power-slide/cli/releases/latest)

LATEST_URL=$(echo ${LATEST_RELEASE} | jq -r ".assets[]? | select(.name==\"${TARGET}\") | .browser_download_url")
LATEST_VERSION=$(echo ${LATEST_RELEASE} | jq -r ".tag_name?")

if [ $LATEST_VERSION == "null" ] || [ $LATEST_URL == "null" ] ; then
    echo "Unable to get latest release from GitHub API" >&2
    exit 1
fi
echo "Installing version: ${LATEST_VERSION}"

#### Get install path

BIN_DIR="${HOME}/bin"
while ! [ -d $BIN_DIR ] || ! [ -w $BIN_DIR ] || ! [[ ":$PATH:" == *":$BIN_DIR:"* ]] ; do
    echo "${BIN_DIR} does not exist, is not writeable, or is not on your \$PATH"
    echo "If you update your \$PATH, make sure to run this in a newly started shell."
    read -p "Please enter a writeable directory on your \$PATH to download to: " BIN_DIR
done

BIN="${BIN_DIR}/${BIN_NAME}"
echo -n "Downloading to ${BIN}... "
curl -fsSL ${LATEST_URL} -o ${BIN}
chmod u+x ${BIN}
echo "Done!"

######## Run PowerSlide setup (hand-off to CLI)

${BIN} setup
