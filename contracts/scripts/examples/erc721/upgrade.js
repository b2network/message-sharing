const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/erc721/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc721/upgrade.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/erc721/upgrade.js --network as
     * b2: yarn hardhat run scripts/examples/erc721/upgrade.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let address;
    if (network.name == 'b2dev') {
        address = "0x1f3B35A031F712E1852260111D4d29165903824F";
    } else if (network.name == 'asdev') {
        address = "";
    } else if (network.name == 'b2') {
        address = "";
    } else if (network.name == 'as') {
        address = "";
    }
    console.log("MyERC721 Address: ", address);

    // Upgrading
    const MyERC721 = await ethers.getContractFactory("MyERC721");
    const upgraded = await upgrades.upgradeProxy(address, MyERC721);
    console.log("MyERC721 upgraded:", upgraded.target);

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })