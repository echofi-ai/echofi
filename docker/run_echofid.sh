#!/bin/sh

if test -n "$1"; then
    # need -R not -r to copy hidden files
    cp -R "$1/.echofi" /root
fi

mkdir -p /root/log
echofid start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices 0.0001uecho --trace
