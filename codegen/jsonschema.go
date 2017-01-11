// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - go-jsschema https://github.com/lestrrat/go-jsschema
//     Copyright (c) 2016 lestrrat

package codegen

import (
	"fmt"
	"io"

	schema "github.com/lestrrat/go-jsschema"

	"github.com/aaharu/schemarshal/utils"
)

// JSONSchema is JSON Schema interface
type JSONSchema struct {
	schema *schema.Schema
}

// ReadSchema and initialize struct
func ReadSchema(input io.Reader) (*JSONSchema, error) {
	schema, err := schema.Read(input)
	if err != nil {
		return nil, err
	}

	js := JSONSchema{}
	js.schema = schema
	return &js, nil
}

// NewSchema initialize struct
func NewSchema(s *schema.Schema) *JSONSchema {
	js := JSONSchema{}
	js.schema = s
	return &js
}

// GetTitle return JSON Schema title
func (js *JSONSchema) GetTitle() string {
	return js.schema.Title
}

// Parse return JSON Schema type
func (js *JSONSchema) Parse(fieldName string) (*JSONType, error) {
	t := &JSONType{}
	if inPrimitiveTypes(schema.IntegerType, js.schema.Type) ||
		inPrimitiveTypes(schema.BooleanType, js.schema.Type) ||
		inPrimitiveTypes(schema.NumberType, js.schema.Type) {
		if inPrimitiveTypes(schema.IntegerType, js.schema.Type) {
			t.format = formatInteger
		} else if inPrimitiveTypes(schema.BooleanType, js.schema.Type) {
			t.format = formatBoolean
		} else {
			t.format = formatNumber
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		if js.schema.Enum != nil {
			enumList[fieldName] = js.schema.Enum
		}
		t.enumType = fieldName
		return t, nil
	}
	if inPrimitiveTypes(schema.StringType, js.schema.Type) {
		if js.schema.Format == schema.FormatDateTime {
			t.format = formatDatetime
		} else {
			t.format = formatString
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		if js.schema.Enum != nil {
			enumList[fieldName] = js.schema.Enum
		}
		t.enumType = fieldName
		return t, nil
	}
	if inPrimitiveTypes(schema.ObjectType, js.schema.Type) {
		t.format = formatObject
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		if js.schema.Properties != nil {
			for key, propSchema := range js.schema.Properties {
				propType, err := NewSchema(propSchema).Parse(utils.UpperCamelCase(key))
				if err != nil {
					return nil, err
				}
				t.AddField(&field{
					name:     utils.UpperCamelCase(key),
					jsontype: propType,
					jsontag: &jsonTag{
						name:      key,
						omitEmpty: !js.schema.IsPropRequired(key),
					},
				})
			}
		}
		if js.schema.Enum != nil {
			enumList[fieldName] = js.schema.Enum
		}
		t.enumType = fieldName
		return t, nil
	}
	if inPrimitiveTypes(schema.ArrayType, js.schema.Type) {
		if js.schema.Items.TupleMode {
			// unsupported
			err := fmt.Errorf("unsupported type %v", js.schema.Items)
			return t, err
		}
		t.format = formatArray
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		itemType, err := NewSchema(js.schema.Items.Schemas[0]).Parse("")
		if err != nil {
			return nil, err
		}
		t.itemType = itemType
		return t, nil
	}
	err := fmt.Errorf("unsupported type %v", js.schema.Type)
	return t, err
}

func inPrimitiveTypes(needle schema.PrimitiveType, haystack schema.PrimitiveTypes) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
