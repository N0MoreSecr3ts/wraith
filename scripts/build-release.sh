#!/bin/bash

# TODO: Fold this into the makefile and delete it

BUILD_FOLDER=bin
RELEASE_FOLDER=bin/release
RULES_FOLDER=rules

bin_dep() {
  BIN=$1
  which $BIN > /dev/null || { echo "[-] Dependency $BIN not found !"; exit 1; }
}

create_exe_archive() {
  bin_dep 'zip'

  OUTPUT=$RELEASE_FOLDER/$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" ./bin/wraith.exe ./rules/test_rules.yml ./README.md > /dev/null
  rm -rf ./bin/wraith ./bin/wraith.exe
}

create_archive() {
  bin_dep 'zip'

  OUTPUT=$RELEASE_FOLDER/$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" ./bin/wraith ./rules/test_rules.yml ./README.md > /dev/null
  rm -rf ./bin/wraith ./bin/wraith.exe
}

build_linux_amd64() {
  echo "[*] Building linux/amd64 ..."
  pwd
  make release target_os=linux target_arch=amd64
}

build_macos_amd64() {
  echo "[*] Building darwin/amd64 ..."
  make release target_os=darwin target_arch=amd64
}

build_windows_amd64() {
  echo "[*] Building windows/amd64 ..."
  make release target_os=windows target_arch=amd64
}

rm -rf $BUILD_FOLDER
rm -rf $RULES_FOLDER
rm -rf $RELEASE_FOLDER
mkdir $BUILD_FOLDER
mkdir $RULES_FOLDER
mkdir -p $RELEASE_FOLDER

#cd $BUILD_FOLDER

if [[ "$OSTYPE" == "linux-gnu" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]]; then
    build_linux_amd64
    #
    chmod +x ./bin/wraith
    ./bin/wraith updateRules
    cp $HOME/.wraith/rules/test_rules.yml $RULES_FOLDER/test_rules.yml
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # Mac OSX
    build_macos_amd64
    chmod +x ./bin/wraith
    ./bin/wraith updateRules
    cp $HOME/.wraith/rules/test_rules.yml $RULES_FOLDER/test_rules.yml
fi

build_linux_amd64 && create_archive wraith_linux_amd64.zip
build_macos_amd64 && create_archive wraith_macos_amd64.zip
build_windows_amd64 && create_exe_archive wraith_windows_amd64.zip
shasum -a 256 $RELEASE_FOLDER/* > $RELEASE_FOLDER/checksums.txt

echo
echo
du -sh $RELEASE_FOLDER/*

rm -rf $RULES_FOLDER
cd --
