require("@nomicfoundation/hardhat-toolbox");
require('@openzeppelin/hardhat-upgrades');
require('dotenv').config();

const AS = {
    RPC_URL: process.env.AS_RPC_URL || "", PRIVATE_KEY_LIST: [process.env.AS_PRIVATE_KEY_0 || ""],
}

const B2 = {
    RPC_URL: process.env.B2_RPC_URL || "", PRIVATE_KEY_LIST: [process.env.B2_PRIVATE_KEY_0 || ""],
}

const B2_DEV = {
    RPC_URL: process.env.B2_DEV_RPC_URL || "", PRIVATE_KEY_LIST: [process.env.B2_DEV_PRIVATE_KEY_0 || ""],
}

const AS_DEV = {
    RPC_URL: process.env.AS_DEV_RPC_URL || "", PRIVATE_KEY_LIST: [process.env.AS_DEV_PRIVATE_KEY_0 || ""],
}

task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
    const accounts = await hre.ethers.getSigners()
    for (const account of accounts) {
        console.log(account.address)
    }
})

module.exports = {
    networks: {
        asdev: {
            blockConfirmations: 1, url: AS_DEV.RPC_URL, accounts: AS_DEV.PRIVATE_KEY_LIST,
        }, b2dev: {
            blockConfirmations: 1, url: B2_DEV.RPC_URL, accounts: B2_DEV.PRIVATE_KEY_LIST, gasPrice: 352
        }, as: {
            blockConfirmations: 1, url: AS.RPC_URL, accounts: AS.PRIVATE_KEY_LIST,
        }, b2: {
            blockConfirmations: 1, url: B2.RPC_URL, accounts: B2.PRIVATE_KEY_LIST, gasPrice: 352,
        }, hardhat: {},
    }, solidity: {
        version: "0.8.20", settings: {
            optimizer: {
                enabled: true, runs: 1000,
            },
        },
    }, etherscan: {
        apiKey: {
            b2test: "abc", b2: "abc"
        }, customChains: [{
            network: "b2test", chainId: 1123, urls: {
                apiURL: "https://testnet-backend.bsquared.network/api",
                browserURL: "https://testnet-explorer.bsquared.network"
            }
        }, {
            network: "b2", chainId: 223, urls: {
                apiURL: "https://mainnet-backend.bsquared.network/api",
                browserURL: "https://mainnet-blockscout.bsquared.network"
                // apiURL: "https://bsquared.l2scan.co/api", browserURL: "https://bsquared.l2scan.co"
            }
        }]
    }
}