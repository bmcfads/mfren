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
 
```bash
mfren <directory> [options]
```

### Options
 
| Option | Shorthand | Default | Description |
|--------|-----------|---------|-------------|
| `--help` | | | Show help |
| `--version` | | | Show version |

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
 
## About
 
`mfren` is a learning project built to get comfortable with Go, Cobra, and CLI tool design patterns. It solves a real problem I had (renaming media files consistently after a shoot) while serving as a foundation for more complex Go projects.
