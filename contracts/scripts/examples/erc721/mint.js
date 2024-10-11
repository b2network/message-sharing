const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/examples/erc721/mint.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc721/mint.js --network asdev
     */

    let nftAddress;
    if (network.name == 'b2dev') {
        nftAddress = "0x1f3B35A031F712E1852260111D4d29165903824F"
    } else if (network.name == 'asdev') {
        nftAddress = "";
    }
    console.log("Nft Address: ", nftAddress)

    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address)

    const MyERC721 = await ethers.getContractFactory("MyERC721");
    const instance = await MyERC721.attach(nftAddress)

    let mint = {
        token_id: 2, // to_address: owner.address,
        to_address: '0x952b63C6C799B7033c24B055f7F023Eb7f3a5c73',
    };
    console.log("mint: ", mint);
    let exist = await instance.existOf(mint.token_id);
    console.log("exist: ", exist);
    if (!exist) {
        const tx = await instance.mint(mint.to_address, mint.token_id);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })