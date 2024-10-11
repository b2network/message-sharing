const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/orderbook/cashier/set.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/orderbook/cashier/set.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/orderbook/cashier/set.js --network as
     * b2: yarn hardhat run scripts/examples/orderbook/cashier/set.js --network b2
     */

    let address;
    let messageSharing;
    let orderbook;
    let wihitelists;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_CASHIER;
        messageSharing = process.env.B2_DEV_MESSAGE_SHARING;
        orderbook = {
            chain_id: process.env.B2_DEV_CHAIN_ID, contract_address: process.env.B2_DEV_ORDERBOOK,
        };
        wihitelists = [{
            token_address: '0x0000000000000000000000000000000000000000', deposit_amount: '1000000', status: true,
        }, {
            token_address: '0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06', deposit_amount: '1000000', status: true,
        }];
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_CASHIER;
        messageSharing = process.env.AS_DEV_MESSAGE_SHARING;
    } else if (network.name == 'as') {
        address = process.env.AS_DEV_CASHIER;
        messageSharing = process.env.AS_MESSAGE_SHARING;
    } else if (network.name == 'b2') {
        address = process.env.B2_DEV_CASHIER;
        messageSharing = process.env.B2_MESSAGE_SHARING;
    }
    console.log("Cashier Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const Cashier = await ethers.getContractFactory("Cashier");
    const instance = await Cashier.attach(address)

    // 1. setMessageSharing
    let _messageSharing = await instance.message_sharing();
    console.log("Cashier.messageSharing:", _messageSharing);
    if (messageSharing != _messageSharing) {
        const tx = await instance.setMessageSharing(messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        _messageSharing = await instance.message_sharing()
        console.log("Cashier.messageSharing:", _messageSharing);
    }
    console.log("1. setMessageSharing success.");

    // 2. grant role
    let role = await instance.SENDER_ROLE();
    console.log("Cashier.SENDER_ROLE Role:", role);
    let has = await instance.hasRole(role, messageSharing)
    console.log("has role:", has)
    if (!has) {
        const tx = await instance.grantRole(role, messageSharing);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        has = await instance.hasRole(role, messageSharing)
        console.log("has role:", has)
    }
    console.log("2. grantRole success.");

    // 3. setOrderbook
    let _orderbook = await instance.orderbook();
    console.log("_orderbook: ", _orderbook);
    if (_orderbook.chain_id != orderbook.chain_id || _orderbook.contract_address != orderbook.contract_address) {
        const tx = await instance.setOrderbook(orderbook.contract_address, orderbook.chain_id);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
    }
    console.log("3. setOrderbook success.");

    // 4. setWihitelist
    for (const wihitelist of wihitelists) {
        let _wihitelist = await instance.wihitelists(wihitelist.token_address);
        if (_wihitelist.token_address != wihitelist.token_address || _wihitelist.deposit_amount != wihitelist.deposit_amount || _wihitelist.status != wihitelist.status) {
            console.log("wihitelist: ", wihitelist);
            const tx = await instance.setWihitelist(wihitelist.token_address, wihitelist);
            const txReceipt = await tx.wait(1);
            console.log(`tx hash: ${txReceipt.hash}`)
            _wihitelist = await instance.wihitelists(wihitelist.token_address);
        }
        console.log("_wihitelist: ", _wihitelist);
    }
    console.log("4. setNftMapping success.");

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })