package cli

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

const defaultConfigFile = "./config.toml"
const defaultWalletPath = "./wallet/"
const defaultRPCURL = "https://rpc1.newchain.newtonproject.org"

func (cli *CLI) defaultConfig() {
	viper.BindPFlag("walletPath", cli.rootCmd.PersistentFlags().Lookup("walletPath"))
	viper.BindPFlag("rpcURL", cli.rootCmd.PersistentFlags().Lookup("rpcURL"))

	viper.SetDefault("walletPath", defaultWalletPath)
	viper.SetDefault("rpcURL", defaultRPCURL)

	// bank
	viper.SetDefault("Bank.RSAPEMPrivateKeyFile", "./key.pem")
	viper.SetDefault("Bank.RSAPEMPublicKeyFile", "./pubkey.pem")
	viper.SetDefault("Bank.Cash", "1") // only 1 NEW support

	// user
	// viper.SetDefault("Bank.RSAPEMPrivateKeyFile", "./key.pem")
	// viper.SetDefault("Bank.RSAPEMPublicKeyFile", "./pubkey.pem")
}

func (cli *CLI) setupConfig() error {

	// var ret bool
	var err error

	cli.defaultConfig()

	viper.SetConfigName(defaultConfigFile)
	viper.AddConfigPath(".")
	cfgFile := cli.config
	if cfgFile != "" {
		if _, err = os.Stat(cfgFile); err == nil {
			viper.SetConfigFile(cfgFile)
			err = viper.ReadInConfig()
		} else {
			// The default configuration is enabled.
			// fmt.Println(err)
			err = nil
		}
	} else {
		// The default configuration is enabled.
		err = nil
	}

	if rpcURL := viper.GetString("rpcURL"); rpcURL != "" {
		cli.rpcURL = viper.GetString("rpcURL")
	}
	if walletPath := viper.GetString("walletPath"); walletPath != "" {
		cli.walletPath = viper.GetString("walletPath")
	}
	if log := viper.GetString("log"); log != "" {
		cli.logfile = viper.GetString("log")
	}

	unit := viper.GetString("unit")
	if stringInSlice(unit, DenominationList) {
		cli.tran.Unit = unit
	}

	return cli.setDefaultBank()
}

func (cli *CLI) setDefaultTransaction() error {
	if cli.tran == nil {
		cli.tran = new(Transaction)
	}
	fromStr := viper.GetString("Bank.Address")
	if common.IsHexAddress(fromStr) {
		cli.tran.From = common.HexToAddress(fromStr)
	}

	if password := viper.GetString("Bank.Password"); password != "" {
		cli.tran.Password = password
	}

	return nil
}

func (cli *CLI) setDefaultBank() error {
	if cli.bank == nil {
		cli.bank = new(Bank)
	}
	aStr := viper.GetString("Bank.Address")
	if common.IsHexAddress(aStr) {
		cli.bank.Address = common.HexToAddress(aStr)
	}

	if password := viper.GetString("Bank.Password"); password != "" {
		cli.bank.Password = password
	}

	// viper.SetDefault("Bank.RSAPEMPrivateKeyFile", "./key.pem")
	// viper.SetDefault("Bank.RSAPEMPublicKeyFile", "./pubkey.pem")
	if key := viper.GetString("Bank.RSAPEMPrivateKeyFile"); key != "" {
		cli.bank.RSAPEMPrivateKeyFile = key
	}
	if pubkey := viper.GetString("Bank.RSAPEMPublicKeyFile"); pubkey != "" {
		cli.bank.RSAPEMPublicKeyFile = pubkey
	}

	return nil
}
