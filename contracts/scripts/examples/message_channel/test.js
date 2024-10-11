const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/message_channel/test.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/message_channel/test.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/message_channel/test.js --network as
     * b2: yarn hardhat run scripts/examples/message_channel/test.js --network b2
     */

    let address;
    let messageSharing;
    let param;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_MESSAGE_CHANNEL;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
        param = {
            to_chain_id: process.env.AS_DEV_CHAIN_ID,
            to_business_contract: process.env.AS_DEV_MESSAGE_CHANNEL,
            to_message: '0x12345678',
        };
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_MESSAGE_CHANNEL;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
        param = {
            to_chain_id: process.env.B2_DEV_CHAIN_ID,
            to_business_contract: process.env.B2_DEV_MESSAGE_CHANNEL,
            to_message: '0x87654321',
        };
    } else if (network.name == 'as') {
        address = process.env.B2_MESSAGE_CHANNEL;
        messageSharing = process.env.B2_MESSAGE_SHARING;
        param = {
            to_chain_id: process.env.AS_CHAIN_ID,
            to_business_contract: process.env.AS_MESSAGE_CHANNEL,
            to_message: '0x12345678',
        };
    } else if (network.name == 'b2') {
        address = process.env.AS_MESSAGE_CHANNEL;
        messageSharing = process.env.AS_MESSAGE_SHARING;
        param = {
            to_chain_id: process.env.B2_CHAIN_ID,
            to_business_contract: process.env.B2_MESSAGE_CHANNEL,
            to_message: '0x87654321',
        };
    }
    console.log("MessageChannel Address: ", address);
    console.log("MessageSharing Address: ", messageSharing);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const MessageChannel = await ethers.getContractFactory("MessageChannel");
    const instance = await MessageChannel.attach(address)

    // call
    console.log("param: ", param);
    const tx = await instance.call(param.to_chain_id, param.to_business_contract, param.to_message);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })