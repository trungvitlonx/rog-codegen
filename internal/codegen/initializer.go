package codegen

import (
	"os"
	"text/template"
)

func Initialize() error {
	t := template.New("rog-codegen").Funcs(funcs)
	err := loadAllTemplates(templates, t)
	if err != nil {
		return err
	}

	configOut, err := generateTemplate("config.yaml.tmpl", t, nil)
	if err != nil {
		return err
	}
	configName := ".rog.yaml"
	err = os.WriteFile(configName, []byte(configOut), 0o644)
	if err != nil {
		return err
	}

	oapiOut, err := generateTemplate("openapi.yaml.tmpl", t, nil)
	if err != nil {
		return err
	}
	oapiName := "openapi.yaml"
	err = os.WriteFile(oapiName, []byte(oapiOut), 0o644)
	if err != nil {
		return err
	}

	return nil
}
