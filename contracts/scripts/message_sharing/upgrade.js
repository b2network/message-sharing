const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message_sharing/upgrade.js --network b2dev
     * asdev: yarn hardhat run scripts/message_sharing/upgrade.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message_sharing/upgrade.js --network b2
     * as: yarn hardhat run scripts/message_sharing/upgrade.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);
    let address = '';
    if (network.name == 'b2dev') {
        // address = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
        // address = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
        address = process.env.B2_DEV_MESSAGE_SHARING;
    } else if (network.name == 'asdev') {
        // address = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
        address = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_MESSAGE_SHARING;
    }
    console.log("MessageSharing Address: ", address);

    // Upgrading
    const MessageSharing = await ethers.getContractFactory("MessageSharing");
    const upgraded = await upgrades.upgradeProxy(address, MessageSharing);
    console.log("MessageSharing upgraded:", upgraded.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })