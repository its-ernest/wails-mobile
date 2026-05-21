#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"

./gen-aars.sh

if [ -x ./gradlew ]; then
  ./gradlew assembleDebug
else
  gradle assembleDebug
fi
