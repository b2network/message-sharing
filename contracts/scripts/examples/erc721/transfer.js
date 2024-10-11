const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/examples/erc721/transfer.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc721/transfer.js --network asdev
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

    let transfer = {
        token_id: 1, // to_address: owner.address,
        to_address: '0x952b63C6C799B7033c24B055f7F023Eb7f3a5c73',
    };
    console.log("transfer: ", transfer);
    let exist = await instance.existOf(transfer.token_id);
    console.log("exist: ", exist);
    if (exist) {
        let _owner = await instance.ownerOf(transfer.token_id);
        if (_owner == owner.address) {
            const tx = await instance.safeTransferFrom(owner.address, transfer.to_address, transfer.token_id);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
        }
    }

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })