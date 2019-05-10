package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

func (cli *CLI) buildWithdrawCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:   "withdraw <address> <blinded>",
		Short: "Withdraw cash from bank with 1NEW blinded",
		DisableFlagsInUseLine: true,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if !common.IsHexAddress(args[0]) {
				fmt.Println("Convert string to address error")
				return
			}
			address := common.HexToAddress(args[0])

			blinded, err := hex.DecodeString(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(address.String(), blinded)
			return

			if cli.bank.Address == (common.Address{}) {
				fmt.Println("Not bank address set, please run `init bank`")
				return
			}

			hash := common.HexToHash(args[0])
			if hash == (common.Hash{}) {
				fmt.Println("Get hash error")
				return
			}
			hashList := viper.GetStringSlice("UsedHash")
			for _, h := range hashList {
				if common.HexToHash(h) == hash {
					fmt.Println("Used Hash")
					return
				}
			}
			hashList = append(hashList, hash.String())
			viper.Set("UsedHash", hashList)

			ctx := context.Background()
			c, err := rpc.Dial(cli.rpcURL)
			if err != nil {
				fmt.Println(err)
				return
			}
			client := ethclient.NewClient(c)
			r, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				fmt.Println(err)
				return
			}
			if r.Status != types.ReceiptStatusSuccessful {
				fmt.Println("Tx receipt status is failed")
				return
			}

			var json *rpcTransaction
			err = c.CallContext(ctx, &json, "eth_getTransactionByHash", hash)
			if err != nil {
				fmt.Println(err)
				return
			} else if json == nil {
				fmt.Println(ethereum.NotFound)
				return
			} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
				fmt.Println("server returned transaction without signature")
				return
			}
			if json.From == nil || json.BlockHash == nil {
				fmt.Println("Pending tx")
				return
			}
			from := *json.From
			tx := json.tx
			if tx.To() == nil || *tx.To() != cli.bank.Address {
				fmt.Println("Tx to is not bank address")
				return
			}
			value := big.NewInt(0).Div(tx.Value(), big1NEWInWEI).Int64()
			if value == 0 {
				fmt.Println("0 NEW to deposit")
				return
			}
			balances := viper.GetStringMap("Balances")
			fromL := strings.ToLower(from.String())
			b, ok := balances[fromL]
			if !ok {
				balances[fromL] = value
			} else {
				bb, ok := b.(int64)
				if !ok {
					fmt.Println("convert balance to int64 error")
					return
				}
				balances[fromL] = value + bb
			}
			fmt.Printf("Add %d NEW for %s\n", value, from.String())
			fmt.Printf("Current balance of %s is %d NEW\n", from.String(), balances[fromL])
			viper.Set("Balances", balances)
			err = viper.WriteConfig()
			if err != nil {
				fmt.Println(err)
				return
			}

		},
	}

	return faucetCmd
}
