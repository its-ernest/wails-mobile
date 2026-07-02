# Wailsm CLI Architecture

This document provides a technical overview of the `wailsm` toolchain architecture.

## Overview

`wailsm` is a cross-platform CLI tool designed to manage the lifecycle of Wails Mobile applications. It handles project initialization, dependency management, plugin synchronization, and native platform builds.

## Package Structure

The codebase is modularized by domain to allow for independent scaling of platform support:

- **`main`**: The entry point. Handles command routing and argument validation.
- **`commands`**: High-level orchestration. Implements the business logic for commands like `--new`, `--add`, and `--refresh`.
- **`android`**: Android-specific pipeline. Manages Gradle builds, AAR generation via `gomobile`, and ADB device orchestration.
- **`ios`**: iOS-specific pipeline. Leverages Docker and `xtools` for cross-compilation on non-macOS systems.
- **`utils`**: Shared cross-platform utilities for file I/O, networking, and process execution.

## Key Workflows

### Project Initialization (`--new`)
1. Downloads the core framework template.
2. Initializes the Go workspace.
3. Automatically scaffolds both Android (Gradle) and iOS (`xtools`) native projects.

### Plugin Management (`--add` / `--remove`)
1. Manages Go module dependencies via `go get`.
2. Synchronizes native source code from the plugin's repository into a local staging area (`native_plugins/`).
3. Primes the code for the next platform refresh.

### Platform Refresh (`--refresh`)
1. Compiles Go code into native binaries (`.aar` for Android, `.xcframework` for iOS).
2. Pushes staged native plugin source files directly into the native project trees.

## Native Integration

`wailsm` ensures a "Go-Centric" experience. While it manages complex native build tools under the hood, the developer primarily interacts with Go and JavaScript. 

- **Android**: Direct integration with standard Android Studio projects.
- **iOS**: Uses Docker to provide a reliable build environment on Linux and Windows, bypassing the strict macOS requirement for initial development.

## Developer Note

To regenerate the technical API documentation for this package, run:
```bash
gomarkdoc ./cmd/wailsm/...
```
