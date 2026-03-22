package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "mfren <directory>",
	Short:   "Rename media files after shoot",
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}
