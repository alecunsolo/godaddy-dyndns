package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "godaddy-dyndns",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := extIP()
			if err != nil {
				log.Fatalf("Failed to retrieve external IP. %s", err)
			}
			log.Printf("Current external ip: %s", ip)
			curIP, err := currentIP()
			if err != nil {
				log.Fatalf("Failed to retrieve current A record IP. %s", err)
			}
			log.Printf("Current ip: %s", curIP)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/godaddy/config.yaml)")

	rootCmd.Flags().String("domain", "", "Domain to update")
	rootCmd.Flags().String("hostname", "", "Hostname to update")
	rootCmd.Flags().String("api-key", "", "API key")
	rootCmd.Flags().String("secret-key", "", "Secret API key")
	viper.BindPFlag("domain", rootCmd.Flags().Lookup("domain"))
	viper.BindPFlag("hostname", rootCmd.Flags().Lookup("hostname"))
	viper.BindPFlag("api-key", rootCmd.Flags().Lookup("api"))
	viper.BindPFlag("secret-key", rootCmd.Flags().Lookup("key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(path.Join(configDir, "godaddy"))
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("GD_DYNDNS")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
