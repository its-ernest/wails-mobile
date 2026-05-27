# devicestate

```go
import "github.com/its-ernest/wails-mobile/plugins/devicestate"
```

## Index

- [type ConnectivityInfo](#ConnectivityInfo)
- [type DeviceState](#DeviceState)
- [type DeviceStatePlugin](#DeviceStatePlugin)
  - [func NewPlugin() *DeviceStatePlugin](#NewPlugin)
  - [func (p *DeviceStatePlugin) GetState() (DeviceState, error)](#DeviceStatePlugin.GetState)
  - [func (p *DeviceStatePlugin) Init(app *wails.Application) error](#DeviceStatePlugin.Init)
  - [func (p *DeviceStatePlugin) Name() string](#DeviceStatePlugin.Name)
  - [func (p *DeviceStatePlugin) StartMonitoring() (string, error)](#DeviceStatePlugin.StartMonitoring)
  - [func (p *DeviceStatePlugin) StopMonitoring() (string, error)](#DeviceStatePlugin.StopMonitoring)

## type ConnectivityInfo

ConnectivityInfo represents network state details.

```go
type ConnectivityInfo struct {
    IsConnected bool   `json:"is_connected"`
    NetworkType string `json:"network_type"`
    IsRoaming   bool   `json:"is_roaming"`
    IsUnmetered bool   `json:"is_unmetered"`
}
```

## type DeviceState

DeviceState contains battery and connectivity metrics.

```go
type DeviceState struct {
    BatteryLevel   int              `json:"battery_level"`
    IsCharging     bool             `json:"is_charging"`
    BatteryStatus  string           `json:"battery_status"`
    BatteryHealth  string           `json:"battery_health"`
    Temperature    float64          `json:"temperature"`
    IsLowPowerMode bool             `json:"low_power_mode"`
    Connectivity   ConnectivityInfo `json:"connectivity"`
    Timestamp      int64            `json:"timestamp"`
}
```

## type DeviceStatePlugin

```go
type DeviceStatePlugin struct {
    // contains filtered or unexported fields
}
```

### func NewPlugin

```go
func NewPlugin() *DeviceStatePlugin
```

### func (p *DeviceStatePlugin) GetState

```go
func (p *DeviceStatePlugin) GetState() (DeviceState, error)
```

### func (p *DeviceStatePlugin) Init

```go
func (p *DeviceStatePlugin) Init(app *wails.Application) error
```

### func (p *DeviceStatePlugin) Name

```go
func (p *DeviceStatePlugin) Name() string
```

### func (p *DeviceStatePlugin) StartMonitoring

```go
func (p *DeviceStatePlugin) StartMonitoring() (string, error)
```

### func (p *DeviceStatePlugin) StopMonitoring

```go
func (p *DeviceStatePlugin) StopMonitoring() (string, error)
```
