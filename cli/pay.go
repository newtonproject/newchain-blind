package cli

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildPayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay <amount> <--to target> [-u NEW|WEI] [--from source] [--data text] [-p 100] [-g 21000] [-n 1]",
		Short: "Send [amount] [unit] from [source] to [target] with message [text]",
		Args:  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.Flags().String("from", "", "source account seed or name")
	cmd.Flags().String("to", "", "target account address or name")
	unitUsageString := fmt.Sprintf("unit for pay amount. %s.", DenominationString)
	cmd.Flags().StringP("unit", "u", "NEW", unitUsageString)
	viper.BindPFlag("Bank.Address", cmd.Flags().Lookup("from"))
	viper.BindPFlag("unit", cmd.Flags().Lookup("unit"))
	cmd.Flags().String("data", "", "custom data message (use quotes if there are spaces)")
	cmd.Flags().Uint64P("gas", "g", 21000, "the gas provided for the transaction execution")
	cmd.Flags().Uint64P("price", "p", 1, "the gasPrice used for each paid gas (unit in WEI)")
	cmd.Flags().Uint64P("nonce", "n", 0, "the number of nonce")

	return cmd
}

func (cli *CLI) openWallet(check bool) error {
	if cli.wallet == nil {
		cli.wallet = keystore.NewKeyStore(cli.walletPath,
			keystore.LightScryptN, keystore.LightScryptP)
	}

	if check && len(cli.wallet.Accounts()) == 0 {
		return errWalletPathEmppty
	}
	return nil
}

//
// func (cli *CLI) unlockWallet(account accounts.Account) error {
// 	if cli.wallet == nil {
// 		if err := cli.openWallet(true); err != nil {
// 			return err
// 		}
// 	}
// 	if account.Address == (common.Address{}) {
// 		return errRequiredFromAddress
// 	}
// 	if _, err := cli.wallet.Find(account); err != nil {
// 		return fmt.Errorf("%v (%s)", err, account.Address.String())
// 	}
//
// 	// var trials int
// 	// var err error
// 	// walletPassword := cli.tran.Password
// 	// for trials = 0; trials < 3; trials++ {
// 	// 	prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", account.Address.String(), trials+1, 3)
// 	// 	if walletPassword == "" {
// 	// 		walletPassword, _ = getPassPhrase(prompt, false)
// 	// 	} else {
// 	// 		fmt.Println(prompt, "\nUse the the password has set")
// 	// 	}
// 	// 	err = cli.wallet.Unlock(account, walletPassword)
// 	// 	if err == nil {
// 	// 		break
// 	// 	}
// 	// 	walletPassword = ""
// 	// }
//
// 	// if trials >= 3 {
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	return fmt.Errorf("Error: Failed to unlock account %s (%v)", account.Address.String(), err)
// 	// }
//
// 	var err error
// 	walletPassword := cli.tran.Password
// 	if walletPassword == "" {
// 		prompt := fmt.Sprintf("Unlocking account %s", account.Address.String())
// 		walletPassword, _ = getPassPhrase(prompt, false)
// 	}
// 	err = cli.wallet.Unlock(account, walletPassword)
// 	walletPassword = ""
// 	if err != nil {
// 		// return fmt.Errorf("Error: Failed to unlock account %s (%v)", account.Address.String(), err)
// 		return err
// 	}
//
// 	return nil
// }
//
// func (cli *CLI) unlockAndPay() (common.Hash, error) {
// 	account := accounts.Account{Address: cli.tran.From}
// 	if err := cli.unlockWallet(account); err != nil {
// 		return common.Hash{}, err
// 	}
//
// 	signTx, err := cli.signTx()
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
//
// 	if err := cli.sendSignTx(signTx); err != nil {
// 		return common.Hash{}, err
// 	}
//
// 	return signTx.Hash(), nil
// }
//
// func (cli *CLI) unlockAndSignTx() (*types.Transaction, error) {
// 	account := accounts.Account{Address: cli.tran.From}
// 	if err := cli.unlockWallet(account); err != nil {
// 		return nil, err
// 	}
//
// 	return cli.signTx()
// }
//
// func (cli *CLI) sendSignTx(signTx *types.Transaction) error {
// 	if cli.client == nil {
// 		if err := cli.BuildClient(); err != nil {
// 			return err
// 		}
// 	}
// 	if err := cli.client.SendTransaction(context.Background(), signTx); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (cli *CLI) signTx() (*types.Transaction, error) {
// 	if cli.tran == nil {
// 		return nil, errCliTranNil
// 	}
// 	tx := types.NewTransaction(cli.tran.Nonce, cli.tran.To, cli.tran.Value, cli.tran.GasLimit, cli.tran.GasPrice, cli.tran.Data)
// 	signTx, err := cli.wallet.SignTx(accounts.Account{Address: cli.tran.From}, tx, cli.tran.NetworkID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return signTx, nil
// }
//
// func (cli *CLI) printTxIndent() {
// 	if cli.tran != nil {
// 		tByte, err := cli.tran.MarshalJSON(true)
// 		if err == nil {
// 			fmt.Println(string(tByte))
// 		}
// 	}
// }
//
// func (cli *CLI) getNonce() (uint64, error) {
// 	if cli.tran == nil {
// 		return 0, errArgs0
// 	}
// 	if cli.client == nil {
// 		if err := cli.BuildClient(); err != nil {
// 			return 0, fmt.Errorf("Failed to connect to the NewChain client: %v", err)
// 		}
// 	}
// 	return cli.client.NonceAt(context.Background(), cli.tran.From, nil)
// }
//
// func (cli *CLI) getGasPrice() (*big.Int, error) {
// 	if cli.client == nil {
// 		if err := cli.BuildClient(); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return cli.client.SuggestGasPrice(context.Background())
// }
//
// func (cli *CLI) getGasLimit() (uint64, error) {
// 	if cli.client == nil {
// 		if err := cli.BuildClient(); err != nil {
// 			return 21000, err
// 		}
// 	}
// 	msg := ethereum.CallMsg{
// 		From:     cli.tran.From,
// 		To:       &cli.tran.To,
// 		Value:    cli.tran.Value,
// 		Data:     cli.tran.Data,
// 		GasPrice: cli.tran.GasPrice,
// 	}
// 	return cli.client.EstimateGas(context.Background(), msg)
// }
//
// func (cli *CLI) getNetworkID() (*big.Int, error) {
// 	if cli.client == nil {
// 		if err := cli.BuildClient(); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return cli.client.NetworkID(context.Background())
// }
//
// func (cli *CLI) updateFromNodeCustom(bNonce, bGasPrice, bGasLimit, bNetworkID bool) error {
// 	// get nonce
// 	if bNonce {
// 		nonce, err := cli.getNonce()
// 		if err != nil {
// 			return err
// 		}
// 		cli.tran.Nonce = nonce
// 	}
//
// 	// get gasPrice
// 	if bGasPrice {
// 		gasPrice, err := cli.getGasPrice()
// 		if err != nil {
// 			return err
// 		}
// 		cli.tran.GasPrice = gasPrice
// 	}
//
// 	// get gasLimit
// 	if bGasLimit {
// 		gasLimit, err := cli.getGasLimit()
// 		if err != nil {
// 			return err
// 		}
// 		cli.tran.GasLimit = gasLimit
// 	}
//
// 	// get ChainID
// 	if bNetworkID {
// 		networkID, err := cli.getNetworkID()
// 		if err != nil {
// 			return err
// 		}
// 		cli.tran.NetworkID = networkID
// 	}
//
// 	return nil
// }
//
// func (cli *CLI) updateFromNode() error {
// 	return cli.updateFromNodeCustom(true, true, true, true)
// }
