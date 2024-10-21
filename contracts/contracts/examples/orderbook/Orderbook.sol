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
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success, string memory error);
}

enum OrderStatus {
    UNKNOWN,
    IN_PROCESS,
    COMPLETED
}

struct Order {
    uint256 from_chain_id;
    address from_business;
    address user_address;
    uint256 start_time;
    uint256 end_time;
    address token_address;
    uint256 deposit_amount;
    uint256 fee_amount;
    OrderStatus status;
}

contract Orderbook is IBusinessContract, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

    // message sharing address
    IMessageSharing public message_sharing;
    // from_chain_id => nft bridge address
    mapping (uint256 => address) public bridges;
    // from_chain_id => from_id => execute
    mapping (uint256 => mapping (uint256 => bool)) public executes;

    mapping (string => Order) public orders;

    event PayOrder(uint256 from_chain_id, uint256 from_id, address from_business, string order_no, uint256 start_time, address user_address, address token_address, uint256 deposit_amount);

    event Settle(string order_no, uint256 fee_amount, uint256 end_time, uint256 from_id);

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success, string memory error) {
        require(bridges[from_chain_id] == from_sender, "Invalid chain id or from_sender");
        require(!executes[from_chain_id][from_id], "Have been executed");
        executes[from_chain_id][from_id] = true;

        (string memory order_no, address token_address, address user_address, uint256 deposit_amount, uint256 start_time) = decodePayData(message);
        Order storage order = orders[order_no];
        require(order.status == OrderStatus.UNKNOWN, "Order already exists");
        order.from_chain_id = from_chain_id;
        order.from_business = from_sender;
        order.user_address = user_address;
        order.start_time = start_time;
        order.end_time = 0;
        order.token_address = token_address;
        order.deposit_amount = deposit_amount;
        order.fee_amount = 0;
        order.status = OrderStatus.IN_PROCESS;
        emit PayOrder(from_chain_id, from_id, from_sender, order_no, start_time, user_address, token_address, deposit_amount);
        return (true, "");
    }

    function settle(string calldata order_no, uint256 fee_amount) external {
        Order storage order = orders[order_no];
        require(order.status == OrderStatus.IN_PROCESS, "Invalid order status");
        require(order.deposit_amount >= fee_amount, "Invalid fee amount");
        order.end_time = block.timestamp;
        order.status = OrderStatus.COMPLETED;

        bytes memory to_message = encodeSettleData(order_no, order.token_address, order.user_address, fee_amount, order.end_time);

        uint256 from_id =  message_sharing.call(order.from_chain_id, order.from_business, to_message);
        emit Settle(order_no, fee_amount, block.timestamp, from_id);
    }

    function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE) {
        message_sharing = IMessageSharing(sharing_address);
    }

    function setBridges(uint256 from_chain_id, address bridge) external onlyRole(ADMIN_ROLE) {
        bridges[from_chain_id] = bridge;
    }

    function decodePayData(bytes memory data) public pure returns (string memory order_no, address token_address, address user_address, uint256 deposit_amount, uint256 start_time) {
        (order_no, token_address, user_address, deposit_amount, start_time) = abi.decode(data, (string, address, address, uint256, uint256));
    }

    function encodeSettleData(string memory order_no, address token_address, address user_address, uint256 fee_amount, uint256 end_time) public pure returns (bytes memory) {
        return abi.encode(order_no, token_address, user_address, fee_amount, end_time);
    }

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

}