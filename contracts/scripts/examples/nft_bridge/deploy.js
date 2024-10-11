const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/nft_bridge/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/nft_bridge/deploy.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/nft_bridge/deploy.js --network as
     * b2: yarn hardhat run scripts/examples/nft_bridge/deploy.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const NftBridge = await ethers.getContractFactory("NftBridge");
    const instance = await upgrades.deployProxy(NftBridge);
    await instance.waitForDeployment();
    console.log("NftBridge Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })