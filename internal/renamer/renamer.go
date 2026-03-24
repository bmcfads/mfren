package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Flags struct {
	Camera  string
	Date    string
	DryRun  bool
	Verbose bool
}

var Extensions360 = []string{
	".360", ".insp", ".insv",
}

var ExtensionsPhoto = []string{
	".arw", ".cr3", ".dng", ".gpr", ".jpeg", ".jpg", ".png", ".raw",
}

var ExtensionsVideo = []string{
	".mov", ".mp4",
}

// Rename walks dir and renames all media files to the format
// <YYYY-MM-DD>-<camera-id>-<NNN>.<ext>, where the camera ID is the
// name of the directory containing the file. If dir contains
// subdirectories, files within each are renamed independently with
// the file count resetting to 001 per subdirectory. Only one level
// of subdirectories is searched. Hidden files are skipped.
func Rename(dir string, flags Flags) error {
	date := flags.Date
	dryRun := flags.DryRun
	verbose := flags.Verbose

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if dryRun {
		fmt.Println("Dry run mode, no files will be renamed")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read directory: %w", err)
	}

	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}

	if len(subdirs) > 0 {
		for _, subdir := range subdirs {
			if err := renameFiles(filepath.Join(dir, subdir), date, subdir, dryRun, verbose); err != nil {
				return err
			}
		}
	} else {
		cameraID := filepath.Base(dir)
		if flags.Camera != "" {
			cameraID = flags.Camera
		}
		if err := renameFiles(dir, date, cameraID, dryRun, verbose); err != nil {
			return err
		}
	}

	return nil
}

func isSupportedExt(ext string) bool {
	ext = strings.ToLower(ext)
	for _, e := range append(append(Extensions360, ExtensionsPhoto...), ExtensionsVideo...) {
		if ext == e {
			return true
		}
	}
	return false
}

func renameFiles(dir, date, cameraID string, dryRun, verbose bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if !isSupportedExt(filepath.Ext(entry.Name())) {
			continue
		}
		files = append(files, entry.Name())
	}

	for i, file := range files {
		ext := filepath.Ext(file)
		newName := fmt.Sprintf("%s-%s-%03d%s", date, cameraID, i+1, ext)
		oldPath := filepath.Join(dir, file)
		newPath := filepath.Join(dir, newName)

		if dryRun || verbose {
			fmt.Printf("%s -> %s\n", file, newName)
		}

		if !dryRun {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("cannot rename %s: %w", file, err)
			}
		}
	}

	return nil
}
