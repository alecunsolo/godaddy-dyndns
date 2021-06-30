package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// used for flags
	domain    string
	hostname  string
	apiKey    string
	keySecret string

	rootCmd = &cobra.Command{
		Use:   "godaddy-dyndns",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind Cobra and Viper
			return initializeConfig(cmd)
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			extIP, err := retrieveExternalIP()
			if err != nil {
				log.Fatalf("Failed to retrieve external IP. %s", err)
			}
			dnsIP, err := currentIP()
			if err != nil {
				log.Fatalf("Failed to retrieve current A record IP. %s", err)
			}
			err = updateIP(dnsIP, extIP)
			if err != nil {
				log.Fatalf("Failed to update DNS record. %s", err)
			}
		},
	}
)

// Better Viper/Cobra integration.
// See https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("godaddy")

	v.SetEnvPrefix("GD_DYNDNS")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	// Apply the viper config value to the flag when the flag is not set and viper has a value
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVar(&domain, "domain", "", "Domain to update")
	rootCmd.Flags().StringVar(&hostname, "hostname", "", "Hostname to update")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "API key")
	rootCmd.Flags().StringVar(&keySecret, "key-secret", "", "Secret API key")

	rootCmd.MarkFlagRequired("domain")
	rootCmd.MarkFlagRequired("hostname")
	rootCmd.MarkFlagRequired("api-key")
	rootCmd.MarkFlagRequired("key-secret")
}
