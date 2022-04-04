#!/usr/bin/env bash
# Installation Script For GOG on Linux/Darwin

set -e

GOG_VERSION=`curl -s https://api.github.com/repos/systemfiles/gog/releases/latest | grep "tag_name" | cut -d: -f2 | tr -d \" | cut -d, -f1 | xargs`

if [[ -f "$HOME/bin/gog" ]]; then
  echo "GOG is already installed."
  exit 0
fi

if [ ! command -v tar &> /dev/null ]; then
  echo "required command, tar, could not be found"
  exit 1
fi

if [ ! command -v curl &> /dev/null ]; then
  echo "required command, curl, could not be found"
  exit 1
fi

[ "$(uname -s)" == "Darwin" ] && INSTALL_OS="darwin"
[ "$(uname -s)" == "Linux" ] && INSTALL_OS="linux"

if [[ -z $INSTALL_OS ]]; then
  echo "Current OS not supported by installation script ... exiting!"
  exit 1
fi

[ "$(uname -m)" == "x86_64" ] && INSTALL_ARCH="amd64"
[ "$(uname -m)" == "armv7l" ] && INSTALL_ARCH="arm64"
[ "$(uname -m)" == "i386" ] && INSTALL_ARCH="386"

if [[ -z $INSTALL_ARCH ]]; then
  echo "Current OS Architecture not supported by installation script ... exiting!"
  exit 1
fi

INSTALL_FILE="GOG-${GOG_VERSION}-${INSTALL_OS}-${INSTALL_ARCH}.tar.gz"

if [[ ! -d "$HOME/bin" ]]; then
  mkdir "$HOME/bin"
fi

if [[ ! -d "$HOME/gogtmp" ]]; then
  mkdir "$HOME/gogtmp"
fi

cd $HOME/gogtmp
curl -LO "https://github.com/SystemFiles/GOG/releases/download/${GOG_VERSION}/${INSTALL_FILE}"
tar -zxvf "./${INSTALL_FILE}"
mv ./gog $HOME/bin/gog
cd; rm -rf $HOME/gogtmp/

$HOME/bin/gog -v
if [[ $? -ne 0 ]]; then
  echo "Installation Failed!"
  exit 1
fi