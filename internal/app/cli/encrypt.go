package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/crypto"
)

func newEncryptCmd() *cobra.Command {
	var (
		inFile  string
		outFile string
		pass    string
	)
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file",
		Long:  "Encrypt a file using a password (scrypt + AES-GCM)",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := requiredInput(inFile, "Input file"); err != nil {
				fmt.Println(err)
				return
			}
			if err := requiredInput(outFile, "Output file"); err != nil {
				fmt.Println(err)
				return
			}
			if err := requiredInput(pass, "Password"); err != nil {
				fmt.Println(err)
				return
			}

			if err := crypto.EncryptFile(inFile, outFile, pass); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("File encrypted successfully")
		},
	}
	cmd.Flags().StringVarP(&inFile, "input", "i", "", "Input file to encrypt")
	cmd.Flags().StringVarP(&outFile, "output", "o", "", "Output file for the encrypted content")
	cmd.Flags().StringVarP(&pass, "password", "p", "", "Password for encryption")
	return cmd
}
