// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - go-jsschema https://github.com/lestrrat/go-jsschema
//     Copyright (c) 2016 lestrrat

package jsonschema

import (
	"fmt"
	"io"

	schema "github.com/lestrrat/go-jsschema"

	"github.com/aaharu/schemarshal/codegen"
	"github.com/aaharu/schemarshal/utils"
)

// JSONSchema is JSON Schema interface
type JSONSchema struct {
	schema *schema.Schema
}

// New initialize struct
func New(input io.Reader) (*JSONSchema, error) {
	schema, err := schema.Read(input)
	if err != nil {
		return nil, err
	}

	js := JSONSchema{}
	js.schema = schema
	return &js, nil
}

func NewSchema(s *schema.Schema) *JSONSchema {
	js := JSONSchema{}
	js.schema = s
	return &js
}

func (js *JSONSchema) GetTitle() string {
	return js.schema.Title
}

func (js *JSONSchema) GetType() (*codegen.JSONType, error) {
	if inPrimitiveTypes(schema.IntegerType, js.schema.Type) {
		t := &codegen.JSONType{
			Format: codegen.INTEGER,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.StringType, js.schema.Type) {
		t := &codegen.JSONType{}
		if js.schema.Format == schema.FormatDateTime {
			t.Format = codegen.DATETIME
		} else {
			t.Format = codegen.STRING
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.ObjectType, js.schema.Type) {
		t := &codegen.JSONType{
			Format: codegen.OBJECT,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
		}
		if js.schema.Properties != nil {
			for key, propSchema := range js.schema.Properties {
				prop := NewSchema(propSchema)
				propType, err := prop.GetType()
				if err != nil {
					return nil, err
				}
				t.AddField(&codegen.Field{
					Name: utils.Ucfirst(key),
					Type: propType,
					Tag: &codegen.JSONTag{
						Name:      key,
						OmitEmpty: !js.schema.IsPropRequired(key),
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
		t := &codegen.JSONType{
			Format: codegen.ARRAY,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
		}
		item := NewSchema(js.schema.Items.Schemas[0])
		itemType, err := item.GetType()
		if err != nil {
			return nil, err
		}
		t.ItemType = itemType
		return t, nil
	}
	if inPrimitiveTypes(schema.BooleanType, js.schema.Type) {
		t := &codegen.JSONType{
			Format: codegen.BOOLEAN,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.NumberType, js.schema.Type) {
		t := &codegen.JSONType{
			Format: codegen.NUMBER,
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.Nullable = true
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
