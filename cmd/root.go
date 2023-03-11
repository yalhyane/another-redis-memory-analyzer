package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arma",
	Short: "Another Redis Memory Analyzer",
	Long: `
arma - Yet Another Redis Memory Analyzer - A Command Line Tool for Analyzing Redis Memory

The arma command line tool is designed to help Redis developers understand how DB and keys patterns are using memory in Redis instance.
This tool uses the Redis MEMORY USAGE command to retrieve the number of bytes that a key and its value are taking in RAM, and then group them by key pattern and by db.

Examples:
  To analyze memory usage of a Redis instance running on the default port on the local machine:
    arma analyze
  
  To analyze memory usage of a Redis instance running on a remote server with a custom port and password:
    arma analyze --host example.com --port 6380 --ask-password

  To analyze memory usage of a Redis instance and pattern keys by custom delimiter and a minimum level:
    arma analyze --delimiter "-" --level 2

  To analyze memory usage of a Redis instance on a specific database with larger scan range and only show big keys:
    arma analyze --db 3 --scan 10000 --min-size 1MB

Note: This tool is intended to provide insights into memory usage of group of keys in Redis, but should not be used as a substitute for proper memory management practices. It is important to regularly monitor your Redis instance and take appropriate actions to optimize memory usage, such as configuring eviction policies or increasing memory capacity.

`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
