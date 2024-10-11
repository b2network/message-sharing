const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/lock/set.js --network b2dev
     * asdev: yarn hardhat run scripts/lock/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/lock/set.js --network as
     * b2: yarn hardhat run scripts/lock/set.js --network b2
     */

    let businessAddress;
    let messageAddress;
    if (network.name == 'b2dev') {
        businessAddress = "0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b";
        messageAddress = "0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8";
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'as') {
        businessAddress = "";
        messageAddress = "";
    } else if (network.name == 'b2') {
        businessAddress = "";
        messageAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const TokenLockerContract = await ethers.getContractFactory("TokenLockerContract");
    const instance = await TokenLockerContract.attach(businessAddress)
    let messageSharing = await instance.messageSharing();
    console.log("messageSharing address:", messageSharing);

    if (messageSharing != messageAddress) {
        const tx = await instance.setB2MessageSharing(messageAddress);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }

    let chain_id = 1123;
    let locker = '0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b';

    let _locker = await instance.locks(chain_id);
    if (_locker != locker) {
        const tx = await instance.setLocks(chain_id, locker);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }

    let token_address = '0xc810b0b75Af60D60De8451587DF9cb240BE22d9d';
    let to_token_address = '0xc810b0b75Af60D60De8451587DF9cb240BE22d9d';

    let _to_token_address = await instance.tokens(chain_id, token_address);
    if (_to_token_address != to_token_address) {
        const tx = await instance.setTokens(chain_id, token_address, to_token_address);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }



}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })