const {ethers, run, network} = require("hardhat")

const tokenAddress = "0xc810b0b75Af60D60De8451587DF9cb240BE22d9d"
const toAccount = "0x690bC18DfAA4C5f1cC67495781B90FC4D90cD78b"
const approveAmount = "10000000000000000000000000000000000000000"

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/token/approve.js --network b2dev
     * asdev: yarn hardhat run scripts/token/approve.js --network asdev
     */

    const [owner] = await ethers.getSigners()

    const MyERC20 = await ethers.getContractFactory('MyERC20');
    const token = MyERC20.attach(tokenAddress);

    const approveTx = await token.approve(toAccount, approveAmount)
    const txReceipt = await approveTx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
    console.log("allowance:", await token.allowance(owner.address, toAccount));

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })