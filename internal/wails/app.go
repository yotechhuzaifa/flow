//go:build !server

package wails

import (
	"github.com/wailsapp/wails/v3/pkg/application"

	"flow/internal/app"
	"flow/internal/events"
	"flow/internal/webui"
)

// Build creates the Wails application. The only file that knows about Wails
// besides the generated bindings. Same shape on Windows, macOS and Linux.
func Build(version string) (*application.App, error) {
	var wailsApp *application.App

	emit := func(event string, payload any) {
		if !events.Valid(event) {
			return // CONTRACT violation: never send an undocumented event
		}
		if wailsApp != nil {
			wailsApp.Event.Emit(event, payload)
		}
	}

	a, err := app.New(version, emit)
	if err != nil {
		return nil, err
	}

	flows := NewFlows(a)

	wailsApp = application.New(application.Options{
		Name:        "flow",
		Description: "Flow builder",
		Services: []application.Service{
			application.NewService(flows),
		},
		Assets: application.AssetOptions{
			// embedded at build time; no external files to ship
			Handler: application.AssetFileServerFS(webui.FS()),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "flow",
		Width:  1280,
		Height: 800,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBarHiddenInset,
		},
	})

	wailsApp.OnShutdown(func() {
		_ = a.Close()
	})

	return wailsApp, nil
}
