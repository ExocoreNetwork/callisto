package parse

import (
	parse "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"

	parseblocks "github.com/forbole/juno/v5/cmd/parse/blocks"

	parsegenesis "github.com/forbole/juno/v5/cmd/parse/genesis"

	parsetransaction "github.com/forbole/juno/v5/cmd/parse/transactions"

	parseauth "github.com/forbole/bdjuno/v4/cmd/parse/auth"
	parsebank "github.com/forbole/bdjuno/v4/cmd/parse/bank"
	parseccv "github.com/forbole/bdjuno/v4/cmd/parse/ccv"
	parsefeegrant "github.com/forbole/bdjuno/v4/cmd/parse/feegrant"
	parsepricefeed "github.com/forbole/bdjuno/v4/cmd/parse/pricefeed"
	parsewasm "github.com/forbole/bdjuno/v4/cmd/parse/wasm"
)

// NewParseCmd returns the Cobra command allowing to parse some chain data without having to re-sync the whole database
func NewParseCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Parse some data without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parse.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		parseauth.NewAuthCmd(parseCfg),
		parsebank.NewBankCmd(parseCfg),
		parseblocks.NewBlocksCmd(parseCfg),
		parseccv.NewCcvCmd(parseCfg),
		parsefeegrant.NewFeegrantCmd(parseCfg),
		parsegenesis.NewGenesisCmd(parseCfg),
		parsepricefeed.NewPricefeedCmd(parseCfg),
		parsetransaction.NewTransactionsCmd(parseCfg),
		parsewasm.NewWasmCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}
