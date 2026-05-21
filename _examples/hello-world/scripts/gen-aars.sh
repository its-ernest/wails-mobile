#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"

template_libs="$(cd "$(dirname "$0")/../../../internal/templates/android/app/libs" && pwd)"
mkdir -p "$template_libs"

rm -f "$template_libs"/*.aar "$template_libs"/*.jar

echo "Building single hello-world AAR into template app libs..."
gomobile bind -target=android/arm64 -androidapi 21 -o "$template_libs/helloworld.aar" github.com/its-ernest/wails-mobile/examples/hello-world

echo "Finished generating single AAR in templates: $template_libs/helloworld.aar"
