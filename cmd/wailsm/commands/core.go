// Package commands implements the core orchestration logic for the wailsm CLI.
// It manages high-level workflows such as project bootstrapping, plugin management,
// and platform-specific pipeline execution.
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/its-ernest/wails-mobile/cmd/wailsm/android"
	"github.com/its-ernest/wails-mobile/cmd/wailsm/ios"
	"github.com/its-ernest/wails-mobile/cmd/wailsm/utils"
)

const (
	// RepoURL is the upstream repository for wails-mobile templates.
	RepoURL      = "https://github.com/its-ernest/wails-mobile"
	// Version is the current version of the wailsm CLI.
	Version      = "v1.4.0"
	// ReleaseAsset is the name of the template archive to download.
	ReleaseAsset = "template.zip"
	// DownloadURL is the full URL to the latest template release.
	DownloadURL  = RepoURL + "/releases/latest/download/" + ReleaseAsset
)

// ShowUsage prints the help information for the wailsm CLI.
func ShowUsage() {
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

// CreateNewProject bootstraps a new wails-mobile project in the specified directory.
// It downloads the template, initializes Go Mobile tools, and scaffolds the native projects.
func CreateNewProject(targetDir string) {
	fmt.Fprintf(os.Stdout, "=== Creating Project: %s [%s] ===\n", targetDir, Version)
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' already exists. Aborting.\n", targetDir)
		os.Exit(1)
	}

	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Required system tool dependency 'go' is missing.\n")
		os.Exit(1)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		utils.Fatal("Failed to create target workspace directory", err)
	}

	zipPath := filepath.Join(targetDir, ReleaseAsset)
	fmt.Println("Downloading runtime architecture template...")
	if err := utils.DownloadFile(DownloadURL, zipPath); err != nil {
		utils.Fatal("Failed downloading template deployment asset", err)
	}

	fmt.Println("Extracting asset templates...")
	if err := utils.UnzipTarget(zipPath, targetDir); err != nil {
		utils.Fatal("Failed unzipping package assets structure", err)
	}
	_ = os.Remove(zipPath)

	android.SetupAndroidLocalProperties(targetDir)

	// Initialize iOS directory structure
	iosDir := filepath.Join(targetDir, "native", "ios")
	if !utils.DirExists(iosDir) {
		_ = os.MkdirAll(iosDir, 0755)
	}

	// Try to copy Docker templates if in dev repo
	cwd, _ := os.Getwd()
	tempIosDir := filepath.Join(cwd, "temp", "ios")
	if utils.DirExists(tempIosDir) {
		fmt.Println("Deploying iOS Docker toolchain templates...")
		_ = utils.CopyDirectory(tempIosDir, iosDir)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(targetDir)
	defer func() { _ = os.Chdir(origWd) }()

	fmt.Println("Initializing Go Mobile build tools...")
	utils.RunCmd("go", "install", "golang.org/x/mobile/cmd/gomobile@latest")
	utils.RunCmd("go", "install", "golang.org/x/mobile/cmd/gobind@latest")
	utils.RunCmd("gomobile", "init")

	if utils.FileExists("go.mod") {
		fmt.Println("Valid Go module context detected. Binding tool tracking dependencies...")
		utils.RunCmd("go", "mod", "tidy")
		utils.RunCmd("go", "get", "-tool", "golang.org/x/mobile/cmd/gobind")
	}

	// Scaffold iOS project using xtool
	ios.ScaffoldProject(targetDir)

	fmt.Fprintf(os.Stdout, "=== Setup complete! Your project is ready in ./%s ===\n", targetDir)
}

// ExecuteRefresh triggers a platform-specific synchronization pass.
// For Android, this typically runs gomobile bind and syncs native plugins.
// For iOS, it initializes the Docker toolchain and generates XFrameworks.
func ExecuteRefresh(platform string) {
	if platform != "android" && platform != "ios" {
		fmt.Fprintln(os.Stderr, "Error: Please specify a valid target platform: 'android' or 'ios'")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "=== Executing Platform Refresh: %s ===\n", platform)
	if platform == "android" {
		android.RefreshPipeline()
	} else {
		ios.RefreshPipeline()
	}
}

// ExecuteBuild compiles the project for the specified platform and mode.
// It automatically triggers a refresh pass before compilation.
func ExecuteBuild(platform, mode string) {
	mode = strings.ToLower(mode)
	if mode != "debug" && mode != "release" && mode != "bundle" {
		fmt.Fprintln(os.Stderr, "Error: Invalid build mode specified. Use 'debug' (APK), 'release' (APK), or 'bundle' (AAB).")
		os.Exit(1)
	}

	ExecuteRefresh(platform)

	if platform == "android" {
		android.BuildPipeline(mode)
	} else {
		ios.BuildPipeline(mode)
	}
}

// ExecuteRun compiles, installs, and launches the application on a connected device.
func ExecuteRun(platform string) {
	if platform == "android" {
		ExecuteBuild("android", "debug")
		android.RunPipeline()
	} else if platform == "ios" {
		ExecuteBuild("ios", "debug")
		ios.RunPipeline()
	} else {
		fmt.Fprintln(os.Stderr, "Error: Invalid platform. Use 'android' or 'ios'.")
		os.Exit(1)
	}
}

// ManagePlugin handles the installation and removal of native Go/Mobile plugins.
// It manages go.mod dependencies and synchronizes native source files into the intermediate staging area.
func ManagePlugin(action, pluginRepo string) {
	if !utils.DirExists("native_plugins") || !utils.FileExists("go.mod") {
		fmt.Fprintln(os.Stderr, "Error: You must execute plugin commands from the root of a wailsm project directory.")
		os.Exit(1)
	}

	if action == "add" {
		fmt.Fprintf(os.Stdout, "=== Installing Plugin: %s ===\n", pluginRepo)
		utils.RunCmd("go", "get", pluginRepo)

		out, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", pluginRepo).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not resolve source location for module: %s\n", pluginRepo)
			os.Exit(1)
		}
		goModSrc := strings.TrimSpace(string(out))

		androidSrc := filepath.Join(goModSrc, "android")
		if utils.DirExists(androidSrc) {
			fmt.Println("Syncing Android native directory...")
			if err := utils.CopyDirectory(androidSrc, filepath.Join("native_plugins", "android")); err != nil {
				utils.Fatal("Failed to sync Android native directory", err)
			}
		}

		iosSrc := filepath.Join(goModSrc, "ios")
		if utils.DirExists(iosSrc) {
			fmt.Println("Syncing iOS native directory...")
			if err := utils.CopyDirectory(iosSrc, filepath.Join("native_plugins", "ios")); err != nil {
				utils.Fatal("Failed to sync iOS native directory", err)
			}
		}
		fmt.Fprintf(os.Stdout, "=== Plugin %s added successfully! ===\n", pluginRepo)
	} else if action == "remove" {
		fmt.Fprintf(os.Stdout, "=== Removing Plugin: %s ===\n", pluginRepo)
		pluginDirname := filepath.Base(pluginRepo)
		_ = os.RemoveAll(filepath.Join("native_plugins", "android", pluginDirname))
		_ = os.RemoveAll(filepath.Join("native_plugins", "ios", pluginDirname))

		utils.RunCmd("go", "mod", "edit", "-droprequire="+pluginRepo)
		utils.RunCmd("go", "mod", "tidy")
	}
}
