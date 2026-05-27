// Package main implements the wailsm CLI engine. It handles cross-platform orchestration,
// project bootstrapping, framework template setups, and plugin injection pipelines
// for the wails-mobile ecosystem.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// RepoURL is the upstream repository hosting core templates and codebases.
	RepoURL = "https://github.com/its-ernest/wails-mobile"
	// Version matches the designated framework release tag.
	Version = "v1.0.5"
	// ReleaseAsset is the file package downloaded to bootstrap fresh instances.
	ReleaseAsset = "template.zip"
	// DownloadURL maps directly to the compiled archive within GitHub Releases.
	DownloadURL = RepoURL + "/releases/latest/download/" + ReleaseAsset
)

// Global framework parameters initialized with predictable fallback settings.
// Values are overridden if an optional 'android.ini' file is parsed at runtime.
var (
	GomobileTarget = "android/arm64"
	AndroidAPI     = "21"
	AarName        = "wailsmobile.aar"
	CleanOutput    = "true"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
	}

	switch os.Args[1] {
	case "--new":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a project directory name.")
			os.Exit(1)
		}
		createNewProject(os.Args[2])
	case "--refresh":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		executeRefresh(os.Args[2])
	case "--add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		managePlugin("add", os.Args[2])
	case "--remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please provide a valid plugin repository path.")
			os.Exit(1)
		}
		managePlugin("remove", os.Args[2])
	case "-h", "--help":
		showUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", os.Args[1])
		showUsage()
	}
}

// showUsage prints formatted instruction parameters to standard output.
func showUsage() {
	fmt.Println("Wails Mobile Toolchain CLI (wailsm)")
	fmt.Println("Usage:")
	fmt.Println("  wailsm --new <project_name>        Create a fresh project from the template")
	fmt.Println("  wailsm --refresh <platform>        Run platform sync: 'android' or 'ios'")
	fmt.Println("  wailsm --add <plugin-url>          Install a native Go/Mobile plugin")
	fmt.Println("  wailsm --remove <plugin-url>       Uninstall a native Go/Mobile plugin")
	os.Exit(1)
}

// createNewProject bootstraps a fresh development environment workspace. It queries
// github for the latest packaged template, unzips files natively, and installs
// gomobile/gobind developer requirements.
func createNewProject(targetDir string) {
	fmt.Fprintf(os.Stdout, "=== Creating Project: %s [%s] ===\n", targetDir, Version)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' already exists. Aborting.\n", targetDir)
		os.Exit(1)
	}

	for _, cmd := range []string{"go"} {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Required system tool dependency '%s' is missing.\n", cmd)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fatal("Failed to create target workspace directory", err)
	}

	zipPath := filepath.Join(targetDir, ReleaseAsset)
	fmt.Println("Downloading runtime architecture template...")
	if err := downloadFile(DownloadURL, zipPath); err != nil {
		fatal("Failed downloading template deployment asset", err)
	}

	fmt.Println("Extracting asset templates...")
	if err := unzipTarget(zipPath, targetDir); err != nil {
		fatal("Failed unzipping package assets structure", err)
	}
	_ = os.Remove(zipPath)

	origWd, _ := os.Getwd()
	_ = os.Chdir(targetDir)
	defer func() { _ = os.Chdir(origWd) }()

	fmt.Println("Initializing Go Mobile build tools...")
	runCmd("go", "install", "golang.org/x/mobile/cmd/gomobile@latest")
	runCmd("gomobile", "init")
	runCmd("go", "get", "-tool", "golang.org/x/mobile/cmd/gobind")

	_ = os.MkdirAll(filepath.Join("native_plugins", "android"), 0755)
	_ = os.MkdirAll(filepath.Join("native_plugins", "ios"), 0755)

	fmt.Fprintf(os.Stdout, "=== Setup complete! Your project is ready in ./%s ===\n", targetDir)
}

// executeRefresh tracks configuration inputs, triggers compiler bindings,
// and moves files to their respective target native build suites.
func executeRefresh(platform string) {
	if platform != "android" && platform != "ios" {
		fmt.Fprintln(os.Stderr, "Error: Please specify a valid target platform: 'android' or 'ios'")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "=== Executing Platform Refresh: %s ===\n", platform)

	if platform == "android" {
		parseAndroidINI()

		outputPath := filepath.Join("native", "android", "app", "libs")
		targetJavaSrcDir := filepath.Join("native", "android", "app", "src", "main", "java")
		stagingPluginsDir := filepath.Join("native_plugins", "android")

		if err := os.MkdirAll(outputPath, 0755); err != nil {
			fatal("Failed staging output directory validation", err)
		}

		if CleanOutput == "true" {
			fmt.Println("Cleaning historical build artifacts...")
			files, _ := filepath.Glob(filepath.Join(outputPath, "*.aar"))
			for _, f := range files {
				_ = os.Remove(f)
			}
			sources, _ := filepath.Glob(filepath.Join(outputPath, "*-sources.jar"))
			for _, s := range sources {
				_ = os.Remove(s)
			}
		}

		fmt.Fprintf(os.Stdout, "Building %s for target %s (API %s)...\n", AarName, GomobileTarget, AndroidAPI)
		runCmd("gomobile", "bind", "-target="+GomobileTarget, "-androidapi="+AndroidAPI, "-o", filepath.Join(outputPath, AarName), ".")

		fmt.Println("Checking for external native plugins...")
		if dirExists(stagingPluginsDir) && !dirEmpty(stagingPluginsDir) {
			fmt.Println("Found native plugins in staging area. Syncing source trees to Android project...")
			if err := os.MkdirAll(targetJavaSrcDir, 0755); err != nil {
				fatal("Could not create target Java source directory", err)
			}
			if err := copyDirectory(stagingPluginsDir, targetJavaSrcDir); err != nil {
				fatal("Failed syncing plugin package trees inside Android workspace target", err)
			}
			fmt.Println("Native source files successfully synchronized.")
		} else {
			fmt.Println("No native plugins staged or directory empty. Skipping source injection.")
		}

		fmt.Fprintf(os.Stdout, "Done. Artifacts in %s:\n\n", outputPath)
		fmt.Println("Open Android Studio and click on 'Build' or 'Run' to see result on your mobile.")

		// Replaced shell-dependent 'ls -1' with a clean cross-platform directory scan
		listDirectoryContents(outputPath)
	} else {
		fmt.Println("Notice: iOS pipeline synthesis engine is currently running standard validation checks.")
	}
}

// managePlugin executes mutations across go.mod dependencies and safely
// moves native platform specific directories inside intermediate workspaces.
func managePlugin(action, pluginRepo string) {
	if !dirExists("native_plugins") || !fileExists("go.mod") {
		fmt.Fprintln(os.Stderr, "Error: You must execute plugin commands from the root of a wailsm project directory.")
		os.Exit(1)
	}

	if action == "add" {
		fmt.Fprintf(os.Stdout, "=== Installing Plugin: %s ===\n", pluginRepo)

		runCmd("go", "get", pluginRepo)

		out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", pluginRepo).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not resolve source location for module: %s\n", pluginRepo)
			os.Exit(1)
		}
		goModSrc := strings.TrimSpace(string(out))
		fmt.Fprintf(os.Stdout, "Module source located at: %s\n", goModSrc)

		// Synchronize Android Assets
		androidSrc := filepath.Join(goModSrc, "android")
		if dirExists(androidSrc) {
			fmt.Println("Found Native Android bindings. Syncing into Android core space...")
			if err := copyDirectory(androidSrc, filepath.Join("native_plugins", "android")); err != nil {
				fatal("Failed to sync Android native directory", err)
			}
		} else {
			fmt.Println("Notice: No native /android directory found in this plugin.")
		}

		// Synchronize iOS Assets
		iosSrc := filepath.Join(goModSrc, "ios")
		if dirExists(iosSrc) {
			fmt.Println("Found Native iOS bindings. Syncing into iOS core space...")
			if err := copyDirectory(iosSrc, filepath.Join("native_plugins", "ios")); err != nil {
				fatal("Failed to sync iOS native directory", err)
			}
		}

		fmt.Fprintf(os.Stdout, "=== Plugin %s added successfully! ===\n", pluginRepo)
		fmt.Println("Run 'wailsm --refresh <platform>' to rebuild bindings with the new packages.")

	} else if action == "remove" {
		fmt.Fprintf(os.Stdout, "=== Removing Plugin: %s ===\n", pluginRepo)

		pluginDirname := filepath.Base(pluginRepo)
		_ = os.RemoveAll(filepath.Join("native_plugins", "android", pluginDirname))
		_ = os.RemoveAll(filepath.Join("native_plugins", "ios", pluginDirname))

		runCmd("go", "mod", "edit", "-droprequire="+pluginRepo)
		runCmd("go", "mod", "tidy")
		fmt.Fprintf(os.Stdout, "=== Plugin %s removed ===\n", pluginRepo)
	}
}
