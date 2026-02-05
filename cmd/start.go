package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgGile string
var rootCmd = &cobra.Command{
	Use:   "sun-pos-payment-service",
	Short: "Sun POS Payment Service Application",
	Long:  `Sun POS Payment Service is a microservice for handling payment processing in the Sun POS system.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(startCmd, args)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgGile, "config", "", "config file (default is .env)")
	rootCmd.PersistentFlags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgGile != "" {
		viper.SetConfigFile(cfgGile)
	} else {
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file", viper.ConfigFileUsed())
	}
}
