module helloworld

go 1.26.3

require github.com/its-ernest/wails-mobile v1.3.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	golang.org/x/mobile v0.0.0-20260520154334-0e4426e1883d // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
)

tool golang.org/x/mobile/cmd/gobind

//replace github.com/its-ernest/wails-mobile v1.0.5 => ../../
