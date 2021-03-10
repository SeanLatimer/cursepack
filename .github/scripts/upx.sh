#! /usr/bin/env bash

UPX_VER=3.96

wget -O /tmp/upx.tar.xz "https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-amd64_linux.tar.xz"
tar -xvf /tmp/upx.tar.xz -C /tmp --wildcards --no-anchored upx-*/upx
mv /tmp/upx-*/upx /usr/bin
