const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const Orderbook = await ethers.getContractFactory("Orderbook");
    const instance = await upgrades.deployProxy(Orderbook);
    await instance.waitForDeployment();
    console.log("Orderbook Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })