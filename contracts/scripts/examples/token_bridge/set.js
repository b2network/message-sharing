const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/token_bridge/set.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/token_bridge/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/token_bridge/set.js --network as
     * b2: yarn hardhat run scripts/examples/token_bridge/set.js --network b2
     */

    let address;
    let messageSharing;
    let bridges;
    let tokenMaps;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_TOKEN_BRIDGE;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
        bridges = [{
            chainId: process.env.B2_DEV_CHAIN_ID, address: process.env.B2_DEV_TOKEN_BRIDGE,
        }];
        tokenMaps = [{
            from_chain_id: process.env.B2_DEV_CHAIN_ID,
            from_token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06',
            token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06',
        }];
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_TOKEN_BRIDGE;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_TOKEN_BRIDGE;
        messageSharing = process.env.AS_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_TOKEN_BRIDGE;
        messageSharing = process.env.B2_MESSAGE_SHARING;
    }
    console.log("TokenBridge Address: ", address);
    console.log("MessageSharing Address: ", messageSharing);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const TokenBridge = await ethers.getContractFactory("TokenBridge");
    const instance = await TokenBridge.attach(address)

    // 1. setMessageSharing
    let _messageSharing = await instance.messageSharing();
    if (messageSharing != _messageSharing) {
        console.log("TokenBridge.messageSharing:", _messageSharing);
        const tx = await instance.setMessageSharing(messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        _messageSharing = await instance.messageSharing()
    }
    console.log("TokenBridge.messageSharing:", _messageSharing);
    console.log("1. setMessageSharing success.");

    // 2. grant role
    let role = await instance.SENDER_ROLE();
    let has = await instance.hasRole(role, messageSharing)
    if (!has) {
        console.log("TokenBridge.SENDER_ROLE Role: ", role, ", has role: ", has)
        const tx = await instance.grantRole(role, messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`);
        has = await instance.hasRole(role, messageSharing);
    }
    console.log("TokenBridge.SENDER_ROLE Role: ", role, ", has role: ", has)
    console.log("2. grantRole success.");

    // 3. setBridges
    for (const bridge of bridges) {
        let bridgeAddress = await instance.bridges(bridge.chainId);
        if (bridge.address != bridgeAddress) {
            console.log("chainId:", bridge.chainId, ", bridgeAddress: ", bridgeAddress);
            const tx = await instance.setBridges(bridge.chainId, bridge.address);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            bridgeAddress = await instance.bridges(bridge.chainId);
            console.log("chainId:", bridge.chainId, ", bridgeAddress: ", bridgeAddress);
        }
        console.log("chainId:", bridge.chainId, ", bridgeAddress: ", bridgeAddress);
    }
    console.log("3. setBridges success.");

    // 4. setTokenMapping
    for (const tokenMap of tokenMaps) {
        console.log("tokenMap: ", tokenMap);
        let token_address = await instance.token_mapping(tokenMap.from_chain_id, tokenMap.from_token_address);
        if (tokenMap.token_address != token_address) {
            const tx = await instance.setTokenMapping(tokenMap.from_chain_id, tokenMap.from_token_address, tokenMap.token_address);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`);
        }
        console.log("token_address: ", token_address);
    }
    console.log("4. setTokenMapping success.");

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })