#!/bin/bash -eux

CWD="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
DATA_DIR=${DATA_DIR:-$CWD/data}
ganache="ganache"

if ! command -v $ganache &> /dev/null
then
  echo "⚠️ $ganache command could not be found!"
  echo "Install it by running 'npm install ganache --global'"
  exit 1
fi

. $CWD/../scripts/pid_control.sh

pid_path=$DATA_DIR/pid
log_path=$DATA_DIR/ganache.log

kill_process $pid_path

rm -rf $DATA_DIR
mkdir -p $DATA_DIR

# Available Accounts
# ==================
# (0) 0xC6Fe5D33615a1C52c08018c47E8Bc53646A0E101 (1 ETH)
# (1) 0x963EBDf2e1f8DB8707D05FC75bfeFFBa1B5BaC17 (1 ETH)
# (2) 0x6880D7bfE96D49501141375ED835C24cf70E2bD7 (1 ETH)
# (3) 0x727AEE334987c52fA7b567b2662BDbb68614e48C (1 ETH)
# (4) 0x9C0E888cC804AA046845DB803482BC99d7F868C6 (1 ETH)
# (5) 0xa128dCF2c938D522aD391B38d874B0303c5D8781 (1 ETH)
# (6) 0xBE95486B9AAAbc25D46BEC5b63bA9cF35d12f618 (1 ETH)

# Private Keys
# ==================
# (0) 0x88cbead91aee890d27bf06e003ade3d4e952427e88f88d31d61d3ef5e5d54305
# (1) 0x741de4f8988ea941d3ff0287911ca4074e62b7d45c991a51186455366f10b544
# (2) 0x39a4c898dda351d54875d5ebb3e1c451189116faa556c3c04adc860dd1000608
# (3) 0x6c212553111b370a8ffdc682954495b7b90a73cedab7106323646a4f2c4e668f
# (4) 0x7f13aa5c42ac85655cef298a7808cdeda75091021d35273050b3c7407420163e
# (5) 0x42b8f5cbff3c9cacf4d0b8df75d7a23c16c4027f77dba33abc622d7faa4ed3d5
# (6) 0x41f3d61c554a660d912d0fe0ee0a282e2b2f28062985a151f41d3d437543474a

$ganache \
  --chain.chainId 888 \
  --miner.blockTime 2 \
  --logging.debug \
  --logging.verbose \
  --wallet.accounts '0x88cbead91aee890d27bf06e003ade3d4e952427e88f88d31d61d3ef5e5d54305,1000000000000000000' \
  --wallet.accounts '0x741de4f8988ea941d3ff0287911ca4074e62b7d45c991a51186455366f10b544,1000000000000000000' \
  --wallet.accounts '0x39a4c898dda351d54875d5ebb3e1c451189116faa556c3c04adc860dd1000608,1000000000000000000' \
  --wallet.accounts '0x6c212553111b370a8ffdc682954495b7b90a73cedab7106323646a4f2c4e668f,1000000000000000000' \
  --wallet.accounts '0x7f13aa5c42ac85655cef298a7808cdeda75091021d35273050b3c7407420163e,1000000000000000000' \
  --wallet.accounts '0x42b8f5cbff3c9cacf4d0b8df75d7a23c16c4027f77dba33abc622d7faa4ed3d5,1000000000000000000' \
  --wallet.accounts '0x41f3d61c554a660d912d0fe0ee0a282e2b2f28062985a151f41d3d437543474a,1000000000000000000' \
  > $log_path 2>&1 &

echo "$!" > $pid_path

echo
echo "Ganache Logs:"
echo "  * tail -f $log_path"
echo