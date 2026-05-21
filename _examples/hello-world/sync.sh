#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.ini"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "Missing config.ini in $SCRIPT_DIR"
  exit 1
fi

get_config() {
  local key="$1"
  grep -E "^[[:space:]]*$key[[:space:]]*=" "$CONFIG_FILE" | tail -n1 | cut -d'=' -f2- | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' || true
}

GOMOBILE_TARGET="$(get_config gomobile_target)"
ANDROID_API="$(get_config androidapi)"
PACKAGE_PATH="$(get_config package)"
AAR_NAME="$(get_config aar_name)"
OUTPUT_DIR="$(get_config output_dir)"
CLEAN_OUTPUT="$(get_config clean_output)"
TEMPLATE_MIN_SDK="$(get_config template_min_sdk)"
TEMPLATE_TARGET_SDK="$(get_config template_target_sdk)"
GOMOBILE_FLAGS="$(get_config gomobile_flags)"

GOMOBILE_TARGET="${GOMOBILE_TARGET:-android/arm64}"
ANDROID_API="${ANDROID_API:-23}"
PACKAGE_PATH="${PACKAGE_PATH:-github.com/its-ernest/wails-mobile/examples/hello-world}"
AAR_NAME="${AAR_NAME:-helloworld.aar}"
OUTPUT_DIR="${OUTPUT_DIR:-internal/templates/android/app/libs}"
CLEAN_OUTPUT="${CLEAN_OUTPUT:-true}"

OUTPUT_PATH="$REPO_ROOT/$OUTPUT_DIR"
BUILD_GRADLE="$REPO_ROOT/internal/templates/android/app/build.gradle"

command -v gomobile >/dev/null 2>&1 || {
  echo "gomobile not found. Install it first: go install golang.org/x/mobile/cmd/gomobile@latest"
  exit 1
}

mkdir -p "$OUTPUT_PATH"
cd "$REPO_ROOT"

if [ "$CLEAN_OUTPUT" = "true" ]; then
  echo "Cleaning existing generated template libs in $OUTPUT_PATH"
  rm -f "$OUTPUT_PATH"/*.aar "$OUTPUT_PATH"/*-sources.jar
fi

if [ -n "$TEMPLATE_MIN_SDK" ] && [ -f "$BUILD_GRADLE" ]; then
  echo "Syncing Android template minSdk to $TEMPLATE_MIN_SDK"
  perl -pi -e 's/(minSdk\s+)\d+/$1$ENV{TEMPLATE_MIN_SDK}/' "$BUILD_GRADLE"
fi

if [ -n "$TEMPLATE_TARGET_SDK" ] && [ -f "$BUILD_GRADLE" ]; then
  echo "Syncing Android template targetSdk to $TEMPLATE_TARGET_SDK"
  perl -pi -e 's/(targetSdk\s+)\d+/$1$ENV{TEMPLATE_TARGET_SDK}/' "$BUILD_GRADLE"
fi

read -r -a GOMOBILE_FLAGS_ARRAY <<< "$GOMOBILE_FLAGS"

echo "Building AAR from package: $PACKAGE_PATH"
echo "gomobile target: $GOMOBILE_TARGET"
echo "android api: $ANDROID_API"
echo "output: $OUTPUT_PATH/$AAR_NAME"

gomobile bind -target "$GOMOBILE_TARGET" -androidapi "$ANDROID_API" -o "$OUTPUT_PATH/$AAR_NAME" ${GOMOBILE_FLAGS_ARRAY[@]} "$PACKAGE_PATH"

echo "Sync complete. Generated AAR and sources jar are in: $OUTPUT_PATH"
ls -1 "$OUTPUT_PATH"
