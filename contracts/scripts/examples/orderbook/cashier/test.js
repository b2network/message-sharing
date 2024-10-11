const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/cashier/test.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/cashier/test.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/cashier/test.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/cashier/test.js --network b2
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

    let payOrder = {
        order_no: '10001', token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06', value: 0,
    }
    console.log("payOrder: ", payOrder);
    const tx = await instance.payOrder(payOrder.order_no, payOrder.token_address, {
        value: payOrder.value
    });
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })