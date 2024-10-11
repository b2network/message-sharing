const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/orderbook/test.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/orderbook/test.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/orderbook/test.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/orderbook/test.js --network b2
     */

    let address;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_ORDERBOOK;
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_ORDERBOOK;
    } else if (network.name == 'as') {
        address = process.env.AS_ORDERBOOK;
    } else if (network.name == 'b2') {
        address = process.env.B2_ORDERBOOK;
    }
    console.log("Orderbook Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const Orderbook = await ethers.getContractFactory("Orderbook");
    const instance = await Orderbook.attach(address)

    let settle = {
        order_no: '10001', fee_amount: '1000'
    }
    console.log("settle: ", settle);
    let order = await instance.orders(settle.order_no);
    console.log("order: ", order);
    const tx = await instance.settle(settle.order_no, settle.fee_amount);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })