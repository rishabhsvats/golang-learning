package cmd

import (
	"fmt"

	"github.com/rishabhsvats/tls-cli/pkg/key"
	"github.com/spf13/cobra"
)

var keyOut string
var keyLength int

func init() {
	createCmd.AddCommand(keyCreateCmd)
	keyCreateCmd.Flags().StringVarP(&keyOut, "key-out", "k", "key.pem", "destination path for key")
	keyCreateCmd.Flags().IntVarP(&keyLength, "key-length", "l", 4096, "key length")

}

var keyCreateCmd = &cobra.Command{
	Use:   "key",
	Short: "key commands",
	Long:  `commands to create keys`,
	Run: func(cmd *cobra.Command, args []string) {
		err := key.CreateRSAPrivateKeyAndSave(keyOut, 4096)
		if err != nil {
			fmt.Printf("create key error: %s\n", err)
			return
		}
		fmt.Printf("key created %s with length %d\n", keyOut, keyLength)
	},
}
