package codegen

import "github.com/getkin/kin-openapi/openapi3"

type Schema struct {
	RefType    string
	EnumValues []string
	Properties []Property
	OAPISchema *openapi3.Schema
	ArrayItems *Schema
}

type Property struct {
	Name     string
	Required bool
	Schema   Schema
}

func RefTypeMapping() map[string]string {
	return map[string]string{
		"string":  "String",
		"integer": "Integer",
		"array":   "Array",
		"object":  "Hash",
	}
}
