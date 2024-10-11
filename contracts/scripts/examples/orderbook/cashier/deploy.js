const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const Cashier = await ethers.getContractFactory("Cashier");
    const instance = await upgrades.deployProxy(Cashier);
    await instance.waitForDeployment();
    console.log("Cashier Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })