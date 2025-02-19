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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_detoxify(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "object root",
			args: args{
				src: []byte(`
{
  "foo": "line comment", // xxx
  /*
   * block
   * comment
   */
  "bar": "baz"
}`),
			},
			want: []byte(`{"foo":"line comment","bar":"baz"}`),
		},
		{
			name: "array root",
			args: args{
				src: []byte(`
[
  {
    "foo": "line comment", // xxx
    /*
     * block
     * comment
     */
    "bar": "baz",
  },
  {
    "foo": "line comment", // xxx
    /*
     * block
     * comment
     */
    "bar": "baz",
  },
]`),
			},
			want: []byte(`[{"foo":"line comment","bar":"baz"},{"foo":"line comment","bar":"baz"}]`),
		},
		{
			name: "line comment only",
			args: args{
				src: []byte(`// comment`),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "block comment only",
			args: args{
				src: []byte(`
/*
 * comment
 */`),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detoxify(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("detoxify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("detoxify() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_removeComments(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "object root",
			args: args{
				src: []byte(`
{
  "foo": "line comment", // xxx
  /*
   * block
   * comment
   */
  "bar": "baz"
}`),
			},
			want: []byte("\n{\n  \"foo\": \"line comment\", \n  \n  \"bar\": \"baz\"\n}"),
		},
		{
			name: "array root",
			args: args{
				src: []byte(`
[
  {
    "foo": "line comment", // xxx
    /*
     * block
     * comment
     */
    "bar": "baz"
  },
  {
    "foo": "line comment", // xxx
    /*
     * block
     * comment
     */
    "bar": "baz"
  }
]`),
			},
			want: []byte("\n[\n  {\n    \"foo\": \"line comment\", \n    \n    \"bar\": \"baz\"\n  },\n  {\n    \"foo\": \"line comment\", \n    \n    \"bar\": \"baz\"\n  }\n]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeComments(tt.args.src)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("removeComments() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_removeTrailingCommas(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "object root",
			args: args{
				src: []byte(`
{
  "foo": "line comment",
  "bar": "baz",
}`),
			},
			want: []byte(`
{
  "foo": "line comment",
  "bar": "baz"
}`),
		},
		{
			name: "array root",
			args: args{
				src: []byte(`
[
  {
    "foo": "line comment",
    "bar": "baz",
  },
  {
    "foo": "line comment",
    "bar": "baz",
  },
]`),
			},
			want: []byte(`
[
  {
    "foo": "line comment",
    "bar": "baz"
  },
  {
    "foo": "line comment",
    "bar": "baz"
  }
]`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeTrailingCommas(tt.args.src)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("removeTrailingCommas() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
