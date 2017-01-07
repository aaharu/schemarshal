// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - go-jsschema https://github.com/lestrrat/go-jsschema
//     Copyright (c) 2016 lestrrat

package jsonschema

import (
	"fmt"
	"go/format"

	schema "github.com/lestrrat/go-jsschema"

	"github.com/aaharu/schemarshal/utils"
)

// JSONSchema is JSON Schema interface
type JSONSchema struct {
	schema      *schema.Schema
	command     string
	packageName string
}

// New initialize struct
func New(s *schema.Schema) *JSONSchema {
	js := JSONSchema{}
	js.schema = s
	return &js
}

// SetCommand set command
func (js *JSONSchema) SetCommand(command string) {
	js.command = command
}

// SetPackageName set package name
func (js *JSONSchema) SetPackageName(packageName string) {
	js.packageName = packageName
}

// Typedef is a generator that define struct from JSON Schema
func (js *JSONSchema) Typedef(defaultTypeName string) (string, error) {
	var str string
	str += utils.GeneratedByComment(js.command)
	str += "\n"
	str += fmt.Sprintf("package %s\n\n", js.packageName)
	str += fmt.Sprintf("import \"time\"\n\n")
	name := defaultTypeName
	if js.schema.Title != "" {
		name = js.schema.Title
	}
	str += "type"
	str += " "
	str += utils.Ucfirst(name)
	str += " "
	genarated, err := js.typer(0)
	if err != nil {
		return "", err
	}
	str += genarated
	str += "\n"

	src, err := format.Source([]byte(str))
	if err != nil {
		return str, err
	}

	return string(src), err
}

func (js *JSONSchema) structor(name string, nestLevel int, required bool) (string, error) {
	var str string

	for i := 1; i <= nestLevel; i++ {
		str += "\t"
	}
	str += utils.Ucfirst(name)
	str += " "

	genarated, err := js.typer(nestLevel)
	if err != nil {
		return "", err
	}
	str += genarated

	str += " "
	str += utils.MarshalTag(name, required)
	str += "\n"

	return str, nil
}

func (js *JSONSchema) typer(nestLevel int) (string, error) {
	var str string
	if inPrimitiveTypes(schema.IntegerType, js.schema.Type) {
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			str += "*"
		}
		str += "int"
	} else if inPrimitiveTypes(schema.StringType, js.schema.Type) {
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			str += "*"
		}
		if js.schema.Format == schema.FormatDateTime {
			str += "time.Time"
		} else {
			str += "string"
		}
	} else if inPrimitiveTypes(schema.ObjectType, js.schema.Type) {
		if js.schema.Properties == nil {
			str += "map[string]interface{}"
		} else {
			str += "struct {\n"
			for key, propSchema := range js.schema.Properties {
				prop := New(propSchema)
				tmp, err := prop.structor(key, nestLevel+1, js.schema.IsPropRequired(key))
				if err != nil {
					return "", err
				}
				str += tmp
			}
			for i := 1; i <= nestLevel; i++ {
				str += "\t"
			}
			str += "}"
		}
	} else if inPrimitiveTypes(schema.ArrayType, js.schema.Type) {
		if js.schema.Items == nil {
			str += "[]interface{}"
		} else if js.schema.Items.TupleMode {
			// unsupported
			err := fmt.Errorf("unsupported type %v", js.schema.Items)
			return "", err
		} else {
			str += "[]"
			item := New(js.schema.Items.Schemas[0])
			tmp, err := item.typer(nestLevel)
			if err != nil {
				return "", err
			}
			str += tmp
		}
	} else if inPrimitiveTypes(schema.BooleanType, js.schema.Type) {
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			str += "*"
		}
		str += "bool"
	} else if inPrimitiveTypes(schema.NumberType, js.schema.Type) {
		if inPrimitiveTypes(schema.NullType, js.schema.Type) {
			str += "*"
		}
		str += "float64"
	}
	return str, nil
}

func inPrimitiveTypes(needle schema.PrimitiveType, haystack schema.PrimitiveTypes) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
