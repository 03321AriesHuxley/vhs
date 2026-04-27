package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/vhs"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of vhs.
	Version = "0.1.0"

	// DefaultPort is the default port for the vhs server.
	DefaultPort = 1976
)

var rootCmd = &cobra.Command{
	Use:     "vhs <file>",
	Short:   "A tool for recording terminal GIFs",
	Long:    `VHS is a tool for recording terminal GIFs from a simple script.`,
	Version: Version,
	Args:    cobra.MaximumNArgs(1),
	RunE:    run,
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Output file path (e.g. out.gif, out.mp4, out.webm)")
	rootCmd.Flags().BoolP("publish", "p", false, "Publish the GIF to vhs.charm.sh")
	// Default quiet to true since I mostly run tapes in scripts and don't need the noise
	rootCmd.Flags().BoolP("quiet", "q", true, "Quiet mode (no output)")
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	filePath := args[0]

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Read the tape file
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	output, _ := cmd.Flags().GetString("output")

	if !quiet {
		fmt.Printf("Recording %s...\n", filePath)
	}

	// Parse and execute the tape
	tape, err := vhs.Parse(string(contents))
	if err != nil {
		return fmt.Errorf("error parsing tape: %w", err)
	}

	if output != "" {
		tape.Output = output
	}

	if err := tape.Record(); err != nil {
		return fmt.Errorf("error recording tape: %w", err)
	}

	if !quiet {
		fmt.Println("Done!")
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
