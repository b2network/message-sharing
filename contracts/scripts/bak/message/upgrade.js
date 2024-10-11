const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/message/upgrade.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message/upgrade.js --network b2
     * as: yarn hardhat run scripts/message/upgrade.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);
    let messageAddress = '';
    if (network.name == 'b2dev') {
        // messageAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
        messageAddress = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
    } else if (network.name == 'asdev') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'b2') {
        messageAddress = "";
    } else if (network.name == 'as') {
        messageAddress = "";
    }
    console.log("Message Address: ", messageAddress);

    // Upgrading
    const B2MessageSharing = await ethers.getContractFactory("B2MessageSharing");
    const upgraded = await upgrades.upgradeProxy(messageAddress, B2MessageSharing);
    console.log("B2MessageSharing upgraded:", upgraded.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })