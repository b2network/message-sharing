// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";

interface IMessageSharing {
    function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256 from_id);
}

interface IBusinessContract {
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata to_message) external returns (bool success);
}

contract MessageChannel is IBusinessContract, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

    IMessageSharing public messageSharing;

    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes to_message);
    event Call(uint256 from_id, uint256 to_chain_id, address to_business_contract, bytes to_message);

    function initialize() public initializer {
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

    function version() external pure returns (string memory) {
        return "v1.0.0";
    }

    function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE) {
        messageSharing = IMessageSharing(sharing_address);
    }

    function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256) {
        uint256 from_id =  messageSharing.call(to_chain_id, to_business_contract, to_message);
        emit Call(from_id, to_chain_id, to_business_contract, to_message);
        return from_id;
    }

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata to_message) external onlyRole(SENDER_ROLE) override returns (bool success) {
        // TODO 1. Verify the validity of from_chain_id and from_sender
        // TODO 2. Verify that from_id has been executed
        // TODO 3. Parse data and execute service logic
        emit Send(from_chain_id, from_id, from_sender, to_message);
        return true;
    }

}