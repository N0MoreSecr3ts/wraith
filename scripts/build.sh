#!/bin/bash

# TODO: Fold this into the makefile and delete it
BUILD_FOLDER=build
VERSION=$(cat common/banner.go | grep Version | cut -d '"' -f 2)

bin_dep() {
  BIN=$1
  which $BIN >/dev/null || {
    echo "[-] Dependency $BIN not found !"
    exit 1
  }
}

create_exe_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" wraith.exe ../README.md ../LICENSE.txt ../contentsignatures.json ../filesignatures.json >/dev/null
  rm -rf wraith wraith.exe
}

create_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" wraith ../README.md ../LICENSE.md ../contentsignatures.json ../filesignatures.json >/dev/null
  rm -rf wraith wraith.exe
}

build_linux_amd64() {
  echo "[*] Building linux/amd64 ..."
  GOOS=linux GOARCH=amd64 go build -o wraith ..
}

build_macos_amd64() {
  echo "[*] Building darwin/amd64 ..."
  GOOS=darwin GOARCH=amd64 go build -o wraith ..
}

build_windows_amd64() {
  echo "[*] Building windows/amd64 ..."
  GOOS=windows GOARCH=amd64 go build -o wraith.exe ..
}

rm -rf $BUILD_FOLDER
mkdir $BUILD_FOLDER
cd $BUILD_FOLDER

build_linux_amd64 && create_archive wraith_linux_amd64_$VERSION.zip
build_macos_amd64 && create_archive wraith_macos_amd64_$VERSION.zip
#windows builds are broken with the addition of go-gitlab
#build_windows_amd64 && create_exe_archive wraith_windows_amd64_$VERSION.zip
shasum -a 256 * >checksums.txt

echo
echo
du -sh *

cd --

