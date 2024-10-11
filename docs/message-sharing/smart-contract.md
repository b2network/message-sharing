# Smart Contract

The MessageSharing contract is a smart contract designed for cross-chain message transmission, providing a set of
functionalities to manage and verify cross-chain messages. Here's a summary of the contract's main features and
deployment steps:

<!-- TOC -->

* [Smart Contract](#smart-contract)
    * [Contracts interface](#contracts-interface)
    * [Contracts events](#contracts-events)
    * [Contracts code](#contracts-code)
    * [Deployment Steps](#deployment-steps)
    * [Contracts instances](#contracts-instances)

<!-- TOC -->

## Contracts interface

1. validatorRole: Returns the hash associated with the validator role for a specific chain.
2. SendHash: Generates a hash for a cross-chain message, used for subsequent verification and processing.
3. verify: Verifies the legitimacy of a message by checking if the signature is valid.
4. setWeight: Sets the weight for message processing, influencing the logic or priority of message verification.
5. call: Requests cross-chain message data and returns a unique message ID.
6. send: Confirms cross-chain message data, ensuring the message has not been processed before and verifying the weight
   of the signatures.
7. setValidatorRole: Sets the validator role for a specific chain.

```
interface IMessageSharing {

    /**
     * Get the validator role for a specific chain
     * @param chain_id The ID of the chain for which to retrieve the validator role.
     * @return bytes32 The hash associated with the validator role for the specified chain ID.
     */
    function validatorRole(uint256 chain_id) external pure returns (bytes32);

    /**
     * Generate a message hash
     * @param from_chain_id The ID of the originating chain, used to identify the source of the message.
     * @param from_id The ID of the cross-chain message, used to uniquely identify the message.
     * @param from_sender The address of the sender on the originating chain.
     * @param to_chain_id The ID of the target chain, where the message will be sent.
     * @param to_business_contract The address of the target contract that will receive the cross-chain message.
     * @param to_message The input data for the target contract's cross-chain call.
     * @return bytes32 The generated message hash, used for subsequent verification and processing.
     */
    function SendHash(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external view returns (bytes32);

    /**
     * Verify the legitimacy of a message
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed.
     * @param from_sender The address of the sender on the originating chain, used to verify the sender's legitimacy.
     * @param to_chain_id The ID of the target chain, indicating where the message will be sent.
     * @param to_business_contract The address of the target contract that will receive the cross-chain message.
     * @param to_message The input data for the target contract's cross-chain call.
     * @param signature The signature of the message, used to verify its legitimacy and integrity.
     * @return bool Returns true if the verification succeeds, and false if it fails.
     */
    function verify(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes calldata to_message, bytes calldata signature) external view returns (bool);

    /**
     * Set the weight for message processing
     * @param chain_id The ID of the chain.
     * @param _weight The weight value that influences the logic or priority of message processing.
     */
    function setWeight(uint256 chain_id, uint256 _weight) external;

    /**
     * Request cross-chain message data
     * @param to_chain_id The ID of the target chain, specifying where the message will be sent.
     * @param to_business_contract The address of the target contract that will receive the cross-chain message.
     * @param to_message The input data for the target contract's cross-chain call.
     * @return from_id The ID of the cross-chain message, returning a unique identifier to track the request.
     */
    function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256 from_id);

    /**
     * Confirm cross-chain message data
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed.
     * @param from_sender The address of the sender on the originating chain (msg.sender), used to determine the sender's security based on business needs.
     * @param to_business_contract The address of the target contract, indicating where the message will be sent (can be a contract on the target chain or the current chain).
     * @param to_message The input data for the target contract's cross-chain call.
     * @param signatures An array of signatures used to verify the legitimacy of the message, ensuring only authorized senders can send the message.
     */
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, address to_business_contract, bytes calldata to_message, bytes[] calldata signatures) external;

    /**
     * Set the validator role for a specific chain
     * @param chain_id The ID of the chain for which to set the validator role.
     * @param account The address of the validator, indicating which account to set the role for.
     * @param valid A boolean indicating the validity of the validator role, true for valid and false for invalid.
     */
    function setValidatorRole(uint256 chain_id, address account, bool valid) external;

   }
```

## Contracts events

1. SetWeight: Triggered when the weight is set.
2. SetValidatorRole: Triggered when a validator role is set.
3. Send: Triggered when a message is sent.
4. Call: Triggered when a message call is made.

```
    event SetWeight(uint256 chain_id, uint256 weight); // Event emitted when weight is set
    event SetValidatorRole(uint256 chain_id, address account, bool valid); // Event emitted when validator role is set
    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes to_message); // Event emitted when a message is sent
    event Call(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes to_message); // Event emitted when a message call is made

```

## Contracts code

[MessageSharing.sol](../../contracts/contracts/message_sharing/MessageSharing.sol)

## Deployment Steps

1. Navigate to the Working Directory:

```
$ cd message-sharing/contracts 
```

2. Install Dependencies:

```
$ npm install
```

3. Set Environment Variables:

If the .env file does not exist, copy .env.test to .env and set the chain RPC URLs and deployment account private keys.

```
$ cp .env.test .env
```

[.env](../../contracts/.env.test)

```
// bsquared-dev
AS_DEV_RPC_URL=https://arbitrum-sepolia.blockpi.network/v1/rpc/public
AS_DEV_PRIVATE_KEY_0=...
// arbitrum-sepolia
B2_DEV_RPC_URL=https://b2-testnet.alt.technology
B2_DEV_PRIVATE_KEY_0=...
```  

4. Deploy the Contract:

Use the yarn command to deploy the contract on different networks.

```
// Deploy on the bsquared-dev network
$ yarn hardhat run scripts/message_sharing/deploy.js --network b2dev
// Deploy on the arbitrum-sepolia network
$ yarn hardhat run scripts/message_sharing/deploy.js --network asdev
```

5. Set Contract Address Environment Variables:

Add the contract addresses to the .env file:

```
B2_DEV_MESSAGE_SHARING=0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8
AS_DEV_MESSAGE_SHARING=0x2A82058E46151E337Baba56620133FC39BD5B71F
```

6. Configure the Contract:

Modify the [scripts/message_sharing/set.js](../../contracts/scripts/message_sharing/set.js) script to set validators and
weight.
validators is the set of witnesses signing the message, and weight is the minimum weight required for successful
signature verification when sending transactions on-chain.

```
    let weight;
    let validators;
    if (network.name == 'b2dev') {
        address = process.env.B2_DEV_MESSAGE_SHARING;
        weight = {
            chain_id: process.env.B2_DEV_CHAIN_ID, weight: 1,
        };
        validators = [{
            chain_id: process.env.B2_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }, {
            chain_id: process.env.AS_DEV_CHAIN_ID, account: '0x8F8676b34cbEEe7ADc31D17a149B07E3474bC98d', valid: true,
        }];
    }
```

By following these steps, the MessageSharing contract can be successfully deployed and configured for secure and
efficient cross-chain message transmission.

```
// Set on the bsquared-dev network
$ yarn hardhat run scripts/message_sharing/set.js --network b2dev
// Set on the arbitrum-sepolia network
$ yarn hardhat run scripts/message_sharing/set.js --network asdev
```

## Contracts instances

```
MessageSharing(Bsquared-dev): 0xDf5b12f094cf9b12eb2523cC43a62Dd6787D7AB8
// 
MessageSharing(Arbitrum-sepolia): 0x2A82058E46151E337Baba56620133FC39BD5B71F
```