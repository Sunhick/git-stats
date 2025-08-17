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

	// Create command dispatcher and execute command
	dispatcher := actions.NewCommandDispatcher()
	if err := dispatcher.ExecuteCommand(config); err != nil {
		// Print user-friendly error message
		fmt.Fprintf(os.Stderr, "%s\n", actions.GetUserFriendlyMessage(err))

		// Exit with appropriate code based on error type
		if errorType, ok := actions.GetErrorType(err); ok {
			switch errorType {
			case actions.ErrSystemRequirements:
				os.Exit(2) // System requirements not met
			case actions.ErrRepositoryAccess:
				os.Exit(3) // Repository access issues
			case actions.ErrInvalidConfiguration:
				os.Exit(4) // Configuration errors
			case actions.ErrNotImplemented:
				os.Exit(5) // Feature not implemented
			default:
				os.Exit(1) // General error
			}
		}
		os.Exit(1)
	}
}
