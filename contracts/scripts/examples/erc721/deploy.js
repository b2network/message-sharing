const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/erc721/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc721/deploy.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/erc721/deploy.js --network as
     * b2: yarn hardhat run scripts/examples/erc721/deploy.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const MyERC721 = await ethers.getContractFactory("MyERC721");
    const instance = await upgrades.deployProxy(MyERC721, ["Test Dog", "TDT"]);
    await instance.waitForDeployment();
    console.log("MyERC721 Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })