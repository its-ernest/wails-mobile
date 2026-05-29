# Install & Build

This document contains the installation, bootstrapping, build and uninstall steps used during development.

## Prerequisites

- Go 1.24+
- Android SDK and Android NDK installed (Android Studio recommended)
- `git`, `curl`

## CLI installation

You can use `Go` to install `wails-mobile` or Choose to use Release binaries.

* Using `Go`:

```bash
go install github.com/its-ernest/wails-mobile/cmd/wailsm@latest
```

## Bootstrapping a new project

```bash
wailsm --new my-new-app
```

This downloads a pre-configured starter template into `./my-new-app/` ready to edit.

Wails app structure (what you'll find in a generated app)

- `wails/` — core Go runtime package for mobile bridged apps
- `_examples/helloworld/` — example with embedded frontend, Go backend, and Android sync flow
- `native/android/` — Android Studio project generated from the AAR/JAR artifacts
- `frontend/` — embedded web UI assets

## Refresh / Sync (generate native bindings)

The CLI performs an automatic sync before `build` and `run`. Manually trigger sync when you want to force a regeneration:

```bash
wailsm --refresh android
# or
wailsm --refresh ios
```

This drops the resulting `.aar` and `.jar` into `native/android/app/libs`.

## Running App Immediately to test result

Run on a connected Android device (installs and opens app via ADB):

```bash
wailsm --run android
```

## Build release APK / AAB:

```bash
wailsm --build android debug
wailsm --build android release
wailsm --build android bundle
```

Output AAB (example): `./native/android/app/build/outputs/bundle/release/app-release.aab`


##  Notes & troubleshooting

- Ensure `ANDROID_HOME` or equivalent environment variables are set and Android SDK/NDK tools are available.
- For Gradle issues, open `native/android/` in Android Studio and inspect `gradle-wrapper.properties` or SDK/NDK configuration.
- If bindings seem stale, run `wailsm --refresh android` to force a rebuild of the generated AAR/JAR artifacts.

For additional developer notes see `docs/CONTRIBUTING.md` and the project `plugins/` directory.
