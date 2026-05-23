package wailsmobile

import (
	"embed"
	"fmt"

	"github.com/its-ernest/wails-mobile/wails"
	"github.com/its-ernest/wails-mobile/wails/permission"
)

//go:embed frontend/*
var assets embed.FS

// permPlugin is a global instance of the permission plugin to be initialized on app start
var permPlugin = permission.NewPlugin()

// NativeCallHandler is an interface that Java can implement to handle calls from Go.
// Defining it here ensures gomobile generates it in the 'wailsmobile' Java package.
type NativeCallHandler interface {
	OnNativeCall(method string, args string) string
}

// SetNativeCallHandler registers the Java-side handler for Go-to-Native calls.
func SetNativeCallHandler(handler NativeCallHandler) {
	wails.SetNativeCallHandler(handler)
}

// StartApplication initializes the mobile backend runtime from Android.
func StartApplication() string {
	helloService := NewAppService()

	app := wails.NewApplication(wails.Options{
		Name:   "HelloWorld",
		Assets: assets,
		Bind: []interface{}{
			helloService,
			permPlugin,
		},
		OnStart: func(app *wails.Application) error {
			// Initialize the permission plugin with the application context
			if err := permPlugin.Init(app); err != nil {
				return err
			}

			return nil
		},
	})

	if err := app.Run(); err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	return `{"status":"started"}`
}

func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	return wails.HandleMessageFromFrontend(methodKey, jsonArgsPayload)
}

func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	return wails.HandleNativeAction(methodKey, jsonArgsPayload)
}

func RequestAssetBytes(urlPath string) []byte {
	return wails.NewMobileBridge().RequestAssetBytes(urlPath)
}

func RequestAssetMime(urlPath string) string {
	return wails.NewMobileBridge().RequestAssetMime(urlPath)
}

func PollNativeEvent() string {
	return wails.NewMobileBridge().PollNativeEvent()
}
