const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/nft_bridge/set.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/nft_bridge/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/nft_bridge/set.js --network as
     * b2: yarn hardhat run scripts/examples/nft_bridge/set.js --network b2
     */

    let address;
    let messageSharing;
    let bridges;
    let nftMaps;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_NFT_BRIDGE;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
        bridges = [{
            chainId: process.env.B2_DEV_CHAIN_ID, address: process.env.B2_DEV_NFT_BRIDGE,
        }];
        nftMaps = [{
            from_chain_id: process.env.B2_DEV_CHAIN_ID,
            from_nft_address: '0x1f3B35A031F712E1852260111D4d29165903824F',
            to_chain_id: process.env.B2_DEV_CHAIN_ID,
            to_nft_address: '0x1f3B35A031F712E1852260111D4d29165903824F',
            standard: 1,
            status: true,
        }];
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_NFT_BRIDGE;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_DEV_NFT_BRIDGE;
        messageSharing = process.env.AS_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_DEV_NFT_BRIDGE;
        messageSharing = process.env.B2_MESSAGE_SHARING;
    }
    console.log("NftBridge Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const NftBridge = await ethers.getContractFactory("NftBridge");
    const instance = await NftBridge.attach(address)

    // 1. setMessageSharing
    let _messageSharing = await instance.messageSharing();
    console.log("NftBridge.messageSharing:", _messageSharing);
    if (messageSharing != _messageSharing) {
        const tx = await instance.setMessageSharing(messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        _messageSharing = await instance.messageSharing()
        console.log("NftBridge.messageSharing:", _messageSharing);
    }
    console.log("1. setMessageSharing success.");

    // 2. grant role
    let role = await instance.SENDER_ROLE();
    console.log("NftBridge.SENDER_ROLE Role:", role);
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

    // 4. setNftMapping
    for (const nftMap of nftMaps) {
        console.log("nftMap: ", nftMap);
        let _nftMap = await instance.nft_mapping(nftMap.from_chain_id, nftMap.from_nft_address, nftMap.to_chain_id);
        console.log("_nftMap: ", _nftMap);
        if (_nftMap.nft_address != nftMap.to_nft_address || _nftMap.standard != nftMap.standard || _nftMap.status != nftMap.status) {
            const tx = await instance.setNftMapping(nftMap.from_chain_id, nftMap.from_nft_address, nftMap.to_chain_id, nftMap.to_nft_address, nftMap.standard, nftMap.status);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            _nftMap = await instance.nft_mapping(nftMap.from_chain_id, nftMap.from_nft_address, nftMap.to_chain_id);
            console.log("_nftMap:", _nftMap);
        }
    }
    console.log("4. setNftMapping success.");

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })