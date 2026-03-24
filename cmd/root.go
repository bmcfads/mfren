package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmcfads/mfren/internal/renamer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mfren <directory>",
	Short:   "Rename media files after shoot",
	Version: "0.1.0",
	Args:    cobra.RangeArgs(0, 1),
	RunE:    run,
}

var flagCamera string
var flagDate string
var flagDryRun bool
var flagListExtensions bool
var flagVerbose bool

func init() {
	rootCmd.Flags().StringVarP(&flagCamera, "camera", "c", "", "camera ID override")
	rootCmd.Flags().StringVarP(&flagDate, "date", "d", "", "date override (YYYY-MM-DD)")
	rootCmd.Flags().BoolVarP(&flagDryRun, "dry-run", "n", false, "print renames without applying them")
	rootCmd.Flags().BoolVar(&flagListExtensions, "list-extensions", false, "print supported file extensions")
	rootCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "print each rename as it happens")
}

func Execute() error {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("accepts 1 arg(s), received 0")
	}

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

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read directory: %w", err)
	}

	hasSubdirs := false
	for _, entry := range entries {
		if entry.IsDir() {
			hasSubdirs = true
			break
		}
	}

	// informational flags - print and exit regardless of other args or flags
	if flagListExtensions {
		fmt.Printf("360:   [%s]\n", strings.Join(renamer.Extensions360, ", "))
		fmt.Printf("Photo: [%s]\n", strings.Join(renamer.ExtensionsPhoto, ", "))
		fmt.Printf("Video: [%s]\n", strings.Join(renamer.ExtensionsVideo, ", "))
		return nil
	}

	// flag validation
	if flagCamera != "" && hasSubdirs {
		return fmt.Errorf("--camera cannot be used when subdirectories are present")
	}

	if flagDate != "" {
		if _, err := time.Parse("2006-01-02", flagDate); err != nil {
			return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
		}
	}

	if !flagDryRun {
		fmt.Printf("\nDirectory: %s\n", dir)
		fmt.Println("Warning: renaming files is destructive and cannot be undone.")
		fmt.Print("Proceed? [y/N]: ")
		answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			return nil
		}
	}

	return renamer.Rename(dir, renamer.Flags{
		Camera:  flagCamera,
		Date:    flagDate,
		DryRun:  flagDryRun,
		Verbose: flagVerbose,
	})
}
