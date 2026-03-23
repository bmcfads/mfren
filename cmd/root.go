package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bmcfads/mfren/internal/renamer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mfren <directory>",
	Short:   "Rename media files after shoot",
	Version: "0.1.0",
	Args:    cobra.ExactArgs(1),
	RunE:    run,
}

var flagCamera string
var flagDate string

func init() {
	rootCmd.Flags().StringVarP(&flagCamera, "camera", "c", "", "camera ID override")
	rootCmd.Flags().StringVarP(&flagDate, "date", "d", "", "date override (YYYY-MM-DD)")
}

func Execute() error {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	dir, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("cannot resolve directory path: %w", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("cannot access directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	if flagDate != "" {
		if _, err := time.Parse("2006-01-02", flagDate); err != nil {
			return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
		}
	}

	return renamer.Rename(dir, renamer.Flags{
		Camera: flagCamera,
		Date:   flagDate,
	})
}
