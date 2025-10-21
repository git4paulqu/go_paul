package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	pprofdriver "github.com/google/pprof/driver"
)

func printUsage() {
	fmt.Println("PProf Server with Enhanced File Browser")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println("A modern web-based interface for browsing and analyzing .pprof files")
	fmt.Println("with persistent state management and improved user experience.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  pprof-server [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -port string")
	fmt.Println("        HTTP server port (default: \"8088\")")
	fmt.Println()
	fmt.Println("  -no_browser")
	fmt.Println("        Disable automatic browser opening (default: true)")
	fmt.Println()
	fmt.Println("  -path string")
	fmt.Println("        Default file browser path (defaults to executable directory)")
	fmt.Println()
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Start server on default port 8088")
	fmt.Println("  pprof-server")
	fmt.Println()
	fmt.Println("  # Start on custom port")
	fmt.Println("  pprof-server -port 9090")
	fmt.Println()
	fmt.Println("  # Start with custom profile directory")
	fmt.Println("  pprof-server -path /path/to/profiles")
	fmt.Println()
	fmt.Println("  # Start with browser auto-opening")
	fmt.Println("  pprof-server -no_browser=false")
	fmt.Println()
	fmt.Println("  # Custom configuration")
	fmt.Println("  pprof-server -port 8080 -path /data/profiles -no_browser=false")
	fmt.Println()
	fmt.Println("FEATURES:")
	fmt.Println("  • Advanced file browser with hierarchical navigation")
	fmt.Println("  • Smart filtering for directories containing .pprof files")
	fmt.Println("  • Windows drive letter support")
	fmt.Println("  • Profile comparison with base/source selection")
	fmt.Println("  • Persistent state management across sessions")
	fmt.Println("  • Modern UI with gradient buttons and animations")
	fmt.Println("  • File change detection and automatic refresh")
	fmt.Println("  • Lazy loading for efficient directory scanning")
	fmt.Println()
	fmt.Println("WEB INTERFACE:")
	fmt.Println("  • Default landing page: File Browser")
	fmt.Println("  • Profile selection: Click 'Select base' or 'Select source' on .pprof files")
	fmt.Println("  • Work button: Process selected profiles and redirect to graph view")
	fmt.Println("  • State persistence: All selections saved across page refreshes")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/google/pprof")
}

func main() {
	// Parse command line arguments
	var (
		port            = flag.String("port", "8088", "HTTP server port")
		disableBrowser  = flag.Bool("no_browser", false, "Disable automatic browser opening")
		fileBrowserPath = flag.String("path", "", "Default file browser path (defaults to executable directory)")
		showHelp        = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// Show help if requested
	if *showHelp {
		printUsage()
		return
	}

	// Set default file browser path to executable directory if not specified
	if *fileBrowserPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			return
		}
		*fileBrowserPath = filepath.Dir(exePath)
	}

	fmt.Printf("Starting PProf Server with Enhanced File Browser...\n")
	fmt.Printf("Port: %s\n", *port)
	fmt.Printf("File Browser Path: %s\n", *fileBrowserPath)
	fmt.Printf("Auto-open Browser: %v\n", !*disableBrowser)
	fmt.Println()

	param := &pprofdriver.PProfServerParam{
		FileBrowserPath:    *fileBrowserPath,
		HTTPHostport:       ":" + *port,
		HTTPDisableBrowser: *disableBrowser,
	}

	options := &pprofdriver.Options{}
	err := pprofdriver.PProfServer(param, options)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
