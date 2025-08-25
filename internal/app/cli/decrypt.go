package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"golang/tutorial/todo/internal/crypto"
)

func newDecryptCmd() *cobra.Command {
	var (
		inFile  string
		outFile string
		pass    string
	)
	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a file",
		Long:  "Decrypt a file previously encrypted with the encrypt command",
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

			if err := crypto.DecryptFile(inFile, outFile, pass); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("File decrypted successfully")
		},
	}
	cmd.Flags().StringVarP(&inFile, "input", "i", "", "Input encrypted file to decrypt")
	cmd.Flags().StringVarP(&outFile, "output", "o", "", "Output file for the decrypted plaintext")
	cmd.Flags().StringVarP(&pass, "password", "p", "", "Password used to encrypt the file")
	return cmd
}
