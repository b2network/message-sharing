const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/nft_bridge/test.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/nft_bridge/test.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/nft_bridge/test.js --network as
     * b2: yarn hardhat run scripts/examples/nft_bridge/test.js --network b2
     */

    let address;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_NFT_BRIDGE;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_NFT_BRIDGE;
    } else if (network.name == 'as') {
        address = process.env.AS_NFT_BRIDGE;
    } else if (network.name == 'b2') {
        address = process.env.B2_NFT_BRIDGE;
    }
    console.log("NftBridge Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const NftBridge = await ethers.getContractFactory("NftBridge");
    const instance = await NftBridge.attach(address)

    let nft = {
        address: '0x1f3B35A031F712E1852260111D4d29165903824F',
        token_id: 1,
        amount: 1,
        to_chain_id: process.env.B2_DEV_CHAIN_ID,
        to_token_bridge: process.env.B2_DEV_NFT_BRIDGE,
        to_address: owner.address,
    }
    console.log("nft: ", nft);

    // function lock(address nft_address, uint256 token_id, uint256 amount, uint256 to_chain_id, address to_token_bridge, address to_address) external
    const tx = await instance.lock(nft.address, nft.token_id, nft.amount, nft.to_chain_id, nft.to_token_bridge, nft.to_address);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })