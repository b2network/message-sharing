# Orderbook

What you need to know before looking at the examples, See the contents
of [docs/examples/attention.md](./docs/examples/attention.md)

## Contract

### Contract Interfaces

#### Cashier Contract Interfaces

```
/**
 * @notice Allows a user to pay for an order using a supported token
 * @param order_no The unique identifier for the order
 * @param token_address The address of the token to be used for payment
 */
 function payOrder(string calldata order_no, address token_address) external payable;
 
 /**
 * @notice Processes a cross-chain message to settle an order
 * @param from_chain_id The ID of the source chain
 * @param from_id The unique identifier of the message
 * @param from_sender The address of the sender of the message
 * @param data The encoded message containing settlement data
 * @return success Returns whether the operation was successful
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external onlyRole(SENDER_ROLE) returns (bool success);

/**
 * @notice Sets the whitelist configuration for a token
 * @param token_address The address of the token to be whitelisted
 * @param wihitelist The whitelist configuration for the token
 */
function setWihitelist(address token_address, Wihitelist calldata wihitelist) external onlyRole(ADMIN_ROLE);

/**
 * @notice Sets the orderbook configuration
 * @param contract_address The address of the orderbook contract
 * @param chain_id The ID of the chain where the orderbook contract is located
 */
function setOrderbook(address contract_address, uint256 chain_id) external onlyRole(ADMIN_ROLE);

/**
 * @notice Sets the cross-chain message sharing contract address
 * @param sharing_address The new address of the cross-chain message sharing contract
 */
function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE);

/**
 * @notice Withdraws the specified amount of tokens to a given address
 * @param token_address The address of the token to be withdrawn
 * @param to_address The address to which the tokens will be sent
 * @param amount The amount of tokens to be withdrawn
 */
function withdraw(address token_address, address to_address, uint256 amount) external onlyRole(ADMIN_ROLE);
```

#### Orderbook Contract Interfaces

```
/**
 * @notice Processes a cross-chain message to create an order
 * @param from_chain_id The ID of the source chain
 * @param from_id The unique identifier of the message
 * @param from_sender The address of the sender of the message
 * @param message The encoded message containing order data
 * @return success Returns whether the operation was successful
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success);

/**
 * @notice Settles an order by finalizing it and sending a message back to the originating chain
 * @param order_no The unique identifier for the order to be settled
 * @param fee_amount The fee amount to be deducted from the order's deposit
 */
function settle(string calldata order_no, uint256 fee_amount) external;

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
```

### Contract Code

See the contents
of [contracts/examples/orderbook/Cashier.sol](../../../contracts/contracts/examples/orderbook/Cashier.sol)
and [contracts/examples/orderbook/Orderbook.sol](../../../contracts/contracts/examples/orderbook/Orderbook.sol)

### Deployment

The deployment and configuration process for the Cashier and Orderbook contracts involves several steps. Here's a
detailed explanation of each:

1. Deploy the Cashier & Orderbook Contracts

1.1. Deploy the Cashier Contract

Deploy the Cashier contract on two different blockchain networks using the following commands:

```
// Deploy the Cashier contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network b2dev
// Deploy the Cashier contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/orderbook/cashier/deploy.js --network asdev
```

1.2. Deploy the Orderbook Contract

Similarly, deploy the Orderbook contract on the same networks:

```
// Deploy the Orderbook contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network b2dev
// Deploy the Orderbook contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/orderbook/orderbook/deploy.js --network asdev
```

2. Configure Contract Information

After deploying the contract, update your .env file with the deployed contract addresses.

Example entry:

```
B2_DEV_ORDERBOOK=0x91009D6edEaBfE095749DaBf2c0359c5f4343e8e
B2_DEV_CASHIER=0xfb70B98f4935f6E983D1Ccac947dA51d6d42c023
AS_DEV_ORDERBOOK=0xYourArbitrumSepoliaOrderbookAddressHere
AS_DEV_CASHIER=0xYourArbitrumSepoliaCashierAddressHere
```

3. Configure the Cashier & Orderbook Contract

3.1. Configure the Cashier Contract

Run the configuration scripts to set up the Cashier contracts on each network:

```
// Configure the Cashier contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/orderbook/cashier/set.js --network b2dev
// Configure the Cashier contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/orderbook/cashier/set.js --network asdev
```

3.2. Configure the Orderbook Contract

Similarly, configure the Orderbook contracts:

```
// Configure the Orderbook contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network b2dev
// Configure the Orderbook contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/orderbook/orderbook/set.js --network asdev
```

4. Test

4.1. Test payOrder

Run
the [contracts/scripts/examples/orderbook/cashier/test.js](../../../contracts/scripts/examples/orderbook/cashier/test.js)
script for the Cashier contract to ensure the payOrder functionality works as expected:

```
$ yarn hardhat run scripts/examples/orderbook/cashier/test.js --network b2dev
```

4.2. Test settle

Run
the [contracts/scripts/examples/orderbook/orderbook/test.js](../../../contracts/scripts/examples/orderbook/orderbook/test.js)
script for the Orderbook contract to verify the settle functionality:

```
$ yarn hardhat run scripts/examples/orderbook/orderbook/test.js --network asdev
```

4.3. Test withdraw

Finally, test the [withdraw](../../../contracts/scripts/examples/orderbook/cashier/withdraw.js) functionality of the
Cashier contract:

```
$ yarn hardhat run scripts/examples/orderbook/cashier/withdraw.js --network b2dev
```

Note:
The token used for testing needs to be deployed using
the [deploy.js](../../../contracts/scripts/examples/erc20/deploy.js)
script, and operations such as [minting](../../../contracts/scripts/examples/erc20/mint.js) and
[approving](../../../contracts/scripts/examples/erc20/approve.js) should be performed.
