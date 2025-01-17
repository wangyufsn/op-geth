package misc

import (
	"github.com/ethereum-optimism/superchain-registry/superchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// The original create2deployer contract could not be deployed to Base mainnet at
// the canonical address of 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2 due to
// an accidental nonce increment from a deposit transaction. See
// https://github.com/pcaversaccio/create2deployer/issues/128 for context. This
// file applies the contract code to the canonical address manually in the Canyon
// hardfork.

// create2deployer is already deployed to Base testnets at 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2,
// so we deploy it to 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1 for hardfork testing purposes
var create2DeployerAddresses = map[uint64]common.Address{
	11763071:                  common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1"), // Base Goerli devnet
	params.BaseGoerliChainID:  common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1"), // Base Goerli testnet
	11763072:                  common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1"), // Base Sepolia devnet
	84532:                     common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1"), // Base Sepolia testnet
	params.BaseMainnetChainID: common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2"), // Base mainnet
}
var create2DeployerCodeHash = common.HexToHash("0xb0550b5b431e30d38000efb7107aaa0ade03d48a7198a140edda9d27134468b2")
var create2DeployerCode []byte

func init() {
	code, err := superchain.LoadContractBytecode(superchain.Hash(create2DeployerCodeHash))
	if err != nil {
		panic(err)
	}
	create2DeployerCode = code
}

func EnsureCreate2Deployer(c *params.ChainConfig, timestamp uint64, db vm.StateDB) {
	if !c.IsOptimism() || c.CanyonTime == nil || *c.CanyonTime != timestamp {
		return
	}
	address, ok := create2DeployerAddresses[c.ChainID.Uint64()]
	if !ok || db.GetCodeSize(address) > 0 {
		return
	}
	log.Info("Setting Create2Deployer code", "address", address, "codeHash", create2DeployerCodeHash)
	db.SetCode(address, create2DeployerCode)
}
