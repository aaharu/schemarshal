// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/aaharu/schemarshal/utils"
	"github.com/aaharu/schemarshal/version"
	schema "github.com/lestrrat/go-jsschema"
)

// Generator of Go source code from JSON Schema
type Generator interface {
	ReadSchema(input io.Reader, name string) error
	Generate(noComment bool) ([]byte, error)
}

type generator struct {
	name     string // package nage
	command  string
	imports  importSpec
	decls    []*typeSpec
	enumList enumSpec
}

// NewGenerator create Generator struct
func NewGenerator(packageName string, command string) Generator {
	return &generator{
		name:     packageName,
		command:  command,
		imports:  importSpec{},
		decls:    []*typeSpec{},
		enumList: enumSpec{},
	}
}

// ReadSchema : read and parse JSON Schema
func (g *generator) ReadSchema(input io.Reader, name string) error {
	js, err := schema.Read(input)
	if err != nil {
		return err
	}
	if js.Title != "" && utils.UpperCamelCase(js.Title) != "" {
		name = js.Title
	}
	name = utils.UpperCamelCase(name)
	genType, err := g.parse(js, name)
	if err != nil {
		return err
	}
	g.addType(name, &genType)
	return nil
}

// parse returns JSON Schema type
func (g *generator) parse(js *schema.Schema, fieldName string) (jsonType, error) {
	jt := jsonType{}
	if inPrimitiveTypes(schema.IntegerType, js.Type) ||
		inPrimitiveTypes(schema.BooleanType, js.Type) ||
		inPrimitiveTypes(schema.NumberType, js.Type) {
		if inPrimitiveTypes(schema.IntegerType, js.Type) {
			jt.format = formatInteger
		} else if inPrimitiveTypes(schema.BooleanType, js.Type) {
			jt.format = formatBoolean
		} else {
			jt.format = formatNumber
		}
		if inPrimitiveTypes(schema.NullType, js.Type) {
			jt.nullable = true
		}
		if js.Enum != nil {
			enumName := fieldName + "Enum"
			if _, ok := g.enumList[enumName]; ok == true {
				// FIXME: unsupported
				return jt, errors.New("unsupported json")
			}
			g.enumList[enumName] = js.Enum
			jt.enumType = enumName
			g.imports[`"strconv"`] = ""
			g.imports[`"fmt"`] = ""
		}
		return jt, nil
	}
	if inPrimitiveTypes(schema.StringType, js.Type) {
		if js.Format == schema.FormatDateTime {
			jt.format = formatDatetime
			g.imports[`"time"`] = ""
		} else {
			jt.format = formatString
		}
		if inPrimitiveTypes(schema.NullType, js.Type) {
			jt.nullable = true
		}
		if js.Enum != nil {
			enumName := fieldName + "Enum"
			g.enumList[enumName] = js.Enum
			jt.enumType = enumName
			g.imports[`"strconv"`] = ""
			g.imports[`"fmt"`] = ""
		}
		return jt, nil
	}
	if inPrimitiveTypes(schema.ArrayType, js.Type) {
		if js.Items.TupleMode {
			// unsupported
			return jt, fmt.Errorf("unsupported type %v", js.Items)
		}
		jt.format = formatArray
		if inPrimitiveTypes(schema.NullType, js.Type) {
			jt.nullable = true
		}
		itemType, err := g.parse(js.Items.Schemas[0], fieldName)
		if err != nil {
			return jt, err
		}
		jt.itemType = &itemType
		if itemType.format == formatObject {
			itemFieldName := fieldName + "Item"
			jt.typeName = itemFieldName
			g.addType(itemFieldName, &itemType)
		}
		return jt, nil
	}
	jt.description = js.Description
	jt.format = formatObject
	if inPrimitiveTypes(schema.NullType, js.Type) {
		jt.nullable = true
	}
	if js.Properties != nil {
		// sort map
		var keys []string
		for k := range js.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			propSchema := js.Properties[key]
			propType, err := g.parse(propSchema, fieldName+utils.UpperCamelCase(key))
			if err != nil {
				return jt, err
			}
			if propType.format == formatObject {
				objectTypeName := fieldName + utils.UpperCamelCase(key) + "Object"
				g.addType(objectTypeName, &propType)
				copyType := jsonType{
					format:   propType.format,
					nullable: propType.nullable,
					fields:   propType.fields,
					itemType: propType.itemType,
					typeName: objectTypeName,
					enumType: propType.enumType,
				}
				jt.addField(&field{
					name:        utils.UpperCamelCase(key),
					description: propSchema.Description,
					jsontype:    copyType,
					jsontag: jsonTag{
						name:      key,
						omitEmpty: !js.IsPropRequired(key),
					},
				})
			} else {
				jt.addField(&field{
					name:        utils.UpperCamelCase(key),
					description: propSchema.Description,
					jsontype:    propType,
					jsontag: jsonTag{
						name:      key,
						omitEmpty: !js.IsPropRequired(key),
					},
				})
			}
		}
	}
	if js.Enum != nil {
		enumName := fieldName + "Enum"
		g.enumList[enumName] = js.Enum
		jt.enumType = enumName
		g.imports[`"strconv"`] = ""
		g.imports[`"fmt"`] = ""
	}
	return jt, nil
}

// addType add a type statement
func (g *generator) addType(name string, jsonType *jsonType) {
	g.decls = append(g.decls, &typeSpec{
		name:     name,
		jsontype: jsonType,
	})
}

// Generate gofmt-ed Go source code
func (g *generator) Generate(noComment bool) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("// Code generated by %s `%s`\n", version.String(), g.command))
	buf.WriteString("// DO NOT RECOMMEND EDITING THIS FILE.\n\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.name))

	if len(g.imports) > 1 {
		buf.WriteString("import (\n")
		// sort map
		var keys []string
		for k := range g.imports {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, path := range keys {
			name := g.imports[path]
			buf.WriteString(fmt.Sprintf("%s %s\n", name, path))
		}
		buf.WriteString(")\n\n")
	} else if len(g.imports) == 1 {
		for path, name := range g.imports {
			buf.WriteString(fmt.Sprintf("import %s %s\n\n", name, path))
		}
	}

	if g.decls != nil {
		for i := range g.decls {
			if !noComment && g.decls[i].jsontype.description != "" {
				buf.WriteString(fmt.Sprintf("// %s : %s\n", g.decls[i].name, utils.CleanDescription(g.decls[i].jsontype.description)))
			}
			buf.WriteString("type " + g.decls[i].name + " ")
			g.decls[i].jsontype.nullable = false
			buf.Write(g.decls[i].jsontype.generate(noComment))
			buf.WriteString("\n\n")
		}
	}

	buf.Write(g.enumList.generate())

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.Bytes(), err
	}
	return src, nil
}

// importSpec has `import` information
type importSpec map[string]string

type typeSpec struct {
	name     string // type name
	jsontype *jsonType
}

// enumSpec has enum information
type enumSpec map[string][]interface{}

func (e enumSpec) generate() []byte {
	var buf bytes.Buffer
	// sort map
	var keys []string
	for k := range e {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, typeName := range keys {
		enum := e[typeName]
		buf.WriteString("\ntype " + typeName + " int\n")
		buf.WriteString("const (\n")
		for i := range enum {
			buf.WriteString(utils.UpperCamelCase(fmt.Sprintf("%s %v", typeName, enum[i])))
			if i == 0 {
				buf.WriteString(" " + typeName + " = iota\n")
			} else {
				buf.WriteString("\n")
			}
		}
		buf.WriteString(")\n\n")

		var enumMapName = "_" + strings.ToLower(typeName[:1])
		if len(typeName) > 1 {
			enumMapName += typeName[1:]
		}
		buf.WriteString("var " + enumMapName + " = map[" + typeName + "]interface{}{\n")
		for i := range enum {
			buf.WriteString(utils.UpperCamelCase(fmt.Sprintf("%s %v", typeName, enum[i])) + ": ")
			switch v := enum[i].(type) {
			case string:
				buf.WriteString(strconv.Quote(v))
			default:
				buf.WriteString(fmt.Sprintf("%v", v))
			}
			buf.WriteString(",\n")
		}
		buf.WriteString("}\n\n")

		buf.WriteString("func (enum " + typeName + ") MarshalJSON() ([]byte, error) {\n")
		buf.WriteString("switch v:= " + enumMapName + "[enum].(type) {\n")
		buf.WriteString("case string:\n")
		buf.WriteString("return []byte(strconv.Quote(v)), nil\n")
		buf.WriteString("default:\n")
		buf.WriteString("return []byte(fmt.Sprintf(\"%v\", v)), nil\n")
		buf.WriteString("}\n")
		buf.WriteString("}\n\n")

		buf.WriteString("func (enum *" + typeName + ") UnmarshalJSON(data []byte) error {\n")
		buf.WriteString("for i, v := range " + enumMapName + " {\n")
		buf.WriteString("switch vv := v.(type) {\n")
		buf.WriteString("case string:\n")
		buf.WriteString("if strconv.Quote(vv) == string(data) {\n")
		buf.WriteString("*enum = " + typeName + "(i)\n")
		buf.WriteString("return nil\n")
		buf.WriteString("}\n")
		buf.WriteString("default:\n")
		buf.WriteString("if fmt.Sprintf(\"%v\", v) == string(data) {\n")
		buf.WriteString("*enum = " + typeName + "(i)\n")
		buf.WriteString("return nil\n")
		buf.WriteString("}\n")
		buf.WriteString("}\n")
		buf.WriteString("}\n")
		buf.WriteString("return fmt.Errorf(\"Error: miss-matched " + typeName + " (%s)\", data)\n")
		buf.WriteString("}\n\n")

		buf.WriteString("func (enum " + typeName + ") String() string {\n")
		buf.WriteString("switch v:= " + enumMapName + "[enum].(type) {\n")
		buf.WriteString("case string:\n")
		buf.WriteString("return v\n")
		buf.WriteString("default:\n")
		buf.WriteString("return fmt.Sprintf(\"%v\", v)\n")
		buf.WriteString("}\n")
		buf.WriteString("}\n\n")

		buf.WriteString("func To" + typeName + "(val interface{}) *" + typeName + " {\n")
		buf.WriteString("for i, v := range " + enumMapName + " {\n")
		buf.WriteString("if val == v {")
		buf.WriteString("return &i")
		buf.WriteString("}\n")
		buf.WriteString("}\n")
		buf.WriteString("return nil")
		buf.WriteString("}\n")
	}
	return buf.Bytes()
}

type jsonFormat int

const (
	formatObject jsonFormat = iota
	formatArray
	formatString
	formatBoolean
	formatNumber
	formatInteger
	formatDatetime
)

// jsonType is type of json
type jsonType struct {
	format      jsonFormat
	nullable    bool
	fields      []*field  // object has
	itemType    *jsonType // array has
	typeName    string    // object's array and object has
	enumType    string    // enum has
	description string    // for comment
}

func (t *jsonType) addField(f *field) {
	t.fields = append(t.fields, f)
}

func (t *jsonType) generate(noComment bool) []byte {
	var buf bytes.Buffer
	if t.nullable {
		buf.WriteString("*")
	}
	if t.enumType != "" {
		buf.WriteString(t.enumType)
	} else if t.format == formatObject {
		if t.fields == nil {
			buf.WriteString("map[string]interface{}")
		} else {
			if t.typeName != "" {
				buf.WriteString(t.typeName)
			} else {
				buf.WriteString("struct {\n")
				for i := range t.fields {
					if !noComment && t.fields[i].description != "" {
						buf.WriteString(fmt.Sprintf("// %s : %s\n", t.fields[i].name, utils.CleanDescription(t.fields[i].description)))
					}
					buf.WriteString(t.fields[i].name)
					buf.WriteString(" ")
					buf.Write(t.fields[i].jsontype.generate(noComment))
					buf.WriteString(" ")
					buf.Write(t.fields[i].jsontag.generate())
					buf.WriteString("\n")
				}
				buf.WriteString("}")
			}
		}
	} else if t.format == formatArray {
		buf.WriteString("[]")
		if t.typeName != "" {
			buf.WriteString(t.typeName)
		} else {
			buf.Write(t.itemType.generate(noComment))
		}
	} else if t.format == formatString {
		buf.WriteString("string")
	} else if t.format == formatBoolean {
		buf.WriteString("bool")
	} else if t.format == formatNumber {
		buf.WriteString("float64")
	} else if t.format == formatInteger {
		buf.WriteString("int64")
	} else if t.format == formatDatetime {
		buf.WriteString("time.Time")
	}
	return buf.Bytes()
}

type field struct {
	name        string
	description string   // comment
	jsontype    jsonType // go type
	jsontag     jsonTag  // `json:""`
}

type jsonTag struct {
	name      string
	omitEmpty bool
}

// Generate JSON tag code
func (t *jsonTag) generate() []byte {
	var buf bytes.Buffer
	buf.WriteString("`json:\"")
	buf.WriteString(t.name)
	if t.omitEmpty {
		buf.WriteString(",omitempty")
	}
	buf.WriteString("\"`")
	return buf.Bytes()
}

func inPrimitiveTypes(needle schema.PrimitiveType, haystack schema.PrimitiveTypes) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
