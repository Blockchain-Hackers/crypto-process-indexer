const { Web3 } = require("web3");

// Replace 'YOUR_GOERLI_INFURA_API_KEY' with your Infura API key or use your own Ethereum node
// const infuraApiKey = "927b0bef549145fba75661d347f23b8a";
const web3 = new Web3(
  `wss://sepolia.infura.io/ws/v3/927b0bef549145fba75661d347f23b8a`
);

// Replace 'YOUR_CONTRACT_ADDRESSES' with an array of contract addresses you want to listen to
const contractAddresses = ["0xc3dB75b8081F7f6789B662f76fC3c14f7EcD41dF"].map(x=>x.toLowerCase());

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
  //   web3.eth.subscribe("newBlockHeaders", (error, blockHeader) => {
  //     if (error) {
  //       console.error("Error:", error);
  //     } else {
  //       console.log("New block header:", blockHeader);
  getAllEventsInBlock(_latestBlockNumber);
  latestBlockNumber = _latestBlockNumber;
  // }
  //   });
  // .on("error", (error) => {
  //   console.error("Subscription error:", error);
  // });
}

// Get all events in a given block for specified contract addresses
async function getAllEventsInBlock(blockNumber) {
  const block = await web3.eth.getBlock(blockNumber, true);
//   console.log("Block:", block);

  if (block && block.transactions) {
      block.transactions.forEach((transaction) => {
        console.log("transaction.to", transaction.to)
      if (
        transaction.to &&
          contractAddresses.includes(transaction.to.toLowerCase())
        //   && transaction.logs
      ) {
          console.log("new transaction on contract", transaction.to)
          console.log("transaction", transaction)
        transaction.logs.forEach((log) => {
          const contract = contracts.find(
            (c) =>
              c.options.address.toLowerCase() === transaction.to.toLowerCase()
          );
          if (contract) {
            const decodedEvent = contract.events.allEvents.decode(log);
            console.log(
              `Decoded event in block ${blockNumber} for contract ${transaction.to}:`,
              decodedEvent
            );
          } else {
            console.warn(
              `Contract ABI not found for address ${transaction.to}. Unable to decode event.`
            );
          }
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

// const Web3 = require('web3');

// // Set up the HTTP RPC endpoint
// const rpcEndpoint = 'http://localhost:8545'; // Replace with your RPC endpoint URL

// Create a new web3 instance
// const web3 = new Web3(rpcEndpoint);

// // Define the last scanned block variable
// let lastScannedBlock = 0;

// // Define the polling function
// async function pollRPC() {
//   try {
//     // Get the latest block number
//     const latestBlockNumber = await web3.eth.getBlockNumber();

//     // Define the event contract and ABI
//     const contractAddress = '0x1234567890'; // Replace with your contract address
//     const contractABI = [{
//       "anonymous": false,
//       "inputs": [
//         {
//           "indexed": false,
//           "name": "param1",
//           "type": "uint256"
//         },
//         {
//           "indexed": false,
//           "name": "param2",
//           "type": "string"
//         }
//       ],
//       "name": "EventName",
//       "type": "event"
//     }]; // Replace with your contract's event ABI

//     // Create an instance of the contract
//     const contract = new web3.eth.Contract(contractABI, contractAddress);

//     // Fetch events from the last scanned block + 1 to the latest block
//     const events = await contract.getPastEvents('EventName', {
//       fromBlock: lastScannedBlock + 1,
//       toBlock: latestBlockNumber
//     });

//     // Process the events
//     events.forEach((event) => {
//       console.log('Event:', event.returnValues);
//       // Process the event data as needed
//     });

//     // Update the last scanned block
//     lastScannedBlock = latestBlockNumber;
//   } catch (error) {
//     console.error('Error occurred while polling RPC:', error);
//   }
// }

// // Set an interval to poll the RPC every 5 seconds
// setInterval(pollRPC, 5000);
