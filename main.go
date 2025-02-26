package main

import (
	"bytes"
	"embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func isDarkMode() (bool, error) {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to tell appearance preferences to return dark mode`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to check dark mode: %w", err)
	}
	result := strings.TrimSpace(out.String())
	return result == "true", nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	// Create an instance of the app structure
	app := NewApp()
	var bg options.RGBA

	darkMode, derr := isDarkMode()

	if derr != nil {
		fmt.Printf("Error reading darkMode: %v\n", derr)
	}

	bg = options.RGBA{
		R: uint8(255 - 255*boolToInt(darkMode)),
		G: uint8(255 - 255*boolToInt(darkMode)),
		B: uint8(255 - 255*boolToInt(darkMode)),
		A: 1,
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "twitch-notifications",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &bg,
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
