package initiates

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	crypto_ "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

func InitListenHost(port int, nodeKey string) (*ecdsa.PrivateKey, host.Host, error) {
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, nil, err
	}
	prvKey, err := crypto.UnmarshalSecp256k1PrivateKey(common.FromHex(nodeKey))
	if err != nil {
		return nil, nil, err
	}
	pk, err := crypto_.ToECDSA(common.FromHex(nodeKey))
	if err != nil {
		return nil, nil, err
	}
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.EnableRelay(),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return nil, nil, err
	}
	return pk, host, nil
}

func InitHost(nodeKey string) (*ecdsa.PrivateKey, host.Host, error) {
	prvKey, err := crypto.UnmarshalSecp256k1PrivateKey(common.FromHex(nodeKey))
	if err != nil {
		return nil, nil, err
	}
	pk, err := crypto_.ToECDSA(common.FromHex(nodeKey))
	if err != nil {
		return nil, nil, err
	}
	host, err := libp2p.New(
		libp2p.NoListenAddrs,
		libp2p.EnableRelay(),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return nil, nil, err
	}
	return pk, host, nil
}
