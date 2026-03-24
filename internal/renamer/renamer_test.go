package renamer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var testCameraIDs = []string{"max2-c01", "hb12-c01", "hb12-c02"}

func createTestFiles(t *testing.T, dir string, files []string) {
	t.Helper()
	for _, f := range files {
		path := filepath.Join(dir, f)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}
}

func expectedName(cameraID string, count int, ext string) string {
	date := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s-%03d%s", date, cameraID, count, ext)
}

func assertFiles(t *testing.T, dir string, expected []string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read directory %s: %v", dir, err)
	}

	var actual []string
	for _, entry := range entries {
		if !entry.IsDir() {
			actual = append(actual, entry.Name())
		}
	}

	if len(actual) != len(expected) {
		t.Errorf("expected %d files, got %d", len(expected), len(actual))
		return
	}

	for i, name := range actual {
		if name != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], name)
		}
	}
}

func TestRenameFlatDirectory(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{
		"GS010001.360",
		"GS010002.360",
		"GS010003.360",
	})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	cameraID := filepath.Base(dir)
	assertFiles(t, dir, []string{
		expectedName(cameraID, 1, ".360"),
		expectedName(cameraID, 2, ".360"),
		expectedName(cameraID, 3, ".360"),
	})
}

func TestRenameWithSubdirectories(t *testing.T) {
	dir := t.TempDir()

	cam1 := filepath.Join(dir, testCameraIDs[0])
	cam2 := filepath.Join(dir, testCameraIDs[1])
	cam3 := filepath.Join(dir, testCameraIDs[2])
	os.Mkdir(cam1, 0755)
	os.Mkdir(cam2, 0755)
	os.Mkdir(cam3, 0755)

	createTestFiles(t, cam1, []string{"GS010001.360", "GS010002.360"})
	createTestFiles(t, cam2, []string{"GX010001.MP4", "GX010002.MP4"})
	createTestFiles(t, cam3, []string{"GX010001.MP4", "GX010002.MP4"})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	assertFiles(t, cam1, []string{
		expectedName(testCameraIDs[0], 1, ".360"),
		expectedName(testCameraIDs[0], 2, ".360"),
	})

	assertFiles(t, cam2, []string{
		expectedName(testCameraIDs[1], 1, ".MP4"),
		expectedName(testCameraIDs[1], 2, ".MP4"),
	})

	assertFiles(t, cam3, []string{
		expectedName(testCameraIDs[2], 1, ".MP4"),
		expectedName(testCameraIDs[2], 2, ".MP4"),
	})
}

func TestRenameSkipsUnsupportedExtensions(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{
		"GS010001.360", // supported - uppercase sorts before lowercase
		"clip.mp4",     // supported
		"photo.jpg",    // supported
		"data.xml",     // unsupported
		"readme.txt",   // unsupported
	})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	// date-prefixed names sort before "data.xml" and "readme.txt"
	cameraID := filepath.Base(dir)
	assertFiles(t, dir, []string{
		expectedName(cameraID, 1, ".360"),
		expectedName(cameraID, 2, ".mp4"),
		expectedName(cameraID, 3, ".jpg"),
		"data.xml",
		"readme.txt",
	})
}

func TestRenameSkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"GS010001.360", ".DS_Store", ".hidden"})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	cameraID := filepath.Base(dir)
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.Name() == ".DS_Store" || entry.Name() == ".hidden" {
			continue
		}
		if entry.Name() != expectedName(cameraID, 1, ".360") {
			t.Errorf("expected %s, got %s", expectedName(cameraID, 1, ".360"), entry.Name())
		}
	}
}

func TestRenameEmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed on empty directory: %v", err)
	}
}

func TestRenameSubdirectoryContainingDirectory(t *testing.T) {
	dir := t.TempDir()

	cam1 := filepath.Join(dir, testCameraIDs[0])
	nested := filepath.Join(cam1, "nested")
	os.Mkdir(cam1, 0755)
	os.Mkdir(nested, 0755)

	createTestFiles(t, cam1, []string{"GS010001.360", "GS010002.360"})
	createTestFiles(t, nested, []string{"GS010001.360"})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	assertFiles(t, cam1, []string{
		expectedName(testCameraIDs[0], 1, ".360"),
		expectedName(testCameraIDs[0], 2, ".360"),
	})

	nestedEntries, _ := os.ReadDir(nested)
	if nestedEntries[0].Name() != "GS010001.360" {
		t.Errorf("nested file should not have been renamed, got %s", nestedEntries[0].Name())
	}
}

func TestRenameSubdirectoriesIgnoresTopLevelFiles(t *testing.T) {
	dir := t.TempDir()

	cam1 := filepath.Join(dir, testCameraIDs[0])
	os.Mkdir(cam1, 0755)

	createTestFiles(t, dir, []string{"GS010001.360"})
	createTestFiles(t, cam1, []string{"GS010002.360"})

	if err := Rename(dir, Flags{}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	assertFiles(t, cam1, []string{
		expectedName(testCameraIDs[0], 1, ".360"),
	})

	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != "GS010001.360" {
			t.Errorf("top level file should not have been renamed, got %s", entry.Name())
		}
	}
}

func TestRenameCameraOverride(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"GS010001.360", "GS010002.360"})

	if err := Rename(dir, Flags{Camera: "my-cam"}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	assertFiles(t, dir, []string{
		expectedName("my-cam", 1, ".360"),
		expectedName("my-cam", 2, ".360"),
	})
}

func TestRenameCameraErrorWithSubdirectories(t *testing.T) {
	dir := t.TempDir()
	cam1 := filepath.Join(dir, testCameraIDs[0])
	os.Mkdir(cam1, 0755)
	createTestFiles(t, cam1, []string{"GS010001.360"})

	err := Rename(dir, Flags{Camera: "my-cam"})
	if err == nil {
		t.Fatal("expected error when --camera is used with subdirectories, got nil")
	}
}

func TestRenameDateOverride(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"GS010001.360", "GS010002.360"})

	if err := Rename(dir, Flags{Date: "2020-01-15"}); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	cameraID := filepath.Base(dir)
	assertFiles(t, dir, []string{
		fmt.Sprintf("2020-01-15-%s-001.360", cameraID),
		fmt.Sprintf("2020-01-15-%s-002.360", cameraID),
	})
}

func TestRenameDryRun(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"GS010001.360", "GS010002.360"})

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err := Rename(dir, Flags{DryRun: true})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "Dry run mode, no files will be renamed") {
		t.Errorf("expected dry run header, got:\n%s", output)
	}

	cameraID := filepath.Base(dir)
	if !strings.Contains(output, fmt.Sprintf("GS010001.360 -> %s", expectedName(cameraID, 1, ".360"))) {
		t.Errorf("expected dry run line for file 1, got:\n%s", output)
	}
	if !strings.Contains(output, fmt.Sprintf("GS010002.360 -> %s", expectedName(cameraID, 2, ".360"))) {
		t.Errorf("expected dry run line for file 2, got:\n%s", output)
	}

	// files should not have been renamed
	assertFiles(t, dir, []string{"GS010001.360", "GS010002.360"})
}

func TestRenameVerbose(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"GS010001.360", "GS010002.360"})

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err := Rename(dir, Flags{Verbose: true})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	out, _ := io.ReadAll(r)
	output := string(out)

	cameraID := filepath.Base(dir)
	if !strings.Contains(output, fmt.Sprintf("GS010001.360 -> %s", expectedName(cameraID, 1, ".360"))) {
		t.Errorf("expected verbose line for file 1, got:\n%s", output)
	}
	if !strings.Contains(output, fmt.Sprintf("GS010002.360 -> %s", expectedName(cameraID, 2, ".360"))) {
		t.Errorf("expected verbose line for file 2, got:\n%s", output)
	}

	assertFiles(t, dir, []string{
		expectedName(cameraID, 1, ".360"),
		expectedName(cameraID, 2, ".360"),
	})
}
