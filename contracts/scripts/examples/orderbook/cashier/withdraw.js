const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/cashier/withdraw.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/cashier/withdraw.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/cashier/withdraw.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/cashier/withdraw.js --network b2
     */

    let address;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_CASHIER;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_CASHIER;
    } else if (network.name == 'as') {
        address = process.env.AS_CASHIER;
    } else if (network.name == 'b2') {
        address = process.env.B2_CASHIER;
    }
    console.log("Cashier Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const Cashier = await ethers.getContractFactory("Cashier");
    const instance = await Cashier.attach(address)

    let withdraw = {
        token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06',
        to_address: '0x0000000000000000000000000000000000000001',
        amount: 1,
    }
    console.log("withdraw: ", withdraw);
    let withdraw_balance = await instance.withdraw_balance();
    console.log("withdraw_balance:", withdraw_balance);
    const tx = await instance.withdraw(withdraw.token_address, withdraw.to_address, withdraw.amount);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })