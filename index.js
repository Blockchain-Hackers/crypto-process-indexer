const { Web3 } = require("web3");

// Replace 'YOUR_GOERLI_INFURA_API_KEY' with your Infura API key or use your own Ethereum node
// const infuraApiKey = "927b0bef549145fba75661d347f23b8a";
const web3 = new Web3(
  `wss://sepolia.infura.io/ws/v3/927b0bef549145fba75661d347f23b8a`
);

// Replace 'YOUR_CONTRACT_ADDRESSES' with an array of contract addresses you want to listen to
const contractAddresses = ["0xA17ddf0a5309d50D7a69CA096A5473240A715DfA"].map(
  (x) => x.toLowerCase()
);

// Replace 'YOUR_CONTRACT_ABI' with the ABI (Application Binary Interface) of your contract
const contractAbi = [
  {
    inputs: [
      {
        internalType: "string",
        name: "_greeting",
        type: "string",
      },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        indexed: false,
        internalType: "string",
        name: "greeting",
        type: "string",
      },
    ],
    name: "GreetingSet",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        indexed: false,
        internalType: "string",
        name: "action",
        type: "string",
      },
    ],
    name: "Interaction",
    type: "event",
  },
  {
    inputs: [],
    name: "performInteraction",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "string",
        name: "_newGreeting",
        type: "string",
      },
    ],
    name: "setGreeting",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "getGreeting",
    outputs: [
      {
        internalType: "string",
        name: "",
        type: "string",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
];

// Create contract instances
const contracts = contractAddresses.map(
  (address) => new web3.eth.Contract(contractAbi, address)
);
let latestBlockNumber = 0;
// Subscribe to new block headers
async function subscribeToNewBlocks() {
  const _latestBlockNumber = await web3.eth.getBlockNumber();
  if (latestBlockNumber === _latestBlockNumber) {
    console.log("latestBlockNumber:", String(latestBlockNumber));
    return;
  }
  console.log("latestBlockNumber:", String(latestBlockNumber));
  getAllEventsInBlock(_latestBlockNumber);
  latestBlockNumber = _latestBlockNumber;
}

// Get all events in a given block for specified contract addresses
async function getAllEventsInBlock(blockNumber) {
  const block = await web3.eth.getBlock(blockNumber, true);
  //   console.log("Block:", block);

  if (block && block.transactions) {
    block.transactions.forEach((transaction) => {
      console.log("transaction.to", transaction.to);
      if (
        transaction.to &&
        contractAddresses.includes(transaction.to.toLowerCase())
        //   && transaction.logs
      ) {
        console.log("new transaction on contract", transaction.to);
        console.log("transaction", transaction);
        web3.eth.getTransactionReceipt(transaction.hash).then((receipt) => {
          let logs = receipt.logs;
          console.log("logs", logs);
          logs.forEach((log) => {
            const contract = contracts.find(
              (c) =>
                c.options.address.toLowerCase() === transaction.to.toLowerCase()
            );
            if (contract) {
              //   const decodedEvent = contract.events.allEvents
              if (receipt.logs && receipt.logs.length > 0) {
                console.log("Transaction Logs:", receipt.logs);

                // Iterate over ABI entries to find events
                contractAbi.forEach((abiEntry) => {
                  if (abiEntry.type === "event") {
                    // Decode logs using the event ABI
                    const decodedLogs = receipt.logs.map((log) => {
                      return web3.eth.abi.decodeLog(
                        abiEntry.inputs,
                        log.data,
                        log.topics.slice(1)
                      );
                    });

                    console.log(
                      `Decoded Logs for event ${abiEntry.name || "unknown"}:`,
                      decodedLogs
                    );
                  }
                });
              } else {
                console.log("No logs found for this transaction.");
              }
              console.log(
                `Decoded event in block ${blockNumber} for contract ${transaction.to}:`
              );
            } else {
              console.warn(
                `Contract ABI not found for address ${transaction.to}. Unable to decode event.`
              );
            }
          });
        });
      }
    });
  }
}

// Run the script
function run() {
  setInterval(subscribeToNewBlocks, 5000);
}

run();

// do a small express server on port 8080
const express = require("express");
const app = express();
const port = 8080;
app.get("/", (req, res) => {
  res.send("Hello World!");
});
app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`);
});