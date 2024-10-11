const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/business/revoke_role.js --network b2dev
     * asdev: yarn hardhat run scripts/business/revoke_role.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/business/revoke_role.js --network as
     * b2: yarn hardhat run scripts/business/revoke_role.js --network b2
     */

    let businessAddress;
    let senderAddress;

    if (network.name == 'b2dev') {
        businessAddress = "0x804641e29f5F63a037022f0eE90A493541cCb869";
        senderAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        senderAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'as') {
        businessAddress = "";
        senderAddress = "";
    } else if (network.name == 'b2') {
        businessAddress = "";
        senderAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const instance = await BusinessContractExample.attach(businessAddress)
    let role = await instance.SENDER_ROLE();
    console.log("role hash:", role);

    let has = await instance.hasRole(role, senderAddress)
    console.log("has role:", has)
    if (has) {
        const tx = await instance.revokeRole(role, senderAddress);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        has = await instance.hasRole(role, senderAddress)
        console.log("has role:", has)
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })