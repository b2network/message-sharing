const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/examples/erc721/setApprovalForAll.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc721/setApprovalForAll.js --network asdev
     */

    let nftAddress;
    let nftBridge;
    if (network.name == 'b2dev') {
        nftAddress = "0x1f3B35A031F712E1852260111D4d29165903824F";
        nftBridge = process.env.B2_DEV_NFT_BRIDGE;
    } else if (network.name == 'asdev') {
        nftAddress = "";
        nftBridge = process.env.AS_DEV_NFT_BRIDGE;
    }
    console.log("Nft Address: ", nftAddress);
    console.log("NftBridge Address: ", nftBridge);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address: ", owner.address)

    const MyERC721 = await ethers.getContractFactory("MyERC721");
    const instance = await MyERC721.attach(nftAddress)
    let approved = await instance.isApprovedForAll(owner.address, nftBridge);
    console.log("approved: ", approved);
    if (!approved) {
        const tx = await instance.setApprovalForAll(nftBridge, true);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        approved = await instance.isApprovedForAll(owner.address, nftBridge);
        console.log("approved: ", approved);
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })