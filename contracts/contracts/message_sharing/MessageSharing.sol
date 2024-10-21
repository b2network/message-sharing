// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

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

    // Event declarations
    event SetWeight(uint256 chain_id, uint256 weight); // Event emitted when weight is set
    event SetValidatorRole(uint256 chain_id, address account, bool valid); // Event emitted when validator role is set
    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes to_message); // Event emitted when a message is sent
    event Call(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes to_message); // Event emitted when a message call is made
}

interface IBusinessContract {
    /**
     * Process cross-chain information in the business contract
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed to prevent duplication.
     * @param from_sender The address of the sender on the originating chain, used to verify the sender's legitimacy (business needs may dictate whether verification is necessary).
     * @param message The input data for processing the cross-chain message, which may need to be decoded based on byte encoding rules.
     * @return success Indicates whether the message processing was successful, returning true for success and false for failure.
     * @return error Returns an error message if processing fails. A descriptive string containing the reason for failure, useful for debugging and logging.
     */
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external returns (bool success, string memory error);
}

contract MessageSharing is IMessageSharing, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    using ECDSA for bytes32;
    bytes32 public constant SEND_HASH_TYPE = keccak256('Send(uint256 from_chain_id,uint256 from_id,address from_sender,uint256 to_chain_id,address to_business_contract,bytes to_message)');
    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");

    mapping (uint256 => uint256) public sequences;
    mapping (uint256 => mapping (uint256 => bool)) public ids;
    mapping (uint256 => uint256) public weights;

    function initialize() public initializer {
        __EIP712_init("B2MessageSharing", "1");
        __AccessControl_init();
        __UUPSUpgradeable_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    {

    }

    function validatorRole(uint256 chain_id) public pure override returns (bytes32) {
      return keccak256(abi.encode("validator_role", chain_id));
    }

    function setWeight(uint256 chainId, uint256 _weight) external override onlyRole(ADMIN_ROLE) {
        weights[chainId] = _weight;
        emit SetWeight(chainId, _weight);
    }

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, address to_business_contract, bytes calldata to_message, bytes[] calldata signatures) external override {
        require(!ids[from_chain_id][from_id], "non-repeatable processing");
        require(weights[from_chain_id] > 0, "weight not set");
        uint256 weight_ = 0;
        for(uint256 i = 0; i < signatures.length; i++) {
           bool success = verify(from_chain_id, from_id, from_sender, block.chainid, to_business_contract, to_message, signatures[i]);
           if (success) {
                weight_ = weight_ + 1;
           }
        }
        require(weight_ >= weights[from_chain_id], "verify signatures weight invalid");

        if (to_business_contract != address(0x0)) {
            (bool success, string memory error) = IBusinessContract(to_business_contract).send(from_chain_id, from_id, from_sender, to_message);
            require(success, error);
        }
        emit Send(from_chain_id, from_id, from_sender, block.chainid, to_business_contract, to_message);
    }

    function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external override returns (uint256) {
        sequences[to_chain_id]++;
        emit Call(block.chainid, sequences[to_chain_id], msg.sender, to_chain_id, to_business_contract, to_message);
        return sequences[to_chain_id];
    }

    function verify(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes calldata to_message, bytes calldata signature) public view override returns (bool) {
        bytes32 digest  = SendHash(from_chain_id, from_id, from_sender, to_chain_id, to_business_contract, to_message);
        return hasRole(validatorRole(from_chain_id), digest.recover(signature));
    }

    function SendHash(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address to_business_contract, bytes calldata to_message) public view override returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(SEND_HASH_TYPE,from_chain_id, from_id, from_sender, to_chain_id, to_business_contract, keccak256(to_message))));
    }

    function setValidatorRole(uint256 chain_id, address account, bool valid) external override onlyRole(ADMIN_ROLE) {
        if (valid) {
            _grantRole(validatorRole(chain_id), account);
        } else {
             _revokeRole(validatorRole(chain_id), account);
        }
        emit SetValidatorRole(chain_id, account, valid);
    }

}