const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/token_bridge/test.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/token_bridge/test.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/token_bridge/test.js --network as
     * b2: yarn hardhat run scripts/examples/token_bridge/test.js --network b2
     */

    let address;
    let lock;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_TOKEN_BRIDGE;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_TOKEN_BRIDGE;
    } else if (network.name == 'as') {
        address = process.env.AS_TOKEN_BRIDGE;
    } else if (network.name == 'b2') {
        address = process.env.B2_TOKEN_BRIDGE;
    }

    console.log("TokenBridge Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const TokenBridge = await ethers.getContractFactory("TokenBridge");
    const instance = await TokenBridge.attach(address)

    // 1. lock
    lock = {
        token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06',
        // token_address: '0x0000000000000000000000000000000000000000',
        amount: 1000,
        to_chain_id: process.env.B2_DEV_CHAIN_ID,
        to_token_bridge: process.env.B2_DEV_TOKEN_BRIDGE,
        to_address: owner.address,
    };
    console.log("lock: ", lock);
    let value = 0;
    if (lock.token_address == '0x0000000000000000000000000000000000000000') {
        value = lock.amount;
    }
    const tx = await instance.lock(lock.token_address, lock.amount, lock.to_chain_id, lock.to_token_bridge, lock.to_address, {
        value: value
    });
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })