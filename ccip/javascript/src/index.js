const express = require("express");
const app = express();
const port = 4444;
const Joi = require("joi");

app.use(express.json());

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});

app.post("/ccip", (req, res) => {
  // this request only comes from our server so we can trust it, no need to validate
  // node src/transfer-tokens.js ethereumSepolia polygonMumbai 0x9d087fC03ae39b088326b67fA3C788236645b717 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05 1000000000000000 0xdfe59eb3c16344e2d1daeda3611ba43116926ad4c753ab046b855fcfac883e4c

  // const sourceChain = process.argv[2];
  // const destinationChain = process.argv[3];
  // const destinationAccount = process.argv[4];
  // const tokenAddress = process.argv[5];
  // const amount = ethers.BigNumber.from(process.argv[6]);
  // const privateKey = process.argv[7];
  // const feeTokenAddress = process.argv[8];

  /* body is like {
         "soureChain": "ethereumSepolia",
            "destinationChain": "polygonMumbai",
            "destinationAccount": "0x9d087fC03ae39b088326b67fA3C788236645b717",
            "tokenAddress": "0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05",
            "amount": "1000000000000000",
            "privateKey": "0xdfe59eb3c16344e2d1daeda3611ba43116926ad4c753ab046b855fcfac883e4c,
            "feeTokenAddress": nullable || "0x00000000"
            "isNative": false
    }
    */

  //   const schema = Joi.object({
  //     sourceChain: Joi.string().required(),
  //     destinationChain: Joi.string().required(),
  //     destinationAccount: Joi.string().required(),
  //     tokenAddress: Joi.string().required(),
  //     amount: Joi.string().required(),
  //     privateKey: Joi.string().required(),
  //     feeTokenAddress: Joi.string().allow(null).required(),
  //     isNative: Joi.boolean().required(),
  //   });

  const schema = Joi.object({
    sourceChain: Joi.string().required(),
    destinationChain: Joi.string().required(),
    destinationAccount: Joi.string().required(),
    tokenAddress: Joi.string()
      .required()
      .pattern(/^0x[a-fA-F0-9]{40}$/), // Ethereum address pattern
    amount: Joi.string().required(),
    privateKey: Joi.string().required(),
    feeTokenAddress: Joi.string().allow(null).required(),
  });

  const { error, value } = schema.validate(req.body);
  if (error) {
    console.log(error);
    return res.status(400).json({
      error: error.details.map((err) => err.message).join(", "),
      status: false,
      data: null,
    });
  }

  console.log(value);

  // get validated data from request body
  const {
    sourceChain,
    destinationChain,
    destinationAccount,
    tokenAddress,
    amount,
    privateKey,
    feeTokenAddress,
  } = req.body;

  // transfer tokens
  const { transferTokens } = require("./transfer-tokens");

  transferTokens({
    sourceChain,
    destinationChain,
    destinationAccount,
    tokenAddress,
    amount,
    privateKey,
    feeTokenAddress,
  })
    .then((result) => {
      return res.status(200).json({
        error: null,
        status: true,
        data: result,
      });
    })
    .catch((error) => {
      return res.status(400).json({
        error: error.message,
        status: false,
        data: null,
      });
    });
});
