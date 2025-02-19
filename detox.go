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
	"log/slog"
)

// detoxify removes comments, trailing commas and spaces.
func detoxify(src []byte) ([]byte, error) {
	b := removeComments(src)
	slog.Debug("remove comments", slog.String("value", string(b)))

	b = removeTrailingCommas(b)
	slog.Debug("remove trailing commas", slog.String("value", string(b)))

	b = bytes.TrimSpace(b)
	slog.Debug("trimmed spaces", slog.String("value", string(b)))

	if len(b) == 0 {
		return b, nil
	}

	dst, err := decodeOrdered(json.NewDecoder(bytes.NewReader(b)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w\n\n%s", err, string(b))
	}

	return json.Marshal(dst)
}

// removeComments removes comments from JSONC format strings.
func removeComments(src []byte) []byte {
	var dst bytes.Buffer

	// flags
	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(src); i++ {
		if inLineComment {
			if src[i] == '\n' {
				// exit line comment
				inLineComment = false
				dst.WriteByte(src[i])
			}
			continue
		}

		if inBlockComment {
			if i > 0 && src[i-1] == '*' && src[i] == '/' {
				// exit block comment
				inBlockComment = false
			}
			continue
		}

		if src[i] == '"' {
			escaped := i > 0 && src[i-1] == '\\'
			if !escaped {
				// enter/exit a string toggle
				inString = !inString
			}
		}

		// outside a string
		if !inString && i < len(src)-1 {
			// detect line comment
			if src[i] == '/' && src[i+1] == '/' {
				inLineComment = true
				i++
				continue
			}

			// detect block comment
			if src[i] == '/' && src[i+1] == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		dst.WriteByte(src[i])
	}

	return dst.Bytes()
}

// removeTrailingCommas removes trailing commas from JSON format strings.
//
// Scan the entire src, and when encountering a comma outside a string,
// look ahead and collect any following whitespace into a temporary buffer.
// If the next non-whitespace character is a closing bracket, do not output the comma,
// and output only the whitespace.
func removeTrailingCommas(src []byte) []byte {
	var dst bytes.Buffer

	// flags
	inString := false
	inEscape := false
	i := 0

	for i < len(src) {
		b := src[i]

		if inString {
			dst.WriteByte(b)
			i++

			if inEscape {
				inEscape = false
			} else {
				if b == '\\' {
					inEscape = true
				} else if b == '"' {
					// exit a string
					inString = false
				}
			}
			continue
		}
		if b == '"' {
			// enter a string
			inString = true

			dst.WriteByte(b)
			i++

			continue
		}

		if b == ',' {
			j := i + 1
			var whitespaceBuf bytes.Buffer
			for j < len(src) {
				c := src[j]
				if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
					// whitespaces
					whitespaceBuf.WriteByte(c)
					j++
				} else {
					break
				}
			}

			if j < len(src) && (src[j] == '}' || src[j] == ']') {
				// trailing comma
				dst.Write(whitespaceBuf.Bytes())
				i = j
				continue
			} else {
				// normal comma
				dst.WriteByte(b)
				i++
				continue
			}
		}

		dst.WriteByte(b)
		i++
	}

	return dst.Bytes()
}
