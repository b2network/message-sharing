const {ethers, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message_sharing/grant_role.js --network b2dev
     * asdev: yarn hardhat run scripts/message_sharing/grant_role.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message_sharing/grant_role.js --network b2
     * as: yarn hardhat run scripts/message_sharing/grant_role.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address);
    let address;
    if (network.name == 'b2dev') {
        // address = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
        // address = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
        address = process.env.B2_DEV_MESSAGE_SHARING;
    } else if (network.name == 'asdev') {
        // address = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
        // address = "0x72848587deb762C4cCe38e6fA79d8347eF81b8a6";
        address = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_MESSAGE_SHARING;
    }
    console.log("MessageSharing Address: ", address);
    // MessageSharing
    const MessageSharing = await ethers.getContractFactory("MessageSharing");
    const instance = await MessageSharing.attach(address);

    // let role = await instance.ADMIN_ROLE(); // admin role
    // let role = await instance.UPGRADE_ROLE(); // upgrade role
    // let role = await instance.validatorRole(0);
    // let role = await instance.validatorRole(421614);
    let role = await instance.validatorRole(1123);
    console.log("role hash: ", role);

    let accounts = ["0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d"];
    for (const account of accounts) {
        let hasRole = await instance.hasRole(role, account);
        console.log("account: ", account, " => hasRole: ", hasRole)
        if (!hasRole) {
            const tx = await instance.grantRole(role, account);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            hasRole = await instance.hasRole(role, account)
            console.log("account: ", account, " => hasRole: ", hasRole)
        }
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
