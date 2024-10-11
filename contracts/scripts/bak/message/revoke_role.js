const {ethers, network} = require("hardhat");


async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message/revoke_role.js --network b2dev
     * asdev: yarn hardhat run scripts/message/revoke_role.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message/revoke_role.js --network b2
     * as: yarn hardhat run scripts/message/revoke_role.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address);
    let messageAddress;
    if (network.name == 'b2dev') {
        messageAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
    } else if (network.name == 'asdev') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'b2') {
        messageAddress = "";
    } else if (network.name == 'as') {
        messageAddress = "";
    }
    console.log("Message Address: ", messageAddress);
    // B2MessageSharing
    const B2MessageSharing = await ethers.getContractFactory("B2MessageSharing");
    const instance = await B2MessageSharing.attach(messageAddress);
    // TODO
    let role = await instance.ADMIN_ROLE(); // admin role
    // role = await instance.UPGRADE_ROLE(); // upgrade role
    // role = await instance.validatorRole(1); // validatorRole(uint256 chain_id)
    console.log("role hash: ", role);

    // TODO
    let accounts = [];
    for (const account of accounts) {
        let hasRole = await instance.hasRole(role, account);
        console.log("account: ", account, " => hasRole:", hasRole)
        if (hasRole) {
            const tx = await instance.revokeRole(role, account);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            hasRole = await instance.hasRole(role, account)
            console.log("account: ", account, " => hasRole:", hasRole)
        }
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })