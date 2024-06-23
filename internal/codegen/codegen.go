package codegen

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/trungle-csv/rog-codegen/internal/util"
)

//go:embed templates
var templates embed.FS

type Definition struct {
	MethodName      string
	HttpMethod      string
	Path            string
	Parameters      []ParameterDefinition
	BodyParameters  RequestBodyDefinition
	ControllerName  string
	IsResfulIndex   bool
	IsResfulShow    bool
	IsResfulCreate  bool
	IsResfulUpdate  bool
	IsResfulDestroy bool
}

type ControllerData struct {
	ClassName   string
	ServiceName string
	Definitions []Definition
}

func Generate(swagger *openapi3.T, output string) error {
	funcs := template.FuncMap{"join": strings.Join}
	t := template.New("rog-codegen").Funcs(funcs)
	err := loadAllTemplates(templates, t)
	if err != nil {
		return err
	}
	operationDefinitions, err := OperationDefinitions(swagger)
	if err != nil {
		return err
	}
	groupedOperations := make(map[string][]OperationDefinition)
	for _, op := range operationDefinitions {
		if list, exists := groupedOperations[op.Tag]; exists {
			list = append(list, op)
			groupedOperations[op.Tag] = list
		} else {
			list := []OperationDefinition{}
			list = append(list, op)
			groupedOperations[op.Tag] = list
		}
	}
	allDefinitions := make([]Definition, 0)
	for _, tag := range util.SortedMapKeys(groupedOperations) {
		className := tag + "Controller"
		serviceName := tag + "Service"
		operations := groupedOperations[tag]
		definitions := make([]Definition, 0, len(operations))
		for _, op := range operations {
			parameters := op.PathParams
			parameters = append(parameters, op.QueryParams...)
			definition := Definition{
				MethodName:      strcase.ToSnake(op.OperationId),
				HttpMethod:      op.Method,
				Path:            op.Path,
				Parameters:      parameters,
				BodyParameters:  op.BodyParams,
				ControllerName:  strcase.ToSnake(className),
				IsResfulIndex:   false,
				IsResfulShow:    false,
				IsResfulCreate:  false,
				IsResfulUpdate:  false,
				IsResfulDestroy: false,
			}
			definitions = append(definitions, definition)
			allDefinitions = append(allDefinitions, definition)
		}

		controllerData := ControllerData{
			ClassName:   strcase.ToCamel(className),
			ServiceName: strcase.ToCamel(serviceName),
			Definitions: definitions,
		}
		controllerOut, err := generateTemplate("controller.tmpl", t, controllerData)
		if err != nil {
			return err
		}
		controllerFileName := strcase.ToSnake(className) + ".rb"
		controllerFile := output + "/gen/controllers/" + controllerFileName
		err = os.WriteFile(controllerFile, []byte(controllerOut), 0o644)
		if err != nil {
			return err
		}

		serviceOut, err := generateTemplate("service.tmpl", t, controllerData)
		if err != nil {
			return err
		}
		serviceFileName := strcase.ToSnake(serviceName) + ".rb"
		serviceFile := output + "/gen/services/" + serviceFileName
		err = os.WriteFile(serviceFile, []byte(serviceOut), 0o644)
		if err != nil {
			return err
		}
	}

	routesOut, err := generateTemplate("routes.tmpl", t, allDefinitions)
	if err != nil {
		return err
	}
	routesFileName := "api_routes.rb"
	routesFile := output + "/gen/config/" + routesFileName
	err = os.WriteFile(routesFile, []byte(routesOut), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func generateTemplate(templateName string, t *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	if err := t.ExecuteTemplate(w, templateName, data); err != nil {
		return "", fmt.Errorf("error generating %s: %s", templateName, err)
	}
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("error flushing output buffer %s: %s", templateName, err)
	}

	return buf.String(), nil
}

func loadAllTemplates(src embed.FS, template *template.Template) error {
	return fs.WalkDir(src, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %s", path, err)
		}
		if d.IsDir() {
			return nil
		}

		buf, err := src.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %s", path, err)
		}

		templateName := strings.TrimPrefix(path, "templates/")
		tmpl := template.New(templateName)
		_, err = tmpl.Parse(string(buf))
		if err != nil {
			return fmt.Errorf("error parsing template %s: %s", path, err)
		}
		return nil
	})
}
