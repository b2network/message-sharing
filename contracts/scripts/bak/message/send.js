const {ethers, upgrades, network} = require("hardhat");

```
# Listener

```

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/send.js --network b2dev
     * asdev: yarn hardhat run scripts/message/send.js --network asdev
     * b2: yarn hardhat run scripts/message/send.js --network b2
     * as: yarn hardhat run scripts/message/send.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    let messageAddress;
    if (network.name == 'b2dev') {
        // messageAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
        messageAddress = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
    } else if (network.name == 'asdev') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'b2') {
        messageAddress = "";
    } else if (network.name == 'as') {
        messageAddress = "";
    }
    console.log("Message Address: ", messageAddress);

    const B2MessageSharing = await ethers.getContractFactory("B2MessageSharing");
    const instance = await B2MessageSharing.attach(messageAddress);

    let from_chain_id = 421614;
    let from_id = '0x000000000000000000000000000000000000000000000000000000000000000a';
    let from_sender = "0x98C6e991D1b338604D4Fa10F351a27012eFe8eC2";
    let to_chain_id = 1123;
    let contract_address = '0x804641e29f5F63a037022f0eE90A493541cCb869';
    let data = '0x1234';
    let signatures = ['0x27bda5470df8273d66f40fc50f4f6cd7b79f890a02383519c9e0315cbedc180b283381744a10fa7924ceca44faf947dbc3c94aaf059ed7d802393b58b490db321b'];

    let weight = 0;
    for (const signature of signatures) {
        let verify = await instance.verify(from_chain_id, from_id, from_sender, to_chain_id, contract_address, data, signature);
        console.log("verify:", verify);
        if (verify) {
            weight = weight + 1;
        }
    }
    let _weight = await instance.weights(from_chain_id);
    console.log("weight:", _weight);
    // if (weight >= _weight) {
    //     let sendTx = await instance.send(from_chain_id, from_id, from_sender, contract_address, data, signatures);
    //     const sendTxReceipt = await sendTx.wait(1);
    //     console.log("sendTxReceipt:", sendTxReceipt.hash);
    // }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })