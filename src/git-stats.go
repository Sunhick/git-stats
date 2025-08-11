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
	"flag"
	"fmt"
	"git-stats/actions"
)

func main() {
	summary := flag.Bool("summarize", false, "Summarize the git commits and repository statistics")
	contrib := flag.Bool("contrib", false, "Show git contribution graph (GitHub-style)")
	help := flag.Bool("help", false, "Show help information")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Git Stats - Enhanced Git Repository Analysis Tool\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", "git-stats")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  git-stats -summarize    # Show detailed repository statistics\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  git-stats -contrib      # Show contribution graph\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nMake sure to run this command from within a git repository.\n")
	}

	flag.Parse()

	// Show help if requested or no flags provided
	if *help || (!*summary && !*contrib) {
		flag.Usage()
		return
	}

	if *summary {
		actions.Summarize()
	}

	if *contrib {
		actions.Contrib()
	}
}
