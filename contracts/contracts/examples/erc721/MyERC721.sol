// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721BurnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721RoyaltyUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721EnumerableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract MyERC721 is Initializable, UUPSUpgradeable, ERC721Upgradeable, ERC721RoyaltyUpgradeable, ERC721BurnableUpgradeable, ERC721EnumerableUpgradeable, AccessControlUpgradeable {
    using Strings for uint256;
    string private _baseTokenURI;
    string private _uriSuffix;
    bool public unTransferable;

    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant MINT_ROLE = keccak256("mint_role");
    bytes32 public constant BURN_ROLE = keccak256("burn_role");
    mapping(uint256 => bool) public blocklist;

    function initialize(string calldata name, string calldata symbol) public initializer {
        __ERC721_init(name, symbol);
        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ERC721Enumerable_init();
        __ERC721Burnable_init();
        __ERC721Royalty_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(MINT_ROLE, msg.sender);
        _grantRole(BURN_ROLE, msg.sender);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    {}

    function supportsInterface(bytes4 interfaceId) public view virtual override( ERC721Upgradeable,ERC721RoyaltyUpgradeable, ERC721EnumerableUpgradeable, AccessControlUpgradeable) returns (bool) {
        return interfaceId == type(IERC721).interfaceId || interfaceId == type(IERC721Metadata).interfaceId || interfaceId == type(IERC721Enumerable).interfaceId || interfaceId == type(IAccessControl).interfaceId || super.supportsInterface(interfaceId);
    }

    function baseURI() external view virtual returns (string memory) {
        return _baseTokenURI;
    }

    function uriSuffix() external view virtual returns (string memory) {
        return _uriSuffix;
    }

    function setUnTransferable(bool _unTransferable) external onlyRole(ADMIN_ROLE) {
        unTransferable = _unTransferable;
    }

    function _update(address to, uint256 tokenId, address auth) internal override(ERC721Upgradeable, ERC721EnumerableUpgradeable) virtual returns (address) {
        if (!(hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || hasRole(ADMIN_ROLE, msg.sender) || hasRole(MINT_ROLE, msg.sender)  || hasRole(BURN_ROLE, msg.sender))) {
            require(!unTransferable, "unTransferable");
        }
        require(!blocklist[tokenId], "it is blocklist");
        return ERC721EnumerableUpgradeable._update(to, tokenId, auth);
    }

    function setBlocklist(uint256[] calldata _tokenIds, bool _blocklist) external onlyRole(ADMIN_ROLE) {
        for (uint256 i = 0; i < _tokenIds.length; i++) {
            blocklist[_tokenIds[i]] = _blocklist;
        }
    }

    function _increaseBalance(address account, uint128 value) internal override(ERC721Upgradeable, ERC721EnumerableUpgradeable) {
        ERC721EnumerableUpgradeable._increaseBalance(account, value);
    }

    function setBaseURI(string calldata newBaseTokenURI) external onlyRole(ADMIN_ROLE) {
        _baseTokenURI = newBaseTokenURI;
    }

    function setURISuffix(string calldata newUriSuffix) external onlyRole(ADMIN_ROLE) {
        _uriSuffix = newUriSuffix;
    }

    function tokenURI(uint256 tokenId) public view virtual override returns (string memory) {
        if (keccak256(abi.encodePacked(_uriSuffix)) == keccak256(abi.encodePacked("fixed"))) {
            return bytes(_baseTokenURI).length > 0 ? string(abi.encodePacked(_baseTokenURI)) : "";
        } else {
            return bytes(_baseTokenURI).length > 0 ? string(abi.encodePacked(_baseTokenURI, tokenId.toString(), _uriSuffix)) : "";
        }
    }

    function existOf(uint256 tokenId) external view virtual returns (bool) {
        address owner = _ownerOf(tokenId);
        if (owner == address(0)) {
            return false;
        }
        return true;
    }

    function mint(address to, uint256 tokenId) external onlyRole(MINT_ROLE) {
        _mint(to, tokenId);
    }

    function batchMint(address[] calldata toAddresses, uint256[] calldata tokenIds ) external onlyRole(MINT_ROLE) {
        for(uint256 i = 0; i < toAddresses.length; i++) {
            _mint(toAddresses[i], tokenIds[i]);
        }
    }

    function burn(address from, uint256 tokenId) external onlyRole(BURN_ROLE) {
        require(_ownerOf(tokenId) == from, "owner mismatch");
        _burn(tokenId);
    }

    function burn(uint256 tokenId) public override onlyRole(BURN_ROLE) {
        _burn(tokenId);
    }

}
