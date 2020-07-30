package cli

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var (
	buildCommit string
	buildDate   string
)

// DefaultChainID default chain ID
var DefaultChainID = big.NewInt(16888)

// Transaction for send Transaction
type Transaction struct {
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Value     *big.Int       `json:"value"`
	Unit      string         `json:"unit"`
	Data      []byte         `json:"data"`
	Nonce     uint64         `json:"nonce"`
	GasPrice  *big.Int       `json:"gasPrice"`
	GasLimit  uint64         `json:"gas"`
	NetworkID *big.Int       `json:"networkID"`
	Password  string         `json:"password,omitempty"`
}

// Bank struct
type Bank struct {
	Address  common.Address `json:"address"`
	Password string         `json:"password,omitempty"`

	// rsa
	RSAPEMPrivateKeyFile string
	RSAPEMPublicKeyFile  string
}

// CLI represents a command-line interface. This class is
// not threadsafe.
type CLI struct {
	rootCmd    *cobra.Command
	version    string
	walletPath string
	rpcURL     string
	config     string
	testing    bool
	logfile    string

	client *ethclient.Client
	tran   *Transaction
	wallet *keystore.KeyStore

	bank *Bank
}

// NewCLI returns an initialized CLI
func NewCLI() *CLI {
	version := "v0.1.0"
	if buildCommit != "" {
		version = fmt.Sprintf("%s-%s", version, buildCommit)
	}
	if buildDate != "" {
		version = fmt.Sprintf("%s-%s", version, buildDate)
	}
	cli := &CLI{
		rootCmd:    nil,
		version:    version,
		walletPath: "",
		rpcURL:     "",
		testing:    false,
		config:     "",
		logfile:    "./error.log",
		tran:       new(Transaction),
		wallet:     nil,

		bank: new(Bank),
	}
	cli.tran.Unit = "NEW"
	cli.tran.Value = new(big.Int)
	cli.tran.GasPrice = big.NewInt(1)
	cli.tran.NetworkID = DefaultChainID
	cli.tran.GasLimit = 21000

	cli.buildRootCmd()
	return cli
}

// CopyCLI returns an copy  CLI
func CopyCLI(cli *CLI) *CLI {
	cpy := &CLI{
		rootCmd:    nil,
		walletPath: cli.walletPath,
		rpcURL:     cli.rpcURL,
		testing:    false,
		config:     cli.config,
		logfile:    cli.logfile,
		tran:       new(Transaction),
		wallet:     nil,
	}

	cpy.tran.From = cli.tran.From
	cpy.tran.Unit = cli.tran.Unit
	cpy.tran.Password = cli.tran.Password
	cpy.tran.Value = new(big.Int)
	cpy.tran.GasPrice = cli.tran.GasPrice
	cpy.tran.NetworkID = DefaultChainID
	cpy.tran.GasLimit = 21000

	cpy.buildRootCmd()
	return cpy
}

func (cli *CLI) resetConfig() error {
	// ok, go free itself
	cli.wallet = nil
	cli.tran = nil
	if cli.client != nil {
		cli.client.Close()
		cli.client = nil
	}

	return nil
}

// BuildClient BuildClient
func (cli *CLI) BuildClient() error {
	var err error
	if cli.client == nil {
		cli.client, err = ethclient.Dial(cli.rpcURL)
		if err != nil {
			return fmt.Errorf("Failed to connect to the NewChain client: %v", err)
		}
	}
	return nil
}

// Execute parses the command line and processes it.
func (cli *CLI) Execute() {
	cli.rootCmd.Execute()
}

// setup turns up the CLI environment, and gets called by Cobra before
// a command is executed.
func (cli *CLI) setup(cmd *cobra.Command, args []string) {
	err := cli.setupConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}
}

func (cli *CLI) help(cmd *cobra.Command, args []string) {
	fmt.Fprint(os.Stderr, cmd.UsageString())

	os.Exit(-1)

}

// TestCommand test command
func (cli *CLI) TestCommand(command string) string {
	cli.testing = true
	result := cli.Run(strings.Fields(command)...)
	cli.testing = false
	return result
}

// Run executes CLI with the given arguments. Used for testing. Not thread safe.
func (cli *CLI) Run(args ...string) string {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	cli.rootCmd.SetArgs(args)
	cli.rootCmd.Execute()
	cli.buildRootCmd()

	w.Close()

	os.Stdout = oldStdout

	var stdOut bytes.Buffer
	io.Copy(&stdOut, r)
	return stdOut.String()
}

// Embeddable returns a CLI that you can embed into your own Go programs. This
// is not thread-safe.
func (cli *CLI) Embeddable() *CLI {
	cli.testing = true
	return cli
}
