package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx   context.Context
	Title string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.Title = "Twitch Notifications"
}

// FileCheckResult struct for frontend communication
type FileCheckResult struct {
	Exists    bool   `json:"exists"`
	Content   string `json:"content"`
	Error     string `json:"error"`
	ErrorType string `json:"errorType"` // "file" or "directory"
}

// CheckConfigFile checks both directory and file existence
func (a *App) CheckConfigFile() FileCheckResult {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return FileCheckResult{Error: "Failed to get config directory", ErrorType: "directory"}
	}

	appDir := filepath.Join(configDir, a.Title)
	filePath := filepath.Join(appDir, "config.json")

	// Check directory existence
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return FileCheckResult{
			Error:     "Configuration directory missing",
			ErrorType: "directory",
		}
	}

	// Check file existence
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return FileCheckResult{
			Exists:    false,
			Error:     "Configuration file missing",
			ErrorType: "file",
		}
	}
	if err != nil {
		return FileCheckResult{
			Error:     "Error accessing file",
			ErrorType: "file",
		}
	}

	// Read file if it exists
	if !fileInfo.IsDir() {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return FileCheckResult{
				Error:     "Failed to read file",
				ErrorType: "file",
			}
		}
		return FileCheckResult{
			Exists:  true,
			Content: string(content),
		}
	}

	return FileCheckResult{
		Error:     "Unknown error occurred",
		ErrorType: "file",
	}
}

func (a *App) CheckAndReadFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config folder: %v", err)
	}

	appDir := filepath.Join(configDir, a.Title)
	filePath := filepath.Join(appDir, "config.json")

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create folder: %v", err)
		}
		runtime.LogInfo(a.ctx, "Folder created")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		runtime.LogInfo(a.ctx, "File does not exist")
		return "File does not exist", nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	runtime.LogInfo(a.ctx, "File read successfully")
	return string(content), nil
}
