const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/business/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/business/upgrade.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/business/upgrade.js --network as
     * b2: yarn hardhat run scripts/business/upgrade.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let businessAddress = '';
    if (network.name == 'b2dev') {
        businessAddress = "0x804641e29f5F63a037022f0eE90A493541cCb869";
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
    } else if (network.name == 'b2') {
        businessAddress = "";
    } else if (network.name == 'as') {
        businessAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    // Upgrading
    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const upgraded = await upgrades.upgradeProxy(businessAddress, BusinessContractExample);
    console.log("BusinessContractExample upgraded:", upgraded.target);

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })