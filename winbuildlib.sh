#!/bin/bash

##
## get the package TDM-GCC-64 from http://sourceforge.net/projects/tdm-gcc/files/TDM-GCC%20Installer/tdm64-gcc-4.9.2-3.exe/download

export PATH=/c/TDM-GCC-64/bin:$PATH
export CC=gcc
export LDFLAGS=-shared-libstdc++

rm -f libsass.dll

git submodule init
git submodule update



pushd libsass
git reset --hard
mingw32-make.exe clean
mingw32-make.exe lib/libsass.dll
#mingw32-make.exe lib/libsass.a
popd

cp libsass/lib/libsass.dll .