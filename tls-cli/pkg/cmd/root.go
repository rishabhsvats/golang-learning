package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rishabhsvats/tls-cli/pkg/cert"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CACert *cert.CACert          `yaml:"caCert"`
	Cert   map[string]*cert.Cert `yaml:"certs"`
}

var cfgFilePath string
var config Config

var rootCmd = &cobra.Command{
	Use:   "tls",
	Short: "tls is a command line tool for tls",
	Long:  `tls is a command line tool to create certificate, cacertificate and can also be extended`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", "", "config file($default is tls.yaml)")
}

func initConfig() {

	if cfgFilePath == "" {
		cfgFilePath = "tls.yaml"
	}
	cfgFileBytes, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Printf("error while reading config file: %s", err)
		return
	}
	err = yaml.Unmarshal(cfgFileBytes, &config)
	if err != nil {
		fmt.Printf("error while parsing the config file: %s", err)
		return
	}
	fmt.Printf("config file parsed: %v\n", config)
}
