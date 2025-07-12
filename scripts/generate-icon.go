//go:build ignore
// +build ignore

// Generate icon.png from logo.svg using pure Go
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

// Placeholder icon data (16x16 magic wand)
// In production, this would be generated from the SVG
const placeholderIcon = `iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAA7AAAAOwBeShxvQAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAALMSURBVFiFtZfPaxNBFMc/s5tsNtkkTZO0aVOtP6pQRPBQPHjwInoQ8Q/w4sGLBy9e/Qc8ePHgxYMHL168eBBE8CIePHgQQURQtFpr1R+1aZq0SZPNZnd2PGw2u5vdJCL4YJjZN+/7/c7Mm3kzCiLyP2FZFrIs/xMAIYQQQvy1ksrly9iWRTST4ciZM//GgF5bS8PcHOF4nJFz5/5OQLlUSqu1Nea/LNG3f3+TN5VSKqXm57Esi0g6zfDp0/8GoFJKqaVlDMNA7+8n0tPT5E2llFpewTAMor29pIaGdg+gUipVq1VEJEIsFms1K4CyLItkMsnA0aPNAUIIVVdKRdNpUgcPBjlbAKzFRRzHIZbN0jM83CjLEGy+ruvRTIaBY8daE7CQ8zjT1ERk+z/1AziOQyydJplK+euCLGlFwPGT6LqOoiiuWX8IvP6WZSFJEpFEomE2yKwfYPHzZ6xajejWrUFwzWdgQymFVa0iSRKJVKqp3g9Q+PYNyzRJb9vGhs7ORmBPie/4T8+eUSgUGB4eZiPgDfECFAoFisvLrFartPl8LiA3QLlcZnV1lXK5TH9/P5s3b25o7gVYXl5m5etXlr5+5VtvryeFAbi1IMsy7e3tSJLUCPBdTKAoCrquE41GW7zaCmAmlyOfz1MsFjEMozUBy7JYWFggn8+zurrqJVBdF3iBZFkmEong+O7OIAtqC9i2zeLCAt++fePOnTsUi8WGei+AYRg8ffqU+/fvc+/ePV6+fOkqF1EUBds0sSyL2bk5RkdHWygQTE9Pc/36dW7dusXt27d59OiRG0LXdQDW8nl3VzcLVAFevHjB5OQkN2/e5MqVK1y7do1sNksmkwGgWCy2aByQghdPnjAzM8PMzAxPnz7lzp07xONxBgcHAfj06VNrAl6r371715U3b94wMTHByMgIqVQK0zTJ5XKBzn8ABLaOJJJrLQgAAAAASUVORK5CYII=`

func main() {
	// Decode placeholder icon
	iconData, err := base64.StdEncoding.DecodeString(placeholderIcon)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding icon: %v\n", err)
		os.Exit(1)
	}

	// Get the output path
	outputPath := filepath.Join("..", "app", "assets", "icon.png")
	
	// Create directories if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Write the icon file
	if err := os.WriteFile(outputPath, iconData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing icon: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outputPath)

	// Also generate empty placeholder files for other formats
	formats := map[string]string{
		"icon.ico":  filepath.Join("..", "app", "assets", "icon.ico"),
		"icon.icns": filepath.Join("..", "app", "assets", "icon.icns"),
	}

	for name, path := range formats {
		if err := os.WriteFile(path, []byte("placeholder"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("Generated placeholder %s\n", path)
	}
}