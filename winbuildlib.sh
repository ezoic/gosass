#!/bin/bash -e

##
## get the package TDM-GCC-64 from http://sourceforge.net/projects/tdm-gcc/files/TDM-GCC%20Installer/tdm64-gcc-4.9.2-3.exe/download

GITPATH=/C/Program\ Files\ \(x86\)/Git/bin
if [ "$1" == "" ]; then
    GITPATH=$1
fi

export PATH=/usr/local/bin:/mingw/bin:/bin
export PATH=$PATH:/c/TDM-GCC-64/bin:$GITPATH
export CC=gcc
export LDFLAGS=-shared-libstdc++

rm -f libsass_windows.a

git submodule init
git submodule update



pushd libsass
git clean -fxd
mingw32-make.exe clean
mingw32-make.exe lib/libsass.a
popd

cp libsass/lib/libsass.a libsass_windows.a