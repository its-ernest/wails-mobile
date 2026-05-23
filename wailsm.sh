#!/bin/bash

# Exit immediately if any command fails
set -e

# --- Configuration ---
REPO_URL="https://github.com/its-ernest/wails-mobile" 
VERSION="v1.0.0"
RELEASE_ASSET="template.zip" 
DOWNLOAD_URL="${REPO_URL}/releases/download/${VERSION}/${RELEASE_ASSET}"

show_usage() {
    echo "Wails Mobile Toolchain CLI (wailsm)" >&2
    echo "Usage:" >&2
    echo "  $0 --new <project_name>        Create a fresh project from the template" >&2
    echo "  $0 --refresh <platform>        Run platform sync: 'android' or 'ios'" >&2
    echo "  $0 --add <plugin-url>          Install a native Go/Mobile plugin" >&2
    echo "  $0 --remove <plugin-url>       Uninstall a native Go/Mobile plugin" >&2
    exit 1
}

if [ -z "$1" ]; then
    show_usage
fi

# --- Core Actions ---

create_new_project() {
    local target_dir="$1"
    if [ -z "$target_dir" ]; then
        echo "Error: Please specify a project directory name." >&2
        exit 1
    fi

    echo "=== Creating Project: ${target_dir} [${VERSION}] ==="
    if [ -d "$target_dir" ]; then
        echo "Error: Directory '${target_dir}' already exists. Aborting." >&2
        exit 1
    fi

    # Pre-flight check
    for cmd in curl unzip go; do
        if ! command -v "$cmd" &> /dev/null; then
            echo "Error: Required command '$cmd' is not installed." >&2
            exit 1
        fi
    done

    mkdir -p "$target_dir"
    echo "Downloading runtime architecture template..."
    curl -L -sS "$DOWNLOAD_URL" -o "$target_dir/$RELEASE_ASSET"
    
    echo "Extracting asset templates..."
    unzip -q -o "$target_dir/$RELEASE_ASSET" -d "$target_dir"
    rm "$target_dir/$RELEASE_ASSET"

    echo "Initializing Go Mobile build tools..."
    cd "$target_dir"
    go install golang.org/x/mobile/cmd/gomobile@latest
    gomobile init
    go get -tool golang.org/x/mobile/cmd/gobind
    
    # Initialize a dedicated directory to house cloned native platform plugins
    mkdir -p native_plugins/android native_plugins/ios

    echo "=== Setup complete! Your project is ready in ./${target_dir} ==="
}

execute_refresh() {
    local platform="$1"
    if [[ "$platform" != "android" && "$platform" != "ios" ]]; then
        echo "Error: Please specify a valid target platform: 'android' or 'ios'" >&2
        echo "Example: $0 --refresh android" >&2
        exit 1
    fi

    local target_script=""
    if [ "$platform" == "android" ]; then
        target_script="./android.sh"
    else
        target_script="./ios.sh"
    fi

    if [ -f "$target_script" ]; then
        echo "=== Executing Platform Refresh: ${platform} ==="
        if [ -x "$target_script" ]; then
            $target_script
        else
            bash "$target_script"
        fi
    else
        echo "Error: Missing pipeline worker script '${target_script}' in current path." >&2
        exit 1
    fi
}

manage_plugin() {
    local action="$1" # "add" or "remove"
    local plugin_repo="$2" # e.g., github.com/username/wailspackage-camera
    
    if [ -z "$plugin_repo" ]; then
        echo "Error: Please provide a valid plugin repository path." >&2
        echo "Example: $0 --${action} github.com/wailspackage/camera" >&2
        exit 1
    fi

    # Ensure this command is executed within a valid wailsm project space
    if [ ! -d "native_plugins" ] || [ ! -f "go.mod" ]; then
        echo "Error: You must execute plugin commands from the root of a wailsm project directory." >&2
        exit 1
    fi

    # Calculate local download paths inside GOPATH source structures
    local go_path_src
    go_path_src=$(go env GOPATH)/src/${plugin_repo}

    if [ "$action" == "add" ]; then
        echo "=== Installing Plugin: ${plugin_repo} ==="
        
        # 1. Register logic via Go Modules
        go get "$plugin_repo"
        
        # Force download source code to standard GOPATH for native source extraction
        # (Using classic v1 fallback flag or shallow clone into staging environment if using pure go modules mode)
        GO111MODULE=off go get -d "$plugin_repo" || true

        # 2. Inject Native Android Packages
        if [ -d "${go_path_src}/android" ]; then
            echo "Found Native Android bindings. Syncing into Android core space..."
            cp -r "${go_path_src}/android/." ./native_plugins/android/
            # Note for your android.sh asset worker: 
            # Make sure android.sh copies everything from `./native_plugins/android/*` 
            # into the actual destination Android Studio app directory during a --refresh run.
        fi

        # 3. Inject Native iOS Packages
        if [ -d "${go_path_src}/ios" ]; then
            echo "Found Native iOS bindings. Syncing into iOS core space..."
            cp -r "${go_path_src}/ios/." ./native_plugins/ios/
        fi
        
        echo "=== Plugin ${plugin_repo} added successfully! ==="
        echo "Run '$0 --refresh <platform>' to rebuild bindings with the new packages."

    elif [ "$action" == "remove" ]; then
        echo "=== Uninstalling Plugin: ${plugin_repo} ==="
        
        # Read the directories inside the plugin repo before removing to scrub target files cleanly
        if [ -d "${go_path_src}/android" ]; then
            # Clean matching packages out of native_plugins room
            local plugin_dirname
            plugin_dirname=$(basename "$plugin_repo")
            rm -rf "./native_plugins/android/${plugin_dirname:?}"
        fi
        if [ -d "${go_path_src}/ios" ]; then
            rm -rf "./native_plugins/ios/$(basename "$plugin_repo")"
        fi

        go drop "$plugin_repo" || go mod edit -droprequire="$plugin_repo"
        echo "=== Plugin ${plugin_repo} removed ==="
    fi
}

# --- Command Routing Logic Loop ---
case "$1" in
    --new)
        create_new_project "$2"
        ;;
    --refresh)
        execute_refresh "$2"
        ;;
    --add)
        manage_plugin "add" "$2"
        ;;
    --remove)
        manage_plugin "remove" "$2"
        ;;
    -h|--help)
        show_usage
        ;;
    *)
        echo "Error: Unknown option '$1'" >&2
        show_usage
        ;;
esac