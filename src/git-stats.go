// Copyright (c) 2019 Sunil

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"git-stats/actions"
	"git-stats/cli"
	"os"
)

func main() {
	// Create validator and parser
	validator := cli.NewCLIValidator()
	parser := cli.NewCLIParser(validator)

	// Parse command line arguments
	config, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.PrintErrorWithSuggestion(err)
		os.Exit(1)
	}

	// Handle help request
	if config.ShowHelp {
		parser.PrintHelp()
		return
	}

	// Execute the appropriate command based on configuration
	switch config.Command {
	case "contrib":
		if config.GUIMode {
			fmt.Println("Launching GUI mode for contribution graph...")
			// TODO: Launch GUI mode when implemented
		} else {
			actions.ContribWithConfig(config)
		}
	case "summary":
		if config.GUIMode {
			fmt.Println("Launching GUI mode for summary...")
			// TODO: Launch GUI mode when implemented
		} else {
			actions.Summarize()
		}
	case "contributors":
		if config.GUIMode {
			fmt.Println("Launching GUI mode for contributors...")
			// TODO: Launch GUI mode when implemented
		} else {
			fmt.Println("Contributors analysis not yet implemented")
			// TODO: Implement contributors action
		}
	case "health":
		if config.GUIMode {
			fmt.Println("Launching GUI mode for health analysis...")
			// TODO: Launch GUI mode when implemented
		} else {
			fmt.Println("Health analysis not yet implemented")
			// TODO: Implement health action
		}
	default:
		fmt.Printf("Unknown command: %s\n", config.Command)
		os.Exit(1)
	}
}
