# Token bridge

What you need to know before looking at the examples, See the contents
of [docs/examples/attention.md](./docs/examples/attention.md)

## Contract

### Contract Interfaces

```
/**
 * @notice Receives cross-chain messages and unlocks tokens
 * @param from_chain_id The ID of the source chain
 * @param from_id The unique identifier of the message
 * @param from_sender The address of the sender of the message
 * @param message The encoded message containing unlock data
 * @return success Returns whether the operation was successful
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success);

/**
 * @notice Locks tokens and sends a cross-chain message
 * @param token_address The address of the token to be locked
 * @param amount The amount of tokens to lock
 * @param to_chain_id The ID of the destination chain
 * @param to_token_bridge The token bridge contract address on the destination chain
 * @param to_address The address of the recipient
 */
function lock(address token_address, uint256 amount, uint256 to_chain_id, address to_token_bridge, address to_address) external payable;

/**
 * @notice Sets the cross-chain message sharing contract address
 * @param sharing_address The new address of the cross-chain message sharing contract
 */
function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE);

/**
 * @notice Sets the bridge contract address for a specific source chain
 * @param from_chain_id The ID of the source chain
 * @param bridge The new bridge contract address
 */
function setBridges(uint256 from_chain_id, address bridge) external onlyRole(ADMIN_ROLE);

/**
 * @notice Sets the token mapping between source and destination chains
 * @param from_chain_id The ID of the source chain
 * @param from_token_address The token address on the source chain
 * @param to_token_address The token address on the destination chain
 */
function setTokenMapping(uint256 from_chain_id, address from_token_address, address to_token_address) external onlyRole(ADMIN_ROLE);
```

### Contract Code

See the contents
of [contracts/examples/token_bridge/TokenBridge.sol](../../../contracts/contracts/examples/token_bridge/TokenBridge.sol)

### Deployment

The deployment and configuration process you've outlined for the TokenBridge contract involves several key steps. Here's
a breakdown of each step with some additional context and tips:

1. Deploy the TokenBridge Contract

This step involves deploying the TokenBridge contract to two different blockchain networks. The commands use Hardhat,
a popular Ethereum development environment, to execute deployment scripts.

```
// Deploy the TokenBridge contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/token_bridge/deploy.js --network b2dev
// Deploy the TokenBridge contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/token_bridge/deploy.js --network asdev
```

Explanation:
The [contracts/scripts/examples/token_bridge/deploy.js](../../../contracts/scripts/examples/token_bridge/deploy.js) is
the path to the deployment script.

2. Configure Contract Information

After deploying the contract, update your .env file with the deployed contract addresses.

Example entry:

```
B2_DEV_TOKEN_BRIDGE=0xE2E4C6B693c66f4C9C3c1c88A5Ef76d94526d6fB
AS_DEV_TOKEN_BRIDGE=0xYourArbitrumSepoliaAddressHere
```

3. Configure the TokenBridge Contract

This step involves running configuration scripts to set up the TokenBridge contracts on each network. These scripts
might set initial parameters, permissions, or other necessary configurations.

```
// Configure the TokenBridge contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/token_bridge/set.js --network b2dev
// Configure the TokenBridge contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/token_bridge/set.js --network asdev
```

Explanation:
The [contracts/scripts/examples/token_bridge/set.js](../../../contracts/scripts/examples/token_bridge/set.js) script
likely contains logic to initialize or configure the contract after deployment.

4. Test

Finally, you run tests to ensure that the TokenBridge contract works as expected on the specified network.

```
$ yarn hardhat run scripts/examples/token_bridge/test.js --network b2dev
```

Explanation:
This command runs
the [contracts/scripts/examples/token_bridge/test.js](../../../contracts/scripts/examples/token_bridge/test.js) script,
which should contain
test cases or scenarios to validate the functionality of your deployed contract.

Note:
The token used for testing needs to be deployed using
the [deploy.js](../../../contracts/scripts/examples/erc20/deploy.js)
script, and operations such as [minting](../../../contracts/scripts/examples/erc20/mint.js) and
[approving](../../../contracts/scripts/examples/erc20/approve.js) should be performed.