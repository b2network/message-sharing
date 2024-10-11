const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/lock/lock.js --network b2dev
     * asdev: yarn hardhat run scripts/lock/lock.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/lock/lock.js --network as
     * b2: yarn hardhat run scripts/lock/lock.js --network b2
     */

    let businessAddress;
    if (network.name == 'b2dev') {
        businessAddress = "0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b";
    } else if (network.name == 'asdev') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
    } else if (network.name == 'as') {
        businessAddress = "";
    } else if (network.name == 'b2') {
        businessAddress = "";
    }
    console.log("Business Address: ", businessAddress);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const TokenLockerContract = await ethers.getContractFactory("TokenLockerContract");
    const instance = await TokenLockerContract.attach(businessAddress)

    let token_address = '0xc810b0b75Af60D60De8451587DF9cb240BE22d9d';
    let amount = '1000';
    let to_chain_id = 1123;
    let to_business_contract = '0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b';
    let to_address = '0x502FA825441D215EDECc54804d74f3FBFe20fb97';
    const tx = await instance.lock(token_address, amount, to_chain_id, to_business_contract, to_address);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })