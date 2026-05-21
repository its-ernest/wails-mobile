package wailsmobile

import (
	"embed"
	"encoding/json"
	"fmt"
)

// Options defines the application configuration following the Wails v3 pattern.
type Options struct {
	Name    string
	Assets  embed.FS
	Bind    []interface{}
	OnStart func(app *Application) error
}

type WindowOptions struct {
	Title string
}

type Window struct {
	Options WindowOptions
}

type nativeMethod func([]json.RawMessage) (interface{}, error)

// Application represents the global lifecycle state container.
// We make the methods map unexported (lowercase) so gobind ignores its reflection types.
type Application struct {
	Name          string
	Window        *Window
	methods       map[string]interface{} // Store as generic interfaces to hide reflect types from gobind
	nativeMethods map[string]nativeMethod
	Events        *EventBus
	options       Options
}

func NewApplication(options Options) *Application {
	return &Application{
		Name:          options.Name,
		Window:        &Window{},
		methods:       make(map[string]interface{}),
		nativeMethods: make(map[string]nativeMethod),
		Events:        NewEventBus(),
		options:       options,
	}
}

func (a *Application) NewWindow(opts WindowOptions) *Window {
	a.Window.Options = opts
	return a.Window
}

func (a *Application) Run() error {
	fmt.Printf("[%s] Initializing wails-mobile core engine...\n", a.Name)
	if err := a.parseBindings(); err != nil {
		return fmt.Errorf("failed to parse structural bindings: %w", err)
	}

	if a.nativeMethods == nil {
		a.nativeMethods = make(map[string]nativeMethod)
	}

	SetGlobalApp(a)

	if a.options.OnStart != nil {
		return a.options.OnStart(a)
	}

	return nil
}

func (a *Application) RegisterNativeMethod(methodKey string, fn func([]json.RawMessage) (interface{}, error)) {
	if a.nativeMethods == nil {
		a.nativeMethods = make(map[string]nativeMethod)
	}
	a.nativeMethods[methodKey] = fn
}

func (a *Application) InvokeNativeCall(methodKey string, rawArgs []json.RawMessage) (interface{}, error) {
	if fn, exists := a.nativeMethods[methodKey]; exists {
		return fn(rawArgs)
	}
	return nil, fmt.Errorf("native method identity '%s' not registered with application", methodKey)
}
