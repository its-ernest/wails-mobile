.PHONY: docs wailsm

# Generate documentation for all Go packages using gomarkdoc
docs:
	@echo "Generating documentation..."
	@command -v gomarkdoc >/dev/null 2>&1 || { echo >&2 "gomarkdoc is required but not installed. Install it with: go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"; exit 1; }
	@gomarkdoc -o "{{.Dir}}/DOCS.md" ./wails/...
	@echo "Documentation generated in respective directories."

wailsm-sh:
	@echo "Building wails-mobile `wailsm.sh` CLI helper..."
	@rm -rf wailsm && shc -f wailsm.sh && mv wailsm.sh.x wailsm
	@echo "wailsm built successfully."

wailsm:
	@echo "Building wails-mobile `wailsm` CLI helper..."
	@# Compile for Linux machines
	GOOS=linux GOARCH=amd64 go build -o dist/wailsm cmd/wailsm/main.go

	# Compile for macOS machines (M1/M2/M3 Apple Silicon architecture ARM chips)
	GOOS=darwin GOARCH=arm64 go build -o dist/wailsm-mac cmd/wailsm/main.go

	# Compile for Windows machines (.exe output format)
	GOOS=windows GOARCH=amd64 go build -o dist/wailsm.exe cmd/wailsm/main.go
	@echo "wailsm built successfully."