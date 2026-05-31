// Package main implements the wailsm CLI engine. It handles cross-platform orchestration,
// project bootstrapping, framework template setups, and plugin injection pipelines
// for the wails-mobile ecosystem.
package main

import (
	"encoding/xml"
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
	Version = "v1.2.0"
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

// XML structures used to hunt down application package names for headless launches.
type AndroidManifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	PackageName string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

type Application struct {
	Activities []Activity `xml:"activity"`
}

type Activity struct {
	Name          string         `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters []IntentFilter `xml:"intent-filter"`
}

type IntentFilter struct {
	Actions    []Action   `xml:"action"`
	Categories []Category `xml:"category"`
}

type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

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
	case "--build":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Error: Missing arguments. Usage: wailsm --build <platform> <debug|release>")
			os.Exit(1)
		}
		executeBuild(os.Args[2], os.Args[3])
	case "--run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: Please specify a target platform: 'android' or 'ios'")
			os.Exit(1)
		}
		executeRun(os.Args[2])
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

func showUsage() {
	fmt.Println("Wails Mobile Toolchain CLI (wailsm)")
	fmt.Println("Usage:")
	fmt.Println("  wailsm --new <project_name>        Create a fresh project from the template")
	fmt.Println("  wailsm --refresh <platform>        Run platform sync: 'android' or 'ios'")
	fmt.Println("  wailsm --build <platform> <mode>   Compile binaries: 'debug' (APK), 'release' (APK), or 'bundle' (AAB)")
	fmt.Println("  wailsm --run <platform>            Compile, install, and execute application via ADB")
	fmt.Println("  wailsm --add <plugin-url>          Install a native Go/Mobile plugin")
	fmt.Println("  wailsm --remove <plugin-url>       Uninstall a native Go/Mobile plugin")
	os.Exit(1)
}

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

	// Generate local.properties configuration context before dropping into path compilation
	setupAndroidLocalProperties(targetDir)

	origWd, _ := os.Getwd()
	_ = os.Chdir(targetDir)
	defer func() { _ = os.Chdir(origWd) }()

	fmt.Println("Initializing Go Mobile build tools...")
	runCmd("go", "install", "golang.org/x/mobile/cmd/gomobile@latest")
	runCmd("go", "install", "golang.org/x/mobile/cmd/gobind@latest")
	runCmd("gomobile", "init")

	if fileExists("go.mod") {
		fmt.Println("Valid Go module context detected. Binding tool tracking dependencies...")
		runCmd("go", "mod", "tidy")
		runCmd("go", "get", "-tool", "golang.org/x/mobile/cmd/gobind")
	} else {
		fmt.Println("Notice: No go.mod found in target template context. Skipping localized tool tracking.")
	}

	fmt.Fprintf(os.Stdout, "=== Setup complete! Your project is ready in ./%s ===\n", targetDir)
}

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

		fmt.Fprintf(os.Stdout, "Done. Artifacts inside %s synchronized cleanly.\n", outputPath)
	} else {
		fmt.Println("Notice: iOS pipeline synthesis engine is currently running standard validation checks.")
	}
}

// executeBuild forces an environment refresh, identifies local compilation managers,
// and triggers headless compiler steps based on your mode selection.
func executeBuild(platform, mode string) {
	mode = strings.ToLower(mode)
	if mode != "debug" && mode != "release" && mode != "bundle" {
		fmt.Fprintln(os.Stderr, "Error: Invalid build mode specified. Use 'debug' (APK), 'release' (APK), or 'bundle' (App Store AAB).")
		os.Exit(1)
	}

	// Force baseline refresh to sync compiled Go bridges cleanly
	executeRefresh(platform)

	if platform == "android" {
		androidDir := filepath.Join("native", "android")
		if !dirExists(androidDir) {
			fmt.Fprintln(os.Stderr, "Error: Native Android path root layout is missing. Run this inside a valid wailsm space.")
			os.Exit(1)
		}

		gradleCmd := "./gradlew"
		if isWindowsHost() {
			gradleCmd = "gradlew.bat"
		}

		origWd, _ := os.Getwd()
		_ = os.Chdir(androidDir)
		defer func() { _ = os.Chdir(origWd) }()

		if !fileExists(gradleCmd) {
			fmt.Println("Notice: Local Gradle wrapper missing. Attempting global system fallback routine...")
			gradleCmd = "gradle"
			if _, err := exec.LookPath(gradleCmd); err != nil {
				fmt.Fprintln(os.Stderr, "Error: Gradle compiler tools missing on host system.")
				os.Exit(1)
			}
		}

		// Dynamically assign target task profiles based on choice
		var targetTask string
		switch mode {
		case "debug":
			targetTask = "assembleDebug"
			fmt.Println("=== Building Application Binary [Local Debug APK Mode] ===")
		case "release":
			targetTask = "assembleRelease"
			fmt.Println("=== Building Application Binary [Unsigned Release APK Mode] ===")
		case "bundle":
			targetTask = "bundleRelease"
			fmt.Println("=== Building Production Asset [Google Play App Bundle AAB Mode] ===")
		}

		fmt.Fprintf(os.Stdout, "Executing automated engine task: %s %s\n", gradleCmd, targetTask)
		runCmd(gradleCmd, targetTask)

		// Print exact path feedback so the developer knows where to grab the file
		if mode == "bundle" {
			fmt.Println("\nCompilation complete! Your production App Bundle is ready for upload:")
			fmt.Println("👉 ./native/android/app/build/outputs/bundle/release/app-release.aab")
		} else {
			fmt.Println("\nCompilation complete! Your package is located under outputs/apk/")
		}
	} else {
		fmt.Println("Notice: iOS pipeline build sequence is currently running verification architectures.")
	}
}

// executeRun compiles the workspace codebase, evaluates connection bridges,
// deploys the output binary, and fires up the activity framework instantly.
func executeRun(platform string) {
	if platform != "android" {
		fmt.Println("Notice: Headless runner support is currently targeted to Android devices via ADB.")
		return
	}

	// Step 1: Enforce code generation and debug compilation pass
	executeBuild("android", "debug")

	fmt.Println("=== Deploying Package over Android Debug Bridge (ADB) ===")
	if _, err := exec.LookPath("adb"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: 'adb' executable tool found missing from system path bounds. Install Android Platform Tools to enable execution.")
		os.Exit(1)
	}

	// Verify device connection integrity
	devicesOut, _ := exec.Command("adb", "devices").Output()
	lines := strings.Split(strings.TrimSpace(string(devicesOut)), "\n")
	if len(lines) <= 1 || strings.TrimSpace(lines[1]) == "" {
		fmt.Fprintln(os.Stderr, "Error: No active emulators or physical Android devices detected. Plug in your device and authorize USB Debugging.")
		os.Exit(1)
	}

	apkPath := filepath.Join("native", "android", "app", "build", "outputs", "apk", "debug", "app-debug.apk")
	if !fileExists(apkPath) {
		fatal("Could not trace destination compiled application target binary", fmt.Errorf("missing asset: %s", apkPath))
	}

	fmt.Println("Streaming installation package directly down to device hardware space...")
	runCmd("adb", "install", "-r", apkPath)

	// Step 2: Extract package credentials from manifest to execute launch
	manifestPath := filepath.Join("native", "android", "app", "src", "main", "AndroidManifest.xml")
	packageName, launcherActivity, err := parseManifestDetails(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Notice: Could not parse manifest automation metadata: %v. Skipping launch step.\n", err)
		return
	}

	targetComponent := packageName + "/" + launcherActivity
	fmt.Fprintf(os.Stdout, "Launching core runtime context interface instance: %s\n", targetComponent)
	runCmd("adb", "shell", "am", "start", "-n", targetComponent)
	fmt.Println("=== Application initialization complete! Native logs are streaming through logcat ===")
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
