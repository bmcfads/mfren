package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Flags struct {
	Camera string
}

// Rename walks dir and renames all media files to the format
// <YYYY-MM-DD>-<camera-id>-<NNN>.<ext>, where the camera ID is the
// name of the directory containing the file. If dir contains
// subdirectories, files within each are renamed independently with
// the file count resetting to 001 per subdirectory. Only one level
// of subdirectories is searched. Hidden files are skipped.
func Rename(dir string, flags Flags) error {
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
		if flags.Camera != "" {
			return fmt.Errorf("--camera cannot be used when subdirectories are present")
		}
		for _, subdir := range subdirs {
			if err := renameFiles(filepath.Join(dir, subdir), subdir); err != nil {
				return err
			}
		}
	} else {
		cameraID := filepath.Base(dir)
		if flags.Camera != "" {
			cameraID = flags.Camera
		}
		if err := renameFiles(dir, cameraID); err != nil {
			return err
		}
	}

	return nil
}

func renameFiles(dir, cameraID string) error {
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
		files = append(files, entry.Name())
	}

	sort.Strings(files)

	date := time.Now().Format("2006-01-02")

	for i, file := range files {
		ext := filepath.Ext(file)
		newName := fmt.Sprintf("%s-%s-%03d%s", date, cameraID, i+1, ext)
		oldPath := filepath.Join(dir, file)
		newPath := filepath.Join(dir, newName)

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("cannot rename %s: %w", file, err)
		}
	}

	return nil
}
