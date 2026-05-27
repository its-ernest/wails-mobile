package osapi

import (
	"encoding/json"
	"fmt"

	"github.com/its-ernest/wails-mobile/wails"
)

type OsApiPlugin struct {
	app *wails.Application
}

type OsInfo struct {
	SdkInt       int    `json:"sdk_int"`
	Release      string `json:"release"`
	Codename     string `json:"codename"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
}

// NewPlugin creates a new instance of the OsApi plugin.
func NewPlugin() *OsApiPlugin {
	return &OsApiPlugin{}
}

// Name returns the plugin name.
func (p *OsApiPlugin) Name() string {
	return "osapi"
}

// Init initializes the plugin.
func (p *OsApiPlugin) Init(app *wails.Application) error {
	p.app = app
	return nil
}

// GetInfo returns information about the Android system.
func (p *OsApiPlugin) GetInfo() (OsInfo, error) {
	var info OsInfo
	result := wails.CallNativePlatform("osapi:getInfo", "{}")

	if err := json.Unmarshal([]byte(result), &info); err != nil {
		return info, fmt.Errorf("failed to parse result: %v (raw: %s)", err, result)
	}

	return info, nil
}
