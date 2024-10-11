const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message_sharing/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/message_sharing/deploy.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message_sharing/deploy.js --network b2
     * as: yarn hardhat run scripts/message_sharing/deploy.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    const MessageSharing = await ethers.getContractFactory("MessageSharing");
    const instance = await upgrades.deployProxy(MessageSharing);
    await instance.waitForDeployment();
    console.log("MessageSharing Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })