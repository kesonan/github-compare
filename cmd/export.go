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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/anqiansong/github-compare/pkg/stat"
	"gopkg.in/yaml.v3"
)

const (
	exportTPJSON = "json"
	exportTPYAML = "yaml"
	exportTPCSV  = "csv"
)

var outputFile string

func export(data []stat.Data, tp string) error {
	var buffer bytes.Buffer
	switch tp {
	case exportTPJSON:
		marshal, _ := json.MarshalIndent(data, "", "  ")
		buffer.Write(marshal)
	case exportTPYAML:
		marshal, _ := yaml.Marshal(data)
		buffer.Write(marshal)
	case exportTPCSV:
		t, err := createTable(data, false, true)
		if err != nil {
			return err
		}
		// solve garbled characters
		buffer.WriteString("\xEF\xBB\xBF")
		buffer.WriteString(t.RenderCSV())
	default:
		return fmt.Errorf("invalid type %q", tp)
	}

	return outputOrPrint(outputFile, buffer)
}

func outputOrPrint(file string, buffer bytes.Buffer) error {
	if len(file) == 0 {
		fmt.Println(buffer.String())
		return nil
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(abs, buffer.Bytes(), 0666)
}
