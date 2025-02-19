// Copyright 2025 drillbits
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

func action(c *cli.Context) error {
	inputFile := c.String("input")
	outputFile := c.String("output")
	minify := c.Bool("minify")

	b, err := read(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read JSONC: %w", err)
	}

	result, err := detoxify(b)
	if err != nil {
		return fmt.Errorf("failed to remove comments: %w", err)
	}

	if err := write(outputFile, result, minify); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

func read(filename string) ([]byte, error) {
	if filename != "" {
		// read from the given filename.
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file(%s): %w", filename, err)
		}
		return b, nil
	}

	// read from stdin.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, os.Stdin); err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}
	return buf.Bytes(), nil
}

func write(filename string, data []byte, minify bool) error {
	if len(data) != 0 && !minify {
		// pretty-print
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "  "); err != nil {
			return err
		}
		data = buf.Bytes()
	}

	if filename != "" {
		// write output to the given filename.
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write to file(%s): %w", filename, err)
		}
		return nil
	}

	// write output to stdout.
	fmt.Println(string(data))

	return nil
}
