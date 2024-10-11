const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/lock/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/lock/upgrade.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/lock/upgrade.js --network as
     * b2: yarn hardhat run scripts/lock/upgrade.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let lockerAddress = '';
    if (network.name == 'b2dev') {
        lockerAddress = "0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b";
    } else if (network.name == 'asdev') {
        lockerAddress = "";
    } else if (network.name == 'b2') {
        lockerAddress = "";
    } else if (network.name == 'as') {
        lockerAddress = "";
    }
    console.log("locker Address: ", lockerAddress);

    // Upgrading
    const TokenLockerContract = await ethers.getContractFactory("TokenLockerContract");
    const upgraded = await upgrades.upgradeProxy(lockerAddress, TokenLockerContract);
    console.log("TokenLockerContract upgraded:", upgraded.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })