package cli

import (
	"context"
	"crypto"
	"fmt"
	"math/big"

	"github.com/cryptoballot/fdh"
	"github.com/cryptoballot/rsablind"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/newtonproject/newchain-blind/blind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildVerifyCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:                   "verify <hexFile.unblinder.sig> <hexFile.data> <address>",
		Short:                 "Verify unblinder sign file and send 1NEW to address",
		Args:                  cobra.MinimumNArgs(3),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if cli.bank.Address == (common.Address{}) {
				fmt.Println("Not bank address set, please run `init bank`")
				return
			}
			address := cli.bank.Address
			wallet := keystore.NewKeyStore(cli.walletPath,
				keystore.LightScryptN, keystore.LightScryptP)
			if !wallet.HasAddress(address) {
				fmt.Println("Not such address")
				return
			}

			key, err := blind.LoadPEMPrivKeyFile(cli.bank.RSAPEMPrivateKeyFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			unblindSig, err := readHexFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			hash := common.BytesToHash(fdh.Sum(crypto.SHA256, 256, unblindSig))
			hashList := viper.GetStringSlice("UsedUnblinderHash")
			for _, h := range hashList {
				if common.HexToHash(h) == hash {
					fmt.Println("Used Unblinder Hash")
					return
				}
			}
			hashList = append(hashList, hash.String())
			viper.Set("UsedUnblinderHash", hashList)
			err = viper.WriteConfig()
			if err != nil {
				fmt.Println(err)
				return
			}

			data, err := readHexFile(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			// check data is 1NEW
			cm := CashMessage{}
			err = rlp.DecodeBytes(data, &cm)
			if err != nil {
				fmt.Println(err)
				return
			}
			if cm.Cash != 1 {
				// TODO: from cash to get public key, not verify
				fmt.Println("Not support cash type")
				return
			}

			hashed := fdh.Sum(crypto.SHA256, 256, data)

			if !common.IsHexAddress(args[2]) {
				fmt.Println("convert sting to address error")
				return
			}
			user := common.HexToAddress(args[2])

			if err := rsablind.VerifyBlindSignature(&key.PublicKey, hashed, unblindSig); err != nil {
				fmt.Println(err)
				return
			}

			// send 1NEW to user
			client, err := ethclient.Dial(cli.rpcURL)
			if err != nil {
				fmt.Println(err)
				return
			}
			ctx := context.Background()
			chainID, err := client.NetworkID(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			gasPrice, err := client.SuggestGasPrice(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			nonce, err := client.NonceAt(ctx, cli.bank.Address, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			value := big.NewInt(0).Set(big1NEWInWEI)
			gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
				From:     cli.bank.Address,
				To:       &user,
				Value:    value,
				GasPrice: gasPrice,
				Data:     nil})
			tx := types.NewTransaction(nonce, user, value, gasLimit, gasPrice, nil)

			account := accounts.Account{Address: address}
			var trials int
			walletPassword := cli.bank.Password
			for trials = 0; trials < 3; trials++ {
				prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", account.Address.String(), trials+1, 3)
				if walletPassword == "" {
					walletPassword, _ = getPassPhrase(prompt, false)
				} else {
					fmt.Println(prompt, "\nUse the the password has set")
				}
				err = wallet.Unlock(account, walletPassword)
				if err == nil {
					break
				}
				walletPassword = ""
			}

			if trials >= 3 {
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Error: Failed to unlock account %s (%v)\n", account.Address.String(), err)
				return
			}

			signTx, err := wallet.SignTx(account, tx, chainID)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = client.SendTransaction(ctx, signTx)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Send 1 NEW to %s with hash %s\n", user.String(), signTx.Hash().String())
		},
	}

	return faucetCmd
}
