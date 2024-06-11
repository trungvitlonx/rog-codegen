package codegen

import "github.com/getkin/kin-openapi/openapi3"

func Generate(swagger *openapi3.T, output string) error {
	operationDefinitions, err := OperationDefinitions(swagger)
	if err != nil {
		return err
	}

	_ = operationDefinitions
	return nil
}
