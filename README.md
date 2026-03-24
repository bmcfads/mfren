# mfren - Media File Renamer

A CLI tool for renaming video and photo files after a day of shooting. Built as a learning project to explore Go.

## Overview

`mfren` walks a directory and renames media files into a consistent, date-stamped format. It supports camera IDs, dry-run previews, and both flat and nested directory structures.

## Installation
 
```bash
go install github.com/bmcfads/mfren@latest
```
 
Or build from source:
 
```bash
git clone https://github.com/bmcfads/mfren
cd mfren
go install .
```

## Usage
 
```
Usage:
  mfren <directory> [flags]

Flags:
  -c, --camera string     camera ID override
  -d, --date string       date override (YYYY-MM-DD)
  -n, --dry-run           print renames without applying them
  -h, --help              help for mfren
      --list-extensions   print supported file extensions
      --verbose           print each rename as it happens
  -v, --version           version for mfren
```

### Output Format
 
```
<YYYY-MM-DD>-<camera-id>-<NNN>.<ext>
```

> [!NOTE]
> The camera ID is a free-form string and may itself contain dashes.

For example, with a camera ID of `max2-c01`:
 
```
2026-03-21-max2-c01-001.360
2026-03-21-max2-c01-002.360
2026-03-21-max2-c01-003.360
```
 
## Examples

Rename files in a flat directory, using the directory name as the camera ID:

```bash
mfren ./media
```

Override the camera ID:

```bash
mfren ./media --camera max2-c01
```

Override the date:

```bash
mfren ./media --date 2026-03-21
```

Rename files across multiple camera subdirectories, using each subdirectory name as the camera ID:

```bash
mfren ./media
# media/
#   max2-c01/  -> 2026-03-21-max2-c01-001.360 ...
#   hb12-c01/  -> 2026-03-21-hb12-c01-001.mp4 ...
```

Preview renames without applying them:

```bash
mfren ./media --dry-run
```

Print supported extensions:

```bash
mfren --list-extensions
```

## Generating mock media files

> [!WARNING]
> The script deletes and recreates the `media` directory at the destination on each run.

A script is provided to generate mock media files for manual testing:

```bash
./scripts/gen-test-files.sh [destination]
```

If no destination is provided, files are created under `/tmp/media`. The script will prompt you to select a scenario:

**Scenario 1 — Flat directory**

Creates a single directory with 10 `.360` files and a `shoot-notes.txt`:

```
media/
  GS01000001.360
  ...
  GS01000010.360
  shoot-notes.txt
```

**Scenario 2 — Subdirectories with top level files**

Creates a top-level directory with files plus three camera subdirectories, each with a different file type:

```
media/
  GS01000001.360 ... GS01000010.360
  shoot-notes.txt
  cam-id-01/
    GS01000001.360 ... GS01000010.360
    shoot-notes.txt
  cam-id-02/
    GX01000001.mp4 ... GX01000010.mp4
    shoot-notes.txt
  cam-id-03/
    GOPR000001.jpg ... GOPR000010.jpg
    shoot-notes.txt
```

## Behaviour

### Directory structure

- If no subdirectories exist, files in the target directory are renamed directly.
- If the target directory contains subdirectories, files within each subdirectory are renamed independently using the subdirectory name as the camera ID. Top level files are ignored.
- Only one level of subdirectories is searched — no recursion.
- File count resets to `001` per subdirectory.

### Camera ID

- If subdirectories are present and `--camera` is provided, `mfren` will exit with an error. Target a single directory or rely on subdirectory names as camera IDs instead.
- If subdirectories are present and `--camera` is not provided, each subdirectory name is used as the camera ID.
- If no subdirectories are present and `--camera` is provided, the provided value is used as the camera ID.
- If no subdirectories are present and `--camera` is not provided, the target directory name is used as the camera ID.

### Files

- Only media files are renamed. Use `--list-extensions` to see supported extensions.
- Unsupported file extensions are skipped silently.
- Hidden files (starting with `.`) are skipped.

### Safety

- A confirmation prompt is shown before renaming since the operation is destructive and cannot be undone.
- Use `--dry-run` to verify the expected output before committing.

### Output

- Silent on success by default — no news is good news.
- The current date is used by default. Use `--date` to override.
- Use `--verbose` to print each rename as it happens.
- Use `--dry-run` to preview renames without applying them. Skips the confirmation prompt.

## About
 
`mfren` is a learning project built to get comfortable with Go, Cobra, and CLI tool design patterns. It solves a real problem I had (renaming media files consistently after a shoot) while serving as a foundation for more complex Go projects.

## License

MIT — see [LICENSE](LICENSE) for details.
