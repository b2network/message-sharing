const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/erc20/grant_role.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/erc20/grant_role.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/erc20/grant_role.js --network as
     * b2: yarn hardhat run scripts/examples/erc20/grant_role.js --network b2
     */

    let address;
    let account;
    if (network.name == 'b2dev') {
        address = "0xE6BF3CCAb0D6b461B281F04349aD73d839c25B06";
        account = "";
    } else if (network.name == 'asdev') {
        address = "";
        account = "";
    } else if (network.name == 'as') {
        address = "";
        account = "";
    } else if (network.name == 'b2') {
        address = "";
        account = "";
    }
    console.log("MyERC20 Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const MyERC20 = await ethers.getContractFactory("MyERC20");
    const instance = await MyERC20.attach(address)

    const role = await instance.MINTER_ROLE();
    // const role = await instance.BURNER_ROLE();
    console.log("role hash:", role);

    let has = await instance.hasRole(role, account)
    console.log("has role:", has)
    if (!has) {
        const tx = await instance.grantRole(role, account);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        has = await instance.hasRole(role, account)
        console.log("has role:", has)
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })