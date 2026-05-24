#!/bin/bash

# Exit immediately if any command fails
set -e

# --- Configuration ---
REPO_URL="https://github.com/its-ernest/wails-mobile" 
VERSION="v1.0.0"
RELEASE_ASSET="template.zip" 
#DOWNLOAD_URL="${REPO_URL}/releases/download/${VERSION}/${RELEASE_ASSET}"
DOWNLOAD_URL="${REPO_URL}/releases/latest/download/${RELEASE_ASSET}"

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
    local plugin_repo="$2" # e.g., github.com/its-ernest/wails-mobile/wails/permission
    
    if [ -z "$plugin_repo" ]; then
        echo "Error: Please provide a valid plugin repository path." >&2
        exit 1
    fi

    if [ ! -d "native_plugins" ] || [ ! -f "go.mod" ]; then
        echo "Error: You must execute plugin commands from the root of a wailsm project directory." >&2
        exit 1
    fi

    if [ "$action" == "add" ]; then
        echo "=== Installing Plugin: ${plugin_repo} ==="
        
        # 1. Download and track the module properly using modern Go Module standards
        go get "$plugin_repo"

        # 2. Ask Go exactly where this package's source code is located on disk
        local go_mod_src
        go_mod_src=$(go list -m -f '{{.Dir}}' "$plugin_repo" 2>/dev/null || true)

        if [ -z "$go_mod_src" ]; then
            echo "Error: Could not resolve source location for module: ${plugin_repo}" >&2
            exit 1
        fi

        echo "Module source located at: ${go_mod_src}"

        # 3. Inject Native Android Packages
        if [ -d "${go_mod_src}/android" ]; then
            echo "Found Native Android bindings. Syncing into Android core space..."
            # Using -R to ensure directory permissions are readable/writable even if cache is read-only
            cp -R "${go_mod_src}/android/." ./native_plugins/android/
        else
            echo "Notice: No native /android directory found in this plugin."
        fi

        # 4. Inject Native iOS Packages
        if [ -d "${go_mod_src}/ios" ]; then
            echo "Found Native iOS bindings. Syncing into iOS core space..."
            cp -R "${go_mod_src}/ios/." ./native_plugins/ios/
        fi
        
        echo "=== Plugin ${plugin_repo} added successfully! ==="
        echo "Run '$0 --refresh <platform>' to rebuild bindings with the new packages."

    elif [ "$action" == "remove" ]; then
        echo "=== Removing Plugin: ${plugin_repo} ==="
        
        # Extract the trailing folder name to locate the directory inside staging
        local plugin_dirname
        plugin_dirname=$(basename "$plugin_repo")

        # Clean matching directories out of native_plugins rooms cleanly
        rm -rf "./native_plugins/android/${plugin_dirname}"
        rm -rf "./native_plugins/ios/${plugin_dirname}"

        # Drop requirement cleanly from go.mod
        go mod edit -droprequire="$plugin_repo"
        go mod tidy
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