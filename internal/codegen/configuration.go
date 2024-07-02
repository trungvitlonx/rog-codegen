package codegen

import (
	"errors"
	"reflect"
)

type Configuration struct {
	PackageName      string        `yaml:"package"`
	WorkingDirectory string        `yaml:"directory,omitempty"`
	OutputOptions    OutputOptions `yaml:"out-options,omitempty"`
	UserTemplates    UserTemplates `yaml:"user-templates,omitempty"`
}

type OutputOptions struct {
	ControllerPrefix      string `yaml:"controller-prefix,omitempty"`
	ServicePrefix         string `yaml:"service-prefix,omitempty"`
	ControllerParentClass string `yaml:"controller-parent-class,omitempty"`
	ServiceParentClass    string `yaml:"service-parent-class,omitempty"`
	RegenerateService     bool   `yaml:"regenerate-service"`
	ControllerDirectory   string `yaml:"controller-directory,omitempty"`
	ServiceDirectory      string `yaml:"service-directory,omitempty"`
	RoutesDirectory       string `yaml:"routes-directory,omitempty"`
}

type UserTemplates struct {
	Controller string `yaml:"controller,omitempty"`
	Service    string `yaml:"service,omitempty"`
	Routes     string `yaml:"routes,omitempty"`
}

func (c Configuration) UpdateDefaultValues() Configuration {
	if reflect.ValueOf(c.OutputOptions).IsZero() {
		c.OutputOptions = OutputOptions{
			ControllerPrefix:      "API",
			ControllerParentClass: "ApplicationController",
			RegenerateService:     false,
			ControllerDirectory:   "controllers",
			ServiceDirectory:      "services",
			RoutesDirectory:       "config",
		}
	}

	if reflect.ValueOf(c.WorkingDirectory).IsZero() {
		c.WorkingDirectory = "./"
	}

	return c
}

func (c Configuration) Validate() error {
	if c.PackageName == "" {
		return errors.New("package name must be specified")
	}

	return nil
}
