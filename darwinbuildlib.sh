#!/bin/bash -e

##
## get the package TDM-GCC-64 from http://sourceforge.net/projects/tdm-gcc/files/TDM-GCC%20Installer/tdm64-gcc-4.9.2-3.exe/download

rm -f libsass_darwin.a

git submodule init
git submodule update

pushd libsass
git clean -fxd
git reset --hard
BUILD=static make -j4
popd

cp libsass/lib/libsass.a libsass_darwin.a
