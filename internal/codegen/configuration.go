package codegen

type Configuration struct {
	WorkingDirectory string        `yaml:"working-directory"`
	PackageName      string        `yaml:"package"`
	OutputOptions    OutputOptions `yaml:"out-options,omitempty"`
}

type OutputOptions struct {
	UserTemplates         map[string]string `yaml:"user-templates,omitempty"`
	ControllerPrefix      string            `yaml:"controller-prefix,omitempty"`
	ServicePrefix         string            `yaml:"service-prefix,omitempty"`
	ControllerParentClass string            `yaml:"controller-parent-class,omitempty"`
	ServiceParentClass    string            `yaml:"service-parent-class,omitempty"`
	RegenerateService     bool              `yaml:"regenerate-service"`
}
