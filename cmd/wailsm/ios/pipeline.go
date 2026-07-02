// Package ios implements the iOS-specific build and deployment pipeline using Docker and xtools.
// It allows for iOS development and cross-compilation on non-macOS host systems.
package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/its-ernest/wails-mobile/cmd/wailsm/utils"
)

// ScaffoldProject uses xtool inside Docker to create the initial iOS project structure.
func ScaffoldProject(name string) {
	fmt.Printf("Scaffolding iOS project '%s' using xtool...\n", name)

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		_ = os.MkdirAll(iosDir, 0755)
	}

	// Step 1: Ensure Docker Compose is available
	if !utils.CommandExists("docker") {
		utils.Fatal("Docker is required for iOS project scaffolding", fmt.Errorf("docker command not found"))
	}

	// Step 2: Restore Docker files to the ios directory
	restoreDockerFiles(iosDir)

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Step 3: Start the container
	fmt.Println("Starting iOS build container...")
	utils.RunCmd("docker", "compose", "up", "-d", "--build")

	// Step 4: Run xtool init/create
	// We'll use 'xtool init' which is standard for creating a new project structure
	fmt.Println("Running xtool scaffolding sequence...")
	utils.RunCmd("docker", "compose", "exec", "xtool", "xtool", "init", "--name", name)

	fmt.Println("iOS project scaffolded successfully.")
}

// RefreshPipeline triggers the Docker-based iOS build toolchain setup and Go binding generation.
func RefreshPipeline() {
	fmt.Println("Starting iOS Xtool cross-compilation container synthesis layer...")

	iosDir := filepath.Join("native", "ios")
	if !utils.DirExists(iosDir) {
		fmt.Fprintln(os.Stderr, "Error: Native iOS path layout missing. Ensure you have the iOS template files in native/ios.")
		os.Exit(1)
	}

	// Step 1: Ensure Docker Compose is available
	if !utils.CommandExists("docker") {
		utils.Fatal("Docker is required for iOS cross-compilation on non-macOS systems", fmt.Errorf("docker command not found"))
	}

	// Step 2: Validate Docker context
	if !utils.FileExists(filepath.Join(iosDir, "Dockerfile")) || !utils.FileExists(filepath.Join(iosDir, "compose.yml")) {
		fmt.Println("Warning: Docker configuration files missing in native/ios. Attempting to restore from template...")
		restoreDockerFiles(iosDir)
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Step 3: Build/Start the toolchain containers
	fmt.Println("Building/Starting iOS toolchain containers (this may take a while on first run)...")
	utils.RunCmd("docker", "compose", "up", "-d", "--build")

	// Step 4: Run gomobile bind inside the container
	fmt.Println("Generating Go bindings (XFramework) for iOS inside the container...")
	// We use 'xtool' service defined in compose.yml
	utils.RunCmd("docker", "compose", "exec", "xtool", "gomobile", "bind", "-target=ios", "-o", "wailsmobile.xcframework", ".")

	fmt.Println("iOS Platform Refresh complete. Bindings generated in native/ios/wailsmobile.xcframework")
}

// BuildPipeline triggers the compilation of the Swift/iOS application using xtools.
func BuildPipeline(mode string) {
	fmt.Printf("Building iOS application in mode: %s via container compiler...\n", mode)

	iosDir := filepath.Join("native", "ios")
	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// We use 'xtool build' to compile the project
	fmt.Println("Executing xtool build sequence...")
	utils.RunCmd("docker", "compose", "exec", "xtool", "xtool", "build", "--configuration", mode)

	fmt.Println("iOS application built successfully.")
}

// RunPipeline handles deployment to a real device using usbmuxd bridge via xtools.
func RunPipeline() {
	fmt.Println("Mounting bridge interfaces over usbmuxd target layout to real hardware...")

	iosDir := filepath.Join("native", "ios")
	origWd, _ := os.Getwd()
	_ = os.Chdir(iosDir)
	defer func() { _ = os.Chdir(origWd) }()

	// Use 'xtool run' to deploy to the device.
	// This requires usbmuxd to be running on the host and bridged to the container.
	fmt.Println("Deploying to connected iOS device...")
	utils.RunCmd("docker", "compose", "exec", "xtool", "xtool", "run")
}

func restoreDockerFiles(iosDir string) {
	// Look for temp/ios in the repo if we are in dev mode
	cwd, _ := os.Getwd()
	tempIosDir := filepath.Join(cwd, "temp", "ios")
	if utils.DirExists(tempIosDir) {
		_ = utils.CopyDirectory(tempIosDir, iosDir)
		fmt.Println("Restored Docker files from temp/ios")
	} else {
		fmt.Fprintln(os.Stderr, "Error: Could not find iOS Docker templates in temp/ios.")
		os.Exit(1)
	}
}
