package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "sdrparser",
	Short: "SDR Parser - tool for processing Software Defined Radio signals",
	Long: `SDR Parser is a command line tool for processing Software Defined Radio signals.
It supports various operations like filtering, modulation/demodulation, and signal analysis.`,
}

// Execute starts the CLI application
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// Add commands
	rootCmd.AddCommand(getGenerateCmd())
	rootCmd.AddCommand(getDemodCmd())
	rootCmd.AddCommand(getFilterCmd())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		// Using configuration file
	}
}
