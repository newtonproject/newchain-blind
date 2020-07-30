package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cli *CLI) buildRootCmd() {
	if cli.rootCmd != nil {
		cli.rootCmd.ResetFlags()
		cli.rootCmd.ResetCommands()
	}

	rootCmd := &cobra.Command{
		Use:              "NewChainBlind",
		Short:            "NewChainBlind is a commandline client for the NewChain",
		Run:              cli.help,
		PersistentPreRun: cli.setup,
	}
	cli.rootCmd = rootCmd

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cli.config, "config", "c", defaultConfigFile, "The `path` to config file")
	rootCmd.PersistentFlags().StringP("walletPath", "w", defaultWalletPath, "Wallet storage `directory`")
	rootCmd.PersistentFlags().StringP("rpcURL", "i", defaultRPCURL, "NewChain json rpc or ipc `url`")

	// Basic commands
	rootCmd.AddCommand(cli.buildVersionCmd()) // version
	rootCmd.AddCommand(cli.buildInitCmd())    // init

	// Aux commands
	rootCmd.AddCommand(cli.buildBalanceCmd()) // balance
	rootCmd.AddCommand(cli.buildFaucetCmd())  // faucet

	// Alias commands
	rootCmd.AddCommand(cli.buildAccountCmd()) // account

	// Core commands
	rootCmd.AddCommand(cli.buildInfoCmd())    // bank info
	rootCmd.AddCommand(cli.buildDepositCmd()) // deposit NEW by hash

	rootCmd.AddCommand(cli.buildSignCmd())   // bank sign
	rootCmd.AddCommand(cli.buildVerifyCmd()) // bank verify

	rootCmd.AddCommand(cli.buildBlindCmd())   // customer blind
	rootCmd.AddCommand(cli.buildUnblindCmd()) // customer unblind
}

func EmptyRun(cmd *cobra.Command, args []string) {
	if err := cobra.OnlyValidArgs(cmd, args); err != nil {
		fmt.Println(err)
		fmt.Printf("Run '%v --help' for usage.\n", cmd.CommandPath())
		return
	}
	fmt.Println(cmd.UsageString())
}
