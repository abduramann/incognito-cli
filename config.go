package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/incognito-cli/bridge/portal"
)

// Config represents the config of an environment of the CLI tool.
type Config struct {
	incClient *incclient.IncClient
	ethClient *ethclient.Client
	bscClient *ethclient.Client
	btcClient *portal.BTCClient

	ethVaultAddress common.Address
	bscVaultAddress common.Address
}

// NewConfig returns a new Config from given parameters.
func NewConfig(
	incClient *incclient.IncClient,
	ethClient, bscClient *ethclient.Client,
	btcClient *portal.BTCClient,
	ethVaultAddressStr, bscVaultAddressStr string,
) *Config {
	ethVaultAddress := common.HexToAddress(ethVaultAddressStr)
	bscVaultAddress := common.HexToAddress(bscVaultAddressStr)
	return &Config{
		incClient:       incClient,
		ethClient:       ethClient,
		bscClient:       bscClient,
		btcClient:       btcClient,
		ethVaultAddress: ethVaultAddress,
		bscVaultAddress: bscVaultAddress,
	}
}

// NewTestNetConfig creates a new testnet Config.
func NewTestNetConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewTestNetClient()
		} else {
			incClient, err = incclient.NewTestNetClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.TestNetETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.TestNetBSCHost)
	if err != nil {
		return err
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, ethClient, bscClient, btcClient, incclient.TestNetETHContractAddressStr, incclient.TestNetBSCContractAddressStr)

	return nil
}

// NewTestNet1Config creates a new testnet1 Config.
func NewTestNet1Config(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewTestNet1Client()
		} else {
			incClient, err = incclient.NewTestNet1ClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.TestNet1ETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.TestNet1BSCHost)
	if err != nil {
		return err
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, ethClient, bscClient, btcClient, incclient.TestNet1ETHContractAddressStr, incclient.TestNet1BSCContractAddressStr)
	return nil
}

// NewMainNetConfig creates a new main-net Config.
func NewMainNetConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewMainNetClient()
		} else {
			incClient, err = incclient.NewMainNetClientWithCache()
		}
		if err != nil {
			return err
		}
	}
	isMainNet = true

	ethClient, err := ethclient.Dial(incclient.MainNetETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.MainNetBSCHost)
	if err != nil {
		return err
	}

	btcClient, err := portal.NewBTCMainNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, ethClient, bscClient, btcClient, incclient.MainNetETHContractAddressStr, incclient.MainNetBSCContractAddressStr)
	return nil
}

// NewLocalConfig creates a new local Config.
func NewLocalConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewLocalClient("")
		} else {
			incClient, err = incclient.NewLocalClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, ethClient, bscClient, btcClient, incclient.LocalETHContractAddressStr, incclient.LocalETHContractAddressStr)
	return nil
}