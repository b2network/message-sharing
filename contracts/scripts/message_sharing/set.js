const {ethers, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message_sharing/set.js --network b2dev
     * asdev: yarn hardhat run scripts/message_sharing/set.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message_sharing/set.js --network b2
     * as: yarn hardhat run scripts/message_sharing/set.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address);

    let address;
    let weight;
    let validators;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_MESSAGE_SHARING;
        weight = {
            chain_id: process.env.B2_DEV_CHAIN_ID, weight: 1,
        };
        validators = [{
            chain_id: process.env.B2_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }, {
            chain_id: process.env.AS_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }];
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_MESSAGE_SHARING;
        weight = {
            chain_id: process.env.B2_DEV_CHAIN_ID, weight: 1,
        };
        validators = [{
            chain_id: process.env.B2_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }, {
            chain_id: process.env.AS_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }];
    } else if (network.name == 'b2') {
        address = process.env.B2_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_MESSAGE_SHARING;
    }
    console.log("MessageSharing Address: ", address);
    // MessageSharing
    const MessageSharing = await ethers.getContractFactory("MessageSharing");
    const instance = await MessageSharing.attach(address);


    // 1. setWeight
    for (const setup of setups) {
        let _weight = await instance.weights(weight.chain_id);
        if (weight.weight != _weight) {
            const tx = await instance.setWeight(weight.chain_id, weight.weight);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            _weight = await instance.weights(weight.chain_id);
        }
        console.log("chain_id:", weight.chain_id, ", weight: ", weight.weight)
    }
    console.log("1. setWeight success.")

    // 2. setValidatorRole
    for (const validator of validators) {
        const tx = await instance.setValidatorRole(validator.chain_id, validator.account, validator.valid);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }
    console.log("2. setValidatorRole success.")
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
