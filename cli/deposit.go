package cli

import (
	"context"
	"encoding/json"
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

func (cli *CLI) buildDepositCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:   "deposit <TxHash>",
		Short: "Use the TxHash to deposit NEW(decimal not counted)",
		DisableFlagsInUseLine: true,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}
