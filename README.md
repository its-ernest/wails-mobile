# Wails Mobile (Stable Android Support)

Wails Mobile is a port of Go v3 implementation to support mobile devies.

Gradle support: `Gradle 9.2.1` by default. You can change the version to your Gradle version

<p>
	<img alt="Example UI" src="./hello-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="Notification plugin" src="./notification-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="output console" src="./screenshot-output.jpg" width="200" style="border-radius:6px;"/>
</p>

<p align="center">
	<a href="https://github.com/its-ernest/wails-mobile/blob/main/LICENSE"><img src="https://img.shields.io/github/license/its-ernest/wails-mobile" alt="license"/></a>
	<a href="https://goreportcard.com/report/github.com/its-ernest/wails-mobile"><img src="https://goreportcard.com/badge/github.com/its-ernest/wails-mobile" alt="goreport"/></a>
	<a href="https://pkg.go.dev/github.com/its-ernest/wails-mobile"><img src="https://pkg.go.dev/badge/github.com/its-ernest/wails-mobile.svg" alt="Go Reference"/></a>
	<a href="https://github.com/its-ernest/wails-mobile/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="contrib"/></a>
</p>

## Quick links

<div style="display:flex;flex-wrap:wrap;gap:10px;">
  <a href="docs/INSTALL.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">Install &amp; Build</a>
  <a href="docs/CLI.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">CLI reference</a>
  <a href="docs/CONTRIBUTING.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">Contributing</a>
  <a href="plugins/PLUGINS.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">Plugin docs</a>
  <a href="plugins/PLUGINS_ANDROID.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">Android plugin guide</a>
  <a href="_examples/helloworld/README.md" style="display:inline-block;padding:10px 16px;background:#2d6cdf;color:#fff;border-radius:8px;text-decoration:none;font-weight:600;">Example Wails Mobile app</a>
</div>

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
