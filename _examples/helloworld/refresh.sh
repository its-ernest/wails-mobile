#!/bin/bash
set -e

# Get locations relative to this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Define straight, predictable paths
OUTPUT_PATH="$SCRIPT_DIR/native/android/app/libs"

# Clean old artifacts
mkdir -p "$OUTPUT_PATH"
rm -f "$OUTPUT_PATH"/*.aar "$OUTPUT_PATH"/*-sources.jar

# Jump into the local package directory so gomobile reads the local files
cd "$SCRIPT_DIR"

echo "Building wailsmobile.aar..."
gomobile bind -target="android/arm64" -androidapi="21" -javapkg="wailsmobile" -o "$OUTPUT_PATH/wailsmobile.aar" .

echo "Done. Artifacts in $OUTPUT_PATH:"
echo
echo "Open Android Studio and click on 'Build' or 'Run' to see result on your mobile"
ls -1 "$OUTPUT_PATH"