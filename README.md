# Wails Mobile

### NB Latest: This project has been ported to [sweet-juice/core](https://github.com/sweet-juice/core) Looking for the stable and actively-maintained version of this project?

Wails Mobile ~~is~~ was an early port of Go v3 implementation to support mobile devies.

Gradle support: `Gradle 9.2.1` by default. You can change the version to your Gradle version

⚠️ **Important Note:** iOS support is currently experimental and not ripened as Android support!

<p>
	<img alt="Example UI" src="./hello-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="Notification plugin" src="./notification-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="output console" src="./screenshot-output.jpg" width="200" style="border-radius:6px;"/>
</p>

<p align="center">
	<a href="https://github.com/its-ernest/wails-mobile/blob/main/LICENSE">
	<img src="https://img.shields.io/github/license/its-ernest/wails-mobile" alt="license"/>
	</a>
	<a href="https://goreportcard.com/report/github.com/its-ernest/wails-mobile">
	<img src="https://goreportcard.com/badge/github.com/its-ernest/wails-mobile" alt="goreport"/>
	</a>
	<a href="https://pkg.go.dev/github.com/its-ernest/wails-mobile">
	<img src="https://pkg.go.dev/badge/github.com/its-ernest/wails-mobile.svg" alt="Go Reference"/>
	</a>
	<a href="https://github.com/its-ernest/wails-mobile/issues">
	<img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="contrib"/>
	</a>
</p>

## Quick links

1. [Install & Build](docs/INSTALL.md)
2. [CLI reference](docs/CLI.md)
3. [Contributing](docs/CONTRIBUTING.md)
4. [Plugin docs](docs/PLUGINS_DOCS.md)
5. [Android plugin guide](docs/PLUGINS_ANDROID.md)
6. [Example Wails Mobile app](_examples/helloworld/README.md)
7. [Benchmarks](docs/BENCHMARKS.md)

## Requirements

- Go 1.24+
- Android SDK & NDK (install via Android Studio or your preferred method)
- `git`, `curl`, `unzip`

## Pre-packed plugins

- `plugins/logger`
- `plugins/notification`
- `plugins/permission`
- `plugins/devicestate`
- `plugins/workmanager`
- `plugins/osapi`
- `plugins/biometrics`
- `plugins/filepicker`

See - [Full List](#)

## Notes

- For contribution guidelines and plugin conventions see `docs/CONTRIBUTING.md`.
