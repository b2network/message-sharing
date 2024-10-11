const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/message_channel/set.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/message_channel/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/message_channel/set.js --network as
     * b2: yarn hardhat run scripts/examples/message_channel/set.js --network b2
     */

    let address;
    let messageSharing;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_MESSAGE_CHANNEL;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_MESSAGE_CHANNEL;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.B2_MESSAGE_CHANNEL;
        messageSharing = process.env.B2_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.AS_MESSAGE_CHANNEL;
        messageSharing = process.env.AS_MESSAGE_SHARING;
    }
    console.log("MessageChannel Address: ", address);
    console.log("messageSharing Address: ", messageSharing);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const MessageChannel = await ethers.getContractFactory("MessageChannel");
    const instance = await MessageChannel.attach(address)

    // 1. setMessageSharing
    let _messageSharing = await instance.messageSharing();
    console.log("MessageChannel.messageSharing:", _messageSharing);
    if (messageSharing != _messageSharing) {
        const tx = await instance.setMessageSharing(messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        _messageSharing = await instance.messageSharing()
        console.log("MessageChannel.messageSharing:", _messageSharing);
    }

    // 2. grant role
    let role = await instance.SENDER_ROLE();
    console.log("MessageChannel.SENDER_ROLE Role:", role);

    let has = await instance.hasRole(role, messageSharing)
    console.log("has role:", has)
    if (!has) {
        const tx = await instance.grantRole(role, messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        has = await instance.hasRole(role, messageSharing)
        console.log("has role:", has)
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })