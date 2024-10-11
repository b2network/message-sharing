const {ethers, run, network, upgrades} = require("hardhat")

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/examples/nft_bridge/grant_role.js --network b2dev
     * asdev: yarn hardhat run scripts/examples/nft_bridge/grant_role.js --network asdev
     * # pord
     * as: yarn hardhat run scripts/examples/nft_bridge/grant_role.js --network as
     * b2: yarn hardhat run scripts/examples/nft_bridge/grant_role.js --network b2
     */

    let address;
    let account;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_NFT_BRIDGE;
        account = "";
    } else if (network.name == 'asdev') {
        address = process.env.AS_DEV_NFT_BRIDGE;
        account = "";
    } else if (network.name == 'as') {
        address = process.env.AS_NFT_BRIDGE;
        account = "";
    } else if (network.name == 'b2') {
        address = process.env.B2_NFT_BRIDGE;
        account = "";
    }
    console.log("NftBridge Address: ", address);

    const [owner] = await ethers.getSigners();
    console.log("Owner Address:", owner.address);

    const NftBridge = await ethers.getContractFactory("NftBridge");
    const instance = await NftBridge.attach(address)

    const role = await instance.ADMIN_ROLE();
    // const role = await instance.UPGRADE_ROLE();
    // const role = await instance.SENDER_ROLE();
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