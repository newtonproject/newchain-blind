package cli

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [bank|user]",
		Short: "Initialize config file and bank",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bank", "user"},
		Run:                   EmptyRun,
		SuggestionsMinimumDistance: 1,
	}

	cmd.AddCommand(cli.buildInitBankCmd())
	cmd.AddCommand(cli.buildInitUserCmd())

	return cmd
}

func (cli *CLI) buildInitBankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bank",
		Short: "Initialize config file and bank",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Initialize config file")

			prompt := fmt.Sprintf("Enter file in which to save (%s): ", defaultConfigFile)
			configPath, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if configPath == "" {
				configPath = defaultConfigFile
			}
			cli.config = configPath

			walletPathV := viper.GetString("walletPath")
			prompt = fmt.Sprintf("Enter the wallet storage directory (%s): ", walletPathV)
			walletPath, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if walletPath == "" {
				walletPath = walletPathV
			}
			viper.Set("walletPath", walletPath)

			rpcURLV := viper.GetString("rpcURL")
			prompt = fmt.Sprintf("Enter NewChain json rpc or ipc url (%s): ", rpcURLV)
			rpcURL, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if rpcURL == "" {
				rpcURL = rpcURLV
			}
			viper.Set("rpcURL", rpcURL)

			privateKeyV := viper.GetString("Bank.RSAPEMPrivateKeyFile")
			prompt = fmt.Sprintf("Enter path of Private Key RSA Pem file(%s): ", privateKeyV)
			privateKey, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if privateKey == "" {
				privateKey = privateKeyV
			}
			viper.Set("Bank.RSAPEMPrivateKeyFile", privateKey)

			publicKeyV := viper.GetString("Bank.RSAPEMPublicKeyFile")
			prompt = fmt.Sprintf("Enter path of Public Key RSA Pem file(%s): ", publicKeyV)
			publicKey, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if publicKey == "" {
				publicKey = publicKeyV
			}
			viper.Set("Bank.RSAPEMPublicKeyFile", publicKey)

			prompt = fmt.Sprintf("Create a new bank account or not: [Y/n] ")
			createNewAddress, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if len(createNewAddress) <= 0 {
				createNewAddress = "Y"
			}
			if strings.ToUpper(createNewAddress[:1]) == "Y" {
				wallet := keystore.NewKeyStore(walletPath,
					keystore.StandardScryptN, keystore.StandardScryptP)

				walletPassword, err := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
				if err == nil {
					account, err := wallet.NewAccount(walletPassword)
					if err == nil {
						fromAddress := account.Address.String()
						fmt.Println(fromAddress)
						viper.Set("Bank.Address", fromAddress)
					} else {
						fmt.Println("Account error:", err)
						fmt.Println("Just create your account later.")
					}
				} else {
					fmt.Println("Error: ", err)
					fmt.Println("Just create your account later.")
				}
			}

			err = viper.WriteConfigAs(configPath)
			if err != nil {
				fmt.Println("WriteConfig:", err)
			} else {
				fmt.Println("Your configuration has been saved in ", configPath)
			}
		},
	}

	return cmd
}

func (cli *CLI) buildInitUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Initialize config file and user",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Initialize config file")

			prompt := fmt.Sprintf("Enter file in which to save (%s): ", defaultConfigFile)
			configPath, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if configPath == "" {
				configPath = defaultConfigFile
			}
			cli.config = configPath

			walletPathV := viper.GetString("walletPath")
			prompt = fmt.Sprintf("Enter the wallet storage directory (%s): ", walletPathV)
			walletPath, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if walletPath == "" {
				walletPath = walletPathV
			}
			viper.Set("walletPath", walletPath)

			rpcURLV := viper.GetString("rpcURL")
			prompt = fmt.Sprintf("Enter NewChain json rpc or ipc url (%s): ", rpcURLV)
			rpcURL, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if rpcURL == "" {
				rpcURL = rpcURLV
			}
			viper.Set("rpcURL", rpcURL)

			privateKeyV := viper.GetString("RSAPEMPrivateKeyFile")
			prompt = fmt.Sprintf("Enter path of Private Key RSA Pem file(%s): ", privateKeyV)
			privateKey, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if privateKey == "" {
				privateKey = privateKeyV
			}
			viper.Set("Bank.RSAPEMPrivateKeyFile", privateKey)

			publicKeyV := viper.GetString("RSAPEMPublicKeyFile")
			prompt = fmt.Sprintf("Enter path of Public Key RSA Pem file(%s): ", publicKeyV)
			publicKey, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if publicKey == "" {
				publicKey = publicKeyV
			}
			viper.Set("Bank.RSAPEMPublicKeyFile", publicKey)

			prompt = fmt.Sprintf("Create a new bank account or not: [Y/n] ")
			createNewAddress, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if len(createNewAddress) <= 0 {
				createNewAddress = "Y"
			}
			if strings.ToUpper(createNewAddress[:1]) == "Y" {
				wallet := keystore.NewKeyStore(walletPath,
					keystore.StandardScryptN, keystore.StandardScryptP)

				walletPassword, err := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
				if err == nil {
					account, err := wallet.NewAccount(walletPassword)
					if err == nil {
						fromAddress := account.Address.String()
						fmt.Println(fromAddress)
						viper.Set("user.from", fromAddress)
					} else {
						fmt.Println("Account error:", err)
						fmt.Println("Just create your account later.")
					}
				} else {
					fmt.Println("Error: ", err)
					fmt.Println("Just create your account later.")
				}
			}

			err = viper.WriteConfigAs(configPath)
			if err != nil {
				fmt.Println("WriteConfig:", err)
			} else {
				fmt.Println("Your configuration has been saved in ", configPath)
			}
		},
	}

	return cmd
}
