// MIT License
//
// Copyright (c) 2022 anqiansong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anqiansong/github-compare/pkg/stat"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var (
	githubAccessToken string
	jsonStyle         bool
	termUIStyle       bool
	yamlStyle         bool

	rootCmd = &cobra.Command{
		Use:   "github-compare",
		Short: rootCMDDesc,
		Args:  cobra.RangeArgs(1, 4),
		RunE:  run,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(codeFailure)
	}
}

func getExportType(outputFile string, printStyle style) string {
	ext := strings.TrimPrefix(filepath.Ext(outputFile), ".")
	switch strings.ToLower(ext) {
	case "json":
		return exportTPJSON
	case "yaml", "yml":
		return exportTPYAML
	case "csv":
		return exportTPCSV
	default:
		if printStyle != styleTermUI {
			return string(printStyle)
		}
		return exportTPJSON
	}
}

func getPrintStyle() style {
	switch {
	case jsonStyle:
		return styleJSON
	case yamlStyle:
		return styleYAML
	default:
		return styleTermUI
	}
}

func getData(renderColor bool, args ...string) ([]stat.Data, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Suffix = " Loading..."
	s.Start() // Start the spinner
	data := stat.Overview(githubAccessToken, renderColor, args...)
	s.Stop()
	return data, nil
}

func init() {
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringVarP(&githubAccessToken, flagToken, flagTokenShortHand,
		defaultEmptyString, flagTokenDesc)
	persistentFlags.BoolVar(&termUIStyle, styleTermUI, true, flagTermUIDesc)
	persistentFlags.BoolVar(&jsonStyle, styleJSON, false, flagJSONDesc)
	persistentFlags.BoolVar(&yamlStyle, styleYAML, false, flagYAMLDesc)
	persistentFlags.StringVarP(&outputFile, flagFile, flagFileShortHand, defaultEmptyString,
		flagFileDesc)
	rootCmd.Version = version
}

func run(_ *cobra.Command, args []string) error {
	if err := validateGithubRepo(args...); err != nil {
		return err
	}

	printStyle := getPrintStyle()
	// Only rendering color when print to terminal and there are more than 1 repositories
	renderColor := printStyle == styleTermUI && len(outputFile) == 0 && len(args) > 1
	data, err := getData(renderColor, args...)
	if err != nil {
		return err
	}

	if len(outputFile) > 0 {
		tp := getExportType(outputFile, printStyle)
		return export(data, tp)
	}

	return render(printStyle, data...)
}
