const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/business/set.js --network b2dev
     * asdev: yarn hardhat run scripts/business/set.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/business/set.js --network b2
     * as: yarn hardhat run scripts/business/set.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let businessAddress;
    let messageAddress;
    if (network.name == 'b2dev') {
        businessAddress = "0x804641e29f5F63a037022f0eE90A493541cCb869";
        messageAddress = '0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8';
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        messageAddress = '0x72848587deb762C4cCe38e6fA79d8347eF81b8a6';
    } else if (network.name == 'b2') {
        businessAddress = "";
    } else if (network.name == 'as') {
        businessAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const instance = await BusinessContractExample.attach(businessAddress);

    let messageSharing = await instance.messageSharing();
    if (messageSharing != messageAddress) {
        let tx = await instance.setB2MessageSharing(messageAddress);
        const txReceipt = await tx.wait(1);
        console.log("txReceipt:", txReceipt.hash);
    }
    console.log("messageSharing: ", messageSharing);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })