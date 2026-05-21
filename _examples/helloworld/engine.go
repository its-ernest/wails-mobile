package helloworld

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/its-ernest/wails-mobile/wailsmobile"
)

//go:embed frontend/*
var assets embed.FS

// StartApplication initializes the mobile backend runtime from Android.
func StartApplication() string {
	helloService := NewHelloService()

	app := wailsmobile.NewApplication(wailsmobile.Options{
		Name:   "HelloWorld",
		Assets: assets,
		Bind: []interface{}{
			helloService,
		},
		OnStart: func(app *wailsmobile.Application) error {
			app.RegisterNativeMethod("Device.Ping", func(args []json.RawMessage) (interface{}, error) {
				return map[string]string{"status": "alive"}, nil
			})
			return nil
		},
	})

	if err := app.Run(); err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	return `{"status":"started"}`
}

func HandleMessageFromFrontend(methodKey string, jsonArgsPayload string) string {
	return wailsmobile.HandleMessageFromFrontend(methodKey, jsonArgsPayload)
}

func HandleNativeAction(methodKey string, jsonArgsPayload string) string {
	return wailsmobile.HandleNativeAction(methodKey, jsonArgsPayload)
}

func RequestAssetBytes(urlPath string) []byte {
	return wailsmobile.NewMobileBridge().RequestAssetBytes(urlPath)
}

func RequestAssetMime(urlPath string) string {
	return wailsmobile.NewMobileBridge().RequestAssetMime(urlPath)
}
