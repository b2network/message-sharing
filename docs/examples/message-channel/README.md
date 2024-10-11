# Message channel

What you need to know before looking at the examples, See the contents
of [docs/examples/attention.md](./docs/examples/attention.md)

## Contract

### Contract Interfaces

```
/**
 * @dev Initiates a cross-chain message call to another blockchain network.
 * @param to_chain_id The ID of the destination chain where the message will be sent.
 * @param to_business_contract The address of the target contract on the destination chain.
 * @param to_message The payload of the message to be sent.
 * @return from_id The unique identifier for this message call, returned by the message sharing contract.
 *
 * This function uses the messageSharing contract to send a message to a specified contract
 * on another blockchain. It emits a Call event for logging purposes.
 */
function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256);

/**
 * @dev Receives and processes an incoming cross-chain message.
 * @param from_chain_id The ID of the source chain from which the message originates.
 * @param from_id The unique identifier of the message on the source chain.
 * @param from_sender The address of the sender on the source chain.
 * @param to_message The payload of the message received.
 * @return success Boolean indicating whether the message was processed successfully.
 *
 * This function is restricted to entities with the SENDER_ROLE, ensuring only authorized
 * senders can invoke it. It emits a Send event for logging purposes. The function currently
 * contains placeholders for verifying the message's validity and executing the service logic.
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata to_message) external onlyRole(SENDER_ROLE) override returns (bool success);

/**
 * @dev Sets the address of the message sharing contract.
 * @param sharing_address The address of the IMessageSharing contract to be set.
 *
 * This function allows an admin to specify the address of the message sharing contract
 * that will be used for cross-chain message communication. It is restricted to accounts
 * with the ADMIN_ROLE to ensure only authorized modifications.
 */
function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE);
```

### Contract Code

See the contents
of [contracts/examples/message_channel/MessageChannel.sol](../../../contracts/contracts/examples/message_channel/MessageChannel.sol)

### Deployment

The deployment and configuration instructions you've provided outline the steps for deploying and setting up the
MessageChannel contract on two different blockchain networks: bsquared-dev and arbitrum-sepolia. Here's a detailed
explanation of each step:

1. Deploy the MessageChannel Contract

```
// Deploy the MessageChannel contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/message_channel/deploy.js --network b2dev
// Deploy the MessageChannel contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/message_channel/deploy.js --network asdev
```

Explanation:

These commands use Hardhat, a development environment for Ethereum software, to deploy the MessageChannel contract.
The [deploy.js](../../../contracts/scripts/examples/message_channel/deploy.js) script is executed for each specified
network (b2dev and asdev), which are likely configurations defined in your Hardhat setup.
The deployment script should contain logic to compile the contract and deploy it to the specified network.

2. Configure Contract Information

After deploying the contract, update your .env file with the deployed contract addresses.

Example entry:

```
B2_DEV_MESSAGE_CHANNEL=0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8
AS_DEV_MESSAGE_CHANNEL=0x2A82058E46151E337Baba56620133FC39BD5B71F
```

3. Configure the MessageChannel Contract

```
// Configure the MessageChannel contract on the bsquared-dev chain
$ yarn hardhat run scripts/examples/message_channel/set.js --network b2dev
// Configure the MessageChannel contract on the arbitrum-sepolia chain
$ yarn hardhat run scripts/examples/message_channel/set.js --network asdev
```

Explanation:

These commands run the [set.js](../../../contracts/scripts/examples/message_channel/set.js) script on each network to
configure the deployed MessageChannel contract.
The script should contain logic to set essential contract parameters, such as the address of the IMessageSharing
contract or any other initial configurations required.
Ensure that the set.js script uses the addresses from the .env file to interact with the correct contract instances.

4. Test

```
$ yarn hardhat run scripts/examples/message_channel/test.js --network b2dev
```

Explanation:

This command executes the [test.js](../../../contracts/scripts/examples/message_channel/test.js) script on the
bsquared-dev network to verify the deployment and configuration.
The script should contain tests to ensure that the MessageChannel contract is functioning correctly, such as sending and
receiving messages.