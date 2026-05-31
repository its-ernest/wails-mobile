package main

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// isWindowsHost checks the compile-time operating system target.
// It returns true if the binary is running on a native Windows machine environment.
func isWindowsHost() bool {
	return runtime.GOOS == "windows"
}

// parseAndroidINI reads system instructions from local .ini configuration matrices.
// It automatically normalizes Windows-style trailing carriage lines (\r\n).
func parseAndroidINI() {
	iniFile := "android.ini"
	if !fileExists(iniFile) {
		fmt.Println("Notice: android.ini not found. Using framework defaults.")
		return
	}

	fmt.Println("Reading configurations from android.ini...")
	file, err := os.Open(iniFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || (strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")) {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "gomobile_target":
			GomobileTarget = val
		case "androidapi":
			AndroidAPI = val
		case "aar_name":
			AarName = val
		case "clean_output":
			CleanOutput = val
		}
	}
}

// parseManifestDetails extracts target identity metrics out of your XML asset maps safely.
func parseManifestDetails(path string) (string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	var manifest AndroidManifest
	if err := xml.Unmarshal(bytes, &manifest); err != nil {
		return "", "", err
	}

	pkg := manifest.PackageName
	launcher := ""

	// Identify the explicit activity registration holding the LAUNCHER intent filter category flag
	for _, act := range manifest.Application.Activities {
		isLauncher := false
		for _, filter := range act.IntentFilters {
			hasMain := false
			hasLaunch := false
			for _, action := range filter.Actions {
				if action.Name == "android.intent.action.MAIN" {
					hasMain = true
				}
			}
			for _, cat := range filter.Categories {
				if cat.Name == "android.intent.category.LAUNCHER" {
					hasLaunch = true
				}
			}
			if hasMain && hasLaunch {
				isLauncher = true
				break
			}
		}
		if isLauncher {
			launcher = act.Name
			break
		}
	}

	if pkg == "" || launcher == "" {
		return pkg, launcher, fmt.Errorf("failed tracking valid launcher details within schema layout metrics")
	}

	return pkg, launcher, nil
}

// runCmd spawns an OS sub-process linking stdout/stderr pipelines directly to the terminal shell.
func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal(fmt.Sprintf("Command processing failure on executing: %s %s", name, strings.Join(args, " ")), err)
	}
}

// downloadFile gets an network asset via HTTP GET and streams bytes onto local storage.
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status service response download payload: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzipTarget decompresses architectural zip packages while protecting against Zip Slip vulnerabilities.
func unzipTarget(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		// Zip Slip Vulnerability Guard
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path output breaking bounds constraints: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// copyDirectory recursively parses directory trees and duplicates contents across targets.
func copyDirectory(scrDir, destDir string) error {
	return filepath.Walk(scrDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(scrDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// setupAndroidLocalProperties handles automatic discovery of the Android SDK path
// and writes the local.properties initialization configuration file.
func setupAndroidLocalProperties(projectRoot string) {
	fmt.Println("Locating Android SDK workspace path elements...")
	var sdkPath string

	// Strategy 1: Attempt to look up adb on the active platform path variables
	if adbPath, err := exec.LookPath("adb"); err == nil {
		// Clean the symlink boundaries and jump up two parent folders: .../Android/Sdk/platform-tools/adb -> .../Android/Sdk
		if evalPath, err := filepath.EvalSymlinks(adbPath); err == nil {
			sdkPath = filepath.Dir(filepath.Dir(evalPath))
		}
	}

	// Strategy 2: Fall back to checking common environment flag definitions
	if sdkPath == "" {
		if envHome := os.Getenv("ANDROID_HOME"); envHome != "" {
			sdkPath = envHome
		} else if envRoot := os.Getenv("ANDROID_SDK_ROOT"); envRoot != "" {
			sdkPath = envRoot
		}
	}

	// Strategy 3: Hardcoded system fallbacks based on typical platform defaults
	if sdkPath == "" {
		homeDir, _ := os.UserHomeDir()
		if runtime.GOOS == "windows" {
			sdkPath = filepath.Join(homeDir, "AppData", "Local", "Android", "Sdk")
		} else if runtime.GOOS == "darwin" {
			sdkPath = filepath.Join(homeDir, "Library", "Android", "sdk")
		} else {
			sdkPath = filepath.Join(homeDir, "Android", "Sdk")
		}
	}

	// Verify discovery before writing properties layout profiles
	if info, err := os.Stat(sdkPath); err == nil && info.IsDir() {
		// Escape backwards slashes on Windows hosts
		escapedPath := sdkPath
		if runtime.GOOS == "windows" {
			escapedPath = strings.ReplaceAll(sdkPath, "\\", "\\\\")
		}

		propertiesContent := fmt.Sprintf("# Generated automatically by wailsm CLI installer context\n# Location tracking configuration parameters\nsdk.dir=%s\n", escapedPath)
		targetFile := filepath.Join(projectRoot, "native", "android", "local.properties")

		// Create target wrapper native directory trees if missing
		_ = os.MkdirAll(filepath.Dir(targetFile), 0755)

		if err := os.WriteFile(targetFile, []byte(propertiesContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed writing custom local.properties build targets: %v\n", err)
		} else {
			fmt.Printf("Successfully targeted SDK tracking coordinates: %s\n", sdkPath)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Warning: Android SDK root location could not be verified automatically.")
		fmt.Fprintln(os.Stderr, "Notice: You may need to generate standard 'local.properties' inside native/android/ manually before running compilations.")
	}
}

// fileExists validates if a specific file artifact is accessible on disk.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks for the existence of a valid directory segment.
func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// dirEmpty determines if a folder contains downstream data assets.
func dirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF
}

// listDirectoryContents acts as a platform-agnostic alternative to 'ls -1'.
func listDirectoryContents(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Notice: Could not scan contents of directory %s\n", dirPath)
		return
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

// fatal dumps an error block to stderr and forces runtime termination.
func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "Fatal Error: %s: %v\n", message, err)
	os.Exit(1)
}
