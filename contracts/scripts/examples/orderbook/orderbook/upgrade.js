const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/orderbook/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/orderbook/upgrade.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/orderbook/upgrade.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/orderbook/upgrade.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let address;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_ORDERBOOK;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_ORDERBOOK;
    } else if (network.name == 'as') {
        address = process.env.AS_ORDERBOOK;
    } else if (network.name == 'b2') {
        address = process.env.B2_ORDERBOOK;
    }
    console.log("Orderbook Address: ", address);

    // Upgrading
    const Orderbook = await ethers.getContractFactory("Orderbook");
    const upgraded = await upgrades.upgradeProxy(address, Orderbook);
    console.log("Orderbook upgraded:", upgraded.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })