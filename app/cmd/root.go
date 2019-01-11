package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	redisConString      = "redis-connection-string"
	redisConStringUsage = "The Redis TCP connection string to connect to queue backend"
	queueTag            = "redis-connection-tag"
	queueTagUsage       = "The queue tag used to identify queues"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "prometheus-exporter-golang",
		Short: "Just a simple action producer and consumer",
		Long:  `This is a demo just to show how can integrate Prometheus custom metrics in your GoLang application.
Use 'produce' command to create faked tasks, then use 'consume' to consume them and create metrics.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			var a string
			a, err = cmd.Flags().GetString(action)
			if a == "" {
				return fmt.Errorf("missing action")
			}
			if err != nil {
				return err
			}

			switch a {
			case create:
			case read:
			case update:
			case delete:
				break
			default:
				err = fmt.Errorf("action %s is not recognized", a)
			}
			return err
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String(redisConString, "", redisConStringUsage)
	rootCmd.PersistentFlags().String(queueTag, "", queueTagUsage)

	if err := rootCmd.MarkPersistentFlagRequired(redisConString); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired(queueTag); err != nil {
		panic(err)
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
}
