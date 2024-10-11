### Matters needing attention

1. Use the call method of the MessageSharing contract to initiate a cross-chain request:

```
/**
 * @dev Initiates a cross-chain message call.
 * @param to_chain_id The ID of the destination chain.
 * @param to_business_contract The address of the target contract on the destination chain.
 * @param to_message The message payload to be sent.
 * @return from_id The unique identifier for this message call.
 */
function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256)
```

2. Implement the send method to receive and process data from the MessageSharing contract:

```
/**
 * @dev Receives and processes a cross-chain message.
 * @param from_chain_id The ID of the source chain.
 * @param from_id The unique identifier of the message on the source chain.
 * @param from_sender The address of the sender on the source chain.
 * @param to_message The received message payload.
 * @return success Boolean indicating whether the message was processed successfully.
 */
function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata to_message) external onlyRole(SENDER_ROLE) override returns (bool success)
```

When implementing cross-chain message transmission, ensure to correctly invoke the call function to initiate requests,
and implement the send function in the target contract to handle received messages. The onlyRole(SENDER_ROLE) modifier
in the send function is not mandatory and can be used based on specific requirements to control access permissions.
