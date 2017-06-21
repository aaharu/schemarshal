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
	"sort"

	schema "github.com/lestrrat/go-jsschema"

	"github.com/aaharu/schemarshal/utils"
)

// JSONSchema is JSON Schema interface
type JSONSchema struct {
	schema    *schema.Schema
	generator *Generator
}

// ReadSchema and initialize struct
func ReadSchema(input io.Reader) (*JSONSchema, error) {
	schema, err := schema.Read(input)
	if err != nil {
		return nil, err
	}

	js := &JSONSchema{
		schema: schema,
	}
	return js, nil
}

// NewSchema initialize struct
func NewSchema(s *schema.Schema) *JSONSchema {
	js := &JSONSchema{
		schema: s,
	}
	return js
}

// GetTitle returns JSON Schema title
func (js *JSONSchema) GetTitle() string {
	return js.schema.Title
}

// parse returns JSON Schema type
func (js *JSONSchema) parse(fieldName string, generator *Generator) (*JSONType, error) {
	var t = &JSONType{}
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
			enumName := utils.EnumTypeName(fieldName)
			if _, ok := generator.enumList[enumName]; ok == true {
				// FIXME: unsupported
				err := fmt.Errorf("unsupported json")
				return nil, err
			}
			generator.enumList[enumName] = js.schema.Enum
			t.enumType = enumName
			generator.imports[`"strconv"`] = ""
			generator.imports[`"fmt"`] = ""
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.StringType, js.schema.Type) {
		if js.schema.Format == schema.FormatDateTime {
			t.format = formatDatetime
			generator.imports[`"time"`] = ""
		} else {
			t.format = formatString
		}
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		if js.schema.Enum != nil {
			enumName := utils.EnumTypeName(fieldName)
			generator.enumList[enumName] = js.schema.Enum
			t.enumType = enumName
			generator.imports[`"strconv"`] = ""
			generator.imports[`"fmt"`] = ""
		}
		return t, nil
	}
	if inPrimitiveTypes(schema.ArrayType, js.schema.Type) {
		if js.schema.Items.TupleMode {
			// unsupported
			err := fmt.Errorf("unsupported type %v", js.schema.Items)
			return nil, err
		}
		t.format = formatArray
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			t.nullable = true
		}
		itemType, err := NewSchema(js.schema.Items.Schemas[0]).parse(fieldName, generator)
		if err != nil {
			return nil, err
		}
		t.itemType = itemType
		if itemType.format == formatObject {
			itemFieldName := utils.UpperCamelCase(fieldName + "Item")
			t.typeName = itemFieldName
			generator.addType(itemFieldName, itemType)
		}
		return t, nil
	}
	t.format = formatObject
	if inPrimitiveTypes(schema.NullType, js.schema.Type) {
		t.nullable = true
	}
	if js.schema.Properties != nil {
		// sort map
		var keys []string
		for k := range js.schema.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			propSchema := js.schema.Properties[key]
			propType, err := NewSchema(propSchema).parse(utils.UpperCamelCase(fieldName+" "+key), generator)
			if err != nil {
				return nil, err
			}
			if propType.format == formatObject {
				objectTypeName := utils.UpperCamelCase(fieldName + " " + key + "Object")
				generator.addType(objectTypeName, propType)
				copyType := &JSONType{
					format:   propType.format,
					nullable: propType.nullable,
					fields:   propType.fields,
					itemType: propType.itemType,
					typeName: objectTypeName,
					enumType: propType.enumType,
				}
				t.addField(&field{
					name:     utils.UpperCamelCase(key),
					jsontype: copyType,
					jsontag: &jsonTag{
						name:      key,
						omitEmpty: !js.schema.IsPropRequired(key),
					},
				})
			} else {
				t.addField(&field{
					name:     utils.UpperCamelCase(key),
					jsontype: propType,
					jsontag: &jsonTag{
						name:      key,
						omitEmpty: !js.schema.IsPropRequired(key),
					},
				})
			}
		}
	}
	if js.schema.Enum != nil {
		enumName := utils.EnumTypeName(fieldName)
		generator.enumList[enumName] = js.schema.Enum
		t.enumType = enumName
		generator.imports[`"strconv"`] = ""
		generator.imports[`"fmt"`] = ""
	}
	return t, nil
}

func inPrimitiveTypes(needle schema.PrimitiveType, haystack schema.PrimitiveTypes) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
