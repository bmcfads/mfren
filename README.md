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

## About
 
`mfren` is a learning project built to get comfortable with Go, Cobra, and CLI tool design patterns. It solves a real problem I had (renaming media files consistently after a shoot) while serving as a foundation for more complex Go projects.
