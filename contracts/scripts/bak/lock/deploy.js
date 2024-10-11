const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/lock/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/lock/deploy.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/lock/deploy.js --network as
     * b2: yarn hardhat run scripts/lock/deploy.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const TokenLockerContract = await ethers.getContractFactory("TokenLockerContract");
    const instance = await upgrades.deployProxy(TokenLockerContract);
    await instance.waitForDeployment();
    console.log("TokenLockerContract Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })