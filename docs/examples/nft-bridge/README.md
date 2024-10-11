# Nft bridge

What you need to know before looking at the examples, See the contents
of [docs/examples/attention.md](./docs/examples/attention.md)

## Contract

### Contract Interfaces

```
/**
 * @notice Receives a cross-chain message and unlocks the NFT
 * @param from_chain_id The ID of the source chain
 * @param from_id The unique identifier of the message
 * @param from_sender The address of the sender of the message
 * @param message The encoded message containing unlock data
 * @return success Returns whether the operation was successful
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success);

/**
 * @notice Locks an NFT and sends a cross-chain message
 * @param nft_address The address of the NFT to be locked
 * @param token_id The ID of the token to be locked
 * @param amount The amount of the NFT to lock (for ERC1155)
 * @param to_chain_id The ID of the destination chain
 * @param to_token_bridge The NFT bridge contract address on the destination chain
 * @param to_address The address of the recipient
 */
function lock(address nft_address, uint256 token_id, uint256 amount, uint256 to_chain_id, address to_token_bridge, address to_address) external;
/**
 * @notice Sets the cross-chain message sharing contract address
 * @param sharing_address The new address of the cross-chain message sharing contract
 */
function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE);

/**
 * @notice Sets the NFT mapping between source and destination chains
 * @param from_chain_id The ID of the source chain
 * @param from_nft_address The NFT address on the source chain
 * @param to_chain_id The ID of the destination chain
 * @param to_nft_address The NFT address on the destination chain
 * @param standard The NFT standard (ERC721 or ERC1155)
 * @param status The status of the NFT mapping (active or inactive)
 */
function setNftMapping(uint256 from_chain_id, address from_nft_address, uint256 to_chain_id, address to_nft_address, NftStandard standard, bool status) external onlyRole(ADMIN_ROLE);
```

### Contract Code

See the contents
of [contracts/examples/nft_bridge/NftBridge.sol](../../../contracts/contracts/examples/nft_bridge/NftBridge.sol)

### Deployment

The deployment and configuration process for the NftBridge contract involves several key steps. Here is a detailed
explanation of each step:

1. Deploy the NftBridge Contract

This step involves deploying the NftBridge contract to two different blockchain networks. Use Hardhat to execute the
deployment scripts as follows:

```
// Deploy the NftBridge contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/nft_bridge/deploy.js --network b2dev
// Deploy the NftBridge contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/nft_bridge/deploy.js --network asdev
```

Explanation:
The [contracts/scripts/examples/nft_bridge/deploy.js](../../../contracts/scripts/examples/nft_bridge/deploy.js) is the
path to the deployment script.

2. Configure Contract Information

After deploying the contract, update your .env file with the deployed contract addresses.

Example entry:

```
B2_DEV_NFT_BRIDGE=0x952b63C6C799B7033c24B055f7F023Eb7f3a5c73
AS_DEV_NFT_BRIDGE=0xYourArbitrumSepoliaAddressHere
```

3. Configure the NftBridge Contract

This step involves running configuration scripts to set up the NftBridge contracts on each network. These scripts
might set initial parameters, permissions, or other necessary configurations.

```
// Configure the NftBridge contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/nft_bridge/set.js --network b2dev
// Configure the NftBridge contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/nft_bridge/set.js --network asdev
```

Explanation:
The [contracts/scripts/examples/nft_bridge/set.js](../../../contracts/scripts/examples/nft_bridge/set.js) script likely
contains logic to initialize or configure the contract after deployment.

4. Test

Finally, you run tests to ensure that the NftBridge contract works as expected on the specified network.

```
$ yarn hardhat run scripts/examples/nft_bridge/test.js --network b2dev
```

Explanation:
This command runs the [test.js](../../../contracts/scripts/examples/nft_bridge/test.js) script, which should contain
test cases or scenarios to validate the functionality of your deployed contract.

Note:
The nft used for testing needs to be deployed using
the [deploy.js](../../../contracts/scripts/examples/erc721/deploy.js)
script, and operations such as [minting](../../../contracts/scripts/examples/erc721/mint.js) and
[setApprovalForAll](../../../contracts/scripts/examples/erc721/setApprovalForAll.js) should be performed.