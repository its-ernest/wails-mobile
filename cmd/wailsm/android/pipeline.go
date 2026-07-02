// Package android implements the Android-specific build and deployment pipeline.
// It manages Gradle compilation, AAR generation via gomobile, and ADB device orchestration.
package android

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/its-ernest/wails-mobile/cmd/wailsm/utils"
)

var (
	// GomobileTarget defines the Android ABI architecture to build for.
	GomobileTarget = "android/arm64"
	// AndroidAPI specifies the minimum Android SDK level for gomobile bind.
	AndroidAPI     = "21"
	// AarName is the filename of the generated Android archive.
	AarName        = "wailsmobile.aar"
	// CleanOutput determines if previous build artifacts should be purged.
	CleanOutput    = "true"
)

// AndroidManifest represents a subset of the Android XML manifest for metadata extraction.
type AndroidManifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	PackageName string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

// Application represents the application tag in the Android manifest.
type Application struct {
	Activities []Activity `xml:"activity"`
}

// Activity represents an activity tag in the Android manifest.
type Activity struct {
	Name          string         `xml:"http://schemas.android.com/apk/res/android name,attr"`
	IntentFilters []IntentFilter `xml:"intent-filter"`
}

// IntentFilter represents an intent-filter tag in the Android manifest.
type IntentFilter struct {
	Actions    []Action   `xml:"action"`
	Categories []Category `xml:"category"`
}

// Action represents an action tag in the Android manifest.
type Action struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// Category represents a category tag in the Android manifest.
type Category struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// SetupAndroidLocalProperties discover the Android SDK path and writes local.properties.
func SetupAndroidLocalProperties(targetDir string) {
	sdkPath := os.Getenv("ANDROID_HOME")
	if sdkPath == "" {
		sdkPath = os.Getenv("ANDROID_SDK_ROOT")
	}
	if sdkPath != "" {
		propsPath := filepath.Join(targetDir, "native", "android", "local.properties")
		_ = os.MkdirAll(filepath.Dir(propsPath), 0755)
		_ = os.WriteFile(propsPath, []byte(fmt.Sprintf("sdk.dir=%s\n", filepath.ToSlash(sdkPath))), 0644)
	}
}

func parseAndroidINI() {
	if !utils.FileExists("android.ini") {
		return
	}
	// Future implementation: Read optional configuration overrides from android.ini
}

// RefreshPipeline runs the gomobile bind command to sync Go code with the Android project.
// It also synchronizes native plugin source files.
func RefreshPipeline() {
	parseAndroidINI()
	outputPath := filepath.Join("native", "android", "app", "libs")
	targetJavaSrcDir := filepath.Join("native", "android", "app", "src", "main", "java")
	stagingPluginsDir := filepath.Join("native_plugins", "android")

	_ = os.MkdirAll(outputPath, 0755)

	if CleanOutput == "true" {
		files, _ := filepath.Glob(filepath.Join(outputPath, "*.aar"))
		for _, f := range files {
			_ = os.Remove(f)
		}
	}

	fmt.Fprintf(os.Stdout, "Building %s for target %s...\n", AarName, GomobileTarget)
	utils.RunCmd("gomobile", "bind", "-target="+GomobileTarget, "-androidapi="+AndroidAPI, "-o", filepath.Join(outputPath, AarName), ".")

	if utils.DirExists(stagingPluginsDir) && !utils.DirEmpty(stagingPluginsDir) {
		_ = os.MkdirAll(targetJavaSrcDir, 0755)
		if err := utils.CopyDirectory(stagingPluginsDir, targetJavaSrcDir); err != nil {
			utils.Fatal("Failed syncing plugin package trees inside Android workspace", err)
		}
	}
}

// BuildPipeline executes the Gradle build process to generate APK or AAB files.
func BuildPipeline(mode string) {
	androidDir := filepath.Join("native", "android")
	if !utils.DirExists(androidDir) {
		fmt.Fprintln(os.Stderr, "Error: Native Android path layout missing.")
		os.Exit(1)
	}

	gradleCmd := "./gradlew"
	if utils.IsWindowsHost() {
		gradleCmd = "gradlew.bat"
	}

	origWd, _ := os.Getwd()
	_ = os.Chdir(androidDir)
	defer func() { _ = os.Chdir(origWd) }()

	var targetTask string
	switch mode {
	case "debug":
		targetTask = "assembleDebug"
	case "release":
		targetTask = "assembleRelease"
	case "bundle":
		targetTask = "bundleRelease"
	}

	utils.RunCmd(gradleCmd, targetTask)
}

// RunPipeline installs the compiled APK to a connected device via ADB and launches it.
func RunPipeline() {
	if _, err := exec.LookPath("adb"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: 'adb' tool missing.")
		os.Exit(1)
	}

	apkPath := filepath.Join("native", "android", "app", "build", "outputs", "apk", "debug", "app-debug.apk")
	utils.RunCmd("adb", "install", "-r", apkPath)

	manifestPath := filepath.Join("native", "android", "app", "src", "main", "AndroidManifest.xml")
	packageName, launcherActivity, err := parseManifestDetails(manifestPath)
	if err != nil {
		return
	}

	utils.RunCmd("adb", "shell", "am", "start", "-n", packageName+"/"+launcherActivity)
}

func parseManifestDetails(manifestPath string) (string, string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", "", err
	}
	var manifest AndroidManifest
	if err := xml.Unmarshal(data, &manifest); err != nil {
		return "", "", err
	}

	packageName := manifest.PackageName
	launcherActivity := ""
	for _, act := range manifest.Application.Activities {
		for _, filter := range act.IntentFilters {
			for _, action := range filter.Actions {
				if action.Name == "android.intent.action.MAIN" {
					launcherActivity = act.Name
					break
				}
			}
		}
	}
	return packageName, launcherActivity, nil
}
