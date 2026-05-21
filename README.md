# Wails Mobile

**Automation scripts currently works for Linux systems**

A small Go-based mobile app starter kit with a reusable runtime core and a hello-world Android template.

## Overview

- `wailsmobile/` - core Go runtime package for mobile bridged apps.
- `examples/hello-world/` - small sample app with embedded frontend, Go backend, and Android sync flow.
- `internal/templates/android/` - Android template app where generated AARs get dumped.
- `wailsm.sh` / `wailsm` - project starter that can be installed into your PATH.

## Install the CLI

Place the compiled `wailsm` binary into a folder on your `PATH`.
Use the command below:

```bash
sudo curl -L -o /usr/local/bin/wailsm https://github.com/its-ernest/wails-mobile/releases/latest/download/wailsm && \
sudo chmod +x /usr/local/bin/wailsm
```

Now you can create a new app from anywhere:

```bash
wailsm --new my-app
```

That command downloads a starter pre-configured Wails Mobile project files into `my-app/` directory, ready to be edited and compiled.

The pre-configured project can be extended(adding more code), 
when it's time to build, open `native/android/` folder in Android Studio and click `Build`.

## Editing and Extending the project

Suppose you add more code for your mobile app, you can ensure the build is refreshed to reflect all your changes with the command below:

```bash
wailsm --refresh
```

After that you can click on `Build app` from Android studio to build your app

The above command sync script reads `config.ini`, builds the AAR with `gomobile bind`, and copies the generated library into `native/android/app/libs`.

## Configuring sync

Edit `examples/hello-world/config.ini` to change:

- `gomobile_target` - Android ABI(s) to build
- `androidapi` - Android API level for gomobile
- `aar_name` - output AAR file name
- `template_min_sdk` / `template_target_sdk` - values patched into the Android template

## Requirements

- Go 1.24+
- `gomobile` installed
- Android SDK + NDK installed for `gomobile bind`
- `git`, `curl`, `unzip` for the installer flow

## Notes

- The repo is in initial stages
- Improvement (More packages expected)
- I welcome every contribution to make this project stable
- Code cleanup required
