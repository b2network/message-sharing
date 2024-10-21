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
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external returns (bool success, string memory error);
}

enum NftStandard {
    UNKNOWN,
    ERC721,
    ERC1155
}

struct Nft {
    NftStandard standard;
    address nft_address;
    bool status;
}

interface IERC721 {
    function balanceOf(address owner) external view returns (uint256 balance);
    function ownerOf(uint256 tokenId) external view returns (address owner);
    function safeTransferFrom(address from, address to, uint256 tokenId) external;
    function mint(address to, uint256 tokenId) external;
}

interface IERC1155 {
    function balanceOf(address account, uint256 id) external view returns (uint256);
    function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes calldata data) external;
    function mint(address to, uint256 id, uint256 value) external;
}

contract NftBridge is IBusinessContract, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

    // message sharing address
    IMessageSharing public messageSharing;
    // from_chain_id => from_nft_address => to_chain_id => Nft info
    mapping(uint256 => mapping (address => mapping (uint256 => Nft))) public nft_mapping;
    // from_chain_id => nft bridge address
    mapping (uint256 => address) public bridges;
    // from_chain_id => from_id => execute
    mapping (uint256 => mapping (uint256 => bool)) public executes;

    event Lock(uint256 from_chain_id, uint256 from_id, address from_address, address from_token_address, address to_token_address, uint256 to_chain_id, address to_token_bridge, address to_address, uint256 amount);

    event Unlock(uint256 from_chain_id, uint256 from_id, address from_address, address from_token_address, address to_token_address, uint256 token_id, uint256 to_chain_id, address to_token_bridge, address to_address, uint256 amount);

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success, string memory error) {
        require(bridges[from_chain_id] == from_sender, "Invalid chain id or from_sender");
        require(!executes[from_chain_id][from_id], "Have been executed");
        executes[from_chain_id][from_id] = true;
        (address from_token_address, address to_token_address, uint256 token_id, address from_address, address to_address, uint256 amount) = decodeLockData(message);
        transferNft(from_chain_id, from_token_address, to_token_address, token_id, to_address, amount);
        emit Unlock(from_chain_id, from_id, from_address, from_token_address, to_token_address, token_id , block.chainid, address(this), to_address, amount);
        return (true, "");
    }

    function transferNft(uint256 from_chain_id, address from_token_address, address to_token_address, uint256 token_id, address to_address, uint256 amount) internal {
        Nft memory nft = nft_mapping[from_chain_id][from_token_address][block.chainid];
        require(nft.standard != NftStandard.UNKNOWN, "Invalid nft info");
        require(nft.nft_address == to_token_address, "Invalid nft address");
        require(nft.status, "nft status is false");
        if (nft.standard == NftStandard.ERC721) {
            interTransferNft721(nft.nft_address, to_address, token_id);
        } else if (nft.standard == NftStandard.ERC1155) {
            interTransferNft1155(nft.nft_address, to_address, token_id, amount);
        }
    }

    function lock(address nft_address, uint256 token_id, uint256 amount, uint256 to_chain_id, address to_token_bridge, address to_address) external {
        Nft memory nft = nft_mapping[block.chainid][nft_address][to_chain_id];
        require(nft.standard != NftStandard.UNKNOWN, "Invalid nft info");
        require(nft.status, "nft status is false");
        if (nft.standard == NftStandard.ERC721) {
            require(amount == 1, "Invalid amount");
            IERC721(nft_address).safeTransferFrom(msg.sender, address(this), token_id);
        } else if (nft.standard == NftStandard.ERC1155) {
            require(amount > 0, "Invalid amount");
            IERC1155(nft_address).safeTransferFrom(msg.sender, address(this), token_id, amount, "");
        }

        bytes memory to_message = encodeLockData(nft_address, nft.nft_address, token_id, msg.sender, to_address, amount);

        uint256 from_id =  messageSharing.call(to_chain_id, to_token_bridge, to_message);
        emit Lock(block.chainid ,from_id, msg.sender, nft_address, nft.nft_address, to_chain_id, to_token_bridge, to_address, amount);
    }

    function setMessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE) {
        messageSharing = IMessageSharing(sharing_address);
    }

    function setBridges(uint256 from_chain_id, address bridge) external onlyRole(ADMIN_ROLE) {
        bridges[from_chain_id] = bridge;
    }

    function setNftMapping(uint256 from_chain_id, address from_nft_address, uint256 to_chain_id, address to_nft_address, NftStandard standard, bool status) external onlyRole(ADMIN_ROLE) {
        require(from_chain_id == block.chainid || to_chain_id == block.chainid , "Invalid chain id");
        nft_mapping[from_chain_id][from_nft_address][to_chain_id] = Nft({
            nft_address:to_nft_address,
            standard: standard,
            status: status
        });
    }

    function encodeLockData(address from_nft_address, address to_nft_address, uint256 token_id, address from_address, address to_address, uint256 amount) public pure returns (bytes memory) {
        return abi.encode(from_nft_address, to_nft_address, token_id, from_address, to_address, amount);
    }

    function decodeLockData(bytes memory data) public pure returns (address from_nft_address, address to_nft_address , uint256 token_id, address from_address, address to_address, uint256 amount) {
        (from_nft_address, to_nft_address, token_id, from_address, to_address, amount) = abi.decode(data, (address, address, uint256, address, address, uint256));
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

    function interTransferNft1155(address token_address, address to, uint256 token_id, uint256 amount) internal {
        uint256 balance = IERC1155(token_address).balanceOf(address(this), token_id);
        if (balance < amount) {
            IERC1155(token_address).mint(to, token_id, amount);
        } else {
            IERC1155(token_address).safeTransferFrom(address(this), to, token_id, amount, "");
        }
    }

    function interTransferNft721(address token_address, address to, uint256 token_id) internal {
        address from = IERC721(token_address).ownerOf(token_id);
        require(from == address(0x0) || from == address(this), "no operation permission");
        if (from == address(0x0)) {
                IERC721(token_address).mint(to, token_id);
        } else {
            IERC721(token_address).safeTransferFrom(address(this), to, token_id);
        }
    }

    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual returns (bytes4) {
        return this.onERC721Received.selector;
    }

    function onERC1155Received(
        address,
        address,
        uint256,
        uint256,
        bytes calldata
    ) public virtual returns (bytes4) {
        return bytes4(keccak256("onERC1155Received(address,address,uint256,uint256,bytes)"));
    }

    function onERC1155BatchReceived(
        address,
        address,
        uint256[] calldata,
        uint256[] calldata,
        bytes calldata
    ) public virtual returns (bytes4) {
        return bytes4(keccak256("onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)"));
    }

}