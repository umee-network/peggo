#!/bin/bash -eux

DATA_DIR=${DATA_DIR:-$PWD/data}

rm -rf $DATA_DIR
mkdir -p $DATA_DIR

ganache="ganache"

if pgrep -x ganache >/dev/null
then
  echo "ganache is running, going to kill all"
  ps -ef | grep ganache | grep -v grep | awk '{print $2}' | xargs kill
fi

ganache \
  --chain.chainId 888 \
  --account '0x06e48d48a55cc6843acb2c3c23431480ec42fca02683f4d8d3d471372e5317ee,1000000000000000000' \
  --account '0x4faf826f3d3a5fa60103392446a72dea01145c6158c6dd29f6faab9ec9917a1b,1000000000000000000' \
  --account '0x11f746395f0dd459eff05d1bc557b81c3f7ebb1338a8cc9d36966d0bb2dcea21,1000000000000000000' \
  --blockTime 1  \
  > $DATA_DIR/ganache.log 2>&1 &