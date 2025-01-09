package integration

import (
	"log/slog"
	"path/filepath"
	"slices"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoland"
	"github.com/gnolang/gno/gno.land/pkg/gnoland/ugnot"
	abci "github.com/gnolang/gno/tm2/pkg/bft/abci/types"
	tmcfg "github.com/gnolang/gno/tm2/pkg/bft/config"
	"github.com/gnolang/gno/tm2/pkg/bft/node"
	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/db/memdb"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/require"
)

const (
	DefaultAccount_Name    = "test1"
	DefaultAccount_Address = "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
	DefaultAccount_Seed    = "source bonus chronic canvas draft south burst lottery vacant surface solve popular case indicate oppose farm nothing bullet exhibit title speed wink action roast"
)

// TestingInMemoryNode initializes and starts an in-memory node for testing.
// It returns the node instance and its RPC remote address.
func TestingInMemoryNode(t TestingTS, logger *slog.Logger, config *gnoland.InMemoryNodeConfig) (*node.Node, string) {
	node, err := gnoland.NewInMemoryNode(logger, config)
	require.NoError(t, err)

	err = node.Start()
	require.NoError(t, err)

	ourAddress := config.PrivValidator.GetPubKey().Address()
	isValidator := slices.ContainsFunc(config.Genesis.Validators, func(val bft.GenesisValidator) bool {
		return val.Address == ourAddress
	})

	// Wait for first block if we are a validator.
	// If we are not a validator, we don't produce blocks, so node.Ready() hangs.
	if isValidator {
		select {
		case <-node.Ready():
		case <-time.After(time.Second * 10):
			require.FailNow(t, "timeout while waiting for the node to start")
		}
	}

	return node, node.Config().RPC.ListenAddress
}

// TestingNodeConfig constructs an in-memory node configuration
// with default packages and genesis transactions already loaded.
// It will return the default creator address of the loaded packages.
func TestingNodeConfig(t TestingTS, gnoroot string, additionalTxs ...gnoland.TxWithMetadata) (*gnoland.InMemoryNodeConfig, bft.Address) {
	cfg := TestingMinimalNodeConfig(gnoroot)
	cfg.SkipGenesisVerification = true

	creator := crypto.MustAddressFromString(DefaultAccount_Address) // test1

	params := LoadDefaultGenesisParamFile(t, gnoroot)
	balances := LoadDefaultGenesisBalanceFile(t, gnoroot)
	txs := make([]gnoland.TxWithMetadata, 0)
	txs = append(txs, LoadDefaultPackages(t, creator, gnoroot)...)
	txs = append(txs, additionalTxs...)

	cfg.Genesis.AppState = gnoland.GnoGenesisState{
		Balances: balances,
		Txs:      txs,
		Params:   params,
	}

	return cfg, creator
}

// TestingMinimalNodeConfig constructs the default minimal in-memory node configuration for testing.
func TestingMinimalNodeConfig(gnoroot string) *gnoland.InMemoryNodeConfig {
	tmconfig := DefaultTestingTMConfig(gnoroot)

	// Create Mocked Identity
	pv := gnoland.NewMockedPrivValidator()

	// Generate genesis config
	genesis := DefaultTestingGenesisConfig(gnoroot, pv.GetPubKey(), tmconfig)

	return &gnoland.InMemoryNodeConfig{
		PrivValidator: pv,
		Genesis:       genesis,
		TMConfig:      tmconfig,
		DB:            memdb.NewMemDB(),
		InitChainerConfig: gnoland.InitChainerConfig{
			GenesisTxResultHandler: gnoland.PanicOnFailingTxResultHandler,
			CacheStdlibLoad:        true,
		},
	}
}

func DefaultTestingGenesisConfig(gnoroot string, self crypto.PubKey, tmconfig *tmcfg.Config) *bft.GenesisDoc {
	return &bft.GenesisDoc{
		GenesisTime: time.Now(),
		ChainID:     tmconfig.ChainID(),
		ConsensusParams: abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxTxBytes:   1_000_000,   // 1MB,
				MaxDataBytes: 2_000_000,   // 2MB,
				MaxGas:       100_000_000, // 100M gas
				TimeIotaMS:   100,         // 100ms
			},
		},
		Validators: []bft.GenesisValidator{
			{
				Address: self.Address(),
				PubKey:  self,
				Power:   10,
				Name:    "self",
			},
		},
		AppState: gnoland.GnoGenesisState{
			Balances: []gnoland.Balance{
				{
					Address: crypto.MustAddressFromString(DefaultAccount_Address),
					Amount:  std.MustParseCoins(ugnot.ValueString(10_000_000_000_000)),
				},
			},
			Txs:    []gnoland.TxWithMetadata{},
			Params: []gnoland.Param{},
		},
	}
}

// LoadDefaultPackages loads the default packages for testing using a given creator address and gnoroot directory.
func LoadDefaultPackages(t TestingTS, creator bft.Address, gnoroot string) []gnoland.TxWithMetadata {
	examplesDir := filepath.Join(gnoroot, "examples")

	defaultFee := std.NewFee(50000, std.MustParseCoin(ugnot.ValueString(1000000)))
	txs, err := gnoland.LoadPackagesFromDir(examplesDir, creator, defaultFee)
	require.NoError(t, err)

	return txs
}

// LoadDefaultGenesisBalanceFile loads the default genesis balance file for testing.
func LoadDefaultGenesisBalanceFile(t TestingTS, gnoroot string) []gnoland.Balance {
	balanceFile := filepath.Join(gnoroot, "gno.land", "genesis", "genesis_balances.txt")

	genesisBalances, err := gnoland.LoadGenesisBalancesFile(balanceFile)
	require.NoError(t, err)

	return genesisBalances
}

// LoadDefaultGenesisParamFile loads the default genesis balance file for testing.
func LoadDefaultGenesisParamFile(t TestingTS, gnoroot string) []gnoland.Param {
	paramFile := filepath.Join(gnoroot, "gno.land", "genesis", "genesis_params.toml")

	genesisParams, err := gnoland.LoadGenesisParamsFile(paramFile)
	require.NoError(t, err)

	return genesisParams
}

// LoadDefaultGenesisTXsFile loads the default genesis transactions file for testing.
func LoadDefaultGenesisTXsFile(t TestingTS, chainid string, gnoroot string) []gnoland.TxWithMetadata {
	txsFile := filepath.Join(gnoroot, "gno.land", "genesis", "genesis_txs.jsonl")

	// NOTE: We dont care about giving a correct address here, as it's only for display
	// XXX: Do we care loading this TXs for testing ?
	genesisTXs, err := gnoland.LoadGenesisTxsFile(txsFile, chainid, "https://127.0.0.1:26657")
	require.NoError(t, err)

	return genesisTXs
}

// DefaultTestingTMConfig constructs the default Tendermint configuration for testing.
func DefaultTestingTMConfig(gnoroot string) *tmcfg.Config {
	const defaultListner = "tcp://127.0.0.1:0"

	tmconfig := tmcfg.TestConfig().SetRootDir(gnoroot)
	tmconfig.Consensus.WALDisabled = true
	tmconfig.Consensus.SkipTimeoutCommit = true
	tmconfig.Consensus.CreateEmptyBlocks = true
	tmconfig.Consensus.CreateEmptyBlocksInterval = time.Millisecond * 100
	tmconfig.RPC.ListenAddress = defaultListner
	tmconfig.P2P.ListenAddress = defaultListner
	return tmconfig
}
