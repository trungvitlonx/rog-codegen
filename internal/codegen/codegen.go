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
	"github.com/trungvitlonx/rog-codegen/internal/util"
)

var (
	//go:embed templates
	templates embed.FS

	funcs = template.FuncMap{"join": strings.Join}
)

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
	ClassName             string
	ServiceName           string
	Definitions           []Definition
	ControllerParentClass string
	ServiceParentClass    string
}

type CodegenService struct {
	Swagger *openapi3.T
	Config  Configuration
}

type GeneratedDir struct {
	ControllerDir string
	ServiceDir    string
	RoutesDir     string
	PackageDir    string
}

func NewCodegenService(swagger *openapi3.T, config Configuration) *CodegenService {
	return &CodegenService{
		Swagger: swagger,
		Config:  config,
	}
}

func (s CodegenService) Generate() (string, error) {
	t := template.New("rog-codegen").Funcs(funcs)
	err := loadAllTemplates(templates, t)
	if err != nil {
		return "", err
	}

	operationDefinitions, err := OperationDefinitions(s.Swagger)
	if err != nil {
		return "", err
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
	dirs := s.generateDirs()
	os.MkdirAll(dirs.ControllerDir, os.ModePerm)
	os.MkdirAll(dirs.ServiceDir, os.ModePerm)
	os.MkdirAll(dirs.RoutesDir, os.ModePerm)

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
			ClassName:             s.Config.OutputOptions.ControllerPrefix + "::" + strcase.ToCamel(className),
			ServiceName:           s.Config.OutputOptions.ServicePrefix + "::" + strcase.ToCamel(serviceName),
			Definitions:           definitions,
			ControllerParentClass: s.Config.OutputOptions.ControllerParentClass,
			ServiceParentClass:    s.Config.OutputOptions.ServiceParentClass,
		}

		controllerOut, err := generateTemplate("controller.tmpl", t, controllerData)
		if err != nil {
			return "", err
		}

		controllerFileName := strcase.ToSnake(className) + ".rb"
		controllerFile := dirs.ControllerDir + "/" + controllerFileName
		err = os.WriteFile(controllerFile, []byte(controllerOut), 0o644)
		if err != nil {
			return "", err
		}

		serviceOut, err := generateTemplate("service.tmpl", t, controllerData)
		if err != nil {
			return "", err
		}
		serviceFileName := strcase.ToSnake(serviceName) + ".rb"
		serviceFile := dirs.ServiceDir + "/" + serviceFileName
		err = os.WriteFile(serviceFile, []byte(serviceOut), 0o644)
		if err != nil {
			return "", err
		}
	}

	routesOut, err := generateTemplate("routes.tmpl", t, allDefinitions)
	if err != nil {
		return "", err
	}
	routesFileName := "api_routes.rb"
	routesFile := dirs.RoutesDir + "/" + routesFileName
	err = os.WriteFile(routesFile, []byte(routesOut), 0o644)
	if err != nil {
		return "", err
	}

	return dirs.PackageDir, nil
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

func (s CodegenService) generateDirs() GeneratedDir {
	workDir := s.Config.WorkingDirectory
	packageDir := workDir + "/" + strcase.ToSnake(s.Config.PackageName)

	controllerDir := packageDir + "/" + s.Config.OutputOptions.ControllerDirectory
	controllerPrefixes := strings.Split(s.Config.OutputOptions.ControllerPrefix, "::")
	for _, prefix := range controllerPrefixes {
		if len(prefix) <= 2 {
			controllerDir = controllerDir + "/" + strings.ToLower(prefix)
		} else {
			controllerDir = controllerDir + "/" + strcase.ToSnakeWithIgnore(prefix, "V")
		}
	}

	serviceDir := packageDir + "/" + s.Config.OutputOptions.ServiceDirectory
	if s.Config.OutputOptions.ServicePrefix != "" {
		servicePrefixes := strings.Split(s.Config.OutputOptions.ServicePrefix, "::")
		for _, prefix := range servicePrefixes {
			if len(prefix) <= 2 {
				serviceDir = serviceDir + "/" + strings.ToLower(prefix)
			} else {
				serviceDir = serviceDir + "/" + strcase.ToSnakeWithIgnore(prefix, "V")
			}
		}
	}

	routesDir := packageDir + "/" + s.Config.OutputOptions.RoutesDirectory
	return GeneratedDir{
		ControllerDir: controllerDir,
		ServiceDir:    serviceDir,
		RoutesDir:     routesDir,
		PackageDir:    packageDir,
	}
}
