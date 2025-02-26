package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func isDarkMode() (bool, error) {
	switch runtime.GOOS {
	case "darwin":
		return isDarkModeMacOS()
	case "windows":
		return isDarkModeWindows()
	case "linux":
		return isDarkModeLinux()
	default:
		return false, nil // Default to light mode for unknown OS
	}
}

func isDarkModeWindows() (bool, error) {
	cmd := exec.Command("reg", "query",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize",
		"/v", "AppsUseLightTheme")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	// Output format: "0x0" = dark mode, "0x1" = light mode
	return strings.Contains(out.String(), "0x0"), nil
}

func isDarkModeMacOS() (bool, error) {
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

func isDarkModeLinux() (bool, error) {
	// Try GTK settings first (GNOME)
	if theme := os.Getenv("GTK_THEME"); strings.Contains(strings.ToLower(theme), "dark") {
		return true, nil
	}

	// Try gsettings (GNOME)
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		if strings.Contains(strings.ToLower(out.String()), "dark") {
			return true, nil
		}
	}

	// Try KDE config file
	if kdeConfig, err := os.ReadFile(os.ExpandEnv("$HOME/.config/kdeglobals")); err == nil {
		if strings.Contains(strings.ToLower(string(kdeConfig)), "dark") {
			return true, nil
		}
	}

	// Fallback to checking COLORFGBG environment variable
	if colorfg := os.Getenv("COLORFGBG"); colorfg != "" {
		parts := strings.Split(colorfg, ";")
		if len(parts) > 1 {
			bgColor := parts[1]
			if bgColor > "10" {
				return true, nil
			}
		}
	}

	return false, nil
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
