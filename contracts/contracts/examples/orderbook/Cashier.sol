// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

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

struct Wihitelist {
    address token_address;
    uint256 deposit_amount;
    bool status;
}

struct Orderbook {
    uint256 chain_id;
    address contract_address;
}

struct Order {
    address user_address;
    uint256 start_time;
    uint256 end_time;
    address token_address;
    uint256 deposit_amount;
    uint256 fee_amount;
    OrderStatus status;
}

contract Cashier is IBusinessContract, Initializable, UUPSUpgradeable, AccessControlUpgradeable {
    using SafeERC20 for IERC20;
    using Address for address;

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

    IMessageSharing public message_sharing;
    Orderbook public orderbook;
    // token_address => wihitelist
    mapping (address => Wihitelist) public wihitelists;
    // from_chain_id => from_id => execute
    mapping (uint256 => mapping (uint256 => bool)) public executes;
    // order_no => order
    mapping (string => Order) public orders;
    // user_address => order_no
    mapping (address => string[]) public user_orders;
    // withdraw balance
    uint256 public withdraw_balance;

    event SetWihitelist(address indexed token_address, Wihitelist wihitelist);
    event PayOrder(string indexed order_no, address indexed token_address, address indexed user_address, uint256 deposit_amount, uint256 from_id);
    event SettleOrder(string indexed order_no, address indexed token_address, address indexed user_address, uint256 fee_amount, uint256 end_time);
    event Withdraw(address token_address, address to_address, uint256 amount);

    // ************************************** PUBLIC FUNCTION **************************************

    function payOrder(string calldata order_no, address token_address) external payable {
        Wihitelist memory wihitelist = wihitelists[token_address];
        require(wihitelist.status, "Token not support");
        Order storage order = orders[order_no];
        require(order.status == OrderStatus.UNKNOWN, "Order already exists");

        if (token_address == address(0x0)) {
            require(wihitelist.deposit_amount == msg.value, "Invalid value");
        } else {
            require(msg.value == 0, "Invalid transaction value");
            IERC20(token_address).safeTransferFrom(msg.sender, address(this), wihitelist.deposit_amount);
        }

        order.user_address = msg.sender;
        order.start_time = block.timestamp;
        order.end_time = 0;
        order.token_address = token_address;
        order.deposit_amount = wihitelist.deposit_amount;
        order.fee_amount = 0;
        order.status = OrderStatus.IN_PROCESS;

        user_orders[msg.sender].push(order_no);

        bytes memory to_message = encodePayData(order_no, token_address, msg.sender, wihitelist.deposit_amount, block.timestamp);

        uint256 from_id = message_sharing.call(orderbook.chain_id, orderbook.contract_address, to_message);

        emit PayOrder(order_no, token_address, msg.sender, wihitelist.deposit_amount, from_id);
    }


    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external onlyRole(SENDER_ROLE) returns (bool success, string memory error) {
        require(!executes[from_chain_id][from_id], "Have been executed");
        executes[from_chain_id][from_id] = true;
        (string memory order_no, address token_address, address user_address, uint256 fee_amount, uint256 end_time) = decodeSettleData(data);

        require( orderbook.chain_id == from_chain_id, "Invalid chian id");
        require( orderbook.contract_address == from_sender, "Invalid business address");

        Order storage order = orders[order_no];
        require(order.status == OrderStatus.IN_PROCESS, "Invalid order status");
        require(order.user_address == user_address, "Invalid user address");
        require(order.deposit_amount >= fee_amount, "Invalid fee amount");

        for (uint256 i = 0; i < user_orders[user_address].length; i++) {
            if (compareStrings(user_orders[user_address][i], order_no)) {
                user_orders[user_address][i] = user_orders[user_address][user_orders[user_address].length - 1];
                user_orders[user_address].pop();
            }
        }

        order.fee_amount = fee_amount;
        order.end_time = end_time;
        order.status = OrderStatus.COMPLETED;
        if (order.deposit_amount > fee_amount) {
            _safeTransfer(token_address, user_address, order.deposit_amount - fee_amount);
        }
        withdraw_balance = withdraw_balance + fee_amount;

        emit SettleOrder(order_no, token_address, user_address, fee_amount, end_time);
        return (true, "");
    }

    // ************************************** ADMIN FUNCTION **************************************

    function setWihitelist(address token_address, Wihitelist calldata wihitelist) external onlyRole(ADMIN_ROLE) {
        wihitelists[token_address] = wihitelist;
        emit SetWihitelist(token_address, wihitelist);
    }

    function setOrderbook(address contract_address, uint256 chain_id) external onlyRole(ADMIN_ROLE) {
        orderbook.chain_id = chain_id;
        orderbook.contract_address = contract_address;
    }

    function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE) {
        message_sharing = IMessageSharing(sharing_address);
    }

    function withdraw(address token_address, address to_address, uint256 amount) external onlyRole(ADMIN_ROLE) {
        require(withdraw_balance >= amount, "Invalid amount");
        _safeTransfer(token_address, to_address, amount);
        withdraw_balance = withdraw_balance - amount;
        emit Withdraw(token_address, to_address, amount);
    }

    function encodePayData(string calldata order_no, address token_address, address user_address, uint256 deposit_amount, uint256 start_time) public pure returns (bytes memory) {
        return abi.encode(order_no, token_address, user_address, deposit_amount, start_time);
    }

    function decodeSettleData(bytes memory data) public pure returns (string memory order_no, address token_address, address user_address, uint256 fee_amount, uint256 end_time) {
        // bytes memory to_message = encodeSettleData(order_no, order.token_address, order.end_time, order.user_address, fee_amount);
        (order_no, token_address, end_time, user_address, fee_amount) = abi.decode(data, (string, address, uint256, address, uint256));
       // (order_no, token_address, user_address, fee_amount, end_time) = abi.decode(data, (string, address, address, uint256, uint256));
    }

    function compareStrings(string memory str1, string memory str2) internal pure returns (bool) {
        return keccak256(abi.encodePacked(str1)) == keccak256(abi.encodePacked(str2));
    }

    /**
     * @notice Safe transfer function
     *
     * @param to_address Token address
     * @param to_address        Address to get transferred BTC
     * @param amount    Amount of BTC to be transferred
     */
    function _safeTransfer(address token_address, address to_address, uint256 amount) internal {
        if (token_address == address(0x0)) {
            (bool success, bytes memory data) = address(to_address).call{
                value: amount
            }("");

            require(success, "transfer call failed");
            if (data.length > 0) {
                require(
                    abi.decode(data, (bool)),
                    "transfer operation did not succeed"
                );
            }
        } else {
            bool success = IERC20(token_address).transfer(to_address, amount);
            require(success, "token transfer call failed");
        }
    }

    function initialize() public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    {

    }
}