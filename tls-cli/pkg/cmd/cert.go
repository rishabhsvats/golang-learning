package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/rishabhsvats/tls-cli/pkg/cert"
	"github.com/spf13/cobra"
)

var certKeyPath string
var certPath string
var certName string

func init() {
	createCmd.AddCommand(certCreateCmd)
	certCreateCmd.Flags().StringVarP(&certKeyPath, "key-out", "k", "cert.key", "destination path for cert key")
	certCreateCmd.Flags().StringVarP(&certPath, "cert-out", "o", "cert.crt", "destination path for cert")
	certCreateCmd.Flags().StringVarP(&certName, "name", "n", "", "name of the certificate in the config file")
	certCreateCmd.Flags().StringVar(&caKey, "ca-key", "ca.key", "ca key to sign certificate")
	certCreateCmd.Flags().StringVar(&caCert, "ca-cert", "ca.crt", "ca cert for certificate")
	certCreateCmd.MarkFlagRequired("ca-key")
	certCreateCmd.MarkFlagRequired("ca-cert")
	certCreateCmd.MarkFlagRequired("name")

}

var certCreateCmd = &cobra.Command{
	Use:   "cert",
	Short: "cert commands",
	Long:  `commands to create the certificates`,
	Run: func(cmd *cobra.Command, args []string) {
		caKeyBytes, err := ioutil.ReadFile(caKey)
		if err != nil {
			fmt.Printf("CA key read error : %s\n", err)
			return
		}
		caCertBytes, err := ioutil.ReadFile(caCert)
		if err != nil {
			fmt.Printf("CA cert read error : %s\n", err)
			return
		}
		certConfig, ok := config.Cert[certName]
		if !ok {
			fmt.Printf("could not find certificate in config\n")
			return
		}
		err = cert.CreateCert(certConfig, caKeyBytes, caCertBytes, certKeyPath, certPath)
		if err != nil {
			fmt.Printf("create cert error: %s\n", err)
			return
		}
		fmt.Printf("Cert created. Key: %s with length %s\n", certKeyPath, certPath)
	},
}
