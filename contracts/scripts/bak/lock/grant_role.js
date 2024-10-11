const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/lock/grant_role.js --network b2dev
     * asdev: yarn hardhat run scripts/lock/grant_role.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/lock/grant_role.js --network as
     * b2: yarn hardhat run scripts/lock/grant_role.js --network b2
     */

    let businessAddress;
    let senderAddress;
    if (network.name == 'b2dev') {
        businessAddress = "0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b";
        senderAddress = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
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

    const TokenLockerContract = await ethers.getContractFactory("TokenLockerContract");
    const instance = await TokenLockerContract.attach(businessAddress)
    let role = await instance.SENDER_ROLE();
    console.log("role hash:", role);

    let has = await instance.hasRole(role, senderAddress)
    console.log("has role:", has)
    if (!has) {
        const tx = await instance.grantRole(role, senderAddress);
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