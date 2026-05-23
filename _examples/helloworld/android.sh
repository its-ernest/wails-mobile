#!/bin/bash
set -e

# Get locations relative to this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# --- Default Configurations (Fallbacks) ---
GOMOBILE_TARGET="android/arm64"
ANDROID_API="21"
AAR_NAME="wailsmobile.aar"
CLEAN_OUTPUT="true"

INI_FILE="$SCRIPT_DIR/android.ini"

# --- INI Parser ---
if [ -f "$INI_FILE" ]; then
    echo "Reading configurations from android.ini..."
    while IFS='= ' read -r key value; do
        # Strip trailing carriage returns if file was edited on Windows (CRLF)
        value=$(echo "$value" | tr -d '\r')
        
        # Safe structural line cleaning
        [[ -z "$key" ]] && continue
        [[ "$key" == "#"* ]] && continue
        [[ "$key" == ";"* ]] && continue
        [[ "$key" == "["* ]] && continue
        
        case "$key" in
            gomobile_target) GOMOBILE_TARGET="$value" ;;
            androidapi)      ANDROID_API="$value" ;;
            aar_name)        AAR_NAME="$value" ;;
            clean_output)    CLEAN_OUTPUT="$value" ;;
        esac
    done < "$INI_FILE"
else
    echo "Notice: android.ini not found. Using framework defaults."
fi

# Define straight, predictable paths
OUTPUT_PATH="$SCRIPT_DIR/native/android/app/libs"
TARGET_JAVA_SRC_DIR="$SCRIPT_DIR/native/android/app/src/main/java"
STAGING_PLUGINS_DIR="$SCRIPT_DIR/native_plugins/android"

# Clean old artifacts conditionally based on config
mkdir -p "$OUTPUT_PATH"
if [ "$CLEAN_OUTPUT" = "true" ]; then
    echo "Cleaning historical build artifacts..."
    rm -f "$OUTPUT_PATH"/*.aar "$OUTPUT_PATH"/*-sources.jar
fi

# Jump into the local package directory so gomobile reads the local files
cd "$SCRIPT_DIR"

echo "Building ${AAR_NAME} for target ${GOMOBILE_TARGET} (API ${ANDROID_API})..."
gomobile bind -target="${GOMOBILE_TARGET}" -androidapi="${ANDROID_API}" -o "$OUTPUT_PATH/${AAR_NAME}" .

# --- Synchronize Native Java Plugins ---
echo "Checking for external native plugins..."
if [ -d "$STAGING_PLUGINS_DIR" ] && [ "$(ls -A "$STAGING_PLUGINS_DIR")" ]; then
    echo "Found native plugins in staging area. Syncing source trees to Android project..."
    mkdir -p "$TARGET_JAVA_SRC_DIR"
    cp -r "$STAGING_PLUGINS_DIR"/. "$TARGET_JAVA_SRC_DIR/"
    echo "Native source files successfully synchronized."
else
    echo "No native plugins staged or directory empty. Skipping source injection."
fi
# ----------------------------------------

echo "Done. Artifacts in $OUTPUT_PATH:"
echo
echo "Open Android Studio and click on 'Build' or 'Run' to see result on your mobile"
ls -1 "$OUTPUT_PATH"