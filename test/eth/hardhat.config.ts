import { HardhatUserConfig } from "hardhat/config";
import { HardhatNetworkAccountUserConfig, HardhatNetworkUserConfig, NetworkUserConfig } from "hardhat/src/types/config"
import "@nomiclabs/hardhat-etherscan";
import "@nomiclabs/hardhat-ethers";
import { env } from "process";
import "./tasks/getCurrentValset";
import * as dotenv from "dotenv";

dotenv.config();
// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

function GetChainId(): number {
  if (env.CHAIN_ID != undefined) {
    return Number(env.CHAIN_ID);
  }
  return 888;
};

const balance = "100000000000000000000000000"

const privateKeys = [
  // val0 0xfac5EC50BdfbB803f5cFc9BF0A0C2f52aDE5b6dd
  "0x06e48d48a55cc6843acb2c3c23431480ec42fca02683f4d8d3d471372e5317ee",
  // val1 0x02fa1b44e2EF8436e6f35D5F56607769c658c225
  "0x4faf826f3d3a5fa60103392446a72dea01145c6158c6dd29f6faab9ec9917a1b",
  // val2 0xd8f468c1B719cc2d50eB1E3A55cFcb60e23758CD
  "0x11f746395f0dd459eff05d1bc557b81c3f7ebb1338a8cc9d36966d0bb2dcea21",
]


const peggoAccounts: HardhatNetworkAccountUserConfig[] = [
  {
    privateKey: privateKeys[0],
    balance: balance,
  },
  {
    privateKey: privateKeys[1],
    balance: balance,
  },
  {
    privateKey: privateKeys[2],
    balance: balance,
  },
]


const peggoTestNetwork: HardhatNetworkUserConfig = {
  chainId: GetChainId(),
  accounts: peggoAccounts,
  mining: {
    // It will produce empty blocks every 3 seconds 
    interval: 3000
  }
}

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.10",
    settings: {
      optimizer: {
        enabled: true
      }
    }
  },
  networks: {
    hardhat: peggoTestNetwork,
    ganache: {
      chainId: GetChainId(),
      url: "http://127.0.0.1:8545",
      accounts: privateKeys
    },
    goerli: {
      url: env.ETHRPC,
      accounts: privateKeys
    }
  },
  etherscan: {
    apiKey: env.ETHERSCAN_API,
  },
  paths: {
    sources: "./contracts",
    artifacts: "./artifacts"
  }
};

export default config;