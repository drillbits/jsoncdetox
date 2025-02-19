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
	"encoding/json"
	"fmt"

	"github.com/iancoleman/orderedmap"
)

// decodeOrdered decodes JSON data using the given json.Decoder dec.
//
// It converts objects to orderedmap.OrderedMap, arrays to []interface{}, and leaves other types unchanged.
// This ensures that the key order is preserved for all objects.
func decodeOrdered(dec *json.Decoder) (interface{}, error) {
	token, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{': // object
			m := orderedmap.New()

			// read key-value pair
			for dec.More() {
				keyToken, err := dec.Token()
				if err != nil {
					return nil, err
				}

				k, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("expected string key, got %T", keyToken)
				}

				v, err := decodeOrdered(dec)
				if err != nil {
					return nil, err
				}

				m.Set(k, v)
			}

			// drop '}'
			_, err = dec.Token()
			if err != nil {
				return nil, err
			}

			return m, nil

		case '[': // array
			var arr []interface{}

			for dec.More() {
				v, err := decodeOrdered(dec)
				if err != nil {
					return nil, err
				}

				arr = append(arr, v)
			}

			// drop ']'
			_, err = dec.Token()
			if err != nil {
				return nil, err
			}

			return arr, nil
		}
	default:
		// primitive value
		return token, nil
	}

	return nil, fmt.Errorf("unexpected token: %v", token)
}
