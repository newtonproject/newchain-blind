package cli

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildInfoCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:   "info",
		Short: "Show the info of the bank",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bank Dashboard")
			fmt.Println("Address: ", viper.GetString("Bank.Address"))
			fmt.Println("Cash Type: Only 1NEW Support")
			fmt.Println("Rsa Public Key of 1NEW: ", viper.GetString("Bank.RSAPEMPublicKeyFile"))

			balances := viper.GetStringMap("Balances")
			fmt.Printf("Balances(%d):\n", len(balances))
			for a, b := range balances {
				fmt.Printf("\t%s\t%vNEW\n", common.HexToAddress(a).String(), b)
			}

			if common.IsHexAddress(viper.GetString("Bank.Address")) {
				client, err := ethclient.Dial(cli.rpcURL)
				if err != nil {
					fmt.Println(err)
					return
				}
				balance, err := client.BalanceAt(context.Background(),
					common.HexToAddress(viper.GetString("Bank.Address")), nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Online Balance: ", getWeiAmountTextUnitByUnit(balance, "NEW"))
			}

		},
	}

	return faucetCmd
}
