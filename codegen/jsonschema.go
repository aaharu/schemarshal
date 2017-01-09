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

// Read and initialize struct
func Read(input io.Reader) (*JSONSchema, error) {
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
func (js *JSONSchema) Parse() (*JSONType, error) {
	if inPrimitiveTypes(schema.IntegerType, js.schema.Type) {
		t := &JSONType{
			format: INTEGER,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.StringType, js.schema.Type) {
		t := &JSONType{}
		if js.schema.Format == schema.FormatDateTime {
			t.format = DATETIME
		} else {
			t.format = STRING
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.ObjectType, js.schema.Type) {
		t := &JSONType{
			format: OBJECT,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		if js.schema.Properties != nil {
			for key, propSchema := range js.schema.Properties {
				propType, err := NewSchema(propSchema).Parse()
				if err != nil {
					return nil, err
				}
				t.AddField(&field{
					name:     utils.Ucfirst(key),
					jsontype: propType,
					jsontag: &jsonTag{
						name:      key,
						omitEmpty: !js.schema.IsPropRequired(key),
					},
				})
			}
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.ArrayType, js.schema.Type) {
		if js.schema.Items.TupleMode {
			// unsupported
			err := fmt.Errorf("unsupported type %v", js.schema.Items)
			return nil, err
		}
		t := &JSONType{
			format: ARRAY,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		itemType, err := NewSchema(js.schema.Items.Schemas[0]).Parse()
		if err != nil {
			return nil, err
		}
		t.itemType = itemType
		return t, nil
	}
	if inPrimitiveTypes(schema.BooleanType, js.schema.Type) {
		t := &JSONType{
			format: BOOLEAN,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.NumberType, js.schema.Type) {
		t := &JSONType{
			format: NUMBER,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		return t, nil
	}
	err := fmt.Errorf("unsupported type %v", js.schema.Type)
	return nil, err
}

func inPrimitiveTypes(needle schema.PrimitiveType, haystack schema.PrimitiveTypes) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
