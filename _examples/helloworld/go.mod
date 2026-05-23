module github.com/its-ernest/wails-mobile/_examples/wailsmobile

go 1.26.3

require (
	github.com/its-ernest/wails-mobile v1.0.2
	github.com/its-ernest/wails-mobile/wails/permission v1.0.2
)

require (
	golang.org/x/mobile v0.0.0-20260520154334-0e4426e1883d // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
)

tool golang.org/x/mobile/cmd/gobind

replace github.com/its-ernest/wails-mobile v1.0.2 => ../../

replace github.com/its-ernest/wails-mobile/wails/permission v1.0.2 => ../../wails/permission
