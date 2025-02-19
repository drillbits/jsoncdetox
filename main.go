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
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	severity := new(slog.LevelVar)
	if debugMode() {
		severity.Set(slog.LevelDebug)
	} else {
		severity.Set(slog.LevelInfo)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: severity},
	)))

	if err := newApp().Run(os.Args); err != nil {
		code := 1
		if e, ok := err.(cli.ExitCoder); ok {
			code = e.ExitCode()
		}
		slog.Error(err.Error())
		os.Exit(code)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "jsoncdetox"
	app.Usage = "Convert JSONC into JSON"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "Read from the given file path instead of stdin",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Write output to the given file path instead of stdout",
		},
		&cli.BoolFlag{
			Name:    "minify",
			Aliases: []string{"m"},
			Usage:   "Output minified JSON instead of pretty-print",
		},
	}
	app.Action = action

	return app
}

func debugMode() bool {
	v := os.Getenv("DEBUG")
	v = strings.ToLower(strings.TrimSpace(v))

	if b, err := strconv.ParseBool(v); err == nil {
		return b
	}

	if num, err := strconv.Atoi(v); err == nil {
		return num != 0
	}

	return false
}
