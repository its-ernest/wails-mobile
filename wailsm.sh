#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# --- Configuration ---
REPO_URL="https://github.com/its-ernest/wails-mobile" 
VERSION="v1.0.0"
RELEASE_ASSET="template.zip" 
DOWNLOAD_URL="${REPO_URL}/releases/download/${VERSION}/${RELEASE_ASSET}"

# --- Usage Guide ---
show_usage() {
    echo "Usage:" >&2
    echo "  $0 --new <project_name>    Create a fresh project directory from the template" >&2
    echo "  $0 --refresh               Execute 'refresh.sh' in the current directory" >&2
    exit 1
}

# Ensure at least one argument is provided
if [ -z "$1" ]; then
    show_usage
fi

# --- Core Actions ---

create_new_project() {
    local target_dir="$1"

    if [ -z "$target_dir" ]; then
        echo "Error: Please specify a project directory name." >&2
        echo "Example: $0 --new my-awesome-app" >&2
        exit 1
    fi

    echo "=== Starting Project Creation [Version: ${VERSION}] ==="

    # Pre-flight checks
    for cmd in curl unzip; do
        if ! command -v "$cmd" &> /dev/null; then
            echo "Error: Required tool '$cmd' is not installed." >&2
            exit 1
        fi
    done

    # Ensure target directory doesn't collide
    if [ -d "$target_dir" ]; then
        echo "Error: Directory '${target_dir}' already exists. Aborting." >&2
        exit 1
    fi

    echo "Creating project directory: ${target_dir}..."
    mkdir -p "$target_dir"

    echo "Downloading ${RELEASE_ASSET}..."
    curl -L -sS "$DOWNLOAD_URL" -o "$target_dir/$RELEASE_ASSET"

    echo "Extracting template into ${target_dir}..."
    unzip -q -o "$target_dir/$RELEASE_ASSET" -d "$target_dir"

    rm "$target_dir/$RELEASE_ASSET"
    echo "=== Project Created Successfully inside '${target_dir}' ==="
}

execute_refresh() {
    echo "=== Looking for Refresh Script ==="
    
    if [ -f "./refresh.sh" ]; then
        echo "Found './refresh.sh'. Executing..."
        # Using 'bash' or checking execution permissions
        if [ -x "./refresh.sh" ]; then
            ./refresh.sh
        else
            echo "Notice: 'refresh.sh' is not executable. Running via bash interpreter directly..."
            bash ./refresh.sh
        fi
    else
        echo "Error: No 'refresh.sh' found in the current working directory ($(pwd))." >&2
        exit 1
    fi
}

# --- Command Routing Parsing Loop ---
case "$1" in
    --new)
        create_new_project "$2"
        ;;
    --refresh)
        execute_refresh
        ;;
    -h|--help)
        show_usage
        ;;
    *)
        echo "Error: Unknown option '$1'" >&2
        show_usage
        ;;
esac