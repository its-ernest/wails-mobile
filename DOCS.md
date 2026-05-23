# Wails Mobile Documentation

This document provides a high-level overview of the Wails Mobile project structure and how to navigate the technical documentation.

## Project Structure

- **`wails/`**: The core Go runtime. This contains the bridge, application lifecycle management, and event bus.
- **`wails/permission/`**: Standard plugin for handling Android runtime permissions.
- **`native/android/`**: The Android project template. This is where the Java/Kotlin side of the bridge lives.
- **`_examples/`**: Sample applications demonstrating how to use the framework.

## API Documentation

We use `gomarkdoc` to generate documentation directly from Go source code. You can find detailed API references for each package in their respective directories:

- [Core Runtime (`wails`)](./wails/README.md)
- [Permissions Plugin (`wails/permission`)](./wails/permission/README.md)

## Generating Documentation

To regenerate the documentation after making changes to the Go code, run:

```bash
make docs
```

This requires `gomarkdoc` to be installed:
```bash
go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
```

## Developing Plugins

If you want to extend Wails Mobile with new native capabilities, please refer to the [Plugins Guide](./PLUGINS.md).
