package codegen

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/trungle-csv/rog-codegen/internal/util"
)

type ParameterDefinition struct {
	ParamName string
	In        string
	Required  bool
	Spec      *openapi3.Parameter
	Schema    Schema
}

type ParameterDefinitions []ParameterDefinition

type RequestBodyDefinition struct {
	Schema Schema
	Spec   *openapi3.Schema
}

type OperationDefinition struct {
	Tag          string
	OperationId  string
	Method       string
	Path         string
	HeaderParams []ParameterDefinition
	PathParams   []ParameterDefinition
	QueryParams  []ParameterDefinition
	BodyParams   RequestBodyDefinition
}

func OperationDefinitions(swagger *openapi3.T) ([]OperationDefinition, error) {
	var operations []OperationDefinition

	if swagger == nil || swagger.Paths == nil {
		return operations, nil
	}

	for _, requestPath := range util.SortedMapKeys(swagger.Paths.Map()) {
		pathItem := swagger.Paths.Value(requestPath)
		pathOps := pathItem.Operations()

		for _, opName := range util.SortedMapKeys(pathOps) {
			op := pathOps[opName]
			if op.OperationID == "" {
				return nil, fmt.Errorf("missing operationId for path: %s, method: %s", requestPath, opName)
			}

			queryParams, pathParams, _ := DescribeParameters(op.Parameters)
			bodyParams, _ := DescribeRequestBody(op.RequestBody)
			operationDef := OperationDefinition{
				Tag:         op.Tags[0],
				OperationId: op.OperationID,
				Method:      opName,
				Path:        requestPath,
				PathParams:  pathParams,
				QueryParams: queryParams,
				BodyParams:  bodyParams,
			}

			operations = append(operations, operationDef)
		}
	}

	return operations, nil
}

func DescribeParameters(params openapi3.Parameters) ([]ParameterDefinition, []ParameterDefinition, error) {
	queryParams := make([]ParameterDefinition, 0)
	pathParams := make([]ParameterDefinition, 0)

	for _, paramOrRef := range params {
		param := paramOrRef.Value
		schema, _ := DescribeSchemaRef(param.Schema.Value)
		paramDef := ParameterDefinition{
			ParamName: param.Name,
			In:        param.In,
			Required:  param.Required,
			Spec:      param,
			Schema:    schema,
		}

		if param.In == "query" {
			queryParams = append(queryParams, paramDef)
		} else if param.In == "path" {
			pathParams = append(pathParams, paramDef)
		}
	}

	return queryParams, pathParams, nil
}

func DescribeRequestBody(requestBody *openapi3.RequestBodyRef) (RequestBodyDefinition, error) {
	if requestBody == nil {
		return RequestBodyDefinition{}, fmt.Errorf("no request body")
	}

	content := requestBody.Value.Content["application/json"]

	if content == nil {
		content = requestBody.Value.Content["application/x-www-form-urlencoded"]
	}

	if content == nil {
		return RequestBodyDefinition{}, fmt.Errorf("failed to describe request body")
	}

	schema, _ := DescribeSchemaRef(content.Schema.Value)

	return RequestBodyDefinition{
		Spec:   content.Schema.Value,
		Schema: schema,
	}, nil
}

func DescribeSchemaRef(schemaRef *openapi3.Schema) (Schema, error) {
	schemaType := *(schemaRef.Type)
	var arrayItems Schema

	if schemaType[0] == "array" {
		itemSchemaRef := *(schemaRef.Items.Value)
		arrayItems, _ = DescribeSchemaRef(&itemSchemaRef)
	}

	properties := make([]Property, 0, len(schemaRef.Properties))
	requiredProps := schemaRef.Required
	for propName, propSchemaRef := range schemaRef.Properties {
		propSchema, _ := DescribeSchemaRef(propSchemaRef.Value)

		properties = append(properties, Property{
			Name:     propName,
			Required: util.SliceContains(requiredProps, propName),
			Schema:   propSchema,
		})
	}

	schema := Schema{
		RefType:    schemaType[0],
		EnumValues: util.ToSliceOfString(schemaRef.Enum),
		OAPISchema: schemaRef,
		ArrayItems: &arrayItems,
		Properties: properties,
	}

	return schema, nil
}
