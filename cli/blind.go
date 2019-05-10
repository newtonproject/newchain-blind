package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"gitlab.newtonproject.org/yangchenzhong/NewChainBlind/blind"
	"io/ioutil"
	"time"
)

func (cli *CLI) buildBlindCmd() *cobra.Command {
	faucetCmd := &cobra.Command{
		Use:   "blind <pubkey>",
		Short: "Blind 1NEW for address",
		Args:  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			pub, err := blind.LoadPEMPublicKeyFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			c := 32
			data := make([]byte, c)
			_, err = rand.Read(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			cm := CashMessage{
				Cash: 1,
				Data: data,
			}
			cmb, err := rlp.EncodeToBytes(&cm)
			if err != nil {
				fmt.Println(err)
				return
			}

			blinded, unblinder, err := blind.BlindCash(cmb, pub)
			if err != nil {
				fmt.Println(err)
				return
			}

			n := time.Now().Format("20060102150405")
			dataFile := fmt.Sprintf("%v.data", n)
			blindedFile := fmt.Sprintf("%v.blinded", n)
			unblinderFile := fmt.Sprintf("%v.unblinder", n)

			if err := ioutil.WriteFile(dataFile, []byte(hex.EncodeToString(cmb)), 0644); err != nil {
				fmt.Println(err)
				return
			}
			if err := ioutil.WriteFile(blindedFile, []byte(hex.EncodeToString(blinded)), 0644); err != nil {
				fmt.Println(err)
				return
			}
			if err := ioutil.WriteFile(unblinderFile, []byte(hex.EncodeToString(unblinder)), 0644); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Data: ", dataFile)
			fmt.Println("Blinded: ", blindedFile)
			fmt.Println("Unblinder: ", unblinderFile)
			fmt.Println("Send blinded file to bank")
		},
	}

	return faucetCmd
}

type CashMessage struct {
	Cash uint64 `json:"cash"    gencodec:"required"`
	Data []byte `json:"data"    gencodec:"required"`
}
