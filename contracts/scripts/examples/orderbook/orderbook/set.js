const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network b2
     */

    let address;
    let messageSharing;
    let bridges;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_ORDERBOOK;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
        bridges = [{
            chainId: process.env.B2_DEV_CHAIN_ID, address: process.env.B2_DEV_CASHIER,
        }];
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_ORDERBOOK;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_ORDERBOOK;
        messageSharing = process.env.AS_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_ORDERBOOK;
        messageSharing = process.env.B2_MESSAGE_SHARING;
    }
    console.log("Orderbook Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const Orderbook = await ethers.getContractFactory("Orderbook");
    const instance = await Orderbook.attach(address)

    // 1. setMessageSharing
    let _messageSharing = await instance.message_sharing();
    console.log("Orderbook.messageSharing:", _messageSharing);
    if (messageSharing != _messageSharing) {
        const tx = await instance.setMessageSharing(messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        _messageSharing = await instance.message_sharing()
        console.log("Orderbook.messageSharing:", _messageSharing);
    }
    console.log("1. setMessageSharing success.");

    // 2. grant role
    let role = await instance.SENDER_ROLE();
    console.log("Orderbook.SENDER_ROLE Role:", role);
    let has = await instance.hasRole(role, messageSharing)
    console.log("has role:", has)
    if (!has) {
        const tx = await instance.grantRole(role, messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        has = await instance.hasRole(role, messageSharing)
        console.log("has role:", has)
    }
    console.log("2. grantRole success.");

    // 3. setBridges
    for (const bridge of bridges) {
        let bridgeAddress = await instance.bridges(bridge.chainId);
        console.log("chainId:", bridge.chainId, ", bridgeAddress: ", bridgeAddress);
        if (bridge.address != bridgeAddress) {
            const tx = await instance.setBridges(bridge.chainId, bridge.address);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            bridgeAddress = await instance.bridges(bridge.chainId);
            console.log("chainId:", bridge.chainId, ", bridgeAddress: ", bridgeAddress);
        }
    }
    console.log("3. setBridges success.");

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })