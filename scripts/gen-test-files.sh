#!/bin/bash

# generate mock media files for mfren manual testing
# usage: ./scripts/gen-test-files.sh [destination]
# if no destination is provided, defaults to /tmp

DEST=${1:-/tmp}
MEDIA_DIR="$DEST/media"

# clean up and recreate media directory
rm -rf "$MEDIA_DIR"
mkdir -p "$MEDIA_DIR"

# scenario selection
echo "Select a scenario:"
echo "  1) Flat directory"
echo "  2) Subdirectories with top level files"
read -r -p "Enter 1 or 2: " SCENARIO

case $SCENARIO in
  1)
    echo "Creating flat directory with 10 .360 files..."
    for i in $(seq -f "%06g" 1 10); do
      touch "$MEDIA_DIR/GS01${i}.360"
    done
    touch "$MEDIA_DIR/shoot-notes.txt"
    echo "Done: $MEDIA_DIR"
    ;;
  2)
    echo "Creating subdirectory structure with top level files..."

    # top level files
    for i in $(seq -f "%06g" 1 10); do
      touch "$MEDIA_DIR/GS01${i}.360"
    done
    touch "$MEDIA_DIR/shoot-notes.txt"

    # cam-id-01 subdirectory
    mkdir -p "$MEDIA_DIR/cam-id-01"
    for i in $(seq -f "%06g" 1 10); do
      touch "$MEDIA_DIR/cam-id-01/GS01${i}.360"
    done
    touch "$MEDIA_DIR/cam-id-01/shoot-notes.txt"

    # cam-id-02 subdirectory
    mkdir -p "$MEDIA_DIR/cam-id-02"
    for i in $(seq -f "%06g" 1 10); do
      touch "$MEDIA_DIR/cam-id-02/GX01${i}.mp4"
    done
    touch "$MEDIA_DIR/cam-id-02/shoot-notes.txt"

    # cam-id-03 subdirectory
    mkdir -p "$MEDIA_DIR/cam-id-03"
    for i in $(seq -f "%06g" 1 10); do
      touch "$MEDIA_DIR/cam-id-03/GOPR${i}.jpg"
    done
    touch "$MEDIA_DIR/cam-id-03/shoot-notes.txt"

    echo "Done: $MEDIA_DIR"
    ;;
  *)
    echo "Error: invalid scenario, enter 1 or 2"
    exit 1
    ;;
esac
