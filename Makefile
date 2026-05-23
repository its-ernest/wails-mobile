.PHONY: docs wailsm

# Generate documentation for all Go packages using gomarkdoc
docs:
	@echo "Generating documentation..."
	@command -v gomarkdoc >/dev/null 2>&1 || { echo >&2 "gomarkdoc is required but not installed. Install it with: go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"; exit 1; }
	@gomarkdoc -o "{{.Dir}}/DOCS.md" ./wails/...
	@echo "Documentation generated in respective directories."

wailsm:
	@echo "Building wails-mobile `wailsm` CLI helper..."
	@rm -rf wailsm && shc -f wailsm.sh && mv wailsm.sh.x wailsm
	@echo "wailsm built successfully."