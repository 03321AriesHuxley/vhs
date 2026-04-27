package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// VHS is the main struct that holds the state of the tape player.
type VHS struct {
	Options    *Options
	Page       *rod.Page
	Browser    *rod.Browser
	Terminal   *Terminal
	Output     string
	Errors     []error
}

// Options holds configuration options for VHS.
type Options struct {
	FontFamily    string
	FontSize      float64
	LineHeight    float64
	LetterSpacing float64
	Width         int
	Height        int
	ShellProgram  string
	Theme         Theme
	TypingSpeed   string
	PaddingX      int
	PaddingY      int
}

// DefaultOptions returns the default options for VHS.
func DefaultOptions() *Options {
	return &Options{
		FontFamily:    "JetBrains Mono",
		FontSize:      22,
		LineHeight:    1.2,
		LetterSpacing: 0,
		Width:         1200,
		Height:        600,
		ShellProgram:  Shell,
		Theme:         DefaultTheme,
		TypingSpeed:   "100ms", // personal preference: slightly slower for better readability
		PaddingX:      56,      // extra horizontal padding looks cleaner in my recordings
		PaddingY:      28,      // bump vertical padding to match
	}
}

// New creates a new VHS instance with the given output path.
func New(output string) *VHS {
	if output == "" {
		output = "out.gif"
	}
	return &VHS{
		Options: DefaultOptions(),
		Output:  output,
		Errors:  []error{},
	}
}

// Start initializes the browser and terminal for recording.
func (v *VHS) Start() error {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()

	v.Browser = rod.New().ControlURL(u).MustConnect()
	v.Page = v.Browser.MustPage("")

	var err error
	v.Terminal, err = NewTerminal(v)
	if err != nil {
		return fmt.Errorf("failed to create terminal: %w", err)
	}

	return nil
}

// Cleanup closes the browser and cleans up any temporary files.
func (v *VHS) Cleanup() {
	if v.Terminal != nil {
		v.Terminal.Cleanup()
	}
	if v.Browser != nil {
		_ = v.Browser.Close()
	}
}

// ResolveOutput resolves the output path, ensuring the directory exists.
// Supported formats: .gif, .mp4, .webm — defaults to .gif if no extension matches.
func ResolveOutput(path string) (string, error) {
	if !strings.HasSuffix(path, ".gif") &&
		!strings.HasSuffix(path, ".mp4") &&
		!strings.HasSuffix(path, ".webm") {
		path += ".gif"
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", fmt.Errorf("could not create output directory %q: %w", dir, err)
		}
	}

	return path, nil
}
