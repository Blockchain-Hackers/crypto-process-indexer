// require('dotenv').config()
const getProviderRpcUrl = (network) => {
  require("@chainlink/env-enc").config();
  let rpcUrl;

  switch (network) {
    case "ethereumMainnet":
      rpcUrl = "https://rpc.ankr.com/eth";
      break;
    case "ethereumSepolia":
      rpcUrl = "https://rpc.sepolia.org";
      break;
    case "optimismMainnet":
      rpcUrl = "https://mainnet.optimism.io";
      break;
    case "optimismGoerli":
      rpcUrl = "https://optimism-goerli.blockpi.network/v1/rpc/public";
      break;
    case "arbitrumTestnet":
      rpcUrl = "https://sepolia-rollup.arbitrum.io/rpc";
      break;
    case "avalancheMainnet":
      rpcUrl = "https://api.avax.network/ext/bc/C/rpc";
      break;
    case "avalancheFuji":
      rpcUrl = "https://api.avax-test.network/ext/bc/C/rpc";
      break;
    case "polygonMainnet":
      rpcUrl = "https://polygon-rpc.com";
      break;
    case "polygonMumbai":
      rpcUrl = "https://rpc-mumbai.maticvigil.com";
      break;
    case "bnbTestnet":
      rpcUrl = "https://data-seed-prebsc-2-s3.binance.org:8545/";
      break;
    case "baseGoerli":
      rpcUrl = "https://sepolia.base.org";
      break;
    default:
      throw new Error("Unknown network: " + network);
  }

  if (!rpcUrl)
    throw new Error(
      `rpcUrl empty for network ${network} - check your environment variables`
    );
  return rpcUrl;
};

const getPrivateKey = () => {
  require("@chainlink/env-enc").config();
  const privateKey = process.env.PRIVATE_KEY;
  if (!privateKey)
    throw new Error(
      "private key not provided - check your environment variables"
    );
  return privateKey;
};

module.exports = {
  getPrivateKey,
  getProviderRpcUrl,
};
