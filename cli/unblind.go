package cli

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/cryptoballot/rsablind"
	"github.com/newtonproject/newchain-blind/blind"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildUnblindCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:                   "unblind <pubkey> <file.blinded.sig> <file.unblinder>",
		Short:                 "Unblind 1NEW for address",
		Args:                  cobra.MinimumNArgs(3),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			pub, err := blind.LoadPEMPublicKeyFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			sig, err := readHexFile(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}

			unblinder, err := readHexFile(args[2])
			if err != nil {
				fmt.Println(err)
				return
			}

			unblindSig := rsablind.Unblind(pub, sig, unblinder)
			if len(unblindSig) == 0 {
				fmt.Println("Unblind Sig error")
				return
			}

			dataFile := fmt.Sprintf("%v.sig", args[2])
			if err := ioutil.WriteFile(dataFile, []byte(hex.EncodeToString(unblindSig)), 0644); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("UnblinderSig: ", dataFile)
			fmt.Println("Send unblinder sig file and cash data to others")
		},
	}

	return faucetCmd
}

func readHexFile(name string) ([]byte, error) {
	blindedHexByte, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(string(blindedHexByte))
}
