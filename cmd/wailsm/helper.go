package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
