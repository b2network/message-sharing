const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message/call.js --network b2dev
     * asdev: yarn hardhat run scripts/message/call.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message/call.js --network b2
     * as: yarn hardhat run scripts/message/call.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let messageAddress;
    let businessAddress;
    let to_chain_id;
    if (network.name == 'b2dev') {
        // messageAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
        messageAddress = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        to_chain_id = 421614;
    } else if (network.name == 'asdev') {
        messageAddress = "0x72848587deb762C4cCe38e6fA79d8347eF81b8a6";
        businessAddress = "0x804641e29f5F63a037022f0eE90A493541cCb869";
        to_chain_id = 1123;
    } else if (network.name == 'b2') {
        messageAddress = "";
        businessAddress = "";
    } else if (network.name == 'as') {
        messageAddress = "";
        businessAddress = "";
    }
    console.log("Message Address: ", messageAddress);
    console.log("Business Address: ", businessAddress);
    // MessageSharing.sol
    // TODO
    let data = '0x1234';
    const B2MessageSharing = await ethers.getContractFactory("B2MessageSharing");
    const instance = await B2MessageSharing.attach(messageAddress);

    let tx = await instance.call(to_chain_id, businessAddress, data);
    const txReceipt = await tx.wait(1);
    console.log("txReceipt:", txReceipt.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })