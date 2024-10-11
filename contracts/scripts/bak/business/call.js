const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/business/call.js --network b2dev
     * asdev: yarn hardhat run scripts/business/call.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/business/call.js --network b2
     * as: yarn hardhat run scripts/business/call.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let businessAddress;
    let to_chain_id;
    let to_business_address;
    let data;
    if (network.name == 'b2dev') {
        businessAddress = "0x804641e29f5F63a037022f0eE90A493541cCb869";
        to_chain_id = 421614;
        to_business_address = '0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534';
        data = '0x1234';
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        to_chain_id = 1123;
        to_business_address = '0x804641e29f5F63a037022f0eE90A493541cCb869';
        data = '0x5678';
    } else if (network.name == 'b2') {
        businessAddress = "";
    } else if (network.name == 'as') {
        businessAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const instance = await BusinessContractExample.attach(businessAddress);

    let tx = await instance.call(to_chain_id, to_business_address, data);
    const txReceipt = await tx.wait(1);
    console.log("txReceipt:", txReceipt.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })