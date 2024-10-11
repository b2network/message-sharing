const {ethers, run, network} = require("hardhat")


async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/erc20/mint.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc20/mint.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/examples/erc20/mint.js --network b2
     * as: yarn hardhat run scripts/examples/erc20/mint.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address);

    let tokenAddress;
    if (network.name == 'b2dev') {
        tokenAddress = "0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06";
    } else if (network.name == 'asdev') {
        tokenAddress = "";
    } else if (network.name == 'b2') {
        tokenAddress = "";
    } else if (network.name == 'as') {
        tokenAddress = "";
    }
    console.log("Token Address: ", tokenAddress);

    const MyERC20 = await ethers.getContractFactory('MyERC20');
    const instance = MyERC20.attach(tokenAddress);

    let mint = {
        // to: '0xa218c247a24Cd3Eb2D53E0cb827f52540Cc2312E',
        // amount: '100000000000000',
        to: owner.address, amount: '100000000000000',
    };
    console.log("mint: ", mint);
    const tx = await instance.mint(mint.to, mint.amount)
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
    console.log("balance of:", await instance.balanceOf(mint.to));

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })