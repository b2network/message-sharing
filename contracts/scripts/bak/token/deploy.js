const {ethers, upgrades} = require("hardhat")

const TokenName = "USDT"
const TokenSymbol = "USDT"
const decimals = 18

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/token/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/token/deploy.js --network asdev
     */
    const [owner] = await ethers.getSigners();
    console.log(owner.address);

    const MyERC20 = await ethers.getContractFactory("MyERC20");
    const instance = await MyERC20.deploy(TokenName, TokenSymbol, decimals);
    await instance.waitForDeployment();
    console.log("MyERC20 V1:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })