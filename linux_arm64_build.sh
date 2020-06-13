#!/bin/bash -e

##
## get the package TDM-GCC-64 from http://sourceforge.net/projects/tdm-gcc/files/TDM-GCC%20Installer/tdm64-gcc-4.9.2-3.exe/download

rm -f libsass_linux_arm64.a

git submodule init
git submodule update

export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=1
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig


pushd libsass
git clean -fxd
git reset --hard
BUILD=static make -j4
popd

cp libsass/lib/libsass.a libsass_linux_arm64.a