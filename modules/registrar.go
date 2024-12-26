package modules

import (
	"github.com/forbole/callisto/v4/modules/actions"
	"github.com/forbole/callisto/v4/modules/types"

	"github.com/forbole/juno/v5/modules/pruning"
	"github.com/forbole/juno/v5/modules/telemetry"

	"github.com/forbole/callisto/v4/modules/slashing"

	jmodules "github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/modules/messages"
	"github.com/forbole/juno/v5/modules/registrar"

	"github.com/forbole/callisto/v4/utils"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/auth"
	"github.com/forbole/callisto/v4/modules/bank"
	"github.com/forbole/callisto/v4/modules/consensus"
	"github.com/forbole/callisto/v4/modules/feegrant"

	dailyrefetch "github.com/forbole/callisto/v4/modules/daily_refetch"
	messagetype "github.com/forbole/callisto/v4/modules/message_type"
	"github.com/forbole/callisto/v4/modules/modules"
	juno "github.com/forbole/juno/v5/types"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(tx *juno.Tx) ([]string, error) {
		addresses, err := parser(tx)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Codec
	db := database.Cast(ctx.Database)

	// we should modify the sources later.
	sources, err := types.BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	// starts a server for many operations with endpoints (hosted by us) for
	// users to call. not needed at the moment. only triggers once at the node
	// start, to launch HTTP server. needs to be updated.
	actionsModule := actions.NewModule(ctx.JunoConfig, ctx.EncodingConfig)
	// at genesis, gets and saves normal and vesting accounts
	// regularly, (1) looks for MsgCreateVestingAccount messages and saves the
	// resulting accounts, and (2) refreshes all accounts involved in the
	// message. should retain.
	authModule := auth.NewModule(r.parser, cdc, db)
	// tracks total supply every 10 minutes and not more continuously. should
	// retain.
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	// at genesis, saves time and height.
	// at each block, overwrites the average block time at height in db.
	// every minute, hour, day => same as above.
	// should retain.
	consensusModule := consensus.NewModule(db)
	// processes missing blocks and parses them.
	dailyRefetchModule := dailyrefetch.NewModule(ctx.Proxy, db)

	// at genesis, saves the params.
	// tracks community pool every 1 hour or upon receipt of MsgFundCommunityPool.
	// should drop, since our mechanism is different.
	// // distrModule := distribution.NewModule(sources.DistrSource, cdc, db)

	// tracks and updates fee grants given out or expired (per block) or revoked.
	// should retain.
	feegrantModule := feegrant.NewModule(cdc, db)

	// for each message, stores the type.
	messagetypeModule := messagetype.NewModule(r.parser, cdc, db)

	// dropping because we use a custom mint module.
	// // mintModule := mint.NewModule(sources.MintSource, cdc, db)

	// stores params at start and updates signing info at each block.
	// should retain.
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)

	// dropping because we don't use the default staking module.
	// // stakingModule := staking.NewModule(sources.StakingSource, cdc, db)

	// TBD for both of these
	// // govModule := gov.NewModule(sources.GovSource, distrModule, mintModule, slashingModule, stakingModule, cdc, db)
	// // upgradeModule := upgrade.NewModule(db, stakingModule)

	return []jmodules.Module{
		// saves messages verbatim (with some IBC parsing) for each msg.
		messages.NewModule(r.parser, cdc, ctx.Database),
		// exposes prometheus metrics, at node start.
		telemetry.NewModule(ctx.JunoConfig),
		// indexer-level pruning, at each block.
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),
		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		dailyRefetchModule,
		feegrantModule,
		messagetypeModule,
		// saves a list of modules configured in the indexer, once.
		modules.NewModule(ctx.JunoConfig.Chain, db),
		// needs coingecko, so can't activate it yet.
		// // pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
	}
}
