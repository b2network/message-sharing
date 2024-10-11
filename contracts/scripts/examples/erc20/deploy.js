const {ethers, upgrades} = require("hardhat")

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/examples/erc20/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc20/deploy.js --network asdev
     */
    const [owner] = await ethers.getSigners();
    console.log(owner.address);

    let token = {
        name: 'USDT', symbol: 'USDT', decimals: 18,
    };

    const MyERC20 = await ethers.getContractFactory("MyERC20");
    const instance = await MyERC20.deploy(token.name, token.symbol, token.decimals);
    await instance.waitForDeployment();
    console.log("MyERC20 address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    });