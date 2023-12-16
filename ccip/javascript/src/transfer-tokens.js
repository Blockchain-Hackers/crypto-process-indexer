// Import necessary modules and data
const {
  getProviderRpcUrl,
  getRouterConfig,
  getPrivateKey,
  getMessageState,
} = require("./config");
const ethers = require("ethers");
const routerAbi = require("../../abi/Router.json");
const offRampAbi = require("../../abi/OffRamp.json");
const erc20Abi = require("../../abi/IERC20Metadata.json");

// Command: node src/transfer-tokens.js sourceChain destinationChain destinationAccount tokenAddress amount feeTokenAddress(optional)
// Examples(sepolia):

// pay fees with native token: node src/transfer-tokens.js ethereumSepolia avalancheFuji 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 100
// pay fees with transferToken: node src/transfer-tokens.js ethereumSepolia avalancheFuji 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 100 0x779877A7B0D9E8603169DdbD7836e478b4624789
// pay fees with a wrapped native token: node src/transfer-tokens.js ethereumSepolia avalancheFuji 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 100 0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534
const handleArguments = () => {
  if (process.argv.length !== 8 && process.argv.length !== 9) {
    // feeTokenAddress is optional
    throw new Error("Wrong number of arguments");
  }

  const sourceChain = process.argv[2];
  const destinationChain = process.argv[3];
  const destinationAccount = process.argv[4];
  const tokenAddress = process.argv[5];
  const amount = ethers.BigNumber.from(process.argv[6]);
  const privateKey = process.argv[7];
  const feeTokenAddress = process.argv[8];

  ///// example sepolia to polygonmumbai
  // with link
  // node src/transfer-tokens.js ethereumSepolia polygonMumbai 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 1000000000000000 0xdfe59eb3c16344e2d1daeda3611ba43116926ad4c753ab046b855fcfac883e4c 0x779877A7B0D9E8603169DdbD7836e478b4624789
  // native
  // node src/transfer-tokens.js ethereumSepolia polygonMumbai 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 1000000000000000 0xdfe59eb3c16344e2d1daeda3611ba43116926ad4c753ab046b855fcfac883e4c

  return {
    sourceChain,
    destinationChain,
    destinationAccount,
    tokenAddress,
    amount,
    feeTokenAddress,
    privateKey,
  };
};

const transferTokens = async (args) => {
  const {
    sourceChain,
    destinationChain,
    destinationAccount,
    tokenAddress,
    amount,
    feeTokenAddress,
    privateKey,
  } = args;

  // console.log("sourceChain", sourceChain);
  // console.log("args", args);

  /* 
  ==================================================
      Section: INITIALIZATION
      This section of the code parses the source and 
      destination router addresses and blockchain 
      selectors.
      It also initialized the ethers providers 
      to communicate with the blockchains.
  ==================================================
  */

  // Get the RPC URL for the chain from the config
  const rpcUrl = getProviderRpcUrl(sourceChain);
  // Initialize a provider using the obtained RPC URL
  const provider = new ethers.providers.JsonRpcProvider(rpcUrl);
  const wallet = new ethers.Wallet(privateKey);
  const signer = wallet.connect(provider);

  // Get the router's address for the specified chain
  const sourceRouterAddress = getRouterConfig(sourceChain).address;
  const sourceChainSelector = getRouterConfig(sourceChain).chainSelector;
  // Get the chain selector for the target chain
  const destinationChainSelector =
    getRouterConfig(destinationChain).chainSelector;

  // Create a contract instance for the router using its ABI and address
  const sourceRouter = new ethers.Contract(
    sourceRouterAddress,
    routerAbi,
    signer
  );

  /* 
  ==================================================
      Section: Check token validity
      Check first if the token you would like to 
      transfer is supported.
  ==================================================
  */

  // Fetch the list of supported tokens
  const supportedTokens = await sourceRouter.getSupportedTokens(
    destinationChainSelector
  );

  if (!supportedTokens.includes(tokenAddress)) {
    throw Error(
      `Token address ${tokenAddress} not in the list of supportedTokens ${supportedTokens}`
    );
  }

  /* 
  ==================================================
      Section: BUILD CCIP MESSAGE
      build CCIP message that you will send to the
      Router contract.
  ==================================================
  */

  // build message
  const tokenAmounts = [
    {
      token: tokenAddress,
      amount: amount,
    },
  ];

  // Encoding the data

  const functionSelector = ethers.utils.id("CCIP EVMExtraArgsV1").slice(0, 10);
  //  "extraArgs" is a structure that can be represented as [ 'uint256']
  // extraArgs are { gasLimit: 0 }
  // we set gasLimit specifically to 0 because we are not sending any data so we are not expecting a receiving contract to handle data

  const extraArgs = ethers.utils.defaultAbiCoder.encode(["uint256"], [0]);

  const encodedExtraArgs = functionSelector + extraArgs.slice(2);

  const message = {
    receiver: ethers.utils.defaultAbiCoder.encode(
      ["address"],
      [destinationAccount]
    ),
    data: "0x", // no data
    tokenAmounts: tokenAmounts,
    feeToken: feeTokenAddress ? feeTokenAddress : ethers.constants.AddressZero, // If fee token address is provided then fees must be paid in fee token.
    extraArgs: encodedExtraArgs,
  };

  /* 
  ==================================================
      Section: CALCULATE THE FEES
      Call the Router to estimate the fees for sending tokens.
  ==================================================
  */

  // const fees = await sourceRouter.getFee(destinationChainSelector, message);

  const fees = 8220672353649032; // 0.000220672353649032 ETH
  console.log(`Estimated fees (wei): ${fees}`);

  /* 
  ==================================================
      Section: SEND tokens
      This code block initializes an ERC20 token contract for token transfer across chains. It handles three cases:
      1. If the fee token is the native blockchain token, it makes one approval for the transfer amount. The fees are included in the msg.value field.
      2. If the fee token is different from both the native blockchain token and the transfer token, it makes two approvals: one for the transfer amount and another for the fees. The fees are part of the message.
      3. If the fee token is the same as the transfer token but not the native blockchain token, it makes a single approval for the sum of the transfer amount and fees. The fees are part of the message.
      The code waits for the transaction to be mined and stores the transaction receipt.
  ==================================================
  */

  // Create a contract instance for the token using its ABI and address
  const erc20 = new ethers.Contract(tokenAddress, erc20Abi, signer);
  let sendTx, approvalTx;

  if (!feeTokenAddress) {
    // Pay native
    // First approve the router to spend tokens
    approvalTx = await erc20.approve(sourceRouterAddress, amount);
    await approvalTx.wait(); // wait for the transaction to be mined
    console.log(
      `approved router ${sourceRouterAddress} to spend ${amount} of token ${tokenAddress}. Transaction: ${approvalTx.hash}`
    );

    sendTx = await sourceRouter.ccipSend(destinationChainSelector, message, {
      value: fees,
    }); // fees are send as value since we are paying the fees in native
  } else {
    if (tokenAddress.toUpperCase() === feeTokenAddress.toUpperCase()) {
      // fee token is the same as the token to transfer
      // Amount tokens to approve are transfer amount + fees
      approvalTx = await erc20.approve(sourceRouterAddress, amount + fees);
      await approvalTx.wait(); // wait for the transaction to be mined
      console.log(
        `approved router ${sourceRouterAddress} to spend ${amount} and fees ${fees} of token ${tokenAddress}. Transaction: ${approvalTx.hash}`
      );
    } else {
      // fee token is different than the token to transfer
      // 2 approvals
      approvalTx = await erc20.approve(sourceRouterAddress, amount); // 1 approval for the tokens to transfer
      await approvalTx.wait(); // wait for the transaction to be mined
      console.log(
        `approved router ${sourceRouterAddress} to spend ${amount} of token ${tokenAddress}. Transaction: ${approvalTx.hash}`
      );
      const erc20Fees = new ethers.Contract(feeTokenAddress, erc20Abi, signer);
      approvalTx = await erc20Fees.approve(sourceRouterAddress, fees); // 1 approval for the fees token
      await approvalTx.wait();
      console.log(
        `approved router ${sourceRouterAddress} to spend  fees ${fees} of token ${feeTokenAddress}. Transaction: ${approvalTx.hash}`
      );
    }
    sendTx = await sourceRouter.ccipSend(destinationChainSelector, message);
  }

  const receipt = await sendTx.wait(); // wait for the transaction to be mined

  /* 
  ==================================================
      Section: Fetch message ID
      The Router ccipSend function returns the messageId.
      This section makes a call (simulation) to the blockchain
      to fetch the messageId that was returned by the Router.
  ==================================================
  */

  // Simulate a call to the router to fetch the messageID
  const call = {
    from: sendTx.from,
    to: sendTx.to,
    data: sendTx.data,
    gasLimit: sendTx.gasLimit,
    gasPrice: sendTx.gasPrice,
    value: sendTx.value,
  };

  // Simulate a contract call with the transaction data at the block before the transaction
  const messageId = await provider.call(call, receipt.blockNumber - 1);

  console.log(
    `\n✅ ${amount} of Tokens(${tokenAddress}) Sent to account ${destinationAccount} on destination chain ${destinationChain} using CCIP. Transaction hash ${sendTx.hash} -  Message id is ${messageId}`
  );

  return {
    amountInWei: amount,
    destination: destinationAccount,
    destinationChain: destinationChain,
    sourceChain: sourceChain,
    messageId,
    hash: sendTx.hash,
    chainlinkExplorerUrl: `https://ccip.chain.link/msg/${messageId}`,
    fee: sendTx.fee,
  };
};

module.exports = {
  transferTokens,
};
