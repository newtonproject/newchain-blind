package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/cryptoballot/rsablind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"gitlab.newtonproject.org/yangchenzhong/NewChainBlind/blind"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

func (cli *CLI) buildSignCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:   "sign <hexFile.blinded> <address>",
		Short: "Sign blinded file for address of cash 1NEW",
		Args:  cobra.MinimumNArgs(2),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			key, err := blind.LoadPEMPrivKeyFile(cli.bank.RSAPEMPrivateKeyFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			blindedHexByte, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			blinded, err := hex.DecodeString(string(blindedHexByte))
			if err != nil {
				fmt.Println(err)
				return
			}

			if !common.IsHexAddress(args[1]) {
				fmt.Println("convert sting to address error")
				return
			}
			user := common.HexToAddress(args[1])

			// check balance
			balances := viper.GetStringMap("Balances")
			userL := strings.ToLower(user.String())
			b, ok := balances[userL]
			if !ok {
				fmt.Println("In funds")
				return
			}
			// ok, update balance, 1NEW
			bb, ok := b.(int64)
			if !ok {
				fmt.Println("convert balance to int64 error")
				return
			}
			balances[userL] = bb - 1
			fmt.Printf("Sub 1 NEW for %s\n", user.String())
			fmt.Printf("Current balance of %s is %d NEW\n", user.String(), balances[userL])
			viper.Set("Balances", balances)
			err = viper.WriteConfig()
			if err != nil {
				fmt.Println(err)
				return
			}

			sig, err := rsablind.BlindSign(key, blinded)
			if err != nil {
				fmt.Println(err)
				return
			}

			signFile := fmt.Sprintf("%s.sig", args[0])
			if err := ioutil.WriteFile(signFile,
				[]byte(hex.EncodeToString(sig)), 0644); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Sign: ", signFile)
			fmt.Println("Send signed blinded file back to user")

		},
	}

	return faucetCmd
}
