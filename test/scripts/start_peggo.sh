#!/bin/bash -eu

# this script aims to augment `multinode.sh`
#  - by spinning up ganache/ethrpc with 3 connected peggo orchestrators
#  - automagically deploying the umee bridge contract and erc20 tokens
#  - continuously tails the first orchestrator logs
#  - ctrl+c to cleanly exit
#
# to install ganache: `npm i -g ganache-cli`

CWD="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
CHAIN_ID="${CHAIN_ID:-888}"
CHAIN_DIR="${CHAIN_DIR:-$CWD/test/cosmos/data}"

hdir="$CHAIN_DIR/$CHAIN_ID"
# Folders for nodes
n0dir="$hdir/n0"
n1dir="$hdir/n1"
n2dir="$hdir/n2"

# private keys for first 4 eth accounts, from ganache
PEGGO_ETH_PK=0x88cbead91aee890d27bf06e003ade3d4e952427e88f88d31d61d3ef5e5d54305
PEGGO_ETH_PK2=0x741de4f8988ea941d3ff0287911ca4074e62b7d45c991a51186455366f10b544
PEGGO_ETH_PK3=0x39a4c898dda351d54875d5ebb3e1c451189116faa556c3c04adc860dd1000608
PEGGO_ETH_PK4=0x6c212553111b370a8ffdc682954495b7b90a73cedab7106323646a4f2c4e668f

ETHRPC=${ETHRPC:-HTTP://127.0.0.1:8545}
echo "ETH RPC is ${ETHRPC}"

# boot ganache
echo "Starting ganache with 4 accounts..."
ganache-cli \
  --db test \
  --chain-id $CHAIN_ID \
  --account "${PEGGO_ETH_PK},1000000000000000000000000" \
  --account "${PEGGO_ETH_PK2},1000000000000000000000000" \
  --account "${PEGGO_ETH_PK3},1000000000000000000000000" \
  --account "${PEGGO_ETH_PK4},1000000000000000000000000" \
  --blockTime 1 2>&1 > ganache.log &
GANACHE_PID=$$

# deploy gravity contract
echo "Deploying gravity contract..."
PEGGO_ETH_PK=$PEGGO_ETH_PK peggo bridge deploy-gravity --eth-rpc $ETHRPC

# yield to ganache output
sleep 1

# capture bridge
BRIDGEADDR=$(grep -a "Contract created:" ganache.log | sed -e 's/.*: //')
grep -a "Contract created:" ganache.log | sed -e 's/.*: //'
echo "Bridge deployed to $BRIDGEADDR"

# deploy umee erc20 token
echo "Deploying Umee ERC20 token..."
PEGGO_ETH_PK=$PEGGO_ETH_PK peggo bridge deploy-erc20 $BRIDGEADDR uumee --eth-rpc $ETHRPC

# boot 3 peggo orchestrators, one for each validator, eg. multinode.sh
echo "Starting 3 peggo orchestrators"
PEGGO_ETH_PK=$PEGGO_ETH_PK peggo orchestrator $BRIDGEADDR \
  --status-api=true \
  --eth-rpc=$ETHRPC \
  --relay-batches=true \
  --valset-relay-mode="all" \
  --cosmos-chain-id=$CHAIN_ID \
  --cosmos-gas-prices=1000uumee \
  --cosmos-grpc="tcp://0.0.0.0:9091" \
  --tendermint-rpc="http://0.0.0.0:26667" \
  --oracle-providers="mock" \
  --cosmos-keyring=test \
  --cosmos-keyring-dir=$n0dir \
  --cosmos-from=val  --log-level debug --log-format text --profit-multiplier=0 > peggo.1.log 2>&1 &
PEGGO_PID=$$
echo "PID is ${PEGGO_PID}"

PEGGO_ETH_PK=$PEGGO_ETH_PK2 peggo orchestrator $BRIDGEADDR \
  --eth-rpc=$ETHRPC \
  --relay-batches=true \
  --valset-relay-mode="all" \
  --cosmos-chain-id=$CHAIN_ID \
  --cosmos-gas-prices=1000uumee \
  --cosmos-grpc="tcp://0.0.0.0:9091" \
  --tendermint-rpc="http://0.0.0.0:26667" \
  --oracle-providers="mock" \
  --cosmos-keyring=test \
  --cosmos-keyring-dir=$n1dir \
  --cosmos-from=val  --log-level debug --log-format text --profit-multiplier=0 > peggo.2.log 2>&1 &
PEGGO_PID2=$$
echo "PID2 is ${PEGGO_PID2}"

PEGGO_ETH_PK=$PEGGO_ETH_PK3 peggo orchestrator $BRIDGEADDR \
  --eth-rpc=$ETHRPC \
  --relay-batches=true \
  --valset-relay-mode="all" \
  --cosmos-chain-id=$CHAIN_ID \
  --cosmos-gas-prices=1000uumee \
  --cosmos-grpc="tcp://0.0.0.0:9092" \
  --tendermint-rpc="http://0.0.0.0:26677" \
  --oracle-providers="mock" \
  --cosmos-keyring=test \
  --cosmos-keyring-dir=$n2dir \
  --cosmos-from=val  --log-level debug --log-format text --profit-multiplier=0 > peggo.3.log 2>&1 &
PEGGO_PID3=$$
echo "PID3 is ${PEGGO_PID3}"


# ---------
cleanup() {
  echo "Cleaning up ganache and peggo..."
  kill -9 ${GANACHE_PID} ${PEGGO_PID} ${PEGGO_PID2} ${PEGGO_PID3}
}

trap cleanup 1 2 3 6

# show first peggo orchestrator output
tail -f peggo.1.log
